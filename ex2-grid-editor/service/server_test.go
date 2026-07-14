package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
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

func TestServerRejectsRemotePublishMutation(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/publish", bytes.NewBufferString(`{
		"participant_id":"browser-a",
		"source_kind":"current",
		"title":"Demo exchange",
		"text_base64":"IyBkZW1vCgpoZWxsbw==",
		"replica_base64":"CQgHBg==",
		"embodiment":"browser"
	}`))
	request.RemoteAddr = "198.51.100.20:4123"
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: got %d want %d", response.Code, http.StatusForbidden)
	}
}

func TestServerPublishedListStartsEmpty(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodGet, "/api/local/documents/demo/published", nil)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", response.Code, http.StatusOK)
	}
	assertBodyContains(t, response.Body.String(), `"published":[]`)
}

func TestServerRejectsInvalidPublishBase64(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/publish", bytes.NewBufferString(`{
		"participant_id":"browser-a",
		"source_kind":"current",
		"title":"Demo exchange",
		"text_base64":"%%%bad%%%",
		"replica_base64":"CQgHBg==",
		"embodiment":"browser"
	}`))
	request.RemoteAddr = "127.0.0.1:4123"
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d want %d body=%s", response.Code, http.StatusBadRequest, response.Body.String())
	}
}

func TestServerRejectsPublishWithoutTitle(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/publish", bytes.NewBufferString(`{
		"participant_id":"browser-a",
		"source_kind":"current",
		"title":"",
		"text_base64":"IyBkZW1vCgpoZWxsbw==",
		"replica_base64":"CQgHBg==",
		"embodiment":"browser"
	}`))
	request.RemoteAddr = "127.0.0.1:4123"
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d want %d body=%s", response.Code, http.StatusBadRequest, response.Body.String())
	}
}

func TestServerPublishesAndResolvesExchangeManifest(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/publish", bytes.NewBufferString(`{
		"participant_id":"browser-a",
		"source_kind":"current",
		"title":"Demo exchange",
		"summary":"Phase 4 smoke",
		"text_base64":"IyBkZW1vCgpoZWxsbw==",
		"replica_base64":"CQgHBg==",
		"embodiment":"browser"
	}`))
	request.RemoteAddr = "127.0.0.1:4123"
	request.Host = "127.0.0.1:7017"
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected publish status: got %d want %d body=%s", response.Code, http.StatusOK, response.Body.String())
	}

	var publishPayload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &publishPayload); err != nil {
		t.Fatalf("decode publish payload: %v", err)
	}
	manifestURL, ok := publishPayload["manifest_url"].(string)
	if !ok || manifestURL == "" {
		t.Fatalf("missing manifest url in publish payload: %#v", publishPayload)
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/api/local/documents/demo/published", nil)
	listResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(listResponse, listRequest)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("unexpected published list status: got %d want %d body=%s", listResponse.Code, http.StatusOK, listResponse.Body.String())
	}

	resolveRequest := httptest.NewRequest(http.MethodGet, "/api/published/"+manifestURL[strings.LastIndex(manifestURL, "/")+1:], nil)
	resolveResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(resolveResponse, resolveRequest)
	if resolveResponse.Code != http.StatusOK {
		t.Fatalf("unexpected resolve status: got %d want %d body=%s", resolveResponse.Code, http.StatusOK, resolveResponse.Body.String())
	}

	var resolvedPayload map[string]any
	if err := json.Unmarshal(resolveResponse.Body.Bytes(), &resolvedPayload); err != nil {
		t.Fatalf("decode resolved payload: %v", err)
	}
	if got, _ := resolvedPayload["text_base64"].(string); got != "IyBkZW1vCgpoZWxsbw==" {
		t.Fatalf("resolved text base64 mismatch: got %q", got)
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

func assertBodyContains(t *testing.T, body string, want string) {
	t.Helper()
	if !strings.Contains(body, want) {
		t.Fatalf("response body %q does not contain %q", body, want)
	}
}
