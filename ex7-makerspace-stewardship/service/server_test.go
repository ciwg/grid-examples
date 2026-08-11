package service

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
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
	request := httptest.NewRequest(http.MethodPost, "/api/records", bytes.NewBufferString(`{"records":["`+base64.StdEncoding.EncodeToString(record)+`"]}`))
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
