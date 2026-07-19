package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/awareness"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/crdt"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/metadata"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocol"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/publish"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/store"
)

type TraceFeed struct {
	DocumentID string       `json:"document_id"`
	Entries    []TraceEntry `json:"entries"`
}

type TraceEntry struct {
	Offset          uint64         `json:"offset"`
	EnvelopeCID     string         `json:"envelope_cid"`
	Protocol        string         `json:"protocol"`
	PCID            string         `json:"pcid"`
	Kind            string         `json:"kind"`
	DocumentID      string         `json:"document_id"`
	ParticipantID   string         `json:"participant_id"`
	Author          string         `json:"author"`
	Embodiment      string         `json:"embodiment,omitempty"`
	Lamport         uint64         `json:"lamport,omitempty"`
	ReceivedAt      string         `json:"received_at"`
	Summary         string         `json:"summary"`
	EnvelopeBase64  string         `json:"envelope_base64"`
	PayloadBase64   string         `json:"payload_base64"`
	ProofAlgorithm  string         `json:"proof_algorithm"`
	ProofKeyID      string         `json:"proof_key_id"`
	DecodedPayload  map[string]any `json:"decoded_payload"`
	ChangeBytesLen  int            `json:"change_bytes_len,omitempty"`
	RecipientID     string         `json:"recipient_id,omitempty"`
	Cursor          int            `json:"cursor,omitempty"`
	Head            int            `json:"head,omitempty"`
	Typing          bool           `json:"typing,omitempty"`
	MetadataTitle   string         `json:"metadata_title,omitempty"`
	PublishTitle    string         `json:"publish_title,omitempty"`
	SourceKind      string         `json:"source_kind,omitempty"`
	SourceVersionID string         `json:"source_version_id,omitempty"`
}

// Intent: Expose real relay-observed PromiseGrid envelopes per document so the
// conference demo can show the signed messaging path directly instead of
// inventing browser-local faux events. Source: DI-holoz
func (app *App) Trace(documentID string, limit int) TraceFeed {
	app.mu.Lock()
	defer app.mu.Unlock()

	entries := app.log.All()
	maxEntries := clampFeedLimit(limit)
	traced := make([]TraceEntry, 0, maxEntries)
	for idx := len(entries) - 1; idx >= 0 && len(traced) < maxEntries; idx-- {
		entry, ok := app.traceEntryLocked(documentID, entries[idx])
		if !ok {
			continue
		}
		traced = append(traced, entry)
	}
	slices.Reverse(traced)
	return TraceFeed{
		DocumentID: documentID,
		Entries:    traced,
	}
}

