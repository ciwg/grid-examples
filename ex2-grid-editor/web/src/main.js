import { createEditor } from "./editor.js";
import { RelayAwarenessClient } from "./relay-awareness.js";
import { AutomergeRelayClient } from "./automerge-relay.js";
import { PreferencesStore, defaultShortcuts, formatShortcut } from "./preferences.js";

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
const peerCountEl = document.getElementById("peer-count");
const presenceLegendEl = document.getElementById("presence-legend");
const profilePillEl = document.getElementById("profile-pill");
const shareLinkEl = document.getElementById("share-link");
const participantEl = document.getElementById("participant-id");
const toastStackEl = document.getElementById("toast-stack");
const helpShortcutsEl = document.getElementById("help-shortcuts");
const welcomeBannerEl = document.getElementById("welcome-banner");
const docIDInput = document.getElementById("doc-id");
const displayNameInput = document.getElementById("display-name");
const colorInput = document.getElementById("color");
const editorRoot = document.getElementById("editor");

const settingsPanel = document.getElementById("settings-panel");
const helpPanel = document.getElementById("help-panel");
const settingsFields = {
  theme: document.getElementById("theme-select"),
  lineNumbers: document.getElementById("line-numbers-toggle"),
  fontSize: document.getElementById("font-size-range"),
  dyslexiaMode: document.getElementById("dyslexia-toggle"),
  profile: document.getElementById("profile-select"),
  shortcuts: {
    search: document.getElementById("shortcut-search"),
    bold: document.getElementById("shortcut-bold"),
    italic: document.getElementById("shortcut-italic"),
    underline: document.getElementById("shortcut-underline"),
    settings: document.getElementById("shortcut-settings"),
    help: document.getElementById("shortcut-help"),
  },
};

const preferences = new PreferencesStore();

const state = {
  documentID: new URLSearchParams(window.location.search).get("doc") || "demo",
  participantID: getOrCreateParticipantID(),
  editor: null,
  awareness: null,
  relay: null,
  prefs: preferences.get(),
  visiblePeers: new Map(),
};

participantEl.textContent = state.participantID;
docIDInput.value = state.documentID;

// Intent: Keep the first Phase 1 browser preferences local while routing every
// UI surface through one abstraction, so later PromiseGrid sync does not have
// to be threaded back through ad hoc storage calls. Source: DI-vasul
function applyPreferences(nextPrefs, options = {}) {
  state.prefs = nextPrefs;
  document.body.dataset.theme = nextPrefs.theme;
  document.body.dataset.dyslexia = String(nextPrefs.dyslexiaMode);
  document.documentElement.style.setProperty("--editor-size", `${nextPrefs.fontSize}px`);
  displayNameInput.value = nextPrefs.displayName;
  colorInput.value = nextPrefs.color;
  profilePillEl.textContent = `${nextPrefs.profile} profile`;
  if (!options.skipFormSync) {
    syncSettingsForm();
  }
  state.editor?.setLineNumbers(nextPrefs.lineNumbers);
  if (state.awareness) {
    state.awareness.setName(nextPrefs.displayName);
    state.awareness.setColor(nextPrefs.color);
  }
  renderHelp();
  renderPeers(state.awareness?.getStates() || new Map());
}

function syncSettingsForm() {
  settingsFields.theme.value = state.prefs.theme;
  settingsFields.lineNumbers.checked = state.prefs.lineNumbers;
  settingsFields.fontSize.value = String(state.prefs.fontSize);
  settingsFields.dyslexiaMode.checked = state.prefs.dyslexiaMode;
  settingsFields.profile.value = state.prefs.profile;
  for (const [action, field] of Object.entries(settingsFields.shortcuts)) {
    field.value = state.prefs.shortcuts[action] || defaultShortcuts()[action] || "";
  }
}

function updatePreferences(patch, options = {}) {
  const nextPrefs = preferences.update(patch);
  applyPreferences(nextPrefs, options);
}

async function loadMeta() {
  const response = await fetch("/api/meta");
  const meta = await response.json();
  metaEls.localID.textContent = meta.local_id;
  metaEls.docPCID.textContent = meta.document_pcid;
  metaEls.awarenessPCID.textContent = meta.awareness_pcid;
}

