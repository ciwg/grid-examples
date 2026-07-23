package service

import (
	"fmt"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

type SignedKnowledgeLinkRecord = records.SignedKnowledgeLinkRecord

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
	return records.BuildSignedKnowledgeLinkRecord(identity, records.Event(event))
}

func verifySignedKnowledgeLinkRecords(events []OperationalEvent, in []SignedKnowledgeLinkRecord) error {
	return records.VerifySignedKnowledgeLinkRecords(events, in)
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
