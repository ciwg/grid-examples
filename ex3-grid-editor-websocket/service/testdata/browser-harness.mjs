import * as readline from "node:readline";

if (!globalThis.window) {
  globalThis.window = globalThis;
}
if (!globalThis.window.setInterval) {
  globalThis.window.setInterval = setInterval;
}
if (!globalThis.window.clearInterval) {
  globalThis.window.clearInterval = clearInterval;
}
if (!globalThis.window.btoa) {
  globalThis.window.btoa = (value) => Buffer.from(value, "binary").toString("base64");
}
if (!globalThis.window.atob) {
  globalThis.window.atob = (value) => Buffer.from(value, "base64").toString("binary");
}

const { RelayAwarenessClient } = await import("../../web/src/relay-awareness.js");
const { AutomergeRelayClient } = await import("../../web/src/automerge-relay.js");

const state = {
  awareness: null,
  relay: null,
  participantId: "",
  documentId: "",
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
  pendingMessage = pendingMessage.then(() => handleMessage(message)).catch((error) => {
    send({ type: "error", message: error.stack || error.message });
  });
});

function send(message) {
  process.stdout.write(`${JSON.stringify(message)}\n`);
}

async function handleMessage(message) {
  switch (message.type) {
    case "connect":
      await connect(message);
      break;
    case "set_text":
      applyLocalText(message.content || "");
      break;
    case "set_cursor":
      if (!state.awareness) {
        return;
      }
      state.awareness.updateSelection(Number(message.anchor) || 0, Number(message.head) || 0);
      state.awareness.setTyping(Boolean(message.typing));
      break;
    case "snapshot":
      sendSnapshot();
      break;
    case "close":
      close();
      send({ type: "closed" });
      break;
    default:
      send({ type: "error", message: `unknown message type ${message.type}` });
  }
}

async function connect(message) {
  close();
  state.participantId = message.participant_id;
  state.documentId = message.doc_id || "demo";
  const basePath = `${message.relay_url.replace(/\/$/, "")}/api/local/documents/${encodeURIComponent(state.documentId)}`;
  state.awareness = new RelayAwarenessClient({
    basePath,
    participantID: state.participantId,
    documentID: state.documentId,
    displayName: message.display_name || "Browser User",
    color: message.color || "#1d6fd6",
  });
  state.relay = new AutomergeRelayClient({
    basePath,
    participantID: state.participantId,
    documentID: state.documentId,
    awareness: state.awareness,
  });
  state.relay.on("document", (text) => {
    send({ type: "document", content: text });
  });
  state.relay.on("error", (error) => {
    send({ type: "error", message: error.message || String(error) });
  });
  state.awareness.on("change", () => {
    send({ type: "awareness", peers: peerSnapshot() });
  });
  state.awareness.on("error", (error) => {
    send({ type: "error", message: error.message || String(error) });
  });
  await state.awareness.connect();
  await state.relay.connect();
  send({
    type: "opened",
    doc_id: state.documentId,
    content: state.relay.getText(),
    relay_transport: state.relay.transportMode(),
    awareness_transport: state.awareness.transportMode(),
  });
}

function applyLocalText(content) {
  if (!state.relay) {
    return;
  }
  const previous = state.relay.getText();
  if (content === previous) {
    return;
  }
  const prefix = commonPrefix(previous, content);
  const suffix = commonSuffix(previous, content, prefix);
  const deleteCount = previous.length - prefix - suffix;
  const insertText = content.slice(prefix, content.length - suffix);
  state.relay.applyLocalUpdate({
    changes: {
      iterChanges(callback) {
        callback(prefix, prefix + deleteCount, 0, 0, {
          toString() {
            return insertText;
          },
        });
      },
    },
  });
  send({ type: "local_change", content: state.relay.getText() });
}

function sendSnapshot() {
  send({
    type: "snapshot",
    content: state.relay ? state.relay.getText() : "",
    peers: peerSnapshot(),
  });
}

function peerSnapshot() {
  if (!state.awareness) {
    return [];
  }
  const states = state.awareness.getStates();
  const peers = [];
  states.forEach((value, key) => {
    peers.push({
      participant_id: key,
      name: value.user?.name || key,
      color: value.user?.color || "#999999",
      anchor: value.selection?.anchor || 0,
      head: value.selection?.head || 0,
      typing: Boolean(value.typing),
      last_seen_at: value.lastSeenAt || null,
      embodiment: value.embodiment || "",
    });
  });
  peers.sort((left, right) => left.participant_id.localeCompare(right.participant_id));
  return peers;
}

function close() {
  state.relay?.disconnect();
  state.awareness?.disconnect();
  state.relay = null;
  state.awareness = null;
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

send({ type: "ready" });
