package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
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

func TestProblemReviewCommandUsesExpectedRoute(t *testing.T) {
	var path string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path = request.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"place_groups":[],"resource_groups":[],"problem_runs":[]}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"problem-review"})
	if err != nil {
		t.Fatalf("problem-review command: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if path != "/api/problem-review" {
		t.Fatalf("unexpected path: %s", path)
	}
}

func TestPendingReviewCommandUsesSharedSearchRoutes(t *testing.T) {
	paths := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		paths = append(paths, request.URL.RequestURI())
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.RawQuery {
		case "status=draft":
			_, _ = writer.Write([]byte(`{"items":[{"id":"ITEM-0001"}]}`))
		case "":
			_, _ = writer.Write([]byte(`{"runs":[{"id":"RUN-0001","approvals":[]},{"id":"RUN-0009","approvals":[{"actor":"carol"}]}]}`))
		case "problem=true":
			_, _ = writer.Write([]byte(`{"runs":[{"id":"RUN-0002"}]}`))
		default:
			t.Fatalf("unexpected query: %q", request.URL.RawQuery)
		}
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	stdout, restoreStdout, err := captureStdout(t)
	if err != nil {
		t.Fatalf("capture stdout: %v", err)
	}

	exitCode, runErr := cli.run([]string{"pending-review"})
	if runErr != nil {
		restoreStdout()
		t.Fatalf("pending-review command: %v", runErr)
	}
	if exitCode != 0 {
		restoreStdout()
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	restoreStdout()
	expected := []string{
		"/api/search?status=draft",
		"/api/search",
		"/api/search?problem=true",
	}
	if !slices.Equal(paths, expected) {
		t.Fatalf("unexpected paths: got=%#v want=%#v", paths, expected)
	}
	output := stdout.String()
	if !strings.Contains(output, `"draft_items"`) || !strings.Contains(output, `"unreviewed_runs"`) || !strings.Contains(output, `"problem_runs"`) {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, `"RUN-0001"`) || strings.Contains(output, `"RUN-0009"`) {
		t.Fatalf("unexpected unreviewed run output: %s", output)
	}
}

func TestPendingReviewCommandRejectsMalformedRunApprovals(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.RawQuery {
		case "status=draft":
			_, _ = writer.Write([]byte(`{"items":[]}`))
		case "":
			_, _ = writer.Write([]byte(`{"runs":[{"id":"RUN-0001","approvals":{"actor":"carol"}}]}`))
		case "problem=true":
			_, _ = writer.Write([]byte(`{"runs":[]}`))
		default:
			t.Fatalf("unexpected query: %q", request.URL.RawQuery)
		}
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"pending-review"})
	if exitCode != 1 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if err == nil || !strings.Contains(err.Error(), `/api/search runs entry "approvals" field is not an array`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPendingReviewCommandReportsFailingSharedRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.RawQuery {
		case "status=draft":
			_, _ = writer.Write([]byte(`{"items":[]}`))
		case "":
			http.Error(writer, "search unavailable", http.StatusServiceUnavailable)
		case "problem=true":
			_, _ = writer.Write([]byte(`{"runs":[]}`))
		default:
			t.Fatalf("unexpected query: %q", request.URL.RawQuery)
		}
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"pending-review"})
	if exitCode != 1 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if err == nil || !strings.Contains(err.Error(), "search unavailable") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShowResponsibilityCommandUsesExpectedRoute(t *testing.T) {
	var path string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path = request.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"id":"RESP-0001","title":"Line lead","links":[]}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"show-responsibility", "RESP-0001"})
	if err != nil {
		t.Fatalf("show-responsibility command: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if path != "/api/responsibilities/RESP-0001" {
		t.Fatalf("unexpected path: %s", path)
	}
}

func TestShowPlaceCommandRendersDrilldownSummary(t *testing.T) {
	var path string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path = request.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"id":"PLACE-0001",
			"kind":"area",
			"name":"Receiving",
			"summary":"Inbound inspection area",
			"parent_id":"PLACE-0000",
			"child_place_ids":["PLACE-0002"],
			"resource_ids":["RES-0001"],
			"related_runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"ITEM-0001","outcome":"accepted_with_notes","notes":"Outer wrap torn","resource_ids":["RES-0001"]}],
			"links":[{"relation":"stores","from_type":"place","from_id":"PLACE-0001","to_type":"resource","to_id":"RES-0001","notes":"Receiving area stores the connector bin"}]
		}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	stdout, restoreStdout, err := captureStdout(t)
	if err != nil {
		t.Fatalf("capture stdout: %v", err)
	}
	exitCode, runErr := cli.run([]string{"show-place", "PLACE-0001"})
	if runErr != nil {
		restoreStdout()
		t.Fatalf("show-place command: %v", runErr)
	}
	if exitCode != 0 {
		restoreStdout()
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	restoreStdout()
	if path != "/api/places/PLACE-0001" {
		t.Fatalf("unexpected path: %s", path)
	}
	output := stdout.String()
	if !strings.Contains(output, "# Place PLACE-0001") ||
		!strings.Contains(output, "child places: PLACE-0002") ||
		!strings.Contains(output, "resources: RES-0001") ||
		!strings.Contains(output, "show: oks-cli show-run RUN-0001") ||
		!strings.Contains(output, "stores place PLACE-0001 -> resource RES-0001") {
		t.Fatalf("unexpected show-place output: %s", output)
	}
}

