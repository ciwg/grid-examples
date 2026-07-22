package service

import (
	"encoding/base64"
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

	if _, err := NewApp(root); err == nil || !strings.Contains(err.Error(), "envelope CAS cid mismatch") {
		t.Fatalf("expected authoritative CAS cid mismatch after tampering, got %v", err)
	}
}

func TestAppWritesAndReplaysSignedKnowledgeApprovalRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "Startup flow", "# Start line", []string{"startup"}, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.RecordApproval("bob", "knowledge_item", item.ID, item.CurrentRevision, "reviewer", DecisionApproved, "ready"); err != nil {
		t.Fatalf("approve item: %v", err)
	}
	run, err := app.RecordRun("carol", RunKindProcedure, item.ID, item.CurrentRevision, "completed", "Executed startup flow", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if err := app.RecordApproval("dave", "run", run.ID, 0, "auditor", DecisionNoted, "captured"); err != nil {
		t.Fatalf("approve run: %v", err)
	}

	meta := app.Meta()
	if meta.KnowledgeApprovalPCID != protocols.KnowledgeApprovalProfile.CID.String() {
		t.Fatalf("unexpected knowledge-approval pCID in meta: got %q want %q", meta.KnowledgeApprovalPCID, protocols.KnowledgeApprovalProfile.CID.String())
	}

	recordBody, err := os.ReadFile(filepath.Join(root, "knowledge-approval-messages.jsonl"))
	if err != nil {
		t.Fatalf("read knowledge-approval messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 signed knowledge-approval records, got %d", len(lines))
	}
	var firstRecord SignedKnowledgeApprovalRecord
	if err := json.Unmarshal([]byte(lines[0]), &firstRecord); err != nil {
		t.Fatalf("decode first approval record: %v", err)
	}
	if firstRecord.PCID != protocols.KnowledgeApprovalProfile.CID.String() {
		t.Fatalf("unexpected approval record pCID: got %q want %q", firstRecord.PCID, protocols.KnowledgeApprovalProfile.CID.String())
	}
	if firstRecord.TargetType != "knowledge_item" {
		t.Fatalf("unexpected first approval record target type %q", firstRecord.TargetType)
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
	if len(reloadedItem.Approvals) != 1 || reloadedItem.Approvals[0].Decision != DecisionApproved {
		t.Fatalf("unexpected item approvals after signed replay verification: %+v", reloadedItem.Approvals)
	}
	reloadedRun, err := reloaded.GetRun(run.ID)
	if err != nil {
		t.Fatalf("get reloaded run: %v", err)
	}
	if len(reloadedRun.Approvals) != 1 || reloadedRun.Approvals[0].Decision != DecisionNoted {
		t.Fatalf("unexpected run approvals after signed replay verification: %+v", reloadedRun.Approvals)
	}
}

func TestAppRejectsTamperedSignedKnowledgeApprovalRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "Startup flow", "# Start line", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.RecordApproval("bob", "knowledge_item", item.ID, item.CurrentRevision, "reviewer", DecisionApproved, "ready"); err != nil {
		t.Fatalf("approve item: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	recordPath := filepath.Join(root, "knowledge-approval-messages.jsonl")
	recordBody, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read knowledge-approval messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected at least one signed knowledge-approval record")
	}
	var record SignedKnowledgeApprovalRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode approval record: %v", err)
	}
	record.EnvelopeCID = "bafkreiapprovaltampered"
	tamperedLine, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("encode tampered approval record: %v", err)
	}
	lines[0] = string(tamperedLine)
	tampered := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(recordPath, []byte(tampered), 0o644); err != nil {
		t.Fatalf("rewrite tampered approval record: %v", err)
	}

	if _, err := NewApp(root); err == nil || !strings.Contains(err.Error(), "knowledge-approval") || !strings.Contains(err.Error(), "envelope CAS cid mismatch") {
		t.Fatalf("expected knowledge-approval authoritative CAS cid mismatch after tampering, got %v", err)
	}
}

func TestAppWritesAndReplaysSignedKnowledgeEvidenceRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = app.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{"variance": "-2", "condition": "wrap torn"}, "photo.txt", []byte("ok"))
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if len(run.Evidence) != 1 || run.Evidence[0].ID == "" {
		t.Fatalf("expected stable evidence id, got %+v", run.Evidence)
	}

	meta := app.Meta()
	if meta.KnowledgeEvidencePCID != protocols.KnowledgeEvidenceProfile.CID.String() {
		t.Fatalf("unexpected knowledge-evidence pCID in meta: got %q want %q", meta.KnowledgeEvidencePCID, protocols.KnowledgeEvidenceProfile.CID.String())
	}
	if !meta.CASObjectsEnabled || !meta.CASAttachmentBlobsEnabled || !meta.CASDraftBodiesEnabled {
		t.Fatalf("expected CAS capability flags in meta, got %+v", meta)
	}
	if meta.PrimaryEmbodimentAdapter != "local_http" {
		t.Fatalf("expected local_http embodiment adapter in meta, got %+v", meta)
	}

	recordBody, err := os.ReadFile(filepath.Join(root, "knowledge-evidence-messages.jsonl"))
	if err != nil {
		t.Fatalf("read knowledge-evidence messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 signed knowledge-evidence record, got %d", len(lines))
	}
	var record SignedKnowledgeEvidenceRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode evidence record: %v", err)
	}
	if record.PCID != protocols.KnowledgeEvidenceProfile.CID.String() {
		t.Fatalf("unexpected evidence record pCID: got %q want %q", record.PCID, protocols.KnowledgeEvidenceProfile.CID.String())
	}
	if record.EvidenceID != run.Evidence[0].AliasID {
		t.Fatalf("unexpected evidence alias in record: got %q want %q", record.EvidenceID, run.Evidence[0].AliasID)
	}

	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	reloaded, err := NewApp(root)
	if err != nil {
		t.Fatalf("reload app: %v", err)
	}
	reloadedRun, err := reloaded.GetRun(run.ID)
	if err != nil {
		t.Fatalf("get reloaded run: %v", err)
	}
	if len(reloadedRun.Evidence) != 1 {
		t.Fatalf("unexpected evidence after signed replay verification: %+v", reloadedRun.Evidence)
	}
	if reloadedRun.Evidence[0].ID != run.Evidence[0].ID {
		t.Fatalf("expected stable evidence id after replay, got %q want %q", reloadedRun.Evidence[0].ID, run.Evidence[0].ID)
	}
	if reloadedRun.Evidence[0].AttachmentPath == "" || reloadedRun.Evidence[0].AttachmentCID == "" || reloadedRun.Evidence[0].AttachmentSize != int64(len([]byte("ok"))) {
		t.Fatalf("unexpected attachment reference after replay: %+v", reloadedRun.Evidence[0])
	}
}

