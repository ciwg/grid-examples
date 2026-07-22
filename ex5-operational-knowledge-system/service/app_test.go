package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
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

func TestAppEvidenceAttachmentsStayImmutableAcrossRepeatedFilenames(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "Startup flow", "# Start line", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindProcedure, item.ID, 1, "completed", "Normal startup", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}

	run, err = app.AddEvidence("bob", run.ID, "Photo one", map[string]string{"shot": "1"}, "photo.txt", []byte("first-bytes"))
	if err != nil {
		t.Fatalf("add first evidence: %v", err)
	}
	firstPath := run.Evidence[0].AttachmentPath

	run, err = app.AddEvidence("bob", run.ID, "Photo two", map[string]string{"shot": "2"}, "photo.txt", []byte("second-bytes"))
	if err != nil {
		t.Fatalf("add second evidence: %v", err)
	}
	if len(run.Evidence) != 2 {
		t.Fatalf("expected two evidence records, got %+v", run.Evidence)
	}
	secondPath := run.Evidence[1].AttachmentPath
	if firstPath == secondPath {
		t.Fatalf("expected unique attachment paths, got %q", firstPath)
	}

	firstBody, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatalf("read first attachment: %v", err)
	}
	secondBody, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatalf("read second attachment: %v", err)
	}
	if string(firstBody) != "first-bytes" {
		t.Fatalf("unexpected first attachment body: %q", string(firstBody))
	}
	if string(secondBody) != "second-bytes" {
		t.Fatalf("unexpected second attachment body: %q", string(secondBody))
	}
}

