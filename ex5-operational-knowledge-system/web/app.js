const statsEl = document.getElementById("stats");
const placeListEl = document.getElementById("place-list");
const resourceListEl = document.getElementById("resource-list");
const responsibilityListEl = document.getElementById("responsibility-list");
const itemListEl = document.getElementById("item-list");
const runListEl = document.getElementById("run-list");
const searchResultsEl = document.getElementById("search-results");
const searchActiveEl = document.getElementById("search-active");
const searchRawEl = document.getElementById("search-raw");
const toastEl = document.getElementById("toast");
const detailMetaEl = document.getElementById("detail-meta");
const detailSummaryEl = document.getElementById("detail-summary");
const detailActionsEl = document.getElementById("detail-actions");
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
};

// Intent: Keep the browser as an equal operational embodiment while making
// knowledge-item drafting collaborative in the browser without collapsing the
// durable revision and approval workflow into ephemeral UI state. Source:
// DI-lusov; DI-zoruk
document.getElementById("place-form").addEventListener("submit", async (event) => {
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
});

document.getElementById("resource-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  await postJSON("/api/resources", {
    actor: form.actor.value,
    kind: form.kind.value,
    name: form.name.value,
    summary: form.summary.value,
    place_id: form.place_id.value,
    tags: splitCSV(form.tags.value),
  });
  form.reset();
  form.actor.value = "alice";
  await refresh();
});

document.getElementById("responsibility-form").addEventListener("submit", async (event) => {
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
});

document.getElementById("item-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const item = await postJSON("/api/items", {
    actor: form.actor.value,
    kind: form.kind.value,
    title: form.title.value,
    summary: form.summary.value,
    body: form.body.value,
    tags: splitCSV(form.tags.value),
    responsibility_ids: splitCSV(form.responsibility_ids.value),
  });
  form.reset();
  form.actor.value = "alice";
  form.kind.value = "procedure";
  await refresh(item.id);
});

document.getElementById("run-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  await postJSON("/api/runs", {
    actor: form.actor.value,
    kind: form.kind.value,
    item_id: form.item_id.value,
    revision: Number(form.revision.value || 1),
    outcome: form.outcome.value,
    notes: form.notes.value,
    machine: form.machine.value,
    location: form.location.value,
    place_id: form.place_id.value,
    resource_ids: splitCSV(form.resource_ids.value),
    responsibility_ids: splitCSV(form.responsibility_ids.value),
  });
  form.reset();
  form.actor.value = "bob";
  form.kind.value = "procedure";
  form.revision.value = "1";
  await refresh();
});

document.getElementById("approval-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const targetID = form.target_id.value;
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
  await refresh(editorState.itemID);
});

document.getElementById("evidence-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const body = new FormData();
  body.set("actor", form.actor.value);
  body.set("summary", form.summary.value);
  body.set("facts_json", form.facts_json.value || "{}");
  if (form.attachment.files[0]) {
    body.set("attachment", form.attachment.files[0]);
  }
  const response = await fetch(`/api/runs/${form.run_id.value}/evidence`, { method: "POST", body });
  if (!response.ok) {
    throw new Error(await response.text());
  }
  form.reset();
  form.actor.value = "bob";
  showToast("Evidence attached");
  await refresh();
});

document.getElementById("search-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const filters = getSearchFilters(event.currentTarget);
  const response = await fetch(`/api/search?${buildSearchParams(filters).toString()}`);
  const payload = await response.json();
  renderSearchResults(filters, payload);
});

editorItemIDEl.addEventListener("change", async () => {
  await loadEditorItem(editorItemIDEl.value);
});

editorBodyEl.addEventListener("input", () => {
  editorState.dirty = true;
  scheduleLivePush();
});

editorBodyEl.addEventListener("click", scheduleLivePush);
editorBodyEl.addEventListener("keyup", scheduleLivePush);

editorRefreshEl.addEventListener("click", async () => {
  if (!editorState.itemID) {
    return;
  }
  await pullLiveState(editorState.itemID, true);
});

editorSnapshotEl.addEventListener("click", async () => {
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
});

editorApproveEl.addEventListener("click", async () => {
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
});

editorSupersedeEl.addEventListener("click", async () => {
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
});

