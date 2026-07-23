package nvim

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNeovimSupersedeItemRefreshesInspector(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	type supersedePayload struct {
		Actor string `json:"actor"`
		Notes string `json:"notes"`
	}

	postSeen := false
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/items/ITEM-0001" && request.URL.Path != "/api/items/ITEM-0001/supersede" {
			http.NotFound(response, request)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/items/ITEM-0001":
			if _, err := fmt.Fprint(response, `{
				"id":"ITEM-0001",
				"kind":"procedure",
				"title":"Start line A",
				"summary":"Startup checklist",
				"status":"approved",
				"current_revision":2,
				"working_version":4,
				"responsibility_ids":[],
				"revisions":[
					{"number":1,"title":"Start line A","summary":"Initial","created_at":"2026-07-21T09:00:00Z"},
					{"number":2,"title":"Start line A","summary":"Updated","created_at":"2026-07-21T10:00:00Z"}
				],
				"approvals":[],
				"related_runs":[],
				"links":[]
			}`); err != nil {
				t.Fatalf("write get response: %v", err)
			}
		case request.Method == http.MethodPost && request.URL.Path == "/api/items/ITEM-0001/supersede":
			var payload supersedePayload
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode supersede payload: %v", err)
			}
			if payload.Actor != "Boss" || payload.Notes != "replaced by audited startup flow" {
				t.Fatalf("unexpected supersede payload: %#v", payload)
			}
			postSeen = true
			if _, err := fmt.Fprint(response, `{
				"id":"ITEM-0001",
				"kind":"procedure",
				"title":"Start line A",
				"summary":"Startup checklist",
				"status":"superseded",
				"current_revision":2,
				"working_version":4,
				"responsibility_ids":[],
				"revisions":[
					{"number":1,"title":"Start line A","summary":"Initial","created_at":"2026-07-21T09:00:00Z"},
					{"number":2,"title":"Start line A","summary":"Updated","created_at":"2026-07-21T10:00:00Z"}
				],
				"approvals":[],
				"related_runs":[],
				"links":[]
			}`); err != nil {
				t.Fatalf("write post response: %v", err)
			}
		default:
			http.Error(response, "unexpected request", http.StatusBadRequest)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "supersede_item.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
vim.env.OKS_SOCKET = "off"
vim.env.OKS_DISPLAY_NAME = "Boss"
local oks = require("oks")
oks.setup()

vim.cmd("OksInspect ITEM-0001")
vim.cmd("OksSupersedeItem replaced by audited startup flow")

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "status: superseded", 1, true) then
  error("missing superseded status")
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
		t.Fatalf("nvim supersede regression: %v\n%s", err, string(output))
	}
	if !postSeen {
		t.Fatalf("supersede POST was not observed")
	}
	if strings.Contains(string(output), "missing ") {
		t.Fatalf("unexpected supersede output: %s", string(output))
	}
}
