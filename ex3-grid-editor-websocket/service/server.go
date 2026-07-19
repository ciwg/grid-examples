package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/web"
)

type Server struct {
	app *App

	mu                   sync.Mutex
	syncSubscribers      map[string]map[chan struct{}]struct{}
	awarenessSubscribers map[string]map[chan struct{}]struct{}
}

func NewServer(app *App) *Server {
	server := &Server{
		app:                  app,
		syncSubscribers:      map[string]map[chan struct{}]struct{}{},
		awarenessSubscribers: map[string]map[chan struct{}]struct{}{},
	}
	app.SetLiveChangeHooks(server.notifySyncSubscribers, server.notifyAwarenessSubscribers)
	return server
}

func (server *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handleIndex)
	mux.HandleFunc("/app.js", server.handleAppJS)
	mux.HandleFunc("/style.css", server.handleStyleCSS)
	mux.HandleFunc("/api/meta", server.handleMeta)
	mux.HandleFunc("/api/peer/messages", server.handlePeerMessages)
	mux.HandleFunc("/api/local/metadata/search", server.handleMetadataSearch)
	mux.HandleFunc("/api/published/", server.handlePublished)
	mux.HandleFunc("/api/local/documents/", server.handleLocalDocument)
	return mux
}

func (server *Server) handleIndex(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	setStaticAssetHeaders(writer)
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := writer.Write(web.MustRead("index.html")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) handleAppJS(writer http.ResponseWriter, request *http.Request) {
	setStaticAssetHeaders(writer)
	writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	if _, err := writer.Write(web.MustRead("app.js")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) handleStyleCSS(writer http.ResponseWriter, request *http.Request) {
	setStaticAssetHeaders(writer)
	writer.Header().Set("Content-Type", "text/css; charset=utf-8")
	if _, err := writer.Write(web.MustRead("style.css")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func setStaticAssetHeaders(writer http.ResponseWriter) {
	// Intent: Force browsers to revalidate ex3's embedded UI assets after relay
	// upgrades so stale cached bundles do not keep calling the retired
	// bootstrap path and fail before the user can interact. Source: DI-povip
	writer.Header().Set("Cache-Control", "no-store, max-age=0")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
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
	case "trace":
		if request.Method != http.MethodGet {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		limit := defaultFeedLimit
		if raw := request.URL.Query().Get("limit"); raw != "" {
			if _, err := fmt.Sscanf(raw, "%d", &limit); err != nil {
				http.Error(writer, "invalid limit", http.StatusBadRequest)
				return
			}
		}
		writeJSON(writer, http.StatusOK, server.app.Trace(documentID, limit))
	case "sync":
		server.handleSync(writer, request, documentID)
	case "session":
		server.handleSession(writer, request, documentID)
	case "sync-socket":
		server.handleSyncSocket(writer, request, documentID)
	case "awareness":
		server.handleAwareness(writer, request, documentID)
	case "awareness-socket":
		server.handleAwarenessSocket(writer, request, documentID)
	case "published":
		server.handlePublishedList(writer, request, documentID)
	case "publish":
		server.handlePublish(writer, request, documentID)
	case "metadata":
		server.handleMetadata(writer, request, documentID)
	default:
		http.NotFound(writer, request)
	}
}

func (server *Server) handleMetadataSearch(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := request.URL.Query().Get("q")
	includeArchived := request.URL.Query().Get("include_archived") == "true"
	writeJSON(writer, http.StatusOK, map[string]any{
		"query":            query,
		"include_archived": includeArchived,
		"results":          server.app.SearchMetadata(query, includeArchived),
	})
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
		if err := server.requireMutationAccess(request, payload.ParticipantID, documentID, server.app.documentPCID.String()); err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
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
		if err := server.requireMutationAccess(request, payload.ParticipantID, documentID, server.app.awarenessPCID.String()); err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
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

func (server *Server) handleSession(writer http.ResponseWriter, request *http.Request, documentID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 8*1024)
	var payload struct {
		ParticipantID string `json:"participant_id"`
	}
	if err := decodeJSONBody(writer, request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := server.requireSessionBootstrap(request); err != nil {
		http.Error(writer, err.Error(), http.StatusForbidden)
		return
	}
	session, err := server.app.IssueRemoteSession(documentID, payload.ParticipantID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, session)
}

func (server *Server) handleMetadata(writer http.ResponseWriter, request *http.Request, documentID string) {
	switch request.Method {
	case http.MethodGet:
		writeJSON(writer, http.StatusOK, server.app.Metadata(documentID))
	case http.MethodPost:
		request.Body = http.MaxBytesReader(writer, request.Body, 128*1024)
		var payload struct {
			ParticipantID string   `json:"participant_id"`
			Title         string   `json:"title"`
			Description   string   `json:"description"`
			Summary       string   `json:"summary"`
			Tags          []string `json:"tags"`
			Collections   []string `json:"collections"`
			Favorite      bool     `json:"favorite"`
			Archived      bool     `json:"archived"`
			Embodiment    string   `json:"embodiment"`
		}
		if err := decodeJSONBody(writer, request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err := server.requireMutationAccess(request, payload.ParticipantID, documentID, server.app.metadataPCID.String()); err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		record, err := server.app.UpdateMetadata(
			documentID,
			payload.ParticipantID,
			payload.Title,
			payload.Description,
			payload.Summary,
			payload.Tags,
			payload.Collections,
			payload.Favorite,
			payload.Archived,
			payload.Embodiment,
		)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusOK, record)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handlePublish(writer http.ResponseWriter, request *http.Request, documentID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
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
	if err := server.requireMutationAccess(request, payload.ParticipantID, documentID, server.app.publishPCID.String()); err != nil {
		http.Error(writer, err.Error(), http.StatusForbidden)
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
	// Intent: Publish through a separate relay endpoint so clients can request
	// durable exchange manifests without overloading live CRDT sync, while
	// reusing the same remote-capability gate as the other document mutation
	// surfaces. Source: DI-tavul; DI-gosaf; DI-povip
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
		"manifest_url": fmt.Sprintf("%s/api/published/%s", requestBaseURL(request), record.EnvelopeCID),
	})
}

func requestBaseURL(request *http.Request) string {
	scheme := "http"
	if forwarded := strings.TrimSpace(request.Header.Get("X-Forwarded-Proto")); forwarded != "" {
		if comma := strings.Index(forwarded, ","); comma >= 0 {
			forwarded = forwarded[:comma]
		}
		forwarded = strings.TrimSpace(forwarded)
		if forwarded != "" {
			scheme = forwarded
		}
	} else if request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, request.Host)
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) requireMutationAccess(request *http.Request, participantID string, documentID string, protocolCID string) error {
	// Intent: Preserve loopback's zero-config local workflow while requiring a
	// short-lived relay-signed capability for non-loopback mutation, so ex3 can
	// run across machines without exposing an unauthenticated signer. Source:
	// DI-povip
	if requestIsLoopback(request) {
		return nil
	}
	if participantID == "" {
		return fmt.Errorf("participant_id is required for remote mutation")
	}
	token := bearerToken(request)
	if token == "" {
		return fmt.Errorf("remote mutation requires a bearer capability")
	}
	if _, err := verifyMutationCapability(token, participantID, documentID, protocolCID, "mutate", time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

func (server *Server) requireSessionBootstrap(request *http.Request) error {
	// Intent: Keep the long-lived admission secret out of ex3's steady-state
	// mutation traffic by using it only for remote session bootstrap, while
	// allowing loopback clients to keep the old zero-config local path. Source:
	// DI-povip
	if requestIsLoopback(request) {
		return nil
	}
	if !server.app.RemoteAccessEnabled() {
		return fmt.Errorf("remote mutation bootstrap is disabled")
	}
	if server.app.ValidateRemoteAccessToken(request.Header.Get("X-Grid-Access-Token")) {
		return nil
	}
	if server.app.ValidateRemoteAccessToken(request.URL.Query().Get("access_token")) {
		return nil
	}
	return fmt.Errorf("remote session bootstrap requires X-Grid-Access-Token")
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

func (server *Server) subscribeSync(documentID string) (chan struct{}, func()) {
	return server.subscribe(server.syncSubscribers, documentID)
}

func (server *Server) subscribeAwareness(documentID string) (chan struct{}, func()) {
	return server.subscribe(server.awarenessSubscribers, documentID)
}

func (server *Server) subscribe(index map[string]map[chan struct{}]struct{}, documentID string) (chan struct{}, func()) {
	channel := make(chan struct{}, 1)
	server.mu.Lock()
	if index[documentID] == nil {
		index[documentID] = map[chan struct{}]struct{}{}
	}
	index[documentID][channel] = struct{}{}
	server.mu.Unlock()
	return channel, func() {
		server.mu.Lock()
		subscribers := index[documentID]
		if subscribers != nil {
			delete(subscribers, channel)
			if len(subscribers) == 0 {
				delete(index, documentID)
			}
		}
		server.mu.Unlock()
		close(channel)
	}
}

func (server *Server) notifySyncSubscribers(documentID string) {
	server.notifySubscribers(server.syncSubscribers, documentID)
}

func (server *Server) notifyAwarenessSubscribers(documentID string) {
	server.notifySubscribers(server.awarenessSubscribers, documentID)
}

func (server *Server) notifySubscribers(index map[string]map[chan struct{}]struct{}, documentID string) {
	server.mu.Lock()
	subscribers := make([]chan struct{}, 0, len(index[documentID]))
	for channel := range index[documentID] {
		subscribers = append(subscribers, channel)
	}
	server.mu.Unlock()
	for _, channel := range subscribers {
		select {
		case channel <- struct{}{}:
		default:
		}
	}
}
