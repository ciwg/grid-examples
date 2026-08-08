package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
	"github.com/ipfs/go-cid"
)

func TestWorkflowRunsOrderByDurableEventTimeAndRetainV1Heads(t *testing.T) {
	newerTime := time.Unix(200, 0).UTC()
	olderTime := time.Unix(100, 0).UTC()
	newer := WorkflowRun{ID: "newer", Workflow: "newer", UpdatedAt: &newerTime}
	older := WorkflowRun{ID: "older", Workflow: "older", UpdatedAt: &olderTime}
	legacy := WorkflowRun{ID: "legacy", Workflow: "legacy"}
	runtime := &Runtime{workflowRuns: &WorkflowRunRegistry{runs: map[string]WorkflowRun{
		newer.ID:  newer,
		older.ID:  older,
		legacy.ID: legacy,
	}}}
	runs := runtime.WorkflowRuns()
	if len(runs) != 3 || runs[0].ID != newer.ID || runs[1].ID != older.ID || runs[2].ID != legacy.ID {
		t.Fatalf("recent run order = %#v", runs)
	}
	legacyJSON, err := json.Marshal(legacy)
	if err != nil {
		t.Fatalf("marshal retained v1 run: %v", err)
	}
	if strings.Contains(string(legacyJSON), "updated_at") {
		t.Fatalf("retained v1 JSON reports an invented time: %s", legacyJSON)
	}

	root := t.TempDir()
	cas, err := store.OpenCAS(root)
	if err != nil {
		t.Fatalf("open CAS: %v", err)
	}
	nonce, err := records.EncodeGrid(records.GridEnvelope{ProtocolPCID: workflowRunProtocolV1CID, Slots: []any{"workflow-run-nonce", make([]byte, 32)}})
	if err != nil {
		t.Fatalf("encode retained v1 nonce: %v", err)
	}
	runCID, err := cas.PutCID(nonce)
	if err != nil {
		t.Fatalf("store retained v1 nonce: %v", err)
	}
	handoff, err := EncodeWorkflowHandoff(WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "legacy"}})
	if err != nil {
		t.Fatalf("encode retained v1 handoff: %v", err)
	}
	inputCID, err := cas.PutCID(handoff)
	if err != nil {
		t.Fatalf("store retained v1 handoff: %v", err)
	}
	legacyRaw, err := records.EncodeGrid(records.GridEnvelope{ProtocolPCID: workflowRunProtocolV1CID, Slots: []any{
		string(WorkflowRunRunning), runCID.Bytes(), "legacy", inputCID.Bytes(), []byte{}, "", []any{},
	}})
	if err != nil {
		t.Fatalf("encode retained v1 event: %v", err)
	}
	decoded, err := decodeWorkflowRunEvent(legacyRaw)
	if err != nil {
		t.Fatalf("decode retained v1 event: %v", err)
	}
	if !decoded.RecordedAt.IsZero() {
		t.Fatalf("retained v1 time = %s, want zero", decoded.RecordedAt)
	}
	if _, err := cas.PutCID(legacyRaw); err != nil {
		t.Fatalf("store retained v1 event: %v", err)
	}
	registry, err := OpenWorkflowRunRegistry(filepath.Join(root, "state"), cas)
	if err != nil {
		t.Fatalf("replay retained v1 event: %v", err)
	}
	if run, ok := registry.get(runCID.String()); !ok || run.UpdatedAt != nil {
		t.Fatalf("replayed retained v1 run = %#v, found=%t", run, ok)
	}
}

