local M = {}

local uv = vim.uv or vim.loop

local function default_participant_id()
  local pid = tostring(vim.fn.getpid())
  local host = (uv.os_gethostname and uv.os_gethostname()) or "host"
  host = host:gsub("[^%w%-_]", "-")
  return table.concat({ "oks", "nvim", host, pid }, "-")
end

local function default_repo_root()
  local source = debug.getinfo(1, "S").source
  if vim.startswith(source, "@") then
    source = source:sub(2)
  end
  local path = vim.fs.normalize(source)
  return vim.fs.dirname(vim.fs.dirname(vim.fs.dirname(vim.fs.dirname(path))))
end

M.config = {
  repo_root = default_repo_root(),
  base_url = vim.env.OKS_BASE_URL or "http://127.0.0.1:7045",
  display_name = vim.env.OKS_DISPLAY_NAME or "Neovim User",
  color = vim.env.OKS_COLOR or "#d66f1d",
  poll_ms = 3000,
}

M.state = {
  item_id = nil,
  bufnr = nil,
  title = "",
  status = "",
  version = 0,
  current_revision = 0,
  participants = {},
  participant_id = default_participant_id(),
  applying_remote = false,
  poll_timer = nil,
  augroup = nil,
}

local inspector = {
  bufnr = nil,
  winid = nil,
}

local function notify(message, level)
  vim.notify("oks: " .. message, level or vim.log.levels.INFO)
end

local function split_lines(body)
  if body == "" then
    return { "" }
  end
  return vim.split(body, "\n", { plain = true })
end

local function current_body()
  if not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return ""
  end
  return table.concat(vim.api.nvim_buf_get_lines(M.state.bufnr, 0, -1, false), "\n")
end

local function current_cursor_offset()
  if not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return 0
  end
  local cursor = vim.api.nvim_win_get_cursor(0)
  local row = cursor[1]
  local col = cursor[2]
  local total = 0
  local lines = vim.api.nvim_buf_get_lines(M.state.bufnr, 0, -1, false)
  for index = 1, row - 1 do
    total = total + #(lines[index] or "") + 1
  end
  return total + col
end

local function json_decode(text)
  if text == "" then
    return nil
  end
  if vim.json and vim.json.decode then
    return vim.json.decode(text)
  end
  return vim.fn.json_decode(text)
end

local function json_encode(value)
  if vim.json and vim.json.encode then
    return vim.json.encode(value)
  end
  return vim.fn.json_encode(value)
end

local function sorted_keys(values)
  local out = {}
  for key, _ in pairs(values or {}) do
    table.insert(out, key)
  end
  table.sort(out)
  return out
end

-- Intent: Reuse the existing ex5 live-draft HTTP surface from Neovim instead
-- of inventing a separate transport for the first embodiment phase. Source:
-- DI-fudok
local function request(method, path, payload)
  local argv = {
    "curl",
    "-sS",
    "-o",
    "-",
    "-w",
    "\n%{http_code}",
  }
  if method == "POST" then
    table.insert(argv, "-X")
    table.insert(argv, "POST")
    table.insert(argv, "-H")
    table.insert(argv, "Content-Type: application/json")
    table.insert(argv, "--data-binary")
    table.insert(argv, "@-")
  end
  table.insert(argv, M.config.base_url .. path)
  local input = payload and json_encode(payload) or ""
  local raw = vim.fn.system(argv, input)
  local shell_error = vim.v.shell_error
  local status_text = raw:match("\n(%d%d%d)%s*$")
  if not status_text then
    return nil, nil, string.format("request failed: %s", raw)
  end
  local status = tonumber(status_text)
  local body = raw:gsub("\n%d%d%d%s*$", "")
  if shell_error ~= 0 and status == nil then
    return nil, nil, string.format("request failed: %s", raw)
  end
  return status, body, nil
end

local function apply_live_state(state, force)
  M.state.title = state.title or M.state.title
  M.state.status = state.status or M.state.status
  M.state.version = state.version or M.state.version
  M.state.current_revision = state.current_revision or M.state.current_revision
  M.state.participants = state.participants or {}
  if not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return
  end
  if force or not vim.bo[M.state.bufnr].modified then
    M.state.applying_remote = true
    vim.api.nvim_buf_set_lines(M.state.bufnr, 0, -1, false, split_lines(state.body or ""))
    vim.bo[M.state.bufnr].modified = false
    M.state.applying_remote = false
  end
