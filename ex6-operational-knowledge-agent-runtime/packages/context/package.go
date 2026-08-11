package contextpkg

import (
	stdctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	pkgmeta "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
)

var (
	PackageID = "context"

	ResponsibilityFamily   = "moks.context.responsibility.v1"
	ResponsibilityProtocol = records.PackageProtocolPCID(ResponsibilityFamily)

	PlaceFamily   = "moks.context.place.v1"
	PlaceProtocol = records.PackageProtocolPCID(PlaceFamily)

	ResourceFamily   = "moks.context.resource.v1"
	ResourceProtocol = records.PackageProtocolPCID(ResourceFamily)
)

type responsibilityPayload struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type placePayload struct {
	Name     string `json:"name"`
	Summary  string `json:"summary"`
	ParentID string `json:"parent_id,omitempty"`
}

type resourcePayload struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	PlaceID string `json:"place_id,omitempty"`
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party operational context package",
		Commands: []pkgmeta.Command{
			{Path: []string{"context", "responsibility", "create"}, Summary: "Create a responsibility"},
			{Path: []string{"context", "responsibility", "list"}, Summary: "List responsibilities"},
			{Path: []string{"context", "responsibility", "inspect"}, Summary: "Inspect a responsibility"},
			{Path: []string{"context", "place", "create"}, Summary: "Create a place"},
			{Path: []string{"context", "place", "list"}, Summary: "List places"},
			{Path: []string{"context", "place", "inspect"}, Summary: "Inspect a place"},
			{Path: []string{"context", "resource", "create"}, Summary: "Create a resource"},
			{Path: []string{"context", "resource", "list"}, Summary: "List resources"},
			{Path: []string{"context", "resource", "inspect"}, Summary: "Inspect a resource"},
		},
		Families: []pkgmeta.Family{
			{Name: ResponsibilityFamily, ProtocolPCID: ResponsibilityProtocol},
			{Name: PlaceFamily, ProtocolPCID: PlaceProtocol},
			{Name: ResourceFamily, ProtocolPCID: ResourceProtocol},
		},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: ResponsibilityProtocol, Role: "family-validator", Summary: "Validates responsibility records."},
			{ProtocolPCID: PlaceProtocol, Role: "family-validator", Summary: "Validates place records."},
			{ProtocolPCID: ResourceProtocol, Role: "family-validator", Summary: "Validates resource records."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"context responsibility create":  createResponsibility,
			"context responsibility list":    listResponsibilities,
			"context responsibility inspect": inspectResponsibility,
			"context place create":           createPlace,
			"context place list":             listPlaces,
			"context place inspect":          inspectPlace,
			"context resource create":        createResource,
			"context resource list":          listResources,
			"context resource inspect":       inspectResource,
		},
		Validators: map[string]kernel.BuiltinValidator{
			ResponsibilityFamily: validateResponsibilityEnvelope,
			PlaceFamily:          validatePlaceEnvelope,
			ResourceFamily:       validateResourceEnvelope,
		},
	}
}

func createResponsibility(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: context responsibility create <record-id> <title> <summary>")
	}
	envelope, err := kernel.NewEnvelope(ResponsibilityFamily, ResponsibilityProtocol, args[0], PackageID, responsibilityPayload{
		Title:   args[1],
		Summary: strings.Join(args[2:], " "),
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored responsibility %s", args[0]), nil
}

func listResponsibilities(_ stdctx.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	return listFamily(runtime, ResponsibilityFamily, func(payload json.RawMessage, recordID string) (string, error) {
		var body responsibilityPayload
		if err := json.Unmarshal(payload, &body); err != nil {
			return "", err
		}
		return fmt.Sprintf("%s\t%s\t%s", recordID, body.Title, body.Summary), nil
	})
}

func inspectResponsibility(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: context responsibility inspect <record-id>")
	}
	return inspectFamily(runtime, ResponsibilityFamily, args[0], func(payload json.RawMessage, recordID string) (string, error) {
		var body responsibilityPayload
		if err := json.Unmarshal(payload, &body); err != nil {
			return "", err
		}
		return fmt.Sprintf("id: %s\ntitle: %s\nsummary: %s", recordID, body.Title, body.Summary), nil
	})
}

