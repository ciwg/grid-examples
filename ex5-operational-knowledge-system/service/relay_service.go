package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

const (
	relayServiceName = "operational-relay"
	relayRoutePrefix = "/relay/v1"
)

type Relay struct {
	mu                             sync.Mutex
	store                          *Store
	cursorPath                     string
	originCursors                  map[string]uint64
	events                         []OperationalEvent
	knowledgeItemRecords           []SignedKnowledgeItemRecord
	knowledgeApprovalRecords       []SignedKnowledgeApprovalRecord
	knowledgeEvidenceRecords       []SignedKnowledgeEvidenceRecord
	operationalRunRecords          []SignedOperationalRunRecord
	operationalPlaceRecords        []SignedOperationalPlaceRecord
	operationalResourceRecords     []SignedOperationalResourceRecord
	knowledgeLinkRecords           []SignedKnowledgeLinkRecord
	knowledgeResponsibilityRecords []SignedKnowledgeResponsibilityRecord
}

type relayCursorState struct {
	KnownOrigins map[string]uint64 `json:"known_origins"`
}

// Intent: Keep the first remote relay as a role-pure durable store-and-forward
// service with its own state root instead of rehosting the local embodiment
// adapter surface. Source: DI-rovik; DI-tasov; DI-nulav
func NewRelay(root string) (*Relay, error) {
	store, events, itemRecords, approvalRecords, evidenceRecords, runRecords, placeRecords, resourceRecords, linkRecords, responsibilityRecords, err := openRelayStore(root)
	if err != nil {
		return nil, err
	}
	relay := &Relay{
		store:                          store,
		cursorPath:                     filepath.Join(root, "origin-cursors.json"),
		events:                         events,
		knowledgeItemRecords:           itemRecords,
		knowledgeApprovalRecords:       approvalRecords,
		knowledgeEvidenceRecords:       evidenceRecords,
		operationalRunRecords:          runRecords,
		operationalPlaceRecords:        placeRecords,
		operationalResourceRecords:     resourceRecords,
		knowledgeLinkRecords:           linkRecords,
		knowledgeResponsibilityRecords: responsibilityRecords,
	}
	if err := relay.verifyStoredHistory(); err != nil {
		_ = relay.store.Close()
		return nil, err
	}
	computed, err := computeRelayOriginCursors(events)
	if err != nil {
		_ = relay.store.Close()
		return nil, err
	}
	saved, err := loadRelayOriginCursors(relay.cursorPath)
	if err != nil {
		_ = relay.store.Close()
		return nil, err
	}
	if !relayCursorMapsEqual(saved, computed) {
		if err := saveRelayOriginCursors(relay.cursorPath, computed); err != nil {
			_ = relay.store.Close()
			return nil, err
		}
	}
	relay.originCursors = computed
	return relay, nil
}

func (relay *Relay) Close() error {
	if relay == nil {
		return nil
	}
	return relay.store.Close()
}

func (relay *Relay) Meta() RelayMeta {
	relay.mu.Lock()
	defer relay.mu.Unlock()
	return RelayMeta{
		DataRoot:                   relay.store.root,
		ServiceName:                relayServiceName,
		RoutePrefix:                relayRoutePrefix,
		RelayFeedFormat:            relayFeedFormat,
		RelayFeedFamilies:          relayFeedFamilies(),
		RelayBlobTransferEnabled:   true,
		PublishRequiresStagedBlobs: true,
	}
}

