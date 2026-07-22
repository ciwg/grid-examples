package service

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestServerCreatesAndListsKnowledgeItems(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewServer(app)

	body := bytes.NewBufferString(`{"actor":"alice","kind":"procedure","title":"Start line","summary":"startup","body":"# Start","tags":["ops"],"responsibility_ids":[]}`)
	request := httptest.NewRequest(http.MethodPost, "/api/items", body)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("unexpected status: %d %s", response.Code, response.Body.String())
	}

	response = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/api/items?kind=procedure", nil)
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected list status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"kind":"procedure"`) {
		t.Fatalf("missing procedure list body: %s", response.Body.String())
	}
}

func TestServerMetaIncludesRuntimeCapabilities(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/meta", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected meta status: %d %s", response.Code, response.Body.String())
	}
	var meta Meta
	if err := json.Unmarshal(response.Body.Bytes(), &meta); err != nil {
		t.Fatalf("decode meta: %v", err)
	}
	if meta.PeerExchangeFormat != peerExchangeBundleFormat {
		t.Fatalf("unexpected peer exchange format in meta: %+v", meta)
	}
	if meta.OperationalRunPCID == "" {
		t.Fatalf("expected operational-run pCID in meta: %+v", meta)
	}
	if !meta.CASObjectsEnabled || !meta.CASAttachmentBlobsEnabled {
		t.Fatalf("expected CAS capability flags in meta: %+v", meta)
	}
}

func TestServerUploadsEvidenceAttachment(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindProcedure, item.ID, 1, "completed", "", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}

	server := NewServer(app)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("actor", "bob"); err != nil {
		t.Fatalf("actor field: %v", err)
	}
	if err := writer.WriteField("summary", "photo"); err != nil {
		t.Fatalf("summary field: %v", err)
	}
	if err := writer.WriteField("facts_json", `{"result":"ok"}`); err != nil {
		t.Fatalf("facts field: %v", err)
	}
	part, err := writer.CreateFormFile("attachment", "evidence.txt")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write([]byte("hello")); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/runs/"+run.ID+"/evidence", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d %s", response.Code, response.Body.String())
	}

	var decoded RunRecord
	if err := json.Unmarshal(response.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(decoded.Evidence) != 1 || decoded.Evidence[0].AttachmentName != "evidence.txt" {
		t.Fatalf("unexpected evidence payload: %+v", decoded)
	}
	if decoded.Evidence[0].AttachmentCID == "" {
		t.Fatalf("expected attachment CID in evidence payload: %+v", decoded.Evidence[0])
	}
}

func TestServerRejectsOversizedEvidenceAttachment(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindProcedure, item.ID, 1, "completed", "", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}

	server := NewServer(app)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("actor", "bob"); err != nil {
		t.Fatalf("actor field: %v", err)
	}
	if err := writer.WriteField("summary", "oversized photo"); err != nil {
		t.Fatalf("summary field: %v", err)
	}
	part, err := writer.CreateFormFile("attachment", "too-large.bin")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(bytes.Repeat([]byte("a"), maxEvidenceAttachmentBytes+1)); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/runs/"+run.ID+"/evidence", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "attachment exceeds") {
		t.Fatalf("unexpected error body: %s", response.Body.String())
	}

	runAfter, err := app.GetRun(run.ID)
	if err != nil {
		t.Fatalf("get run: %v", err)
	}
	if len(runAfter.Evidence) != 0 {
		t.Fatalf("oversized attachment should not create evidence: %+v", runAfter.Evidence)
	}
}

func TestServerWorkflowSearchAndDashboard(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewServer(app)

	responsibilityBody := bytes.NewBufferString(`{"actor":"alice","title":"Training lead","summary":"Owns onboarding records","role_keys":["trainer"],"tags":["training"]}`)
	request := httptest.NewRequest(http.MethodPost, "/api/responsibilities", responsibilityBody)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create responsibility status: %d %s", response.Code, response.Body.String())
	}
	var responsibility Responsibility
	if err := json.Unmarshal(response.Body.Bytes(), &responsibility); err != nil {
		t.Fatalf("decode responsibility: %v", err)
	}

	itemBody := bytes.NewBufferString(`{"actor":"alice","kind":"training","title":"Forklift onboarding","summary":"Training record for new operators","body":"# Forklift onboarding","responsibility_ids":["` + responsibility.ID + `"]}`)
	request = httptest.NewRequest(http.MethodPost, "/api/items", itemBody)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create item status: %d %s", response.Code, response.Body.String())
	}
	var item KnowledgeItem
	if err := json.Unmarshal(response.Body.Bytes(), &item); err != nil {
		t.Fatalf("decode item: %v", err)
	}

	runBody := bytes.NewBufferString(`{"actor":"bob","kind":"training","item_id":"` + item.ID + `","revision":1,"outcome":"passed","notes":"Forklift onboarding completed on first shift","responsibility_ids":["` + responsibility.ID + `"]}`)
	request = httptest.NewRequest(http.MethodPost, "/api/runs", runBody)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create run status: %d %s", response.Code, response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/search?q=onboarding", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("search status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), responsibility.ID) {
		t.Fatalf("search missing responsibility: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), item.ID) {
		t.Fatalf("search missing item: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"kind":"training"`) {
		t.Fatalf("search missing training run/item: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("dashboard status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"responsibilities":1`) || !strings.Contains(response.Body.String(), `"training_items":1`) || !strings.Contains(response.Body.String(), `"training_runs":1`) {
		t.Fatalf("unexpected dashboard body: %s", response.Body.String())
	}
}

func TestServerPlacesResourcesAndLiveItems(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewServer(app)

	placeBody := bytes.NewBufferString(`{"actor":"alice","kind":"area","name":"Receiving","summary":"Inbound inspection area","tags":["inventory"]}`)
	request := httptest.NewRequest(http.MethodPost, "/api/places", placeBody)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create place status: %d %s", response.Code, response.Body.String())
	}
	var place Place
	if err := json.Unmarshal(response.Body.Bytes(), &place); err != nil {
		t.Fatalf("decode place: %v", err)
	}

	resourceBody := bytes.NewBufferString(`{"actor":"alice","kind":"container","name":"RJ45 Bin","summary":"Bin for connectors","place_id":"` + place.ID + `","tags":["inventory","parts"]}`)
	request = httptest.NewRequest(http.MethodPost, "/api/resources", resourceBody)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create resource status: %d %s", response.Code, response.Body.String())
	}
	var resource Resource
	if err := json.Unmarshal(response.Body.Bytes(), &resource); err != nil {
		t.Fatalf("decode resource: %v", err)
	}

	itemBody := bytes.NewBufferString(`{"actor":"alice","kind":"inventory_audit","title":"Count receiving bin","summary":"Cycle count steps","body":"# Count receiving bin"}`)
	request = httptest.NewRequest(http.MethodPost, "/api/items", itemBody)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create item status: %d %s", response.Code, response.Body.String())
	}
	var item KnowledgeItem
	if err := json.Unmarshal(response.Body.Bytes(), &item); err != nil {
		t.Fatalf("decode item: %v", err)
	}

	liveBody := bytes.NewBufferString(`{"participant_id":"browser-a","display_name":"Alice","color":"#1d6fd6","cursor":5,"head":5,"typing":true,"base_version":1,"update_body":true,"body":"# Count receiving bin\n\nObserved 12 connectors"}`)
	request = httptest.NewRequest(http.MethodPost, "/api/items/"+item.ID+"/live", liveBody)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("live update status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"version":2`) || !strings.Contains(response.Body.String(), `"participant_id":"browser-a"`) {
		t.Fatalf("unexpected live update body: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/api/items/"+item.ID+"/supersede", bytes.NewBufferString(`{"actor":"boss","notes":"Replaced by audited count procedure"}`))
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("supersede item status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"status":"superseded"`) {
		t.Fatalf("unexpected supersede body: %s", response.Body.String())
	}

	runBody := bytes.NewBufferString(`{"actor":"bob","kind":"inventory_audit","item_id":"` + item.ID + `","revision":1,"outcome":"completed","notes":"Counted receiving bin","place_id":"` + place.ID + `","resource_ids":["` + resource.ID + `"]}`)
	request = httptest.NewRequest(http.MethodPost, "/api/runs", runBody)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create run status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"place_id":"`+place.ID+`"`) || !strings.Contains(response.Body.String(), resource.ID) {
		t.Fatalf("unexpected run body: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/search?q=receiving", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("search status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), place.ID) || !strings.Contains(response.Body.String(), resource.ID) {
		t.Fatalf("search missing place/resource: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("dashboard status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"inventory_items":1`) {
		t.Fatalf("dashboard missing inventory count: %s", response.Body.String())
	}
}

