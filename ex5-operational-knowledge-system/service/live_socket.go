package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const liveSocketHeartbeatInterval = 2 * time.Second

type liveSocketRequest struct {
	Type          string `json:"type"`
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

type liveSocketEnvelope struct {
	Type    string        `json:"type"`
	State   LiveItemState `json:"state"`
	Message string        `json:"message,omitempty"`
}

// Intent: Keep websocket live drafting on the same item-level collaboration
// model as the existing `/live` route so transport changes do not invent a
// second draft semantics path. Source: DI-noruv
func (server *Server) handleItemLiveSocket(writer http.ResponseWriter, request *http.Request, itemID string) {
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
	updates, unsubscribe, err := server.app.SubscribeLiveItem(itemID)
	if err != nil {
		_ = socket.WriteJSON(map[string]any{
			"type":    "error",
			"message": err.Error(),
		})
		return
	}
	defer unsubscribe()

	connectedParticipantID := ""
	defer func() {
		if connectedParticipantID != "" {
			_ = server.app.RemoveLiveParticipant(itemID, connectedParticipantID)
		}
	}()

	initialState, err := server.app.LiveItemState(itemID)
	if err != nil {
		_ = socket.WriteJSON(map[string]any{
			"type":    "error",
			"message": err.Error(),
		})
		return
	}
	lastSnapshot, err := writeLiveStateSnapshot(socket, liveSocketEnvelope{Type: "live-state", State: initialState}, nil)
	if err != nil {
		return
	}

	readErr := make(chan error, 1)
	go func() {
		readErr <- server.readItemLiveSocket(socket, itemID, &connectedParticipantID)
	}()

	ticker := time.NewTicker(liveSocketHeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case err := <-readErr:
			if err == nil || err == io.EOF {
				return
			}
			return
		case <-updates:
			state, err := server.app.LiveItemState(itemID)
			if err != nil {
				return
			}
			lastSnapshot, err = writeLiveStateSnapshot(socket, liveSocketEnvelope{Type: "live-state", State: state}, lastSnapshot)
			if err != nil {
				return
			}
		case <-ticker.C:
			state, err := server.app.LiveItemState(itemID)
			if err != nil {
				return
			}
			lastSnapshot, err = writeLiveStateSnapshot(socket, liveSocketEnvelope{Type: "live-state", State: state}, lastSnapshot)
			if err != nil {
				return
			}
		}
	}
}

func (server *Server) readItemLiveSocket(socket *websocketConn, itemID string, connectedParticipantID *string) error {
	for {
		var payload liveSocketRequest
		if err := socket.ReadJSON(&payload); err != nil {
			return err
		}
		if payload.Type != "live-update" {
			if err := socket.WriteJSON(map[string]any{
				"type":    "error",
				"message": fmt.Sprintf("unknown live socket message type %q", payload.Type),
			}); err != nil {
				return err
			}
			continue
		}
		if connectedParticipantID != nil && payload.ParticipantID != "" {
			*connectedParticipantID = payload.ParticipantID
		}
		state, conflict, err := server.app.UpdateLiveItem(itemID, payload.ParticipantID, payload.DisplayName, payload.Color, payload.Cursor, payload.Head, payload.Typing, payload.BaseVersion, payload.UpdateBody, payload.Body)
		if err != nil {
			if writeErr := socket.WriteJSON(map[string]any{
				"type":    "error",
				"message": err.Error(),
			}); writeErr != nil {
				return writeErr
			}
			continue
		}
		if conflict {
			if err := socket.WriteJSON(liveSocketEnvelope{Type: "live-conflict", State: state}); err != nil {
				return err
			}
		}
	}
}

func writeLiveStateSnapshot(socket *websocketConn, envelope liveSocketEnvelope, lastSnapshot []byte) ([]byte, error) {
	currentSnapshot, err := json.Marshal(envelope)
	if err != nil {
		return lastSnapshot, err
	}
	if string(currentSnapshot) == string(lastSnapshot) {
		return lastSnapshot, nil
	}
	if err := socket.WriteJSON(envelope); err != nil {
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
