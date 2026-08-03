package kernel

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/ipfs/go-cid"
)

// WorkflowReceipt is local metadata for a peer-authenticated workflow
// transfer. The exact lifecycle envelope remains unchanged in workflow-evidence
// CAS; this sidecar records the authenticated peer identities that the envelope
// itself does not carry.
// Intent: Preserve receipt provenance without treating local metadata as
// lifecycle authority. Source: DI-rufir
type WorkflowReceipt struct {
	ArtifactCID string   `json:"artifact_cid"`
	EvidenceCID string   `json:"evidence_cid"`
	PeerIDs     []string `json:"peer_ids"`
}

// WorkflowInboxEvidence describes one retained lifecycle statement for an
// artifact. Valid is derived during each scan, never trusted from the sidecar.
type WorkflowInboxEvidence struct {
	EvidenceCID string   `json:"evidence_cid"`
	PeerIDs     []string `json:"peer_ids,omitempty"`
	Valid       bool     `json:"valid"`
	Reason      string   `json:"reason,omitempty"`
}

// WorkflowInboxEntry is the operator-facing view of one received artifact.
// Receipt evidence and artifact availability are deliberately reported apart
// from local workflow lifecycle state.
type WorkflowInboxEntry struct {
	ArtifactCID       string                  `json:"artifact_cid"`
	ArtifactAvailable bool                    `json:"artifact_available"`
	Evidence          []WorkflowInboxEvidence `json:"evidence"`
	AlreadyImported   bool                    `json:"already_imported"`
	ReadyToImport     bool                    `json:"ready_to_import"`
	Reason            string                  `json:"reason,omitempty"`
}

type workflowReceiptStore struct {
	path     string
	mu       sync.RWMutex
	receipts map[string]WorkflowReceipt
	loadErr  error
}

type workflowReceiptDisk struct {
	ArtifactCID string   `json:"artifact_cid"`
	EvidenceCID string   `json:"evidence_cid"`
	PeerID      string   `json:"peer_id,omitempty"`
	PeerIDs     []string `json:"peer_ids,omitempty"`
}

func openWorkflowReceiptStore(path string) (*workflowReceiptStore, error) {
	store := &workflowReceiptStore{path: path, receipts: map[string]WorkflowReceipt{}}
	body, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}
	var diskReceipts []workflowReceiptDisk
	if err := json.Unmarshal(body, &diskReceipts); err != nil {
		// Intent: A corrupt local provenance projection must not hide retained
		// evidence or prevent the runtime from serving unrelated work. Source: DI-rufir
		store.loadErr = fmt.Errorf("decode workflow receipt metadata: %w", err)
		return store, nil
	}
	for _, diskReceipt := range diskReceipts {
		receipt := WorkflowReceipt{
			ArtifactCID: diskReceipt.ArtifactCID,
			EvidenceCID: diskReceipt.EvidenceCID,
			PeerIDs:     append([]string{}, diskReceipt.PeerIDs...),
		}
		if len(receipt.PeerIDs) == 0 && strings.TrimSpace(diskReceipt.PeerID) != "" {
			receipt.PeerIDs = []string{diskReceipt.PeerID}
		}
		receipt.normalizePeers()
		if err := receipt.validate(); err != nil {
			store.loadErr = fmt.Errorf("validate workflow receipt metadata: %w", err)
			return store, nil
		}
		if _, duplicate := store.receipts[receipt.EvidenceCID]; duplicate {
			store.loadErr = fmt.Errorf("duplicate workflow receipt evidence CID: %s", receipt.EvidenceCID)
			return store, nil
		}
		store.receipts[receipt.EvidenceCID] = receipt
	}
	return store, nil
}

