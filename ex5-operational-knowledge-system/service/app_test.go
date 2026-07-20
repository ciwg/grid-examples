package service

import (
	"path/filepath"
	"testing"
)

func TestAppCreatesResponsibilitiesItemsRunsAndEvidence(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	responsibility, err := app.CreateResponsibility("alice", "Line lead", "Owns line startup", []string{"author", "approver"}, []string{"ops"})
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	if responsibility.ID == "" {
		t.Fatalf("expected responsibility id")
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line A", "Startup procedure", "# Start line A", []string{"startup"}, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create knowledge item: %v", err)
	}
	if item.CurrentRevision != 1 {
		t.Fatalf("expected revision 1, got %d", item.CurrentRevision)
	}

	item, err = app.AddRevision("bob", item.ID, "Start line A", "Startup procedure revised", "# Start line A\n\nUpdated", []string{"startup", "v2"})
	if err != nil {
		t.Fatalf("add revision: %v", err)
	}
	if item.CurrentRevision != 2 {
		t.Fatalf("expected revision 2, got %d", item.CurrentRevision)
	}

	place, err := app.CreatePlace("alice", "line", "Line A", "Primary startup area", "", []string{"startup"})
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "machine", "Machine A", "Main production machine", place.ID, []string{"line-a"})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}

	run, err := app.RecordRun("carol", RunKindProcedure, item.ID, 2, "completed", "Line started cleanly", "Machine A", "Bay 1", place.ID, []string{resource.ID}, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if run.Revision != 2 {
		t.Fatalf("expected linked revision 2, got %d", run.Revision)
	}
	if run.PlaceID != place.ID || len(run.ResourceIDs) != 1 || run.ResourceIDs[0] != resource.ID {
		t.Fatalf("unexpected run context: %+v", run)
	}

	run, err = app.AddEvidence("carol", run.ID, "Photo of startup checklist", map[string]string{"checklist": "passed"}, "checklist.txt", []byte("ok"))
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if len(run.Evidence) != 1 {
		t.Fatalf("expected one evidence record, got %d", len(run.Evidence))
	}
	if run.Evidence[0].AttachmentName != "checklist.txt" {
		t.Fatalf("unexpected attachment name %q", run.Evidence[0].AttachmentName)
	}

	if err := app.RecordApproval("boss", "knowledge_item", item.ID, 2, "reviewer", DecisionApproved, "ready for use"); err != nil {
		t.Fatalf("record item approval: %v", err)
	}
	if err := app.RecordApproval("boss", "run", run.ID, 0, "approver", DecisionNoted, "noted in shift handoff"); err != nil {
		t.Fatalf("record run approval: %v", err)
	}
	if err := app.AddLink("alice", "responsibility", responsibility.ID, "knowledge_item", item.ID, "owns", "line lead owns startup"); err != nil {
		t.Fatalf("add link: %v", err)
	}

	dashboard := app.Dashboard()
	if dashboard.Responsibilities != 1 || dashboard.Procedures != 1 || dashboard.ProcedureRuns != 1 {
		t.Fatalf("unexpected dashboard: %+v", dashboard)
	}
	if dashboard.Approvals != 2 {
		t.Fatalf("unexpected approvals count: %+v", dashboard)
	}
}

