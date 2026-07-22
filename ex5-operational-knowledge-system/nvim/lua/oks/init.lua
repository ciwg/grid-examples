local M = {}

local uv = vim.uv or vim.loop
local bitops = bit or bit32

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
  socket_path = vim.env.OKS_SOCKET_PATH or ".operational-knowledge-system/embodiment.sock",
  base_url = vim.env.OKS_BASE_URL or "http://127.0.0.1:7045",
  display_name = vim.env.OKS_DISPLAY_NAME or "Neovim User",
  color = vim.env.OKS_COLOR or "#d66f1d",
  poll_ms = 3000,
}

M.state = {
  item_id = nil,
  bufnr = nil,
  winid = nil,
  title = "",
  status = "",
  version = 0,
  current_revision = 0,
  participants = {},
  participant_id = default_participant_id(),
  applying_remote = false,
  poll_timer = nil,
  heartbeat_timer = nil,
  reconnect_timer = nil,
  socket = nil,
  socket_connected = false,
  socket_buffer = "",
  socket_generation = 0,
  transport = "http-poll",
  augroup = nil,
}

local inspector = {
  bufnr = nil,
  winid = nil,
}

local apply_live_state
local stop_poll_loop
local start_poll_loop

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

-- Intent: Keep Neovim presence and push offsets tied to the actual live-draft
-- window instead of whichever split happens to be focused after inspectors open.
-- Source: DI-pazud
local function live_draft_winid()
  if M.state.winid and vim.api.nvim_win_is_valid(M.state.winid) then
    if vim.api.nvim_win_get_buf(M.state.winid) == M.state.bufnr then
      return M.state.winid
    end
  end
  for _, winid in ipairs(vim.api.nvim_list_wins()) do
    if vim.api.nvim_win_get_buf(winid) == M.state.bufnr then
      M.state.winid = winid
      return winid
    end
  end
  return nil
end

local function current_cursor_offset()
  if not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return 0
  end
  local winid = live_draft_winid()
  if not winid then
    return 0
  end
  local cursor = vim.api.nvim_win_get_cursor(winid)
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

local function url_encode(value)
  if vim.uri_encode then
    return vim.uri_encode(value)
  end
  return (value:gsub("([^%w%-_%.~])", function(char)
    return string.format("%%%02X", string.byte(char))
  end))
end

