package kernel

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
	"github.com/ipfs/go-cid"
)

const WorkflowHandoffProtocolPCID = "bafkreiahdp34nto2rnnqde26jw3xnkd6xnlalnr72sug3w7tjb3bhhoj4q"
const WorkflowRunProtocolPCID = "bafkreifmttp5fwt3yvxvkb7ni6kwg3j3arl7mbjsyzszf7s7crxrncch24"

var workflowHandoffProtocolCID = cid.MustParse(WorkflowHandoffProtocolPCID)
var workflowRunProtocolCID = cid.MustParse(WorkflowRunProtocolPCID)

type WorkflowHandoff struct {
	PCID   string            `json:"pcid"`
	Values map[string]string `json:"values"`
}

// EncodeWorkflowHandoff encodes a deterministic pCID-selected handoff.
// Intent: Make workflow inputs and outputs portable binary contracts rather
// than parser-dependent command text. Source: DI-lumek
func EncodeWorkflowHandoff(handoff WorkflowHandoff) ([]byte, error) {
	if handoff.PCID == "" || len(handoff.Values) == 0 {
		return nil, errors.New("workflow handoff pCID and values are required")
	}
	if err := validateWorkflowPCID(handoff.PCID); err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(handoff.Values))
	for key, value := range handoff.Values {
		if strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return nil, errors.New("workflow handoff keys and values are required")
		}
		keys = append(keys, key)
	}
	slices.Sort(keys)
	values := make([]any, 0, len(keys))
	for _, key := range keys {
		values = append(values, []any{key, handoff.Values[key]})
	}
	return records.EncodeGrid(records.GridEnvelope{ProtocolPCID: workflowHandoffProtocolCID, Slots: []any{handoff.PCID, values}})
}

func DecodeWorkflowHandoff(raw []byte) (WorkflowHandoff, error) {
	envelope, err := records.DecodeGrid(raw)
	if err != nil {
		return WorkflowHandoff{}, err
	}
	if envelope.ProtocolPCID != workflowHandoffProtocolCID || len(envelope.Slots) != 2 {
		return WorkflowHandoff{}, errors.New("invalid workflow handoff envelope")
	}
	pcid, ok := envelope.Slots[0].(string)
	if !ok || strings.TrimSpace(pcid) == "" {
		return WorkflowHandoff{}, errors.New("workflow handoff pCID must be text")
	}
	if err := validateWorkflowPCID(pcid); err != nil {
		return WorkflowHandoff{}, err
	}
	entries, ok := envelope.Slots[1].([]any)
	if !ok || len(entries) == 0 {
		return WorkflowHandoff{}, errors.New("workflow handoff values must be a non-empty array")
	}
	values := make(map[string]string, len(entries))
	previous := ""
	for _, entry := range entries {
		pair, ok := entry.([]any)
		if !ok || len(pair) != 2 {
			return WorkflowHandoff{}, errors.New("workflow handoff value must be a key/value pair")
		}
		key, keyOK := pair[0].(string)
		value, valueOK := pair[1].(string)
		if !keyOK || !valueOK || key == "" || value == "" || key <= previous {
			return WorkflowHandoff{}, errors.New("workflow handoff values must be sorted unique non-empty text")
		}
		values[key] = value
		previous = key
	}
	return WorkflowHandoff{PCID: pcid, Values: values}, nil
}

type WorkflowRunState string

const (
	WorkflowRunRunning   WorkflowRunState = "running"
	WorkflowRunWaiting   WorkflowRunState = "waiting-for-input"
	WorkflowRunCompleted WorkflowRunState = "completed"
	WorkflowRunFailed    WorkflowRunState = "failed"
)

type WorkflowRun struct {
	ID        string           `json:"id"`
	Workflow  string           `json:"workflow"`
	State     WorkflowRunState `json:"state"`
	InputCID  string           `json:"input_cid"`
	OutputCID string           `json:"output_cid,omitempty"`
	Reason    string           `json:"reason,omitempty"`
	EventCID  string           `json:"event_cid"`
}
type workflowRunEvent struct {
	RunCID   cid.Cid
	Workflow string
	State    WorkflowRunState
	Input    cid.Cid
	Output   cid.Cid
	Reason   string
	Parents  []cid.Cid
}

func encodeWorkflowRunEvent(event workflowRunEvent) ([]byte, error) {
	if event.Workflow == "" || event.RunCID.Version() != 1 || event.Input.Version() != 1 || len(event.Parents) > 1 || !validWorkflowRunState(event.State) || (event.State != WorkflowRunRunning && len(event.Parents) != 1) {
		return nil, errors.New("invalid workflow run event")
	}
	if event.State == WorkflowRunCompleted && !event.Output.Defined() {
		return nil, errors.New("completed workflow run event requires output")
	}
	parents := make([]any, 0, len(event.Parents))
	for _, parent := range event.Parents {
		parents = append(parents, parent.Bytes())
	}
	output := []byte{}
	if event.Output.Defined() {
		output = event.Output.Bytes()
	}
	return records.EncodeGrid(records.GridEnvelope{ProtocolPCID: workflowRunProtocolCID, Slots: []any{string(event.State), event.RunCID.Bytes(), event.Workflow, event.Input.Bytes(), output, event.Reason, parents}})
}

