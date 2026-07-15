import * as readline from "node:readline";
import * as http from "node:http";
import * as https from "node:https";
import * as Automerge from "@automerge/automerge";

const AutomergeNext = Automerge.next;
const EMPTY_DOCUMENT_BYTES = Uint8Array.from(Buffer.from("hW9Kg8HDZmEAdQEQUDnUuZsuTLOKK6EtAqSUwAF91ThR16b5XY1P61eTHXkwnJNTicqZ35V+jMImBQWmigYBAgMCEwIjBkACVgIHFQkhAiMCNAFCAlYCgAECfwB/AX8Bf8660dIGfwB/B38HY29udGVudH8AfwEBfwR/AH8AAA==", "base64"));

const state = {
  relayUrl: "",
  participantId: "",
  displayName: "Neovim User",
  color: "#d66f1d",
  documentId: "",
  doc: ensureDocument(null),
  offset: 0,
  selection: { anchor: 0, head: 0 },
  syncTimer: null,
  awarenessTimer: null,
  syncSocket: null,
  awarenessSocket: null,
  relayConnected: false,
  relayTransport: "polling",
  awarenessTransport: "polling",
};

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: false,
});

let pendingMessage = Promise.resolve();

rl.on("line", (line) => {
  if (!line.trim()) {
    return;
  }
  let message;
  try {
    message = JSON.parse(line);
  } catch (error) {
    send({ type: "error", message: `invalid json: ${error.message}` });
    return;
  }
  // Intent: Preserve Neovim command ordering so `open`, `set_text`, cursor, and
  // `close` messages cannot race each other through overlapping async relay
  // work. Source: DI-sulod; DI-gafit
  pendingMessage = pendingMessage.then(() => handleMessage(message)).catch((error) => {
    send({ type: "error", message: error.stack || error.message });
  });
});

function send(message) {
  process.stdout.write(`${JSON.stringify(message)}\n`);
}

function sendInfo(message) {
  send({ type: "info", message });
}

async function handleMessage(message) {
  switch (message.type) {
    case "connect":
      state.relayUrl = message.relay_url;
      state.participantId = message.participant_id;
      state.displayName = message.display_name || state.displayName;
      state.color = message.color || state.color;
      send({ type: "connected", participant_id: state.participantId });
      break;
    case "open":
      await openDocument(message.doc_id || "demo");
      break;
    case "set_text":
      if (!state.documentId) {
        return;
      }
      applyLocalText(message.content || "");
      break;
    case "set_cursor":
      state.selection = {
        anchor: Number.isFinite(message.anchor) ? message.anchor : 0,
        head: Number.isFinite(message.head) ? message.head : (Number.isFinite(message.anchor) ? message.anchor : 0),
      };
      await postAwareness(Boolean(message.typing));
      break;
    case "set_user":
      if (message.display_name) {
        state.displayName = message.display_name;
      }
      if (message.color) {
        state.color = message.color;
      }
      await postAwareness(false);
      break;
    case "close":
      closeDocument();
      send({ type: "closed" });
      break;
    default:
      send({ type: "error", message: `unknown message type ${message.type}` });
  }
}

async function openDocument(documentId) {
  closeDocument();
  state.documentId = documentId;
  state.doc = ensureDocument(null);
  state.offset = 0;
  // Intent: Move the Neovim embodiment onto ex3's websocket live transport
  // while preserving the same stdin/stdout sidecar contract that the plugin
  // already speaks. Source: DI-bitus
  if (websocketCapable()) {
    await connectSyncSocket();
    await connectAwarenessSocket();
  } else {
    state.relayTransport = "polling";
    state.awarenessTransport = "polling";
    await pollSync();
    await pollAwareness();
    state.syncTimer = setInterval(() => {
      pollSync().catch((error) => send({ type: "error", message: error.stack || error.message }));
    }, 250);
    state.awarenessTimer = setInterval(() => {
      pollAwareness().catch((error) => send({ type: "error", message: error.stack || error.message }));
    }, 350);
    await postAwareness(false);
  }
  send({
    type: "opened",
    doc_id: state.documentId,
    content: getText(),
    relay_transport: state.relayTransport,
    awareness_transport: state.awarenessTransport,
  });
}

