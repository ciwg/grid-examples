package service

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const embodimentSocketFilename = "embodiment.sock"

// Intent: Canonicalize the runtime root once before deriving local embodiment
// contract paths so metadata and listeners agree on one absolute location.
// Source: DI-rakuv
func CanonicalDataRoot(dataRoot string) (string, error) {
	return filepath.Abs(dataRoot)
}

func EmbodimentSocketPath(dataRoot string) string {
	return filepath.Join(dataRoot, embodimentSocketFilename)
}

type LocalEmbodimentServer struct {
	app        *App
	httpServer http.Handler
	socketPath string
	listener   net.Listener
}

func NewLocalEmbodimentServer(app *App, socketPath string) *LocalEmbodimentServer {
	resolvedSocketPath, err := filepath.Abs(socketPath)
	if err != nil {
		resolvedSocketPath = socketPath
	}
	return &LocalEmbodimentServer{
		app:        app,
		httpServer: NewServer(app).Handler(),
		socketPath: resolvedSocketPath,
	}
}

// Intent: Publish one direct local Unix-socket contract for terminal
// embodiments without splitting durable runtime ownership away from the
// existing ex5 process. Source: DI-favel
func (server *LocalEmbodimentServer) ListenAndServe() error {
	listener, err := ensureOwnedSocketListener(server.socketPath)
	if err != nil {
		return err
	}
	server.listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		go server.handleConn(conn)
	}
}

func (server *LocalEmbodimentServer) Close() error {
	if server.listener != nil {
		if err := server.listener.Close(); err != nil {
			return err
		}
	}
	if err := os.Remove(server.socketPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// Intent: Refuse to steal a live runtime's embodiment socket while still
// reclaiming clearly stale socket files left behind by earlier exits. Source:
// DI-fegom
func ensureOwnedSocketListener(socketPath string) (net.Listener, error) {
	if err := os.MkdirAll(filepath.Dir(socketPath), 0o755); err != nil {
		return nil, err
	}
	info, err := os.Lstat(socketPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return net.Listen("unix", socketPath)
	}
	if info.Mode()&os.ModeSocket == 0 {
		return nil, fmt.Errorf("local embodiment socket path exists and is not a socket: %s", socketPath)
	}
	conn, err := net.DialTimeout("unix", socketPath, 250*time.Millisecond)
	if err == nil {
		_ = conn.Close()
		return nil, fmt.Errorf("local embodiment socket already owned by an active runtime: %s", socketPath)
	}
	if err := os.Remove(socketPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return net.Listen("unix", socketPath)
}

func (server *LocalEmbodimentServer) handleConn(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	decoder := json.NewDecoder(bufio.NewReader(conn))
	var request LocalEmbodimentRequest
	if err := decoder.Decode(&request); err != nil {
		if writeErr := server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "error", Message: err.Error()}); writeErr != nil {
			return
		}
		return
	}
	switch request.Type {
	case "operation":
		server.handleOperation(conn, request)
	case "request":
		server.handleRequest(conn, request)
	case "live-open":
		server.handleLiveStream(conn, decoder, request)
	default:
		if err := server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "error", Message: fmt.Sprintf("unknown local socket request type %q", request.Type)}); err != nil {
			return
		}
	}
}

// Intent: Move the first terminal inspect/read slice onto typed local runtime
// operations so the direct socket contract stops tunneling selected workflows
// through route-shaped HTTP forwarding. Source: DI-monuv
func (server *LocalEmbodimentServer) handleOperation(conn net.Conn, request LocalEmbodimentRequest) {
	switch request.Operation {
	case "inspect_item":
		item, err := server.app.GetKnowledgeItem(request.ItemID)
		if err != nil {
			_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "response", Status: http.StatusNotFound, Body: err.Error()})
			return
		}
		server.writeSocketProjection(conn, item)
	case "inspect_run":
		run, err := server.app.GetRun(request.RunID)
		if err != nil {
			_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "response", Status: http.StatusNotFound, Body: err.Error()})
			return
		}
		server.writeSocketProjection(conn, run)
	case "inspect_entity":
		server.handleInspectEntityOperation(conn, request)
	case "search":
		options := SearchOptions{}
		if request.SearchOptions != nil {
			options = *request.SearchOptions
		}
		server.writeSocketProjection(conn, server.app.SearchWithOptions(options))
	case "pending_review":
		server.writeSocketProjection(conn, server.app.PendingReview())
	case "problem_review":
		server.writeSocketProjection(conn, server.app.ProblemReview())
	default:
		_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{
			Type:   "response",
			Status: http.StatusBadRequest,
			Body:   fmt.Sprintf("unknown local socket operation %q", request.Operation),
		})
	}
}

