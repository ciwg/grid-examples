package service

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	pgtransport "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/transport"
)

// Intent: Expose a narrower incremental relay feed over the existing signed
// family layer so ex5 can advance beyond whole-bundle exchange without
// collapsing back to local projection sequence as the relay contract.
// Source: DI-pazek
func (app *App) ExportRelayFeed(request RelayFeedRequest) (RelayFeedBatch, error) {
	bundle, err := app.ExportPeerExchangeBundle()
	if err != nil {
		return RelayFeedBatch{}, err
	}
	unseenEvents, unseenKeys := pgtransport.FilterRelayFeedEvents(bundle.Events, request.KnownOrigins)
	requiredBlobCIDs := pgtransport.RequiredBlobCIDsForEvents(unseenEvents)
	return RelayFeedBatch{
		Format:                         pgtransport.RelayFeedFormat,
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
		KnowledgeItemRecords:           pgtransport.FilterKnowledgeItemRecordsByOrigin(bundle.KnowledgeItemRecords, unseenKeys),
		KnowledgeApprovalRecords:       pgtransport.FilterKnowledgeApprovalRecordsByOrigin(bundle.KnowledgeApprovalRecords, unseenKeys),
		KnowledgeEvidenceRecords:       pgtransport.FilterKnowledgeEvidenceRecordsByOrigin(bundle.KnowledgeEvidenceRecords, unseenKeys),
		OperationalRunRecords:          pgtransport.FilterOperationalRunRecordsByOrigin(bundle.OperationalRunRecords, unseenKeys),
		OperationalPlaceRecords:        pgtransport.FilterOperationalPlaceRecordsByOrigin(bundle.OperationalPlaceRecords, unseenKeys),
		OperationalResourceRecords:     pgtransport.FilterOperationalResourceRecordsByOrigin(bundle.OperationalResourceRecords, unseenKeys),
		KnowledgeLinkRecords:           pgtransport.FilterKnowledgeLinkRecordsByOrigin(bundle.KnowledgeLinkRecords, unseenKeys),
		KnowledgeResponsibilityRecords: pgtransport.FilterKnowledgeResponsibilityRecordsByOrigin(bundle.KnowledgeResponsibilityRecords, unseenKeys),
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
		Format:                         pgtransport.PeerExchangeBundleFormat,
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
	return app.store.cas.LoadObject(strings.TrimSpace(cid))
}

// Intent: Stage relay-delivered blobs into local CAS only when the uploaded
// bytes actually match the CID named by the relay route.
// Source: DI-pazek
func (app *App) StoreRelayBlob(cid string, body []byte) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	writtenCID, err := app.store.cas.WriteObject(body)
	if err != nil {
		return err
	}
	if writtenCID != strings.TrimSpace(cid) {
		return fmt.Errorf("relay blob cid mismatch: wrote %q want %q", writtenCID, strings.TrimSpace(cid))
	}
	return nil
}

func validateRelayFeedBatch(batch RelayFeedBatch) error {
	if strings.TrimSpace(batch.Format) != pgtransport.RelayFeedFormat {
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
		body, err := app.store.cas.LoadObject(cid)
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
