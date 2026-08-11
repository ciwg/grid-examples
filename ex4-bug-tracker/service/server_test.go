package service_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/service"
)

func TestServerDisablesCachingForIndex(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	assertHeaderEquals(t, response, "Cache-Control", "no-store, max-age=0")
	assertHeaderEquals(t, response, "Pragma", "no-cache")
	assertHeaderEquals(t, response, "Expires", "0")
}

func TestServerServesBrowserSignerAdapter(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)
	request := httptest.NewRequest(http.MethodGet, "/promise.js", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), "submitPromise") {
		t.Fatal("served signer adapter omitted submitPromise")
	}
}

func TestServerRejectsUnsignedIssueMutation(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	create := httptest.NewRequest(http.MethodPost, "/api/issues", strings.NewReader(`{"title":"Crash","description":"App crashes on upload.","severity":"High"}`))
	create.Header.Set("X-Bug-User", "reporter")
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, create)
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("create status = %d, want %d body=%s", response.Code, http.StatusMethodNotAllowed, response.Body.String())
	}
}

func TestServerRejectsUnsignedAttachmentReference(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)
	upload := httptest.NewRequest(http.MethodPost, "/api/issues/BUG-0001/attachments", strings.NewReader("legacy"))
	uploadResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(uploadResponse, upload)
	if uploadResponse.Code != http.StatusMethodNotAllowed {
		t.Fatalf("upload status = %d, want %d body=%s", uploadResponse.Code, http.StatusMethodNotAllowed, uploadResponse.Body.String())
	}
}

func newTestServer(t *testing.T) *service.Server {
	t.Helper()
	app, err := service.NewApp(filepath.Join(t.TempDir(), ".bug-tracker"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	return service.NewServer(app)
}

func assertHeaderEquals(t *testing.T, response *httptest.ResponseRecorder, key string, want string) {
	t.Helper()
	if got := response.Header().Get(key); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}