func (app *App) traceEntryLocked(documentID string, entry store.Entry) (TraceEntry, bool) {
	envelopeBytes, err := base64.StdEncoding.DecodeString(entry.RawBase64)
	if err != nil {
		return TraceEntry{}, false
	}
	envelope, err := protocol.ParseEnvelope(envelopeBytes)
	if err != nil {
		return TraceEntry{}, false
	}

	switch entry.PCID {
	case app.documentPCID.String():
		var message crdt.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return TraceEntry{}, false
		}
		if message.DocumentID != documentID {
			return TraceEntry{}, false
		}
		return TraceEntry{
			Offset:         entry.Offset,
			EnvelopeCID:    entry.EnvelopeCID,
			Protocol:       "live-document",
			PCID:           entry.PCID,
			Kind:           message.Kind,
			DocumentID:     message.DocumentID,
			ParticipantID:  message.ParticipantID,
			RecipientID:    message.RecipientID,
			Author:         message.Author,
			Embodiment:     message.Embodiment,
			Lamport:        message.Lamport,
			ReceivedAt:     entry.ReceivedAt,
			Summary:        fmt.Sprintf("%s sent a CRDT change (%d bytes)", traceActorLabel(message.Embodiment, message.ParticipantID), len(message.ChangeBytes)),
			EnvelopeBase64: entry.RawBase64,
			PayloadBase64:  base64.StdEncoding.EncodeToString(envelope.PayloadBytes),
			ProofAlgorithm: envelope.Proof.Algorithm,
			ProofKeyID:     envelope.Proof.KeyID,
			DecodedPayload: decodedPayloadMap(envelope.PayloadBytes),
			ChangeBytesLen: len(message.ChangeBytes),
		}, true
	case app.awarenessPCID.String():
		var message awareness.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return TraceEntry{}, false
		}
		if message.DocumentID != documentID {
			return TraceEntry{}, false
		}
		return TraceEntry{
			Offset:         entry.Offset,
			EnvelopeCID:    entry.EnvelopeCID,
			Protocol:       "live-awareness",
			PCID:           entry.PCID,
			Kind:           message.Kind,
			DocumentID:     message.DocumentID,
			ParticipantID:  message.ParticipantID,
			Author:         message.Author,
			Embodiment:     message.Embodiment,
			Lamport:        message.Lamport,
			ReceivedAt:     entry.ReceivedAt,
			Summary:        fmt.Sprintf("%s moved to %d:%d", traceActorLabel(message.Embodiment, message.ParticipantID), message.Cursor, message.Head),
			EnvelopeBase64: entry.RawBase64,
			PayloadBase64:  base64.StdEncoding.EncodeToString(envelope.PayloadBytes),
			ProofAlgorithm: envelope.Proof.Algorithm,
			ProofKeyID:     envelope.Proof.KeyID,
			DecodedPayload: decodedPayloadMap(envelope.PayloadBytes),
			Cursor:         message.Cursor,
			Head:           message.Head,
			Typing:         message.Typing,
		}, true
	case app.metadataPCID.String():
		var message metadata.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return TraceEntry{}, false
		}
		if message.DocumentID != documentID {
			return TraceEntry{}, false
		}
		return TraceEntry{
			Offset:         entry.Offset,
			EnvelopeCID:    entry.EnvelopeCID,
			Protocol:       "document-metadata",
			PCID:           entry.PCID,
			Kind:           message.Kind,
			DocumentID:     message.DocumentID,
			ParticipantID:  message.ParticipantID,
			Author:         message.Author,
			Embodiment:     message.Embodiment,
			Lamport:        message.Lamport,
			ReceivedAt:     entry.ReceivedAt,
			Summary:        fmt.Sprintf("%s updated metadata", traceActorLabel(message.Embodiment, message.ParticipantID)),
			EnvelopeBase64: entry.RawBase64,
			PayloadBase64:  base64.StdEncoding.EncodeToString(envelope.PayloadBytes),
			ProofAlgorithm: envelope.Proof.Algorithm,
			ProofKeyID:     envelope.Proof.KeyID,
			DecodedPayload: decodedPayloadMap(envelope.PayloadBytes),
			MetadataTitle:  message.Title,
		}, true
	case app.publishPCID.String():
		var message publish.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return TraceEntry{}, false
		}
		if message.DocumentID != documentID {
			return TraceEntry{}, false
		}
		return TraceEntry{
			Offset:          entry.Offset,
			EnvelopeCID:     entry.EnvelopeCID,
			Protocol:        "publish-document",
			PCID:            entry.PCID,
			Kind:            message.Kind,
			DocumentID:      message.DocumentID,
			ParticipantID:   message.ParticipantID,
			Author:          message.Author,
			Embodiment:      message.Embodiment,
			Lamport:         message.Lamport,
			ReceivedAt:      entry.ReceivedAt,
			Summary:         fmt.Sprintf("%s published %q", traceActorLabel(message.Embodiment, message.ParticipantID), message.Title),
			EnvelopeBase64:  entry.RawBase64,
			PayloadBase64:   base64.StdEncoding.EncodeToString(envelope.PayloadBytes),
			ProofAlgorithm:  envelope.Proof.Algorithm,
			ProofKeyID:      envelope.Proof.KeyID,
			DecodedPayload:  decodedPayloadMap(envelope.PayloadBytes),
			PublishTitle:    message.Title,
			SourceKind:      message.SourceKind,
			SourceVersionID: message.SourceVersionID,
		}, true
	default:
		return TraceEntry{}, false
	}
}

func decodedPayloadMap(payloadBytes []byte) map[string]any {
	var decoded map[string]any
	if err := protocol.Unmarshal(payloadBytes, &decoded); err != nil {
		return map[string]any{"error": err.Error()}
	}
	jsonBytes, err := json.Marshal(decoded)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		return map[string]any{"error": err.Error()}
	}
	return decoded
}

func traceActorLabel(embodiment string, participantID string) string {
	if embodiment != "" {
		return fmt.Sprintf("%s %s", embodiment, participantID)
	}
	return participantID
}
