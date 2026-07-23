package service

import (
	"fmt"
	"sort"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

type SignedKnowledgeEvidenceRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
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
		Sequence:       effectiveOriginSequence(event),
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
	record, ok, err := records.BuildSignedKnowledgeEvidenceRecord(identity, records.Event(event))
	return SignedKnowledgeEvidenceRecord(record), ok, err
}

func verifySignedKnowledgeEvidenceRecords(events []OperationalEvent, in []SignedKnowledgeEvidenceRecord) error {
	eventSlice := make([]records.Event, len(events))
	recordSlice := make([]records.SignedKnowledgeEvidenceRecord, len(in))
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	for i, record := range in {
		recordSlice[i] = records.SignedKnowledgeEvidenceRecord(record)
	}
	return records.VerifySignedKnowledgeEvidenceRecords(eventSlice, recordSlice)
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
