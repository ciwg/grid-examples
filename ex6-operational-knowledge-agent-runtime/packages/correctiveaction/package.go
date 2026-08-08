package correctiveaction

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

const (
	PackageID     = "correctiveaction"
	EventFamily   = "moks.correctiveaction.event.v1"
	EventProtocol = "pcid:moks.correctiveaction.event.v1"
)

type event struct {
	ActionID         string `json:"action_id"`
	QuarantineCaseID string `json:"quarantine_case_id"`
	Actor            string `json:"actor"`
	EvidenceID       string `json:"evidence_id"`
	Summary          string `json:"summary"`
	Notes            string `json:"notes,omitempty"`
}

func Package() kernel.BuiltinPackage {
	return kernel.BuiltinPackage{Manifest: pkgmeta.Manifest{ID: PackageID, Version: "0.1.0", Description: "First-party corrective-action package", Commands: []pkgmeta.Command{{Path: []string{"correctiveaction", "open"}, Summary: "Open a corrective action"}, {Path: []string{"correctiveaction", "list"}, Summary: "List corrective actions"}, {Path: []string{"correctiveaction", "inspect"}, Summary: "Inspect a corrective action"}}, Families: []pkgmeta.Family{{Name: EventFamily, ProtocolPCID: EventProtocol}}, Claims: []pkgmeta.ImplementationClaim{{ProtocolPCID: EventProtocol, Role: "family-validator", Summary: "Validates corrective-action events."}, {ProtocolPCID: EventProtocol, Role: "domain-behavior", Summary: "Declares corrective-action openings."}}}, Commands: map[string]kernel.BuiltinCommand{"correctiveaction open": openAction, "correctiveaction list": listActions, "correctiveaction inspect": inspectAction}, Validators: map[string]kernel.BuiltinValidator{EventFamily: validateEvent}}
}

// Intent: Open accountable follow-up only from a rejected quarantine case. Source: DI-hiboj
func openAction(ctx stdctx.Context, r *kernel.Runtime, a []string) (string, error) {
	if len(a) < 5 {
		return "", errors.New("usage: correctiveaction open <action-id> <quarantine-case-id> <actor> <evidence-id> <summary> [notes...]")
	}
	if _, e := r.RunCommand(ctx, []string{"quarantine", "inspect", a[1]}); e != nil {
		return "", e
	}
	state, _ := r.RunCommand(ctx, []string{"quarantine", "inspect", a[1]})
	if !strings.Contains(state, "transition: reject") {
		return "", errors.New("corrective action requires a rejected quarantine case")
	}
	e := event{a[0], a[1], a[2], a[3], a[4], strings.Join(a[5:], " ")}
	env, err := kernel.NewEnvelope(EventFamily, EventProtocol, a[0], PackageID, e)
	if err != nil {
		return "", err
	}
	if _, err = r.AppendRecord(ctx, records.MustMarshal(env)); err != nil {
		return "", err
	}
	return "stored corrective action " + a[0], nil
}
func actions(r *kernel.Runtime) (map[string]event, error) {
	out := map[string]event{}
	for _, h := range r.History() {
		if h.Envelope.Family != EventFamily {
			continue
		}
		var e event
		if err := json.Unmarshal(h.Envelope.Payload, &e); err != nil {
			return nil, err
		}
		out[e.ActionID] = e
	}
	return out, nil
}
func listActions(_ stdctx.Context, r *kernel.Runtime, _ []string) (string, error) {
	s, e := actions(r)
	if e != nil {
		return "", e
	}
	ids := make([]string, 0, len(s))
	for id := range s {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		v := s[id]
		out = append(out, fmt.Sprintf("%s\t%s\t%s", v.ActionID, v.QuarantineCaseID, v.Summary))
	}
	return strings.Join(out, "\n"), nil
}
func inspectAction(_ stdctx.Context, r *kernel.Runtime, a []string) (string, error) {
	if len(a) != 1 {
		return "", errors.New("usage: correctiveaction inspect <action-id>")
	}
	s, e := actions(r)
	if e != nil {
		return "", e
	}
	v, ok := s[a[0]]
	if !ok {
		return "", fmt.Errorf("unknown corrective action: %s", a[0])
	}
	return fmt.Sprintf("action_id: %s\nquarantine_case_id: %s\nsummary: %s", v.ActionID, v.QuarantineCaseID, v.Summary), nil
}
func validateEvent(env records.Envelope) error {
	var e event
	if err := json.Unmarshal(env.Payload, &e); err != nil {
		return err
	}
	if e.ActionID == "" || e.QuarantineCaseID == "" || e.Actor == "" || e.EvidenceID == "" || e.Summary == "" || env.RecordID != e.ActionID {
		return errors.New("corrective action event requires action_id, quarantine_case_id, actor, evidence_id, summary, and matching record ID")
	}
	return nil
}
