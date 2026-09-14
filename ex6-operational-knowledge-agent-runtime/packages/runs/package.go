package runspkg

import (
	stdctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	pkgmeta "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
)

var (
	PackageID = "runs"

	RunFamily   = "moks.runs.run.v1"
	RunProtocol = records.PackageProtocolPCID(RunFamily)

	EvidenceFamily   = "moks.runs.evidence.v1"
	EvidenceProtocol = records.PackageProtocolPCID(EvidenceFamily)

	ApprovalFamily   = "moks.runs.approval.v1"
	ApprovalProtocol = records.PackageProtocolPCID(ApprovalFamily)
)

type runPayload struct {
	ItemID  string `json:"item_id"`
	Actor   string `json:"actor"`
	Outcome string `json:"outcome"`
	Notes   string `json:"notes"`
}

type evidencePayload struct {
	RunID   string            `json:"run_id"`
	Summary string            `json:"summary"`
	Facts   map[string]string `json:"facts,omitempty"`
	BodyRef string            `json:"body_ref,omitempty"`
}

type approvalPayload struct {
	RunID    string `json:"run_id"`
	Decision string `json:"decision"`
	Notes    string `json:"notes,omitempty"`
}

type runState struct {
	RunID     string
	ItemID    string
	Actor     string
	Outcome   string
	Notes     string
	Evidence  []string
	Approvals []string
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party runs package",
		Commands: []pkgmeta.Command{
			{Path: []string{"runs", "record"}, Summary: "Record performed work"},
			{Path: []string{"runs", "inspect"}, Summary: "Inspect a run"},
			{Path: []string{"runs", "evidence", "add"}, Summary: "Add evidence to a run"},
			{Path: []string{"runs", "approve"}, Summary: "Add a run approval"},
		},
		Families: []pkgmeta.Family{
			{Name: RunFamily, ProtocolPCID: RunProtocol},
			{Name: EvidenceFamily, ProtocolPCID: EvidenceProtocol},
			{Name: ApprovalFamily, ProtocolPCID: ApprovalProtocol},
		},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: RunProtocol, Role: "family-validator", Summary: "Validates run records."},
			{ProtocolPCID: EvidenceProtocol, Role: "family-validator", Summary: "Validates run evidence records."},
			{ProtocolPCID: ApprovalProtocol, Role: "family-validator", Summary: "Validates run approval records."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"runs record":       recordRun,
			"runs inspect":      inspectRun,
			"runs evidence add": addEvidence,
			"runs approve":      approveRun,
		},
		Validators: map[string]kernel.BuiltinValidator{
			RunFamily:      validateRunEnvelope,
			EvidenceFamily: validateEvidenceEnvelope,
			ApprovalFamily: validateApprovalEnvelope,
		},
	}
}

