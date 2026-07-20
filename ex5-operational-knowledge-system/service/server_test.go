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

func TestServerUploadsEvidenceAttachment(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindProcedure, item.ID, 1, "completed", "", "", "", nil)
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
