package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestRelayServerMeta(t *testing.T) {
	relay, err := NewRelay(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new relay: %v", err)
	}
	defer func() {
		if err := relay.Close(); err != nil {
			t.Fatalf("close relay: %v", err)
		}
	}()

	server := NewRelayServer(relay)
	request := httptest.NewRequest(http.MethodGet, "/relay/v1/meta", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected relay meta status: %d %s", response.Code, response.Body.String())
	}
	var meta RelayMeta
	if err := json.Unmarshal(response.Body.Bytes(), &meta); err != nil {
		t.Fatalf("decode relay meta: %v", err)
	}
	if meta.ServiceName != relayServiceName {
		t.Fatalf("unexpected relay service name: %+v", meta)
	}
	if meta.RoutePrefix != relayRoutePrefix {
		t.Fatalf("unexpected relay route prefix: %+v", meta)
	}
	if meta.RelayFeedFormat != relayFeedFormat || !meta.RelayBlobTransferEnabled || !meta.PublishRequiresStagedBlobs {
		t.Fatalf("unexpected relay meta capabilities: %+v", meta)
	}
}

func TestRelayServerPublishRequiresBlobStaging(t *testing.T) {
	sourceApp, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	item, err := sourceApp.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := sourceApp.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if _, err := sourceApp.AddEvidence("carol", run.ID, "photo", map[string]string{"result": "ok"}, "evidence.txt", []byte("photo")); err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	batch, err := sourceApp.ExportRelayFeed(RelayFeedRequest{KnownOrigins: map[string]uint64{}})
	if err != nil {
		t.Fatalf("export relay feed: %v", err)
	}

	relay, err := NewRelay(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new relay: %v", err)
	}
	defer func() {
		if err := relay.Close(); err != nil {
			t.Fatalf("close relay: %v", err)
		}
	}()
	server := NewRelayServer(relay)

	body, err := json.Marshal(batch)
	if err != nil {
		t.Fatalf("encode relay batch: %v", err)
	}
	publishRequest := httptest.NewRequest(http.MethodPost, "/relay/v1/feed/publish", bytes.NewReader(body))
	publishResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(publishResponse, publishRequest)
	if publishResponse.Code != http.StatusConflict {
		t.Fatalf("expected relay publish conflict, got %d %s", publishResponse.Code, publishResponse.Body.String())
	}
	if !strings.Contains(publishResponse.Body.String(), batch.RequiredBlobCIDs[0]) {
		t.Fatalf("expected missing blob cid in relay publish response: %s", publishResponse.Body.String())
	}

	blobBody, err := sourceApp.RelayBlob(batch.RequiredBlobCIDs[0])
	if err != nil {
		t.Fatalf("load source blob: %v", err)
	}
	putRequest := httptest.NewRequest(http.MethodPut, "/relay/v1/blobs/"+batch.RequiredBlobCIDs[0], bytes.NewReader(blobBody))
	putResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(putResponse, putRequest)
	if putResponse.Code != http.StatusCreated {
		t.Fatalf("unexpected relay blob put status: %d %s", putResponse.Code, putResponse.Body.String())
	}

	publishRequest = httptest.NewRequest(http.MethodPost, "/relay/v1/feed/publish", bytes.NewReader(body))
	publishResponse = httptest.NewRecorder()
	server.Handler().ServeHTTP(publishResponse, publishRequest)
	if publishResponse.Code != http.StatusCreated {
		t.Fatalf("unexpected relay publish status after blob staging: %d %s", publishResponse.Code, publishResponse.Body.String())
	}
	if !strings.Contains(publishResponse.Body.String(), `"published_knowledge_evidence":1`) {
		t.Fatalf("unexpected relay publish response: %s", publishResponse.Body.String())
	}
}

