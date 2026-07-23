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

func TestNeovimApproveItemUsesCurrentRevisionAndRefreshesInspector(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	type approvalPayload struct {
		Actor    string `json:"actor"`
		Revision int    `json:"revision"`
		Role     string `json:"role"`
		Decision string `json:"decision"`
		Notes    string `json:"notes"`
	}

	postSeen := false
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/items/ITEM-0001" && request.URL.Path != "/api/items/ITEM-0001/approvals" {
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
				"status":"draft",
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
		case request.Method == http.MethodPost && request.URL.Path == "/api/items/ITEM-0001/approvals":
			var payload approvalPayload
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode approval payload: %v", err)
			}
			if payload.Actor != "Boss" {
				t.Fatalf("unexpected actor: %#v", payload)
			}
			if payload.Revision != 2 || payload.Role != "reviewer" || payload.Decision != "approved" || payload.Notes != "ready for use" {
				t.Fatalf("unexpected approval payload: %#v", payload)
			}
			postSeen = true
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
				"approvals":[
					{"actor":"Boss","revision":2,"role":"reviewer","decision":"approved","notes":"ready for use"}
				],
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

	script := filepath.Join(t.TempDir(), "approve_item.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
vim.env.OKS_SOCKET = "off"
vim.env.OKS_DISPLAY_NAME = "Boss"
local oks = require("oks")
oks.setup()

vim.cmd("OksInspect ITEM-0001")
vim.cmd("OksApproveItem reviewer approved ready for use")

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "status: approved", 1, true) then
  error("missing approved status")
end
if not string.find(body, "approved by Boss role=reviewer revision=2", 1, true) then
  error("missing refreshed approval line")
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
		t.Fatalf("nvim approve regression: %v\n%s", err, string(output))
	}
	if !postSeen {
		t.Fatalf("approval POST was not observed")
	}
	if strings.Contains(string(output), "missing ") {
		t.Fatalf("unexpected approval output: %s", string(output))
	}
}
