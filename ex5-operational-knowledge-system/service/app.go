package service

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const defaultTeam = "OPS"

type App struct {
	dataRoot string
	store    *Store

	mu               sync.Mutex
	responsibilities map[string]*Responsibility
	items            map[string]*KnowledgeItem
	runs             map[string]*RunRecord
	links            map[string]*Link
	approvals        map[string]*Approval
	nextSequence     uint64
	nextNumbers      map[string]int
}

func NewApp(dataRoot string) (*App, error) {
	store, events, err := OpenStore(dataRoot)
	if err != nil {
		return nil, err
	}
	app := &App{
		dataRoot:         dataRoot,
		store:            store,
		responsibilities: map[string]*Responsibility{},
		items:            map[string]*KnowledgeItem{},
		runs:             map[string]*RunRecord{},
		links:            map[string]*Link{},
		approvals:        map[string]*Approval{},
		nextNumbers:      map[string]int{},
	}
	for _, event := range events {
		if err := app.applyEventLocked(event); err != nil {
			return nil, fmt.Errorf("replay event %d: %w", event.Sequence, err)
		}
		if event.Sequence > app.nextSequence {
			app.nextSequence = event.Sequence
		}
	}
	return app, nil
}

func (app *App) Meta() Meta {
	return Meta{
		DataRoot:          app.dataRoot,
		KnowledgeKinds:    []string{KnowledgeKindProcedure, KnowledgeKindTraining, KnowledgeKindMaintenance},
		RunKinds:          []string{RunKindProcedure, RunKindTraining, RunKindMaintenance},
		ApprovalDecisions: []string{DecisionApproved, DecisionRejected, DecisionNoted},
	}
}

func (app *App) Dashboard() Dashboard {
	app.mu.Lock()
	defer app.mu.Unlock()
	out := Dashboard{
		Responsibilities: len(app.responsibilities),
		Approvals:        len(app.approvals),
		Links:            len(app.links),
	}
	for _, item := range app.items {
		switch item.Kind {
		case KnowledgeKindProcedure:
			out.Procedures++
		case KnowledgeKindTraining:
			out.TrainingItems++
		case KnowledgeKindMaintenance:
			out.MaintenanceItems++
		}
	}
	for _, run := range app.runs {
		switch run.Kind {
		case RunKindProcedure:
			out.ProcedureRuns++
		case RunKindTraining:
			out.TrainingRuns++
		case RunKindMaintenance:
			out.MaintenanceRuns++
		}
		out.Evidence += len(run.Evidence)
	}
	return out
}

