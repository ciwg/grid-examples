package kernel

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
)

// WorkflowTransfer is the exact-byte bundle carried by the dedicated workflow
// relay endpoint. Its lifecycle envelope is evidence of the sender's decision,
// not a receiver lifecycle instruction.
// Intent: Exchange independently verifiable workflow bytes without granting a
// remote runtime authority over local workflow availability. Source: DI-novuk
type WorkflowTransfer struct {
	ArtifactCID    string `json:"artifact_cid"`
	Artifact       []byte `json:"artifact"`
	LifecycleEvent []byte `json:"lifecycle_event"`
	Signature      string `json:"signature"`
}

func (transfer WorkflowTransfer) signingBytes() ([]byte, error) {
	unsigned := struct {
		ArtifactCID    string `json:"artifact_cid"`
		Artifact       []byte `json:"artifact"`
		LifecycleEvent []byte `json:"lifecycle_event"`
	}{transfer.ArtifactCID, transfer.Artifact, transfer.LifecycleEvent}
	return json.Marshal(unsigned)
}

// ExportWorkflowTransfer prepares the sender's currently accepted lifecycle
// evidence and exact CAS artifact for peer-authenticated transfer.
func (runtime *Runtime) ExportWorkflowTransfer(alias string) (WorkflowTransfer, error) {
	workflow, err := runtime.workflow(alias)
	if err != nil {
		return WorkflowTransfer{}, err
	}
	artifactCID, err := cid.Decode(workflow.ArtifactCID)
	if err != nil {
		return WorkflowTransfer{}, err
	}
	artifact, err := runtime.cas.GetCID(artifactCID)
	if err != nil {
		return WorkflowTransfer{}, err
	}
	head, ok := runtime.workflows.headCID(alias)
	if !ok {
		return WorkflowTransfer{}, fmt.Errorf("workflow lifecycle head is missing: %s", alias)
	}
	event, err := runtime.cas.GetCID(head)
	if err != nil {
		return WorkflowTransfer{}, err
	}
	transfer := WorkflowTransfer{ArtifactCID: artifactCID.String(), Artifact: artifact, LifecycleEvent: event}
	bytes, err := transfer.signingBytes()
	if err != nil {
		return WorkflowTransfer{}, err
	}
	transfer.Signature, err = runtime.peers.SignBytes(bytes)
	return transfer, err
}

// ImportWorkflowTransferFromPeer verifies and retains a peer's exact artifact
// and lifecycle evidence without changing this runtime's workflow lifecycle.
func (runtime *Runtime) ImportWorkflowTransferFromPeer(peerID string, transfer WorkflowTransfer) error {
	if strings.TrimSpace(peerID) == "" {
		return errors.New("peer ID is required")
	}
	if !runtime.peers.AllowsPush(peerID) {
		return fmt.Errorf("peer is not allowed to push workflow transfers: %s", peerID)
	}
	bytes, err := transfer.signingBytes()
	if err != nil {
		return err
	}
	if err := runtime.peers.VerifyPeerBytes(peerID, bytes, transfer.Signature); err != nil {
		return err
	}
	artifactCID, err := cid.Decode(transfer.ArtifactCID)
	if err != nil {
		return fmt.Errorf("workflow transfer artifact CID: %w", err)
	}
	event, err := DecodeWorkflowLifecycleEvent(transfer.LifecycleEvent)
	if err != nil {
		return fmt.Errorf("workflow transfer lifecycle event: %w", err)
	}
	if event.ArtifactCID != artifactCID {
		return errors.New("workflow transfer lifecycle event does not describe the artifact")
	}
	storedArtifactCID, err := runtime.cas.PutCID(transfer.Artifact)
	if err != nil {
		return err
	}
	if storedArtifactCID != artifactCID {
		return errors.New("workflow transfer artifact CID does not match artifact bytes")
	}
	// Intent: Retain remote lifecycle evidence outside the replayed local CAS so
	// a received statement cannot become a local lifecycle decision. Source: DI-novuk
	if _, err := runtime.workflowEvidence.PutCID(transfer.LifecycleEvent); err != nil {
		return err
	}
	return nil
}
