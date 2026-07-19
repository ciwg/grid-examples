package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/awareness"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/crdt"
)

const liveSocketFallbackInterval = 2 * time.Second

type syncSocketRequest struct {
	Type          string `json:"type"`
	ParticipantID string `json:"participant_id"`
	RecipientID   string `json:"recipient_id"`
	MessageBase64 string `json:"message_base64"`
	Embodiment    string `json:"embodiment"`
}

type socketAuthRequest struct {
	Type       string `json:"type"`
	Capability string `json:"capability"`
}

type syncSocketFeed struct {
	Type       string            `json:"type"`
	DocumentID string            `json:"document_id"`
	Messages   []crdt.SyncRecord `json:"messages"`
	NextOffset uint64            `json:"next_offset"`
}

type awarenessSocketRequest struct {
	Type          string `json:"type"`
	ParticipantID string `json:"participant_id"`
	Cursor        int    `json:"cursor"`
	Head          int    `json:"head"`
	Typing        bool   `json:"typing"`
	DisplayName   string `json:"display_name"`
	Color         string `json:"color"`
	Embodiment    string `json:"embodiment"`
}

type awarenessSocketState struct {
	Type       string                `json:"type"`
	DocumentID string                `json:"document_id"`
	Awareness  []awareness.PeerState `json:"awareness"`
}

func (server *Server) handleSyncSocket(writer http.ResponseWriter, request *http.Request, documentID string) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	since := uint64(0)
	if raw := request.URL.Query().Get("since"); raw != "" {
		if _, err := fmt.Sscanf(raw, "%d", &since); err != nil {
			http.Error(writer, "invalid since", http.StatusBadRequest)
			return
		}
	}
	socket, err := upgradeWebSocket(writer, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer closeSocket(socket)
	expectedParticipantID := ""
	if !requestIsLoopback(request) {
		participantID, err := server.authorizeSocket(socket, documentID, server.app.documentPCID.String())
		if err != nil {
			return
		}
		expectedParticipantID = participantID
	}

	readErr := make(chan error, 1)
	go func() {
		readErr <- server.readSyncSocket(socket, documentID, expectedParticipantID)
	}()
	updates, unsubscribe := server.subscribeSync(documentID)
	defer unsubscribe()

	// Intent: Move the browser live-document transport onto websocket while
	// keeping the relay's signed feed model and HTTP metadata/publish surfaces
	// unchanged. Source: DI-vubih
	offset, err := server.writeSyncFeed(socket, documentID, since)
	if err != nil {
		return
	}
	if err := socket.WriteJSON(map[string]any{
		"type":        "sync-ready",
		"document_id": documentID,
		"next_offset": offset,
	}); err != nil {
		return
	}

	ticker := time.NewTicker(liveSocketFallbackInterval)
	defer ticker.Stop()
	for {
		select {
		case err := <-readErr:
			if err == nil || err == io.EOF {
				return
			}
			return
		case <-updates:
			offset, err = server.writeSyncFeed(socket, documentID, offset)
			if err != nil {
				return
			}
		case <-ticker.C:
			offset, err = server.writeSyncFeed(socket, documentID, offset)
			if err != nil {
				return
			}
		}
	}
}

func (server *Server) readSyncSocket(socket *websocketConn, documentID string, expectedParticipantID string) error {
	for {
		var payload syncSocketRequest
		if err := socket.ReadJSON(&payload); err != nil {
			return err
		}
		if payload.Type != "post-sync" {
			if err := socket.WriteJSON(map[string]any{
				"type":    "error",
				"message": fmt.Sprintf("unknown sync socket message type %q", payload.Type),
			}); err != nil {
				return err
			}
			continue
		}
		if expectedParticipantID != "" && payload.ParticipantID != expectedParticipantID {
			if err := socket.WriteJSON(map[string]any{
				"type":    "error",
				"message": fmt.Sprintf("capability audience %q does not match participant %q", expectedParticipantID, payload.ParticipantID),
			}); err != nil {
				return err
			}
			continue
		}
		record, err := server.app.PostSync(documentID, payload.ParticipantID, payload.RecipientID, payload.MessageBase64, payload.Embodiment)
		if err != nil {
			if writeErr := socket.WriteJSON(map[string]any{
				"type":    "error",
				"message": err.Error(),
			}); writeErr != nil {
				return writeErr
			}
			continue
		}
		if err := socket.WriteJSON(map[string]any{
			"type":   "sync-posted",
			"record": record,
		}); err != nil {
			return err
		}
	}
}

func (server *Server) writeSyncFeed(socket *websocketConn, documentID string, offset uint64) (uint64, error) {
	for {
		feed := server.app.SyncFeed(documentID, offset, defaultFeedLimit)
		if len(feed.Messages) == 0 || feed.NextOffset <= offset {
			return offset, nil
		}
		if err := socket.WriteJSON(syncSocketFeed{
			Type:       "sync-feed",
			DocumentID: documentID,
			Messages:   feed.Messages,
			NextOffset: feed.NextOffset,
		}); err != nil {
			return offset, err
		}
		offset = feed.NextOffset
		if len(feed.Messages) < defaultFeedLimit {
			return offset, nil
		}
	}
}

