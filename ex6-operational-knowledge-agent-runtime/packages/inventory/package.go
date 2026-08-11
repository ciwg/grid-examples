package inventorypkg

import (
	stdctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	pkgmeta "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages"
	contextpkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/context"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/knowledge"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/runs"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
)

var (
	PackageID = "inventory"

	ItemFamily   = "moks.inventory.item.v1"
	ItemProtocol = records.PackageProtocolPCID(ItemFamily)

	CountFamily   = "moks.inventory.count.v1"
	CountProtocol = records.PackageProtocolPCID(CountFamily)

	ReconcileFamily   = "moks.inventory.reconcile.v1"
	ReconcileProtocol = records.PackageProtocolPCID(ReconcileFamily)
)

type inventoryItemPayload struct {
	PlaceID string `json:"place_id"`
}

type countPayload struct {
	InventoryID string `json:"inventory_id"`
	RunID       string `json:"run_id"`
	PlaceID     string `json:"place_id"`
	Counter     string `json:"counter"`
	Quantity    string `json:"quantity"`
}

type reconcilePayload struct {
	InventoryID string `json:"inventory_id"`
	ResourceID  string `json:"resource_id,omitempty"`
	Decision    string `json:"decision"`
	Notes       string `json:"notes,omitempty"`
}

type inventoryState struct {
	InventoryID string
	PlaceID     string
	Title       string
	Summary     string
	Counts      []string
	Reconciles  []string
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party inventory package",
		Commands: []pkgmeta.Command{
			{Path: []string{"inventory", "create"}, Summary: "Create an inventory item"},
			{Path: []string{"inventory", "list"}, Summary: "List inventory items"},
			{Path: []string{"inventory", "inspect"}, Summary: "Inspect an inventory item"},
			{Path: []string{"inventory", "record-count"}, Summary: "Record an inventory count session"},
			{Path: []string{"inventory", "record-reconcile"}, Summary: "Record an inventory reconciliation"},
		},
		Families: []pkgmeta.Family{
			{Name: ItemFamily, ProtocolPCID: ItemProtocol},
			{Name: CountFamily, ProtocolPCID: CountProtocol},
			{Name: ReconcileFamily, ProtocolPCID: ReconcileProtocol},
		},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: ItemProtocol, Role: "family-validator", Summary: "Validates inventory item records."},
			{ProtocolPCID: ItemProtocol, Role: "domain-behavior", Summary: "Declares inventory items over the knowledge family."},
			{ProtocolPCID: CountProtocol, Role: "family-validator", Summary: "Validates inventory count records."},
			{ProtocolPCID: CountProtocol, Role: "domain-behavior", Summary: "Declares inventory count records over the runs family."},
			{ProtocolPCID: ReconcileProtocol, Role: "family-validator", Summary: "Validates inventory reconciliation records."},
			{ProtocolPCID: ReconcileProtocol, Role: "domain-behavior", Summary: "Declares inventory reconciliation records tied to place and resource context."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"inventory create":           createInventory,
			"inventory list":             listInventory,
			"inventory inspect":          inspectInventory,
			"inventory record-count":     recordCount,
			"inventory record-reconcile": recordReconcile,
		},
		Validators: map[string]kernel.BuiltinValidator{
			ItemFamily:      validateInventoryItemEnvelope,
			CountFamily:     validateCountEnvelope,
			ReconcileFamily: validateReconcileEnvelope,
		},
	}
}

