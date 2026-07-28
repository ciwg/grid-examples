package receivingpkg

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

const (
	PackageID = "receiving"

	ItemFamily   = "moks.receiving.item.v1"
	ItemProtocol = "pcid:moks.receiving.item.v1"

	ReceiptFamily   = "moks.receiving.receipt.v1"
	ReceiptProtocol = "pcid:moks.receiving.receipt.v1"

	DispositionFamily   = "moks.receiving.disposition.v1"
	DispositionProtocol = "pcid:moks.receiving.disposition.v1"
)

type receivingItemPayload struct {
	PlaceID string `json:"place_id"`
}

type receiptPayload struct {
	ReceivingID string `json:"receiving_id"`
	RunID       string `json:"run_id"`
	PlaceID     string `json:"place_id"`
	Receiver    string `json:"receiver"`
}

type dispositionPayload struct {
	ReceivingID string `json:"receiving_id"`
	ResourceID  string `json:"resource_id,omitempty"`
	Decision    string `json:"decision"`
	Notes       string `json:"notes,omitempty"`
}

type receivingState struct {
	ReceivingID  string
	PlaceID      string
	Title        string
	Summary      string
	Receipts     []string
	Dispositions []string
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party receiving package",
		Commands: []pkgmeta.Command{
			{Path: []string{"receiving", "create"}, Summary: "Create a receiving item"},
			{Path: []string{"receiving", "list"}, Summary: "List receiving items"},
			{Path: []string{"receiving", "inspect"}, Summary: "Inspect a receiving item"},
			{Path: []string{"receiving", "record-receipt"}, Summary: "Record a receiving session"},
			{Path: []string{"receiving", "record-disposition"}, Summary: "Record a receiving disposition"},
		},
		Families: []pkgmeta.Family{
			{Name: ItemFamily, ProtocolPCID: ItemProtocol},
			{Name: ReceiptFamily, ProtocolPCID: ReceiptProtocol},
			{Name: DispositionFamily, ProtocolPCID: DispositionProtocol},
		},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: ItemProtocol, Role: "family-validator", Summary: "Validates receiving item records."},
			{ProtocolPCID: ItemProtocol, Role: "domain-behavior", Summary: "Declares receiving items over the knowledge family."},
			{ProtocolPCID: ReceiptProtocol, Role: "family-validator", Summary: "Validates receipt records."},
			{ProtocolPCID: ReceiptProtocol, Role: "domain-behavior", Summary: "Declares receipt records over the runs family."},
			{ProtocolPCID: DispositionProtocol, Role: "family-validator", Summary: "Validates receiving disposition records."},
			{ProtocolPCID: DispositionProtocol, Role: "domain-behavior", Summary: "Declares receiving dispositions tied to place and resource context."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"receiving create":             createReceiving,
			"receiving list":               listReceiving,
			"receiving inspect":            inspectReceiving,
			"receiving record-receipt":     recordReceipt,
			"receiving record-disposition": recordDisposition,
		},
		Validators: map[string]kernel.BuiltinValidator{
			ItemFamily:        validateReceivingItemEnvelope,
			ReceiptFamily:     validateReceiptEnvelope,
			DispositionFamily: validateDispositionEnvelope,
		},
	}
}

