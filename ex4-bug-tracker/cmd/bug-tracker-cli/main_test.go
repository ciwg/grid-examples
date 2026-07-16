package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/service"
)

func TestCLIFlowsAgainstServerHandler(t *testing.T) {
	t.Parallel()
	server := newCLITestServer(t)
	app := server.app
	issue, err := app.CreateIssue("reporter", "CLI issue", "Issue created for CLI testing.", service.SeverityMedium)
	if err != nil {
		t.Fatalf("create issue: %v", err)
	}
	if _, err := app.ChangeStatus("triage", issue.ID, service.StatusTriaged); err != nil {
		t.Fatalf("triage issue: %v", err)
	}
	if _, err := app.AssignIssue("triage", issue.ID, "engineer"); err != nil {
		t.Fatalf("assign issue: %v", err)
	}

	cli := &CLI{
		ServerURL:  "http://example.test",
		User:       "engineer",
		HTTPClient: server.client(),
	}

	assignedOutput := captureStdout(t, func() {
		if err := cli.Assigned(); err != nil {
			t.Fatalf("assigned: %v", err)
		}
	})
	if !strings.Contains(assignedOutput, issue.ID) {
		t.Fatalf("assigned output = %q, want issue id %s", assignedOutput, issue.ID)
	}

	showOutput := captureStdout(t, func() {
		if err := cli.Show(issue.ID); err != nil {
			t.Fatalf("show: %v", err)
		}
	})
	if !strings.Contains(showOutput, "Title: CLI issue") {
		t.Fatalf("show output = %q, want title", showOutput)
	}

	if err := cli.ChangeStatus(issue.ID, service.StatusInProgress); err != nil {
		t.Fatalf("start issue: %v", err)
	}
	updated, err := app.GetIssue(issue.ID)
	if err != nil {
		t.Fatalf("get issue after start: %v", err)
	}
	if updated.Status != service.StatusInProgress {
		t.Fatalf("status after cli start = %q, want %q", updated.Status, service.StatusInProgress)
	}

	if err := cli.Comment(issue.ID, "cli integration note"); err != nil {
		t.Fatalf("comment issue: %v", err)
	}
	updated, err = app.GetIssue(issue.ID)
	if err != nil {
		t.Fatalf("get issue after comment: %v", err)
	}
	last := updated.Timeline[len(updated.Timeline)-1]
	if last.Type != "commented" || last.Comment != "cli integration note" {
		t.Fatalf("last event = %#v, want commented cli note", last)
	}
}

type cliTestServer struct {
	app     *service.App
	handler http.Handler
}

func newCLITestServer(t *testing.T) *cliTestServer {
	t.Helper()
	app, err := service.NewApp(filepath.Join(t.TempDir(), ".bug-tracker"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := service.NewServer(app)
	return &cliTestServer{app: app, handler: server.Handler()}
}

func (server *cliTestServer) client() *http.Client {
	return &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()
		server.handler.ServeHTTP(recorder, request)
		return recorder.Result(), nil
	})}
}

type roundTripFunc func(request *http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	os.Stdout = writer
	defer func() {
		os.Stdout = oldStdout
	}()
	fn()
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, reader); err != nil {
		t.Fatalf("copy stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close reader: %v", err)
	}
	return buffer.String()
}