function closeDocument() {
  if (state.syncTimer) {
    clearInterval(state.syncTimer);
    state.syncTimer = null;
  }
  if (state.awarenessTimer) {
    clearInterval(state.awarenessTimer);
    state.awarenessTimer = null;
  }
  if (state.syncSocket) {
    state.syncSocket.close();
    state.syncSocket = null;
  }
  if (state.awarenessSocket) {
    state.awarenessSocket.close();
    state.awarenessSocket = null;
  }
  state.documentId = "";
  state.doc = ensureDocument(null);
  state.offset = 0;
  setRelayConnected(false);
}

function applyLocalText(content) {
  const previous = getText();
  if (content === previous) {
    return;
  }
  const prefix = commonPrefix(previous, content);
  const suffix = commonSuffix(previous, content, prefix);
  const deleteCount = previous.length - prefix - suffix;
  const insertText = content.slice(prefix, content.length - suffix);
  // Intent: Keep the sidecar replica writable after remote sync has advanced
  // it, because Automerge returns immutable historical snapshots and Neovim
  // local edits must always fork from the current head. Source: DI-sulod
  state.doc = Automerge.change(Automerge.clone(state.doc), (draft) => {
    AutomergeNext.splice(draft, ["content"], prefix, deleteCount, insertText);
  });
  // Intent: Send durable Automerge change packets so the relay log can replay
  // full document history without reconstructing per-peer sync sessions.
  // Source: DI-sulod; DI-larok
  const change = Automerge.getLastLocalChange(state.doc);
  if (change) {
    postChange(change).catch((error) => send({ type: "error", message: error.stack || error.message }));
  }
}

async function pollSync() {
  if (!state.documentId) {
    return;
  }
  while (true) {
    const payload = await getJSON(`${basePath()}/sync?since=${state.offset}&limit=256`);
    setRelayConnected(true);
    for (const record of payload.messages || []) {
      // Intent: Treat the relay as an exchange surface for peer sync, not as a
      // loopback transport for this participant's own signed messages. Replaying
      // self-authored sync records can perturb local replica state without adding
      // new information. Source: DI-sulod; DI-gafit
      if (record.participant_id === state.participantId) {
        continue;
      }
      if (record.recipient_id && record.recipient_id !== state.participantId) {
        continue;
      }
      await receive(record);
    }
    const nextOffset = payload.next_offset || state.offset;
    if (nextOffset <= state.offset) {
      break;
    }
    state.offset = nextOffset;
    if ((payload.messages || []).length < 256) {
      break;
    }
  }
}

async function pollAwareness() {
  if (!state.documentId) {
    return;
  }
  const payload = await getJSON(`${basePath()}/awareness`);
  const peers = [];
  for (const peer of payload.awareness || []) {
    if (!peer.participant_id || peer.participant_id === state.participantId) {
      continue;
    }
    peers.push({
      participant_id: peer.participant_id,
      name: peer.display_name || peer.author || peer.participant_id,
      color: peer.color || "#999999",
      anchor: peer.cursor || 0,
      head: peer.head || peer.cursor || 0,
      typing: Boolean(peer.typing),
      last_seen_at: peer.last_seen_at || null,
    });
  }
  send({ type: "awareness", peers });
}

async function postAwareness(typing) {
  if (!state.documentId) {
    return;
  }
  // Intent: Keep live-awareness as its own sidecar transport channel so Neovim
  // presence and cursor updates stay separate from durable document sync.
  // Source: DI-bitus
  if (state.awarenessSocket && state.awarenessSocket.readyState === WebSocket.OPEN) {
    state.awarenessSocket.send(JSON.stringify({
      type: "post-awareness",
      participant_id: state.participantId,
      cursor: state.selection.anchor,
      head: state.selection.head,
      typing,
      display_name: state.displayName,
      color: state.color,
      embodiment: "nvim",
    }));
    return;
  }
  if (state.awarenessTransport === "websocket") {
    return;
  }
  await postJSON(`${basePath()}/awareness`, {
    participant_id: state.participantId,
    cursor: state.selection.anchor,
    head: state.selection.head,
    typing,
    display_name: state.displayName,
    color: state.color,
    embodiment: "nvim",
  });
}

