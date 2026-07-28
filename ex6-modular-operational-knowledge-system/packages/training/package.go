package trainingpkg

import (
	stdctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
	pkgmeta "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/knowledge"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/runs"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/records"
)

const (
	PackageID = "training"

	ItemFamily   = "moks.training.item.v1"
	ItemProtocol = "pcid:moks.training.item.v1"

	SessionFamily   = "moks.training.session.v1"
	SessionProtocol = "pcid:moks.training.session.v1"

	CompletionFamily   = "moks.training.completion.v1"
	CompletionProtocol = "pcid:moks.training.completion.v1"
)

type trainingItemPayload struct {
	ItemID string `json:"item_id"`
}

type trainingSessionPayload struct {
	TrainingID string `json:"training_id"`
	RunID      string `json:"run_id"`
	Trainee    string `json:"trainee"`
	Instructor string `json:"instructor"`
}

type completionPayload struct {
	TrainingID string `json:"training_id"`
	Person     string `json:"person"`
	Decision   string `json:"decision"`
	Notes      string `json:"notes,omitempty"`
}

type trainingState struct {
	TrainingID  string
	Title       string
	Summary     string
	Sessions    []string
	Completions []string
}

func Package() kernel.BuiltinPackage {
	manifest := pkgmeta.Manifest{
		ID:          PackageID,
		Version:     "0.1.0",
		Description: "First-party training package",
		Commands: []pkgmeta.Command{
			{Path: []string{"training", "create"}, Summary: "Create a training item"},
			{Path: []string{"training", "list"}, Summary: "List training items"},
			{Path: []string{"training", "inspect"}, Summary: "Inspect a training item"},
			{Path: []string{"training", "record-session"}, Summary: "Record a training session"},
			{Path: []string{"training", "certify"}, Summary: "Record a training completion or certification"},
		},
		Families: []pkgmeta.Family{
			{Name: ItemFamily, ProtocolPCID: ItemProtocol},
			{Name: SessionFamily, ProtocolPCID: SessionProtocol},
			{Name: CompletionFamily, ProtocolPCID: CompletionProtocol},
		},
		Claims: []pkgmeta.ImplementationClaim{
			{ProtocolPCID: ItemProtocol, Role: "domain-behavior", Summary: "Declares training items over the knowledge family."},
			{ProtocolPCID: SessionProtocol, Role: "domain-behavior", Summary: "Declares training sessions over the runs family."},
			{ProtocolPCID: CompletionProtocol, Role: "domain-behavior", Summary: "Declares training completion records."},
		},
	}
	return kernel.BuiltinPackage{
		Manifest: manifest,
		Commands: map[string]kernel.BuiltinCommand{
			"training create":         createTraining,
			"training list":           listTraining,
			"training inspect":        inspectTraining,
			"training record-session": recordSession,
			"training certify":        certifyTraining,
		},
		Validators: map[string]kernel.BuiltinValidator{
			ItemFamily:       validateTrainingItemEnvelope,
			SessionFamily:    validateTrainingSessionEnvelope,
			CompletionFamily: validateCompletionEnvelope,
		},
	}
}