func (server *LocalEmbodimentServer) handleInspectEntityOperation(conn net.Conn, request LocalEmbodimentRequest) {
	switch request.EntityType {
	case "item":
		item, err := server.app.GetKnowledgeItem(request.EntityID)
		if err != nil {
			_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "response", Status: http.StatusNotFound, Body: err.Error()})
			return
		}
		server.writeSocketProjection(conn, item)
	case "run":
		run, err := server.app.GetRun(request.EntityID)
		if err != nil {
			_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "response", Status: http.StatusNotFound, Body: err.Error()})
			return
		}
		server.writeSocketProjection(conn, run)
	case "place":
		place, err := server.app.GetPlace(request.EntityID)
		if err != nil {
			_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "response", Status: http.StatusNotFound, Body: err.Error()})
			return
		}
		server.writeSocketProjection(conn, place)
	case "resource":
		resource, err := server.app.GetResource(request.EntityID)
		if err != nil {
			_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "response", Status: http.StatusNotFound, Body: err.Error()})
			return
		}
		server.writeSocketProjection(conn, resource)
	case "responsibility":
		record, err := server.app.GetResponsibility(request.EntityID)
		if err != nil {
			_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "response", Status: http.StatusNotFound, Body: err.Error()})
			return
		}
		server.writeSocketProjection(conn, record)
	default:
		_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{
			Type:   "response",
			Status: http.StatusBadRequest,
			Body:   fmt.Sprintf("unsupported inspect entity type %q", request.EntityType),
		})
	}
}

// Intent: Keep the first direct local contract aligned with the existing app
// semantics by forwarding request/response calls through the current handler
// layer rather than inventing a second mutation model. Source: DI-favel
func (server *LocalEmbodimentServer) handleRequest(conn net.Conn, request LocalEmbodimentRequest) {
	body, err := requestBodyBytes(request)
	if err != nil {
		if writeErr := server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "error", Message: err.Error()}); writeErr != nil {
			return
		}
		return
	}
	httpRequest := httptest.NewRequest(request.Method, "http://local"+request.Path, bytes.NewReader(body))
	for key, value := range request.Headers {
		httpRequest.Header.Set(key, value)
	}
	recorder := httptest.NewRecorder()
	server.httpServer.ServeHTTP(recorder, httpRequest)
	response := recorder.Result()
	payload, err := readLocalEmbodimentResponseBody(response)
	if err != nil {
		if writeErr := server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "error", Message: err.Error()}); writeErr != nil {
			return
		}
		return
	}
	if err := server.writeSocketMessage(conn, LocalEmbodimentResponse{
		Type:    "response",
		Status:  response.StatusCode,
		Headers: map[string]any{"content_type": response.Header.Get("Content-Type")},
		Body:    string(payload),
	}); err != nil {
		return
	}
}

func readLocalEmbodimentResponseBody(response *http.Response) ([]byte, error) {
	defer func() {
		_ = response.Body.Close()
	}()
	return io.ReadAll(response.Body)
}

