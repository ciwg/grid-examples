package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/web"
)

const maxAttachmentBytes = 8 << 20

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
	mux.HandleFunc("/api/issues", server.handleIssues)
	mux.HandleFunc("/api/issues/", server.handleIssue)
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

func writeNoStoreHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Cache-Control", "no-store, max-age=0")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
}

func (server *Server) handleMeta(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, server.app.Meta())
}

func (server *Server) handleIssues(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		issues, err := server.app.ListIssues(request.URL.Query().Get("status"), request.URL.Query().Get("assignee"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{"issues": issues})
	case http.MethodPost:
		actor, err := requestActor(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		request.Body = http.MaxBytesReader(writer, request.Body, 32*1024)
		var payload struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Severity    string `json:"severity"`
		}
		if err := decodeJSONBody(request, &payload); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		issue, err := server.app.CreateIssue(actor, payload.Title, payload.Description, payload.Severity)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusCreated, issue)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) handleIssue(writer http.ResponseWriter, request *http.Request) {
	path := strings.TrimPrefix(request.URL.Path, "/api/issues/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(writer, request)
		return
	}
	issueID := parts[0]
	if len(parts) == 1 {
		if request.Method != http.MethodGet {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		issue, err := server.app.GetIssue(issueID)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(writer, http.StatusOK, issue)
		return
	}
	switch parts[1] {
	case "comments":
		server.handleComment(writer, request, issueID)
	case "assignment":
		server.handleAssignment(writer, request, issueID)
	case "status":
		server.handleStatus(writer, request, issueID)
	case "attachments":
		if len(parts) == 2 {
			server.handleAttachmentUpload(writer, request, issueID)
			return
		}
		if len(parts) == 3 {
			server.handleAttachmentDownload(writer, request, issueID, parts[2])
			return
		}
		http.NotFound(writer, request)
	default:
		http.NotFound(writer, request)
	}
}

func (server *Server) handleComment(writer http.ResponseWriter, request *http.Request, issueID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	actor, err := requestActor(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 16*1024)
	var payload struct {
		Comment string `json:"comment"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	issue, err := server.app.AddComment(actor, issueID, payload.Comment)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, issue)
}

func (server *Server) handleAssignment(writer http.ResponseWriter, request *http.Request, issueID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	actor, err := requestActor(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 8*1024)
	var payload struct {
		Assignee string `json:"assignee"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	issue, err := server.app.AssignIssue(actor, issueID, payload.Assignee)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, issue)
}

func (server *Server) handleStatus(writer http.ResponseWriter, request *http.Request, issueID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	actor, err := requestActor(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 8*1024)
	var payload struct {
		Status string `json:"status"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	issue, err := server.app.ChangeStatus(actor, issueID, payload.Status)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, issue)
}

func (server *Server) handleAttachmentUpload(writer http.ResponseWriter, request *http.Request, issueID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	actor, err := requestActor(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxAttachmentBytes+4096)
	if err := request.ParseMultipartForm(maxAttachmentBytes); err != nil {
		http.Error(writer, fmt.Sprintf("parse multipart form: %v", err), http.StatusBadRequest)
		return
	}
	if request.MultipartForm != nil {
		defer func() {
			if err := request.MultipartForm.RemoveAll(); err != nil {
			}
		}()
	}
	file, header, err := request.FormFile("attachment")
	if err != nil {
		http.Error(writer, fmt.Sprintf("read attachment: %v", err), http.StatusBadRequest)
		return
	}
	bytes, err := io.ReadAll(io.LimitReader(file, maxAttachmentBytes+1))
	if closeErr := file.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	if err != nil {
		http.Error(writer, fmt.Sprintf("read attachment body: %v", err), http.StatusBadRequest)
		return
	}
	if len(bytes) > maxAttachmentBytes {
		http.Error(writer, "attachment is too large", http.StatusBadRequest)
		return
	}
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(bytes)
	}
	issue, err := server.app.AddAttachment(actor, issueID, header.Filename, contentType, bytes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, issue)
}

func (server *Server) handleAttachmentDownload(writer http.ResponseWriter, request *http.Request, issueID string, attachmentID string) {
	if _, err := requestActor(request); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	attachment, err := server.app.DownloadAttachment(issueID, attachmentID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", attachment.ContentType)
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", attachment.Name))
	if _, err := writer.Write(attachment.Bytes); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func requestActor(request *http.Request) (string, error) {
	actor := strings.TrimSpace(request.Header.Get("X-Bug-User"))
	if actor == "" {
		actor = strings.TrimSpace(request.URL.Query().Get("user"))
	}
	if actor == "" {
		return "", fmt.Errorf("missing user identity")
	}
	if _, err := validateIdentity(actor); err != nil {
		return "", err
	}
	return actor, nil
}

func decodeJSONBody(request *http.Request, target any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode json body: %w", err)
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, statusCode int, value any) {
	bytes, err := json.Marshal(value)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(statusCode)
	if _, err := writer.Write(append(bytes, '\n')); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
