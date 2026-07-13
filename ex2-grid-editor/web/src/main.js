import { createEditor } from "./editor.js";
import { RelayAwarenessClient } from "./relay-awareness.js";
import { AutomergeRelayClient } from "./automerge-relay.js";

const metaEls = {
  localID: document.getElementById("local-id"),
  docPCID: document.getElementById("doc-pcid"),
  awarenessPCID: document.getElementById("awareness-pcid"),
};

const statusEl = document.getElementById("status");
const revisionEl = document.getElementById("revision");
const contentCIDEl = document.getElementById("content-cid");
const peerListEl = document.getElementById("peer-list");
const peerBadgesEl = document.getElementById("peer-badges");
const docIDInput = document.getElementById("doc-id");
const displayNameInput = document.getElementById("display-name");
const colorInput = document.getElementById("color");
const editorRoot = document.getElementById("editor");

const state = {
  documentID: new URLSearchParams(window.location.search).get("doc") || "demo",
  participantID: getOrCreateParticipantID(),
  editor: null,
  awareness: null,
  relay: null,
};

docIDInput.value = state.documentID;
displayNameInput.value = getStoredPreference("display-name", "Browser User");
colorInput.value = getStoredPreference("color", "#1d6fd6");

async function loadMeta() {
  const response = await fetch("/api/meta");
  const meta = await response.json();
  metaEls.localID.textContent = meta.local_id;
  metaEls.docPCID.textContent = meta.document_pcid;
  metaEls.awarenessPCID.textContent = meta.awareness_pcid;
}

async function bootDocument(documentID) {
  statusEl.textContent = "connecting…";
  state.documentID = documentID;
  window.history.replaceState(null, "", `/?doc=${encodeURIComponent(documentID)}`);

  state.relay?.disconnect();
  state.awareness?.disconnect();
  state.editor?.destroy();
  editorRoot.innerHTML = "";

  const basePath = `/api/local/documents/${encodeURIComponent(documentID)}`;
  const awareness = new RelayAwarenessClient({
    basePath,
    participantID: state.participantID,
    documentID,
    displayName: displayNameInput.value || "Browser User",
    color: colorInput.value || "#1d6fd6",
  });
  const relay = new AutomergeRelayClient({
    basePath,
    participantID: state.participantID,
    documentID,
    awareness,
  });
  const editor = createEditor(
    editorRoot,
    awareness,
    state.participantID,
    (update) => {
      relay.applyLocalUpdate(update);
      awareness.setTyping(true);
      scheduleTypingStop(awareness);
      contentCIDEl.textContent = `local replica: ${relay.getReplicaCID()}`;
    },
    (anchor, head) => {
      awareness.updateSelection(anchor, head);
    },
  );

  state.awareness = awareness;
  state.relay = relay;
  state.editor = editor;
  editor.setText(relay.getText());
  contentCIDEl.textContent = `local replica: ${relay.getReplicaCID()}`;

  relay.on("document", (text) => {
    editor.setText(text);
    contentCIDEl.textContent = `local replica: ${relay.getReplicaCID()}`;
  });
  relay.on("error", (error) => {
    statusEl.textContent = error.message;
  });
  awareness.on("error", (error) => {
    statusEl.textContent = error.message;
  });
  awareness.on("change", () => {
    renderPeers(awareness.getStates());
    const peers = Array.from(awareness.getStates().keys()).filter((id) => id !== state.participantID);
    relay.observePeers(peers.map((participantID) => ({ participant_id: participantID })));
  });

  await awareness.connect();
  await relay.connect();
  renderPeers(awareness.getStates());
  statusEl.textContent = "connected";
  await refreshState(basePath);
}

async function refreshState(basePath) {
  const response = await fetch(`${basePath}/state`);
  if (!response.ok) {
    statusEl.textContent = `state GET failed: ${response.status}`;
    return;
  }
  const payload = await response.json();
  revisionEl.textContent = `messages: ${payload.message_count || 0}`;
}

function renderPeers(states) {
  peerListEl.innerHTML = "";
  peerBadgesEl.innerHTML = "";
  const remotePeers = Array.from(states.entries()).filter(([participantID]) => participantID !== state.participantID);
  if (remotePeers.length === 0) {
    const li = document.createElement("li");
    li.className = "muted";
    li.textContent = "No remote peers yet";
    peerListEl.appendChild(li);
    return;
  }
  for (const [participantID, peer] of remotePeers) {
    const name = peer.user?.name || participantID;
    const color = peer.user?.color || "#999999";
    const cursor = peer.selection?.anchor ?? 0;

    const li = document.createElement("li");
    li.innerHTML = `<span class="swatch" style="background:${color}"></span><span>${name} @ ${cursor}</span>`;
    peerListEl.appendChild(li);

    const badge = document.createElement("div");
    badge.className = "peer-badge";
    badge.innerHTML = `<span class="swatch" style="background:${color}"></span><strong>${name}</strong><span>cursor ${cursor}</span>`;
    peerBadgesEl.appendChild(badge);
  }
}

function scheduleTypingStop(awareness) {
  window.clearTimeout(scheduleTypingStop.timer);
  scheduleTypingStop.timer = window.setTimeout(() => {
    awareness.setTyping(false);
  }, 800);
}

function getOrCreateParticipantID() {
  const key = "grid-editor-participant-id";
  const existing = window.sessionStorage.getItem(key);
  if (existing) {
    return existing;
  }
  const created = `browser-${crypto.randomUUID()}`;
  window.sessionStorage.setItem(key, created);
  return created;
}

function getStoredPreference(key, fallback) {
  return window.localStorage.getItem(`grid-editor-${key}`) || fallback;
}

function setStoredPreference(key, value) {
  window.localStorage.setItem(`grid-editor-${key}`, value);
}

document.getElementById("open-doc").addEventListener("click", () => {
  bootDocument(docIDInput.value.trim() || "demo").catch((error) => {
    statusEl.textContent = error.message;
  });
});

displayNameInput.addEventListener("change", () => {
  setStoredPreference("display-name", displayNameInput.value || "Browser User");
  state.awareness?.setName(displayNameInput.value || "Browser User");
});

colorInput.addEventListener("change", () => {
  setStoredPreference("color", colorInput.value || "#1d6fd6");
  state.awareness?.setColor(colorInput.value || "#1d6fd6");
});

loadMeta()
  .then(() => bootDocument(state.documentID))
  .catch((error) => {
    statusEl.textContent = error.message;
  });
