package service

import (
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

type SignedKnowledgeEvidenceRecord struct {
	Sequence       uint64 `json:"sequence"`
	EvidenceID     string `json:"evidence_id"`
	RunID          string `json:"run_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type knowledgeEvidencePayload struct {
	EvidenceID     string            `cbor:"evidence_id"`
	RunID          string            `cbor:"run_id"`
	EventType      string            `cbor:"event_type"`
	Sequence       uint64            `cbor:"sequence"`
	Timestamp      string            `cbor:"timestamp"`
	Actor          string            `cbor:"actor"`
	Summary        string            `cbor:"summary"`
	Facts          map[string]string `cbor:"facts,omitempty"`
	AttachmentName string            `cbor:"attachment_name,omitempty"`
	AttachmentPath string            `cbor:"attachment_path,omitempty"`
	AttachmentCID  string            `cbor:"attachment_cid,omitempty"`
	AttachmentSize int64             `cbor:"attachment_size,omitempty"`
}

func knowledgeEvidencePayloadForEvent(event OperationalEvent) (knowledgeEvidencePayload, bool) {
	if event.Type != "evidence_added" {
		return knowledgeEvidencePayload{}, false
	}
	return knowledgeEvidencePayload{
		EvidenceID:     event.EvidenceID,
		RunID:          event.EntityID,
		EventType:      event.Type,
		Sequence:       event.Sequence,
		Timestamp:      event.Timestamp,
		Actor:          event.Actor,
		Summary:        event.Summary,
		Facts:          cloneFacts(event.Facts),
		AttachmentName: event.AttachmentName,
		AttachmentPath: event.AttachmentPath,
		AttachmentCID:  event.AttachmentCID,
		AttachmentSize: event.AttachmentSize,
	}, true
}

// Intent: Freeze structured evidence plus attachment references as the third
// ex5 PromiseGrid-native family while leaving copied attachment bytes on the
// current runtime storage path. Source: DI-kavup; DI-ribof
func buildSignedKnowledgeEvidenceRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedKnowledgeEvidenceRecord, bool, error) {
	payload, ok := knowledgeEvidencePayloadForEvent(event)
	if !ok {
		return SignedKnowledgeEvidenceRecord{}, false, nil
	}
	payloadBytes, err := protocols.Marshal(payload)
	if err != nil {
		return SignedKnowledgeEvidenceRecord{}, false, fmt.Errorf("marshal evidence payload: %w", err)
	}
	envelope := protocols.NewEnvelope(protocols.KnowledgeEvidenceProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return SignedKnowledgeEvidenceRecord{}, false, fmt.Errorf("build signable evidence envelope: %w", err)
	}
	proofBytes, err := identity.SignProof(signable)
	if err != nil {
		return SignedKnowledgeEvidenceRecord{}, false, fmt.Errorf("sign evidence envelope: %w", err)
	}
	envelope = protocols.NewEnvelope(protocols.KnowledgeEvidenceProfile.CID, payloadBytes, proofBytes)
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return SignedKnowledgeEvidenceRecord{}, false, fmt.Errorf("encode evidence envelope: %w", err)
	}
	envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
	if err != nil {
		return SignedKnowledgeEvidenceRecord{}, false, fmt.Errorf("cid evidence envelope: %w", err)
	}
	return SignedKnowledgeEvidenceRecord{
		Sequence:       event.Sequence,
		EvidenceID:     event.EvidenceID,
		RunID:          event.EntityID,
		PCID:           protocols.KnowledgeEvidenceProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func verifySignedKnowledgeEvidenceRecords(events []OperationalEvent, records []SignedKnowledgeEvidenceRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[uint64]knowledgeEvidencePayload{}
	for _, event := range events {
		payload, ok := knowledgeEvidencePayloadForEvent(event)
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
		if record.PCID != protocols.KnowledgeEvidenceProfile.CID.String() {
			return fmt.Errorf("knowledge-evidence record %d uses unexpected pCID %q", record.Sequence, record.PCID)
		}
		envelopeBytes, err := base64.StdEncoding.DecodeString(record.EnvelopeBase64)
		if err != nil {
			return fmt.Errorf("decode knowledge-evidence record %d envelope: %w", record.Sequence, err)
		}
		envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
		if err != nil {
			return fmt.Errorf("cid knowledge-evidence record %d envelope: %w", record.Sequence, err)
		}
		if envelopeCID.String() != record.EnvelopeCID {
			return fmt.Errorf("knowledge-evidence record %d envelope cid mismatch", record.Sequence)
		}
		envelope, err := protocols.ParseEnvelope(envelopeBytes)
		if err != nil {
			return fmt.Errorf("parse knowledge-evidence record %d envelope: %w", record.Sequence, err)
		}
		if envelope.PCID.String() != protocols.KnowledgeEvidenceProfile.CID.String() {
			return fmt.Errorf("knowledge-evidence record %d envelope pCID mismatch", record.Sequence)
		}
		signable, err := envelope.SignableBytes()
		if err != nil {
			return fmt.Errorf("build knowledge-evidence record %d signable bytes: %w", record.Sequence, err)
		}
		if err := VerifyRuntimeProof(signable, envelope.ProofBytes); err != nil {
			return fmt.Errorf("verify knowledge-evidence record %d proof: %w", record.Sequence, err)
		}
		var got knowledgeEvidencePayload
		if err := protocols.Unmarshal(envelope.PayloadBytes, &got); err != nil {
			return fmt.Errorf("decode knowledge-evidence record %d payload: %w", record.Sequence, err)
		}
		if err := compareKnowledgeEvidencePayload(payload, got); err != nil {
			return fmt.Errorf("knowledge-evidence record %d payload mismatch: %w", record.Sequence, err)
		}
	}
	return nil
}

func compareKnowledgeEvidencePayload(expected knowledgeEvidencePayload, got knowledgeEvidencePayload) error {
	switch {
	case expected.EvidenceID != got.EvidenceID:
		return fmt.Errorf("evidence_id mismatch")
	case expected.RunID != got.RunID:
		return fmt.Errorf("run_id mismatch")
	case expected.EventType != got.EventType:
		return fmt.Errorf("event_type mismatch")
	case expected.Sequence != got.Sequence:
		return fmt.Errorf("sequence mismatch")
	case expected.Timestamp != got.Timestamp:
		return fmt.Errorf("timestamp mismatch")
	case expected.Actor != got.Actor:
		return fmt.Errorf("actor mismatch")
	case expected.Summary != got.Summary:
		return fmt.Errorf("summary mismatch")
	case expected.AttachmentName != got.AttachmentName:
		return fmt.Errorf("attachment_name mismatch")
	case expected.AttachmentPath != got.AttachmentPath:
		return fmt.Errorf("attachment_path mismatch")
	case expected.AttachmentCID != got.AttachmentCID:
		return fmt.Errorf("attachment_cid mismatch")
	case expected.AttachmentSize != got.AttachmentSize:
		return fmt.Errorf("attachment_size mismatch")
	}
	if len(expected.Facts) != len(got.Facts) {
		return fmt.Errorf("facts length mismatch")
	}
	keys := make([]string, 0, len(expected.Facts))
	for key := range expected.Facts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if expected.Facts[key] != got.Facts[key] {
			return fmt.Errorf("fact mismatch for %q", key)
		}
	}
	return nil
}
