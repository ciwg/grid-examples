package maintenancepkg

import (
	stdctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
	pkgmeta "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages"
	contextpkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/context"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/knowledge"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/runs"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/records"
)

const (
	PackageID = "maintenance"

	ItemFamily   = "moks.maintenance.item.v1"
	ItemProtocol = "pcid:moks.maintenance.item.v1"

	ServiceFamily   = "moks.maintenance.service.v1"
	ServiceProtocol = "pcid:moks.maintenance.service.v1"

	FindingFamily   = "moks.maintenance.finding.v1"
	FindingProtocol = "pcid:moks.maintenance.finding.v1"
)

type maintenanceItemPayload struct {
	ResourceID string `json:"resource_id"`
}

type servicePayload struct {
	MaintenanceID string `json:"maintenance_id"`
	RunID         string `json:"run_id"`
	ResourceID    string `json:"resource_id"`
	Performer     string `json:"performer"`
}

type findingPayload struct {
	MaintenanceID string `json:"maintenance_id"`
	ResourceID    string `json:"resource_id"`
	Decision      string `json:"decision"`
	Notes         string `json:"notes,omitempty"`
}

type maintenanceState struct {
	MaintenanceID string
	ResourceID    string
	Title         string
	Summary       string
	Services      []string
	Findings      []string
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party maintenance package",
		Commands: []pkgmeta.Command{
			{Path: []string{"maintenance", "create"}, Summary: "Create a maintenance item"},
			{Path: []string{"maintenance", "list"}, Summary: "List maintenance items"},
			{Path: []string{"maintenance", "inspect"}, Summary: "Inspect a maintenance item"},
			{Path: []string{"maintenance", "record-service"}, Summary: "Record a maintenance service session"},
			{Path: []string{"maintenance", "record-finding"}, Summary: "Record a maintenance finding"},
		},
		Families: []pkgmeta.Family{
			{Name: ItemFamily, ProtocolPCID: ItemProtocol},
			{Name: ServiceFamily, ProtocolPCID: ServiceProtocol},
			{Name: FindingFamily, ProtocolPCID: FindingProtocol},
		},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: ItemProtocol, Role: "family-validator", Summary: "Validates maintenance item records."},
			{ProtocolPCID: ItemProtocol, Role: "domain-behavior", Summary: "Declares maintenance items over the knowledge family."},
			{ProtocolPCID: ServiceProtocol, Role: "family-validator", Summary: "Validates maintenance service records."},
			{ProtocolPCID: ServiceProtocol, Role: "domain-behavior", Summary: "Declares maintenance service records over the runs family."},
			{ProtocolPCID: FindingProtocol, Role: "family-validator", Summary: "Validates maintenance finding records."},
			{ProtocolPCID: FindingProtocol, Role: "domain-behavior", Summary: "Declares maintenance findings tied to a resource."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"maintenance create":         createMaintenance,
			"maintenance list":           listMaintenance,
			"maintenance inspect":        inspectMaintenance,
			"maintenance record-service": recordService,
			"maintenance record-finding": recordFinding,
		},
		Validators: map[string]kernel.BuiltinValidator{
			ItemFamily:    validateMaintenanceItemEnvelope,
			ServiceFamily: validateServiceEnvelope,
			FindingFamily: validateFindingEnvelope,
		},
	}
}

