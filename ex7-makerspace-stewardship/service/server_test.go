package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPWorkflowRecordsSafetyHoldLoanAndReturn(t *testing.T) {
	app := NewDemoApp()
	handler := app.Handler()

	request := httptest.NewRequest(http.MethodPost, "/api/tools/table-saw/observations", bytes.NewBufferString(`{"reporterId":"alice","text":"Guard is loose","safetyHold":true,"photos":[{"name":"guard.png","dataUrl":"data:image/png;base64,aGVsbG8="}]}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("record safety hold status = %d, body = %s", response.Code, response.Body.String())
	}

	dueAt := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	request = httptest.NewRequest(http.MethodPost, "/api/tools/cordless-drill/loans", bytes.NewBufferString(`{"memberId":"alice","dueAt":"`+dueAt+`"}`))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("create loan status = %d, body = %s", response.Code, response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/api/tools/cordless-drill/returns", bytes.NewBufferString(`{"memberId":"alice","condition":"Returned with charger"}`))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("return loan status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestMalformedJSONStopsRequest(t *testing.T) {
	app := NewDemoApp()
	request := httptest.NewRequest(http.MethodPost, "/api/tools/table-saw/observations", bytes.NewBufferString(`not-json`))
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("malformed JSON status = %d, body = %s", response.Code, response.Body.String())
	}
}
