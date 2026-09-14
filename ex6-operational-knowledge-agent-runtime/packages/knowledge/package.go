package knowledgepkg

import (
	stdctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	pkgmeta "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
)

var (
	PackageID = "knowledge"

	ItemFamily   = "moks.knowledge.item.v1"
	ItemProtocol = records.PackageProtocolPCID(ItemFamily)

	RevisionFamily   = "moks.knowledge.revision.v1"
	RevisionProtocol = records.PackageProtocolPCID(RevisionFamily)

	LifecycleFamily   = "moks.knowledge.lifecycle.v1"
	LifecycleProtocol = records.PackageProtocolPCID(LifecycleFamily)
)

type itemPayload struct {
	Kind    string `json:"kind"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type revisionPayload struct {
	ItemID   string `json:"item_id"`
	Revision int    `json:"revision"`
	Title    string `json:"title"`
	BodyRef  string `json:"body_ref"`
}

type lifecyclePayload struct {
	ItemID string `json:"item_id"`
	Status string `json:"status"`
	Notes  string `json:"notes,omitempty"`
}

type itemState struct {
	ItemID        string
	Kind          string
	Title         string
	Summary       string
	Status        string
	Revision      int
	RevisionBody  string
	RevisionTitle string
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party knowledge package",
		Commands: []pkgmeta.Command{
			{Path: []string{"knowledge", "item", "create"}, Summary: "Create a knowledge item"},
			{Path: []string{"knowledge", "item", "list"}, Summary: "List knowledge items"},
			{Path: []string{"knowledge", "item", "inspect"}, Summary: "Inspect a knowledge item"},
			{Path: []string{"knowledge", "revision", "snapshot"}, Summary: "Snapshot a knowledge revision"},
			{Path: []string{"knowledge", "item", "approve"}, Summary: "Approve a knowledge item"},
			{Path: []string{"knowledge", "item", "supersede"}, Summary: "Supersede a knowledge item"},
		},
		Families: []pkgmeta.Family{
			{Name: ItemFamily, ProtocolPCID: ItemProtocol},
			{Name: RevisionFamily, ProtocolPCID: RevisionProtocol},
			{Name: LifecycleFamily, ProtocolPCID: LifecycleProtocol},
		},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: ItemProtocol, Role: "family-validator", Summary: "Validates knowledge item records."},
			{ProtocolPCID: RevisionProtocol, Role: "family-validator", Summary: "Validates knowledge revision records."},
			{ProtocolPCID: LifecycleProtocol, Role: "family-validator", Summary: "Validates knowledge lifecycle records."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"knowledge item create":       createItem,
			"knowledge item list":         listItems,
			"knowledge item inspect":      inspectItem,
			"knowledge revision snapshot": snapshotRevision,
			"knowledge item approve":      approveItem,
			"knowledge item supersede":    supersedeItem,
		},
		Validators: map[string]kernel.BuiltinValidator{
			ItemFamily:      validateItemEnvelope,
			RevisionFamily:  validateRevisionEnvelope,
			LifecycleFamily: validateLifecycleEnvelope,
		},
	}
}

func createItem(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("usage: knowledge item create <item-id> <kind> <title> <summary>")
	}
	envelope, err := kernel.NewEnvelope(ItemFamily, ItemProtocol, args[0], PackageID, itemPayload{
		Kind:    args[1],
		Title:   args[2],
		Summary: strings.Join(args[3:], " "),
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored item %s", args[0]), nil
}

func listItems(_ stdctx.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	state, err := loadItemState(runtime)
	if err != nil {
		return "", err
	}
	lines := []string{}
	for _, item := range state {
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s\trev=%d", item.ItemID, item.Kind, item.Title, item.Status, item.Revision))
	}
	return strings.Join(lines, "\n"), nil
}

func inspectItem(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: knowledge item inspect <item-id>")
	}
	state, err := loadItemState(runtime)
	if err != nil {
		return "", err
	}
	item, ok := state[args[0]]
	if !ok {
		return "", fmt.Errorf("unknown knowledge item: %s", args[0])
	}
	return fmt.Sprintf(
		"id: %s\nkind: %s\ntitle: %s\nsummary: %s\nstatus: %s\nrevision: %d\nrevision_title: %s\nbody_ref: %s",
		item.ItemID,
		item.Kind,
		item.Title,
		item.Summary,
		item.Status,
		item.Revision,
		item.RevisionTitle,
		item.RevisionBody,
	), nil
}

func snapshotRevision(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 5 {
		return "", errors.New("usage: knowledge revision snapshot <item-id> <revision-record-id> <revision-number> <title> <body>")
	}
	revisionNumber, err := strconv.Atoi(args[2])
	if err != nil {
		return "", err
	}
	bodyRef, err := runtime.PutCAS([]byte(strings.Join(args[4:], " ")))
	if err != nil {
		return "", err
	}
	envelope, err := kernel.NewEnvelope(RevisionFamily, RevisionProtocol, args[1], PackageID, revisionPayload{
		ItemID:   args[0],
		Revision: revisionNumber,
		Title:    args[3],
		BodyRef:  bodyRef,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored revision %s for %s", args[1], args[0]), nil
}

func approveItem(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	return recordLifecycle(ctx, runtime, "approved", "usage: knowledge item approve <item-id> <event-id> [notes...]", args)
}

func supersedeItem(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	return recordLifecycle(ctx, runtime, "superseded", "usage: knowledge item supersede <item-id> <event-id> [notes...]", args)
}

func recordLifecycle(ctx stdctx.Context, runtime *kernel.Runtime, status string, usage string, args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New(usage)
	}
	notes := ""
	if len(args) > 2 {
		notes = strings.Join(args[2:], " ")
	}
	envelope, err := kernel.NewEnvelope(LifecycleFamily, LifecycleProtocol, args[1], PackageID, lifecyclePayload{
		ItemID: args[0],
		Status: status,
		Notes:  notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(envelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored lifecycle %s for %s", status, args[0]), nil
}

func loadItemState(runtime *kernel.Runtime) (map[string]*itemState, error) {
	items := map[string]*itemState{}
	for _, entry := range runtime.History() {
		switch entry.Envelope.Family {
		case ItemFamily:
			var body itemPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			items[entry.Envelope.RecordID] = &itemState{
				ItemID:  entry.Envelope.RecordID,
				Kind:    body.Kind,
				Title:   body.Title,
				Summary: body.Summary,
				Status:  "draft",
			}
		case RevisionFamily:
			var body revisionPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureItem(items, body.ItemID)
			if body.Revision >= item.Revision {
				item.Revision = body.Revision
				item.RevisionBody = body.BodyRef
				item.RevisionTitle = body.Title
			}
		case LifecycleFamily:
			var body lifecyclePayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			item := ensureItem(items, body.ItemID)
			item.Status = body.Status
		}
	}
	return items, nil
}

func ensureItem(items map[string]*itemState, itemID string) *itemState {
	item := items[itemID]
	if item == nil {
		item = &itemState{ItemID: itemID, Status: "draft"}
		items[itemID] = item
	}
	return item
}

// Intent: Keep knowledge-family payload validation in the owning package so the
// runtime remains generic while the package owns revision and lifecycle rules.
// Source: DI-vakod
func validateItemEnvelope(envelope records.Envelope) error {
	var body itemPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.Kind) == "" {
		return errors.New("kind is required")
	}
	if strings.TrimSpace(body.Title) == "" {
		return errors.New("title is required")
	}
	return nil
}

// Intent: Keep knowledge-family payload validation in the owning package so the
// runtime remains generic while the package owns revision and lifecycle rules.
// Source: DI-vakod
func validateRevisionEnvelope(envelope records.Envelope) error {
	var body revisionPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.ItemID) == "" {
		return errors.New("item_id is required")
	}
	if body.Revision < 1 {
		return errors.New("revision must be >= 1")
	}
	if strings.TrimSpace(body.BodyRef) == "" {
		return errors.New("body_ref is required")
	}
	return nil
}

// Intent: Keep knowledge-family payload validation in the owning package so the
// runtime remains generic while the package owns revision and lifecycle rules.
// Source: DI-vakod
func validateLifecycleEnvelope(envelope records.Envelope) error {
	var body lifecyclePayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.ItemID) == "" {
		return errors.New("item_id is required")
	}
	if body.Status != "approved" && body.Status != "superseded" {
		return fmt.Errorf("unsupported status %q", body.Status)
	}
	return nil
}