// Intent: Publish only contiguous unseen origin-aware relay history after all
// referenced blobs are already staged, so the remote relay never advertises
// evidence-bearing durable history without the corresponding CAS objects.
// Source: DI-tasov; DI-nulav
func (relay *Relay) Publish(batch RelayFeedBatch) (RelayPublishResult, error) {
	relay.mu.Lock()
	defer relay.mu.Unlock()

	if err := validateRelayFeedBatch(batch); err != nil {
		return RelayPublishResult{}, err
	}
	missing, err := relay.missingBlobCIDs(batch.RequiredBlobCIDs)
	if err != nil {
		return RelayPublishResult{}, err
	}
	if len(missing) > 0 {
		return RelayPublishResult{MissingBlobCIDs: missing}, nil
	}
	if err := relay.validateRemoteRelayBatch(batch); err != nil {
		return RelayPublishResult{}, err
	}

	unseenEvents, unseenKeys, nextCursors, err := relay.collectRelayPublishEvents(batch.Events)
	if err != nil {
		return RelayPublishResult{}, err
	}
	if len(unseenEvents) == 0 {
		return RelayPublishResult{}, nil
	}

	itemRecords := filterKnowledgeItemRecordsByOrigin(batch.KnowledgeItemRecords, unseenKeys)
	approvalRecords := filterKnowledgeApprovalRecordsByOrigin(batch.KnowledgeApprovalRecords, unseenKeys)
	evidenceRecords := filterKnowledgeEvidenceRecordsByOrigin(batch.KnowledgeEvidenceRecords, unseenKeys)
	runRecords := filterOperationalRunRecordsByOrigin(batch.OperationalRunRecords, unseenKeys)
	placeRecords := filterOperationalPlaceRecordsByOrigin(batch.OperationalPlaceRecords, unseenKeys)
	resourceRecords := filterOperationalResourceRecordsByOrigin(batch.OperationalResourceRecords, unseenKeys)
	linkRecords := filterKnowledgeLinkRecordsByOrigin(batch.KnowledgeLinkRecords, unseenKeys)
	responsibilityRecords := filterKnowledgeResponsibilityRecordsByOrigin(batch.KnowledgeResponsibilityRecords, unseenKeys)

	for _, record := range itemRecords {
		if err := relay.store.AppendSignedKnowledgeItemRecord(record); err != nil {
			return RelayPublishResult{}, fmt.Errorf("append knowledge-item relay record: %w", err)
		}
	}
	for _, record := range approvalRecords {
		if err := relay.store.AppendSignedKnowledgeApprovalRecord(record); err != nil {
			return RelayPublishResult{}, fmt.Errorf("append knowledge-approval relay record: %w", err)
		}
	}
	for _, record := range evidenceRecords {
		if err := relay.store.AppendSignedKnowledgeEvidenceRecord(record); err != nil {
			return RelayPublishResult{}, fmt.Errorf("append knowledge-evidence relay record: %w", err)
		}
	}
	for _, record := range runRecords {
		if err := relay.store.AppendSignedOperationalRunRecord(record); err != nil {
			return RelayPublishResult{}, fmt.Errorf("append operational-run relay record: %w", err)
		}
	}
	for _, record := range placeRecords {
		if err := relay.store.AppendSignedOperationalPlaceRecord(record); err != nil {
			return RelayPublishResult{}, fmt.Errorf("append operational-place relay record: %w", err)
		}
	}
	for _, record := range resourceRecords {
		if err := relay.store.AppendSignedOperationalResourceRecord(record); err != nil {
			return RelayPublishResult{}, fmt.Errorf("append operational-resource relay record: %w", err)
		}
	}
	for _, record := range linkRecords {
		if err := relay.store.AppendSignedKnowledgeLinkRecord(record); err != nil {
			return RelayPublishResult{}, fmt.Errorf("append knowledge-link relay record: %w", err)
		}
	}
	for _, record := range responsibilityRecords {
		if err := relay.store.AppendSignedKnowledgeResponsibilityRecord(record); err != nil {
			return RelayPublishResult{}, fmt.Errorf("append knowledge-responsibility relay record: %w", err)
		}
	}
	for _, event := range unseenEvents {
		if err := relay.store.AppendEvent(event); err != nil {
			return RelayPublishResult{}, fmt.Errorf("append relay event: %w", err)
		}
	}
	if err := saveRelayOriginCursors(relay.cursorPath, nextCursors); err != nil {
		return RelayPublishResult{}, err
	}

	relay.events = append(relay.events, unseenEvents...)
	relay.knowledgeItemRecords = append(relay.knowledgeItemRecords, itemRecords...)
	relay.knowledgeApprovalRecords = append(relay.knowledgeApprovalRecords, approvalRecords...)
	relay.knowledgeEvidenceRecords = append(relay.knowledgeEvidenceRecords, evidenceRecords...)
	relay.operationalRunRecords = append(relay.operationalRunRecords, runRecords...)
	relay.operationalPlaceRecords = append(relay.operationalPlaceRecords, placeRecords...)
	relay.operationalResourceRecords = append(relay.operationalResourceRecords, resourceRecords...)
	relay.knowledgeLinkRecords = append(relay.knowledgeLinkRecords, linkRecords...)
	relay.knowledgeResponsibilityRecords = append(relay.knowledgeResponsibilityRecords, responsibilityRecords...)
	relay.originCursors = nextCursors

	return RelayPublishResult{
		PublishedEvents:               len(unseenEvents),
		PublishedKnowledgeItems:       len(itemRecords),
		PublishedKnowledgeApprovals:   len(approvalRecords),
		PublishedKnowledgeEvidence:    len(evidenceRecords),
		PublishedOperationalRuns:      len(runRecords),
		PublishedOperationalPlaces:    len(placeRecords),
		PublishedOperationalResources: len(resourceRecords),
		PublishedKnowledgeLinks:       len(linkRecords),
		PublishedResponsibilities:     len(responsibilityRecords),
	}, nil
}

