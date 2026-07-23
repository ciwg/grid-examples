const statsEl = document.getElementById("stats");
const placeListEl = document.getElementById("place-list");
const resourceListEl = document.getElementById("resource-list");
const responsibilityListEl = document.getElementById("responsibility-list");
const itemListEl = document.getElementById("item-list");
const runListEl = document.getElementById("run-list");
const draftReviewEl = document.getElementById("draft-review");
const problemReviewEl = document.getElementById("problem-review");
const searchResultsEl = document.getElementById("search-results");
const searchActiveEl = document.getElementById("search-active");
const searchRawEl = document.getElementById("search-raw");
const searchDebugEl = document.getElementById("search-debug");
const searchClearEl = document.getElementById("search-clear");
const searchAdvancedEl = document.getElementById("search-advanced");
const searchPlaceSelectEl = document.getElementById("search-place-select");
const searchResourceSelectEl = document.getElementById("search-resource-select");
const searchResponsibilitySelectEl = document.getElementById("search-responsibility-select");
const reviewLaneDraftsEl = document.getElementById("review-lane-drafts");
const reviewLaneHotspotsEl = document.getElementById("review-lane-hotspots");
const reviewLaneSearchEl = document.getElementById("review-lane-search");
const searchPresetDraftsEl = document.getElementById("search-preset-drafts");
const searchPresetProblemsEl = document.getElementById("search-preset-problems");
const searchPresetCountsEl = document.getElementById("search-preset-counts");
const searchPresetRunsEl = document.getElementById("search-preset-runs");
const focusProblemsEl = document.getElementById("focus-problems");
const focusSearchEl = document.getElementById("focus-search");
const focusDraftsEl = document.getElementById("focus-drafts");
const modeReviewEl = document.getElementById("mode-review");
const modeAuthorEl = document.getElementById("mode-author");
const modeOperateEl = document.getElementById("mode-operate");
const modeCreateEl = document.getElementById("mode-create");
const modeBrowseEl = document.getElementById("mode-browse");
const toastEl = document.getElementById("toast");
const workspaceStatusEl = document.getElementById("workspace-status");
const detailMetaEl = document.getElementById("detail-meta");
const detailSummaryEl = document.getElementById("detail-summary");
const detailPrimaryEl = document.getElementById("detail-primary");
const detailActionsEl = document.getElementById("detail-actions");
const detailReviewEl = document.getElementById("detail-review");
const detailTimelineEl = document.getElementById("detail-timeline");
const detailJSONEl = document.getElementById("detail-json");

const editorItemIDEl = document.getElementById("editor-item-id");
const editorActorEl = document.getElementById("editor-actor");
const editorDisplayNameEl = document.getElementById("editor-display-name");
const editorColorEl = document.getElementById("editor-color");
const editorCollabDetailsEl = document.getElementById("editor-collab-details");
const editorLifecycleDetailsEl = document.getElementById("editor-lifecycle-details");
const editorMetaEl = document.getElementById("editor-meta");
const editorParticipantsEl = document.getElementById("editor-participants");
const editorStatusCardsEl = document.getElementById("editor-status-cards");
const editorBodyEl = document.getElementById("editor-body");
const editorFocusWritingEl = document.getElementById("editor-focus-writing");
const editorRefreshEl = document.getElementById("editor-refresh");
const editorSnapshotEl = document.getElementById("editor-snapshot");
const editorApproveEl = document.getElementById("editor-approve");
const editorSupersedeEl = document.getElementById("editor-supersede");
const editorSupportDetailsEl = document.getElementById("editor-support-details");
const approvalFormEl = document.getElementById("approval-form");
const createDetailsEl = document.getElementById("create-details");
const browseDetailsEl = document.getElementById("browse-details");
const resourcePlaceSelectEl = document.getElementById("resource-place-select");
const itemResponsibilitySelectEl = document.getElementById("item-responsibility-select");
const runItemSelectEl = document.getElementById("run-item-select");
const runItemIDEl = document.getElementById("run-item-id");
const runPlaceSelectEl = document.getElementById("run-place-select");
const runResourceSelectEl = document.getElementById("run-resource-select");
const runResponsibilitySelectEl = document.getElementById("run-responsibility-select");
const evidenceRunSelectEl = document.getElementById("evidence-run-select");
const approvalTargetSelectEl = document.getElementById("approval-target-select");
const approvalTargetSummaryEl = document.getElementById("approval-target-summary");
const operateContextEl = document.getElementById("operate-context");
const operateChooseRunEl = document.getElementById("operate-choose-run");
const operateChooseEvidenceEl = document.getElementById("operate-choose-evidence");
const operateChooseApproveEl = document.getElementById("operate-choose-approve");
const operateRunCurrentEl = document.getElementById("operate-run-current");
const operateEvidenceCurrentEl = document.getElementById("operate-evidence-current");
const operateApproveCurrentEl = document.getElementById("operate-approve-current");
const operateRunDetailsEl = document.getElementById("operate-run-details");
const operateEvidenceDetailsEl = document.getElementById("operate-evidence-details");
const operateApprovalDetailsEl = document.getElementById("operate-approval-details");

const participantID = getParticipantID();
const bridgePendingRPC = new Map();
const bridgePendingHandshakes = new Map();
const bridgeRPCDeadlineMS = 1000;
let bridgeSequence = 0;
const editorState = {
  itemID: "",
  version: 0,
  title: "",
  status: "",
  currentRevision: 0,
  participantCount: 0,
  dirty: false,
  pushing: false,
  lastRenderedBody: "",
  pushTimer: 0,
  pollTimer: 0,
  heartbeatTimer: 0,
  reconnectTimer: 0,
  socket: null,
  socketConnected: false,
  transport: "http-poll",
  socketGeneration: 0,
  bridgeRequestID: "",
};

const browserBridgeState = {
  ready: false,
  supported: false,
  meta: null,
  socketPath: "",
};

const detailState = {
  type: "",
  id: "",
  record: null,
};

const catalogState = {
  places: [],
  resources: [],
  responsibilities: [],
  items: [],
  runs: [],
};

const workspaceEls = Array.from(document.querySelectorAll(".workspace[data-mode]"));
const modeButtons = {
  review: modeReviewEl,
  author: modeAuthorEl,
  operate: modeOperateEl,
  create: modeCreateEl,
  browse: modeBrowseEl,
};

const reviewLaneEls = {
  drafts: document.getElementById("review-drafts-lane"),
  hotspots: document.getElementById("review-hotspots-lane"),
  search: document.getElementById("review-search-lane"),
};

const reviewLaneButtons = {
  drafts: reviewLaneDraftsEl,
  hotspots: reviewLaneHotspotsEl,
  search: reviewLaneSearchEl,
};

const operateStageDetails = {
  run: operateRunDetailsEl,
  evidence: operateEvidenceDetailsEl,
  approval: operateApprovalDetailsEl,
};

const operateStageButtons = {
  run: operateChooseRunEl,
  evidence: operateChooseEvidenceEl,
  approval: operateChooseApproveEl,
};

window.addEventListener("message", (event) => {
  if (event.source !== window || !event.data || event.data.__oks_bridge !== true || event.data.direction !== "bridge->page") {
    return;
  }
  const message = event.data;
  if (message.kind === "handshake") {
    const resolve = bridgePendingHandshakes.get(message.request_id);
    if (resolve) {
      bridgePendingHandshakes.delete(message.request_id);
      resolve(!!message.ok);
    }
    return;
  }
  if (message.kind === "rpc-response") {
    const entry = bridgePendingRPC.get(message.request_id);
    if (!entry) {
      return;
    }
    bridgePendingRPC.delete(message.request_id);
    entry.resolve(message.response);
    return;
  }
  if (message.kind === "error") {
    const entry = bridgePendingRPC.get(message.request_id);
    if (entry) {
      bridgePendingRPC.delete(message.request_id);
      entry.reject(new Error(message.error || "Browser bridge error"));
      return;
    }
    if (message.request_id === editorState.bridgeRequestID) {
      editorState.socketConnected = false;
      if (editorState.itemID) {
        scheduleLiveReconnect(editorState.itemID, editorState.socketGeneration);
      }
      showToast(message.error || "Browser bridge error");
      setWorkspaceStatus("Direct browser embodiment unavailable", "error", message.error || "Browser bridge error");
    }
    return;
  }
  if (message.kind === "live-message" && message.request_id === editorState.bridgeRequestID && message.response) {
    handleLiveBridgeResponse(message.response);
  }
});

function nextBridgeRequestID(prefix = "bridge") {
  bridgeSequence += 1;
  return `${prefix}-${bridgeSequence}`;
}

function postBridgeMessage(message) {
  window.postMessage({
    __oks_bridge: true,
    direction: "page->bridge",
    ...message,
  }, window.location.origin);
}

function runHandled(action, context) {
  return (...args) => {
    clearWorkspaceStatus();
    Promise.resolve(action(...args)).catch((error) => handleError(error, context));
  };
}

function displayRecordID(record) {
  return record?.alias_id || record?.id || "";
}

// Intent: Keep the single-page browser readable as distinct workflow modes by
// letting operators explicitly activate Review, Author, Operate, Create, or
// Browse without hiding any shipped surface, while letting inactive modes
// recede more aggressively than simple tinting alone. Source: DI-bavum;
// DI-nabek
function setActiveMode(mode) {
  for (const workspace of workspaceEls) {
    const active = workspace.dataset.mode === mode;
    workspace.classList.toggle("is-active", active);
    workspace.classList.toggle("is-muted", !active);
  }
  for (const [key, button] of Object.entries(modeButtons)) {
    button.classList.toggle("is-active", key === mode);
  }
}

function focusMode(mode) {
  const workspace = document.getElementById(`workspace-${mode}`);
  if (!workspace) {
    return;
  }
  createDetailsEl.open = mode === "create";
  browseDetailsEl.open = mode === "browse";
  if (mode === "create") {
    createDetailsEl.open = true;
  }
  if (mode === "browse") {
    browseDetailsEl.open = true;
  }
  setActiveMode(mode);
  workspace.scrollIntoView({ behavior: "smooth", block: "start" });
}

// Intent: Keep Review feeling like one working state instead of three peer
// panels by showing exactly one review queue lane at a time while preserving
// drafts, hotspots, and search as reachable entry paths. Source: DI-rabok;
// DI-javik
function setReviewLane(lane) {
  for (const [key, section] of Object.entries(reviewLaneEls)) {
    if (!section) {
      continue;
    }
    section.hidden = key !== lane;
  }
  for (const [key, button] of Object.entries(reviewLaneButtons)) {
    if (!button) {
      continue;
    }
    const active = key === lane;
    button.classList.toggle("is-active", active);
    button.classList.toggle("button-secondary", !active);
  }
}

// Intent: Make Operate feel like a staged action workspace instead of one
// large transaction console by opening only the selected operation form while
// keeping every generic form reachable. Source: DI-zumor
function openOperateStage(stage) {
  for (const [key, details] of Object.entries(operateStageDetails)) {
    if (details) {
      details.open = key === stage;
    }
  }
  for (const [key, button] of Object.entries(operateStageButtons)) {
    if (!button) {
      continue;
    }
    const active = key === stage;
    button.classList.toggle("is-active", active);
    button.classList.toggle("button-secondary", !active);
  }
}

// Intent: Keep the browser as an equal operational embodiment while making
// knowledge-item drafting collaborative in the browser without collapsing the
// durable revision and approval workflow into ephemeral UI state. Source:
// DI-lusov; DI-zoruk
document.getElementById("place-form").addEventListener("submit", runHandled(async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  await postJSON("/api/places", {
    actor: form.actor.value,
    kind: form.kind.value,
    name: form.name.value,
    summary: form.summary.value,
    parent_id: form.parent_id.value,
    tags: splitCSV(form.tags.value),
  });
  form.reset();
  form.actor.value = "alice";
  await refresh();
}, "Create Place"));

document.getElementById("resource-form").addEventListener("submit", runHandled(async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const placeID = firstPresent(form.place_id.value, form.place_id_select.value);
  await postJSON("/api/resources", {
    actor: form.actor.value,
    kind: form.kind.value,
    name: form.name.value,
    summary: form.summary.value,
    place_id: placeID,
    tags: splitCSV(form.tags.value),
  });
  form.reset();
  form.actor.value = "alice";
  await refresh();
}, "Create Resource"));

