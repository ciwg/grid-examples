package kernel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
)

func TestShippedWorkflowSchemaPCIDsMatchCanonicalSpecifications(t *testing.T) {
	root := filepath.Join("..", "workflows")
	specifications := filepath.Join("..", "docs", "protocols", "workflow-adapter-schemas")
	cas, err := store.OpenCAS(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, workflow := range []string{"inventory-discrepancy-review", "inventory-receipt", "knowledge-review", "maintenance-round", "procedure-execution", "receiving-check", "training-qualification"} {
		manifestRaw, err := os.ReadFile(filepath.Join(root, workflow, "workflow.json"))
		if err != nil {
			t.Fatal(err)
		}
		var manifest WorkflowManifest
		if err := json.Unmarshal(manifestRaw, &manifest); err != nil {
			t.Fatal(err)
		}
		if err := manifest.Validate(); err != nil {
			t.Fatal(err)
		}
		for _, schema := range []struct{ suffix, pcid string }{{"input", manifest.InputPCID}, {"output", manifest.OutputPCID}} {
			raw, err := os.ReadFile(filepath.Join(specifications, workflow+"-"+schema.suffix+"-v1.json"))
			if err != nil {
				t.Fatal(err)
			}
			actual, err := cas.PutCID(raw)
			if err != nil {
				t.Fatal(err)
			}
			if actual.String() != schema.pcid {
				t.Fatalf("%s %s pCID = %s, want %s", workflow, schema.suffix, actual, schema.pcid)
			}
		}
	}
}

func TestWorkflowManifestRejectsPartialExecutableDeclaration(t *testing.T) {
	manifest := WorkflowManifest{ID: "test", Version: "1", Summary: "test", Adapter: "test", InputPCID: WorkflowHandoffProtocolPCID}
	if err := manifest.Validate(); err == nil {
		t.Fatal("partial executable declaration accepted")
	}
}

func TestCaptureWorkflowDirRejectsRetiredAdapterPCIDs(t *testing.T) {
	runtime, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	for pcid := range legacyWorkflowAdapterPCIDs {
		directory := t.TempDir()
		manifest := `{"id":"retired","version":"1","summary":"retired","required_packages":[],"required_protocols":[],"adapter":"retired","input_pcid":"` + pcid + `","output_pcid":"` + WorkflowHandoffProtocolPCID + `"}`
		if err := os.WriteFile(filepath.Join(directory, "workflow.json"), []byte(manifest), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := runtime.CaptureWorkflowDir(directory, "retired"); err == nil {
			t.Fatalf("capture accepted retired adapter pCID %s", pcid)
		}
	}
	ids, err := runtime.cas.ListCIDs()
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 0 {
		t.Fatalf("retired capture persisted %d CAS objects", len(ids))
	}
}
