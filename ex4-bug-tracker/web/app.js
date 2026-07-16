const state = {
  meta: null,
  issues: [],
  currentIssue: null,
  currentIssueID: new URLSearchParams(window.location.search).get("issue") || "",
  user: localStorage.getItem("ex4-bug-user") || "reporter",
};

const userSelect = document.getElementById("user-select");
const newIssueButton = document.getElementById("new-issue-button");
const newIssuePanel = document.getElementById("new-issue-panel");
const newIssueForm = document.getElementById("new-issue-form");
const cancelNewIssueButton = document.getElementById("cancel-new-issue");
const issueList = document.getElementById("issue-list");
const filterStatus = document.getElementById("filter-status");
const filterAssignee = document.getElementById("filter-assignee");
const issueSeverity = document.getElementById("issue-severity");
const issueDetailPanel = document.getElementById("issue-detail-panel");
const emptyState = document.getElementById("empty-state");
const messageBanner = document.getElementById("message-banner");
const issueTitleDisplay = document.getElementById("issue-title-display");
const issueIDDisplay = document.getElementById("issue-id");
const issueStatusBadge = document.getElementById("issue-status-badge");
const issueSeverityDisplay = document.getElementById("issue-severity-display");
const issueReporterDisplay = document.getElementById("issue-reporter-display");
const issueAssigneeDisplay = document.getElementById("issue-assignee-display");
const issueUpdatedDisplay = document.getElementById("issue-updated-display");
const issueDescriptionDisplay = document.getElementById("issue-description-display");
const assignmentForm = document.getElementById("assignment-form");
const assignmentSelect = document.getElementById("assignment-select");
const statusForm = document.getElementById("status-form");
const statusSelect = document.getElementById("status-select");
const timelineList = document.getElementById("timeline-list");
const commentForm = document.getElementById("comment-form");
const commentBody = document.getElementById("comment-body");
const attachmentForm = document.getElementById("attachment-form");
const attachmentFile = document.getElementById("attachment-file");

function currentIdentity() {
  return state.meta.identities.find((identity) => identity.id === state.user);
}

function currentRole() {
  return currentIdentity()?.role || "";
}