func createInventory(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("usage: inventory create <inventory-id> <place-id> <title> <summary>")
	}
	if !placeExists(runtime, args[1]) {
		return "", fmt.Errorf("unknown place: %s", args[1])
	}
	knowledgeEnvelope, err := kernel.NewEnvelope(knowledgepkg.ItemFamily, knowledgepkg.ItemProtocol, args[0], PackageID, map[string]any{
		"kind":    "inventory",
		"title":   args[2],
		"summary": strings.Join(args[3:], " "),
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(knowledgeEnvelope)); err != nil {
		return "", err
	}
	itemEnvelope, err := kernel.NewEnvelope(ItemFamily, ItemProtocol, args[0], PackageID, inventoryItemPayload{
		PlaceID: args[1],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(itemEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored inventory %s", args[0]), nil
}

func listInventory(_ stdctx.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	state, err := loadInventoryState(runtime)
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
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\tcounts=%d\treconciles=%d", item.InventoryID, item.PlaceID, item.Title, len(item.Counts), len(item.Reconciles)))
	}
	return strings.Join(lines, "\n"), nil
}

func inspectInventory(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: inventory inspect <inventory-id>")
	}
	state, err := loadInventoryState(runtime)
	if err != nil {
		return "", err
	}
	item, ok := state[args[0]]
	if !ok {
		return "", fmt.Errorf("unknown inventory item: %s", args[0])
	}
	return fmt.Sprintf(
		"id: %s\nplace_id: %s\ntitle: %s\nsummary: %s\ncounts: %s\nreconciles: %s",
		item.InventoryID,
		item.PlaceID,
		item.Title,
		item.Summary,
		strings.Join(item.Counts, ", "),
		strings.Join(item.Reconciles, ", "),
	), nil
}

func recordCount(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 6 {
		return "", errors.New("usage: inventory record-count <inventory-id> <run-id> <place-id> <counter> <quantity> <outcome> [notes...]")
	}
	if !placeExists(runtime, args[2]) {
		return "", fmt.Errorf("unknown place: %s", args[2])
	}
	notes := ""
	if len(args) > 6 {
		notes = strings.Join(args[6:], " ")
	}
	runEnvelope, err := kernel.NewEnvelope(runspkg.RunFamily, runspkg.RunProtocol, args[1], PackageID, map[string]any{
		"item_id": args[0],
		"actor":   args[3],
		"outcome": args[5],
		"notes":   notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(runEnvelope)); err != nil {
		return "", err
	}
	countEnvelope, err := kernel.NewEnvelope(CountFamily, CountProtocol, args[1], PackageID, countPayload{
		InventoryID: args[0],
		RunID:       args[1],
		PlaceID:     args[2],
		Counter:     args[3],
		Quantity:    args[4],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(countEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored inventory count %s", args[1]), nil
}

func recordReconcile(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: inventory record-reconcile <inventory-id> <event-id> <decision> [resource-id] [notes...]")
	}
	resourceID := ""
	notesStart := 3
	if len(args) > 3 {
		resourceID = args[3]
		if resourceID != "" && !resourceExists(runtime, resourceID) {
			return "", fmt.Errorf("unknown resource: %s", resourceID)
		}
		notesStart = 4
	}
	notes := ""
	if len(args) > notesStart {
		notes = strings.Join(args[notesStart:], " ")
	}
	reconcileEnvelope, err := kernel.NewEnvelope(ReconcileFamily, ReconcileProtocol, args[1], PackageID, reconcilePayload{
		InventoryID: args[0],
		ResourceID:  resourceID,
		Decision:    args[2],
		Notes:       notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(reconcileEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored inventory reconciliation %s", args[1]), nil
}

func loadInventoryState(runtime *kernel.Runtime) (map[string]*inventoryState, error) {
	state := map[string]*inventoryState{}
	for _, entry := range runtime.History() {
		switch entry.Envelope.Family {
		case ItemFamily:
			var body inventoryItemPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureInventory(state, entry.Envelope.RecordID)
			item.InventoryID = entry.Envelope.RecordID
			item.PlaceID = body.PlaceID
		case knowledgepkg.ItemFamily:
			var body map[string]string
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			if body["kind"] != "inventory" {
				continue
			}
			item := ensureInventory(state, entry.Envelope.RecordID)
			item.InventoryID = entry.Envelope.RecordID
			item.Title = body["title"]
			item.Summary = body["summary"]
		case CountFamily:
			var body countPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureInventory(state, body.InventoryID)
			item.Counts = append(item.Counts, fmt.Sprintf("%s:%s=%s@%s", body.RunID, body.Counter, body.Quantity, body.PlaceID))
		case ReconcileFamily:
			var body reconcilePayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureInventory(state, body.InventoryID)
			if strings.TrimSpace(body.ResourceID) != "" {
				item.Reconciles = append(item.Reconciles, fmt.Sprintf("%s:%s@%s", entry.Envelope.RecordID, body.Decision, body.ResourceID))
			} else {
				item.Reconciles = append(item.Reconciles, fmt.Sprintf("%s:%s", entry.Envelope.RecordID, body.Decision))
			}
		}
	}
	for _, item := range state {
		sort.Strings(item.Counts)
		sort.Strings(item.Reconciles)
	}
	return state, nil
}

func ensureInventory(state map[string]*inventoryState, id string) *inventoryState {
	item := state[id]
	if item == nil {
		item = &inventoryState{InventoryID: id}
		state[id] = item
	}
	return item
}

func placeExists(runtime *kernel.Runtime, placeID string) bool {
	for _, entry := range runtime.History() {
		if entry.Envelope.Family == contextpkg.PlaceFamily && entry.Envelope.RecordID == placeID {
			return true
		}
	}
	return false
}

func resourceExists(runtime *kernel.Runtime, resourceID string) bool {
	for _, entry := range runtime.History() {
		if entry.Envelope.Family == contextpkg.ResourceFamily && entry.Envelope.RecordID == resourceID {
			return true
		}
	}
	return false
}

// Intent: Keep inventory payload validation in the owning package so the
// runtime stays generic while the package adds domain behavior over shared families.
// Source: DI-lavom
func validateInventoryItemEnvelope(envelope records.Envelope) error {
	var body inventoryItemPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.PlaceID) == "" {
		return errors.New("place_id is required")
	}
	return nil
}

// Intent: Keep inventory payload validation in the owning package so the
// runtime stays generic while the package adds domain behavior over shared families.
// Source: DI-lavom
func validateCountEnvelope(envelope records.Envelope) error {
	var body countPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.InventoryID) == "" || strings.TrimSpace(body.RunID) == "" || strings.TrimSpace(body.PlaceID) == "" || strings.TrimSpace(body.Counter) == "" || strings.TrimSpace(body.Quantity) == "" {
		return errors.New("inventory_id, run_id, place_id, counter, and quantity are required")
	}
	return nil
}

// Intent: Keep inventory payload validation in the owning package so the
// runtime stays generic while the package adds domain behavior over shared families.
// Source: DI-lavom
func validateReconcileEnvelope(envelope records.Envelope) error {
	var body reconcilePayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.InventoryID) == "" || strings.TrimSpace(body.Decision) == "" {
		return errors.New("inventory_id and decision are required")
	}
	return nil
}