func TestWorkflowRunRequiresActiveArtifact(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := runtime.RegisterWorkflowOperation("test", func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) { return input, nil }); err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(root, "workflow")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"id":"test","version":"1","summary":"test","required_packages":[],"required_protocols":[],"adapter":"test","input_pcid":"` + WorkflowHandoffProtocolPCID + `","output_pcid":"` + WorkflowHandoffProtocolPCID + `"}`
	if err := os.WriteFile(filepath.Join(directory, "workflow.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.CaptureWorkflowDir(directory, "test"); err != nil {
		t.Fatal(err)
	}
	input := WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}}
	if _, err := runtime.StartWorkflowRun(context.Background(), "test", input); err == nil {
		t.Fatal("inactive workflow started")
	}
	if err := runtime.ActivateWorkflow("test"); err != nil {
		t.Fatal(err)
	}
	run, err := runtime.StartWorkflowRun(context.Background(), "test", input)
	if err != nil {
		t.Fatal(err)
	}
	if run.State != WorkflowRunCompleted {
		t.Fatalf("state = %s", run.State)
	}
	if err := runtime.DeactivateWorkflow("test"); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.StartWorkflowRun(context.Background(), "test", input); err == nil {
		t.Fatal("deactivated workflow started")
	}
}

func TestRevokedWorkflowRemainsIneligibleAfterRestart(t *testing.T) {
	// Intent: A revoked workflow remains retained evidence, but restart must not
	// restore local activation or execution eligibility. Source: DI-rupit.
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	operation := func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) { return input, nil }
	if err := runtime.RegisterWorkflowOperation("test", operation); err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(root, "workflow")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"id":"test","version":"1","summary":"test","required_packages":[],"required_protocols":[],"adapter":"test","input_pcid":"` + WorkflowHandoffProtocolPCID + `","output_pcid":"` + WorkflowHandoffProtocolPCID + `"}`
	if err := os.WriteFile(filepath.Join(directory, "workflow.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	imported, err := runtime.CaptureWorkflowDir(directory, "test")
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.ActivateWorkflow("test"); err != nil {
		t.Fatal(err)
	}
	if err := runtime.RevokeWorkflow("test"); err != nil {
		t.Fatal(err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := reopened.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := reopened.RegisterWorkflowOperation("test", operation); err != nil {
		t.Fatal(err)
	}
	workflows := reopened.Workflows()
	if len(workflows) != 1 || workflows[0].ArtifactCID != imported.ArtifactCID || workflows[0].State != WorkflowRevoked {
		t.Fatalf("replayed workflows = %#v, want retained revoked artifact %q", workflows, imported.ArtifactCID)
	}
	if err := reopened.ActivateWorkflow("test"); err == nil {
		t.Fatal("revoked workflow activated after restart")
	}
	input := WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}}
	if _, err := reopened.StartWorkflowRun(context.Background(), "test", input); err == nil {
		t.Fatal("revoked workflow started after restart")
	}
}

func TestWorkflowHandoffRequiresActiveTarget(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	operation := func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) { return input, nil }
	for _, name := range []string{"source", "target"} {
		if err := runtime.RegisterWorkflowOperation(name, operation); err != nil {
			t.Fatal(err)
		}
		directory := filepath.Join(root, name)
		if err := os.Mkdir(directory, 0755); err != nil {
			t.Fatal(err)
		}
		manifest := `{"id":"` + name + `","version":"1","summary":"test","required_packages":[],"required_protocols":[],"adapter":"` + name + `","input_pcid":"` + WorkflowHandoffProtocolPCID + `","output_pcid":"` + WorkflowHandoffProtocolPCID + `"}`
		if err := os.WriteFile(filepath.Join(directory, "workflow.json"), []byte(manifest), 0644); err != nil {
			t.Fatal(err)
		}
		if _, err := runtime.CaptureWorkflowDir(directory, name); err != nil {
			t.Fatal(err)
		}
	}
	if err := runtime.ActivateWorkflow("source"); err != nil {
		t.Fatal(err)
	}
	source, err := runtime.StartWorkflowRun(context.Background(), "source", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.HandoffWorkflowRun(context.Background(), source.ID, "target"); err == nil {
		t.Fatal("inactive target accepted handoff")
	}
	if err := runtime.ActivateWorkflow("target"); err != nil {
		t.Fatal(err)
	}
	target, err := runtime.HandoffWorkflowRun(context.Background(), source.ID, "target")
	if err != nil {
		t.Fatal(err)
	}
	if target.State != WorkflowRunCompleted {
		t.Fatalf("target state = %s", target.State)
	}
	if err := runtime.DeactivateWorkflow("source"); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.HandoffWorkflowRun(context.Background(), source.ID, "target"); err == nil {
		t.Fatal("deactivated source handed off work")
	}
}

func TestWorkflowRunRebuildsAfterRestartAndDeletedCache(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.RegisterWorkflowOperation("test", func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) { return input, nil }); err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(root, "workflow")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"id":"test","version":"1","summary":"test","required_packages":[],"required_protocols":[],"adapter":"test","input_pcid":"` + WorkflowHandoffProtocolPCID + `","output_pcid":"` + WorkflowHandoffProtocolPCID + `"}`
	if err := os.WriteFile(filepath.Join(directory, "workflow.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.CaptureWorkflowDir(directory, "test"); err != nil {
		t.Fatal(err)
	}
	if err := runtime.ActivateWorkflow("test"); err != nil {
		t.Fatal(err)
	}
	run, err := runtime.StartWorkflowRun(context.Background(), "test", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, "state", "workflow-run-cache.json")); err != nil {
		t.Fatal(err)
	}
	runtime, err = Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	replayed, err := runtime.WorkflowRun(run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if replayed.State != WorkflowRunCompleted || replayed.OutputCID == "" || replayed.UpdatedAt == nil {
		t.Fatalf("replayed run = %#v", replayed)
	}
}

func FuzzDecodeWorkflowHandoff(f *testing.F) {
	raw, err := EncodeWorkflowHandoff(WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		f.Fatal(err)
	}
	f.Add(raw)
	f.Fuzz(func(t *testing.T, body []byte) { _, _ = DecodeWorkflowHandoff(body) })
}

func TestWorkflowRunWaitSupplyRetryAndPolicyHandoff(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	for _, name := range []string{"source", "target"} {
		adapter := func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) {
			if input.Values["subject"] == "" {
				return WorkflowHandoff{}, WaitingForWorkflowInput("subject required")
			}
			return input, nil
		}
		if err := runtime.RegisterWorkflowOperation(name, adapter); err != nil {
			t.Fatal(err)
		}
		writeExecutableWorkflow(t, root, runtime, name)
		if err := runtime.ActivateWorkflow(name); err != nil {
			t.Fatal(err)
		}
	}
	waiting, err := runtime.StartWorkflowRun(context.Background(), "source", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"placeholder": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if waiting.State != WorkflowRunWaiting {
		t.Fatalf("waiting state = %s", waiting.State)
	}
	completed, err := runtime.SupplyWorkflowRun(context.Background(), waiting.ID, WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if completed.State != WorkflowRunCompleted {
		t.Fatalf("supplied state = %s", completed.State)
	}
	if err := runtime.SetWorkflowHandoffPolicy(WorkflowHandoffPolicy{SourceWorkflow: "source", OutputPCID: WorkflowHandoffProtocolPCID, TargetWorkflow: "target", InputPCID: WorkflowHandoffProtocolPCID}); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.StartWorkflowRun(context.Background(), "source", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "policy egg"}}); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, run := range runtime.WorkflowRuns() {
		if run.Workflow == "target" && run.State == WorkflowRunCompleted {
			found = true
		}
	}
	if !found {
		t.Fatal("policy did not start target")
	}
}

func TestWorkflowAllToAllHandoffMatrix(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	workflows := []string{"inventory-discrepancy-review", "inventory-receipt", "knowledge-review", "maintenance-round", "procedure-execution", "receiving-check", "training-qualification"}
	inputs := []string{"bafkreigwhgyyvdxkjrckvimjhesxg2ms2wtahqcogu7276xqovcjoxif3e", "bafkreie3xn5cs7in24a5aenl7kpyaa22e346wr4tcqm4evxgcn2v55yvne", "bafkreigurwnuwdbri3ntiuzseyoft7tazq74rrxridxbpekqcqxd2rzvhi", "bafkreicuxiyha56khoiwfktrh3pqrq6v47uovjm52kmuskmkbeml6qyxsu", "bafkreiawxq2i7q57tks6f5viofxkko2jf2txmlurbp3i33svynytyjswfq", "bafkreicubq2eqovbcdvlj3ggpmrnwlzohp5no5yuwehk4eyucicia4kehy", "bafkreigdqszsce7qvcihohhsgr4r5wh7l3lgbqm3gd55befrf7d4hxpwyi"}
	outputs := []string{"bafkreichcquo564amype7v6locjdhc7xl5kgb6i7oyo25k65o677kyztey", "bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne", "bafkreigxtnvprgom4b6ftkavmvfeqv2teo45b7s3n57k3wn7ipziwlad4e", "bafkreifj3qpjinq4vanr5vo22dvazuybg6jrrawnfhw43pooqnzp62vtg4", "bafkreiamprv3apzowjzqbkp3hnrhrla5aq7lp5kyzbca5j3iv4v5jmhwa4", "bafkreiddsvg5v7a2dwa4omzicgobed5yqheehbkjw263kb3efvv2yfgnz4", "bafkreiauh6xo45sp3zhmhhjvehiqstp36loiyu26smicojfntl7ek75chy"}
	for index, name := range workflows {
		output := outputs[index]
		if err := runtime.RegisterWorkflowOperation(name, func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) {
			return WorkflowHandoff{PCID: output, Values: input.Values}, nil
		}); err != nil {
			t.Fatal(err)
		}
		writeSchemaWorkflow(t, root, runtime, name, inputs[index], output)
		if err := runtime.ActivateWorkflow(name); err != nil {
			t.Fatal(err)
		}
	}
	for sourceIndex, source := range workflows {
		run, err := runtime.StartWorkflowRun(context.Background(), source, WorkflowHandoff{PCID: inputs[sourceIndex], Values: map[string]string{"subject": source}})
		if err != nil {
			t.Fatalf("start %s: %v", source, err)
		}
		for _, target := range workflows {
			if source == target {
				continue
			}
			handed, err := runtime.HandoffWorkflowRun(context.Background(), run.ID, target)
			if err != nil {
				t.Fatalf("%s -> %s: %v", source, target, err)
			}
			if handed.State != WorkflowRunWaiting {
				t.Fatalf("%s -> %s state = %s", source, target, handed.State)
			}
		}
	}
}

func TestWorkflowPolicyRoutesDistinctSchemasToWaitingInput(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	const sourceInput = "bafkreigwhgyyvdxkjrckvimjhesxg2ms2wtahqcogu7276xqovcjoxif3e"
	const sourceOutput = "bafkreichcquo564amype7v6locjdhc7xl5kgb6i7oyo25k65o677kyztey"
	const targetInput = "bafkreie3xn5cs7in24a5aenl7kpyaa22e346wr4tcqm4evxgcn2v55yvne"
	const targetOutput = "bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne"
	for _, schema := range []struct{ name, input, output string }{{"source", sourceInput, sourceOutput}, {"target", targetInput, targetOutput}} {
		output := schema.output
		if err := runtime.RegisterWorkflowOperation(schema.name, func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) {
			return WorkflowHandoff{PCID: output, Values: input.Values}, nil
		}); err != nil {
			t.Fatal(err)
		}
		writeSchemaWorkflow(t, root, runtime, schema.name, schema.input, schema.output)
		if err := runtime.ActivateWorkflow(schema.name); err != nil {
			t.Fatal(err)
		}
	}
	if err := runtime.SetWorkflowHandoffPolicy(WorkflowHandoffPolicy{SourceWorkflow: "source", OutputPCID: sourceOutput, TargetWorkflow: "target", InputPCID: targetInput}); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.StartWorkflowRun(context.Background(), "source", WorkflowHandoff{PCID: sourceInput, Values: map[string]string{"subject": "egg"}}); err != nil {
		t.Fatal(err)
	}
	var waiting WorkflowRun
	for _, run := range runtime.WorkflowRuns() {
		if run.Workflow == "target" {
			waiting = run
		}
	}
	if waiting.State != WorkflowRunWaiting {
		t.Fatalf("policy target state = %s", waiting.State)
	}
	completed, err := runtime.SupplyWorkflowRun(context.Background(), waiting.ID, WorkflowHandoff{PCID: targetInput, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if completed.State != WorkflowRunCompleted {
		t.Fatalf("supplied policy target state = %s", completed.State)
	}
}

func TestWorkflowIncompatibleHandoffRetainsExactSourceOutput(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	const sourceInput = "bafkreigwhgyyvdxkjrckvimjhesxg2ms2wtahqcogu7276xqovcjoxif3e"
	const sourceOutput = "bafkreichcquo564amype7v6locjdhc7xl5kgb6i7oyo25k65o677kyztey"
	const targetInput = "bafkreie3xn5cs7in24a5aenl7kpyaa22e346wr4tcqm4evxgcn2v55yvne"
	const targetOutput = "bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne"
	for _, schema := range []struct{ name, input, output string }{{"source", sourceInput, sourceOutput}, {"target", targetInput, targetOutput}} {
		output := schema.output
		if err := runtime.RegisterWorkflowOperation(schema.name, func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) {
			return WorkflowHandoff{PCID: output, Values: input.Values}, nil
		}); err != nil {
			t.Fatal(err)
		}
		writeSchemaWorkflow(t, root, runtime, schema.name, schema.input, schema.output)
		if err := runtime.ActivateWorkflow(schema.name); err != nil {
			t.Fatal(err)
		}
	}
	source, err := runtime.StartWorkflowRun(context.Background(), "source", WorkflowHandoff{PCID: sourceInput, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	waiting, err := runtime.HandoffWorkflowRun(context.Background(), source.ID, "target")
	if err != nil {
		t.Fatal(err)
	}
	if waiting.State != WorkflowRunWaiting {
		t.Fatalf("handoff state = %s", waiting.State)
	}
	if waiting.InputCID != source.OutputCID {
		t.Fatalf("waiting input %s did not retain source output %s", waiting.InputCID, source.OutputCID)
	}
}

func TestWorkflowIncompatibleHandoffRequiresAvailableTargetAdapter(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	const sourceInput = "bafkreigwhgyyvdxkjrckvimjhesxg2ms2wtahqcogu7276xqovcjoxif3e"
	const sourceOutput = "bafkreichcquo564amype7v6locjdhc7xl5kgb6i7oyo25k65o677kyztey"
	const targetInput = "bafkreie3xn5cs7in24a5aenl7kpyaa22e346wr4tcqm4evxgcn2v55yvne"
	const targetOutput = "bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne"
	if err := runtime.RegisterWorkflowOperation("source", func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) {
		return WorkflowHandoff{PCID: sourceOutput, Values: input.Values}, nil
	}); err != nil {
		t.Fatal(err)
	}
	writeSchemaWorkflow(t, root, runtime, "source", sourceInput, sourceOutput)
	writeSchemaWorkflow(t, root, runtime, "target", targetInput, targetOutput)
	if err := runtime.ActivateWorkflow("source"); err != nil {
		t.Fatal(err)
	}
	if err := runtime.ActivateWorkflow("target"); err != nil {
		t.Fatal(err)
	}
	source, err := runtime.StartWorkflowRun(context.Background(), "source", WorkflowHandoff{PCID: sourceInput, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.HandoffWorkflowRun(context.Background(), source.ID, "target"); err == nil {
		t.Fatal("incompatible handoff queued target without an available adapter")
	}
}

func TestWorkflowPolicyRejectsManifestPCIDMismatch(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	const sourceInput = "bafkreigwhgyyvdxkjrckvimjhesxg2ms2wtahqcogu7276xqovcjoxif3e"
	const sourceOutput = "bafkreichcquo564amype7v6locjdhc7xl5kgb6i7oyo25k65o677kyztey"
	const targetInput = "bafkreie3xn5cs7in24a5aenl7kpyaa22e346wr4tcqm4evxgcn2v55yvne"
	const targetOutput = "bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne"
	for _, schema := range []struct{ name, input, output string }{{"source", sourceInput, sourceOutput}, {"target", targetInput, targetOutput}} {
		if err := runtime.RegisterWorkflowOperation(schema.name, func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) {
			return input, nil
		}); err != nil {
			t.Fatal(err)
		}
		writeSchemaWorkflow(t, root, runtime, schema.name, schema.input, schema.output)
	}
	if err := runtime.SetWorkflowHandoffPolicy(WorkflowHandoffPolicy{SourceWorkflow: "source", OutputPCID: sourceInput, TargetWorkflow: "target", InputPCID: targetInput}); err == nil {
		t.Fatal("policy accepted source input as its output schema")
	}
	if err := runtime.SetWorkflowHandoffPolicy(WorkflowHandoffPolicy{SourceWorkflow: "source", OutputPCID: sourceOutput, TargetWorkflow: "target", InputPCID: targetOutput}); err == nil {
		t.Fatal("policy accepted target output as its input schema")
	}
}

func TestWorkflowPolicyPersistsAcrossRestart(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	operation := func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) { return input, nil }
	for _, name := range []string{"source", "target"} {
		if err := runtime.RegisterWorkflowOperation(name, operation); err != nil {
			t.Fatal(err)
		}
		writeExecutableWorkflow(t, root, runtime, name)
		if err := runtime.ActivateWorkflow(name); err != nil {
			t.Fatal(err)
		}
	}
	policy := WorkflowHandoffPolicy{SourceWorkflow: "source", OutputPCID: WorkflowHandoffProtocolPCID, TargetWorkflow: "target", InputPCID: WorkflowHandoffProtocolPCID}
	if err := runtime.SetWorkflowHandoffPolicy(policy); err != nil {
		t.Fatal(err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}
	runtime, err = Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	for _, name := range []string{"source", "target"} {
		if err := runtime.RegisterWorkflowOperation(name, operation); err != nil {
			t.Fatal(err)
		}
	}
	if policies := runtime.WorkflowHandoffPolicies(); len(policies) != 1 || policies[0] != policy {
		t.Fatalf("replayed policies = %#v", policies)
	}
	if _, err := runtime.StartWorkflowRun(context.Background(), "source", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}}); err != nil {
		t.Fatal(err)
	}
	for _, run := range runtime.WorkflowRuns() {
		if run.Workflow == "target" && run.State == WorkflowRunCompleted {
			return
		}
	}
	t.Fatal("replayed policy did not route the target workflow")
}

func TestWorkflowRunRetryRetainsInput(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	attempts := 0
	if err := runtime.RegisterWorkflowOperation("retry", func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) {
		attempts++
		if attempts == 1 {
			return WorkflowHandoff{}, errors.New("temporary adapter failure")
		}
		return input, nil
	}); err != nil {
		t.Fatal(err)
	}
	writeExecutableWorkflow(t, root, runtime, "retry")
	if err := runtime.ActivateWorkflow("retry"); err != nil {
		t.Fatal(err)
	}
	failed, err := runtime.StartWorkflowRun(context.Background(), "retry", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if failed.State != WorkflowRunFailed {
		t.Fatalf("failed state = %s", failed.State)
	}
	retried, err := runtime.RetryWorkflowRun(context.Background(), failed.ID)
	if err != nil {
		t.Fatal(err)
	}
	if retried.State != WorkflowRunCompleted {
		t.Fatalf("retry state = %s", retried.State)
	}
	if retried.InputCID != failed.InputCID {
		t.Fatalf("retry input = %s, want %s", retried.InputCID, failed.InputCID)
	}
}

func TestWorkflowRunsWithIdenticalInputRemainDistinct(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := runtime.RegisterWorkflowOperation("repeat", func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) { return input, nil }); err != nil {
		t.Fatal(err)
	}
	writeExecutableWorkflow(t, root, runtime, "repeat")
	if err := runtime.ActivateWorkflow("repeat"); err != nil {
		t.Fatal(err)
	}
	input := WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "same egg"}}
	first, err := runtime.StartWorkflowRun(context.Background(), "repeat", input)
	if err != nil {
		t.Fatal(err)
	}
	second, err := runtime.StartWorkflowRun(context.Background(), "repeat", input)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID == second.ID {
		t.Fatal("identical input collapsed into one run")
	}
}

func TestWorkflowHandoffPolicyRejectsCycle(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	for _, name := range []string{"a", "b"} {
		if err := runtime.RegisterWorkflowOperation(name, func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) { return input, nil }); err != nil {
			t.Fatal(err)
		}
		writeExecutableWorkflow(t, root, runtime, name)
	}
	first := WorkflowHandoffPolicy{SourceWorkflow: "a", OutputPCID: WorkflowHandoffProtocolPCID, TargetWorkflow: "b", InputPCID: WorkflowHandoffProtocolPCID}
	if err := runtime.SetWorkflowHandoffPolicy(first); err != nil {
		t.Fatal(err)
	}
	second := WorkflowHandoffPolicy{SourceWorkflow: "b", OutputPCID: WorkflowHandoffProtocolPCID, TargetWorkflow: "a", InputPCID: WorkflowHandoffProtocolPCID}
	if err := runtime.SetWorkflowHandoffPolicy(second); err == nil {
		t.Fatal("cycle policy accepted")
	}
}

func TestWorkflowRunReplayRejectsForeignParent(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	rootCID, err := newWorkflowRunCID(runtime.cas)
	if err != nil {
		t.Fatal(err)
	}
	inputCID, err := runtime.cas.PutCID(mustWorkflowHandoff(t))
	if err != nil {
		t.Fatal(err)
	}
	startRaw, err := encodeWorkflowRunEvent(workflowRunEvent{RunCID: rootCID, Workflow: "source", State: WorkflowRunRunning, Input: inputCID})
	if err != nil {
		t.Fatal(err)
	}
	startCID, err := runtime.cas.PutCID(startRaw)
	if err != nil {
		t.Fatal(err)
	}
	foreignCID, err := newWorkflowRunCID(runtime.cas)
	if err != nil {
		t.Fatal(err)
	}
	foreignRaw, err := encodeWorkflowRunEvent(workflowRunEvent{RunCID: foreignCID, Workflow: "target", State: WorkflowRunFailed, Input: inputCID, Reason: "forged", Parents: []cid.Cid{startCID}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.cas.PutCID(foreignRaw); err != nil {
		t.Fatal(err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}
	runtime, err = Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if _, err := runtime.WorkflowRun(foreignCID.String()); err == nil {
		t.Fatal("foreign-parent run replayed")
	}
	if _, err := runtime.WorkflowRun(rootCID.String()); err != nil {
		t.Fatal(err)
	}
}

func TestWorkflowRunReplayRequiresRetainedNonce(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	inputCID, err := runtime.cas.PutCID(mustWorkflowHandoff(t))
	if err != nil {
		t.Fatal(err)
	}
	forgedRaw, err := encodeWorkflowRunEvent(workflowRunEvent{RunCID: inputCID, Workflow: "forged", State: WorkflowRunRunning, Input: inputCID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.cas.PutCID(forgedRaw); err != nil {
		t.Fatal(err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}
	runtime, err = Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if _, err := runtime.WorkflowRun(inputCID.String()); err == nil {
		t.Fatal("root without a retained nonce replayed")
	}
}

func TestFailWorkflowRunPersistsTerminalFailure(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	runCID, err := newWorkflowRunCID(runtime.cas)
	if err != nil {
		t.Fatal(err)
	}
	inputCID, err := runtime.cas.PutCID(mustWorkflowHandoff(t))
	if err != nil {
		t.Fatal(err)
	}
	running, err := runtime.workflowRuns.append(workflowRunEvent{RunCID: runCID, Workflow: "test", State: WorkflowRunRunning, Input: inputCID})
	if err != nil {
		t.Fatal(err)
	}
	failed, err := runtime.failWorkflowRun(running, "output persistence", errors.New("disk full"))
	if err != nil {
		t.Fatal(err)
	}
	if failed.State != WorkflowRunFailed {
		t.Fatalf("state = %s", failed.State)
	}
}

func TestWorkflowRunManualRecoveryUsesDedicatedRunningEdge(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	runCID, err := newWorkflowRunCID(runtime.cas)
	if err != nil {
		t.Fatal(err)
	}
	inputCID, err := runtime.cas.PutCID(mustWorkflowHandoff(t))
	if err != nil {
		t.Fatal(err)
	}
	running, err := runtime.workflowRuns.append(workflowRunEvent{RunCID: runCID, Workflow: "test", State: WorkflowRunRunning, Input: inputCID})
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := runtime.resumeWorkflowRun(running, inputCID)
	if err != nil {
		t.Fatal(err)
	}
	if recovered.State != WorkflowRunRunning || recovered.Reason != "manual recovery" {
		t.Fatalf("recovery = %#v", recovered)
	}
	if validWorkflowRunTransition(WorkflowRunRunning, WorkflowRunRunning, "") {
		t.Fatal("unmarked running transition accepted")
	}
}

func TestOpenRejectsMalformedPersistedHandoffPolicy(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "state"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "state", "workflow-handoff-policy.json"), []byte(`[{"source_workflow":"a","output_pcid":"x","target_workflow":"b","input_pcid":"y"}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(root); err == nil {
		t.Fatal("malformed persisted policy opened")
	}
}