end

local function open_scratch_buffer(name)
  local bufnr = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_name(bufnr, name)
  vim.bo[bufnr].buftype = "nofile"
  vim.bo[bufnr].bufhidden = "wipe"
  vim.bo[bufnr].swapfile = false
  vim.bo[bufnr].modifiable = true
  vim.bo[bufnr].filetype = "markdown"
  vim.cmd("botright vsplit")
  local winid = vim.api.nvim_get_current_win()
  vim.api.nvim_win_set_buf(winid, bufnr)
  vim.api.nvim_win_set_option(winid, "wrap", false)
  return bufnr, winid
end

local function refresh_live_state(force)
  if not M.state.item_id then
    return false
  end
  local status, body, err = request("GET", "/api/items/" .. M.state.item_id .. "/live")
  if err then
    notify(err, vim.log.levels.ERROR)
    return false
  end
  if status ~= 200 then
    notify("refresh failed: " .. tostring(status), vim.log.levels.ERROR)
    return false
  end
  local decoded = json_decode(body)
  if not decoded then
    notify("refresh decode failed", vim.log.levels.ERROR)
    return false
  end
  if not force and vim.bo[M.state.bufnr].modified and decoded.version ~= M.state.version then
    M.state.participants = decoded.participants or M.state.participants
    M.state.title = decoded.title or M.state.title
    M.state.status = decoded.status or M.state.status
    M.state.current_revision = decoded.current_revision or M.state.current_revision
    M.state.version = decoded.version or M.state.version
    notify("remote live draft advanced; local buffer left unchanged", vim.log.levels.WARN)
    return false
  end
  apply_live_state(decoded, force)
  return true
end

-- Intent: Keep Neovim visible in the shared participant roster without
-- advancing the draft version when the local editor is only reporting presence
-- and cursor state. Source: DI-fudok
local function push_presence(typing)
  if not M.state.item_id then
    return
  end
  local cursor = current_cursor_offset()
  local status, body, err = request("POST", "/api/items/" .. M.state.item_id .. "/live", {
    participant_id = M.state.participant_id,
    display_name = M.config.display_name,
    color = M.config.color,
    cursor = cursor,
    head = cursor,
    typing = typing,
    base_version = M.state.version,
    body = "",
  })
  if err or status ~= 200 then
    return
  end
  local decoded = json_decode(body)
  if decoded then
    M.state.participants = decoded.participants or M.state.participants
    M.state.version = decoded.version or M.state.version
  end
end

-- Intent: Make :write participate in the same optimistic live-draft flow as
-- the browser so Neovim does not blindly overwrite a newer shared body.
-- Source: DI-fudok
local function push_body()
  if not M.state.item_id or not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return
  end
  local cursor = current_cursor_offset()
  local status, body, err = request("POST", "/api/items/" .. M.state.item_id .. "/live", {
    participant_id = M.state.participant_id,
    display_name = M.config.display_name,
    color = M.config.color,
    cursor = cursor,
    head = cursor,
    typing = false,
    base_version = M.state.version,
    body = current_body(),
  })
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  local decoded = json_decode(body)
  if status == 409 then
    if decoded and decoded.state then
      M.state.participants = decoded.state.participants or M.state.participants
      M.state.version = decoded.state.version or M.state.version
      M.state.current_revision = decoded.state.current_revision or M.state.current_revision
      M.state.status = decoded.state.status or M.state.status
      M.state.title = decoded.state.title or M.state.title
    end
    notify("live draft conflict; refresh or merge before retrying", vim.log.levels.WARN)
    return
  end
  if status ~= 200 or not decoded then
    notify("push failed: " .. tostring(status), vim.log.levels.ERROR)
    return
  end
  apply_live_state(decoded, true)
  notify("pushed " .. M.state.item_id)
end

local function stop_poll_loop()
  if M.state.poll_timer then
    M.state.poll_timer:stop()
    M.state.poll_timer:close()
    M.state.poll_timer = nil
  end
end