func createReceiving(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("usage: receiving create <receiving-id> <place-id> <title> <summary>")
	}
	if !placeExists(runtime, args[1]) {
		return "", fmt.Errorf("unknown place: %s", args[1])
	}
	knowledgeEnvelope, err := kernel.NewEnvelope(knowledgepkg.ItemFamily, knowledgepkg.ItemProtocol, args[0], PackageID, map[string]any{
		"kind":    "receiving",
		"title":   args[2],
		"summary": strings.Join(args[3:], " "),
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(knowledgeEnvelope)); err != nil {
		return "", err
	}
	itemEnvelope, err := kernel.NewEnvelope(ItemFamily, ItemProtocol, args[0], PackageID, receivingItemPayload{
		PlaceID: args[1],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(itemEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored receiving %s", args[0]), nil
}

func listReceiving(_ stdctx.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	state, err := loadReceivingState(runtime)
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
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\treceipts=%d\tdispositions=%d", item.ReceivingID, item.PlaceID, item.Title, len(item.Receipts), len(item.Dispositions)))
	}
	return strings.Join(lines, "\n"), nil
}

func inspectReceiving(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: receiving inspect <receiving-id>")
	}
	state, err := loadReceivingState(runtime)
	if err != nil {
		return "", err
	}
	item, ok := state[args[0]]
	if !ok {
		return "", fmt.Errorf("unknown receiving item: %s", args[0])
	}
	return fmt.Sprintf(
		"id: %s\nplace_id: %s\ntitle: %s\nsummary: %s\nreceipts: %s\ndispositions: %s",
		item.ReceivingID,
		item.PlaceID,
		item.Title,
		item.Summary,
		strings.Join(item.Receipts, ", "),
		strings.Join(item.Dispositions, ", "),
	), nil
}

func recordReceipt(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 5 {
		return "", errors.New("usage: receiving record-receipt <receiving-id> <run-id> <place-id> <receiver> <outcome> [notes...]")
	}
	if !placeExists(runtime, args[2]) {
		return "", fmt.Errorf("unknown place: %s", args[2])
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
	receiptEnvelope, err := kernel.NewEnvelope(ReceiptFamily, ReceiptProtocol, args[1], PackageID, receiptPayload{
		ReceivingID: args[0],
		RunID:       args[1],
		PlaceID:     args[2],
		Receiver:    args[3],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(receiptEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored receipt %s", args[1]), nil
}

func recordDisposition(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: receiving record-disposition <receiving-id> <event-id> <decision> [resource-id] [notes...]")
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
	dispositionEnvelope, err := kernel.NewEnvelope(DispositionFamily, DispositionProtocol, args[1], PackageID, dispositionPayload{
		ReceivingID: args[0],
		ResourceID:  resourceID,
		Decision:    args[2],
		Notes:       notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(dispositionEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored disposition %s", args[1]), nil
}

func loadReceivingState(runtime *kernel.Runtime) (map[string]*receivingState, error) {
	state := map[string]*receivingState{}
	for _, entry := range runtime.History() {
		switch entry.Envelope.Family {
		case ItemFamily:
			var body receivingItemPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureReceiving(state, entry.Envelope.RecordID)
			item.ReceivingID = entry.Envelope.RecordID
			item.PlaceID = body.PlaceID
		case knowledgepkg.ItemFamily:
			var body map[string]string
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			if body["kind"] != "receiving" {
				continue
			}
			item := ensureReceiving(state, entry.Envelope.RecordID)
			item.ReceivingID = entry.Envelope.RecordID
			item.Title = body["title"]
			item.Summary = body["summary"]
		case ReceiptFamily:
			var body receiptPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureReceiving(state, body.ReceivingID)
			item.Receipts = append(item.Receipts, fmt.Sprintf("%s:%s@%s", body.RunID, body.Receiver, body.PlaceID))
		case DispositionFamily:
			var body dispositionPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureReceiving(state, body.ReceivingID)
			if strings.TrimSpace(body.ResourceID) != "" {
				item.Dispositions = append(item.Dispositions, fmt.Sprintf("%s:%s@%s", entry.Envelope.RecordID, body.Decision, body.ResourceID))
			} else {
				item.Dispositions = append(item.Dispositions, fmt.Sprintf("%s:%s", entry.Envelope.RecordID, body.Decision))
			}
		}
	}
	for _, item := range state {
		sort.Strings(item.Receipts)
		sort.Strings(item.Dispositions)
	}
	return state, nil
}

func ensureReceiving(state map[string]*receivingState, id string) *receivingState {
	item := state[id]
	if item == nil {
		item = &receivingState{ReceivingID: id}
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

// Intent: Keep receiving payload validation in the owning package so the
// runtime stays generic while the package adds domain behavior over shared families.
// Source: DI-zibek
func validateReceivingItemEnvelope(envelope records.Envelope) error {
	var body receivingItemPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.PlaceID) == "" {
		return errors.New("place_id is required")
	}
	return nil
}

// Intent: Keep receiving payload validation in the owning package so the
// runtime stays generic while the package adds domain behavior over shared families.
// Source: DI-zibek
func validateReceiptEnvelope(envelope records.Envelope) error {
	var body receiptPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.ReceivingID) == "" || strings.TrimSpace(body.RunID) == "" || strings.TrimSpace(body.PlaceID) == "" || strings.TrimSpace(body.Receiver) == "" {
		return errors.New("receiving_id, run_id, place_id, and receiver are required")
	}
	return nil
}

// Intent: Keep receiving payload validation in the owning package so the
// runtime stays generic while the package adds domain behavior over shared families.
// Source: DI-zibek
func validateDispositionEnvelope(envelope records.Envelope) error {
	var body dispositionPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.ReceivingID) == "" || strings.TrimSpace(body.Decision) == "" {
		return errors.New("receiving_id and decision are required")
	}
	return nil
}
