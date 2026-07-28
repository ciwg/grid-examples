package linkspkg

import (
	stdctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
	pkgmeta "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/records"
)

const (
	PackageID = "links"

	TypedFamily   = "moks.links.typed.v1"
	TypedProtocol = "pcid:moks.links.typed.v1"
)

type typedLinkPayload struct {
	FromType string `json:"from_type"`
	FromID   string `json:"from_id"`
	ToType   string `json:"to_type"`
	ToID     string `json:"to_id"`
	Relation string `json:"relation"`
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party typed links package",
		Commands: []pkgmeta.Command{
			{Path: []string{"links", "create"}, Summary: "Create a typed link"},
			{Path: []string{"links", "inspect"}, Summary: "Inspect a typed link"},
		},
		Families: []pkgmeta.Family{
			{Name: TypedFamily, ProtocolPCID: TypedProtocol},
		},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: TypedProtocol, Role: "family-validator", Summary: "Validates typed link records."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"links create":  createLink,
			"links inspect": inspectLink,
		},
		Validators: map[string]kernel.BuiltinValidator{
			TypedFamily: validateLinkEnvelope,
		},
	}
}

func createLink(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 6 {
		return "", errors.New("usage: links create <link-id> <from-type> <from-id> <to-type> <to-id> <relation>")
	}
	envelope, err := kernel.NewEnvelope(TypedFamily, TypedProtocol, args[0], PackageID, typedLinkPayload{
		FromType: args[1],
		FromID:   args[2],
		ToType:   args[3],
		ToID:     args[4],
		Relation: args[5],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored link %s", args[0]), nil
}

func inspectLink(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: links inspect <link-id>")
	}
	for _, entry := range runtime.History() {
		if entry.Envelope.Family != TypedFamily || entry.Envelope.RecordID != args[0] {
			continue
		}
		var body typedLinkPayload
		if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"id: %s\nfrom: %s %s\nto: %s %s\nrelation: %s",
			args[0],
			body.FromType,
			body.FromID,
			body.ToType,
			body.ToID,
			body.Relation,
		), nil
	}
	return "", fmt.Errorf("unknown link: %s", args[0])
}

// Intent: Keep typed-link payload validation in the owning egg so the basket
// remains generic while the package owns relationship semantics. Source:
// DI-figar
func validateLinkEnvelope(envelope records.Envelope) error {
	var body typedLinkPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.FromType) == "" || strings.TrimSpace(body.FromID) == "" {
		return errors.New("from_type and from_id are required")
	}
	if strings.TrimSpace(body.ToType) == "" || strings.TrimSpace(body.ToID) == "" {
		return errors.New("to_type and to_id are required")
	}
	if strings.TrimSpace(body.Relation) == "" {
		return errors.New("relation is required")
	}
	return nil
}
