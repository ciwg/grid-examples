package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
	"github.com/ipfs/go-cid"
)

type BuiltinCommand func(context.Context, *Runtime, []string) (string, error)
type BuiltinValidator func(records.Envelope) error
type WorkflowAdapterWorker func(context.Context, packages.WorkflowAdapter, []byte) ([]byte, error)

type workflowChainContextKey struct{}

func withWorkflowChain(ctx context.Context, workflow string) context.Context {
	chain, _ := ctx.Value(workflowChainContextKey{}).([]string)
	return context.WithValue(ctx, workflowChainContextKey{}, append(append([]string{}, chain...), workflow))
}
func workflowChainContains(ctx context.Context, workflow string) bool {
	chain, _ := ctx.Value(workflowChainContextKey{}).([]string)
	for _, entry := range chain {
		if entry == workflow {
			return true
		}
	}
	return false
}

type BuiltinPackage struct {
	Manifest   packages.Manifest
	Commands   map[string]BuiltinCommand
	Validators map[string]BuiltinValidator
}

type activePackage struct {
	manifest    packages.Manifest
	builtin     bool
	commands    map[string]BuiltinCommand
	validators  map[string]BuiltinValidator
	external    packages.Runner
	packageRoot string
}

type registeredFamily struct {
	owner        *activePackage
	protocolPCID string
}

type Runtime struct {
	root                    string
	packagesRoot            string
	history                 *store.History
	cas                     *store.CAS
	workflowEvidence        *store.CAS
	workflowReceipts        *workflowReceiptStore
	peers                   *grid.PeerStore
	policies                *grid.PolicyStore
	packages                map[string]*activePackage
	commands                map[string]*activePackage
	families                map[string]registeredFamily
	routes                  []registeredRoute
	workflows               *WorkflowRegistry
	workflowRuns            *WorkflowRunRegistry
	workflowOps             map[string]WorkflowOperation
	workflowAdapters        map[string]packages.WorkflowAdapter
	workflowAdapterPackages map[string]string
	workflowWorker          WorkflowAdapterWorker
	handoffPolicies         *WorkflowHandoffPolicies
	routePromises           *RoutePromiseRegistry
}

func Open(root string) (*Runtime, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	history, err := store.OpenHistory(filepath.Join(root, "state"))
	if err != nil {
		return nil, err
	}
	casStore, err := store.OpenCAS(filepath.Join(root, "cas"))
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	// Intent: Keep received lifecycle evidence out of the CAS replayed as local
	// workflow authority. Source: DI-novuk
	evidenceStore, err := store.OpenCAS(filepath.Join(root, "workflow-evidence"))
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	// Intent: Retain authenticated sender provenance separately from raw remote
	// evidence so receipt inspection cannot mutate local workflow lifecycle.
	// Source: DI-rufir
	receiptStore, err := openWorkflowReceiptStore(workflowReceiptPath(root))
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	peerStore, err := grid.OpenPeerStore(filepath.Join(root, "state"), root)
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	policyStore, err := grid.OpenPolicyStore(filepath.Join(root, "state"))
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	// Intent: Rebuild local workflow availability from durable lifecycle events
	// before any installed package can participate in the runtime. Source: DI-lovek
	workflowRegistry, err := OpenWorkflowRegistry(filepath.Join(root, "state"), casStore)
	if err != nil {
		if closeErr := history.Close(); closeErr != nil {
			return nil, errors.Join(err, closeErr)
		}
		return nil, err
	}
	workflowRuns, err := OpenWorkflowRunRegistry(filepath.Join(root, "state"), casStore)
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	handoffPolicies, err := OpenWorkflowHandoffPolicies(filepath.Join(root, "state"))
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	routePromises, err := OpenRoutePromiseRegistry(casStore)
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	runtime := &Runtime{
		root:                    root,
		packagesRoot:            filepath.Join(root, "packages"),
		history:                 history,
		cas:                     casStore,
		workflowEvidence:        evidenceStore,
		workflowReceipts:        receiptStore,
		peers:                   peerStore,
		policies:                policyStore,
		packages:                map[string]*activePackage{},
		commands:                map[string]*activePackage{},
		families:                map[string]registeredFamily{},
		routes:                  []registeredRoute{},
		workflows:               workflowRegistry,
		workflowRuns:            workflowRuns,
		workflowOps:             map[string]WorkflowOperation{},
		workflowAdapters:        map[string]packages.WorkflowAdapter{},
		workflowAdapterPackages: map[string]string{},
		workflowWorker:          runDockerWorkflowAdapter,
		handoffPolicies:         handoffPolicies,
		routePromises:           routePromises,
	}
	if err := os.MkdirAll(runtime.packagesRoot, 0o755); err != nil {
		_ = history.Close()
		return nil, err
	}
	if err := runtime.activateInstalledFromRoot(context.Background()); err != nil {
		_ = history.Close()
		return nil, err
	}
	return runtime, nil
}

// runDockerWorkflowAdapter is the only production executable-adapter path.
// Intent: Prevent an installed workflow adapter from receiving ambient host
// authority or a direct-process fallback. Source: DI-fofuh
func runDockerWorkflowAdapter(ctx context.Context, adapter packages.WorkflowAdapter, input []byte) ([]byte, error) {
	timeout, err := time.ParseDuration(adapter.Timeout)
	if err != nil {
		return nil, fmt.Errorf("workflow adapter timeout: %w", err)
	}
	return (packages.DockerWorker{
		Image:     adapter.Image,
		Args:      adapter.Command,
		CPUs:      adapter.CPUs,
		Memory:    adapter.Memory,
		PIDsLimit: adapter.PIDsLimit,
		Timeout:   timeout,
	}).Run(ctx, input)
}

func (runtime *Runtime) RegistryAllowList() []string     { return runtime.policies.RegistryAllowList() }
func (runtime *Runtime) AllowRegistry(host string) error { return runtime.policies.AllowRegistry(host) }
func (runtime *Runtime) RemoveRegistry(host string) error {
	return runtime.policies.RemoveRegistry(host)
}

// BindAgent records the local implementation adapter selected for an explicitly
// named app agent. Intent: Keep package installation from silently becoming an
// app identity or a live receive promise. Source: DI-komaz; DI-butam
func (runtime *Runtime) BindAgent(binding AgentBinding) error {
	return runtime.routePromises.BindAgent(binding)
}

