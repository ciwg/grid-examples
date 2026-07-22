package service

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

const peerExchangeBundleFormat = "ex5-peer-exchange-v3"

// Intent: Expose the current relay-visible ex5 PromiseGrid slice as
// whole-family peer exchange over the current local HTTP adapter so peers can
// exchange signed item, approval, evidence, operational-run, link, and
// responsibility artifacts with origin-aware identity, while preserving the
// current family history and compatibility projections. Source: DI-voruk;
// DI-vamok; DI-faruv; DI-ruzok; DI-rumek
func (app *App) ExportPeerExchangeBundle() (PeerExchangeBundle, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	events, err := readEvents(app.store.events)
	if err != nil {
		return PeerExchangeBundle{}, fmt.Errorf("read events: %w", err)
	}
	events = normalizeOperationalEvents(events, app.localPeerID)
	itemRecords, err := app.store.LoadSignedKnowledgeItemRecordsAuthoritative()
	if err != nil {
		return PeerExchangeBundle{}, fmt.Errorf("read knowledge-item records: %w", err)
	}
	approvalRecords, err := app.store.LoadSignedKnowledgeApprovalRecordsAuthoritative()
	if err != nil {
		return PeerExchangeBundle{}, fmt.Errorf("read knowledge-approval records: %w", err)
	}
	evidenceRecords, err := app.store.LoadSignedKnowledgeEvidenceRecordsAuthoritative()
	if err != nil {
		return PeerExchangeBundle{}, fmt.Errorf("read knowledge-evidence records: %w", err)
	}
	runRecords, err := app.store.LoadSignedOperationalRunRecordsAuthoritative()
	if err != nil {
		return PeerExchangeBundle{}, fmt.Errorf("read operational-run records: %w", err)
	}
	linkRecords, err := app.store.LoadSignedKnowledgeLinkRecordsAuthoritative()
	if err != nil {
		return PeerExchangeBundle{}, fmt.Errorf("read knowledge-link records: %w", err)
	}
	responsibilityRecords, err := app.store.LoadSignedKnowledgeResponsibilityRecordsAuthoritative()
	if err != nil {
		return PeerExchangeBundle{}, fmt.Errorf("read knowledge-responsibility records: %w", err)
	}
	itemRecords = normalizeKnowledgeItemRecordOrigins(itemRecords, events)
	approvalRecords = normalizeKnowledgeApprovalRecordOrigins(approvalRecords, events)
	evidenceRecords = normalizeKnowledgeEvidenceRecordOrigins(evidenceRecords, events)
	runRecords = normalizeOperationalRunRecordOrigins(runRecords, events)
	linkRecords = normalizeKnowledgeLinkRecordOrigins(linkRecords, events)
	responsibilityRecords = normalizeKnowledgeResponsibilityRecordOrigins(responsibilityRecords, events)

	filteredEvents := make([]OperationalEvent, 0, len(events))
	for _, event := range events {
		if peerExchangeSupportsEvent(event.Type) {
			filteredEvents = append(filteredEvents, event)
		}
	}
	sort.Slice(filteredEvents, func(i, j int) bool { return filteredEvents[i].Sequence < filteredEvents[j].Sequence })
	sort.Slice(itemRecords, func(i, j int) bool { return itemRecords[i].Sequence < itemRecords[j].Sequence })
	sort.Slice(approvalRecords, func(i, j int) bool { return approvalRecords[i].Sequence < approvalRecords[j].Sequence })
	sort.Slice(evidenceRecords, func(i, j int) bool { return evidenceRecords[i].Sequence < evidenceRecords[j].Sequence })
	sort.Slice(runRecords, func(i, j int) bool { return runRecords[i].Sequence < runRecords[j].Sequence })
	sort.Slice(linkRecords, func(i, j int) bool { return linkRecords[i].Sequence < linkRecords[j].Sequence })
	sort.Slice(responsibilityRecords, func(i, j int) bool { return responsibilityRecords[i].Sequence < responsibilityRecords[j].Sequence })

	casBlobObjects := map[string]string{}
	for _, event := range filteredEvents {
		if event.Type != "evidence_added" || strings.TrimSpace(event.AttachmentCID) == "" {
			continue
		}
		blobBytes, err := app.store.loadCASObject(event.AttachmentCID)
		if err != nil {
			return PeerExchangeBundle{}, fmt.Errorf("read evidence blob %q from cas: %w", event.AttachmentCID, err)
		}
		casBlobObjects[event.AttachmentCID] = base64.StdEncoding.EncodeToString(blobBytes)
	}

	return PeerExchangeBundle{
		Format:                         peerExchangeBundleFormat,
		ExportedAt:                     time.Now().Format(time.RFC3339),
		Implementation:                 "ex5-local-runtime",
		ExportingPeerID:                app.localPeerID,
		KnowledgeItemPCID:              protocols.KnowledgeItemProfile.CID.String(),
		KnowledgeApprovalPCID:          protocols.KnowledgeApprovalProfile.CID.String(),
		KnowledgeEvidencePCID:          protocols.KnowledgeEvidenceProfile.CID.String(),
		KnowledgeLinkPCID:              protocols.KnowledgeLinkProfile.CID.String(),
		KnowledgeResponsibilityPCID:    protocols.KnowledgeResponsibilityProfile.CID.String(),
		OperationalRunPCID:             protocols.OperationalRunProfile.CID.String(),
		Events:                         filteredEvents,
		KnowledgeItemRecords:           itemRecords,
		KnowledgeApprovalRecords:       approvalRecords,
		KnowledgeEvidenceRecords:       evidenceRecords,
		OperationalRunRecords:          runRecords,
		KnowledgeLinkRecords:           linkRecords,
		KnowledgeResponsibilityRecords: responsibilityRecords,
		CASBlobObjects:                 casBlobObjects,
	}, nil
}

