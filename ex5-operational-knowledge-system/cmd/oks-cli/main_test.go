package main

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestSearchCommandEncodesQueryString(t *testing.T) {
	var (
		rawQuery string
		query    string
	)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		rawQuery = request.URL.RawQuery
		query = request.URL.Query().Get("q")
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"search", "supplier: Acme Parts & variance=-2"})
	if err != nil {
		t.Fatalf("search command: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if query != "supplier: Acme Parts & variance=-2" {
		t.Fatalf("unexpected decoded query: %q", query)
	}
	if !strings.Contains(rawQuery, "q=supplier%3A+Acme+Parts+%26+variance%3D-2") {
		t.Fatalf("unexpected raw query: %q", rawQuery)
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

func TestApproveItemCommandUsesExplicitActor(t *testing.T) {
	var (
		received map[string]any
		path     string
	)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path = request.URL.Path
		if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"approve-item", "ITEM-0001", "2", "carol", "reviewer", "approved", "ready for use"})
	if err != nil {
		t.Fatalf("approve-item command: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if path != "/api/items/ITEM-0001/approvals" {
		t.Fatalf("unexpected path: %s", path)
	}
	if received["actor"] != "carol" || received["revision"] != float64(2) || received["role"] != "reviewer" || received["decision"] != "approved" {
		t.Fatalf("unexpected approve-item payload: %#v", received)
	}
	if received["notes"] != "ready for use" {
		t.Fatalf("unexpected approve-item notes: %#v", received["notes"])
	}
}

func TestApproveRunCommandUsesExplicitActor(t *testing.T) {
	var (
		received map[string]any
		path     string
	)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path = request.URL.Path
		if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"approve-run", "RUN-0001", "dave", "approver", "noted", "handoff recorded"})
	if err != nil {
		t.Fatalf("approve-run command: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if path != "/api/runs/RUN-0001/approvals" {
		t.Fatalf("unexpected path: %s", path)
	}
	if received["actor"] != "dave" || received["revision"] != float64(0) || received["role"] != "approver" || received["decision"] != "noted" {
		t.Fatalf("unexpected approve-run payload: %#v", received)
	}
	if received["notes"] != "handoff recorded" {
		t.Fatalf("unexpected approve-run notes: %#v", received["notes"])
	}
}

func TestAddEvidenceCommandSupportsFactsAndAttachment(t *testing.T) {
	var (
		path           string
		contentType    string
		fields         map[string]string
		attachmentName string
		attachmentBody string
	)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path = request.URL.Path
		contentType = request.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			t.Fatalf("parse content type: %v", err)
		}
		if mediaType != "multipart/form-data" {
			t.Fatalf("unexpected media type: %s", mediaType)
		}
		reader := multipart.NewReader(request.Body, params["boundary"])
		fields = map[string]string{}
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("read multipart part: %v", err)
			}
			body, err := io.ReadAll(part)
			if err != nil {
				t.Fatalf("read part body: %v", err)
			}
			if part.FormName() == "attachment" {
				attachmentName = part.FileName()
				attachmentBody = string(body)
				continue
			}
			fields[part.FormName()] = string(body)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	attachmentPath := filepath.Join(tempDir, "evidence.txt")
	if err := os.WriteFile(attachmentPath, []byte("hello evidence"), 0o644); err != nil {
		t.Fatalf("write attachment: %v", err)
	}

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"add-evidence", "RUN-0001", "dave", "dock photo", `{"result":"ok"}`, attachmentPath})
	if err != nil {
		t.Fatalf("add-evidence command: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if path != "/api/runs/RUN-0001/evidence" {
		t.Fatalf("unexpected path: %s", path)
	}
	if fields["actor"] != "dave" || fields["summary"] != "dock photo" || fields["facts_json"] != `{"result":"ok"}` {
		t.Fatalf("unexpected multipart fields: %#v", fields)
	}
	if attachmentName != "evidence.txt" || attachmentBody != "hello evidence" {
		t.Fatalf("unexpected attachment payload: %q %q", attachmentName, attachmentBody)
	}
}

func TestAddEvidenceCommandSupportsSummaryOnly(t *testing.T) {
	var (
		path           string
		fields         map[string]string
		attachmentSeen bool
	)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path = request.URL.Path
		mediaType, params, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
		if err != nil {
			t.Fatalf("parse content type: %v", err)
		}
		if mediaType != "multipart/form-data" {
			t.Fatalf("unexpected media type: %s", mediaType)
		}
		reader := multipart.NewReader(request.Body, params["boundary"])
		fields = map[string]string{}
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("read multipart part: %v", err)
			}
			body, err := io.ReadAll(part)
			if err != nil {
				t.Fatalf("read part body: %v", err)
			}
			if part.FormName() == "attachment" {
				attachmentSeen = true
			}
			fields[part.FormName()] = string(body)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"add-evidence", "RUN-0002", "alice", "verbal handoff"})
	if err != nil {
		t.Fatalf("add-evidence summary-only command: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if path != "/api/runs/RUN-0002/evidence" {
		t.Fatalf("unexpected path: %s", path)
	}
	if fields["actor"] != "alice" || fields["summary"] != "verbal handoff" {
		t.Fatalf("unexpected fields: %#v", fields)
	}
	if _, ok := fields["facts_json"]; ok {
		t.Fatalf("did not expect facts_json field: %#v", fields)
	}
	if attachmentSeen {
		t.Fatalf("did not expect attachment part")
	}
}

func TestApproveCommandsRequireExplicitActor(t *testing.T) {
	cli := &CLI{ServerURL: "http://127.0.0.1:7045"}
	exitCode, err := cli.run([]string{"approve-item", "ITEM-0001", "1", "reviewer", "approved", "missing actor"})
	if exitCode != 2 {
		t.Fatalf("unexpected approve-item exit code: %d", exitCode)
	}
	if err == nil || !strings.Contains(err.Error(), "usage: oks-cli") {
		t.Fatalf("unexpected approve-item error: %v", err)
	}

	exitCode, err = cli.run([]string{"approve-run", "RUN-0001", "approver", "approved", "missing actor"})
	if exitCode != 2 {
		t.Fatalf("unexpected approve-run exit code: %d", exitCode)
	}
	if err == nil || !strings.Contains(err.Error(), "usage: oks-cli") {
		t.Fatalf("unexpected approve-run error: %v", err)
	}
}