function splitCSV(input) {
  return input.split(",").map((value) => value.trim()).filter(Boolean);
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
    ["Inventory", data.inventory_items],
    ["Runs", data.procedure_runs + data.training_runs + data.maintenance_runs + data.inventory_runs],
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

// Intent: Let browser operators combine structured filters with free-text
// search so they can drill into one slice of the operational graph without
// manually cross-referencing IDs outside the app. Source: DI-honus
function getSearchFilters(form) {
  return {
    q: form.q.value.trim(),
    kind: form.kind.value.trim(),
    status: form.status.value.trim(),
    place_id: form.place_id.value.trim(),
    resource_id: form.resource_id.value.trim(),
    responsibility_id: form.responsibility_id.value.trim(),
  };
}

function buildSearchParams(filters) {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(filters)) {
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
      const card = document.createElement("button");
      card.type = "button";
      card.className = "card selectable-card";
      card.innerHTML = `<div class="kind">${type}</div><h3>${item.id}</h3><div class="meta">${searchSummary(type, item)}</div>`;
      card.addEventListener("click", () => {
        inspectRecord(type, item.id).catch(handleError);
      });
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
  if (filters.place_id) {
    labels.push(`place: ${filters.place_id}`);
  }
  if (filters.resource_id) {
    labels.push(`resource: ${filters.resource_id}`);
  }
  if (filters.responsibility_id) {
    labels.push(`responsibility: ${filters.responsibility_id}`);
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
  detailMetaEl.textContent = `Loading ${type} ${id}...`;
  detailSummaryEl.innerHTML = "";
  detailActionsEl.innerHTML = "";
  detailTimelineEl.innerHTML = "";
  const record = await getJSON(detailPath(type, id));
  detailMetaEl.textContent = detailSummary(type, record);
  renderDetailSummary(type, record);
  renderDetailTimeline(record.timeline || []);
  detailJSONEl.textContent = JSON.stringify(record, null, 2);
  renderDetailActions(type, record);
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
        ["Events", (record.timeline || []).length],
      ];
    case "resource":
      return [
        ["Place", record.place_id || "-"],
        ["Tags", (record.tags || []).length],
        ["Links", (record.links || []).length],
        ["Events", (record.timeline || []).length],
      ];
    case "responsibility":
      return [
        ["Team", record.team || "-"],
        ["Items", (record.linked_item_ids || []).length],
        ["Runs", (record.linked_run_ids || []).length],
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

function renderDetailActions(type, record) {
  detailActionsEl.innerHTML = "";
  const links = [];
  if (type === "resource" && record.place_id) {
    links.push(["Open place", "place", record.place_id]);
  }
  if (type === "item") {
    links.push(["Open live draft", "item", record.id]);
    for (const responsibilityID of record.responsibility_ids || []) {
      links.push([`Responsibility ${responsibilityID}`, "responsibility", responsibilityID]);
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
    for (const resourceID of record.resource_ids || []) {
      links.push([`Resource ${resourceID}`, "resource", resourceID]);
    }
    for (const childID of record.child_place_ids || []) {
      links.push([`Child place ${childID}`, "place", childID]);
    }
  }
  if (type === "responsibility") {
    for (const itemID of record.linked_item_ids || []) {
      links.push([`Item ${itemID}`, "item", itemID]);
    }
    for (const runID of record.linked_run_ids || []) {
      links.push([`Run ${runID}`, "run", runID]);
    }
  }
  for (const [label, nextType, nextID] of links) {
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = label;
    button.addEventListener("click", () => {
      inspectRecord(nextType, nextID).catch(handleError);
    });
    detailActionsEl.appendChild(button);
  }
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
  const [dashboard, places, resources, responsibilities, items, runs] = await Promise.all([
    getJSON("/api/dashboard"),
    getJSON("/api/places"),
    getJSON("/api/resources"),
    getJSON("/api/responsibilities"),
    getJSON("/api/items"),
    getJSON("/api/runs"),
  ]);
  renderStats(dashboard);
  renderPlaces(places.places || []);
  renderResources(resources.resources || []);
  renderResponsibilities(responsibilities.responsibilities || []);
  renderKnowledgeItems(items.items || []);
  renderRuns(runs.runs || []);
  if (!selectedItemID && (items.items || []).length > 0) {
    selectedItemID = items.items[0].id;
  }
  if (selectedItemID) {
    await Promise.all([
      loadEditorItem(selectedItemID),
      inspectRecord("item", selectedItemID),
    ]);
  }
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

function getParticipantID() {
  const storageKey = "oks.participant_id";
  const existing = window.localStorage.getItem(storageKey);
  if (existing) {
    return existing;
  }
  const created = `browser-${crypto.randomUUID()}`;
  window.localStorage.setItem(storageKey, created);
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

function handleError(error) {
  showToast(error.message);
  searchResultsEl.textContent = error.stack || error.message;
}

refresh().catch(handleError);