// Intent: Accept origin-aware unseen peer history into already-populated
// runtimes while preserving whole-family signed artifacts, keeping local
// sequence as a projection order, and rejecting namespace collisions that the
// current runtime still cannot reconcile honestly. Source: DI-ruzok; DI-rumek
func (app *App) ImportPeerExchangeBundle(bundle PeerExchangeBundle) (PeerExchangeImportResult, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	if err := validatePeerExchangeBundle(bundle); err != nil {
		return PeerExchangeImportResult{}, err
	}

	result := summarizePeerExchangeImport(bundle)
	unseenEvents, unseenEventKeys, err := app.collectUnseenPeerExchangeEventsLocked(bundle)
	if err != nil {
		return PeerExchangeImportResult{}, err
	}
	if len(unseenEvents) == 0 {
		result.ImportedEvents = 0
		result.ImportedKnowledgeItems = 0
		result.ImportedKnowledgeApprovals = 0
		result.ImportedKnowledgeEvidence = 0
		result.ImportedOperationalRuns = 0
		result.ImportedKnowledgeLinks = 0
		result.ImportedResponsibilities = 0
		return result, nil
	}
	for _, record := range filterKnowledgeItemRecordsByOrigin(bundle.KnowledgeItemRecords, unseenEventKeys) {
		if err := app.store.AppendSignedKnowledgeItemRecord(record); err != nil {
			return PeerExchangeImportResult{}, fmt.Errorf("append knowledge-item record: %w", err)
		}
	}
	for _, record := range filterKnowledgeApprovalRecordsByOrigin(bundle.KnowledgeApprovalRecords, unseenEventKeys) {
		if err := app.store.AppendSignedKnowledgeApprovalRecord(record); err != nil {
			return PeerExchangeImportResult{}, fmt.Errorf("append knowledge-approval record: %w", err)
		}
	}
	for cid, blobBase64 := range bundle.CASBlobObjects {
		blobBytes, err := base64.StdEncoding.DecodeString(blobBase64)
		if err != nil {
			return PeerExchangeImportResult{}, fmt.Errorf("decode cas blob %q: %w", cid, err)
		}
		writtenCID, err := app.store.writeCASObject(blobBytes)
		if err != nil {
			return PeerExchangeImportResult{}, fmt.Errorf("write cas blob %q: %w", cid, err)
		}
		if writtenCID != cid {
			return PeerExchangeImportResult{}, fmt.Errorf("cas blob cid mismatch for %q", cid)
		}
	}
	for _, record := range filterKnowledgeEvidenceRecordsByOrigin(bundle.KnowledgeEvidenceRecords, unseenEventKeys) {
		if err := app.store.AppendSignedKnowledgeEvidenceRecord(record); err != nil {
			return PeerExchangeImportResult{}, fmt.Errorf("append knowledge-evidence record: %w", err)
		}
	}
	for _, record := range filterOperationalRunRecordsByOrigin(bundle.OperationalRunRecords, unseenEventKeys) {
		if err := app.store.AppendSignedOperationalRunRecord(record); err != nil {
			return PeerExchangeImportResult{}, fmt.Errorf("append operational-run record: %w", err)
		}
	}
	for _, record := range filterKnowledgeLinkRecordsByOrigin(bundle.KnowledgeLinkRecords, unseenEventKeys) {
		if err := app.store.AppendSignedKnowledgeLinkRecord(record); err != nil {
			return PeerExchangeImportResult{}, fmt.Errorf("append knowledge-link record: %w", err)
		}
	}
	for _, record := range filterKnowledgeResponsibilityRecordsByOrigin(bundle.KnowledgeResponsibilityRecords, unseenEventKeys) {
		if err := app.store.AppendSignedKnowledgeResponsibilityRecord(record); err != nil {
			return PeerExchangeImportResult{}, fmt.Errorf("append knowledge-responsibility record: %w", err)
		}
	}
	for _, event := range unseenEvents {
		app.nextSequence++
		accepted := event
		accepted.Sequence = app.nextSequence
		if err := app.store.AppendEvent(accepted); err != nil {
			app.nextSequence--
			return PeerExchangeImportResult{}, fmt.Errorf("append peer-exchange event: %w", err)
		}
		if err := app.applyEventLocked(accepted); err != nil {
			app.nextSequence--
			return PeerExchangeImportResult{}, fmt.Errorf("apply peer-exchange event %s: %w", originEventKey(accepted.OriginPeerID, accepted.OriginSequence), err)
		}
	}
	result.ImportedEvents = len(unseenEvents)
	result.ImportedKnowledgeItems = len(filterKnowledgeItemRecordsByOrigin(bundle.KnowledgeItemRecords, unseenEventKeys))
	result.ImportedKnowledgeApprovals = len(filterKnowledgeApprovalRecordsByOrigin(bundle.KnowledgeApprovalRecords, unseenEventKeys))
	result.ImportedKnowledgeEvidence = len(filterKnowledgeEvidenceRecordsByOrigin(bundle.KnowledgeEvidenceRecords, unseenEventKeys))
	result.ImportedOperationalRuns = len(filterOperationalRunRecordsByOrigin(bundle.OperationalRunRecords, unseenEventKeys))
	result.ImportedKnowledgeLinks = len(filterKnowledgeLinkRecordsByOrigin(bundle.KnowledgeLinkRecords, unseenEventKeys))
	result.ImportedResponsibilities = len(filterKnowledgeResponsibilityRecordsByOrigin(bundle.KnowledgeResponsibilityRecords, unseenEventKeys))
	return result, nil
}