async function postChange(changeBytes) {
  if (!state.documentId) {
    return;
  }
  // Intent: Send sidecar live-document traffic over websocket in ex3 without
  // changing the signed change payloads or the relay's feed semantics. Source:
  // DI-bitus
  if (state.syncSocket && state.syncSocket.readyState === WebSocket.OPEN) {
    state.syncSocket.send(JSON.stringify({
      type: "post-sync",
      participant_id: state.participantId,
      recipient_id: "",
      message_base64: bytesToBase64(changeBytes),
      embodiment: "nvim",
    }));
    return;
  }
  await postJSON(`${basePath()}/sync`, {
    participant_id: state.participantId,
    recipient_id: "",
    message_base64: bytesToBase64(changeBytes),
    embodiment: "nvim",
  });
}

async function getJSON(url) {
  return requestJSON("GET", url);
}

async function postJSON(url, body) {
  return requestJSON("POST", url, body);
}

function requestJSON(method, rawURL, body) {
  return new Promise((resolve, reject) => {
    const url = new URL(rawURL);
    const client = url.protocol === "https:" ? https : http;
    const payload = body ? JSON.stringify(body) : null;
    const request = client.request({
      protocol: url.protocol,
      hostname: url.hostname,
      port: url.port,
      path: `${url.pathname}${url.search}`,
      method,
      headers: payload
        ? {
            "Content-Type": "application/json",
            "Content-Length": Buffer.byteLength(payload),
          }
        : undefined,
    }, (response) => {
      let data = "";
      response.setEncoding("utf8");
      response.on("data", (chunk) => {
        data += chunk;
      });
      response.on("end", () => {
        const statusCode = response.statusCode || 0;
        if (statusCode < 200 || statusCode >= 300) {
          setRelayConnected(false);
          reject(new Error(`${method} ${rawURL} failed: ${statusCode} ${data}`));
          return;
        }
        if (data === "") {
          resolve({});
          return;
        }
        try {
          resolve(JSON.parse(data));
        } catch (error) {
          reject(error);
        }
      });
    });
    request.on("error", (error) => {
      setRelayConnected(false);
      reject(error);
    });
    if (payload) {
      request.write(payload);
    }
    request.end();
  });
}

async function receive(record) {
  if (!record.participant_id) {
    return;
  }
  const previous = Automerge.clone(state.doc);
  const [nextDoc] = Automerge.applyChanges(previous, [base64ToBytes(record.message_base64)]);
  if (!Automerge.equals(previous, nextDoc)) {
    state.doc = ensureDocument(nextDoc);
    send({
      type: "changed",
      content: getText(),
    });
  }
}

function basePath() {
  return `${state.relayUrl.replace(/\/$/, "")}/api/local/documents/${encodeURIComponent(state.documentId)}`;
}

async function connectSyncSocket() {
  state.relayTransport = "websocket";
  const socket = new WebSocket(toWebSocketURL("sync-socket", { since: state.offset }));
  state.syncSocket = socket;
  await new Promise((resolve, reject) => {
    let settled = false;
    socket.addEventListener("open", () => {
      setRelayConnected(true);
    });
    socket.addEventListener("message", (event) => {
      handleSyncSocketMessage(event.data)
        .then(() => {
          if (!settled) {
            settled = true;
            resolve();
          }
        })
        .catch((error) => {
          if (!settled) {
            settled = true;
            reject(error);
            return;
          }
          send({ type: "error", message: error.stack || error.message });
        });
    });
    socket.addEventListener("error", () => {
      const error = new Error("sync websocket failed");
      setRelayConnected(false);
      if (!settled) {
        settled = true;
        reject(error);
        return;
      }
      send({ type: "error", message: error.message });
    });
    socket.addEventListener("close", () => {
      if (state.syncSocket === socket) {
        state.syncSocket = null;
      }
      setRelayConnected(false);
    });
  });
}

