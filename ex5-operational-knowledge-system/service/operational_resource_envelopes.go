package service

import (
	"fmt"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

type SignedOperationalResourceRecord = records.SignedOperationalResourceRecord

type operationalResourcePayload struct {
	EntityID  string   `cbor:"entity_id"`
	EventType string   `cbor:"event_type"`
	Sequence  uint64   `cbor:"sequence"`
	Timestamp string   `cbor:"timestamp"`
	Actor     string   `cbor:"actor"`
	Kind      string   `cbor:"kind"`
	Name      string   `cbor:"name"`
	Summary   string   `cbor:"summary,omitempty"`
	PlaceID   string   `cbor:"place_id,omitempty"`
	Tags      []string `cbor:"tags,omitempty"`
}

func operationalResourcePayloadForEvent(event OperationalEvent) (operationalResourcePayload, bool) {
	if event.Type != "resource_created" {
		return operationalResourcePayload{}, false
	}
	return operationalResourcePayload{
		EntityID:  event.EntityID,
		EventType: event.Type,
		Sequence:  effectiveOriginSequence(event),
		Timestamp: event.Timestamp,
		Actor:     event.Actor,
		Kind:      event.Kind,
		Name:      event.Name,
		Summary:   event.Summary,
		PlaceID:   event.PlaceID,
		Tags:      append([]string(nil), event.Tags...),
	}, true
}

// Intent: Freeze first-class operational resources as signed durable context
// so exchanged runs and links can resolve their resource references without
// falling back to unresolved local-only context. Source: DI-pivul
func buildSignedOperationalResourceRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedOperationalResourceRecord, bool, error) {
	return records.BuildSignedOperationalResourceRecord(identity, records.Event(event))
}

func verifySignedOperationalResourceRecords(events []OperationalEvent, in []SignedOperationalResourceRecord) error {
	return records.VerifySignedOperationalResourceRecords(events, in)
}

func compareOperationalResourcePayload(expected operationalResourcePayload, got operationalResourcePayload) error {
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
	case expected.Name != got.Name:
		return fmt.Errorf("name mismatch")
	case expected.Summary != got.Summary:
		return fmt.Errorf("summary mismatch")
	case expected.PlaceID != got.PlaceID:
		return fmt.Errorf("place_id mismatch")
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
