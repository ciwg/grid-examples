package kernel

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
)

func TestWorkflowTransferRetainsEvidenceWithoutLocalLifecycleChange(t *testing.T) {
	alice, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open Alice: %v", err)
	}
	defer func() {
		if closeErr := alice.Close(); closeErr != nil {
			t.Errorf("close Alice: %v", closeErr)
		}
	}()
	bobRoot := t.TempDir()
	bob, err := Open(bobRoot)
	if err != nil {
		t.Fatalf("open Bob: %v", err)
	}
	defer func() {
		if closeErr := bob.Close(); closeErr != nil {
			t.Errorf("close Bob: %v", closeErr)
		}
	}()
	if err := bob.AllowPeer(grid.AllowedPeer{
		PeerID:            alice.LocalPeerID(),
		BatchURL:          "https://alice.invalid/relay/batch",
		ImportURL:         "https://alice.invalid/relay/import",
		PublicKey:         alice.LocalPeerPublicKey(),
		AllowPush:         true,
		AttesterClass:     "peer",
		AttestationWeight: 1,
		Federation:        "independent",
	}); err != nil {
		t.Fatalf("allow Alice: %v", err)
	}
	artifact, err := alice.cas.PutCID([]byte("workflow artifact bytes"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	if err := alice.ImportWorkflow(Workflow{ID: "handoff", ArtifactCID: artifact.String()}); err != nil {
		t.Fatalf("import Alice workflow: %v", err)
	}
	transfer, err := alice.ExportWorkflowTransfer("handoff")
	if err != nil {
		t.Fatalf("export transfer: %v", err)
	}
	if err := bob.ImportWorkflowTransferFromPeer(alice.LocalPeerID(), transfer); err != nil {
		t.Fatalf("import transfer: %v", err)
	}
	receivedArtifact, err := bob.cas.GetCID(artifact)
	if err != nil {
		t.Fatalf("get Bob artifact: %v", err)
	}
	if !bytes.Equal(receivedArtifact, []byte("workflow artifact bytes")) {
		t.Fatalf("Bob artifact = %q", receivedArtifact)
	}
	if workflows := bob.Workflows(); len(workflows) != 0 {
		t.Fatalf("received evidence changed Bob lifecycle: %#v", workflows)
	}
	inbox, err := bob.ScanWorkflowInbox()
	if err != nil {
		t.Fatalf("scan Bob inbox: %v", err)
	}
	if len(inbox) != 1 || inbox[0].ArtifactCID != artifact.String() || !inbox[0].ReadyToImport {
		t.Fatalf("Bob inbox = %#v", inbox)
	}
	if len(inbox[0].Evidence) != 1 || len(inbox[0].Evidence[0].PeerIDs) != 1 || inbox[0].Evidence[0].PeerIDs[0] != alice.LocalPeerID() || !inbox[0].Evidence[0].Valid {
		t.Fatalf("Bob inbox evidence = %#v", inbox[0].Evidence)
	}
	if closeErr := bob.Close(); closeErr != nil {
		t.Fatalf("close Bob before restart: %v", closeErr)
	}
	bob, err = Open(bobRoot)
	if err != nil {
		t.Fatalf("reopen Bob: %v", err)
	}
	inbox, err = bob.ScanWorkflowInbox()
	if err != nil {
		t.Fatalf("scan restarted Bob inbox: %v", err)
	}
	if len(inbox) != 1 || !inbox[0].ReadyToImport || len(inbox[0].Evidence[0].PeerIDs) != 1 || inbox[0].Evidence[0].PeerIDs[0] != alice.LocalPeerID() {
		t.Fatalf("restarted Bob inbox = %#v", inbox)
	}
	if err := bob.ImportWorkflowInbox(artifact.String(), "received-handoff"); err != nil {
		t.Fatalf("import Bob inbox artifact: %v", err)
	}
	if workflows := bob.Workflows(); len(workflows) != 1 || workflows[0].ID != "received-handoff" || workflows[0].State != WorkflowImported {
		t.Fatalf("imported Bob workflow = %#v", workflows)
	}
	inbox, err = bob.ScanWorkflowInbox()
	if err != nil || len(inbox) != 1 || !inbox[0].AlreadyImported {
		t.Fatalf("inbox after local import = %#v, %v", inbox, err)
	}
	evidence, err := store.OpenCAS(bobRoot + "/workflow-evidence")
	if err != nil {
		t.Fatalf("open evidence CAS: %v", err)
	}
	ids, err := evidence.ListCIDs()
	if err != nil {
		t.Fatalf("list evidence: %v", err)
	}
	if len(ids) != 1 {
		t.Fatalf("evidence object count = %d, want 1", len(ids))
	}
	receivedEvent, err := evidence.GetCID(ids[0])
	if err != nil {
		t.Fatalf("get evidence: %v", err)
	}
	if !bytes.Equal(receivedEvent, transfer.LifecycleEvent) {
		t.Fatal("received lifecycle evidence does not preserve exact bytes")
	}
}

func TestWorkflowInboxRefusesEvidenceWithoutReceiptMetadata(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer func() {
		if closeErr := runtime.Close(); closeErr != nil {
			t.Errorf("close runtime: %v", closeErr)
		}
	}()
	artifact, err := runtime.cas.PutCID([]byte("unattributed workflow artifact"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	evidence, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{
		State:         WorkflowImported,
		WorkflowAlias: "unattributed",
		ArtifactCID:   artifact,
	})
	if err != nil {
		t.Fatalf("encode evidence: %v", err)
	}
	if _, err := runtime.workflowEvidence.PutCID(evidence); err != nil {
		t.Fatalf("store evidence: %v", err)
	}
	inbox, err := runtime.ScanWorkflowInbox()
	if err != nil {
		t.Fatalf("scan inbox: %v", err)
	}
	if len(inbox) != 1 || inbox[0].ReadyToImport || len(inbox[0].Evidence) != 1 || inbox[0].Evidence[0].Reason != "receipt metadata is missing" {
		t.Fatalf("unattributed inbox = %#v", inbox)
	}
	if err := runtime.ImportWorkflowInbox(artifact.String(), "unattributed"); err == nil {
		t.Fatal("expected unattributed inbox artifact to be rejected")
	}
}

func TestWorkflowInboxGroupsIdenticalEvidenceFromMultiplePeers(t *testing.T) {
	alice, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open Alice: %v", err)
	}
	defer func() {
		if closeErr := alice.Close(); closeErr != nil {
			t.Errorf("close Alice: %v", closeErr)
		}
	}()
	carol, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open Carol: %v", err)
	}
	defer func() {
		if closeErr := carol.Close(); closeErr != nil {
			t.Errorf("close Carol: %v", closeErr)
		}
	}()
	bob, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open Bob: %v", err)
	}
	defer func() {
		if closeErr := bob.Close(); closeErr != nil {
			t.Errorf("close Bob: %v", closeErr)
		}
	}()
	for _, peer := range []*Runtime{alice, carol} {
		if err := bob.AllowPeer(grid.AllowedPeer{
			PeerID:            peer.LocalPeerID(),
			BatchURL:          "https://peer.invalid/relay/batch",
			ImportURL:         "https://peer.invalid/relay/import",
			PublicKey:         peer.LocalPeerPublicKey(),
			AllowPush:         true,
			AttesterClass:     "peer",
			AttestationWeight: 1,
			Federation:        "independent",
		}); err != nil {
			t.Fatalf("allow peer: %v", err)
		}
	}
	artifact, err := alice.cas.PutCID([]byte("shared artifact"))
	if err != nil {
		t.Fatalf("store Alice artifact: %v", err)
	}
	if _, err := carol.cas.PutCID([]byte("shared artifact")); err != nil {
		t.Fatalf("store Carol artifact: %v", err)
	}
	for _, peer := range []*Runtime{alice, carol} {
		if err := peer.ImportWorkflow(Workflow{ID: "shared", ArtifactCID: artifact.String()}); err != nil {
			t.Fatalf("import peer workflow: %v", err)
		}
	}
	for _, peer := range []*Runtime{alice, carol} {
		transfer, err := peer.ExportWorkflowTransfer("shared")
		if err != nil {
			t.Fatalf("export peer transfer: %v", err)
		}
		if err := bob.ImportWorkflowTransferFromPeer(peer.LocalPeerID(), transfer); err != nil {
			t.Fatalf("receive peer transfer: %v", err)
		}
	}
	inbox, err := bob.ScanWorkflowInbox()
	if err != nil {
		t.Fatalf("scan inbox: %v", err)
	}
	if len(inbox) != 1 || len(inbox[0].Evidence) != 1 || len(inbox[0].Evidence[0].PeerIDs) != 2 {
		t.Fatalf("multi-peer inbox = %#v", inbox)
	}
}

