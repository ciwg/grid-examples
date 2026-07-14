import { createEditor } from "./editor.js";
import { RelayAwarenessClient } from "./relay-awareness.js";
import { AutomergeRelayClient } from "./automerge-relay.js";
import { PreferencesStore, defaultShortcuts, formatShortcut } from "./preferences.js";
import { DocumentRegistry, generateDemoText, templateCatalog } from "./documents.js";
import { extractHeadings, renderMarkdown } from "./markdown.js";
import { extractMentions, inferVersionName, summarizeDocument } from "./review.js";
import { buildExportArtifact, buildPublishSource, parsePublishedURL } from "./exchange.js";

const metaEls = {
  localID: document.getElementById("local-id"),
  docPCID: document.getElementById("doc-pcid"),
  awarenessPCID: document.getElementById("awareness-pcid"),
  publishPCID: document.getElementById("publish-pcid"),
};

const statusEl = document.getElementById("status");
const autosaveEl = document.getElementById("autosave-indicator");
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
const docTitleInput = document.getElementById("doc-title");
const displayNameInput = document.getElementById("display-name");
const colorInput = document.getElementById("color");
const editorRoot = document.getElementById("editor");
const previewPaneEl = document.getElementById("preview-pane");
const previewBodyEl = document.getElementById("preview-body");
const previewMetaEl = document.getElementById("preview-meta");
const docTabBarEl = document.getElementById("doc-tab-bar");
const recentDocsEl = document.getElementById("recent-docs");
const openTabsEl = document.getElementById("open-tabs");
const templateGalleryEl = document.getElementById("template-gallery");
const editorStageEl = document.getElementById("editor-stage");
const fileImportEl = document.getElementById("file-import");
const focusButtonEl = document.getElementById("focus-button");
const outlineListEl = document.getElementById("outline-list");
const savedVersionsEl = document.getElementById("saved-versions");
const recentParticipantsEl = document.getElementById("recent-participants");
const activityFeedEl = document.getElementById("activity-feed");
const publishedListEl = document.getElementById("published-list");

const timestampEls = {
  created: document.getElementById("doc-created"),
  viewed: document.getElementById("doc-viewed"),
  edited: document.getElementById("doc-edited"),
  exported: document.getElementById("doc-exported"),
};

const settingsPanel = document.getElementById("settings-panel");
const helpPanel = document.getElementById("help-panel");
const searchPanel = document.getElementById("search-panel");
const exportPanel = document.getElementById("export-panel");
const commentPanel = document.getElementById("comment-panel");
const summaryPanel = document.getElementById("summary-panel");
const debugPanel = document.getElementById("debug-panel");
const commentSelectionEl = document.getElementById("comment-selection");
const commentBodyEl = document.getElementById("comment-body");
const commentListEl = document.getElementById("comment-list");
const summaryTextEl = document.getElementById("summary-text");
const debugOutputEl = document.getElementById("debug-output");
const searchFields = {
  query: document.getElementById("search-query"),
  replace: document.getElementById("replace-query"),
  caseSensitive: document.getElementById("search-case"),
  regex: document.getElementById("search-regex"),
};
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
const registry = new DocumentRegistry();

const state = {
  documentID: new URLSearchParams(window.location.search).get("doc") || "demo",
  participantID: getOrCreateParticipantID(),
  editor: null,
  awareness: null,
  relay: null,
  prefs: preferences.get(),
  visiblePeers: new Map(),
  previewEnabled: false,
  splitEnabled: false,
  focusMode: false,
  pendingCommentSelection: null,
};

participantEl.textContent = state.participantID;
docIDInput.value = state.documentID;

// Intent: Keep the first workflow-heavy Phase 2 metadata, recent-doc, and
// export surfaces local while preserving a clean seam for later
// PromiseGrid-native publishing and document-exchange work. Source: DI-nuvif;
// DI-dovoz
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
  metaEls.publishPCID.textContent = meta.publish_pcid;
}

async function bootDocument(documentID) {
  setStatus("connecting", "connecting…");
  state.documentID = documentID;
  docIDInput.value = documentID;
  updateDocumentURL(documentID);
  updateShareLink(documentID);
  registry.openTab(documentID);
  const documentMeta = registry.touchViewed(documentID);
  docTitleInput.value = documentMeta.title;
  renderRegistry();
  renderDocumentMeta(documentMeta);
  renderReview();

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
      registry.touchEdited(documentID);
      renderDocumentMeta(registry.get(documentID));
      renderPreview();
      renderReview();
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
    updateDerivedTitle(text);
    renderPreview();
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
    noteRecentParticipants(states);
    const peers = Array.from(states.keys()).filter((id) => id !== state.participantID);
    relay.observePeers(peers.map((participantID) => ({ participant_id: participantID })));
  });

  await awareness.connect();
  await relay.connect();
  renderPeers(awareness.getStates());
  noteRecentParticipants(awareness.getStates());
  await refreshState(basePath);
  await applySeedIfNeeded(documentID);
  await refreshPublished(documentID);
  renderPreview();
  renderReview();
}