func TestAppRejectsTamperedSignedKnowledgeEvidenceRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if _, err := app.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{"variance": "-2"}, "photo.txt", []byte("ok")); err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	recordPath := filepath.Join(root, "knowledge-evidence-messages.jsonl")
	recordBody, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read knowledge-evidence messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected at least one signed knowledge-evidence record")
	}
	var record SignedKnowledgeEvidenceRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode evidence record: %v", err)
	}
	record.EnvelopeCID = "bafkreievidencetampered"
	tamperedLine, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("encode tampered evidence record: %v", err)
	}
	lines[0] = string(tamperedLine)
	tampered := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(recordPath, []byte(tampered), 0o644); err != nil {
		t.Fatalf("rewrite tampered evidence record: %v", err)
	}

	if _, err := NewApp(root); err == nil || !strings.Contains(err.Error(), "knowledge-evidence") || !strings.Contains(err.Error(), "envelope CAS cid mismatch") {
		t.Fatalf("expected knowledge-evidence authoritative CAS cid mismatch after tampering, got %v", err)
	}
}

func TestAppWritesAndReplaysSignedOperationalRunRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "Dock bin", "Holds inbound parts", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Outer wrap torn", "scanner-1", "dock-a", place.ID, []string{resource.ID}, []string{resp.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	recordBody, err := os.ReadFile(filepath.Join(root, "operational-run-messages.jsonl"))
	if err != nil {
		t.Fatalf("read operational-run messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 signed operational-run record, got %d", len(lines))
	}
	var record SignedOperationalRunRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode operational-run record: %v", err)
	}
	if record.RunID != run.AliasID || record.ItemID != item.ID {
		t.Fatalf("unexpected operational-run record %+v", record)
	}
	if record.PCID != protocols.OperationalRunProfile.CID.String() {
		t.Fatalf("unexpected operational-run pCID %q", record.PCID)
	}

	reloaded, err := NewApp(root)
	if err != nil {
		t.Fatalf("reload app: %v", err)
	}
	reloadedRun, err := reloaded.GetRun(run.ID)
	if err != nil {
		t.Fatalf("get reloaded run: %v", err)
	}
	if reloadedRun.PlaceID != place.ID || len(reloadedRun.ResourceIDs) != 1 || reloadedRun.ResourceIDs[0] != resource.ID {
		t.Fatalf("unexpected replayed run context: %+v", reloadedRun)
	}
	if len(reloadedRun.ResponsibilityIDs) != 1 || reloadedRun.ResponsibilityIDs[0] != resp.ID {
		t.Fatalf("unexpected replayed run responsibility context: %+v", reloadedRun)
	}
}

func TestAppRejectsTamperedSignedOperationalRunRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if _, err := app.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil); err != nil {
		t.Fatalf("record run: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	recordPath := filepath.Join(root, "operational-run-messages.jsonl")
	recordBody, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read operational-run messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected at least one signed operational-run record")
	}
	var record SignedOperationalRunRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode operational-run record: %v", err)
	}
	record.EnvelopeCID = "bafkreioperationalruntampered"
	tamperedLine, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("encode tampered operational-run record: %v", err)
	}
	lines[0] = string(tamperedLine)
	if err := os.WriteFile(recordPath, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("rewrite tampered operational-run record: %v", err)
	}

	if _, err := NewApp(root); err == nil || !strings.Contains(err.Error(), "operational-run") || !strings.Contains(err.Error(), "envelope CAS cid mismatch") {
		t.Fatalf("expected operational-run authoritative CAS cid mismatch after tampering, got %v", err)
	}
}

