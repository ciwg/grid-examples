package service

import (
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
	mux.HandleFunc("POST /api/tools/{id}/observations", a.handleObservation)
	mux.HandleFunc("POST /api/tools/{id}/clear-safety-hold", a.handleClearSafetyHold)
	mux.HandleFunc("POST /api/tools/{id}/loans", a.handleLoan)
	mux.HandleFunc("POST /api/tools/{id}/returns", a.handleReturn)
	return mux
}

func (a *App) handleState(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.State())
}
func (a *App) handleObservation(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ReporterID string  `json:"reporterId"`
		Text       string  `json:"text"`
		SafetyHold bool    `json:"safetyHold"`
		Photos     []Photo `json:"photos"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	tool, err := a.AddObservation(r.PathValue("id"), input.ReporterID, strings.TrimSpace(input.Text), input.SafetyHold, input.Photos)
	writeResult(w, tool, err)
}
func (a *App) handleClearSafetyHold(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StewardID  string `json:"stewardId"`
		Assessment string `json:"assessment"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	tool, err := a.ClearSafetyHold(r.PathValue("id"), input.StewardID, strings.TrimSpace(input.Assessment))
	writeResult(w, tool, err)
}
func (a *App) handleLoan(w http.ResponseWriter, r *http.Request) {
	var input struct {
		MemberID string `json:"memberId"`
		DueAt    string `json:"dueAt"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	dueAt, err := time.Parse(time.RFC3339, input.DueAt)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "dueAt must be RFC3339"})
		return
	}
	tool, err := a.CreateLoan(r.PathValue("id"), input.MemberID, dueAt)
	writeResult(w, tool, err)
}
func (a *App) handleReturn(w http.ResponseWriter, r *http.Request) {
	var input struct {
		MemberID  string `json:"memberId"`
		Condition string `json:"condition"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	tool, err := a.ReturnLoan(r.PathValue("id"), input.MemberID, strings.TrimSpace(input.Condition))
	writeResult(w, tool, err)
}
func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return false
	}
	return true
}
func writeResult(w http.ResponseWriter, tool Tool, err error) {
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, tool)
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}
