package service

import (
	"fmt"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

type SignedKnowledgeApprovalRecord = records.SignedKnowledgeApprovalRecord

type knowledgeApprovalPayload struct {
	EntityID   string `cbor:"entity_id"`
	EventType  string `cbor:"event_type"`
	Sequence   uint64 `cbor:"sequence"`
	Timestamp  string `cbor:"timestamp"`
	Actor      string `cbor:"actor"`
	TargetType string `cbor:"target_type"`
	TargetID   string `cbor:"target_id"`
	Revision   int    `cbor:"revision,omitempty"`
	Role       string `cbor:"role"`
	Decision   string `cbor:"decision"`
	Notes      string `cbor:"notes,omitempty"`
}

func knowledgeApprovalPayloadForEvent(event OperationalEvent) (knowledgeApprovalPayload, bool) {
	if event.Type != "approval_recorded" {
		return knowledgeApprovalPayload{}, false
	}
	return knowledgeApprovalPayload{
		EntityID:   event.EntityID,
		EventType:  event.Type,
		Sequence:   effectiveOriginSequence(event),
		Timestamp:  event.Timestamp,
		Actor:      event.Actor,
		TargetType: event.TargetType,
		TargetID:   event.TargetID,
		Revision:   event.Revision,
		Role:       event.Role,
		Decision:   event.Decision,
		Notes:      event.Notes,
	}, true
}

// Intent: Freeze named-role review outcomes as the second ex5
// PromiseGrid-native family so both item and run approvals become signed
// durable artifacts under one approval contract. Source: DI-vosul
func buildSignedKnowledgeApprovalRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedKnowledgeApprovalRecord, bool, error) {
	return records.BuildSignedKnowledgeApprovalRecord(identity, records.Event(event))
}

func verifySignedKnowledgeApprovalRecords(events []OperationalEvent, in []SignedKnowledgeApprovalRecord) error {
	return records.VerifySignedKnowledgeApprovalRecords(events, in)
}

func compareKnowledgeApprovalPayload(expected knowledgeApprovalPayload, got knowledgeApprovalPayload) error {
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
	case expected.TargetType != got.TargetType:
		return fmt.Errorf("target_type mismatch")
	case expected.TargetID != got.TargetID:
		return fmt.Errorf("target_id mismatch")
	case expected.Revision != got.Revision:
		return fmt.Errorf("revision mismatch")
	case expected.Role != got.Role:
		return fmt.Errorf("role mismatch")
	case expected.Decision != got.Decision:
		return fmt.Errorf("decision mismatch")
	case expected.Notes != got.Notes:
		return fmt.Errorf("notes mismatch")
	}
	return nil
}
