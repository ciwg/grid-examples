package nvim

import (
	"encoding/json"
	"fmt"
	"net"
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