// Intent: Pull remote relay history by origin-aware cursor map while emitting a
// fresh batch-local compatibility sequence so mixed-peer relay slices do not
// leak exporter-local sequence collisions back into import validation.
// Source: DI-tasov
func (relay *Relay) Pull(request RelayFeedRequest) (RelayFeedBatch, error) {
	relay.mu.Lock()
	defer relay.mu.Unlock()

	unseenEvents, unseenKeys := filterRelayFeedEvents(relay.events, request.KnownOrigins)
	return RelayFeedBatch{
		Format:                         relayFeedFormat,
		ExportedAt:                     time.Now().Format(time.RFC3339),
		Implementation:                 relayServiceName,
		ExportingPeerID:                relayServiceName,
		KnowledgeItemPCID:              protocols.KnowledgeItemProfile.CID.String(),
		KnowledgeApprovalPCID:          protocols.KnowledgeApprovalProfile.CID.String(),
		KnowledgeEvidencePCID:          protocols.KnowledgeEvidenceProfile.CID.String(),
		KnowledgeLinkPCID:              protocols.KnowledgeLinkProfile.CID.String(),
		KnowledgeResponsibilityPCID:    protocols.KnowledgeResponsibilityProfile.CID.String(),
		OperationalRunPCID:             protocols.OperationalRunProfile.CID.String(),
		OperationalPlacePCID:           protocols.OperationalPlaceProfile.CID.String(),
		OperationalResourcePCID:        protocols.OperationalResourceProfile.CID.String(),
		Events:                         renumberRelayBatchEvents(unseenEvents),
		KnowledgeItemRecords:           filterKnowledgeItemRecordsByOrigin(relay.knowledgeItemRecords, unseenKeys),
		KnowledgeApprovalRecords:       filterKnowledgeApprovalRecordsByOrigin(relay.knowledgeApprovalRecords, unseenKeys),
		KnowledgeEvidenceRecords:       filterKnowledgeEvidenceRecordsByOrigin(relay.knowledgeEvidenceRecords, unseenKeys),
		OperationalRunRecords:          filterOperationalRunRecordsByOrigin(relay.operationalRunRecords, unseenKeys),
		OperationalPlaceRecords:        filterOperationalPlaceRecordsByOrigin(relay.operationalPlaceRecords, unseenKeys),
		OperationalResourceRecords:     filterOperationalResourceRecordsByOrigin(relay.operationalResourceRecords, unseenKeys),
		KnowledgeLinkRecords:           filterKnowledgeLinkRecordsByOrigin(relay.knowledgeLinkRecords, unseenKeys),
		KnowledgeResponsibilityRecords: filterKnowledgeResponsibilityRecordsByOrigin(relay.knowledgeResponsibilityRecords, unseenKeys),
		RequiredBlobCIDs:               requiredBlobCIDsForEvents(unseenEvents),
	}, nil
}