-- Intent: Keep the first Neovim phase close to the browser draft studio by
-- polling the same runtime state and refreshing presence on a short interval.
-- Source: DI-fudok
local function start_poll_loop()
  stop_poll_loop()
  M.state.poll_timer = uv.new_timer()
  M.state.poll_timer:start(M.config.poll_ms, M.config.poll_ms, function()
    vim.schedule(function()
      if not M.state.item_id or not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
        stop_poll_loop()
        return
      end
      refresh_live_state(false)
      push_presence(false)
    end)
  end)
end

local function session_lines()
  local lines = {
    "oks live draft",
    "item: " .. (M.state.item_id or "-"),
    "title: " .. (M.state.title or "-"),
    "status: " .. (M.state.status or "-"),
    "version: " .. tostring(M.state.version or 0),
    "revision: " .. tostring(M.state.current_revision or 0),
    "participant: " .. M.state.participant_id,
    "base_url: " .. M.config.base_url,
    "participants:",
  }
  if #(M.state.participants or {}) == 0 then
    table.insert(lines, "  - none")
  else
    for _, participant in ipairs(M.state.participants) do
      table.insert(lines, string.format("  - %s cursor=%d typing=%s", participant.display_name or participant.participant_id or "peer", participant.cursor or 0, tostring(participant.typing or false)))
    end
  end
  return lines
end

local function append_list(lines, heading, values)
  table.insert(lines, "")
  table.insert(lines, "## " .. heading)
  if not values or #values == 0 then
    table.insert(lines, "- none")
    return
  end
  for _, value in ipairs(values) do
    table.insert(lines, "- " .. value)
  end
end

local function evidence_fact_summary(evidence)
  local facts = evidence and evidence.facts or {}
  local parts = {}
  for _, key in ipairs(sorted_keys(facts)) do
    table.insert(parts, key .. "=" .. tostring(facts[key]))
  end
  return table.concat(parts, ", ")
end

local function approval_summary(approval)
  return string.format("- %s by %s role=%s revision=%d", approval.decision or "-", approval.actor or "-", approval.role or "-", approval.revision or 0)
end

-- Intent: Let Neovim users inspect the durable item record around the live
-- draft without leaving the editor or attempting write actions that this phase
-- does not support yet. Source: DI-lonuk
local function item_detail_lines(detail)
  local lines = {
    "# " .. (detail.title or M.state.title or M.state.item_id or "knowledge item"),
    "",
    "- id: " .. (detail.id or M.state.item_id or "-"),
    "- kind: " .. (detail.kind or "-"),
    "- status: " .. (detail.status or "-"),
    "- current revision: " .. tostring(detail.current_revision or 0),
    "- working version: " .. tostring(detail.working_version or 0),
    "- summary: " .. (detail.summary or ""),
  }

  append_list(lines, "Responsibilities", detail.responsibility_ids or {})

  table.insert(lines, "")
  table.insert(lines, "## Revisions")
  if not detail.revisions or #detail.revisions == 0 then
    table.insert(lines, "- none")
  else
    for _, revision in ipairs(detail.revisions) do
      table.insert(lines, string.format("- r%d %s — %s", revision.number or 0, revision.title or "-", revision.created_at or "-"))
      if revision.summary and revision.summary ~= "" then
        table.insert(lines, "  " .. revision.summary)
      end
    end
  end

  table.insert(lines, "")
  table.insert(lines, "## Approvals")
  if not detail.approvals or #detail.approvals == 0 then
    table.insert(lines, "- none")
  else
    for _, approval in ipairs(detail.approvals) do
      table.insert(lines, string.format("- %s by %s role=%s revision=%d", approval.decision or "-", approval.actor or "-", approval.role or "-", approval.revision or 0))
      if approval.notes and approval.notes ~= "" then
        table.insert(lines, "  " .. approval.notes)
      end
    end
  end

  table.insert(lines, "")
  table.insert(lines, "## Related runs")
  if not detail.related_runs or #detail.related_runs == 0 then
    table.insert(lines, "- none")
  else
    for _, run in ipairs(detail.related_runs) do
      table.insert(lines, string.format("- %s kind=%s revision=%d outcome=%s", run.id or "-", run.kind or "-", run.revision or 0, run.outcome or "-"))
      if run.notes and run.notes ~= "" then
        table.insert(lines, "  notes: " .. run.notes)
      end
      if run.place_id and run.place_id ~= "" then
        table.insert(lines, "  place: " .. run.place_id)
      end
      if run.resource_ids and #run.resource_ids > 0 then
        table.insert(lines, "  resources: " .. table.concat(run.resource_ids, ", "))
      end
      if run.evidence and #run.evidence > 0 then
        for _, evidence in ipairs(run.evidence) do
          local summary = evidence_fact_summary(evidence)
          if summary ~= "" then
            table.insert(lines, "  evidence: " .. summary)
          end
        end
      end
    end
  end

  return lines
