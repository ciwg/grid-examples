import * as readline from "node:readline";
import * as http from "node:http";
import * as https from "node:https";
import * as Automerge from "@automerge/automerge";

const AutomergeNext = Automerge.next;
const EMPTY_DOCUMENT_BYTES = Uint8Array.from(Buffer.from("hW9Kg8HDZmEAdQEQUDnUuZsuTLOKK6EtAqSUwAF91ThR16b5XY1P61eTHXkwnJNTicqZ35V+jMImBQWmigYBAgMCEwIjBkACVgIHFQkhAiMCNAFCAlYCgAECfwB/AX8Bf8660dIGfwB/B38HY29udGVudH8AfwEBfwR/AH8AAA==", "base64"));

const state = {
  relayUrl: "",
  accessToken: "",
  capabilities: {},
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
  initialSyncReady: false,
  startupTransportsReady: false,
  openedSent: false,
  queuedLocalText: null,
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
      state.accessToken = message.access_token || "";
      state.capabilities = {};
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
  state.initialSyncReady = false;
  state.startupTransportsReady = false;
  state.openedSent = false;
  state.queuedLocalText = null;
  if (state.accessToken) {
    const session = await requestJSON("POST", `${basePath()}/session`, {
      participant_id: state.participantId,
    }, {
      "X-Grid-Access-Token": state.accessToken,
    });
    state.capabilities = session.capabilities || {};
  }
  await hydrateFromSnapshot();
  // Intent: Move the Neovim embodiment onto ex3's websocket live transport
  // while preserving the same stdin/stdout sidecar contract that the plugin
  // already speaks. Source: DI-bitus
  if (websocketCapable()) {
    await connectSyncSocket();
    await connectAwarenessSocket();
    state.startupTransportsReady = true;
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
    state.initialSyncReady = true;
    state.startupTransportsReady = true;
  }
  completeInitialOpen();
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
  state.initialSyncReady = false;
  state.startupTransportsReady = false;
  state.openedSent = false;
  state.queuedLocalText = null;
  setRelayConnected(false);
}