// PublishReceivePromise retains an app agent's explicit local promise to
// receive one pCID. Source: DI-kojab
func (runtime *Runtime) PublishReceivePromise(promise ReceivePromise) error {
	return runtime.routePromises.PublishReceivePromise(promise)
}

// PublishDeliveryPromise retains a routing-role agent's explicit local promise
// to deliver one pCID to one named app agent. Source: DI-kojab
func (runtime *Runtime) PublishDeliveryPromise(promise DeliveryPromise) error {
	return runtime.routePromises.PublishDeliveryPromise(promise)
}

// PullWorkflowImage explicitly acquires only the active package adapter image
// matching a verified workflow. It does not change workflow lifecycle state.
// Source: DI-zivut
func (runtime *Runtime) PullWorkflowImage(ctx context.Context, alias string) error {
	manifest, err := runtime.VerifyWorkflow(alias)
	if err != nil {
		return err
	}
	adapter, ok := runtime.workflowAdapters[manifest.Adapter]
	if !ok {
		return errors.New("workflow adapter is unavailable")
	}
	host, err := packages.RegistryHostFromImage(adapter.Image)
	if err != nil {
		return err
	}
	allowed := false
	for _, entry := range runtime.RegistryAllowList() {
		if entry == host {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("adapter registry is not allowed: %s", host)
	}
	return packages.PullImage(ctx, adapter.Image)
}

// RegisterWorkflowOperation binds an active artifact's declared adapter to
// trusted built-in behavior. Intent: Keep execution behind explicit artifact
// activation rather than making installed capabilities implicitly runnable.
// Source: DI-lumek
func (runtime *Runtime) RegisterWorkflowOperation(name string, operation WorkflowOperation) error {
	if strings.TrimSpace(name) == "" || operation == nil {
		return errors.New("workflow operation name and handler are required")
	}
	if _, exists := runtime.workflowOps[name]; exists {
		return errors.New("workflow operation is already registered")
	}
	runtime.workflowOps[name] = operation
	return nil
}

func (runtime *Runtime) StartWorkflowRun(ctx context.Context, workflowID string, input WorkflowHandoff) (WorkflowRun, error) {
	workflow, err := runtime.workflow(workflowID)
	if err != nil {
		return WorkflowRun{}, err
	}
	if workflow.State != WorkflowActive {
		return WorkflowRun{}, errors.New("workflow is not active")
	}
	manifest, err := runtime.VerifyWorkflow(workflowID)
	if err != nil {
		return WorkflowRun{}, err
	}
	if manifest.Adapter == "" || manifest.InputPCID == "" || manifest.OutputPCID == "" {
		return WorkflowRun{}, errors.New("workflow does not declare an executable adapter")
	}
	if input.PCID != manifest.InputPCID {
		return WorkflowRun{}, errors.New("workflow input pCID is not accepted")
	}
	if _, installed := runtime.workflowAdapters[manifest.Adapter]; installed {
		packageID, known := runtime.workflowAdapterPackages[manifest.Adapter]
		// Intent: An active package and Docker confinement are not an app's
		// current promise to receive work. Refuse before creating a run unless
		// the complete local binding, receive, and delivery evidence is enabled.
		// Source: DI-bidam; DI-guraj
		if !known || !runtime.routePromises.routeExecutable(packageID, manifest.InputPCID) {
			return WorkflowRun{}, errors.New("workflow adapter is not currently promised for its input pCID")
		}
	}
	raw, err := EncodeWorkflowHandoff(input)
	if err != nil {
		return WorkflowRun{}, err
	}
	inputCID, err := runtime.cas.PutCID(raw)
	if err != nil {
		return WorkflowRun{}, err
	}
	// Intent: A fresh retained nonce distinguishes repeated identical inputs,
	// without making mutable local state part of the durable run identity.
	// Source: DI-lumek
	runCID, err := newWorkflowRunCID(runtime.cas)
	if err != nil {
		return WorkflowRun{}, err
	}
	seed := workflowRunEvent{RunCID: runCID, Workflow: workflowID, State: WorkflowRunRunning, Input: inputCID}
	run, err := runtime.workflowRuns.append(seed)
	if err != nil {
		return WorkflowRun{}, err
	}
	return runtime.executeWorkflowRun(ctx, run, manifest, input)
}

func (runtime *Runtime) executeWorkflowRun(ctx context.Context, run WorkflowRun, manifest WorkflowManifest, input WorkflowHandoff) (WorkflowRun, error) {
	if adapter, installed := runtime.workflowAdapters[manifest.Adapter]; installed {
		if adapter.InputPCID != manifest.InputPCID || adapter.OutputPCID != manifest.OutputPCID {
			return runtime.transitionWorkflowRun(run, WorkflowRunFailed, cid.Undef, "active package adapter contract does not match workflow")
		}
		inputBytes, err := EncodeWorkflowHandoff(input)
		if err != nil {
			return runtime.failWorkflowRun(run, "workflow input encoding: "+err.Error(), err)
		}
		workerOutput, err := runtime.workflowWorker(ctx, adapter, inputBytes)
		if err != nil {
			return runtime.transitionWorkflowRun(run, WorkflowRunFailed, cid.Undef, err.Error())
		}
		proposal, err := DecodeWorkflowAdapterResult(workerOutput)
		if err != nil {
			return runtime.transitionWorkflowRun(run, WorkflowRunFailed, cid.Undef, "workflow worker result: "+err.Error())
		}
		if proposal.Output.PCID != manifest.OutputPCID {
			return runtime.transitionWorkflowRun(run, WorkflowRunFailed, cid.Undef, "workflow adapter emitted an undeclared output pCID")
		}
		if _, err := runtime.applyWorkflowAdapterResult(ctx, packages.CommandResult{CAS: proposal.CAS, Records: proposal.Records}); err != nil {
			return runtime.failWorkflowRun(run, "workflow worker proposal: "+err.Error(), err)
		}
		return runtime.completeWorkflowRun(ctx, run, manifest, proposal.Output)
	}
	operation := runtime.workflowOps[manifest.Adapter]
	if operation == nil {
		return runtime.transitionWorkflowRun(run, WorkflowRunFailed, cid.Undef, "workflow adapter is unavailable")
	}
	output, err := operation(ctx, runtime, input)
	if err != nil {
		var waiting *WorkflowWaitingError
		if errors.As(err, &waiting) {
			return runtime.transitionWorkflowRun(run, WorkflowRunWaiting, cid.Undef, waiting.Error())
		}
		return runtime.transitionWorkflowRun(run, WorkflowRunFailed, cid.Undef, err.Error())
	}
	return runtime.completeWorkflowRun(ctx, run, manifest, output)
}

func (runtime *Runtime) completeWorkflowRun(ctx context.Context, run WorkflowRun, manifest WorkflowManifest, output WorkflowHandoff) (WorkflowRun, error) {
	if output.PCID != manifest.OutputPCID {
		return runtime.transitionWorkflowRun(run, WorkflowRunFailed, cid.Undef, "workflow adapter emitted an undeclared output pCID")
	}
	raw, err := EncodeWorkflowHandoff(output)
	if err != nil {
		return runtime.failWorkflowRun(run, "workflow output encoding: "+err.Error(), err)
	}
	outputCID, err := runtime.cas.PutCID(raw)
	if err != nil {
		return runtime.failWorkflowRun(run, "workflow output persistence: "+err.Error(), err)
	}
	completed, err := runtime.transitionWorkflowRun(run, WorkflowRunCompleted, outputCID, "")
	if err != nil {
		return WorkflowRun{}, err
	}
	if policy, ok := runtime.handoffPolicies.Find(completed.Workflow, output.PCID); ok {
		// Intent: Policy routing is convenience, never a bypass of the target's
		// active-artifact and input-pCID checks. Source: DI-lumek
		if workflowChainContains(ctx, policy.TargetWorkflow) {
			return runtime.transitionWorkflowRun(completed, WorkflowRunFailed, outputCID, "policy handoff: cycle detected")
		}
		if _, err := runtime.HandoffWorkflowRun(withWorkflowChain(ctx, completed.Workflow), completed.ID, policy.TargetWorkflow); err != nil {
			return runtime.transitionWorkflowRun(completed, WorkflowRunFailed, outputCID, "policy handoff: "+err.Error())
		}
	}
	return completed, nil
}

// failWorkflowRun records a terminal failure after adapter-side work reaches a
// persistence boundary. Intent: Prefer a terminal failure, while an explicit
// retry remains available if physical CAS failure prevents that transition.
// Source: DI-lumek
func (runtime *Runtime) failWorkflowRun(run WorkflowRun, reason string, cause error) (WorkflowRun, error) {
	failed, err := runtime.transitionWorkflowRun(run, WorkflowRunFailed, cid.Undef, reason)
	if err != nil {
		return WorkflowRun{}, errors.Join(cause, err)
	}
	return failed, nil
}

func (runtime *Runtime) transitionWorkflowRun(run WorkflowRun, state WorkflowRunState, output cid.Cid, reason string) (WorkflowRun, error) {
	parent, err := cid.Decode(run.EventCID)
	if err != nil {
		return WorkflowRun{}, err
	}
	input, err := cid.Decode(run.InputCID)
	if err != nil {
		return WorkflowRun{}, err
	}
	runCID, err := cid.Decode(run.ID)
	if err != nil {
		return WorkflowRun{}, err
	}
	next, err := runtime.workflowRuns.append(workflowRunEvent{RunCID: runCID, Workflow: run.Workflow, State: state, Input: input, Output: output, Reason: reason, Parents: []cid.Cid{parent}})
	if err != nil {
		return WorkflowRun{}, err
	}
	return next, nil
}

func (runtime *Runtime) WorkflowRun(id string) (WorkflowRun, error) {
	run, ok := runtime.workflowRuns.get(id)
	if !ok {
		return WorkflowRun{}, errors.New("workflow run is not found")
	}
	return run, nil
}

// HandoffWorkflowRun transfers an exact completed output into another active
// artifact's declared input contract. Intent: Preserve the basket boundary at
// handoff time so an inactive artifact cannot receive work. Source: DI-lumek
func (runtime *Runtime) HandoffWorkflowRun(ctx context.Context, runID string, targetWorkflow string) (WorkflowRun, error) {
	source, err := runtime.WorkflowRun(runID)
	if err != nil {
		return WorkflowRun{}, err
	}
	if source.State != WorkflowRunCompleted || source.OutputCID == "" {
		return WorkflowRun{}, errors.New("workflow run has no completed output to hand off")
	}
	sourceWorkflow, err := runtime.workflow(source.Workflow)
	if err != nil {
		return WorkflowRun{}, err
	}
	if sourceWorkflow.State != WorkflowActive {
		return WorkflowRun{}, errors.New("source workflow is not active")
	}
	outputCID, err := cid.Decode(source.OutputCID)
	if err != nil {
		return WorkflowRun{}, err
	}
	raw, err := runtime.cas.GetCID(outputCID)
	if err != nil {
		return WorkflowRun{}, err
	}
	handoff, err := DecodeWorkflowHandoff(raw)
	if err != nil {
		return WorkflowRun{}, err
	}
	target, err := runtime.workflow(targetWorkflow)
	if err != nil {
		return WorkflowRun{}, err
	}
	if target.State != WorkflowActive {
		return WorkflowRun{}, errors.New("workflow is not active")
	}
	manifest, err := runtime.VerifyWorkflow(targetWorkflow)
	if err != nil {
		return WorkflowRun{}, err
	}
	if handoff.PCID != manifest.InputPCID {
		if manifest.Adapter == "" || manifest.InputPCID == "" || manifest.OutputPCID == "" || runtime.workflowOps[manifest.Adapter] == nil {
			return WorkflowRun{}, errors.New("target workflow does not declare an available executable adapter")
		}
		// Intent: Preserve the source envelope as evidence when a policy or
		// operator selects a different target schema. The target must receive
		// an explicit input under its own pCID before its adapter can execute.
		// Source: DI-lumek
		return runtime.queueWorkflowRun(targetWorkflow, handoff, "handoff output pCID does not satisfy target input pCID")
	}
	return runtime.StartWorkflowRun(ctx, targetWorkflow, handoff)
}

// queueWorkflowRun records an incompatible handoff without treating source
// data as if it already satisfied the target artifact's declared schema.
func (runtime *Runtime) queueWorkflowRun(workflowID string, input WorkflowHandoff, reason string) (WorkflowRun, error) {
	raw, err := EncodeWorkflowHandoff(input)
	if err != nil {
		return WorkflowRun{}, err
	}
	inputCID, err := runtime.cas.PutCID(raw)
	if err != nil {
		return WorkflowRun{}, err
	}
	runCID, err := newWorkflowRunCID(runtime.cas)
	if err != nil {
		return WorkflowRun{}, err
	}
	run, err := runtime.workflowRuns.append(workflowRunEvent{RunCID: runCID, Workflow: workflowID, State: WorkflowRunRunning, Input: inputCID})
	if err != nil {
		return WorkflowRun{}, err
	}
	return runtime.transitionWorkflowRun(run, WorkflowRunWaiting, cid.Undef, reason)
}

// WorkflowWaitingError tells the runtime that an adapter needs more input, not
// that its execution failed. Source: DI-lumek
type WorkflowWaitingError struct{ Reason string }

func (err *WorkflowWaitingError) Error() string   { return err.Reason }
func WaitingForWorkflowInput(reason string) error { return &WorkflowWaitingError{Reason: reason} }

func (runtime *Runtime) SupplyWorkflowRun(ctx context.Context, runID string, input WorkflowHandoff) (WorkflowRun, error) {
	run, err := runtime.WorkflowRun(runID)
	if err != nil {
		return WorkflowRun{}, err
	}
	if run.State != WorkflowRunWaiting {
		return WorkflowRun{}, errors.New("workflow run is not waiting for input")
	}
	workflow, err := runtime.workflow(run.Workflow)
	if err != nil {
		return WorkflowRun{}, err
	}
	if workflow.State != WorkflowActive {
		return WorkflowRun{}, errors.New("workflow is not active")
	}
	manifest, err := runtime.VerifyWorkflow(run.Workflow)
	if err != nil {
		return WorkflowRun{}, err
	}
	if input.PCID != manifest.InputPCID {
		return WorkflowRun{}, errors.New("workflow input pCID is not accepted")
	}
	raw, err := EncodeWorkflowHandoff(input)
	if err != nil {
		return WorkflowRun{}, err
	}
	inputCID, err := runtime.cas.PutCID(raw)
	if err != nil {
		return WorkflowRun{}, err
	}
	resumed, err := runtime.resumeWorkflowRun(run, inputCID)
	if err != nil {
		return WorkflowRun{}, err
	}
	return runtime.executeWorkflowRun(ctx, resumed, manifest, input)
}
func (runtime *Runtime) RetryWorkflowRun(ctx context.Context, runID string) (WorkflowRun, error) {
	run, err := runtime.WorkflowRun(runID)
	if err != nil {
		return WorkflowRun{}, err
	}
	if run.State != WorkflowRunFailed && run.State != WorkflowRunRunning {
		return WorkflowRun{}, errors.New("workflow run is not failed or recovery-running")
	}
	workflow, err := runtime.workflow(run.Workflow)
	if err != nil {
		return WorkflowRun{}, err
	}
	if workflow.State != WorkflowActive {
		return WorkflowRun{}, errors.New("workflow is not active")
	}
	manifest, err := runtime.VerifyWorkflow(run.Workflow)
	if err != nil {
		return WorkflowRun{}, err
	}
	inputCID, err := cid.Decode(run.InputCID)
	if err != nil {
		return WorkflowRun{}, err
	}
	raw, err := runtime.cas.GetCID(inputCID)
	if err != nil {
		return WorkflowRun{}, err
	}
	input, err := DecodeWorkflowHandoff(raw)
	if err != nil {
		return WorkflowRun{}, err
	}
	resumed, err := runtime.resumeWorkflowRun(run, inputCID)
	if err != nil {
		return WorkflowRun{}, err
	}
	return runtime.executeWorkflowRun(ctx, resumed, manifest, input)
}
func (runtime *Runtime) resumeWorkflowRun(run WorkflowRun, input cid.Cid) (WorkflowRun, error) {
	head, err := cid.Decode(run.EventCID)
	if err != nil {
		return WorkflowRun{}, err
	}
	root, err := cid.Decode(run.ID)
	if err != nil {
		return WorkflowRun{}, err
	}
	reason := ""
	if run.State == WorkflowRunRunning {
		// Intent: Preserve an auditable, operator-only recovery edge after a
		// physical persistence interruption; ordinary running transitions remain
		// invalid. Source: DI-lumek
		reason = "manual recovery"
	}
	return runtime.workflowRuns.append(workflowRunEvent{RunCID: root, Workflow: run.Workflow, State: WorkflowRunRunning, Input: input, Reason: reason, Parents: []cid.Cid{head}})
}
func (runtime *Runtime) SetWorkflowHandoffPolicy(policy WorkflowHandoffPolicy) error {
	sourceWorkflow, err := runtime.workflow(policy.SourceWorkflow)
	if err != nil {
		return err
	}
	source, err := runtime.WorkflowManifest(sourceWorkflow.ArtifactCID)
	if err != nil {
		return err
	}
	targetWorkflow, err := runtime.workflow(policy.TargetWorkflow)
	if err != nil {
		return err
	}
	target, err := runtime.WorkflowManifest(targetWorkflow.ArtifactCID)
	if err != nil {
		return err
	}
	if source.OutputPCID != policy.OutputPCID || target.InputPCID != policy.InputPCID {
		return errors.New("workflow handoff policy does not match declared pCIDs")
	}
	policies := []WorkflowHandoffPolicy{}
	for _, existing := range runtime.handoffPolicies.List() {
		if existing.SourceWorkflow != policy.SourceWorkflow || existing.OutputPCID != policy.OutputPCID {
			policies = append(policies, existing)
		}
	}
	policies = append(policies, policy)
	if workflowPolicyReaches(policies, policy.TargetWorkflow, policy.SourceWorkflow) {
		return errors.New("workflow handoff policy creates a cycle")
	}
	return runtime.handoffPolicies.Set(policy)
}

func (runtime *Runtime) RemoveWorkflowHandoffPolicy(source, output string) error {
	return runtime.handoffPolicies.Remove(source, output)
}
func (runtime *Runtime) WorkflowHandoffPolicies() []WorkflowHandoffPolicy {
	return runtime.handoffPolicies.List()
}
func (runtime *Runtime) WorkflowRuns() []WorkflowRun {
	runtime.workflowRuns.mu.RLock()
	defer runtime.workflowRuns.mu.RUnlock()
	runs := make([]WorkflowRun, 0, len(runtime.workflowRuns.runs))
	for _, run := range runtime.workflowRuns.runs {
		runs = append(runs, run)
	}
	// Intent: Present current run heads in operator-relevant event order while
	// retaining a stable identifier tie break and readable v1 history. Source: DI-gihor
	slices.SortFunc(runs, func(a, b WorkflowRun) int {
		switch {
		case a.UpdatedAt == nil && b.UpdatedAt != nil:
			return 1
		case a.UpdatedAt != nil && b.UpdatedAt == nil:
			return -1
		case a.UpdatedAt != nil && b.UpdatedAt != nil:
			if comparison := b.UpdatedAt.Compare(*a.UpdatedAt); comparison != 0 {
				return comparison
			}
		}
		return strings.Compare(a.ID, b.ID)
	})
	return runs
}

func (runtime *Runtime) Close() error {
	return runtime.history.Close()
}

func (runtime *Runtime) RegisterBuiltin(pkg BuiltinPackage) error {
	if err := pkg.Manifest.Validate(); err != nil {
		return err
	}
	registered := &activePackage{
		manifest:   pkg.Manifest,
		builtin:    true,
		commands:   pkg.Commands,
		validators: pkg.Validators,
	}
	return runtime.activatePackage(registered)
}

func (runtime *Runtime) ActivateInstalled(ctx context.Context, manifestPath string) error {
	manifest, err := packages.LoadManifest(manifestPath)
	if err != nil {
		return err
	}
	if strings.TrimSpace(manifest.Executable) == "" {
		return fmt.Errorf("package %s executable is required for installed packages", manifest.ID)
	}
	described, err := (packages.Runner{Executable: manifest.Executable}).Describe(ctx)
	if err != nil {
		return err
	}
	described.Executable = manifest.Executable
	if !manifest.Equal(described) {
		return fmt.Errorf("manifest self-check mismatch for package %s", manifest.ID)
	}
	return runtime.activatePackage(&activePackage{
		manifest:    manifest,
		external:    packages.Runner{Executable: manifest.Executable},
		packageRoot: filepath.Dir(manifestPath),
	})
}

func (runtime *Runtime) activatePackage(pkg *activePackage) error {
	if _, exists := runtime.packages[pkg.manifest.ID]; exists {
		return fmt.Errorf("package already active: %s", pkg.manifest.ID)
	}
	for _, command := range pkg.manifest.Commands {
		key := command.Key()
		if _, exists := runtime.commands[key]; exists {
			return fmt.Errorf("command already registered: %s", key)
		}
	}
	for _, family := range pkg.manifest.Families {
		if _, exists := runtime.families[family.Name]; exists {
			return fmt.Errorf("family already registered: %s", family.Name)
		}
		// Intent: Treat family validation as an explicit routing promise so the
		// current runtime moves toward pCID-declared routing roles instead of
		// leaving family handling implicit in local package ownership alone.
		// Source: DI-rutom
		if !pkg.manifest.HasClaim(family.ProtocolPCID, "family-validator") {
			return fmt.Errorf("family %s requires family-validator claim for protocol %s", family.Name, family.ProtocolPCID)
		}
	}
	for _, adapter := range pkg.manifest.WorkflowAdapters {
		if _, exists := runtime.workflowAdapters[adapter.Name]; exists {
			return fmt.Errorf("workflow adapter already registered: %s", adapter.Name)
		}
	}
	runtime.packages[pkg.manifest.ID] = pkg
	for _, command := range pkg.manifest.Commands {
		runtime.commands[command.Key()] = pkg
	}
	for _, family := range pkg.manifest.Families {
		runtime.families[family.Name] = registeredFamily{
			owner:        pkg,
			protocolPCID: family.ProtocolPCID,
		}
	}
	for _, claim := range pkg.manifest.Claims {
		// Intent: Keep the route table aligned with explicit multi-hop route
		// metadata so parser/transform roles are visible next to direct handlers.
		// Source: DI-lafek
		runtime.routes = append(runtime.routes, registeredRoute{
			owner:        pkg,
			protocolPCID: claim.ProtocolPCID,
			role:         claim.Role,
			routeType:    claim.NormalizedRouteType(),
			emits:        claim.SortedEmitsProtocols(),
			summary:      claim.Summary,
			families:     pkg.manifest.FamiliesForProtocol(claim.ProtocolPCID),
		})
	}
	for _, adapter := range pkg.manifest.WorkflowAdapters {
		runtime.workflowAdapters[adapter.Name] = adapter
		runtime.workflowAdapterPackages[adapter.Name] = pkg.manifest.ID
	}
	return nil
}

func (runtime *Runtime) PackageManifests() []packages.Manifest {
	manifests := make([]packages.Manifest, 0, len(runtime.packages))
	for _, pkg := range runtime.packages {
		manifests = append(manifests, pkg.manifest)
	}
	packages.SortManifests(manifests)
	return manifests
}

func (runtime *Runtime) PackageManifest(id string) (packages.Manifest, bool) {
	pkg, ok := runtime.packages[id]
	if !ok {
		return packages.Manifest{}, false
	}
	return pkg.manifest, true
}

func (runtime *Runtime) LocalPeerID() string {
	return runtime.peers.LocalPeerID()
}

func (runtime *Runtime) LocalPeerPublicKey() string {
	return runtime.peers.LocalPublicKey()
}

func (runtime *Runtime) AllowedPeers() []grid.AllowedPeer {
	return runtime.peers.AllowedPeers()
}

func (runtime *Runtime) AllowPeer(peer grid.AllowedPeer) error {
	return runtime.peers.SetAllowedPeer(peer)
}

func (runtime *Runtime) SetPeerTrust(peerID string, attesterClass string, weight int) error {
	return runtime.peers.SetPeerTrust(peerID, attesterClass, weight)
}

func (runtime *Runtime) SetPeerFederation(peerID string, federation string) error {
	return runtime.peers.SetPeerFederation(peerID, federation)
}

func (runtime *Runtime) RevokePeer(peerID string) error {
	return runtime.peers.RemoveAllowedPeer(peerID)
}

func (runtime *Runtime) LookupPeer(peerID string) (grid.AllowedPeer, bool) {
	return runtime.peers.Lookup(peerID)
}

func (runtime *Runtime) ClaimPolicies() []grid.ClaimTrustPolicy {
	return runtime.policies.ClaimPolicies()
}

func (runtime *Runtime) RoutePlanPolicy() grid.RoutePlanPolicy {
	return runtime.policies.RoutePlanPolicy()
}

func (runtime *Runtime) ProtocolRoutePlanPolicies() []grid.ProtocolRoutePlanPolicy {
	return runtime.policies.ProtocolRoutePlanPolicies()
}

func (runtime *Runtime) ProtocolRoutePlanPolicy(protocolPCID string) (grid.RoutePlanPolicy, bool) {
	return runtime.policies.ProtocolRoutePlanPolicy(protocolPCID)
}

func (runtime *Runtime) ProtocolRoleRoutePlanPolicies() []grid.ProtocolRoleRoutePlanPolicy {
	return runtime.policies.ProtocolRoleRoutePlanPolicies()
}

func (runtime *Runtime) ProtocolRoleRoutePlanPolicy(protocolPCID string, role string) (grid.RoutePlanPolicy, bool) {
	return runtime.policies.ProtocolRoleRoutePlanPolicy(protocolPCID, role)
}

func (runtime *Runtime) TraceScopeAliases() []grid.TraceScopeAlias {
	return runtime.policies.TraceScopeAliases()
}

func (runtime *Runtime) TraceScopeAlias(name string) (grid.TraceScopeAlias, bool) {
	return runtime.policies.TraceScopeAlias(name)
}

func (runtime *Runtime) EffectiveRoutePlanPolicy(protocolPCID string) grid.RoutePlanPolicy {
	return runtime.policies.EffectiveRoutePlanPolicy(protocolPCID)
}

func (runtime *Runtime) EffectiveRoutePlanPolicyForRole(protocolPCID string, role string) grid.RoutePlanPolicy {
	return runtime.policies.EffectiveRoutePlanPolicyForRole(protocolPCID, role)
}

func (runtime *Runtime) SetRoutePlanPolicy(policy grid.RoutePlanPolicy) error {
	return runtime.policies.SetRoutePlanPolicy(policy)
}

func (runtime *Runtime) SetProtocolRoutePlanPolicy(protocolPCID string, policy grid.RoutePlanPolicy) error {
	return runtime.policies.SetProtocolRoutePlanPolicy(protocolPCID, policy)
}

func (runtime *Runtime) SetProtocolRoleRoutePlanPolicy(protocolPCID string, role string, policy grid.RoutePlanPolicy) error {
	return runtime.policies.SetProtocolRoleRoutePlanPolicy(protocolPCID, role, policy)
}

// Intent: Keep reusable trace-view naming under runtime-owned policy state so
// operators can define local routing inspection vocabularies without editing
// code or changing the built-in scope presets. Source: DI-bemok
func (runtime *Runtime) SetTraceScopeAlias(alias grid.TraceScopeAlias) error {
	if _, builtIn := builtinTraceScopeClauses(alias.Name); builtIn {
		return fmt.Errorf("trace scope alias conflicts with built-in scope: %s", alias.Name)
	}
	return runtime.policies.SetTraceScopeAlias(alias)
}

func (runtime *Runtime) RemoveProtocolRoutePlanPolicy(protocolPCID string) error {
	return runtime.policies.RemoveProtocolRoutePlanPolicy(protocolPCID)
}

func (runtime *Runtime) RemoveProtocolRoleRoutePlanPolicy(protocolPCID string, role string) error {
	return runtime.policies.RemoveProtocolRoleRoutePlanPolicy(protocolPCID, role)
}

func (runtime *Runtime) RemoveTraceScopeAlias(name string) error {
	return runtime.policies.RemoveTraceScopeAlias(name)
}

func (runtime *Runtime) SetClaimPolicy(policy grid.ClaimTrustPolicy) error {
	for _, peerID := range policy.AllowedAttesters {
		if _, ok := runtime.LookupPeer(peerID); !ok {
			return fmt.Errorf("unknown attester peer: %s", peerID)
		}
	}
	return runtime.policies.SetClaimPolicy(policy)
}

func (runtime *Runtime) RemoveClaimPolicy(protocolPCID string, role string) error {
	return runtime.policies.RemoveClaimPolicy(protocolPCID, role)
}

func (runtime *Runtime) PutCAS(body []byte) (string, error) {
	return runtime.cas.Put(body)
}

func (runtime *Runtime) GetCAS(objectID string) ([]byte, error) {
	return runtime.cas.Get(objectID)
}

// Intent: Validate known families through their owning package while still
// preserving unknown-family carriage as durable exact bytes for the grid.
// Source: DI-lupok
func (runtime *Runtime) AppendRecord(ctx context.Context, raw []byte) (records.Envelope, error) {
	return runtime.appendRecord(ctx, raw, true)
}

func (runtime *Runtime) appendRecord(ctx context.Context, raw []byte, signLocal bool) (records.Envelope, error) {
	_, prepared, err := runtime.prepareRecord(ctx, raw, signLocal, true)
	if err != nil {
		return records.Envelope{}, err
	}
	envelope, _, err := runtime.history.AppendRaw(prepared)
	return envelope, err
}

func (runtime *Runtime) prepareRecord(ctx context.Context, raw []byte, signLocal bool, allowExternalValidation bool) (records.Envelope, []byte, error) {
	envelope, err := records.Parse(raw)
	if err != nil {
		return records.Envelope{}, nil, err
	}
	if signLocal && !envelope.HasAuthorSignature() {
		// Intent: Sign locally authored durable records once at creation time so
		// later relay trust can distinguish semantic authoring from carriage.
		// Source: DI-sovem
		envelope, err = runtime.peers.SignAuthorEnvelope(envelope)
		if err != nil {
			return records.Envelope{}, nil, err
		}
		raw = records.MustMarshal(envelope)
	}
	// Intent: Verify embedded semantic author signatures before durable storage
	// so bad author proofs are rejected even outside relay exchange.
	// Source: DI-sovem
	if err := runtime.peers.VerifyAuthorEnvelope(envelope); err != nil {
		return records.Envelope{}, nil, err
	}
	if registered, exists := runtime.families[envelope.Family]; exists {
		if envelope.ProtocolPCID != registered.protocolPCID {
			return records.Envelope{}, nil, fmt.Errorf("family %s expects protocol_pcid %s, got %s", envelope.Family, registered.protocolPCID, envelope.ProtocolPCID)
		}
		owner := registered.owner
		if owner.builtin {
			if validator := owner.validators[envelope.Family]; validator != nil {
				if err := validator(envelope); err != nil {
					return records.Envelope{}, nil, err
				}
			}
		} else {
			if !allowExternalValidation {
				return records.Envelope{}, nil, fmt.Errorf("workflow adapter proposal cannot write externally validated family %s", envelope.Family)
			}
			if err := owner.external.ValidateEnvelope(ctx, raw); err != nil {
				return records.Envelope{}, nil, err
			}
		}
	}
	return envelope, raw, nil
}

func (runtime *Runtime) History() []store.StoredEnvelope {
	return runtime.history.Entries()
}

func (runtime *Runtime) ExportBatch() (grid.Batch, error) {
	entries := runtime.history.Entries()
	rawRecords := make([][]byte, 0, len(entries))
	for _, entry := range entries {
		rawRecords = append(rawRecords, append([]byte{}, entry.Raw...))
	}
	claims := runtime.ImplementationClaims()
	claimProofs, err := runtime.peers.SignClaims(claims)
	if err != nil {
		return grid.Batch{}, err
	}
	recordSignatures, err := runtime.peers.SignRecords(rawRecords)
	if err != nil {
		return grid.Batch{}, err
	}
	// Intent: Export the claim-derived route table alongside claims so relay
	// peers can see the current routing-role model instead of inferring it from
	// local package activation only. Source: DI-ruvot
	return grid.Batch{
		Format:               grid.RelayBatchFormat,
		Implementation:       runtime.LocalPeerID(),
		ExportedAt:           time.Now().UTC().Format(time.RFC3339),
		ImplementationClaims: claims,
		Routes:               runtime.ProtocolRoutes(),
		ClaimProofs:          claimProofs,
		Records:              rawRecords,
		RecordProofs:         grid.ProofsForRecords(rawRecords),
		RecordSignatures:     recordSignatures,
	}, nil
}

func (runtime *Runtime) SignedExportBatch() (grid.Batch, error) {
	batch, err := runtime.ExportBatch()
	if err != nil {
		return grid.Batch{}, err
	}
	return runtime.peers.SignBatch(batch)
}

func (runtime *Runtime) AttestBatchClaims(batch grid.Batch) (grid.Batch, error) {
	if batch.Implementation == runtime.LocalPeerID() {
		return grid.Batch{}, errors.New("third-party claim attestation requires a different peer")
	}
	attestations, err := runtime.peers.AttestClaims(batch.ImplementationClaims)
	if err != nil {
		return grid.Batch{}, err
	}
	batch.ClaimAttestations = append(batch.ClaimAttestations, attestations...)
	return batch, nil
}

func (runtime *Runtime) ImportBatch(ctx context.Context, batch grid.Batch) error {
	// Intent: Treat the current relay shell as idempotent exact-byte carriage so
	// repeated imports stop re-appending identical durable records while malformed
	// batch metadata is rejected before touching local history. Source: DI-sibok
	if err := batch.Validate(); err != nil {
		return err
	}
	// Intent: Verify per-record digests before durable import so receivers can
	// reject tampered relay contents even when they do not yet understand the
	// record family semantics. Source: DI-zumep
	if err := batch.VerifyClaimProofs(); err != nil {
		return err
	}
	if err := batch.VerifyClaimAttestations(); err != nil {
		return err
	}
	// Intent: Verify exported implementation claims against the exporting peer's
	// key material before treating those claims as trustworthy batch metadata.
	// Source: DI-luzef
	if err := runtime.peers.VerifyClaimProofs(batch); err != nil {
		return err
	}
	// Intent: Verify outside countersigners for implementation claims so import
	// can distinguish exporter self-claims from third-party attestation.
	// Source: DI-fogem
	if err := runtime.peers.VerifyClaimAttestations(batch); err != nil {
		return err
	}
	// Intent: Require local claim-attestation quorum only when local runtime
	// policy says a claim needs it, so trust remains explicit and operator-owned.
	// Source: DI-movek
	if err := runtime.policies.VerifyClaimPolicies(batch, runtime.AllowedPeers()); err != nil {
		return err
	}
	if err := batch.VerifyRecordProofs(); err != nil {
		return err
	}
	// Intent: Verify relay-carriage signatures per record so import does not rely
	// solely on the enclosing batch signature for transport trust. Source: DI-ravud
	if err := runtime.peers.VerifyRecordSignatures(batch); err != nil {
		return err
	}
	for _, raw := range batch.Records {
		if _, err := runtime.appendRecord(ctx, raw, false); err != nil {
			return err
		}
	}
	return nil
}

// Intent: Keep live multi-peer exchange behind explicit allow rules and reject
// batches whose claimed peer identity or signature does not match the allowed
// peer entry.
// Source: DI-zotem
func (runtime *Runtime) ImportBatchFromPeer(ctx context.Context, peerID string, batch grid.Batch, direction string) error {
	if strings.TrimSpace(peerID) == "" {
		return errors.New("peer id is required")
	}
	peer, ok := runtime.LookupPeer(peerID)
	if !ok {
		return fmt.Errorf("peer not allowed: %s", peerID)
	}
	switch direction {
	case "pull":
		if !peer.AllowPull {
			return fmt.Errorf("peer %s is not allowed for pull", peerID)
		}
	case "push":
		if !peer.AllowPush {
			return fmt.Errorf("peer %s is not allowed for push", peerID)
		}
	default:
		return fmt.Errorf("unknown peer direction: %s", direction)
	}
	if batch.Implementation != peerID {
		return fmt.Errorf("batch implementation %s does not match peer %s", batch.Implementation, peerID)
	}
	if err := runtime.peers.VerifyPeerBatch(peerID, batch); err != nil {
		return err
	}
	return runtime.ImportBatch(ctx, batch)
}

// Intent: Publish explicit per-package protocol claims so relay exports say
// what the active packages believe they implement instead of implying that through
// local layout or family names alone. Source: DI-lupok
func (runtime *Runtime) ImplementationClaims() []grid.ImplementationClaim {
	claims := []grid.ImplementationClaim{}
	for _, pkg := range runtime.PackageManifests() {
		for _, claim := range pkg.Claims {
			claims = append(claims, grid.ImplementationClaim{
				PackageID:      pkg.ID,
				PackageVersion: pkg.Version,
				ProtocolPCID:   claim.ProtocolPCID,
				Role:           claim.Role,
				RouteType:      claim.NormalizedRouteType(),
				EmitsProtocols: claim.SortedEmitsProtocols(),
				Summary:        claim.Summary,
			})
		}
	}
	return claims
}

func (runtime *Runtime) RunCommand(ctx context.Context, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("command is required")
	}
	for width := len(args); width > 0; width-- {
		key := strings.Join(args[:width], " ")
		owner := runtime.commands[key]
		if owner == nil {
			continue
		}
		if owner.builtin {
			handler := owner.commands[key]
			if handler == nil {
				return "", fmt.Errorf("missing builtin handler for %s", key)
			}
			return handler(ctx, runtime, args[width:])
		}
		result, err := owner.external.RunCommand(ctx, key, args[width:])
		if err != nil {
			return "", err
		}
		return runtime.applyExternalCommandResult(ctx, result)
	}
	return "", fmt.Errorf("unknown command: %s", strings.Join(args, " "))
}

