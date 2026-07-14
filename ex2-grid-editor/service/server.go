package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/web"
)

type Server struct {
	app *App
}

func NewServer(app *App) *Server {
	return &Server{app: app}
}

func (server *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handleIndex)
	mux.HandleFunc("/app.js", server.handleAppJS)
	mux.HandleFunc("/style.css", server.handleStyleCSS)
	mux.HandleFunc("/api/meta", server.handleMeta)
	mux.HandleFunc("/api/peer/messages", server.handlePeerMessages)
	mux.HandleFunc("/api/published/", server.handlePublished)
	mux.HandleFunc("/api/local/documents/", server.handleLocalDocument)
	return mux
}

func (server *Server) handleIndex(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := writer.Write(web.MustRead("index.html")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) handleAppJS(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	if _, err := writer.Write(web.MustRead("app.js")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) handleStyleCSS(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/css; charset=utf-8")
	if _, err := writer.Write(web.MustRead("style.css")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) handleMeta(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, server.app.Meta())
}

func (server *Server) handlePeerMessages(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		since := uint64(0)
		limit := defaultFeedLimit
		if raw := request.URL.Query().Get("since"); raw != "" {
			if _, err := fmt.Sscanf(raw, "%d", &since); err != nil {
				http.Error(writer, "invalid since", http.StatusBadRequest)
				return
			}
		}
		if raw := request.URL.Query().Get("limit"); raw != "" {
			if _, err := fmt.Sscanf(raw, "%d", &limit); err != nil {
				http.Error(writer, "invalid limit", http.StatusBadRequest)
				return
			}
		}
		messages, nextOffset := server.app.PeerMessagesSince(since, limit)
		writeJSON(writer, http.StatusOK, peerResponse{Messages: messages, NextOffset: nextOffset})
	case http.MethodPost:
		request.Body = http.MaxBytesReader(writer, request.Body, maxChangeBytesLen*8)
		var payload struct {
			Messages []string `json:"messages"`
		}
		if err := decodeJSONBody(writer, request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		for _, raw := range payload.Messages {
			if err := server.app.IngestRawBase64(raw); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handleLocalDocument(writer http.ResponseWriter, request *http.Request) {
	path := strings.TrimPrefix(request.URL.Path, "/api/local/documents/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 {
		http.NotFound(writer, request)
		return
	}
	documentID := parts[0]
	if err := validateDocumentID(documentID); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	action := parts[1]
	switch action {
	case "state":
		if request.Method != http.MethodGet {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(writer, http.StatusOK, server.app.State(documentID))
	case "sync":
		server.handleSync(writer, request, documentID)
	case "awareness":
		server.handleAwareness(writer, request, documentID)
	case "published":
		server.handlePublishedList(writer, request, documentID)
	case "publish":
		server.handlePublish(writer, request, documentID)
	default:
		http.NotFound(writer, request)
	}
}

func (server *Server) handlePublished(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Intent: Keep publish/import exchange fetchable by explicit manifest URL so
	// importers resolve a durable handoff artifact instead of silently joining
	// the live relay feed. Source: DI-tavul; DI-gosaf
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	envelopeCID := strings.TrimPrefix(request.URL.Path, "/api/published/")
	if envelopeCID == "" {
		http.NotFound(writer, request)
		return
	}
	resolved, err := server.app.ResolvePublished(envelopeCID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(writer, http.StatusOK, resolved)
}

func (server *Server) handleSync(writer http.ResponseWriter, request *http.Request, documentID string) {
	switch request.Method {
	case http.MethodGet:
		since := uint64(0)
		limit := defaultFeedLimit
		if raw := request.URL.Query().Get("since"); raw != "" {
			if _, err := fmt.Sscanf(raw, "%d", &since); err != nil {
				http.Error(writer, "invalid since", http.StatusBadRequest)
				return
			}
		}
		if raw := request.URL.Query().Get("limit"); raw != "" {
			if _, err := fmt.Sscanf(raw, "%d", &limit); err != nil {
				http.Error(writer, "invalid limit", http.StatusBadRequest)
				return
			}
		}
		writeJSON(writer, http.StatusOK, server.app.SyncFeed(documentID, since, limit))
	case http.MethodPost:
		if err := requireLoopbackMutation(request); err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		request.Body = http.MaxBytesReader(writer, request.Body, maxChangeBytesLen*4)
		var payload struct {
			ParticipantID string `json:"participant_id"`
			RecipientID   string `json:"recipient_id"`
			MessageBase64 string `json:"message_base64"`
			Embodiment    string `json:"embodiment"`
		}
		if err := decodeJSONBody(writer, request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		// Intent: Keep the relay HTTP surface explicit about CRDT sync so browser
		// and sidecar clients exchange signed Automerge messages through a stable
		// endpoint instead of overloading snapshot-oriented routes. Source:
		// DI-ramuv; DI-lumek
		record, err := server.app.PostSync(documentID, payload.ParticipantID, payload.RecipientID, payload.MessageBase64, payload.Embodiment)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusOK, record)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handleAwareness(writer http.ResponseWriter, request *http.Request, documentID string) {
	switch request.Method {
	case http.MethodGet:
		writeJSON(writer, http.StatusOK, map[string]any{
			"document_id": documentID,
			"awareness":   server.app.AwarenessState(documentID),
		})
	case http.MethodPost:
		if err := requireLoopbackMutation(request); err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		request.Body = http.MaxBytesReader(writer, request.Body, 16*1024)
		var payload struct {
			ParticipantID string `json:"participant_id"`
			Cursor        int    `json:"cursor"`
			Head          int    `json:"head"`
			Typing        bool   `json:"typing"`
			DisplayName   string `json:"display_name"`
			Color         string `json:"color"`
			Embodiment    string `json:"embodiment"`
		}
		if err := decodeJSONBody(writer, request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err := server.app.UpdateAwareness(documentID, payload.ParticipantID, payload.Cursor, payload.Head, payload.Typing, payload.DisplayName, payload.Color, payload.Embodiment); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handlePublishedList(writer http.ResponseWriter, request *http.Request, documentID string) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{
		"document_id": documentID,
		"published":   server.app.Published(documentID),
	})
}

func (server *Server) handlePublish(writer http.ResponseWriter, request *http.Request, documentID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := requireLoopbackMutation(request); err != nil {
		http.Error(writer, err.Error(), http.StatusForbidden)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxPublishBytesLen*12)
	var payload struct {
		ParticipantID     string `json:"participant_id"`
		SourceKind        string `json:"source_kind"`
		SourceVersionID   string `json:"source_version_id"`
		SourceVersionName string `json:"source_version_name"`
		Title             string `json:"title"`
		Summary           string `json:"summary"`
		TextBase64        string `json:"text_base64"`
		ReplicaBase64     string `json:"replica_base64"`
		Embodiment        string `json:"embodiment"`
	}
	if err := decodeJSONBody(writer, request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	textBytes, err := base64.StdEncoding.DecodeString(payload.TextBase64)
	if err != nil {
		http.Error(writer, "invalid text base64", http.StatusBadRequest)
		return
	}
	replicaBytes, err := base64.StdEncoding.DecodeString(payload.ReplicaBase64)
	if err != nil {
		http.Error(writer, "invalid replica base64", http.StatusBadRequest)
		return
	}
	// Intent: Publish through a separate loopback-only relay endpoint so local
	// clients can request durable exchange manifests without overloading live
	// CRDT sync or exposing an unauthenticated remote signer. Source: DI-tavul;
	// DI-gosaf; DI-rabod
	record, err := server.app.PublishDocument(
		documentID,
		payload.ParticipantID,
		payload.SourceKind,
		payload.SourceVersionID,
		payload.SourceVersionName,
		payload.Title,
		payload.Summary,
		textBytes,
		replicaBytes,
		payload.Embodiment,
	)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{
		"record":       record,
		"manifest_url": fmt.Sprintf("http://%s/api/published/%s", request.Host, record.EnvelopeCID),
	})
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func requireLoopbackMutation(request *http.Request) error {
	// Intent: Keep the relay from acting as an open network signer by allowing
	// local mutation requests only from loopback clients unless a later
	// authenticated remote-client mode is explicitly designed. Source: DI-rabod
	if requestIsLoopback(request) {
		return nil
	}
	return fmt.Errorf("local mutation endpoints require a loopback client")
}

func decodeJSONBody(writer http.ResponseWriter, request *http.Request, value any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("unexpected trailing JSON")
	}
	return nil
}
