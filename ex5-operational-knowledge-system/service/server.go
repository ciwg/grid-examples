package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/web"
)

const maxEvidenceAttachmentBytes = 8 << 20

type Server struct {
	app *App
}

func NewServer(app *App) *Server {
	return &Server{app: app}
}

// Intent: Expose the shared operational model through a local HTTP adapter so
// the browser and CLI can be equal embodiments over one runtime without making
// HTTP itself the durable PromiseGrid-facing contract. Source: DI-radok;
// DI-zuvob
func (server *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handleIndex)
	mux.HandleFunc("/app.js", server.handleAppJS)
	mux.HandleFunc("/style.css", server.handleStyleCSS)
	mux.HandleFunc("/api/meta", server.handleMeta)
	mux.HandleFunc("/api/peer-exchange/export", server.handlePeerExchangeExport)
	mux.HandleFunc("/api/peer-exchange/import", server.handlePeerExchangeImport)
	mux.HandleFunc("/api/dashboard", server.handleDashboard)
	mux.HandleFunc("/api/problem-review", server.handleProblemReview)
	mux.HandleFunc("/api/search", server.handleSearch)
	mux.HandleFunc("/api/places", server.handlePlaces)
	mux.HandleFunc("/api/places/", server.handlePlace)
	mux.HandleFunc("/api/resources", server.handleResources)
	mux.HandleFunc("/api/resources/", server.handleResource)
	mux.HandleFunc("/api/responsibilities", server.handleResponsibilities)
	mux.HandleFunc("/api/responsibilities/", server.handleResponsibility)
	mux.HandleFunc("/api/items", server.handleItems)
	mux.HandleFunc("/api/items/", server.handleItem)
	mux.HandleFunc("/api/runs", server.handleRuns)
	mux.HandleFunc("/api/runs/", server.handleRun)
	mux.HandleFunc("/api/links", server.handleLinks)
	return mux
}

func (server *Server) handleIndex(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	writeNoStoreHeaders(writer)
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := writer.Write(web.MustRead("index.html")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) handleAppJS(writer http.ResponseWriter, request *http.Request) {
	writeNoStoreHeaders(writer)
	writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	if _, err := writer.Write(web.MustRead("app.js")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) handleStyleCSS(writer http.ResponseWriter, request *http.Request) {
	writeNoStoreHeaders(writer)
	writer.Header().Set("Content-Type", "text/css; charset=utf-8")
	if _, err := writer.Write(web.MustRead("style.css")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) handleMeta(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, server.app.Meta())
}

// Intent: Keep the first peer-visible ex5 exchange on the same local adapter
// surface as the current embodiments so browser, CLI, and test tooling can
// bootstrap signed-family exchange before any later embodiment/runtime
// tightening. Source: DI-voruk
func (server *Server) handlePeerExchangeExport(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	bundle, err := server.app.ExportPeerExchangeBundle()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(writer, http.StatusOK, bundle)
}

func (server *Server) handlePeerExchangeImport(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 8*1024*1024)
	var bundle PeerExchangeBundle
	if err := decodeJSONBody(request, &bundle); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := server.app.ImportPeerExchangeBundle(bundle)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, result)
}

func (server *Server) handleDashboard(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, server.app.Dashboard())
}

// Intent: Expose grouped receiving/count problem hotspots through the same
// local HTTP surface so browser operators can review repeated issues without
// reconstructing them by hand. Source: DI-pogul
func (server *Server) handleProblemReview(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, server.app.ProblemReview())
}

// Intent: Keep the local search endpoint useful for real operator drilldown by
// accepting structured filters, including problem-only review, not only one
// free-text query. Source: DI-honus; DI-vafuk; DI-vemur
func (server *Server) handleSearch(writer http.ResponseWriter, request *http.Request) {
	options := SearchOptions{
		Query:            request.URL.Query().Get("q"),
		Kind:             request.URL.Query().Get("kind"),
		Status:           request.URL.Query().Get("status"),
		Outcome:          request.URL.Query().Get("outcome"),
		PlaceID:          request.URL.Query().Get("place_id"),
		ResourceID:       request.URL.Query().Get("resource_id"),
		ResponsibilityID: request.URL.Query().Get("responsibility_id"),
		Problem:          request.URL.Query().Get("problem") == "true",
	}
	writeJSON(writer, http.StatusOK, server.app.SearchWithOptions(options))
}

