package service

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

type SignedKnowledgeResponsibilityRecord struct {
	Sequence         uint64 `json:"sequence"`
	OriginPeerID     string `json:"origin_peer_id"`
	OriginSequence   uint64 `json:"origin_sequence"`
	ResponsibilityID string `json:"responsibility_id"`
	PCID             string `json:"pcid"`
	EnvelopeCID      string `json:"envelope_cid"`
	EnvelopeBase64   string `json:"envelope_base64"`
	RecordedAt       string `json:"recorded_at"`
	Implementation   string `json:"implementation"`
}

type knowledgeResponsibilityPayload struct {
	EntityID  string   `cbor:"entity_id"`
	EventType string   `cbor:"event_type"`
	Sequence  uint64   `cbor:"sequence"`
	Timestamp string   `cbor:"timestamp"`
	Actor     string   `cbor:"actor"`
	Title     string   `cbor:"title"`
	Summary   string   `cbor:"summary,omitempty"`
	Team      string   `cbor:"team,omitempty"`
	RoleKeys  []string `cbor:"role_keys,omitempty"`
	Tags      []string `cbor:"tags,omitempty"`
}

func knowledgeResponsibilityPayloadForEvent(event OperationalEvent) (knowledgeResponsibilityPayload, bool) {
	if event.Type != "responsibility_created" {
		return knowledgeResponsibilityPayload{}, false
	}
	return knowledgeResponsibilityPayload{
		EntityID:  event.EntityID,
		EventType: event.Type,
		Sequence:  effectiveOriginSequence(event),
		Timestamp: event.Timestamp,
		Actor:     event.Actor,
		Title:     event.Title,
		Summary:   event.Summary,
		Team:      event.Team,
		RoleKeys:  append([]string(nil), event.RoleKeys...),
		Tags:      append([]string(nil), event.Tags...),
	}, true
}

// Intent: Freeze first-class responsibilities as the fifth ex5
// PromiseGrid-native family so durable role-bearing operational duties become
// signed and replay-verifiable without changing the current embodiment adapter
// surfaces. Source: DI-sarib
func buildSignedKnowledgeResponsibilityRecord(identity *RuntimeIdentity, event OperationalEvent) (SignedKnowledgeResponsibilityRecord, bool, error) {
	payload, ok := knowledgeResponsibilityPayloadForEvent(event)
	if !ok {
		return SignedKnowledgeResponsibilityRecord{}, false, nil
	}
	payloadBytes, err := protocols.Marshal(payload)
	if err != nil {
		return SignedKnowledgeResponsibilityRecord{}, false, fmt.Errorf("marshal responsibility payload: %w", err)
	}
	envelope := protocols.NewEnvelope(protocols.KnowledgeResponsibilityProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return SignedKnowledgeResponsibilityRecord{}, false, fmt.Errorf("build signable responsibility envelope: %w", err)
	}
	proofBytes, err := identity.SignProof(signable)
	if err != nil {
		return SignedKnowledgeResponsibilityRecord{}, false, fmt.Errorf("sign responsibility envelope: %w", err)
	}
	envelope = protocols.NewEnvelope(protocols.KnowledgeResponsibilityProfile.CID, payloadBytes, proofBytes)
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return SignedKnowledgeResponsibilityRecord{}, false, fmt.Errorf("encode responsibility envelope: %w", err)
	}
	envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
	if err != nil {
		return SignedKnowledgeResponsibilityRecord{}, false, fmt.Errorf("cid responsibility envelope: %w", err)
	}
	return SignedKnowledgeResponsibilityRecord{
		Sequence:         event.Sequence,
		OriginPeerID:     effectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence:   effectiveOriginSequence(event),
		ResponsibilityID: event.EntityID,
		PCID:             protocols.KnowledgeResponsibilityProfile.CID.String(),
		EnvelopeCID:      envelopeCID.String(),
		EnvelopeBase64:   base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:       event.Timestamp,
		Implementation:   "ex5-local-runtime",
	}, true, nil
}

func verifySignedKnowledgeResponsibilityRecords(events []OperationalEvent, records []SignedKnowledgeResponsibilityRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]knowledgeResponsibilityPayload{}
	for _, event := range events {
		payload, ok := knowledgeResponsibilityPayloadForEvent(event)
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
		if record.PCID != protocols.KnowledgeResponsibilityProfile.CID.String() {
			return fmt.Errorf("knowledge-responsibility record %d uses unexpected pCID %q", record.Sequence, record.PCID)
		}
		envelopeBytes, err := base64.StdEncoding.DecodeString(record.EnvelopeBase64)
		if err != nil {
			return fmt.Errorf("decode knowledge-responsibility record %d envelope: %w", record.Sequence, err)
		}
		envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
		if err != nil {
			return fmt.Errorf("cid knowledge-responsibility record %d envelope: %w", record.Sequence, err)
		}
		if envelopeCID.String() != record.EnvelopeCID {
			return fmt.Errorf("knowledge-responsibility record %d envelope cid mismatch", record.Sequence)
		}
		envelope, err := protocols.ParseEnvelope(envelopeBytes)
		if err != nil {
			return fmt.Errorf("parse knowledge-responsibility record %d envelope: %w", record.Sequence, err)
		}
		if envelope.PCID.String() != protocols.KnowledgeResponsibilityProfile.CID.String() {
			return fmt.Errorf("knowledge-responsibility record %d envelope pCID mismatch", record.Sequence)
		}
		signable, err := envelope.SignableBytes()
		if err != nil {
			return fmt.Errorf("build knowledge-responsibility record %d signable bytes: %w", record.Sequence, err)
		}
		if err := VerifyRuntimeProof(signable, envelope.ProofBytes); err != nil {
			return fmt.Errorf("verify knowledge-responsibility record %d proof: %w", record.Sequence, err)
		}
		var got knowledgeResponsibilityPayload
		if err := protocols.Unmarshal(envelope.PayloadBytes, &got); err != nil {
			return fmt.Errorf("decode knowledge-responsibility record %d payload: %w", record.Sequence, err)
		}
		if err := compareKnowledgeResponsibilityPayload(payload, got); err != nil {
			return fmt.Errorf("knowledge-responsibility record %d payload mismatch: %w", record.Sequence, err)
		}
	}
	return nil
}

func compareKnowledgeResponsibilityPayload(expected knowledgeResponsibilityPayload, got knowledgeResponsibilityPayload) error {
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
	case expected.Title != got.Title:
		return fmt.Errorf("title mismatch")
	case expected.Summary != got.Summary:
		return fmt.Errorf("summary mismatch")
	case expected.Team != got.Team:
		return fmt.Errorf("team mismatch")
	}
	if len(expected.RoleKeys) != len(got.RoleKeys) {
		return fmt.Errorf("role_keys length mismatch")
	}
	for i := range expected.RoleKeys {
		if expected.RoleKeys[i] != got.RoleKeys[i] {
			return fmt.Errorf("role_keys mismatch")
		}
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