// Intent: Keep durable writes runtime-mediated even for installed agents so
// external packages can extend the system without bypassing runtime-owned CAS
// and append-only history. Source: DI-rovum
func (runtime *Runtime) applyExternalCommandResult(ctx context.Context, result packages.CommandResult) (string, error) {
	return runtime.applyProposedCommandResult(ctx, result, true)
}

// applyWorkflowAdapterResult never starts an installed package executable on
// the host while accepting Docker-worker output. Intent: Confined adapter
// execution must not regain host authority through record validation. Source:
// DI-fofuh
func (runtime *Runtime) applyWorkflowAdapterResult(ctx context.Context, result packages.CommandResult) (string, error) {
	return runtime.applyProposedCommandResult(ctx, result, false)
}

func (runtime *Runtime) applyProposedCommandResult(ctx context.Context, result packages.CommandResult, allowExternalValidation bool) (string, error) {
	prepared, err := runtime.prepareExternalCommandResult(ctx, result, allowExternalValidation)
	if err != nil {
		return "", err
	}
	for _, write := range prepared.CAS {
		if _, err := runtime.PutCAS([]byte(write.Body)); err != nil {
			return "", err
		}
	}
	for _, raw := range prepared.Records {
		if _, _, err := runtime.history.AppendRaw(raw); err != nil {
			return "", err
		}
	}
	return prepared.Output, nil
}

