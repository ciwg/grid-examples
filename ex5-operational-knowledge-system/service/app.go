package service

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

const defaultTeam = "OPS"

type App struct {
	dataRoot string
	store    *Store

	mu               sync.Mutex
	responsibilities map[string]*Responsibility
	places           map[string]*Place
	resources        map[string]*Resource
	items            map[string]*KnowledgeItem
	runs             map[string]*RunRecord
	links            map[string]*Link
	approvals        map[string]*Approval
	presence         map[string]map[string]*LivePresence
	nextSequence     uint64
	nextNumbers      map[string]int
}

func NewApp(dataRoot string) (*App, error) {
	store, events, knowledgeItemRecords, knowledgeApprovalRecords, knowledgeEvidenceRecords, knowledgeLinkRecords, knowledgeResponsibilityRecords, err := OpenStore(dataRoot)
	if err != nil {
		return nil, err
	}
	app := &App{
		dataRoot:         dataRoot,
		store:            store,
		responsibilities: map[string]*Responsibility{},
		places:           map[string]*Place{},
		resources:        map[string]*Resource{},
		items:            map[string]*KnowledgeItem{},
		runs:             map[string]*RunRecord{},
		links:            map[string]*Link{},
		approvals:        map[string]*Approval{},
		presence:         map[string]map[string]*LivePresence{},
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
	if err := verifySignedKnowledgeItemRecords(events, knowledgeItemRecords); err != nil {
		return nil, fmt.Errorf("verify knowledge-item envelopes: %w", err)
	}
	if err := verifySignedKnowledgeApprovalRecords(events, knowledgeApprovalRecords); err != nil {
		return nil, fmt.Errorf("verify knowledge-approval envelopes: %w", err)
	}
	if err := verifySignedKnowledgeEvidenceRecords(events, knowledgeEvidenceRecords); err != nil {
		return nil, fmt.Errorf("verify knowledge-evidence envelopes: %w", err)
	}
	if err := verifySignedKnowledgeLinkRecords(events, knowledgeLinkRecords); err != nil {
		return nil, fmt.Errorf("verify knowledge-link envelopes: %w", err)
	}
	if err := verifySignedKnowledgeResponsibilityRecords(events, knowledgeResponsibilityRecords); err != nil {
		return nil, fmt.Errorf("verify knowledge-responsibility envelopes: %w", err)
	}
	drafts, err := store.LoadDrafts()
	if err != nil {
		return nil, err
	}
	for itemID, draft := range drafts {
		item, ok := app.items[itemID]
		if !ok {
			continue
		}
		item.WorkingBody = draft.Body
		item.WorkingVersion = draft.Version
		item.WorkingUpdatedAt = draft.UpdatedAt
	}
	return app, nil
}

func (app *App) Meta() Meta {
	return Meta{
		DataRoot:                    app.dataRoot,
		KnowledgeKinds:              []string{KnowledgeKindProcedure, KnowledgeKindTraining, KnowledgeKindMaintenance, KnowledgeKindReceiving, KnowledgeKindInventory},
		RunKinds:                    []string{RunKindProcedure, RunKindTraining, RunKindMaintenance, RunKindReceiving, RunKindInventory},
		ApprovalDecisions:           []string{DecisionApproved, DecisionRejected, DecisionNoted},
		ItemStatuses:                []string{ItemStatusDraft, ItemStatusApproved, ItemStatusSuperseded},
		KnowledgeItemPCID:           protocols.KnowledgeItemProfile.CID.String(),
		KnowledgeApprovalPCID:       protocols.KnowledgeApprovalProfile.CID.String(),
		KnowledgeEvidencePCID:       protocols.KnowledgeEvidenceProfile.CID.String(),
		KnowledgeLinkPCID:           protocols.KnowledgeLinkProfile.CID.String(),
		KnowledgeResponsibilityPCID: protocols.KnowledgeResponsibilityProfile.CID.String(),
		PeerExchangeFormat:          peerExchangeBundleFormat,
		PeerExchangeFamilies: []string{
			"knowledge-item",
			"knowledge-approval",
			"knowledge-link",
			"knowledge-responsibility",
		},
		CASObjectsEnabled:         true,
		CASAttachmentBlobsEnabled: true,
	}
}

func (app *App) Dashboard() Dashboard {
	app.mu.Lock()
	defer app.mu.Unlock()
	out := Dashboard{
		Responsibilities: len(app.responsibilities),
		Places:           len(app.places),
		Resources:        len(app.resources),
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
		case KnowledgeKindReceiving:
			out.ReceivingItems++
		case KnowledgeKindInventory:
			out.InventoryItems++
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
		case RunKindReceiving:
			out.ReceivingRuns++
		case RunKindInventory:
			out.InventoryRuns++
		}
		out.Evidence += len(run.Evidence)
	}
	return out
}

// Intent: Give operators one grouped view of repeated receiving and count
// problems so they can spot hotspots by place and resource without rebuilding
// the pattern from individual runs. Source: DI-pogul
func (app *App) ProblemReview() ProblemReview {
	app.mu.Lock()
	defer app.mu.Unlock()

	placeGroups := map[string]*ProblemReviewGroup{}
	resourceGroups := map[string]*ProblemReviewGroup{}
	out := ProblemReview{}
	for _, run := range app.runs {
		highlights := problemHighlightsForRun(run)
		if len(highlights) == 0 {
			continue
		}
		out.ProblemRuns++
		if place, ok := app.places[run.PlaceID]; ok {
			group := placeGroups[place.ID]
			if group == nil {
				group = &ProblemReviewGroup{
					GroupType: "place",
					GroupID:   place.ID,
					Kind:      place.Kind,
					Name:      place.Name,
				}
				placeGroups[place.ID] = group
			}
			addProblemRunToGroup(group, run, highlights)
		}
		for _, resourceID := range run.ResourceIDs {
			resource, ok := app.resources[resourceID]
			if !ok {
				continue
			}
			group := resourceGroups[resource.ID]
			if group == nil {
				group = &ProblemReviewGroup{
					GroupType: "resource",
					GroupID:   resource.ID,
					Kind:      resource.Kind,
					Name:      resource.Name,
				}
				resourceGroups[resource.ID] = group
			}
			addProblemRunToGroup(group, run, highlights)
		}
	}
	out.PlaceGroups = flattenProblemGroups(placeGroups)
	out.ResourceGroups = flattenProblemGroups(resourceGroups)
	return out
}