func TestWorkflowHandoffPolicyWriteFailureLeavesLivePolicyUnchanged(t *testing.T) {
	root := t.TempDir()
	policies, err := OpenWorkflowHandoffPolicies(root)
	if err != nil {
		t.Fatal(err)
	}
	first := WorkflowHandoffPolicy{SourceWorkflow: "source", OutputPCID: WorkflowHandoffProtocolPCID, TargetWorkflow: "target", InputPCID: WorkflowHandoffProtocolPCID}
	if err := policies.Set(first); err != nil {
		t.Fatal(err)
	}
	policies.path = root
	second := WorkflowHandoffPolicy{SourceWorkflow: "other", OutputPCID: WorkflowHandoffProtocolPCID, TargetWorkflow: "final", InputPCID: WorkflowHandoffProtocolPCID}
	if err := policies.Set(second); err == nil {
		t.Fatal("policy write to directory succeeded")
	}
	if actual := policies.List(); len(actual) != 1 || actual[0] != first {
		t.Fatalf("live policies changed after write failure: %#v", actual)
	}
	if err := policies.Remove(first.SourceWorkflow, first.OutputPCID); err == nil {
		t.Fatal("policy removal through directory path succeeded")
	}
	if actual := policies.List(); len(actual) != 1 || actual[0] != first {
		t.Fatalf("live policies changed after failed removal: %#v", actual)
	}
}

