const statsEl = document.getElementById("stats");
const responsibilityListEl = document.getElementById("responsibility-list");
const itemListEl = document.getElementById("item-list");
const runListEl = document.getElementById("run-list");
const searchResultsEl = document.getElementById("search-results");
const toastEl = document.getElementById("toast");

// Intent: Keep the browser as a first-class embodiment over the same local
// operational runtime as the CLI, with local form state layered over durable
// event-backed records served by the HTTP adapter. Source: DI-radok; DI-zuvob
document.getElementById("responsibility-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const payload = {
    actor: form.actor.value,
    title: form.title.value,
    summary: form.summary.value,
    role_keys: splitCSV(form.role_keys.value),
    tags: splitCSV(form.tags.value),
  };
  await postJSON("/api/responsibilities", payload);
  form.reset();
  form.actor.value = "alice";
  await refresh();
});

document.getElementById("item-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const payload = {
    actor: form.actor.value,
    kind: form.kind.value,
    title: form.title.value,
    summary: form.summary.value,
    body: form.body.value,
    tags: splitCSV(form.tags.value),
    responsibility_ids: splitCSV(form.responsibility_ids.value),
  };
  await postJSON("/api/items", payload);
  form.reset();
  form.actor.value = "alice";
  form.kind.value = "procedure";
  await refresh();
});

document.getElementById("run-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const payload = {
    actor: form.actor.value,
    kind: form.kind.value,
    item_id: form.item_id.value,
    revision: Number(form.revision.value || 1),
    outcome: form.outcome.value,
    notes: form.notes.value,
    machine: form.machine.value,
    location: form.location.value,
    responsibility_ids: splitCSV(form.responsibility_ids.value),
  };
  await postJSON("/api/runs", payload);
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
  const payload = {
    actor: form.actor.value,
    revision: Number(form.revision.value || 0),
    role: form.role.value,
    decision: form.decision.value,
    notes: form.notes.value,
  };
  const base = form.target_type.value === "run" ? `/api/runs/${targetID}/approvals` : `/api/items/${targetID}/approvals`;
  await postJSON(base, payload);
  form.reset();
  form.actor.value = "boss";
  form.target_type.value = "knowledge_item";
  form.decision.value = "approved";
  form.revision.value = "0";
  await refresh();
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
  await refresh();
});

document.getElementById("search-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const query = event.currentTarget.q.value;
  const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
  const payload = await response.json();
  searchResultsEl.textContent = JSON.stringify(payload, null, 2);
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
  if (!response.ok) {
    throw new Error(await response.text());
  }
  showToast(`Saved via ${path}`);
}

function renderStats(data) {
  statsEl.innerHTML = "";
  const fields = [
    ["Responsibilities", data.responsibilities],
    ["Procedures", data.procedures],
    ["Training", data.training_items],
    ["Maintenance", data.maintenance_items],
    ["Runs", data.procedure_runs + data.training_runs + data.maintenance_runs],
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

function renderResponsibilities(items) {
  responsibilityListEl.innerHTML = "";
  for (const item of items) {
    const card = document.createElement("div");
    card.className = "card";
    card.innerHTML = `<h3>${item.id} · ${item.title}</h3><div class="meta">${item.summary || ""}\nroles: ${(item.linked_role_keys || []).join(", ") || "-"}</div>`;
    responsibilityListEl.appendChild(card);
  }
}

function renderKnowledgeItems(items) {
  itemListEl.innerHTML = "";
  for (const item of items) {
    const card = document.createElement("div");
    card.className = "card";
    card.innerHTML = `<div class="kind">${item.kind}</div><h3>${item.id} · ${item.title}</h3><div class="meta">revision ${item.current_revision}\n${item.summary || ""}</div>`;
    itemListEl.appendChild(card);
  }
}

function renderRuns(items) {
  runListEl.innerHTML = "";
  for (const item of items) {
    const card = document.createElement("div");
    card.className = "card";
    card.innerHTML = `<div class="kind">${item.kind} run</div><h3>${item.id} · ${item.item_id}</h3><div class="meta">revision ${item.revision}\noutcome: ${item.outcome || "-"}\n${item.notes || ""}</div>`;
    runListEl.appendChild(card);
  }
}

function showToast(message) {
  toastEl.hidden = false;
  toastEl.textContent = message;
  clearTimeout(showToast.timer);
  showToast.timer = setTimeout(() => {
    toastEl.hidden = true;
  }, 1800);
}

async function refresh() {
  const [dashboard, responsibilities, items, runs] = await Promise.all([
    fetch("/api/dashboard").then((response) => response.json()),
    fetch("/api/responsibilities").then((response) => response.json()),
    fetch("/api/items").then((response) => response.json()),
    fetch("/api/runs").then((response) => response.json()),
  ]);
  renderStats(dashboard);
  renderResponsibilities(responsibilities.responsibilities || []);
  renderKnowledgeItems(items.items || []);
  renderRuns(runs.runs || []);
}

refresh().catch((error) => {
  showToast(error.message);
  searchResultsEl.textContent = error.stack || error.message;
});
