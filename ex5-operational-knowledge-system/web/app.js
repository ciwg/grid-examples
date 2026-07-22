const statsEl = document.getElementById("stats");
const placeListEl = document.getElementById("place-list");
const resourceListEl = document.getElementById("resource-list");
const responsibilityListEl = document.getElementById("responsibility-list");
const itemListEl = document.getElementById("item-list");
const runListEl = document.getElementById("run-list");
const problemReviewEl = document.getElementById("problem-review");
const searchResultsEl = document.getElementById("search-results");
const searchActiveEl = document.getElementById("search-active");
const searchRawEl = document.getElementById("search-raw");
const searchDebugEl = document.getElementById("search-debug");
const searchClearEl = document.getElementById("search-clear");
const toastEl = document.getElementById("toast");
const workspaceStatusEl = document.getElementById("workspace-status");
const detailMetaEl = document.getElementById("detail-meta");
const detailSummaryEl = document.getElementById("detail-summary");
const detailActionsEl = document.getElementById("detail-actions");
const detailReviewEl = document.getElementById("detail-review");
const detailTimelineEl = document.getElementById("detail-timeline");
const detailJSONEl = document.getElementById("detail-json");

const editorItemIDEl = document.getElementById("editor-item-id");
const editorActorEl = document.getElementById("editor-actor");
const editorDisplayNameEl = document.getElementById("editor-display-name");
const editorColorEl = document.getElementById("editor-color");
const editorMetaEl = document.getElementById("editor-meta");
const editorParticipantsEl = document.getElementById("editor-participants");
const editorBodyEl = document.getElementById("editor-body");
const editorRefreshEl = document.getElementById("editor-refresh");
const editorSnapshotEl = document.getElementById("editor-snapshot");
const editorApproveEl = document.getElementById("editor-approve");
const editorSupersedeEl = document.getElementById("editor-supersede");
const approvalFormEl = document.getElementById("approval-form");
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

