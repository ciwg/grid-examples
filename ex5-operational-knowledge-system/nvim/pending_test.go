package nvim

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNeovimPendingRendersDraftAndReviewQueues(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/search" {
			http.NotFound(response, request)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		switch request.URL.RawQuery {
		case "status=draft":
			if _, err := fmt.Fprint(response, `{
				"filters":{"status":"draft"},
				"places":[],
				"resources":[],
				"responsibilities":[],
				"items":[{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Start line A","summary":"Startup checklist"}],
				"runs":[]
			}`); err != nil {
				t.Fatalf("write draft response: %v", err)
			}
		case "":
			if _, err := fmt.Fprint(response, `{
				"filters":{},
				"places":[],
				"resources":[],
				"responsibilities":[],
				"items":[],
				"runs":[
					{"id":"RUN-0001","kind":"procedure","item_id":"ITEM-0001","outcome":"completed","notes":"Shift startup finished","approvals":[]},
					{"id":"RUN-0002","kind":"receiving_check","item_id":"ITEM-0002","outcome":"accepted_with_notes","notes":"Outer wrap torn","approvals":[{"decision":"noted","actor":"carol"}]}
				]
			}`); err != nil {
				t.Fatalf("write all-runs response: %v", err)
			}
		case "problem=true":
			if _, err := fmt.Fprint(response, `{
				"filters":{"problem":true},
				"places":[],
				"resources":[],
				"responsibilities":[],
				"items":[],
				"runs":[
					{"id":"RUN-0002","kind":"receiving_check","item_id":"ITEM-0002","outcome":"accepted_with_notes","notes":"Outer wrap torn","approvals":[{"decision":"noted","actor":"carol"}]}
				]
			}`); err != nil {
				t.Fatalf("write problem response: %v", err)
			}
		default:
			http.Error(response, fmt.Sprintf("unexpected query %q", request.URL.RawQuery), http.StatusBadRequest)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "pending.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
local oks = require("oks")

oks.pending()

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "## Draft items", 1, true) then
  error("missing draft section")
end
if not string.find(body, "inspect: :OksInspect ITEM-0001", 1, true) then
  error("missing item inspect hint")
end
if not string.find(body, "## Unreviewed runs", 1, true) then
  error("missing unreviewed runs section")
end
if not string.find(body, "inspect: :OksInspectRun RUN-0001", 1, true) then
  error("missing unreviewed run inspect hint")
end
if not string.find(body, "## Problem runs", 1, true) then
  error("missing problem runs section")
end
if not string.find(body, "inspect: :OksInspectRun RUN-0002", 1, true) then
  error("missing problem run inspect hint")
end
if not string.find(vim.api.nvim_buf_get_name(0), "oks-pending://review", 1, true) then
  error("unexpected pending buffer name: " .. vim.api.nvim_buf_get_name(0))
end
vim.cmd("qa!")
`, server.URL)
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
		t.Fatalf("nvim pending regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing ") || strings.Contains(string(output), "unexpected pending buffer name") {
		t.Fatalf("unexpected pending output: %s", string(output))
	}
}