document.getElementById("responsibility-form").addEventListener("submit", runHandled(async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  await postJSON("/api/responsibilities", {
    actor: form.actor.value,
    title: form.title.value,
    summary: form.summary.value,
    role_keys: splitCSV(form.role_keys.value),
    tags: splitCSV(form.tags.value),
  });
  form.reset();
  form.actor.value = "alice";
  await refresh();
}, "Create Responsibility"));

document.getElementById("item-form").addEventListener("submit", runHandled(async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const responsibilityIDs = mergeSelectedIDs(
    form.responsibility_ids.value,
    form.responsibility_id_select.value,
  );
  const item = await postJSON("/api/items", {
    actor: form.actor.value,
    kind: form.kind.value,
    title: form.title.value,
    summary: form.summary.value,
    body: form.body.value,
    tags: splitCSV(form.tags.value),
    responsibility_ids: responsibilityIDs,
  });
  form.reset();
  form.actor.value = "alice";
  form.kind.value = "procedure";
  await refresh(item.id);
}, "Create Knowledge Item"));

document.getElementById("run-form").addEventListener("submit", runHandled(async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const itemID = firstPresent(form.item_id.value, form.item_id_select.value);
  const placeID = firstPresent(form.place_id.value, form.place_id_select.value);
  const resourceIDs = mergeSelectedIDs(form.resource_ids.value, form.resource_id_select.value);
  const responsibilityIDs = mergeSelectedIDs(form.responsibility_ids.value, form.responsibility_id_select.value);
  await postJSON("/api/runs", {
    actor: form.actor.value,
    kind: form.kind.value,
    item_id: itemID,
    revision: Number(form.revision.value || 1),
    outcome: form.outcome.value,
    notes: form.notes.value,
    machine: form.machine.value,
    location: form.location.value,
    place_id: placeID,
    resource_ids: resourceIDs,
    responsibility_ids: responsibilityIDs,
  });
  form.reset();
  form.actor.value = "bob";
  form.kind.value = "procedure";
  form.revision.value = "1";
  await refresh();
}, "Record Run"));

document.getElementById("approval-form").addEventListener("submit", runHandled(async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const targetID = firstPresent(form.target_id.value, form.target_id_select.value);
  const targetType = form.target_type.value;
  const path = targetType === "run" ? `/api/runs/${targetID}/approvals` : `/api/items/${targetID}/approvals`;
  await postJSON(path, {
    actor: form.actor.value,
    revision: Number(form.revision.value || 0),
    role: form.role.value,
    decision: form.decision.value,
    notes: form.notes.value,
  });
  form.reset();
  form.actor.value = "boss";
  form.target_type.value = "knowledge_item";
  form.decision.value = "approved";
  form.revision.value = "0";
  renderApprovalTargetOptions();
  await refresh(editorState.itemID);
}, "Record Approval"));

document.getElementById("evidence-form").addEventListener("submit", runHandled(async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const runID = firstPresent(form.run_id.value, form.run_id_select.value);
  let facts = {};
  if (form.facts_json.value.trim()) {
    facts = JSON.parse(form.facts_json.value);
  }
  let attachmentName = "";
  let attachmentBodyBase64 = "";
  if (form.attachment.files[0]) {
    attachmentName = form.attachment.files[0].name;
    attachmentBodyBase64 = bytesToBase64(await form.attachment.files[0].arrayBuffer());
  }
  const response = await bridgeRPC({
    type: "operation",
    operation: "add_evidence",
    actor: form.actor.value,
    run_id: runID,
    summary: form.summary.value,
    facts,
    attachment_name: attachmentName,
    attachment_body_base64: attachmentBodyBase64,
  });
  if (!response || response.status >= 400) {
    throw new Error(response && response.body ? response.body : "Add evidence failed");
  }
  form.reset();
  form.actor.value = "bob";
  showToast("Evidence attached");
  await refresh();
}, "Add Evidence"));

document.getElementById("search-form").addEventListener("submit", runHandled(async (event) => {
  event.preventDefault();
  event.currentTarget.dataset.problem = "false";
  const filters = getSearchFilters(event.currentTarget);
  await runSearch(filters);
}, "Search"));

searchClearEl.addEventListener("click", () => {
  clearSearch();
});

searchPresetDraftsEl.addEventListener("click", runHandled(async () => {
  setReviewLane("search");
  await runPresetSearch({ status: "draft" }, "Draft review is loaded. Open one draft record and keep the next action attached to that record.");
}, "Search"));

searchPresetProblemsEl.addEventListener("click", runHandled(async () => {
  setReviewLane("search");
  await runPresetSearch({ kind: "receiving_check", problem: true }, "Receiving problems are loaded. Start with one hotspot or run and stay in that review thread.");
}, "Search"));

searchPresetCountsEl.addEventListener("click", runHandled(async () => {
  setReviewLane("search");
  await runPresetSearch({ kind: "inventory_audit" }, "Inventory counts are loaded. Open one run or context record and review the count history there.");
}, "Search"));

searchPresetRunsEl.addEventListener("click", runHandled(async () => {
  setReviewLane("search");
  await runPresetSearch({}, "Recent runs are loaded. Search broadly when you know the work but not yet the exact record.");
}, "Search"));

reviewLaneDraftsEl.addEventListener("click", () => {
  setActiveMode("review");
  setReviewLane("drafts");
});

reviewLaneHotspotsEl.addEventListener("click", () => {
  setActiveMode("review");
  setReviewLane("hotspots");
});

reviewLaneSearchEl.addEventListener("click", () => {
  setActiveMode("review");
  setReviewLane("search");
});

modeReviewEl.addEventListener("click", () => {
  focusMode("review");
});

// Intent: Keep Review-mode inspection calm by only loading the live draft when
// the operator explicitly enters Author and there is a current item context to
// carry forward. Source: DI-suvor
modeAuthorEl.addEventListener("click", runHandled(async () => {
  focusMode("author");
  if (!editorState.itemID && detailState.type === "item" && detailState.record) {
    await loadEditorItem(detailState.record.id, { activateMode: false });
  }
}, "Author"));

modeOperateEl.addEventListener("click", () => {
  focusMode("operate");
});

modeCreateEl.addEventListener("click", () => {
  focusMode("create");
});

modeBrowseEl.addEventListener("click", () => {
  focusMode("browse");
});

focusProblemsEl.addEventListener("click", () => {
  setActiveMode("review");
  setReviewLane("hotspots");
  problemReviewEl.scrollIntoView({ behavior: "smooth", block: "start" });
  setWorkspaceStatus("Review hotspots first when you need the fastest path into repeated receiving or count problems.", "info");
});

focusSearchEl.addEventListener("click", () => {
  setActiveMode("review");
  setReviewLane("search");
  const queryInput = document.querySelector("#search-form input[name='q']");
  if (queryInput) {
    queryInput.scrollIntoView({ behavior: "smooth", block: "center" });
    queryInput.focus({ preventScroll: true });
  }
  setWorkspaceStatus("Search is the main path when you already know the item, run, place, resource, or responsibility you need.", "info");
});

focusDraftsEl.addEventListener("click", runHandled(async () => {
  setActiveMode("review");
  setReviewLane("drafts");
  draftReviewEl.scrollIntoView({ behavior: "smooth", block: "start" });
  setWorkspaceStatus("Draft items are loaded. Open one record, then use the inspector to draft, approve, or record work from that item.", "info");
}, "Primary Flow"));

// Intent: Keep collaboration metadata reachable without letting it compete
// with the draft body once the operator explicitly chooses writing focus.
// Source: DI-tavul
editorFocusWritingEl.addEventListener("click", () => {
  setActiveMode("author");
  editorCollabDetailsEl.open = false;
  editorLifecycleDetailsEl.open = false;
  editorSupportDetailsEl.open = false;
  editorBodyEl.scrollIntoView({ behavior: "smooth", block: "center" });
  editorBodyEl.focus({ preventScroll: true });
  setWorkspaceStatus("Writing focus is active. Stay in one draft, then snapshot when the change is coherent.", "info");
});

operateChooseRunEl.addEventListener("click", () => {
  setActiveMode("operate");
  openOperateStage("run");
  operateRunDetailsEl.scrollIntoView({ behavior: "smooth", block: "start" });
});

operateChooseEvidenceEl.addEventListener("click", () => {
  setActiveMode("operate");
  openOperateStage("evidence");
  operateEvidenceDetailsEl.scrollIntoView({ behavior: "smooth", block: "start" });
});

operateChooseApproveEl.addEventListener("click", () => {
  setActiveMode("operate");
  openOperateStage("approval");
  operateApprovalDetailsEl.scrollIntoView({ behavior: "smooth", block: "start" });
});

operateRunCurrentEl.addEventListener("click", runHandled(async () => {
  triggerOperateContext("run");
}, "Operate"));

operateEvidenceCurrentEl.addEventListener("click", runHandled(async () => {
  triggerOperateContext("evidence");
}, "Operate"));

operateApproveCurrentEl.addEventListener("click", runHandled(async () => {
  triggerOperateContext("approve");
}, "Operate"));

editorItemIDEl.addEventListener("change", runHandled(async () => {
  await loadEditorItem(editorItemIDEl.value, { activateMode: true });
}, "Live Draft Studio"));

editorBodyEl.addEventListener("input", () => {
  editorState.dirty = true;
  renderEditorStatusCards({
    version: editorState.version,
    current_revision: editorState.currentRevision,
    participants: new Array(editorState.participantCount),
    status: editorState.status,
  }, editorBodyEl.value);
  scheduleLivePush();
});

editorBodyEl.addEventListener("click", scheduleLivePush);
editorBodyEl.addEventListener("keyup", scheduleLivePush);

editorRefreshEl.addEventListener("click", runHandled(async () => {
  if (!editorState.itemID) {
    return;
  }
  await pullLiveState(editorState.itemID, true);
}, "Live Draft Studio"));

editorSnapshotEl.addEventListener("click", runHandled(async () => {
  if (!editorState.itemID) {
    showToast("Select an item before snapshotting");
    return;
  }
  await flushLivePush();
  const item = await getJSON(`/api/items/${editorState.itemID}`);
  const updated = await postJSON(`/api/items/${editorState.itemID}/revisions`, {
    actor: editorActorEl.value,
    title: item.title,
    summary: item.summary,
    body: editorBodyEl.value,
    tags: item.tags || [],
  });
  showToast(`Snapshot created as revision ${updated.current_revision}`);
  await refresh(editorState.itemID);
}, "Live Draft Studio"));

editorApproveEl.addEventListener("click", runHandled(async () => {
  if (!editorState.itemID) {
    showToast("Select an item before approving");
    return;
  }
  await postJSON(`/api/items/${editorState.itemID}/approvals`, {
    actor: editorActorEl.value,
    revision: editorState.currentRevision,
    role: "reviewer",
    decision: "approved",
    notes: "Approved from live draft studio",
  });
  showToast("Current revision approved");
  await refresh(editorState.itemID);
}, "Live Draft Studio"));

editorSupersedeEl.addEventListener("click", runHandled(async () => {
  if (!editorState.itemID) {
    showToast("Select an item before superseding");
    return;
  }
  await postJSON(`/api/items/${editorState.itemID}/supersede`, {
    actor: editorActorEl.value,
    notes: "Superseded from live draft studio",
  });
  showToast("Item superseded");
  await refresh(editorState.itemID);
}, "Live Draft Studio"));

approvalFormEl.target_type.addEventListener("change", () => {
  renderApprovalTargetOptions();
  syncApprovalDefaults();
});

runItemSelectEl.addEventListener("change", () => {
  if (runItemSelectEl.value) {
    runItemIDEl.value = runItemSelectEl.value;
    const item = catalogState.items.find((value) => value.id === runItemSelectEl.value);
    if (item && item.current_revision) {
      document.getElementById("run-form").revision.value = String(item.current_revision);
    }
  }
});

evidenceRunSelectEl.addEventListener("change", () => {
  if (evidenceRunSelectEl.value) {
    document.getElementById("evidence-form").run_id.value = evidenceRunSelectEl.value;
  }
});

