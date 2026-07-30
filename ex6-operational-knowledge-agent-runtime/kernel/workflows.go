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

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
	"github.com/ipfs/go-cid"
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

// WorkflowRegistry projects retained lifecycle envelopes into local availability.
// Intent: CAS events are authoritative; this cache is disposable. Source: DI-bavuk
type WorkflowRegistry struct {
	cachePath string
	cas       *store.CAS
	mu        sync.RWMutex
	workflows map[string]Workflow
	heads     map[string]cid.Cid
}

func OpenWorkflowRegistry(root string, cas *store.CAS) (*WorkflowRegistry, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	r := &WorkflowRegistry{cachePath: filepath.Join(root, "workflow-lifecycle-cache.json"), cas: cas, workflows: map[string]Workflow{}, heads: map[string]cid.Cid{}}
	if err := r.rebuild(); err != nil {
		return nil, err
	}
	return r, nil
}
func (r *WorkflowRegistry) rebuild() error {
	w, h, e := r.scan()
	if e != nil {
		return e
	}
	r.workflows, r.heads = w, h
	return r.cache()
}
func (r *WorkflowRegistry) scan() (map[string]Workflow, map[string]cid.Cid, error) {
	ids, e := r.cas.ListCIDs()
	if e != nil {
		return nil, nil, e
	}
	candidates := map[string]WorkflowLifecycleEvent{}
	for _, id := range ids {
		b, e := r.cas.GetCID(id)
		if e != nil {
			return nil, nil, e
		}
		event, e := DecodeWorkflowLifecycleEvent(b)
		if e == nil {
			candidates[id.String()] = event
		}
	}
	accepted := map[string]WorkflowLifecycleEvent{}
	for progress := true; progress; {
		progress = false
		for _, id := range ids {
			k := id.String()
			event, ok := candidates[k]
			if !ok {
				continue
			}
			if event.State == WorkflowImported {
				accepted[k] = event
				delete(candidates, k)
				progress = true
				continue
			}
			parent, ok := accepted[event.Parents[0].String()]
			if !ok {
				continue
			}
			delete(candidates, k)
			progress = true
			if parent.WorkflowAlias == event.WorkflowAlias && parent.ArtifactCID == event.ArtifactCID {
				accepted[k] = event
			}
		}
	}
	child := map[string]bool{}
	for _, event := range accepted {
		for _, parent := range event.Parents {
			child[parent.String()] = true
		}
	}
	w := map[string]Workflow{}
	h := map[string]cid.Cid{}
	for _, id := range ids {
		k := id.String()
		event, ok := accepted[k]
		if !ok || child[k] {
			continue
		}
		if _, ok := w[event.WorkflowAlias]; ok {
			return nil, nil, fmt.Errorf("workflow alias %q has competing lifecycle heads", event.WorkflowAlias)
		}
		w[event.WorkflowAlias] = Workflow{event.WorkflowAlias, event.ArtifactCID.String(), event.State}
		h[event.WorkflowAlias] = id
	}
	return w, h, nil
}
func (r *WorkflowRegistry) importWorkflow(w Workflow) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.workflows[w.ID]; ok {
		return errors.New("workflow is already imported")
	}
	a, e := cid.Decode(w.ArtifactCID)
	if e != nil {
		return e
	}
	w.State = WorkflowImported
	w.ArtifactCID = a.String()
	return r.append(w, a, nil)
}
func (r *WorkflowRegistry) activateWorkflow(id string) error { return r.transition(id, WorkflowActive) }
func (r *WorkflowRegistry) deactivateWorkflow(id string) error {
	return r.transition(id, WorkflowDeactivated)
}
func (r *WorkflowRegistry) revokeWorkflow(id string) error { return r.transition(id, WorkflowRevoked) }
func (r *WorkflowRegistry) transition(id string, s WorkflowState) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.workflows[id]
	if !ok {
		return errors.New("workflow is not imported")
	}
	if w.State == WorkflowRevoked && s == WorkflowActive {
		return errors.New("revoked workflow cannot be activated")
	}
	a, e := cid.Decode(w.ArtifactCID)
	if e != nil {
		return e
	}
	w.State = s
	return r.append(w, a, []cid.Cid{r.heads[id]})
}
func (r *WorkflowRegistry) append(w Workflow, a cid.Cid, p []cid.Cid) error {
	b, e := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{State: w.State, WorkflowAlias: w.ID, ArtifactCID: a, Parents: p})
	if e != nil {
		return e
	}
	id, e := r.cas.PutCID(b)
	if e != nil {
		return e
	}
	r.workflows[w.ID], r.heads[w.ID] = w, id
	return r.cache()
}
func (r *WorkflowRegistry) cache() error {
	b, e := json.Marshal(r.listLocked())
	if e != nil {
		return e
	}
	return os.WriteFile(r.cachePath, append(b, '\n'), 0644)
}
func (r *WorkflowRegistry) workflowsList() []Workflow {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.listLocked()
}
func (r *WorkflowRegistry) listLocked() []Workflow {
	w := make([]Workflow, 0, len(r.workflows))
	for _, x := range r.workflows {
		w = append(w, x)
	}
	slices.SortFunc(w, func(a, b Workflow) int { return strings.Compare(a.ID, b.ID) })
	return w
}
func (runtime *Runtime) ImportWorkflow(w Workflow) error {
	if strings.TrimSpace(w.ID) == "" || strings.TrimSpace(w.ArtifactCID) == "" {
		return errors.New("workflow ID and artifact CID are required")
	}
	a, e := runtime.workflowArtifactCID(w.ArtifactCID)
	if e != nil {
		return e
	}
	w.ArtifactCID = a.String()
	return runtime.workflows.importWorkflow(w)
}
func (runtime *Runtime) workflowArtifactCID(id string) (cid.Cid, error) {
	if c, e := cid.Decode(id); e == nil {
		if _, e = runtime.cas.GetCID(c); e != nil {
			return cid.Undef, e
		}
		return c, nil
	}
	b, e := runtime.GetCAS(id)
	if e != nil {
		return cid.Undef, e
	}
	return runtime.cas.PutCID(b)
}
func (runtime *Runtime) ActivateWorkflow(id string) error {
	return runtime.workflows.activateWorkflow(id)
}
func (runtime *Runtime) DeactivateWorkflow(id string) error {
	return runtime.workflows.deactivateWorkflow(id)
}
func (runtime *Runtime) RevokeWorkflow(id string) error { return runtime.workflows.revokeWorkflow(id) }
func (runtime *Runtime) Workflows() []Workflow          { return runtime.workflows.workflowsList() }