function applyLocalText(content) {
  const previous = getText();
  if (content === previous) {
    return;
  }
  if (!state.initialSyncReady) {
    // Intent: Hold Neovim-originated text edits until the sidecar has finished
    // its initial relay catch-up so opening an existing shared document cannot
    // overwrite remote content with an empty local buffer. Source: DI-gafit
    state.queuedLocalText = content;
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

function completeInitialOpen() {
  if (state.openedSent || !state.documentId || !state.initialSyncReady || !state.startupTransportsReady) {
    return;
  }
  state.openedSent = true;
  send({
    type: "opened",
    doc_id: state.documentId,
    content: getText(),
    relay_transport: state.relayTransport,
    awareness_transport: state.awarenessTransport,
  });
  if (state.queuedLocalText !== null) {
    const queued = state.queuedLocalText;
    state.queuedLocalText = null;
    applyLocalText(queued);
  }
}

async function pollSync() {
  if (!state.documentId) {
    return;
  }
  while (true) {
    const payload = await getJSON(`${basePath()}/sync?since=${state.offset}&limit=256`);
    setRelayConnected(true);
    await receiveMany(payload.messages || []);
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
  const replicaBase64 = bytesToBase64(Automerge.save(state.doc));
  const textBase64 = bytesToBase64(new TextEncoder().encode(getText()));
  // Intent: Send sidecar live-document traffic over websocket in ex3 without
  // changing the signed change payloads or the relay's feed semantics. Source:
  // DI-bitus
  if (state.syncSocket && state.syncSocket.readyState === WebSocket.OPEN) {
    state.syncSocket.send(JSON.stringify({
      type: "post-sync",
      participant_id: state.participantId,
      recipient_id: "",
      message_base64: bytesToBase64(changeBytes),
      text_base64: textBase64,
      replica_base64: replicaBase64,
      embodiment: "nvim",
    }));
    return;
  }
  await postJSON(`${basePath()}/sync`, {
    participant_id: state.participantId,
    recipient_id: "",
    message_base64: bytesToBase64(changeBytes),
    text_base64: textBase64,
    replica_base64: replicaBase64,
    embodiment: "nvim",
  });
}

async function getJSON(url) {
  return requestJSON("GET", url);
}

async function postJSON(url, body) {
  return requestJSON("POST", url, body);
}

function requestJSON(method, rawURL, body, extraHeaders = {}) {
  return new Promise((resolve, reject) => {
    const url = new URL(rawURL);
    const client = url.protocol === "https:" ? https : http;
    const payload = body ? JSON.stringify(body) : null;
    const authHeaders = mutationHeaders(method, url.pathname);
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
            ...authHeaders,
            ...extraHeaders,
          }
        : (Object.keys(authHeaders).length > 0 || Object.keys(extraHeaders).length > 0)
          ? {
              ...authHeaders,
              ...extraHeaders,
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
    return false;
  }
  if (record.participant_id === state.participantId) {
    return false;
  }
  if (record.recipient_id && record.recipient_id !== state.participantId) {
    return false;
  }
  const [nextDoc] = Automerge.applyChanges(state.doc, [base64ToBytes(record.message_base64)]);
  const nextText = nextDoc.content?.toString() || "";
  if (nextText !== getText()) {
    state.doc = ensureDocument(nextDoc);
    return true;
  }
  state.doc = ensureDocument(nextDoc);
  return false;
}

async function receiveMany(records) {
  let changed = false;
  for (const record of records) {
    // Intent: Treat the relay as an exchange surface for peer sync, not as a
    // loopback transport for this participant's own signed messages, and apply
    // long replay batches against the current replica instead of cloning the
    // full document for every historical record. Source: DI-sulod; DI-gafit
    if (await receive(record)) {
      changed = true;
    }
  }
  if (changed && state.openedSent) {
    send({
      type: "changed",
      content: getText(),
    });
  }
}

async function hydrateFromSnapshot() {
  const snapshot = await getJSON(`${basePath()}/state`);
  if (!snapshot.replica_base64) {
    return;
  }
  if (!snapshot.text_base64 && Number(snapshot.message_count || 0) > 0) {
    return;
  }
  // Intent: Let late-joining Neovim clients start from the relay's latest
  // replica snapshot so long-lived demo documents do not need a full history
  // replay before cursor/presence becomes usable. Source: DI-gafit
  state.doc = ensureDocument(Automerge.load(base64ToBytes(snapshot.replica_base64)));
  state.offset = Number(snapshot.snapshot_offset || snapshot.next_offset || 0);
}

async function handleSyncSocketRecords(records) {
  if (!Array.isArray(records) || records.length === 0) {
    return;
  }
  await receiveMany(records);
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
    let socketWork = Promise.resolve();
    socket.addEventListener("open", () => {
      setRelayConnected(true);
      if (state.capabilities.sync) {
        socket.send(JSON.stringify({
          type: "auth",
          capability: state.capabilities.sync,
        }));
      }
    });
    socket.addEventListener("message", (event) => {
      // Intent: Process websocket sync frames in arrival order so `sync-ready`
      // cannot race ahead of earlier `sync-feed` batches while a long document
      // history is still being applied. Source: DI-gafit
      socketWork = socketWork.then(() => handleSyncSocketMessage(event.data));
      socketWork
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
    let socketWork = Promise.resolve();
    socket.addEventListener("open", () => {
      if (state.capabilities.awareness) {
        socket.send(JSON.stringify({
          type: "auth",
          capability: state.capabilities.awareness,
        }));
      }
      Promise.resolve(postAwareness(false)).catch((error) => {
        if (!settled) {
          settled = true;
          reject(error);
        }
      });
    });
    socket.addEventListener("message", (event) => {
      socketWork = socketWork.then(() => {
        handleAwarenessSocketMessage(event.data);
      });
      socketWork.then(() => {
        if (!settled) {
          settled = true;
          resolve();
        }
      }).catch((error) => {
        if (!settled) {
          settled = true;
          reject(error);
          return;
        }
        send({ type: "error", message: error.stack || error.message });
      });
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
    await handleSyncSocketRecords(payload.messages || []);
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
    state.initialSyncReady = true;
    completeInitialOpen();
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

function mutationHeaders(method, pathname) {
  if (method !== "POST") {
    return {};
  }
  let capability = "";
  if (pathname.endsWith("/sync")) {
    capability = state.capabilities.sync || "";
  } else if (pathname.endsWith("/awareness")) {
    capability = state.capabilities.awareness || "";
  } else if (pathname.endsWith("/metadata")) {
    capability = state.capabilities.metadata || "";
  } else if (pathname.endsWith("/publish")) {
    capability = state.capabilities.publish || "";
  }
  if (!capability) {
    return {};
  }
  return {
    Authorization: `Bearer ${capability}`,
  };
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
