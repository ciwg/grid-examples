package service_test

import (
	"bytes"
	"mime/multipart"
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

func TestServerCreatesAndFetchesIssue(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	create := httptest.NewRequest(http.MethodPost, "/api/issues", strings.NewReader(`{"title":"Crash","description":"App crashes on upload.","severity":"High"}`))
	create.Header.Set("X-Bug-User", "reporter")
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, create)
	if response.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d body=%s", response.Code, http.StatusCreated, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"id":"BUG-0001"`) {
		t.Fatalf("create body = %s, want BUG-0001", response.Body.String())
	}

	get := httptest.NewRequest(http.MethodGet, "/api/issues/BUG-0001", nil)
	getResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(getResponse, get)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d body=%s", getResponse.Code, http.StatusOK, getResponse.Body.String())
	}
	if !strings.Contains(getResponse.Body.String(), `"status":"New"`) {
		t.Fatalf("get body = %s, want status New", getResponse.Body.String())
	}
}

func TestServerUploadsAndDownloadsAttachment(t *testing.T) {
	t.Parallel()
	server := newTestServer(t)

	create := httptest.NewRequest(http.MethodPost, "/api/issues", strings.NewReader(`{"title":"Crash","description":"App crashes on upload.","severity":"High"}`))
	create.Header.Set("X-Bug-User", "reporter")
	createResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(createResponse, create)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d body=%s", createResponse.Code, http.StatusCreated, createResponse.Body.String())
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("attachment", "trace.log")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write([]byte("stack trace")); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	upload := httptest.NewRequest(http.MethodPost, "/api/issues/BUG-0001/attachments", body)
	upload.Header.Set("Content-Type", writer.FormDataContentType())
	upload.Header.Set("X-Bug-User", "reporter")
	uploadResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(uploadResponse, upload)
	if uploadResponse.Code != http.StatusOK {
		t.Fatalf("upload status = %d, want %d body=%s", uploadResponse.Code, http.StatusOK, uploadResponse.Body.String())
	}
	if !strings.Contains(uploadResponse.Body.String(), `"attachment_id":"ATT-000002"`) {
		t.Fatalf("upload body = %s, want ATT-000002", uploadResponse.Body.String())
	}

	download := httptest.NewRequest(http.MethodGet, "/api/issues/BUG-0001/attachments/ATT-000002?user=reporter", nil)
	downloadResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(downloadResponse, download)
	if downloadResponse.Code != http.StatusOK {
		t.Fatalf("download status = %d, want %d body=%s", downloadResponse.Code, http.StatusOK, downloadResponse.Body.String())
	}
	if body := downloadResponse.Body.String(); body != "stack trace" {
		t.Fatalf("download body = %q, want stack trace", body)
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