func TestAppReplaysLargeKnowledgeBodiesAfterRestart(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	initialBody := strings.Repeat("A", 70*1024)
	revisionBody := strings.Repeat("B", 80*1024)
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Large procedure", "Large startup flow", initialBody, nil, nil)
	if err != nil {
		t.Fatalf("create large item: %v", err)
	}
	item, err = app.AddRevision("alice", item.ID, "Large procedure", "Large startup flow revised", revisionBody, nil)
	if err != nil {
		t.Fatalf("add large revision: %v", err)
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
	if reloadedItem.CurrentRevision != 2 {
		t.Fatalf("unexpected current revision after reload: %d", reloadedItem.CurrentRevision)
	}
	if reloadedItem.WorkingBody != revisionBody {
		t.Fatalf("unexpected working body size after reload: got %d want %d", len(reloadedItem.WorkingBody), len(revisionBody))
	}
	if len(reloadedItem.Revisions) != 2 {
		t.Fatalf("unexpected revision count after reload: %d", len(reloadedItem.Revisions))
	}
	if reloadedItem.Revisions[0].Body != initialBody {
		t.Fatalf("initial large body did not replay correctly: got %d want %d", len(reloadedItem.Revisions[0].Body), len(initialBody))
	}
	if reloadedItem.Revisions[1].Body != revisionBody {
		t.Fatalf("revised large body did not replay correctly: got %d want %d", len(reloadedItem.Revisions[1].Body), len(revisionBody))
	}
}

func TestAppWritesAndReplaysSignedKnowledgeItemRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "Startup flow", "# Start line", []string{"startup"}, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	item, err = app.AddRevision("alice", item.ID, "Start line", "Startup flow v2", "# Start line\n\nv2", []string{"startup", "v2"})
	if err != nil {
		t.Fatalf("add revision: %v", err)
	}
	if err := app.RecordApproval("bob", "knowledge_item", item.ID, item.CurrentRevision, "reviewer", DecisionApproved, "ready"); err != nil {
		t.Fatalf("approve item: %v", err)
	}

	meta := app.Meta()
	if meta.KnowledgeItemPCID != protocols.KnowledgeItemProfile.CID.String() {
		t.Fatalf("unexpected knowledge-item pCID in meta: got %q want %q", meta.KnowledgeItemPCID, protocols.KnowledgeItemProfile.CID.String())
	}

	recordBody, err := os.ReadFile(filepath.Join(root, "knowledge-item-messages.jsonl"))
	if err != nil {
		t.Fatalf("read knowledge-item messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 signed knowledge-item records, got %d", len(lines))
	}
	var firstRecord SignedKnowledgeItemRecord
	if err := json.Unmarshal([]byte(lines[0]), &firstRecord); err != nil {
		t.Fatalf("decode first record: %v", err)
	}
	if firstRecord.PCID != protocols.KnowledgeItemProfile.CID.String() {
		t.Fatalf("unexpected record pCID: got %q want %q", firstRecord.PCID, protocols.KnowledgeItemProfile.CID.String())
	}
	if firstRecord.EventType != "knowledge_item_created" {
		t.Fatalf("unexpected first record event type %q", firstRecord.EventType)
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
	if reloadedItem.Status != ItemStatusApproved {
		t.Fatalf("unexpected item status after signed replay verification: %s", reloadedItem.Status)
	}
	if reloadedItem.CurrentRevision != 2 {
		t.Fatalf("unexpected revision after signed replay verification: %d", reloadedItem.CurrentRevision)
	}
}

func TestAppRejectsTamperedSignedKnowledgeItemRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "Startup flow", "# Start line", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if _, err := app.AddRevision("alice", item.ID, "Start line", "Startup flow v2", "# Start line\n\nv2", nil); err != nil {
		t.Fatalf("add revision: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	recordPath := filepath.Join(root, "knowledge-item-messages.jsonl")
	recordBody, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read knowledge-item messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected at least one signed knowledge-item record")
	}
	var record SignedKnowledgeItemRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode record: %v", err)
	}
	record.EnvelopeCID = "bafkreitampered"
	tamperedLine, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("encode tampered record: %v", err)
	}
	lines[0] = string(tamperedLine)
	tampered := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(recordPath, []byte(tampered), 0o644); err != nil {
		t.Fatalf("rewrite tampered record: %v", err)
	}

	if _, err := NewApp(root); err == nil || !strings.Contains(err.Error(), "envelope cid mismatch") {
		t.Fatalf("expected envelope cid mismatch after tampering, got %v", err)
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
	if len(updatedResponsibility.Links) != 0 {
		t.Fatalf("did not expect implicit responsibility links from run linkage: %+v", updatedResponsibility.Links)
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

func TestAppValidatesTypedLinkEndpointsAndProjectsResponsibilityLinks(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns receiving checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Dock receipt", "# Inspect inbound pallet", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}

	if err := app.AddLink("alice", "responsibility", resp.ID, "knowledge_item", item.ID, "owns", "Receiving lead owns dock receipt"); err != nil {
		t.Fatalf("add responsibility->item link: %v", err)
	}
	if err := app.AddLink("alice", "responsibility", resp.ID, "run", run.ID, "reviews", "Receiving lead reviews the receiving run"); err != nil {
		t.Fatalf("add responsibility->run link: %v", err)
	}

	loadedResp, err := app.GetResponsibility(resp.ID)
	if err != nil {
		t.Fatalf("get responsibility: %v", err)
	}
	if len(loadedResp.Links) != 2 {
		t.Fatalf("expected projected responsibility links, got %+v", loadedResp.Links)
	}
	if loadedResp.Links[0].FromType != "responsibility" || loadedResp.Links[1].FromType != "responsibility" {
		t.Fatalf("unexpected responsibility link projection: %+v", loadedResp.Links)
	}

	err = app.AddLink("alice", "responsibility", "RESP-9999", "knowledge_item", item.ID, "owns", "bad endpoint")
	if err == nil {
		t.Fatalf("expected dangling responsibility link to fail")
	}
	if !strings.Contains(err.Error(), "from endpoint invalid") {
		t.Fatalf("unexpected dangling endpoint error: %v", err)
	}

	err = app.AddLink("alice", "item", item.ID, "run", run.ID, "uses", "bad type")
	if err == nil {
		t.Fatalf("expected unsupported endpoint type to fail")
	}
	if !strings.Contains(err.Error(), "unsupported link endpoint type") {
		t.Fatalf("unexpected unsupported type error: %v", err)
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

	state, conflict, err := app.UpdateLiveItem(item.ID, "browser-alice", "Alice", "#0c6d62", 4, 4, true, state.Version, true, "# Count bins\n\nUpdated notes")
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

	state, conflict, err := app.UpdateLiveItem(item.ID, "browser-a", "Alice", "#123456", 3, 3, true, 1, true, "# Start shift\n\nChecked PPE.")
	if err != nil {
		t.Fatalf("first live update: %v", err)
	}
	if conflict || state.Version != 2 {
		t.Fatalf("unexpected first live update result: conflict=%v state=%+v", conflict, state)
	}

	staleState, conflict, err := app.UpdateLiveItem(item.ID, "browser-b", "Bob", "#654321", 2, 2, true, 1, true, "# Start shift\n\nStale overwrite.")
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

func TestAppRejectsStaleRevisionApprovalForKnowledgeItems(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start shift", "Shift startup", "# Start shift", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	item, err = app.AddRevision("alice", item.ID, "Start shift", "Shift startup revised", "# Start shift\n\nRevision 2", nil)
	if err != nil {
		t.Fatalf("add revision: %v", err)
	}
	err = app.RecordApproval("boss", "knowledge_item", item.ID, 1, "reviewer", DecisionApproved, "stale approval")
	if err == nil {
		t.Fatalf("expected stale revision approval to fail")
	}
	if !strings.Contains(err.Error(), "stale revision") {
		t.Fatalf("unexpected stale approval error: %v", err)
	}

	loaded, err := app.GetKnowledgeItem(item.ID)
	if err != nil {
		t.Fatalf("get item: %v", err)
	}
	if loaded.Status != ItemStatusDraft {
		t.Fatalf("stale approval should not change status: %+v", loaded)
	}
	if len(loaded.Approvals) != 0 {
		t.Fatalf("stale approval should not append approval record: %+v", loaded.Approvals)
	}
}

func TestAppAllowsClearingLiveDraftBody(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start shift", "Shift startup", "# Start shift", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	state, conflict, err := app.UpdateLiveItem(item.ID, "browser-a", "Alice", "#123456", 3, 3, true, 1, true, "# Start shift\n\nChecked PPE.")
	if err != nil {
		t.Fatalf("first live update: %v", err)
	}
	if conflict || state.Version != 2 {
		t.Fatalf("unexpected first live update result: conflict=%v state=%+v", conflict, state)
	}

	state, conflict, err = app.UpdateLiveItem(item.ID, "browser-a", "Alice", "#123456", 0, 0, true, 2, true, "")
	if err != nil {
		t.Fatalf("clear live body: %v", err)
	}
	if conflict {
		t.Fatalf("did not expect clear-body conflict: %+v", state)
	}
	if state.Version != 3 || state.Body != "" {
		t.Fatalf("expected cleared live body at version 3, got %+v", state)
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
	if reloadedItem.WorkingBody != "" || reloadedItem.WorkingVersion != 3 {
		t.Fatalf("expected cleared live draft after reload, got %+v", reloadedItem)
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

	outcomeSearch := app.SearchWithOptions(SearchOptions{
		Kind:    RunKindInventory,
		Outcome: "completed",
		PlaceID: place.ID,
	})
	outcomeRuns := outcomeSearch["runs"].([]RunRecord)
	if len(outcomeRuns) != 1 || outcomeRuns[0].ID != run.ID {
		t.Fatalf("unexpected outcome-filtered run result: %+v", outcomeRuns)
	}
}

func TestAppSearchIncludesRunEvidenceAndApprovalHistory(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = app.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{
		"supplier":     "Acme Parts",
		"packing_slip": "PS-1234",
		"condition":    "wrap torn",
	}, "", nil)
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if err := app.RecordApproval("boss", "run", run.ID, 0, "reviewer", DecisionApproved, "Reviewed at dock"); err != nil {
		t.Fatalf("record approval: %v", err)
	}

	evidenceSearch := app.SearchWithOptions(SearchOptions{Query: "acme parts"})
	evidenceRuns := evidenceSearch["runs"].([]RunRecord)
	if len(evidenceRuns) != 1 || evidenceRuns[0].ID != run.ID {
		t.Fatalf("expected evidence search to find run by facts, got %+v", evidenceRuns)
	}

	approvalSearch := app.SearchWithOptions(SearchOptions{Query: "reviewed at dock"})
	approvalRuns := approvalSearch["runs"].([]RunRecord)
	if len(approvalRuns) != 1 || approvalRuns[0].ID != run.ID {
		t.Fatalf("expected approval search to find run by approval notes, got %+v", approvalRuns)
	}
}

func TestGetKnowledgeItemIncludesRelatedRuns(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "Startup flow", "# Start line", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	_, err = app.RecordRun("bob", RunKindProcedure, item.ID, 1, "completed", "Normal startup", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record first run: %v", err)
	}
	secondRun, err := app.RecordRun("carol", RunKindProcedure, item.ID, 1, "completed", "Second startup", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record second run: %v", err)
	}

	loaded, err := app.GetKnowledgeItem(item.ID)
	if err != nil {
		t.Fatalf("get item: %v", err)
	}
	if len(loaded.RelatedRuns) != 2 {
		t.Fatalf("expected related runs on item, got %+v", loaded.RelatedRuns)
	}
	if loaded.RelatedRuns[1].ID != secondRun.ID {
		t.Fatalf("unexpected related run ordering: %+v", loaded.RelatedRuns)
	}
}

func TestContextRecordsIncludeRelatedRuns(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns receiving checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count flow", "# Count receiving bin", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindInventory, item.ID, 1, "completed", "Counted receiving bin", "", "", place.ID, []string{resource.ID}, []string{resp.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = app.AddEvidence("bob", run.ID, "Cycle count", map[string]string{
		"expected_count": "12",
		"actual_count":   "10",
		"discrepancy":    "-2",
	}, "", nil)
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}

	loadedPlace, err := app.GetPlace(place.ID)
	if err != nil {
		t.Fatalf("get place: %v", err)
	}
	if len(loadedPlace.RelatedRuns) != 1 || loadedPlace.RelatedRuns[0].ID != run.ID {
		t.Fatalf("unexpected place related runs: %+v", loadedPlace.RelatedRuns)
	}
	if len(loadedPlace.RelatedRuns[0].Evidence) != 1 || loadedPlace.RelatedRuns[0].Evidence[0].Facts["discrepancy"] != "-2" {
		t.Fatalf("expected place related run evidence facts, got %+v", loadedPlace.RelatedRuns[0])
	}

	loadedResource, err := app.GetResource(resource.ID)
	if err != nil {
		t.Fatalf("get resource: %v", err)
	}
	if len(loadedResource.RelatedRuns) != 1 || loadedResource.RelatedRuns[0].ID != run.ID {
		t.Fatalf("unexpected resource related runs: %+v", loadedResource.RelatedRuns)
	}
	if len(loadedResource.RelatedRuns[0].Evidence) != 1 || loadedResource.RelatedRuns[0].Evidence[0].Facts["actual_count"] != "10" {
		t.Fatalf("expected resource related run evidence facts, got %+v", loadedResource.RelatedRuns[0])
	}

	loadedResp, err := app.GetResponsibility(resp.ID)
	if err != nil {
		t.Fatalf("get responsibility: %v", err)
	}
	if len(loadedResp.RelatedRuns) != 1 || loadedResp.RelatedRuns[0].ID != run.ID {
		t.Fatalf("unexpected responsibility related runs: %+v", loadedResp.RelatedRuns)
	}
	if len(loadedResp.RelatedRuns[0].Evidence) != 1 || loadedResp.RelatedRuns[0].Evidence[0].Facts["expected_count"] != "12" {
		t.Fatalf("expected responsibility related run evidence facts, got %+v", loadedResp.RelatedRuns[0])
	}
}

func TestAppTracksReceivingCheckKindsAndDashboardCounts(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	responsibility, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, []string{"receiving"})
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	place, err := app.CreatePlace("alice", "dock", "Dock A", "Inbound receiving dock", "", []string{"receiving"})
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "Intake pallet", "Pallet staged for receipt", place.ID, []string{"inbound"})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check for inbound pallet", "# Inspect inbound pallet", []string{"receiving"}, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", place.ID, []string{resource.ID}, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = app.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{
		"supplier":       "Acme Parts",
		"packing_slip":   "PS-1234",
		"received_units": "18",
		"expected_units": "20",
		"variance":       "-2",
		"condition":      "wrap torn",
	}, "", nil)
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}

	meta := app.Meta()
	if got := meta.KnowledgeKinds[3]; got != KnowledgeKindReceiving {
		t.Fatalf("expected receiving kind in meta, got %+v", meta.KnowledgeKinds)
	}
	if got := meta.RunKinds[3]; got != RunKindReceiving {
		t.Fatalf("expected receiving run kind in meta, got %+v", meta.RunKinds)
	}

	dashboard := app.Dashboard()
	if dashboard.ReceivingItems != 1 || dashboard.ReceivingRuns != 1 {
		t.Fatalf("unexpected receiving dashboard counts: %+v", dashboard)
	}
	if len(run.Evidence) != 1 {
		t.Fatalf("expected receiving evidence, got %+v", run)
	}

	search := app.SearchWithOptions(SearchOptions{Kind: KnowledgeKindReceiving, PlaceID: place.ID})
	items := search["items"].([]KnowledgeItem)
	runs := search["runs"].([]RunRecord)
	if len(items) != 1 || items[0].ID != item.ID {
		t.Fatalf("unexpected receiving item search result: %+v", items)
	}
	if len(runs) != 1 || runs[0].ID != run.ID {
		t.Fatalf("unexpected receiving run search result: %+v", runs)
	}
}

