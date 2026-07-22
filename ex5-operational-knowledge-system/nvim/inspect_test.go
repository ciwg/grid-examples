package nvim

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/service"
)

func TestNeovimInspectItemRendersProjectedDetail(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	socketPath := filepath.Join(t.TempDir(), "embodiment.sock")
	go serveSingleNvimSocketResponse(t, socketPath, service.LocalEmbodimentResponse{
		Type:   "response",
		Status: 200,
		Body: `{
			"id":"ITEM-0001",
			"kind":"procedure",
			"status":"approved",
			"title":"Start line A",
			"summary":"Startup checklist",
			"current_revision":2,
			"working_version":3,
			"responsibility_ids":["RESP-0001"],
			"revisions":[{"number":2,"title":"Start line A","summary":"Adds safety latch check","created_at":"2026-07-21T10:00:00Z"}],
			"approvals":[{"decision":"approved","actor":"carol","role":"reviewer","revision":2,"notes":"Ready for operators"}],
			"related_runs":[{"id":"RUN-0001","kind":"procedure","revision":2,"outcome":"completed","notes":"Shift startup finished","place_id":"PLACE-0001","resource_ids":["RES-0001"],"evidence":[{"facts":{"result":"ok"}}]}],
			"links":[{"from_type":"knowledge_item","from_id":"ITEM-0001","to_type":"resource","to_id":"RES-0001","relation":"uses","notes":"References startup key rack"}]
		}`,
	})
	waitForNvimUnixSocket(t, socketPath)

	script := filepath.Join(t.TempDir(), "inspect_item.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_SOCKET_PATH = %q
local oks = require("oks")
oks.setup()

vim.cmd("OksInspect ITEM-0001")

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "# Start line A", 1, true) then
  error("missing item header")
end
if not string.find(body, "## Revisions", 1, true) then
  error("missing revisions section")
end
if not string.find(body, "## Approvals", 1, true) then
  error("missing approvals section")
end
if not string.find(body, "## Related runs", 1, true) then
  error("missing related runs section")
end
if not string.find(body, "inspect: :OksInspectRun RUN-0001", 1, true) then
  error("missing related run inspect hint")
end
if not string.find(body, "resources: RES-0001", 1, true) then
  error("missing related run resource summary")
end
if not string.find(vim.api.nvim_buf_get_name(0), "oks-inspect://ITEM-0001", 1, true) then
  error("unexpected item buffer name: " .. vim.api.nvim_buf_get_name(0))
end
vim.cmd("qa!")
`, socketPath)
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
		t.Fatalf("nvim inspect item regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing ") || strings.Contains(string(output), "unexpected item buffer name") {
		t.Fatalf("unexpected inspect item output: %s", string(output))
	}
}

func TestNeovimInspectRunRendersContextHandoffs(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/runs/RUN-0001" {
			http.NotFound(response, request)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprint(response, `{
			"id":"RUN-0001",
			"kind":"receiving_check",
			"item_id":"ITEM-0001",
			"item_kind":"receiving_check",
			"revision":2,
			"actor":"dave",
			"outcome":"accepted_with_notes",
			"notes":"Outer wrap torn",
			"place_id":"PLACE-0001",
			"resource_ids":["RES-0001"],
			"responsibility_ids":["RESP-0001"],
			"evidence":[{"summary":"Dock photo","facts":{"supplier":"Acme"}}],
			"approvals":[{"decision":"noted","actor":"ellen","role":"reviewer","revision":2,"notes":"Accepted with note"}],
			"links":[{"from_type":"run","from_id":"RUN-0001","to_type":"place","to_id":"PLACE-0001","relation":"observed_at","notes":"Recorded at receiving dock"}]
		}`); err != nil {
			t.Fatalf("write run response: %v", err)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "inspect_run.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
local oks = require("oks")
oks.setup()

vim.cmd("OksInspectRun RUN-0001")

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "## Evidence", 1, true) then
  error("missing evidence section")
end
if not string.find(body, "## Approvals", 1, true) then
  error("missing approvals section")
end
if not string.find(body, "## Handoffs", 1, true) then
  error("missing handoffs section")
end
if not string.find(body, ":OksInspect ITEM-0001", 1, true) then
  error("missing item handoff")
end
if not string.find(body, ":OksInspectEntity place PLACE-0001", 1, true) then
  error("missing place handoff")
end
if not string.find(body, ":OksInspectEntity resource RES-0001", 1, true) then
  error("missing resource handoff")
end
if not string.find(body, ":OksInspectEntity responsibility RESP-0001", 1, true) then
  error("missing responsibility handoff")
end
if not string.find(vim.api.nvim_buf_get_name(0), "oks-run://RUN-0001", 1, true) then
  error("unexpected run buffer name: " .. vim.api.nvim_buf_get_name(0))
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
		t.Fatalf("nvim inspect run regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing ") || strings.Contains(string(output), "unexpected run buffer name") {
		t.Fatalf("unexpected inspect run output: %s", string(output))
	}
}

func TestNeovimInspectEntityRendersResponsibilityDetail(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/responsibilities/RESP-0001" {
			http.NotFound(response, request)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprint(response, `{
			"id":"RESP-0001",
			"title":"Receiving lead",
			"summary":"Owns intake review",
			"team":"ops",
			"linked_item_ids":["ITEM-0001"],
			"linked_run_ids":["RUN-0001"],
			"linked_role_keys":["reviewer"],
			"related_runs":[{"id":"RUN-0001","kind":"receiving_check","outcome":"accepted_with_notes"}],
			"links":[{"from_type":"responsibility","from_id":"RESP-0001","to_type":"knowledge_item","to_id":"ITEM-0001","relation":"owns","notes":"Primary intake procedure"}]
		}`); err != nil {
			t.Fatalf("write responsibility response: %v", err)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "inspect_entity.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
local oks = require("oks")
oks.setup()

vim.cmd("OksInspectEntity responsibility RESP-0001")

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "# Receiving lead", 1, true) then
  error("missing responsibility header")
end
if not string.find(body, "## Related runs", 1, true) then
  error("missing responsibility related runs section")
end
if not string.find(body, "inspect: :OksInspectRun RUN-0001", 1, true) then
  error("missing responsibility related run inspect hint")
end
if not string.find(body, "## Links", 1, true) then
  error("missing responsibility links section")
end
if not string.find(vim.api.nvim_buf_get_name(0), "oks-responsibility://RESP-0001", 1, true) then
  error("unexpected responsibility buffer name: " .. vim.api.nvim_buf_get_name(0))
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
		t.Fatalf("nvim inspect entity regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing ") || strings.Contains(string(output), "unexpected responsibility buffer name") {
		t.Fatalf("unexpected inspect entity output: %s", string(output))
	}
}

func serveSingleNvimSocketResponse(t *testing.T, socketPath string, response service.LocalEmbodimentResponse) {
	t.Helper()
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("listen unix socket: %v", err)
	}
	defer func() {
		_ = listener.Close()
	}()
	conn, err := listener.Accept()
	if err != nil {
		t.Fatalf("accept unix socket: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	var request service.LocalEmbodimentRequest
	if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&request); err != nil {
		t.Fatalf("decode unix socket request: %v", err)
	}
	if request.Type != "request" {
		t.Fatalf("unexpected socket request type: %+v", request)
	}
	if err := json.NewEncoder(conn).Encode(response); err != nil {
		t.Fatalf("encode unix socket response: %v", err)
	}
}

func waitForNvimUnixSocket(t *testing.T, socketPath string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(socketPath); err == nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("socket did not become ready: %s", socketPath)
}
