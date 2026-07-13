local M = {}

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
  relay_url = vim.env.GRID_EDITOR_RELAY_URL or 'http://127.0.0.1:7001',
  sidecar_cmd = nil,
  display_name = vim.env.GRID_EDITOR_DISPLAY_NAME or 'Neovim User',
  color = vim.env.GRID_EDITOR_COLOR or '#d66f1d',
}

M.state = {
  bufnr = nil,
  doc_id = nil,
  suppress = false,
  cursor_ns = nil,
  augroup = nil,
  job_id = nil,
  participant_id = 'nvim-' .. tostring(vim.fn.getpid()),
  peers = {},
  relay_connected = false,
  info_bufnr = nil,
  info_winid = nil,
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

local function draw_peers(peers)
  if not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return
  end
  if not M.state.cursor_ns then
    M.state.cursor_ns = vim.api.nvim_create_namespace('grid_editor_peers')
  end
  vim.api.nvim_buf_clear_namespace(M.state.bufnr, M.state.cursor_ns, 0, -1)
  local lines = vim.api.nvim_buf_get_lines(M.state.bufnr, 0, -1, false)
  for _, peer in ipairs(peers or {}) do
    local row, col = offset_to_pos(lines, peer.anchor or 0)
    local group = peer_highlight(peer.color)
    -- Intent: Make remote awareness visible in-buffer with per-peer colors so
    -- Neovim users can tell who else is present without opening extra UI.
    -- Source: DI-gafit; DI-samuv
    vim.api.nvim_buf_set_extmark(M.state.bufnr, M.state.cursor_ns, row, col, {
      virt_text = {
        { '▏', group },
        { ' ' .. (peer.name or peer.participant_id or 'peer') .. ' ', group },
      },
      virt_text_pos = 'overlay',
      end_col = math.min(col + 1, #(lines[row + 1] or '')),
    })
  end
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
    'relay status: ' .. (M.state.relay_connected and 'connected' or 'disconnected'),
    'sidecar: ' .. table.concat(sidecar_argv(), ' '),
    '',
    'peers:',
  }
  if #(M.state.peers or {}) == 0 then
    table.insert(lines, '  (none)')
  else
    for _, peer in ipairs(M.state.peers) do
      table.insert(lines, string.format('  - %s  cursor=%d color=%s typing=%s', peer.name or peer.participant_id or 'peer', peer.anchor or 0, peer.color or '#999999', tostring(peer.typing or false)))
    end
  end
  return lines
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
  local lines = vim.split(content or '', '\n', { plain = true })
  vim.api.nvim_buf_set_lines(M.state.bufnr, 0, -1, false, lines)
  M.state.suppress = false
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
  if message.type == 'opened' or message.type == 'changed' then
    set_buffer_content(message.content or '')
  elseif message.type == 'awareness' then
    M.state.peers = message.peers or {}
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

  local stdout_chunks = {}
  local function flush_stdout()
    while #stdout_chunks > 0 do
      local line = table.remove(stdout_chunks, 1)
      if line ~= '' then
        local ok, decoded = pcall(vim.json.decode, line)
        if ok then
          vim.schedule(function()
            handle_sidecar_message(decoded)
          end)
        end
      end
    end
  end

  M.state.job_id = vim.fn.jobstart(sidecar_argv(), {
    rpc = false,
    stdout_buffered = false,
    stderr_buffered = false,
    on_stdout = function(_, data)
      for _, line in ipairs(data) do
        if line ~= '' then
          table.insert(stdout_chunks, line)
        end
      end
      flush_stdout()
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
  M.state.bufnr = vim.api.nvim_create_buf(true, false)
  vim.api.nvim_buf_set_name(M.state.bufnr, 'grid-editor://' .. doc_id)
  vim.api.nvim_set_current_buf(M.state.bufnr)

  if M.state.augroup then
    pcall(vim.api.nvim_del_augroup_by_id, M.state.augroup)
  end
  M.state.augroup = vim.api.nvim_create_augroup('GridEditor', { clear = true })

  vim.api.nvim_create_autocmd({ 'TextChanged', 'TextChangedI' }, {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function()
      if M.state.suppress then
        return
      end
      send_sidecar({
        type = 'set_text',
        content = join_lines(),
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
  update_cursor(false)
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
  M.state.bufnr = nil
  M.state.peers = {}
  M.state.relay_connected = false
  if M.state.info_winid and vim.api.nvim_win_is_valid(M.state.info_winid) then
    vim.api.nvim_win_close(M.state.info_winid, true)
  end
  M.state.info_winid = nil
  M.state.info_bufnr = nil
end

function M.info()
  vim.notify(table.concat(session_lines(), '\n'))
end

function M.peers()
  open_info_window()
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
end

return M