function apiFetch(path, options = {}) {
  const headers = new Headers(options.headers || {});
  headers.set("X-Bug-User", state.user);
  if (options.body && !(options.body instanceof FormData) && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  return fetch(path, { ...options, headers });
}

function setMessage(kind, text) {
  if (!text) {
    messageBanner.hidden = true;
    messageBanner.className = "message-banner";
    messageBanner.textContent = "";
    return;
  }
  messageBanner.hidden = false;
  messageBanner.className = `message-banner ${kind}`;
  messageBanner.textContent = text;
}

function setCurrentIssueID(issueID) {
  state.currentIssueID = issueID || "";
  const url = new URL(window.location.href);
  if (state.currentIssueID) {
    url.searchParams.set("issue", state.currentIssueID);
  } else {
    url.searchParams.delete("issue");
  }
  window.history.replaceState({}, "", url);
}

function openNewIssuePanel() {
  if (currentRole() !== "reporter") {
    return;
  }
  newIssuePanel.hidden = false;
  issueDetailPanel.hidden = true;
  emptyState.hidden = true;
}

function closeNewIssuePanel() {
  newIssuePanel.hidden = true;
  newIssueForm.reset();
  renderDetail();
}

function formatTimestamp(value) {
  const timestamp = new Date(value);
  if (Number.isNaN(timestamp.getTime())) {
    return value;
  }
  return timestamp.toLocaleString();
}

function issueLink(attachmentID) {
  return `/api/issues/${encodeURIComponent(state.currentIssueID)}/attachments/${encodeURIComponent(attachmentID)}?user=${encodeURIComponent(state.user)}`;
}

function describeEvent(event) {
  switch (event.type) {
    case "created":
      return `<strong>${event.actor}</strong> opened the issue at <strong>${event.severity}</strong> severity.`;
    case "commented":
      return `<p>${escapeHTML(event.comment)}</p>`;
    case "assigned":
      if (event.assignee) {
        return `<strong>${event.actor}</strong> assigned the issue to <strong>${event.assignee}</strong>.`;
      }
      return `<strong>${event.actor}</strong> cleared the assignee.`;
    case "status_changed":
      if (event.previous_status === "Resolved" && event.status === "Triaged") {
        return `<strong>${event.actor}</strong> reopened the issue and moved it back to <strong>${event.status}</strong>.`;
      }
      return `<strong>${event.actor}</strong> changed status from <strong>${event.previous_status}</strong> to <strong>${event.status}</strong>.`;
    case "attachment_added":
      return `<strong>${event.actor}</strong> uploaded <a class="attachment-link" href="${issueLink(event.attachment_id)}">${escapeHTML(event.attachment_name)}</a> (${event.attachment_size} bytes).`;
    default:
      return `<em>Unknown event ${escapeHTML(event.type)}</em>`;
  }
}

function escapeHTML(value) {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

function renderIssueList() {
  issueList.innerHTML = "";
  for (const issue of state.issues) {
    const li = document.createElement("li");
    const button = document.createElement("button");
    if (issue.id === state.currentIssueID) {
      button.classList.add("active");
    }
    button.type = "button";
    button.innerHTML = `
      <strong>${issue.id}</strong>
      <div>${escapeHTML(issue.title)}</div>
      <div class="meta">
        <span>${issue.status}</span>
        <span>${issue.severity}</span>
        <span>${issue.assignee || "unassigned"}</span>
      </div>
    `;
    button.addEventListener("click", async () => {
      setCurrentIssueID(issue.id);
      closeNewIssuePanel();
      try {
        await loadIssue(issue.id);
      } catch (error) {
        setMessage("error", error instanceof Error ? error.message : String(error));
      }
    });
    li.appendChild(button);
    issueList.appendChild(li);
  }
}

function renderDetail() {
  const issue = state.currentIssue;
  if (!issue) {
    issueDetailPanel.hidden = true;
    emptyState.hidden = !newIssuePanel.hidden;
    return;
  }
  if (!newIssuePanel.hidden) {
    issueDetailPanel.hidden = true;
    emptyState.hidden = true;
    return;
  }
  issueDetailPanel.hidden = false;
  emptyState.hidden = true;
  issueIDDisplay.textContent = issue.id;
  issueTitleDisplay.textContent = issue.title;
  issueStatusBadge.textContent = issue.status;
  issueSeverityDisplay.textContent = issue.severity;
  issueReporterDisplay.textContent = issue.reporter;
  issueAssigneeDisplay.textContent = issue.assignee || "unassigned";
  issueUpdatedDisplay.textContent = formatTimestamp(issue.updated_at);
  issueDescriptionDisplay.textContent = issue.description;
  renderAssignmentOptions(issue);
  renderStatusOptions(issue);
  timelineList.innerHTML = "";
  for (const event of issue.timeline) {
    const item = document.createElement("li");
    item.innerHTML = `
      <div class="timeline-header">
        <span>${event.actor}</span>
        <span>${formatTimestamp(event.timestamp)}</span>
      </div>
      <div class="timeline-body">${describeEvent(event)}</div>
    `;
    timelineList.appendChild(item);
  }
}

function renderAssignmentOptions(issue) {
  assignmentSelect.innerHTML = "";
  const blank = document.createElement("option");
  blank.value = "";
  blank.textContent = "Unassigned";
  assignmentSelect.appendChild(blank);
  for (const identity of state.meta.identities.filter((identity) => identity.role === "engineer")) {
    const option = document.createElement("option");
    option.value = identity.id;
    option.textContent = identity.id;
    assignmentSelect.appendChild(option);
  }
  assignmentSelect.value = issue.assignee || "";
  assignmentForm.hidden = currentRole() !== "triage";
}

function renderStatusOptions(issue) {
  statusSelect.innerHTML = "";
  const allowedStatuses = allowedStatusTargets(issue);
  for (const status of allowedStatuses) {
    const option = document.createElement("option");
    option.value = status;
    option.textContent = status;
    statusSelect.appendChild(option);
  }
  if (allowedStatuses.length === 0) {
    statusForm.hidden = true;
    return;
  }
  statusForm.hidden = false;
  statusSelect.value = allowedStatuses[0];
}

// Intent: Mirror the server's fixed workflow in the browser so users only see
// actions that their current role can actually perform, instead of learning by
// avoidable server-side rejection banners. Source: DI-zumog
function allowedStatusTargets(issue) {
  const role = currentRole();
  if (!issue) {
    return [];
  }
  if (role === "triage" && issue.status === "New") {
    return ["Triaged"];
  }
  if (role === "engineer" && issue.status === "Triaged" && issue.assignee === state.user) {
    return ["In Progress"];
  }
  if (role === "engineer" && issue.status === "In Progress" && issue.assignee === state.user) {
    return ["Resolved"];
  }
  if ((role === "reporter" || role === "triage") && issue.status === "Resolved") {
    return ["Triaged"];
  }
  return [];
}

async function loadMeta() {
  const response = await fetch("/api/meta");
  if (!response.ok) {
    throw new Error(await response.text() || "failed to load metadata");
  }
  state.meta = await response.json();
  userSelect.innerHTML = "";
  for (const identity of state.meta.identities) {
    const option = document.createElement("option");
    option.value = identity.id;
    option.textContent = `${identity.id} (${identity.role})`;
    userSelect.appendChild(option);
  }
  userSelect.value = state.user;
  issueSeverity.innerHTML = "";
  for (const severity of state.meta.severities) {
    const option = document.createElement("option");
    option.value = severity;
    option.textContent = severity;
    issueSeverity.appendChild(option);
  }
  filterStatus.innerHTML = `<option value="">All statuses</option>`;
  for (const status of state.meta.statuses) {
    const option = document.createElement("option");
    option.value = status;
    option.textContent = status;
    filterStatus.appendChild(option);
  }
  filterAssignee.innerHTML = `<option value="">Anyone</option>`;
  for (const identity of state.meta.identities.filter((identity) => identity.role === "engineer")) {
    const option = document.createElement("option");
    option.value = identity.id;
    option.textContent = identity.id;
    filterAssignee.appendChild(option);
  }
  newIssueButton.hidden = currentRole() !== "reporter";
}

async function loadIssues() {
  const query = new URLSearchParams();
  if (filterStatus.value) {
    query.set("status", filterStatus.value);
  }
  if (filterAssignee.value) {
    query.set("assignee", filterAssignee.value);
  }
  const response = await fetch(`/api/issues?${query.toString()}`);
  if (!response.ok) {
    throw new Error(await response.text() || "failed to load issues");
  }
  const payload = await response.json();
  state.issues = payload.issues || [];
  renderIssueList();
}

async function loadIssue(issueID) {
  if (!issueID) {
    state.currentIssue = null;
    renderDetail();
    return;
  }
  const response = await fetch(`/api/issues/${encodeURIComponent(issueID)}`);
  if (!response.ok) {
    state.currentIssue = null;
    renderDetail();
    throw new Error(await response.text() || `failed to load ${issueID}`);
  }
  state.currentIssue = await response.json();
  renderIssueList();
  renderDetail();
}

// Intent: Keep the browser stable when API calls fail by surfacing one clear
// error banner and leaving the page in a coherent queue/detail state rather
// than letting startup or refresh throw the whole module into a broken view.
// Source: DI-fakuv
async function refreshAll(issueID = state.currentIssueID) {
  try {
    await loadIssues();
    if (issueID) {
      await loadIssue(issueID);
    } else {
      renderDetail();
    }
  } catch (error) {
    setMessage("error", error instanceof Error ? error.message : String(error));
  }
}

newIssueButton.addEventListener("click", () => {
  setMessage("", "");
  openNewIssuePanel();
});

cancelNewIssueButton.addEventListener("click", () => {
  closeNewIssuePanel();
});

userSelect.addEventListener("change", () => {
  state.user = userSelect.value;
  localStorage.setItem("ex4-bug-user", state.user);
  newIssueButton.hidden = currentRole() !== "reporter";
  if (currentRole() !== "reporter") {
    closeNewIssuePanel();
  }
  renderDetail();
});

filterStatus.addEventListener("change", () => {
  refreshAll();
});

filterAssignee.addEventListener("change", () => {
  refreshAll();
});

newIssueForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    setMessage("", "");
    const payload = {
      title: document.getElementById("issue-title").value,
      description: document.getElementById("issue-description").value,
      severity: issueSeverity.value,
    };
    const response = await apiFetch("/api/issues", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    if (!response.ok) {
      setMessage("error", await response.text());
      return;
    }
    const issue = await response.json();
    setCurrentIssueID(issue.id);
    closeNewIssuePanel();
    await refreshAll(issue.id);
    setMessage("ok", `Created ${issue.id}.`);
  } catch (error) {
    setMessage("error", error instanceof Error ? error.message : String(error));
  }
});

assignmentForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    const response = await apiFetch(`/api/issues/${encodeURIComponent(state.currentIssueID)}/assignment`, {
      method: "POST",
      body: JSON.stringify({ assignee: assignmentSelect.value }),
    });
    if (!response.ok) {
      setMessage("error", await response.text());
      return;
    }
    await refreshAll(state.currentIssueID);
    setMessage("ok", "Assignment updated.");
  } catch (error) {
    setMessage("error", error instanceof Error ? error.message : String(error));
  }
});

statusForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    const response = await apiFetch(`/api/issues/${encodeURIComponent(state.currentIssueID)}/status`, {
      method: "POST",
      body: JSON.stringify({ status: statusSelect.value }),
    });
    if (!response.ok) {
      setMessage("error", await response.text());
      return;
    }
    await refreshAll(state.currentIssueID);
    setMessage("ok", "Status updated.");
  } catch (error) {
    setMessage("error", error instanceof Error ? error.message : String(error));
  }
});

commentForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    const response = await apiFetch(`/api/issues/${encodeURIComponent(state.currentIssueID)}/comments`, {
      method: "POST",
      body: JSON.stringify({ comment: commentBody.value }),
    });
    if (!response.ok) {
      setMessage("error", await response.text());
      return;
    }
    commentBody.value = "";
    await refreshAll(state.currentIssueID);
    setMessage("ok", "Comment posted.");
  } catch (error) {
    setMessage("error", error instanceof Error ? error.message : String(error));
  }
});

attachmentForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    if (!attachmentFile.files.length) {
      setMessage("error", "Choose a file first.");
      return;
    }
    const formData = new FormData();
    formData.append("attachment", attachmentFile.files[0]);
    const response = await apiFetch(`/api/issues/${encodeURIComponent(state.currentIssueID)}/attachments`, {
      method: "POST",
      body: formData,
    });
    if (!response.ok) {
      setMessage("error", await response.text());
      return;
    }
    attachmentFile.value = "";
    await refreshAll(state.currentIssueID);
    setMessage("ok", "Attachment uploaded.");
  } catch (error) {
    setMessage("error", error instanceof Error ? error.message : String(error));
  }
});

try {
  await loadMeta();
  await refreshAll(state.currentIssueID);
} catch (error) {
  setMessage("error", error instanceof Error ? error.message : String(error));
}
