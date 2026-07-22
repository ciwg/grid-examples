package service

import "strings"

// Intent: Make peer-visible create artifacts use the create-envelope CID as
// their durable identity while preserving the historical short ID as an alias
// for compatibility replay, verification, and embodiment migration. Source:
// DI-loruk
func decoratePeerVisibleEventCanonicalIDs(
	events []OperationalEvent,
	itemRecords []SignedKnowledgeItemRecord,
	approvalRecords []SignedKnowledgeApprovalRecord,
	evidenceRecords []SignedKnowledgeEvidenceRecord,
	runRecords []SignedOperationalRunRecord,
	linkRecords []SignedKnowledgeLinkRecord,
	responsibilityRecords []SignedKnowledgeResponsibilityRecord,
) []OperationalEvent {
	itemCreateCIDs := map[string]string{}
	for _, record := range itemRecords {
		if record.EventType != "knowledge_item_created" {
			continue
		}
		itemCreateCIDs[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	approvalCreateCIDs := map[string]string{}
	for _, record := range approvalRecords {
		approvalCreateCIDs[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	evidenceCreateCIDs := map[string]string{}
	for _, record := range evidenceRecords {
		evidenceCreateCIDs[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	runCreateCIDs := map[string]string{}
	for _, record := range runRecords {
		runCreateCIDs[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	linkCreateCIDs := map[string]string{}
	for _, record := range linkRecords {
		linkCreateCIDs[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	responsibilityCreateCIDs := map[string]string{}
	for _, record := range responsibilityRecords {
		responsibilityCreateCIDs[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	out := make([]OperationalEvent, 0, len(events))
	for _, event := range events {
		out = append(out, decoratePeerVisibleEventCanonicalID(
			event,
			itemCreateCIDs,
			approvalCreateCIDs,
			evidenceCreateCIDs,
			runCreateCIDs,
			linkCreateCIDs,
			responsibilityCreateCIDs,
		))
	}
	return out
}

func decoratePeerVisibleEventCanonicalID(
	event OperationalEvent,
	itemCreateCIDs map[string]string,
	approvalCreateCIDs map[string]string,
	evidenceCreateCIDs map[string]string,
	runCreateCIDs map[string]string,
	linkCreateCIDs map[string]string,
	responsibilityCreateCIDs map[string]string,
) OperationalEvent {
	key := originEventKey(event.OriginPeerID, event.OriginSequence)
	switch event.Type {
	case "knowledge_item_created":
		if strings.TrimSpace(event.DisplayID) == "" {
			event.DisplayID = event.EntityID
		}
		if itemCreateCIDs[key] != "" {
			event.CanonicalID = itemCreateCIDs[key]
		}
	case "responsibility_created":
		if strings.TrimSpace(event.DisplayID) == "" {
			event.DisplayID = event.EntityID
		}
		if responsibilityCreateCIDs[key] != "" {
			event.CanonicalID = responsibilityCreateCIDs[key]
		}
	case "run_recorded":
		if strings.TrimSpace(event.DisplayID) == "" {
			event.DisplayID = event.EntityID
		}
		if runCreateCIDs[key] != "" {
			event.CanonicalID = runCreateCIDs[key]
		}
	case "approval_recorded":
		if strings.TrimSpace(event.DisplayID) == "" {
			event.DisplayID = event.EntityID
		}
		if approvalCreateCIDs[key] != "" {
			event.CanonicalID = approvalCreateCIDs[key]
		}
	case "link_added":
		if strings.TrimSpace(event.DisplayID) == "" {
			event.DisplayID = event.EntityID
		}
		if linkCreateCIDs[key] != "" {
			event.CanonicalID = linkCreateCIDs[key]
		}
	case "evidence_added":
		if strings.TrimSpace(event.DisplayID) == "" {
			event.DisplayID = event.EvidenceID
		}
		if evidenceCreateCIDs[key] != "" {
			event.CanonicalID = evidenceCreateCIDs[key]
		}
	}
	return event
}

func canonicalOrAliasID(canonicalID string, aliasID string) string {
	if strings.TrimSpace(canonicalID) != "" {
		return canonicalID
	}
	return aliasID
}

func peerVisibleEntityTypeForEvent(event OperationalEvent) string {
	switch event.Type {
	case "knowledge_item_created", "revision_added", "knowledge_item_status_changed", "knowledge_item_superseded":
		return "knowledge_item"
	case "responsibility_created":
		return "responsibility"
	case "run_recorded":
		return "run"
	case "approval_recorded":
		return "approval"
	case "link_added":
		return "link"
	case "evidence_added":
		return "evidence"
	default:
		return ""
	}
}