async function bootDocument(documentID) {
  setStatus("connecting", "connecting…");
  state.documentID = documentID;
  docIDInput.value = documentID;
  updateDocumentURL(documentID);
  updateShareLink(documentID);

  state.relay?.disconnect();
  state.awareness?.disconnect();
  state.editor?.destroy();
  state.visiblePeers = new Map();
  editorRoot.innerHTML = "";

  const basePath = `/api/local/documents/${encodeURIComponent(documentID)}`;
  const awareness = new RelayAwarenessClient({
    basePath,
    participantID: state.participantID,
    documentID,
    displayName: state.prefs.displayName,
    color: state.prefs.color,
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
  applyPreferences(state.prefs);
  editor.setText(relay.getText());
  contentCIDEl.textContent = `local replica: ${relay.getReplicaCID()}`;

  relay.on("document", (text) => {
    editor.setText(text);
    contentCIDEl.textContent = `local replica: ${relay.getReplicaCID()}`;
  });
  relay.on("status", (status) => {
    renderStatus(status);
  });
  relay.on("error", (error) => {
    setStatus("error", error.message);
    showToast(error.message);
  });
  awareness.on("error", (error) => {
    setStatus("error", error.message);
    showToast(error.message);
  });
  awareness.on("change", () => {
    const states = awareness.getStates();
    renderPeers(states);
    const peers = Array.from(states.keys()).filter((id) => id !== state.participantID);
    relay.observePeers(peers.map((participantID) => ({ participant_id: participantID })));
  });

  await awareness.connect();
  await relay.connect();
  renderPeers(awareness.getStates());
  await refreshState(basePath);
}

async function refreshState(basePath) {
  const response = await fetch(`${basePath}/state`);
  if (!response.ok) {
    setStatus("error", `state GET failed: ${response.status}`);
    return;
  }
  const payload = await response.json();
  revisionEl.textContent = `messages: ${payload.message_count || 0}`;
}

function renderPeers(states) {
  const remotePeers = Array.from(states.entries())
    .filter(([participantID]) => participantID !== state.participantID)
    .map(([participantID, peer]) => ({ participantID, ...peer }))
    .map((peer) => ({ ...peer, presenceState: presenceState(peer.lastSeenAt, state.prefs.profile) }))
    .filter((peer) => peer.presenceState !== "gone");

  emitPeerNotifications(remotePeers);
  peerListEl.innerHTML = "";
  peerBadgesEl.innerHTML = "";
  peerCountEl.textContent = `${remotePeers.length} peer${remotePeers.length === 1 ? "" : "s"}`;

  const counts = { live: 0, stale: 0, offline: 0 };
  for (const peer of remotePeers) {
    counts[peer.presenceState] += 1;
    const name = peer.user?.name || peer.participantID;
    const color = peer.user?.color || "#999999";
    const cursor = peer.selection?.anchor ?? 0;

    const li = document.createElement("li");
    li.innerHTML = `
      <span class="swatch" style="background:${color}"></span>
      <span class="peer-meta">
        <strong>${escapeHTML(name)}</strong>
        <span class="peer-state ${peer.presenceState}">${peer.presenceState}</span>
      </span>
      <span class="tiny muted">@ ${cursor}</span>
    `;
    peerListEl.appendChild(li);

    const badge = document.createElement("div");
    badge.className = "peer-badge";
    badge.dataset.presenceState = peer.presenceState;
    badge.innerHTML = `<span class="swatch" style="background:${color}"></span><strong>${escapeHTML(name)}</strong><span>${peer.presenceState}</span>`;
    peerBadgesEl.appendChild(badge);
  }

  if (remotePeers.length === 0) {
    const li = document.createElement("li");
    li.className = "muted";
    li.textContent = "No remote peers yet";
    peerListEl.appendChild(li);
  }

  presenceLegendEl.textContent = `Live ${counts.live} · Stale ${counts.stale} · Offline ${counts.offline}`;
}

function emitPeerNotifications(peers) {
  const next = new Map(peers.map((peer) => [peer.participantID, peer]));
  for (const [participantID, peer] of next.entries()) {
    if (!state.visiblePeers.has(participantID)) {
      showToast(`${peer.user?.name || participantID} joined ${state.documentID}`);
    }
  }
  for (const [participantID, peer] of state.visiblePeers.entries()) {
    if (!next.has(participantID)) {
      showToast(`${peer.user?.name || participantID} left ${state.documentID}`);
    }
  }
  state.visiblePeers = next;
}

function renderStatus(status) {
  if (status.phase === "disconnected") {
    showToast("relay disconnected; retrying");
  }
  if (status.phase === "ready" && status.connected) {
    setStatus("ready", "ready");
    return;
  }
  if (status.phase === "syncing" || status.phase === "unsynced") {
    setStatus(status.phase, status.phase === "unsynced" ? "unsynced local changes" : "syncing…");
    return;
  }
  if (status.phase === "connecting") {
    setStatus("connecting", "connecting…");
    return;
  }
  if (status.phase === "disconnected") {
    setStatus("disconnected", status.message || "relay disconnected");
    return;
  }
  if (status.phase === "error") {
    setStatus("error", status.message || "sync error");
  }
}

function setStatus(phase, text) {
  statusEl.textContent = text;
  statusEl.className = `status-pill status-${phase}`;
}

function openSettings() {
  syncSettingsForm();
  openOverlay(settingsPanel);
}

function renderHelp() {
  helpShortcutsEl.innerHTML = "";
  const shortcuts = state.prefs.shortcuts;
  for (const [action, label] of Object.entries({
    search: "Search",
    bold: "Bold",
    italic: "Italic",
    underline: "Underline",
    settings: "Settings",
    help: "Help",
  })) {
    const row = document.createElement("div");
    row.className = "card";
    row.innerHTML = `<strong>${label}</strong><div class="tiny muted">${escapeHTML(formatShortcut(shortcuts[action]))}</div>`;
    helpShortcutsEl.appendChild(row);
  }
}

function openHelp() {
  renderHelp();
  openOverlay(helpPanel);
}

function openOverlay(element) {
  element.classList.remove("hidden");
  element.setAttribute("aria-hidden", "false");
}

function closeOverlay(element) {
  element.classList.add("hidden");
  element.setAttribute("aria-hidden", "true");
}

function createNewDocument() {
  const nextDoc = `doc-${crypto.randomUUID().slice(0, 8)}`;
  bootDocument(nextDoc).catch((error) => setStatus("error", error.message));
}

function openFromPromptedLink() {
  const raw = window.prompt("Paste a grid-editor doc link or document ID");
  if (!raw) {
    return;
  }
  const docID = parseDocumentReference(raw);
  if (!docID) {
    showToast("Could not parse a document from that input");
    return;
  }
  bootDocument(docID).catch((error) => setStatus("error", error.message));
}

function parseDocumentReference(raw) {
  const trimmed = raw.trim();
  if (!trimmed) {
    return "";
  }
  try {
    const url = new URL(trimmed, window.location.origin);
    return url.searchParams.get("doc") || url.pathname.replaceAll("/", "");
  } catch {
    return trimmed;
  }
}

function updateDocumentURL(documentID) {
  window.history.replaceState(null, "", `/?doc=${encodeURIComponent(documentID)}`);
}

function updateShareLink(documentID) {
  shareLinkEl.textContent = `Current link: ${window.location.origin}/?doc=${encodeURIComponent(documentID)}`;
}

function scheduleTypingStop(awareness) {
  window.clearTimeout(scheduleTypingStop.timer);
  scheduleTypingStop.timer = window.setTimeout(() => {
    awareness.setTyping(false);
  }, 900);
}

function presenceState(lastSeenAt, profile) {
  if (!lastSeenAt) {
    return "live";
  }
  const ageMs = Date.now() - new Date(lastSeenAt).getTime();
  // Intent: Render awareness using the approved demo/normal lifecycle windows
  // so the main peer list answers "who is here now?" while still giving
  // demos enough time before a peer is dimmed or removed. Source: DI-mivor;
  // DI-vasul
  const thresholds = profile === "normal"
    ? { live: 60_000, stale: 5 * 60_000, offline: 15 * 60_000 }
    : { live: 5 * 60_000, stale: 15 * 60_000, offline: 30 * 60_000 };
  if (ageMs <= thresholds.live) {
    return "live";
  }
  if (ageMs <= thresholds.stale) {
    return "stale";
  }
  if (ageMs <= thresholds.offline) {
    return "offline";
  }
  return "gone";
}

function showToast(message) {
  const toast = document.createElement("div");
  toast.className = "toast";
  toast.textContent = message;
  toastStackEl.appendChild(toast);
  window.setTimeout(() => toast.remove(), 3200);
}

function applyFormat(action) {
  const wrappers = {
    bold: ["**", "**"],
    italic: ["*", "*"],
    underline: ["<u>", "</u>"],
  };
  const [prefix, suffix] = wrappers[action];
  state.editor?.wrapSelection(prefix, suffix);
}

function searchDocument() {
  const query = window.prompt("Find text");
  if (!query) {
    return;
  }
  const found = state.editor?.findNext(query);
  if (!found) {
    showToast(`No match for “${query}”`);
  }
}

function registerEvents() {
  document.getElementById("open-doc").addEventListener("click", () => {
    bootDocument(docIDInput.value.trim() || "demo").catch((error) => setStatus("error", error.message));
  });
  document.getElementById("new-doc").addEventListener("click", createNewDocument);
  document.getElementById("paste-link").addEventListener("click", openFromPromptedLink);
  document.getElementById("search-button").addEventListener("click", searchDocument);
  document.getElementById("bold-button").addEventListener("click", () => applyFormat("bold"));
  document.getElementById("italic-button").addEventListener("click", () => applyFormat("italic"));
  document.getElementById("underline-button").addEventListener("click", () => applyFormat("underline"));
  document.getElementById("settings-button").addEventListener("click", openSettings);
  document.getElementById("help-button").addEventListener("click", openHelp);
  document.getElementById("settings-close").addEventListener("click", () => closeOverlay(settingsPanel));
  document.getElementById("help-close").addEventListener("click", () => closeOverlay(helpPanel));
  document.getElementById("welcome-open-settings").addEventListener("click", openSettings);
  document.getElementById("welcome-dismiss").addEventListener("click", () => {
    window.localStorage.setItem("grid-editor-dismissed-welcome", "true");
    welcomeBannerEl.classList.add("hidden");
  });

  displayNameInput.addEventListener("change", () => {
    updatePreferences({ displayName: displayNameInput.value || "Browser User" }, { skipFormSync: true });
  });
  colorInput.addEventListener("change", () => {
    updatePreferences({ color: colorInput.value || "#1d6fd6" }, { skipFormSync: true });
  });

  settingsFields.theme.addEventListener("change", () => updatePreferences({ theme: settingsFields.theme.value }, { skipFormSync: true }));
  settingsFields.lineNumbers.addEventListener("change", () => updatePreferences({ lineNumbers: settingsFields.lineNumbers.checked }, { skipFormSync: true }));
  settingsFields.fontSize.addEventListener("input", () => updatePreferences({ fontSize: Number(settingsFields.fontSize.value) }, { skipFormSync: true }));
  settingsFields.dyslexiaMode.addEventListener("change", () => updatePreferences({ dyslexiaMode: settingsFields.dyslexiaMode.checked }, { skipFormSync: true }));
  settingsFields.profile.addEventListener("change", () => updatePreferences({ profile: settingsFields.profile.value }, { skipFormSync: true }));
  for (const [action, field] of Object.entries(settingsFields.shortcuts)) {
    field.addEventListener("change", () => updatePreferences({ shortcuts: { [action]: field.value.trim() } }, { skipFormSync: true }));
  }

  document.addEventListener("keydown", (event) => {
    if (isTypingTarget(event.target) && event.key !== "Escape") {
      return;
    }
    const shortcuts = state.prefs.shortcuts;
    if (matchesShortcut(event, shortcuts.search)) {
      event.preventDefault();
      searchDocument();
    } else if (matchesShortcut(event, shortcuts.bold)) {
      event.preventDefault();
      applyFormat("bold");
    } else if (matchesShortcut(event, shortcuts.italic)) {
      event.preventDefault();
      applyFormat("italic");
    } else if (matchesShortcut(event, shortcuts.underline)) {
      event.preventDefault();
      applyFormat("underline");
    } else if (matchesShortcut(event, shortcuts.settings)) {
      event.preventDefault();
      openSettings();
    } else if (matchesShortcut(event, shortcuts.help)) {
      event.preventDefault();
      openHelp();
    } else if (event.key === "Escape") {
      closeOverlay(settingsPanel);
      closeOverlay(helpPanel);
    }
  });
}

function matchesShortcut(event, shortcut) {
  if (!shortcut) {
    return false;
  }
  const parts = shortcut.toLowerCase().split("-").filter(Boolean);
  const key = parts.pop();
  const wantsMod = parts.includes("mod");
  const wantsShift = parts.includes("shift");
  const wantsAlt = parts.includes("alt");
  const wantsCtrl = parts.includes("ctrl");
  const modActive = navigator.platform.includes("Mac") ? event.metaKey : event.ctrlKey;
  if (wantsMod !== modActive) {
    return false;
  }
  if (wantsShift !== event.shiftKey) {
    return false;
  }
  if (wantsAlt !== event.altKey) {
    return false;
  }
  if (wantsCtrl && !event.ctrlKey) {
    return false;
  }
  return event.key.toLowerCase() === key;
}

function isTypingTarget(target) {
  return target instanceof HTMLElement && ["INPUT", "TEXTAREA", "SELECT"].includes(target.tagName);
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

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

registerEvents();
applyPreferences(state.prefs);
if (window.localStorage.getItem("grid-editor-dismissed-welcome") === "true") {
  welcomeBannerEl.classList.add("hidden");
}

loadMeta()
  .then(() => bootDocument(state.documentID))
  .catch((error) => {
    setStatus("error", error.message);
  });