func (store *workflowReceiptStore) record(receipt WorkflowReceipt) error {
	if err := receipt.validate(); err != nil {
		return err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.loadErr != nil {
		return fmt.Errorf("workflow receipt metadata is unavailable: %w", store.loadErr)
	}
	if existing, ok := store.receipts[receipt.EvidenceCID]; ok {
		if existing.ArtifactCID != receipt.ArtifactCID {
			return fmt.Errorf("workflow receipt evidence CID has conflicting artifact CID: %s", receipt.EvidenceCID)
		}
		updated := WorkflowReceipt{
			ArtifactCID: existing.ArtifactCID,
			EvidenceCID: existing.EvidenceCID,
			PeerIDs:     append([]string{}, existing.PeerIDs...),
		}
		if updated.addPeers(receipt.PeerIDs) {
			store.receipts[receipt.EvidenceCID] = updated
			if err := store.persistLocked(); err != nil {
				// Intent: A failed sidecar update must not make a retry believe new
				// provenance was durable when only the prior receipt was retained.
				// Source: DI-rufir
				store.receipts[receipt.EvidenceCID] = existing
				return err
			}
		}
		return nil
	}
	store.receipts[receipt.EvidenceCID] = receipt
	if err := store.persistLocked(); err != nil {
		delete(store.receipts, receipt.EvidenceCID)
		return err
	}
	return nil
}

func (store *workflowReceiptStore) list() []WorkflowReceipt {
	store.mu.RLock()
	defer store.mu.RUnlock()
	receipts := make([]WorkflowReceipt, 0, len(store.receipts))
	for _, receipt := range store.receipts {
		receipts = append(receipts, receipt)
	}
	slices.SortFunc(receipts, func(left, right WorkflowReceipt) int {
		return strings.Compare(left.EvidenceCID, right.EvidenceCID)
	})
	return receipts
}

func (store *workflowReceiptStore) persistLocked() error {
	receipts := make([]workflowReceiptDisk, 0, len(store.receipts))
	for _, receipt := range store.receipts {
		receipts = append(receipts, workflowReceiptDisk{
			ArtifactCID: receipt.ArtifactCID,
			EvidenceCID: receipt.EvidenceCID,
			PeerIDs:     receipt.PeerIDs,
		})
	}
	slices.SortFunc(receipts, func(left, right workflowReceiptDisk) int {
		return strings.Compare(left.EvidenceCID, right.EvidenceCID)
	})
	body, err := json.MarshalIndent(receipts, "", "  ")
	if err != nil {
		return err
	}
	return writeWorkflowReceiptFile(store.path, append(body, '\n'))
}

// writeWorkflowReceiptFile replaces metadata atomically after syncing the new
// bytes, so an interrupted receipt update leaves either the prior complete
// sidecar or the new complete sidecar.
// Intent: Keep local provenance metadata from becoming a recurring relay
// availability failure after a crash. Source: DI-rufir
func writeWorkflowReceiptFile(path string, body []byte) (result error) {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".workflow-receipts-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() {
		if temporaryPath == "" {
			return
		}
		if err := os.Remove(temporaryPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			result = errors.Join(result, err)
		}
	}()
	if err := temporary.Chmod(0o644); err != nil {
		if closeErr := temporary.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	if _, err := temporary.Write(body); err != nil {
		if closeErr := temporary.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	if err := temporary.Sync(); err != nil {
		if closeErr := temporary.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return err
	}
	temporaryPath = ""
	directory, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	if err := directory.Sync(); err != nil {
		if closeErr := directory.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	return directory.Close()
}

func (receipt WorkflowReceipt) validate() error {
	if len(receipt.PeerIDs) == 0 {
		return errors.New("workflow receipt peer IDs are required")
	}
	for _, peerID := range receipt.PeerIDs {
		if strings.TrimSpace(peerID) == "" {
			return errors.New("workflow receipt peer ID is required")
		}
	}
	if _, err := cid.Decode(receipt.ArtifactCID); err != nil {
		return fmt.Errorf("workflow receipt artifact CID: %w", err)
	}
	if _, err := cid.Decode(receipt.EvidenceCID); err != nil {
		return fmt.Errorf("workflow receipt evidence CID: %w", err)
	}
	return nil
}

func (receipt *WorkflowReceipt) addPeers(peerIDs []string) bool {
	known := map[string]bool{}
	for _, peerID := range receipt.PeerIDs {
		known[peerID] = true
	}
	changed := false
	for _, peerID := range peerIDs {
		if !known[peerID] {
			receipt.PeerIDs = append(receipt.PeerIDs, peerID)
			known[peerID] = true
			changed = true
		}
	}
	receipt.normalizePeers()
	return changed
}

func (receipt *WorkflowReceipt) normalizePeers() {
	known := map[string]bool{}
	normalized := make([]string, 0, len(receipt.PeerIDs))
	for _, peerID := range receipt.PeerIDs {
		if !known[peerID] {
			normalized = append(normalized, peerID)
			known[peerID] = true
		}
	}
	slices.Sort(normalized)
	receipt.PeerIDs = normalized
}

// ScanWorkflowInbox derives the receipt inbox from retained evidence on every
// call. The sidecar contributes only sender provenance; it cannot make an
// artifact importable unless the corresponding raw evidence validates again.
// Intent: Keep the inbox a CAS-derived staging view rather than another
// lifecycle ledger. Source: DI-jifuk; DI-rufir
func (runtime *Runtime) ScanWorkflowInbox() ([]WorkflowInboxEntry, error) {
	entries := map[string]*WorkflowInboxEntry{}
	receipts := runtime.workflowReceipts.list()
	for _, receipt := range receipts {
		entry := workflowInboxEntry(entries, receipt.ArtifactCID)
		entry.Evidence = append(entry.Evidence, runtime.scanWorkflowReceipt(receipt))
	}
	ids, err := runtime.workflowEvidence.ListCIDs()
	if err != nil {
		return nil, err
	}
	known := map[string]bool{}
	for _, receipt := range receipts {
		known[receipt.EvidenceCID] = true
	}
	for _, evidenceCID := range ids {
		if known[evidenceCID.String()] {
			continue
		}
		raw, err := runtime.workflowEvidence.GetCID(evidenceCID)
		if err != nil {
			continue
		}
		event, err := DecodeWorkflowLifecycleEvent(raw)
		if err != nil {
			continue
		}
		entry := workflowInboxEntry(entries, event.ArtifactCID.String())
		entry.Evidence = append(entry.Evidence, WorkflowInboxEvidence{
			EvidenceCID: evidenceCID.String(),
			Reason:      "receipt metadata is missing",
		})
	}
	result := make([]WorkflowInboxEntry, 0, len(entries))
	for _, entry := range entries {
		artifactCID, err := cid.Decode(entry.ArtifactCID)
		if err == nil {
			_, getErr := runtime.cas.GetCID(artifactCID)
			entry.ArtifactAvailable = getErr == nil
		}
		for _, evidence := range entry.Evidence {
			if evidence.Valid {
				entry.ReadyToImport = entry.ArtifactAvailable
				break
			}
		}
		for _, workflow := range runtime.Workflows() {
			if workflow.ArtifactCID == entry.ArtifactCID {
				entry.AlreadyImported = true
				break
			}
		}
		if !entry.ReadyToImport {
			switch {
			case !entry.ArtifactAvailable:
				entry.Reason = "artifact bytes are unavailable"
			case len(entry.Evidence) == 0:
				entry.Reason = "matching lifecycle evidence is unavailable"
			default:
				entry.Reason = "no valid peer-authenticated lifecycle evidence"
			}
		}
		slices.SortFunc(entry.Evidence, func(left, right WorkflowInboxEvidence) int {
			return strings.Compare(left.EvidenceCID, right.EvidenceCID)
		})
		result = append(result, *entry)
	}
	slices.SortFunc(result, func(left, right WorkflowInboxEntry) int {
		return strings.Compare(left.ArtifactCID, right.ArtifactCID)
	})
	return result, nil
}

func workflowInboxEntry(entries map[string]*WorkflowInboxEntry, artifactCID string) *WorkflowInboxEntry {
	if entry, ok := entries[artifactCID]; ok {
		return entry
	}
	entry := &WorkflowInboxEntry{ArtifactCID: artifactCID, Evidence: []WorkflowInboxEvidence{}}
	entries[artifactCID] = entry
	return entry
}

func (runtime *Runtime) scanWorkflowReceipt(receipt WorkflowReceipt) WorkflowInboxEvidence {
	evidence := WorkflowInboxEvidence{EvidenceCID: receipt.EvidenceCID, PeerIDs: receipt.PeerIDs}
	evidenceCID, err := cid.Decode(receipt.EvidenceCID)
	if err != nil {
		evidence.Reason = "receipt metadata has an invalid evidence CID"
		return evidence
	}
	raw, err := runtime.workflowEvidence.GetCID(evidenceCID)
	if err != nil {
		evidence.Reason = "lifecycle evidence is unavailable"
		return evidence
	}
	event, err := DecodeWorkflowLifecycleEvent(raw)
	if err != nil {
		evidence.Reason = "lifecycle evidence is invalid"
		return evidence
	}
	if event.ArtifactCID.String() != receipt.ArtifactCID {
		evidence.Reason = "lifecycle evidence describes a different artifact"
		return evidence
	}
	evidence.Valid = true
	return evidence
}

// InspectWorkflowInbox returns one CAS-derived received artifact entry.
func (runtime *Runtime) InspectWorkflowInbox(artifactCID string) (WorkflowInboxEntry, error) {
	canonical, err := cid.Decode(artifactCID)
	if err != nil {
		return WorkflowInboxEntry{}, err
	}
	entries, err := runtime.ScanWorkflowInbox()
	if err != nil {
		return WorkflowInboxEntry{}, err
	}
	for _, entry := range entries {
		if entry.ArtifactCID == canonical.String() {
			return entry, nil
		}
	}
	return WorkflowInboxEntry{}, errors.New("received workflow artifact is not in the inbox")
}

// ImportWorkflowInbox creates an explicit local lifecycle import only after a
// scan confirms both retained artifact bytes and peer-authenticated evidence.
// Intent: Require a local operator alias before received bytes gain local
// lifecycle meaning. Source: DI-jifuk
func (runtime *Runtime) ImportWorkflowInbox(artifactCID string, alias string) error {
	entry, err := runtime.InspectWorkflowInbox(artifactCID)
	if err != nil {
		return err
	}
	if !entry.ReadyToImport {
		return fmt.Errorf("received workflow artifact is not ready to import: %s", entry.Reason)
	}
	return runtime.ImportWorkflow(Workflow{ID: alias, ArtifactCID: entry.ArtifactCID})
}

func workflowReceiptPath(root string) string {
	return filepath.Join(root, "state", "workflow-receipts.json")
}
