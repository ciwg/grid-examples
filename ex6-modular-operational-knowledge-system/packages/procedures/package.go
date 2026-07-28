package procedurespkg

import (
	stdctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
	pkgmeta "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/knowledge"
	linkspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/links"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/runs"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/records"
)

const (
	PackageID = "procedures"

	ItemFamily   = "moks.procedures.item.v1"
	ItemProtocol = "pcid:moks.procedures.item.v1"

	UseFamily   = "moks.procedures.use.v1"
	UseProtocol = "pcid:moks.procedures.use.v1"
)

type procedureItemPayload struct {
	ItemID string `json:"item_id"`
}

type procedureUsePayload struct {
	ProcedureID string `json:"procedure_id"`
	RunID       string `json:"run_id"`
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party procedures package",
		Commands: []pkgmeta.Command{
			{Path: []string{"procedures", "create"}, Summary: "Create a procedure item"},
			{Path: []string{"procedures", "inspect"}, Summary: "Inspect a procedure item"},
			{Path: []string{"procedures", "record-use"}, Summary: "Record a procedure use run"},
		},
		Families: []pkgmeta.Family{
			{Name: ItemFamily, ProtocolPCID: ItemProtocol},
			{Name: UseFamily, ProtocolPCID: UseProtocol},
		},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: ItemProtocol, Role: "domain-behavior", Summary: "Declares procedure items over the knowledge family."},
			{ProtocolPCID: UseProtocol, Role: "domain-behavior", Summary: "Declares procedure use records over the runs family."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"procedures create":     createProcedure,
			"procedures inspect":    inspectProcedure,
			"procedures record-use": recordUse,
		},
		Validators: map[string]kernel.BuiltinValidator{
			ItemFamily: validateProcedureItemEnvelope,
			UseFamily:  validateProcedureUseEnvelope,
		},
	}
}

func createProcedure(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: procedures create <procedure-id> <title> <summary>")
	}
	knowledgeEnvelope, err := kernel.NewEnvelope(knowledgepkg.ItemFamily, knowledgepkg.ItemProtocol, args[0], PackageID, map[string]any{
		"kind":    "procedure",
		"title":   args[1],
		"summary": strings.Join(args[2:], " "),
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(knowledgeEnvelope)); err != nil {
		return "", err
	}
	procedureEnvelope, err := kernel.NewEnvelope(ItemFamily, ItemProtocol, args[0], PackageID, procedureItemPayload{
		ItemID: args[0],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(procedureEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored procedure %s", args[0]), nil
}

func inspectProcedure(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: procedures inspect <procedure-id>")
	}
	foundProcedure := false
	title := ""
	summary := ""
	uses := []string{}
	links := []string{}
	for _, entry := range runtime.History() {
		switch entry.Envelope.Family {
		case ItemFamily:
			if entry.Envelope.RecordID != args[0] {
				continue
			}
			foundProcedure = true
		case knowledgepkg.ItemFamily:
			if entry.Envelope.RecordID != args[0] {
				continue
			}
			var body map[string]string
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return "", err
			}
			if body["kind"] == "procedure" {
				title = body["title"]
				summary = body["summary"]
			}
		case UseFamily:
			var body procedureUsePayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return "", err
			}
			if body.ProcedureID == args[0] {
				uses = append(uses, body.RunID)
			}
		case linkspkg.TypedFamily:
			var body map[string]string
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return "", err
			}
			if body["from_id"] == args[0] || body["to_id"] == args[0] {
				links = append(links, entry.Envelope.RecordID)
			}
		}
	}
	if !foundProcedure {
		return "", fmt.Errorf("unknown procedure: %s", args[0])
	}
	return fmt.Sprintf(
		"id: %s\ntitle: %s\nsummary: %s\nuses: %s\nlinks: %s",
		args[0],
		title,
		summary,
		strings.Join(uses, ", "),
		strings.Join(links, ", "),
	), nil
}

func recordUse(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("usage: procedures record-use <procedure-id> <run-id> <actor> <outcome> [notes...]")
	}
	notes := ""
	if len(args) > 4 {
		notes = strings.Join(args[4:], " ")
	}
	runEnvelope, err := kernel.NewEnvelope(runspkg.RunFamily, runspkg.RunProtocol, args[1], PackageID, map[string]any{
		"item_id": args[0],
		"actor":   args[2],
		"outcome": args[3],
		"notes":   notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(runEnvelope)); err != nil {
		return "", err
	}
	useEnvelope, err := kernel.NewEnvelope(UseFamily, UseProtocol, args[1], PackageID, procedureUsePayload{
		ProcedureID: args[0],
		RunID:       args[1],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(useEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored procedure use %s", args[1]), nil
}

// Intent: Keep procedure payload validation in the owning package so the
// runtime remains generic while the package adds domain meaning above shared families.
// Source: DI-tusav
func validateProcedureItemEnvelope(envelope records.Envelope) error {
	var body procedureItemPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.ItemID) == "" {
		return errors.New("item_id is required")
	}
	return nil
}

// Intent: Keep procedure payload validation in the owning package so the
// runtime remains generic while the package adds domain meaning above shared families.
// Source: DI-tusav
func validateProcedureUseEnvelope(envelope records.Envelope) error {
	var body procedureUsePayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.ProcedureID) == "" || strings.TrimSpace(body.RunID) == "" {
		return errors.New("procedure_id and run_id are required")
	}
	return nil
}
