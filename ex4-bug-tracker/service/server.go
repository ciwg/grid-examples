package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/identity"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocol"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/web"
)

const maxAttachmentBytes = 8 << 20
const maxPromiseBytes = 64 << 10

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
	mux.HandleFunc("/promise.js", server.handlePromiseJS)
	mux.HandleFunc("/style.css", server.handleStyleCSS)
	mux.HandleFunc("/api/meta", server.handleMeta)
	mux.HandleFunc("/api/agents/enroll", server.handleEnrollAgent)
	mux.HandleFunc("/api/agents/enroll/prepare", server.handlePrepareEnrollment)
	mux.HandleFunc("/api/agents/enroll/finalize", server.handleFinalizeEnrollment)
	mux.HandleFunc("/api/promises/prepare", server.handlePreparePromise)
	mux.HandleFunc("/api/promises/finalize", server.handleFinalizePromise)
	mux.HandleFunc("/api/promises", server.handleSubmitPromise)
	mux.HandleFunc("/api/attachments", server.handleAttachmentObject)
	mux.HandleFunc("/api/issues", server.handleIssues)
	mux.HandleFunc("/api/issues/", server.handleIssue)
	return mux
}

func (server *Server) handleEnrollAgent(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isCBORRequest(request) {
		http.Error(writer, "Content-Type application/cbor is required", http.StatusUnsupportedMediaType)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxPromiseBytes)
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	var payload struct {
		Enrollment identity.Enrollment      `cbor:"enrollment"`
		Proof      identity.EnrollmentProof `cbor:"proof"`
	}
	if err := protocol.Unmarshal(bytes, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := server.app.EnrollAgent(payload.Enrollment, payload.Proof); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusCreated, payload.Enrollment)
}

func (server *Server) handlePrepareEnrollment(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxPromiseBytes)
	var payload struct {
		PublicKey []byte `json:"public_key"`
		Role      string `json:"role"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	enrollment, signableBytes, err := server.app.PrepareEnrollment(payload.PublicKey, payload.Role)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{"enrollment": enrollment, "signable_bytes": signableBytes})
}

func (server *Server) handleFinalizeEnrollment(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxPromiseBytes)
	var payload struct {
		Enrollment identity.Enrollment      `json:"enrollment"`
		Proof      identity.EnrollmentProof `json:"proof"`
	}
	if err := decodeJSONBody(request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := server.app.EnrollAgent(payload.Enrollment, payload.Proof); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusCreated, payload.Enrollment)
}

func (server *Server) handleSubmitPromise(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isCBORRequest(request) {
		http.Error(writer, "Content-Type application/cbor is required", http.StatusUnsupportedMediaType)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxPromiseBytes)
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	issue, err := server.app.SubmitPromise(bytes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, issue)
}

func (server *Server) handlePreparePromise(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxPromiseBytes)
	var draft PromiseDraft
	if err := decodeJSONBody(request, &draft); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	prepared, err := server.app.PreparePromise(draft)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, prepared)
}

func (server *Server) handleFinalizePromise(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxPromiseBytes)
	var proof PromiseProof
	if err := decodeJSONBody(request, &proof); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	wire, err := server.app.FinalizePromise(proof)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, map[string][]byte{"envelope": wire})
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

func (server *Server) handlePromiseJS(writer http.ResponseWriter, request *http.Request) {
	writeNoStoreHeaders(writer)
	writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	if _, err := writer.Write(web.MustRead("promise.js")); err != nil {
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
	default:
		http.Error(writer, "signed promise submission is required", http.StatusMethodNotAllowed)
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
	case "attachments":
		if len(parts) == 2 {
			http.Error(writer, "signed attachment-reference promise is required", http.StatusMethodNotAllowed)
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

// handleAttachmentObject accepts opaque bytes only; a signed attachment
// reference is still required before the object appears in an issue timeline.
func (server *Server) handleAttachmentObject(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maxAttachmentBytes)
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	cid, err := server.app.StoreAttachmentObject(bytes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]any{"cid": cid, "size": len(bytes), "content_type": request.Header.Get("Content-Type")})
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

func isCBORRequest(request *http.Request) bool {
	contentType := request.Header.Get("Content-Type")
	return strings.EqualFold(strings.TrimSpace(strings.Split(contentType, ";")[0]), "application/cbor")
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
