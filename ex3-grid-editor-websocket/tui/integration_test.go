package tui

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	charmLog "github.com/charmbracelet/log"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/service"
)

// Intent: Exercise the actual Grid TUI sidecar and model adapter against the
// browser embodiment, proving fresh-document text and awareness convergence
// without making a broader interoperability or protocol claim. Source: DI-gufuj.
func TestBrowserAndGridTUIConvergeOnFreshDocument(t *testing.T) {
	t.Chdir(integrationRepoRoot(t))
	server := integrationRelay(t)
	browser := startIntegrationProcess(t, integrationRepoRoot(t), "node", filepath.Join(integrationRepoRoot(t), "service", "testdata", "browser-harness.mjs"))
	defer browser.Close(t)
	browser.WaitForType(t, "ready")

	client, err := startSidecar(server.URL, "", "grid-tui-a", "Grid TUI A", "#22c55e", charmLog.New(io.Discard))
	if err != nil {
		t.Fatalf("start Grid TUI sidecar: %v", err)
	}
	defer client.stop()
	waitSidecarEvent(t, client, func(event sidecarEvent) bool { return event.Type == "info" })

	browser.Send(t, map[string]any{
		"type":           "connect",
		"relay_url":      server.URL,
		"participant_id": "browser-a",
		"doc_id":         "fresh-browser-grid-tui",
		"display_name":   "Browser A",
		"color":          "#1d6fd6",
	})
	browser.WaitForType(t, "opened")
	if err := client.send(map[string]any{"type": "open", "doc_id": "fresh-browser-grid-tui"}); err != nil {
		t.Fatalf("open Grid TUI document: %v", err)
	}
	waitSidecarEvent(t, client, func(event sidecarEvent) bool { return event.Type == "opened" })

	browser.Send(t, map[string]any{"type": "set_text", "content": "hello from browser"})
	changed := waitSidecarEvent(t, client, func(event sidecarEvent) bool {
		return event.Type == "changed" && event.Content == "hello from browser"
	})
	browser.Send(t, map[string]any{"type": "set_cursor", "anchor": 18, "head": 18, "typing": true})
	awareness := waitSidecarEvent(t, client, func(event sidecarEvent) bool {
		peer, ok := peerByID(event.Peers, "browser-a")
		return event.Type == "awareness" && ok && peer.Anchor == 18 && peer.Typing
	})

	state := newModel(Config{Relay: server.URL, DocumentID: "fresh-browser-grid-tui", Name: "Grid TUI A", Color: "#22c55e"}, client)
	state.applySidecar(changed)
	state.applySidecar(awareness)
	if got := state.editor.Value(); got != "hello from browser" {
		t.Fatalf("Grid TUI text = %q, want browser change", got)
	}
	peer, ok := peerByID(state.peers, "browser-a")
	if !ok || peer.Name != "Browser A" || peer.Color != "#1d6fd6" || !peer.Typing {
		t.Fatalf("Grid TUI peer state = %#v, want Browser A awareness", state.peers)
	}
	if rendered := renderLine(state.editor.Value(), state.peers, 0, state.editor.Value(), -1); rendered == state.editor.Value() {
		t.Fatal("Grid TUI did not render the browser cursor")
	}
	if sidebar := state.sidebar(); !strings.Contains(sidebar, "Browser A") || !strings.Contains(sidebar, "typing") {
		t.Fatalf("Grid TUI sidebar did not render browser awareness: %q", sidebar)
	}
}

// Intent: Exercise two real Grid TUI sidecars through the relay, proving that
// terminal-local helpers share the same document and awareness behavior. Source: DI-gufuj.
func TestGridTUIClientsConvergeOnFreshDocument(t *testing.T) {
	t.Chdir(integrationRepoRoot(t))
	server := integrationRelay(t)
	first, err := startSidecar(server.URL, "", "grid-tui-a", "Grid TUI A", "#22c55e", charmLog.New(io.Discard))
	if err != nil {
		t.Fatalf("start first Grid TUI sidecar: %v", err)
	}
	defer first.stop()
	second, err := startSidecar(server.URL, "", "grid-tui-b", "Grid TUI B", "#f97316", charmLog.New(io.Discard))
	if err != nil {
		t.Fatalf("start second Grid TUI sidecar: %v", err)
	}
	defer second.stop()
	waitSidecarEvent(t, first, func(event sidecarEvent) bool { return event.Type == "info" })
	waitSidecarEvent(t, second, func(event sidecarEvent) bool { return event.Type == "info" })

	const documentID = "fresh-grid-tui-pair"
	if err := first.send(map[string]any{"type": "open", "doc_id": documentID}); err != nil {
		t.Fatalf("open first Grid TUI document: %v", err)
	}
	if err := second.send(map[string]any{"type": "open", "doc_id": documentID}); err != nil {
		t.Fatalf("open second Grid TUI document: %v", err)
	}
	waitSidecarEvent(t, first, func(event sidecarEvent) bool { return event.Type == "opened" })
	waitSidecarEvent(t, second, func(event sidecarEvent) bool { return event.Type == "opened" })

	if err := first.send(map[string]any{"type": "set_text", "content": "hello from Grid TUI A"}); err != nil {
		t.Fatalf("send Grid TUI text: %v", err)
	}
	changed := waitSidecarEvent(t, second, func(event sidecarEvent) bool {
		return event.Type == "changed" && event.Content == "hello from Grid TUI A"
	})
	if err := first.send(map[string]any{"type": "set_cursor", "anchor": 21, "head": 21, "typing": true}); err != nil {
		t.Fatalf("send Grid TUI cursor: %v", err)
	}
	awareness := waitSidecarEvent(t, second, func(event sidecarEvent) bool {
		peer, ok := peerByID(event.Peers, "grid-tui-a")
		return event.Type == "awareness" && ok && peer.Anchor == 21 && peer.Typing
	})

	state := newModel(Config{Relay: server.URL, DocumentID: documentID, Name: "Grid TUI B", Color: "#f97316"}, second)
	state.applySidecar(changed)
	state.applySidecar(awareness)
	peer, ok := peerByID(state.peers, "grid-tui-a")
	if !ok || peer.Name != "Grid TUI A" || peer.Color != "#22c55e" || !peer.Typing {
		t.Fatalf("second Grid TUI peer state = %#v, want first Grid TUI awareness", state.peers)
	}
	if got := state.editor.Value(); got != "hello from Grid TUI A" {
		t.Fatalf("second Grid TUI text = %q, want first Grid TUI change", got)
	}
}