func TestProblemReviewGroupsByPlaceAndResource(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", []string{"receiving"})
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, []string{"parts"})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	recvItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create receiving item: %v", err)
	}
	recvRun, err := app.RecordRun("bob", RunKindReceiving, recvItem.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record receiving run: %v", err)
	}
	if _, err := app.AddEvidence("bob", recvRun.ID, "Receiving inspection", map[string]string{
		"supplier":  "Acme Parts",
		"variance":  "-2",
		"condition": "wrap torn",
	}, "", nil); err != nil {
		t.Fatalf("add receiving evidence: %v", err)
	}
	invItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count", "# Count receiving bin", nil, nil)
	if err != nil {
		t.Fatalf("create inventory item: %v", err)
	}
	invRun, err := app.RecordRun("bob", RunKindInventory, invItem.ID, 1, "completed", "Counted receiving bin", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record inventory run: %v", err)
	}
	if _, err := app.AddEvidence("bob", invRun.ID, "Cycle count", map[string]string{
		"expected_count": "12",
		"actual_count":   "10",
		"discrepancy":    "-2",
	}, "", nil); err != nil {
		t.Fatalf("add inventory evidence: %v", err)
	}

	review := app.ProblemReview()
	if review.ProblemRuns != 2 {
		t.Fatalf("unexpected problem run count: %+v", review)
	}
	if len(review.PlaceGroups) != 1 || review.PlaceGroups[0].GroupID != place.ID {
		t.Fatalf("unexpected place problem groups: %+v", review.PlaceGroups)
	}
	if review.PlaceGroups[0].ProblemCount != 2 || review.PlaceGroups[0].ReceivingProblems != 1 || review.PlaceGroups[0].InventoryProblems != 1 {
		t.Fatalf("unexpected place group counts: %+v", review.PlaceGroups[0])
	}
	if len(review.ResourceGroups) != 1 || review.ResourceGroups[0].GroupID != resource.ID {
		t.Fatalf("unexpected resource problem groups: %+v", review.ResourceGroups)
	}
	if review.ResourceGroups[0].ProblemCount != 2 {
		t.Fatalf("unexpected resource group counts: %+v", review.ResourceGroups[0])
	}
	highlights := strings.Join(review.PlaceGroups[0].HighlightExamples, " | ")
	if !strings.Contains(highlights, "outcome: accepted_with_notes") || !strings.Contains(highlights, "discrepancy: -2") {
		t.Fatalf("unexpected grouped highlights: %s", highlights)
	}
}

