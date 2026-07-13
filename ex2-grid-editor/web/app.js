const metaEls = {
  localID: document.getElementById("local-id"),
  docPCID: document.getElementById("doc-pcid"),
  awarenessPCID: document.getElementById("awareness-pcid"),
};

const editor = document.getElementById("editor");
const statusEl = document.getElementById("status");
const revisionEl = document.getElementById("revision");
const contentCIDEl = document.getElementById("content-cid");
const peerListEl = document.getElementById("peer-list");
const docIDInput = document.getElementById("doc-id");
const displayNameInput = document.getElementById("display-name");
const colorInput = document.getElementById("color");

const state = {
  documentID: new URLSearchParams(window.location.search).get("doc") || docIDInput.value || "demo",
  localID: "",
  lastMessageCID: "",
  typingTimer: null,
  replaceTimer: null,
};

docIDInput.value = state.documentID;

async function getJSON(path) {
  const response = await fetch(path);
  if (!response.ok) {
    throw new Error(`GET ${path} failed: ${response.status}`);
  }
  return response.json();
}

async function postJSON(path, body) {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!response.ok) {
    const text = await response.text();
    throw new Error(`POST ${path} failed: ${response.status} ${text}`);
  }
  return response.json();
}

async function loadMeta() {
  const meta = await getJSON("/api/meta");
  state.localID = meta.local_id;
  metaEls.localID.textContent = meta.local_id;
  metaEls.docPCID.textContent = meta.document_pcid;
  metaEls.awarenessPCID.textContent = meta.awareness_pcid;
}

function scheduleReplace() {
  clearTimeout(state.replaceTimer);
  state.replaceTimer = setTimeout(async () => {
    try {
      await postJSON(`/api/local/documents/${encodeURIComponent(state.documentID)}/replace`, {
        content: editor.value,
        embodiment: "browser",
      });
      statusEl.textContent = "saved locally";
    } catch (error) {
      statusEl.textContent = error.message;
    }
  }, 250);
}

function scheduleTypingStop() {
  clearTimeout(state.typingTimer);
  state.typingTimer = setTimeout(() => sendAwareness(false), 1200);
}

async function sendAwareness(typing = false) {
  const anchor = editor.selectionStart || 0;
  const head = editor.selectionEnd || anchor;
  try {
    await postJSON(`/api/local/documents/${encodeURIComponent(state.documentID)}/awareness`, {
      cursor: anchor,
      head,
      typing,
      display_name: displayNameInput.value || "Browser User",
      color: colorInput.value || "#1d6fd6",
      embodiment: "browser",
    });
  } catch (error) {
    statusEl.textContent = error.message;
  }
}

function renderPeers(peers) {
  peerListEl.innerHTML = "";
  const remotePeers = peers.filter((peer) => peer.author !== state.localID);
  if (remotePeers.length === 0) {
    const li = document.createElement("li");
    li.textContent = "No remote peers yet";
    li.className = "muted";
    peerListEl.appendChild(li);
    return;
  }
  for (const peer of remotePeers) {
    const li = document.createElement("li");
    const swatch = document.createElement("span");
    swatch.className = "swatch";
    swatch.style.background = peer.color || "#999999";
    const text = document.createElement("span");
    text.textContent = `${peer.display_name || peer.author} @ ${peer.cursor}`;
    li.appendChild(swatch);
    li.appendChild(text);
    peerListEl.appendChild(li);
  }
}

async function pollState() {
  try {
    const response = await getJSON(`/api/local/documents/${encodeURIComponent(state.documentID)}/state`);
    if (response.message_cid !== state.lastMessageCID && typeof response.content === "string" && response.content !== editor.value) {
      editor.value = response.content;
      state.lastMessageCID = response.message_cid || "";
    }
    revisionEl.textContent = `revision: ${response.lamport || 0}`;
    contentCIDEl.textContent = `content: ${response.content_cid || "-"}`;
    renderPeers(response.awareness || []);
    statusEl.textContent = "connected";
  } catch (error) {
    statusEl.textContent = error.message;
  }
}

document.getElementById("open-doc").addEventListener("click", () => {
  state.documentID = docIDInput.value.trim() || "demo";
  window.history.replaceState(null, "", `/?doc=${encodeURIComponent(state.documentID)}`);
  state.lastMessageCID = "";
  pollState();
});

editor.addEventListener("input", () => {
  scheduleReplace();
  sendAwareness(true);
  scheduleTypingStop();
});

editor.addEventListener("keyup", () => sendAwareness(false));
editor.addEventListener("mouseup", () => sendAwareness(false));
editor.addEventListener("focus", () => sendAwareness(false));

window.setInterval(pollState, 900);

loadMeta().then(pollState).catch((error) => {
  statusEl.textContent = error.message;
});