func TestServerExportsAndImportsPeerExchangeBundle(t *testing.T) {
	sourceApp, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	responsibility, err := sourceApp.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := sourceApp.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := sourceApp.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if err := sourceApp.RecordApproval("carol", "run", run.ID, 0, "auditor", DecisionNoted, "captured"); err != nil {
		t.Fatalf("approve run: %v", err)
	}
	if err := sourceApp.AddLink("alice", "responsibility", responsibility.ID, "knowledge_item", item.ID, "owns", "Receiving lead owns intake"); err != nil {
		t.Fatalf("add link: %v", err)
	}
	sourceServer := NewServer(sourceApp)

	exportRequest := httptest.NewRequest(http.MethodGet, "/api/peer-exchange/export", nil)
	exportResponse := httptest.NewRecorder()
	sourceServer.Handler().ServeHTTP(exportResponse, exportRequest)
	if exportResponse.Code != http.StatusOK {
		t.Fatalf("unexpected export status: %d %s", exportResponse.Code, exportResponse.Body.String())
	}
	var bundle PeerExchangeBundle
	if err := json.Unmarshal(exportResponse.Body.Bytes(), &bundle); err != nil {
		t.Fatalf("decode export bundle: %v", err)
	}
	if bundle.Format != peerExchangeBundleFormat {
		t.Fatalf("unexpected bundle format %q", bundle.Format)
	}

	targetApp, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	targetServer := NewServer(targetApp)
	body, err := json.Marshal(bundle)
	if err != nil {
		t.Fatalf("encode import bundle: %v", err)
	}
	importRequest := httptest.NewRequest(http.MethodPost, "/api/peer-exchange/import", bytes.NewReader(body))
	importResponse := httptest.NewRecorder()
	targetServer.Handler().ServeHTTP(importResponse, importRequest)
	if importResponse.Code != http.StatusOK {
		t.Fatalf("unexpected import status: %d %s", importResponse.Code, importResponse.Body.String())
	}
	var result PeerExchangeImportResult
	if err := json.Unmarshal(importResponse.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode import result: %v", err)
	}
	if result.ImportedKnowledgeItems != len(bundle.KnowledgeItemRecords) {
		t.Fatalf("unexpected import item count: %+v", result)
	}
	if result.ImportedOperationalRuns != 1 {
		t.Fatalf("expected imported operational run count after bootstrap import, got %+v", result)
	}
	if len(result.UnresolvedReferences) != 0 {
		t.Fatalf("expected no unresolved references after run-family bootstrap import, got %+v", result.UnresolvedReferences)
	}
}

