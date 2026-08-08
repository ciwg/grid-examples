package builtin

import (
	"context"
	"errors"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
)

// WorkflowOperations supplies the trusted in-process adapter for every shipped
// artifact. Intent: Keep v1 orchestration explicit and pCID-typed while Docker
// workers remain a separately constrained future backend. Source: DI-lumek
func WorkflowOperations() map[string]kernel.WorkflowOperation {
	operations := map[string]kernel.WorkflowOperation{}
	for name, specification := range workflowOperationSpecifications {
		operations[name] = commandWorkflowOperation(specification)
	}
	return operations
}

type workflowOperationSpecification struct {
	command, fields                   []string
	inputPCID, outputPCID             string
	legacyInputPCID, legacyOutputPCID string
	arguments                         []string
	outputFields                      map[string]string
}

// Intent: Keep trusted adapter requirements mechanically aligned with the
// pCID-defined workflow schema shipped by each artifact, including the
// receiving-exception contract's explicit opening-event output. Source:
// DI-lumek; DI-hogid
var workflowOperationSpecifications = map[string]workflowOperationSpecification{
	"inventory-discrepancy-review": {[]string{"inventory", "record-reconcile"}, []string{"inventory_id", "event_id", "decision", "resource_id", "notes"}, "bafkreigwhgyyvdxkjrckvimjhesxg2ms2wtahqcogu7276xqovcjoxif3e", "bafkreichcquo564amype7v6locjdhc7xl5kgb6i7oyo25k65o677kyztey", "bafkreihjwthblvvsaxlngupwghkshl2lnwgcj5txrr3qpelxtb76stg7ae", "bafkreihjhnfom2j2avcjjujcbvy22dbayjkdmkjj6ca3fbfjlm7vm23nxy", nil, nil},
	"inventory-receipt":            {[]string{"inventory", "record-count"}, []string{"inventory_id", "run_id", "place_id", "counter", "quantity", "outcome", "notes"}, "bafkreie3xn5cs7in24a5aenl7kpyaa22e346wr4tcqm4evxgcn2v55yvne", "bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne", "bafkreidndb65kuarxuv3eue6ij3qupblgfuvjmm6v4s5vcfo2d7acbbcwq", "bafkreie7k5xcmmvygwh5fqsbvruh5iivxsyduost7bzrwphhfhudx7ga4q", nil, nil},
	"knowledge-review":             {[]string{"knowledge", "item", "approve"}, []string{"item_id", "event_id", "notes"}, "bafkreigurwnuwdbri3ntiuzseyoft7tazq74rrxridxbpekqcqxd2rzvhi", "bafkreigxtnvprgom4b6ftkavmvfeqv2teo45b7s3n57k3wn7ipziwlad4e", "bafkreig4pegj6ckn6hg2yci7i5tt4vcknd4e25ul7a7ugyel5gejjz7iti", "bafkreigkvvjof6vhurueeod4mqtghfosiwunlskzkuhjuskiy66cnx5yua", nil, nil},
	"maintenance-round":            {[]string{"maintenance", "record-service"}, []string{"maintenance_id", "run_id", "resource_id", "performer", "outcome", "notes"}, "bafkreicuxiyha56khoiwfktrh3pqrq6v47uovjm52kmuskmkbeml6qyxsu", "bafkreifj3qpjinq4vanr5vo22dvazuybg6jrrawnfhw43pooqnzp62vtg4", "bafkreiboes7s6tcaebjnlibd7fkwj62typezjdsipskyafvtuzf74ypx3i", "bafkreicfalhnnj67rctw63c6j4w6x5l7ntmqtvyndmqi26vc46ukavmiha", nil, nil},
	"procedure-execution":          {[]string{"procedures", "record-use"}, []string{"procedure_id", "run_id", "actor", "outcome", "notes"}, "bafkreiawxq2i7q57tks6f5viofxkko2jf2txmlurbp3i33svynytyjswfq", "bafkreiamprv3apzowjzqbkp3hnrhrla5aq7lp5kyzbca5j3iv4v5jmhwa4", "bafkreih2mechxf4slowhcag6xac5fqn7wy7tw6pw2lj5nvhc5nthmrx54e", "bafkreib6qcz4g3lsc4yzfulqihsbczc4wkpo3fwm5f7dvgrznv4qubwppe", nil, nil},
	"receiving-check":              {[]string{"receiving", "record-receipt"}, []string{"receiving_id", "run_id", "place_id", "receiver", "outcome", "notes"}, "bafkreicubq2eqovbcdvlj3ggpmrnwlzohp5no5yuwehk4eyucicia4kehy", "bafkreiddsvg5v7a2dwa4omzicgobed5yqheehbkjw263kb3efvv2yfgnz4", "bafkreibyegjb3p52b3hzf4lw3jwqu3hktxockiamrt2a3bee2gwdl46ja4", "bafkreib3rq3zyljjn4v7tunm2xqpy26i7mpr24bko2sdof524p7xnqnjo4", nil, nil},
	"receiving-exception":          {[]string{"quarantine", "open"}, []string{"receiving_id", "receipt_run_id", "case_id", "actor", "evidence_id", "exception", "notes"}, "bafkreifpihgcjflvxrz3kevjjnfdcnn4nzzdc7rbhvch6exs2cfz3xxc7i", "bafkreid4334xddsgvuvhmlw75dsj7o5pjn2q4mrebqqkasc2coqfijupn4", "", "", []string{"case_id", "receiving_id", "receipt_run_id", "actor", "evidence_id", "exception", "notes"}, map[string]string{"opening_event_id": "case_id"}},
	"training-qualification":       {[]string{"training", "record-session"}, []string{"training_id", "run_id", "trainee", "instructor", "outcome", "notes"}, "bafkreigdqszsce7qvcihohhsgr4r5wh7l3lgbqm3gd55befrf7d4hxpwyi", "bafkreiauh6xo45sp3zhmhhjvehiqstp36loiyu26smicojfntl7ek75chy", "bafkreigkwagey45deeh6cc2hirr53avia3rbt2vgixxcbvwiccctnobi2y", "bafkreifrf4xznrekx4lohyueahudtz6ju7qhgodpnbyrd4rytjauo7qgsm", nil, nil},
}