func (app *App) ListPlaces() []Place {
	app.mu.Lock()
	defer app.mu.Unlock()
	out := make([]Place, 0, len(app.places))
	for _, place := range app.places {
		out = append(out, clonePlace(place))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (app *App) GetPlace(id string) (Place, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	place, ok := app.places[id]
	if !ok {
		return Place{}, fmt.Errorf("place %q not found", id)
	}
	return app.placeWithRelatedRunsLocked(place), nil
}

// Intent: Model physical context through generic hierarchical places instead
// of a warehouse-only vocabulary so ex5 can cover rooms, benches, racks,
// stations, bins, and similar operational spaces in one shared workflow model.
// Source: DI-foluk
func (app *App) CreatePlace(actor string, kind string, name string, summary string, parentID string, tags []string) (Place, error) {
	if strings.TrimSpace(actor) == "" {
		return Place{}, fmt.Errorf("actor is required")
	}
	if strings.TrimSpace(kind) == "" {
		return Place{}, fmt.Errorf("kind is required")
	}
	if strings.TrimSpace(name) == "" {
		return Place{}, fmt.Errorf("name is required")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	parentID = strings.TrimSpace(parentID)
	if parentID != "" {
		if _, ok := app.places[parentID]; !ok {
			return Place{}, fmt.Errorf("parent place %q not found", parentID)
		}
	}
	id := app.nextIDLocked("PLACE")
	event := OperationalEvent{
		EntityType: "place",
		EntityID:   id,
		Type:       "place_created",
		Actor:      actor,
		Kind:       strings.TrimSpace(kind),
		Name:       strings.TrimSpace(name),
		Summary:    strings.TrimSpace(summary),
		ParentID:   parentID,
		Tags:       normalizeStrings(tags),
	}
	if err := app.appendEventLocked(event); err != nil {
		return Place{}, err
	}
	return clonePlace(app.places[id]), nil
}

func (app *App) ListResources() []Resource {
	app.mu.Lock()
	defer app.mu.Unlock()
	out := make([]Resource, 0, len(app.resources))
	for _, resource := range app.resources {
		out = append(out, cloneResource(resource))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (app *App) GetResource(id string) (Resource, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	resource, ok := app.resources[id]
	if !ok {
		return Resource{}, fmt.Errorf("resource %q not found", id)
	}
	return app.resourceWithRelatedRunsLocked(resource), nil
}

// Intent: Keep tracked operational things generic so ex5 can represent
// machines, tools, parts, containers, and similar resources without splitting
// into separate domain-specific engines too early. Source: DI-foluk
func (app *App) CreateResource(actor string, kind string, name string, summary string, placeID string, tags []string) (Resource, error) {
	if strings.TrimSpace(actor) == "" {
		return Resource{}, fmt.Errorf("actor is required")
	}
	if strings.TrimSpace(kind) == "" {
		return Resource{}, fmt.Errorf("kind is required")
	}
	if strings.TrimSpace(name) == "" {
		return Resource{}, fmt.Errorf("name is required")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	placeID = strings.TrimSpace(placeID)
	if placeID != "" {
		if _, ok := app.places[placeID]; !ok {
			return Resource{}, fmt.Errorf("place %q not found", placeID)
		}
	}
	id := app.nextIDLocked("RES")
	event := OperationalEvent{
		EntityType: "resource",
		EntityID:   id,
		Type:       "resource_created",
		Actor:      actor,
		Kind:       strings.TrimSpace(kind),
		Name:       strings.TrimSpace(name),
		Summary:    strings.TrimSpace(summary),
		PlaceID:    placeID,
		Tags:       normalizeStrings(tags),
	}
	if err := app.appendEventLocked(event); err != nil {
		return Resource{}, err
	}
	return cloneResource(app.resources[id]), nil
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
	return app.responsibilityWithRelatedRunsLocked(record), nil
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
	return app.itemWithRelatedRunsLocked(item), nil
}

// Intent: Keep procedures, training, maintenance, receiving, and inventory
// work as hybrid knowledge items with structured metadata plus revisioned
// shared text, so operational records and collaborative knowledge can coexist
// in one tool. Source: DI-kovup; DI-vemok
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
		Status:            ItemStatusDraft,
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
		Status:     ItemStatusDraft,
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

// Intent: Make supersedence explicit in the durable item lifecycle so ex5 can
// preserve historical revisions without pretending replaced procedures or
// records are still current. Source: DI-zoruk
func (app *App) SupersedeKnowledgeItem(actor string, itemID string, notes string) (KnowledgeItem, error) {
	if strings.TrimSpace(actor) == "" {
		return KnowledgeItem{}, fmt.Errorf("actor is required")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	if _, ok := app.items[itemID]; !ok {
		return KnowledgeItem{}, fmt.Errorf("knowledge item %q not found", itemID)
	}
	event := OperationalEvent{
		EntityType: "knowledge_item",
		EntityID:   itemID,
		Type:       "knowledge_item_superseded",
		Actor:      actor,
		Status:     ItemStatusSuperseded,
		Notes:      strings.TrimSpace(notes),
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
// procedure, training, maintenance, receiving, and inventory execution can
// point back to the exact revision, evidence, and responsibilities involved.
// Source: DI-kovup; DI-zuvob; DI-vemok
func (app *App) RecordRun(actor string, kind string, itemID string, revision int, outcome string, notes string, machine string, location string, placeID string, resourceIDs []string, responsibilityIDs []string) (RunRecord, error) {
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
	if err := app.validatePlaceLocked(placeID); err != nil {
		return RunRecord{}, err
	}
	if err := app.validateResourcesLocked(resourceIDs); err != nil {
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
		PlaceID:           strings.TrimSpace(placeID),
		ResourceIDs:       normalizeStrings(resourceIDs),
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
	evidenceID := app.nextIDLocked("EVID")
	event := OperationalEvent{
		EntityType: "run",
		EntityID:   runID,
		Type:       "evidence_added",
		EvidenceID: evidenceID,
		Actor:      actor,
		Summary:    strings.TrimSpace(summary),
		Facts:      normalizeFacts(facts),
	}
	if len(attachmentBody) > 0 {
		path, cid, size, err := app.store.SaveAttachment(runID, attachmentName, attachmentBody)
		if err != nil {
			return RunRecord{}, err
		}
		event.AttachmentName = filepath.Base(attachmentName)
		event.AttachmentPath = path
		event.AttachmentCID = cid
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
	var item *KnowledgeItem
	switch targetType {
	case "knowledge_item":
		var ok bool
		item, ok = app.items[targetID]
		if !ok {
			return fmt.Errorf("knowledge item %q not found", targetID)
		}
		if revision <= 0 || revision > len(item.Revisions) {
			return fmt.Errorf("revision %d not found for %s", revision, targetID)
		}
		// Intent: Only let a knowledge-item approval move lifecycle state when it
		// approves the currently active revision, so an old approval cannot mark a
		// newer draft as approved by accident. Source: DI-dazim
		if decision == DecisionApproved && revision != item.CurrentRevision {
			return fmt.Errorf("cannot approve stale revision %d for current revision %d", revision, item.CurrentRevision)
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
	if err := app.appendEventLocked(event); err != nil {
		return err
	}
	if targetType == "knowledge_item" && decision == DecisionApproved {
		statusEvent := OperationalEvent{
			EntityType: "knowledge_item",
			EntityID:   targetID,
			Type:       "knowledge_item_status_changed",
			Actor:      actor,
			Status:     ItemStatusApproved,
			Notes:      strings.TrimSpace(notes),
		}
		return app.appendEventLocked(statusEvent)
	}
	return nil
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
	fromType = strings.TrimSpace(strings.ToLower(fromType))
	fromID = strings.TrimSpace(fromID)
	toType = strings.TrimSpace(strings.ToLower(toType))
	toID = strings.TrimSpace(toID)
	// Intent: Keep the typed-link graph trustworthy by rejecting dangling or
	// mistyped endpoints before they enter the append-only history. Source:
	// DI-luzaf
	if err := app.validateLinkEndpointLocked(fromType, fromID); err != nil {
		return fmt.Errorf("from endpoint invalid: %w", err)
	}
	if err := app.validateLinkEndpointLocked(toType, toID); err != nil {
		return fmt.Errorf("to endpoint invalid: %w", err)
	}
	event := OperationalEvent{
		EntityType: "link",
		EntityID:   app.nextIDLocked("LINK"),
		Type:       "link_added",
		Actor:      actor,
		FromType:   fromType,
		FromID:     fromID,
		ToType:     toType,
		ToID:       toID,
		Relation:   strings.TrimSpace(relation),
		Notes:      strings.TrimSpace(notes),
	}
	return app.appendEventLocked(event)
}

func (app *App) Search(query string) map[string]any {
	return app.SearchWithOptions(SearchOptions{Query: query})
}

// Intent: Let operators narrow the operational graph by structured context
// such as kind, status, outcome, place, resource, responsibility, and
// problem-only review state without forcing them to rely on one free-text
// query string. Source: DI-honus; DI-vafuk; DI-vemur
func (app *App) SearchWithOptions(options SearchOptions) map[string]any {
	app.mu.Lock()
	defer app.mu.Unlock()
	options = normalizeSearchOptions(options)
	problemContext := app.problemSearchContextLocked(options)
	places := []Place{}
	resources := []Resource{}
	resp := []Responsibility{}
	items := []KnowledgeItem{}
	runs := []RunRecord{}
	for _, record := range app.places {
		if matchesPlaceSearch(record, options) && problemContext.allowsPlace(record.ID) {
			places = append(places, clonePlace(record))
		}
	}
	for _, record := range app.resources {
		if matchesResourceSearch(record, options) && problemContext.allowsResource(record.ID) {
			resources = append(resources, cloneResource(record))
		}
	}
	for _, record := range app.responsibilities {
		if matchesResponsibilitySearch(record, options) && problemContext.allowsResponsibility(record.ID) {
			resp = append(resp, cloneResponsibility(record))
		}
	}
	for _, record := range app.items {
		if matchesItemSearch(record, options) && problemContext.allowsItem(record.ID) {
			items = append(items, cloneKnowledgeItem(record))
		}
	}
	for _, record := range app.runs {
		searchBlob := runSearchBlob(record, app.places, app.resources)
		if matchesRunSearch(record, searchBlob, options) {
			runs = append(runs, cloneRun(record))
		}
	}
	sort.Slice(places, func(i, j int) bool { return places[i].ID < places[j].ID })
	sort.Slice(resources, func(i, j int) bool { return resources[i].ID < resources[j].ID })
	sort.Slice(resp, func(i, j int) bool { return resp[i].ID < resp[j].ID })
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	sort.Slice(runs, func(i, j int) bool { return runs[i].ID < runs[j].ID })
	return map[string]any{
		"filters": map[string]any{
			"query":             options.Query,
			"kind":              options.Kind,
			"status":            options.Status,
			"outcome":           options.Outcome,
			"place_id":          options.PlaceID,
			"resource_id":       options.ResourceID,
			"responsibility_id": options.ResponsibilityID,
			"problem":           options.Problem,
		},
		"places":           places,
		"resources":        resources,
		"responsibilities": resp,
		"items":            items,
		"runs":             runs,
	}
}

func normalizeSearchOptions(options SearchOptions) SearchOptions {
	options.Query = strings.ToLower(strings.TrimSpace(options.Query))
	options.Kind = strings.ToLower(strings.TrimSpace(options.Kind))
	options.Status = strings.ToLower(strings.TrimSpace(options.Status))
	options.Outcome = strings.ToLower(strings.TrimSpace(options.Outcome))
	options.PlaceID = strings.TrimSpace(options.PlaceID)
	options.ResourceID = strings.TrimSpace(options.ResourceID)
	options.ResponsibilityID = strings.TrimSpace(options.ResponsibilityID)
	return options
}

type problemSearchContext struct {
	enabled           bool
	placeIDs          map[string]bool
	resourceIDs       map[string]bool
	responsibilityIDs map[string]bool
	itemIDs           map[string]bool
}

func (context problemSearchContext) allowsPlace(id string) bool {
	if !context.enabled {
		return true
	}
	return context.placeIDs[id]
}

func (context problemSearchContext) allowsResource(id string) bool {
	if !context.enabled {
		return true
	}
	return context.resourceIDs[id]
}

func (context problemSearchContext) allowsResponsibility(id string) bool {
	if !context.enabled {
		return true
	}
	return context.responsibilityIDs[id]
}

func (context problemSearchContext) allowsItem(id string) bool {
	if !context.enabled {
		return true
	}
	return context.itemIDs[id]
}

// Intent: Keep `problem=true` search honest across all result groups by only
// returning places, resources, responsibilities, and items that are actually
// connected to problem runs after the same structured run filters are applied.
// Source: DI-ralek
func (app *App) problemSearchContextLocked(options SearchOptions) problemSearchContext {
	if !options.Problem {
		return problemSearchContext{}
	}
	context := problemSearchContext{
		enabled:           true,
		placeIDs:          map[string]bool{},
		resourceIDs:       map[string]bool{},
		responsibilityIDs: map[string]bool{},
		itemIDs:           map[string]bool{},
	}
	for _, run := range app.runs {
		if !problemRunMatchesOptions(run, options) {
			continue
		}
		if run.PlaceID != "" {
			context.placeIDs[run.PlaceID] = true
		}
		context.itemIDs[run.ItemID] = true
		for _, resourceID := range run.ResourceIDs {
			context.resourceIDs[resourceID] = true
		}
		for _, responsibilityID := range run.ResponsibilityIDs {
			context.responsibilityIDs[responsibilityID] = true
		}
	}
	return context
}

func matchesPlaceSearch(record *Place, options SearchOptions) bool {
	if options.Kind != "" && strings.ToLower(record.Kind) != options.Kind {
		return false
	}
	if options.PlaceID != "" && record.ID != options.PlaceID {
		return false
	}
	return matchesQuery(record.ID+" "+record.Name+" "+record.Summary+" "+record.Kind, options.Query)
}

func matchesResourceSearch(record *Resource, options SearchOptions) bool {
	if options.Kind != "" && strings.ToLower(record.Kind) != options.Kind {
		return false
	}
	if options.PlaceID != "" && record.PlaceID != options.PlaceID {
		return false
	}
	if options.ResourceID != "" && record.ID != options.ResourceID {
		return false
	}
	return matchesQuery(record.ID+" "+record.PlaceID+" "+record.Name+" "+record.Summary+" "+record.Kind, options.Query)
}

func matchesResponsibilitySearch(record *Responsibility, options SearchOptions) bool {
	if options.ResponsibilityID != "" && record.ID != options.ResponsibilityID {
		return false
	}
	return matchesQuery(record.ID+" "+record.Title+" "+record.Summary+" "+strings.Join(record.LinkedRoleKeys, " "), options.Query)
}

func matchesItemSearch(record *KnowledgeItem, options SearchOptions) bool {
	if options.Kind != "" && strings.ToLower(record.Kind) != options.Kind {
		return false
	}
	if options.Status != "" && strings.ToLower(record.Status) != options.Status {
		return false
	}
	if options.ResponsibilityID != "" && !containsValue(record.ResponsibilityIDs, options.ResponsibilityID) {
		return false
	}
	return matchesQuery(record.ID+" "+record.Title+" "+record.Summary+" "+record.WorkingBody+" "+record.Status+" "+strings.Join(record.ResponsibilityIDs, " "), options.Query)
}

func matchesRunSearch(record *RunRecord, searchBlob string, options SearchOptions) bool {
	if !matchesRunContextFilters(record, options) {
		return false
	}
	if options.Problem && len(problemHighlightsForRun(record)) == 0 {
		return false
	}
	return matchesQuery(searchBlob, options.Query)
}

func matchesRunContextFilters(record *RunRecord, options SearchOptions) bool {
	if options.Kind != "" && strings.ToLower(record.Kind) != options.Kind {
		return false
	}
	if options.Outcome != "" && strings.ToLower(record.Outcome) != options.Outcome {
		return false
	}
	if options.PlaceID != "" && record.PlaceID != options.PlaceID {
		return false
	}
	if options.ResourceID != "" && !containsValue(record.ResourceIDs, options.ResourceID) {
		return false
	}
	if options.ResponsibilityID != "" && !containsValue(record.ResponsibilityIDs, options.ResponsibilityID) {
		return false
	}
	return true
}

func problemRunMatchesOptions(run *RunRecord, options SearchOptions) bool {
	if len(problemHighlightsForRun(run)) == 0 {
		return false
	}
	return matchesRunContextFilters(run, options)
}

// Intent: Keep ex5 free-text search aligned with the operational-memory story
// by indexing evidence summaries/facts and approval review notes alongside the
// basic run fields. Source: DI-farun
func runSearchBlob(record *RunRecord, places map[string]*Place, resources map[string]*Resource) string {
	var parts []string
	parts = append(parts, record.ID, record.ItemID, record.PlaceID, record.Outcome, record.Notes, record.Machine, record.Location, strings.Join(record.ResourceIDs, " "), strings.Join(record.ResponsibilityIDs, " "))
	if place, ok := places[record.PlaceID]; ok {
		parts = append(parts, place.ID, place.Name, place.Summary, place.Kind)
	}
	for _, resourceID := range record.ResourceIDs {
		if resource, ok := resources[resourceID]; ok {
			parts = append(parts, resource.ID, resource.Name, resource.Summary, resource.Kind)
		}
	}
	for _, evidence := range record.Evidence {
		parts = append(parts, evidence.Summary, evidence.Actor, evidence.AttachmentName)
		for key, value := range evidence.Facts {
			parts = append(parts, key, value)
		}
	}
	for _, approval := range record.Approvals {
		parts = append(parts, approval.Actor, approval.Role, approval.Decision, approval.Notes)
	}
	return strings.Join(parts, " ")
}

func matchesQuery(value string, query string) bool {
	if query == "" {
		return true
	}
	return strings.Contains(strings.ToLower(value), query)
}

func (app *App) LiveItemState(itemID string) (LiveItemState, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	item, ok := app.items[itemID]
	if !ok {
		return LiveItemState{}, fmt.Errorf("knowledge item %q not found", itemID)
	}
	app.cleanupPresenceLocked(itemID)
	return app.liveStateLocked(item), nil
}

// Intent: Keep collaborative item drafting available in the browser while
// preserving a separate durable revision workflow for approvals and historical
// reproduction. Source: DI-lusov; DI-zoruk
func (app *App) UpdateLiveItem(itemID string, participantID string, displayName string, color string, cursor int, head int, typing bool, baseVersion int, updateBody bool, body string) (LiveItemState, bool, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	item, ok := app.items[itemID]
	if !ok {
		return LiveItemState{}, false, fmt.Errorf("knowledge item %q not found", itemID)
	}
	if strings.TrimSpace(participantID) == "" {
		return LiveItemState{}, false, fmt.Errorf("participant_id is required")
	}
	app.ensurePresenceLocked(itemID)
	now := time.Now().Format(time.RFC3339)
	app.presence[itemID][participantID] = &LivePresence{
		ParticipantID: participantID,
		DisplayName:   fallbackDisplayName(displayName, participantID),
		Color:         fallbackColor(color),
		Cursor:        maxInt(0, cursor),
		Head:          maxInt(0, head),
		Typing:        typing,
		LastSeenAt:    now,
	}
	body = normalizeBody(body)
	// Intent: Distinguish presence-only posts from deliberate body writes so the
	// shared draft can legitimately become empty without breaking Neovim/browser
	// presence heartbeats. Source: DI-dazim
	if updateBody && body != item.WorkingBody {
		if baseVersion != item.WorkingVersion {
			app.cleanupPresenceLocked(itemID)
			return app.liveStateLocked(item), true, nil
		}
		item.WorkingBody = body
		item.WorkingVersion++
		item.WorkingUpdatedAt = now
		if err := app.store.SaveDraft(itemID, PersistedDraft{
			Body:      item.WorkingBody,
			Version:   item.WorkingVersion,
			UpdatedAt: item.WorkingUpdatedAt,
		}); err != nil {
			return LiveItemState{}, false, err
		}
	}
	app.cleanupPresenceLocked(itemID)
	return app.liveStateLocked(item), false, nil
}

func (app *App) appendEventLocked(event OperationalEvent) error {
	app.nextSequence++
	event.Sequence = app.nextSequence
	event.Timestamp = time.Now().Format(time.RFC3339)
	itemRecord, hasItemRecord, err := buildSignedKnowledgeItemRecord(app.store.identity, event)
	if err != nil {
		app.nextSequence--
		return err
	}
	approvalRecord, hasApprovalRecord, err := buildSignedKnowledgeApprovalRecord(app.store.identity, event)
	if err != nil {
		app.nextSequence--
		return err
	}
	evidenceRecord, hasEvidenceRecord, err := buildSignedKnowledgeEvidenceRecord(app.store.identity, event)
	if err != nil {
		app.nextSequence--
		return err
	}
	linkRecord, hasLinkRecord, err := buildSignedKnowledgeLinkRecord(app.store.identity, event)
	if err != nil {
		app.nextSequence--
		return err
	}
	responsibilityRecord, hasResponsibilityRecord, err := buildSignedKnowledgeResponsibilityRecord(app.store.identity, event)
	if err != nil {
		app.nextSequence--
		return err
	}
	if hasItemRecord {
		// Intent: Persist the first PromiseGrid-native knowledge-item artifact
		// before the compatibility event log so ex5 never claims a new signed item
		// event that the runtime failed to materialize. Source: DI-mibor
		if err := app.store.AppendSignedKnowledgeItemRecord(itemRecord); err != nil {
			app.nextSequence--
			return err
		}
	}
	if hasApprovalRecord {
		// Intent: Persist the second PromiseGrid-native knowledge-approval
		// artifact before the compatibility event log so ex5 never claims a new
		// signed approval that the runtime failed to materialize. Source: DI-vosul
		if err := app.store.AppendSignedKnowledgeApprovalRecord(approvalRecord); err != nil {
			app.nextSequence--
			return err
		}
	}
	if hasEvidenceRecord {
		// Intent: Persist the third PromiseGrid-native knowledge-evidence
		// artifact before the compatibility event log so ex5 never claims a new
		// signed evidence record that the runtime failed to materialize. Source:
		// DI-kavup; DI-ribof
		if err := app.store.AppendSignedKnowledgeEvidenceRecord(evidenceRecord); err != nil {
			app.nextSequence--
			return err
		}
	}
	if hasLinkRecord {
		// Intent: Persist the fourth PromiseGrid-native knowledge-link artifact
		// before the compatibility event log so ex5 never claims a new signed
		// link record that the runtime failed to materialize. Source: DI-votek
		if err := app.store.AppendSignedKnowledgeLinkRecord(linkRecord); err != nil {
			app.nextSequence--
			return err
		}
	}
	if hasResponsibilityRecord {
		// Intent: Persist the fifth PromiseGrid-native
		// knowledge-responsibility artifact before the compatibility event log so
		// ex5 never claims a new signed responsibility record that the runtime
		// failed to materialize. Source: DI-sarib
		if err := app.store.AppendSignedKnowledgeResponsibilityRecord(responsibilityRecord); err != nil {
			app.nextSequence--
			return err
		}
	}
	if err := app.store.AppendEvent(event); err != nil {
		app.nextSequence--
		return err
	}
	return app.applyEventLocked(event)
}

func (app *App) applyEventLocked(event OperationalEvent) error {
	switch event.Type {
	case "place_created":
		place := &Place{
			ID:            event.EntityID,
			Kind:          event.Kind,
			Name:          event.Name,
			Summary:       event.Summary,
			ParentID:      event.ParentID,
			Tags:          append([]string(nil), event.Tags...),
			CreatedAt:     event.Timestamp,
			UpdatedAt:     event.Timestamp,
			Timeline:      []OperationalEvent{event},
			ChildPlaceIDs: []string{},
			ResourceIDs:   []string{},
		}
		app.places[event.EntityID] = place
		if parent, ok := app.places[event.ParentID]; ok {
			parent.ChildPlaceIDs = appendUnique(parent.ChildPlaceIDs, event.EntityID)
			parent.UpdatedAt = event.Timestamp
		}
	case "resource_created":
		resource := &Resource{
			ID:        event.EntityID,
			Kind:      event.Kind,
			Name:      event.Name,
			Summary:   event.Summary,
			PlaceID:   event.PlaceID,
			Tags:      append([]string(nil), event.Tags...),
			CreatedAt: event.Timestamp,
			UpdatedAt: event.Timestamp,
			Timeline:  []OperationalEvent{event},
		}
		app.resources[event.EntityID] = resource
		if place, ok := app.places[event.PlaceID]; ok {
			place.ResourceIDs = appendUnique(place.ResourceIDs, event.EntityID)
			place.UpdatedAt = event.Timestamp
		}
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
			Status:            ItemStatusDraft,
			Title:             event.Title,
			Summary:           event.Summary,
			Tags:              append([]string(nil), event.Tags...),
			ResponsibilityIDs: append([]string(nil), event.ResponsibilityIDs...),
			CreatedAt:         event.Timestamp,
			UpdatedAt:         event.Timestamp,
			WorkingBody:       event.Body,
			WorkingVersion:    1,
			WorkingUpdatedAt:  event.Timestamp,
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
		item.Status = ItemStatusDraft
		item.Title = event.Title
		item.Summary = event.Summary
		item.Tags = append([]string(nil), event.Tags...)
		item.UpdatedAt = event.Timestamp
		item.WorkingBody = event.Body
		item.WorkingVersion++
		item.WorkingUpdatedAt = event.Timestamp
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
			PlaceID:           event.PlaceID,
			ResourceIDs:       append([]string(nil), event.ResourceIDs...),
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
		id := event.EvidenceID
		if strings.TrimSpace(id) == "" {
			id = app.nextDerivedIDLocked("EVID")
		}
		run.Evidence = append(run.Evidence, Evidence{
			ID:             id,
			Summary:        event.Summary,
			Facts:          cloneFacts(event.Facts),
			AttachmentName: event.AttachmentName,
			AttachmentPath: event.AttachmentPath,
			AttachmentCID:  event.AttachmentCID,
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
	case "knowledge_item_status_changed", "knowledge_item_superseded":
		item, ok := app.items[event.EntityID]
		if !ok {
			return fmt.Errorf("knowledge item %q not found for status change", event.EntityID)
		}
		item.Status = event.Status
		item.UpdatedAt = event.Timestamp
		item.Timeline = append(item.Timeline, event)
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
	app.attachLinkEndpointLocked(link.FromType, link.FromID, link, event)
	app.attachLinkEndpointLocked(link.ToType, link.ToID, link, event)
}

func (app *App) validateResponsibilitiesLocked(ids []string) error {
	for _, id := range ids {
		if _, ok := app.responsibilities[id]; !ok {
			return fmt.Errorf("responsibility %q not found", id)
		}
	}
	return nil
}

func (app *App) validatePlaceLocked(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	if _, ok := app.places[id]; !ok {
		return fmt.Errorf("place %q not found", id)
	}
	return nil
}

func (app *App) validateLinkEndpointLocked(entityType string, entityID string) error {
	switch entityType {
	case "place":
		if _, ok := app.places[entityID]; !ok {
			return fmt.Errorf("place %q not found", entityID)
		}
	case "resource":
		if _, ok := app.resources[entityID]; !ok {
			return fmt.Errorf("resource %q not found", entityID)
		}
	case "responsibility":
		if _, ok := app.responsibilities[entityID]; !ok {
			return fmt.Errorf("responsibility %q not found", entityID)
		}
	case "knowledge_item":
		if _, ok := app.items[entityID]; !ok {
			return fmt.Errorf("knowledge item %q not found", entityID)
		}
	case "run":
		if _, ok := app.runs[entityID]; !ok {
			return fmt.Errorf("run %q not found", entityID)
		}
	default:
		return fmt.Errorf("unsupported link endpoint type %q", entityType)
	}
	return nil
}

func (app *App) attachLinkEndpointLocked(entityType string, entityID string, link Link, event OperationalEvent) {
	switch entityType {
	case "responsibility":
		if responsibility, ok := app.responsibilities[entityID]; ok {
			responsibility.Links = append(responsibility.Links, link)
			responsibility.UpdatedAt = event.Timestamp
			responsibility.Timeline = append(responsibility.Timeline, event)
		}
	case "place":
		if place, ok := app.places[entityID]; ok {
			place.Links = append(place.Links, link)
			place.UpdatedAt = event.Timestamp
			place.Timeline = append(place.Timeline, event)
		}
	case "resource":
		if resource, ok := app.resources[entityID]; ok {
			resource.Links = append(resource.Links, link)
			resource.UpdatedAt = event.Timestamp
			resource.Timeline = append(resource.Timeline, event)
		}
	case "knowledge_item":
		if item, ok := app.items[entityID]; ok {
			item.Links = append(item.Links, link)
			item.UpdatedAt = event.Timestamp
			item.Timeline = append(item.Timeline, event)
		}
	case "run":
		if run, ok := app.runs[entityID]; ok {
			run.Links = append(run.Links, link)
			run.UpdatedAt = event.Timestamp
			run.Timeline = append(run.Timeline, event)
		}
	}
}

func (app *App) validateResourcesLocked(ids []string) error {
	for _, id := range ids {
		if _, ok := app.resources[id]; !ok {
			return fmt.Errorf("resource %q not found", id)
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
	case KnowledgeKindReceiving:
		return "RECV"
	case KnowledgeKindInventory:
		return "INV"
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
	case RunKindReceiving:
		return KnowledgeKindReceiving
	case RunKindInventory:
		return KnowledgeKindInventory
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
	if _, err := fmt.Sscanf(parts[1], "%d", &number); err != nil {
		return
	}
	if number > app.nextNumbers[prefix] {
		app.nextNumbers[prefix] = number
	}
}

func (app *App) ensurePresenceLocked(itemID string) {
	if app.presence[itemID] == nil {
		app.presence[itemID] = map[string]*LivePresence{}
	}
}

func (app *App) cleanupPresenceLocked(itemID string) {
	peers := app.presence[itemID]
	if len(peers) == 0 {
		return
	}
	cutoff := time.Now().Add(-15 * time.Second)
	for participantID, peer := range peers {
		seenAt, err := time.Parse(time.RFC3339, peer.LastSeenAt)
		if err != nil || seenAt.Before(cutoff) {
			delete(peers, participantID)
		}
	}
}

func (app *App) liveStateLocked(item *KnowledgeItem) LiveItemState {
	participants := make([]LivePresence, 0, len(app.presence[item.ID]))
	for _, peer := range app.presence[item.ID] {
		participants = append(participants, *peer)
	}
	sort.Slice(participants, func(i, j int) bool { return participants[i].ParticipantID < participants[j].ParticipantID })
	return LiveItemState{
		ItemID:          item.ID,
		Title:           item.Title,
		Status:          item.Status,
		Body:            item.WorkingBody,
		Version:         item.WorkingVersion,
		CurrentRevision: item.CurrentRevision,
		Participants:    participants,
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
	return strings.ReplaceAll(strings.ReplaceAll(body, "\r\n", "\n"), "\r", "\n")
}

func validateKnowledgeKind(kind string) error {
	switch kind {
	case KnowledgeKindProcedure, KnowledgeKindTraining, KnowledgeKindMaintenance, KnowledgeKindReceiving, KnowledgeKindInventory:
		return nil
	default:
		return fmt.Errorf("unsupported knowledge kind %q", kind)
	}
}

func validateRunKind(kind string) error {
	switch kind {
	case RunKindProcedure, RunKindTraining, RunKindMaintenance, RunKindReceiving, RunKindInventory:
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

func fallbackDisplayName(name string, participantID string) string {
	name = strings.TrimSpace(name)
	if name != "" {
		return name
	}
	return participantID
}

func fallbackColor(color string) string {
	color = strings.TrimSpace(color)
	if color == "" {
		return "#1d6fd6"
	}
	return color
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
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

func flattenProblemGroups(groups map[string]*ProblemReviewGroup) []ProblemReviewGroup {
	out := make([]ProblemReviewGroup, 0, len(groups))
	for _, group := range groups {
		sort.Slice(group.Runs, func(i, j int) bool { return group.Runs[i].ID < group.Runs[j].ID })
		sort.Strings(group.HighlightExamples)
		out = append(out, *group)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ProblemCount == out[j].ProblemCount {
			return out[i].GroupID < out[j].GroupID
		}
		return out[i].ProblemCount > out[j].ProblemCount
	})
	return out
}

func addProblemRunToGroup(group *ProblemReviewGroup, run *RunRecord, highlights []string) {
	group.ProblemCount++
	switch run.Kind {
	case RunKindReceiving:
		group.ReceivingProblems++
	case RunKindInventory:
		group.InventoryProblems++
	}
	group.Runs = append(group.Runs, cloneRun(run))
	for _, highlight := range highlights {
		group.HighlightExamples = appendUnique(group.HighlightExamples, highlight)
	}
	if len(group.HighlightExamples) > 6 {
		group.HighlightExamples = group.HighlightExamples[:6]
	}
}

func problemHighlightsForRun(run *RunRecord) []string {
	if run.Kind != RunKindReceiving && run.Kind != RunKindInventory {
		return nil
	}
	highlights := []string{}
	if problemOutcome(run.Outcome) {
		highlights = append(highlights, "outcome: "+run.Outcome)
	}
	for _, evidence := range run.Evidence {
		highlights = append(highlights, problemHighlightsFromFacts(evidence.Facts)...)
	}
	if len(highlights) == 0 {
		return nil
	}
	sort.Strings(highlights)
	return normalizeStrings(highlights)
}

func problemOutcome(outcome string) bool {
	outcome = strings.ToLower(strings.TrimSpace(outcome))
	switch outcome {
	case "", "completed", "accepted", "approved", "ok", "passed", "done":
		return false
	default:
		return true
	}
}

func problemHighlightsFromFacts(facts map[string]string) []string {
	highlights := []string{}
	expectedCount := strings.TrimSpace(facts["expected_count"])
	actualCount := strings.TrimSpace(facts["actual_count"])
	if expectedCount != "" && actualCount != "" && expectedCount != actualCount {
		highlights = append(highlights, "count mismatch: "+expectedCount+" -> "+actualCount)
	}
	keys := make([]string, 0, len(facts))
	for key := range facts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := strings.TrimSpace(facts[key])
		lowerKey := strings.ToLower(strings.TrimSpace(key))
		switch lowerKey {
		case "variance", "discrepancy", "delta":
			if isNonZeroFact(value) {
				highlights = append(highlights, lowerKey+": "+value)
			}
		case "condition":
			if isProblemCondition(value) {
				highlights = append(highlights, lowerKey+": "+value)
			}
		}
	}
	return highlights
}

func isNonZeroFact(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "", "0", "0.0", "+0", "-0", "none", "no", "false":
		return false
	default:
		return true
	}
}

func isProblemCondition(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "", "ok", "good", "acceptable", "intact", "clean", "none":
		return false
	default:
		return true
	}
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func containsValue(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func cloneResponsibility(in *Responsibility) Responsibility {
	out := *in
	out.Tags = append([]string(nil), in.Tags...)
	out.LinkedItemIDs = append([]string(nil), in.LinkedItemIDs...)
	out.LinkedRunIDs = append([]string(nil), in.LinkedRunIDs...)
	out.RelatedRuns = append([]RunRecord(nil), in.RelatedRuns...)
	out.LinkedRoleKeys = append([]string(nil), in.LinkedRoleKeys...)
	out.Links = append([]Link(nil), in.Links...)
	out.Timeline = append([]OperationalEvent(nil), in.Timeline...)
	return out
}

func clonePlace(in *Place) Place {
	out := *in
	out.Tags = append([]string(nil), in.Tags...)
	out.ChildPlaceIDs = append([]string(nil), in.ChildPlaceIDs...)
	out.ResourceIDs = append([]string(nil), in.ResourceIDs...)
	out.RelatedRuns = append([]RunRecord(nil), in.RelatedRuns...)
	out.Links = append([]Link(nil), in.Links...)
	out.Timeline = append([]OperationalEvent(nil), in.Timeline...)
	return out
}

func cloneResource(in *Resource) Resource {
	out := *in
	out.Tags = append([]string(nil), in.Tags...)
	out.RelatedRuns = append([]RunRecord(nil), in.RelatedRuns...)
	out.Links = append([]Link(nil), in.Links...)
	out.Timeline = append([]OperationalEvent(nil), in.Timeline...)
	return out
}

func cloneKnowledgeItem(in *KnowledgeItem) KnowledgeItem {
	out := *in
	out.Tags = append([]string(nil), in.Tags...)
	out.ResponsibilityIDs = append([]string(nil), in.ResponsibilityIDs...)
	out.Revisions = append([]KnowledgeRevision(nil), in.Revisions...)
	out.RelatedRuns = append([]RunRecord(nil), in.RelatedRuns...)
	out.Approvals = append([]Approval(nil), in.Approvals...)
	out.Links = append([]Link(nil), in.Links...)
	out.Timeline = append([]OperationalEvent(nil), in.Timeline...)
	return out
}

func (app *App) itemWithRelatedRunsLocked(item *KnowledgeItem) KnowledgeItem {
	out := cloneKnowledgeItem(item)
	out.RelatedRuns = []RunRecord{}
	for _, run := range app.runs {
		if run.ItemID != item.ID {
			continue
		}
		out.RelatedRuns = append(out.RelatedRuns, cloneRun(run))
	}
	sort.Slice(out.RelatedRuns, func(i, j int) bool {
		return out.RelatedRuns[i].ID < out.RelatedRuns[j].ID
	})
	return out
}

func cloneRun(in *RunRecord) RunRecord {
	out := *in
	out.ResourceIDs = append([]string(nil), in.ResourceIDs...)
	out.ResponsibilityIDs = append([]string(nil), in.ResponsibilityIDs...)
	out.Evidence = append([]Evidence(nil), in.Evidence...)
	out.Approvals = append([]Approval(nil), in.Approvals...)
	out.Links = append([]Link(nil), in.Links...)
	out.Timeline = append([]OperationalEvent(nil), in.Timeline...)
	return out
}

// Intent: Let operators inspect operational history from contextual anchors
// like places, resources, and responsibilities instead of forcing all history
// lookup to start from the knowledge item alone. Source: DI-julos
func (app *App) placeWithRelatedRunsLocked(place *Place) Place {
	out := clonePlace(place)
	out.RelatedRuns = []RunRecord{}
	for _, run := range app.runs {
		if run.PlaceID != place.ID {
			continue
		}
		out.RelatedRuns = append(out.RelatedRuns, cloneRun(run))
	}
	sort.Slice(out.RelatedRuns, func(i, j int) bool { return out.RelatedRuns[i].ID < out.RelatedRuns[j].ID })
	return out
}

func (app *App) resourceWithRelatedRunsLocked(resource *Resource) Resource {
	out := cloneResource(resource)
	out.RelatedRuns = []RunRecord{}
	for _, run := range app.runs {
		if !containsValue(run.ResourceIDs, resource.ID) {
			continue
		}
		out.RelatedRuns = append(out.RelatedRuns, cloneRun(run))
	}
	sort.Slice(out.RelatedRuns, func(i, j int) bool { return out.RelatedRuns[i].ID < out.RelatedRuns[j].ID })
	return out
}

func (app *App) responsibilityWithRelatedRunsLocked(record *Responsibility) Responsibility {
	out := cloneResponsibility(record)
	out.RelatedRuns = []RunRecord{}
	for _, run := range app.runs {
		if !containsValue(run.ResponsibilityIDs, record.ID) {
			continue
		}
		out.RelatedRuns = append(out.RelatedRuns, cloneRun(run))
	}
	sort.Slice(out.RelatedRuns, func(i, j int) bool { return out.RelatedRuns[i].ID < out.RelatedRuns[j].ID })
	return out
}
