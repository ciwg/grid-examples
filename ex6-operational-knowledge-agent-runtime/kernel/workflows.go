package kernel

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
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

// WorkflowLifecycleEvent records one durable local workflow state transition.
// Intent: Keep local lifecycle authority reconstructible without deleting the
// content-addressed workflow artifact or its prior history. Source: DI-lovek
type WorkflowLifecycleEvent struct {
	Workflow Workflow `json:"workflow"`
}

// WorkflowRegistry projects append-only local lifecycle events into current state.
// Intent: Runtime restart must recover local workflow policy from durable events,
// while artifact identity remains in CAS. Source: DI-lovek
type WorkflowRegistry struct {
	path      string
	mu        sync.RWMutex
	workflows map[string]Workflow
}

func OpenWorkflowRegistry(stateRoot string) (*WorkflowRegistry, error) {
	if err := os.MkdirAll(stateRoot, 0o755); err != nil {
		return nil, err
	}
	registry := &WorkflowRegistry{
		path:      filepath.Join(stateRoot, "workflow-lifecycle.jsonl"),
		workflows: map[string]Workflow{},
	}
	if err := registry.replay(); err != nil {
		return nil, err
	}
	return registry, nil
}

func (registry *WorkflowRegistry) replay() (err error) {
	file, err := os.OpenFile(registry.path, os.O_CREATE|os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var event WorkflowLifecycleEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return err
		}
		if err := validateWorkflow(event.Workflow); err != nil {
			return err
		}
		registry.workflows[event.Workflow.ID] = event.Workflow
	}
	return scanner.Err()
}

func (registry *WorkflowRegistry) importWorkflow(workflow Workflow) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if _, exists := registry.workflows[workflow.ID]; exists {
		return errors.New("workflow is already imported")
	}
	workflow.State = WorkflowImported
	return registry.appendLocked(workflow)
}

func (registry *WorkflowRegistry) activateWorkflow(id string) error {
	return registry.transitionWorkflow(id, WorkflowActive)
}

func (registry *WorkflowRegistry) deactivateWorkflow(id string) error {
	return registry.transitionWorkflow(id, WorkflowDeactivated)
}

func (registry *WorkflowRegistry) revokeWorkflow(id string) error {
	return registry.transitionWorkflow(id, WorkflowRevoked)
}

func (registry *WorkflowRegistry) transitionWorkflow(id string, state WorkflowState) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	workflow, exists := registry.workflows[id]
	if !exists {
		return errors.New("workflow is not imported")
	}
	if workflow.State == WorkflowRevoked && state == WorkflowActive {
		return errors.New("revoked workflow cannot be activated")
	}
	workflow.State = state
	return registry.appendLocked(workflow)
}

func (registry *WorkflowRegistry) appendLocked(workflow Workflow) error {
	payload, err := json.Marshal(WorkflowLifecycleEvent{Workflow: workflow})
	if err != nil {
		return err
	}
	file, err := os.OpenFile(registry.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	if _, err := file.Write(append(payload, '\n')); err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	if err := file.Sync(); err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	registry.workflows[workflow.ID] = workflow
	return nil
}

func (registry *WorkflowRegistry) workflowsList() []Workflow {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	workflows := make([]Workflow, 0, len(registry.workflows))
	for _, workflow := range registry.workflows {
		workflows = append(workflows, workflow)
	}
	slices.SortFunc(workflows, func(left, right Workflow) int { return strings.Compare(left.ID, right.ID) })
	return workflows
}

func validateWorkflow(workflow Workflow) error {
	if strings.TrimSpace(workflow.ID) == "" || strings.TrimSpace(workflow.ArtifactCID) == "" {
		return errors.New("workflow ID and artifact CID are required")
	}
	switch workflow.State {
	case WorkflowImported, WorkflowActive, WorkflowDeactivated, WorkflowRevoked:
		return nil
	default:
		return errors.New("workflow state is invalid")
	}
}

func (runtime *Runtime) ImportWorkflow(workflow Workflow) error {
	if strings.TrimSpace(workflow.ID) == "" || strings.TrimSpace(workflow.ArtifactCID) == "" {
		return errors.New("workflow ID and artifact CID are required")
	}
	if _, err := runtime.GetCAS(workflow.ArtifactCID); err != nil {
		return err
	}
	return runtime.workflows.importWorkflow(workflow)
}

func (runtime *Runtime) ActivateWorkflow(id string) error {
	return runtime.workflows.activateWorkflow(id)
}

func (runtime *Runtime) DeactivateWorkflow(id string) error {
	return runtime.workflows.deactivateWorkflow(id)
}

func (runtime *Runtime) RevokeWorkflow(id string) error {
	return runtime.workflows.revokeWorkflow(id)
}

func (runtime *Runtime) Workflows() []Workflow {
	return runtime.workflows.workflowsList()
}
