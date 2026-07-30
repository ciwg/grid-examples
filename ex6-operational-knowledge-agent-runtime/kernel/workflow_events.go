package kernel

import (
	"errors"
	"fmt"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
	"github.com/ipfs/go-cid"
)

const WorkflowLifecycleProtocolPCID = "bafkreicuygjo7udgzvopsv6bsvx5vrcbwo22bhd2izmkzuyaucvts7ignq"

var workflowLifecycleProtocolCID = cid.MustParse(WorkflowLifecycleProtocolPCID)

// WorkflowLifecycleEvent is a pCID-selected local lifecycle decision retained
// as exact canonical `grid()` bytes in CAS.
// Intent: Keep artifact-scoped lifecycle history independently replayable from
// mutable registry projections. Source: DI-bavuk
type WorkflowLifecycleEvent struct {
	State         WorkflowState
	WorkflowAlias string
	ArtifactCID   cid.Cid
	Parents       []cid.Cid
}

func EncodeWorkflowLifecycleEvent(event WorkflowLifecycleEvent) ([]byte, error) {
	if err := event.Validate(); err != nil {
		return nil, err
	}
	parents := make([]any, 0, len(event.Parents))
	for _, parent := range event.Parents {
		parents = append(parents, parent.Bytes())
	}
	return records.EncodeGrid(records.GridEnvelope{
		ProtocolPCID: workflowLifecycleProtocolCID,
		Slots: []any{
			workflowOperation(event.State),
			event.WorkflowAlias,
			event.ArtifactCID.Bytes(),
			parents,
		},
	})
}

func DecodeWorkflowLifecycleEvent(raw []byte) (WorkflowLifecycleEvent, error) {
	envelope, err := records.DecodeGrid(raw)
	if err != nil {
		return WorkflowLifecycleEvent{}, err
	}
	if envelope.ProtocolPCID != workflowLifecycleProtocolCID {
		return WorkflowLifecycleEvent{}, errors.New("grid envelope does not select the workflow lifecycle protocol")
	}
	if len(envelope.Slots) != 4 {
		return WorkflowLifecycleEvent{}, errors.New("workflow lifecycle envelope must contain four protocol slots")
	}
	operation, ok := envelope.Slots[0].(uint64)
	if !ok {
		return WorkflowLifecycleEvent{}, errors.New("workflow lifecycle operation must be an unsigned integer")
	}
	state, err := workflowState(operation)
	if err != nil {
		return WorkflowLifecycleEvent{}, err
	}
	alias, ok := envelope.Slots[1].(string)
	if !ok {
		return WorkflowLifecycleEvent{}, errors.New("workflow lifecycle alias must be text")
	}
	artifactBytes, ok := envelope.Slots[2].([]byte)
	if !ok {
		return WorkflowLifecycleEvent{}, errors.New("workflow lifecycle artifact must be CID bytes")
	}
	artifactCID, err := cid.Cast(artifactBytes)
	if err != nil {
		return WorkflowLifecycleEvent{}, fmt.Errorf("workflow lifecycle artifact CID: %w", err)
	}
	parentValues, ok := envelope.Slots[3].([]any)
	if !ok {
		return WorkflowLifecycleEvent{}, errors.New("workflow lifecycle parents must be an array")
	}
	parents := make([]cid.Cid, 0, len(parentValues))
	for _, parentValue := range parentValues {
		parentBytes, ok := parentValue.([]byte)
		if !ok {
			return WorkflowLifecycleEvent{}, errors.New("workflow lifecycle parent must be CID bytes")
		}
		parentCID, err := cid.Cast(parentBytes)
		if err != nil {
			return WorkflowLifecycleEvent{}, fmt.Errorf("workflow lifecycle parent CID: %w", err)
		}
		parents = append(parents, parentCID)
	}
	event := WorkflowLifecycleEvent{
		State:         state,
		WorkflowAlias: alias,
		ArtifactCID:   artifactCID,
		Parents:       parents,
	}
	if err := event.Validate(); err != nil {
		return WorkflowLifecycleEvent{}, err
	}
	return event, nil
}

func (event WorkflowLifecycleEvent) Validate() error {
	if strings.TrimSpace(event.WorkflowAlias) == "" {
		return errors.New("workflow lifecycle alias is required")
	}
	if event.ArtifactCID.Version() != 1 {
		return errors.New("workflow lifecycle artifact must be CIDv1")
	}
	for _, parent := range event.Parents {
		if parent.Version() != 1 {
			return errors.New("workflow lifecycle parent must be CIDv1")
		}
	}
	switch event.State {
	case WorkflowImported:
		if len(event.Parents) != 0 {
			return errors.New("workflow import must not have parents")
		}
	case WorkflowActive, WorkflowDeactivated, WorkflowRevoked:
		if len(event.Parents) != 1 {
			return errors.New("workflow lifecycle transition must have one parent")
		}
	default:
		return errors.New("workflow lifecycle state is invalid")
	}
	return nil
}

func workflowOperation(state WorkflowState) uint64 {
	switch state {
	case WorkflowImported:
		return 0
	case WorkflowActive:
		return 1
	case WorkflowDeactivated:
		return 2
	case WorkflowRevoked:
		return 3
	default:
		return 255
	}
}

func workflowState(operation uint64) (WorkflowState, error) {
	switch operation {
	case 0:
		return WorkflowImported, nil
	case 1:
		return WorkflowActive, nil
	case 2:
		return WorkflowDeactivated, nil
	case 3:
		return WorkflowRevoked, nil
	default:
		return "", errors.New("workflow lifecycle operation is invalid")
	}
}