type preparedExternalCommandResult struct {
	Output  string
	CAS     []packages.CASWrite
	Records [][]byte
}

// prepareExternalCommandResult validates every worker-proposed durable write
// before the first CAS or history mutation. Intent: A rejected proposal must
// not leave valid-looking partial worker state behind. Source: DI-fofuh
func (runtime *Runtime) prepareExternalCommandResult(ctx context.Context, result packages.CommandResult, allowExternalValidation bool) (preparedExternalCommandResult, error) {
	replacements := map[string]string{}
	seenAliases := map[string]struct{}{}
	for _, write := range result.CAS {
		if strings.TrimSpace(write.Alias) == "" {
			return preparedExternalCommandResult{}, errors.New("cas alias is required")
		}
		if _, exists := seenAliases[write.Alias]; exists {
			return preparedExternalCommandResult{}, fmt.Errorf("duplicate cas alias: %s", write.Alias)
		}
		seenAliases[write.Alias] = struct{}{}
		replacements["$cas:"+write.Alias] = store.LegacyObjectID([]byte(write.Body))
	}
	prepared := preparedExternalCommandResult{Output: result.Output, CAS: append([]packages.CASWrite{}, result.CAS...), Records: make([][]byte, 0, len(result.Records))}
	for _, raw := range result.Records {
		replaced, err := replaceCASAliases(raw, replacements)
		if err != nil {
			return preparedExternalCommandResult{}, err
		}
		if _, signed, err := runtime.prepareRecord(ctx, replaced, true, allowExternalValidation); err != nil {
			return preparedExternalCommandResult{}, err
		} else {
			prepared.Records = append(prepared.Records, signed)
		}
	}
	return prepared, nil
}