func TestAppDualWritesCASObjectsForSignedFamiliesAndEvidenceBlobs(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "Dock bin", "Holds inbound parts", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	responsibility, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.RecordApproval("bob", "knowledge_item", item.ID, item.CurrentRevision, "reviewer", DecisionApproved, "ready"); err != nil {
		t.Fatalf("approve item: %v", err)
	}
	run, err := app.RecordRun("carol", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = app.AddEvidence("carol", run.ID, "Receiving inspection", map[string]string{"variance": "-2"}, "photo.txt", []byte("ok"))
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if run.Evidence[0].AttachmentCID == "" {
		t.Fatalf("expected attachment CID on evidence: %+v", run.Evidence[0])
	}
	if err := app.AddLink("alice", "resource", resource.ID, "knowledge_item", item.ID, "used_by", "Dock bin supports check"); err != nil {
		t.Fatalf("add link: %v", err)
	}

	checkSignedRecordCAS := func(path string, decode func([]byte) (string, string)) {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read signed records %s: %v", path, err)
		}
		lines := strings.Split(strings.TrimSpace(string(body)), "\n")
		if len(lines) == 0 {
			t.Fatalf("expected signed records in %s", path)
		}
		cid, envelopeBase64 := decode([]byte(lines[0]))
		envelopeBytes, err := base64.StdEncoding.DecodeString(envelopeBase64)
		if err != nil {
			t.Fatalf("decode envelope base64 for %s: %v", path, err)
		}
		casBytes, err := os.ReadFile(app.store.casObjectPath(cid))
		if err != nil {
			t.Fatalf("read cas object for %s: %v", cid, err)
		}
		if string(casBytes) != string(envelopeBytes) {
			t.Fatalf("cas envelope bytes mismatch for %s", cid)
		}
	}

	checkSignedRecordCAS(filepath.Join(root, "knowledge-item-messages.jsonl"), func(line []byte) (string, string) {
		var record SignedKnowledgeItemRecord
		if err := json.Unmarshal(line, &record); err != nil {
			t.Fatalf("decode knowledge-item record: %v", err)
		}
		return record.EnvelopeCID, record.EnvelopeBase64
	})
	checkSignedRecordCAS(filepath.Join(root, "knowledge-approval-messages.jsonl"), func(line []byte) (string, string) {
		var record SignedKnowledgeApprovalRecord
		if err := json.Unmarshal(line, &record); err != nil {
			t.Fatalf("decode knowledge-approval record: %v", err)
		}
		return record.EnvelopeCID, record.EnvelopeBase64
	})
	checkSignedRecordCAS(filepath.Join(root, "knowledge-evidence-messages.jsonl"), func(line []byte) (string, string) {
		var record SignedKnowledgeEvidenceRecord
		if err := json.Unmarshal(line, &record); err != nil {
			t.Fatalf("decode knowledge-evidence record: %v", err)
		}
		return record.EnvelopeCID, record.EnvelopeBase64
	})
	checkSignedRecordCAS(filepath.Join(root, "operational-run-messages.jsonl"), func(line []byte) (string, string) {
		var record SignedOperationalRunRecord
		if err := json.Unmarshal(line, &record); err != nil {
			t.Fatalf("decode operational-run record: %v", err)
		}
		return record.EnvelopeCID, record.EnvelopeBase64
	})
	checkSignedRecordCAS(filepath.Join(root, "operational-place-messages.jsonl"), func(line []byte) (string, string) {
		var record SignedOperationalPlaceRecord
		if err := json.Unmarshal(line, &record); err != nil {
			t.Fatalf("decode operational-place record: %v", err)
		}
		return record.EnvelopeCID, record.EnvelopeBase64
	})
	checkSignedRecordCAS(filepath.Join(root, "operational-resource-messages.jsonl"), func(line []byte) (string, string) {
		var record SignedOperationalResourceRecord
		if err := json.Unmarshal(line, &record); err != nil {
			t.Fatalf("decode operational-resource record: %v", err)
		}
		return record.EnvelopeCID, record.EnvelopeBase64
	})
	checkSignedRecordCAS(filepath.Join(root, "knowledge-link-messages.jsonl"), func(line []byte) (string, string) {
		var record SignedKnowledgeLinkRecord
		if err := json.Unmarshal(line, &record); err != nil {
			t.Fatalf("decode knowledge-link record: %v", err)
		}
		return record.EnvelopeCID, record.EnvelopeBase64
	})
	checkSignedRecordCAS(filepath.Join(root, "knowledge-responsibility-messages.jsonl"), func(line []byte) (string, string) {
		var record SignedKnowledgeResponsibilityRecord
		if err := json.Unmarshal(line, &record); err != nil {
			t.Fatalf("decode knowledge-responsibility record: %v", err)
		}
		return record.EnvelopeCID, record.EnvelopeBase64
	})

	attachmentBytes, err := os.ReadFile(run.Evidence[0].AttachmentPath)
	if err != nil {
		t.Fatalf("read compatibility attachment: %v", err)
	}
	casAttachmentBytes, err := os.ReadFile(app.store.casObjectPath(run.Evidence[0].AttachmentCID))
	if err != nil {
		t.Fatalf("read cas attachment: %v", err)
	}
	if string(attachmentBytes) != string(casAttachmentBytes) {
		t.Fatalf("cas attachment bytes mismatch")
	}
}

func TestAppExportsAndImportsPeerExchangeBootstrap(t *testing.T) {
	source, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}

	place, err := source.CreatePlace("alice", "area", "Receiving", "Inbound area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := source.CreateResource("alice", "container", "Dock bin", "Holds inbound parts", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	responsibility, err := source.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := source.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := source.RecordApproval("bob", "knowledge_item", item.ID, item.CurrentRevision, "reviewer", DecisionApproved, "ready"); err != nil {
		t.Fatalf("approve item: %v", err)
	}
	run, err := source.RecordRun("carol", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if err := source.RecordApproval("dave", "run", run.ID, 0, "auditor", DecisionNoted, "captured"); err != nil {
		t.Fatalf("approve run: %v", err)
	}
	if err := source.AddLink("alice", "responsibility", responsibility.ID, "knowledge_item", item.ID, "owns", "Receiving lead owns intake"); err != nil {
		t.Fatalf("add responsibility link: %v", err)
	}
	if err := source.AddLink("alice", "resource", resource.ID, "knowledge_item", item.ID, "used_by", "Dock bin supports check"); err != nil {
		t.Fatalf("add resource link: %v", err)
	}

	bundle, err := source.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}
	if bundle.Format != peerExchangeBundleFormat {
		t.Fatalf("unexpected bundle format %q", bundle.Format)
	}
	if len(bundle.KnowledgeApprovalRecords) != 2 {
		t.Fatalf("expected 2 approval records in bundle, got %d", len(bundle.KnowledgeApprovalRecords))
	}
	if len(bundle.OperationalRunRecords) != 1 {
		t.Fatalf("expected 1 operational-run record in bundle, got %d", len(bundle.OperationalRunRecords))
	}
	if len(bundle.OperationalPlaceRecords) != 1 {
		t.Fatalf("expected 1 operational-place record in bundle, got %d", len(bundle.OperationalPlaceRecords))
	}
	if len(bundle.OperationalResourceRecords) != 1 {
		t.Fatalf("expected 1 operational-resource record in bundle, got %d", len(bundle.OperationalResourceRecords))
	}
	if len(bundle.KnowledgeLinkRecords) != 2 {
		t.Fatalf("expected 2 link records in bundle, got %d", len(bundle.KnowledgeLinkRecords))
	}

	target, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	result, err := target.ImportPeerExchangeBundle(bundle)
	if err != nil {
		t.Fatalf("import bundle: %v", err)
	}
	if result.ImportedKnowledgeApprovals != 2 || result.ImportedOperationalRuns != 1 || result.ImportedOperationalPlaces != 1 || result.ImportedOperationalResources != 1 || result.ImportedKnowledgeLinks != 2 {
		t.Fatalf("unexpected import counts: %+v", result)
	}
	if len(result.UnresolvedReferences) != 0 {
		t.Fatalf("expected no unresolved references, got %+v", result.UnresolvedReferences)
	}
	if len(target.items) != 1 || len(target.responsibilities) != 1 {
		t.Fatalf("expected item and responsibility after import, got items=%d responsibilities=%d", len(target.items), len(target.responsibilities))
	}
	if len(target.runs) != 1 || len(target.places) != 1 || len(target.resources) != 1 {
		t.Fatalf("bootstrap import should materialize runs, places, and resources now: runs=%d places=%d resources=%d", len(target.runs), len(target.places), len(target.resources))
	}
	if len(target.approvals) != 2 {
		t.Fatalf("expected both approvals preserved after import, got %d", len(target.approvals))
	}
	if len(target.links) != 2 {
		t.Fatalf("expected both links preserved after import, got %d", len(target.links))
	}
	importedItem, err := target.GetKnowledgeItem(item.ID)
	if err != nil {
		t.Fatalf("get imported item: %v", err)
	}
	if len(importedItem.Approvals) != 1 || importedItem.Approvals[0].Decision != DecisionApproved {
		t.Fatalf("expected only the item approval to attach to imported item, got %+v", importedItem.Approvals)
	}
	if len(importedItem.Links) != 2 {
		t.Fatalf("expected both item-facing links to remain visible on imported item, got %+v", importedItem.Links)
	}
	importedResponsibility, err := target.GetResponsibility(responsibility.ID)
	if err != nil {
		t.Fatalf("get imported responsibility: %v", err)
	}
	if len(importedResponsibility.Links) != 1 {
		t.Fatalf("expected only the responsibility link to attach, got %+v", importedResponsibility.Links)
	}
	importedRun, err := target.GetRun(run.ID)
	if err != nil {
		t.Fatalf("get imported run: %v", err)
	}
	if importedRun.ItemID != item.ID || importedRun.Outcome != "accepted_with_notes" {
		t.Fatalf("unexpected imported run: %+v", importedRun)
	}
}