func (relay *Relay) Blob(cid string) ([]byte, error) {
	relay.mu.Lock()
	defer relay.mu.Unlock()
	return relay.store.loadCASObject(strings.TrimSpace(cid))
}

// Intent: Stage raw relay blobs by CID before any evidence-bearing feed publish
// succeeds so remote relay durability matches the signed history it accepts.
// Source: DI-nulav
func (relay *Relay) StoreBlob(cid string, body []byte) error {
	relay.mu.Lock()
	defer relay.mu.Unlock()
	writtenCID, err := relay.store.writeCASObject(body)
	if err != nil {
		return err
	}
	if writtenCID != strings.TrimSpace(cid) {
		return fmt.Errorf("relay blob cid mismatch: wrote %q want %q", writtenCID, strings.TrimSpace(cid))
	}
	return nil
}

func (relay *Relay) verifyStoredHistory() error {
	if err := verifySignedKnowledgeItemRecords(relay.events, relay.knowledgeItemRecords); err != nil {
		return err
	}
	if err := verifySignedKnowledgeApprovalRecords(relay.events, relay.knowledgeApprovalRecords); err != nil {
		return err
	}
	if err := verifySignedKnowledgeEvidenceRecords(relay.events, relay.knowledgeEvidenceRecords); err != nil {
		return err
	}
	if err := verifySignedOperationalRunRecords(relay.events, relay.operationalRunRecords); err != nil {
		return err
	}
	if err := verifySignedOperationalPlaceRecords(relay.events, relay.operationalPlaceRecords); err != nil {
		return err
	}
	if err := verifySignedOperationalResourceRecords(relay.events, relay.operationalResourceRecords); err != nil {
		return err
	}
	if err := verifySignedKnowledgeLinkRecords(relay.events, relay.knowledgeLinkRecords); err != nil {
		return err
	}
	if err := verifySignedKnowledgeResponsibilityRecords(relay.events, relay.knowledgeResponsibilityRecords); err != nil {
		return err
	}
	return nil
}

func (relay *Relay) missingBlobCIDs(cids []string) ([]string, error) {
	missing := []string{}
	for _, cid := range cids {
		cid = strings.TrimSpace(cid)
		if cid == "" {
			continue
		}
		if _, err := relay.store.loadCASObject(cid); err != nil {
			if os.IsNotExist(err) {
				missing = append(missing, cid)
				continue
			}
			return nil, err
		}
	}
	sort.Strings(missing)
	return missing, nil
}

func (relay *Relay) collectRelayPublishEvents(events []OperationalEvent) ([]OperationalEvent, map[string]bool, map[string]uint64, error) {
	perOrigin := map[string][]uint64{}
	for _, event := range events {
		key := strings.TrimSpace(event.OriginPeerID)
		perOrigin[key] = append(perOrigin[key], event.OriginSequence)
	}
	next := cloneKnownOrigins(relay.originCursors)
	for originPeerID, sequences := range perOrigin {
		sort.Slice(sequences, func(i, j int) bool { return sequences[i] < sequences[j] })
		cursor := next[originPeerID]
		firstUnseen := uint64(0)
		last := cursor
		for _, sequence := range sequences {
			if sequence <= cursor {
				continue
			}
			if firstUnseen == 0 {
				firstUnseen = sequence
			}
			if sequence != last+1 {
				return nil, nil, nil, fmt.Errorf("relay publish for %q must be contiguous from %d, got %d", originPeerID, last+1, sequence)
			}
			last = sequence
		}
		if firstUnseen == 0 {
			continue
		}
		if cursor == 0 && firstUnseen != 1 {
			return nil, nil, nil, fmt.Errorf("relay publish for new origin %q must start at 1, got %d", originPeerID, firstUnseen)
		}
		next[originPeerID] = last
	}

	unseenEvents := make([]OperationalEvent, 0, len(events))
	unseenKeys := map[string]bool{}
	for _, event := range events {
		if event.OriginSequence <= relay.originCursors[event.OriginPeerID] {
			continue
		}
		unseenEvents = append(unseenEvents, event)
		unseenKeys[originEventKey(event.OriginPeerID, event.OriginSequence)] = true
	}
	return unseenEvents, unseenKeys, next, nil
}