func validWorkflowRunState(state WorkflowRunState) bool {
	return state == WorkflowRunRunning || state == WorkflowRunWaiting || state == WorkflowRunCompleted || state == WorkflowRunFailed
}
func validateWorkflowPCID(value string) error {
	protocol, err := cid.Decode(value)
	if err != nil || protocol.Version() != 1 {
		return errors.New("workflow pCID must be a CIDv1")
	}
	return nil
}
func newWorkflowRunCID(cas *store.CAS) (cid.Cid, error) {
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return cid.Undef, err
	}
	body, err := records.EncodeGrid(records.GridEnvelope{ProtocolPCID: workflowRunProtocolCID, Slots: []any{"workflow-run-nonce", nonce}})
	if err != nil {
		return cid.Undef, err
	}
	return cas.PutCID(body)
}
func isRetainedWorkflowRunNonce(cas *store.CAS, runCID cid.Cid) bool {
	body, err := cas.GetCID(runCID)
	if err != nil {
		return false
	}
	envelope, err := records.DecodeGrid(body)
	if err != nil || envelope.ProtocolPCID != workflowRunProtocolCID || len(envelope.Slots) != 2 {
		return false
	}
	label, labelOK := envelope.Slots[0].(string)
	nonce, nonceOK := envelope.Slots[1].([]byte)
	return labelOK && label == "workflow-run-nonce" && nonceOK && len(nonce) == 32
}
func isRetainedWorkflowHandoff(cas *store.CAS, objectCID cid.Cid) bool {
	body, err := cas.GetCID(objectCID)
	if err != nil {
		return false
	}
	_, err = DecodeWorkflowHandoff(body)
	return err == nil
}
func decodeWorkflowRunEvent(raw []byte) (workflowRunEvent, error) {
	envelope, err := records.DecodeGrid(raw)
	if err != nil {
		return workflowRunEvent{}, err
	}
	if envelope.ProtocolPCID != workflowRunProtocolCID || len(envelope.Slots) != 7 {
		return workflowRunEvent{}, errors.New("invalid workflow run envelope")
	}
	state, ok := envelope.Slots[0].(string)
	if !ok {
		return workflowRunEvent{}, errors.New("workflow run state must be text")
	}
	runBytes, ok := envelope.Slots[1].([]byte)
	if !ok {
		return workflowRunEvent{}, errors.New("workflow run CID must be bytes")
	}
	runCID, err := cid.Cast(runBytes)
	if err != nil {
		return workflowRunEvent{}, err
	}
	workflow, ok := envelope.Slots[2].(string)
	if !ok {
		return workflowRunEvent{}, errors.New("workflow run workflow must be text")
	}
	inputBytes, ok := envelope.Slots[3].([]byte)
	if !ok {
		return workflowRunEvent{}, errors.New("workflow run input must be bytes")
	}
	input, err := cid.Cast(inputBytes)
	if err != nil {
		return workflowRunEvent{}, err
	}
	outputBytes, ok := envelope.Slots[4].([]byte)
	if !ok {
		return workflowRunEvent{}, errors.New("workflow run output must be bytes")
	}
	output := cid.Undef
	if len(outputBytes) > 0 {
		output, err = cid.Cast(outputBytes)
		if err != nil {
			return workflowRunEvent{}, err
		}
	}
	reason, ok := envelope.Slots[5].(string)
	if !ok {
		return workflowRunEvent{}, errors.New("workflow run reason must be text")
	}
	parentValues, ok := envelope.Slots[6].([]any)
	if !ok {
		return workflowRunEvent{}, errors.New("workflow run parents must be array")
	}
	parents := make([]cid.Cid, 0, len(parentValues))
	for _, value := range parentValues {
		rawParent, ok := value.([]byte)
		if !ok {
			return workflowRunEvent{}, errors.New("workflow run parent must be bytes")
		}
		parent, err := cid.Cast(rawParent)
		if err != nil {
			return workflowRunEvent{}, err
		}
		parents = append(parents, parent)
	}
	event := workflowRunEvent{RunCID: runCID, Workflow: workflow, State: WorkflowRunState(state), Input: input, Output: output, Reason: reason, Parents: parents}
	if _, err = encodeWorkflowRunEvent(event); err != nil {
		return workflowRunEvent{}, err
	}
	return event, nil
}

type WorkflowOperation func(context.Context, *Runtime, WorkflowHandoff) (WorkflowHandoff, error)
type WorkflowRunRegistry struct {
	cas       *store.CAS
	cachePath string
	mu        sync.RWMutex
	runs      map[string]WorkflowRun
	heads     map[string]cid.Cid
}

