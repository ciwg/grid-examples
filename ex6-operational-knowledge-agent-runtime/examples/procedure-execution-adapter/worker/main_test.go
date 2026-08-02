package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages"
)

const procedureExecutionAdapterImage = "sha256:e02dfedf0daa5770bb785d11b4c4c8f51e377ad8144d73bf0846d5e64fb9410d"

func TestProcedureExecutionResultProposesDurableProcedureUse(t *testing.T) {
	resultBytes, err := procedureExecutionResult(kernel.WorkflowHandoff{
		PCID: procedureExecutionInputPCID,
		Values: map[string]string{
			"procedure_id": "proc-1",
			"run_id":       "run-1",
			"actor":        "alice",
			"outcome":      "completed",
			"notes":        "followed procedure",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := kernel.DecodeWorkflowAdapterResult(resultBytes)
	if err != nil {
		t.Fatal(err)
	}
	if result.Output.PCID != procedureExecutionOutputPCID || result.Output.Values["stage"] != "completed" {
		t.Fatalf("unexpected output: %#v", result.Output)
	}
	if len(result.Records) != 2 {
		t.Fatalf("record proposals = %d, want 2", len(result.Records))
	}
}

func TestDockerProcedureExecutionAdapter(t *testing.T) {
	if os.Getenv("MOKS_DOCKER_INTEGRATION") != "1" {
		t.Skip("set MOKS_DOCKER_INTEGRATION=1 after building the pinned procedure-execution adapter image")
	}
	input, err := kernel.EncodeWorkflowHandoff(kernel.WorkflowHandoff{
		PCID: procedureExecutionInputPCID,
		Values: map[string]string{
			"procedure_id": "proc-1",
			"run_id":       "run-1",
			"actor":        "alice",
			"outcome":      "completed",
			"notes":        "followed procedure",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	output, err := (packages.DockerWorker{Image: procedureExecutionAdapterImage, CPUs: "0.5", Memory: "128m", PIDsLimit: 64, Timeout: 30 * time.Second}).Run(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	result, err := kernel.DecodeWorkflowAdapterResult(output)
	if err != nil {
		t.Fatal(err)
	}
	if result.Output.PCID != procedureExecutionOutputPCID || result.Output.Values["stage"] != "completed" || len(result.Records) != 2 {
		t.Fatalf("unexpected Docker result: %#v", result)
	}
}
