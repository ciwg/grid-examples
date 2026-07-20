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

	run, err := app.RecordRun("carol", RunKindProcedure, item.ID, 2, "completed", "Line started cleanly", "Machine A", "Bay 1", []string{responsibility.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if run.Revision != 2 {
		t.Fatalf("expected linked revision 2, got %d", run.Revision)
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
	run, err := app.RecordRun("dave", RunKindTraining, item.ID, 1, "passed", "Completed onboarding", "", "", []string{resp.ID})
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
	run, err := app.RecordRun("bob", RunKindMaintenance, item.ID, 1, "completed", "Press startup stayed smooth", "Press 7", "Cell A", []string{responsibility.ID})
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
	if len(foundResponsibilities) != 1 || foundResponsibilities[0].ID != responsibility.ID {
		t.Fatalf("unexpected responsibility search result: %+v", foundResponsibilities)
	}
	if len(foundItems) != 1 || foundItems[0].ID != item.ID {
		t.Fatalf("unexpected item search result: %+v", foundItems)
	}
	if len(foundRuns) != 1 || foundRuns[0].ID != run.ID {
		t.Fatalf("unexpected run search result: %+v", foundRuns)
	}
}
