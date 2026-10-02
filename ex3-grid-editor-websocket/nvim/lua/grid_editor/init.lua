local M = {}

local function default_participant_id()
  local pid = tostring(vim.fn.getpid())
  local host = ((vim.uv or vim.loop).os_gethostname and (vim.uv or vim.loop).os_gethostname()) or 'host'
  host = host:gsub('[^%w%-_]', '-')
  local nonce = tostring(math.floor(((vim.uv or vim.loop).hrtime() or 0) % 1000000))
  -- Intent: Keep Neovim awareness participants distinct across machines and
  -- restarts so relay presence state and peer filtering do not collapse two
  -- separate editors onto the same participant ID. Source: DI-samuv; DI-gafit
  return table.concat({ 'nvim', host, pid, nonce }, '-')
end

local function default_repo_root()
  local source = debug.getinfo(1, 'S').source
  if vim.startswith(source, '@') then
    source = source:sub(2)
  end
  local path = vim.fs.normalize(source)
  return vim.fs.dirname(vim.fs.dirname(vim.fs.dirname(vim.fs.dirname(path))))
end

local function default_sidecar_cmd(repo_root)
  -- Intent: Make Neovim usable without a preinstalled sidecar binary by
  -- defaulting the embodiment-local launcher to `go run` inside the checked-out
  -- repo. Source: DI-samuv
  return { 'go', 'run', repo_root .. '/cmd/grid-nvim-sidecar' }
end

M.config = {
  repo_root = default_repo_root(),
  -- Intent: Keep the copied Neovim embodiment on ex3's relay default so it
  -- does not collide with ex2 during local parallel runs. Source: DI-vatub
  relay_url = vim.env.GRID_EDITOR_RELAY_URL or 'http://127.0.0.1:7025',
  access_token = vim.env.GRID_EDITOR_ACCESS_TOKEN or '',
  sidecar_cmd = nil,
  display_name = vim.env.GRID_EDITOR_DISPLAY_NAME or 'Neovim User',
  color = vim.env.GRID_EDITOR_COLOR or '#d66f1d',
  show_line_numbers = true,
  presence_profile = 'demo',
}

M.state = {
  bufnr = nil,
  doc_id = nil,
  suppress = false,
  expected_content = nil,
  cursor_ns = nil,
  selection_ns = nil,
  augroup = nil,
  job_id = nil,
  participant_id = default_participant_id(),
  peers = {},
  peer_index = {},
  relay_connected = false,
  session_ready = false,
  info_bufnr = nil,
  info_winid = nil,
  help_bufnr = nil,
  help_winid = nil,
  remote_cursors = {},
  remote_selections = {},
}

local function join_lines()
  if not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return ''
  end
  return table.concat(vim.api.nvim_buf_get_lines(M.state.bufnr, 0, -1, false), '\n')
end