local function websocket_url(path)
  local base = M.config.base_url
  if vim.startswith(base, "https://") then
    return "wss://" .. base:sub(#("https://") + 1) .. path
  end
  if vim.startswith(base, "http://") then
    return "ws://" .. base:sub(#("http://") + 1) .. path
  end
  return nil
end

local function socket_path_available()
  return M.config.socket_path and M.config.socket_path ~= "" and uv.fs_stat(M.config.socket_path) ~= nil
end

local function split_json_lines(buffer)
  local messages = {}
  while true do
    local newline = buffer:find("\n", 1, true)
    if not newline then
      break
    end
    local line = vim.trim(buffer:sub(1, newline - 1))
    buffer = buffer:sub(newline + 1)
    if line ~= "" then
      table.insert(messages, line)
    end
  end
  return messages, buffer
end

local function parse_websocket_url(url)
  if not url then
    return nil, "websocket url required"
  end
  local scheme, rest = url:match("^(wss?)://(.+)$")
  if not scheme then
    return nil, "unsupported websocket url"
  end
  local host_port, path = rest:match("^([^/]+)(/.*)$")
  if not host_port then
    host_port = rest
    path = "/"
  end
  local host, port = host_port:match("^([^:]+):(%d+)$")
  if not host then
    host = host_port
    port = scheme == "wss" and "443" or "80"
  end
  return {
    scheme = scheme,
    host = host,
    port = tonumber(port),
    path = path,
  }, nil
end

local function encode_websocket_text_frame(text)
  local length = #text
  if length < 126 then
    return string.char(0x81, length) .. text
  end
  if length <= 0xFFFF then
    local hi = math.floor(length / 256)
    local lo = length % 256
    return string.char(0x81, 126, hi, lo) .. text
  end
  return nil, "websocket frame too large"
end

local function decode_websocket_frames(buffer)
  local frames = {}
  local index = 1
  while true do
    if #buffer - index + 1 < 2 then
      break
    end
    local first = string.byte(buffer, index)
    local second = string.byte(buffer, index + 1)
    if bitops.band(first, 0x80) == 0 then
      return nil, nil, "fragmented websocket frames are unsupported"
    end
    local opcode = bitops.band(first, 0x0f)
    local masked = bitops.band(second, 0x80) ~= 0
    local length = bitops.band(second, 0x7f)
    local header = 2
    if length == 126 then
      if #buffer - index + 1 < 4 then
        break
      end
      local hi = string.byte(buffer, index + 2)
      local lo = string.byte(buffer, index + 3)
      length = hi * 256 + lo
      header = 4
    elseif length == 127 then
      return nil, nil, "64-bit websocket frames are unsupported"
    end
    local mask = nil
    if masked then
      if #buffer - index + 1 < header + 4 then
        break
      end
      mask = { string.byte(buffer, index + header, index + header + 3) }
      header = header + 4
    end
    if #buffer - index + 1 < header + length then
      break
    end
    local payload = buffer:sub(index + header, index + header + length - 1)
    if mask then
      local bytes = { string.byte(payload, 1, #payload) }
      for i = 1, #bytes do
        bytes[i] = bitops.bxor(bytes[i], mask[((i - 1) % 4) + 1])
      end
      payload = string.char(unpack(bytes))
    end
    table.insert(frames, { opcode = opcode, payload = payload })
    index = index + header + length
  end
  return frames, buffer:sub(index), nil
end

local allowed_search_filters = {
  kind = true,
  status = true,
  outcome = true,
  place_id = true,
  resource_id = true,
  responsibility_id = true,
  problem = true,
}

local function socket_request(method, path, payload)
  local pipe = uv.new_pipe(false)
  local done = false
  local response_text = ""
  local connect_err = nil
  local write_err = nil
  local request = {
    type = "request",
    method = method,
    path = path,
  }
  if payload ~= nil then
    request.body = json_encode(payload)
    request.headers = { ["Content-Type"] = "application/json" }
  end
  pipe:connect(M.config.socket_path, function(err)
    if err then
      connect_err = err
      done = true
      return
    end
    pipe:write(json_encode(request) .. "\n", function(write_error)
      if write_error then
        write_err = write_error
        done = true
        return
      end
      pipe:shutdown(function()
      end)
    end)
    pipe:read_start(function(read_err, chunk)
      if read_err then
        connect_err = read_err
        done = true
        return
      end
      if not chunk then
        done = true
        return
      end
      response_text = response_text .. chunk
    end)
  end)
  local finished = vim.wait(5000, function()
    return done
  end, 10)
  pcall(pipe.read_stop, pipe)
  pcall(pipe.close, pipe)
  if not finished then
    return nil, nil, "socket request timed out"
  end
  if connect_err then
    return nil, nil, "socket request failed: " .. tostring(connect_err)
  end
  if write_err then
    return nil, nil, "socket write failed: " .. tostring(write_err)
  end
  local decoded = json_decode(response_text)
  if not decoded then
    return nil, nil, "socket response decode failed"
  end
  if decoded.type == "error" then
    return nil, nil, decoded.message or "socket request failed"
  end
  return decoded.status, decoded.body or "", nil
end

-- Intent: Prefer the direct local Unix-socket embodiment contract for
-- Neovim requests while preserving HTTP fallback so editor workflows stay
-- operable when only the browser adapter is available. Source: DI-favel
local function request(method, path, payload)
  if socket_path_available() then
    local status, body, err = socket_request(method, path, payload)
    if err == nil then
      return status, body, nil
    end
  end
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

local function apply_socket_live_state(decoded, conflict)
  if conflict then
    if decoded.state then
      M.state.participants = decoded.state.participants or M.state.participants
      M.state.version = decoded.state.version or M.state.version
      M.state.current_revision = decoded.state.current_revision or M.state.current_revision
      M.state.status = decoded.state.status or M.state.status
      M.state.title = decoded.state.title or M.state.title
      apply_live_state(decoded.state, true)
    end
    notify("live draft conflict; refresh or merge before retrying", vim.log.levels.WARN)
    return
  end
  local state = decoded.state
  if not state then
    return
  end
  if current_body() == (state.body or "") then
    vim.bo[M.state.bufnr].modified = false
  end
  apply_live_state(state, not vim.bo[M.state.bufnr].modified)
end

apply_live_state = function(state, force)
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

local function wipe_buffer(bufnr)
  if not bufnr or not vim.api.nvim_buf_is_valid(bufnr) then
    return
  end
  pcall(vim.api.nvim_buf_delete, bufnr, { force = true })
end

local function buffer_name(bufnr)
  if not bufnr or not vim.api.nvim_buf_is_valid(bufnr) then
    return ""
  end
  return vim.api.nvim_buf_get_name(bufnr)
end

local function item_id_from_buffer_name(name)
  local item_id = name:match("^oks%-inspect://(.+)$")
  if item_id and item_id ~= "" then
    return item_id
  end
  item_id = name:match("^oks://(.+)$")
  if item_id and item_id ~= "" then
    return item_id
  end
  return nil
end

local function run_id_from_buffer_name(name)
  local run_id = name:match("^oks%-run://(.+)$")
  if run_id and run_id ~= "" then
    return run_id
  end
  return nil
end

local function current_context_item_id()
  local current = item_id_from_buffer_name(buffer_name(vim.api.nvim_get_current_buf()))
  if current then
    return current
  end
  if M.state.item_id and M.state.item_id ~= "" then
    return M.state.item_id
  end
  return nil
end

local function current_context_run_id()
  return run_id_from_buffer_name(buffer_name(vim.api.nvim_get_current_buf()))
end

local function valid_decision(decision)
  return decision == "approved" or decision == "rejected" or decision == "noted"
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

local function clear_reconnect_loop()
  if M.state.reconnect_timer then
    M.state.reconnect_timer:stop()
    M.state.reconnect_timer:close()
    M.state.reconnect_timer = nil
  end
end

local function close_live_socket()
  if M.state.socket then
    pcall(M.state.socket.read_stop, M.state.socket)
    pcall(M.state.socket.close, M.state.socket)
    M.state.socket = nil
  end
  M.state.socket_connected = false
  M.state.socket_buffer = ""
  M.state.transport = "http-poll"
end

local function handle_live_socket_payload(payload)
  local decoded = json_decode(payload)
  if not decoded then
    return
  end
  if decoded.type == "error" then
    notify(decoded.message or "live transport error", vim.log.levels.WARN)
    return
  end
  if decoded.type == "live-conflict" then
    apply_socket_live_state(decoded, true)
    return
  end
  if decoded.type == "live-state" then
    apply_socket_live_state(decoded, false)
  end
end

local function schedule_socket_reconnect(item_id, generation)
  clear_reconnect_loop()
  M.state.reconnect_timer = uv.new_timer()
  M.state.reconnect_timer:start(5000, 0, vim.schedule_wrap(function()
    if M.state.item_id ~= item_id or M.state.socket_generation ~= generation then
      return
    end
    M.connect_live_socket(item_id)
  end))
end

local function send_live_socket_update(update_body, typing)
  if not M.state.socket_connected or not M.state.socket then
    return false
  end
  local payload = json_encode({
    type = "live-update",
    participant_id = M.state.participant_id,
    display_name = M.config.display_name,
    color = M.config.color,
    cursor = current_cursor_offset(),
    head = current_cursor_offset(),
    typing = typing,
    base_version = M.state.version,
    update_body = update_body,
    body = update_body and current_body() or "",
  })
  M.state.socket:write(payload .. "\n")
  return true
end

function M.connect_live_socket(item_id)
  if socket_path_available() then
    local generation = M.state.socket_generation + 1
    M.state.socket_generation = generation
    close_live_socket()
    local socket = uv.new_pipe(false)
    M.state.socket = socket
    socket:connect(M.config.socket_path, vim.schedule_wrap(function(err)
      if err then
        close_live_socket()
        start_poll_loop()
        schedule_socket_reconnect(item_id, generation)
        return
      end
      socket:write(json_encode({
        type = "live-open",
        item_id = item_id,
        participant_id = M.state.participant_id,
        display_name = M.config.display_name,
        color = M.config.color,
        cursor = current_cursor_offset(),
        head = current_cursor_offset(),
        typing = false,
      }) .. "\n")
      socket:read_start(vim.schedule_wrap(function(read_err, chunk)
        if generation ~= M.state.socket_generation then
          return
        end
        if read_err or not chunk then
          close_live_socket()
          start_poll_loop()
          schedule_socket_reconnect(item_id, generation)
          return
        end
        M.state.socket_connected = true
        M.state.transport = "local-socket"
        clear_reconnect_loop()
        stop_poll_loop()
        M.state.socket_buffer = (M.state.socket_buffer or "") .. chunk
        local messages, rest = split_json_lines(M.state.socket_buffer)
        M.state.socket_buffer = rest or ""
        for _, message in ipairs(messages) do
          handle_live_socket_payload(message)
        end
      end))
    end))
    return
  end
  local parsed, parse_err = parse_websocket_url(websocket_url("/api/items/" .. item_id .. "/live/socket"))
  if parse_err or not parsed or parsed.scheme ~= "ws" then
    start_poll_loop()
    return
  end
  local generation = M.state.socket_generation + 1
  M.state.socket_generation = generation
  close_live_socket()
  local socket = uv.new_tcp()
  M.state.socket = socket
  socket:connect(parsed.host, parsed.port, vim.schedule_wrap(function(err)
    if err then
      close_live_socket()
      start_poll_loop()
      schedule_socket_reconnect(item_id, generation)
      return
    end
    local request_text = table.concat({
      "GET " .. parsed.path .. " HTTP/1.1",
      "Host: " .. parsed.host .. ":" .. tostring(parsed.port),
      "Upgrade: websocket",
      "Connection: Upgrade",
      "Sec-WebSocket-Version: 13",
      "Sec-WebSocket-Key: b2tzLW5laW92aW0td3M=",
      "",
      "",
    }, "\r\n")
    socket:write(request_text)
    socket:read_start(vim.schedule_wrap(function(read_err, chunk)
      if generation ~= M.state.socket_generation then
        return
      end
      if read_err or not chunk then
        close_live_socket()
        start_poll_loop()
        schedule_socket_reconnect(item_id, generation)
        return
      end
      M.state.socket_buffer = (M.state.socket_buffer or "") .. chunk
      if not M.state.socket_connected then
        local header_end = M.state.socket_buffer:find("\r\n\r\n", 1, true)
        if not header_end then
          return
        end
        local header = M.state.socket_buffer:sub(1, header_end + 3)
        if not header:match("^HTTP/1%.1 101 ") then
          close_live_socket()
          start_poll_loop()
          schedule_socket_reconnect(item_id, generation)
          return
        end
        M.state.socket_buffer = M.state.socket_buffer:sub(header_end + 4)
        M.state.socket_connected = true
        M.state.transport = "websocket"
        clear_reconnect_loop()
        stop_poll_loop()
        send_live_socket_update(false, false)
      end
      local frames, rest, frame_err = decode_websocket_frames(M.state.socket_buffer)
      if frame_err then
        notify(frame_err, vim.log.levels.WARN)
        close_live_socket()
        start_poll_loop()
        schedule_socket_reconnect(item_id, generation)
        return
      end
      M.state.socket_buffer = rest or ""
      for _, frame in ipairs(frames or {}) do
        if frame.opcode == 0x1 then
          handle_live_socket_payload(frame.payload)
        elseif frame.opcode == 0x8 then
          close_live_socket()
          start_poll_loop()
          schedule_socket_reconnect(item_id, generation)
          return
        end
      end
    end))
  end))
end

-- Intent: Keep Neovim visible in the shared participant roster without
-- advancing the draft version when the local editor is only reporting presence
-- and cursor state. Source: DI-fudok
local function push_presence(typing)
  if not M.state.item_id then
    return
  end
  if send_live_socket_update(false, typing) then
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
    update_body = false,
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
    return false
  end
  if send_live_socket_update(true, false) then
    return true
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
    update_body = true,
    body = current_body(),
  })
  if err then
    notify(err, vim.log.levels.ERROR)
    return false
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
    return false
  end
  if status ~= 200 or not decoded then
    notify("push failed: " .. tostring(status), vim.log.levels.ERROR)
    return false
  end
  apply_live_state(decoded, true)
  notify("pushed " .. M.state.item_id)
  return true
end

stop_poll_loop = function()
  if M.state.poll_timer then
    M.state.poll_timer:stop()
    M.state.poll_timer:close()
    M.state.poll_timer = nil
  end
end

local function stop_heartbeat_loop()
  if M.state.heartbeat_timer then
    M.state.heartbeat_timer:stop()
    M.state.heartbeat_timer:close()
    M.state.heartbeat_timer = nil
  end
end

-- Intent: Keep the first Neovim phase close to the browser draft studio by
-- polling the same runtime state and refreshing presence on a short interval.
-- Source: DI-fudok
start_poll_loop = function()
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

local function start_heartbeat_loop()
  stop_heartbeat_loop()
  M.state.heartbeat_timer = uv.new_timer()
  M.state.heartbeat_timer:start(M.config.poll_ms, M.config.poll_ms, function()
    vim.schedule(function()
      if not M.state.item_id or not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
        stop_heartbeat_loop()
        return
      end
      push_presence(false)
    end)
  end)
end

local function start_live_transport()
  stop_poll_loop()
  start_heartbeat_loop()
  -- Intent: Make Neovim prefer the direct Unix-socket live-draft contract
  -- while preserving HTTP polling fallback when the local socket surface is
  -- unavailable. Source: DI-favel
  M.connect_live_socket(M.state.item_id)
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
    "socket_path: " .. (M.config.socket_path or "-"),
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

local function link_summary(link)
  return string.format("- %s:%s -> %s:%s relation=%s", link.from_type or "-", link.from_id or "-", link.to_type or "-", link.to_id or "-", link.relation or "-")
end

local function append_links(lines, links)
  table.insert(lines, "")
  table.insert(lines, "## Links")
  if not links or #links == 0 then
    table.insert(lines, "- none")
    return
  end
  for _, link in ipairs(links) do
    table.insert(lines, link_summary(link))
    if link.notes and link.notes ~= "" then
      table.insert(lines, "  " .. link.notes)
    end
  end
end

local function append_inspect_handoff(lines, label, command)
  table.insert(lines, "- " .. label .. ": " .. command)
end

-- Intent: Keep terminal run review flowing into the surrounding operational
-- context instead of stranding the user at one run record after triage.
-- Source: DI-vunep
local function append_run_context_handoffs(lines, detail)
  table.insert(lines, "")
  table.insert(lines, "## Handoffs")
  append_inspect_handoff(lines, "item", ":OksInspect " .. (detail.item_id or "-"))
  if detail.place_id and detail.place_id ~= "" then
    append_inspect_handoff(lines, "place", ":OksInspectEntity place " .. detail.place_id)
  end
  if detail.resource_ids and #detail.resource_ids > 0 then
    for _, resource_id in ipairs(detail.resource_ids) do
      if resource_id and resource_id ~= "" then
        append_inspect_handoff(lines, "resource", ":OksInspectEntity resource " .. resource_id)
      end
    end
  end
  if detail.responsibility_ids and #detail.responsibility_ids > 0 then
    for _, responsibility_id in ipairs(detail.responsibility_ids) do
      if responsibility_id and responsibility_id ~= "" then
        append_inspect_handoff(lines, "responsibility", ":OksInspectEntity responsibility " .. responsibility_id)
      end
    end
  end
end

local function append_filter_summary(lines, filters)
  table.insert(lines, "")
  table.insert(lines, "## Filters")
  local wrote = false
  for _, key in ipairs({ "q", "kind", "status", "outcome", "place_id", "resource_id", "responsibility_id", "problem" }) do
    local value = filters and filters[key]
    if value ~= nil and value ~= "" and value ~= false then
      table.insert(lines, "- " .. key .. ": " .. tostring(value))
      wrote = true
    end
  end
  if not wrote then
    table.insert(lines, "- none")
  end
end

local function append_search_section(lines, heading, results, render)
  table.insert(lines, "")
  table.insert(lines, "## " .. heading)
  if not results or #results == 0 then
    table.insert(lines, "- none")
    return
  end
  for _, result in ipairs(results) do
    for _, line in ipairs(render(result)) do
      table.insert(lines, line)
    end
  end
end

local function search_result_lines(query, decoded)
  local lines = {
    "# Search",
    "",
    "- query: " .. ((query and query ~= "") and query or "(filters only)"),
  }

  append_filter_summary(lines, decoded.filters)

  append_search_section(lines, "Places", decoded.places, function(place)
    local out = {
      string.format("- %s kind=%s name=%s", place.id or "-", place.kind or "-", place.name or "-"),
      "  inspect: :OksInspectEntity place " .. (place.id or "-"),
    }
    if place.summary and place.summary ~= "" then
      table.insert(out, 2, "  " .. place.summary)
    end
    return out
  end)

  append_search_section(lines, "Resources", decoded.resources, function(resource)
    local out = {
      string.format("- %s kind=%s name=%s", resource.id or "-", resource.kind or "-", resource.name or "-"),
      "  inspect: :OksInspectEntity resource " .. (resource.id or "-"),
    }
    if resource.summary and resource.summary ~= "" then
      table.insert(out, 2, "  " .. resource.summary)
    end
    if resource.place_id and resource.place_id ~= "" then
      table.insert(out, "  place: " .. resource.place_id)
    end
    return out
  end)

  append_search_section(lines, "Responsibilities", decoded.responsibilities, function(responsibility)
    local out = {
      string.format("- %s title=%s", responsibility.id or "-", responsibility.title or "-"),
      "  inspect: :OksInspectEntity responsibility " .. (responsibility.id or "-"),
    }
    if responsibility.summary and responsibility.summary ~= "" then
      table.insert(out, 2, "  " .. responsibility.summary)
    end
    return out
  end)

  append_search_section(lines, "Items", decoded.items, function(item)
    local out = {
      string.format("- %s kind=%s status=%s title=%s", item.id or "-", item.kind or "-", item.status or "-", item.title or "-"),
      "  inspect: :OksInspect " .. (item.id or "-"),
    }
    if item.summary and item.summary ~= "" then
      table.insert(out, 2, "  " .. item.summary)
    end
    return out
  end)

  append_search_section(lines, "Runs", decoded.runs, function(run)
    local out = {
      string.format("- %s kind=%s outcome=%s item=%s", run.id or "-", run.kind or "-", run.outcome or "-", run.item_id or "-"),
      "  inspect: :OksInspectRun " .. (run.id or "-"),
    }
    if run.notes and run.notes ~= "" then
      table.insert(out, 2, "  " .. run.notes)
    end
    return out
  end)

  return lines
end

-- Intent: Keep Neovim search on the shared `/api/search` contract while
-- adding the same structured-filter vocabulary the CLI and browser already
-- use, instead of inventing a second editor-only search shape. Source:
-- DI-fanub
local function build_search_request(command_args)
  local tokens = vim.split(vim.trim(command_args or ""), "%s+", { trimempty = true })
  if #tokens == 0 then
    return nil, nil, "search query or filter is required"
  end

  local query_tokens = {}
  local filters = {}
  local saw_filter = false

  for _, token in ipairs(tokens) do
    local key, value = token:match("^([%a_][%w_]*)=(.*)$")
    if key ~= nil then
      saw_filter = true
      if not allowed_search_filters[key] then
        return nil, nil, 'unsupported search filter "' .. key .. '"'
      end
      filters[key] = value
    else
      if saw_filter then
        return nil, nil, "free-text query words must come before key=value filters"
      end
      table.insert(query_tokens, token)
    end
  end

  local query = table.concat(query_tokens, " ")
  if query == "" and next(filters) == nil then
    return nil, nil, "search query or filter is required"
  end

  local params = {}
  if query ~= "" then
    table.insert(params, "q=" .. url_encode(query))
  end
  for _, key in ipairs({ "kind", "status", "outcome", "place_id", "resource_id", "responsibility_id", "problem" }) do
    local value = filters[key]
    if value ~= nil and value ~= "" then
      table.insert(params, key .. "=" .. url_encode(value))
    end
  end

  local path = "/api/search"
  if #params > 0 then
    path = path .. "?" .. table.concat(params, "&")
  end
  return path, query, nil
end

local function search_projection(path, failure_prefix)
  local status, body, err = request("GET", path)
  if err then
    notify(err, vim.log.levels.ERROR)
    return nil
  end
  if status ~= 200 then
    notify(failure_prefix .. ": " .. tostring(status), vim.log.levels.ERROR)
    return nil
  end
  local decoded = json_decode(body)
  if not decoded then
    notify(failure_prefix .. " decode failed", vim.log.levels.ERROR)
    return nil
  end
  return decoded
end

local function pending_review_lines(draft_items, all_runs, problem_runs)
  local lines = {
    "# Pending review",
    "",
    "- inspect next from draft items, unreviewed runs, or problem runs",
  }

  table.insert(lines, "")
  table.insert(lines, "## Draft items")
  if not draft_items or #draft_items == 0 then
    table.insert(lines, "- none")
  else
    for _, item in ipairs(draft_items) do
      table.insert(lines, string.format("- %s kind=%s status=%s title=%s", item.id or "-", item.kind or "-", item.status or "-", item.title or "-"))
      table.insert(lines, "  inspect: :OksInspect " .. (item.id or "-"))
      if item.summary and item.summary ~= "" then
        table.insert(lines, "  " .. item.summary)
      end
    end
  end

  table.insert(lines, "")
  table.insert(lines, "## Unreviewed runs")
  local unreviewed = {}
  for _, run in ipairs(all_runs or {}) do
    if not run.approvals or #run.approvals == 0 then
      table.insert(unreviewed, run)
    end
  end
  if #unreviewed == 0 then
    table.insert(lines, "- none")
  else
    for _, run in ipairs(unreviewed) do
      table.insert(lines, string.format("- %s kind=%s outcome=%s item=%s", run.id or "-", run.kind or "-", run.outcome or "-", run.item_id or "-"))
      table.insert(lines, "  inspect: :OksInspectRun " .. (run.id or "-"))
      if run.notes and run.notes ~= "" then
        table.insert(lines, "  " .. run.notes)
      end
    end
  end

  table.insert(lines, "")
  table.insert(lines, "## Problem runs")
  if not problem_runs or #problem_runs == 0 then
    table.insert(lines, "- none")
  else
    for _, run in ipairs(problem_runs) do
      table.insert(lines, string.format("- %s kind=%s outcome=%s item=%s", run.id or "-", run.kind or "-", run.outcome or "-", run.item_id or "-"))
      table.insert(lines, "  inspect: :OksInspectRun " .. (run.id or "-"))
      if run.notes and run.notes ~= "" then
        table.insert(lines, "  " .. run.notes)
      end
    end
  end

  return lines
end

local function append_problem_group_section(lines, heading, group_type, groups)
  table.insert(lines, "")
  table.insert(lines, "## " .. heading)
  if not groups or #groups == 0 then
    table.insert(lines, "- none")
    return
  end
  for _, group in ipairs(groups) do
    table.insert(lines, string.format(
      "- %s kind=%s name=%s problems=%d receiving=%d inventory=%d",
      group.group_id or "-",
      group.kind or "-",
      group.name or "-",
      group.problem_count or 0,
      group.receiving_problems or 0,
      group.inventory_problems or 0
    ))
    table.insert(lines, "  inspect: :OksInspectEntity " .. group_type .. " " .. (group.group_id or "-"))
    if group.highlights and #group.highlights > 0 then
      table.insert(lines, "  highlights: " .. table.concat(group.highlights, " | "))
    end
    table.insert(lines, "  runs:")
    if not group.runs or #group.runs == 0 then
      table.insert(lines, "  - none")
    else
      for _, run in ipairs(group.runs) do
        table.insert(lines, string.format("  - %s kind=%s outcome=%s item=%s", run.id or "-", run.kind or "-", run.outcome or "-", run.item_id or "-"))
        table.insert(lines, "    inspect: :OksInspectRun " .. (run.id or "-"))
        if run.notes and run.notes ~= "" then
          table.insert(lines, "    " .. run.notes)
        end
        if run.resource_ids and #run.resource_ids > 0 then
          table.insert(lines, "    resources: " .. table.concat(run.resource_ids, ", "))
        end
      end
    end
  end
end

local function problem_review_lines(review)
  local lines = {
    "# Problem review",
    "",
    "- grouped receiving and inventory hotspots over /api/problem-review",
    "- problem runs: " .. tostring(review.problem_runs or 0),
  }

  append_problem_group_section(lines, "Place groups", "place", review.place_groups)
  append_problem_group_section(lines, "Resource groups", "resource", review.resource_groups)

  return lines
end

-- Intent: Treat missing approvals arrays in shared pending-review run payloads
-- as contract drift so terminal review queues fail loudly instead of inventing
-- fake unreviewed work from omitted fields. Source: DI-davur
local function validate_pending_runs(runs, label)
  if runs == nil then
    return true
  end
  if type(runs) ~= "table" or not vim.tbl_islist(runs) then
    notify(label .. " is not an array", vim.log.levels.ERROR)
    return false
  end
  for index, run in ipairs(runs) do
    if type(run) ~= "table" then
      notify(string.format("%s[%d] is not an object", label, index), vim.log.levels.ERROR)
      return false
    end
    if run.approvals == nil then
      notify(string.format('%s[%d] missing "approvals" array', label, index), vim.log.levels.ERROR)
      return false
    end
    if type(run.approvals) ~= "table" or not vim.tbl_islist(run.approvals) then
      notify(string.format('%s[%d] "approvals" field is not an array', label, index), vim.log.levels.ERROR)
      return false
    end
  end
  return true
end

-- Intent: Keep the Neovim grouped problem-review view aligned with the shared
-- hotspot projection by rejecting malformed group payloads instead of silently
-- dropping grouped review context. Source: DI-sivok
local function validate_problem_groups(groups, label)
  if groups == nil then
    notify(label .. ' missing array', vim.log.levels.ERROR)
    return false
  end
  if type(groups) ~= "table" or not vim.tbl_islist(groups) then
    notify(label .. " is not an array", vim.log.levels.ERROR)
    return false
  end
  for index, group in ipairs(groups) do
    if type(group) ~= "table" then
      notify(string.format("%s[%d] is not an object", label, index), vim.log.levels.ERROR)
      return false
    end
    if group.runs == nil then
      notify(string.format('%s[%d] missing "runs" array', label, index), vim.log.levels.ERROR)
      return false
    end
    if type(group.runs) ~= "table" or not vim.tbl_islist(group.runs) then
      notify(string.format('%s[%d] "runs" field is not an array', label, index), vim.log.levels.ERROR)
      return false
    end
    if group.highlights ~= nil and (type(group.highlights) ~= "table" or not vim.tbl_islist(group.highlights)) then
      notify(string.format('%s[%d] "highlights" field is not an array', label, index), vim.log.levels.ERROR)
      return false
    end
  end
  return true
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
    -- Intent: Make older item inspectors as actionable as the newer search and
    -- pending views by emitting direct run-inspection handoffs inline with the
    -- projected related-run summaries. Source: DI-josav
    for _, run in ipairs(detail.related_runs) do
      table.insert(lines, string.format("- %s kind=%s revision=%d outcome=%s", run.id or "-", run.kind or "-", run.revision or 0, run.outcome or "-"))
      table.insert(lines, "  inspect: :OksInspectRun " .. (run.id or "-"))
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

  append_links(lines, detail.links)

  return lines
end

local function show_item_detail(item_id, detail)
  if inspector.winid and vim.api.nvim_win_is_valid(inspector.winid) and inspector.bufnr and vim.api.nvim_buf_is_valid(inspector.bufnr) then
    vim.api.nvim_set_current_win(inspector.winid)
  else
    inspector.bufnr, inspector.winid = open_scratch_buffer("oks-inspect://" .. item_id)
  end
  vim.bo[inspector.bufnr].modifiable = true
  vim.api.nvim_buf_set_name(inspector.bufnr, "oks-inspect://" .. item_id)
  vim.api.nvim_buf_set_lines(inspector.bufnr, 0, -1, false, item_detail_lines(detail))
  vim.bo[inspector.bufnr].modifiable = false
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
  show_item_detail(item_id, decoded)
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

  append_run_context_handoffs(lines, detail)
  append_links(lines, detail.links)

  return lines
end

local function show_run_detail(run_id, detail)
  if inspector.winid and vim.api.nvim_win_is_valid(inspector.winid) and inspector.bufnr and vim.api.nvim_buf_is_valid(inspector.bufnr) then
    vim.api.nvim_set_current_win(inspector.winid)
  else
    inspector.bufnr, inspector.winid = open_scratch_buffer("oks-run://" .. run_id)
  end
  vim.bo[inspector.bufnr].modifiable = true
  vim.api.nvim_buf_set_name(inspector.bufnr, "oks-run://" .. run_id)
  vim.api.nvim_buf_set_lines(inspector.bufnr, 0, -1, false, run_detail_lines(detail))
  vim.bo[inspector.bufnr].modifiable = false
end

local function generic_entity_lines(entity_type, detail)
  if entity_type == "place" then
    local lines = {
      "# " .. (detail.name or detail.id or "place"),
      "",
      "- id: " .. (detail.id or "-"),
      "- kind: " .. (detail.kind or "-"),
      "- summary: " .. (detail.summary or ""),
      "- parent_id: " .. (detail.parent_id or "-"),
    }
    append_list(lines, "Child places", detail.child_place_ids or {})
    append_list(lines, "Resources", detail.resource_ids or {})
    if detail.related_runs then
      table.insert(lines, "")
      table.insert(lines, "## Related runs")
      if #detail.related_runs == 0 then
        table.insert(lines, "- none")
      else
        -- Intent: Keep entity inspectors consistent with terminal queue and
        -- run-inspector navigation by giving every related-run summary a
        -- direct :OksInspectRun handoff. Source: DI-josav
        for _, run in ipairs(detail.related_runs) do
          table.insert(lines, string.format("- %s kind=%s outcome=%s", run.id or "-", run.kind or "-", run.outcome or "-"))
          table.insert(lines, "  inspect: :OksInspectRun " .. (run.id or "-"))
        end
      end
    end
    append_links(lines, detail.links)
    return lines
  end
  if entity_type == "resource" then
    local lines = {
      "# " .. (detail.name or detail.id or "resource"),
      "",
      "- id: " .. (detail.id or "-"),
      "- kind: " .. (detail.kind or "-"),
      "- summary: " .. (detail.summary or ""),
      "- place_id: " .. (detail.place_id or "-"),
    }
    if detail.related_runs then
      table.insert(lines, "")
      table.insert(lines, "## Related runs")
      if #detail.related_runs == 0 then
        table.insert(lines, "- none")
      else
        -- Intent: Keep entity inspectors consistent with terminal queue and
        -- run-inspector navigation by giving every related-run summary a
        -- direct :OksInspectRun handoff. Source: DI-josav
        for _, run in ipairs(detail.related_runs) do
          table.insert(lines, string.format("- %s kind=%s outcome=%s", run.id or "-", run.kind or "-", run.outcome or "-"))
          table.insert(lines, "  inspect: :OksInspectRun " .. (run.id or "-"))
        end
      end
    end
    append_links(lines, detail.links)
    return lines
  end
  if entity_type == "responsibility" then
    local lines = {
      "# " .. (detail.title or detail.id or "responsibility"),
      "",
      "- id: " .. (detail.id or "-"),
      "- summary: " .. (detail.summary or ""),
      "- team: " .. (detail.team or "-"),
    }
    append_list(lines, "Linked items", detail.linked_item_ids or {})
    append_list(lines, "Linked runs", detail.linked_run_ids or {})
    append_list(lines, "Roles", detail.linked_role_keys or {})
    if detail.related_runs then
      table.insert(lines, "")
      table.insert(lines, "## Related runs")
      if #detail.related_runs == 0 then
        table.insert(lines, "- none")
      else
        -- Intent: Keep entity inspectors consistent with terminal queue and
        -- run-inspector navigation by giving every related-run summary a
        -- direct :OksInspectRun handoff. Source: DI-josav
        for _, run in ipairs(detail.related_runs) do
          table.insert(lines, string.format("- %s kind=%s outcome=%s", run.id or "-", run.kind or "-", run.outcome or "-"))
          table.insert(lines, "  inspect: :OksInspectRun " .. (run.id or "-"))
        end
      end
    end
    append_links(lines, detail.links)
    return lines
  end
  return {
    "# " .. (detail.id or entity_type),
    "",
    "unsupported entity type: " .. entity_type,
  }
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
  show_run_detail(run_id, decoded)
  notify("inspected run " .. run_id)
end

local function inspect_buffer(name, lines)
  if inspector.winid and vim.api.nvim_win_is_valid(inspector.winid) and inspector.bufnr and vim.api.nvim_buf_is_valid(inspector.bufnr) then
    vim.api.nvim_set_current_win(inspector.winid)
  else
    inspector.bufnr, inspector.winid = open_scratch_buffer(name)
  end
  vim.bo[inspector.bufnr].modifiable = true
  vim.api.nvim_buf_set_name(inspector.bufnr, name)
  vim.api.nvim_buf_set_lines(inspector.bufnr, 0, -1, false, lines)
  vim.bo[inspector.bufnr].modifiable = false
end

-- Intent: Let Neovim users browse linked operational entities through the same
-- projected detail APIs as the browser and CLI, without adding mutation
-- actions to the editor yet. Source: DI-zalor
local function inspect_entity(entity_type, entity_id)
  entity_type = vim.trim(entity_type or "")
  entity_id = vim.trim(entity_id or "")
  if entity_type == "" or entity_id == "" then
    notify("entity type and id are required", vim.log.levels.WARN)
    return
  end

  local path
  local name_prefix
  if entity_type == "item" then
    path = "/api/items/" .. entity_id
    name_prefix = "oks-inspect://"
  elseif entity_type == "run" then
    path = "/api/runs/" .. entity_id
    name_prefix = "oks-run://"
  elseif entity_type == "place" then
    path = "/api/places/" .. entity_id
    name_prefix = "oks-place://"
  elseif entity_type == "resource" then
    path = "/api/resources/" .. entity_id
    name_prefix = "oks-resource://"
  elseif entity_type == "responsibility" then
    path = "/api/responsibilities/" .. entity_id
    name_prefix = "oks-responsibility://"
  else
    notify("unsupported entity type: " .. entity_type, vim.log.levels.WARN)
    return
  end

  local status, body, err = request("GET", path)
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    notify("inspect entity failed: " .. tostring(status), vim.log.levels.ERROR)
    return
  end
  local decoded = json_decode(body)
  if not decoded then
    notify("inspect entity decode failed", vim.log.levels.ERROR)
    return
  end

  local lines
  if entity_type == "item" then
    lines = item_detail_lines(decoded)
  elseif entity_type == "run" then
    lines = run_detail_lines(decoded)
  else
    lines = generic_entity_lines(entity_type, decoded)
  end
  inspect_buffer(name_prefix .. entity_id, lines)
  notify("inspected " .. entity_type .. " " .. entity_id)
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

function M.inspect_entity(entity_type, entity_id)
  inspect_entity(entity_type, entity_id)
end

-- Intent: Let Neovim users discover operational records through the same
-- search projection as the browser and CLI, while keeping the editor-side
-- surface read-only and routed into the existing inspectors for deeper
-- browsing. Source: DI-givot
function M.search(command_args)
  local path, query, build_err = build_search_request(command_args)
  if build_err then
    notify(build_err, vim.log.levels.WARN)
    return
  end

  local status, body, err = request("GET", path)
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    notify("search failed: " .. tostring(status), vim.log.levels.ERROR)
    return
  end
  local decoded = json_decode(body)
  if not decoded then
    notify("search decode failed", vim.log.levels.ERROR)
    return
  end
  local suffix = query ~= "" and url_encode(query) or "filters"
  inspect_buffer("oks-search://" .. suffix, search_result_lines(query, decoded))
  notify("searched " .. ((query ~= "") and query or "filters"))
end

-- Intent: Give terminal users a compact “what needs attention next” surface
-- by reusing the existing search projections for draft items and review-worthy
-- runs before adding write-side approval actions in Neovim. Source: DI-lorav
function M.pending()
  local draft = search_projection("/api/search?status=draft", "pending draft search failed")
  if not draft then
    return
  end
  local all = search_projection("/api/search", "pending run search failed")
  if not all then
    return
  end
  local problems = search_projection("/api/search?problem=true", "pending problem search failed")
  if not problems then
    return
  end
  if not validate_pending_runs(all.runs, "/api/search runs") then
    return
  end
  if not validate_pending_runs(problems.runs, "/api/search?problem=true runs") then
    return
  end
  inspect_buffer("oks-pending://review", pending_review_lines(draft.items or {}, all.runs or {}, problems.runs or {}))
  notify("opened pending review")
end

-- Intent: Close the biggest remaining terminal-side grouped-review gap by
-- giving Neovim the same hotspot projection the browser and CLI already use,
-- without adding a separate editor-only review API. Source: DI-sivok
function M.problem_review()
  local review = search_projection("/api/problem-review", "problem review failed")
  if not review then
    return
  end
  if not validate_problem_groups(review.place_groups, "/api/problem-review place_groups") then
    return
  end
  if not validate_problem_groups(review.resource_groups, "/api/problem-review resource_groups") then
    return
  end
  inspect_buffer("oks-problem-review://review", problem_review_lines(review))
  notify("opened problem review")
end

local function refresh_after_item_approval(item_id, detail)
  local current_name = buffer_name(vim.api.nvim_get_current_buf())
  local inspector_name = buffer_name(inspector.bufnr)
  if current_name == "oks-pending://review" or inspector_name == "oks-pending://review" then
    M.pending()
    return
  end
  if M.state.item_id == item_id and M.state.bufnr and vim.api.nvim_buf_is_valid(M.state.bufnr) then
    M.state.title = detail.title or M.state.title
    M.state.status = detail.status or M.state.status
    M.state.current_revision = detail.current_revision or M.state.current_revision
    refresh_live_state(true)
    return
  end
  show_item_detail(item_id, detail)
end

local function refresh_after_item_supersede(item_id, detail)
  local current_name = buffer_name(vim.api.nvim_get_current_buf())
  local inspector_name = buffer_name(inspector.bufnr)
  if current_name == "oks-pending://review" or inspector_name == "oks-pending://review" then
    M.pending()
    return
  end
  if M.state.item_id == item_id and M.state.bufnr and vim.api.nvim_buf_is_valid(M.state.bufnr) then
    M.state.title = detail.title or M.state.title
    M.state.status = detail.status or M.state.status
    M.state.current_revision = detail.current_revision or M.state.current_revision
    refresh_live_state(true)
    return
  end
  show_item_detail(item_id, detail)
end

-- Intent: Advance the durable revision metadata after a successful snapshot
-- without clobbering the already-pushed live draft body, while still
-- refreshing any visible item inspector state. Source: DI-jabup
local function refresh_after_revision_snapshot(item_id, detail)
  local current_name = buffer_name(vim.api.nvim_get_current_buf())
  local inspector_name = buffer_name(inspector.bufnr)
  if current_name == "oks-pending://review" or inspector_name == "oks-pending://review" then
    M.pending()
    return
  end
  if M.state.item_id == item_id and M.state.bufnr and vim.api.nvim_buf_is_valid(M.state.bufnr) then
    M.state.title = detail.title or M.state.title
    M.state.status = detail.status or M.state.status
    M.state.current_revision = detail.current_revision or M.state.current_revision
    vim.bo[M.state.bufnr].modified = false
  end
  if current_name == "oks-inspect://" .. item_id or inspector_name == "oks-inspect://" .. item_id then
    show_item_detail(item_id, detail)
  end
end

local function refresh_after_run_approval(run_id, detail)
  local current_name = buffer_name(vim.api.nvim_get_current_buf())
  local inspector_name = buffer_name(inspector.bufnr)
  if current_name == "oks-pending://review" or inspector_name == "oks-pending://review" then
    M.pending()
    return
  end
  show_run_detail(run_id, detail)
end

-- Intent: Let terminal-first reviewers approve the current item from Neovim
-- through the existing revision-aware item approval API, while refreshing the
-- relevant live, inspect, or pending-review surface afterward. Source: DI-vamor
function M.approve_item(command_args)
  local args = vim.split(vim.trim(command_args or ""), "%s+", { trimempty = true })
  local item_id = current_context_item_id()
  local role
  local decision
  local notes = ""

  if item_id and #args >= 2 then
    role = args[1]
    decision = args[2]
    notes = table.concat(args, " ", 3)
  elseif #args >= 3 then
    item_id = args[1]
    role = args[2]
    decision = args[3]
    notes = table.concat(args, " ", 4)
  else
    notify("usage: :OksApproveItem [ITEM_ID] ROLE DECISION [NOTES...]", vim.log.levels.WARN)
    return
  end

  if not valid_decision(decision) then
    notify("decision must be approved, rejected, or noted", vim.log.levels.WARN)
    return
  end

  local status, body, err = request("GET", "/api/items/" .. item_id)
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    notify("load item for approval failed: " .. tostring(status), vim.log.levels.ERROR)
    return
  end
  local item = json_decode(body)
  if not item then
    notify("item approval preflight decode failed", vim.log.levels.ERROR)
    return
  end

  status, body, err = request("POST", "/api/items/" .. item_id .. "/approvals", {
    actor = M.config.display_name,
    revision = item.current_revision,
    role = role,
    decision = decision,
    notes = notes,
  })
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    local message = vim.trim(body or "")
    if message == "" then
      message = tostring(status)
    end
    notify("item approval failed: " .. message, vim.log.levels.ERROR)
    return
  end
  local approved = json_decode(body)
  if not approved then
    notify("item approval decode failed", vim.log.levels.ERROR)
    return
  end
  refresh_after_item_approval(item_id, approved)
  notify("approved item " .. item_id .. " as " .. decision)
end

-- Intent: Let terminal-first reviewers approve the current run from Neovim
-- through the existing run approval API, while refreshing the relevant
-- inspect or pending-review surface afterward. Source: DI-bafor
function M.approve_run(command_args)
  local args = vim.split(vim.trim(command_args or ""), "%s+", { trimempty = true })
  local run_id = current_context_run_id()
  local role
  local decision
  local notes = ""

  if run_id and #args >= 2 then
    role = args[1]
    decision = args[2]
    notes = table.concat(args, " ", 3)
  elseif #args >= 3 then
    run_id = args[1]
    role = args[2]
    decision = args[3]
    notes = table.concat(args, " ", 4)
  else
    notify("usage: :OksApproveRun [RUN_ID] ROLE DECISION [NOTES...]", vim.log.levels.WARN)
    return
  end

  if not valid_decision(decision) then
    notify("decision must be approved, rejected, or noted", vim.log.levels.WARN)
    return
  end

  local status, body, err = request("POST", "/api/runs/" .. run_id .. "/approvals", {
    actor = M.config.display_name,
    revision = 0,
    role = role,
    decision = decision,
    notes = notes,
  })
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    local message = vim.trim(body or "")
    if message == "" then
      message = tostring(status)
    end
    notify("run approval failed: " .. message, vim.log.levels.ERROR)
    return
  end
  local approved = json_decode(body)
  if not approved then
    notify("run approval decode failed", vim.log.levels.ERROR)
    return
  end
  refresh_after_run_approval(run_id, approved)
  notify("approved run " .. run_id .. " as " .. decision)
end

-- Intent: Let terminal-first reviewers supersede the current item from Neovim
-- through the existing item supersede API, while refreshing the relevant
-- live, inspect, or pending-review surface afterward. Source: DI-pudor
function M.supersede_item(command_args)
  local args = vim.split(vim.trim(command_args or ""), "%s+", { trimempty = true })
  local item_id = current_context_item_id()
  local notes = ""

  if item_id and #args >= 0 then
    notes = table.concat(args, " ")
  elseif #args >= 1 then
    item_id = args[1]
    notes = table.concat(args, " ", 2)
  else
    notify("usage: :OksSupersedeItem [ITEM_ID] [NOTES...]", vim.log.levels.WARN)
    return
  end

  if not item_id or item_id == "" then
    notify("usage: :OksSupersedeItem [ITEM_ID] [NOTES...]", vim.log.levels.WARN)
    return
  end

  local status, body, err = request("POST", "/api/items/" .. item_id .. "/supersede", {
    actor = M.config.display_name,
    notes = notes,
  })
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    local message = vim.trim(body or "")
    if message == "" then
      message = tostring(status)
    end
    notify("item supersede failed: " .. message, vim.log.levels.ERROR)
    return
  end
  local superseded = json_decode(body)
  if not superseded then
    notify("item supersede decode failed", vim.log.levels.ERROR)
    return
  end
  refresh_after_item_supersede(item_id, superseded)
  notify("superseded item " .. item_id)
end

-- Intent: Let terminal-first authors cut a durable revision from the current
-- live-draft session through the existing item revision API, while keeping the
-- local live-draft and optional inspector state in sync afterward. Source:
-- DI-jabup
function M.snapshot()
  if not M.state.item_id or not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    notify("open a live draft before snapshotting", vim.log.levels.WARN)
    return
  end

  if not push_body() then
    return
  end

  local item_id = M.state.item_id
  local status, body, err = request("GET", "/api/items/" .. item_id)
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    notify("load item for snapshot failed: " .. tostring(status), vim.log.levels.ERROR)
    return
  end
  local item = json_decode(body)
  if not item then
    notify("snapshot preflight decode failed", vim.log.levels.ERROR)
    return
  end

  status, body, err = request("POST", "/api/items/" .. item_id .. "/revisions", {
    actor = M.config.display_name,
    title = item.title,
    summary = item.summary,
    body = current_body(),
    tags = item.tags or {},
  })
  if err then
    notify(err, vim.log.levels.ERROR)
    return
  end
  if status ~= 200 then
    local message = vim.trim(body or "")
    if message == "" then
      message = tostring(status)
    end
    notify("snapshot failed: " .. message, vim.log.levels.ERROR)
    return
  end
  local updated = json_decode(body)
  if not updated then
    notify("snapshot decode failed", vim.log.levels.ERROR)
    return
  end
  refresh_after_revision_snapshot(item_id, updated)
  notify("snapshot created as revision " .. tostring(updated.current_revision or "?"))
end

-- Intent: Let the headless regression harness seed inspector state without
-- going through HTTP-backed inspect flows, so :OksClose teardown can be
-- validated locally. Source: DI-mabek
function M._test_set_inspector(bufnr, winid)
  inspector.bufnr = bufnr
  inspector.winid = winid
end

function M.close()
  stop_poll_loop()
  stop_heartbeat_loop()
  clear_reconnect_loop()
  close_live_socket()
  if M.state.augroup then
    pcall(vim.api.nvim_del_augroup_by_id, M.state.augroup)
    M.state.augroup = nil
  end
  -- Intent: Make :OksClose behave like a real session teardown by wiping the
  -- live-draft and inspector buffers instead of leaving detached scratch state
  -- visible after the session hooks are gone. Source: DI-mabek
  wipe_buffer(inspector.bufnr)
  wipe_buffer(M.state.bufnr)
  M.state.item_id = nil
  M.state.bufnr = nil
  M.state.winid = nil
  M.state.title = ""
  M.state.status = ""
  M.state.version = 0
  M.state.current_revision = 0
  M.state.participants = {}
  inspector.bufnr = nil
  inspector.winid = nil
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
  M.state.winid = vim.api.nvim_get_current_win()
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
  start_live_transport()
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
  vim.api.nvim_create_user_command("OksInspectEntity", function(command)
    local args = vim.split(vim.trim(command.args), "%s+", { trimempty = true })
    M.inspect_entity(args[1], args[2])
  end, { nargs = "+" })
  -- Intent: Let the shipped :OksSearch command accept the documented
  -- multi-token query-plus-filter syntax instead of truncating input at a
  -- single Ex argument before the shared search parser runs. Source: DI-lavup
  vim.api.nvim_create_user_command("OksSearch", function(command)
    M.search(command.args)
  end, { nargs = "+" })
  vim.api.nvim_create_user_command("OksPending", function()
    M.pending()
  end, {})
  vim.api.nvim_create_user_command("OksProblemReview", function()
    M.problem_review()
  end, {})
  vim.api.nvim_create_user_command("OksSnapshot", function()
    M.snapshot()
  end, {})
  vim.api.nvim_create_user_command("OksApproveItem", function(command)
    M.approve_item(command.args)
  end, { nargs = "+" })
  vim.api.nvim_create_user_command("OksApproveRun", function(command)
    M.approve_run(command.args)
  end, { nargs = "+" })
  vim.api.nvim_create_user_command("OksSupersedeItem", function(command)
    M.supersede_item(command.args)
  end, { nargs = "*" })
  vim.api.nvim_create_user_command("OksClose", function()
    M.close()
  end, {})
end

return M