approvalTargetSelectEl.addEventListener("change", () => {
  const form = document.getElementById("approval-form");
  if (approvalTargetSelectEl.value) {
    form.target_id.value = approvalTargetSelectEl.value;
    if (form.target_type.value === "knowledge_item") {
      const item = catalogState.items.find((value) => value.id === approvalTargetSelectEl.value);
      if (item && item.current_revision) {
        form.revision.value = String(item.current_revision);
      }
    } else {
      form.revision.value = "0";
    }
  }
  syncApprovalDefaults();
});

function splitCSV(input) {
  return input.split(",").map((value) => value.trim()).filter(Boolean);
}

function firstPresent(...values) {
  for (const value of values) {
    if (value && value.trim()) {
      return value.trim();
    }
  }
  return "";
}

function mergeSelectedIDs(csvValue, selectedValue) {
  const values = splitCSV(csvValue || "");
  if (selectedValue && !values.includes(selectedValue)) {
    values.unshift(selectedValue);
  }
  return values;
}

function setWorkspaceStatus(message, tone = "error", detail = "") {
  workspaceStatusEl.hidden = false;
  workspaceStatusEl.dataset.tone = tone;
  workspaceStatusEl.innerHTML = "";
  const headline = document.createElement("strong");
  headline.textContent = message;
  workspaceStatusEl.appendChild(headline);
  if (detail) {
    const block = document.createElement("pre");
    block.textContent = detail;
    workspaceStatusEl.appendChild(block);
  }
}

function clearWorkspaceStatus() {
  workspaceStatusEl.hidden = true;
  workspaceStatusEl.dataset.tone = "";
  workspaceStatusEl.innerHTML = "";
}

function browserEmbodimentUnavailableMessage() {
  return "This embodiment currently requires Chrome or Chromium with the ex5 browser extension installed.";
}

function isChromeOrChromiumBrowser() {
  const agent = navigator.userAgent || "";
  if (agent.includes("Edg/")) {
    return false;
  }
  return agent.includes("Chrome/") || agent.includes("Chromium/");
}

// Intent: Make the browser embodiment state its Chrome/Chromium direct-contract
// requirement honestly up front and prove the extension/native-host/runtime
// chain before startup marks the browser ready, instead of silently demoting
// back into the older HTTP browser path or overstating readiness after only a
// page-bridge check. Source: DI-punek; DI-salov
async function initializeBrowserEmbodiment() {
  browserBridgeState.supported = isChromeOrChromiumBrowser();
  if (!browserBridgeState.supported) {
    setWorkspaceStatus("Chrome or Chromium is required", "error", browserEmbodimentUnavailableMessage());
    return false;
  }
  const metaResponse = await fetch("/api/meta");
  if (!metaResponse.ok) {
    throw new Error(await metaResponse.text());
  }
  browserBridgeState.meta = await metaResponse.json();
  browserBridgeState.socketPath = browserBridgeState.meta?.local_unix_socket_path || "";
  const browserEmbodiment = browserBridgeState.meta?.embodiments?.browser || {};
  if (browserEmbodiment.primary_adapter !== "chrome_native_messaging") {
    throw new Error("runtime does not advertise the direct Chrome/Chromium browser embodiment");
  }
  const handshakeOK = await bridgeHandshake();
  if (!handshakeOK) {
    setWorkspaceStatus("Direct browser embodiment unavailable", "error", browserEmbodimentUnavailableMessage());
    return false;
  }
  browserBridgeState.ready = true;
  return true;
}

function bridgeHandshake() {
  const requestID = nextBridgeRequestID("bridge-handshake");
  return new Promise((resolve) => {
    const timer = setTimeout(() => {
      bridgePendingHandshakes.delete(requestID);
      resolve(false);
    }, 300);
    bridgePendingHandshakes.set(requestID, (ok) => {
      clearTimeout(timer);
      resolve(ok);
    });
    postBridgeMessage({
      kind: "handshake",
      request_id: requestID,
      socket_path: browserBridgeState.socketPath,
    });
  });
}

function bridgeRPC(request) {
  if (!browserBridgeState.ready) {
    throw new Error(browserEmbodimentUnavailableMessage());
  }
  const requestID = nextBridgeRequestID("bridge-rpc");
  return new Promise((resolve, reject) => {
    // Intent: Bound browser one-shot direct-contract waits at the page layer so
    // lost bridge replies fail closed and clean up pending promise state
    // instead of hanging the UI indefinitely. Source: DI-zabem
    const timer = setTimeout(() => {
      const entry = bridgePendingRPC.get(requestID);
      if (!entry) {
        return;
      }
      bridgePendingRPC.delete(requestID);
      entry.reject(new Error(`Direct browser RPC timed out after ${bridgeRPCDeadlineMS}ms.`));
    }, bridgeRPCDeadlineMS);
    bridgePendingRPC.set(requestID, {
      resolve(value) {
        clearTimeout(timer);
        resolve(value);
      },
      reject(error) {
        clearTimeout(timer);
        reject(error);
      },
    });
    postBridgeMessage({
      kind: "rpc",
      request_id: requestID,
      socket_path: browserBridgeState.socketPath,
      request,
    });
  });
}

function bytesToBase64(buffer) {
  const bytes = new Uint8Array(buffer);
  let binary = "";
  const chunkSize = 0x8000;
  for (let index = 0; index < bytes.length; index += chunkSize) {
    const chunk = bytes.subarray(index, index + chunkSize);
    for (const value of chunk) {
      binary += String.fromCharCode(value);
    }
  }
  return btoa(binary);
}

function directOperationForGET(path) {
  const url = new URL(path, window.location.href);
  const { pathname } = url;
  if (pathname === "/api/dashboard") {
    return { type: "operation", operation: "dashboard" };
  }
  if (pathname === "/api/places") {
    return { type: "operation", operation: "list_places" };
  }
  if (pathname === "/api/resources") {
    return { type: "operation", operation: "list_resources" };
  }
  if (pathname === "/api/responsibilities") {
    return { type: "operation", operation: "list_responsibilities" };
  }
  if (pathname === "/api/items") {
    return { type: "operation", operation: "list_items" };
  }
  if (pathname === "/api/runs") {
    return { type: "operation", operation: "list_runs" };
  }
  const liveStateMatch = pathname.match(/^\/api\/items\/([^/]+)\/live$/);
  if (liveStateMatch) {
    return { type: "operation", operation: "load_live_state", item_id: liveStateMatch[1] };
  }
  const itemMatch = pathname.match(/^\/api\/items\/([^/]+)$/);
  if (itemMatch) {
    return { type: "operation", operation: "inspect_item", item_id: itemMatch[1] };
  }
  const runMatch = pathname.match(/^\/api\/runs\/([^/]+)$/);
  if (runMatch) {
    return { type: "operation", operation: "inspect_run", run_id: runMatch[1] };
  }
  const placeMatch = pathname.match(/^\/api\/places\/([^/]+)$/);
  if (placeMatch) {
    return { type: "operation", operation: "inspect_entity", entity_type: "place", entity_id: placeMatch[1] };
  }
  const resourceMatch = pathname.match(/^\/api\/resources\/([^/]+)$/);
  if (resourceMatch) {
    return { type: "operation", operation: "inspect_entity", entity_type: "resource", entity_id: resourceMatch[1] };
  }
  const responsibilityMatch = pathname.match(/^\/api\/responsibilities\/([^/]+)$/);
  if (responsibilityMatch) {
    return { type: "operation", operation: "inspect_entity", entity_type: "responsibility", entity_id: responsibilityMatch[1] };
  }
  if (pathname === "/api/search") {
    return {
      type: "operation",
      operation: "search",
      search_options: {
        query: url.searchParams.get("q") || "",
        kind: url.searchParams.get("kind") || "",
        status: url.searchParams.get("status") || "",
        outcome: url.searchParams.get("outcome") || "",
        problem: url.searchParams.get("problem") === "true",
        place_id: url.searchParams.get("place_id") || "",
        resource_id: url.searchParams.get("resource_id") || "",
        responsibility_id: url.searchParams.get("responsibility_id") || "",
      },
    };
  }
  if (pathname === "/api/problem-review") {
    return { type: "operation", operation: "problem_review" };
  }
  return null;
}

// Intent: Raise the browser create/operate mutation slice above route-shaped
// forwarding by naming each durable workflow directly instead of tunneling
// those writes back through generic method-and-path requests. Source: DI-rumav
function directOperationForPOST(path, payload) {
  if (path === "/api/places") {
    return {
      type: "operation",
      operation: "create_place",
      actor: payload.actor,
      kind: payload.kind,
      name: payload.name,
      summary: payload.summary,
      parent_id: payload.parent_id,
      tags: payload.tags || [],
    };
  }
  if (path === "/api/resources") {
    return {
      type: "operation",
      operation: "create_resource",
      actor: payload.actor,
      kind: payload.kind,
      name: payload.name,
      summary: payload.summary,
      place_id: payload.place_id,
      tags: payload.tags || [],
    };
  }
  if (path === "/api/responsibilities") {
    return {
      type: "operation",
      operation: "create_responsibility",
      actor: payload.actor,
      title: payload.title,
      summary: payload.summary,
      role_keys: payload.role_keys || [],
      tags: payload.tags || [],
    };
  }
  if (path === "/api/items") {
    return {
      type: "operation",
      operation: "create_item",
      actor: payload.actor,
      kind: payload.kind,
      title: payload.title,
      summary: payload.summary,
      body: payload.body,
      tags: payload.tags || [],
      responsibility_ids: payload.responsibility_ids || [],
    };
  }
  if (path === "/api/runs") {
    return {
      type: "operation",
      operation: "record_run",
      actor: payload.actor,
      kind: payload.kind,
      item_id: payload.item_id,
      revision: payload.revision,
      outcome: payload.outcome,
      notes: payload.notes,
      machine: payload.machine,
      location: payload.location,
      place_id: payload.place_id,
      resource_ids: payload.resource_ids || [],
      responsibility_ids: payload.responsibility_ids || [],
    };
  }
  const itemRevisionMatch = path.match(/^\/api\/items\/([^/]+)\/revisions$/);
  if (itemRevisionMatch) {
    return {
      type: "operation",
      operation: "add_revision",
      actor: payload.actor,
      item_id: itemRevisionMatch[1],
      title: payload.title,
      summary: payload.summary,
      body: payload.body,
      tags: payload.tags || [],
    };
  }
  const itemApprovalMatch = path.match(/^\/api\/items\/([^/]+)\/approvals$/);
  if (itemApprovalMatch) {
    return {
      type: "operation",
      operation: "record_item_approval",
      actor: payload.actor,
      item_id: itemApprovalMatch[1],
      revision: payload.revision,
      role: payload.role,
      decision: payload.decision,
      notes: payload.notes,
    };
  }
  const runApprovalMatch = path.match(/^\/api\/runs\/([^/]+)\/approvals$/);
  if (runApprovalMatch) {
    return {
      type: "operation",
      operation: "record_run_approval",
      actor: payload.actor,
      run_id: runApprovalMatch[1],
      role: payload.role,
      decision: payload.decision,
      notes: payload.notes,
    };
  }
  const itemSupersedeMatch = path.match(/^\/api\/items\/([^/]+)\/supersede$/);
  if (itemSupersedeMatch) {
    return {
      type: "operation",
      operation: "supersede_item",
      actor: payload.actor,
      item_id: itemSupersedeMatch[1],
      notes: payload.notes,
    };
  }
  return null;
}

async function requestBrowserJSON(method, path, payload) {
  const directOperation = method === "GET" ? directOperationForGET(path) : null;
  const directWriteOperation = method === "POST" ? directOperationForPOST(path, payload || {}) : null;
  const response = directOperation ? await bridgeRPC(directOperation) : directWriteOperation ? await bridgeRPC(directWriteOperation) : await bridgeRPC({
    type: "request",
    method,
    path,
    headers: payload ? { "Content-Type": "application/json" } : {},
    body: payload ? JSON.stringify(payload) : "",
  });
  if (!response) {
    return null;
  }
  if (response.status >= 400) {
    throw new Error(response.body || `browser direct request failed: ${response.status}`);
  }
  return response.body ? JSON.parse(response.body) : null;
}

async function postJSON(path, payload) {
  const result = await requestBrowserJSON("POST", path, payload);
  showToast(`Saved via ${path}`);
  return result;
}

