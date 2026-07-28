package builtin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/records"
)

const (
	OpsPackageID      = "ops-note"
	OpsFamily         = "moks.ops.note.v1"
	OpsFamilyProtocol = "pcid:moks.ops.note.v1"
)

type notePayload struct {
	Title   string `json:"title"`
	BodyRef string `json:"body_ref"`
}

func OpsPackage() kernel.BuiltinPackage {
	manifest := packages.Manifest{
		ID:          OpsPackageID,
		Version:     "0.1.0",
		Description: "Built-in operational note example package",
		Commands: []packages.Command{
			{Path: []string{"ops", "note", "add"}, Summary: "Store an operational note"},
			{Path: []string{"ops", "note", "list"}, Summary: "List operational notes"},
		},
		Families: []packages.Family{
			{Name: OpsFamily, ProtocolPCID: OpsFamilyProtocol},
		},
		Claims: []packages.ImplementationClaim{
			{
				ProtocolPCID: OpsFamilyProtocol,
				Role:         "family-validator",
				Summary:      "Validates and stores operational note envelopes for the built-in note family.",
			},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"ops note add":  addNote,
			"ops note list": listNotes,
		},
		Validators: map[string]kernel.BuiltinValidator{
			OpsFamily: validateNoteEnvelope,
		},
	}
}

func addNote(ctx context.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: ops note add <record-id> <title> <body>")
	}
	recordID := args[0]
	title := args[1]
	body := strings.Join(args[2:], " ")
	bodyRef, err := runtime.PutCAS([]byte(body))
	if err != nil {
		return "", err
	}
	envelope, err := kernel.NewEnvelope(OpsFamily, OpsFamilyProtocol, recordID, OpsPackageID, notePayload{
		Title:   title,
		BodyRef: bodyRef,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored %s", recordID), nil
}

func listNotes(_ context.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	lines := []string{}
	for _, entry := range runtime.History() {
		if entry.Envelope.Family != OpsFamily {
			continue
		}
		var payload notePayload
		if err := json.Unmarshal(entry.Envelope.Payload, &payload); err != nil {
			return "", err
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s", entry.Envelope.RecordID, payload.Title, payload.BodyRef))
	}
	return strings.Join(lines, "\n"), nil
}

// Intent: Keep family-specific payload checks in the owning package while the
// runtime remains responsible for generic durable carriage and relay. Source:
// DI-moksu
func validateNoteEnvelope(envelope records.Envelope) error {
	var payload notePayload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		return err
	}
	if strings.TrimSpace(payload.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(payload.BodyRef) == "" {
		return errors.New("body_ref is required")
	}
	return nil
}