func createPlace(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: context place create <record-id> <name> <summary> [parent-id]")
	}
	parentID := ""
	if len(args) > 3 {
		parentID = args[3]
	}
	envelope, err := kernel.NewEnvelope(PlaceFamily, PlaceProtocol, args[0], PackageID, placePayload{
		Name:     args[1],
		Summary:  args[2],
		ParentID: parentID,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored place %s", args[0]), nil
}

func listPlaces(_ stdctx.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	return listFamily(runtime, PlaceFamily, func(payload json.RawMessage, recordID string) (string, error) {
		var body placePayload
		if err := json.Unmarshal(payload, &body); err != nil {
			return "", err
		}
		return fmt.Sprintf("%s\t%s\t%s\t%s", recordID, body.Name, body.Summary, body.ParentID), nil
	})
}

func inspectPlace(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: context place inspect <record-id>")
	}
	return inspectFamily(runtime, PlaceFamily, args[0], func(payload json.RawMessage, recordID string) (string, error) {
		var body placePayload
		if err := json.Unmarshal(payload, &body); err != nil {
			return "", err
		}
		return fmt.Sprintf("id: %s\nname: %s\nsummary: %s\nparent_id: %s", recordID, body.Name, body.Summary, body.ParentID), nil
	})
}

func createResource(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: context resource create <record-id> <name> <summary> [place-id]")
	}
	placeID := ""
	if len(args) > 3 {
		placeID = args[3]
	}
	envelope, err := kernel.NewEnvelope(ResourceFamily, ResourceProtocol, args[0], PackageID, resourcePayload{
		Name:    args[1],
		Summary: args[2],
		PlaceID: placeID,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored resource %s", args[0]), nil
}

func listResources(_ stdctx.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	return listFamily(runtime, ResourceFamily, func(payload json.RawMessage, recordID string) (string, error) {
		var body resourcePayload
		if err := json.Unmarshal(payload, &body); err != nil {
			return "", err
		}
		return fmt.Sprintf("%s\t%s\t%s\t%s", recordID, body.Name, body.Summary, body.PlaceID), nil
	})
}

func inspectResource(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: context resource inspect <record-id>")
	}
	return inspectFamily(runtime, ResourceFamily, args[0], func(payload json.RawMessage, recordID string) (string, error) {
		var body resourcePayload
		if err := json.Unmarshal(payload, &body); err != nil {
			return "", err
		}
		return fmt.Sprintf("id: %s\nname: %s\nsummary: %s\nplace_id: %s", recordID, body.Name, body.Summary, body.PlaceID), nil
	})
}

func listFamily(runtime *kernel.Runtime, family string, render func(json.RawMessage, string) (string, error)) (string, error) {
	lines := []string{}
	for _, entry := range runtime.History() {
		if entry.Envelope.Family != family {
			continue
		}
		line, err := render(entry.Envelope.Payload, entry.Envelope.RecordID)
		if err != nil {
			return "", err
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), nil
}

func inspectFamily(runtime *kernel.Runtime, family string, recordID string, render func(json.RawMessage, string) (string, error)) (string, error) {
	for _, entry := range runtime.History() {
		if entry.Envelope.Family == family && entry.Envelope.RecordID == recordID {
			return render(entry.Envelope.Payload, recordID)
		}
	}
	return "", fmt.Errorf("unknown %s record: %s", family, recordID)
}

// Intent: Keep operational context payload validation in the owning package so
// the runtime carries the family contract without becoming domain-specific. Source:
// DI-lorup
func validateResponsibilityEnvelope(envelope records.Envelope) error {
	var body responsibilityPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.Title) == "" {
		return errors.New("title is required")
	}
	return nil
}

// Intent: Keep operational context payload validation in the owning package so
// the runtime carries the family contract without becoming domain-specific. Source:
// DI-lorup
func validatePlaceEnvelope(envelope records.Envelope) error {
	var body placePayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}

// Intent: Keep operational context payload validation in the owning package so
// the runtime carries the family contract without becoming domain-specific. Source:
// DI-lorup
func validateResourceEnvelope(envelope records.Envelope) error {
	var body resourcePayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}