func TestAppPersistsAndReplaysState(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	resp, err := app.CreateResponsibility("alice", "Trainer", "Owns onboarding", []string{"trainer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindTraining, "Forklift intro", "Training guide", "Intro", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create knowledge item: %v", err)
	}
	run, err := app.RecordRun("dave", RunKindTraining, item.ID, 1, "passed", "Completed onboarding", "", "", "", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if _, err := app.AddEvidence("dave", run.ID, "Signed checklist", map[string]string{"result": "pass"}, "", nil); err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	reloaded, err := NewApp(root)
	if err != nil {
		t.Fatalf("reload app: %v", err)
	}
	reloadedItem, err := reloaded.GetKnowledgeItem(item.ID)
	if err != nil {
		t.Fatalf("get reloaded item: %v", err)
	}
	if len(reloadedItem.ResponsibilityIDs) != 1 || reloadedItem.ResponsibilityIDs[0] != resp.ID {
		t.Fatalf("unexpected responsibility links: %+v", reloadedItem.ResponsibilityIDs)
	}
	reloadedRun, err := reloaded.GetRun(run.ID)
	if err != nil {
		t.Fatalf("get reloaded run: %v", err)
	}
	if len(reloadedRun.Evidence) != 1 {
		t.Fatalf("expected evidence after reload, got %+v", reloadedRun)
	}
	if reloadedItem.WorkingBody != "Intro" || reloadedItem.Status != ItemStatusDraft {
		t.Fatalf("unexpected item live state after reload: %+v", reloadedItem)
	}
}

func TestAppSearchAndLinkingReflectOperationalFlow(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	responsibility, err := app.CreateResponsibility("alice", "Maintenance lead", "Owns press upkeep", []string{"maintainer"}, []string{"press"})
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindMaintenance, "Press lubrication", "Weekly press upkeep", "# Lubricate press", []string{"press", "weekly"}, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create knowledge item: %v", err)
	}
	place, err := app.CreatePlace("alice", "cell", "Cell A", "Press cell", "", []string{"press"})
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "machine", "Press 7", "Stamping press", place.ID, []string{"press"})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindMaintenance, item.ID, 1, "completed", "Press startup stayed smooth", "Press 7", "Cell A", place.ID, []string{resource.ID}, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if err := app.AddLink("alice", "knowledge_item", item.ID, "run", run.ID, "performed_by", "This run used the weekly press upkeep doc"); err != nil {
		t.Fatalf("add link: %v", err)
	}
	if err := app.RecordApproval("carol", "run", run.ID, 0, "reviewer", DecisionApproved, "Approved for handoff history"); err != nil {
		t.Fatalf("record run approval: %v", err)
	}

	updatedResponsibility, err := app.GetResponsibility(responsibility.ID)
	if err != nil {
		t.Fatalf("get responsibility: %v", err)
	}
	if len(updatedResponsibility.LinkedItemIDs) != 1 || updatedResponsibility.LinkedItemIDs[0] != item.ID {
		t.Fatalf("unexpected linked item ids: %+v", updatedResponsibility.LinkedItemIDs)
	}
	if len(updatedResponsibility.LinkedRunIDs) != 1 || updatedResponsibility.LinkedRunIDs[0] != run.ID {
		t.Fatalf("unexpected linked run ids: %+v", updatedResponsibility.LinkedRunIDs)
	}

	updatedRun, err := app.GetRun(run.ID)
	if err != nil {
		t.Fatalf("get run: %v", err)
	}
	if len(updatedRun.Approvals) != 1 || updatedRun.Approvals[0].Decision != DecisionApproved {
		t.Fatalf("unexpected run approvals: %+v", updatedRun.Approvals)
	}
	if len(updatedRun.Links) != 1 || updatedRun.Links[0].Relation != "performed_by" {
		t.Fatalf("unexpected run links: %+v", updatedRun.Links)
	}

	search := app.Search("press")
	foundResponsibilities := search["responsibilities"].([]Responsibility)
	foundItems := search["items"].([]KnowledgeItem)
	foundRuns := search["runs"].([]RunRecord)
	foundPlaces := search["places"].([]Place)
	foundResources := search["resources"].([]Resource)
	if len(foundResponsibilities) != 1 || foundResponsibilities[0].ID != responsibility.ID {
		t.Fatalf("unexpected responsibility search result: %+v", foundResponsibilities)
	}
	if len(foundItems) != 1 || foundItems[0].ID != item.ID {
		t.Fatalf("unexpected item search result: %+v", foundItems)
	}
	if len(foundRuns) != 1 || foundRuns[0].ID != run.ID {
		t.Fatalf("unexpected run search result: %+v", foundRuns)
	}
	if len(foundPlaces) != 1 || foundPlaces[0].ID != place.ID {
		t.Fatalf("unexpected place search result: %+v", foundPlaces)
	}
	if len(foundResources) != 1 || foundResources[0].ID != resource.ID {
		t.Fatalf("unexpected resource search result: %+v", foundResources)
	}
}

func TestAppTracksLiveDraftsStatusesAndSupersedence(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count bins", "Cycle count flow", "# Count bins", []string{"inventory"}, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}

	state, err := app.LiveItemState(item.ID)
	if err != nil {
		t.Fatalf("live item state: %v", err)
	}
	if state.Version != 1 || state.Status != ItemStatusDraft || state.Body != "# Count bins" {
		t.Fatalf("unexpected initial live state: %+v", state)
	}

	state, conflict, err := app.UpdateLiveItem(item.ID, "browser-alice", "Alice", "#0c6d62", 4, 4, true, state.Version, "# Count bins\n\nUpdated notes")
	if err != nil {
		t.Fatalf("update live item: %v", err)
	}
	if conflict {
		t.Fatalf("did not expect conflict: %+v", state)
	}
	if state.Version != 2 || len(state.Participants) != 1 || state.Body != "# Count bins\n\nUpdated notes" {
		t.Fatalf("unexpected updated live state: %+v", state)
	}

	if err := app.RecordApproval("boss", "knowledge_item", item.ID, 1, "reviewer", DecisionApproved, "ready for use"); err != nil {
		t.Fatalf("record approval: %v", err)
	}
	approved, err := app.GetKnowledgeItem(item.ID)
	if err != nil {
		t.Fatalf("get approved item: %v", err)
	}
	if approved.Status != ItemStatusApproved {
		t.Fatalf("expected approved status, got %+v", approved)
	}

	superseded, err := app.SupersedeKnowledgeItem("boss", item.ID, "Replaced by new count flow")
	if err != nil {
		t.Fatalf("supersede item: %v", err)
	}
	if superseded.Status != ItemStatusSuperseded {
		t.Fatalf("expected superseded status, got %+v", superseded)
	}
}