async function connectAwarenessSocket() {
  state.awarenessTransport = "websocket";
  const socket = new WebSocket(toWebSocketURL("awareness-socket"));
  state.awarenessSocket = socket;
  await new Promise((resolve, reject) => {
    let settled = false;
    socket.addEventListener("open", () => {
      Promise.resolve(postAwareness(false)).catch((error) => {
        if (!settled) {
          settled = true;
          reject(error);
        }
      });
    });
    socket.addEventListener("message", (event) => {
      try {
        handleAwarenessSocketMessage(event.data);
        if (!settled) {
          settled = true;
          resolve();
        }
      } catch (error) {
        if (!settled) {
          settled = true;
          reject(error);
          return;
        }
        send({ type: "error", message: error.stack || error.message });
      }
    });
    socket.addEventListener("error", () => {
      const error = new Error("awareness websocket failed");
      setRelayConnected(false);
      if (!settled) {
        settled = true;
        reject(error);
        return;
      }
      send({ type: "error", message: error.message });
    });
    socket.addEventListener("close", () => {
      if (state.awarenessSocket === socket) {
        state.awarenessSocket = null;
      }
    });
  });
}

async function handleSyncSocketMessage(raw) {
  const payload = JSON.parse(raw);
  if (payload.type === "sync-feed") {
    setRelayConnected(true);
    for (const record of payload.messages || []) {
      if (record.participant_id === state.participantId) {
        continue;
      }
      if (record.recipient_id && record.recipient_id !== state.participantId) {
        continue;
      }
      await receive(record);
    }
    const nextOffset = payload.next_offset || state.offset;
    if (nextOffset > state.offset) {
      state.offset = nextOffset;
    }
    return;
  }
  if (payload.type === "sync-ready") {
    setRelayConnected(true);
    const nextOffset = payload.next_offset || state.offset;
    if (nextOffset > state.offset) {
      state.offset = nextOffset;
    }
    return;
  }
  if (payload.type === "sync-posted") {
    setRelayConnected(true);
    return;
  }
  if (payload.type === "error") {
    throw new Error(payload.message || "sync websocket error");
  }
}

function handleAwarenessSocketMessage(raw) {
  const payload = JSON.parse(raw);
  if (payload.type === "awareness-state") {
    setRelayConnected(true);
    const peers = [];
    for (const peer of payload.awareness || []) {
      if (!peer.participant_id || peer.participant_id === state.participantId) {
        continue;
      }
      peers.push({
        participant_id: peer.participant_id,
        name: peer.display_name || peer.author || peer.participant_id,
        color: peer.color || "#999999",
        anchor: peer.cursor || 0,
        head: peer.head || peer.cursor || 0,
        typing: Boolean(peer.typing),
        last_seen_at: peer.last_seen_at || null,
      });
    }
    send({ type: "awareness", peers });
    return;
  }
  if (payload.type === "error") {
    throw new Error(payload.message || "awareness websocket error");
  }
}

function toWebSocketURL(action, params = {}) {
  const url = new URL(`${basePath()}/${action}`);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null || value === "") {
      continue;
    }
    url.searchParams.set(key, String(value));
  }
  return url.toString();
}

function websocketCapable() {
  return typeof globalThis.WebSocket === "function";
}

function setRelayConnected(connected) {
  if (state.relayConnected === connected) {
    return;
  }
  state.relayConnected = connected;
  send({ type: "relay_status", connected });
}

function ensureDocument(doc) {
  if (doc?.content !== undefined) {
    return doc;
  }
  // Intent: Reuse one canonical serialized Automerge base across all sidecar
  // replicas so append-only change packets replay against matching initial
  // actor/object identities instead of per-process freshly minted ones.
  // Source: DI-sulod; DI-larok
  return Automerge.load(EMPTY_DOCUMENT_BYTES);
}

function getText() {
  return state.doc.content?.toString() || "";
}

function commonPrefix(left, right) {
  const size = Math.min(left.length, right.length);
  let index = 0;
  while (index < size && left[index] === right[index]) {
    index += 1;
  }
  return index;
}

function commonSuffix(left, right, prefix) {
  const leftRemain = left.length - prefix;
  const rightRemain = right.length - prefix;
  const size = Math.min(leftRemain, rightRemain);
  let index = 0;
  while (index < size && left[left.length - 1 - index] === right[right.length - 1 - index]) {
    index += 1;
  }
  return index;
}

function bytesToBase64(bytes) {
  return Buffer.from(bytes).toString("base64");
}

function base64ToBytes(value) {
  return Uint8Array.from(Buffer.from(value, "base64"));
}

sendInfo("grid-nvim-sidecar helper ready");
