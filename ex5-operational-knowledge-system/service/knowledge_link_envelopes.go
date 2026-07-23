package service

import (
	"fmt"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

type SignedKnowledgeLinkRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	LinkID         string `json:"link_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type knowledgeLinkPayload struct {
	EntityID  string `cbor:"entity_id"`
	EventType string `cbor:"event_type"`
	Sequence  uint64 `cbor:"sequence"`
	Timestamp string `cbor:"timestamp"`
	Actor     string `cbor:"actor"`
	FromType  string `cbor:"from_type"`
	FromID    string `cbor:"from_id"`
	ToType    string `cbor:"to_type"`
	ToID      string `cbor:"to_id"`
	Relation  string `cbor:"relation"`
	Notes     string `cbor:"notes,omitempty"`
}

func knowledgeLinkPayloadForEvent(event OperationalEvent) (knowledgeLinkPayload, bool) {
	if event.Type != "link_added" {
		return knowledgeLinkPayload{}, false
	}
	return knowledgeLinkPayload{
		EntityID:  event.EntityID,
		EventType: event.Type,
		Sequence:  effectiveOriginSequence(event),
		Timestamp: event.Timestamp,
		Actor:     event.Actor,
		FromType:  event.FromType,
		FromID:    event.FromID,
		ToType:    event.ToType,
		ToID:      event.ToID,
		Relation:  event.Relation,
		Notes:     event.Notes,
	}, true
}

// Intent: Freeze typed operational links as the fourth ex5 PromiseGrid-native
// family so the durable link graph becomes signed and replay-verifiable without
// changing the current embodiment adapter surfaces. Source: DI-votek
func buildSignedKnowledgeLinkRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedKnowledgeLinkRecord, bool, error) {
	record, ok, err := records.BuildSignedKnowledgeLinkRecord(identity, records.Event(event))
	return SignedKnowledgeLinkRecord(record), ok, err
}

func verifySignedKnowledgeLinkRecords(events []OperationalEvent, in []SignedKnowledgeLinkRecord) error {
	eventSlice := make([]records.Event, len(events))
	recordSlice := make([]records.SignedKnowledgeLinkRecord, len(in))
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	for i, record := range in {
		recordSlice[i] = records.SignedKnowledgeLinkRecord(record)
	}
	return records.VerifySignedKnowledgeLinkRecords(eventSlice, recordSlice)
}

func compareKnowledgeLinkPayload(expected knowledgeLinkPayload, got knowledgeLinkPayload) error {
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
	case expected.FromType != got.FromType:
		return fmt.Errorf("from_type mismatch")
	case expected.FromID != got.FromID:
		return fmt.Errorf("from_id mismatch")
	case expected.ToType != got.ToType:
		return fmt.Errorf("to_type mismatch")
	case expected.ToID != got.ToID:
		return fmt.Errorf("to_id mismatch")
	case expected.Relation != got.Relation:
		return fmt.Errorf("relation mismatch")
	case expected.Notes != got.Notes:
		return fmt.Errorf("notes mismatch")
	}
	return nil
}