func (server *Server) handlePlaces(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		writeJSON(writer, http.StatusOK, map[string]any{"places": server.app.ListPlaces()})
	case http.MethodPost:
		request.Body = http.MaxBytesReader(writer, request.Body, 64*1024)
		var payload struct {
			Actor    string   `json:"actor"`
			Kind     string   `json:"kind"`
			Name     string   `json:"name"`
			Summary  string   `json:"summary"`
			ParentID string   `json:"parent_id"`
			Tags     []string `json:"tags"`
		}
		if err := decodeJSONBody(request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		place, err := server.app.CreatePlace(payload.Actor, payload.Kind, payload.Name, payload.Summary, payload.ParentID, payload.Tags)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusCreated, place)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handlePlace(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(request.URL.Path, "/api/places/")
	place, err := server.app.GetPlace(strings.Trim(id, "/"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(writer, http.StatusOK, place)
}

func (server *Server) handleResources(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		writeJSON(writer, http.StatusOK, map[string]any{"resources": server.app.ListResources()})
	case http.MethodPost:
		request.Body = http.MaxBytesReader(writer, request.Body, 64*1024)
		var payload struct {
			Actor   string   `json:"actor"`
			Kind    string   `json:"kind"`
			Name    string   `json:"name"`
			Summary string   `json:"summary"`
			PlaceID string   `json:"place_id"`
			Tags    []string `json:"tags"`
		}
		if err := decodeJSONBody(request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		resource, err := server.app.CreateResource(payload.Actor, payload.Kind, payload.Name, payload.Summary, payload.PlaceID, payload.Tags)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusCreated, resource)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handleResource(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(request.URL.Path, "/api/resources/")
	resource, err := server.app.GetResource(strings.Trim(id, "/"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(writer, http.StatusOK, resource)
}

func (server *Server) handleResponsibilities(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		writeJSON(writer, http.StatusOK, map[string]any{"responsibilities": server.app.ListResponsibilities()})
	case http.MethodPost:
		request.Body = http.MaxBytesReader(writer, request.Body, 64*1024)
		var payload struct {
			Actor    string   `json:"actor"`
			Title    string   `json:"title"`
			Summary  string   `json:"summary"`
			RoleKeys []string `json:"role_keys"`
			Tags     []string `json:"tags"`
		}
		if err := decodeJSONBody(request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		record, err := server.app.CreateResponsibility(payload.Actor, payload.Title, payload.Summary, payload.RoleKeys, payload.Tags)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusCreated, record)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handleResponsibility(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(request.URL.Path, "/api/responsibilities/")
	record, err := server.app.GetResponsibility(strings.Trim(id, "/"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(writer, http.StatusOK, record)
}

func (server *Server) handleItems(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		items, err := server.app.ListKnowledgeItems(request.URL.Query().Get("kind"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{"items": items})
	case http.MethodPost:
		request.Body = http.MaxBytesReader(writer, request.Body, 256*1024)
		var payload struct {
			Actor             string   `json:"actor"`
			Kind              string   `json:"kind"`
			Title             string   `json:"title"`
			Summary           string   `json:"summary"`
			Body              string   `json:"body"`
			Tags              []string `json:"tags"`
			ResponsibilityIDs []string `json:"responsibility_ids"`
		}
		if err := decodeJSONBody(request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		item, err := server.app.CreateKnowledgeItem(payload.Actor, payload.Kind, payload.Title, payload.Summary, payload.Body, payload.Tags, payload.ResponsibilityIDs)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusCreated, item)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handleItem(writer http.ResponseWriter, request *http.Request) {
	path := strings.Trim(strings.TrimPrefix(request.URL.Path, "/api/items/"), "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(writer, request)
		return
	}
	itemID := parts[0]
	if len(parts) == 1 {
		if request.Method != http.MethodGet {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		item, err := server.app.GetKnowledgeItem(itemID)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(writer, http.StatusOK, item)
		return
	}
	switch parts[1] {
	case "revisions":
		server.handleRevision(writer, request, itemID)
	case "approvals":
		server.handleItemApproval(writer, request, itemID)
	case "supersede":
		server.handleItemSupersede(writer, request, itemID)
	case "live":
		if len(parts) == 3 && parts[2] == "socket" {
			server.handleItemLiveSocket(writer, request, itemID)
			return
		}
		server.handleItemLive(writer, request, itemID)
	default:
		http.NotFound(writer, request)
	}
}

func (server *Server) handleRevision(writer http.ResponseWriter, request *http.Request, itemID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 256*1024)
	var payload struct {
		Actor   string   `json:"actor"`
		Title   string   `json:"title"`
		Summary string   `json:"summary"`
		Body    string   `json:"body"`
		Tags    []string `json:"tags"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := server.app.AddRevision(payload.Actor, itemID, payload.Title, payload.Summary, payload.Body, payload.Tags)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, item)
}

func (server *Server) handleItemApproval(writer http.ResponseWriter, request *http.Request, itemID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 32*1024)
	var payload struct {
		Actor    string `json:"actor"`
		Revision int    `json:"revision"`
		Role     string `json:"role"`
		Decision string `json:"decision"`
		Notes    string `json:"notes"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := server.app.RecordApproval(payload.Actor, "knowledge_item", itemID, payload.Revision, payload.Role, payload.Decision, payload.Notes); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := server.app.GetKnowledgeItem(itemID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, item)
}

func (server *Server) handleItemSupersede(writer http.ResponseWriter, request *http.Request, itemID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 16*1024)
	var payload struct {
		Actor string `json:"actor"`
		Notes string `json:"notes"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := server.app.SupersedeKnowledgeItem(payload.Actor, itemID, payload.Notes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, item)
}

// Intent: Expose the shared working draft and participant presence for a
// knowledge item through the same local adapter as the rest of the workflow,
// while keeping durable revisions and approvals separate from live draft
// traffic. Source: DI-lusov; DI-zoruk
func (server *Server) handleItemLive(writer http.ResponseWriter, request *http.Request, itemID string) {
	switch request.Method {
	case http.MethodGet:
		state, err := server.app.LiveItemState(itemID)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(writer, http.StatusOK, state)
	case http.MethodPost:
		request.Body = http.MaxBytesReader(writer, request.Body, 256*1024)
		var payload struct {
			ParticipantID string `json:"participant_id"`
			DisplayName   string `json:"display_name"`
			Color         string `json:"color"`
			Cursor        int    `json:"cursor"`
			Head          int    `json:"head"`
			Typing        bool   `json:"typing"`
			BaseVersion   int    `json:"base_version"`
			UpdateBody    bool   `json:"update_body"`
			Body          string `json:"body"`
		}
		if err := decodeJSONBody(request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		state, conflict, err := server.app.UpdateLiveItem(itemID, payload.ParticipantID, payload.DisplayName, payload.Color, payload.Cursor, payload.Head, payload.Typing, payload.BaseVersion, payload.UpdateBody, payload.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if conflict {
			writeJSON(writer, http.StatusConflict, map[string]any{
				"conflict": true,
				"state":    state,
			})
			return
		}
		writeJSON(writer, http.StatusOK, state)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handleRuns(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		runs, err := server.app.ListRuns(request.URL.Query().Get("kind"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{"runs": runs})
	case http.MethodPost:
		request.Body = http.MaxBytesReader(writer, request.Body, 64*1024)
		var payload struct {
			Actor             string   `json:"actor"`
			Kind              string   `json:"kind"`
			ItemID            string   `json:"item_id"`
			Revision          int      `json:"revision"`
			Outcome           string   `json:"outcome"`
			Notes             string   `json:"notes"`
			PlaceID           string   `json:"place_id"`
			ResourceIDs       []string `json:"resource_ids"`
			Machine           string   `json:"machine"`
			Location          string   `json:"location"`
			ResponsibilityIDs []string `json:"responsibility_ids"`
		}
		if err := decodeJSONBody(request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		run, err := server.app.RecordRun(payload.Actor, payload.Kind, payload.ItemID, payload.Revision, payload.Outcome, payload.Notes, payload.Machine, payload.Location, payload.PlaceID, payload.ResourceIDs, payload.ResponsibilityIDs)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusCreated, run)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handleRun(writer http.ResponseWriter, request *http.Request) {
	path := strings.Trim(strings.TrimPrefix(request.URL.Path, "/api/runs/"), "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(writer, request)
		return
	}
	runID := parts[0]
	if len(parts) == 1 {
		if request.Method != http.MethodGet {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		run, err := server.app.GetRun(runID)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(writer, http.StatusOK, run)
		return
	}
	switch parts[1] {
	case "evidence":
		server.handleEvidence(writer, request, runID)
	case "approvals":
		server.handleRunApproval(writer, request, runID)
	default:
		http.NotFound(writer, request)
	}
}

func (server *Server) handleEvidence(writer http.ResponseWriter, request *http.Request, runID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxEvidenceAttachmentBytes+64*1024)
	if err := request.ParseMultipartForm(maxEvidenceAttachmentBytes); err != nil {
		http.Error(writer, fmt.Sprintf("parse multipart form: %v", err), http.StatusBadRequest)
		return
	}
	if request.MultipartForm != nil {
		defer func() {
			if err := request.MultipartForm.RemoveAll(); err != nil {
				fmt.Fprintf(os.Stderr, "oks multipart cleanup: %v\n", err)
			}
		}()
	}
	actor := request.FormValue("actor")
	summary := request.FormValue("summary")
	facts := map[string]string{}
	if rawFacts := strings.TrimSpace(request.FormValue("facts_json")); rawFacts != "" {
		if err := json.Unmarshal([]byte(rawFacts), &facts); err != nil {
			http.Error(writer, fmt.Sprintf("decode facts_json: %v", err), http.StatusBadRequest)
			return
		}
	}
	var attachmentName string
	var attachmentBody []byte
	file, header, err := request.FormFile("attachment")
	if err == nil {
		attachmentName = header.Filename
		body, readErr := io.ReadAll(io.LimitReader(file, maxEvidenceAttachmentBytes+1))
		closeErr := file.Close()
		if readErr != nil {
			http.Error(writer, fmt.Sprintf("read attachment: %v", readErr), http.StatusBadRequest)
			return
		}
		if closeErr != nil {
			http.Error(writer, fmt.Sprintf("close attachment: %v", closeErr), http.StatusBadRequest)
			return
		}
		// Intent: Reject evidence uploads that exceed the documented attachment
		// limit instead of silently accepting truncated overflow bytes. Source:
		// DI-navos
		if len(body) > maxEvidenceAttachmentBytes {
			http.Error(writer, fmt.Sprintf("attachment exceeds %d bytes", maxEvidenceAttachmentBytes), http.StatusBadRequest)
			return
		}
		attachmentBody = body
	}
	run, err := server.app.AddEvidence(actor, runID, summary, facts, attachmentName, attachmentBody)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, run)
}

func (server *Server) handleRunApproval(writer http.ResponseWriter, request *http.Request, runID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 32*1024)
	var payload struct {
		Actor    string `json:"actor"`
		Role     string `json:"role"`
		Decision string `json:"decision"`
		Notes    string `json:"notes"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := server.app.RecordApproval(payload.Actor, "run", runID, 0, payload.Role, payload.Decision, payload.Notes); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	run, err := server.app.GetRun(runID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, run)
}

func (server *Server) handleLinks(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 32*1024)
	var payload struct {
		Actor    string `json:"actor"`
		FromType string `json:"from_type"`
		FromID   string `json:"from_id"`
		ToType   string `json:"to_type"`
		ToID     string `json:"to_id"`
		Relation string `json:"relation"`
		Notes    string `json:"notes"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := server.app.AddLink(payload.Actor, payload.FromType, payload.FromID, payload.ToType, payload.ToID, payload.Relation, payload.Notes); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]any{"ok": true})
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(payload); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func decodeJSONBody(request *http.Request, target any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return nil
}

func writeNoStoreHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Cache-Control", "no-store, max-age=0")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
}
