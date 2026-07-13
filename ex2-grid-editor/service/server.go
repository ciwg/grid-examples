package service

import (
	"encoding/json"
	"fmt"
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
	mux.HandleFunc("/api/local/documents/", server.handleLocalDocument)
	return mux
}

func (server *Server) handleIndex(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = writer.Write(web.MustRead("index.html"))
}

func (server *Server) handleAppJS(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	_, _ = writer.Write(web.MustRead("app.js"))
}

func (server *Server) handleStyleCSS(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/css; charset=utf-8")
	_, _ = writer.Write(web.MustRead("style.css"))
}

func (server *Server) handleMeta(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, server.app.Meta())
}

func (server *Server) handlePeerMessages(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		since := uint64(0)
		if raw := request.URL.Query().Get("since"); raw != "" {
			if _, err := fmt.Sscanf(raw, "%d", &since); err != nil {
				http.Error(writer, "invalid since", http.StatusBadRequest)
				return
			}
		}
		messages, nextOffset := server.app.PeerMessagesSince(since)
		writeJSON(writer, http.StatusOK, peerResponse{Messages: messages, NextOffset: nextOffset})
	case http.MethodPost:
		var payload struct {
			Messages []string `json:"messages"`
		}
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
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
	action := parts[1]
	switch action {
	case "state":
		if request.Method != http.MethodGet {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(writer, http.StatusOK, server.app.DocumentState(documentID))
	case "replace":
		if request.Method != http.MethodPost {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			Content    string `json:"content"`
			Embodiment string `json:"embodiment"`
		}
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err := server.app.ReplaceDocument(documentID, payload.Content, payload.Embodiment); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
	case "awareness":
		if request.Method != http.MethodPost {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			Cursor      int    `json:"cursor"`
			Head        int    `json:"head"`
			Typing      bool   `json:"typing"`
			DisplayName string `json:"display_name"`
			Color       string `json:"color"`
			Embodiment  string `json:"embodiment"`
		}
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err := server.app.UpdateAwareness(documentID, payload.Cursor, payload.Head, payload.Typing, payload.DisplayName, payload.Color, payload.Embodiment); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
	default:
		http.NotFound(writer, request)
	}
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