async function applySeedIfNeeded(documentID) {
  const seed = registry.seedContent(documentID);
  if (!seed) {
    return;
  }
  if (state.relay.getText() !== "") {
    registry.registerSeedContent(documentID, "");
    return;
  }
  state.relay.replaceText(seed);
  state.editor.setText(seed);
  registry.registerSeedContent(documentID, "");
  registry.touchEdited(documentID);
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
    autosaveEl.textContent = "auto-save synced";
    return;
  }
  if (status.phase === "syncing" || status.phase === "unsynced") {
    setStatus(status.phase, status.phase === "unsynced" ? "unsynced local changes" : "syncing…");
    autosaveEl.textContent = status.phase === "unsynced" ? "auto-save pending" : "auto-save syncing";
    return;
  }
  if (status.phase === "connecting") {
    setStatus("connecting", "connecting…");
    autosaveEl.textContent = "auto-save connecting";
    return;
  }
  if (status.phase === "disconnected") {
    setStatus("disconnected", status.message || "relay disconnected");
    autosaveEl.textContent = "auto-save paused";
    return;
  }
  if (status.phase === "error") {
    setStatus("error", status.message || "sync error");
    autosaveEl.textContent = "auto-save error";
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
  const items = {
    search: "Find / Replace",
    bold: "Bold",
    italic: "Italic",
    underline: "Underline",
    settings: "Settings",
    help: "Help",
  };
  for (const [action, label] of Object.entries(items)) {
    const row = document.createElement("div");
    row.className = "card";
    row.innerHTML = `<strong>${label}</strong><div class="tiny muted">${escapeHTML(formatShortcut(shortcuts[action]))}</div>`;
    helpShortcutsEl.appendChild(row);
  }
}

function renderRegistry() {
  renderDocumentList(openTabsEl, registry.listOpenTabs(), (documentMeta) => documentMeta.documentID === state.documentID ? "active" : "", (documentMeta) => bootDocument(documentMeta.documentID));
  renderDocumentList(recentDocsEl, registry.listRecent(), "", (documentMeta) => bootDocument(documentMeta.documentID));
  renderDocTabs();
  renderTemplateGallery();
}

// Intent: Keep the first Phase 3 review/history surfaces browser-local and
// optional so comments, outline, saved versions, and diagnostics can mature
// without changing the relay contract. Source: DI-safor; DI-lapek
function renderReview() {
  renderOutline();
  renderSavedVersions();
  renderRecentParticipants();
  renderActivity();
  renderComments();
}

async function refreshPublished(documentID) {
  const response = await fetch(`/api/local/documents/${encodeURIComponent(documentID)}/published`);
  if (!response.ok) {
    publishedListEl.innerHTML = "";
    appendEmptyState(publishedListEl, "Published exchanges unavailable");
    showToast(`Published exchanges failed: ${response.status}`);
    return;
  }
  const payload = await response.json();
  renderPublished(payload.published || []);
}

function renderDocumentList(target, documents, extraClass, onClick) {
  target.innerHTML = "";
  for (const documentMeta of documents) {
    const li = document.createElement("li");
    const button = document.createElement("button");
    button.className = extraClass;
    button.type = "button";
    button.textContent = `${documentMeta.title} (${documentMeta.documentID})`;
    button.addEventListener("click", () => onClick(documentMeta));
    li.appendChild(button);
    target.appendChild(li);
  }
  if (documents.length === 0) {
    const li = document.createElement("li");
    li.className = "tiny muted";
    li.textContent = "None yet";
    target.appendChild(li);
  }
}

function renderDocTabs() {
  docTabBarEl.innerHTML = "";
  for (const documentMeta of registry.listOpenTabs()) {
    const tab = document.createElement("button");
    tab.type = "button";
    tab.className = `doc-tab ${documentMeta.documentID === state.documentID ? "active" : ""}`;
    tab.innerHTML = `<strong>${escapeHTML(documentMeta.title)}</strong><span class="tiny muted">${escapeHTML(documentMeta.documentID)}</span>`;
    tab.addEventListener("click", () => bootDocument(documentMeta.documentID));
    docTabBarEl.appendChild(tab);
  }
}

