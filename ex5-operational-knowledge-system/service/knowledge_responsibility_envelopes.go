package service

import (
	"fmt"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

type SignedKnowledgeResponsibilityRecord struct {
	Sequence         uint64 `json:"sequence"`
	OriginPeerID     string `json:"origin_peer_id"`
	OriginSequence   uint64 `json:"origin_sequence"`
	ResponsibilityID string `json:"responsibility_id"`
	PCID             string `json:"pcid"`
	EnvelopeCID      string `json:"envelope_cid"`
	EnvelopeBase64   string `json:"envelope_base64"`
	RecordedAt       string `json:"recorded_at"`
	Implementation   string `json:"implementation"`
}

type knowledgeResponsibilityPayload struct {
	EntityID  string   `cbor:"entity_id"`
	EventType string   `cbor:"event_type"`
	Sequence  uint64   `cbor:"sequence"`
	Timestamp string   `cbor:"timestamp"`
	Actor     string   `cbor:"actor"`
	Title     string   `cbor:"title"`
	Summary   string   `cbor:"summary,omitempty"`
	Team      string   `cbor:"team,omitempty"`
	RoleKeys  []string `cbor:"role_keys,omitempty"`
	Tags      []string `cbor:"tags,omitempty"`
}

func knowledgeResponsibilityPayloadForEvent(event OperationalEvent) (knowledgeResponsibilityPayload, bool) {
	if event.Type != "responsibility_created" {
		return knowledgeResponsibilityPayload{}, false
	}
	return knowledgeResponsibilityPayload{
		EntityID:  event.EntityID,
		EventType: event.Type,
		Sequence:  effectiveOriginSequence(event),
		Timestamp: event.Timestamp,
		Actor:     event.Actor,
		Title:     event.Title,
		Summary:   event.Summary,
		Team:      event.Team,
		RoleKeys:  append([]string(nil), event.RoleKeys...),
		Tags:      append([]string(nil), event.Tags...),
	}, true
}

// Intent: Freeze first-class responsibilities as the fifth ex5
// PromiseGrid-native family so durable role-bearing operational duties become
// signed and replay-verifiable without changing the current embodiment adapter
// surfaces. Source: DI-sarib
func buildSignedKnowledgeResponsibilityRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedKnowledgeResponsibilityRecord, bool, error) {
	record, ok, err := records.BuildSignedKnowledgeResponsibilityRecord(identity, records.Event(event))
	return SignedKnowledgeResponsibilityRecord(record), ok, err
}

func verifySignedKnowledgeResponsibilityRecords(events []OperationalEvent, in []SignedKnowledgeResponsibilityRecord) error {
	eventSlice := make([]records.Event, len(events))
	recordSlice := make([]records.SignedKnowledgeResponsibilityRecord, len(in))
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	for i, record := range in {
		recordSlice[i] = records.SignedKnowledgeResponsibilityRecord(record)
	}
	return records.VerifySignedKnowledgeResponsibilityRecords(eventSlice, recordSlice)
}

func compareKnowledgeResponsibilityPayload(expected knowledgeResponsibilityPayload, got knowledgeResponsibilityPayload) error {
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
	case expected.Title != got.Title:
		return fmt.Errorf("title mismatch")
	case expected.Summary != got.Summary:
		return fmt.Errorf("summary mismatch")
	case expected.Team != got.Team:
		return fmt.Errorf("team mismatch")
	}
	if len(expected.RoleKeys) != len(got.RoleKeys) {
		return fmt.Errorf("role_keys length mismatch")
	}
	for i := range expected.RoleKeys {
		if expected.RoleKeys[i] != got.RoleKeys[i] {
			return fmt.Errorf("role_keys mismatch")
		}
	}
	if len(expected.Tags) != len(got.Tags) {
		return fmt.Errorf("tags length mismatch")
	}
	for i := range expected.Tags {
		if expected.Tags[i] != got.Tags[i] {
			return fmt.Errorf("tags mismatch")
		}
	}
	return nil
}
