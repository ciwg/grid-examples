package service

import (
	"bufio"
	"encoding/json"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLocalEmbodimentServerHandlesRequestResponse(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewLocalEmbodimentServer(app, EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket")))
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, server.socketPath)

	response := localSocketRequest(t, server.socketPath, LocalEmbodimentRequest{
		Type:   "request",
		Method: "GET",
		Path:   "/api/meta",
	})
	if response.Type != "response" || response.Status != 200 {
		t.Fatalf("unexpected local socket response: %+v", response)
	}
	if !strings.Contains(response.Body, `"local_unix_socket_enabled":true`) {
		t.Fatalf("missing local socket capability metadata: %s", response.Body)
	}
}

func TestLocalEmbodimentServerStreamsLiveDraftUpdates(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	server := NewLocalEmbodimentServer(app, EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket")))
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, server.socketPath)

	conn, err := net.DialTimeout("unix", server.socketPath, 2*time.Second)
	if err != nil {
		t.Fatalf("dial local socket: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	if err := json.NewEncoder(conn).Encode(LocalEmbodimentRequest{
		Type:          "live-open",
		ItemID:        item.ID,
		ParticipantID: "nvim-a",
		DisplayName:   "Alice",
		Color:         "#123456",
	}); err != nil {
		t.Fatalf("encode live open: %v", err)
	}
	reader := bufio.NewReader(conn)
	initial := readLocalSocketResponse(t, reader)
	if initial.Type != "live-state" || initial.State.Body != "# Start" {
		t.Fatalf("unexpected initial live state: %+v", initial)
	}
	if err := json.NewEncoder(conn).Encode(LocalEmbodimentRequest{
		Type:          "live-update",
		ItemID:        item.ID,
		ParticipantID: "nvim-a",
		DisplayName:   "Alice",
		Color:         "#123456",
		Cursor:        4,
		Head:          4,
		BaseVersion:   initial.State.Version,
		UpdateBody:    true,
		Body:          "# Start\n\nEdited",
	}); err != nil {
		t.Fatalf("encode live update: %v", err)
	}
	updated := readLocalSocketResponse(t, reader)
	if updated.Type != "live-state" || updated.State.Body != "# Start\n\nEdited" {
		t.Fatalf("unexpected updated live state: %+v", updated)
	}
}

func waitForUnixSocket(t *testing.T, socketPath string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := net.DialTimeout("unix", socketPath, 100*time.Millisecond); err == nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("socket did not become ready: %s", socketPath)
}

func localSocketRequest(t *testing.T, socketPath string, request LocalEmbodimentRequest) LocalEmbodimentResponse {
	t.Helper()
	conn, err := net.DialTimeout("unix", socketPath, 2*time.Second)
	if err != nil {
		t.Fatalf("dial local socket: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	if err := json.NewEncoder(conn).Encode(request); err != nil {
		t.Fatalf("encode request: %v", err)
	}
	return readLocalSocketResponse(t, bufio.NewReader(conn))
}

func readLocalSocketResponse(t *testing.T, reader *bufio.Reader) LocalEmbodimentResponse {
	t.Helper()
	var response LocalEmbodimentResponse
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		t.Fatalf("decode local socket response: %v", err)
	}
	return response
}