func (server *Server) handleAwarenessSocket(writer http.ResponseWriter, request *http.Request, documentID string) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	socket, err := upgradeWebSocket(writer, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer closeSocket(socket)
	expectedParticipantID := ""
	if !requestIsLoopback(request) {
		participantID, err := server.authorizeSocket(socket, documentID, server.app.awarenessPCID.String())
		if err != nil {
			return
		}
		expectedParticipantID = participantID
	}

	readErr := make(chan error, 1)
	go func() {
		readErr <- server.readAwarenessSocket(socket, documentID, expectedParticipantID)
	}()
	updates, unsubscribe := server.subscribeAwareness(documentID)
	defer unsubscribe()

	// Intent: Keep live-awareness as its own websocket-fed stream so cursor and
	// presence updates do not get collapsed into the document-sync channel.
	// Source: DI-vubih
	lastSnapshot, err := server.writeAwarenessSnapshot(socket, documentID, nil)
	if err != nil {
		return
	}

	ticker := time.NewTicker(liveSocketFallbackInterval)
	defer ticker.Stop()
	for {
		select {
		case err := <-readErr:
			if err == nil || err == io.EOF {
				return
			}
			return
		case <-updates:
			lastSnapshot, err = server.writeAwarenessSnapshot(socket, documentID, lastSnapshot)
			if err != nil {
				return
			}
		case <-ticker.C:
			lastSnapshot, err = server.writeAwarenessSnapshot(socket, documentID, lastSnapshot)
			if err != nil {
				return
			}
		}
	}
}

func (server *Server) readAwarenessSocket(socket *websocketConn, documentID string, expectedParticipantID string) error {
	for {
		var payload awarenessSocketRequest
		if err := socket.ReadJSON(&payload); err != nil {
			return err
		}
		if payload.Type != "post-awareness" {
			if err := socket.WriteJSON(map[string]any{
				"type":    "error",
				"message": fmt.Sprintf("unknown awareness socket message type %q", payload.Type),
			}); err != nil {
				return err
			}
			continue
		}
		if expectedParticipantID != "" && payload.ParticipantID != expectedParticipantID {
			if err := socket.WriteJSON(map[string]any{
				"type":    "error",
				"message": fmt.Sprintf("capability audience %q does not match participant %q", expectedParticipantID, payload.ParticipantID),
			}); err != nil {
				return err
			}
			continue
		}
		if err := server.app.UpdateAwareness(documentID, payload.ParticipantID, payload.Cursor, payload.Head, payload.Typing, payload.DisplayName, payload.Color, payload.Embodiment); err != nil {
			if writeErr := socket.WriteJSON(map[string]any{
				"type":    "error",
				"message": err.Error(),
			}); writeErr != nil {
				return writeErr
			}
			continue
		}
	}
}

func (server *Server) writeAwarenessSnapshot(socket *websocketConn, documentID string, lastSnapshot []byte) ([]byte, error) {
	payload := awarenessSocketState{
		Type:       "awareness-state",
		DocumentID: documentID,
		Awareness:  server.app.AwarenessState(documentID),
	}
	currentSnapshot, err := json.Marshal(payload)
	if err != nil {
		return lastSnapshot, err
	}
	if string(currentSnapshot) == string(lastSnapshot) {
		return lastSnapshot, nil
	}
	if err := socket.WriteJSON(payload); err != nil {
		return lastSnapshot, err
	}
	return currentSnapshot, nil
}

func closeSocket(socket *websocketConn) {
	if socket == nil {
		return
	}
	if err := socket.Close(); err != nil {
	}
}

func (server *Server) authorizeSocket(socket *websocketConn, documentID string, protocolCID string) (string, error) {
	var auth socketAuthRequest
	if err := socket.ReadJSON(&auth); err != nil {
		if writeErr := socket.WriteJSON(map[string]any{
			"type":    "error",
			"message": "remote websocket auth required",
		}); writeErr != nil {
			return "", writeErr
		}
		return "", err
	}
	if auth.Type != "auth" || auth.Capability == "" {
		if err := socket.WriteJSON(map[string]any{
			"type":    "error",
			"message": "remote websocket auth required",
		}); err != nil {
			return "", err
		}
		return "", fmt.Errorf("remote websocket auth required")
	}
	claims, err := verifyMutationCapability(auth.Capability, "", documentID, protocolCID, "mutate", time.Now().UTC())
	if err != nil {
		if writeErr := socket.WriteJSON(map[string]any{
			"type":    "error",
			"message": err.Error(),
		}); writeErr != nil {
			return "", writeErr
		}
		return "", err
	}
	if claims.Audience == "" {
		if err := socket.WriteJSON(map[string]any{
			"type":    "error",
			"message": "capability audience missing",
		}); err != nil {
			return "", err
		}
		return "", fmt.Errorf("capability audience missing")
	}
	return claims.Audience, nil
}