func TestAppPeerExchangeBundleCarriesOriginIdentity(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if _, err := app.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil); err != nil {
		t.Fatalf("record run: %v", err)
	}

	bundle, err := app.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}
	if bundle.ExportingPeerID == "" {
		t.Fatalf("expected exporting peer id")
	}
	for _, event := range bundle.Events {
		if event.OriginPeerID == "" || event.OriginSequence == 0 {
			t.Fatalf("expected origin identity on exported event: %+v", event)
		}
	}
	if len(bundle.KnowledgeItemRecords) == 0 || bundle.KnowledgeItemRecords[0].OriginPeerID == "" || bundle.KnowledgeItemRecords[0].OriginSequence == 0 {
		t.Fatalf("expected origin identity on exported knowledge-item record: %+v", bundle.KnowledgeItemRecords)
	}
	foundCanonicalCreate := false
	for _, event := range bundle.Events {
		if event.Type != "knowledge_item_created" {
			continue
		}
		foundCanonicalCreate = true
		if event.CanonicalID != bundle.KnowledgeItemRecords[0].EnvelopeCID {
			t.Fatalf("expected create event canonical id %q, got %+v", bundle.KnowledgeItemRecords[0].EnvelopeCID, event)
		}
		if event.DisplayID == "" || event.DisplayID == event.CanonicalID {
			t.Fatalf("expected create event display alias alongside canonical id, got %+v", event)
		}
	}
	if !foundCanonicalCreate {
		t.Fatalf("expected exported knowledge_item_created event")
	}
}

func TestAppImportsPeerExchangeIntoNonEmptyRuntimeAndDedupesByOrigin(t *testing.T) {
	source, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	if _, err := source.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil); err != nil {
		t.Fatalf("create source item: %v", err)
	}
	bundle, err := source.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}

	target, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	if _, err := target.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil); err != nil {
		t.Fatalf("seed target responsibility: %v", err)
	}
	result, err := target.ImportPeerExchangeBundle(bundle)
	if err != nil {
		t.Fatalf("import bundle into non-empty runtime: %v", err)
	}
	if result.ImportedKnowledgeItems != 1 || result.ImportedEvents == 0 {
		t.Fatalf("unexpected first import result: %+v", result)
	}
	second, err := target.ImportPeerExchangeBundle(bundle)
	if err != nil {
		t.Fatalf("re-import bundle: %v", err)
	}
	if second.ImportedEvents != 0 || second.ImportedKnowledgeItems != 0 {
		t.Fatalf("expected duplicate import to dedupe by origin, got %+v", second)
	}
}

func TestAppAllowsPeerExchangeAliasReuseAcrossPeers(t *testing.T) {
	source, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	sourceItem, err := source.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create source item: %v", err)
	}
	bundle, err := source.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}

	target, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	if _, err := target.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet again", "Local receiving check", "# Inspect inbound pallet again", nil, nil); err != nil {
		t.Fatalf("seed target item: %v", err)
	}
	result, err := target.ImportPeerExchangeBundle(bundle)
	if err != nil {
		t.Fatalf("expected alias reuse import success, got %v", err)
	}
	if result.ImportedKnowledgeItems != 1 {
		t.Fatalf("unexpected import result: %+v", result)
	}
	if len(target.items) != 2 {
		t.Fatalf("expected both canonical items after import, got %d", len(target.items))
	}
	importedItem, err := target.GetKnowledgeItem(sourceItem.ID)
	if err != nil {
		t.Fatalf("get imported canonical item: %v", err)
	}
	if importedItem.AliasID == "" {
		t.Fatalf("expected imported canonical item to preserve alias id: %+v", importedItem)
	}
}

