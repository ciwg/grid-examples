package kernel

import (
	"errors"
	"slices"
	"strings"
)

type WorkflowState string

const (
	WorkflowImported    WorkflowState = "imported"
	WorkflowRevoked     WorkflowState = "revoked"
	WorkflowActive      WorkflowState = "active"
	WorkflowDeactivated WorkflowState = "deactivated"
)

type Workflow struct {
	ID          string        `json:"id"`
	ArtifactCID string        `json:"artifact_cid"`
	State       WorkflowState `json:"state"`
}

func (runtime *Runtime) ImportWorkflow(workflow Workflow) error {
	if strings.TrimSpace(workflow.ID) == "" || strings.TrimSpace(workflow.ArtifactCID) == "" {
		return errors.New("workflow ID and artifact CID are required")
	}
	if _, err := runtime.GetCAS(workflow.ArtifactCID); err != nil {
		return err
	}
	workflow.State = WorkflowImported
	runtime.workflows[workflow.ID] = workflow
	return nil
}

func (runtime *Runtime) ActivateWorkflow(id string) error {
	workflow, ok := runtime.workflows[id]
	if !ok {
		return errors.New("workflow is not imported")
	}
	workflow.State = WorkflowActive
	runtime.workflows[id] = workflow
	return nil
}

func (runtime *Runtime) DeactivateWorkflow(id string) error {
	workflow, ok := runtime.workflows[id]
	if !ok {
		return errors.New("workflow is not imported")
	}
	workflow.State = WorkflowDeactivated
	runtime.workflows[id] = workflow
	return nil
}

func (runtime *Runtime) Workflows() []Workflow {
	workflows := make([]Workflow, 0, len(runtime.workflows))
	for _, workflow := range runtime.workflows {
		workflows = append(workflows, workflow)
	}
	slices.SortFunc(workflows, func(left, right Workflow) int { return strings.Compare(left.ID, right.ID) })
	return workflows
}