func (app *App) ListResponsibilities() []Responsibility {
	app.mu.Lock()
	defer app.mu.Unlock()
	out := make([]Responsibility, 0, len(app.responsibilities))
	for _, responsibility := range app.responsibilities {
		out = append(out, cloneResponsibility(responsibility))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (app *App) GetResponsibility(id string) (Responsibility, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	record, ok := app.responsibilities[id]
	if !ok {
		return Responsibility{}, fmt.Errorf("responsibility %q not found", id)
	}
	return cloneResponsibility(record), nil
}

// Intent: Keep responsibilities as first-class durable records instead of
// hiding them inside procedure metadata, so workflow roles and operational
// duties stay independently linkable across procedures, training, and
// maintenance. Source: DI-kovup
func (app *App) CreateResponsibility(actor string, title string, summary string, roleKeys []string, tags []string) (Responsibility, error) {
	if strings.TrimSpace(actor) == "" {
		return Responsibility{}, fmt.Errorf("actor is required")
	}
	if strings.TrimSpace(title) == "" {
		return Responsibility{}, fmt.Errorf("title is required")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	id := app.nextIDLocked("RESP")
	event := OperationalEvent{
		EntityType: "responsibility",
		EntityID:   id,
		Type:       "responsibility_created",
		Actor:      actor,
		Title:      strings.TrimSpace(title),
		Summary:    strings.TrimSpace(summary),
		RoleKeys:   normalizeStrings(roleKeys),
		Tags:       normalizeStrings(tags),
		Team:       defaultTeam,
	}
	if err := app.appendEventLocked(event); err != nil {
		return Responsibility{}, err
	}
	return cloneResponsibility(app.responsibilities[id]), nil
}

func (app *App) ListKnowledgeItems(kind string) ([]KnowledgeItem, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	out := []KnowledgeItem{}
	for _, item := range app.items {
		if kind != "" && item.Kind != kind {
			continue
		}
		out = append(out, cloneKnowledgeItem(item))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (app *App) GetKnowledgeItem(id string) (KnowledgeItem, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	item, ok := app.items[id]
	if !ok {
		return KnowledgeItem{}, fmt.Errorf("knowledge item %q not found", id)
	}
	return cloneKnowledgeItem(item), nil
}

// Intent: Keep procedures, training content, and maintenance content as
// hybrid knowledge items with structured metadata plus revisioned shared text,
// so operational records and collaborative knowledge can coexist in one tool.
// Source: DI-kovup
func (app *App) CreateKnowledgeItem(actor string, kind string, title string, summary string, body string, tags []string, responsibilityIDs []string) (KnowledgeItem, error) {
	if err := validateKnowledgeKind(kind); err != nil {
		return KnowledgeItem{}, err
	}
	if strings.TrimSpace(actor) == "" {
		return KnowledgeItem{}, fmt.Errorf("actor is required")
	}
	if strings.TrimSpace(title) == "" {
		return KnowledgeItem{}, fmt.Errorf("title is required")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	if err := app.validateResponsibilitiesLocked(responsibilityIDs); err != nil {
		return KnowledgeItem{}, err
	}
	id := app.nextIDLocked(kindPrefix(kind))
	event := OperationalEvent{
		EntityType:        "knowledge_item",
		EntityID:          id,
		Type:              "knowledge_item_created",
		Actor:             actor,
		Kind:              kind,
		Title:             strings.TrimSpace(title),
		Summary:           strings.TrimSpace(summary),
		Body:              normalizeBody(body),
		Tags:              normalizeStrings(tags),
		ResponsibilityIDs: normalizeStrings(responsibilityIDs),
		Revision:          1,
	}
	if err := app.appendEventLocked(event); err != nil {
		return KnowledgeItem{}, err
	}
	return cloneKnowledgeItem(app.items[id]), nil
}

func (app *App) AddRevision(actor string, itemID string, title string, summary string, body string, tags []string) (KnowledgeItem, error) {
	if strings.TrimSpace(actor) == "" {
		return KnowledgeItem{}, fmt.Errorf("actor is required")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	item, ok := app.items[itemID]
	if !ok {
		return KnowledgeItem{}, fmt.Errorf("knowledge item %q not found", itemID)
	}
	revision := item.CurrentRevision + 1
	event := OperationalEvent{
		EntityType: item.Kind + "_revision",
		EntityID:   itemID,
		Type:       "revision_added",
		Actor:      actor,
		Kind:       item.Kind,
		Title:      strings.TrimSpace(title),
		Summary:    strings.TrimSpace(summary),
		Body:       normalizeBody(body),
		Tags:       normalizeStrings(tags),
		Revision:   revision,
	}
	if err := app.appendEventLocked(event); err != nil {
		return KnowledgeItem{}, err
	}
	return cloneKnowledgeItem(app.items[itemID]), nil
}

func (app *App) ListRuns(kind string) ([]RunRecord, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	out := []RunRecord{}
	for _, run := range app.runs {
		if kind != "" && run.Kind != kind {
			continue
		}
		out = append(out, cloneRun(run))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (app *App) GetRun(id string) (RunRecord, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	run, ok := app.runs[id]
	if !ok {
		return RunRecord{}, fmt.Errorf("run %q not found", id)
	}
	return cloneRun(run), nil
}

// Intent: Use performed runs as the durable anchor for completed work so every
// procedure/training/maintenance execution can point back to the exact
// revision, evidence, and responsibilities involved. Source: DI-kovup;
// DI-zuvob
func (app *App) RecordRun(actor string, kind string, itemID string, revision int, outcome string, notes string, machine string, location string, responsibilityIDs []string) (RunRecord, error) {
	if err := validateRunKind(kind); err != nil {
		return RunRecord{}, err
	}
	if strings.TrimSpace(actor) == "" {
		return RunRecord{}, fmt.Errorf("actor is required")
	}
	if strings.TrimSpace(itemID) == "" {
		return RunRecord{}, fmt.Errorf("item_id is required")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	item, ok := app.items[itemID]
	if !ok {
		return RunRecord{}, fmt.Errorf("knowledge item %q not found", itemID)
	}
	if itemKindForRun(kind) != item.Kind {
		return RunRecord{}, fmt.Errorf("item %q has kind %s, not %s", itemID, item.Kind, itemKindForRun(kind))
	}
	if revision <= 0 || revision > len(item.Revisions) {
		return RunRecord{}, fmt.Errorf("revision %d not found for %s", revision, itemID)
	}
	if err := app.validateResponsibilitiesLocked(responsibilityIDs); err != nil {
		return RunRecord{}, err
	}
	id := app.nextIDLocked("RUN")
	event := OperationalEvent{
		EntityType:        "run",
		EntityID:          id,
		Type:              "run_recorded",
		Actor:             actor,
		Kind:              kind,
		TargetID:          itemID,
		Revision:          revision,
		Outcome:           strings.TrimSpace(outcome),
		Notes:             strings.TrimSpace(notes),
		Machine:           strings.TrimSpace(machine),
		Location:          strings.TrimSpace(location),
		ResponsibilityIDs: normalizeStrings(responsibilityIDs),
	}
	if err := app.appendEventLocked(event); err != nil {
		return RunRecord{}, err
	}
	return cloneRun(app.runs[id]), nil
}

func (app *App) AddEvidence(actor string, runID string, summary string, facts map[string]string, attachmentName string, attachmentBody []byte) (RunRecord, error) {
	if strings.TrimSpace(actor) == "" {
		return RunRecord{}, fmt.Errorf("actor is required")
	}
	if strings.TrimSpace(summary) == "" {
		return RunRecord{}, fmt.Errorf("summary is required")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	run, ok := app.runs[runID]
	if !ok {
		return RunRecord{}, fmt.Errorf("run %q not found", runID)
	}
	event := OperationalEvent{
		EntityType: "run",
		EntityID:   runID,
		Type:       "evidence_added",
		Actor:      actor,
		Summary:    strings.TrimSpace(summary),
		Facts:      normalizeFacts(facts),
	}
	if len(attachmentBody) > 0 {
		path, size, err := app.store.SaveAttachment(runID, attachmentName, attachmentBody)
		if err != nil {
			return RunRecord{}, err
		}
		event.AttachmentName = filepath.Base(attachmentName)
		event.AttachmentPath = path
		event.AttachmentSize = size
	}
	if err := app.appendEventLocked(event); err != nil {
		return RunRecord{}, err
	}
	return cloneRun(run), nil
}

func (app *App) RecordApproval(actor string, targetType string, targetID string, revision int, role string, decision string, notes string) error {
	if strings.TrimSpace(actor) == "" {
		return fmt.Errorf("actor is required")
	}
	if strings.TrimSpace(targetType) == "" || strings.TrimSpace(targetID) == "" {
		return fmt.Errorf("target_type and target_id are required")
	}
	if strings.TrimSpace(role) == "" {
		return fmt.Errorf("role is required")
	}
	if err := validateDecision(decision); err != nil {
		return err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	switch targetType {
	case "knowledge_item":
		if _, ok := app.items[targetID]; !ok {
			return fmt.Errorf("knowledge item %q not found", targetID)
		}
	case "run":
		if _, ok := app.runs[targetID]; !ok {
			return fmt.Errorf("run %q not found", targetID)
		}
	default:
		return fmt.Errorf("unsupported target_type %q", targetType)
	}
	event := OperationalEvent{
		EntityType: "approval",
		EntityID:   app.nextIDLocked("APR"),
		Type:       "approval_recorded",
		Actor:      actor,
		TargetType: targetType,
		TargetID:   targetID,
		Revision:   revision,
		Role:       strings.TrimSpace(role),
		Decision:   decision,
		Notes:      strings.TrimSpace(notes),
	}
	return app.appendEventLocked(event)
}

func (app *App) AddLink(actor string, fromType string, fromID string, toType string, toID string, relation string, notes string) error {
	if strings.TrimSpace(actor) == "" {
		return fmt.Errorf("actor is required")
	}
	if strings.TrimSpace(fromType) == "" || strings.TrimSpace(fromID) == "" || strings.TrimSpace(toType) == "" || strings.TrimSpace(toID) == "" {
		return fmt.Errorf("link endpoints are required")
	}
	if strings.TrimSpace(relation) == "" {
		return fmt.Errorf("relation is required")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	event := OperationalEvent{
		EntityType: "link",
		EntityID:   app.nextIDLocked("LINK"),
		Type:       "link_added",
		Actor:      actor,
		FromType:   strings.TrimSpace(fromType),
		FromID:     strings.TrimSpace(fromID),
		ToType:     strings.TrimSpace(toType),
		ToID:       strings.TrimSpace(toID),
		Relation:   strings.TrimSpace(relation),
		Notes:      strings.TrimSpace(notes),
	}
	return app.appendEventLocked(event)
}

func (app *App) Search(query string) map[string]any {
	app.mu.Lock()
	defer app.mu.Unlock()
	query = strings.ToLower(strings.TrimSpace(query))
	resp := []Responsibility{}
	items := []KnowledgeItem{}
	runs := []RunRecord{}
	for _, record := range app.responsibilities {
		if query == "" || strings.Contains(strings.ToLower(record.Title+" "+record.Summary), query) {
			resp = append(resp, cloneResponsibility(record))
		}
	}
	for _, record := range app.items {
		if query == "" || strings.Contains(strings.ToLower(record.Title+" "+record.Summary), query) {
			items = append(items, cloneKnowledgeItem(record))
		}
	}
	for _, record := range app.runs {
		if query == "" || strings.Contains(strings.ToLower(record.Outcome+" "+record.Notes+" "+record.Machine+" "+record.Location), query) {
			runs = append(runs, cloneRun(record))
		}
	}
	return map[string]any{
		"responsibilities": resp,
		"items":            items,
		"runs":             runs,
	}
}

func (app *App) appendEventLocked(event OperationalEvent) error {
	app.nextSequence++
	event.Sequence = app.nextSequence
	event.Timestamp = time.Now().Format(time.RFC3339)
	if err := app.store.AppendEvent(event); err != nil {
		app.nextSequence--
		return err
	}
	return app.applyEventLocked(event)
}

func (app *App) applyEventLocked(event OperationalEvent) error {
	switch event.Type {
	case "responsibility_created":
		record := &Responsibility{
			ID:             event.EntityID,
			Title:          event.Title,
			Summary:        event.Summary,
			Team:           defaultTeam,
			Tags:           append([]string(nil), event.Tags...),
			CreatedAt:      event.Timestamp,
			UpdatedAt:      event.Timestamp,
			LinkedRoleKeys: append([]string(nil), event.RoleKeys...),
			Timeline:       []OperationalEvent{event},
		}
		app.responsibilities[event.EntityID] = record
	case "knowledge_item_created":
		item := &KnowledgeItem{
			ID:                event.EntityID,
			Kind:              event.Kind,
			Title:             event.Title,
			Summary:           event.Summary,
			Tags:              append([]string(nil), event.Tags...),
			ResponsibilityIDs: append([]string(nil), event.ResponsibilityIDs...),
			CreatedAt:         event.Timestamp,
			UpdatedAt:         event.Timestamp,
			CurrentRevision:   event.Revision,
			Revisions: []KnowledgeRevision{{
				Number:    event.Revision,
				Title:     event.Title,
				Summary:   event.Summary,
				Body:      event.Body,
				Tags:      append([]string(nil), event.Tags...),
				Author:    event.Actor,
				CreatedAt: event.Timestamp,
			}},
			Timeline: []OperationalEvent{event},
		}
		app.items[event.EntityID] = item
		for _, responsibilityID := range event.ResponsibilityIDs {
			if responsibility, ok := app.responsibilities[responsibilityID]; ok {
				responsibility.LinkedItemIDs = appendUnique(responsibility.LinkedItemIDs, event.EntityID)
				responsibility.UpdatedAt = event.Timestamp
			}
		}
	case "revision_added":
		item, ok := app.items[event.EntityID]
		if !ok {
			return fmt.Errorf("knowledge item %q not found for revision", event.EntityID)
		}
		item.CurrentRevision = event.Revision
		item.Title = event.Title
		item.Summary = event.Summary
		item.Tags = append([]string(nil), event.Tags...)
		item.UpdatedAt = event.Timestamp
		item.Revisions = append(item.Revisions, KnowledgeRevision{
			Number:    event.Revision,
			Title:     event.Title,
			Summary:   event.Summary,
			Body:      event.Body,
			Tags:      append([]string(nil), event.Tags...),
			Author:    event.Actor,
			CreatedAt: event.Timestamp,
		})
		item.Timeline = append(item.Timeline, event)
	case "run_recorded":
		item, ok := app.items[event.TargetID]
		if !ok {
			return fmt.Errorf("knowledge item %q not found for run", event.TargetID)
		}
		run := &RunRecord{
			ID:                event.EntityID,
			Kind:              event.Kind,
			ItemID:            event.TargetID,
			ItemKind:          item.Kind,
			Revision:          event.Revision,
			Actor:             event.Actor,
			Outcome:           event.Outcome,
			Notes:             event.Notes,
			Machine:           event.Machine,
			Location:          event.Location,
			ResponsibilityIDs: append([]string(nil), event.ResponsibilityIDs...),
			CreatedAt:         event.Timestamp,
			UpdatedAt:         event.Timestamp,
			Timeline:          []OperationalEvent{event},
		}
		app.runs[event.EntityID] = run
		for _, responsibilityID := range event.ResponsibilityIDs {
			if responsibility, ok := app.responsibilities[responsibilityID]; ok {
				responsibility.LinkedRunIDs = appendUnique(responsibility.LinkedRunIDs, event.EntityID)
				responsibility.UpdatedAt = event.Timestamp
			}
		}
	case "evidence_added":
		run, ok := app.runs[event.EntityID]
		if !ok {
			return fmt.Errorf("run %q not found for evidence", event.EntityID)
		}
		id := app.nextDerivedIDLocked("EVID")
		run.Evidence = append(run.Evidence, Evidence{
			ID:             id,
			Summary:        event.Summary,
			Facts:          cloneFacts(event.Facts),
			AttachmentName: event.AttachmentName,
			AttachmentPath: event.AttachmentPath,
			AttachmentSize: event.AttachmentSize,
			Actor:          event.Actor,
			CreatedAt:      event.Timestamp,
		})
		run.UpdatedAt = event.Timestamp
		run.Timeline = append(run.Timeline, event)
	case "approval_recorded":
		approval := Approval{
			ID:         event.EntityID,
			TargetType: event.TargetType,
			TargetID:   event.TargetID,
			Revision:   event.Revision,
			RunID:      event.RunID,
			Role:       event.Role,
			Decision:   event.Decision,
			Actor:      event.Actor,
			Notes:      event.Notes,
			CreatedAt:  event.Timestamp,
		}
		app.approvals[event.EntityID] = &approval
		switch event.TargetType {
		case "knowledge_item":
			if item, ok := app.items[event.TargetID]; ok {
				item.Approvals = append(item.Approvals, approval)
				item.UpdatedAt = event.Timestamp
				item.Timeline = append(item.Timeline, event)
			}
		case "run":
			if run, ok := app.runs[event.TargetID]; ok {
				run.Approvals = append(run.Approvals, approval)
				run.UpdatedAt = event.Timestamp
				run.Timeline = append(run.Timeline, event)
			}
		}
	case "link_added":
		link := &Link{
			ID:        event.EntityID,
			FromType:  event.FromType,
			FromID:    event.FromID,
			ToType:    event.ToType,
			ToID:      event.ToID,
			Relation:  event.Relation,
			Notes:     event.Notes,
			Actor:     event.Actor,
			CreatedAt: event.Timestamp,
		}
		app.links[event.EntityID] = link
		app.attachLinkLocked(*link, event)
	default:
		return fmt.Errorf("unknown event type %q", event.Type)
	}
	app.observeIDLocked(event.EntityID)
	return nil
}

func (app *App) attachLinkLocked(link Link, event OperationalEvent) {
	if responsibility, ok := app.responsibilities[link.FromID]; ok {
		responsibility.UpdatedAt = event.Timestamp
		responsibility.Timeline = append(responsibility.Timeline, event)
	}
	if responsibility, ok := app.responsibilities[link.ToID]; ok {
		responsibility.UpdatedAt = event.Timestamp
		responsibility.Timeline = append(responsibility.Timeline, event)
	}
	if item, ok := app.items[link.FromID]; ok {
		item.Links = append(item.Links, link)
		item.UpdatedAt = event.Timestamp
		item.Timeline = append(item.Timeline, event)
	}
	if item, ok := app.items[link.ToID]; ok {
		item.Links = append(item.Links, link)
		item.UpdatedAt = event.Timestamp
		item.Timeline = append(item.Timeline, event)
	}
	if run, ok := app.runs[link.FromID]; ok {
		run.Links = append(run.Links, link)
		run.UpdatedAt = event.Timestamp
		run.Timeline = append(run.Timeline, event)
	}
	if run, ok := app.runs[link.ToID]; ok {
		run.Links = append(run.Links, link)
		run.UpdatedAt = event.Timestamp
		run.Timeline = append(run.Timeline, event)
	}
}

func (app *App) validateResponsibilitiesLocked(ids []string) error {
	for _, id := range ids {
		if _, ok := app.responsibilities[id]; !ok {
			return fmt.Errorf("responsibility %q not found", id)
		}
	}
	return nil
}

func kindPrefix(kind string) string {
	switch kind {
	case KnowledgeKindProcedure:
		return "PROC"
	case KnowledgeKindTraining:
		return "TRAIN"
	case KnowledgeKindMaintenance:
		return "MAINT"
	default:
		return "ITEM"
	}
}

func itemKindForRun(runKind string) string {
	switch runKind {
	case RunKindProcedure:
		return KnowledgeKindProcedure
	case RunKindTraining:
		return KnowledgeKindTraining
	case RunKindMaintenance:
		return KnowledgeKindMaintenance
	default:
		return ""
	}
}

func (app *App) nextIDLocked(prefix string) string {
	app.nextNumbers[prefix]++
	return fmt.Sprintf("%s-%04d", prefix, app.nextNumbers[prefix])
}

func (app *App) nextDerivedIDLocked(prefix string) string {
	app.nextNumbers[prefix]++
	return fmt.Sprintf("%s-%04d", prefix, app.nextNumbers[prefix])
}

func (app *App) observeIDLocked(id string) {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) != 2 {
		return
	}
	prefix := parts[0]
	var number int
	fmt.Sscanf(parts[1], "%d", &number)
	if number > app.nextNumbers[prefix] {
		app.nextNumbers[prefix] = number
	}
}

func normalizeStrings(values []string) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func normalizeBody(body string) string {
	return strings.ReplaceAll(strings.TrimSpace(body), "\r\n", "\n")
}

func validateKnowledgeKind(kind string) error {
	switch kind {
	case KnowledgeKindProcedure, KnowledgeKindTraining, KnowledgeKindMaintenance:
		return nil
	default:
		return fmt.Errorf("unsupported knowledge kind %q", kind)
	}
}

func validateRunKind(kind string) error {
	switch kind {
	case RunKindProcedure, RunKindTraining, RunKindMaintenance:
		return nil
	default:
		return fmt.Errorf("unsupported run kind %q", kind)
	}
}

func validateDecision(decision string) error {
	switch decision {
	case DecisionApproved, DecisionRejected, DecisionNoted:
		return nil
	default:
		return fmt.Errorf("unsupported decision %q", decision)
	}
}

func normalizeFacts(facts map[string]string) map[string]string {
	if len(facts) == 0 {
		return map[string]string{}
	}
	out := map[string]string{}
	for key, value := range facts {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}
		out[key] = value
	}
	return out
}

func cloneFacts(in map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range in {
		out[key] = value
	}
	return out
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func cloneResponsibility(in *Responsibility) Responsibility {
	out := *in
	out.Tags = append([]string(nil), in.Tags...)
	out.LinkedItemIDs = append([]string(nil), in.LinkedItemIDs...)
	out.LinkedRunIDs = append([]string(nil), in.LinkedRunIDs...)
	out.LinkedRoleKeys = append([]string(nil), in.LinkedRoleKeys...)
	out.Timeline = append([]OperationalEvent(nil), in.Timeline...)
	return out
}

func cloneKnowledgeItem(in *KnowledgeItem) KnowledgeItem {
	out := *in
	out.Tags = append([]string(nil), in.Tags...)
	out.ResponsibilityIDs = append([]string(nil), in.ResponsibilityIDs...)
	out.Revisions = append([]KnowledgeRevision(nil), in.Revisions...)
	out.Approvals = append([]Approval(nil), in.Approvals...)
	out.Links = append([]Link(nil), in.Links...)
	out.Timeline = append([]OperationalEvent(nil), in.Timeline...)
	return out
}

func cloneRun(in *RunRecord) RunRecord {
	out := *in
	out.ResponsibilityIDs = append([]string(nil), in.ResponsibilityIDs...)
	out.Evidence = append([]Evidence(nil), in.Evidence...)
	out.Approvals = append([]Approval(nil), in.Approvals...)
	out.Links = append([]Link(nil), in.Links...)
	out.Timeline = append([]OperationalEvent(nil), in.Timeline...)
	return out
}