func TestAppUsesCASAuthoritativelyForFrozenFamilyReplay(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	responsibility, err := app.CreateResponsibility("alice", "Receiver", "Owns receiving checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Check pallet", "Receiving check", "# Check pallet", nil, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.RecordApproval("bob", "knowledge_item", item.ID, item.CurrentRevision, "reviewer", DecisionApproved, "ready"); err != nil {
		t.Fatalf("approve item: %v", err)
	}
	if _, err := app.RecordRun("carol", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Wrap torn", "", "", "", nil, nil); err != nil {
		t.Fatalf("record run: %v", err)
	}
	if err := app.AddLink("alice", "responsibility", responsibility.ID, "knowledge_item", item.ID, "owns", "Receiver owns check"); err != nil {
		t.Fatalf("add link: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	tamperEnvelopeBase64 := func(path string, rewrite func([]byte) ([]byte, error)) {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read manifest %s: %v", path, err)
		}
		lines := strings.Split(strings.TrimSpace(string(body)), "\n")
		if len(lines) == 0 {
			t.Fatalf("expected manifest line in %s", path)
		}
		rewritten, err := rewrite([]byte(lines[0]))
		if err != nil {
			t.Fatalf("rewrite manifest %s: %v", path, err)
		}
		lines[0] = string(rewritten)
		if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
			t.Fatalf("write tampered manifest %s: %v", path, err)
		}
	}

	tamperEnvelopeBase64(filepath.Join(root, "knowledge-item-messages.jsonl"), func(line []byte) ([]byte, error) {
		var record SignedKnowledgeItemRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, err
		}
		record.EnvelopeBase64 = "tampered-item-envelope"
		return json.Marshal(record)
	})
	tamperEnvelopeBase64(filepath.Join(root, "knowledge-approval-messages.jsonl"), func(line []byte) ([]byte, error) {
		var record SignedKnowledgeApprovalRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, err
		}
		record.EnvelopeBase64 = "tampered-approval-envelope"
		return json.Marshal(record)
	})
	tamperEnvelopeBase64(filepath.Join(root, "operational-run-messages.jsonl"), func(line []byte) ([]byte, error) {
		var record SignedOperationalRunRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, err
		}
		record.EnvelopeBase64 = "tampered-operational-run-envelope"
		return json.Marshal(record)
	})
	tamperEnvelopeBase64(filepath.Join(root, "knowledge-link-messages.jsonl"), func(line []byte) ([]byte, error) {
		var record SignedKnowledgeLinkRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, err
		}
		record.EnvelopeBase64 = "tampered-link-envelope"
		return json.Marshal(record)
	})
	tamperEnvelopeBase64(filepath.Join(root, "knowledge-responsibility-messages.jsonl"), func(line []byte) ([]byte, error) {
		var record SignedKnowledgeResponsibilityRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, err
		}
		record.EnvelopeBase64 = "tampered-responsibility-envelope"
		return json.Marshal(record)
	})

	reloaded, err := NewApp(root)
	if err != nil {
		t.Fatalf("reload app with tampered manifests: %v", err)
	}
	records, err := reloaded.store.LoadSignedKnowledgeItemRecordsAuthoritative()
	if err != nil {
		t.Fatalf("load authoritative item records: %v", err)
	}
	if len(records) == 0 || records[0].EnvelopeBase64 == "tampered-item-envelope" {
		t.Fatalf("expected CAS-authoritative envelope bytes, got %+v", records)
	}
}

func TestAppBackfillsMissingCASEnvelopesFromManifestOnce(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "Startup flow", "# Start line", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	body, err := os.ReadFile(filepath.Join(root, "knowledge-item-messages.jsonl"))
	if err != nil {
		t.Fatalf("read knowledge-item manifest: %v", err)
	}
	var record SignedKnowledgeItemRecord
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(body))), &record); err != nil {
		t.Fatalf("decode knowledge-item manifest: %v", err)
	}
	casPath := filepath.Join(root, "cas", "objects", record.EnvelopeCID[:2], record.EnvelopeCID)
	if err := os.Remove(casPath); err != nil {
		t.Fatalf("remove cas object: %v", err)
	}

	reloaded, err := NewApp(root)
	if err != nil {
		t.Fatalf("reload app after removing cas object: %v", err)
	}
	reloadedItem, err := reloaded.GetKnowledgeItem(item.ID)
	if err != nil {
		t.Fatalf("get reloaded item: %v", err)
	}
	if reloadedItem.ID != item.ID {
		t.Fatalf("unexpected reloaded item %+v", reloadedItem)
	}
	casBytes, err := os.ReadFile(casPath)
	if err != nil {
		t.Fatalf("expected backfilled cas object: %v", err)
	}
	envelopeBytes, err := base64.StdEncoding.DecodeString(record.EnvelopeBase64)
	if err != nil {
		t.Fatalf("decode manifest envelope: %v", err)
	}
	if string(casBytes) != string(envelopeBytes) {
		t.Fatalf("unexpected backfilled cas bytes")
	}
}

func TestAppRejectsTamperedPeerExchangeBundle(t *testing.T) {
	source, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	responsibility, err := source.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := source.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := source.RecordApproval("bob", "knowledge_item", item.ID, item.CurrentRevision, "reviewer", DecisionApproved, "ready"); err != nil {
		t.Fatalf("approve item: %v", err)
	}
	bundle, err := source.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}
	if len(bundle.KnowledgeApprovalRecords) == 0 {
		t.Fatalf("expected approval records in bundle")
	}
	bundle.KnowledgeApprovalRecords[0].EnvelopeCID = "bafkreipeerexchangetampered"

	target, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	if _, err := target.ImportPeerExchangeBundle(bundle); err == nil || !strings.Contains(err.Error(), "knowledge-approval") || !strings.Contains(err.Error(), "envelope cid mismatch") {
		t.Fatalf("expected tampered approval envelope rejection, got %v", err)
	}
}

