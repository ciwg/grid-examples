package kernel

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

// WorkflowHandoffPolicy selects one locally approved next artifact for a
// completed output protocol. Intent: Keep automatic routing explicit, local,
// and removable instead of treating every compatible artifact as authority.
// Source: DI-lumek
type WorkflowHandoffPolicy struct {
	SourceWorkflow string `json:"source_workflow"`
	OutputPCID     string `json:"output_pcid"`
	TargetWorkflow string `json:"target_workflow"`
	InputPCID      string `json:"input_pcid"`
}

type WorkflowHandoffPolicies struct {
	path     string
	mu       sync.RWMutex
	policies []WorkflowHandoffPolicy
}

func OpenWorkflowHandoffPolicies(root string) (*WorkflowHandoffPolicies, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	policies := &WorkflowHandoffPolicies{path: filepath.Join(root, "workflow-handoff-policy.json")}
	if raw, err := os.ReadFile(policies.path); err == nil {
		if err := json.Unmarshal(raw, &policies.policies); err != nil {
			return nil, err
		}
		if err := validateWorkflowHandoffPolicySet(policies.policies); err != nil {
			return nil, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return policies, nil
}

func (policies *WorkflowHandoffPolicies) Set(policy WorkflowHandoffPolicy) error {
	if err := policy.Validate(); err != nil {
		return err
	}
	policies.mu.Lock()
	defer policies.mu.Unlock()
	for index, existing := range policies.policies {
		if existing.SourceWorkflow == policy.SourceWorkflow && existing.OutputPCID == policy.OutputPCID {
			policies.policies[index] = policy
			return policies.saveLocked()
		}
	}
	policies.policies = append(policies.policies, policy)
	return policies.saveLocked()
}
func (policies *WorkflowHandoffPolicies) Remove(source, output string) error {
	policies.mu.Lock()
	defer policies.mu.Unlock()
	for index, policy := range policies.policies {
		if policy.SourceWorkflow == source && policy.OutputPCID == output {
			policies.policies = append(policies.policies[:index], policies.policies[index+1:]...)
			return policies.saveLocked()
		}
	}
	return errors.New("workflow handoff policy is not found")
}
func (policies *WorkflowHandoffPolicies) Find(source, output string) (WorkflowHandoffPolicy, bool) {
	policies.mu.RLock()
	defer policies.mu.RUnlock()
	for _, policy := range policies.policies {
		if policy.SourceWorkflow == source && policy.OutputPCID == output {
			return policy, true
		}
	}
	return WorkflowHandoffPolicy{}, false
}
func (policies *WorkflowHandoffPolicies) List() []WorkflowHandoffPolicy {
	policies.mu.RLock()
	defer policies.mu.RUnlock()
	result := append([]WorkflowHandoffPolicy(nil), policies.policies...)
	slices.SortFunc(result, func(a, b WorkflowHandoffPolicy) int {
		return strings.Compare(a.SourceWorkflow+"\x00"+a.OutputPCID, b.SourceWorkflow+"\x00"+b.OutputPCID)
	})
	return result
}
func (policies *WorkflowHandoffPolicies) saveLocked() error {
	raw, err := json.MarshalIndent(policies.policies, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(policies.path, append(raw, '\n'), 0644)
}
func (policy WorkflowHandoffPolicy) Validate() error {
	if strings.TrimSpace(policy.SourceWorkflow) == "" || strings.TrimSpace(policy.OutputPCID) == "" || strings.TrimSpace(policy.TargetWorkflow) == "" || strings.TrimSpace(policy.InputPCID) == "" {
		return errors.New("workflow handoff policy fields are required")
	}
	if policy.SourceWorkflow == policy.TargetWorkflow {
		return errors.New("workflow handoff policy cannot target its source")
	}
	if err := validateWorkflowPCID(policy.OutputPCID); err != nil {
		return err
	}
	if err := validateWorkflowPCID(policy.InputPCID); err != nil {
		return err
	}
	return nil
}

func validateWorkflowHandoffPolicySet(policies []WorkflowHandoffPolicy) error {
	seen := map[string]bool{}
	for _, policy := range policies {
		if err := policy.Validate(); err != nil {
			return err
		}
		key := policy.SourceWorkflow + "\x00" + policy.OutputPCID
		if seen[key] {
			return errors.New("workflow handoff policy is duplicated")
		}
		seen[key] = true
		if workflowPolicyReaches(policies, policy.TargetWorkflow, policy.SourceWorkflow) {
			return errors.New("workflow handoff policy creates a cycle")
		}
	}
	return nil
}

func workflowPolicyReaches(policies []WorkflowHandoffPolicy, source, target string) bool {
	seen := map[string]bool{}
	var visit func(string) bool
	visit = func(current string) bool {
		if current == target {
			return true
		}
		if seen[current] {
			return false
		}
		seen[current] = true
		for _, policy := range policies {
			if policy.SourceWorkflow == current && visit(policy.TargetWorkflow) {
				return true
			}
		}
		return false
	}
	return visit(source)
}