async function getJSON(path) {
  if (path === "/api/meta" && !browserBridgeState.ready) {
    const response = await fetch(path);
    if (!response.ok) {
      throw new Error(await response.text());
    }
    return response.json();
  }
  return requestBrowserJSON("GET", path, null);
}

function renderStats(data) {
  statsEl.innerHTML = "";
  const fields = [
    ["Places", data.places],
    ["Resources", data.resources],
    ["Responsibilities", data.responsibilities],
    ["Procedures", data.procedures],
    ["Training", data.training_items],
    ["Maintenance", data.maintenance_items],
    ["Receiving", data.receiving_items],
    ["Inventory", data.inventory_items],
    ["Runs", data.procedure_runs + data.training_runs + data.maintenance_runs + data.receiving_runs + data.inventory_runs],
    ["Approvals", data.approvals],
    ["Evidence", data.evidence],
    ["Links", data.links],
  ];
  for (const [label, value] of fields) {
    const card = document.createElement("div");
    card.className = "stat";
    card.innerHTML = `<strong>${value}</strong><span>${label}</span>`;
    statsEl.appendChild(card);
  }
}

function renderPlaces(items) {
  catalogState.places = items;
  placeListEl.innerHTML = "";
  for (const item of items) {
    const card = document.createElement("button");
    card.type = "button";
    card.className = "card selectable-card";
    card.innerHTML = `<div class="kind">${item.kind}</div><h3>${item.id} · ${item.name}</h3><div class="meta">${item.summary || ""}\nparent: ${item.parent_id || "-"}\nchildren: ${(item.child_place_ids || []).length}\nresources: ${(item.resource_ids || []).length}</div>`;
    card.addEventListener("click", () => {
      inspectRecord("place", item.id).catch(handleError);
    });
    placeListEl.appendChild(card);
  }
}

function renderResources(items) {
  catalogState.resources = items;
  resourceListEl.innerHTML = "";
  for (const item of items) {
    const card = document.createElement("button");
    card.type = "button";
    card.className = "card selectable-card";
    card.innerHTML = `<div class="kind">${item.kind}</div><h3>${item.id} · ${item.name}</h3><div class="meta">${item.summary || ""}\nplace: ${item.place_id || "-"}</div>`;
    card.addEventListener("click", () => {
      inspectRecord("resource", item.id).catch(handleError);
    });
    resourceListEl.appendChild(card);
  }
}

function renderResponsibilities(items) {
  catalogState.responsibilities = items;
  responsibilityListEl.innerHTML = "";
  for (const item of items) {
    const card = document.createElement("button");
    card.type = "button";
    card.className = "card selectable-card";
    card.innerHTML = `<h3>${item.id} · ${item.title}</h3><div class="meta">${item.summary || ""}\nroles: ${(item.linked_role_keys || []).join(", ") || "-"}</div>`;
    card.addEventListener("click", () => {
      inspectRecord("responsibility", item.id).catch(handleError);
    });
    responsibilityListEl.appendChild(card);
  }
}

function renderKnowledgeItems(items) {
  catalogState.items = items;
  itemListEl.innerHTML = "";
  editorItemIDEl.innerHTML = "";
  const placeholder = document.createElement("option");
  placeholder.value = "";
  placeholder.textContent = "Select a knowledge item";
  editorItemIDEl.appendChild(placeholder);
  for (const item of items) {
    const option = document.createElement("option");
    option.value = item.id;
    option.textContent = `${item.id} · ${item.title}`;
    editorItemIDEl.appendChild(option);

    const card = document.createElement("button");
    card.type = "button";
    card.className = "card selectable-card";
    card.innerHTML = `<div class="kind">${item.kind} · ${item.status}</div><h3>${item.id} · ${item.title}</h3><div class="meta">revision ${item.current_revision} · live v${item.working_version}\n${item.summary || ""}</div>`;
    card.addEventListener("click", () => {
      inspectRecord("item", item.id).catch(handleError);
    });
    itemListEl.appendChild(card);
  }
  if (editorState.itemID) {
    editorItemIDEl.value = editorState.itemID;
  }
}

// Intent: Make draft review the clearest browser home path by surfacing only
// draft items that most directly need revision, run recording, or approval
// before operators branch into broader hotspot or search work. Source:
// DI-rabok; DI-javik
function renderDraftQueue(items) {
  draftReviewEl.innerHTML = "";
  const drafts = items.filter((item) => item.status === "draft");
  if (!drafts.length) {
    const empty = document.createElement("div");
    empty.className = "meta";
    empty.textContent = "No draft items are waiting right now. Switch to hotspots or search when work starts from an observed issue or known record.";
    draftReviewEl.appendChild(empty);
    return;
  }
  for (const item of drafts) {
    const card = document.createElement("article");
    card.className = "card";
    card.innerHTML = `<div class="search-result-head"><div><div class="kind">${item.kind} draft</div><h3>${item.id} · ${item.title}</h3></div><strong>rev ${item.current_revision || 0}</strong></div><div class="meta">${item.summary || "No summary yet."}</div>`;
    const actions = document.createElement("div");
    actions.className = "card-actions";
    actions.appendChild(makeActionButton("Inspect", () => inspectRecord("item", item.id), "Draft Queue"));
    actions.appendChild(makeActionButton("Continue draft", async () => {
      await loadEditorItem(item.id, { activateMode: true });
      await inspectRecord("item", item.id, { activateMode: false });
    }, "Draft Queue"));
    actions.appendChild(makeActionButton("Review item", async () => {
      await inspectRecord("item", item.id);
      startApprovalFromContext("item", detailState.record);
    }, "Draft Queue"));
    card.appendChild(actions);
    draftReviewEl.appendChild(card);
  }
}

function renderRuns(items) {
  catalogState.runs = items;
  runListEl.innerHTML = "";
  for (const item of items) {
    const card = document.createElement("button");
    card.type = "button";
    card.className = "card selectable-card";
    card.innerHTML = `<div class="kind">${item.kind} run</div><h3>${item.id} · ${item.item_id}</h3><div class="meta">revision ${item.revision}\noutcome: ${item.outcome || "-"}\nplace: ${item.place_id || "-"}\nresources: ${(item.resource_ids || []).join(", ") || "-"}\n${item.notes || ""}</div>`;
    card.addEventListener("click", () => {
      inspectRecord("run", item.id).catch(handleError);
    });
    runListEl.appendChild(card);
  }
}

// Intent: Keep browser actions reachable without memorizing raw IDs by
// populating select helpers from the same projections that already drive the
// current lists, while preserving manual override fields for every workflow.
// Source: DI-lafor
function refreshActionCatalog() {
  renderSelect(searchPlaceSelectEl, catalogState.places, "Any place", (place) => `${place.id} · ${place.name}`);
  renderSelect(searchResourceSelectEl, catalogState.resources, "Any tool or resource", (resource) => `${resource.id} · ${resource.name}`);
  renderSelect(searchResponsibilitySelectEl, catalogState.responsibilities, "Any owning role", (responsibility) => `${responsibility.id} · ${responsibility.title}`);
  renderSelect(resourcePlaceSelectEl, catalogState.places, "Select a place", (place) => `${place.id} · ${place.name}`);
  renderSelect(itemResponsibilitySelectEl, catalogState.responsibilities, "Optional owner or review role", (responsibility) => `${responsibility.id} · ${responsibility.title}`);
  renderSelect(runItemSelectEl, catalogState.items, "Select a procedure or checklist", (item) => `${item.id} · ${item.title}`);
  renderSelect(runPlaceSelectEl, catalogState.places, "Optional place", (place) => `${place.id} · ${place.name}`);
  renderSelect(runResourceSelectEl, catalogState.resources, "Optional main tool or resource", (resource) => `${resource.id} · ${resource.name}`);
  renderSelect(runResponsibilitySelectEl, catalogState.responsibilities, "Optional owning role", (responsibility) => `${responsibility.id} · ${responsibility.title}`);
  renderSelect(evidenceRunSelectEl, catalogState.runs, "Select a run", (run) => `${run.id} · ${run.item_id} · ${run.outcome || "-"}`);
  renderApprovalTargetOptions();
  renderOperateContext();
}

function renderSelect(selectEl, items, placeholder, formatter) {
  const previous = selectEl.value;
  selectEl.innerHTML = "";
  const option = document.createElement("option");
  option.value = "";
  option.textContent = placeholder;
  selectEl.appendChild(option);
  for (const item of items) {
    const next = document.createElement("option");
    next.value = item.id;
    next.textContent = formatter(item);
    selectEl.appendChild(next);
  }
  if (previous) {
    selectEl.value = previous;
  }
}

function renderApprovalTargetOptions() {
  const form = approvalFormEl;
  const targetType = form.target_type.value;
  const items = targetType === "run" ? catalogState.runs : catalogState.items;
  const formatter = targetType === "run"
    ? (run) => `${run.id} · ${run.item_id} · ${run.outcome || "-"}`
    : (item) => `${item.id} · ${item.title} · rev ${item.current_revision}`;
  renderSelect(approvalTargetSelectEl, items, `Select a ${targetType === "run" ? "run" : "procedure or checklist"}`, formatter);
  syncApprovalDefaults();
}

// Intent: Keep the operate workspace anchored to the current record so logging
// work, attaching evidence, and capturing review decisions start from the
// operator's current context instead of a blank schema-heavy form. Source:
// DI-matub
function renderOperateContext() {
  operateContextEl.innerHTML = "";
  const title = document.createElement("strong");
  title.textContent = "Current browser context";
  const summary = document.createElement("div");
  summary.className = "meta";
  let canRun = false;
  let canEvidence = false;
  let canApprove = false;

  if (!detailState.record || !detailState.type) {
    summary.textContent = "Open one current record from Review or Browse to launch the most likely operation from that record.";
  } else if (detailState.type === "item") {
    summary.textContent = `${detailState.record.id} · ${detailState.record.title} is open. The fastest next steps are to log work from this revision or capture a review decision for it.`;
    canRun = true;
    canApprove = true;
  } else if (detailState.type === "run") {
    summary.textContent = `${detailState.record.id} is open. The fastest next steps are to attach evidence, review this run, or log the follow-on run with the same context.`;
    canRun = true;
    canEvidence = true;
    canApprove = true;
  } else if (detailState.type === "place") {
    summary.textContent = `${detailState.record.id} is open. Start the next run anchored to this location, then return here to review the resulting history.`;
    canRun = true;
  } else if (detailState.type === "resource") {
    summary.textContent = `${detailState.record.id} is open. Start the next run with this tool or resource already staged in context.`;
    canRun = true;
  } else if (detailState.type === "responsibility") {
    summary.textContent = `${detailState.record.id} is open. Start the next run with this ownership role already staged in context.`;
    canRun = true;
  }

  operateRunCurrentEl.disabled = !canRun;
  operateEvidenceCurrentEl.disabled = !canEvidence;
  operateApproveCurrentEl.disabled = !canApprove;
  operateRunCurrentEl.textContent = detailState.type === "run"
    ? "Log follow-on work from this run"
    : "Log work from current record";
  operateEvidenceCurrentEl.textContent = "Attach evidence from current record";
  operateApproveCurrentEl.textContent = detailState.type === "item"
    ? "Review this item"
    : detailState.type === "run"
      ? "Review this run"
      : "Review current record";

  operateContextEl.appendChild(title);
  operateContextEl.appendChild(summary);
}

function triggerOperateContext(kind) {
  if (!detailState.record || !detailState.type) {
    setActiveMode("review");
    setWorkspaceStatus("Open one current record first, then launch the matching operate action from that record.", "info");
    return;
  }
  if (kind === "run") {
    startRunFromContext(detailState.type, detailState.record);
    return;
  }
  if (kind === "evidence") {
    if (detailState.type !== "run") {
      setWorkspaceStatus("Evidence attaches to a run. Open one run first, then attach evidence from that record.", "info");
      setActiveMode("review");
      return;
    }
    startEvidenceFromRun(detailState.record);
    return;
  }
  if (kind === "approve") {
    if (detailState.type !== "item" && detailState.type !== "run") {
      setWorkspaceStatus("Review decisions apply to items or runs. Open one of those records first.", "info");
      setActiveMode("review");
      return;
    }
    startApprovalFromContext(detailState.type, detailState.record);
  }
}