func TestRelayServerPullReturnsOnlyUnseenOrigins(t *testing.T) {
	sourceApp, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	item, err := sourceApp.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := sourceApp.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if _, err := sourceApp.AddEvidence("carol", run.ID, "photo-1", map[string]string{"result": "ok"}, "evidence-1.txt", []byte("photo-1")); err != nil {
		t.Fatalf("add first evidence: %v", err)
	}
	firstBatch, err := sourceApp.ExportRelayFeed(RelayFeedRequest{KnownOrigins: map[string]uint64{}})
	if err != nil {
		t.Fatalf("export first relay feed: %v", err)
	}

	relay, err := NewRelay(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new relay: %v", err)
	}
	defer func() {
		if err := relay.Close(); err != nil {
			t.Fatalf("close relay: %v", err)
		}
	}()
	server := NewRelayServer(relay)

	for _, cid := range firstBatch.RequiredBlobCIDs {
		blobBody, err := sourceApp.RelayBlob(cid)
		if err != nil {
			t.Fatalf("load first blob %s: %v", cid, err)
		}
		putRequest := httptest.NewRequest(http.MethodPut, "/relay/v1/blobs/"+cid, bytes.NewReader(blobBody))
		putResponse := httptest.NewRecorder()
		server.Handler().ServeHTTP(putResponse, putRequest)
		if putResponse.Code != http.StatusCreated {
			t.Fatalf("unexpected relay blob put status for %s: %d %s", cid, putResponse.Code, putResponse.Body.String())
		}
	}
	firstBody, err := json.Marshal(firstBatch)
	if err != nil {
		t.Fatalf("encode first relay batch: %v", err)
	}
	firstPublish := httptest.NewRequest(http.MethodPost, "/relay/v1/feed/publish", bytes.NewReader(firstBody))
	firstResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(firstResponse, firstPublish)
	if firstResponse.Code != http.StatusCreated {
		t.Fatalf("unexpected first relay publish status: %d %s", firstResponse.Code, firstResponse.Body.String())
	}

	if _, err := sourceApp.AddEvidence("carol", run.ID, "photo-2", map[string]string{"result": "still_ok"}, "evidence-2.txt", []byte("photo-2")); err != nil {
		t.Fatalf("add second evidence: %v", err)
	}
	secondBatch, err := sourceApp.ExportRelayFeed(RelayFeedRequest{KnownOrigins: map[string]uint64{sourceApp.localPeerID: uint64(len(firstBatch.Events))}})
	if err != nil {
		t.Fatalf("export second relay feed: %v", err)
	}
	for _, cid := range secondBatch.RequiredBlobCIDs {
		blobBody, err := sourceApp.RelayBlob(cid)
		if err != nil {
			t.Fatalf("load second blob %s: %v", cid, err)
		}
		putRequest := httptest.NewRequest(http.MethodPut, "/relay/v1/blobs/"+cid, bytes.NewReader(blobBody))
		putResponse := httptest.NewRecorder()
		server.Handler().ServeHTTP(putResponse, putRequest)
		if putResponse.Code != http.StatusCreated {
			t.Fatalf("unexpected relay blob put status for %s: %d %s", cid, putResponse.Code, putResponse.Body.String())
		}
	}
	secondBody, err := json.Marshal(secondBatch)
	if err != nil {
		t.Fatalf("encode second relay batch: %v", err)
	}
	secondPublish := httptest.NewRequest(http.MethodPost, "/relay/v1/feed/publish", bytes.NewReader(secondBody))
	secondResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(secondResponse, secondPublish)
	if secondResponse.Code != http.StatusCreated {
		t.Fatalf("unexpected second relay publish status: %d %s", secondResponse.Code, secondResponse.Body.String())
	}

	pullRequest := httptest.NewRequest(http.MethodPost, "/relay/v1/feed/pull", bytes.NewBufferString(`{"known_origins":{"`+sourceApp.localPeerID+`":3}}`))
	pullResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(pullResponse, pullRequest)
	if pullResponse.Code != http.StatusOK {
		t.Fatalf("unexpected relay pull status: %d %s", pullResponse.Code, pullResponse.Body.String())
	}
	var pulled RelayFeedBatch
	if err := json.Unmarshal(pullResponse.Body.Bytes(), &pulled); err != nil {
		t.Fatalf("decode relay pull batch: %v", err)
	}
	if len(pulled.Events) != 1 || pulled.Events[0].Type != "evidence_added" || pulled.Events[0].OriginSequence != 4 {
		t.Fatalf("unexpected pulled relay events: %+v", pulled.Events)
	}
	if pulled.Events[0].Sequence != 1 {
		t.Fatalf("expected relay pull to renumber compatibility sequence, got %+v", pulled.Events[0])
	}
	if len(pulled.RequiredBlobCIDs) != 1 {
		t.Fatalf("expected one required relay blob cid, got %+v", pulled.RequiredBlobCIDs)
	}
}

func TestRelayServerRejectsFirstPublishGap(t *testing.T) {
	sourceApp, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	item, err := sourceApp.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := sourceApp.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if _, err := sourceApp.AddEvidence("carol", run.ID, "photo", map[string]string{"result": "ok"}, "evidence.txt", []byte("photo")); err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	batch, err := sourceApp.ExportRelayFeed(RelayFeedRequest{KnownOrigins: map[string]uint64{sourceApp.localPeerID: 2}})
	if err != nil {
		t.Fatalf("export gapped relay feed: %v", err)
	}

	relay, err := NewRelay(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new relay: %v", err)
	}
	defer func() {
		if err := relay.Close(); err != nil {
			t.Fatalf("close relay: %v", err)
		}
	}()
	server := NewRelayServer(relay)

	for _, cid := range batch.RequiredBlobCIDs {
		blobBody, err := sourceApp.RelayBlob(cid)
		if err != nil {
			t.Fatalf("load blob %s: %v", cid, err)
		}
		putRequest := httptest.NewRequest(http.MethodPut, "/relay/v1/blobs/"+cid, bytes.NewReader(blobBody))
		putResponse := httptest.NewRecorder()
		server.Handler().ServeHTTP(putResponse, putRequest)
		if putResponse.Code != http.StatusCreated {
			t.Fatalf("unexpected relay blob put status: %d %s", putResponse.Code, putResponse.Body.String())
		}
	}

	body, err := json.Marshal(batch)
	if err != nil {
		t.Fatalf("encode gapped batch: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/relay/v1/feed/publish", bytes.NewReader(body))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected relay publish gap rejection, got %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "must be contiguous from 1") {
		t.Fatalf("expected start-at-1 relay rejection, got %s", response.Body.String())
	}
}