function renderTemplateGallery() {
  templateGalleryEl.innerHTML = "";
  for (const template of templateCatalog()) {
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = template.label;
    button.addEventListener("click", () => createFromTemplate(template));
    templateGalleryEl.appendChild(button);
  }
}

function renderDocumentMeta(documentMeta) {
  docTitleInput.value = documentMeta.title;
  timestampEls.created.textContent = formatTime(documentMeta.createdAt);
  timestampEls.viewed.textContent = formatTime(documentMeta.lastViewedAt);
  timestampEls.edited.textContent = formatTime(documentMeta.lastEditedAt);
  timestampEls.exported.textContent = formatTime(documentMeta.lastExportedAt);
}

function renderOutline() {
  const headings = state.editor ? extractHeadings(state.editor.getText()) : [];
  outlineListEl.innerHTML = "";
  for (const heading of headings) {
    const li = document.createElement("li");
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = `${"".padStart((heading.level - 1) * 2, " ")}${heading.text}`;
    button.addEventListener("click", () => state.editor?.goToLine(heading.line));
    li.appendChild(button);
    outlineListEl.appendChild(li);
  }
  if (headings.length === 0) {
    appendEmptyState(outlineListEl, "No headings yet");
  }
}

function renderSavedVersions() {
  const versions = registry.listSavedVersions(state.documentID);
  savedVersionsEl.innerHTML = "";
  for (const version of versions) {
    const li = document.createElement("li");
    li.innerHTML = `<strong>${escapeHTML(version.name)}</strong><span class="tiny muted">${formatTime(version.createdAt)}</span>`;
    savedVersionsEl.appendChild(li);
  }
  if (versions.length === 0) {
    appendEmptyState(savedVersionsEl, "No saved versions yet");
  }
}

function renderRecentParticipants() {
  const participants = registry.listRecentParticipants(state.documentID);
  recentParticipantsEl.innerHTML = "";
  for (const participant of participants) {
    const li = document.createElement("li");
    li.innerHTML = `<span class="swatch" style="background:${participant.color || "#999999"}"></span><span><strong>${escapeHTML(participant.name || participant.participantID)}</strong><span class="tiny muted"> ${formatTime(participant.lastSeenAt)}</span></span>`;
    recentParticipantsEl.appendChild(li);
  }
  if (participants.length === 0) {
    appendEmptyState(recentParticipantsEl, "No participant history yet");
  }
}

function renderActivity() {
  const events = registry.listActivity(state.documentID);
  activityFeedEl.innerHTML = "";
  for (const event of events.slice(0, 12)) {
    const li = document.createElement("li");
    li.innerHTML = `<strong>${escapeHTML(event.label)}</strong><span class="tiny muted">${formatTime(event.at)}</span>`;
    activityFeedEl.appendChild(li);
  }
  if (events.length === 0) {
    appendEmptyState(activityFeedEl, "No activity yet");
  }
}

function renderComments() {
  const comments = registry.listComments(state.documentID);
  commentListEl.innerHTML = "";
  for (const comment of comments) {
    const li = document.createElement("li");
    li.className = "card";
    const reactions = (comment.reactions || []).map((reaction) => reaction.emoji).join(" ");
    li.innerHTML = `
      <strong>${escapeHTML(comment.authorName)}</strong>
      <div class="tiny muted">${formatTime(comment.createdAt)}${comment.resolved ? " · resolved" : ""}</div>
      <div>${escapeHTML(comment.selectionText || "")}</div>
      <div>${escapeHTML(comment.body)}</div>
      <div class="tiny muted">${reactions || "no reactions"}</div>
    `;
    li.addEventListener("click", () => {
      if (comment.selectionFrom != null && comment.selectionTo != null) {
        state.editor?.selectRange(comment.selectionFrom, comment.selectionTo);
      }
    });
    commentListEl.appendChild(li);
  }
  if (comments.length === 0) {
    appendEmptyState(commentListEl, "No comments yet");
  }
}

function renderPublished(records) {
  publishedListEl.innerHTML = "";
  for (const record of records.slice().reverse()) {
    const li = document.createElement("li");
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = `${record.title} (${record.source_kind})`;
    button.addEventListener("click", async () => {
      await copyText(`${window.location.origin}/api/published/${record.envelope_cid}`, "Exchange link copied");
    });
    li.appendChild(button);
    const meta = document.createElement("div");
    meta.className = "tiny muted";
    meta.textContent = formatTime(record.published_at);
    li.appendChild(meta);
    publishedListEl.appendChild(li);
  }
  if (records.length === 0) {
    appendEmptyState(publishedListEl, "No published exchanges yet");
  }
}