func TestWorkflowHandoffPolicySyncFailureKeepsLivePolicyAlignedWithFile(t *testing.T) {
	root := t.TempDir()
	policies, err := OpenWorkflowHandoffPolicies(root)
	if err != nil {
		t.Fatal(err)
	}
	policies.syncDirectory = func(string) error { return errors.New("directory sync failed") }
	policy := WorkflowHandoffPolicy{SourceWorkflow: "source", OutputPCID: WorkflowHandoffProtocolPCID, TargetWorkflow: "target", InputPCID: WorkflowHandoffProtocolPCID}
	if err := policies.Set(policy); err == nil {
		t.Fatal("policy write with failed directory sync succeeded")
	}
	if actual := policies.List(); len(actual) != 1 || actual[0] != policy {
		t.Fatalf("live policy diverged after rename: %#v", actual)
	}
	reopened, err := OpenWorkflowHandoffPolicies(root)
	if err != nil {
		t.Fatal(err)
	}
	if actual := reopened.List(); len(actual) != 1 || actual[0] != policy {
		t.Fatalf("persisted policy diverged after rename: %#v", actual)
	}
}

func TestRuntimeOpenIgnoresDisposableWorkflowCacheWriteFailures(t *testing.T) {
	root := t.TempDir()
	runtime, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.RegisterWorkflowOperation("test", func(_ context.Context, _ *Runtime, input WorkflowHandoff) (WorkflowHandoff, error) { return input, nil }); err != nil {
		t.Fatal(err)
	}
	writeExecutableWorkflow(t, root, runtime, "test")
	if err := runtime.ActivateWorkflow("test"); err != nil {
		t.Fatal(err)
	}
	run, err := runtime.StartWorkflowRun(context.Background(), "test", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{filepath.Join(root, "state", "workflow-lifecycle-cache.json"), filepath.Join(root, "state", "workflow-run-cache.json")} {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	runtime, err = Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if _, err := runtime.WorkflowRun(run.ID); err != nil {
		t.Fatal(err)
	}
	if workflow, err := runtime.workflow("test"); err != nil || workflow.State != WorkflowActive {
		t.Fatalf("replayed workflow = %#v, %v", workflow, err)
	}
}

func TestCompletedWorkflowRunEventRequiresOutput(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	runCID, err := newWorkflowRunCID(runtime.cas)
	if err != nil {
		t.Fatal(err)
	}
	inputCID, err := runtime.cas.PutCID(mustWorkflowHandoff(t))
	if err != nil {
		t.Fatal(err)
	}
	parentRaw, err := encodeWorkflowRunEvent(workflowRunEvent{RunCID: runCID, Workflow: "test", State: WorkflowRunRunning, Input: inputCID})
	if err != nil {
		t.Fatal(err)
	}
	parentCID, err := runtime.cas.PutCID(parentRaw)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := encodeWorkflowRunEvent(workflowRunEvent{RunCID: runCID, Workflow: "test", State: WorkflowRunCompleted, Input: inputCID, Parents: []cid.Cid{parentCID}}); err == nil {
		t.Fatal("completed event without output encoded")
	}
}

func TestWorkflowRunRegistryRejectsStaleAndIllegalTransitions(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	runCID, err := newWorkflowRunCID(runtime.cas)
	if err != nil {
		t.Fatal(err)
	}
	inputCID, err := runtime.cas.PutCID(mustWorkflowHandoff(t))
	if err != nil {
		t.Fatal(err)
	}
	start, err := runtime.workflowRuns.append(workflowRunEvent{RunCID: runCID, Workflow: "test", State: WorkflowRunRunning, Input: inputCID})
	if err != nil {
		t.Fatal(err)
	}
	startEvent, err := cid.Decode(start.EventCID)
	if err != nil {
		t.Fatal(err)
	}
	waiting, err := runtime.workflowRuns.append(workflowRunEvent{RunCID: runCID, Workflow: "test", State: WorkflowRunWaiting, Input: inputCID, Reason: "need input", Parents: []cid.Cid{startEvent}})
	if err != nil {
		t.Fatal(err)
	}
	waitingEvent, err := cid.Decode(waiting.EventCID)
	if err != nil {
		t.Fatal(err)
	}
	resume := workflowRunEvent{RunCID: runCID, Workflow: "test", State: WorkflowRunRunning, Input: inputCID, Parents: []cid.Cid{waitingEvent}}
	if _, err := runtime.workflowRuns.append(resume); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.workflowRuns.append(resume); err == nil {
		t.Fatal("stale transition accepted")
	}
	running, err := runtime.WorkflowRun(runCID.String())
	if err != nil {
		t.Fatal(err)
	}
	completed, err := runtime.transitionWorkflowRun(running, WorkflowRunCompleted, inputCID, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.transitionWorkflowRun(completed, WorkflowRunWaiting, cid.Undef, "forged"); err == nil {
		t.Fatal("illegal completed transition accepted")
	}
}

func mustWorkflowHandoff(t *testing.T) []byte {
	t.Helper()
	raw, err := EncodeWorkflowHandoff(WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestInstalledWorkflowAdapterExecutesValidatedWorkerProposal(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	writeExecutableWorkflow(t, runtime.root, runtime, "worker")
	if err := runtime.ActivateWorkflow("worker"); err != nil {
		t.Fatal(err)
	}
	runtime.workflowAdapters["worker"] = packages.WorkflowAdapter{Name: "worker", Image: "example/worker:1", InputPCID: WorkflowHandoffProtocolPCID, OutputPCID: WorkflowHandoffProtocolPCID, CPUs: "0.5", Memory: "128m", PIDsLimit: 64, Timeout: "30s"}
	runtime.workflowAdapterPackages["worker"] = "worker-package"
	if _, err := runtime.StartWorkflowRun(context.Background(), "worker", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}}); err == nil || !strings.Contains(err.Error(), "not currently promised") {
		t.Fatalf("missing route evidence error = %v", err)
	}
	if len(runtime.WorkflowRuns()) != 0 {
		t.Fatalf("pre-dispatch refusal created runs: %#v", runtime.WorkflowRuns())
	}
	if err := runtime.BindAgent(AgentBinding{AgentID: "worker-app", PackageID: "worker-package", Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := runtime.PublishReceivePromise(ReceivePromise{AgentID: "worker-app", ProtocolPCID: WorkflowHandoffProtocolPCID, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := runtime.PublishDeliveryPromise(DeliveryPromise{AgentID: "local-router", RecipientAgentID: "worker-app", ProtocolPCID: WorkflowHandoffProtocolPCID, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	runtime.workflowWorker = func(_ context.Context, adapter packages.WorkflowAdapter, raw []byte) ([]byte, error) {
		if adapter.Name != "worker" {
			t.Fatalf("unexpected adapter: %#v", adapter)
		}
		input, err := DecodeWorkflowHandoff(raw)
		if err != nil {
			t.Fatal(err)
		}
		if input.Values["subject"] != "egg" {
			t.Fatalf("unexpected worker input: %#v", input)
		}
		return EncodeWorkflowAdapterResult(WorkflowAdapterResult{
			Output:  input,
			CAS:     []packages.CASWrite{{Alias: "body", Body: "worker evidence"}},
			Records: []json.RawMessage{json.RawMessage(`{"family":"worker.note.v1","protocol_pcid":"pcid:worker.note.v1","record_id":"worker-1","signer":"worker","timestamp":"2026-08-02T00:00:00Z","payload":{"body_ref":"$cas:body"}}`)},
		})
	}
	run, err := runtime.StartWorkflowRun(context.Background(), "worker", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if run.State != WorkflowRunCompleted || run.OutputCID == "" {
		t.Fatalf("worker run = %#v", run)
	}
	if len(runtime.History()) != 1 || runtime.History()[0].Envelope.Family != "worker.note.v1" {
		t.Fatalf("worker proposal was not runtime-mediated: %#v", runtime.History())
	}
}

func TestInstalledWorkflowAdapterRejectsWrongOutputBeforeWrites(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	writeExecutableWorkflow(t, runtime.root, runtime, "worker")
	if err := runtime.ActivateWorkflow("worker"); err != nil {
		t.Fatal(err)
	}
	runtime.workflowAdapters["worker"] = packages.WorkflowAdapter{Name: "worker", Image: "example/worker:1", InputPCID: WorkflowHandoffProtocolPCID, OutputPCID: WorkflowHandoffProtocolPCID, CPUs: "0.5", Memory: "128m", PIDsLimit: 64, Timeout: "30s"}
	runtime.workflowAdapterPackages["worker"] = "worker-package"
	if err := runtime.BindAgent(AgentBinding{AgentID: "worker-app", PackageID: "worker-package", Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := runtime.PublishReceivePromise(ReceivePromise{AgentID: "worker-app", ProtocolPCID: WorkflowHandoffProtocolPCID, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := runtime.PublishDeliveryPromise(DeliveryPromise{AgentID: "local-router", RecipientAgentID: "worker-app", ProtocolPCID: WorkflowHandoffProtocolPCID, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	runtime.workflowWorker = func(_ context.Context, _ packages.WorkflowAdapter, _ []byte) ([]byte, error) {
		return EncodeWorkflowAdapterResult(WorkflowAdapterResult{
			Output:  WorkflowHandoff{PCID: WorkflowRunProtocolPCID, Values: map[string]string{"subject": "wrong"}},
			Records: []json.RawMessage{json.RawMessage(`{"family":"worker.note.v1","protocol_pcid":"pcid:worker.note.v1","record_id":"worker-1","signer":"worker","timestamp":"2026-08-02T00:00:00Z","payload":{}}`)},
		})
	}
	run, err := runtime.StartWorkflowRun(context.Background(), "worker", WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "egg"}})
	if err != nil {
		t.Fatal(err)
	}
	if run.State != WorkflowRunFailed || run.Reason != "workflow adapter emitted an undeclared output pCID" {
		t.Fatalf("worker run = %#v", run)
	}
	if len(runtime.History()) != 0 {
		t.Fatalf("wrong output applied proposed writes: %#v", runtime.History())
	}
}

func TestActivePackageRegistersExactWorkflowAdapterContract(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	adapter := packages.WorkflowAdapter{Name: "worker", Image: "example/worker:1", InputPCID: WorkflowHandoffProtocolPCID, OutputPCID: WorkflowHandoffProtocolPCID, CPUs: "0.5", Memory: "128m", PIDsLimit: 64, Timeout: "30s"}
	pkg := &activePackage{manifest: packages.Manifest{ID: "procedure-execution-adapter", Version: "0.1.0", WorkflowAdapters: []packages.WorkflowAdapter{adapter}}}
	if err := runtime.activatePackage(pkg); err != nil {
		t.Fatal(err)
	}
	if !runtime.workflowAdapterAvailable(WorkflowManifest{Adapter: "worker", InputPCID: WorkflowHandoffProtocolPCID, OutputPCID: WorkflowHandoffProtocolPCID}) {
		t.Fatal("active package adapter was not available for its exact contract")
	}
	if runtime.workflowAdapterAvailable(WorkflowManifest{Adapter: "worker", InputPCID: WorkflowHandoffProtocolPCID, OutputPCID: WorkflowRunProtocolPCID}) {
		t.Fatal("active package adapter accepted a mismatched output contract")
	}
	if err := runtime.activatePackage(&activePackage{manifest: packages.Manifest{ID: "duplicate-adapter", Version: "0.1.0", WorkflowAdapters: []packages.WorkflowAdapter{adapter}}}); err == nil {
		t.Fatal("duplicate active package workflow adapter was accepted")
	}
}

func TestWorkerProposalPreflightRejectsInvalidRecordWithoutCASWrite(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	result := packages.CommandResult{
		CAS:     []packages.CASWrite{{Alias: "body", Body: "must not persist"}},
		Records: []json.RawMessage{json.RawMessage(`{"invalid":`)},
	}
	if _, err := runtime.applyExternalCommandResult(context.Background(), result); err == nil {
		t.Fatal("invalid worker proposal was accepted")
	}
	if _, err := runtime.cas.Get(store.LegacyObjectID([]byte("must not persist"))); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("invalid worker proposal persisted CAS: %v", err)
	}
}

func TestWorkerProposalRejectsExternallyValidatedFamilyWithoutHostProcess(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	runtime.families["external.note.v1"] = registeredFamily{
		owner:        &activePackage{external: packages.Runner{Executable: filepath.Join(t.TempDir(), "must-not-run")}},
		protocolPCID: "pcid:external.note.v1",
	}
	_, err = runtime.applyWorkflowAdapterResult(context.Background(), packages.CommandResult{Records: []json.RawMessage{json.RawMessage(`{"family":"external.note.v1","protocol_pcid":"pcid:external.note.v1","record_id":"one","signer":"worker","timestamp":"2026-08-02T00:00:00Z","payload":{}}`)}})
	if err == nil || !strings.Contains(err.Error(), "cannot write externally validated family") {
		t.Fatalf("external-family worker proposal error = %v", err)
	}
}

func writeExecutableWorkflow(t *testing.T, root string, runtime *Runtime, name string) {
	t.Helper()
	writeSchemaWorkflow(t, root, runtime, name, WorkflowHandoffProtocolPCID, WorkflowHandoffProtocolPCID)
}

func writeSchemaWorkflow(t *testing.T, root string, runtime *Runtime, name, inputPCID, outputPCID string) {
	t.Helper()
	directory := filepath.Join(root, name)
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"id":"` + name + `","version":"1","summary":"test","required_packages":[],"required_protocols":[],"adapter":"` + name + `","input_pcid":"` + inputPCID + `","output_pcid":"` + outputPCID + `"}`
	if err := os.WriteFile(filepath.Join(directory, "workflow.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.CaptureWorkflowDir(directory, name); err != nil {
		t.Fatal(err)
	}
}