func integrationRelay(t *testing.T) *httptest.Server {
	t.Helper()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new relay app: %v", err)
	}
	server := httptest.NewServer(service.NewServer(app).Handler())
	t.Cleanup(server.Close)
	return server
}

func integrationRepoRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller failed")
	}
	return filepath.Dir(filepath.Dir(filename))
}

func peerByID(peers []peer, id string) (peer, bool) {
	for _, candidate := range peers {
		if candidate.ID == id {
			return candidate, true
		}
	}
	return peer{}, false
}

func waitSidecarEvent(t *testing.T, client *sidecar, predicate func(sidecarEvent) bool) sidecarEvent {
	t.Helper()
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	for {
		select {
		case event, ok := <-client.events:
			if !ok {
				t.Fatal("Grid TUI sidecar closed before the expected event")
			}
			if predicate(event) {
				return event
			}
		case <-timer.C:
			t.Fatal("timed out waiting for Grid TUI sidecar event")
		}
	}
}

type integrationProcess struct {
	cmd      *exec.Cmd
	stdin    io.WriteCloser
	messages chan map[string]any
	stderr   strings.Builder
	done     chan error
}

func startIntegrationProcess(t *testing.T, workdir, name string, args ...string) *integrationProcess {
	t.Helper()
	command := exec.Command(name, args...)
	command.Dir = workdir
	stdout, err := command.StdoutPipe()
	if err != nil {
		t.Fatalf("%s stdout pipe: %v", name, err)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		t.Fatalf("%s stderr pipe: %v", name, err)
	}
	stdin, err := command.StdinPipe()
	if err != nil {
		t.Fatalf("%s stdin pipe: %v", name, err)
	}
	process := &integrationProcess{cmd: command, stdin: stdin, messages: make(chan map[string]any, 128), done: make(chan error, 1)}
	if err := command.Start(); err != nil {
		t.Fatalf("start %s: %v", name, err)
	}
	go process.scan(stdout)
	go func() {
		if _, err := io.Copy(&process.stderr, stderr); err != nil {
			process.stderr.WriteString(fmt.Sprintf("stderr capture: %v", err))
		}
	}()
	go func() { process.done <- command.Wait() }()
	return process
}

func (process *integrationProcess) scan(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var message map[string]any
		if err := json.Unmarshal(scanner.Bytes(), &message); err != nil {
			process.messages <- map[string]any{"type": "error", "message": err.Error()}
			continue
		}
		process.messages <- message
	}
	if err := scanner.Err(); err != nil {
		process.messages <- map[string]any{"type": "error", "message": err.Error()}
	}
	close(process.messages)
}

func (process *integrationProcess) Send(t *testing.T, message any) {
	t.Helper()
	if err := json.NewEncoder(process.stdin).Encode(message); err != nil {
		t.Fatalf("send browser harness message: %v", err)
	}
}

func (process *integrationProcess) WaitForType(t *testing.T, messageType string) map[string]any {
	t.Helper()
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	for {
		select {
		case message, ok := <-process.messages:
			if !ok {
				t.Fatalf("browser harness closed early: %s", process.stderr.String())
			}
			if message["type"] == "error" {
				t.Fatalf("browser harness error: %v; stderr=%s", message["message"], process.stderr.String())
			}
			if message["type"] == messageType {
				return message
			}
		case err := <-process.done:
			t.Fatalf("browser harness exited early: %v; stderr=%s", err, process.stderr.String())
		case <-timer.C:
			t.Fatalf("timed out waiting for browser harness %q; stderr=%s", messageType, process.stderr.String())
		}
	}
}

func (process *integrationProcess) Close(t *testing.T) {
	t.Helper()
	if err := process.stdin.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
		t.Errorf("close browser harness stdin: %v", err)
	}
	if process.cmd.Process != nil {
		if err := process.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			t.Errorf("stop browser harness: %v", err)
		}
	}
	select {
	case <-process.done:
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for browser harness shutdown")
	}
}