function noteRecentParticipants(states) {
  states.forEach((peer, participantID) => {
    if (participantID === state.participantID) {
      return;
    }
    registry.noteParticipant(state.documentID, {
      participantID,
      name: peer.user?.name || participantID,
      color: peer.user?.color || "#999999",
      lastSeenAt: peer.lastSeenAt || new Date().toISOString(),
    });
  });
  renderRecentParticipants();
}

function appendEmptyState(target, label) {
  const li = document.createElement("li");
  li.className = "tiny muted";
  li.textContent = label;
  target.appendChild(li);
}

function openCommentPanel() {
  const selection = state.editor?.getSelection();
  if (!selection) {
    return;
  }
  state.pendingCommentSelection = selection;
  commentSelectionEl.value = selection.text || "(no text selected)";
  commentBodyEl.value = "";
  renderComments();
  openOverlay(commentPanel);
  commentBodyEl.focus();
}

function saveComment() {
  const selection = state.pendingCommentSelection || state.editor?.getSelection();
  const body = commentBodyEl.value.trim();
  if (!selection || !body) {
    showToast("Select text and enter a comment");
    return;
  }
  const comment = {
    id: crypto.randomUUID(),
    authorName: state.prefs.displayName,
    authorID: state.participantID,
    createdAt: new Date().toISOString(),
    selectionFrom: selection.from,
    selectionTo: selection.to,
    selectionText: selection.text,
    body,
    resolved: false,
    reactions: [],
    mentions: extractMentions(body, buildParticipantIndex()),
  };
  registry.addComment(state.documentID, comment);
  for (const mention of comment.mentions) {
    showToast(`Mentioned @${mention.label}`);
  }
  renderReview();
  closeOverlay(commentPanel);
}

function resolveSelectedComments() {
  const selection = state.pendingCommentSelection || state.editor?.getSelection();
  if (!selection) {
    return;
  }
  for (const comment of registry.listComments(state.documentID)) {
    if (comment.selectionFrom === selection.from && comment.selectionTo === selection.to && !comment.resolved) {
      registry.toggleCommentResolved(state.documentID, comment.id, state.prefs.displayName);
      registry.addReaction(state.documentID, comment.id, {
        emoji: "✅",
        authorName: state.prefs.displayName,
        createdAt: new Date().toISOString(),
      });
    }
  }
  renderReview();
}

function saveNamedVersion() {
  if (!state.editor) {
    return;
  }
  const suggested = inferVersionName(docTitleInput.value || state.documentID, state.editor.getText());
  const name = window.prompt("Version name", suggested);
  if (!name) {
    return;
  }
  registry.addSavedVersion(state.documentID, {
    id: crypto.randomUUID(),
    name,
    createdAt: new Date().toISOString(),
    summary: summarizeDocument(state.editor.getText()),
    content: state.editor.getText(),
    replicaBase64: bytesToBase64(state.relay.getReplicaBytes()),
  });
  renderSavedVersions();
  renderActivity();
}

function openSummary() {
  if (!state.editor) {
    return;
  }
  summaryTextEl.textContent = summarizeDocument(state.editor.getText());
  openOverlay(summaryPanel);
}

function readSummary() {
  const text = summaryTextEl.textContent || "";
  if (!text) {
    return;
  }
  if (!("speechSynthesis" in window)) {
    showToast("Text-to-speech is not available here");
    return;
  }
  window.speechSynthesis.cancel();
  window.speechSynthesis.speak(new SpeechSynthesisUtterance(text));
}

function dictateComment() {
  const Recognition = window.SpeechRecognition || window.webkitSpeechRecognition;
  if (!Recognition) {
    showToast("Voice dictation is not available here");
    return;
  }
  const recognition = new Recognition();
  recognition.onresult = (event) => {
    const transcript = Array.from(event.results).map((result) => result[0].transcript).join(" ");
    commentBodyEl.value = commentBodyEl.value ? `${commentBodyEl.value} ${transcript}` : transcript;
  };
  recognition.start();
}

function toggleFocusMode() {
  state.focusMode = !state.focusMode;
  document.body.classList.toggle("focus-mode", state.focusMode);
  focusButtonEl.textContent = state.focusMode ? "Exit Focus" : "Focus";
}

