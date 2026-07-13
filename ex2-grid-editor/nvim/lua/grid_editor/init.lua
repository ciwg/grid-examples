local M = {}

M.config = {
  base_url = 'http://127.0.0.1:7001',
  display_name = 'Neovim User',
  color = '#d66f1d',
  poll_ms = 1000,
}

M.state = {
  bufnr = nil,
  doc_id = nil,
  timer = nil,
  suppress = false,
  last_message_cid = nil,
  cursor_ns = nil,
  augroup = nil,
  participant_id = 'nvim-' .. tostring(vim.fn.getpid()),
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

local function request(args, body)
  local cmd = { 'curl', '-s' }
  if body then
    table.insert(cmd, '-X')
    table.insert(cmd, 'POST')
    table.insert(cmd, '-H')
    table.insert(cmd, 'Content-Type: application/json')
    table.insert(cmd, '-d')
    table.insert(cmd, vim.json.encode(body))
  end
  for _, arg in ipairs(args) do
    table.insert(cmd, arg)
  end
  local result = vim.system(cmd, { text = true }):wait()
  if result.code ~= 0 then
    error(result.stderr ~= '' and result.stderr or ('curl exited with ' .. result.code))
  end
  if result.stdout == '' then
    return nil
  end
  return vim.json.decode(result.stdout)
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
    local row, col = offset_to_pos(lines, peer.cursor or 0)
    vim.api.nvim_buf_set_extmark(M.state.bufnr, M.state.cursor_ns, row, col, {
      virt_text = { { ' ' .. (peer.display_name or peer.author or 'peer') .. ' ', 'Search' } },
      virt_text_pos = 'overlay',
      hl_group = 'Visual',
      end_col = math.min(col + 1, #(lines[row + 1] or '')),
    })
  end
end

local function send_replace()
  if M.state.suppress or not M.state.doc_id then
    return
  end
  request({
    M.config.base_url .. '/api/local/documents/' .. M.state.doc_id .. '/replace',
  }, {
    content = join_lines(),
    embodiment = 'nvim',
  })
end

local function send_awareness(typing)
  if not M.state.doc_id or not M.state.bufnr or not vim.api.nvim_buf_is_valid(M.state.bufnr) then
    return
  end
  local cursor = vim.api.nvim_win_get_cursor(0)
  local lines = vim.api.nvim_buf_get_lines(M.state.bufnr, 0, -1, false)
  local offset = pos_to_offset(lines, cursor[1] - 1, cursor[2])
  request({
    M.config.base_url .. '/api/local/documents/' .. M.state.doc_id .. '/awareness',
  }, {
    participant_id = M.state.participant_id,
    cursor = offset,
    head = offset,
    typing = typing or false,
    display_name = M.config.display_name,
    color = M.config.color,
    embodiment = 'nvim',
  })
end

local function poll_once()
  if not M.state.doc_id then
    return
  end
  local payload = request({
    M.config.base_url .. '/api/local/documents/' .. M.state.doc_id .. '/state',
  })
  if not payload then
    return
  end
  if payload.message_cid ~= M.state.last_message_cid then
    M.state.suppress = true
    local lines = vim.split(payload.content or '', '\n', { plain = true })
    vim.api.nvim_buf_set_lines(M.state.bufnr, 0, -1, false, lines)
    M.state.suppress = false
    M.state.last_message_cid = payload.message_cid
  end
  draw_peers(payload.awareness)
end

local function stop_timer()
  if M.state.timer then
    M.state.timer:stop()
    M.state.timer:close()
    M.state.timer = nil
  end
end

local function start_timer()
  stop_timer()
  M.state.timer = vim.loop.new_timer()
  M.state.timer:start(0, M.config.poll_ms, vim.schedule_wrap(function()
    pcall(poll_once)
  end))
end

function M.open(doc_id)
  M.state.doc_id = doc_id
  M.state.last_message_cid = nil
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
      pcall(send_replace)
      pcall(send_awareness, true)
    end,
  })
  vim.api.nvim_create_autocmd({ 'CursorMoved', 'CursorMovedI' }, {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function()
      pcall(send_awareness, false)
    end,
  })
  vim.api.nvim_create_autocmd('BufUnload', {
    group = M.state.augroup,
    buffer = M.state.bufnr,
    callback = function()
      stop_timer()
    end,
  })
  start_timer()
  pcall(poll_once)
  pcall(send_awareness, false)
end

function M.close()
  stop_timer()
  if M.state.augroup then
    pcall(vim.api.nvim_del_augroup_by_id, M.state.augroup)
    M.state.augroup = nil
  end
  M.state.doc_id = nil
  M.state.bufnr = nil
  M.state.last_message_cid = nil
end

function M.info()
  vim.notify(string.format('grid-editor base=%s doc=%s', M.config.base_url, M.state.doc_id or 'none'))
end

function M.setup(opts)
  M.config = vim.tbl_deep_extend('force', M.config, opts or {})
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
end

return M
