package service

import (
	"fmt"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

type SignedOperationalRunRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	RunID          string `json:"run_id"`
	ItemID         string `json:"item_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type operationalRunPayload struct {
	RunID             string   `cbor:"run_id"`
	ItemID            string   `cbor:"item_id"`
	EventType         string   `cbor:"event_type"`
	Sequence          uint64   `cbor:"sequence"`
	Timestamp         string   `cbor:"timestamp"`
	Actor             string   `cbor:"actor"`
	Kind              string   `cbor:"kind"`
	Revision          int      `cbor:"revision"`
	Outcome           string   `cbor:"outcome,omitempty"`
	Notes             string   `cbor:"notes,omitempty"`
	PlaceID           string   `cbor:"place_id,omitempty"`
	ResourceIDs       []string `cbor:"resource_ids,omitempty"`
	Machine           string   `cbor:"machine,omitempty"`
	Location          string   `cbor:"location,omitempty"`
	ResponsibilityIDs []string `cbor:"responsibility_ids,omitempty"`
}

func operationalRunPayloadForEvent(event OperationalEvent) (operationalRunPayload, bool) {
	if event.Type != "run_recorded" {
		return operationalRunPayload{}, false
	}
	return operationalRunPayload{
		RunID:             event.EntityID,
		ItemID:            event.TargetID,
		EventType:         event.Type,
		Sequence:          effectiveOriginSequence(event),
		Timestamp:         event.Timestamp,
		Actor:             event.Actor,
		Kind:              event.Kind,
		Revision:          event.Revision,
		Outcome:           event.Outcome,
		Notes:             event.Notes,
		PlaceID:           event.PlaceID,
		ResourceIDs:       append([]string(nil), event.ResourceIDs...),
		Machine:           event.Machine,
		Location:          event.Location,
		ResponsibilityIDs: append([]string(nil), event.ResponsibilityIDs...),
	}, true
}

// Intent: Freeze performed execution records as the sixth ex5
// PromiseGrid-native family so evidence can anchor to a signed operational run
// contract instead of a compatibility-only local event. Source: DI-vamok
func buildSignedOperationalRunRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedOperationalRunRecord, bool, error) {
	record, ok, err := records.BuildSignedOperationalRunRecord(identity, records.Event(event))
	return SignedOperationalRunRecord(record), ok, err
}

func verifySignedOperationalRunRecords(events []OperationalEvent, in []SignedOperationalRunRecord) error {
	eventSlice := make([]records.Event, len(events))
	recordSlice := make([]records.SignedOperationalRunRecord, len(in))
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	for i, record := range in {
		recordSlice[i] = records.SignedOperationalRunRecord(record)
	}
	return records.VerifySignedOperationalRunRecords(eventSlice, recordSlice)
}

func compareOperationalRunPayload(expected operationalRunPayload, got operationalRunPayload) error {
	switch {
	case expected.RunID != got.RunID:
		return fmt.Errorf("run_id mismatch")
	case expected.ItemID != got.ItemID:
		return fmt.Errorf("item_id mismatch")
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
	case expected.Revision != got.Revision:
		return fmt.Errorf("revision mismatch")
	case expected.Outcome != got.Outcome:
		return fmt.Errorf("outcome mismatch")
	case expected.Notes != got.Notes:
		return fmt.Errorf("notes mismatch")
	case expected.PlaceID != got.PlaceID:
		return fmt.Errorf("place_id mismatch")
	case expected.Machine != got.Machine:
		return fmt.Errorf("machine mismatch")
	case expected.Location != got.Location:
		return fmt.Errorf("location mismatch")
	}
	if len(expected.ResourceIDs) != len(got.ResourceIDs) {
		return fmt.Errorf("resource_ids length mismatch")
	}
	for i := range expected.ResourceIDs {
		if expected.ResourceIDs[i] != got.ResourceIDs[i] {
			return fmt.Errorf("resource_ids mismatch")
		}
	}
	if len(expected.ResponsibilityIDs) != len(got.ResponsibilityIDs) {
		return fmt.Errorf("responsibility_ids length mismatch")
	}
	for i := range expected.ResponsibilityIDs {
		if expected.ResponsibilityIDs[i] != got.ResponsibilityIDs[i] {
			return fmt.Errorf("responsibility_ids mismatch")
		}
	}
	return nil
}
