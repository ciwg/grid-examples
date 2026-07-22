package service

import (
	"encoding/base64"
	"fmt"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

type SignedKnowledgeLinkRecord struct {
	Sequence       uint64 `json:"sequence"`
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
		Sequence:  event.Sequence,
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
	payload, ok := knowledgeLinkPayloadForEvent(event)
	if !ok {
		return SignedKnowledgeLinkRecord{}, false, nil
	}
	payloadBytes, err := protocols.Marshal(payload)
	if err != nil {
		return SignedKnowledgeLinkRecord{}, false, fmt.Errorf("marshal link payload: %w", err)
	}
	envelope := protocols.NewEnvelope(protocols.KnowledgeLinkProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return SignedKnowledgeLinkRecord{}, false, fmt.Errorf("build signable link envelope: %w", err)
	}
	proofBytes, err := identity.SignProof(signable)
	if err != nil {
		return SignedKnowledgeLinkRecord{}, false, fmt.Errorf("sign link envelope: %w", err)
	}
	envelope = protocols.NewEnvelope(protocols.KnowledgeLinkProfile.CID, payloadBytes, proofBytes)
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return SignedKnowledgeLinkRecord{}, false, fmt.Errorf("encode link envelope: %w", err)
	}
	envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
	if err != nil {
		return SignedKnowledgeLinkRecord{}, false, fmt.Errorf("cid link envelope: %w", err)
	}
	return SignedKnowledgeLinkRecord{
		Sequence:       event.Sequence,
		LinkID:         event.EntityID,
		PCID:           protocols.KnowledgeLinkProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func verifySignedKnowledgeLinkRecords(events []OperationalEvent, records []SignedKnowledgeLinkRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[uint64]knowledgeLinkPayload{}
	for _, event := range events {
		payload, ok := knowledgeLinkPayloadForEvent(event)
		if !ok {
			continue
		}
		expected[event.Sequence] = payload
	}
	for _, record := range records {
		payload, ok := expected[record.Sequence]
		if !ok {
			continue
		}
		if record.PCID != protocols.KnowledgeLinkProfile.CID.String() {
			return fmt.Errorf("knowledge-link record %d uses unexpected pCID %q", record.Sequence, record.PCID)
		}
		envelopeBytes, err := base64.StdEncoding.DecodeString(record.EnvelopeBase64)
		if err != nil {
			return fmt.Errorf("decode knowledge-link record %d envelope: %w", record.Sequence, err)
		}
		envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
		if err != nil {
			return fmt.Errorf("cid knowledge-link record %d envelope: %w", record.Sequence, err)
		}
		if envelopeCID.String() != record.EnvelopeCID {
			return fmt.Errorf("knowledge-link record %d envelope cid mismatch", record.Sequence)
		}
		envelope, err := protocols.ParseEnvelope(envelopeBytes)
		if err != nil {
			return fmt.Errorf("parse knowledge-link record %d envelope: %w", record.Sequence, err)
		}
		if envelope.PCID.String() != protocols.KnowledgeLinkProfile.CID.String() {
			return fmt.Errorf("knowledge-link record %d envelope pCID mismatch", record.Sequence)
		}
		signable, err := envelope.SignableBytes()
		if err != nil {
			return fmt.Errorf("build knowledge-link record %d signable bytes: %w", record.Sequence, err)
		}
		if err := VerifyRuntimeProof(signable, envelope.ProofBytes); err != nil {
			return fmt.Errorf("verify knowledge-link record %d proof: %w", record.Sequence, err)
		}
		var got knowledgeLinkPayload
		if err := protocols.Unmarshal(envelope.PayloadBytes, &got); err != nil {
			return fmt.Errorf("decode knowledge-link record %d payload: %w", record.Sequence, err)
		}
		if err := compareKnowledgeLinkPayload(payload, got); err != nil {
			return fmt.Errorf("knowledge-link record %d payload mismatch: %w", record.Sequence, err)
		}
	}
	return nil
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