func (relay *Relay) validateRemoteRelayBatch(batch RelayFeedBatch) error {
	casBlobObjects := map[string]string{}
	for _, cid := range batch.RequiredBlobCIDs {
		cid = strings.TrimSpace(cid)
		if cid == "" {
			continue
		}
		body, err := relay.store.loadCASObject(cid)
		if err != nil {
			return err
		}
		casBlobObjects[cid] = base64.StdEncoding.EncodeToString(body)
	}
	bundle := PeerExchangeBundle{
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
	}
	return validatePeerExchangeBundle(bundle)
}

func renumberRelayBatchEvents(events []OperationalEvent) []OperationalEvent {
	out := append([]OperationalEvent(nil), events...)
	for i := range out {
		out[i].Sequence = uint64(i + 1)
	}
	return out
}

func relayFeedFamilies() []string {
	return []string{
		"knowledge-item",
		"knowledge-approval",
		"knowledge-evidence",
		"operational-run",
		"operational-place",
		"operational-resource",
		"knowledge-link",
		"knowledge-responsibility",
	}
}

func loadRelayOriginCursors(path string) (map[string]uint64, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]uint64{}, nil
		}
		return nil, err
	}
	var state relayCursorState
	if err := json.Unmarshal(body, &state); err != nil {
		return nil, err
	}
	if state.KnownOrigins == nil {
		return map[string]uint64{}, nil
	}
	return state.KnownOrigins, nil
}

func saveRelayOriginCursors(path string, known map[string]uint64) error {
	body, err := json.Marshal(relayCursorState{KnownOrigins: known})
	if err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}

func computeRelayOriginCursors(events []OperationalEvent) (map[string]uint64, error) {
	perOrigin := map[string][]uint64{}
	for _, event := range events {
		perOrigin[event.OriginPeerID] = append(perOrigin[event.OriginPeerID], event.OriginSequence)
	}
	out := map[string]uint64{}
	for originPeerID, sequences := range perOrigin {
		sort.Slice(sequences, func(i, j int) bool { return sequences[i] < sequences[j] })
		expected := uint64(1)
		for _, sequence := range sequences {
			if sequence != expected {
				return nil, fmt.Errorf("stored relay history for %q is missing sequence %d", originPeerID, expected)
			}
			expected++
		}
		out[originPeerID] = uint64(len(sequences))
	}
	return out, nil
}

func relayCursorMapsEqual(left map[string]uint64, right map[string]uint64) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		if right[key] != value {
			return false
		}
	}
	return true
}

func cloneKnownOrigins(in map[string]uint64) map[string]uint64 {
	out := map[string]uint64{}
	for key, value := range in {
		out[key] = value
	}
	return out
}