end

-- Intent: Reuse the existing item detail projection for Neovim inspection so
-- the editor sees the same revision, approval, and related-run truth as the
-- browser and CLI. Source: DI-lonuk
local function inspect_item(item_id)
  item_id = vim.trim(item_id or M.state.item_id or "")
  if item_id == "" then
    notify("item id is required", vim.log.levels.WARN)
    return
  end
  local status, body, err = request("GET", "/api/items/" .. item_id)
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    notify("inspect failed: " .. tostring(status), vim.log.levels.ERROR)
    return
  end
  local decoded = json_decode(body)
  if not decoded then
    notify("inspect decode failed", vim.log.levels.ERROR)
    return
  end

  if inspector.winid and vim.api.nvim_win_is_valid(inspector.winid) and inspector.bufnr and vim.api.nvim_buf_is_valid(inspector.bufnr) then
    vim.api.nvim_set_current_win(inspector.winid)
  else
    inspector.bufnr, inspector.winid = open_scratch_buffer("oks-inspect://" .. item_id)
  end
  vim.bo[inspector.bufnr].modifiable = true
  vim.api.nvim_buf_set_name(inspector.bufnr, "oks-inspect://" .. item_id)
  vim.api.nvim_buf_set_lines(inspector.bufnr, 0, -1, false, item_detail_lines(decoded))
  vim.bo[inspector.bufnr].modifiable = false
  notify("inspected " .. item_id)
end

