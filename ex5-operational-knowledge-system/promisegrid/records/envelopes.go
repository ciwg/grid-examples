package records

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

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

type operationalResourcePayload struct {
	EntityID  string   `cbor:"entity_id"`
	EventType string   `cbor:"event_type"`
	Sequence  uint64   `cbor:"sequence"`
	Timestamp string   `cbor:"timestamp"`
	Actor     string   `cbor:"actor"`
	Kind      string   `cbor:"kind"`
	Name      string   `cbor:"name"`
	Summary   string   `cbor:"summary,omitempty"`
	PlaceID   string   `cbor:"place_id,omitempty"`
	Tags      []string `cbor:"tags,omitempty"`
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

// Intent: Keep frozen-family envelope build and verification logic in the
// reusable record substrate so durable PromiseGrid record truth is no longer
// owned only by ex5 service files. Source: DI-ragiv
func BuildSignedKnowledgeItemRecord(identity Signer, event Event) (SignedKnowledgeItemRecord, bool, error) {
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
		OriginPeerID:   EffectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: EffectiveOriginSequence(event),
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

func BuildSignedKnowledgeApprovalRecord(identity Signer, event Event) (SignedKnowledgeApprovalRecord, bool, error) {
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
		OriginPeerID:   EffectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: EffectiveOriginSequence(event),
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

func BuildSignedKnowledgeEvidenceRecord(identity Signer, event Event) (SignedKnowledgeEvidenceRecord, bool, error) {
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
		OriginPeerID:   EffectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: EffectiveOriginSequence(event),
		EvidenceID:     event.EvidenceID,
		RunID:          event.EntityID,
		PCID:           protocols.KnowledgeEvidenceProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func BuildSignedOperationalRunRecord(identity Signer, event Event) (SignedOperationalRunRecord, bool, error) {
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
		OriginPeerID:   EffectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: EffectiveOriginSequence(event),
		RunID:          event.EntityID,
		ItemID:         event.TargetID,
		PCID:           protocols.OperationalRunProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func BuildSignedOperationalPlaceRecord(identity Signer, event Event) (SignedOperationalPlaceRecord, bool, error) {
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
		OriginPeerID:   EffectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: EffectiveOriginSequence(event),
		PlaceID:        event.EntityID,
		PCID:           protocols.OperationalPlaceProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func BuildSignedOperationalResourceRecord(identity Signer, event Event) (SignedOperationalResourceRecord, bool, error) {
	payload, ok := operationalResourcePayloadForEvent(event)
	if !ok {
		return SignedOperationalResourceRecord{}, false, nil
	}
	payloadBytes, err := protocols.Marshal(payload)
	if err != nil {
		return SignedOperationalResourceRecord{}, false, fmt.Errorf("marshal operational-resource payload: %w", err)
	}
	envelope := protocols.NewEnvelope(protocols.OperationalResourceProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return SignedOperationalResourceRecord{}, false, fmt.Errorf("build signable operational-resource envelope: %w", err)
	}
	proofBytes, err := identity.SignProof(signable)
	if err != nil {
		return SignedOperationalResourceRecord{}, false, fmt.Errorf("sign operational-resource envelope: %w", err)
	}
	envelope = protocols.NewEnvelope(protocols.OperationalResourceProfile.CID, payloadBytes, proofBytes)
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return SignedOperationalResourceRecord{}, false, fmt.Errorf("encode operational-resource envelope: %w", err)
	}
	envelopeCID, err := protocols.CIDForBytes(envelopeBytes)
	if err != nil {
		return SignedOperationalResourceRecord{}, false, fmt.Errorf("cid operational-resource envelope: %w", err)
	}
	return SignedOperationalResourceRecord{
		Sequence:       event.Sequence,
		OriginPeerID:   EffectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: EffectiveOriginSequence(event),
		ResourceID:     event.EntityID,
		PCID:           protocols.OperationalResourceProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func BuildSignedKnowledgeLinkRecord(identity Signer, event Event) (SignedKnowledgeLinkRecord, bool, error) {
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
		OriginPeerID:   EffectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence: EffectiveOriginSequence(event),
		LinkID:         event.EntityID,
		PCID:           protocols.KnowledgeLinkProfile.CID.String(),
		EnvelopeCID:    envelopeCID.String(),
		EnvelopeBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:     event.Timestamp,
		Implementation: "ex5-local-runtime",
	}, true, nil
}

func BuildSignedKnowledgeResponsibilityRecord(identity Signer, event Event) (SignedKnowledgeResponsibilityRecord, bool, error) {
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
		OriginPeerID:     EffectiveOriginPeerID(event, identity.PeerID()),
		OriginSequence:   EffectiveOriginSequence(event),
		ResponsibilityID: event.EntityID,
		PCID:             protocols.KnowledgeResponsibilityProfile.CID.String(),
		EnvelopeCID:      envelopeCID.String(),
		EnvelopeBase64:   base64.StdEncoding.EncodeToString(envelopeBytes),
		RecordedAt:       event.Timestamp,
		Implementation:   "ex5-local-runtime",
	}, true, nil
}

func VerifySignedKnowledgeItemRecords(events []Event, records []SignedKnowledgeItemRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]knowledgeItemPayload{}
	for _, event := range events {
		payload, ok := knowledgeItemPayloadForEvent(event)
		if !ok {
			continue
		}
		expected[OriginEventKey(EffectiveOriginPeerID(event, ""), EffectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := strings.TrimSpace(record.OriginPeerID)
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[OriginEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if err := verifyEnvelopeRecord(record.Sequence, record.PCID, record.EnvelopeCID, record.EnvelopeBase64, protocols.KnowledgeItemProfile.CID.String()); err != nil {
			return fmt.Errorf("knowledge-item record %d %w", record.Sequence, err)
		}
		envelope, err := parseEnvelope(record.Sequence, "knowledge-item", record.EnvelopeBase64)
		if err != nil {
			return err
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

func VerifySignedKnowledgeApprovalRecords(events []Event, records []SignedKnowledgeApprovalRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]knowledgeApprovalPayload{}
	for _, event := range events {
		payload, ok := knowledgeApprovalPayloadForEvent(event)
		if !ok {
			continue
		}
		expected[OriginEventKey(EffectiveOriginPeerID(event, ""), EffectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := strings.TrimSpace(record.OriginPeerID)
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[OriginEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if err := verifyEnvelopeRecord(record.Sequence, record.PCID, record.EnvelopeCID, record.EnvelopeBase64, protocols.KnowledgeApprovalProfile.CID.String()); err != nil {
			return fmt.Errorf("knowledge-approval record %d %w", record.Sequence, err)
		}
		envelope, err := parseEnvelope(record.Sequence, "knowledge-approval", record.EnvelopeBase64)
		if err != nil {
			return err
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

func VerifySignedKnowledgeEvidenceRecords(events []Event, records []SignedKnowledgeEvidenceRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]knowledgeEvidencePayload{}
	for _, event := range events {
		payload, ok := knowledgeEvidencePayloadForEvent(event)
		if !ok {
			continue
		}
		expected[OriginEventKey(EffectiveOriginPeerID(event, ""), EffectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := strings.TrimSpace(record.OriginPeerID)
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[OriginEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if err := verifyEnvelopeRecord(record.Sequence, record.PCID, record.EnvelopeCID, record.EnvelopeBase64, protocols.KnowledgeEvidenceProfile.CID.String()); err != nil {
			return fmt.Errorf("knowledge-evidence record %d %w", record.Sequence, err)
		}
		envelope, err := parseEnvelope(record.Sequence, "knowledge-evidence", record.EnvelopeBase64)
		if err != nil {
			return err
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

func VerifySignedOperationalRunRecords(events []Event, records []SignedOperationalRunRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]operationalRunPayload{}
	for _, event := range events {
		payload, ok := operationalRunPayloadForEvent(event)
		if !ok {
			continue
		}
		expected[OriginEventKey(EffectiveOriginPeerID(event, ""), EffectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := strings.TrimSpace(record.OriginPeerID)
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[OriginEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if err := verifyEnvelopeRecord(record.Sequence, record.PCID, record.EnvelopeCID, record.EnvelopeBase64, protocols.OperationalRunProfile.CID.String()); err != nil {
			return fmt.Errorf("operational-run record %d %w", record.Sequence, err)
		}
		envelope, err := parseEnvelope(record.Sequence, "operational-run", record.EnvelopeBase64)
		if err != nil {
			return err
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

func VerifySignedOperationalPlaceRecords(events []Event, records []SignedOperationalPlaceRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]operationalPlacePayload{}
	for _, event := range events {
		payload, ok := operationalPlacePayloadForEvent(event)
		if !ok {
			continue
		}
		expected[OriginEventKey(EffectiveOriginPeerID(event, ""), EffectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := strings.TrimSpace(record.OriginPeerID)
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[OriginEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if err := verifyEnvelopeRecord(record.Sequence, record.PCID, record.EnvelopeCID, record.EnvelopeBase64, protocols.OperationalPlaceProfile.CID.String()); err != nil {
			return fmt.Errorf("operational-place record %d %w", record.Sequence, err)
		}
		envelope, err := parseEnvelope(record.Sequence, "operational-place", record.EnvelopeBase64)
		if err != nil {
			return err
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

func VerifySignedOperationalResourceRecords(events []Event, records []SignedOperationalResourceRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]operationalResourcePayload{}
	for _, event := range events {
		payload, ok := operationalResourcePayloadForEvent(event)
		if !ok {
			continue
		}
		expected[OriginEventKey(EffectiveOriginPeerID(event, ""), EffectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := strings.TrimSpace(record.OriginPeerID)
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[OriginEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if err := verifyEnvelopeRecord(record.Sequence, record.PCID, record.EnvelopeCID, record.EnvelopeBase64, protocols.OperationalResourceProfile.CID.String()); err != nil {
			return fmt.Errorf("operational-resource record %d %w", record.Sequence, err)
		}
		envelope, err := parseEnvelope(record.Sequence, "operational-resource", record.EnvelopeBase64)
		if err != nil {
			return err
		}
		var got operationalResourcePayload
		if err := protocols.Unmarshal(envelope.PayloadBytes, &got); err != nil {
			return fmt.Errorf("decode operational-resource record %d payload: %w", record.Sequence, err)
		}
		if err := compareOperationalResourcePayload(payload, got); err != nil {
			return fmt.Errorf("operational-resource record %d payload mismatch: %w", record.Sequence, err)
		}
	}
	return nil
}

func VerifySignedKnowledgeLinkRecords(events []Event, records []SignedKnowledgeLinkRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]knowledgeLinkPayload{}
	for _, event := range events {
		payload, ok := knowledgeLinkPayloadForEvent(event)
		if !ok {
			continue
		}
		expected[OriginEventKey(EffectiveOriginPeerID(event, ""), EffectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := strings.TrimSpace(record.OriginPeerID)
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[OriginEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if err := verifyEnvelopeRecord(record.Sequence, record.PCID, record.EnvelopeCID, record.EnvelopeBase64, protocols.KnowledgeLinkProfile.CID.String()); err != nil {
			return fmt.Errorf("knowledge-link record %d %w", record.Sequence, err)
		}
		envelope, err := parseEnvelope(record.Sequence, "knowledge-link", record.EnvelopeBase64)
		if err != nil {
			return err
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

func VerifySignedKnowledgeResponsibilityRecords(events []Event, records []SignedKnowledgeResponsibilityRecord) error {
	if len(records) == 0 {
		return nil
	}
	expected := map[string]knowledgeResponsibilityPayload{}
	for _, event := range events {
		payload, ok := knowledgeResponsibilityPayloadForEvent(event)
		if !ok {
			continue
		}
		expected[OriginEventKey(EffectiveOriginPeerID(event, ""), EffectiveOriginSequence(event))] = payload
	}
	for _, record := range records {
		peerID := strings.TrimSpace(record.OriginPeerID)
		originSequence := record.OriginSequence
		if originSequence == 0 {
			originSequence = record.Sequence
		}
		payload, ok := expected[OriginEventKey(peerID, originSequence)]
		if !ok {
			continue
		}
		if err := verifyEnvelopeRecord(record.Sequence, record.PCID, record.EnvelopeCID, record.EnvelopeBase64, protocols.KnowledgeResponsibilityProfile.CID.String()); err != nil {
			return fmt.Errorf("knowledge-responsibility record %d %w", record.Sequence, err)
		}
		envelope, err := parseEnvelope(record.Sequence, "knowledge-responsibility", record.EnvelopeBase64)
		if err != nil {
			return err
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

func verifyEnvelopeRecord(sequence uint64, recordPCID string, envelopeCID string, envelopeBase64 string, expectedPCID string) error {
	if recordPCID != expectedPCID {
		return fmt.Errorf("uses unexpected pCID %q", recordPCID)
	}
	envelopeBytes, err := base64.StdEncoding.DecodeString(envelopeBase64)
	if err != nil {
		return fmt.Errorf("decode envelope: %w", err)
	}
	gotCID, err := protocols.CIDForBytes(envelopeBytes)
	if err != nil {
		return fmt.Errorf("cid envelope: %w", err)
	}
	if gotCID.String() != envelopeCID {
		return fmt.Errorf("envelope cid mismatch")
	}
	envelope, err := protocols.ParseEnvelope(envelopeBytes)
	if err != nil {
		return fmt.Errorf("parse envelope: %w", err)
	}
	if envelope.PCID.String() != expectedPCID {
		return fmt.Errorf("envelope pCID mismatch")
	}
	signable, err := envelope.SignableBytes()
	if err != nil {
		return fmt.Errorf("build signable bytes: %w", err)
	}
	if err := VerifyRuntimeProof(signable, envelope.ProofBytes); err != nil {
		return fmt.Errorf("verify proof: %w", err)
	}
	_ = sequence
	return nil
}

func parseEnvelope(sequence uint64, family string, envelopeBase64 string) (protocols.Envelope, error) {
	envelopeBytes, err := base64.StdEncoding.DecodeString(envelopeBase64)
	if err != nil {
		return protocols.Envelope{}, fmt.Errorf("decode %s record %d envelope: %w", family, sequence, err)
	}
	envelope, err := protocols.ParseEnvelope(envelopeBytes)
	if err != nil {
		return protocols.Envelope{}, fmt.Errorf("parse %s record %d envelope: %w", family, sequence, err)
	}
	return envelope, nil
}

func knowledgeItemPayloadForEvent(event Event) (knowledgeItemPayload, bool) {
	switch event.Type {
	case "knowledge_item_created", "revision_added", "knowledge_item_status_changed", "knowledge_item_superseded":
		return knowledgeItemPayload{
			EntityID:          event.EntityID,
			EventType:         event.Type,
			Sequence:          EffectiveOriginSequence(event),
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

func knowledgeApprovalPayloadForEvent(event Event) (knowledgeApprovalPayload, bool) {
	if event.Type != "approval_recorded" {
		return knowledgeApprovalPayload{}, false
	}
	return knowledgeApprovalPayload{
		EntityID:   event.EntityID,
		EventType:  event.Type,
		Sequence:   EffectiveOriginSequence(event),
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

func knowledgeEvidencePayloadForEvent(event Event) (knowledgeEvidencePayload, bool) {
	if event.Type != "evidence_added" {
		return knowledgeEvidencePayload{}, false
	}
	return knowledgeEvidencePayload{
		EvidenceID:     event.EvidenceID,
		RunID:          event.EntityID,
		EventType:      event.Type,
		Sequence:       EffectiveOriginSequence(event),
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

func operationalRunPayloadForEvent(event Event) (operationalRunPayload, bool) {
	if event.Type != "run_recorded" {
		return operationalRunPayload{}, false
	}
	return operationalRunPayload{
		RunID:             event.EntityID,
		ItemID:            event.TargetID,
		EventType:         event.Type,
		Sequence:          EffectiveOriginSequence(event),
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

func operationalPlacePayloadForEvent(event Event) (operationalPlacePayload, bool) {
	if event.Type != "place_created" {
		return operationalPlacePayload{}, false
	}
	return operationalPlacePayload{
		EntityID:  event.EntityID,
		EventType: event.Type,
		Sequence:  EffectiveOriginSequence(event),
		Timestamp: event.Timestamp,
		Actor:     event.Actor,
		Kind:      event.Kind,
		Name:      event.Name,
		Summary:   event.Summary,
		ParentID:  event.ParentID,
		Tags:      append([]string(nil), event.Tags...),
	}, true
}

func operationalResourcePayloadForEvent(event Event) (operationalResourcePayload, bool) {
	if event.Type != "resource_created" {
		return operationalResourcePayload{}, false
	}
	return operationalResourcePayload{
		EntityID:  event.EntityID,
		EventType: event.Type,
		Sequence:  EffectiveOriginSequence(event),
		Timestamp: event.Timestamp,
		Actor:     event.Actor,
		Kind:      event.Kind,
		Name:      event.Name,
		Summary:   event.Summary,
		PlaceID:   event.PlaceID,
		Tags:      append([]string(nil), event.Tags...),
	}, true
}

func knowledgeLinkPayloadForEvent(event Event) (knowledgeLinkPayload, bool) {
	if event.Type != "link_added" {
		return knowledgeLinkPayload{}, false
	}
	return knowledgeLinkPayload{
		EntityID:  event.EntityID,
		EventType: event.Type,
		Sequence:  EffectiveOriginSequence(event),
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

func knowledgeResponsibilityPayloadForEvent(event Event) (knowledgeResponsibilityPayload, bool) {
	if event.Type != "responsibility_created" {
		return knowledgeResponsibilityPayload{}, false
	}
	return knowledgeResponsibilityPayload{
		EntityID:  event.EntityID,
		EventType: event.Type,
		Sequence:  EffectiveOriginSequence(event),
		Timestamp: event.Timestamp,
		Actor:     event.Actor,
		Title:     event.Title,
		Summary:   event.Summary,
		Team:      event.Team,
		RoleKeys:  append([]string(nil), event.RoleKeys...),
		Tags:      append([]string(nil), event.Tags...),
	}, true
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

func compareOperationalResourcePayload(expected operationalResourcePayload, got operationalResourcePayload) error {
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
	case expected.PlaceID != got.PlaceID:
		return fmt.Errorf("place_id mismatch")
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