function openDebugPanel() {
  debugOutputEl.textContent = JSON.stringify({
    documentID: state.documentID,
    profile: state.prefs.profile,
    participantID: state.participantID,
    revision: revisionEl.textContent,
    localReplica: contentCIDEl.textContent,
    comments: registry.listComments(state.documentID),
    activity: registry.listActivity(state.documentID).slice(0, 12),
    participants: registry.listRecentParticipants(state.documentID),
    versions: registry.listSavedVersions(state.documentID),
  }, null, 2);
  openOverlay(debugPanel);
}

function buildParticipantIndex() {
  const index = new Map();
  for (const participant of registry.listRecentParticipants(state.documentID)) {
    index.set(String(participant.name || "").toLowerCase(), participant.participantID);
  }
  return index;
}

function renderPreview() {
  if (!state.editor) {
    previewBodyEl.innerHTML = "";
    return;
  }
  previewBodyEl.innerHTML = renderMarkdown(state.editor.getText());
  const headings = extractHeadings(state.editor.getText());
  previewMetaEl.textContent = headings.length > 0 ? `${headings.length} headings` : "same document";
}

function openHelp() {
  renderHelp();
  openOverlay(helpPanel);
}

function openSearch() {
  openOverlay(searchPanel);
  searchFields.query.focus();
}

function openExport() {
  openOverlay(exportPanel);
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
  registry.openTab(nextDoc);
  bootDocument(nextDoc).catch((error) => setStatus("error", error.message));
}

function createFromTemplate(template) {
  const nextDoc = `doc-${crypto.randomUUID().slice(0, 8)}`;
  registry.registerSeedContent(nextDoc, template.content);
  registry.updateTitle(nextDoc, `${template.label}`);
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

function duplicateCurrentDocument() {
  if (!state.editor) {
    return;
  }
  const nextDoc = `doc-${crypto.randomUUID().slice(0, 8)}`;
  const duplicate = registry.duplicateDocument(state.documentID, nextDoc, state.editor.getText());
  bootDocument(duplicate.documentID).catch((error) => setStatus("error", error.message));
}

function updateDocumentURL(documentID) {
  window.history.replaceState(null, "", `/?doc=${encodeURIComponent(documentID)}${window.location.hash}`);
}

function updateShareLink(documentID) {
  shareLinkEl.textContent = `Current link: ${window.location.origin}/?doc=${encodeURIComponent(documentID)}`;
}

function copyShareLink() {
  copyText(`${window.location.origin}/?doc=${encodeURIComponent(state.documentID)}`, "Link copied");
}

function emailShareLink() {
  window.location.href = `mailto:?subject=${encodeURIComponent(`grid-editor document ${state.documentID}`)}&body=${encodeURIComponent(`${window.location.origin}/?doc=${encodeURIComponent(state.documentID)}`)}`;
}

function togglePreview() {
  state.previewEnabled = !state.previewEnabled;
  applyPaneMode();
}

function toggleSplit() {
  state.splitEnabled = !state.splitEnabled;
  state.previewEnabled = true;
  applyPaneMode();
}

function applyPaneMode() {
  previewPaneEl.classList.toggle("hidden", !state.previewEnabled);
  editorStageEl.classList.toggle("split", state.previewEnabled && state.splitEnabled);
  renderPreview();
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

function preserveEditorSelectionOnMouseDown(button) {
  if (!button) {
    return;
  }
  // Intent: Keep toolbar actions from stealing focus away from the CodeMirror
  // selection so stacked formatting actions like bold-plus-underline still
  // apply to the same selected text in the real browser UI. Source: DI-favok
  button.addEventListener("mousedown", (event) => {
    event.preventDefault();
  });
}

function runSearchReplace(replaceMode = false) {
  const query = searchFields.query.value;
  if (!query) {
    showToast("Enter text to search for");
    return;
  }
  const options = {
    caseSensitive: searchFields.caseSensitive.checked,
    regex: searchFields.regex.checked,
  };
  if (replaceMode) {
    const count = state.editor?.replaceAll(query, searchFields.replace.value, options) || 0;
    showToast(count > 0 ? `Replaced ${count} matches` : "No matches to replace");
    return;
  }
  const found = state.editor?.findNext(query, options);
  if (!found) {
    showToast(`No match for “${query}”`);
  }
}

function gotoLine() {
  const raw = window.prompt("Go to line");
  const line = Number(raw);
  if (!Number.isFinite(line) || line <= 0) {
    return;
  }
  state.editor?.goToLine(line);
  closeOverlay(searchPanel);
}

function importFile() {
  fileImportEl.click();
}

async function handleImportedFile(event) {
  const file = event.target.files?.[0];
  if (!file || !state.editor) {
    return;
  }
  if (file.type.startsWith("image/")) {
    const dataURL = await readFileAsDataURL(file);
    state.editor.insertAtCursor(`![${file.name}](${dataURL})`);
    showToast(`Inserted image attachment ${file.name}`);
  } else {
    const text = await file.text();
    state.relay.replaceText(text);
    state.editor.setText(text);
    renderPreview();
    showToast(`Imported ${file.name}`);
  }
  fileImportEl.value = "";
}

function exportDocument(format) {
  if (!state.editor || !state.relay) {
    return;
  }
  const title = safeFilename(docTitleInput.value || state.documentID);
  const text = state.editor.getText();
  const artifact = buildExportArtifact(
    format,
    docTitleInput.value || state.documentID,
    text,
    state.relay.getReplicaBytes(),
    renderMarkdown,
  );
  const blob = new Blob([artifact.body], { type: artifact.mime });
  const extension = artifact.extension;
  triggerDownload(blob, `${title}.${extension}`);
  registry.touchExported(state.documentID);
  renderDocumentMeta(registry.get(state.documentID));
  showToast(`Exported ${extension.toUpperCase()}`);
}

function publishSnapshot() {
  if (!state.editor) {
    return;
  }
  const snapshot = {
    id: crypto.randomUUID(),
    createdAt: new Date().toISOString(),
    title: docTitleInput.value || state.documentID,
    content: state.editor.getText(),
  };
  registry.addSnapshot(state.documentID, snapshot);
  showToast("Snapshot saved locally");
}

async function publishExchange() {
  if (!state.editor || !state.relay) {
    return;
  }
  // Intent: Keep publish/import as an explicit current-time handoff flow that
  // can target either the current editor state or a named saved version,
  // rather than conflating it with live CRDT sync. Source: DI-tavul; DI-gosaf
  const source = choosePublishSource();
  if (!source) {
    return;
  }
  const response = await fetch(`/api/local/documents/${encodeURIComponent(state.documentID)}/publish`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      participant_id: state.participantID,
      source_kind: source.sourceKind,
      source_version_id: source.sourceVersionID,
      source_version_name: source.sourceVersionName,
      title: source.title,
      summary: source.summary,
      text_base64: textToBase64(source.text),
      replica_base64: source.replicaBase64,
      embodiment: "browser",
    }),
  });
  if (!response.ok) {
    const detail = await response.text();
    showToast(`Publish failed: ${response.status}${detail ? ` ${detail}` : ""}`);
    return;
  }
  const payload = await response.json();
  await copyText(payload.manifest_url, "Exchange link copied");
  await refreshPublished(state.documentID);
}