func openRelayStore(root string) (*Store, []OperationalEvent, []SignedKnowledgeItemRecord, []SignedKnowledgeApprovalRecord, []SignedKnowledgeEvidenceRecord, []SignedOperationalRunRecord, []SignedOperationalPlaceRecord, []SignedOperationalResourceRecord, []SignedKnowledgeLinkRecord, []SignedKnowledgeResponsibilityRecord, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "cas", "objects"), 0o755); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	openLog := func(path string) (*os.File, error) {
		return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	}

	eventPath := filepath.Join(root, "events.jsonl")
	eventsFile, err := openLog(eventPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	events, err := readEvents(eventsFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeItemPath := filepath.Join(root, "knowledge-item-messages.jsonl")
	knowledgeItemFile, err := openLog(knowledgeItemPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeItemRecords, err := readSignedKnowledgeItemRecords(knowledgeItemFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeApprovalPath := filepath.Join(root, "knowledge-approval-messages.jsonl")
	knowledgeApprovalFile, err := openLog(knowledgeApprovalPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeApprovalRecords, err := readSignedKnowledgeApprovalRecords(knowledgeApprovalFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeEvidencePath := filepath.Join(root, "knowledge-evidence-messages.jsonl")
	knowledgeEvidenceFile, err := openLog(knowledgeEvidencePath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeEvidenceRecords, err := readSignedKnowledgeEvidenceRecords(knowledgeEvidenceFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	operationalRunPath := filepath.Join(root, "operational-run-messages.jsonl")
	operationalRunFile, err := openLog(operationalRunPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	operationalRunRecords, err := readSignedOperationalRunRecords(operationalRunFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	operationalPlacePath := filepath.Join(root, "operational-place-messages.jsonl")
	operationalPlaceFile, err := openLog(operationalPlacePath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	operationalPlaceRecords, err := readSignedOperationalPlaceRecords(operationalPlaceFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	operationalResourcePath := filepath.Join(root, "operational-resource-messages.jsonl")
	operationalResourceFile, err := openLog(operationalResourcePath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	operationalResourceRecords, err := readSignedOperationalResourceRecords(operationalResourceFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeLinkPath := filepath.Join(root, "knowledge-link-messages.jsonl")
	knowledgeLinkFile, err := openLog(knowledgeLinkPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeLinkRecords, err := readSignedKnowledgeLinkRecords(knowledgeLinkFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeResponsibilityPath := filepath.Join(root, "knowledge-responsibility-messages.jsonl")
	knowledgeResponsibilityFile, err := openLog(knowledgeResponsibilityPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeResponsibilityRecords, err := readSignedKnowledgeResponsibilityRecords(knowledgeResponsibilityFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	store := &Store{
		root:                            root,
		events:                          eventsFile,
		knowledgeItemMessages:           knowledgeItemFile,
		knowledgeApprovalMessages:       knowledgeApprovalFile,
		knowledgeEvidenceMessages:       knowledgeEvidenceFile,
		operationalRunMessages:          operationalRunFile,
		operationalPlaceMessages:        operationalPlaceFile,
		operationalResourceMessages:     operationalResourceFile,
		knowledgeLinkMessages:           knowledgeLinkFile,
		knowledgeResponsibilityMessages: knowledgeResponsibilityFile,
		eventPath:                       eventPath,
		knowledgeItemPath:               knowledgeItemPath,
		knowledgeApprovalPath:           knowledgeApprovalPath,
		knowledgeEvidencePath:           knowledgeEvidencePath,
		operationalRunPath:              operationalRunPath,
		operationalPlacePath:            operationalPlacePath,
		operationalResourcePath:         operationalResourcePath,
		knowledgeLinkPath:               knowledgeLinkPath,
		knowledgeResponsibilityPath:     knowledgeResponsibilityPath,
		casRoot:                         filepath.Join(root, "cas", "objects"),
	}

	itemRecords, err := store.hydrateSignedKnowledgeItemRecords(knowledgeItemRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	approvalRecords, err := store.hydrateSignedKnowledgeApprovalRecords(knowledgeApprovalRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	evidenceRecords, err := store.hydrateSignedKnowledgeEvidenceRecords(knowledgeEvidenceRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	runRecords, err := store.hydrateSignedOperationalRunRecords(operationalRunRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	placeRecords, err := store.hydrateSignedOperationalPlaceRecords(operationalPlaceRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	resourceRecords, err := store.hydrateSignedOperationalResourceRecords(operationalResourceRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	linkRecords, err := store.hydrateSignedKnowledgeLinkRecords(knowledgeLinkRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	responsibilityRecords, err := store.hydrateSignedKnowledgeResponsibilityRecords(knowledgeResponsibilityRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return store, events, itemRecords, approvalRecords, evidenceRecords, runRecords, placeRecords, resourceRecords, linkRecords, responsibilityRecords, nil
}