local function offset_to_pos(lines, offset)
  local total = 0
  for i, line in ipairs(lines) do
    local next_total = total + #line
    if offset <= next_total then
      return i - 1, math.max(0, offset - total)
    end
    total = next_total + 1
  end
  local row = math.max(0, #lines - 1)
  return row, #(lines[#lines] or '')
end

local function pos_to_offset(lines, row, col)
  local total = 0
  for i = 1, row do
    total = total + #(lines[i] or '') + 1
  end
  return total + col
end

local function normalize_color(color)
  if type(color) == 'string' and color:match('^#%x%x%x%x%x%x$') then
    return color
  end
  return '#999999'
end

local function peer_highlight(color)
  local normalized = normalize_color(color)
  local group = 'GridEditorPeer' .. normalized:gsub('#', '')
  if vim.fn.hlexists(group) == 0 then
    vim.api.nvim_set_hl(0, group, {
      fg = '#ffffff',
      bg = normalized,
      bold = true,
    })
  end
  return group
end

local function peer_selection_highlight(color)
  local normalized = normalize_color(color)
  local group = 'GridEditorPeerSelection' .. normalized:gsub('#', '')
  if vim.fn.hlexists(group) == 0 then
    vim.api.nvim_set_hl(0, group, {
      bg = normalized,
      fg = '#ffffff',
      blend = 75,
    })
  end
  return group
end

local function set_cursor_extmark(bufnr, namespace, row, col, options)
  local ok = pcall(vim.api.nvim_buf_set_extmark, bufnr, namespace, row, col, options)
  if ok then
    return
  end
  local fallbacks = {
    function()
      local fallback = vim.deepcopy(options)
      fallback.virt_text_pos = 'eol'
      return fallback
    end,
    function()
      local fallback = vim.deepcopy(options)
      fallback.virt_text_pos = 'overlay'
      return fallback
    end,
  }
  for _, build in ipairs(fallbacks) do
    local fallback = build()
    if pcall(vim.api.nvim_buf_set_extmark, bufnr, namespace, row, col, fallback) then
      return
    end
  end
end

local peer_presence_state
local peer_visible_in_buffer

local function char_to_byte_col(line, char_col)
  if char_col <= 0 then
    return 0
  end
  local text = line or ''
  if char_col >= vim.fn.strchars(text) then
    return #text
  end
  return #vim.fn.strcharpart(text, 0, char_col)
end

local function clear_peer_marks(participant_id)
  if not M.state.cursor_ns then
    return
  end
  if M.state.remote_cursors[participant_id] then
    pcall(vim.api.nvim_buf_del_extmark, M.state.bufnr, M.state.cursor_ns, M.state.remote_cursors[participant_id])
    M.state.remote_cursors[participant_id] = nil
  end
  if M.state.remote_selections[participant_id] then
    pcall(vim.api.nvim_buf_del_extmark, M.state.bufnr, M.state.cursor_ns, M.state.remote_selections[participant_id])
    M.state.remote_selections[participant_id] = nil
  end
end

local function render_peer(peer, lines)
  local participant_id = peer.participant_id or peer.name or 'peer'
  clear_peer_marks(participant_id)

  local row, char_col = offset_to_pos(lines, peer.anchor or 0)
  local head_row, head_char_col = offset_to_pos(lines, peer.head or peer.anchor or 0)
  local line = lines[row + 1] or ''
  local cursor_byte_col = char_to_byte_col(line, char_col)
  local head_line = lines[head_row + 1] or ''
  local head_byte_col = char_to_byte_col(head_line, head_char_col)
  local group = peer_highlight(peer.color)
  local selection_group = peer_selection_highlight(peer.color)

  if row ~= head_row or char_col ~= head_char_col then
    local start_row, start_col = row, cursor_byte_col
    local end_row, end_col = head_row, head_byte_col
    if start_row > end_row or (start_row == end_row and start_col > end_col) then
      start_row, end_row = end_row, start_row
      start_col, end_col = end_col, start_col
    end
    M.state.remote_selections[participant_id] = vim.api.nvim_buf_set_extmark(M.state.bufnr, M.state.cursor_ns, start_row, start_col, {
      end_row = end_row,
      end_col = end_col,
      hl_group = selection_group,
      hl_mode = 'combine',
      priority = 200,
    })
  end

  local cursor_end_col = math.min(cursor_byte_col + 1, #line)
  -- Intent: Keep the remote cursor itself anchored at the exact document
  -- position, but move the peer label off the underlying text so the buffer
  -- remains readable while the cursor owner stays obvious. Source: DI-gafit;
  -- DI-samuv; DI-favok
  M.state.remote_cursors[participant_id] = vim.api.nvim_buf_set_extmark(M.state.bufnr, M.state.cursor_ns, row, cursor_byte_col, {
    virt_text = {
      { '▎ ' .. (peer.name or participant_id or 'peer') .. (peer.typing and ' typing' or '') .. ' ', group },
    },
    virt_text_pos = 'eol',
    hl_group = group,
    end_col = cursor_end_col,
    hl_mode = 'combine',
    priority = 300,
    sign_text = '▎',
    sign_hl_group = group,
  })
end

local function draw_peers(peers)
  if not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return
  end
  if not M.state.cursor_ns then
    M.state.cursor_ns = vim.api.nvim_create_namespace('grid_editor_peers')
  end
  local lines = vim.api.nvim_buf_get_lines(M.state.bufnr, 0, -1, false)
  local visible = {}
  for _, peer in ipairs(peers or {}) do
    if not peer_visible_in_buffer(peer) then
      goto continue
    end
    local participant_id = peer.participant_id or peer.name or 'peer'
    visible[participant_id] = true
    render_peer(peer, lines)
    ::continue::
  end
  for participant_id, _ in pairs(M.state.remote_cursors) do
    if not visible[participant_id] then
      clear_peer_marks(participant_id)
    end
  end
  for participant_id, _ in pairs(M.state.remote_selections) do
    if not visible[participant_id] then
      clear_peer_marks(participant_id)
    end
  end
end

peer_presence_state = function(peer)
  if type(peer.last_seen_at) ~= 'string' or peer.last_seen_at == '' then
    return 'live'
  end
  local observed = vim.fn.strptime('%Y-%m-%dT%H:%M:%SZ', peer.last_seen_at)
  if observed <= 0 then
    return 'live'
  end
  local age = os.time() - observed
  local profile = M.config.presence_profile == 'normal'
    and { live = 60, stale = 5 * 60, offline = 15 * 60 }
    or { live = 12, stale = 45, offline = 2 * 60 }
  if age <= profile.live then
    return 'live'
  end
  if age <= profile.stale then
    return 'stale'
  end
  if age <= profile.offline then
    return 'offline'
  end
  return 'gone'
end

peer_visible_in_buffer = function(peer)
  -- Intent: Only render currently live peers in-buffer so abandoned demo
  -- sessions do not drown out the active browser/editor cursors users expect
  -- to see. The full roster remains available in :GridEditorPeers. Source: DI-gafit
  return peer_presence_state(peer) == 'live'
end

local function sidecar_argv()
  local base = M.config.sidecar_cmd
  if base == nil then
    base = default_sidecar_cmd(M.config.repo_root)
  end
  if type(base) == 'string' then
    return { base, '--relay', M.config.relay_url }
  end
  local argv = vim.deepcopy(base)
  table.insert(argv, '--relay')
  table.insert(argv, M.config.relay_url)
  return argv
end

local function session_lines()
  -- Intent: Expose relay, doc, peer, and sidecar state from the live awareness
  -- session so Neovim users have an explicit place to inspect connection state.
  -- Source: DI-samuv
  local lines = {
    'grid-editor',
    '',
    'relay: ' .. M.config.relay_url,
    'doc: ' .. (M.state.doc_id or 'none'),
    'participant: ' .. M.state.participant_id,
    'display name: ' .. M.config.display_name,
    'color: ' .. M.config.color,
    'remote access token: ' .. (M.config.access_token ~= '' and '(configured)' or '(none)'),
    'relay status: ' .. (M.state.relay_connected and 'connected' or 'disconnected'),
    'sidecar: ' .. table.concat(sidecar_argv(), ' '),
    '',
    'peers:',
  }
  if #(M.state.peers or {}) == 0 then
    table.insert(lines, '  (none)')
  else
    for _, peer in ipairs(M.state.peers) do
      table.insert(lines, string.format('  - %s  cursor=%d color=%s typing=%s state=%s', peer.name or peer.participant_id or 'peer', peer.anchor or 0, peer.color or '#999999', tostring(peer.typing or false), peer_presence_state(peer)))
    end
  end
  return lines
end

local function help_lines()
  return {
    'grid-editor help',
    '',
    ':GridEditorOpen <doc>',
    ':GridEditorClose',
    ':GridEditorInfo',
    ':GridEditorPeers',
    ':GridEditorHelp',
    '',
    'Phase 1 notes:',
    '- peer markers and selections are rendered in-buffer',
    '- relay status changes are notified as they happen',
    '- name/color/relay defaults come from setup() or env vars',
  }
end

local function refresh_info_window()
  if not M.state.info_bufnr or not vim.api.nvim_buf_is_valid(M.state.info_bufnr) then
    return
  end
  vim.bo[M.state.info_bufnr].modifiable = true
  vim.api.nvim_buf_set_lines(M.state.info_bufnr, 0, -1, false, session_lines())
  vim.bo[M.state.info_bufnr].modifiable = false
end

local function open_info_window()
  if M.state.info_winid and vim.api.nvim_win_is_valid(M.state.info_winid) then
    refresh_info_window()
    vim.api.nvim_set_current_win(M.state.info_winid)
    return
  end
  local bufnr = vim.api.nvim_create_buf(false, true)
  M.state.info_bufnr = bufnr
  vim.bo[bufnr].bufhidden = 'wipe'
  vim.bo[bufnr].filetype = 'grid-editor'
  vim.bo[bufnr].modifiable = false
  local width = math.max(60, math.floor(vim.o.columns * 0.48))
  local height = 14 + math.max(1, #(M.state.peers or {}))
  local row = math.max(1, math.floor((vim.o.lines - height) / 2) - 1)
  local col = math.max(1, math.floor((vim.o.columns - width) / 2))
  M.state.info_winid = vim.api.nvim_open_win(bufnr, true, {
    relative = 'editor',
    row = row,
    col = col,
    width = math.min(width, vim.o.columns - 4),
    height = math.min(height, vim.o.lines - 4),
    style = 'minimal',
    border = 'rounded',
    title = ' grid-editor ',
    title_pos = 'center',
  })
  refresh_info_window()
end

local function open_help_window()
  if M.state.help_winid and vim.api.nvim_win_is_valid(M.state.help_winid) then
    vim.api.nvim_set_current_win(M.state.help_winid)
    return
  end
  local bufnr = vim.api.nvim_create_buf(false, true)
  M.state.help_bufnr = bufnr
  vim.bo[bufnr].bufhidden = 'wipe'
  vim.bo[bufnr].filetype = 'grid-editor'
  vim.bo[bufnr].modifiable = false
  local width = math.max(56, math.floor(vim.o.columns * 0.44))
  local height = #help_lines() + 2
  local row = math.max(1, math.floor((vim.o.lines - height) / 2) - 1)
  local col = math.max(1, math.floor((vim.o.columns - width) / 2))
  M.state.help_winid = vim.api.nvim_open_win(bufnr, true, {
    relative = 'editor',
    row = row,
    col = col,
    width = math.min(width, vim.o.columns - 4),
    height = math.min(height, vim.o.lines - 4),
    style = 'minimal',
    border = 'rounded',
    title = ' grid-editor help ',
    title_pos = 'center',
  })
  vim.bo[bufnr].modifiable = true
  vim.api.nvim_buf_set_lines(bufnr, 0, -1, false, help_lines())
  vim.bo[bufnr].modifiable = false
end

local function send_sidecar(message)
  if not M.state.job_id or M.state.job_id <= 0 then
    return
  end
  vim.fn.chansend(M.state.job_id, vim.json.encode(message) .. '\n')
end

local function set_buffer_content(content)
  if not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return
  end
  M.state.suppress = true
  M.state.expected_content = content or ''
  local lines = vim.split(content or '', '\n', { plain = true })
  vim.api.nvim_buf_set_lines(M.state.bufnr, 0, -1, false, lines)
  M.state.suppress = false
  -- Intent: Treat sidecar-driven buffer replacement as a remote replica update,
  -- not as a new local edit to echo back into the sidecar. Source: DI-sulod; DI-gafit
  draw_peers(M.state.peers)
end

local function update_cursor(typing)
  if not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return
  end
  local cursor = vim.api.nvim_win_get_cursor(0)
  local lines = vim.api.nvim_buf_get_lines(M.state.bufnr, 0, -1, false)
  local offset = pos_to_offset(lines, cursor[1] - 1, cursor[2])
  send_sidecar({
    type = 'set_cursor',
    anchor = offset,
    head = offset,
    typing = typing or false,
  })
end

local function handle_sidecar_message(message)
  if message.type == 'opened' then
    M.state.session_ready = true
    set_buffer_content(message.content or '')
    update_cursor(false)
  elseif message.type == 'changed' then
    set_buffer_content(message.content or '')
  elseif message.type == 'awareness' then
    local next_index = {}
    for _, peer in ipairs(message.peers or {}) do
      next_index[peer.participant_id] = peer
    end
    for participant_id, peer in pairs(next_index) do
      if not M.state.peer_index[participant_id] then
        vim.schedule(function()
          vim.notify('grid-editor peer joined: ' .. (peer.name or participant_id))
        end)
      end
    end
    for participant_id, peer in pairs(M.state.peer_index or {}) do
      if not next_index[participant_id] then
        vim.schedule(function()
          vim.notify('grid-editor peer left: ' .. (peer.name or participant_id))
        end)
      end
    end
    M.state.peers = message.peers or {}
    M.state.peer_index = next_index
    draw_peers(M.state.peers)
    refresh_info_window()
  elseif message.type == 'relay_status' then
    M.state.relay_connected = message.connected and true or false
    refresh_info_window()
    local status = message.connected and 'relay up' or 'relay down'
    vim.schedule(function()
      vim.notify('grid-editor ' .. status)
    end)
  elseif message.type == 'error' then
    vim.schedule(function()
      vim.notify('grid-editor sidecar error: ' .. (message.message or 'unknown'), vim.log.levels.ERROR)
    end)
  end
end

local function start_sidecar()
  if M.state.job_id and M.state.job_id > 0 then
    return true
  end

  local stdout_tail = ''
  local function queue_stdout_line(line)
    if line == '' then
      return
    end
    local ok, decoded = pcall(vim.json.decode, line)
    if ok then
      vim.schedule(function()
        handle_sidecar_message(decoded)
      end)
      return
    end
    vim.schedule(function()
      vim.notify('grid-editor dropped malformed sidecar output', vim.log.levels.WARN)
    end)
  end

  local function feed_stdout(data)
    if type(data) ~= 'table' then
      return
    end
    for index, chunk in ipairs(data) do
      local line = chunk
      if index == 1 then
        line = stdout_tail .. line
      end
      if index == #data then
        stdout_tail = line
      else
        queue_stdout_line(line)
      end
    end
  end

  local function flush_stdout_tail()
    if stdout_tail == '' then
      return
    end
    local line = stdout_tail
    stdout_tail = ''
    queue_stdout_line(line)
  end

  M.state.job_id = vim.fn.jobstart(sidecar_argv(), {
    rpc = false,
    stdout_buffered = false,
    stderr_buffered = false,
    on_stdout = function(_, data)
      -- Intent: Reassemble sidecar stdout into complete newline-delimited JSON
      -- records before decoding, because Neovim job callbacks can split large
      -- awareness payloads across multiple chunks. Source: DI-gafit
      feed_stdout(data)
    end,
    on_stderr = function(_, data)
      for _, line in ipairs(data) do
        if line ~= '' then
          vim.schedule(function()
            vim.notify(line, vim.log.levels.INFO)
          end)
        end
      end
    end,
    on_exit = function()
      flush_stdout_tail()
      M.state.job_id = nil
    end,
  })

  if M.state.job_id <= 0 then
    M.state.job_id = nil
    vim.notify('grid-editor failed to start sidecar', vim.log.levels.ERROR)
    return false
  end

  send_sidecar({
    type = 'connect',
    relay_url = M.config.relay_url,
    access_token = M.config.access_token,
    participant_id = M.state.participant_id,
    display_name = M.config.display_name,
    color = M.config.color,
  })
  return true
end

function M.open(doc_id)
  if not start_sidecar() then
    return
  end

  M.state.doc_id = doc_id
  M.state.session_ready = false
  M.state.bufnr = vim.api.nvim_create_buf(true, false)
  vim.api.nvim_buf_set_name(M.state.bufnr, 'grid-editor://' .. doc_id)
  vim.api.nvim_set_current_buf(M.state.bufnr)
  vim.wo.number = M.config.show_line_numbers
  vim.wo.relativenumber = false
  vim.wo.signcolumn = 'yes:2'

  if M.state.augroup then
    pcall(vim.api.nvim_del_augroup_by_id, M.state.augroup)
  end
  M.state.augroup = vim.api.nvim_create_augroup('GridEditor', { clear = true })

  vim.api.nvim_create_autocmd({ 'TextChanged', 'TextChangedI' }, {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function()
      if M.state.suppress or not M.state.session_ready then
        return
      end
      local content = join_lines()
      if M.state.expected_content ~= nil and content == M.state.expected_content then
        M.state.expected_content = nil
        return
      end
      M.state.expected_content = nil
      send_sidecar({
        type = 'set_text',
        content = content,
      })
      update_cursor(true)
    end,
  })

  vim.api.nvim_create_autocmd({ 'CursorMoved', 'CursorMovedI' }, {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function()
      update_cursor(false)
    end,
  })

  vim.api.nvim_create_autocmd('BufUnload', {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function()
      M.close()
    end,
  })

  send_sidecar({
    type = 'open',
    doc_id = doc_id,
  })
end

function M.close()
  if M.state.augroup then
    pcall(vim.api.nvim_del_augroup_by_id, M.state.augroup)
    M.state.augroup = nil
  end
  send_sidecar({ type = 'close' })
  if M.state.job_id and M.state.job_id > 0 then
    vim.fn.jobstop(M.state.job_id)
    M.state.job_id = nil
  end
  M.state.doc_id = nil
  M.state.session_ready = false
  M.state.bufnr = nil
  M.state.peers = {}
  M.state.peer_index = {}
  M.state.relay_connected = false
  if M.state.info_winid and vim.api.nvim_win_is_valid(M.state.info_winid) then
    vim.api.nvim_win_close(M.state.info_winid, true)
  end
  M.state.info_winid = nil
  M.state.info_bufnr = nil
  if M.state.help_winid and vim.api.nvim_win_is_valid(M.state.help_winid) then
    vim.api.nvim_win_close(M.state.help_winid, true)
  end
  M.state.help_winid = nil
  M.state.help_bufnr = nil
end

function M.info()
  vim.notify(table.concat(session_lines(), '\n'))
end

function M.peers()
  open_info_window()
end

function M.help()
  open_help_window()
end

function M.setup(opts)
  M.config = vim.tbl_deep_extend('force', M.config, opts or {})
  if M.config.sidecar_cmd == nil then
    M.config.sidecar_cmd = default_sidecar_cmd(M.config.repo_root)
  end
  vim.api.nvim_create_user_command('GridEditorOpen', function(command)
    local doc_id = command.args ~= '' and command.args or 'demo'
    M.open(doc_id)
  end, { nargs = '?', desc = 'Open a grid-editor document' })
  vim.api.nvim_create_user_command('GridEditorClose', function()
    M.close()
  end, { nargs = 0, desc = 'Close the current grid-editor session' })
  vim.api.nvim_create_user_command('GridEditorInfo', function()
    M.info()
  end, { nargs = 0, desc = 'Show grid-editor connection info' })
  vim.api.nvim_create_user_command('GridEditorPeers', function()
    M.peers()
  end, { nargs = 0, desc = 'Show grid-editor peer roster' })
  vim.api.nvim_create_user_command('GridEditorHelp', function()
    M.help()
  end, { nargs = 0, desc = 'Show grid-editor help' })
end

return M