func TestShowResourceCommandRendersDrilldownSummary(t *testing.T) {
	var path string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path = request.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"id":"RES-0001",
			"kind":"container",
			"name":"RJ45 Bin",
			"summary":"Connector bin",
			"place_id":"PLACE-0001",
			"related_runs":[{"id":"RUN-0002","kind":"inventory_audit","item_id":"ITEM-0002","outcome":"completed","notes":"Counted receiving bin","resource_ids":["RES-0001"]}],
			"links":[{"relation":"used_in","from_type":"resource","from_id":"RES-0001","to_type":"run","to_id":"RUN-0002","notes":"Connector bin was counted during the run"}]
		}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	stdout, restoreStdout, err := captureStdout(t)
	if err != nil {
		t.Fatalf("capture stdout: %v", err)
	}
	exitCode, runErr := cli.run([]string{"show-resource", "RES-0001"})
	if runErr != nil {
		restoreStdout()
		t.Fatalf("show-resource command: %v", runErr)
	}
	if exitCode != 0 {
		restoreStdout()
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	restoreStdout()
	if path != "/api/resources/RES-0001" {
		t.Fatalf("unexpected path: %s", path)
	}
	output := stdout.String()
	if !strings.Contains(output, "# Resource RES-0001") ||
		!strings.Contains(output, "place=PLACE-0001") ||
		!strings.Contains(output, "show: oks-cli show-run RUN-0002") ||
		!strings.Contains(output, "used_in resource RES-0001 -> run RUN-0002") {
		t.Fatalf("unexpected show-resource output: %s", output)
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

func TestSearchCommandAddsStructuredFilters(t *testing.T) {
	var values map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		values = map[string]string{
			"q":                 request.URL.Query().Get("q"),
			"kind":              request.URL.Query().Get("kind"),
			"place_id":          request.URL.Query().Get("place_id"),
			"responsibility_id": request.URL.Query().Get("responsibility_id"),
			"problem":           request.URL.Query().Get("problem"),
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"search", "dock discrepancy", "kind=inventory_audit", "place_id=PLACE-0001", "responsibility_id=RESP-0001", "problem=true"})
	if err != nil {
		t.Fatalf("search with filters: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if values["q"] != "dock discrepancy" || values["kind"] != "inventory_audit" || values["place_id"] != "PLACE-0001" || values["responsibility_id"] != "RESP-0001" || values["problem"] != "true" {
		t.Fatalf("unexpected search query values: %#v", values)
	}
}

func TestSearchCommandRejectsUnsupportedFilterKey(t *testing.T) {
	cli := &CLI{ServerURL: "http://127.0.0.1:7045"}
	exitCode, err := cli.run([]string{"search", "dock discrepancy", "owner=alice"})
	if exitCode != 1 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if err == nil || !strings.Contains(err.Error(), `unsupported search filter "owner"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func captureStdout(t *testing.T) (*bytes.Buffer, func(), error) {
	t.Helper()
	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	os.Stdout = writer
	buffer := &bytes.Buffer{}
	return buffer, func() {
		_ = writer.Close()
		_, _ = io.Copy(buffer, reader)
		_ = reader.Close()
		os.Stdout = originalStdout
	}, nil
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

func TestAddLinkCommandEmitsExpectedPayload(t *testing.T) {
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
	exitCode, err := cli.run([]string{"add-link", "alice", "responsibility", "RESP-0001", "knowledge_item", "ITEM-0001", "uses", "startup ownership"})
	if err != nil {
		t.Fatalf("add-link command: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if path != "/api/links" {
		t.Fatalf("unexpected path: %s", path)
	}
	if received["actor"] != "alice" || received["from_type"] != "responsibility" || received["from_id"] != "RESP-0001" || received["to_type"] != "knowledge_item" || received["to_id"] != "ITEM-0001" || received["relation"] != "uses" {
		t.Fatalf("unexpected add-link payload: %#v", received)
	}
	if received["notes"] != "startup ownership" {
		t.Fatalf("unexpected add-link notes: %#v", received["notes"])
	}
}

func TestAddLinkCommandSurfacesServerValidationFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "to endpoint invalid: knowledge_item \"ITEM-9999\" not found", http.StatusBadRequest)
	}))
	defer server.Close()

	cli := &CLI{ServerURL: server.URL}
	exitCode, err := cli.run([]string{"add-link", "alice", "responsibility", "RESP-0001", "knowledge_item", "ITEM-9999", "uses"})
	if exitCode != 1 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if err == nil || !strings.Contains(err.Error(), `to endpoint invalid: knowledge_item "ITEM-9999" not found`) {
		t.Fatalf("unexpected add-link error: %v", err)
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
