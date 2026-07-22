package service

import (
	"encoding/base64"
	"fmt"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

type SignedKnowledgeItemRecord struct {
	Sequence       uint64 `json:"sequence"`
	ItemID         string `json:"item_id"`
	EventType      string `json:"event_type"`
	Revision       int    `json:"revision"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type knowledgeItemPayload struct {
	EntityID          string   `cbor:"entity_id"`
	EventType         string   `cbor:"event_type"`
	Sequence          uint64   `cbor:"sequence"`
	Timestamp         string   `cbor:"timestamp"`
	Actor             string   `cbor:"actor"`
	Kind              string   `cbor:"kind,omitempty"`
	Status            string   `cbor:"status,omitempty"`
	Title             string   `cbor:"title,omitempty"`
	Summary           string   `cbor:"summary,omitempty"`
	Body              string   `cbor:"body,omitempty"`
	Tags              []string `cbor:"tags,omitempty"`
	ResponsibilityIDs []string `cbor:"responsibility_ids,omitempty"`
	Revision          int      `cbor:"revision,omitempty"`
	Notes             string   `cbor:"notes,omitempty"`
}

func knowledgeItemPayloadForEvent(event OperationalEvent) (knowledgeItemPayload, bool) {
	switch event.Type {
	case "knowledge_item_created", "revision_added", "knowledge_item_status_changed", "knowledge_item_superseded":
		return knowledgeItemPayload{
			EntityID:          event.EntityID,
			EventType:         event.Type,
			Sequence:          event.Sequence,
			Timestamp:         event.Timestamp,
			Actor:             event.Actor,
			Kind:              event.Kind,
			Status:            event.Status,
			Title:             event.Title,
			Summary:           event.Summary,
			Body:              event.Body,
			Tags:              append([]string(nil), event.Tags...),
			ResponsibilityIDs: append([]string(nil), event.ResponsibilityIDs...),
			Revision:          event.Revision,
			Notes:             event.Notes,
		}, true
	default:
		return knowledgeItemPayload{}, false
	}
}

// Intent: Start the ex5 PromiseGrid migration with one signed durable family:
// knowledge-item create/revision/lifecycle events become signed envelopes while
// the rest of the runtime still projects through the existing local model.
// Source: DI-mibor
func buildSignedKnowledgeItemRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedKnowledgeItemRecord, bool, error) {
	payload, ok := knowledgeItemPayloadForEvent(event)
	if !ok {
		return SignedKnowledgeItemRecord{}, false, nil
	}
	payloadBytes, err := protocols.Marshal(payload)
	if err != nil {
		return SignedKnowledgeItemRecord{}, false, fmt.Errorf("marshal item payload: %w", err)
	}
	envelope := protocols.NewEnvelope(protocols.KnowledgeItemProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return SignedKnowledgeItemRecord{}, false, fmt.Errorf("build signable item envelope: %w", err)
	}
	proofBytes, err := identity.SignProof(signable)
	if err != nil {
		return SignedKnowledgeItemRecord{}, false, fmt.Errorf("sign item envelope: %w", err)
	}
	envelope = protocols.NewEnvelope(protocols.KnowledgeItemProfile.CID, payloadBytes, proofBytes)
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return SignedKnowledgeItemRecord{}, false, fmt.Errorf("encode item envelope: %w", err)
	}
	envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
	if err != nil {
		return SignedKnowledgeItemRecord{}, false, fmt.Errorf("cid item envelope: %w", err)
	}
	return SignedKnowledgeItemRecord{
		Sequence:       event.Sequence,
		ItemID:         event.EntityID,
		EventType:      event.Type,
		Revision:       event.Revision,
		PCID:           protocols.KnowledgeItemProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func verifySignedKnowledgeItemRecords(events []OperationalEvent, records []SignedKnowledgeItemRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[uint64]knowledgeItemPayload{}
	for _, event := range events {
		payload, ok := knowledgeItemPayloadForEvent(event)
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
		if record.PCID != protocols.KnowledgeItemProfile.CID.String() {
			return fmt.Errorf("knowledge-item record %d uses unexpected pCID %q", record.Sequence, record.PCID)
		}
		envelopeBytes, err := base64.StdEncoding.DecodeString(record.EnvelopeBase64)
		if err != nil {
			return fmt.Errorf("decode knowledge-item record %d envelope: %w", record.Sequence, err)
		}
		envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
		if err != nil {
			return fmt.Errorf("cid knowledge-item record %d envelope: %w", record.Sequence, err)
		}
		if envelopeCID.String() != record.EnvelopeCID {
			return fmt.Errorf("knowledge-item record %d envelope cid mismatch", record.Sequence)
		}
		envelope, err := protocols.ParseEnvelope(envelopeBytes)
		if err != nil {
			return fmt.Errorf("parse knowledge-item record %d envelope: %w", record.Sequence, err)
		}
		if envelope.PCID.String() != protocols.KnowledgeItemProfile.CID.String() {
			return fmt.Errorf("knowledge-item record %d envelope pCID mismatch", record.Sequence)
		}
		signable, err := envelope.SignableBytes()
		if err != nil {
			return fmt.Errorf("build knowledge-item record %d signable bytes: %w", record.Sequence, err)
		}
		if err := VerifyRuntimeProof(signable, envelope.ProofBytes); err != nil {
			return fmt.Errorf("verify knowledge-item record %d proof: %w", record.Sequence, err)
		}
		var got knowledgeItemPayload
		if err := protocols.Unmarshal(envelope.PayloadBytes, &got); err != nil {
			return fmt.Errorf("decode knowledge-item record %d payload: %w", record.Sequence, err)
		}
		if err := compareKnowledgeItemPayload(payload, got); err != nil {
			return fmt.Errorf("knowledge-item record %d payload mismatch: %w", record.Sequence, err)
		}
	}
	return nil
}

func compareKnowledgeItemPayload(expected knowledgeItemPayload, got knowledgeItemPayload) error {
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
	case expected.Status != got.Status:
		return fmt.Errorf("status mismatch")
	case expected.Title != got.Title:
		return fmt.Errorf("title mismatch")
	case expected.Summary != got.Summary:
		return fmt.Errorf("summary mismatch")
	case expected.Body != got.Body:
		return fmt.Errorf("body mismatch")
	case expected.Revision != got.Revision:
		return fmt.Errorf("revision mismatch")
	case expected.Notes != got.Notes:
		return fmt.Errorf("notes mismatch")
	}
	if len(expected.Tags) != len(got.Tags) {
		return fmt.Errorf("tag length mismatch")
	}
	for i := range expected.Tags {
		if expected.Tags[i] != got.Tags[i] {
			return fmt.Errorf("tag mismatch")
		}
	}
	if len(expected.ResponsibilityIDs) != len(got.ResponsibilityIDs) {
		return fmt.Errorf("responsibility length mismatch")
	}
	for i := range expected.ResponsibilityIDs {
		if expected.ResponsibilityIDs[i] != got.ResponsibilityIDs[i] {
			return fmt.Errorf("responsibility mismatch")
		}
	}
	return nil
}