function syncApprovalDefaults() {
  const form = approvalFormEl;
  if (form.target_type.value === "run" && detailState.type === "run" && detailState.record) {
    approvalTargetSelectEl.value = detailState.record.id;
    form.target_id.value = detailState.record.id;
    form.revision.value = "0";
    approvalTargetSummaryEl.textContent = `${detailState.record.id} · ${detailState.record.item_id} · ${detailState.record.outcome || "-"}`;
    return;
  }
  if (form.target_type.value === "knowledge_item" && detailState.type === "item" && detailState.record) {
    approvalTargetSelectEl.value = detailState.record.id;
    form.target_id.value = detailState.record.id;
    form.revision.value = String(detailState.record.current_revision || 0);
    approvalTargetSummaryEl.textContent = `${detailState.record.id} · ${detailState.record.title} · rev ${detailState.record.current_revision || 0}`;
    return;
  }
  const selected = approvalTargetSelectEl.selectedOptions[0];
  approvalTargetSummaryEl.textContent = selected && approvalTargetSelectEl.value
    ? selected.textContent
    : `Choose the ${form.target_type.value === "run" ? "run" : "procedure or checklist"} you are reviewing.`;
}

// Intent: Let browser operators jump from the current inspected or searched
// record into the matching operate form with the right defaults already staged,
// instead of re-deriving the same context from blank generic forms. Source:
// DI-mitav
function focusOperateForm(formEl, message, focusSelector = "textarea, input, select") {
  setActiveMode("operate");
  if (formEl === document.getElementById("run-form")) {
    openOperateStage("run");
  } else if (formEl === document.getElementById("evidence-form")) {
    openOperateStage("evidence");
  } else if (formEl === approvalFormEl) {
    openOperateStage("approval");
  }
  formEl.scrollIntoView({ behavior: "smooth", block: "start" });
  const target = formEl.querySelector(focusSelector);
  if (target) {
    target.focus({ preventScroll: true });
  }
  setWorkspaceStatus(message, "info");
}

// Intent: Preserve the generic run form while making item, place, resource,
// responsibility, and run context act like first-class launch points for a new
// run record. Source: DI-mitav
function startRunFromContext(type, record) {
  const form = document.getElementById("run-form");
  applyContextDefaults(type, record);
  if (type === "item") {
    const revision = record.current_revision || form.revision.value || 1;
    form.revision.value = String(revision);
    focusOperateForm(form, `Run form primed for ${record.id} revision ${revision}.`, "input[name='outcome']");
    return;
  }
  if (type === "run") {
    focusOperateForm(form, `Run form primed from ${record.id} so you can record the next pass with matching context.`, "input[name='outcome']");
    return;
  }
  if (type === "place") {
    focusOperateForm(form, `Run form primed for work at ${record.id}.`, "select[name='kind']");
    return;
  }
  if (type === "resource") {
    focusOperateForm(form, `Run form primed for work involving ${record.id}.`, "select[name='kind']");
    return;
  }
  if (type === "responsibility") {
    focusOperateForm(form, `Run form primed for work owned by ${record.id}.`, "select[name='kind']");
  }
}

// Intent: Make evidence attachment feel like a follow-on action from a run the
// operator is already looking at, instead of a detached generic upload form.
// Source: DI-mitav
function startEvidenceFromRun(run) {
  applyContextDefaults("run", run);
  const form = document.getElementById("evidence-form");
  focusOperateForm(form, `Evidence form primed for ${run.id}.`, "input[name='summary']");
}

// Intent: Keep approvals reachable as a generic form while making the common
// item/run approval path start from the current record context. Source:
// DI-mitav
function startApprovalFromContext(type, record) {
  if (type !== "item" && type !== "run") {
    throw new Error(`Unsupported approval context ${type}`);
  }
  applyContextDefaults(type, record);
  const form = approvalFormEl;
  focusOperateForm(
    form,
    `${type === "run" ? "Run" : "Item"} approval primed for ${record.id}.`,
    "input[name='role']",
  );
}

// Intent: Surface repeated receiving and count problems as grouped hotspots so
// operators can see where issues cluster before drilling into one record at a
// time. Source: DI-pogul
function renderProblemReview(summary) {
  problemReviewEl.innerHTML = "";
  const groups = [
    ["Places with repeated problems", "place", summary.place_groups || []],
    ["Resources with repeated problems", "resource", summary.resource_groups || []],
  ];
  for (const [title, targetType, items] of groups) {
    if (!items.length) {
      continue;
    }
    const block = document.createElement("div");
    block.className = "search-group";
    const heading = document.createElement("h3");
    heading.textContent = `${title} (${items.length})`;
    block.appendChild(heading);
    for (const item of items) {
      const card = document.createElement("article");
      card.className = "card";
      const highlights = (item.highlights || []).slice(0, 3).join(" · ") || "no highlights";
      card.innerHTML = `<div class="hotspot-head"><div><div class="kind">${item.kind}</div><h3>${item.group_id} · ${item.name}</h3></div><strong>${item.problem_count} problems</strong></div><div class="meta">receiving: ${item.receiving_problems} · inventory: ${item.inventory_problems}\n${highlights}</div>`;
      const actions = document.createElement("div");
      actions.className = "card-actions";
      actions.appendChild(makeActionButton("Inspect context", () => inspectRecord(targetType, item.group_id), "Problem Review"));
      actions.appendChild(makeActionButton("Problem runs here", () => runSearch({ [`${targetType}_id`]: item.group_id, problem: true }), "Problem Review"));
      actions.appendChild(makeActionButton("Receiving here", () => runSearch({ [`${targetType}_id`]: item.group_id, kind: "receiving_check" }), "Problem Review"));
      actions.appendChild(makeActionButton("Inventory here", () => runSearch({ [`${targetType}_id`]: item.group_id, kind: "inventory_audit" }), "Problem Review"));
      card.appendChild(actions);
      block.appendChild(card);
    }
    problemReviewEl.appendChild(block);
  }
  if (!problemReviewEl.children.length) {
    const empty = document.createElement("div");
    empty.className = "meta";
    empty.textContent = "No repeated receiving or count problems recorded yet.";
    problemReviewEl.appendChild(empty);
  }
}

// Intent: Let browser operators combine structured filters with free-text
// search so they can drill into one slice of the operational graph without
// manually cross-referencing IDs outside the app. Source: DI-honus
function getSearchFilters(form) {
  return {
    q: form.q.value.trim(),
    kind: form.kind.value.trim(),
    status: form.status.value.trim(),
    outcome: form.outcome.value.trim(),
    place_id: firstPresent(form.place_id.value, form.place_id_select.value),
    resource_id: firstPresent(form.resource_id.value, form.resource_id_select.value),
    responsibility_id: firstPresent(form.responsibility_id.value, form.responsibility_id_select.value),
    problem: form.dataset.problem === "true",
  };
}

function hasAdvancedSearchFilters(filters) {
  return Boolean(
    filters.kind
    || filters.status
    || filters.outcome
    || filters.place_id
    || filters.resource_id
    || filters.responsibility_id
    || filters.problem,
  );
}

function buildSearchParams(filters) {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(filters)) {
    if (typeof value === "boolean") {
      if (value) {
        params.set(key, "true");
      }
      continue;
    }
    if (value) {
      params.set(key, value);
    }
  }
  return params;
}

function renderSearchResults(filters, payload) {
  setReviewLane("search");
  searchResultsEl.innerHTML = "";
  searchRawEl.hidden = false;
  searchRawEl.textContent = JSON.stringify(payload, null, 2);
  searchDebugEl.hidden = false;
  searchDebugEl.open = false;
  searchActiveEl.textContent = formatSearchFilters(payload.filters || filters);
  const groups = [
    ["places", "place"],
    ["resources", "resource"],
    ["responsibilities", "responsibility"],
    ["items", "item"],
    ["runs", "run"],
  ];
  for (const [key, type] of groups) {
    const items = payload[key] || [];
    if (items.length === 0) {
      continue;
    }
    const block = document.createElement("div");
    block.className = "search-group";
    const heading = document.createElement("h3");
    heading.textContent = `${key} (${items.length})`;
    block.appendChild(heading);
    for (const item of items) {
      const card = document.createElement("article");
      card.className = "card";
      card.innerHTML = `<div class="search-result-head"><div><div class="kind">${type}</div><h3>${displayRecordID(item)}</h3></div></div><div class="meta">${searchSummary(type, item)}</div>`;
      const actions = document.createElement("div");
      actions.className = "card-actions";
      actions.appendChild(makeActionButton("Inspect", () => inspectRecord(type, item.id), "Search"));
      if (type === "item") {
        actions.appendChild(makeActionButton("Open draft", async () => {
          await loadEditorItem(item.id, { activateMode: true });
          await inspectRecord("item", item.id, { activateMode: false });
        }, "Search"));
        actions.appendChild(makeActionButton("Record run", async () => {
          await inspectRecord("item", item.id);
          startRunFromContext("item", detailState.record);
        }, "Search"));
        actions.appendChild(makeActionButton("Approve item", async () => {
          await inspectRecord("item", item.id);
          startApprovalFromContext("item", detailState.record);
        }, "Search"));
      }
      if (type === "run") {
        actions.appendChild(makeActionButton("Item", () => inspectRecord("item", item.item_id), "Search"));
        actions.appendChild(makeActionButton("Add evidence", async () => {
          await inspectRecord("run", item.id);
          startEvidenceFromRun(detailState.record);
        }, "Search"));
        actions.appendChild(makeActionButton("Approve run", async () => {
          await inspectRecord("run", item.id);
          startApprovalFromContext("run", detailState.record);
        }, "Search"));
      }
      card.appendChild(actions);
      block.appendChild(card);
    }
    searchResultsEl.appendChild(block);
  }
  if (!searchResultsEl.children.length) {
    const empty = document.createElement("div");
    empty.className = "meta";
    empty.textContent = filters.q ? `No results for "${filters.q}".` : "No results.";
    searchResultsEl.appendChild(empty);
  }
}

// Intent: Make the browser decisive about the most likely next action after a
// record is open, while preserving the broader related-action list underneath
// and making the dominant review path unambiguous on the landing surface.
// Source: DI-sorik; DI-pafur
function renderDetailPrimary(type, record) {
  detailPrimaryEl.innerHTML = "";
  const title = document.createElement("strong");
  title.textContent = "Next Step";
  const summary = document.createElement("p");
  summary.className = "meta";
  const actions = document.createElement("div");
  actions.className = "detail-primary-actions";

  if (type === "item") {
    summary.textContent = "Revise this item or record work from its current revision, then approve that durable revision when it is ready.";
    actions.appendChild(makePrimaryActionButton("Continue draft", async () => {
      await loadEditorItem(record.id, { activateMode: true });
      await inspectRecord("item", record.id, { activateMode: false });
    }, "Record Inspector"));
    actions.appendChild(makePrimaryActionButton("Record run", () => startRunFromContext("item", record), "Record Inspector"));
    actions.appendChild(makePrimaryActionButton("Approve item", () => startApprovalFromContext("item", record), "Record Inspector"));
  } else if (type === "run") {
    summary.textContent = "Capture supporting evidence first, then record the review decision or open the durable item behind this run.";
    actions.appendChild(makePrimaryActionButton("Add evidence", () => startEvidenceFromRun(record), "Record Inspector"));
    actions.appendChild(makePrimaryActionButton("Approve run", () => startApprovalFromContext("run", record), "Record Inspector"));
    actions.appendChild(makePrimaryActionButton("Open item", () => inspectRecord("item", record.item_id), "Record Inspector"));
  } else if (type === "place") {
    summary.textContent = "Review repeated problems for this place or start a run that is explicitly anchored to this location.";
    actions.appendChild(makePrimaryActionButton("Search problems here", () => runSearch({ place_id: record.id, problem: true }), "Record Inspector"));
    actions.appendChild(makePrimaryActionButton("Record run here", () => startRunFromContext("place", record), "Record Inspector"));
  } else if (type === "resource") {
    summary.textContent = "Use this resource as the anchor for the next run, or review the receiving and count history attached to it.";
    actions.appendChild(makePrimaryActionButton("Search problems here", () => runSearch({ resource_id: record.id, problem: true }), "Record Inspector"));
    actions.appendChild(makePrimaryActionButton("Record run with this resource", () => startRunFromContext("resource", record), "Record Inspector"));
  } else if (type === "responsibility") {
    summary.textContent = "Review the work owned by this responsibility, then record the next run using the same ownership context.";
    actions.appendChild(makePrimaryActionButton("Search receiving problems", () => runSearch({ responsibility_id: record.id, problem: true }), "Record Inspector"));
    actions.appendChild(makePrimaryActionButton("Record run for this responsibility", () => startRunFromContext("responsibility", record), "Record Inspector"));
  }

  if (!actions.children.length) {
    detailPrimaryEl.hidden = true;
    return;
  }
  detailPrimaryEl.hidden = false;
  detailPrimaryEl.appendChild(title);
  detailPrimaryEl.appendChild(summary);
  detailPrimaryEl.appendChild(actions);
}