// Intent: Move Neovim live drafting onto the same direct local socket without
// inventing a second collaboration model separate from the existing item-level
// live state. Source: DI-favel
func (server *LocalEmbodimentServer) handleLiveStream(conn net.Conn, decoder *json.Decoder, request LocalEmbodimentRequest) {
	writer := &socketJSONWriter{encoder: json.NewEncoder(conn)}
	updates, unsubscribe, err := server.app.SubscribeLiveItem(request.ItemID)
	if err != nil {
		_ = writer.Write(LocalEmbodimentResponse{Type: "error", Message: err.Error()})
		return
	}
	defer unsubscribe()
	connectedParticipantID := strings.TrimSpace(request.ParticipantID)
	if connectedParticipantID != "" {
		initialState, stateErr := server.app.LiveItemState(request.ItemID)
		if stateErr != nil {
			_ = writer.Write(LocalEmbodimentResponse{Type: "error", Message: stateErr.Error()})
			return
		}
		_, _, err = server.app.UpdateLiveItem(request.ItemID, connectedParticipantID, request.DisplayName, request.Color, request.Cursor, request.Head, request.Typing, initialState.Version, false, "")
		if err != nil {
			_ = writer.Write(LocalEmbodimentResponse{Type: "error", Message: err.Error()})
			return
		}
	}
	defer func() {
		if connectedParticipantID != "" {
			_ = server.app.RemoveLiveParticipant(request.ItemID, connectedParticipantID)
		}
	}()
	initialState, err := server.app.LiveItemState(request.ItemID)
	if err != nil {
		_ = writer.Write(LocalEmbodimentResponse{Type: "error", Message: err.Error()})
		return
	}
	lastSnapshot, err := writeSocketLiveState(writer, LocalEmbodimentResponse{Type: "live-state", State: initialState}, nil)
	if err != nil {
		return
	}
	readDone := make(chan error, 1)
	go func() {
		readDone <- server.readLiveSocketMessages(decoder, writer, request.ItemID, &connectedParticipantID)
	}()
	for {
		select {
		case err := <-readDone:
			if err != nil {
				_ = writer.Write(LocalEmbodimentResponse{Type: "error", Message: err.Error()})
			}
			return
		case <-updates:
			state, err := server.app.LiveItemState(request.ItemID)
			if err != nil {
				return
			}
			lastSnapshot, err = writeSocketLiveState(writer, LocalEmbodimentResponse{Type: "live-state", State: state}, lastSnapshot)
			if err != nil {
				return
			}
		}
	}
}

func (server *LocalEmbodimentServer) readLiveSocketMessages(decoder *json.Decoder, writer *socketJSONWriter, itemID string, connectedParticipantID *string) error {
	for {
		var payload LocalEmbodimentRequest
		if err := decoder.Decode(&payload); err != nil {
			return nil
		}
		switch payload.Type {
		case "live-close":
			return nil
		case "live-update":
			if connectedParticipantID != nil && strings.TrimSpace(payload.ParticipantID) != "" {
				*connectedParticipantID = payload.ParticipantID
			}
			state, conflict, err := server.app.UpdateLiveItem(itemID, payload.ParticipantID, payload.DisplayName, payload.Color, payload.Cursor, payload.Head, payload.Typing, payload.BaseVersion, payload.UpdateBody, payload.Body)
			if err != nil {
				if writeErr := writer.Write(LocalEmbodimentResponse{Type: "error", Message: err.Error()}); writeErr != nil {
					return writeErr
				}
				continue
			}
			if conflict {
				if err := writer.Write(LocalEmbodimentResponse{Type: "live-conflict", State: state}); err != nil {
					return err
				}
			}
		default:
			if err := writer.Write(LocalEmbodimentResponse{Type: "error", Message: fmt.Sprintf("unknown local live message type %q", payload.Type)}); err != nil {
				return err
			}
		}
	}
}

func (server *LocalEmbodimentServer) writeSocketMessage(conn net.Conn, response LocalEmbodimentResponse) error {
	encoder := json.NewEncoder(conn)
	return encoder.Encode(response)
}

func (server *LocalEmbodimentServer) writeSocketProjection(conn net.Conn, projection any) {
	payload, err := json.Marshal(projection)
	if err != nil {
		_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{Type: "error", Message: err.Error()})
		return
	}
	_ = server.writeSocketMessage(conn, LocalEmbodimentResponse{
		Type:   "response",
		Status: http.StatusOK,
		Body:   string(payload),
	})
}

func requestBodyBytes(request LocalEmbodimentRequest) ([]byte, error) {
	if strings.TrimSpace(request.BodyBase64) != "" {
		return base64.StdEncoding.DecodeString(request.BodyBase64)
	}
	return []byte(request.Body), nil
}

type socketJSONWriter struct {
	encoder *json.Encoder
	mu      sync.Mutex
}

func (writer *socketJSONWriter) Write(response LocalEmbodimentResponse) error {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	return writer.encoder.Encode(response)
}

func writeSocketLiveState(writer *socketJSONWriter, response LocalEmbodimentResponse, lastSnapshot []byte) ([]byte, error) {
	currentSnapshot, err := json.Marshal(response)
	if err != nil {
		return lastSnapshot, err
	}
	if string(currentSnapshot) == string(lastSnapshot) {
		return lastSnapshot, nil
	}
	if err := writer.Write(response); err != nil {
		return lastSnapshot, err
	}
	return currentSnapshot, nil
}
