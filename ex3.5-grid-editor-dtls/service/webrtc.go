package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
)

const webRTCChannelSync = "sync"
const webRTCChannelAwareness = "awareness"

type webRTCOfferRequest struct {
	Channel string                    `json:"channel"`
	Offer   webrtc.SessionDescription `json:"offer"`
}

type webRTCOfferResponse struct {
	Answer webrtc.SessionDescription `json:"answer"`
}

// webRTCSession owns one browser-to-node WebRTC association. Intent: Browser
// JavaScript cannot expose raw DTLS/UDP sockets, so the local embodiment uses
// WebRTC DataChannels, whose reliable data transport is DTLS-backed, while
// application envelopes retain their existing signed meanings. Source:
// DI-higom; DI-pirov.
type webRTCSession struct {
	server        *Server
	peer          *webrtc.PeerConnection
	documentID    string
	since         uint64
	channel       string
	participantID string

	closed    chan struct{}
	closeOnce sync.Once
}

func (server *Server) handleWebRTCOffer(writer http.ResponseWriter, request *http.Request, documentID string) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload webRTCOfferRequest
	if err := decodeJSONBody(writer, request, &payload); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Channel != webRTCChannelSync && payload.Channel != webRTCChannelAwareness {
		http.Error(writer, "invalid webrtc channel", http.StatusBadRequest)
		return
	}
	protocolCID := server.app.documentPCID.String()
	if payload.Channel == webRTCChannelAwareness {
		protocolCID = server.app.awarenessPCID.String()
	}
	expectedParticipantID := ""
	if !requestIsLoopback(request) {
		// Intent: Keep remote WebRTC mutation admission equivalent to the
		// WebSocket and HTTP paths: a browser must present a scoped signed
		// capability before its DTLS-backed DataChannel is established. Source:
		// DI-povip.
		token := bearerToken(request)
		if token == "" {
			err := fmt.Errorf("remote webrtc mutation requires a bearer capability")
			if diagnosticErr := server.app.RecordAdmissionDiagnostic("webrtc", err.Error()); diagnosticErr != nil {
				err = fmt.Errorf("%w (record admission diagnostic: %v)", err, diagnosticErr)
			}
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		claims, err := verifyMutationCapability(token, "", documentID, protocolCID, "mutate", time.Now().UTC())
		if err != nil {
			if diagnosticErr := server.app.RecordAdmissionDiagnostic("webrtc", err.Error()); diagnosticErr != nil {
				err = fmt.Errorf("%w (record admission diagnostic: %v)", err, diagnosticErr)
			}
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		if claims.Audience == "" {
			err := fmt.Errorf("capability audience missing")
			if diagnosticErr := server.app.RecordAdmissionDiagnostic("webrtc", err.Error()); diagnosticErr != nil {
				err = fmt.Errorf("%w (record admission diagnostic: %v)", err, diagnosticErr)
			}
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		expectedParticipantID = claims.Audience
	}
	since := uint64(0)
	if raw := request.URL.Query().Get("since"); raw != "" {
		if _, err := fmt.Sscanf(raw, "%d", &since); err != nil {
			http.Error(writer, "invalid since", http.StatusBadRequest)
			return
		}
	}
	peer, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		http.Error(writer, fmt.Sprintf("create webrtc peer: %v", err), http.StatusInternalServerError)
		return
	}
	session := &webRTCSession{
		server:        server,
		peer:          peer,
		documentID:    documentID,
		since:         since,
		channel:       payload.Channel,
		participantID: expectedParticipantID,
		closed:        make(chan struct{}),
	}
	peer.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state == webrtc.PeerConnectionStateClosed || state == webrtc.PeerConnectionStateFailed || state == webrtc.PeerConnectionStateDisconnected {
			session.close()
		}
	})
	peer.OnDataChannel(session.handleDataChannel)
	if err := peer.SetRemoteDescription(payload.Offer); err != nil {
		session.close()
		http.Error(writer, fmt.Sprintf("set webrtc offer: %v", err), http.StatusBadRequest)
		return
	}
	answer, err := peer.CreateAnswer(nil)
	if err != nil {
		session.close()
		http.Error(writer, fmt.Sprintf("create webrtc answer: %v", err), http.StatusInternalServerError)
		return
	}
	gatherComplete := webrtc.GatheringCompletePromise(peer)
	if err := peer.SetLocalDescription(answer); err != nil {
		session.close()
		http.Error(writer, fmt.Sprintf("set webrtc answer: %v", err), http.StatusInternalServerError)
		return
	}
	<-gatherComplete
	writeJSON(writer, http.StatusOK, webRTCOfferResponse{Answer: *peer.LocalDescription()})
}

func (session *webRTCSession) handleDataChannel(channel *webrtc.DataChannel) {
	if channel.Label() != session.channel {
		session.close()
		return
	}
	channel.OnOpen(func() {
		switch session.channel {
		case webRTCChannelSync:
			go session.runSync(channel)
		case webRTCChannelAwareness:
			go session.runAwareness(channel)
		}
	})
	channel.OnMessage(func(message webrtc.DataChannelMessage) {
		if message.IsString {
			session.handleMessage(channel, message.Data)
		}
	})
	channel.OnClose(session.close)
}