func TestAppSearchWithProblemFilterMatchesProblemReviewLogic(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	recvItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create receiving item: %v", err)
	}
	recvRun, err := app.RecordRun("bob", RunKindReceiving, recvItem.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record receiving run: %v", err)
	}
	if _, err := app.AddEvidence("bob", recvRun.ID, "Receiving inspection", map[string]string{"variance": "-2", "condition": "wrap torn"}, "", nil); err != nil {
		t.Fatalf("add receiving evidence: %v", err)
	}
	invItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindInventory, "Count receiving bin", "Cycle count", "# Count receiving bin", nil, nil)
	if err != nil {
		t.Fatalf("create inventory item: %v", err)
	}
	invRun, err := app.RecordRun("bob", RunKindInventory, invItem.ID, 1, "completed", "Counted receiving bin", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record inventory run: %v", err)
	}
	if _, err := app.AddEvidence("bob", invRun.ID, "Cycle count", map[string]string{"expected_count": "12", "actual_count": "10", "discrepancy": "-2"}, "", nil); err != nil {
		t.Fatalf("add inventory evidence: %v", err)
	}
	normalRun, err := app.RecordRun("bob", RunKindReceiving, recvItem.ID, 1, "accepted", "Non-problem run", "", "", place.ID, []string{resource.ID}, nil)
	if err != nil {
		t.Fatalf("record normal run: %v", err)
	}

	search := app.SearchWithOptions(SearchOptions{PlaceID: place.ID, Problem: true})
	runs := search["runs"].([]RunRecord)
	if len(runs) != 2 {
		t.Fatalf("expected only problem runs, got %+v", runs)
	}
	if runs[0].ID != recvRun.ID && runs[1].ID != recvRun.ID {
		t.Fatalf("receiving problem run missing from problem search: %+v", runs)
	}
	if runs[0].ID != invRun.ID && runs[1].ID != invRun.ID {
		t.Fatalf("inventory problem run missing from problem search: %+v", runs)
	}
	if runs[0].ID == normalRun.ID || runs[1].ID == normalRun.ID {
		t.Fatalf("non-problem run leaked into problem search: %+v", runs)
	}
}

