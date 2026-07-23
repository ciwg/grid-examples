package service

import (
	"fmt"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

type SignedKnowledgeItemRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	ItemID         string `json:"item_id"`
	EventType      string `json:"event_type"`
	Revision       int    `json:"revision"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type knowledgeItemPayload struct {
	EntityID          string   `cbor:"entity_id"`
	EventType         string   `cbor:"event_type"`
	Sequence          uint64   `cbor:"sequence"`
	Timestamp         string   `cbor:"timestamp"`
	Actor             string   `cbor:"actor"`
	Kind              string   `cbor:"kind,omitempty"`
	Status            string   `cbor:"status,omitempty"`
	Title             string   `cbor:"title,omitempty"`
	Summary           string   `cbor:"summary,omitempty"`
	Body              string   `cbor:"body,omitempty"`
	Tags              []string `cbor:"tags,omitempty"`
	ResponsibilityIDs []string `cbor:"responsibility_ids,omitempty"`
	Revision          int      `cbor:"revision,omitempty"`
	Notes             string   `cbor:"notes,omitempty"`
}

func knowledgeItemPayloadForEvent(event OperationalEvent) (knowledgeItemPayload, bool) {
	switch event.Type {
	case "knowledge_item_created", "revision_added", "knowledge_item_status_changed", "knowledge_item_superseded":
		return knowledgeItemPayload{
			EntityID:          event.EntityID,
			EventType:         event.Type,
			Sequence:          effectiveOriginSequence(event),
			Timestamp:         event.Timestamp,
			Actor:             event.Actor,
			Kind:              event.Kind,
			Status:            event.Status,
			Title:             event.Title,
			Summary:           event.Summary,
			Body:              event.Body,
			Tags:              append([]string(nil), event.Tags...),
			ResponsibilityIDs: append([]string(nil), event.ResponsibilityIDs...),
			Revision:          event.Revision,
			Notes:             event.Notes,
		}, true
	default:
		return knowledgeItemPayload{}, false
	}
}

// Intent: Start the ex5 PromiseGrid migration with one signed durable family:
// knowledge-item create/revision/lifecycle events become signed envelopes while
// the rest of the runtime still projects through the existing local model.
// Source: DI-mibor
func buildSignedKnowledgeItemRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedKnowledgeItemRecord, bool, error) {
	record, ok, err := records.BuildSignedKnowledgeItemRecord(identity, records.Event(event))
	return SignedKnowledgeItemRecord(record), ok, err
}

func verifySignedKnowledgeItemRecords(events []OperationalEvent, in []SignedKnowledgeItemRecord) error {
	eventSlice := make([]records.Event, len(events))
	recordSlice := make([]records.SignedKnowledgeItemRecord, len(in))
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	for i, record := range in {
		recordSlice[i] = records.SignedKnowledgeItemRecord(record)
	}
	return records.VerifySignedKnowledgeItemRecords(eventSlice, recordSlice)
}

func compareKnowledgeItemPayload(expected knowledgeItemPayload, got knowledgeItemPayload) error {
	switch {
	case expected.EntityID != got.EntityID:
		return fmt.Errorf("entity_id mismatch")
	case expected.EventType != got.EventType:
		return fmt.Errorf("event_type mismatch")
	case expected.Sequence != got.Sequence:
		return fmt.Errorf("sequence mismatch")
	case expected.Timestamp != got.Timestamp:
		return fmt.Errorf("timestamp mismatch")
	case expected.Actor != got.Actor:
		return fmt.Errorf("actor mismatch")
	case expected.Kind != got.Kind:
		return fmt.Errorf("kind mismatch")
	case expected.Status != got.Status:
		return fmt.Errorf("status mismatch")
	case expected.Title != got.Title:
		return fmt.Errorf("title mismatch")
	case expected.Summary != got.Summary:
		return fmt.Errorf("summary mismatch")
	case expected.Body != got.Body:
		return fmt.Errorf("body mismatch")
	case expected.Revision != got.Revision:
		return fmt.Errorf("revision mismatch")
	case expected.Notes != got.Notes:
		return fmt.Errorf("notes mismatch")
	}
	if len(expected.Tags) != len(got.Tags) {
		return fmt.Errorf("tag length mismatch")
	}
	for i := range expected.Tags {
		if expected.Tags[i] != got.Tags[i] {
			return fmt.Errorf("tag mismatch")
		}
	}
	if len(expected.ResponsibilityIDs) != len(got.ResponsibilityIDs) {
		return fmt.Errorf("responsibility length mismatch")
	}
	for i := range expected.ResponsibilityIDs {
		if expected.ResponsibilityIDs[i] != got.ResponsibilityIDs[i] {
			return fmt.Errorf("responsibility mismatch")
		}
	}
	return nil
}