func (session *webRTCSession) runSync(channel *webrtc.DataChannel) {
	updates, unsubscribe := session.server.subscribeSync(session.documentID)
	defer unsubscribe()
	offset, err := session.writeSyncFeed(channel, session.since)
	if err != nil {
		session.close()
		return
	}
	if err := session.writeJSON(channel, map[string]any{"type": "sync-ready", "document_id": session.documentID, "next_offset": offset}); err != nil {
		session.close()
		return
	}
	ticker := time.NewTicker(liveSocketFallbackInterval)
	defer ticker.Stop()
	for {
		select {
		case <-session.closed:
			return
		case <-updates:
			offset, err = session.writeSyncFeed(channel, offset)
		case <-ticker.C:
			offset, err = session.writeSyncFeed(channel, offset)
		}
		if err != nil {
			session.close()
			return
		}
	}
}

func (session *webRTCSession) runAwareness(channel *webrtc.DataChannel) {
	updates, unsubscribe := session.server.subscribeAwareness(session.documentID)
	defer unsubscribe()
	lastSnapshot, err := session.writeAwarenessSnapshot(channel, nil)
	if err != nil {
		session.close()
		return
	}
	ticker := time.NewTicker(liveSocketFallbackInterval)
	defer ticker.Stop()
	for {
		select {
		case <-session.closed:
			return
		case <-updates:
			lastSnapshot, err = session.writeAwarenessSnapshot(channel, lastSnapshot)
		case <-ticker.C:
			lastSnapshot, err = session.writeAwarenessSnapshot(channel, lastSnapshot)
		}
		if err != nil {
			session.close()
			return
		}
	}
}

func (session *webRTCSession) handleMessage(channel *webrtc.DataChannel, raw []byte) {
	switch session.channel {
	case webRTCChannelSync:
		var payload syncSocketRequest
		if err := json.Unmarshal(raw, &payload); err != nil {
			session.writeError(channel, err)
			return
		}
		if payload.Type != "post-sync" {
			session.writeError(channel, fmt.Errorf("unknown sync datachannel message type %q", payload.Type))
			return
		}
		if session.participantID != "" && payload.ParticipantID != session.participantID {
			session.writeError(channel, fmt.Errorf("capability audience %q does not match participant %q", session.participantID, payload.ParticipantID))
			return
		}
		record, err := session.server.app.PostSync(session.documentID, payload.ParticipantID, payload.RecipientID, payload.MessageBase64, payload.Embodiment, payload.TextBase64, payload.ReplicaBase64)
		if err != nil {
			session.writeError(channel, err)
			return
		}
		if err := session.writeJSON(channel, map[string]any{"type": "sync-posted", "record": record}); err != nil {
			session.close()
		}
	case webRTCChannelAwareness:
		var payload awarenessSocketRequest
		if err := json.Unmarshal(raw, &payload); err != nil {
			session.writeError(channel, err)
			return
		}
		if payload.Type != "post-awareness" {
			session.writeError(channel, fmt.Errorf("unknown awareness datachannel message type %q", payload.Type))
			return
		}
		if session.participantID != "" && payload.ParticipantID != session.participantID {
			session.writeError(channel, fmt.Errorf("capability audience %q does not match participant %q", session.participantID, payload.ParticipantID))
			return
		}
		if err := session.server.app.UpdateAwareness(session.documentID, payload.ParticipantID, payload.Cursor, payload.Head, payload.Typing, payload.DisplayName, payload.Color, payload.Embodiment); err != nil {
			session.writeError(channel, err)
		}
	}
}

func (session *webRTCSession) writeSyncFeed(channel *webrtc.DataChannel, offset uint64) (uint64, error) {
	for {
		feed := session.server.app.SyncFeed(session.documentID, offset, defaultFeedLimit)
		if len(feed.Messages) == 0 || feed.NextOffset <= offset {
			return offset, nil
		}
		if err := session.writeJSON(channel, syncSocketFeed{Type: "sync-feed", DocumentID: session.documentID, Messages: feed.Messages, NextOffset: feed.NextOffset}); err != nil {
			return offset, err
		}
		offset = feed.NextOffset
		if len(feed.Messages) < defaultFeedLimit {
			return offset, nil
		}
	}
}

func (session *webRTCSession) writeAwarenessSnapshot(channel *webrtc.DataChannel, previous []byte) ([]byte, error) {
	payload := awarenessSocketState{Type: "awareness-state", DocumentID: session.documentID, Awareness: session.server.app.AwarenessState(session.documentID)}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return previous, err
	}
	if string(encoded) == string(previous) {
		return previous, nil
	}
	if err := channel.SendText(string(encoded)); err != nil {
		return previous, err
	}
	return encoded, nil
}

func (session *webRTCSession) writeJSON(channel *webrtc.DataChannel, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return channel.SendText(string(encoded))
}

func (session *webRTCSession) writeError(channel *webrtc.DataChannel, err error) {
	if writeErr := session.writeJSON(channel, map[string]any{"type": "error", "message": err.Error()}); writeErr != nil {
		session.close()
	}
}

func (session *webRTCSession) close() {
	session.closeOnce.Do(func() {
		close(session.closed)
		if err := session.peer.Close(); err != nil {
			return
		}
	})
}