func peerExchangeSupportsEvent(eventType string) bool {
	switch eventType {
	case "knowledge_item_created", "revision_added", "knowledge_item_status_changed", "knowledge_item_superseded",
		"approval_recorded",
		"evidence_added",
		"run_recorded",
		"link_added",
		"responsibility_created":
		return true
	default:
		return false
	}
}

func validatePeerExchangeBundle(bundle PeerExchangeBundle) error {
	if strings.TrimSpace(bundle.Format) != peerExchangeBundleFormat {
		return fmt.Errorf("unsupported peer exchange bundle format %q", bundle.Format)
	}
	if bundle.KnowledgeItemPCID != protocols.KnowledgeItemProfile.CID.String() {
		return fmt.Errorf("knowledge-item pCID mismatch")
	}
	if bundle.KnowledgeApprovalPCID != protocols.KnowledgeApprovalProfile.CID.String() {
		return fmt.Errorf("knowledge-approval pCID mismatch")
	}
	if bundle.KnowledgeEvidencePCID != protocols.KnowledgeEvidenceProfile.CID.String() {
		return fmt.Errorf("knowledge-evidence pCID mismatch")
	}
	if bundle.OperationalRunPCID != protocols.OperationalRunProfile.CID.String() {
		return fmt.Errorf("operational-run pCID mismatch")
	}
	if bundle.KnowledgeLinkPCID != protocols.KnowledgeLinkProfile.CID.String() {
		return fmt.Errorf("knowledge-link pCID mismatch")
	}
	if bundle.KnowledgeResponsibilityPCID != protocols.KnowledgeResponsibilityProfile.CID.String() {
		return fmt.Errorf("knowledge-responsibility pCID mismatch")
	}
	if len(bundle.Events) == 0 {
		return fmt.Errorf("peer exchange bundle must contain at least one event")
	}
	lastSequence := uint64(0)
	seenSequences := map[uint64]bool{}
	seenOriginKeys := map[string]bool{}
	itemEventSequences := map[string]bool{}
	approvalEventSequences := map[string]bool{}
	evidenceEventSequences := map[string]bool{}
	runEventSequences := map[string]bool{}
	linkEventSequences := map[string]bool{}
	responsibilityEventSequences := map[string]bool{}
	for i, event := range bundle.Events {
		if !peerExchangeSupportsEvent(event.Type) {
			return fmt.Errorf("unsupported peer exchange event type %q", event.Type)
		}
		if i > 0 && event.Sequence <= lastSequence {
			return fmt.Errorf("peer exchange events must be strictly ascending by sequence")
		}
		lastSequence = event.Sequence
		if seenSequences[event.Sequence] {
			return fmt.Errorf("duplicate peer exchange event sequence %d", event.Sequence)
		}
		seenSequences[event.Sequence] = true
		if strings.TrimSpace(event.OriginPeerID) == "" || event.OriginSequence == 0 {
			return fmt.Errorf("peer exchange event %d is missing origin identity", event.Sequence)
		}
		key := originEventKey(event.OriginPeerID, event.OriginSequence)
		if seenOriginKeys[key] {
			return fmt.Errorf("duplicate peer exchange origin tuple %s", key)
		}
		seenOriginKeys[key] = true
		if _, ok := knowledgeItemPayloadForEvent(event); ok {
			itemEventSequences[key] = true
		}
		if _, ok := knowledgeApprovalPayloadForEvent(event); ok {
			approvalEventSequences[key] = true
		}
		if _, ok := knowledgeEvidencePayloadForEvent(event); ok {
			evidenceEventSequences[key] = true
		}
		if _, ok := operationalRunPayloadForEvent(event); ok {
			runEventSequences[key] = true
		}
		if _, ok := knowledgeLinkPayloadForEvent(event); ok {
			linkEventSequences[key] = true
		}
		if _, ok := knowledgeResponsibilityPayloadForEvent(event); ok {
			responsibilityEventSequences[key] = true
		}
	}
	for _, record := range bundle.KnowledgeItemRecords {
		if !itemEventSequences[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			return fmt.Errorf("knowledge-item record %d has no matching event", record.Sequence)
		}
	}
	for _, record := range bundle.KnowledgeApprovalRecords {
		if !approvalEventSequences[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			return fmt.Errorf("knowledge-approval record %d has no matching event", record.Sequence)
		}
	}
	for _, record := range bundle.KnowledgeEvidenceRecords {
		if !evidenceEventSequences[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			return fmt.Errorf("knowledge-evidence record %d has no matching event", record.Sequence)
		}
	}
	for _, record := range bundle.OperationalRunRecords {
		if !runEventSequences[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			return fmt.Errorf("operational-run record %d has no matching event", record.Sequence)
		}
	}
	for _, record := range bundle.KnowledgeLinkRecords {
		if !linkEventSequences[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			return fmt.Errorf("knowledge-link record %d has no matching event", record.Sequence)
		}
	}
	for _, record := range bundle.KnowledgeResponsibilityRecords {
		if !responsibilityEventSequences[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			return fmt.Errorf("knowledge-responsibility record %d has no matching event", record.Sequence)
		}
	}
	if err := verifySignedKnowledgeItemRecords(bundle.Events, bundle.KnowledgeItemRecords); err != nil {
		return fmt.Errorf("verify knowledge-item records: %w", err)
	}
	if err := verifySignedKnowledgeApprovalRecords(bundle.Events, bundle.KnowledgeApprovalRecords); err != nil {
		return fmt.Errorf("verify knowledge-approval records: %w", err)
	}
	if err := verifySignedKnowledgeEvidenceRecords(bundle.Events, bundle.KnowledgeEvidenceRecords); err != nil {
		return fmt.Errorf("verify knowledge-evidence records: %w", err)
	}
	if err := verifySignedOperationalRunRecords(bundle.Events, bundle.OperationalRunRecords); err != nil {
		return fmt.Errorf("verify operational-run records: %w", err)
	}
	if err := verifySignedKnowledgeLinkRecords(bundle.Events, bundle.KnowledgeLinkRecords); err != nil {
		return fmt.Errorf("verify knowledge-link records: %w", err)
	}
	if err := verifySignedKnowledgeResponsibilityRecords(bundle.Events, bundle.KnowledgeResponsibilityRecords); err != nil {
		return fmt.Errorf("verify knowledge-responsibility records: %w", err)
	}
	for _, event := range bundle.Events {
		if event.Type != "evidence_added" || strings.TrimSpace(event.AttachmentCID) == "" {
			continue
		}
		blobBase64, ok := bundle.CASBlobObjects[event.AttachmentCID]
		if !ok {
			return fmt.Errorf("knowledge-evidence blob %q missing from bundle", event.AttachmentCID)
		}
		blobBytes, err := base64.StdEncoding.DecodeString(blobBase64)
		if err != nil {
			return fmt.Errorf("decode knowledge-evidence blob %q: %w", event.AttachmentCID, err)
		}
		blobCID, err := protocols.CIDForBytes(blobBytes)
		if err != nil {
			return fmt.Errorf("cid knowledge-evidence blob %q: %w", event.AttachmentCID, err)
		}
		if blobCID.String() != event.AttachmentCID {
			return fmt.Errorf("knowledge-evidence blob %q cid mismatch", event.AttachmentCID)
		}
	}
	return nil
}