func createTraining(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("usage: training create <training-id> <title> <summary>")
	}
	knowledgeEnvelope, err := kernel.NewEnvelope(knowledgepkg.ItemFamily, knowledgepkg.ItemProtocol, args[0], PackageID, map[string]any{
		"kind":    "training",
		"title":   args[1],
		"summary": strings.Join(args[2:], " "),
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(knowledgeEnvelope)); err != nil {
		return "", err
	}
	trainingEnvelope, err := kernel.NewEnvelope(ItemFamily, ItemProtocol, args[0], PackageID, trainingItemPayload{
		ItemID: args[0],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(trainingEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored training %s", args[0]), nil
}

func listTraining(_ stdctx.Context, runtime *kernel.Runtime, _ []string) (string, error) {
	state, err := loadTrainingState(runtime)
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
		training := state[id]
		lines = append(lines, fmt.Sprintf("%s\t%s\tsessions=%d\tcompletions=%d", training.TrainingID, training.Title, len(training.Sessions), len(training.Completions)))
	}
	return strings.Join(lines, "\n"), nil
}

func inspectTraining(_ stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: training inspect <training-id>")
	}
	state, err := loadTrainingState(runtime)
	if err != nil {
		return "", err
	}
	training, ok := state[args[0]]
	if !ok {
		return "", fmt.Errorf("unknown training: %s", args[0])
	}
	return fmt.Sprintf(
		"id: %s\ntitle: %s\nsummary: %s\nsessions: %s\ncompletions: %s",
		training.TrainingID,
		training.Title,
		training.Summary,
		strings.Join(training.Sessions, ", "),
		strings.Join(training.Completions, ", "),
	), nil
}

func recordSession(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 5 {
		return "", errors.New("usage: training record-session <training-id> <run-id> <trainee> <instructor> <outcome> [notes...]")
	}
	notes := ""
	if len(args) > 5 {
		notes = strings.Join(args[5:], " ")
	}
	runEnvelope, err := kernel.NewEnvelope(runspkg.RunFamily, runspkg.RunProtocol, args[1], PackageID, map[string]any{
		"item_id": args[0],
		"actor":   args[2],
		"outcome": args[4],
		"notes":   notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(runEnvelope)); err != nil {
		return "", err
	}
	sessionEnvelope, err := kernel.NewEnvelope(SessionFamily, SessionProtocol, args[1], PackageID, trainingSessionPayload{
		TrainingID: args[0],
		RunID:      args[1],
		Trainee:    args[2],
		Instructor: args[3],
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(sessionEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored training session %s", args[1]), nil
}

func certifyTraining(ctx stdctx.Context, runtime *kernel.Runtime, args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("usage: training certify <training-id> <event-id> <person> <decision> [notes...]")
	}
	notes := ""
	if len(args) > 4 {
		notes = strings.Join(args[4:], " ")
	}
	completionEnvelope, err := kernel.NewEnvelope(CompletionFamily, CompletionProtocol, args[1], PackageID, completionPayload{
		TrainingID: args[0],
		Person:     args[2],
		Decision:   args[3],
		Notes:      notes,
	})
	if err != nil {
		return "", err
	}
	if _, err := runtime.AppendRecord(ctx, records.MustMarshal(completionEnvelope)); err != nil {
		return "", err
	}
	return fmt.Sprintf("stored training completion %s", args[1]), nil
}

func loadTrainingState(runtime *kernel.Runtime) (map[string]*trainingState, error) {
	state := map[string]*trainingState{}
	for _, entry := range runtime.History() {
		switch entry.Envelope.Family {
		case ItemFamily:
			var body trainingItemPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			training := ensureTraining(state, entry.Envelope.RecordID)
			training.TrainingID = body.ItemID
		case knowledgepkg.ItemFamily:
			var body map[string]string
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			if body["kind"] != "training" {
				continue
			}
			training := ensureTraining(state, entry.Envelope.RecordID)
			training.TrainingID = entry.Envelope.RecordID
			training.Title = body["title"]
			training.Summary = body["summary"]
		case SessionFamily:
			var body trainingSessionPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			training := ensureTraining(state, body.TrainingID)
			training.Sessions = append(training.Sessions, fmt.Sprintf("%s:%s->%s", body.RunID, body.Trainee, body.Instructor))
		case CompletionFamily:
			var body completionPayload
			if err := json.Unmarshal(entry.Envelope.Payload, &body); err != nil {
				return nil, err
			}
			training := ensureTraining(state, body.TrainingID)
			training.Completions = append(training.Completions, fmt.Sprintf("%s:%s", body.Person, body.Decision))
		}
	}
	for _, training := range state {
		sort.Strings(training.Sessions)
		sort.Strings(training.Completions)
	}
	return state, nil
}

func ensureTraining(state map[string]*trainingState, trainingID string) *trainingState {
	training := state[trainingID]
	if training == nil {
		training = &trainingState{TrainingID: trainingID}
		state[trainingID] = training
	}
	return training
}

// Intent: Keep training payload validation in the owning egg so the basket
// stays generic while the package adds domain behavior over shared families.
// Source: DI-sivuk
func validateTrainingItemEnvelope(envelope records.Envelope) error {
	var body trainingItemPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.ItemID) == "" {
		return errors.New("item_id is required")
	}
	return nil
}

// Intent: Keep training payload validation in the owning egg so the basket
// stays generic while the package adds domain behavior over shared families.
// Source: DI-sivuk
func validateTrainingSessionEnvelope(envelope records.Envelope) error {
	var body trainingSessionPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.TrainingID) == "" || strings.TrimSpace(body.RunID) == "" || strings.TrimSpace(body.Trainee) == "" || strings.TrimSpace(body.Instructor) == "" {
		return errors.New("training_id, run_id, trainee, and instructor are required")
	}
	return nil
}

// Intent: Keep training payload validation in the owning egg so the basket
// stays generic while the package adds domain behavior over shared families.
// Source: DI-sivuk
func validateCompletionEnvelope(envelope records.Envelope) error {
	var body completionPayload
	if err := json.Unmarshal(envelope.Payload, &body); err != nil {
		return err
	}
	if strings.TrimSpace(body.TrainingID) == "" || strings.TrimSpace(body.Person) == "" || strings.TrimSpace(body.Decision) == "" {
		return errors.New("training_id, person, and decision are required")
	}
	return nil
}
