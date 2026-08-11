package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	procedurespkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/procedures"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/runs"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
)

const (
	procedureExecutionInputPCID  = "bafkreiawxq2i7q57tks6f5viofxkko2jf2txmlurbp3i33svynytyjswfq"
	procedureExecutionOutputPCID = "bafkreiamprv3apzowjzqbkp3hnrhrla5aq7lp5kyzbca5j3iv4v5jmhwa4"
)

func main() {
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail(err)
	}
	input, err := kernel.DecodeWorkflowHandoff(inputBytes)
	if err != nil {
		fail(err)
	}
	if input.PCID != procedureExecutionInputPCID {
		fail(fmt.Errorf("unsupported input pCID %s", input.PCID))
	}
	for _, field := range []string{"procedure_id", "run_id", "actor", "outcome", "notes"} {
		if strings.TrimSpace(input.Values[field]) == "" {
			fail(fmt.Errorf("procedure execution requires %s", field))
		}
	}
	result, err := procedureExecutionResult(input)
	if err != nil {
		fail(err)
	}
	if _, err := os.Stdout.Write(result); err != nil {
		fail(err)
	}
}

// Intent: Reproduce procedure execution's durable run and procedure-use
// effects as runtime-validated proposals without granting the worker state
// access. Source: DI-fofuh
func procedureExecutionResult(input kernel.WorkflowHandoff) ([]byte, error) {
	outputValues := make(map[string]string, len(input.Values)+1)
	for key, value := range input.Values {
		outputValues[key] = value
	}
	outputValues["stage"] = "completed"
	timestamp := time.Now().UTC().Format(time.RFC3339)
	runPayload, err := json.Marshal(map[string]string{
		"item_id": input.Values["procedure_id"],
		"actor":   input.Values["actor"],
		"outcome": input.Values["outcome"],
		"notes":   input.Values["notes"],
	})
	if err != nil {
		return nil, err
	}
	usePayload, err := json.Marshal(map[string]string{
		"procedure_id": input.Values["procedure_id"],
		"run_id":       input.Values["run_id"],
	})
	if err != nil {
		return nil, err
	}
	return kernel.EncodeWorkflowAdapterResult(kernel.WorkflowAdapterResult{
		Output: kernel.WorkflowHandoff{PCID: procedureExecutionOutputPCID, Values: outputValues},
		Records: [][]byte{
			records.MustMarshal(records.Envelope{Family: runspkg.RunFamily, ProtocolPCID: runspkg.RunProtocol, RecordID: input.Values["run_id"], Signer: "procedure-execution-adapter", Timestamp: timestamp, Payload: runPayload}),
			records.MustMarshal(records.Envelope{Family: procedurespkg.UseFamily, ProtocolPCID: procedurespkg.UseProtocol, RecordID: input.Values["run_id"], Signer: "procedure-execution-adapter", Timestamp: timestamp, Payload: usePayload}),
		},
	})
}

func fail(err error) {
	if _, writeErr := fmt.Fprintln(os.Stderr, err); writeErr != nil {
		os.Exit(1)
	}
	os.Exit(1)
}
