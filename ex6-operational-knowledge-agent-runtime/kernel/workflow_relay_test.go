package kernel

import (
	"bytes"
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
