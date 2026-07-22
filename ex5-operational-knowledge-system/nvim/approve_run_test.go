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

func TestNeovimApproveRunRefreshesRunInspector(t *testing.T) {
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
		if request.URL.Path != "/api/runs/RUN-0001" && request.URL.Path != "/api/runs/RUN-0001/approvals" {
			http.NotFound(response, request)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/runs/RUN-0001":
			if _, err := fmt.Fprint(response, `{
				"id":"RUN-0001",
				"kind":"procedure",
				"item_id":"ITEM-0001",
				"item_kind":"procedure",
				"revision":2,
				"actor":"alice",
				"outcome":"completed",
				"notes":"Shift startup complete",
				"place_id":"PLACE-0001",
				"resource_ids":["RES-0001"],
				"evidence":[],
				"approvals":[],
				"links":[]
			}`); err != nil {
				t.Fatalf("write get response: %v", err)
			}
		case request.Method == http.MethodPost && request.URL.Path == "/api/runs/RUN-0001/approvals":
			var payload approvalPayload
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode approval payload: %v", err)
			}
			if payload.Actor != "Boss" {
				t.Fatalf("unexpected actor: %#v", payload)
			}
			if payload.Revision != 0 || payload.Role != "approver" || payload.Decision != "noted" || payload.Notes != "shift handoff recorded" {
				t.Fatalf("unexpected approval payload: %#v", payload)
			}
			postSeen = true
			if _, err := fmt.Fprint(response, `{
				"id":"RUN-0001",
				"kind":"procedure",
				"item_id":"ITEM-0001",
				"item_kind":"procedure",
				"revision":2,
				"actor":"alice",
				"outcome":"completed",
				"notes":"Shift startup complete",
				"place_id":"PLACE-0001",
				"resource_ids":["RES-0001"],
				"evidence":[],
				"approvals":[
					{"actor":"Boss","revision":0,"role":"approver","decision":"noted","notes":"shift handoff recorded"}
				],
				"links":[]
			}`); err != nil {
				t.Fatalf("write post response: %v", err)
			}
		default:
			http.Error(response, "unexpected request", http.StatusBadRequest)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "approve_run.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
vim.env.OKS_DISPLAY_NAME = "Boss"
local oks = require("oks")
oks.setup()

vim.cmd("OksInspectRun RUN-0001")
vim.cmd("OksApproveRun approver noted shift handoff recorded")

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "## Approvals", 1, true) then
  error("missing approvals section")
end
if not string.find(body, "noted by Boss role=approver revision=0", 1, true) then
  error("missing refreshed run approval line")
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
		t.Fatalf("nvim approve run regression: %v\n%s", err, string(output))
	}
	if !postSeen {
		t.Fatalf("approval POST was not observed")
	}
	if strings.Contains(string(output), "missing ") {
		t.Fatalf("unexpected run approval output: %s", string(output))
	}
}
