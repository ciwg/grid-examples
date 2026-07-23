package service

import (
	"fmt"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

type SignedOperationalPlaceRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	PlaceID        string `json:"place_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type operationalPlacePayload struct {
	EntityID  string   `cbor:"entity_id"`
	EventType string   `cbor:"event_type"`
	Sequence  uint64   `cbor:"sequence"`
	Timestamp string   `cbor:"timestamp"`
	Actor     string   `cbor:"actor"`
	Kind      string   `cbor:"kind"`
	Name      string   `cbor:"name"`
	Summary   string   `cbor:"summary,omitempty"`
	ParentID  string   `cbor:"parent_id,omitempty"`
	Tags      []string `cbor:"tags,omitempty"`
}

func operationalPlacePayloadForEvent(event OperationalEvent) (operationalPlacePayload, bool) {
	if event.Type != "place_created" {
		return operationalPlacePayload{}, false
	}
	return operationalPlacePayload{
		EntityID:  event.EntityID,
		EventType: event.Type,
		Sequence:  effectiveOriginSequence(event),
		Timestamp: event.Timestamp,
		Actor:     event.Actor,
		Kind:      event.Kind,
		Name:      event.Name,
		Summary:   event.Summary,
		ParentID:  event.ParentID,
		Tags:      append([]string(nil), event.Tags...),
	}, true
}

// Intent: Freeze first-class operational places as signed durable context so
// exchanged runs and links can resolve their place references without falling
// back to unresolved local-only context. Source: DI-pivul
func buildSignedOperationalPlaceRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedOperationalPlaceRecord, bool, error) {
	record, ok, err := records.BuildSignedOperationalPlaceRecord(identity, records.Event(event))
	return SignedOperationalPlaceRecord(record), ok, err
}

func verifySignedOperationalPlaceRecords(events []OperationalEvent, in []SignedOperationalPlaceRecord) error {
	eventSlice := make([]records.Event, len(events))
	recordSlice := make([]records.SignedOperationalPlaceRecord, len(in))
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	for i, record := range in {
		recordSlice[i] = records.SignedOperationalPlaceRecord(record)
	}
	return records.VerifySignedOperationalPlaceRecords(eventSlice, recordSlice)
}

func compareOperationalPlacePayload(expected operationalPlacePayload, got operationalPlacePayload) error {
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
	case expected.ParentID != got.ParentID:
		return fmt.Errorf("parent_id mismatch")
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