func TestWorkflowInboxReportsCorruptEvidenceAsNotReady(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer func() {
		if closeErr := runtime.Close(); closeErr != nil {
			t.Errorf("close runtime: %v", closeErr)
		}
	}()
	artifact, err := runtime.cas.PutCID([]byte("artifact with corrupt evidence"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	evidence, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{State: WorkflowImported, WorkflowAlias: "corrupt", ArtifactCID: artifact})
	if err != nil {
		t.Fatalf("encode evidence: %v", err)
	}
	evidenceCID, err := runtime.workflowEvidence.PutCID(evidence)
	if err != nil {
		t.Fatalf("store evidence: %v", err)
	}
	if err := runtime.workflowReceipts.record(WorkflowReceipt{ArtifactCID: artifact.String(), EvidenceCID: evidenceCID.String(), PeerIDs: []string{"peer-alice"}}); err != nil {
		t.Fatalf("record receipt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runtime.root, "workflow-evidence", evidenceCID.String()), []byte("corrupt"), 0o644); err != nil {
		t.Fatalf("corrupt evidence: %v", err)
	}
	inbox, err := runtime.ScanWorkflowInbox()
	if err != nil {
		t.Fatalf("scan inbox: %v", err)
	}
	if len(inbox) != 1 || inbox[0].ReadyToImport || inbox[0].Evidence[0].Reason != "lifecycle evidence is unavailable" {
		t.Fatalf("corrupt evidence inbox = %#v", inbox)
	}
}

func TestRuntimeOpensWithCorruptReceiptMetadata(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	artifact, err := runtime.cas.PutCID([]byte("artifact with unavailable receipt metadata"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	evidence, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{State: WorkflowImported, WorkflowAlias: "metadata", ArtifactCID: artifact})
	if err != nil {
		t.Fatalf("encode evidence: %v", err)
	}
	if _, err := runtime.workflowEvidence.PutCID(evidence); err != nil {
		t.Fatalf("store evidence: %v", err)
	}
	if closeErr := runtime.Close(); closeErr != nil {
		t.Fatalf("close runtime: %v", closeErr)
	}
	if err := os.WriteFile(filepath.Join(root, "state", "workflow-receipts.json"), []byte("not JSON"), 0o644); err != nil {
		t.Fatalf("corrupt receipt metadata: %v", err)
	}
	runtime, err = Open(root)
	if err != nil {
		t.Fatalf("reopen with corrupt receipt metadata: %v", err)
	}
	defer func() {
		if closeErr := runtime.Close(); closeErr != nil {
			t.Errorf("close reopened runtime: %v", closeErr)
		}
	}()
	inbox, err := runtime.ScanWorkflowInbox()
	if err != nil {
		t.Fatalf("scan inbox: %v", err)
	}
	if len(inbox) != 1 || inbox[0].ReadyToImport || inbox[0].Evidence[0].Reason != "receipt metadata is missing" {
		t.Fatalf("inbox with corrupt receipt metadata = %#v", inbox)
	}
}

func TestWorkflowInboxReadsLegacySinglePeerReceiptMetadata(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	artifact, err := runtime.cas.PutCID([]byte("artifact with legacy receipt metadata"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	evidence, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{State: WorkflowImported, WorkflowAlias: "legacy", ArtifactCID: artifact})
	if err != nil {
		t.Fatalf("encode evidence: %v", err)
	}
	evidenceCID, err := runtime.workflowEvidence.PutCID(evidence)
	if err != nil {
		t.Fatalf("store evidence: %v", err)
	}
	if closeErr := runtime.Close(); closeErr != nil {
		t.Fatalf("close runtime: %v", closeErr)
	}
	legacyBody, err := json.Marshal([]map[string]string{{
		"artifact_cid": artifact.String(),
		"evidence_cid": evidenceCID.String(),
		"peer_id":      "peer-alice",
	}})
	if err != nil {
		t.Fatalf("marshal legacy receipt metadata: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "state", "workflow-receipts.json"), append(legacyBody, '\n'), 0o644); err != nil {
		t.Fatalf("write legacy receipt metadata: %v", err)
	}
	runtime, err = Open(root)
	if err != nil {
		t.Fatalf("open legacy receipt metadata: %v", err)
	}
	defer func() {
		if closeErr := runtime.Close(); closeErr != nil {
			t.Errorf("close reopened runtime: %v", closeErr)
		}
	}()
	inbox, err := runtime.ScanWorkflowInbox()
	if err != nil {
		t.Fatalf("scan inbox: %v", err)
	}
	if len(inbox) != 1 || !inbox[0].ReadyToImport || len(inbox[0].Evidence) != 1 || len(inbox[0].Evidence[0].PeerIDs) != 1 || inbox[0].Evidence[0].PeerIDs[0] != "peer-alice" {
		t.Fatalf("legacy receipt inbox = %#v", inbox)
	}
}

func TestWorkflowReceiptPeerUpdateRollsBackAfterPersistenceFailure(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer func() {
		if closeErr := runtime.Close(); closeErr != nil {
			t.Errorf("close runtime: %v", closeErr)
		}
	}()
	artifact, err := runtime.cas.PutCID([]byte("receipt persistence rollback artifact"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	evidence, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{State: WorkflowImported, WorkflowAlias: "rollback", ArtifactCID: artifact})
	if err != nil {
		t.Fatalf("encode evidence: %v", err)
	}
	evidenceCID, err := runtime.workflowEvidence.PutCID(evidence)
	if err != nil {
		t.Fatalf("store evidence: %v", err)
	}
	receipt := WorkflowReceipt{ArtifactCID: artifact.String(), EvidenceCID: evidenceCID.String(), PeerIDs: []string{"peer-alice"}}
	if err := runtime.workflowReceipts.record(receipt); err != nil {
		t.Fatalf("record first peer: %v", err)
	}
	persistedPath := runtime.workflowReceipts.path
	runtime.workflowReceipts.path = t.TempDir()
	if err := runtime.workflowReceipts.record(WorkflowReceipt{ArtifactCID: artifact.String(), EvidenceCID: evidenceCID.String(), PeerIDs: []string{"peer-carol"}}); err == nil {
		t.Fatal("expected receipt persistence failure")
	}
	receipts := runtime.workflowReceipts.list()
	if len(receipts) != 1 || len(receipts[0].PeerIDs) != 1 || receipts[0].PeerIDs[0] != "peer-alice" {
		t.Fatalf("receipt after failed persistence = %#v", receipts)
	}
	runtime.workflowReceipts.path = persistedPath
	if err := runtime.workflowReceipts.record(WorkflowReceipt{ArtifactCID: artifact.String(), EvidenceCID: evidenceCID.String(), PeerIDs: []string{"peer-carol"}}); err != nil {
		t.Fatalf("retry receipt persistence: %v", err)
	}
	receipts = runtime.workflowReceipts.list()
	if len(receipts) != 1 || len(receipts[0].PeerIDs) != 2 {
		t.Fatalf("receipt after retry = %#v", receipts)
	}
}
