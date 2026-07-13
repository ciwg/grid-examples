package service_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/service"
)

func TestBrowserAndNvimInteroperateThroughRelay(t *testing.T) {
	t.Parallel()

	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := httptest.NewServer(service.NewServer(app).Handler())
	defer server.Close()

	repoRoot := repoRoot(t)
	browser := startJSONProcess(t, repoRoot, "node", filepath.Join(repoRoot, "service", "testdata", "browser-harness.mjs"))
	defer browser.Close()
	sidecar := startJSONProcess(t, repoRoot, "node", filepath.Join(repoRoot, "cmd", "grid-nvim-sidecar", "helper.bundle.cjs"), "--relay", server.URL)
	defer sidecar.Close()

	browser.WaitForType(t, "ready")
	sidecar.WaitForType(t, "info")

	browser.Send(t, map[string]any{
		"type":           "connect",
		"relay_url":      server.URL,
		"participant_id": "browser-a",
		"doc_id":         "demo",
		"display_name":   "Browser A",
		"color":          "#1d6fd6",
	})
	browserOpened := browser.WaitForType(t, "opened")
	if got := stringField(t, browserOpened, "content"); got != "" {
		t.Fatalf("browser opened with unexpected content %q", got)
	}

	sidecar.Send(t, map[string]any{
		"type":           "connect",
		"relay_url":      server.URL,
		"participant_id": "nvim-a",
		"display_name":   "Nvim A",
		"color":          "#d66f1d",
	})
	sidecar.WaitForType(t, "connected")
	sidecar.Send(t, map[string]any{
		"type":   "open",
		"doc_id": "demo",
	})
	sidecarOpened := sidecar.WaitForType(t, "opened")
	if got := stringField(t, sidecarOpened, "content"); got != "" {
		t.Fatalf("sidecar opened with unexpected content %q", got)
	}

	browser.Send(t, map[string]any{
		"type":    "set_text",
		"content": "hello from browser",
	})
	sidecarChanged := sidecar.WaitForMessage(t, func(message map[string]any) bool {
		return stringField(t, message, "type") == "changed" && stringField(t, message, "content") == "hello from browser"
	})
	if got := stringField(t, sidecarChanged, "content"); got != "hello from browser" {
		t.Fatalf("unexpected sidecar content %q", got)
	}

	sidecar.Send(t, map[string]any{
		"type":    "set_text",
		"content": "hello from browser\nand nvim",
	})
	browserChanged := browser.WaitForMessage(t, func(message map[string]any) bool {
		return stringField(t, message, "type") == "document" && stringField(t, message, "content") == "hello from browser\nand nvim"
	})
	if got := stringField(t, browserChanged, "content"); got != "hello from browser\nand nvim" {
		t.Fatalf("unexpected browser content %q", got)
	}

	browser.Send(t, map[string]any{
		"type":   "set_cursor",
		"anchor": 3,
		"head":   3,
		"typing": true,
	})
	sidecarAwareness := sidecar.WaitForMessage(t, func(message map[string]any) bool {
		return stringField(t, message, "type") == "awareness" && hasPeer(message, "browser-a", 3)
	})
	if !hasPeer(sidecarAwareness, "browser-a", 3) {
		t.Fatalf("browser peer did not appear in sidecar awareness: %#v", sidecarAwareness)
	}

	sidecar.Send(t, map[string]any{
		"type":   "set_cursor",
		"anchor": 7,
		"head":   7,
		"typing": false,
	})
	browserAwareness := browser.WaitForMessage(t, func(message map[string]any) bool {
		return stringField(t, message, "type") == "awareness" && hasPeer(message, "nvim-a", 7)
	})
	if !hasPeer(browserAwareness, "nvim-a", 7) {
		t.Fatalf("nvim peer did not appear in browser awareness: %#v", browserAwareness)
	}
}

type jsonProcess struct {
	t        *testing.T
	cmd      *exec.Cmd
	stdin    io.WriteCloser
	messages chan map[string]any
	stderr   strings.Builder
	done     chan error
}

func startJSONProcess(t *testing.T, workdir string, name string, args ...string) *jsonProcess {
	t.Helper()

	command := exec.Command(name, args...)
	command.Dir = workdir
	stdout, err := command.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}
	stdin, err := command.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	process := &jsonProcess{
		t:        t,
		cmd:      command,
		stdin:    stdin,
		messages: make(chan map[string]any, 128),
		done:     make(chan error, 1),
	}
	if err := command.Start(); err != nil {
		t.Fatalf("start %s: %v", name, err)
	}
	go process.scanStdout(stdout)
	go process.captureStderr(stderr)
	go func() {
		process.done <- command.Wait()
		close(process.done)
	}()
	return process
}

func (process *jsonProcess) scanStdout(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		var message map[string]any
		if err := json.Unmarshal([]byte(line), &message); err != nil {
			process.messages <- map[string]any{
				"type":    "error",
				"message": fmt.Sprintf("invalid json from process: %v", err),
			}
			continue
		}
		process.messages <- message
	}
	if err := scanner.Err(); err != nil {
		process.messages <- map[string]any{
			"type":    "error",
			"message": fmt.Sprintf("stdout scan error: %v", err),
		}
	}
}

func (process *jsonProcess) captureStderr(reader io.Reader) {
	_, _ = io.Copy(&process.stderr, reader)
}

func (process *jsonProcess) Send(t *testing.T, value any) {
	t.Helper()
	if err := json.NewEncoder(process.stdin).Encode(value); err != nil {
		t.Fatalf("send message: %v", err)
	}
}

func (process *jsonProcess) WaitForType(t *testing.T, messageType string) map[string]any {
	t.Helper()
	return process.WaitForMessage(t, func(message map[string]any) bool {
		return stringField(t, message, "type") == messageType
	})
}

func (process *jsonProcess) WaitForMessage(t *testing.T, predicate func(map[string]any) bool) map[string]any {
	t.Helper()

	timeout := time.NewTimer(8 * time.Second)
	defer timeout.Stop()
	for {
		select {
		case message := <-process.messages:
			if message == nil {
				t.Fatalf("process closed without expected message; stderr=%s", process.stderr.String())
			}
			if stringField(t, message, "type") == "error" {
				t.Fatalf("process error: %s; stderr=%s", stringField(t, message, "message"), process.stderr.String())
			}
			if predicate(message) {
				return message
			}
		case err := <-process.done:
			t.Fatalf("process exited early: %v; stderr=%s", err, process.stderr.String())
		case <-timeout.C:
			t.Fatalf("timed out waiting for process message; stderr=%s", process.stderr.String())
		}
	}
}

func (process *jsonProcess) Close() {
	if process == nil {
		return
	}
	_ = process.stdin.Close()
	if process.cmd.Process != nil {
		_ = process.cmd.Process.Kill()
	}
	<-process.done
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime caller failed")
	}
	return filepath.Dir(filepath.Dir(filename))
}

func stringField(t *testing.T, message map[string]any, key string) string {
	t.Helper()
	value, _ := message[key]
	if value == nil {
		return ""
	}
	text, ok := value.(string)
	if !ok {
		t.Fatalf("field %q is not a string: %#v", key, value)
	}
	return text
}

func hasPeer(message map[string]any, participantID string, anchor float64) bool {
	peersValue, ok := message["peers"]
	if !ok {
		return false
	}
	peers, ok := peersValue.([]any)
	if !ok {
		return false
	}
	for _, rawPeer := range peers {
		peer, ok := rawPeer.(map[string]any)
		if !ok {
			continue
		}
		if peer["participant_id"] == participantID && peer["anchor"] == anchor {
			return true
		}
	}
	return false
}