func TestServerRejectsPeerExchangeImportIntoNonEmptyRuntime(t *testing.T) {
	sourceApp, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	responsibility, err := sourceApp.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	if _, err := sourceApp.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{responsibility.ID}); err != nil {
		t.Fatalf("create item: %v", err)
	}
	bundle, err := sourceApp.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}

	targetApp, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	if _, err := targetApp.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "Startup flow", "# Start line", nil, nil); err != nil {
		t.Fatalf("seed target item: %v", err)
	}
	targetServer := NewServer(targetApp)
	body, err := json.Marshal(bundle)
	if err != nil {
		t.Fatalf("encode import bundle: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/peer-exchange/import", bytes.NewReader(body))
	response := httptest.NewRecorder()
	targetServer.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected non-empty import success for disjoint ids, got %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"imported_knowledge_items":1`) {
		t.Fatalf("unexpected import body: %s", response.Body.String())
	}
}

func TestServerAllowsPeerExchangeImportEntityAliasReuse(t *testing.T) {
	sourceApp, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	sourceItem, err := sourceApp.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create source item: %v", err)
	}
	bundle, err := sourceApp.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}

	targetApp, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	if _, err := targetApp.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect receiving", "Local receiving flow", "# Inspect receiving", nil, nil); err != nil {
		t.Fatalf("seed target receiving item: %v", err)
	}
	targetServer := NewServer(targetApp)
	body, err := json.Marshal(bundle)
	if err != nil {
		t.Fatalf("encode import bundle: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/peer-exchange/import", bytes.NewReader(body))
	response := httptest.NewRecorder()
	targetServer.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected alias reuse import success, got %d %s", response.Code, response.Body.String())
	}
	imported, err := targetApp.GetKnowledgeItem(sourceItem.ID)
	if err != nil {
		t.Fatalf("get imported canonical item: %v", err)
	}
	if imported.AliasID == "" {
		t.Fatalf("expected imported canonical item alias after server import, got %+v", imported)
	}
}

func TestServerLiveItemGetAndConflictResponse(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Inspect station", "Station checklist", "# Inspect station", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/items/"+item.ID+"/live", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("get live state status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"version":1`) || !strings.Contains(response.Body.String(), `"body":"# Inspect station"`) {
		t.Fatalf("unexpected initial live state: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/api/items/"+item.ID+"/live", bytes.NewBufferString(`{"participant_id":"browser-a","display_name":"Alice","color":"#0c6d62","cursor":4,"head":4,"typing":true,"base_version":1,"update_body":true,"body":"# Inspect station\n\nChecked bins."}`))
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("first live post status: %d %s", response.Code, response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/api/items/"+item.ID+"/live", bytes.NewBufferString(`{"participant_id":"browser-b","display_name":"Bob","color":"#b75c1c","cursor":2,"head":2,"typing":true,"base_version":1,"update_body":true,"body":"# Inspect station\n\nStale body."}`))
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("expected conflict status, got %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"conflict":true`) || !strings.Contains(response.Body.String(), `Checked bins.`) {
		t.Fatalf("unexpected conflict payload: %s", response.Body.String())
	}
}

func TestServerLiveItemPresenceOnlyUpdateKeepsBodyAndVersion(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Inspect station", "Station checklist", "# Inspect station", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodPost, "/api/items/"+item.ID+"/live", bytes.NewBufferString(`{"participant_id":"browser-a","display_name":"Alice","color":"#0c6d62","cursor":4,"head":4,"typing":true,"base_version":1,"update_body":true,"body":"# Inspect station\n\nChecked bins."}`))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("initial live body post status: %d %s", response.Code, response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/api/items/"+item.ID+"/live", bytes.NewBufferString(`{"participant_id":"oks-nvim-host-1234","display_name":"Neovim","color":"#d66f1d","cursor":6,"head":6,"typing":false,"base_version":2,"update_body":false,"body":""}`))
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("presence-only live post status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"version":2`) {
		t.Fatalf("presence-only post should not advance version: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `Checked bins.`) {
		t.Fatalf("presence-only post should keep shared body: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"participant_id":"oks-nvim-host-1234"`) || !strings.Contains(response.Body.String(), `"display_name":"Neovim"`) {
		t.Fatalf("presence-only post missing nvim participant: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/items/"+item.ID+"/live", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("get live state after presence-only post status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"version":2`) || !strings.Contains(response.Body.String(), `Checked bins.`) || !strings.Contains(response.Body.String(), `"participant_id":"oks-nvim-host-1234"`) {
		t.Fatalf("unexpected live state after presence-only post: %s", response.Body.String())
	}
}

func TestServerLiveItemClearBodyPersistsEmptyDraft(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Inspect station", "Station checklist", "# Inspect station", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodPost, "/api/items/"+item.ID+"/live", bytes.NewBufferString(`{"participant_id":"browser-a","display_name":"Alice","color":"#0c6d62","cursor":4,"head":4,"typing":true,"base_version":1,"update_body":true,"body":"# Inspect station\n\nChecked bins."}`))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("initial live body post status: %d %s", response.Code, response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/api/items/"+item.ID+"/live", bytes.NewBufferString(`{"participant_id":"browser-a","display_name":"Alice","color":"#0c6d62","cursor":0,"head":0,"typing":false,"base_version":2,"update_body":true,"body":""}`))
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("clear-body post status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"version":3`) || !strings.Contains(response.Body.String(), `"body":""`) {
		t.Fatalf("unexpected clear-body payload: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/items/"+item.ID+"/live", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("get live state after clear-body status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"version":3`) || !strings.Contains(response.Body.String(), `"body":""`) {
		t.Fatalf("unexpected live state after clear-body post: %s", response.Body.String())
	}
}

func TestServerRejectsStaleKnowledgeItemApproval(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Inspect station", "Station checklist", "# Inspect station", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	item, err = app.AddRevision("alice", item.ID, "Inspect station", "Station checklist revised", "# Inspect station\n\nRevision 2", nil)
	if err != nil {
		t.Fatalf("add revision: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodPost, "/api/items/"+item.ID+"/approvals", bytes.NewBufferString(`{"actor":"boss","revision":1,"role":"reviewer","decision":"approved","notes":"stale approval"}`))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected stale approval rejection, got %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "stale revision") {
		t.Fatalf("unexpected stale approval response: %s", response.Body.String())
	}

	loaded, err := app.GetKnowledgeItem(item.ID)
	if err != nil {
		t.Fatalf("get item: %v", err)
	}
	if loaded.Status != ItemStatusDraft {
		t.Fatalf("stale approval should not change item status: %+v", loaded)
	}
	if len(loaded.Approvals) != 0 {
		t.Fatalf("stale approval should not append approval record: %+v", loaded.Approvals)
	}
}

func TestServerSearchAcceptsStructuredFilters(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewServer(app)

	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns receiving checks", []string{"reviewer"}, []string{"receiving"})
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", []string{"receiving"})
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, []string{"parts"})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count flow", "# Count receiving bin", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.RecordApproval("boss", "knowledge_item", item.ID, 1, "reviewer", DecisionApproved, "Ready to use"); err != nil {
		t.Fatalf("record approval: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindInventory, item.ID, 1, "completed", "Counted receiving bin", "", "", place.ID, []string{resource.ID}, []string{resp.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/search?kind=inventory_audit&status=approved&responsibility_id="+resp.ID, nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("search status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"status":"approved"`) || !strings.Contains(response.Body.String(), item.ID) {
		t.Fatalf("filtered search missing approved item: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"responsibility_id":"`+resp.ID+`"`) {
		t.Fatalf("filtered search missing responsibility filter echo: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/search?place_id="+place.ID+"&resource_id="+resource.ID, nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("place/resource search status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), place.ID) || !strings.Contains(response.Body.String(), resource.ID) || !strings.Contains(response.Body.String(), run.ID) {
		t.Fatalf("place/resource filtered search missing expected records: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/search?kind=inventory_audit&outcome=completed&place_id="+place.ID, nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("outcome search status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"outcome":"completed"`) || !strings.Contains(response.Body.String(), `"kind":"inventory_audit"`) || !strings.Contains(response.Body.String(), run.ID) {
		t.Fatalf("outcome-filtered search missing expected run: %s", response.Body.String())
	}
}

func TestServerSearchIncludesRunEvidenceAndApprovalHistory(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = app.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{
		"supplier":     "Acme Parts",
		"packing_slip": "PS-1234",
		"condition":    "wrap torn",
	}, "", nil)
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if err := app.RecordApproval("boss", "run", run.ID, 0, "reviewer", DecisionApproved, "Reviewed at dock"); err != nil {
		t.Fatalf("record approval: %v", err)
	}

	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/search?q=Acme%20Parts", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("evidence search status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), run.ID) || !strings.Contains(response.Body.String(), `"supplier":"Acme Parts"`) {
		t.Fatalf("evidence search missing run or evidence facts: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/search?q=Reviewed%20at%20dock", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("approval search status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), run.ID) || !strings.Contains(response.Body.String(), `"notes":"Reviewed at dock"`) {
		t.Fatalf("approval search missing run or approval notes: %s", response.Body.String())
	}
}

func TestServerItemDetailIncludesRelatedRuns(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count flow", "# Count receiving bin", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindInventory, item.ID, 1, "completed", "Counted receiving bin", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/items/"+item.ID, nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("item detail status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"related_runs"`) || !strings.Contains(response.Body.String(), run.ID) {
		t.Fatalf("item detail missing related run history: %s", response.Body.String())
	}
}

func TestServerItemDetailIncludesReviewDataForNeovimInspector(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check for pallet", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.RecordApproval("boss", "knowledge_item", item.ID, 1, "reviewer", DecisionApproved, "Ready for receiving"); err != nil {
		t.Fatalf("record approval: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, 1, "accepted_with_notes", "Wrap was torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if _, err := app.AddEvidence("bob", run.ID, "receiving note", map[string]string{"supplier": "Acme Parts", "condition": "wrap torn"}, "", nil); err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/items/"+item.ID, nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("item detail status: %d %s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	if !strings.Contains(body, `"revisions"`) || !strings.Contains(body, `"approvals"`) || !strings.Contains(body, `"related_runs"`) {
		t.Fatalf("item detail missing review sections: %s", body)
	}
	if !strings.Contains(body, `"decision":"approved"`) || !strings.Contains(body, `"notes":"Wrap was torn"`) || !strings.Contains(body, `"condition":"wrap torn"`) {
		t.Fatalf("item detail missing expected inspector review data: %s", body)
	}
}

func TestServerContextDetailIncludesRelatedRuns(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns receiving checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count flow", "# Count receiving bin", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindInventory, item.ID, 1, "completed", "Counted receiving bin", "", "", place.ID, []string{resource.ID}, []string{resp.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	server := NewServer(app)

	for _, path := range []string{
		"/api/places/" + place.ID,
		"/api/resources/" + resource.ID,
		"/api/responsibilities/" + resp.ID,
	} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("detail status for %s: %d %s", path, response.Code, response.Body.String())
		}
		if !strings.Contains(response.Body.String(), `"related_runs"`) || !strings.Contains(response.Body.String(), run.ID) {
			t.Fatalf("context detail missing related runs for %s: %s", path, response.Body.String())
		}
	}
}

func TestServerContextDetailIncludesLinksForNeovimEntityInspector(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns receiving checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Dock receipt", "# Inspect inbound pallet", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", place.ID, []string{resource.ID}, []string{resp.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if err := app.AddLink("alice", "place", place.ID, "resource", resource.ID, "stores", "Receiving area stores the connector bin"); err != nil {
		t.Fatalf("add place link: %v", err)
	}
	if err := app.AddLink("alice", "responsibility", resp.ID, "knowledge_item", item.ID, "owns", "Receiving lead owns the dock receipt item"); err != nil {
		t.Fatalf("add responsibility link: %v", err)
	}
	if err := app.AddLink("alice", "resource", resource.ID, "run", run.ID, "used_in", "Connector bin was counted during the run"); err != nil {
		t.Fatalf("add resource link: %v", err)
	}
	server := NewServer(app)

	for _, path := range []string{
		"/api/places/" + place.ID,
		"/api/resources/" + resource.ID,
		"/api/responsibilities/" + resp.ID,
	} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("detail status for %s: %d %s", path, response.Code, response.Body.String())
		}
		body := response.Body.String()
		if !strings.Contains(body, `"links"`) {
			t.Fatalf("context detail missing links for %s: %s", path, body)
		}
		if path == "/api/responsibilities/"+resp.ID && (!strings.Contains(body, `"linked_item_ids"`) || !strings.Contains(body, `"relation":"owns"`)) {
			t.Fatalf("responsibility detail missing linked item ids for %s: %s", path, body)
		}
	}
}

func TestServerRejectsInvalidTypedLinkEndpoints(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns receiving checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Dock receipt", "# Inspect inbound pallet", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(`{"actor":"alice","from_type":"responsibility","from_id":"RESP-9999","to_type":"knowledge_item","to_id":"`+item.ID+`","relation":"owns","notes":"bad endpoint"}`))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected dangling endpoint rejection, got %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "from endpoint invalid") {
		t.Fatalf("unexpected dangling endpoint response: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(`{"actor":"alice","from_type":"item","from_id":"`+item.ID+`","to_type":"responsibility","to_id":"`+resp.ID+`","relation":"owns","notes":"bad type"}`))
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected unsupported type rejection, got %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "unsupported link endpoint type") {
		t.Fatalf("unexpected unsupported type response: %s", response.Body.String())
	}
}

func TestServerRunDetailIncludesEvidenceFacts(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count flow", "# Count receiving bin", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindInventory, item.ID, 1, "completed", "Counted receiving bin", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = app.AddEvidence("bob", run.ID, "Count sheet", map[string]string{
		"expected_count": "12",
		"actual_count":   "10",
		"discrepancy":    "-2",
	}, "", nil)
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/runs/"+run.ID, nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("run detail status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"expected_count":"12"`) || !strings.Contains(response.Body.String(), `"actual_count":"10"`) || !strings.Contains(response.Body.String(), `"discrepancy":"-2"`) {
		t.Fatalf("run detail missing evidence facts: %s", response.Body.String())
	}
}

func TestServerRunDetailIncludesApprovalsForNeovimInspector(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Dock receipt", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if _, err := app.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{
		"supplier":       "Acme Parts",
		"condition":      "wrap torn",
		"expected_units": "20",
		"received_units": "18",
	}, "", nil); err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if err := app.RecordApproval("lead", "run", run.ID, 1, "reviewer", DecisionNoted, "Investigate damaged wrap"); err != nil {
		t.Fatalf("record run approval: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/runs/"+run.ID, nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("run detail status: %d %s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	if !strings.Contains(body, `"evidence"`) || !strings.Contains(body, `"approvals"`) {
		t.Fatalf("run detail missing review sections: %s", body)
	}
	if !strings.Contains(body, `"condition":"wrap torn"`) || !strings.Contains(body, `"decision":"noted"`) || !strings.Contains(body, `"notes":"Investigate damaged wrap"`) {
		t.Fatalf("run detail missing expected inspector data: %s", body)
	}
}

func TestServerSupportsReceivingCheckKinds(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewServer(app)

	itemBody := bytes.NewBufferString(`{"actor":"alice","kind":"receiving_check","title":"Inspect inbound pallet","summary":"Dock receipt","body":"# Inspect inbound pallet"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/items", itemBody)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create receiving item status: %d %s", response.Code, response.Body.String())
	}
	var item KnowledgeItem
	if err := json.Unmarshal(response.Body.Bytes(), &item); err != nil {
		t.Fatalf("decode receiving item: %v", err)
	}
	if item.Kind != KnowledgeKindReceiving {
		t.Fatalf("unexpected receiving item: %+v", item)
	}

	runBody := bytes.NewBufferString(`{"actor":"bob","kind":"receiving_check","item_id":"` + item.ID + `","revision":1,"outcome":"accepted_with_notes","notes":"Outer wrap torn"}`)
	request = httptest.NewRequest(http.MethodPost, "/api/runs", runBody)
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create receiving run status: %d %s", response.Code, response.Body.String())
	}
	var run RunRecord
	if err := json.Unmarshal(response.Body.Bytes(), &run); err != nil {
		t.Fatalf("decode receiving run: %v", err)
	}
	if run.Kind != RunKindReceiving {
		t.Fatalf("unexpected receiving run: %+v", run)
	}

	request = httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("dashboard status: %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"receiving_items":1`) || !strings.Contains(response.Body.String(), `"receiving_runs":1`) {
		t.Fatalf("dashboard missing receiving counts: %s", response.Body.String())
	}
}

func TestServerReceivingRunDetailIncludesEvidenceFacts(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Dock receipt", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = app.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{
		"supplier":       "Acme Parts",
		"packing_slip":   "PS-1234",
		"received_units": "18",
		"expected_units": "20",
		"variance":       "-2",
	}, "", nil)
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/runs/"+run.ID, nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("run detail status: %d %s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	if !strings.Contains(body, `"supplier":"Acme Parts"`) || !strings.Contains(body, `"packing_slip":"PS-1234"`) || !strings.Contains(body, `"variance":"-2"`) {
		t.Fatalf("receiving run detail missing evidence facts: %s", body)
	}
}

func TestServerProblemReviewGroupsHotspots(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	recvItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create receiving item: %v", err)
	}
	recvRun, err := app.RecordRun("bob", RunKindReceiving, recvItem.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record receiving run: %v", err)
	}
	if _, err := app.AddEvidence("bob", recvRun.ID, "Receiving inspection", map[string]string{"variance": "-2", "condition": "wrap torn"}, "", nil); err != nil {
		t.Fatalf("add receiving evidence: %v", err)
	}
	invItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count", "# Count receiving bin", nil, nil)
	if err != nil {
		t.Fatalf("create inventory item: %v", err)
	}
	invRun, err := app.RecordRun("bob", RunKindInventory, invItem.ID, 1, "completed", "Counted receiving bin", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record inventory run: %v", err)
	}
	if _, err := app.AddEvidence("bob", invRun.ID, "Cycle count", map[string]string{"expected_count": "12", "actual_count": "10", "discrepancy": "-2"}, "", nil); err != nil {
		t.Fatalf("add inventory evidence: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/problem-review", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("problem review status: %d %s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	if !strings.Contains(body, `"problem_runs":2`) || !strings.Contains(body, place.ID) || !strings.Contains(body, resource.ID) {
		t.Fatalf("problem review missing grouped hotspots: %s", body)
	}
	if !strings.Contains(body, `"outcome: accepted_with_notes"`) || !strings.Contains(body, `"discrepancy: -2"`) {
		t.Fatalf("problem review missing expected highlights: %s", body)
	}
}

func TestServerSearchProblemFilterMatchesGroupedReviewLogic(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	recvItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create receiving item: %v", err)
	}
	recvRun, err := app.RecordRun("bob", RunKindReceiving, recvItem.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record receiving run: %v", err)
	}
	if _, err := app.AddEvidence("bob", recvRun.ID, "Receiving inspection", map[string]string{"variance": "-2", "condition": "wrap torn"}, "", nil); err != nil {
		t.Fatalf("add receiving evidence: %v", err)
	}
	invItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count", "# Count receiving bin", nil, nil)
	if err != nil {
		t.Fatalf("create inventory item: %v", err)
	}
	invRun, err := app.RecordRun("bob", RunKindInventory, invItem.ID, 1, "completed", "Counted receiving bin", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record inventory run: %v", err)
	}
	if _, err := app.AddEvidence("bob", invRun.ID, "Cycle count", map[string]string{"expected_count": "12", "actual_count": "10", "discrepancy": "-2"}, "", nil); err != nil {
		t.Fatalf("add inventory evidence: %v", err)
	}
	_, err = app.RecordRun("bob", RunKindReceiving, recvItem.ID, 1, "accepted", "Non-problem run", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record normal run: %v", err)
	}
	server := NewServer(app)

	request := httptest.NewRequest(http.MethodGet, "/api/search?place_id="+place.ID+"&problem=true", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("problem search status: %d %s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	if !strings.Contains(body, `"problem":true`) {
		t.Fatalf("problem search missing filter echo: %s", body)
	}
	if !strings.Contains(body, recvRun.ID) || !strings.Contains(body, invRun.ID) {
		t.Fatalf("problem search missing expected problem runs: %s", body)
	}
	if strings.Contains(body, "Non-problem run") {
		t.Fatalf("problem search leaked non-problem run: %s", body)
	}
}