func TestAppSearchIndexesRecordIDsAcrossAllRecordTypes(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", place.ID, []string{resource.ID}, []string{resp.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}

	cases := []struct {
		name     string
		query    string
		key      string
		expected string
	}{
		{name: "place id", query: place.ID, key: "places", expected: place.ID},
		{name: "resource id", query: resource.ID, key: "resources", expected: resource.ID},
		{name: "responsibility id", query: resp.ID, key: "responsibilities", expected: resp.ID},
		{name: "item id", query: item.ID, key: "items", expected: item.ID},
		{name: "run id", query: run.ID, key: "runs", expected: run.ID},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			search := app.SearchWithOptions(SearchOptions{Query: tc.query})
			switch tc.key {
			case "places":
				results := search["places"].([]Place)
				if len(results) != 1 || results[0].ID != tc.expected {
					t.Fatalf("unexpected place search results: %+v", results)
				}
			case "resources":
				results := search["resources"].([]Resource)
				if len(results) != 1 || results[0].ID != tc.expected {
					t.Fatalf("unexpected resource search results: %+v", results)
				}
			case "responsibilities":
				results := search["responsibilities"].([]Responsibility)
				if len(results) != 1 || results[0].ID != tc.expected {
					t.Fatalf("unexpected responsibility search results: %+v", results)
				}
			case "items":
				results := search["items"].([]KnowledgeItem)
				if len(results) != 1 || results[0].ID != tc.expected {
					t.Fatalf("unexpected item search results: %+v", results)
				}
			case "runs":
				results := search["runs"].([]RunRecord)
				if len(results) != 1 || results[0].ID != tc.expected {
					t.Fatalf("unexpected run search results: %+v", results)
				}
			}
		})
	}
}

