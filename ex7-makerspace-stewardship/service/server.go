package service

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex7-makerspace-stewardship/web"
)

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", web.Handler())
	mux.HandleFunc("GET /api/state", a.handleState)
	mux.HandleFunc("POST /api/records", a.handleRecordIngress)
	mux.HandleFunc("/api/approval-requests", a.handleApprovalRequests)
	mux.HandleFunc("/api/approval-requests/", a.handleApprovalRequest)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://127.0.0.1:7038" && strings.HasPrefix(r.URL.Path, "/api/approval-requests") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == http.MethodOptions && origin == "http://127.0.0.1:7038" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (a *App) handleApprovalRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a.handleApprovalRequestCreate(w, r)
		return
	}
	if r.Method == http.MethodGet && isLoopback(r.RemoteAddr) {
		requests, err := a.PendingApprovalRequests(time.Now())
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, requests)
		return
	}
	writeJSON(w, http.StatusForbidden, map[string]string{"error": "pending requests require loopback"})
}

func (a *App) handleApprovalRequestCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TargetPCID    string `json:"target_pcid"`
		PayloadBase64 string `json:"payload_base64"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	payload, err := base64.StdEncoding.DecodeString(input.PayloadBase64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "payload_base64 is invalid"})
		return
	}
	request, err := a.CreateApprovalRequest(input.TargetPCID, payload, time.Now())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusAccepted, request)
}

func (a *App) handleApprovalRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/approval-requests/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "request not found"})
		return
	}
	if len(parts) == 1 && r.Method == http.MethodGet {
		request, err := a.ApprovalRequest(parts[0], r.URL.Query().Get("token"), time.Now())
		if err != nil {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, request)
		return
	}
	if len(parts) == 2 && parts[1] == "approve" && r.Method == http.MethodPost {
		if !isLoopback(r.RemoteAddr) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "approval requires loopback"})
			return
		}
		request, err := a.ApproveRequest(parts[0], time.Now())
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, request)
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "request not found"})
}

func isLoopback(address string) bool {
	return strings.HasPrefix(address, "127.0.0.1:") || strings.HasPrefix(address, "[::1]:")
}

func (a *App) handleState(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.State())
}
func (a *App) handleRecordIngress(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Records []string `json:"records"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	if len(input.Records) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "records is required"})
		return
	}
	records := make([][]byte, 0, len(input.Records))
	for _, encoded := range input.Records {
		raw, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "record is not base64"})
			return
		}
		records = append(records, raw)
	}
	if err := a.IngestRecords(records); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusAccepted, a.State())
}
func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return false
	}
	return true
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}
