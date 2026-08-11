package service

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/computerscienceiscool/grid-examples/ex7-makerspace-stewardship/web"
)

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", web.Handler())
	mux.HandleFunc("GET /api/state", a.handleState)
	mux.HandleFunc("POST /api/records", a.handleRecordIngress)
	return mux
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
