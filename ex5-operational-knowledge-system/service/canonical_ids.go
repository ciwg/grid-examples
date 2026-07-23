package service

import records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"

// Intent: Preserve create-envelope-CID durable identity semantics while the
// canonical ID decoration logic moves into the reusable PromiseGrid record
// substrate. Source: DI-ragiv
func decoratePeerVisibleEventCanonicalIDs(
	events []OperationalEvent,
	itemRecords []SignedKnowledgeItemRecord,
	approvalRecords []SignedKnowledgeApprovalRecord,
	evidenceRecords []SignedKnowledgeEvidenceRecord,
	runRecords []SignedOperationalRunRecord,
	placeRecords []SignedOperationalPlaceRecord,
	resourceRecords []SignedOperationalResourceRecord,
	linkRecords []SignedKnowledgeLinkRecord,
	responsibilityRecords []SignedKnowledgeResponsibilityRecord,
) []OperationalEvent {
	eventSlice := make([]records.Event, len(events))
	itemSlice := make([]records.SignedKnowledgeItemRecord, len(itemRecords))
	approvalSlice := make([]records.SignedKnowledgeApprovalRecord, len(approvalRecords))
	evidenceSlice := make([]records.SignedKnowledgeEvidenceRecord, len(evidenceRecords))
	runSlice := make([]records.SignedOperationalRunRecord, len(runRecords))
	placeSlice := make([]records.SignedOperationalPlaceRecord, len(placeRecords))
	resourceSlice := make([]records.SignedOperationalResourceRecord, len(resourceRecords))
	linkSlice := make([]records.SignedKnowledgeLinkRecord, len(linkRecords))
	responsibilitySlice := make([]records.SignedKnowledgeResponsibilityRecord, len(responsibilityRecords))
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	for i, record := range itemRecords {
		itemSlice[i] = records.SignedKnowledgeItemRecord(record)
	}
	for i, record := range approvalRecords {
		approvalSlice[i] = records.SignedKnowledgeApprovalRecord(record)
	}
	for i, record := range evidenceRecords {
		evidenceSlice[i] = records.SignedKnowledgeEvidenceRecord(record)
	}
	for i, record := range runRecords {
		runSlice[i] = records.SignedOperationalRunRecord(record)
	}
	for i, record := range placeRecords {
		placeSlice[i] = records.SignedOperationalPlaceRecord(record)
	}
	for i, record := range resourceRecords {
		resourceSlice[i] = records.SignedOperationalResourceRecord(record)
	}
	for i, record := range linkRecords {
		linkSlice[i] = records.SignedKnowledgeLinkRecord(record)
	}
	for i, record := range responsibilityRecords {
		responsibilitySlice[i] = records.SignedKnowledgeResponsibilityRecord(record)
	}
	out := records.DecoratePeerVisibleEventCanonicalIDs(
		eventSlice,
		itemSlice,
		approvalSlice,
		evidenceSlice,
		runSlice,
		placeSlice,
		resourceSlice,
		linkSlice,
		responsibilitySlice,
	)
	converted := make([]OperationalEvent, len(out))
	for i, event := range out {
		converted[i] = OperationalEvent(event)
	}
	return converted
}

func canonicalOrAliasID(canonicalID string, aliasID string) string {
	return records.CanonicalOrAliasID(canonicalID, aliasID)
}

func peerVisibleEntityTypeForEvent(event OperationalEvent) string {
	return records.PeerVisibleEntityTypeForEvent(records.Event(event))
}