-- Intent: Let Neovim users inspect a specific run's evidence and approval
-- record directly from the existing run detail projection before any write-side
-- workflow actions are added to the editor. Source: DI-ravok
local function run_detail_lines(detail)
  local lines = {
    "# " .. (detail.id or "run"),
    "",
    "- kind: " .. (detail.kind or "-"),
    "- item_id: " .. (detail.item_id or "-"),
    "- item_kind: " .. (detail.item_kind or "-"),
    "- revision: " .. tostring(detail.revision or 0),
    "- actor: " .. (detail.actor or "-"),
    "- outcome: " .. (detail.outcome or "-"),
    "- place: " .. (detail.place_id or "-"),
    "- resources: " .. ((detail.resource_ids and #detail.resource_ids > 0) and table.concat(detail.resource_ids, ", ") or "-"),
  }

  if detail.notes and detail.notes ~= "" then
    table.insert(lines, "- notes: " .. detail.notes)
  end

  table.insert(lines, "")
  table.insert(lines, "## Evidence")
  if not detail.evidence or #detail.evidence == 0 then
    table.insert(lines, "- none")
  else
    for _, evidence in ipairs(detail.evidence) do
      table.insert(lines, "- " .. (evidence.summary or "evidence"))
      local summary = evidence_fact_summary(evidence)
      if summary ~= "" then
        table.insert(lines, "  " .. summary)
      end
      if evidence.attachment_name and evidence.attachment_name ~= "" then
        table.insert(lines, "  attachment: " .. evidence.attachment_name)
      end
    end
  end

  table.insert(lines, "")
  table.insert(lines, "## Approvals")
  if not detail.approvals or #detail.approvals == 0 then
    table.insert(lines, "- none")
  else
    for _, approval in ipairs(detail.approvals) do
      table.insert(lines, approval_summary(approval))
      if approval.notes and approval.notes ~= "" then
        table.insert(lines, "  " .. approval.notes)
      end
    end
  end

  return lines
end

-- Intent: Reuse the existing run detail projection for Neovim review so the
-- editor sees the same run evidence and approvals as the browser and CLI.
-- Source: DI-ravok
local function inspect_run(run_id)
  run_id = vim.trim(run_id or "")
  if run_id == "" then
    notify("run id is required", vim.log.levels.WARN)
    return
  end
  local status, body, err = request("GET", "/api/runs/" .. run_id)
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    notify("inspect run failed: " .. tostring(status), vim.log.levels.ERROR)
    return
  end
  local decoded = json_decode(body)
  if not decoded then
    notify("inspect run decode failed", vim.log.levels.ERROR)
    return
  end

  if inspector.winid and vim.api.nvim_win_is_valid(inspector.winid) and inspector.bufnr and vim.api.nvim_buf_is_valid(inspector.bufnr) then
    vim.api.nvim_set_current_win(inspector.winid)
  else
    inspector.bufnr, inspector.winid = open_scratch_buffer("oks-run://" .. run_id)
  end
  vim.bo[inspector.bufnr].modifiable = true
  vim.api.nvim_buf_set_name(inspector.bufnr, "oks-run://" .. run_id)
  vim.api.nvim_buf_set_lines(inspector.bufnr, 0, -1, false, run_detail_lines(decoded))
  vim.bo[inspector.bufnr].modifiable = false
  notify("inspected run " .. run_id)
end

function M.info()
  notify(table.concat(session_lines(), "\n"))
end

function M.refresh(force)
  if not M.state.item_id then
    notify("no active item", vim.log.levels.WARN)
    return
  end
  refresh_live_state(force == true)
end

function M.push()
  push_body()
end

function M.inspect(item_id)
  inspect_item(item_id)
end

function M.inspect_run(run_id)
  inspect_run(run_id)
end

function M.close()
  stop_poll_loop()
  if M.state.augroup then
    pcall(vim.api.nvim_del_augroup_by_id, M.state.augroup)
    M.state.augroup = nil
  end
  M.state.item_id = nil
  M.state.participants = {}
end

function M.open(item_id)
  item_id = vim.trim(item_id or "")
  if item_id == "" then
    notify("item id is required", vim.log.levels.ERROR)
    return
  end
  M.close()
  M.state.item_id = item_id
  M.state.bufnr = vim.api.nvim_create_buf(true, false)
  vim.api.nvim_buf_set_name(M.state.bufnr, "oks://" .. item_id)
  vim.api.nvim_set_current_buf(M.state.bufnr)
  vim.bo[M.state.bufnr].buftype = "acwrite"
  vim.bo[M.state.bufnr].swapfile = false
  vim.bo[M.state.bufnr].filetype = "markdown"

  if not refresh_live_state(true) then
    return
  end
  push_presence(false)

  M.state.augroup = vim.api.nvim_create_augroup("OksLiveDraft", { clear = true })
  vim.api.nvim_create_autocmd("BufWriteCmd", {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function()
      push_body()
    end,
  })
  vim.api.nvim_create_autocmd({ "InsertEnter", "InsertLeave" }, {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function(event)
      push_presence(event.event == "InsertEnter")
    end,
  })
  vim.api.nvim_create_autocmd({ "TextChanged", "TextChangedI" }, {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function()
      if M.state.applying_remote then
        return
      end
      push_presence(true)
    end,
  })
  vim.api.nvim_create_autocmd("BufUnload", {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function()
      M.close()
    end,
  })
  start_poll_loop()
  notify("opened " .. item_id)
end

function M.setup(opts)
  M.config = vim.tbl_deep_extend("force", M.config, opts or {})
  vim.api.nvim_create_user_command("OksOpen", function(command)
    M.open(command.args)
  end, { nargs = 1 })
  vim.api.nvim_create_user_command("OksRefresh", function()
    M.refresh(true)
  end, {})
  vim.api.nvim_create_user_command("OksPush", function()
    M.push()
  end, {})
  vim.api.nvim_create_user_command("OksInfo", function()
    M.info()
  end, {})
  vim.api.nvim_create_user_command("OksInspect", function(command)
    M.inspect(command.args)
  end, { nargs = "?" })
  vim.api.nvim_create_user_command("OksInspectRun", function(command)
    M.inspect_run(command.args)
  end, { nargs = 1 })
  vim.api.nvim_create_user_command("OksClose", function()
    M.close()
  end, {})
end

return M
