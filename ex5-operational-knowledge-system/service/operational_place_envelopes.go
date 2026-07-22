package service

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
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
	payload, ok := operationalPlacePayloadForEvent(event)
	if !ok {
		return SignedOperationalPlaceRecord{}, false, nil
	}
	payloadBytes, err := protocols.Marshal(payload)
	if err != nil {
		return SignedOperationalPlaceRecord{}, false, fmt.Errorf("marshal operational-place payload: %w", err)
	}
	envelope := protocols.NewEnvelope(protocols.OperationalPlaceProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return SignedOperationalPlaceRecord{}, false, fmt.Errorf("build signable operational-place envelope: %w", err)
	}
	proofBytes, err := identity.SignProof(signable)
	if err != nil {
		return SignedOperationalPlaceRecord{}, false, fmt.Errorf("sign operational-place envelope: %w", err)
	}
	envelope = protocols.NewEnvelope(protocols.OperationalPlaceProfile.CID, payloadBytes, proofBytes)
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return SignedOperationalPlaceRecord{}, false, fmt.Errorf("encode operational-place envelope: %w", err)
	}
	envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
	if err != nil {
		return SignedOperationalPlaceRecord{}, false, fmt.Errorf("cid operational-place envelope: %w", err)
	}
	return SignedOperationalPlaceRecord{
		Sequence:       event.Sequence,
		OriginPeerID:   effectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: effectiveOriginSequence(event),
		PlaceID:        event.EntityID,
		PCID:           protocols.OperationalPlaceProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func verifySignedOperationalPlaceRecords(events []OperationalEvent, records []SignedOperationalPlaceRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]operationalPlacePayload{}
	for _, event := range events {
		payload, ok := operationalPlacePayloadForEvent(event)
		if !ok {
			continue
		}
		expected[originEventKey(effectiveOriginPeerID(event, ""), effectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := strings.TrimSpace(record.OriginPeerID)
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[originEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if record.PCID != protocols.OperationalPlaceProfile.CID.String() {
			return fmt.Errorf("operational-place record %d uses unexpected pCID %q", record.Sequence, record.PCID)
		}
		envelopeBytes, err := base64.StdEncoding.DecodeString(record.EnvelopeBase64)
		if err != nil {
			return fmt.Errorf("decode operational-place record %d envelope: %w", record.Sequence, err)
		}
		envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
		if err != nil {
			return fmt.Errorf("cid operational-place record %d envelope: %w", record.Sequence, err)
		}
		if envelopeCID.String() != record.EnvelopeCID {
			return fmt.Errorf("operational-place record %d envelope cid mismatch", record.Sequence)
		}
		envelope, err := protocols.ParseEnvelope(envelopeBytes)
		if err != nil {
			return fmt.Errorf("parse operational-place record %d envelope: %w", record.Sequence, err)
		}
		if envelope.PCID.String() != protocols.OperationalPlaceProfile.CID.String() {
			return fmt.Errorf("operational-place record %d envelope pCID mismatch", record.Sequence)
		}
		signable, err := envelope.SignableBytes()
		if err != nil {
			return fmt.Errorf("build operational-place record %d signable bytes: %w", record.Sequence, err)
		}
		if err := VerifyRuntimeProof(signable, envelope.ProofBytes); err != nil {
			return fmt.Errorf("verify operational-place record %d proof: %w", record.Sequence, err)
		}
		var got operationalPlacePayload
		if err := protocols.Unmarshal(envelope.PayloadBytes, &got); err != nil {
			return fmt.Errorf("decode operational-place record %d payload: %w", record.Sequence, err)
		}
		if err := compareOperationalPlacePayload(payload, got); err != nil {
			return fmt.Errorf("operational-place record %d payload mismatch: %w", record.Sequence, err)
		}
	}
	return nil
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
