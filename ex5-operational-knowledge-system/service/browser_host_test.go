package service

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"path/filepath"
	"testing"
	"time"
)

func TestBrowserHostForwardsTypedOperations(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	socketPath := EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket"))
	server := NewLocalEmbodimentServer(app, socketPath)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, socketPath)

	requestPayload, err := json.Marshal(BrowserHostEnvelope{
		RequestID:  "req-1",
		SocketPath: socketPath,
		Request: LocalEmbodimentRequest{
			Type:      "operation",
			Operation: "inspect_item",
			ItemID:    item.ID,
		},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	input := bytes.NewBuffer(nil)
	if err := binary.Write(input, binary.LittleEndian, uint32(len(requestPayload))); err != nil {
		t.Fatalf("write request size: %v", err)
	}
	if _, err := input.Write(requestPayload); err != nil {
		t.Fatalf("write request payload: %v", err)
	}

	output := bytes.NewBuffer(nil)
	host := NewBrowserHost()
	if err := host.ServeSession(bytes.NewReader(input.Bytes()), output); err != nil {
		t.Fatalf("serve session: %v", err)
	}
	response := decodeBrowserHostResponse(t, output)
	if response.RequestID != "req-1" {
		t.Fatalf("unexpected request id: %+v", response)
	}
	if response.Response.Status != 200 {
		t.Fatalf("unexpected host response: %+v", response)
	}
	var inspected KnowledgeItem
	if err := json.Unmarshal([]byte(response.Response.Body), &inspected); err != nil {
		t.Fatalf("decode inspected item: %v", err)
	}
	if inspected.ID != item.ID {
		t.Fatalf("unexpected inspected item id: %q", inspected.ID)
	}
}

func TestBrowserHostForwardsRuntimeReadyProbe(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	socketPath := EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket"))
	server := NewLocalEmbodimentServer(app, socketPath)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, socketPath)

	requestPayload, err := json.Marshal(BrowserHostEnvelope{
		RequestID:  "ready-1",
		SocketPath: socketPath,
		Request: LocalEmbodimentRequest{
			Type:      "operation",
			Operation: "runtime_ready",
		},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	input := bytes.NewBuffer(nil)
	if err := binary.Write(input, binary.LittleEndian, uint32(len(requestPayload))); err != nil {
		t.Fatalf("write request size: %v", err)
	}
	if _, err := input.Write(requestPayload); err != nil {
		t.Fatalf("write request payload: %v", err)
	}

	output := bytes.NewBuffer(nil)
	host := NewBrowserHost()
	if err := host.ServeSession(bytes.NewReader(input.Bytes()), output); err != nil {
		t.Fatalf("serve session: %v", err)
	}
	response := decodeBrowserHostResponse(t, output)
	if response.RequestID != "ready-1" {
		t.Fatalf("unexpected request id: %+v", response)
	}
	if response.Response.Status != 200 {
		t.Fatalf("unexpected host response: %+v", response)
	}
	if response.Response.Body != `{"ready":true}` {
		t.Fatalf("unexpected readiness payload: %+v", response)
	}
}

func TestBrowserHostForwardsDashboardAndCollectionReads(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	if _, err := app.CreatePlace("alice", "area", "Receiving", "Inbound", "", nil); err != nil {
		t.Fatalf("create place: %v", err)
	}
	socketPath := EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket"))
	server := NewLocalEmbodimentServer(app, socketPath)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, socketPath)

	for _, operation := range []string{"dashboard", "list_places"} {
		requestPayload, err := json.Marshal(BrowserHostEnvelope{
			RequestID:  operation + "-1",
			SocketPath: socketPath,
			Request: LocalEmbodimentRequest{
				Type:      "operation",
				Operation: operation,
			},
		})
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}
		input := bytes.NewBuffer(nil)
		if err := binary.Write(input, binary.LittleEndian, uint32(len(requestPayload))); err != nil {
			t.Fatalf("write request size: %v", err)
		}
		if _, err := input.Write(requestPayload); err != nil {
			t.Fatalf("write request payload: %v", err)
		}
		output := bytes.NewBuffer(nil)
		host := NewBrowserHost()
		if err := host.ServeSession(bytes.NewReader(input.Bytes()), output); err != nil {
			t.Fatalf("serve session: %v", err)
		}
		response := decodeBrowserHostResponse(t, output)
		if response.Response.Status != 200 {
			t.Fatalf("unexpected host response for %s: %+v", operation, response)
		}
	}
}

