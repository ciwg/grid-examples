package quarantinepkg

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
	PackageID = "quarantine"

	EventFamily   = "moks.quarantine.event.v1"
	EventProtocol = records.PackageProtocolPCID(EventFamily)
)

type eventPayload struct {
	CaseID       string `json:"case_id"`
	Transition   string `json:"transition"`
	ReceivingID  string `json:"receiving_id,omitempty"`
	ReceiptRunID string `json:"receipt_run_id,omitempty"`
	Actor        string `json:"actor"`
	EvidenceID   string `json:"evidence_id"`
	Exception    string `json:"exception,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type caseState struct {
	CaseID      string
	ReceivingID string
	Transition  string
	Events      []string
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party quarantine package",
		Commands: []pkgmeta.Command{
			{Path: []string{"quarantine", "open"}, Summary: "Open a quarantine case"},
			{Path: []string{"quarantine", "release"}, Summary: "Release a quarantine case"},
			{Path: []string{"quarantine", "reject"}, Summary: "Reject a quarantine case"},
			{Path: []string{"quarantine", "list"}, Summary: "List quarantine cases"},
			{Path: []string{"quarantine", "inspect"}, Summary: "Inspect a quarantine case"},
		},
		Families: []pkgmeta.Family{{Name: EventFamily, ProtocolPCID: EventProtocol}},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: EventProtocol, Role: "family-validator", Summary: "Validates append-only quarantine events."},
			{ProtocolPCID: EventProtocol, Role: "domain-behavior", Summary: "Declares quarantine case transitions over explicit evidence."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"quarantine open":    openQuarantine,
			"quarantine release": releaseQuarantine,
			"quarantine reject":  rejectQuarantine,
			"quarantine list":    listQuarantines,
			"quarantine inspect": inspectQuarantine,
		},
		Validators: map[string]kernel.BuiltinValidator{EventFamily: validateEventEnvelope},
	}
}

// Intent: A receiving exception creates a distinct durable case rather than a
// free-text receiving disposition, so later resolution can carry its own
// explicit evidence. Source: DI-hogid
func openQuarantine(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 6 {
		return "", errors.New("usage: quarantine open <case-id> <receiving-id> <receipt-run-id> <actor> <evidence-id> <exception> [notes...]")
	}
	state, err := loadQuarantineState(runtime)
	if err != nil {
		return "", err
	}
	if _, exists := state[args[0]]; exists {
		return "", fmt.Errorf("quarantine case already exists: %s", args[0])
	}
	if _, err := runtime.RunCommand(ctx, []string{"receiving", "inspect", args[1]}); err != nil {
		return "", err
	}
	return appendEvent(ctx, runtime, args[0], eventPayload{CaseID: args[0], Transition: "open", ReceivingID: args[1], ReceiptRunID: args[2], Actor: args[3], EvidenceID: args[4], Exception: args[5], Notes: strings.Join(args[6:], " ")})
}

func releaseQuarantine(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	return resolveQuarantine(ctx, runtime, "release", args)
}

func rejectQuarantine(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	return resolveQuarantine(ctx, runtime, "reject", args)
}

func resolveQuarantine(ctx stdctx.Context, runtime *kernel.Runtime, transition string, args []string) (string, error) {
	if len(args) < 4 {
		return "", fmt.Errorf("usage: quarantine %s <case-id> <event-id> <actor> <evidence-id> [notes...]", transition)
	}
	state, err := loadQuarantineState(runtime)
	if err != nil {
		return "", err
	}
	caseRecord, exists := state[args[0]]
	if !exists {
		return "", fmt.Errorf("unknown quarantine case: %s", args[0])
	}
	if caseRecord.Transition != "open" {
		return "", fmt.Errorf("quarantine case %s is already %s", args[0], caseRecord.Transition)
	}
	return appendEvent(ctx, runtime, args[1], eventPayload{CaseID: args[0], Transition: transition, Actor: args[2], EvidenceID: args[3], Notes: strings.Join(args[4:], " ")})
}

func appendEvent(ctx stdctx.Context, runtime *kernel.Runtime, eventID string, payload eventPayload) (string, error) {
	envelope, err := kernel.NewEnvelope(EventFamily, EventProtocol, eventID, PackageID, payload)
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored quarantine %s %s", payload.Transition, eventID), nil
}

func listQuarantines(_ stdctx.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	state, err := loadQuarantineState(runtime)
	if err != nil {
		return "", err
	}
	ids := make([]string, 0, len(state))
	for caseID := range state {
		ids = append(ids, caseID)
	}
	sort.Strings(ids)
	lines := make([]string, 0, len(ids))
	for _, caseID := range ids {
		caseRecord := state[caseID]
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\tevents=%d", caseRecord.CaseID, caseRecord.ReceivingID, caseRecord.Transition, len(caseRecord.Events)))
	}
	return strings.Join(lines, "\n"), nil
}

func inspectQuarantine(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: quarantine inspect <case-id>")
	}
	state, err := loadQuarantineState(runtime)
	if err != nil {
		return "", err
	}
	caseRecord, exists := state[args[0]]
	if !exists {
		return "", fmt.Errorf("unknown quarantine case: %s", args[0])
	}
	return fmt.Sprintf("case_id: %s\nreceiving_id: %s\ntransition: %s\nevents: %s", caseRecord.CaseID, caseRecord.ReceivingID, caseRecord.Transition, strings.Join(caseRecord.Events, ", ")), nil
}

func loadQuarantineState(runtime *kernel.Runtime) (map[string]*caseState, error) {
	state := map[string]*caseState{}
	for _, entry := range runtime.History() {
		if entry.Envelope.Family != EventFamily {
			continue
		}
		var event eventPayload
		if err := json.Unmarshal(entry.Envelope.Payload, &event); err != nil {
			return nil, err
		}
		caseRecord := state[event.CaseID]
		if caseRecord == nil {
			caseRecord = &caseState{CaseID: event.CaseID}
			state[event.CaseID] = caseRecord
		}
		if event.ReceivingID != "" {
			caseRecord.ReceivingID = event.ReceivingID
		}
		caseRecord.Transition = event.Transition
		caseRecord.Events = append(caseRecord.Events, entry.Envelope.RecordID)
	}
	return state, nil
}

func validateEventEnvelope(envelope records.Envelope) error {
	var event eventPayload
	if err := json.Unmarshal(envelope.Payload, &event); err != nil {
		return err
	}
	if event.CaseID == "" || event.Actor == "" || event.EvidenceID == "" {
		return errors.New("quarantine event requires case_id, actor, and evidence_id")
	}
	switch event.Transition {
	case "open":
		if event.ReceivingID == "" || event.ReceiptRunID == "" || event.Exception == "" || envelope.RecordID != event.CaseID {
			return errors.New("quarantine open requires matching case record, receiving_id, receipt_run_id, and exception")
		}
	case "release", "reject":
		if envelope.RecordID == event.CaseID {
			return errors.New("quarantine resolution requires a distinct event record")
		}
	default:
		return fmt.Errorf("unsupported quarantine transition: %s", event.Transition)
	}
	return nil
}