func TestAppExportsAndImportsPeerExchangeEvidenceBundle(t *testing.T) {
	source, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	item, err := source.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := source.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = source.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{"variance": "-2"}, "photo.txt", []byte("ok"))
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	sourceAttachmentPath := run.Evidence[0].AttachmentPath

	bundle, err := source.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}
	if len(bundle.KnowledgeEvidenceRecords) != 1 {
		t.Fatalf("expected 1 evidence record in bundle, got %d", len(bundle.KnowledgeEvidenceRecords))
	}
	if len(bundle.CASBlobObjects) != 1 {
		t.Fatalf("expected 1 evidence blob in bundle, got %d", len(bundle.CASBlobObjects))
	}

	target, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	result, err := target.ImportPeerExchangeBundle(bundle)
	if err != nil {
		t.Fatalf("import bundle: %v", err)
	}
	if result.ImportedKnowledgeEvidence != 1 || len(result.UnresolvedReferences) != 0 {
		t.Fatalf("unexpected evidence import result: %+v", result)
	}
	importedRun, err := target.GetRun(run.ID)
	if err != nil {
		t.Fatalf("get imported run: %v", err)
	}
	if len(importedRun.Evidence) != 1 {
		t.Fatalf("expected imported evidence, got %+v", importedRun.Evidence)
	}
	if importedRun.Evidence[0].AttachmentCID == "" || importedRun.Evidence[0].AttachmentPath == "" {
		t.Fatalf("expected imported evidence attachment refs, got %+v", importedRun.Evidence[0])
	}
	if importedRun.Evidence[0].AttachmentPath == sourceAttachmentPath {
		t.Fatalf("expected imported attachment path to be materialized locally, got source path %q", importedRun.Evidence[0].AttachmentPath)
	}
	body, err := os.ReadFile(importedRun.Evidence[0].AttachmentPath)
	if err != nil {
		t.Fatalf("read imported attachment: %v", err)
	}
	if string(body) != "ok" {
		t.Fatalf("unexpected imported attachment body %q", string(body))
	}
}

func TestAppRejectsPeerExchangeEvidenceBundleMissingBlob(t *testing.T) {
	source, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	item, err := source.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := source.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = source.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{"variance": "-2"}, "photo.txt", []byte("ok"))
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	bundle, err := source.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}
	delete(bundle.CASBlobObjects, run.Evidence[0].AttachmentCID)

	target, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	if _, err := target.ImportPeerExchangeBundle(bundle); err == nil || !strings.Contains(err.Error(), "knowledge-evidence blob") || !strings.Contains(err.Error(), "missing from bundle") {
		t.Fatalf("expected missing evidence blob rejection, got %v", err)
	}
}

func TestAppRejectsPeerExchangeEvidenceBundleTamperedBlob(t *testing.T) {
	source, err := NewApp(filepath.Join(t.TempDir(), "source"))
	if err != nil {
		t.Fatalf("new source app: %v", err)
	}
	item, err := source.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := source.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Outer wrap torn", "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	run, err = source.AddEvidence("bob", run.ID, "Receiving inspection", map[string]string{"variance": "-2"}, "photo.txt", []byte("ok"))
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	bundle, err := source.ExportPeerExchangeBundle()
	if err != nil {
		t.Fatalf("export bundle: %v", err)
	}
	bundle.CASBlobObjects[run.Evidence[0].AttachmentCID] = base64.StdEncoding.EncodeToString([]byte("tampered"))

	target, err := NewApp(filepath.Join(t.TempDir(), "target"))
	if err != nil {
		t.Fatalf("new target app: %v", err)
	}
	if _, err := target.ImportPeerExchangeBundle(bundle); err == nil || !strings.Contains(err.Error(), "knowledge-evidence blob") || !strings.Contains(err.Error(), "cid mismatch") {
		t.Fatalf("expected tampered evidence blob rejection, got %v", err)
	}
}

func TestAppWritesAndReplaysSignedKnowledgeLinkRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.AddLink("alice", "responsibility", resp.ID, "knowledge_item", item.ID, "owns", "Receiving lead owns intake"); err != nil {
		t.Fatalf("add link: %v", err)
	}

	meta := app.Meta()
	if meta.KnowledgeLinkPCID != protocols.KnowledgeLinkProfile.CID.String() {
		t.Fatalf("unexpected knowledge-link pCID in meta: got %q want %q", meta.KnowledgeLinkPCID, protocols.KnowledgeLinkProfile.CID.String())
	}

	recordBody, err := os.ReadFile(filepath.Join(root, "knowledge-link-messages.jsonl"))
	if err != nil {
		t.Fatalf("read knowledge-link messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 signed knowledge-link record, got %d", len(lines))
	}
	var record SignedKnowledgeLinkRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode link record: %v", err)
	}
	if record.PCID != protocols.KnowledgeLinkProfile.CID.String() {
		t.Fatalf("unexpected link record pCID: got %q want %q", record.PCID, protocols.KnowledgeLinkProfile.CID.String())
	}

	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	reloaded, err := NewApp(root)
	if err != nil {
		t.Fatalf("reload app: %v", err)
	}
	reloadedResp, err := reloaded.GetResponsibility(resp.ID)
	if err != nil {
		t.Fatalf("get reloaded responsibility: %v", err)
	}
	if len(reloadedResp.Links) != 1 {
		t.Fatalf("unexpected responsibility links after signed replay verification: %+v", reloadedResp.Links)
	}
	if reloadedResp.Links[0].Relation != "owns" || reloadedResp.Links[0].ToID != item.ID {
		t.Fatalf("unexpected reloaded responsibility link: %+v", reloadedResp.Links[0])
	}
}

func TestAppRejectsTamperedSignedKnowledgeLinkRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Inspect inbound pallet", "Receiving check", "# Inspect inbound pallet", nil, []string{resp.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := app.AddLink("alice", "responsibility", resp.ID, "knowledge_item", item.ID, "owns", "Receiving lead owns intake"); err != nil {
		t.Fatalf("add link: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	recordPath := filepath.Join(root, "knowledge-link-messages.jsonl")
	recordBody, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read knowledge-link messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected at least one signed knowledge-link record")
	}
	var record SignedKnowledgeLinkRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode link record: %v", err)
	}
	record.EnvelopeCID = "bafkreilinktampered"
	tamperedLine, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("encode tampered link record: %v", err)
	}
	lines[0] = string(tamperedLine)
	tampered := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(recordPath, []byte(tampered), 0o644); err != nil {
		t.Fatalf("rewrite tampered link record: %v", err)
	}

	if _, err := NewApp(root); err == nil || !strings.Contains(err.Error(), "knowledge-link") || !strings.Contains(err.Error(), "envelope CAS cid mismatch") {
		t.Fatalf("expected knowledge-link authoritative CAS cid mismatch after tampering, got %v", err)
	}
}