function choosePublishSource() {
  const currentText = state.editor.getText();
  const versions = registry.listSavedVersions(state.documentID);
  const raw = window.prompt("Publish 'current' or enter a saved version name", "current");
  const source = buildPublishSource(
    currentText,
    bytesToBase64(state.relay.getReplicaBytes()),
    docTitleInput.value || state.documentID,
    versions,
    raw,
  );
  if (!source) {
    showToast("Saved version not found or missing content");
    return null;
  }
  return source;
}

async function importExchange() {
  // Intent: Import published exchange artifacts by manifest URL into a new
  // local document so Grid Editor can round-trip durable handoff objects
  // without rewriting past live relay history. Source: DI-tavul; DI-gosaf
  const raw = window.prompt("Paste a published exchange URL");
  if (!raw) {
    return;
  }
  const reference = parsePublishedURL(raw, window.location.origin);
  if (!reference) {
    showToast("Could not parse an exchange URL");
    return;
  }
  const response = await fetch(`${reference.origin}/api/published/${reference.envelopeCID}`);
  if (!response.ok) {
    showToast(`Import failed: ${response.status}`);
    return;
  }
  const payload = await response.json();
  const nextDoc = `doc-${crypto.randomUUID().slice(0, 8)}`;
  registry.registerSeedContent(nextDoc, base64ToText(payload.text_base64));
  registry.updateTitle(nextDoc, `${payload.record.title} imported`);
  await bootDocument(nextDoc);
  showToast(`Imported ${payload.record.title}`);
}

function exportAuditReport() {
  const documentMeta = registry.get(state.documentID);
  const report = {
    document: documentMeta,
    shareLink: `${window.location.origin}/?doc=${encodeURIComponent(state.documentID)}`,
    localID: metaEls.localID.textContent,
    protocol: {
      documentPCID: metaEls.docPCID.textContent,
      awarenessPCID: metaEls.awarenessPCID.textContent,
    },
    generatedAt: new Date().toISOString(),
  };
  triggerDownload(new Blob([JSON.stringify(report, null, 2)], { type: "application/json;charset=utf-8" }), `${safeFilename(documentMeta.title)}-audit.json`);
  registry.touchExported(state.documentID);
  renderDocumentMeta(registry.get(state.documentID));
}

