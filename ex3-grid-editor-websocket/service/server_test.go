package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/service"
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

func TestServerIssuesRemoteSessionWithBootstrapToken(t *testing.T) {
	t.Parallel()
	server := newTestServerWithOptions(t, service.AppOptions{RemoteAccessToken: "ex3-demo-access"})

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/session", bytes.NewBufferString(`{"participant_id":"browser-a"}`))
	request.RemoteAddr = "198.51.100.20:4123"
	request.Header.Set("X-Grid-Access-Token", "ex3-demo-access")
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d body=%s", response.Code, http.StatusOK, response.Body.String())
	}
	assertBodyContains(t, response.Body.String(), `"participant_id":"browser-a"`)
	assertBodyContains(t, response.Body.String(), `"sync":"`)
	assertBodyContains(t, response.Body.String(), `"awareness":"`)
}

func TestServerAllowsRemoteSyncMutationWithCapability(t *testing.T) {
	t.Parallel()
	server := newTestServerWithOptions(t, service.AppOptions{RemoteAccessToken: "ex3-demo-access"})

	sessionRequest := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/session", bytes.NewBufferString(`{"participant_id":"browser-a"}`))
	sessionRequest.RemoteAddr = "198.51.100.20:4123"
	sessionRequest.Header.Set("X-Grid-Access-Token", "ex3-demo-access")
	sessionResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(sessionResponse, sessionRequest)
	if sessionResponse.Code != http.StatusOK {
		t.Fatalf("unexpected session status: got %d want %d body=%s", sessionResponse.Code, http.StatusOK, sessionResponse.Body.String())
	}
	var sessionPayload struct {
		Capabilities map[string]string `json:"capabilities"`
	}
	if err := json.Unmarshal(sessionResponse.Body.Bytes(), &sessionPayload); err != nil {
		t.Fatalf("decode session payload: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/sync", bytes.NewBufferString(`{"participant_id":"browser-a","recipient_id":"","message_base64":"AQID","embodiment":"browser"}`))
	request.RemoteAddr = "198.51.100.20:4123"
	request.Header.Set("Authorization", "Bearer "+sessionPayload.Capabilities["sync"])
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d body=%s", response.Code, http.StatusOK, response.Body.String())
	}
}

func TestServerRejectsRemoteSyncSocketUpgrade(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodGet, "/api/local/documents/demo/sync-socket?since=0", nil)
	request.RemoteAddr = "198.51.100.20:4123"
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Sec-WebSocket-Version", "13")
	request.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code == http.StatusForbidden {
		t.Fatalf("unexpected status: got %d want non-%d", response.Code, http.StatusForbidden)
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

func TestServerRejectsRemoteMetadataMutation(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/metadata", bytes.NewBufferString(`{
		"participant_id":"browser-a",
		"title":"Demo title",
		"description":"Description",
		"summary":"Summary",
		"tags":["grid"],
		"collections":["team"],
		"favorite":true,
		"archived":false,
		"embodiment":"browser"
	}`))
	request.RemoteAddr = "198.51.100.20:4123"
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: got %d want %d", response.Code, http.StatusForbidden)
	}
}

func TestServerStoresAndReadsMetadata(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodPost, "/api/local/documents/demo/metadata", bytes.NewBufferString(`{
		"participant_id":"browser-a",
		"title":"Demo title",
		"description":"Description",
		"summary":"Summary",
		"tags":["grid","docs"],
		"collections":["team"],
		"favorite":true,
		"archived":false,
		"embodiment":"browser"
	}`))
	request.RemoteAddr = "127.0.0.1:4123"
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected metadata post status: got %d want %d body=%s", response.Code, http.StatusOK, response.Body.String())
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/local/documents/demo/metadata", nil)
	getResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(getResponse, getRequest)

	if getResponse.Code != http.StatusOK {
		t.Fatalf("unexpected metadata get status: got %d want %d body=%s", getResponse.Code, http.StatusOK, getResponse.Body.String())
	}
	assertBodyContains(t, getResponse.Body.String(), `"title":"Demo title"`)
	assertBodyContains(t, getResponse.Body.String(), `"favorite":true`)
}

func TestServerSearchesMetadata(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	post := func(documentID string, body string) {
		t.Helper()
		request := httptest.NewRequest(http.MethodPost, "/api/local/documents/"+documentID+"/metadata", bytes.NewBufferString(body))
		request.RemoteAddr = "127.0.0.1:4123"
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("unexpected metadata post status for %s: got %d want %d body=%s", documentID, response.Code, http.StatusOK, response.Body.String())
		}
	}

	post("demo", `{
		"participant_id":"browser-a",
		"title":"Demo title",
		"description":"Description",
		"summary":"Summary",
		"tags":["grid"],
		"collections":["team"],
		"favorite":true,
		"archived":false,
		"embodiment":"browser"
	}`)
	post("archive", `{
		"participant_id":"browser-b",
		"title":"Archive title",
		"description":"Description",
		"summary":"Summary",
		"tags":["grid"],
		"collections":["team"],
		"favorite":false,
		"archived":true,
		"embodiment":"browser"
	}`)

	request := httptest.NewRequest(http.MethodGet, "/api/local/metadata/search?q=grid", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected metadata search status: got %d want %d body=%s", response.Code, http.StatusOK, response.Body.String())
	}
	assertBodyContains(t, response.Body.String(), `"document_id":"demo"`)
	if strings.Contains(response.Body.String(), `"document_id":"archive"`) {
		t.Fatalf("archived metadata unexpectedly returned without include_archived=true: %s", response.Body.String())
	}

	withArchived := httptest.NewRequest(http.MethodGet, "/api/local/metadata/search?q=grid&include_archived=true", nil)
	withArchivedResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(withArchivedResponse, withArchived)
	if withArchivedResponse.Code != http.StatusOK {
		t.Fatalf("unexpected metadata search include_archived status: got %d want %d body=%s", withArchivedResponse.Code, http.StatusOK, withArchivedResponse.Body.String())
	}
	assertBodyContains(t, withArchivedResponse.Body.String(), `"document_id":"archive"`)
}

func newTestServer(t *testing.T) *service.Server {
	return newTestServerWithOptions(t, service.AppOptions{})
}

func newTestServerWithOptions(t *testing.T, options service.AppOptions) *service.Server {
	t.Helper()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"), options)
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
