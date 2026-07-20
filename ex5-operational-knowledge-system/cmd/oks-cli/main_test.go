package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunRecordCommandParsesItemRevisionOutcomeNotesAndContext(t *testing.T) {
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
	exitCode, err := cli.run([]string{"record-run", "bob", "procedure", "PROC-0001", "1", "completed", "Completed startup cleanly", "PLACE-0001", "RES-0001,RES-0002"})
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
	if got := received["place_id"]; got != "PLACE-0001" {
		t.Fatalf("unexpected place_id: %#v", got)
	}
	resourceIDs, ok := received["resource_ids"].([]any)
	if !ok || len(resourceIDs) != 2 || resourceIDs[0] != "RES-0001" || resourceIDs[1] != "RES-0002" {
		t.Fatalf("unexpected resource_ids: %#v", received["resource_ids"])
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

func TestNewPlaceAndResourceCommandsEmitExpectedPayloads(t *testing.T) {
	requests := []map[string]any{}
	paths := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		paths = append(paths, request.URL.Path)
		requests = append(requests, payload)
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	if exitCode, err := cli.run([]string{"new-place", "alice", "area", "Receiving", "Inbound area", "PLACE-0000"}); err != nil || exitCode != 0 {
		t.Fatalf("new-place failed: exit=%d err=%v", exitCode, err)
	}
	if exitCode, err := cli.run([]string{"new-resource", "alice", "container", "RJ45 Bin", "Connectors bin", "PLACE-0001"}); err != nil || exitCode != 0 {
		t.Fatalf("new-resource failed: exit=%d err=%v", exitCode, err)
	}

	if len(paths) != 2 || paths[0] != "/api/places" || paths[1] != "/api/resources" {
		t.Fatalf("unexpected paths: %#v", paths)
	}
	if requests[0]["parent_id"] != "PLACE-0000" || requests[0]["name"] != "Receiving" {
		t.Fatalf("unexpected place payload: %#v", requests[0])
	}
	if requests[1]["place_id"] != "PLACE-0001" || requests[1]["kind"] != "container" {
		t.Fatalf("unexpected resource payload: %#v", requests[1])
	}
}