func recordRun(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("usage: runs record <run-id> <item-id> <actor> <outcome> [notes...]")
	}
	notes := ""
	if len(args) > 4 {
		notes = strings.Join(args[4:], " ")
	}
	envelope, err := kernel.NewEnvelope(RunFamily, RunProtocol, args[0], PackageID, runPayload{
		ItemID:  args[1],
		Actor:   args[2],
		Outcome: args[3],
		Notes:   notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored run %s", args[0]), nil
}

func addEvidence(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: runs evidence add <run-id> <evidence-id> <summary> [facts] [body...]")
	}
	facts := map[string]string{}
	if len(args) > 3 && strings.TrimSpace(args[3]) != "" {
		facts = parseFacts(args[3])
	}
	bodyRef := ""
	if len(args) > 4 {
		var err error
		bodyRef, err = runtime.PutCAS([]byte(strings.Join(args[4:], " ")))
		if err != nil {
			return "", err
		}
	}
	envelope, err := kernel.NewEnvelope(EvidenceFamily, EvidenceProtocol, args[1], PackageID, evidencePayload{
		RunID:   args[0],
		Summary: args[2],
		Facts:   facts,
		BodyRef: bodyRef,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored evidence %s for %s", args[1], args[0]), nil
}

func approveRun(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: runs approve <run-id> <approval-id> <decision> [notes...]")
	}
	notes := ""
	if len(args) > 3 {
		notes = strings.Join(args[3:], " ")
	}
	envelope, err := kernel.NewEnvelope(ApprovalFamily, ApprovalProtocol, args[1], PackageID, approvalPayload{
		RunID:    args[0],
		Decision: args[2],
		Notes:    notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored approval %s for %s", args[1], args[0]), nil
}

func inspectRun(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: runs inspect <run-id>")
	}
	state, err := loadRunState(runtime)
	if err != nil {
		return "", err
	}
	run, ok := state[args[0]]
	if !ok {
		return "", fmt.Errorf("unknown run: %s", args[0])
	}
	return fmt.Sprintf(
		"id: %s\nitem_id: %s\nactor: %s\noutcome: %s\nnotes: %s\nevidence: %s\napprovals: %s",
		run.RunID,
		run.ItemID,
		run.Actor,
		run.Outcome,
		run.Notes,
		strings.Join(run.Evidence, ", "),
		strings.Join(run.Approvals, ", "),
	), nil
}

func loadRunState(runtime *kernel.Runtime) (map[string]*runState, error) {
	runs := map[string]*runState{}
	for _, entry := range runtime.History() {
		switch entry.Envelope.Family {
		case RunFamily:
			var body runPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			runs[entry.Envelope.RecordID] = &runState{
				RunID:   entry.Envelope.RecordID,
				ItemID:  body.ItemID,
				Actor:   body.Actor,
				Outcome: body.Outcome,
				Notes:   body.Notes,
			}
		case EvidenceFamily:
			var body evidencePayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			run := ensureRun(runs, body.RunID)
			run.Evidence = append(run.Evidence, fmt.Sprintf("%s:%s", entry.Envelope.RecordID, body.Summary))
		case ApprovalFamily:
			var body approvalPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			run := ensureRun(runs, body.RunID)
			run.Approvals = append(run.Approvals, fmt.Sprintf("%s:%s", entry.Envelope.RecordID, body.Decision))
		}
	}
	for _, run := range runs {
		sort.Strings(run.Evidence)
		sort.Strings(run.Approvals)
	}
	return runs, nil
}

func ensureRun(runs map[string]*runState, runID string) *runState {
	run := runs[runID]
	if run == nil {
		run = &runState{RunID: runID}
		runs[runID] = run
	}
	return run
}

func parseFacts(raw string) map[string]string {
	facts := map[string]string{}
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		key, value, found := strings.Cut(entry, "=")
		if !found {
			facts[entry] = ""
			continue
		}
		facts[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return facts
}

// Intent: Keep run-family payload validation in the owning package so the
// runtime stays generic while the package owns performed-work semantics. Source:
// DI-pamuk
func validateRunEnvelope(envelope records.Envelope) error {
	var body runPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.ItemID) == "" {
		return errors.New("item_id is required")
	}
	if strings.TrimSpace(body.Actor) == "" {
		return errors.New("actor is required")
	}
	if strings.TrimSpace(body.Outcome) == "" {
		return errors.New("outcome is required")
	}
	return nil
}

// Intent: Keep run-family payload validation in the owning package so the
// runtime stays generic while the package owns performed-work semantics. Source:
// DI-pamuk
func validateEvidenceEnvelope(envelope records.Envelope) error {
	var body evidencePayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.RunID) == "" {
		return errors.New("run_id is required")
	}
	if strings.TrimSpace(body.Summary) == "" {
		return errors.New("summary is required")
	}
	return nil
}

// Intent: Keep run-family payload validation in the owning package so the
// runtime stays generic while the package owns performed-work semantics. Source:
// DI-pamuk
func validateApprovalEnvelope(envelope records.Envelope) error {
	var body approvalPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.RunID) == "" {
		return errors.New("run_id is required")
	}
	if strings.TrimSpace(body.Decision) == "" {
		return errors.New("decision is required")
	}
	return nil
}
