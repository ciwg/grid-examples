package kernel

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ipfs/go-cid"
)

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
	if replayed.State != WorkflowRunCompleted || replayed.OutputCID == "" {
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
	inputs := []string{"bafkreihjwthblvvsaxlngupwghkshl2lnwgcj5txrr3qpelxtb76stg7ae", "bafkreidndb65kuarxuv3eue6ij3qupblgfuvjmm6v4s5vcfo2d7acbbcwq", "bafkreig4pegj6ckn6hg2yci7i5tt4vcknd4e25ul7a7ugyel5gejjz7iti", "bafkreiboes7s6tcaebjnlibd7fkwj62typezjdsipskyafvtuzf74ypx3i", "bafkreih2mechxf4slowhcag6xac5fqn7wy7tw6pw2lj5nvhc5nthmrx54e", "bafkreibyegjb3p52b3hzf4lw3jwqu3hktxockiamrt2a3bee2gwdl46ja4", "bafkreigkwagey45deeh6cc2hirr53avia3rbt2vgixxcbvwiccctnobi2y"}
	outputs := []string{"bafkreihjhnfom2j2avcjjujcbvy22dbayjkdmkjj6ca3fbfjlm7vm23nxy", "bafkreie7k5xcmmvygwh5fqsbvruh5iivxsyduost7bzrwphhfhudx7ga4q", "bafkreigkvvjof6vhurueeod4mqtghfosiwunlskzkuhjuskiy66cnx5yua", "bafkreicfalhnnj67rctw63c6j4w6x5l7ntmqtvyndmqi26vc46ukavmiha", "bafkreib6qcz4g3lsc4yzfulqihsbczc4wkpo3fwm5f7dvgrznv4qubwppe", "bafkreib3rq3zyljjn4v7tunm2xqpy26i7mpr24bko2sdof524p7xnqnjo4", "bafkreifrf4xznrekx4lohyueahudtz6ju7qhgodpnbyrd4rytjauo7qgsm"}
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
	const sourceInput = "bafkreihjwthblvvsaxlngupwghkshl2lnwgcj5txrr3qpelxtb76stg7ae"
	const sourceOutput = "bafkreihjhnfom2j2avcjjujcbvy22dbayjkdmkjj6ca3fbfjlm7vm23nxy"
	const targetInput = "bafkreidndb65kuarxuv3eue6ij3qupblgfuvjmm6v4s5vcfo2d7acbbcwq"
	const targetOutput = "bafkreie7k5xcmmvygwh5fqsbvruh5iivxsyduost7bzrwphhfhudx7ga4q"
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
	const sourceInput = "bafkreihjwthblvvsaxlngupwghkshl2lnwgcj5txrr3qpelxtb76stg7ae"
	const sourceOutput = "bafkreihjhnfom2j2avcjjujcbvy22dbayjkdmkjj6ca3fbfjlm7vm23nxy"
	const targetInput = "bafkreidndb65kuarxuv3eue6ij3qupblgfuvjmm6v4s5vcfo2d7acbbcwq"
	const targetOutput = "bafkreie7k5xcmmvygwh5fqsbvruh5iivxsyduost7bzrwphhfhudx7ga4q"
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
	const sourceInput = "bafkreihjwthblvvsaxlngupwghkshl2lnwgcj5txrr3qpelxtb76stg7ae"
	const sourceOutput = "bafkreihjhnfom2j2avcjjujcbvy22dbayjkdmkjj6ca3fbfjlm7vm23nxy"
	const targetInput = "bafkreidndb65kuarxuv3eue6ij3qupblgfuvjmm6v4s5vcfo2d7acbbcwq"
	const targetOutput = "bafkreie7k5xcmmvygwh5fqsbvruh5iivxsyduost7bzrwphhfhudx7ga4q"
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
	const sourceInput = "bafkreihjwthblvvsaxlngupwghkshl2lnwgcj5txrr3qpelxtb76stg7ae"
	const sourceOutput = "bafkreihjhnfom2j2avcjjujcbvy22dbayjkdmkjj6ca3fbfjlm7vm23nxy"
	const targetInput = "bafkreidndb65kuarxuv3eue6ij3qupblgfuvjmm6v4s5vcfo2d7acbbcwq"
	const targetOutput = "bafkreie7k5xcmmvygwh5fqsbvruh5iivxsyduost7bzrwphhfhudx7ga4q"
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
