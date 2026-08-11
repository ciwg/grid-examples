package service

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPSignedRecordIngressProjectsRecognizedEvidence(t *testing.T) {
	record, aliceKey := signedTestRecord(t, observationPCID, "obs-1", "alice", `{"observation":"Guard is loose","tool_id":"table-saw"}`)
	app, err := NewPersistentRecordApp(t.TempDir(), NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": aliceKey}))
	if err != nil {
		t.Fatalf("create record app: %v", err)
	}
	bootstrap, _ := testParticipantBootstrap(t, "alice")
	encoded := make([]string, 0, len(bootstrap)+1)
	for _, raw := range append(bootstrap, record) {
		encoded = append(encoded, base64.StdEncoding.EncodeToString(raw))
	}
	request := httptest.NewRequest(http.MethodPost, "/api/records", bytes.NewBufferString(`{"records":["`+encoded[0]+`","`+encoded[1]+`","`+encoded[2]+`"]}`))
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusAccepted {
		t.Fatalf("ingress status = %d, body = %s", response.Code, response.Body.String())
	}
	if got := len(app.State().Tools[0].Observations); got != 1 {
		t.Fatalf("projected observations = %d", got)
	}
}

func TestHTTPRecordIngressRejectsMalformedBase64(t *testing.T) {
	app, err := NewPersistentRecordApp(t.TempDir(), RecognitionPolicy{})
	if err != nil {
		t.Fatalf("create record app: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/records", bytes.NewBufferString(`{"records":["not-base64"]}`))
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("malformed record status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestTwoAgentTerminalApprovalHTTPFlow(t *testing.T) {
	aliceIdentity := &ParticipantIdentity{Label: "alice", Private: testPrivate("alice")}
	alice, err := NewPersistentParticipantApp(t.TempDir(), NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": aliceIdentity.Private.Public().(ed25519.PublicKey)}), aliceIdentity)
	if err != nil {
		t.Fatalf("create Alice agent: %v", err)
	}
	bootstrap, _ := testParticipantBootstrap(t, "alice")
	if err := alice.IngestRecords(bootstrap); err != nil {
		t.Fatalf("bootstrap Alice: %v", err)
	}
	terminal, err := NewPersistentRecordApp(t.TempDir(), NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": aliceIdentity.Private.Public().(ed25519.PublicKey)}))
	if err != nil {
		t.Fatalf("create terminal agent: %v", err)
	}
	if err := terminal.IngestRecords(bootstrap); err != nil {
		t.Fatalf("bootstrap terminal view: %v", err)
	}

	create := httptest.NewRequest(http.MethodPost, "/api/approval-requests", bytes.NewBufferString(`{"target_pcid":"`+observationPCID+`","payload_base64":"`+base64.StdEncoding.EncodeToString([]byte(`{"observation":"Terminal E2E","tool_id":"table-saw"}`))+`"}`))
	created := httptest.NewRecorder()
	alice.Handler().ServeHTTP(created, create)
	if created.Code != http.StatusAccepted {
		t.Fatalf("create status = %d: %s", created.Code, created.Body.String())
	}
	var request ApprovalRequest
	if err := json.NewDecoder(created.Body).Decode(&request); err != nil {
		t.Fatalf("decode request: %v", err)
	}
	if got := len(alice.State().Tools[0].Observations); got != 0 {
		t.Fatalf("pending request projected %d", got)
	}

	approve := httptest.NewRequest(http.MethodPost, "/api/approval-requests/"+request.RequestID+"/approve", nil)
	approve.RemoteAddr = "127.0.0.1:7037"
	approved := httptest.NewRecorder()
	alice.Handler().ServeHTTP(approved, approve)
	if approved.Code != http.StatusOK {
		t.Fatalf("approval status = %d: %s", approved.Code, approved.Body.String())
	}
	poll := httptest.NewRequest(http.MethodGet, "/api/approval-requests/"+request.RequestID+"?token="+request.ApprovalToken, nil)
	polled := httptest.NewRecorder()
	alice.Handler().ServeHTTP(polled, poll)
	if polled.Code != http.StatusOK {
		t.Fatalf("poll status = %d: %s", polled.Code, polled.Body.String())
	}
	if err := json.NewDecoder(polled.Body).Decode(&request); err != nil {
		t.Fatalf("decode approved request: %v", err)
	}
	raw, err := base64.StdEncoding.DecodeString(request.SignedRecordBase64)
	if err != nil {
		t.Fatalf("decode signed response: %v", err)
	}
	if err := terminal.IngestRecords([][]byte{raw}); err != nil {
		t.Fatalf("terminal ingress signed response: %v", err)
	}
	if got := len(terminal.State().Tools[0].Observations); got != 1 {
		t.Fatalf("terminal projection = %d", got)
	}
}
