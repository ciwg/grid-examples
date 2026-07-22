package service

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
)

const relayFeedFormat = "ex5-relay-feed-v1"

// Intent: Expose a narrower incremental relay feed over the existing signed
// family layer so ex5 can advance beyond whole-bundle exchange without
// collapsing back to local projection sequence as the relay contract.
// Source: DI-pazek
func (app *App) ExportRelayFeed(request RelayFeedRequest) (RelayFeedBatch, error) {
	bundle, err := app.ExportPeerExchangeBundle()
	if err != nil {
		return RelayFeedBatch{}, err
	}
	unseenEvents, unseenKeys := filterRelayFeedEvents(bundle.Events, request.KnownOrigins)
	requiredBlobCIDs := requiredBlobCIDsForEvents(unseenEvents)
	return RelayFeedBatch{
		Format:                         relayFeedFormat,
		ExportedAt:                     bundle.ExportedAt,
		Implementation:                 "ex5-relay-feed",
		ExportingPeerID:                bundle.ExportingPeerID,
		KnowledgeItemPCID:              bundle.KnowledgeItemPCID,
		KnowledgeApprovalPCID:          bundle.KnowledgeApprovalPCID,
		KnowledgeEvidencePCID:          bundle.KnowledgeEvidencePCID,
		KnowledgeLinkPCID:              bundle.KnowledgeLinkPCID,
		KnowledgeResponsibilityPCID:    bundle.KnowledgeResponsibilityPCID,
		OperationalRunPCID:             bundle.OperationalRunPCID,
		OperationalPlacePCID:           bundle.OperationalPlacePCID,
		OperationalResourcePCID:        bundle.OperationalResourcePCID,
		Events:                         unseenEvents,
		KnowledgeItemRecords:           filterKnowledgeItemRecordsByOrigin(bundle.KnowledgeItemRecords, unseenKeys),
		KnowledgeApprovalRecords:       filterKnowledgeApprovalRecordsByOrigin(bundle.KnowledgeApprovalRecords, unseenKeys),
		KnowledgeEvidenceRecords:       filterKnowledgeEvidenceRecordsByOrigin(bundle.KnowledgeEvidenceRecords, unseenKeys),
		OperationalRunRecords:          filterOperationalRunRecordsByOrigin(bundle.OperationalRunRecords, unseenKeys),
		OperationalPlaceRecords:        filterOperationalPlaceRecordsByOrigin(bundle.OperationalPlaceRecords, unseenKeys),
		OperationalResourceRecords:     filterOperationalResourceRecordsByOrigin(bundle.OperationalResourceRecords, unseenKeys),
		KnowledgeLinkRecords:           filterKnowledgeLinkRecordsByOrigin(bundle.KnowledgeLinkRecords, unseenKeys),
		KnowledgeResponsibilityRecords: filterKnowledgeResponsibilityRecordsByOrigin(bundle.KnowledgeResponsibilityRecords, unseenKeys),
		RequiredBlobCIDs:               requiredBlobCIDs,
	}, nil
}

// Intent: Require evidence blobs to be staged into local CAS before importing
// the incremental relay feed so the feed stays record-oriented and blob
// transfer stays explicitly content-addressed by CID.
// Source: DI-pazek
func (app *App) ImportRelayFeed(batch RelayFeedBatch) (RelayFeedImportResult, error) {
	if err := validateRelayFeedBatch(batch); err != nil {
		return RelayFeedImportResult{}, err
	}
	if len(batch.Events) == 0 {
		return RelayFeedImportResult{}, nil
	}
	casBlobObjects, missing, err := app.relayFeedBlobObjects(batch.RequiredBlobCIDs)
	if err != nil {
		return RelayFeedImportResult{}, err
	}
	if len(missing) > 0 {
		return RelayFeedImportResult{MissingBlobCIDs: missing}, nil
	}
	result, err := app.ImportPeerExchangeBundle(PeerExchangeBundle{
		Format:                         peerExchangeBundleFormat,
		ExportedAt:                     batch.ExportedAt,
		Implementation:                 batch.Implementation,
		ExportingPeerID:                batch.ExportingPeerID,
		KnowledgeItemPCID:              batch.KnowledgeItemPCID,
		KnowledgeApprovalPCID:          batch.KnowledgeApprovalPCID,
		KnowledgeEvidencePCID:          batch.KnowledgeEvidencePCID,
		KnowledgeLinkPCID:              batch.KnowledgeLinkPCID,
		KnowledgeResponsibilityPCID:    batch.KnowledgeResponsibilityPCID,
		OperationalRunPCID:             batch.OperationalRunPCID,
		OperationalPlacePCID:           batch.OperationalPlacePCID,
		OperationalResourcePCID:        batch.OperationalResourcePCID,
		Events:                         batch.Events,
		KnowledgeItemRecords:           batch.KnowledgeItemRecords,
		KnowledgeApprovalRecords:       batch.KnowledgeApprovalRecords,
		KnowledgeEvidenceRecords:       batch.KnowledgeEvidenceRecords,
		OperationalRunRecords:          batch.OperationalRunRecords,
		OperationalPlaceRecords:        batch.OperationalPlaceRecords,
		OperationalResourceRecords:     batch.OperationalResourceRecords,
		KnowledgeLinkRecords:           batch.KnowledgeLinkRecords,
		KnowledgeResponsibilityRecords: batch.KnowledgeResponsibilityRecords,
		CASBlobObjects:                 casBlobObjects,
	})
	if err != nil {
		return RelayFeedImportResult{}, err
	}
	return RelayFeedImportResult{
		ImportedEvents:               result.ImportedEvents,
		ImportedKnowledgeItems:       result.ImportedKnowledgeItems,
		ImportedKnowledgeApprovals:   result.ImportedKnowledgeApprovals,
		ImportedKnowledgeEvidence:    result.ImportedKnowledgeEvidence,
		ImportedOperationalRuns:      result.ImportedOperationalRuns,
		ImportedOperationalPlaces:    result.ImportedOperationalPlaces,
		ImportedOperationalResources: result.ImportedOperationalResources,
		ImportedKnowledgeLinks:       result.ImportedKnowledgeLinks,
		ImportedResponsibilities:     result.ImportedResponsibilities,
		UnresolvedReferences:         result.UnresolvedReferences,
	}, nil
}