func replaceCASAliases(raw []byte, replacements map[string]string) ([]byte, error) {
	if len(replacements) == 0 {
		return raw, nil
	}
	envelope, err := records.Parse(raw)
	if err != nil {
		return nil, err
	}
	var value any
	if err := json.Unmarshal(envelope.Payload, &value); err != nil {
		return nil, err
	}
	value = replaceAliasesRecursive(value, replacements)
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	envelope.Payload, err = records.CanonicalJSON(payload)
	if err != nil {
		return nil, err
	}
	return records.MustMarshal(envelope), nil
}

func replaceAliasesRecursive(value any, replacements map[string]string) any {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			typed[key] = replaceAliasesRecursive(child, replacements)
		}
		return typed
	case []any:
		for index, child := range typed {
			typed[index] = replaceAliasesRecursive(child, replacements)
		}
		return typed
	case string:
		if replacement, ok := replacements[typed]; ok {
			return replacement
		}
		return typed
	default:
		return value
	}
}

func (runtime *Runtime) InstallPackageDir(ctx context.Context, sourceDir string) (packages.Manifest, error) {
	manifest, err := packages.LoadManifest(filepath.Join(sourceDir, "moks-package.json"))
	if err != nil {
		return packages.Manifest{}, err
	}
	destination := filepath.Join(runtime.packagesRoot, manifest.ID)
	if err := copyDirectory(sourceDir, destination); err != nil {
		return packages.Manifest{}, err
	}
	if err := runtime.ActivateInstalled(ctx, filepath.Join(destination, "moks-package.json")); err != nil {
		return packages.Manifest{}, err
	}
	return manifest, nil
}

// Intent: Re-activate installed agents from the runtime-owned package root on
// startup so installation survives later CLI invocations and process restarts.
// Source: DI-rovum
func (runtime *Runtime) activateInstalledFromRoot(ctx context.Context) error {
	entries, err := os.ReadDir(runtime.packagesRoot)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		manifestPath := filepath.Join(runtime.packagesRoot, entry.Name(), "moks-package.json")
		if _, err := os.Stat(manifestPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
		if err := runtime.ActivateInstalled(ctx, manifestPath); err != nil {
			return err
		}
	}
	return nil
}

func NewEnvelope(family string, protocolPCID string, recordID string, signer string, payload any) (records.Envelope, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return records.Envelope{}, err
	}
	canonicalBody, err := records.CanonicalJSON(body)
	if err != nil {
		return records.Envelope{}, err
	}
	envelope := records.Envelope{
		Family:       family,
		ProtocolPCID: protocolPCID,
		RecordID:     recordID,
		Signer:       signer,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Payload:      canonicalBody,
	}
	return envelope, envelope.Validate()
}