function copyFormatted(format) {
  if (!state.editor) {
    return;
  }
  const text = state.editor.getText();
  const value = format === "html" ? wrapHTML(docTitleInput.value || state.documentID, renderMarkdown(text)) : text;
  copyText(value, `Copied ${format}`);
}

function addBookmark() {
  if (!state.editor) {
    return;
  }
  const line = state.editor.getCursorLine();
  registry.addBookmark(state.documentID, {
    id: crypto.randomUUID(),
    line,
    label: `${docTitleInput.value || state.documentID} line ${line}`,
    createdAt: new Date().toISOString(),
  });
  renderRegistry();
  showToast(`Bookmarked line ${line}`);
}

function generateDoc() {
  const nextDoc = `doc-${crypto.randomUUID().slice(0, 8)}`;
  registry.registerSeedContent(nextDoc, generateDemoText());
  registry.updateTitle(nextDoc, "Generated demo document");
  bootDocument(nextDoc).catch((error) => setStatus("error", error.message));
}

function sampleDoc() {
  const template = templateCatalog().find((value) => value.id === "demo");
  createFromTemplate(template);
}

function updateDerivedTitle(text) {
  const match = (text || "").match(/^#\s+(.+)$/m);
  if (!match) {
    return;
  }
  const current = registry.get(state.documentID);
  if (!current.title || current.title === `Document ${state.documentID}` || current.title === state.documentID) {
    const next = registry.updateTitle(state.documentID, match[1].trim());
    renderDocumentMeta(next);
  }
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

function registerEvents() {
  document.getElementById("open-doc").addEventListener("click", () => {
    bootDocument(docIDInput.value.trim() || "demo").catch((error) => setStatus("error", error.message));
  });
  document.getElementById("new-doc").addEventListener("click", createNewDocument);
  document.getElementById("duplicate-doc").addEventListener("click", duplicateCurrentDocument);
  document.getElementById("paste-link").addEventListener("click", openFromPromptedLink);
  document.getElementById("copy-link").addEventListener("click", copyShareLink);
  document.getElementById("email-link").addEventListener("click", emailShareLink);
  document.getElementById("search-button").addEventListener("click", openSearch);
  document.getElementById("preview-button").addEventListener("click", togglePreview);
  document.getElementById("split-button").addEventListener("click", toggleSplit);
  document.getElementById("import-button").addEventListener("click", importFile);
  document.getElementById("export-button").addEventListener("click", openExport);
  document.getElementById("snapshot-button").addEventListener("click", publishSnapshot);
  document.getElementById("bookmark-button").addEventListener("click", addBookmark);
  document.getElementById("comment-button").addEventListener("click", openCommentPanel);
  document.getElementById("save-version-button").addEventListener("click", saveNamedVersion);
  document.getElementById("summary-button").addEventListener("click", openSummary);
  document.getElementById("focus-button").addEventListener("click", toggleFocusMode);
  document.getElementById("debug-button").addEventListener("click", openDebugPanel);
  document.getElementById("bold-button").addEventListener("click", () => applyFormat("bold"));
  document.getElementById("italic-button").addEventListener("click", () => applyFormat("italic"));
  document.getElementById("underline-button").addEventListener("click", () => applyFormat("underline"));
  [
    "search-button",
    "bold-button",
    "italic-button",
    "underline-button",
    "preview-button",
    "split-button",
    "import-button",
    "export-button",
    "snapshot-button",
    "bookmark-button",
    "comment-button",
    "save-version-button",
    "summary-button",
    "focus-button",
    "debug-button",
  ].forEach((id) => preserveEditorSelectionOnMouseDown(document.getElementById(id)));
  document.getElementById("settings-button").addEventListener("click", openSettings);
  document.getElementById("help-button").addEventListener("click", openHelp);
  document.getElementById("settings-close").addEventListener("click", () => closeOverlay(settingsPanel));
  document.getElementById("help-close").addEventListener("click", () => closeOverlay(helpPanel));
  document.getElementById("search-close").addEventListener("click", () => closeOverlay(searchPanel));
  document.getElementById("export-close").addEventListener("click", () => closeOverlay(exportPanel));
  document.getElementById("comment-close").addEventListener("click", () => closeOverlay(commentPanel));
  document.getElementById("summary-close").addEventListener("click", () => closeOverlay(summaryPanel));
  document.getElementById("debug-close").addEventListener("click", () => closeOverlay(debugPanel));
  document.getElementById("welcome-open-settings").addEventListener("click", openSettings);
  document.getElementById("welcome-dismiss").addEventListener("click", () => {
    window.localStorage.setItem("grid-editor-dismissed-welcome", "true");
    welcomeBannerEl.classList.add("hidden");
  });
  document.getElementById("find-next").addEventListener("click", () => runSearchReplace(false));
  document.getElementById("replace-all").addEventListener("click", () => runSearchReplace(true));
  document.getElementById("goto-line").addEventListener("click", gotoLine);
  document.getElementById("export-markdown").addEventListener("click", () => exportDocument("markdown"));
  document.getElementById("export-html").addEventListener("click", () => exportDocument("html"));
  document.getElementById("export-text").addEventListener("click", () => exportDocument("text"));
  document.getElementById("export-automerge").addEventListener("click", () => exportDocument("automerge"));
  document.getElementById("copy-markdown").addEventListener("click", () => copyFormatted("markdown"));
  document.getElementById("copy-html").addEventListener("click", () => copyFormatted("html"));
  document.getElementById("publish-snapshot").addEventListener("click", publishSnapshot);
  document.getElementById("publish-exchange").addEventListener("click", publishExchange);
  document.getElementById("import-exchange").addEventListener("click", importExchange);
  document.getElementById("export-audit").addEventListener("click", exportAuditReport);
  document.getElementById("generate-demo-doc").addEventListener("click", generateDoc);
  document.getElementById("sample-doc").addEventListener("click", sampleDoc);
  document.getElementById("comment-save").addEventListener("click", saveComment);
  document.getElementById("comment-resolve-all").addEventListener("click", resolveSelectedComments);
  document.getElementById("speak-summary").addEventListener("click", readSummary);
  document.getElementById("dictate-summary").addEventListener("click", dictateComment);
  fileImportEl.addEventListener("change", handleImportedFile);

  displayNameInput.addEventListener("change", () => {
    updatePreferences({ displayName: displayNameInput.value || "Browser User" }, { skipFormSync: true });
  });
  colorInput.addEventListener("change", () => {
    updatePreferences({ color: colorInput.value || "#1d6fd6" }, { skipFormSync: true });
  });
  docTitleInput.addEventListener("change", () => {
    renderDocumentMeta(registry.updateTitle(state.documentID, docTitleInput.value));
    renderRegistry();
    renderActivity();
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
      openSearch();
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
      if (state.focusMode) {
        state.focusMode = false;
        document.body.classList.remove("focus-mode");
        focusButtonEl.textContent = "Focus";
      }
      closeOverlay(settingsPanel);
      closeOverlay(helpPanel);
      closeOverlay(searchPanel);
      closeOverlay(exportPanel);
      closeOverlay(commentPanel);
      closeOverlay(summaryPanel);
      closeOverlay(debugPanel);
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

function formatTime(value) {
  if (!value) {
    return "-";
  }
  try {
    return new Date(value).toLocaleString();
  } catch {
    return value;
  }
}

function triggerDownload(blob, name) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = name;
  link.click();
  URL.revokeObjectURL(url);
}

function safeFilename(value) {
  return String(value || "document").toLowerCase().replace(/[^a-z0-9._-]+/g, "-");
}

function wrapHTML(title, body) {
  return `<!doctype html><html><head><meta charset="utf-8"><title>${escapeHTML(title)}</title></head><body>${body}</body></html>`;
}

async function copyText(value, successLabel) {
  try {
    await navigator.clipboard.writeText(value);
    showToast(successLabel);
  } catch {
    showToast("Clipboard write failed");
  }
}

function textToBase64(value) {
  return window.btoa(unescape(encodeURIComponent(value)));
}

function base64ToText(value) {
  return decodeURIComponent(escape(window.atob(value)));
}

function bytesToBase64(bytes) {
  let text = "";
  bytes.forEach((value) => {
    text += String.fromCharCode(value);
  });
  return window.btoa(text);
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

async function readFileAsDataURL(file) {
  return await new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result);
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });
}

registerEvents();
applyPreferences(state.prefs);
renderRegistry();
if (window.localStorage.getItem("grid-editor-dismissed-welcome") === "true") {
  welcomeBannerEl.classList.add("hidden");
}

loadMeta()
  .then(() => bootDocument(state.documentID))
  .catch((error) => {
    setStatus("error", error.message);
  });