func TestAppWritesAndReplaysSignedKnowledgeResponsibilityRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	resp, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, []string{"receiving"})
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}

	meta := app.Meta()
	if meta.KnowledgeResponsibilityPCID != protocols.KnowledgeResponsibilityProfile.CID.String() {
		t.Fatalf("unexpected knowledge-responsibility pCID in meta: got %q want %q", meta.KnowledgeResponsibilityPCID, protocols.KnowledgeResponsibilityProfile.CID.String())
	}

	recordBody, err := os.ReadFile(filepath.Join(root, "knowledge-responsibility-messages.jsonl"))
	if err != nil {
		t.Fatalf("read knowledge-responsibility messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 signed knowledge-responsibility record, got %d", len(lines))
	}
	var record SignedKnowledgeResponsibilityRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode responsibility record: %v", err)
	}
	if record.PCID != protocols.KnowledgeResponsibilityProfile.CID.String() {
		t.Fatalf("unexpected responsibility record pCID: got %q want %q", record.PCID, protocols.KnowledgeResponsibilityProfile.CID.String())
	}
	if record.ResponsibilityID != resp.AliasID {
		t.Fatalf("unexpected responsibility alias in record: got %q want %q", record.ResponsibilityID, resp.AliasID)
	}

	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	reloaded, err := NewApp(root)
	if err != nil {
		t.Fatalf("reload app: %v", err)
	}
	reloadedResp, err := reloaded.GetResponsibility(resp.ID)
	if err != nil {
		t.Fatalf("get reloaded responsibility: %v", err)
	}
	if reloadedResp.Title != resp.Title || len(reloadedResp.LinkedRoleKeys) != 1 || reloadedResp.LinkedRoleKeys[0] != "reviewer" {
		t.Fatalf("unexpected responsibility after signed replay verification: %+v", reloadedResp)
	}
}

func TestAppRejectsTamperedSignedKnowledgeResponsibilityRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	if _, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake checks", []string{"reviewer"}, []string{"receiving"}); err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	recordPath := filepath.Join(root, "knowledge-responsibility-messages.jsonl")
	recordBody, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read knowledge-responsibility messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected at least one signed knowledge-responsibility record")
	}
	var record SignedKnowledgeResponsibilityRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode responsibility record: %v", err)
	}
	record.EnvelopeCID = "bafkreiresponsibilitytampered"
	tamperedLine, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("encode tampered responsibility record: %v", err)
	}
	lines[0] = string(tamperedLine)
	tampered := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(recordPath, []byte(tampered), 0o644); err != nil {
		t.Fatalf("rewrite tampered responsibility record: %v", err)
	}

	if _, err := NewApp(root); err == nil || !strings.Contains(err.Error(), "knowledge-responsibility") || !strings.Contains(err.Error(), "envelope CAS cid mismatch") {
		t.Fatalf("expected knowledge-responsibility authoritative CAS cid mismatch after tampering, got %v", err)
	}
}

func TestAppWritesAndReplaysSignedOperationalPlaceAndResourceRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound area", "", []string{"receiving"})
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "Dock bin", "Holds inbound parts", place.ID, []string{"parts"})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}

	meta := app.Meta()
	if meta.OperationalPlacePCID != protocols.OperationalPlaceProfile.CID.String() {
		t.Fatalf("unexpected operational-place pCID in meta: got %q want %q", meta.OperationalPlacePCID, protocols.OperationalPlaceProfile.CID.String())
	}
	if meta.OperationalResourcePCID != protocols.OperationalResourceProfile.CID.String() {
		t.Fatalf("unexpected operational-resource pCID in meta: got %q want %q", meta.OperationalResourcePCID, protocols.OperationalResourceProfile.CID.String())
	}

	placeBody, err := os.ReadFile(filepath.Join(root, "operational-place-messages.jsonl"))
	if err != nil {
		t.Fatalf("read operational-place messages: %v", err)
	}
	placeLines := strings.Split(strings.TrimSpace(string(placeBody)), "\n")
	if len(placeLines) != 1 {
		t.Fatalf("expected 1 signed operational-place record, got %d", len(placeLines))
	}
	var placeRecord SignedOperationalPlaceRecord
	if err := json.Unmarshal([]byte(placeLines[0]), &placeRecord); err != nil {
		t.Fatalf("decode place record: %v", err)
	}
	if placeRecord.PCID != protocols.OperationalPlaceProfile.CID.String() {
		t.Fatalf("unexpected place record pCID: got %q want %q", placeRecord.PCID, protocols.OperationalPlaceProfile.CID.String())
	}
	if placeRecord.PlaceID != place.AliasID {
		t.Fatalf("unexpected place alias in record: got %q want %q", placeRecord.PlaceID, place.AliasID)
	}

	resourceBody, err := os.ReadFile(filepath.Join(root, "operational-resource-messages.jsonl"))
	if err != nil {
		t.Fatalf("read operational-resource messages: %v", err)
	}
	resourceLines := strings.Split(strings.TrimSpace(string(resourceBody)), "\n")
	if len(resourceLines) != 1 {
		t.Fatalf("expected 1 signed operational-resource record, got %d", len(resourceLines))
	}
	var resourceRecord SignedOperationalResourceRecord
	if err := json.Unmarshal([]byte(resourceLines[0]), &resourceRecord); err != nil {
		t.Fatalf("decode resource record: %v", err)
	}
	if resourceRecord.PCID != protocols.OperationalResourceProfile.CID.String() {
		t.Fatalf("unexpected resource record pCID: got %q want %q", resourceRecord.PCID, protocols.OperationalResourceProfile.CID.String())
	}
	if resourceRecord.ResourceID != resource.AliasID {
		t.Fatalf("unexpected resource alias in record: got %q want %q", resourceRecord.ResourceID, resource.AliasID)
	}

	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	reloaded, err := NewApp(root)
	if err != nil {
		t.Fatalf("reload app: %v", err)
	}
	reloadedPlace, err := reloaded.GetPlace(place.ID)
	if err != nil {
		t.Fatalf("get reloaded place: %v", err)
	}
	if reloadedPlace.Name != place.Name || reloadedPlace.AliasID != place.AliasID {
		t.Fatalf("unexpected place after signed replay verification: %+v", reloadedPlace)
	}
	reloadedResource, err := reloaded.GetResource(resource.ID)
	if err != nil {
		t.Fatalf("get reloaded resource: %v", err)
	}
	if reloadedResource.PlaceID != place.ID || reloadedResource.AliasID != resource.AliasID {
		t.Fatalf("unexpected resource after signed replay verification: %+v", reloadedResource)
	}
}

