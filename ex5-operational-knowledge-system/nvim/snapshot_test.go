package nvim

import (
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

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/service"
)

func TestNeovimSnapshotCreatesDurableRevisionFromLiveDraft(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}
	dataRoot := filepath.Join(t.TempDir(), "runtime")
	app, err := service.NewApp(dataRoot)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", service.KnowledgeKindProcedure, "Start line A", "Startup checklist", "old line", []string{"startup", "audit"}, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if _, err := app.AddRevision("alice", item.ID, "Start line A", "Initial", "old line", []string{"startup", "audit"}); err != nil {
		t.Fatalf("create initial revision: %v", err)
	}
	if _, err := app.AddRevision("alice", item.ID, "Start line A", "Updated", "old line", []string{"startup", "audit"}); err != nil {
		t.Fatalf("create second revision: %v", err)
	}
	socketServer := service.NewLocalEmbodimentServer(app, service.EmbodimentSocketPath(dataRoot))
	go func() {
		if err := socketServer.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = socketServer.Close()
	}()
	waitForNvimUnixSocket(t, service.EmbodimentSocketPath(dataRoot))

	script := filepath.Join(t.TempDir(), "snapshot.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_SOCKET_PATH = %q
vim.env.OKS_DISPLAY_NAME = "Boss"
local oks = require("oks")
oks.setup()

vim.cmd("OksOpen %s")
vim.wait(2000, function()
  return oks.state.transport == "local-socket"
end, 50)
if oks.state.transport ~= "local-socket" then
  error("local socket transport did not connect")
end
vim.api.nvim_buf_set_lines(0, 0, -1, false, { "new line 1", "new line 2" })
vim.cmd("OksSnapshot")

if oks.state.current_revision ~= 4 then
  error("unexpected revision " .. tostring(oks.state.current_revision))
end
local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if body ~= "new line 1\nnew line 2" then
  error("unexpected live body " .. body)
end
vim.cmd("qa!")
`, service.EmbodimentSocketPath(dataRoot), item.ID)
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
	response := localSocketNvimRequest(t, service.EmbodimentSocketPath(dataRoot), service.LocalEmbodimentRequest{
		Type:   "request",
		Method: "GET",
		Path:   "/api/items/" + item.ID,
	})
	var refreshed map[string]any
	if err := json.Unmarshal([]byte(response.Body), &refreshed); err != nil {
		t.Fatalf("decode refreshed item: %v", err)
	}
	if refreshed["current_revision"] != float64(4) {
		t.Fatalf("unexpected current revision after snapshot: %#v", refreshed["current_revision"])
	}
	if strings.Contains(string(output), "unexpected revision") || strings.Contains(string(output), "unexpected live body") {
		t.Fatalf("unexpected snapshot output: %s", string(output))
	}
}

func TestNeovimSnapshotUsesRepoRootDefaultSocketPathAcrossWorkingDirectoryChanges(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}
	nvimRuntimeRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("get wd: %v", err)
	}
	root := t.TempDir()
	repoRoot := root
	socketPath := filepath.Join(repoRoot, ".operational-knowledge-system", "embodiment.sock")
	script := filepath.Join(t.TempDir(), "snapshot_default_socket.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_SOCKET_PATH = nil
local oks = require("oks")
oks.setup({ repo_root = %q })
vim.cmd("cd %s")
if oks.config.socket_path ~= %q then
  error("unexpected socket path " .. tostring(oks.config.socket_path))
end
vim.cmd("qa!")
`, repoRoot, filepath.Join(root, "elsewhere"), socketPath)
	if err := os.MkdirAll(filepath.Join(root, "elsewhere"), 0o755); err != nil {
		t.Fatalf("mkdir elsewhere: %v", err)
	}
	if err := os.WriteFile(script, []byte(scriptBody), 0o644); err != nil {
		t.Fatalf("write script: %v", err)
	}

	command := exec.Command(
		nvimPath,
		"--headless",
		"-u", "NONE",
		"-c", "set runtimepath+="+nvimRuntimeRoot,
		"-l", script,
	)
	command.Dir = filepath.Join(root, "elsewhere")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("nvim default-socket regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "unexpected socket path") {
		t.Fatalf("unexpected default socket output: %s", string(output))
	}
}

func TestNeovimPushFallsBackToHTTPAfterSocketWriteFailure(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}
	var received map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/items/ITEM-0001/live" {
			http.NotFound(writer, request)
			return
		}
		if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
			t.Fatalf("decode fallback request: %v", err)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"item_id":"ITEM-0001",
			"title":"Socket fallback",
			"status":"draft",
			"body":"after failure",
			"version":2,
			"current_revision":1,
			"participants":[]
		}`))
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "push_fallback.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
local oks = require("oks")
oks.setup()
oks.config.socket_path = ""
oks.state.item_id = "ITEM-0001"
oks.state.bufnr = vim.api.nvim_get_current_buf()
oks.state.winid = vim.api.nvim_get_current_win()
oks.state.version = 1
oks.state.current_revision = 1
oks.state.title = "Socket fallback"
oks.state.status = "draft"
vim.api.nvim_buf_set_lines(0, 0, -1, false, { "after failure" })
oks.state.socket_connected = true
oks.state.transport = "local-socket"
oks.state.socket_generation = 7
oks.state.socket = {
  write = function(_, _, cb)
    cb("broken pipe")
  end,
  read_stop = function()
  end,
  close = function()
  end,
}
oks.push()
vim.wait(1000, function()
  return oks.state.version == 2
end, 50)
if oks.state.transport ~= "http-poll" then
  error("expected http-poll fallback, got " .. tostring(oks.state.transport))
end
if oks.state.version ~= 2 then
  error("expected fallback version update")
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
		t.Fatalf("nvim push fallback regression: %v\n%s", err, string(output))
	}
	if received["update_body"] != true || received["body"] != "after failure" {
		t.Fatalf("unexpected fallback payload: %#v", received)
	}
}

func localSocketNvimRequest(t *testing.T, socketPath string, request service.LocalEmbodimentRequest) service.LocalEmbodimentResponse {
	t.Helper()
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("dial unix socket: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	if err := json.NewEncoder(conn).Encode(request); err != nil {
		t.Fatalf("encode unix socket request: %v", err)
	}
	var response service.LocalEmbodimentResponse
	if err := json.NewDecoder(conn).Decode(&response); err != nil {
		t.Fatalf("decode unix socket response: %v", err)
	}
	return response
}