function formatSearchFilters(filters) {
  const labels = [];
  if (filters.query || filters.q) {
    labels.push(`searching for "${filters.query || filters.q}"`);
  }
  if (filters.kind) {
    labels.push(`work type ${filters.kind}`);
  }
  if (filters.status) {
    labels.push(`state ${filters.status}`);
  }
  if (filters.outcome) {
    labels.push(`result ${filters.outcome}`);
  }
  if (filters.place_id) {
    labels.push(`at ${filters.place_id}`);
  }
  if (filters.resource_id) {
    labels.push(`with ${filters.resource_id}`);
  }
  if (filters.responsibility_id) {
    labels.push(`owned by ${filters.responsibility_id}`);
  }
  if (filters.problem) {
    labels.push("problem-focused");
  }
  if (labels.length === 0) {
    return "No active search filters.";
  }
  return `Showing results for ${labels.join(" · ")}.`;
}

function searchSummary(type, item) {
  switch (type) {
    case "place":
      return `${item.name || ""}\n${item.summary || ""}`;
    case "resource":
      return `${item.name || ""}\nlocated at ${item.place_id || "-"}`;
    case "responsibility":
      return `${item.title || ""}\n${item.summary || ""}`;
    case "item":
      return `${item.title || ""}\n${item.kind || ""} · state ${item.status || ""}`;
    case "run":
      return `${item.item_id || ""}\nrevision ${item.revision || 0} · result ${item.outcome || "-"}`;
    default:
      return "";
  }
}

// Intent: Let operators inspect and traverse the current operational graph in
// the browser without manually copying IDs between separate lists, while
// keeping the existing local HTTP runtime and record model unchanged. Source:
// DI-vopuk
async function inspectRecord(type, id, options = {}) {
  if (options.activateMode !== false) {
    setActiveMode("review");
  }
  detailState.type = type;
  detailState.id = id;
  detailState.record = null;
  detailMetaEl.textContent = `Loading ${type} ${id}...`;
  detailSummaryEl.innerHTML = "";
  detailPrimaryEl.hidden = true;
  detailPrimaryEl.innerHTML = "";
  detailActionsEl.innerHTML = "";
  detailReviewEl.innerHTML = "";
  detailTimelineEl.innerHTML = "";
  const record = await getJSON(detailPath(type, id));
  detailState.record = record;
  detailMetaEl.textContent = detailSummary(type, record);
  renderDetailSummary(type, record);
  renderDetailPrimary(type, record);
  renderDetailReview(type, record);
  renderDetailTimeline(record.timeline || []);
  detailJSONEl.textContent = JSON.stringify(record, null, 2);
  renderDetailActions(type, record);
  applyContextDefaults(type, record);
  renderOperateContext();
}

function detailPath(type, id) {
  switch (type) {
    case "place":
      return `/api/places/${id}`;
    case "resource":
      return `/api/resources/${id}`;
    case "responsibility":
      return `/api/responsibilities/${id}`;
    case "item":
      return `/api/items/${id}`;
    case "run":
      return `/api/runs/${id}`;
    default:
      throw new Error(`Unsupported detail type ${type}`);
  }
}

function detailSummary(type, record) {
  switch (type) {
    case "place":
      return `${displayRecordID(record)} · ${record.kind} · ${record.name}`;
    case "resource":
      return `${displayRecordID(record)} · ${record.kind} · ${record.name}`;
    case "responsibility":
      return `${displayRecordID(record)} · responsibility · ${record.title}`;
    case "item":
      return `${displayRecordID(record)} · ${record.kind} · ${record.title} · ${record.status}`;
    case "run":
      return `${displayRecordID(record)} · ${record.kind} run · ${record.item_id}`;
    default:
      return `${type} ${displayRecordID(record) || record.id || ""}`;
  }
}

// Intent: Turn the record inspector into a real operational detail view with
// human-readable summaries and timelines, instead of forcing users to read raw
// JSON for every place, resource, responsibility, item, or run. Source:
// DI-honus
function renderDetailSummary(type, record) {
  detailSummaryEl.innerHTML = "";
  const stats = detailStats(type, record);
  for (const [label, value] of stats) {
    const card = document.createElement("div");
    card.className = "detail-stat";
    card.innerHTML = `<strong>${value}</strong><span>${label}</span>`;
    detailSummaryEl.appendChild(card);
  }
}

function detailStats(type, record) {
  switch (type) {
    case "place":
      return [
        ["Parent", record.parent_id || "-"],
        ["Children", (record.child_place_ids || []).length],
        ["Resources", (record.resource_ids || []).length],
        ["Runs", (record.related_runs || []).length],
        ["Events", (record.timeline || []).length],
      ];
    case "resource":
      return [
        ["Place", record.place_id || "-"],
        ["Tags", (record.tags || []).length],
        ["Runs", (record.related_runs || []).length],
        ["Links", (record.links || []).length],
        ["Events", (record.timeline || []).length],
      ];
    case "responsibility":
      return [
        ["Team", record.team || "-"],
        ["Items", (record.linked_item_ids || []).length],
        ["Runs", (record.linked_run_ids || []).length],
        ["Links", (record.links || []).length],
        ["Related runs", (record.related_runs || []).length],
        ["Events", (record.timeline || []).length],
      ];
    case "item":
      return [
        ["Status", record.status || "-"],
        ["Revision", record.current_revision || 0],
        ["Approvals", (record.approvals || []).length],
        ["Events", (record.timeline || []).length],
      ];
    case "run":
      return [
        ["Outcome", record.outcome || "-"],
        ["Revision", record.revision || 0],
        ["Evidence", (record.evidence || []).length],
        ["Approvals", (record.approvals || []).length],
      ];
    default:
      return [["Events", (record.timeline || []).length]];
  }
}

// Intent: Turn context inspection into actionable history drilldown so browser
// operators can jump directly into receiving/count/problem searches from the
// record they are already reviewing, while keeping problem drilldowns aligned
// with the grouped hotspot logic. Source: DI-vafuk; DI-vemur
function renderDetailActions(type, record) {
  detailActionsEl.innerHTML = "";
  const links = [];
  if (type === "resource") {
    links.push(["Record run with this resource", "operate", () => startRunFromContext("resource", record)]);
    if (record.place_id) {
      links.push(["Open place", "place", record.place_id]);
    }
    links.push(["Search receiving here", "search", { kind: "receiving_check", resource_id: record.id }]);
    links.push(["Search counts here", "search", { kind: "inventory_audit", resource_id: record.id }]);
    links.push(["Search problems here", "search", { resource_id: record.id, problem: true }]);
    for (const run of record.related_runs || []) {
      links.push([`Run ${run.id}`, "run", run.id]);
    }
  }
  if (type === "item") {
    links.push(["Record run for this item", "operate", () => startRunFromContext("item", record)]);
    links.push(["Approve this item", "operate", () => startApprovalFromContext("item", record)]);
    links.push(["Open live draft", "item", record.id]);
    for (const responsibilityID of record.responsibility_ids || []) {
      links.push([`Responsibility ${responsibilityID}`, "responsibility", responsibilityID]);
    }
    for (const run of record.related_runs || []) {
      links.push([`Run ${run.id}`, "run", run.id]);
    }
  }
  if (type === "run") {
    links.push(["Add evidence to this run", "operate", () => startEvidenceFromRun(record)]);
    links.push(["Approve this run", "operate", () => startApprovalFromContext("run", record)]);
    links.push(["Record follow-on run", "operate", () => startRunFromContext("run", record)]);
    links.push(["Open item", "item", record.item_id]);
    if (record.place_id) {
      links.push(["Open place", "place", record.place_id]);
    }
    for (const resourceID of record.resource_ids || []) {
      links.push([`Resource ${resourceID}`, "resource", resourceID]);
    }
    for (const responsibilityID of record.responsibility_ids || []) {
      links.push([`Responsibility ${responsibilityID}`, "responsibility", responsibilityID]);
    }
  }
  if (type === "place") {
    links.push(["Record run here", "operate", () => startRunFromContext("place", record)]);
    links.push(["Search receiving here", "search", { kind: "receiving_check", place_id: record.id }]);
    links.push(["Search counts here", "search", { kind: "inventory_audit", place_id: record.id }]);
    links.push(["Search problems here", "search", { place_id: record.id, problem: true }]);
    for (const resourceID of record.resource_ids || []) {
      links.push([`Resource ${resourceID}`, "resource", resourceID]);
    }
    for (const childID of record.child_place_ids || []) {
      links.push([`Child place ${childID}`, "place", childID]);
    }
    for (const run of record.related_runs || []) {
      links.push([`Run ${run.id}`, "run", run.id]);
    }
  }
  if (type === "responsibility") {
    links.push(["Record run for this responsibility", "operate", () => startRunFromContext("responsibility", record)]);
    links.push(["Search receiving runs", "search", { kind: "receiving_check", responsibility_id: record.id }]);
    links.push(["Search inventory counts", "search", { kind: "inventory_audit", responsibility_id: record.id }]);
    links.push(["Search receiving problems", "search", { responsibility_id: record.id, problem: true }]);
    for (const itemID of record.linked_item_ids || []) {
      links.push([`Item ${itemID}`, "item", itemID]);
    }
    for (const runID of record.linked_run_ids || []) {
      links.push([`Run ${runID}`, "run", runID]);
    }
  }
  for (const [label, nextType, nextID] of links) {
    if (nextType === "operate") {
      detailActionsEl.appendChild(makeActionButton(label, nextID, "Record Inspector"));
      continue;
    }
    if (nextType === "search") {
      detailActionsEl.appendChild(makeActionButton(label, () => runSearch(nextID), "Record Inspector"));
      continue;
    }
    detailActionsEl.appendChild(makeActionButton(label, () => inspectRecord(nextType, nextID), "Record Inspector"));
  }
}

// Intent: Keep every current browser form reachable while prefilling the
// high-frequency run, evidence, and approval flows from the record the operator
// is already reviewing, so the browser stops depending on raw ID memorization.
// Source: DI-lafor
function applyContextDefaults(type, record) {
  const runForm = document.getElementById("run-form");
  const evidenceForm = document.getElementById("evidence-form");
  const approvalForm = approvalFormEl;

  if (type === "item") {
    runItemSelectEl.value = record.id;
    runItemIDEl.value = record.id;
    runForm.revision.value = String(record.current_revision || runForm.revision.value);
    approvalForm.target_type.value = "knowledge_item";
    renderApprovalTargetOptions();
    approvalTargetSelectEl.value = record.id;
    approvalForm.target_id.value = record.id;
    approvalForm.revision.value = String(record.current_revision || 0);
    if ((record.responsibility_ids || [])[0]) {
      runResponsibilitySelectEl.value = record.responsibility_ids[0];
    }
  }
  if (type === "run") {
    evidenceRunSelectEl.value = record.id;
    evidenceForm.run_id.value = record.id;
    approvalForm.target_type.value = "run";
    renderApprovalTargetOptions();
    approvalTargetSelectEl.value = record.id;
    approvalForm.target_id.value = record.id;
    approvalForm.revision.value = "0";
    if (record.item_id) {
      runItemSelectEl.value = record.item_id;
      runItemIDEl.value = record.item_id;
    }
    if (record.place_id) {
      runPlaceSelectEl.value = record.place_id;
      runForm.place_id.value = record.place_id;
    }
    if ((record.resource_ids || [])[0]) {
      runResourceSelectEl.value = record.resource_ids[0];
    }
    if ((record.responsibility_ids || [])[0]) {
      runResponsibilitySelectEl.value = record.responsibility_ids[0];
    }
    if (record.revision) {
      runForm.revision.value = String(record.revision);
    }
  }
  if (type === "place") {
    runPlaceSelectEl.value = record.id;
    runForm.place_id.value = record.id;
    resourcePlaceSelectEl.value = record.id;
  }
  if (type === "resource") {
    runResourceSelectEl.value = record.id;
    if (record.place_id) {
      runPlaceSelectEl.value = record.place_id;
      runForm.place_id.value = record.place_id;
    }
  }
  if (type === "responsibility") {
    runResponsibilitySelectEl.value = record.id;
    itemResponsibilitySelectEl.value = record.id;
  }
  syncApprovalDefaults();
  renderOperateContext();
}