func recordOriginKey(peerID string, originSequence uint64, sequence uint64) string {
	if strings.TrimSpace(peerID) == "" {
		return originEventKey("", sequence)
	}
	if originSequence == 0 {
		return originEventKey(peerID, sequence)
	}
	return originEventKey(peerID, originSequence)
}

func (app *App) collectUnseenPeerExchangeEventsLocked(bundle PeerExchangeBundle) ([]OperationalEvent, map[string]bool, error) {
	unseen := make([]OperationalEvent, 0, len(bundle.Events))
	unseenKeys := map[string]bool{}
	for _, rawEvent := range bundle.Events {
		event := normalizeOperationalEvent(rawEvent, app.localPeerID)
		key := originEventKey(event.OriginPeerID, event.OriginSequence)
		if app.knownOriginEvents[key] {
			continue
		}
		if err := app.validateImportedOriginCollisionLocked(event); err != nil {
			return nil, nil, err
		}
		unseen = append(unseen, event)
		unseenKeys[key] = true
	}
	return unseen, unseenKeys, nil
}

func (app *App) validateImportedOriginCollisionLocked(event OperationalEvent) error {
	switch event.Type {
	case "knowledge_item_created":
		if existing, ok := app.items[event.EntityID]; ok && existing.ID == event.EntityID {
			return fmt.Errorf("knowledge item id %q collides with existing local namespace", event.EntityID)
		}
	case "responsibility_created":
		if existing, ok := app.responsibilities[event.EntityID]; ok && existing.ID == event.EntityID {
			return fmt.Errorf("responsibility id %q collides with existing local namespace", event.EntityID)
		}
	case "run_recorded":
		if existing, ok := app.runs[event.EntityID]; ok && existing.ID == event.EntityID {
			return fmt.Errorf("run id %q collides with existing local namespace", event.EntityID)
		}
	case "approval_recorded":
		if existing, ok := app.approvals[event.EntityID]; ok && existing.ID == event.EntityID {
			return fmt.Errorf("approval id %q collides with existing local namespace", event.EntityID)
		}
	case "link_added":
		if existing, ok := app.links[event.EntityID]; ok && existing.ID == event.EntityID {
			return fmt.Errorf("link id %q collides with existing local namespace", event.EntityID)
		}
	case "evidence_added":
		if strings.TrimSpace(event.EvidenceID) == "" {
			return fmt.Errorf("peer exchange evidence event %s is missing evidence_id", originEventKey(event.OriginPeerID, event.OriginSequence))
		}
		for _, run := range app.runs {
			for _, evidence := range run.Evidence {
				if evidence.ID == event.EvidenceID {
					return fmt.Errorf("evidence id %q collides with existing local namespace", event.EvidenceID)
				}
			}
		}
	}
	return nil
}

