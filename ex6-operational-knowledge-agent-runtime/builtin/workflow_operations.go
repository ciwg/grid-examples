package builtin

import (
	"context"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
)

// WorkflowOperations supplies the trusted in-process adapter for every shipped
// artifact. Intent: Keep v1 orchestration explicit and pCID-typed while Docker
// workers remain a separately constrained future backend. Source: DI-lumek
func WorkflowOperations() map[string]kernel.WorkflowOperation {
	operations := map[string]kernel.WorkflowOperation{}
	for name, specification := range workflowOperationSpecifications {
		operations[name] = commandWorkflowOperation(specification.command, specification.fields, specification.outputPCID)
	}
	return operations
}

type workflowOperationSpecification struct {
	command, fields []string
	outputPCID      string
}

// Intent: Keep trusted adapter requirements mechanically aligned with the
// pCID-defined workflow schema shipped by each artifact. Source: DI-lumek
var workflowOperationSpecifications = map[string]workflowOperationSpecification{
	"inventory-discrepancy-review": {[]string{"inventory", "record-reconcile"}, []string{"inventory_id", "event_id", "decision", "resource_id", "notes"}, "bafkreichcquo564amype7v6locjdhc7xl5kgb6i7oyo25k65o677kyztey"},
	"inventory-receipt":            {[]string{"inventory", "record-count"}, []string{"inventory_id", "run_id", "place_id", "counter", "quantity", "outcome", "notes"}, "bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne"},
	"knowledge-review":             {[]string{"knowledge", "item", "approve"}, []string{"item_id", "event_id", "notes"}, "bafkreigxtnvprgom4b6ftkavmvfeqv2teo45b7s3n57k3wn7ipziwlad4e"},
	"maintenance-round":            {[]string{"maintenance", "record-service"}, []string{"maintenance_id", "run_id", "resource_id", "performer", "outcome", "notes"}, "bafkreifj3qpjinq4vanr5vo22dvazuybg6jrrawnfhw43pooqnzp62vtg4"},
	"procedure-execution":          {[]string{"procedures", "record-use"}, []string{"procedure_id", "run_id", "actor", "outcome", "notes"}, "bafkreiamprv3apzowjzqbkp3hnrhrla5aq7lp5kyzbca5j3iv4v5jmhwa4"},
	"receiving-check":              {[]string{"receiving", "record-receipt"}, []string{"receiving_id", "run_id", "place_id", "receiver", "outcome", "notes"}, "bafkreiddsvg5v7a2dwa4omzicgobed5yqheehbkjw263kb3efvv2yfgnz4"},
	"training-qualification":       {[]string{"training", "record-session"}, []string{"training_id", "run_id", "trainee", "instructor", "outcome", "notes"}, "bafkreiauh6xo45sp3zhmhhjvehiqstp36loiyu26smicojfntl7ek75chy"},
}

func commandWorkflowOperation(command []string, fields []string, outputPCID string) kernel.WorkflowOperation {
	return func(ctx context.Context, runtime *kernel.Runtime, input kernel.WorkflowHandoff) (kernel.WorkflowHandoff, error) {
		args := append([]string{}, command...)
		for _, field := range fields {
			value := strings.TrimSpace(input.Values[field])
			if value == "" {
				return kernel.WorkflowHandoff{}, kernel.WaitingForWorkflowInput("workflow input requires " + field)
			}
			args = append(args, value)
		}
		if _, err := runtime.RunCommand(ctx, args); err != nil {
			return kernel.WorkflowHandoff{}, err
		}
		output := kernel.WorkflowHandoff{PCID: outputPCID, Values: map[string]string{}}
		for key, value := range input.Values {
			output.Values[key] = value
		}
		output.Values["stage"] = "completed"
		return output, nil
	}
}
