package service

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

type SignedKnowledgeApprovalRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	ApprovalID     string `json:"approval_id"`
	TargetType     string `json:"target_type"`
	TargetID       string `json:"target_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type knowledgeApprovalPayload struct {
	EntityID   string `cbor:"entity_id"`
	EventType  string `cbor:"event_type"`
	Sequence   uint64 `cbor:"sequence"`
	Timestamp  string `cbor:"timestamp"`
	Actor      string `cbor:"actor"`
	TargetType string `cbor:"target_type"`
	TargetID   string `cbor:"target_id"`
	Revision   int    `cbor:"revision,omitempty"`
	Role       string `cbor:"role"`
	Decision   string `cbor:"decision"`
	Notes      string `cbor:"notes,omitempty"`
}

func knowledgeApprovalPayloadForEvent(event OperationalEvent) (knowledgeApprovalPayload, bool) {
	if event.Type != "approval_recorded" {
		return knowledgeApprovalPayload{}, false
	}
	return knowledgeApprovalPayload{
		EntityID:   event.EntityID,
		EventType:  event.Type,
		Sequence:   effectiveOriginSequence(event),
		Timestamp:  event.Timestamp,
		Actor:      event.Actor,
		TargetType: event.TargetType,
		TargetID:   event.TargetID,
		Revision:   event.Revision,
		Role:       event.Role,
		Decision:   event.Decision,
		Notes:      event.Notes,
	}, true
}

// Intent: Freeze named-role review outcomes as the second ex5
// PromiseGrid-native family so both item and run approvals become signed
// durable artifacts under one approval contract. Source: DI-vosul
func buildSignedKnowledgeApprovalRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedKnowledgeApprovalRecord, bool, error) {
	payload, ok := knowledgeApprovalPayloadForEvent(event)
	if !ok {
		return SignedKnowledgeApprovalRecord{}, false, nil
	}
	payloadBytes, err := protocols.Marshal(payload)
	if err != nil {
		return SignedKnowledgeApprovalRecord{}, false, fmt.Errorf("marshal approval payload: %w", err)
	}
	envelope := protocols.NewEnvelope(protocols.KnowledgeApprovalProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return SignedKnowledgeApprovalRecord{}, false, fmt.Errorf("build signable approval envelope: %w", err)
	}
	proofBytes, err := identity.SignProof(signable)
	if err != nil {
		return SignedKnowledgeApprovalRecord{}, false, fmt.Errorf("sign approval envelope: %w", err)
	}
	envelope = protocols.NewEnvelope(protocols.KnowledgeApprovalProfile.CID, payloadBytes, proofBytes)
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return SignedKnowledgeApprovalRecord{}, false, fmt.Errorf("encode approval envelope: %w", err)
	}
	envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
	if err != nil {
		return SignedKnowledgeApprovalRecord{}, false, fmt.Errorf("cid approval envelope: %w", err)
	}
	return SignedKnowledgeApprovalRecord{
		Sequence:       event.Sequence,
		OriginPeerID:   effectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: effectiveOriginSequence(event),
		ApprovalID:     event.EntityID,
		TargetType:     event.TargetType,
		TargetID:       event.TargetID,
		PCID:           protocols.KnowledgeApprovalProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func verifySignedKnowledgeApprovalRecords(events []OperationalEvent, records []SignedKnowledgeApprovalRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]knowledgeApprovalPayload{}
	for _, event := range events {
		payload, ok := knowledgeApprovalPayloadForEvent(event)
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
		if record.PCID != protocols.KnowledgeApprovalProfile.CID.String() {
			return fmt.Errorf("knowledge-approval record %d uses unexpected pCID %q", record.Sequence, record.PCID)
		}
		envelopeBytes, err := base64.StdEncoding.DecodeString(record.EnvelopeBase64)
		if err != nil {
			return fmt.Errorf("decode knowledge-approval record %d envelope: %w", record.Sequence, err)
		}
		envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
		if err != nil {
			return fmt.Errorf("cid knowledge-approval record %d envelope: %w", record.Sequence, err)
		}
		if envelopeCID.String() != record.EnvelopeCID {
			return fmt.Errorf("knowledge-approval record %d envelope cid mismatch", record.Sequence)
		}
		envelope, err := protocols.ParseEnvelope(envelopeBytes)
		if err != nil {
			return fmt.Errorf("parse knowledge-approval record %d envelope: %w", record.Sequence, err)
		}
		if envelope.PCID.String() != protocols.KnowledgeApprovalProfile.CID.String() {
			return fmt.Errorf("knowledge-approval record %d envelope pCID mismatch", record.Sequence)
		}
		signable, err := envelope.SignableBytes()
		if err != nil {
			return fmt.Errorf("build knowledge-approval record %d signable bytes: %w", record.Sequence, err)
		}
		if err := VerifyRuntimeProof(signable, envelope.ProofBytes); err != nil {
			return fmt.Errorf("verify knowledge-approval record %d proof: %w", record.Sequence, err)
		}
		var got knowledgeApprovalPayload
		if err := protocols.Unmarshal(envelope.PayloadBytes, &got); err != nil {
			return fmt.Errorf("decode knowledge-approval record %d payload: %w", record.Sequence, err)
		}
		if err := compareKnowledgeApprovalPayload(payload, got); err != nil {
			return fmt.Errorf("knowledge-approval record %d payload mismatch: %w", record.Sequence, err)
		}
	}
	return nil
}

func compareKnowledgeApprovalPayload(expected knowledgeApprovalPayload, got knowledgeApprovalPayload) error {
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
	case expected.TargetType != got.TargetType:
		return fmt.Errorf("target_type mismatch")
	case expected.TargetID != got.TargetID:
		return fmt.Errorf("target_id mismatch")
	case expected.Revision != got.Revision:
		return fmt.Errorf("revision mismatch")
	case expected.Role != got.Role:
		return fmt.Errorf("role mismatch")
	case expected.Decision != got.Decision:
		return fmt.Errorf("decision mismatch")
	case expected.Notes != got.Notes:
		return fmt.Errorf("notes mismatch")
	}
	return nil
}
