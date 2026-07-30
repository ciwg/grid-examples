package kernel

import (
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
	"github.com/ipfs/go-cid"
)

const testWorkflowArtifactCID = "bafkreihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku"

func TestWorkflowLifecycleEventRoundTrip(t *testing.T) {
	artifactCID, err := cid.Decode(testWorkflowArtifactCID)
	if err != nil {
		t.Fatalf("decode artifact CID: %v", err)
	}
	raw, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{
		State:         WorkflowImported,
		WorkflowAlias: "inventory-receipt",
		ArtifactCID:   artifactCID,
	})
	if err != nil {
		t.Fatalf("encode workflow lifecycle event: %v", err)
	}
	event, err := DecodeWorkflowLifecycleEvent(raw)
	if err != nil {
		t.Fatalf("decode workflow lifecycle event: %v", err)
	}
	if event.State != WorkflowImported {
		t.Errorf("state = %q, want %q", event.State, WorkflowImported)
	}
	if event.WorkflowAlias != "inventory-receipt" {
		t.Errorf("alias = %q, want inventory-receipt", event.WorkflowAlias)
	}
	if event.ArtifactCID != artifactCID {
		t.Errorf("artifact CID = %s, want %s", event.ArtifactCID, artifactCID)
	}
}

func TestWorkflowLifecycleEventRejectsInvalidParentCount(t *testing.T) {
	artifactCID, err := cid.Decode(testWorkflowArtifactCID)
	if err != nil {
		t.Fatalf("decode artifact CID: %v", err)
	}
	if _, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{
		State:         WorkflowImported,
		WorkflowAlias: "inventory-receipt",
		ArtifactCID:   artifactCID,
		Parents:       []cid.Cid{artifactCID},
	}); err == nil {
		t.Fatal("encode import with parent succeeded")
	}
}

func TestDecodeWorkflowLifecycleEventRejectsOtherProtocol(t *testing.T) {
	artifactCID, err := cid.Decode(testWorkflowArtifactCID)
	if err != nil {
		t.Fatalf("decode artifact CID: %v", err)
	}
	otherProtocol, err := cid.Decode("bafkreigh2akiscaildc2zk7wxx2z4dqm6qllq7rnyv4z3n3h4qtdwgqwdy")
	if err != nil {
		t.Fatalf("decode other protocol CID: %v", err)
	}
	raw, err := records.EncodeGrid(records.GridEnvelope{
		ProtocolPCID: otherProtocol,
		Slots:        []any{uint64(0), "inventory-receipt", artifactCID.Bytes(), []any{}},
	})
	if err != nil {
		t.Fatalf("encode other protocol event: %v", err)
	}
	if _, err := DecodeWorkflowLifecycleEvent(raw); err == nil {
		t.Fatal("decode other protocol event succeeded")
	}
}