func TestAppRejectsTamperedSignedOperationalPlaceRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	if _, err := app.CreatePlace("alice", "area", "Receiving", "Inbound area", "", nil); err != nil {
		t.Fatalf("create place: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	recordPath := filepath.Join(root, "operational-place-messages.jsonl")
	recordBody, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read operational-place messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected at least one signed operational-place record")
	}
	var record SignedOperationalPlaceRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode place record: %v", err)
	}
	record.EnvelopeCID = "bafkreioperationalplacetampered"
	tamperedLine, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("encode tampered place record: %v", err)
	}
	lines[0] = string(tamperedLine)
	if err := os.WriteFile(recordPath, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("rewrite tampered place record: %v", err)
	}

	if _, err := NewApp(root); err == nil || !strings.Contains(err.Error(), "operational-place") || !strings.Contains(err.Error(), "envelope CAS cid mismatch") {
		t.Fatalf("expected operational-place authoritative CAS cid mismatch after tampering, got %v", err)
	}
}

func TestAppRejectsTamperedSignedOperationalResourceRecords(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	if _, err := app.CreateResource("alice", "container", "Dock bin", "Holds inbound parts", place.ID, nil); err != nil {
		t.Fatalf("create resource: %v", err)
	}
	if err := app.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	recordPath := filepath.Join(root, "operational-resource-messages.jsonl")
	recordBody, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read operational-resource messages: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(recordBody)), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected at least one signed operational-resource record")
	}
	var record SignedOperationalResourceRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("decode resource record: %v", err)
	}
	record.EnvelopeCID = "bafkreioperationalresourcetampered"
	tamperedLine, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("encode tampered resource record: %v", err)
	}
	lines[0] = string(tamperedLine)
	if err := os.WriteFile(recordPath, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("rewrite tampered resource record: %v", err)
	}

	if _, err := NewApp(root); err == nil || !strings.Contains(err.Error(), "operational-resource") || !strings.Contains(err.Error(), "envelope CAS cid mismatch") {
		t.Fatalf("expected operational-resource authoritative CAS cid mismatch after tampering, got %v", err)
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

func TestAppLoadsDraftBodiesAuthoritativelyFromCAS(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	state, conflict, err := app.UpdateLiveItem(item.ID, "browser-a", "Alice", "#123456", 0, 0, true, 1, true, "# CAS body")
	if err != nil {
		t.Fatalf("update live draft: %v", err)
	}
	if conflict || state.Version != 2 {
		t.Fatalf("unexpected live update result: conflict=%v state=%+v", conflict, state)
	}
	manifestPath := filepath.Join(root, "drafts", item.ID+".json")
	body, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read draft manifest: %v", err)
	}
	var draft PersistedDraft
	if err := json.Unmarshal(body, &draft); err != nil {
		t.Fatalf("decode draft manifest: %v", err)
	}
	if draft.BodyCID == "" {
		t.Fatalf("expected body CID in draft manifest: %+v", draft)
	}
	draft.Body = "# tampered body"
	tamperedBody, err := json.Marshal(draft)
	if err != nil {
		t.Fatalf("encode tampered draft manifest: %v", err)
	}
	if err := os.WriteFile(manifestPath, tamperedBody, 0o644); err != nil {
		t.Fatalf("rewrite tampered draft manifest: %v", err)
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
	if reloadedItem.WorkingBody != "# CAS body" || reloadedItem.WorkingVersion != 2 {
		t.Fatalf("expected CAS-authoritative live draft after reload, got %+v", reloadedItem)
	}
}

func TestAppBackfillsLegacyDraftFilesIntoCASManifest(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	legacy := PersistedDraft{
		Body:      "# legacy draft",
		Version:   7,
		UpdatedAt: "2026-07-22T14:00:00Z",
	}
	legacyBody, err := json.Marshal(legacy)
	if err != nil {
		t.Fatalf("encode legacy draft: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "drafts", item.ID+".json"), legacyBody, 0o644); err != nil {
		t.Fatalf("write legacy draft file: %v", err)
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
	if reloadedItem.WorkingBody != "# legacy draft" || reloadedItem.WorkingVersion != 7 {
		t.Fatalf("expected legacy draft to reload through CAS backfill, got %+v", reloadedItem)
	}

	manifestPath := filepath.Join(root, "drafts", item.ID+".json")
	manifestBody, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read backfilled draft manifest: %v", err)
	}
	var manifest PersistedDraft
	if err := json.Unmarshal(manifestBody, &manifest); err != nil {
		t.Fatalf("decode backfilled draft manifest: %v", err)
	}
	if manifest.BodyCID == "" {
		t.Fatalf("expected backfilled body CID in draft manifest: %+v", manifest)
	}
	casBody, err := os.ReadFile(reloaded.store.casObjectPath(manifest.BodyCID))
	if err != nil {
		t.Fatalf("read backfilled draft CAS object: %v", err)
	}
	if string(casBody) != "# legacy draft" {
		t.Fatalf("unexpected backfilled draft CAS body %q", string(casBody))
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
	if meta.PeerExchangeFormat != peerExchangeBundleFormat || len(meta.PeerExchangeFamilies) != 8 || meta.OperationalPlacePCID == "" || meta.OperationalResourcePCID == "" {
		t.Fatalf("unexpected peer-exchange meta: %+v", meta)
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