func OpenWorkflowRunRegistry(root string, cas *store.CAS) (*WorkflowRunRegistry, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	registry := &WorkflowRunRegistry{cas: cas, cachePath: filepath.Join(root, "workflow-run-cache.json"), runs: map[string]WorkflowRun{}, heads: map[string]cid.Cid{}}
	return registry, registry.rebuild()
}
func (registry *WorkflowRunRegistry) rebuild() error {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	ids, err := registry.cas.ListCIDs()
	if err != nil {
		return err
	}
	events := map[string]workflowRunEvent{}
	for _, id := range ids {
		raw, readErr := registry.cas.GetCID(id)
		if readErr != nil {
			continue
		}
		event, decodeErr := decodeWorkflowRunEvent(raw)
		if decodeErr == nil {
			events[id.String()] = event
		}
	}
	accepted := map[string]workflowRunEvent{}
	for progress := true; progress; {
		progress = false
		for _, id := range ids {
			event, ok := events[id.String()]
			if !ok {
				continue
			}
			if event.State == WorkflowRunRunning && len(event.Parents) == 0 && isRetainedWorkflowRunNonce(registry.cas, event.RunCID) && isRetainedWorkflowHandoff(registry.cas, event.Input) {
				accepted[id.String()] = event
				delete(events, id.String())
				progress = true
				continue
			}
			if len(event.Parents) == 0 {
				continue
			}
			if parent, ok := accepted[event.Parents[0].String()]; ok && parent.RunCID == event.RunCID && parent.Workflow == event.Workflow && validWorkflowRunTransition(parent.State, event.State, event.Reason) && isRetainedWorkflowHandoff(registry.cas, event.Input) && (!event.Output.Defined() || isRetainedWorkflowHandoff(registry.cas, event.Output)) {
				accepted[id.String()] = event
				delete(events, id.String())
				progress = true
			}
		}
	}
	child := map[string]bool{}
	for _, event := range accepted {
		for _, parent := range event.Parents {
			child[parent.String()] = true
		}
	}
	runs := map[string]WorkflowRun{}
	heads := map[string]cid.Cid{}
	for _, id := range ids {
		event, ok := accepted[id.String()]
		if !ok || child[id.String()] {
			continue
		}
		runID := event.RunCID.String()
		if _, exists := runs[runID]; exists {
			return errors.New("workflow run has competing heads")
		}
		run := WorkflowRun{ID: runID, Workflow: event.Workflow, State: event.State, InputCID: event.Input.String(), Reason: event.Reason, EventCID: id.String()}
		if event.Output.Defined() {
			run.OutputCID = event.Output.String()
		}
		runs[runID] = run
		heads[runID] = id
	}
	registry.runs, registry.heads = runs, heads
	if err := registry.cacheLocked(); err != nil {
		// Intent: A CAS-derived run projection remains usable when its disposable
		// cache cannot be refreshed during open. Source: DI-lumek
		return nil
	}
	return nil
}
func (registry *WorkflowRunRegistry) cache() error {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	return registry.cacheLocked()
}
func (registry *WorkflowRunRegistry) cacheLocked() error {
	rows := make([]WorkflowRun, 0, len(registry.runs))
	for _, run := range registry.runs {
		rows = append(rows, run)
	}
	slices.SortFunc(rows, func(a, b WorkflowRun) int { return strings.Compare(a.ID, b.ID) })
	raw, err := json.Marshal(rows)
	if err != nil {
		return err
	}
	return os.WriteFile(registry.cachePath, append(raw, '\n'), 0644)
}
func (registry *WorkflowRunRegistry) append(event workflowRunEvent) (WorkflowRun, error) {
	raw, err := encodeWorkflowRunEvent(event)
	if err != nil {
		return WorkflowRun{}, err
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if len(event.Parents) == 0 {
		if _, exists := registry.heads[event.RunCID.String()]; exists {
			return WorkflowRun{}, errors.New("workflow run already has a root")
		}
	} else {
		current, exists := registry.runs[event.RunCID.String()]
		if !exists || current.EventCID != event.Parents[0].String() || !validWorkflowRunTransition(current.State, event.State, event.Reason) {
			return WorkflowRun{}, errors.New("workflow run transition is not from its current head")
		}
	}
	eventCID, err := registry.cas.PutCID(raw)
	if err != nil {
		return WorkflowRun{}, err
	}
	run := WorkflowRun{ID: event.RunCID.String(), Workflow: event.Workflow, State: event.State, InputCID: event.Input.String(), Reason: event.Reason, EventCID: eventCID.String()}
	if event.Output.Defined() {
		run.OutputCID = event.Output.String()
	}
	registry.runs[run.ID], registry.heads[run.ID] = run, eventCID
	if err := registry.cacheLocked(); err != nil {
		return run, nil
	}
	return run, nil
}

func validWorkflowRunTransition(from, to WorkflowRunState, reason string) bool {
	switch from {
	case WorkflowRunRunning:
		return to == WorkflowRunWaiting || to == WorkflowRunCompleted || to == WorkflowRunFailed || (to == WorkflowRunRunning && reason == "manual recovery")
	case WorkflowRunWaiting, WorkflowRunFailed:
		return to == WorkflowRunRunning
	case WorkflowRunCompleted:
		return to == WorkflowRunFailed && strings.HasPrefix(reason, "policy handoff:")
	default:
		return false
	}
}
func (registry *WorkflowRunRegistry) get(id string) (WorkflowRun, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	run, ok := registry.runs[id]
	return run, ok
}