func TestBrowserHostForwardsLiveStateBootstrapOperation(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	socketPath := EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket"))
	server := NewLocalEmbodimentServer(app, socketPath)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, socketPath)

	requestPayload, err := json.Marshal(BrowserHostEnvelope{
		RequestID:  "live-state-1",
		SocketPath: socketPath,
		Request: LocalEmbodimentRequest{
			Type:      "operation",
			Operation: "load_live_state",
			ItemID:    item.ID,
		},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	input := bytes.NewBuffer(nil)
	if err := binary.Write(input, binary.LittleEndian, uint32(len(requestPayload))); err != nil {
		t.Fatalf("write request size: %v", err)
	}
	if _, err := input.Write(requestPayload); err != nil {
		t.Fatalf("write request payload: %v", err)
	}

	output := bytes.NewBuffer(nil)
	host := NewBrowserHost()
	if err := host.ServeSession(bytes.NewReader(input.Bytes()), output); err != nil {
		t.Fatalf("serve session: %v", err)
	}
	response := decodeBrowserHostResponse(t, output)
	if response.Response.Status != 200 {
		t.Fatalf("unexpected host response: %+v", response)
	}
	var state LiveItemState
	if err := json.Unmarshal([]byte(response.Response.Body), &state); err != nil {
		t.Fatalf("decode state: %v", err)
	}
	if state.ItemID != item.ID || state.Body != "# Start" {
		t.Fatalf("unexpected live state: %+v", state)
	}
}

func TestBrowserHostForwardsCreateItemOperation(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	socketPath := EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket"))
	server := NewLocalEmbodimentServer(app, socketPath)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, socketPath)

	requestPayload, err := json.Marshal(BrowserHostEnvelope{
		RequestID:  "create-item-1",
		SocketPath: socketPath,
		Request: LocalEmbodimentRequest{
			Type:              "operation",
			Operation:         "create_item",
			Actor:             "alice",
			Kind:              KnowledgeKindProcedure,
			Title:             "Startup checklist",
			Summary:           "Boot line",
			Body:              "# Startup checklist",
			Tags:              []string{"startup"},
			ResponsibilityIDs: nil,
		},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	input := bytes.NewBuffer(nil)
	if err := binary.Write(input, binary.LittleEndian, uint32(len(requestPayload))); err != nil {
		t.Fatalf("write request size: %v", err)
	}
	if _, err := input.Write(requestPayload); err != nil {
		t.Fatalf("write request payload: %v", err)
	}

	output := bytes.NewBuffer(nil)
	host := NewBrowserHost()
	if err := host.ServeSession(bytes.NewReader(input.Bytes()), output); err != nil {
		t.Fatalf("serve session: %v", err)
	}
	response := decodeBrowserHostResponse(t, output)
	if response.Response.Status != 200 {
		t.Fatalf("unexpected host response: %+v", response)
	}
	var item KnowledgeItem
	if err := json.Unmarshal([]byte(response.Response.Body), &item); err != nil {
		t.Fatalf("decode created item: %v", err)
	}
	if item.Title != "Startup checklist" || item.Kind != KnowledgeKindProcedure {
		t.Fatalf("unexpected created item: %+v", item)
	}
}

func TestBrowserHostForwardsAddEvidenceOperation(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindProcedure, item.ID, 1, "completed", "done", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	socketPath := EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket"))
	server := NewLocalEmbodimentServer(app, socketPath)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, socketPath)

	requestPayload, err := json.Marshal(BrowserHostEnvelope{
		RequestID:  "evidence-1",
		SocketPath: socketPath,
		Request: LocalEmbodimentRequest{
			Type:                 "operation",
			Operation:            "add_evidence",
			Actor:                "bob",
			RunID:                run.ID,
			Summary:              "Checklist photo",
			Facts:                map[string]string{"result": "ok"},
			AttachmentName:       "check.txt",
			AttachmentBodyBase64: base64.StdEncoding.EncodeToString([]byte("ok")),
		},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	input := bytes.NewBuffer(nil)
	if err := binary.Write(input, binary.LittleEndian, uint32(len(requestPayload))); err != nil {
		t.Fatalf("write request size: %v", err)
	}
	if _, err := input.Write(requestPayload); err != nil {
		t.Fatalf("write request payload: %v", err)
	}

	output := bytes.NewBuffer(nil)
	host := NewBrowserHost()
	if err := host.ServeSession(bytes.NewReader(input.Bytes()), output); err != nil {
		t.Fatalf("serve session: %v", err)
	}
	response := decodeBrowserHostResponse(t, output)
	if response.Response.Status != 200 {
		t.Fatalf("unexpected host response: %+v", response)
	}
	var updatedRun RunRecord
	if err := json.Unmarshal([]byte(response.Response.Body), &updatedRun); err != nil {
		t.Fatalf("decode run: %v", err)
	}
	if len(updatedRun.Evidence) != 1 || updatedRun.Evidence[0].AttachmentName != "check.txt" {
		t.Fatalf("unexpected evidence result: %+v", updatedRun)
	}
}

func TestBrowserHostStreamsLiveMessages(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	socketPath := EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket"))
	server := NewLocalEmbodimentServer(app, socketPath)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, socketPath)

	reader, writer := io.Pipe()
	outputReader, outputWriter := io.Pipe()
	host := NewBrowserHost()
	done := make(chan error, 1)
	go func() {
		defer func() {
			_ = outputWriter.Close()
		}()
		done <- host.ServeSession(reader, outputWriter)
	}()

	encodeBrowserHostEnvelope(t, writer, BrowserHostEnvelope{
		RequestID:  "live-1",
		SocketPath: socketPath,
		Request: LocalEmbodimentRequest{
			Type:          "live-open",
			ItemID:        item.ID,
			ParticipantID: "browser-a",
			DisplayName:   "Alice",
			Color:         "#123456",
		},
	})
	initial := decodeBrowserHostResponseFromReader(t, outputReader)
	if initial.Response.Type != "live-state" || initial.Response.State.Version != 1 {
		t.Fatalf("unexpected initial live response: %+v", initial)
	}

	encodeBrowserHostEnvelope(t, writer, BrowserHostEnvelope{
		RequestID:  "live-1",
		SocketPath: socketPath,
		Request: LocalEmbodimentRequest{
			Type:          "live-update",
			ParticipantID: "browser-a",
			DisplayName:   "Alice",
			Color:         "#123456",
			Cursor:        5,
			Head:          5,
			Typing:        true,
			BaseVersion:   1,
			UpdateBody:    true,
			Body:          "# Start\n\nChecked PPE.",
		},
	})
	updated := decodeBrowserHostResponseFromReader(t, outputReader)
	if updated.Response.Type != "live-state" || updated.Response.State.Version != 2 {
		t.Fatalf("unexpected live update response: %+v", updated)
	}
	if updated.Response.State.Body != "# Start\n\nChecked PPE." {
		t.Fatalf("unexpected live body: %+v", updated.Response.State)
	}

	encodeBrowserHostEnvelope(t, writer, BrowserHostEnvelope{
		RequestID: "live-1",
		Request: LocalEmbodimentRequest{
			Type: "live-close",
		},
	})
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("serve live session: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("live browser host session did not exit")
	}
}

func encodeBrowserHostEnvelope(t *testing.T, writer io.Writer, envelope BrowserHostEnvelope) {
	t.Helper()
	payload, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}
	if err := binary.Write(writer, binary.LittleEndian, uint32(len(payload))); err != nil {
		t.Fatalf("write envelope size: %v", err)
	}
	if _, err := writer.Write(payload); err != nil {
		t.Fatalf("write envelope payload: %v", err)
	}
}

func decodeBrowserHostResponse(t *testing.T, buffer *bytes.Buffer) BrowserHostResponse {
	t.Helper()
	var size uint32
	if err := binary.Read(buffer, binary.LittleEndian, &size); err != nil {
		t.Fatalf("read response size: %v", err)
	}
	payload := make([]byte, size)
	if _, err := io.ReadFull(buffer, payload); err != nil {
		t.Fatalf("read response payload: %v", err)
	}
	var response BrowserHostResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return response
}

func decodeBrowserHostResponseFromReader(t *testing.T, reader io.Reader) BrowserHostResponse {
	t.Helper()
	var size uint32
	if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
		t.Fatalf("read response size: %v", err)
	}
	payload := make([]byte, size)
	if _, err := io.ReadFull(reader, payload); err != nil {
		t.Fatalf("read response payload: %v", err)
	}
	var response BrowserHostResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return response
}