// Intent: Read relay-carried blobs from local CAS by CID so incremental feed
// import can stay separate from raw blob transfer.
// Source: DI-pazek
func (app *App) RelayBlob(cid string) ([]byte, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	return app.store.loadCASObject(strings.TrimSpace(cid))
}

// Intent: Stage relay-delivered blobs into local CAS only when the uploaded
// bytes actually match the CID named by the relay route.
// Source: DI-pazek
func (app *App) StoreRelayBlob(cid string, body []byte) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	writtenCID, err := app.store.writeCASObject(body)
	if err != nil {
		return err
	}
	if writtenCID != strings.TrimSpace(cid) {
		return fmt.Errorf("relay blob cid mismatch: wrote %q want %q", writtenCID, strings.TrimSpace(cid))
	}
	return nil
}

func filterRelayFeedEvents(events []OperationalEvent, known map[string]uint64) ([]OperationalEvent, map[string]bool) {
	out := make([]OperationalEvent, 0, len(events))
	wanted := map[string]bool{}
	for _, event := range events {
		seen := known[event.OriginPeerID]
		if event.OriginSequence <= seen {
			continue
		}
		out = append(out, event)
		wanted[originEventKey(event.OriginPeerID, event.OriginSequence)] = true
	}
	return out, wanted
}

func requiredBlobCIDsForEvents(events []OperationalEvent) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, event := range events {
		if event.Type != "evidence_added" || strings.TrimSpace(event.AttachmentCID) == "" {
			continue
		}
		cid := strings.TrimSpace(event.AttachmentCID)
		if seen[cid] {
			continue
		}
		seen[cid] = true
		out = append(out, cid)
	}
	sort.Strings(out)
	return out
}

func validateRelayFeedBatch(batch RelayFeedBatch) error {
	if strings.TrimSpace(batch.Format) != relayFeedFormat {
		return fmt.Errorf("unsupported relay feed format %q", batch.Format)
	}
	required := map[string]bool{}
	for _, cid := range batch.RequiredBlobCIDs {
		cid = strings.TrimSpace(cid)
		if cid == "" {
			continue
		}
		required[cid] = true
	}
	for _, event := range batch.Events {
		if event.Type != "evidence_added" || strings.TrimSpace(event.AttachmentCID) == "" {
			continue
		}
		if !required[strings.TrimSpace(event.AttachmentCID)] {
			return fmt.Errorf("relay feed missing required blob cid %q", event.AttachmentCID)
		}
	}
	return nil
}

func (app *App) relayFeedBlobObjects(cids []string) (map[string]string, []string, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	objects := map[string]string{}
	missing := []string{}
	for _, cid := range cids {
		cid = strings.TrimSpace(cid)
		if cid == "" {
			continue
		}
		body, err := app.store.loadCASObject(cid)
		if err != nil {
			if strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "not found") {
				missing = append(missing, cid)
				continue
			}
			return nil, nil, err
		}
		objects[cid] = base64.StdEncoding.EncodeToString(body)
	}
	sort.Strings(missing)
	return objects, missing, nil
}