func filterKnowledgeItemRecordsByOrigin(records []SignedKnowledgeItemRecord, wanted map[string]bool) []SignedKnowledgeItemRecord {
	out := []SignedKnowledgeItemRecord{}
	for _, record := range records {
		if wanted[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func filterKnowledgeApprovalRecordsByOrigin(records []SignedKnowledgeApprovalRecord, wanted map[string]bool) []SignedKnowledgeApprovalRecord {
	out := []SignedKnowledgeApprovalRecord{}
	for _, record := range records {
		if wanted[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func filterKnowledgeEvidenceRecordsByOrigin(records []SignedKnowledgeEvidenceRecord, wanted map[string]bool) []SignedKnowledgeEvidenceRecord {
	out := []SignedKnowledgeEvidenceRecord{}
	for _, record := range records {
		if wanted[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func filterOperationalRunRecordsByOrigin(records []SignedOperationalRunRecord, wanted map[string]bool) []SignedOperationalRunRecord {
	out := []SignedOperationalRunRecord{}
	for _, record := range records {
		if wanted[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func filterKnowledgeLinkRecordsByOrigin(records []SignedKnowledgeLinkRecord, wanted map[string]bool) []SignedKnowledgeLinkRecord {
	out := []SignedKnowledgeLinkRecord{}
	for _, record := range records {
		if wanted[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func filterKnowledgeResponsibilityRecordsByOrigin(records []SignedKnowledgeResponsibilityRecord, wanted map[string]bool) []SignedKnowledgeResponsibilityRecord {
	out := []SignedKnowledgeResponsibilityRecord{}
	for _, record := range records {
		if wanted[recordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func summarizePeerExchangeImport(bundle PeerExchangeBundle) PeerExchangeImportResult {
	itemIDs := map[string]bool{}
	responsibilityIDs := map[string]bool{}
	runIDs := map[string]bool{}
	for _, event := range bundle.Events {
		switch event.Type {
		case "knowledge_item_created", "revision_added", "knowledge_item_status_changed", "knowledge_item_superseded":
			itemIDs[event.EntityID] = true
		case "responsibility_created":
			responsibilityIDs[event.EntityID] = true
		case "run_recorded":
			runIDs[event.EntityID] = true
		}
	}
	result := PeerExchangeImportResult{
		ImportedEvents:             len(bundle.Events),
		ImportedKnowledgeItems:     len(bundle.KnowledgeItemRecords),
		ImportedKnowledgeApprovals: len(bundle.KnowledgeApprovalRecords),
		ImportedKnowledgeEvidence:  len(bundle.KnowledgeEvidenceRecords),
		ImportedOperationalRuns:    len(bundle.OperationalRunRecords),
		ImportedKnowledgeLinks:     len(bundle.KnowledgeLinkRecords),
		ImportedResponsibilities:   len(bundle.KnowledgeResponsibilityRecords),
	}
	for _, event := range bundle.Events {
		switch event.Type {
		case "run_recorded":
			if strings.TrimSpace(event.PlaceID) != "" {
				result.UnresolvedReferences = append(result.UnresolvedReferences, PeerExchangeImportIssue{
					RecordType: "operational_run",
					RecordID:   event.EntityID,
					Reason:     "place reference is outside the current peer-visible slice",
				})
			}
			for _, resourceID := range event.ResourceIDs {
				result.UnresolvedReferences = append(result.UnresolvedReferences, PeerExchangeImportIssue{
					RecordType: "operational_run",
					RecordID:   event.EntityID,
					Reason:     "resource reference " + resourceID + " is outside the current peer-visible slice",
				})
			}
		case "approval_recorded":
			if event.TargetType == "run" && !runIDs[event.TargetID] {
				result.UnresolvedReferences = append(result.UnresolvedReferences, PeerExchangeImportIssue{
					RecordType: "knowledge_approval",
					RecordID:   event.EntityID,
					Reason:     "run target is missing from the bundle",
				})
			}
		case "link_added":
			if reason := peerExchangeMissingEndpointReason(event.FromType, event.FromID, itemIDs, responsibilityIDs, runIDs); reason != "" {
				result.UnresolvedReferences = append(result.UnresolvedReferences, PeerExchangeImportIssue{
					RecordType: "knowledge_link",
					RecordID:   event.EntityID,
					Reason:     "from endpoint " + reason,
				})
			}
			if reason := peerExchangeMissingEndpointReason(event.ToType, event.ToID, itemIDs, responsibilityIDs, runIDs); reason != "" {
				result.UnresolvedReferences = append(result.UnresolvedReferences, PeerExchangeImportIssue{
					RecordType: "knowledge_link",
					RecordID:   event.EntityID,
					Reason:     "to endpoint " + reason,
				})
			}
		}
	}
	return result
}

func peerExchangeMissingEndpointReason(entityType string, entityID string, itemIDs map[string]bool, responsibilityIDs map[string]bool, runIDs map[string]bool) string {
	switch entityType {
	case "knowledge_item":
		if !itemIDs[entityID] {
			return "knowledge item is missing from the bundle"
		}
	case "responsibility":
		if !responsibilityIDs[entityID] {
			return "responsibility is missing from the bundle"
		}
	case "run":
		if !runIDs[entityID] {
			return "run is missing from the bundle"
		}
	case "place":
		return "place is outside the first peer-exchange slice"
	case "resource":
		return "resource is outside the first peer-exchange slice"
	default:
		return "uses unsupported endpoint type " + entityType
	}
	return ""
}
