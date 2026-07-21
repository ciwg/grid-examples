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
  vim.api.nvim_create_user_command("OksClose", function()
    M.close()
  end, {})
end

return M
