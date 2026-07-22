package service

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
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
	payload, ok := operationalRunPayloadForEvent(event)
	if !ok {
		return SignedOperationalRunRecord{}, false, nil
	}
	payloadBytes, err := protocols.Marshal(payload)
	if err != nil {
		return SignedOperationalRunRecord{}, false, fmt.Errorf("marshal operational-run payload: %w", err)
	}
	envelope := protocols.NewEnvelope(protocols.OperationalRunProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return SignedOperationalRunRecord{}, false, fmt.Errorf("build signable operational-run envelope: %w", err)
	}
	proofBytes, err := identity.SignProof(signable)
	if err != nil {
		return SignedOperationalRunRecord{}, false, fmt.Errorf("sign operational-run envelope: %w", err)
	}
	envelope = protocols.NewEnvelope(protocols.OperationalRunProfile.CID, payloadBytes, proofBytes)
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return SignedOperationalRunRecord{}, false, fmt.Errorf("encode operational-run envelope: %w", err)
	}
	envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
	if err != nil {
		return SignedOperationalRunRecord{}, false, fmt.Errorf("cid operational-run envelope: %w", err)
	}
	return SignedOperationalRunRecord{
		Sequence:       event.Sequence,
		OriginPeerID:   effectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: effectiveOriginSequence(event),
		RunID:          event.EntityID,
		ItemID:         event.TargetID,
		PCID:           protocols.OperationalRunProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func verifySignedOperationalRunRecords(events []OperationalEvent, records []SignedOperationalRunRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]operationalRunPayload{}
	for _, event := range events {
		payload, ok := operationalRunPayloadForEvent(event)
		if !ok {
			continue
		}
		expected[originEventKey(effectiveOriginPeerID(event, ""), effectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := record.OriginPeerID
		if strings.TrimSpace(peerID) == "" {
			peerID = ""
		}
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[originEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if record.PCID != protocols.OperationalRunProfile.CID.String() {
			return fmt.Errorf("operational-run record %d uses unexpected pCID %q", record.Sequence, record.PCID)
		}
		envelopeBytes, err := base64.StdEncoding.DecodeString(record.EnvelopeBase64)
		if err != nil {
			return fmt.Errorf("decode operational-run record %d envelope: %w", record.Sequence, err)
		}
		envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
		if err != nil {
			return fmt.Errorf("cid operational-run record %d envelope: %w", record.Sequence, err)
		}
		if envelopeCID.String() != record.EnvelopeCID {
			return fmt.Errorf("operational-run record %d envelope cid mismatch", record.Sequence)
		}
		envelope, err := protocols.ParseEnvelope(envelopeBytes)
		if err != nil {
			return fmt.Errorf("parse operational-run record %d envelope: %w", record.Sequence, err)
		}
		if envelope.PCID.String() != protocols.OperationalRunProfile.CID.String() {
			return fmt.Errorf("operational-run record %d envelope pCID mismatch", record.Sequence)
		}
		signable, err := envelope.SignableBytes()
		if err != nil {
			return fmt.Errorf("build operational-run record %d signable bytes: %w", record.Sequence, err)
		}
		if err := VerifyRuntimeProof(signable, envelope.ProofBytes); err != nil {
			return fmt.Errorf("verify operational-run record %d proof: %w", record.Sequence, err)
		}
		var got operationalRunPayload
		if err := protocols.Unmarshal(envelope.PayloadBytes, &got); err != nil {
			return fmt.Errorf("decode operational-run record %d payload: %w", record.Sequence, err)
		}
		if err := compareOperationalRunPayload(payload, got); err != nil {
			return fmt.Errorf("operational-run record %d payload mismatch: %w", record.Sequence, err)
		}
	}
	return nil
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
