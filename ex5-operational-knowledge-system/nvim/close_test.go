package nvim

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNeovimCloseWipesLiveAndInspectorBuffers(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	script := filepath.Join(t.TempDir(), "close.lua")
	scriptBody := `
local oks = require("oks")

oks.state.bufnr = vim.api.nvim_get_current_buf()
oks.state.winid = vim.api.nvim_get_current_win()
vim.api.nvim_buf_set_name(oks.state.bufnr, "oks://ITEM-0001")
vim.bo[oks.state.bufnr].buftype = "acwrite"

vim.cmd("vsplit")
local inspector_buf = vim.api.nvim_create_buf(false, true)
vim.api.nvim_buf_set_name(inspector_buf, "oks-inspect://ITEM-0001")
vim.api.nvim_win_set_buf(vim.api.nvim_get_current_win(), inspector_buf)
oks._test_set_inspector(inspector_buf, vim.api.nvim_get_current_win())

oks.close()

local live_valid = vim.api.nvim_buf_is_valid(oks.state.bufnr or -1)
local inspector_valid = vim.api.nvim_buf_is_valid(inspector_buf)
if live_valid or inspector_valid then
  error(string.format("close left buffers alive live=%s inspector=%s", tostring(live_valid), tostring(inspector_valid)))
end
if oks.state.item_id ~= nil or oks.state.bufnr ~= nil or oks.state.winid ~= nil then
  error("close left live state populated")
end
vim.cmd("qa!")
`
	if err := os.WriteFile(script, []byte(scriptBody), 0o644); err != nil {
		t.Fatalf("write script: %v", err)
	}

	command := exec.Command(
		nvimPath,
		"--headless",
		"-u", "NONE",
		"-c", "set runtimepath+=.",
		"-l", script,
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("nvim close regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "left buffers alive") {
		t.Fatalf("unexpected close output: %s", string(output))
	}
}
