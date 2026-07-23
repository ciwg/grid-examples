package records

import (
	"fmt"
	"strings"
)

func EffectiveOriginPeerID(event Event, localPeerID string) string {
	if strings.TrimSpace(event.OriginPeerID) != "" {
		return event.OriginPeerID
	}
	return localPeerID
}

func EffectiveOriginSequence(event Event) uint64 {
	if event.OriginSequence != 0 {
		return event.OriginSequence
	}
	return event.Sequence
}

func NormalizeEvents(events []Event, localPeerID string) []Event {
	out := make([]Event, 0, len(events))
	for _, event := range events {
		out = append(out, NormalizeEvent(event, localPeerID))
	}
	return out
}

func NormalizeEvent(event Event, localPeerID string) Event {
	if strings.TrimSpace(event.OriginPeerID) == "" {
		event.OriginPeerID = localPeerID
	}
	if event.OriginSequence == 0 {
		event.OriginSequence = event.Sequence
	}
	return event
}

func OriginEventKey(peerID string, originSequence uint64) string {
	return peerID + "#" + fmt.Sprintf("%d", originSequence)
}

func RecordOriginKey(peerID string, originSequence uint64, sequence uint64) string {
	if strings.TrimSpace(peerID) == "" {
		return OriginEventKey("", sequence)
	}
	if originSequence == 0 {
		return OriginEventKey(peerID, sequence)
	}
	return OriginEventKey(peerID, originSequence)
}

func NormalizeKnowledgeItemRecordOrigins(records []SignedKnowledgeItemRecord, events []Event) []SignedKnowledgeItemRecord {
	out := make([]SignedKnowledgeItemRecord, 0, len(records))
	eventsBySequence := map[uint64]Event{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func NormalizeKnowledgeApprovalRecordOrigins(records []SignedKnowledgeApprovalRecord, events []Event) []SignedKnowledgeApprovalRecord {
	out := make([]SignedKnowledgeApprovalRecord, 0, len(records))
	eventsBySequence := map[uint64]Event{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func NormalizeKnowledgeEvidenceRecordOrigins(records []SignedKnowledgeEvidenceRecord, events []Event) []SignedKnowledgeEvidenceRecord {
	out := make([]SignedKnowledgeEvidenceRecord, 0, len(records))
	eventsBySequence := map[uint64]Event{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func NormalizeOperationalRunRecordOrigins(records []SignedOperationalRunRecord, events []Event) []SignedOperationalRunRecord {
	out := make([]SignedOperationalRunRecord, 0, len(records))
	eventsBySequence := map[uint64]Event{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func NormalizeOperationalPlaceRecordOrigins(records []SignedOperationalPlaceRecord, events []Event) []SignedOperationalPlaceRecord {
	out := make([]SignedOperationalPlaceRecord, 0, len(records))
	eventsBySequence := map[uint64]Event{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func NormalizeOperationalResourceRecordOrigins(records []SignedOperationalResourceRecord, events []Event) []SignedOperationalResourceRecord {
	out := make([]SignedOperationalResourceRecord, 0, len(records))
	eventsBySequence := map[uint64]Event{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func NormalizeKnowledgeLinkRecordOrigins(records []SignedKnowledgeLinkRecord, events []Event) []SignedKnowledgeLinkRecord {
	out := make([]SignedKnowledgeLinkRecord, 0, len(records))
	eventsBySequence := map[uint64]Event{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func NormalizeKnowledgeResponsibilityRecordOrigins(records []SignedKnowledgeResponsibilityRecord, events []Event) []SignedKnowledgeResponsibilityRecord {
	out := make([]SignedKnowledgeResponsibilityRecord, 0, len(records))
	eventsBySequence := map[uint64]Event{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

// Intent: Keep canonical durable ID derivation in the reusable record
// substrate so create-envelope CIDs define peer-visible durable identity
// independently of ex5 projection code. Source: DI-ragiv
func DecoratePeerVisibleEventCanonicalIDs(
	events []Event,
	itemRecords []SignedKnowledgeItemRecord,
	approvalRecords []SignedKnowledgeApprovalRecord,
	evidenceRecords []SignedKnowledgeEvidenceRecord,
	runRecords []SignedOperationalRunRecord,
	placeRecords []SignedOperationalPlaceRecord,
	resourceRecords []SignedOperationalResourceRecord,
	linkRecords []SignedKnowledgeLinkRecord,
	responsibilityRecords []SignedKnowledgeResponsibilityRecord,
) []Event {
	itemCreateCIDs := map[string]string{}
	for _, record := range itemRecords {
		if record.EventType != "knowledge_item_created" {
			continue
		}
		itemCreateCIDs[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	approvalCreateCIDs := map[string]string{}
	for _, record := range approvalRecords {
		approvalCreateCIDs[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	evidenceCreateCIDs := map[string]string{}
	for _, record := range evidenceRecords {
		evidenceCreateCIDs[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	runCreateCIDs := map[string]string{}
	for _, record := range runRecords {
		runCreateCIDs[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	placeCreateCIDs := map[string]string{}
	for _, record := range placeRecords {
		placeCreateCIDs[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	resourceCreateCIDs := map[string]string{}
	for _, record := range resourceRecords {
		resourceCreateCIDs[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	linkCreateCIDs := map[string]string{}
	for _, record := range linkRecords {
		linkCreateCIDs[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	responsibilityCreateCIDs := map[string]string{}
	for _, record := range responsibilityRecords {
		responsibilityCreateCIDs[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] = record.EnvelopeCID
	}
	out := make([]Event, 0, len(events))
	for _, event := range events {
		out = append(out, decoratePeerVisibleEventCanonicalID(
			event,
			itemCreateCIDs,
			approvalCreateCIDs,
			evidenceCreateCIDs,
			runCreateCIDs,
			placeCreateCIDs,
			resourceCreateCIDs,
			linkCreateCIDs,
			responsibilityCreateCIDs,
		))
	}
	return out
}

func decoratePeerVisibleEventCanonicalID(
	event Event,
	itemCreateCIDs map[string]string,
	approvalCreateCIDs map[string]string,
	evidenceCreateCIDs map[string]string,
	runCreateCIDs map[string]string,
	placeCreateCIDs map[string]string,
	resourceCreateCIDs map[string]string,
	linkCreateCIDs map[string]string,
	responsibilityCreateCIDs map[string]string,
) Event {
	key := OriginEventKey(event.OriginPeerID, event.OriginSequence)
	switch event.Type {
	case "place_created":
		if strings.TrimSpace(event.DisplayID) == "" {
			event.DisplayID = event.EntityID
		}
		if placeCreateCIDs[key] != "" {
			event.CanonicalID = placeCreateCIDs[key]
		}
	case "resource_created":
		if strings.TrimSpace(event.DisplayID) == "" {
			event.DisplayID = event.EntityID
		}
		if resourceCreateCIDs[key] != "" {
			event.CanonicalID = resourceCreateCIDs[key]
		}
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

func CanonicalOrAliasID(canonicalID string, aliasID string) string {
	if strings.TrimSpace(canonicalID) != "" {
		return canonicalID
	}
	return aliasID
}

func PeerVisibleEntityTypeForEvent(event Event) string {
	switch event.Type {
	case "place_created":
		return "place"
	case "resource_created":
		return "resource"
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

func cloneFacts(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