// outputForInput preserves the exact v1 output contract selected by the
// artifact's input pCID. Intent: Retained artifacts must remain executable
// after corrected schema pCIDs are introduced for newly captured artifacts.
// Source: DI-lumek
func (specification workflowOperationSpecification) outputForInput(inputPCID string) (string, bool) {
	switch inputPCID {
	case specification.inputPCID:
		return specification.outputPCID, true
	case specification.legacyInputPCID:
		if specification.legacyInputPCID == "" {
			return "", false
		}
		return specification.legacyOutputPCID, true
	default:
		return "", false
	}
}

func commandWorkflowOperation(specification workflowOperationSpecification) kernel.WorkflowOperation {
	return func(ctx context.Context, runtime *kernel.Runtime, input kernel.WorkflowHandoff) (kernel.WorkflowHandoff, error) {
		outputPCID, ok := specification.outputForInput(input.PCID)
		if !ok {
			return kernel.WorkflowHandoff{}, errors.New("workflow adapter input pCID is not supported")
		}
		args := append([]string{}, specification.command...)
		argumentFields := specification.fields
		if len(specification.arguments) > 0 {
			argumentFields = specification.arguments
		}
		for _, field := range argumentFields {
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
		for outputField, inputField := range specification.outputFields {
			// Intent: Return the durable opening-event identity without requiring a
			// hidden adapter-generated ID. Source: DI-hogid
			output.Values[outputField] = input.Values[inputField]
		}
		output.Values["stage"] = "completed"
		return output, nil
	}
}
