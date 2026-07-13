package service_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/service"
)

func TestServerRejectsRemoteSyncMutation(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/sync", bytes.NewBufferString(`{"participant_id":"browser-a","recipient_id":"","message_base64":"AQID","embodiment":"browser"}`))
	request.RemoteAddr = "198.51.100.20:4123"
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: got %d want %d", response.Code, http.StatusForbidden)
	}
}

func TestServerAllowsLoopbackSyncMutation(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/sync", bytes.NewBufferString(`{"participant_id":"browser-a","recipient_id":"","message_base64":"AQID","embodiment":"browser"}`))
	request.RemoteAddr = "127.0.0.1:4123"
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d body=%s", response.Code, http.StatusOK, response.Body.String())
	}
}

func TestServerRejectsMissingParticipantID(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/awareness", bytes.NewBufferString(`{"participant_id":"","cursor":0,"head":0,"typing":false,"display_name":"Alice","color":"#1d6fd6","embodiment":"browser"}`))
	request.RemoteAddr = "127.0.0.1:4123"
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d want %d", response.Code, http.StatusBadRequest)
	}
}

func newTestServer(t *testing.T) *service.Server {
	t.Helper()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	return service.NewServer(app)
}
