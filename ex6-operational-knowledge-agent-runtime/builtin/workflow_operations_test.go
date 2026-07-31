package builtin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWorkflowOperationSpecificationsMatchEmbeddedSchemas(t *testing.T) {
	for workflow, specification := range workflowOperationSpecifications {
		manifestRaw, err := os.ReadFile(filepath.Join("..", "workflows", workflow, "workflow.json"))
		if err != nil {
			t.Fatal(err)
		}
		var manifest struct {
			InputSchema string `json:"input_schema"`
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
		if specification.outputPCID != manifest.OutputPCID {
			t.Fatalf("%s adapter output pCID = %s, manifest output pCID = %s", workflow, specification.outputPCID, manifest.OutputPCID)
		}
	}
}
