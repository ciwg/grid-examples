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
	return map[string]kernel.WorkflowOperation{
		"inventory-discrepancy-review": commandWorkflowOperation([]string{"inventory", "record-reconcile"}, []string{"inventory_id", "event_id", "decision", "resource_id", "notes"}, "bafkreihjhnfom2j2avcjjujcbvy22dbayjkdmkjj6ca3fbfjlm7vm23nxy"),
		"inventory-receipt":            commandWorkflowOperation([]string{"inventory", "record-count"}, []string{"inventory_id", "run_id", "place_id", "counter", "quantity", "outcome", "notes"}, "bafkreie7k5xcmmvygwh5fqsbvruh5iivxsyduost7bzrwphhfhudx7ga4q"),
		"knowledge-review":             commandWorkflowOperation([]string{"knowledge", "item", "approve"}, []string{"item_id", "event_id", "notes"}, "bafkreigkvvjof6vhurueeod4mqtghfosiwunlskzkuhjuskiy66cnx5yua"),
		"maintenance-round":            commandWorkflowOperation([]string{"maintenance", "record-service"}, []string{"maintenance_id", "run_id", "resource_id", "performer", "outcome", "notes"}, "bafkreicfalhnnj67rctw63c6j4w6x5l7ntmqtvyndmqi26vc46ukavmiha"),
		"procedure-execution":          commandWorkflowOperation([]string{"procedures", "record-use"}, []string{"procedure_id", "run_id", "actor", "outcome", "notes"}, "bafkreib6qcz4g3lsc4yzfulqihsbczc4wkpo3fwm5f7dvgrznv4qubwppe"),
		"receiving-check":              commandWorkflowOperation([]string{"receiving", "record-receipt"}, []string{"receiving_id", "run_id", "place_id", "receiver", "outcome", "notes"}, "bafkreib3rq3zyljjn4v7tunm2xqpy26i7mpr24bko2sdof524p7xnqnjo4"),
		"training-qualification":       commandWorkflowOperation([]string{"training", "record-session"}, []string{"training_id", "run_id", "trainee", "instructor", "outcome", "notes"}, "bafkreifrf4xznrekx4lohyueahudtz6ju7qhgodpnbyrd4rytjauo7qgsm"),
	}
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