// Intent: Reuse the structured search form as the single drilldown path so
// direct inspector actions and manual operator searches stay behaviorally
// identical, including problem-only review drilldowns. Source: DI-vafuk;
// DI-vemur
async function runSearch(filters) {
  clearWorkspaceStatus();
  setReviewLane("search");
  const form = document.getElementById("search-form");
  form.q.value = filters.q || "";
  form.kind.value = filters.kind || "";
  form.status.value = filters.status || "";
  form.outcome.value = filters.outcome || "";
  syncSearchSelectAndOverride(searchPlaceSelectEl, form.place_id, filters.place_id || "");
  syncSearchSelectAndOverride(searchResourceSelectEl, form.resource_id, filters.resource_id || "");
  syncSearchSelectAndOverride(searchResponsibilitySelectEl, form.responsibility_id, filters.responsibility_id || "");
  form.dataset.problem = filters.problem ? "true" : "false";
  searchAdvancedEl.open = hasAdvancedSearchFilters(filters);
  const effectiveFilters = getSearchFilters(form);
  const payload = await getJSON(`/api/search?${buildSearchParams(effectiveFilters).toString()}`);
  renderSearchResults(effectiveFilters, payload);
  clearWorkspaceStatus();
}

function clearSearch() {
  const form = document.getElementById("search-form");
  form.reset();
  form.dataset.problem = "false";
  if (searchPlaceSelectEl) {
    searchPlaceSelectEl.value = "";
  }
  if (searchResourceSelectEl) {
    searchResourceSelectEl.value = "";
  }
  if (searchResponsibilitySelectEl) {
    searchResponsibilitySelectEl.value = "";
  }
  searchAdvancedEl.open = false;
  searchResultsEl.innerHTML = "";
  searchActiveEl.textContent = "No active search filters.";
  searchRawEl.hidden = true;
  searchRawEl.textContent = "";
  searchDebugEl.hidden = true;
  searchDebugEl.open = false;
}

function makeActionButton(label, action, context) {
  const button = document.createElement("button");
  button.type = "button";
  button.className = "button-tertiary";
  button.textContent = label;
  button.addEventListener("click", () => {
    Promise.resolve(action()).catch((error) => handleError(error, context));
  });
  return button;
}

function makePrimaryActionButton(label, action, context) {
  const button = document.createElement("button");
  button.type = "button";
  button.textContent = label;
  button.addEventListener("click", () => {
    Promise.resolve(action()).catch((error) => handleError(error, context));
  });
  return button;
}

function renderDetailTimeline(events) {
  detailTimelineEl.innerHTML = "";
  if (!events.length) {
    const empty = document.createElement("div");
    empty.className = "meta";
    empty.textContent = "No timeline events recorded yet.";
    detailTimelineEl.appendChild(empty);
    return;
  }
  for (const event of events) {
    const card = document.createElement("div");
    card.className = "timeline-entry";
    card.innerHTML = `<div class="timeline-head"><span class="kind">${event.type}</span><span class="meta">${event.timestamp || ""}</span></div><div class="timeline-body">${timelineSummary(event)}</div>`;
    detailTimelineEl.appendChild(card);
  }
}

// Intent: Make run, item, and context review practical in the browser by
// surfacing revisions, approvals, evidence, and related run history as
// first-class panels instead of hiding them inside raw record JSON. Source:
// DI-honus; DI-julos; DI-vemok; DI-zemok
function renderDetailReview(type, record) {
  detailReviewEl.innerHTML = "";
  const sections = [];
  if (type === "item") {
    sections.push(["Revisions", (record.revisions || []).map((revision) => `${revision.number} · ${revision.title} · ${revision.author} · ${revision.created_at}`)]);
    sections.push(["Approvals", (record.approvals || []).map((approval) => `${approval.decision} · ${approval.role} · ${approval.actor} · rev ${approval.revision}${approval.notes ? ` · ${approval.notes}` : ""}`)]);
    sections.push(["Related runs", (record.related_runs || []).map((run) => `${run.id} · rev ${run.revision} · ${run.outcome || "-"} · ${run.created_at}`)]);
    if (record.kind === "receiving_check") {
      sections.push(["Receiving history", receivingContextEntries(record.related_runs || [])]);
    }
    if (record.kind === "inventory_audit") {
      sections.push(["Inventory count history", inventoryContextEntries(record.related_runs || [])]);
    }
  }
  if (type === "run") {
    sections.push(["Evidence", (record.evidence || []).map((evidence) => `${evidence.summary} · ${evidence.actor} · ${evidence.created_at}${evidence.attachment_name ? ` · attachment ${evidence.attachment_name}` : ""}`)]);
    sections.push(["Approvals", (record.approvals || []).map((approval) => `${approval.decision} · ${approval.role} · ${approval.actor}${approval.notes ? ` · ${approval.notes}` : ""}`)]);
    sections.push(["Responsibilities", (record.responsibility_ids || []).map((id) => id)]);
    if (record.kind === "receiving_check") {
      sections.push(["Receiving review", receivingEvidenceEntries(record.evidence || [])]);
    }
    if (record.kind === "inventory_audit") {
      sections.push(["Inventory discrepancy", inventoryEvidenceEntries(record.evidence || [])]);
    }
  }
  if (type === "responsibility") {
    sections.push(["Linked items", (record.linked_item_ids || []).map((id) => id)]);
    sections.push(["Linked runs", (record.linked_run_ids || []).map((id) => id)]);
    sections.push(["Related runs", (record.related_runs || []).map((run) => `${run.id} · ${run.kind} · rev ${run.revision} · ${run.outcome || "-"} · ${run.created_at}`)]);
    sections.push(["Receiving context review", receivingContextEntries(record.related_runs || [])]);
    sections.push(["Inventory count history", inventoryContextEntries(record.related_runs || [])]);
  }
  if (type === "place" || type === "resource") {
    sections.push(["Related runs", (record.related_runs || []).map((run) => `${run.id} · ${run.kind} · rev ${run.revision} · ${run.outcome || "-"} · ${run.created_at}`)]);
    sections.push(["Receiving context review", receivingContextEntries(record.related_runs || [])]);
    sections.push(["Inventory count history", inventoryContextEntries(record.related_runs || [])]);
  }
  for (const [title, entries] of sections) {
    if (!entries.length) {
      continue;
    }
    const panel = document.createElement("section");
    panel.className = "detail-subpanel";
    const heading = document.createElement("h3");
    heading.textContent = title;
    panel.appendChild(heading);
    const list = document.createElement("div");
    list.className = "detail-list";
    for (const entry of entries) {
      const row = document.createElement("div");
      row.className = "detail-list-item";
      row.textContent = entry;
      list.appendChild(row);
    }
    panel.appendChild(list);
    detailReviewEl.appendChild(panel);
  }
}

// Intent: Make inventory audit runs easier to review than generic evidence
// blobs by surfacing discrepancy/count facts and audit history as named review
// sections in the browser inspector. Source: DI-pojul
function inventoryEvidenceEntries(evidenceList) {
  const entries = [];
  for (const evidence of evidenceList) {
    const facts = formatEvidenceFacts(evidence.facts || {});
    entries.push(`${evidence.summary} · ${facts}${evidence.attachment_name ? ` · attachment ${evidence.attachment_name}` : ""}`);
  }
  return entries;
}

function inventoryAuditEntries(runs) {
  return runs
    .filter((run) => run.kind === "inventory_audit")
    .map((run) => `${run.id} · rev ${run.revision} · ${run.outcome || "-"} · ${run.created_at}`);
}

// Intent: Make inventory history useful from context anchors like places,
// resources, and responsibilities by showing run-level count/discrepancy facts
// instead of only bare related-run ids. Source: DI-zemok
function inventoryContextEntries(runs) {
  const entries = [];
  for (const run of runs.filter((value) => value.kind === "inventory_audit")) {
    const evidence = inventoryEvidenceEntries(run.evidence || []);
    if (!evidence.length) {
      entries.push(`${run.id} · rev ${run.revision} · ${run.outcome || "-"} · ${run.created_at} · no evidence facts`);
      continue;
    }
    for (const detail of evidence) {
      entries.push(`${run.id} · rev ${run.revision} · ${run.outcome || "-"} · ${detail}`);
    }
  }
  return entries;
}

function formatEvidenceFacts(facts) {
  const keys = Object.keys(facts).sort();
  if (keys.length === 0) {
    return "no facts";
  }
  return keys.map((key) => `${key}: ${facts[key]}`).join(" · ");
}

// Intent: Make receiving and inbound inspection work readable as its own
// operational flow, not just as generic evidence text or an inventory-only
// special case. Source: DI-vemok
function receivingEvidenceEntries(evidenceList) {
  const entries = [];
  for (const evidence of evidenceList) {
    const facts = formatEvidenceFacts(evidence.facts || {});
    entries.push(`${evidence.summary} · ${facts}${evidence.attachment_name ? ` · attachment ${evidence.attachment_name}` : ""}`);
  }
  return entries;
}

function receivingRunEntries(runs) {
  return runs
    .filter((run) => run.kind === "receiving_check")
    .map((run) => `${run.id} · rev ${run.revision} · ${run.outcome || "-"} · ${run.created_at}`);
}

function receivingContextEntries(runs) {
  const entries = [];
  for (const run of runs.filter((value) => value.kind === "receiving_check")) {
    const evidence = receivingEvidenceEntries(run.evidence || []);
    if (!evidence.length) {
      entries.push(`${run.id} · rev ${run.revision} · ${run.outcome || "-"} · ${run.created_at} · no evidence facts`);
      continue;
    }
    for (const detail of evidence) {
      entries.push(`${run.id} · rev ${run.revision} · ${run.outcome || "-"} · ${detail}`);
    }
  }
  return entries;
}

function timelineSummary(event) {
  const fragments = [];
  if (event.actor) {
    fragments.push(`actor: ${event.actor}`);
  }
  if (event.title) {
    fragments.push(`title: ${event.title}`);
  }
  if (event.summary) {
    fragments.push(`summary: ${event.summary}`);
  }
  if (event.status) {
    fragments.push(`status: ${event.status}`);
  }
  if (event.revision) {
    fragments.push(`revision: ${event.revision}`);
  }
  if (event.decision) {
    fragments.push(`decision: ${event.decision}`);
  }
  if (event.outcome) {
    fragments.push(`outcome: ${event.outcome}`);
  }
  if (event.relation) {
    fragments.push(`relation: ${event.relation}`);
  }
  if (event.notes) {
    fragments.push(event.notes);
  }
  return fragments.join(" · ") || event.entity_id || "event";
}

function renderEditorState(state) {
  editorState.version = state.version;
  editorState.title = state.title;
  editorState.status = state.status;
  editorState.currentRevision = state.current_revision;
  editorState.participantCount = (state.participants || []).length;
  editorState.lastRenderedBody = state.body;
  editorMetaEl.textContent = `${state.title} · status ${state.status} · live v${state.version} · current revision ${state.current_revision}${editorState.dirty ? " · local edits pending" : ""}`;
  renderEditorStatusCards(state, editorBodyEl.value || state.body);
  editorParticipantsEl.innerHTML = "";
  for (const participant of state.participants || []) {
    const pill = document.createElement("span");
    pill.className = "participant-pill";
    pill.style.setProperty("--participant-color", participant.color || "#0c6d62");
    pill.textContent = `${participant.display_name} @ ${participant.cursor}:${participant.head}${participant.typing ? " typing" : ""}`;
    editorParticipantsEl.appendChild(pill);
  }
  if (!editorState.dirty || editorBodyEl.value === "" || editorBodyEl.value === editorState.lastRenderedBody) {
    editorBodyEl.value = state.body;
  }
}