func TestAppProblemSearchFiltersNonRunGroupsToProblemContext(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	problemResp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create problem responsibility: %v", err)
	}
	normalResp, err := app.CreateResponsibility("alice", "Shipping lead", "Owns outbound checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create normal responsibility: %v", err)
	}
	problemPlace, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", nil)
	if err != nil {
		t.Fatalf("create problem place: %v", err)
	}
	normalPlace, err := app.CreatePlace("alice", "area", "Shipping", "Outbound staging", "", nil)
	if err != nil {
		t.Fatalf("create normal place: %v", err)
	}
	problemResource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", problemPlace.ID, nil)
	if err != nil {
		t.Fatalf("create problem resource: %v", err)
	}
	normalResource, err := app.CreateResource("alice", "container", "Outbound tote", "Shipping tote", normalPlace.ID, nil)
	if err != nil {
		t.Fatalf("create normal resource: %v", err)
	}
	problemItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{problemResp.ID})
	if err != nil {
		t.Fatalf("create problem item: %v", err)
	}
	normalItem, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Stage outbound carton", "Shipping checklist", "# Stage outbound carton", nil, []string{normalResp.ID})
	if err != nil {
		t.Fatalf("create normal item: %v", err)
	}
	problemRun, err := app.RecordRun("bob", RunKindReceiving, problemItem.ID, 1, "accepted_with_notes", "Outer wrap torn", "", "", problemPlace.ID, []string{problemResource.ID}, []string{problemResp.ID})
	if err != nil {
		t.Fatalf("record problem run: %v", err)
	}
	if _, err := app.AddEvidence("bob", problemRun.ID, "Receiving inspection", map[string]string{"condition": "wrap torn"}, "", nil); err != nil {
		t.Fatalf("add problem evidence: %v", err)
	}
	if _, err := app.RecordRun("bob", RunKindProcedure, normalItem.ID, 1, "completed", "Shipped carton", "", "", normalPlace.ID, []string{normalResource.ID}, []string{normalResp.ID}); err != nil {
		t.Fatalf("record normal run: %v", err)
	}

	search := app.SearchWithOptions(SearchOptions{Problem: true})
	places := search["places"].([]Place)
	resources := search["resources"].([]Resource)
	responsibilities := search["responsibilities"].([]Responsibility)
	items := search["items"].([]KnowledgeItem)
	runs := search["runs"].([]RunRecord)

	if len(places) != 1 || places[0].ID != problemPlace.ID {
		t.Fatalf("unexpected problem places: %+v", places)
	}
	if len(resources) != 1 || resources[0].ID != problemResource.ID {
		t.Fatalf("unexpected problem resources: %+v", resources)
	}
	if len(responsibilities) != 1 || responsibilities[0].ID != problemResp.ID {
		t.Fatalf("unexpected problem responsibilities: %+v", responsibilities)
	}
	if len(items) != 1 || items[0].ID != problemItem.ID {
		t.Fatalf("unexpected problem items: %+v", items)
	}
	if len(runs) != 1 || runs[0].ID != problemRun.ID {
		t.Fatalf("unexpected problem runs: %+v", runs)
	}
}
