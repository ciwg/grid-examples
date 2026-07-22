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

func TestNeovimSnapshotCreatesDurableRevisionFromLiveDraft(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	type livePayload struct {
		UpdateBody bool   `json:"update_body"`
		Body       string `json:"body"`
	}
	type revisionPayload struct {
		Actor   string   `json:"actor"`
		Title   string   `json:"title"`
		Summary string   `json:"summary"`
		Body    string   `json:"body"`
		Tags    []string `json:"tags"`
	}

	var pushedBody string
	var revisionSeen bool
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/items/ITEM-0001/live":
			if _, err := fmt.Fprint(response, `{
				"item_id":"ITEM-0001",
				"title":"Start line A",
				"status":"draft",
				"version":4,
				"current_revision":2,
				"body":"old line",
				"participants":[]
			}`); err != nil {
				t.Fatalf("write live get response: %v", err)
			}
		case request.Method == http.MethodPost && request.URL.Path == "/api/items/ITEM-0001/live":
			var payload livePayload
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode live payload: %v", err)
			}
			if payload.UpdateBody {
				pushedBody = payload.Body
			}
			if _, err := fmt.Fprintf(response, `{
				"item_id":"ITEM-0001",
				"title":"Start line A",
				"status":"draft",
				"version":5,
				"current_revision":2,
				"body":%q,
				"participants":[]
			}`, map[bool]string{true: payload.Body, false: "old line"}[payload.UpdateBody]); err != nil {
				t.Fatalf("write live post response: %v", err)
			}
		case request.Method == http.MethodGet && request.URL.Path == "/api/items/ITEM-0001":
			if _, err := fmt.Fprint(response, `{
				"id":"ITEM-0001",
				"kind":"procedure",
				"title":"Start line A",
				"summary":"Startup checklist",
				"status":"draft",
				"current_revision":2,
				"working_version":5,
				"tags":["startup","audit"],
				"responsibility_ids":[],
				"revisions":[],
				"approvals":[],
				"related_runs":[],
				"links":[]
			}`); err != nil {
				t.Fatalf("write item get response: %v", err)
			}
		case request.Method == http.MethodPost && request.URL.Path == "/api/items/ITEM-0001/revisions":
			var payload revisionPayload
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode revision payload: %v", err)
			}
			if payload.Actor != "Boss" {
				t.Fatalf("unexpected actor: %#v", payload)
			}
			if payload.Title != "Start line A" || payload.Summary != "Startup checklist" {
				t.Fatalf("unexpected title/summary: %#v", payload)
			}
			if payload.Body != "new line 1\nnew line 2" {
				t.Fatalf("unexpected revision body: %#v", payload)
			}
			if strings.Join(payload.Tags, ",") != "startup,audit" {
				t.Fatalf("unexpected tags: %#v", payload)
			}
			revisionSeen = true
			if _, err := fmt.Fprint(response, `{
				"id":"ITEM-0001",
				"kind":"procedure",
				"title":"Start line A",
				"summary":"Startup checklist",
				"status":"draft",
				"current_revision":3,
				"working_version":5,
				"tags":["startup","audit"],
				"responsibility_ids":[],
				"revisions":[
					{"number":1,"title":"Start line A","summary":"Initial","created_at":"2026-07-21T09:00:00Z"},
					{"number":2,"title":"Start line A","summary":"Updated","created_at":"2026-07-21T10:00:00Z"},
					{"number":3,"title":"Start line A","summary":"Startup checklist","created_at":"2026-07-22T01:15:00Z"}
				],
				"approvals":[],
				"related_runs":[],
				"links":[]
			}`); err != nil {
				t.Fatalf("write revision response: %v", err)
			}
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "snapshot.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
vim.env.OKS_DISPLAY_NAME = "Boss"
local oks = require("oks")
oks.setup()

vim.cmd("OksOpen ITEM-0001")
vim.api.nvim_buf_set_lines(0, 0, -1, false, { "new line 1", "new line 2" })
vim.cmd("OksSnapshot")

if oks.state.current_revision ~= 3 then
  error("unexpected revision " .. tostring(oks.state.current_revision))
end
local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if body ~= "new line 1\nnew line 2" then
  error("unexpected live body " .. body)
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
		t.Fatalf("nvim snapshot regression: %v\n%s", err, string(output))
	}
	if pushedBody != "new line 1\nnew line 2" {
		t.Fatalf("live draft push did not carry updated body: %q", pushedBody)
	}
	if !revisionSeen {
		t.Fatalf("revision POST was not observed")
	}
	if strings.Contains(string(output), "unexpected revision") || strings.Contains(string(output), "unexpected live body") {
		t.Fatalf("unexpected snapshot output: %s", string(output))
	}
}