function liveSocketURL(itemID) {
  const url = new URL(`/api/items/${itemID}/live/socket`, window.location.href);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  return url.toString();
}

// Intent: Make browser authoring feel like sustained drafting work instead of
// only an administrative form by surfacing live draft health, scale, and
// collaboration state beside the editor. Source: DI-rofek
function renderEditorStatusCards(state, body) {
  editorStatusCardsEl.innerHTML = "";
  const words = body.trim() ? body.trim().split(/\s+/).length : 0;
  const paragraphs = body.trim() ? body.trim().split(/\n\s*\n/).length : 0;
  const cards = [
    ["Live Version", `v${state.version}`],
    ["Current Revision", String(state.current_revision)],
    ["Words", String(words)],
    ["Paragraphs", String(paragraphs)],
    ["Participants", String((state.participants || []).length)],
    ["Status", state.status || "-"],
  ];
  for (const [label, value] of cards) {
    const card = document.createElement("div");
    card.className = "author-stat";
    card.innerHTML = `<strong>${value}</strong><span>${label}</span>`;
    editorStatusCardsEl.appendChild(card);
  }
}

async function refresh(selectedItemID = editorState.itemID) {
  const [dashboard, problemReview, places, resources, responsibilities, items, runs] = await Promise.all([
    getJSON("/api/dashboard"),
    getJSON("/api/problem-review"),
    getJSON("/api/places"),
    getJSON("/api/resources"),
    getJSON("/api/responsibilities"),
    getJSON("/api/items"),
    getJSON("/api/runs"),
  ]);
  renderStats(dashboard);
  renderProblemReview(problemReview);
  renderPlaces(places.places || []);
  renderResources(resources.resources || []);
  renderResponsibilities(responsibilities.responsibilities || []);
  renderKnowledgeItems(items.items || []);
  renderDraftQueue(items.items || []);
  renderRuns(runs.runs || []);
  refreshActionCatalog();
  if (!selectedItemID && (items.items || []).length > 0) {
    const drafts = (items.items || []).filter((item) => item.status === "draft");
    selectedItemID = (drafts[0] || items.items[0]).id;
  }
  const tasks = [];
  if (selectedItemID) {
    tasks.push(inspectRecord("item", selectedItemID, { activateMode: false }));
  }
  if (editorState.itemID) {
    tasks.push(loadEditorItem(editorState.itemID, { activateMode: false }));
  }
  if (tasks.length) {
    await Promise.all(tasks);
  }
  clearWorkspaceStatus();
}

async function loadEditorItem(itemID, options = {}) {
  if (options.activateMode) {
    setActiveMode("author");
  }
  editorCollabDetailsEl.open = false;
  editorState.itemID = itemID;
  editorItemIDEl.value = itemID;
  if (!itemID) {
    editorMetaEl.textContent = "Select a draft item to load its live draft.";
    editorParticipantsEl.innerHTML = "";
    editorStatusCardsEl.innerHTML = "";
    editorBodyEl.value = "";
    stopLiveTransport();
    return;
  }
  editorState.dirty = false;
  await pullLiveState(itemID, true);
  startLiveTransport();
}

async function pullLiveState(itemID, replaceBody) {
  const state = await getJSON(`/api/items/${itemID}/live`);
  if (replaceBody) {
    editorState.dirty = false;
  }
  if (replaceBody || !editorState.dirty) {
    editorBodyEl.value = state.body;
  }
  renderEditorState(state);
}

function scheduleLivePush() {
  if (!editorState.itemID) {
    return;
  }
  clearTimeout(editorState.pushTimer);
  editorState.pushTimer = setTimeout(() => {
    flushLivePush().catch(handleError);
  }, 300);
}

async function flushLivePush() {
  if (!editorState.itemID || editorState.pushing) {
    return;
  }
  if (sendLiveSocketUpdate(true)) {
    return;
  }
  if (editorState.bridgeRequestID) {
    throw new Error("Direct browser live transport is still opening.");
  }
  throw new Error(browserEmbodimentUnavailableMessage());
}

function startPollLoop() {
  clearPollLoop();
}

function clearPollLoop() {
  if (editorState.pollTimer) {
    clearInterval(editorState.pollTimer);
    editorState.pollTimer = 0;
  }
}

function startHeartbeatLoop() {
  clearHeartbeatLoop();
  editorState.heartbeatTimer = setInterval(() => {
    if (!editorState.itemID) {
      return;
    }
    if (editorState.socketConnected) {
      sendLiveSocketUpdate(false);
    }
  }, 5000);
}

function clearHeartbeatLoop() {
  if (editorState.heartbeatTimer) {
    clearInterval(editorState.heartbeatTimer);
    editorState.heartbeatTimer = 0;
  }
}

function clearReconnectLoop() {
  if (editorState.reconnectTimer) {
    clearTimeout(editorState.reconnectTimer);
    editorState.reconnectTimer = 0;
  }
}

function stopLiveTransport() {
  clearPollLoop();
  clearHeartbeatLoop();
  clearReconnectLoop();
  if (editorState.bridgeRequestID) {
    postBridgeMessage({
      kind: "live-close",
      request_id: editorState.bridgeRequestID,
      socket_path: browserBridgeState.socketPath,
      request: { type: "live-close" },
    });
  }
  editorState.socketConnected = false;
  editorState.socket = null;
  editorState.bridgeRequestID = "";
  editorState.transport = "native-messaging";
}

function startLiveTransport() {
  stopLiveTransport();
  startHeartbeatLoop();
  if (!browserBridgeState.ready || !editorState.itemID) {
    return;
  }
  // Intent: Keep the browser on one direct local contract family for both
  // request/response and live drafting instead of splitting reads/writes onto
  // native messaging while leaving collaboration on the older HTTP adapter.
  // Source: DI-punek
  connectLiveSocket(editorState.itemID);
}

function connectLiveSocket(itemID) {
  const requestID = nextBridgeRequestID("bridge-live");
  editorState.socketGeneration += 1;
  editorState.bridgeRequestID = requestID;
  editorState.socketConnected = false;
  editorState.transport = "native-messaging";
  postBridgeMessage({
    kind: "live-open",
    request_id: requestID,
    socket_path: browserBridgeState.socketPath,
    request: {
      type: "live-open",
      item_id: itemID,
      participant_id: participantID,
      display_name: editorDisplayNameEl.value,
      color: editorColorEl.value,
      cursor: editorBodyEl.selectionStart || 0,
      head: editorBodyEl.selectionEnd || 0,
      typing: false,
    },
  });
}

function scheduleLiveReconnect(itemID, generation) {
  clearReconnectLoop();
  editorState.reconnectTimer = setTimeout(() => {
    if (editorState.itemID !== itemID || generation !== editorState.socketGeneration) {
      return;
    }
    connectLiveSocket(itemID);
  }, 5000);
}

function handleLiveBridgeResponse(payload) {
  if (payload.type === "error") {
    showToast(payload.message || "Live browser bridge error");
    return;
  }
  if (!payload.state) {
    return;
  }
  // Intent: Treat the first real runtime live reply as the browser's
  // acknowledgement boundary so connected state reflects confirmed direct
  // contract truth instead of an optimistic local send. Source: DI-talik
  if (payload.type === "live-state" || payload.type === "live-conflict") {
    editorState.socketConnected = true;
  }
  if (payload.type === "live-conflict") {
    editorState.dirty = false;
    editorBodyEl.value = payload.state.body;
    renderEditorState(payload.state);
    editorMetaEl.textContent = `${payload.state.title} · status ${payload.state.status} · live v${payload.state.version} · current revision ${payload.state.current_revision} · remote changes replaced your stale base version`;
    showToast("Live draft conflict resolved by reloading the shared body");
    return;
  }
  if (editorBodyEl.value === payload.state.body) {
    editorState.dirty = false;
  }
  renderEditorState(payload.state);
}

function sendLiveSocketUpdate(updateBody) {
  if (!editorState.socketConnected || !editorState.bridgeRequestID) {
    return false;
  }
  postBridgeMessage({
    kind: "live-update",
    request_id: editorState.bridgeRequestID,
    socket_path: browserBridgeState.socketPath,
    request: {
      type: "live-update",
      item_id: editorState.itemID,
      participant_id: participantID,
      display_name: editorDisplayNameEl.value,
      color: editorColorEl.value,
      cursor: editorBodyEl.selectionStart || 0,
      head: editorBodyEl.selectionEnd || 0,
      typing: !!updateBody,
      base_version: editorState.version,
      update_body: updateBody,
      body: updateBody ? editorBodyEl.value : "",
    },
  });
  return true;
}

async function sendLivePresencePOST(typing) {
  if (!sendLiveSocketUpdate(typing)) {
    if (editorState.bridgeRequestID) {
      throw new Error("Direct browser live transport is still opening.");
    }
    throw new Error(browserEmbodimentUnavailableMessage());
  }
}

function createMemoryStorage() {
  const data = new Map();
  return {
    getItem(key) {
      return data.has(key) ? data.get(key) : null;
    },
    setItem(key, value) {
      data.set(key, String(value));
    },
    removeItem(key) {
      data.delete(key);
    },
  };
}

// Intent: Keep ex5 browser startup alive in private or policy-restricted
// environments where localStorage access or UUID helpers may fail, so the live
// draft embodiment still boots and can participate with an ephemeral identity.
// Source: DI-mitob
function safeParticipantStorage() {
  const memory = createMemoryStorage();
  return {
    getItem(key) {
      try {
        return window.localStorage ? window.localStorage.getItem(key) : memory.getItem(key);
      } catch {
        return memory.getItem(key);
      }
    },
    setItem(key, value) {
      try {
        if (window.localStorage) {
          window.localStorage.setItem(key, value);
          return;
        }
      } catch {
        // Intent: Fall back to in-memory participant identity when browser
        // storage is blocked so the UI still starts and joins the shared draft.
        // Source: DI-mitob
      }
      memory.setItem(key, value);
    },
  };
}

function fallbackParticipantSuffix() {
  return `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 10)}`;
}

function createParticipantID() {
  if (window.crypto && typeof window.crypto.randomUUID === "function") {
    return `browser-${window.crypto.randomUUID()}`;
  }
  return `browser-${fallbackParticipantSuffix()}`;
}

function getParticipantID() {
  const storageKey = "oks.participant_id";
  const storage = safeParticipantStorage();
  const existing = storage.getItem(storageKey);
  if (existing) {
    return existing;
  }
  const created = createParticipantID();
  storage.setItem(storageKey, created);
  return created;
}

function showToast(message) {
  toastEl.hidden = false;
  toastEl.textContent = message;
  clearTimeout(showToast.timer);
  showToast.timer = setTimeout(() => {
    toastEl.hidden = true;
  }, 1800);
}

// Intent: Prefer the calmer helper-select path for advanced search filters
// while preserving exact manual ID entry whenever the chosen value is not in
// the current catalog. Source: DI-zovuk; DI-ralek
function syncSearchSelectAndOverride(selectEl, inputEl, value) {
  const normalized = (value || "").trim();
  if (!selectEl || !inputEl) {
    return;
  }
  const match = Array.from(selectEl.options).some((option) => option.value === normalized);
  if (match) {
    selectEl.value = normalized;
    inputEl.value = "";
    return;
  }
  selectEl.value = "";
  inputEl.value = normalized;
}

function handleError(error, context = "Browser") {
  showToast(error.message);
  setWorkspaceStatus(`${context} failed`, "error", error.message);
}

clearSearch();
setActiveMode("review");
setReviewLane("drafts");
openOperateStage("run");
initializeBrowserEmbodiment()
  .then((ready) => {
    if (!ready) {
      return;
    }
    return refresh();
  })
  .catch(handleError);
// Intent: Keep search centered on common review tasks first, then let operators
// drop into structured filters only when they need finer drilldown. Source:
// DI-rovak
async function runPresetSearch(filters, message) {
  setActiveMode("review");
  await runSearch(filters);
  searchResultsEl.scrollIntoView({ behavior: "smooth", block: "start" });
  setWorkspaceStatus(message, "info");
}