func TestAppPersistsDraftsAndRejectsStaleLiveUpdates(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start shift", "Shift startup", "# Start shift", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}

	state, conflict, err := app.UpdateLiveItem(item.ID, "browser-a", "Alice", "#123456", 3, 3, true, 1, "# Start shift\n\nChecked PPE.")
	if err != nil {
		t.Fatalf("first live update: %v", err)
	}
	if conflict || state.Version != 2 {
		t.Fatalf("unexpected first live update result: conflict=%v state=%+v", conflict, state)
	}

	staleState, conflict, err := app.UpdateLiveItem(item.ID, "browser-b", "Bob", "#654321", 2, 2, true, 1, "# Start shift\n\nStale overwrite.")
	if err != nil {
		t.Fatalf("stale live update: %v", err)
	}
	if !conflict {
		t.Fatalf("expected stale update conflict, got state %+v", staleState)
	}
	if staleState.Body != "# Start shift\n\nChecked PPE." {
		t.Fatalf("stale update should not overwrite body: %+v", staleState)
	}

	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	reloaded, err := NewApp(root)
	if err != nil {
		t.Fatalf("reload app: %v", err)
	}
	reloadedItem, err := reloaded.GetKnowledgeItem(item.ID)
	if err != nil {
		t.Fatalf("get reloaded item: %v", err)
	}
	if reloadedItem.WorkingBody != "# Start shift\n\nChecked PPE." || reloadedItem.WorkingVersion != 2 {
		t.Fatalf("expected persisted draft after reload, got %+v", reloadedItem)
	}
}

func TestAppSearchWithStructuredFilters(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	responsibility, err := app.CreateResponsibility("alice", "Receiving lead", "Owns receiving checks", []string{"reviewer"}, []string{"receiving"})
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", []string{"receiving"})
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, []string{"parts"})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count flow", "# Count receiving bin", []string{"inventory"}, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.RecordApproval("boss", "knowledge_item", item.ID, 1, "reviewer", DecisionApproved, "Ready to use"); err != nil {
		t.Fatalf("record approval: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindInventory, item.ID, 1, "completed", "Counted receiving bin", "", "", place.ID, []string{resource.ID}, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}

	approvedSearch := app.SearchWithOptions(SearchOptions{
		Kind:   KnowledgeKindInventory,
		Status: ItemStatusApproved,
	})
	items := approvedSearch["items"].([]KnowledgeItem)
	if len(items) != 1 || items[0].ID != item.ID {
		t.Fatalf("unexpected approved item filter result: %+v", items)
	}

	placeSearch := app.SearchWithOptions(SearchOptions{
		PlaceID: place.ID,
	})
	places := placeSearch["places"].([]Place)
	resources := placeSearch["resources"].([]Resource)
	runs := placeSearch["runs"].([]RunRecord)
	if len(places) != 1 || places[0].ID != place.ID {
		t.Fatalf("unexpected place filter result: %+v", places)
	}
	if len(resources) != 1 || resources[0].ID != resource.ID {
		t.Fatalf("unexpected resource filter result: %+v", resources)
	}
	if len(runs) != 1 || runs[0].ID != run.ID {
		t.Fatalf("unexpected run filter result: %+v", runs)
	}

	respSearch := app.SearchWithOptions(SearchOptions{
		ResponsibilityID: responsibility.ID,
	})
	responsibilities := respSearch["responsibilities"].([]Responsibility)
	respItems := respSearch["items"].([]KnowledgeItem)
	respRuns := respSearch["runs"].([]RunRecord)
	if len(responsibilities) != 1 || responsibilities[0].ID != responsibility.ID {
		t.Fatalf("unexpected responsibility filter result: %+v", responsibilities)
	}
	if len(respItems) != 1 || respItems[0].ID != item.ID {
		t.Fatalf("unexpected responsibility-linked item result: %+v", respItems)
	}
	if len(respRuns) != 1 || respRuns[0].ID != run.ID {
		t.Fatalf("unexpected responsibility-linked run result: %+v", respRuns)
	}
}