func createMaintenance(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("usage: maintenance create <maintenance-id> <resource-id> <title> <summary>")
	}
	if !resourceExists(runtime, args[1]) {
		return "", fmt.Errorf("unknown resource: %s", args[1])
	}
	knowledgeEnvelope, err := kernel.NewEnvelope(knowledgepkg.ItemFamily, knowledgepkg.ItemProtocol, args[0], PackageID, map[string]any{
		"kind":    "maintenance",
		"title":   args[2],
		"summary": strings.Join(args[3:], " "),
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(knowledgeEnvelope)); err != nil {
		return "", err
	}
	maintenanceEnvelope, err := kernel.NewEnvelope(ItemFamily, ItemProtocol, args[0], PackageID, maintenanceItemPayload{
		ResourceID: args[1],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(maintenanceEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored maintenance %s", args[0]), nil
}

func listMaintenance(_ stdctx.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	state, err := loadMaintenanceState(runtime)
	if err != nil {
		return "", err
	}
	ids := make([]string, 0, len(state))
	for id := range state {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	lines := make([]string, 0, len(ids))
	for _, id := range ids {
		item := state[id]
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\tservices=%d\tfindings=%d", item.MaintenanceID, item.ResourceID, item.Title, len(item.Services), len(item.Findings)))
	}
	return strings.Join(lines, "\n"), nil
}

func inspectMaintenance(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: maintenance inspect <maintenance-id>")
	}
	state, err := loadMaintenanceState(runtime)
	if err != nil {
		return "", err
	}
	item, ok := state[args[0]]
	if !ok {
		return "", fmt.Errorf("unknown maintenance: %s", args[0])
	}
	return fmt.Sprintf(
		"id: %s\nresource_id: %s\ntitle: %s\nsummary: %s\nservices: %s\nfindings: %s",
		item.MaintenanceID,
		item.ResourceID,
		item.Title,
		item.Summary,
		strings.Join(item.Services, ", "),
		strings.Join(item.Findings, ", "),
	), nil
}

func recordService(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 5 {
		return "", errors.New("usage: maintenance record-service <maintenance-id> <run-id> <resource-id> <performer> <outcome> [notes...]")
	}
	if !resourceExists(runtime, args[2]) {
		return "", fmt.Errorf("unknown resource: %s", args[2])
	}
	notes := ""
	if len(args) > 5 {
		notes = strings.Join(args[5:], " ")
	}
	runEnvelope, err := kernel.NewEnvelope(runspkg.RunFamily, runspkg.RunProtocol, args[1], PackageID, map[string]any{
		"item_id": args[0],
		"actor":   args[3],
		"outcome": args[4],
		"notes":   notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(runEnvelope)); err != nil {
		return "", err
	}
	serviceEnvelope, err := kernel.NewEnvelope(ServiceFamily, ServiceProtocol, args[1], PackageID, servicePayload{
		MaintenanceID: args[0],
		RunID:         args[1],
		ResourceID:    args[2],
		Performer:     args[3],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(serviceEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored maintenance service %s", args[1]), nil
}

func recordFinding(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("usage: maintenance record-finding <maintenance-id> <finding-id> <resource-id> <decision> [notes...]")
	}
	if !resourceExists(runtime, args[2]) {
		return "", fmt.Errorf("unknown resource: %s", args[2])
	}
	notes := ""
	if len(args) > 4 {
		notes = strings.Join(args[4:], " ")
	}
	findingEnvelope, err := kernel.NewEnvelope(FindingFamily, FindingProtocol, args[1], PackageID, findingPayload{
		MaintenanceID: args[0],
		ResourceID:    args[2],
		Decision:      args[3],
		Notes:         notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(findingEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored maintenance finding %s", args[1]), nil
}

func loadMaintenanceState(runtime *kernel.Runtime) (map[string]*maintenanceState, error) {
	state := map[string]*maintenanceState{}
	for _, entry := range runtime.History() {
		switch entry.Envelope.Family {
		case ItemFamily:
			var body maintenanceItemPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureMaintenance(state, entry.Envelope.RecordID)
			item.MaintenanceID = entry.Envelope.RecordID
			item.ResourceID = body.ResourceID
		case knowledgepkg.ItemFamily:
			var body map[string]string
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			if body["kind"] != "maintenance" {
				continue
			}
			item := ensureMaintenance(state, entry.Envelope.RecordID)
			item.MaintenanceID = entry.Envelope.RecordID
			item.Title = body["title"]
			item.Summary = body["summary"]
		case ServiceFamily:
			var body servicePayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureMaintenance(state, body.MaintenanceID)
			item.Services = append(item.Services, fmt.Sprintf("%s:%s@%s", body.RunID, body.Performer, body.ResourceID))
		case FindingFamily:
			var body findingPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureMaintenance(state, body.MaintenanceID)
			item.Findings = append(item.Findings, fmt.Sprintf("%s:%s", entry.Envelope.RecordID, body.Decision))
		}
	}
	for _, item := range state {
		sort.Strings(item.Services)
		sort.Strings(item.Findings)
	}
	return state, nil
}

func ensureMaintenance(state map[string]*maintenanceState, id string) *maintenanceState {
	item := state[id]
	if item == nil {
		item = &maintenanceState{MaintenanceID: id}
		state[id] = item
	}
	return item
}

func resourceExists(runtime *kernel.Runtime, resourceID string) bool {
	for _, entry := range runtime.History() {
		if entry.Envelope.Family == contextpkg.ResourceFamily && entry.Envelope.RecordID == resourceID {
			return true
		}
	}
	return false
}

// Intent: Keep maintenance payload validation in the owning package so the
// runtime stays generic while the package adds domain behavior over shared families.
// Source: DI-ramek
func validateMaintenanceItemEnvelope(envelope records.Envelope) error {
	var body maintenanceItemPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.ResourceID) == "" {
		return errors.New("resource_id is required")
	}
	return nil
}

// Intent: Keep maintenance payload validation in the owning package so the
// runtime stays generic while the package adds domain behavior over shared families.
// Source: DI-ramek
func validateServiceEnvelope(envelope records.Envelope) error {
	var body servicePayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.MaintenanceID) == "" || strings.TrimSpace(body.RunID) == "" || strings.TrimSpace(body.ResourceID) == "" || strings.TrimSpace(body.Performer) == "" {
		return errors.New("maintenance_id, run_id, resource_id, and performer are required")
	}
	return nil
}

// Intent: Keep maintenance payload validation in the owning package so the
// runtime stays generic while the package adds domain behavior over shared families.
// Source: DI-ramek
func validateFindingEnvelope(envelope records.Envelope) error {
	var body findingPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.MaintenanceID) == "" || strings.TrimSpace(body.ResourceID) == "" || strings.TrimSpace(body.Decision) == "" {
		return errors.New("maintenance_id, resource_id, and decision are required")
	}
	return nil
}
