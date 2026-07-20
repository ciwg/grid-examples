package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunRecordCommandParsesItemRevisionOutcomeAndNotes(t *testing.T) {
	var received map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", request.Method)
		}
		if request.URL.Path != "/api/runs" {
			t.Fatalf("unexpected path: %s", request.URL.Path)
		}
		if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"record-run", "bob", "procedure", "PROC-0001", "1", "completed", "Completed startup cleanly"})
	if err != nil {
		t.Fatalf("run command: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if got := received["actor"]; got != "bob" {
		t.Fatalf("unexpected actor: %#v", got)
	}
	if got := received["kind"]; got != "procedure" {
		t.Fatalf("unexpected kind: %#v", got)
	}
	if got := received["item_id"]; got != "PROC-0001" {
		t.Fatalf("unexpected item_id: %#v", got)
	}
	if got := received["revision"]; got != float64(1) {
		t.Fatalf("unexpected revision: %#v", got)
	}
	if got := received["outcome"]; got != "completed" {
		t.Fatalf("unexpected outcome: %#v", got)
	}
	if got := received["notes"]; got != "Completed startup cleanly" {
		t.Fatalf("unexpected notes: %#v", got)
	}
}

func TestRunSearchRequiresQuery(t *testing.T) {
	cli := &CLI{ServerURL: "http://127.0.0.1:7045"}
	exitCode, err := cli.run([]string{"search"})
	if exitCode != 2 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if err == nil {
		t.Fatalf("expected usage error")
	}
	if !strings.Contains(err.Error(), "usage: oks-cli") {
		t.Fatalf("unexpected error: %v", err)
	}
}