const participantID = getParticipantID();
const editorState = {
  itemID: "",
  version: 0,
  title: "",
  status: "",
  currentRevision: 0,
  dirty: false,
  pushing: false,
  lastRenderedBody: "",
  pushTimer: 0,
  pollTimer: 0,
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

function runHandled(action, context) {
  return (...args) => {
    clearWorkspaceStatus();
    Promise.resolve(action(...args)).catch((error) => handleError(error, context));
  };
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
  const body = new FormData();
  body.set("actor", form.actor.value);
  body.set("summary", form.summary.value);
  body.set("facts_json", form.facts_json.value || "{}");
  if (form.attachment.files[0]) {
    body.set("attachment", form.attachment.files[0]);
  }
  const response = await fetch(`/api/runs/${runID}/evidence`, { method: "POST", body });
  if (!response.ok) {
    throw new Error(await response.text());
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

editorItemIDEl.addEventListener("change", runHandled(async () => {
  await loadEditorItem(editorItemIDEl.value);
}, "Live Draft Studio"));

editorBodyEl.addEventListener("input", () => {
  editorState.dirty = true;
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

async function postJSON(path, payload) {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  const text = await response.text();
  if (!response.ok) {
    throw new Error(text);
  }
  showToast(`Saved via ${path}`);
  return text ? JSON.parse(text) : null;
}

async function getJSON(path) {
  const response = await fetch(path);
  if (!response.ok) {
    throw new Error(await response.text());
  }
  return response.json();
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
      Promise.all([
        loadEditorItem(item.id),
        inspectRecord("item", item.id),
      ]).catch(handleError);
    });
    itemListEl.appendChild(card);
  }
  if (editorState.itemID) {
    editorItemIDEl.value = editorState.itemID;
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
  renderSelect(resourcePlaceSelectEl, catalogState.places, "Select a place", (place) => `${place.id} · ${place.name}`);
  renderSelect(itemResponsibilitySelectEl, catalogState.responsibilities, "Optional responsibility", (responsibility) => `${responsibility.id} · ${responsibility.title}`);
  renderSelect(runItemSelectEl, catalogState.items, "Select a knowledge item", (item) => `${item.id} · ${item.title}`);
  renderSelect(runPlaceSelectEl, catalogState.places, "Optional place", (place) => `${place.id} · ${place.name}`);
  renderSelect(runResourceSelectEl, catalogState.resources, "Optional primary resource", (resource) => `${resource.id} · ${resource.name}`);
  renderSelect(runResponsibilitySelectEl, catalogState.responsibilities, "Optional primary responsibility", (responsibility) => `${responsibility.id} · ${responsibility.title}`);
  renderSelect(evidenceRunSelectEl, catalogState.runs, "Select a run", (run) => `${run.id} · ${run.item_id} · ${run.outcome || "-"}`);
  renderApprovalTargetOptions();
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
  renderSelect(approvalTargetSelectEl, items, `Select a ${targetType === "run" ? "run" : "knowledge item"}`, formatter);
  syncApprovalDefaults();
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
    : `Choose a ${form.target_type.value === "run" ? "run" : "knowledge item"} to load matching records.`;
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
    place_id: form.place_id.value.trim(),
    resource_id: form.resource_id.value.trim(),
    responsibility_id: form.responsibility_id.value.trim(),
    problem: form.dataset.problem === "true",
  };
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
      card.innerHTML = `<div class="search-result-head"><div><div class="kind">${type}</div><h3>${item.id}</h3></div></div><div class="meta">${searchSummary(type, item)}</div>`;
      const actions = document.createElement("div");
      actions.className = "card-actions";
      actions.appendChild(makeActionButton("Inspect", () => inspectRecord(type, item.id), "Search"));
      if (type === "item") {
        actions.appendChild(makeActionButton("Open draft", () => Promise.all([loadEditorItem(item.id), inspectRecord("item", item.id)]), "Search"));
      }
      if (type === "run") {
        actions.appendChild(makeActionButton("Item", () => inspectRecord("item", item.item_id), "Search"));
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

function formatSearchFilters(filters) {
  const labels = [];
  if (filters.query || filters.q) {
    labels.push(`query: ${filters.query || filters.q}`);
  }
  if (filters.kind) {
    labels.push(`kind: ${filters.kind}`);
  }
  if (filters.status) {
    labels.push(`status: ${filters.status}`);
  }
  if (filters.outcome) {
    labels.push(`outcome: ${filters.outcome}`);
  }
  if (filters.place_id) {
    labels.push(`place: ${filters.place_id}`);
  }
  if (filters.resource_id) {
    labels.push(`resource: ${filters.resource_id}`);
  }
  if (filters.responsibility_id) {
    labels.push(`responsibility: ${filters.responsibility_id}`);
  }
  if (filters.problem) {
    labels.push("problems only");
  }
  if (labels.length === 0) {
    return "No active search filters.";
  }
  return `Active filters: ${labels.join(" · ")}`;
}

function searchSummary(type, item) {
  switch (type) {
    case "place":
      return `${item.name || ""}\n${item.summary || ""}`;
    case "resource":
      return `${item.name || ""}\nplace: ${item.place_id || "-"}`;
    case "responsibility":
      return `${item.title || ""}\n${item.summary || ""}`;
    case "item":
      return `${item.title || ""}\n${item.kind || ""} · ${item.status || ""}`;
    case "run":
      return `${item.item_id || ""}\nrevision ${item.revision || 0} · ${item.outcome || "-"}`;
    default:
      return "";
  }
}

// Intent: Let operators inspect and traverse the current operational graph in
// the browser without manually copying IDs between separate lists, while
// keeping the existing local HTTP runtime and record model unchanged. Source:
// DI-vopuk
async function inspectRecord(type, id) {
  detailState.type = type;
  detailState.id = id;
  detailState.record = null;
  detailMetaEl.textContent = `Loading ${type} ${id}...`;
  detailSummaryEl.innerHTML = "";
  detailActionsEl.innerHTML = "";
  detailReviewEl.innerHTML = "";
  detailTimelineEl.innerHTML = "";
  const record = await getJSON(detailPath(type, id));
  detailState.record = record;
  detailMetaEl.textContent = detailSummary(type, record);
  renderDetailSummary(type, record);
  renderDetailReview(type, record);
  renderDetailTimeline(record.timeline || []);
  detailJSONEl.textContent = JSON.stringify(record, null, 2);
  renderDetailActions(type, record);
  applyContextDefaults(type, record);
  if (type === "item") {
    await loadEditorItem(id);
  }
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
      return `${record.id} · ${record.kind} · ${record.name}`;
    case "resource":
      return `${record.id} · ${record.kind} · ${record.name}`;
    case "responsibility":
      return `${record.id} · responsibility · ${record.title}`;
    case "item":
      return `${record.id} · ${record.kind} · ${record.title} · ${record.status}`;
    case "run":
      return `${record.id} · ${record.kind} run · ${record.item_id}`;
    default:
      return `${type} ${record.id || ""}`;
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
  if (type === "resource" && record.place_id) {
    links.push(["Open place", "place", record.place_id]);
    links.push(["Search receiving here", "search", { kind: "receiving_check", resource_id: record.id }]);
    links.push(["Search counts here", "search", { kind: "inventory_audit", resource_id: record.id }]);
    links.push(["Search problems here", "search", { resource_id: record.id, problem: true }]);
    for (const run of record.related_runs || []) {
      links.push([`Run ${run.id}`, "run", run.id]);
    }
  }
  if (type === "item") {
    links.push(["Open live draft", "item", record.id]);
    for (const responsibilityID of record.responsibility_ids || []) {
      links.push([`Responsibility ${responsibilityID}`, "responsibility", responsibilityID]);
    }
    for (const run of record.related_runs || []) {
      links.push([`Run ${run.id}`, "run", run.id]);
    }
  }
  if (type === "run") {
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
}

// Intent: Reuse the structured search form as the single drilldown path so
// direct inspector actions and manual operator searches stay behaviorally
// identical, including problem-only review drilldowns. Source: DI-vafuk;
// DI-vemur
async function runSearch(filters) {
  clearWorkspaceStatus();
  const form = document.getElementById("search-form");
  form.q.value = filters.q || "";
  form.kind.value = filters.kind || "";
  form.status.value = filters.status || "";
  form.outcome.value = filters.outcome || "";
  form.place_id.value = filters.place_id || "";
  form.resource_id.value = filters.resource_id || "";
  form.responsibility_id.value = filters.responsibility_id || "";
  form.dataset.problem = filters.problem ? "true" : "false";
  const effectiveFilters = getSearchFilters(form);
  const response = await fetch(`/api/search?${buildSearchParams(effectiveFilters).toString()}`);
  if (!response.ok) {
    throw new Error(await response.text());
  }
  const payload = await response.json();
  renderSearchResults(effectiveFilters, payload);
  clearWorkspaceStatus();
}

function clearSearch() {
  const form = document.getElementById("search-form");
  form.reset();
  form.dataset.problem = "false";
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
  editorState.lastRenderedBody = state.body;
  editorMetaEl.textContent = `${state.title} · status ${state.status} · live v${state.version} · current revision ${state.current_revision}${editorState.dirty ? " · local edits pending" : ""}`;
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
  renderRuns(runs.runs || []);
  refreshActionCatalog();
  if (!selectedItemID && (items.items || []).length > 0) {
    selectedItemID = items.items[0].id;
  }
  if (selectedItemID) {
    await Promise.all([
      loadEditorItem(selectedItemID),
      inspectRecord("item", selectedItemID),
    ]);
  }
  clearWorkspaceStatus();
}

async function loadEditorItem(itemID) {
  editorState.itemID = itemID;
  editorItemIDEl.value = itemID;
  if (!itemID) {
    editorMetaEl.textContent = "Select a knowledge item to load its live draft.";
    editorParticipantsEl.innerHTML = "";
    editorBodyEl.value = "";
    clearPollLoop();
    return;
  }
  editorState.dirty = false;
  await pullLiveState(itemID, true);
  startPollLoop();
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
  editorState.pushing = true;
  try {
    const response = await fetch(`/api/items/${editorState.itemID}/live`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        participant_id: participantID,
        display_name: editorDisplayNameEl.value,
        color: editorColorEl.value,
        cursor: editorBodyEl.selectionStart || 0,
        head: editorBodyEl.selectionEnd || 0,
        typing: true,
        base_version: editorState.version,
        update_body: true,
        body: editorBodyEl.value,
      }),
    });
    if (response.status === 409) {
      const payload = await response.json();
      editorState.dirty = false;
      editorBodyEl.value = payload.state.body;
      renderEditorState(payload.state);
      editorMetaEl.textContent = `${payload.state.title} · status ${payload.state.status} · live v${payload.state.version} · current revision ${payload.state.current_revision} · remote changes replaced your stale base version`;
      showToast("Live draft conflict resolved by reloading the shared body");
      return;
    }
    if (!response.ok) {
      throw new Error(await response.text());
    }
    const state = await response.json();
    editorState.dirty = false;
    renderEditorState(state);
  } finally {
    editorState.pushing = false;
  }
}

function startPollLoop() {
  clearPollLoop();
  editorState.pollTimer = setInterval(() => {
    if (!editorState.itemID || editorState.pushing) {
      return;
    }
    pullLiveState(editorState.itemID, false).catch(handleError);
  }, 2000);
}

function clearPollLoop() {
  if (editorState.pollTimer) {
    clearInterval(editorState.pollTimer);
    editorState.pollTimer = 0;
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

function handleError(error, context = "Browser") {
  showToast(error.message);
  setWorkspaceStatus(`${context} failed`, "error", error.message);
}

clearSearch();
refresh().catch(handleError);
