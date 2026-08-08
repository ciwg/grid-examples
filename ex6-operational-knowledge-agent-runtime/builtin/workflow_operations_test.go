package builtin

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages"
)

func TestWorkflowOperationSpecificationsMatchEmbeddedSchemas(t *testing.T) {
	for workflow, specification := range workflowOperationSpecifications {
		manifestRaw, err := os.ReadFile(filepath.Join("..", "workflows", workflow, "workflow.json"))
		if err != nil {
			t.Fatal(err)
		}
		var manifest struct {
			InputSchema string `json:"input_schema"`
			InputPCID   string `json:"input_pcid"`
			OutputPCID  string `json:"output_pcid"`
		}
		if err := json.Unmarshal(manifestRaw, &manifest); err != nil {
			t.Fatal(err)
		}
		inputRaw, err := os.ReadFile(filepath.Join("..", "workflows", workflow, manifest.InputSchema))
		if err != nil {
			t.Fatal(err)
		}
		var input struct {
			Fields []string `json:"fields"`
		}
		if err := json.Unmarshal(inputRaw, &input); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(input.Fields, specification.fields) {
			t.Fatalf("%s adapter fields = %#v, schema fields = %#v", workflow, specification.fields, input.Fields)
		}
		if specification.inputPCID != manifest.InputPCID {
			t.Fatalf("%s adapter input pCID = %s, manifest input pCID = %s", workflow, specification.inputPCID, manifest.InputPCID)
		}
		if specification.outputPCID != manifest.OutputPCID {
			t.Fatalf("%s adapter output pCID = %s, manifest output pCID = %s", workflow, specification.outputPCID, manifest.OutputPCID)
		}
	}
}

func TestWorkflowOperationPreservesLegacyPCIDPair(t *testing.T) {
	specification := workflowOperationSpecifications["inventory-receipt"]
	runtime, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := runtime.RegisterBuiltin(kernel.BuiltinPackage{
		Manifest: packages.Manifest{ID: "inventory", Version: "test", Commands: []packages.Command{{Path: specification.command, Summary: "test"}}},
		Commands: map[string]kernel.BuiltinCommand{specification.command[0] + " " + specification.command[1]: func(_ context.Context, _ *kernel.Runtime, _ []string) (string, error) { return "ok", nil }},
	}); err != nil {
		t.Fatal(err)
	}
	operation := commandWorkflowOperation(specification)
	if err := runtime.RegisterWorkflowOperation("inventory-receipt", operation); err != nil {
		t.Fatal(err)
	}
	input := kernel.WorkflowHandoff{PCID: specification.legacyInputPCID, Values: map[string]string{"inventory_id": "inventory-1", "run_id": "run-1", "place_id": "dock", "counter": "alice", "quantity": "1", "outcome": "accepted", "notes": "legacy"}}
	output, err := operation(context.Background(), runtime, input)
	if err != nil {
		t.Fatal(err)
	}
	if output.PCID != specification.legacyOutputPCID {
		t.Fatalf("legacy output pCID = %s, want %s", output.PCID, specification.legacyOutputPCID)
	}
	manifest := `{"id":"inventory-receipt","version":"1.0.0","summary":"legacy","required_packages":[],"required_protocols":[],"adapter":"inventory-receipt","input_pcid":"` + specification.legacyInputPCID + `","output_pcid":"` + specification.legacyOutputPCID + `"}`
	var artifact bytes.Buffer
	archive := tar.NewWriter(&artifact)
	if err := archive.WriteHeader(&tar.Header{Name: "workflow.json", Mode: 0o644, Size: int64(len(manifest))}); err != nil {
		t.Fatal(err)
	}
	if _, err := archive.Write([]byte(manifest)); err != nil {
		t.Fatal(err)
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	artifactCID, err := runtime.PutCAS(artifact.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.ImportWorkflow(kernel.Workflow{ID: "legacy-inventory-receipt", ArtifactCID: artifactCID}); err != nil {
		t.Fatal(err)
	}
	if err := runtime.ActivateWorkflow("legacy-inventory-receipt"); err != nil {
		t.Fatal(err)
	}
	verification, err := runtime.VerifyWorkflowReadiness("legacy-inventory-receipt")
	if err != nil {
		t.Fatal(err)
	}
	if verification.Contract != "retained-v1" || !verification.AdapterAvailable || verification.SchemaCASReady || !verification.EligibleToExecute {
		t.Fatalf("legacy verification = %#v", verification)
	}
	run, err := runtime.StartWorkflowRun(context.Background(), "legacy-inventory-receipt", input)
	if err != nil {
		t.Fatal(err)
	}
	if run.State != kernel.WorkflowRunCompleted {
		t.Fatalf("legacy artifact state = %s", run.State)
	}
}

func TestWorkflowOperationSpecificationsPreserveAllLegacyPCIDPairs(t *testing.T) {
	for workflow, specification := range workflowOperationSpecifications {
		if specification.legacyInputPCID == "" {
			continue
		}
		outputPCID, ok := specification.outputForInput(specification.legacyInputPCID)
		if !ok {
			t.Fatalf("%s does not accept its legacy input pCID", workflow)
		}
		if outputPCID != specification.legacyOutputPCID {
			t.Fatalf("%s legacy output pCID = %s, want %s", workflow, outputPCID, specification.legacyOutputPCID)
		}
	}
}
