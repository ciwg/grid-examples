package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAnyMemberCanPlaceSafetyHoldButOnlyStewardClearsIt(t *testing.T) {
	app := NewDemoApp()
	if _, err := app.AddObservation("table-saw", "alice", "Blade guard is loose", true, nil); err != nil {
		t.Fatalf("place safety hold: %v", err)
	}
	if _, err := app.ClearSafetyHold("table-saw", "alice", "Looks safe"); err == nil {
		t.Fatal("unrecognized member cleared safety hold")
	}
	tool, err := app.ClearSafetyHold("table-saw", "carol", "Guard fastener tightened; available for qualified in-space use")
	if err != nil {
		t.Fatalf("recognized steward clears safety hold: %v", err)
	}
	if tool.SafetyHold {
		t.Fatal("tool remains on safety hold")
	}
}

func TestObservationAcceptsAnOptionalPhoto(t *testing.T) {
	app := NewDemoApp()
	tool, err := app.AddObservation("table-saw", "alice", "Scratch near fence", false, []Photo{{Name: "scratch.png", DataURL: "data:image/png;base64,aGVsbG8="}})
	if err != nil {
		t.Fatalf("record photo observation: %v", err)
	}
	if len(tool.Observations[0].Photos) != 1 {
		t.Fatal("photo was not preserved")
	}
}

func TestLoanPreservesAcceptedPolicyAndReturnEvidence(t *testing.T) {
	app := NewDemoApp()
	loaned, err := app.CreateLoan("cordless-drill", "alice", time.Now().Add(48*time.Hour))
	if err != nil {
		t.Fatalf("create loan: %v", err)
	}
	if loaned.ActiveLoan == nil || !loaned.ActiveLoan.TermsComplete || loaned.ActiveLoan.PolicyVersion == "" {
		t.Fatal("loan does not preserve policy evidence")
	}
	returned, err := app.ReturnLoan("cordless-drill", "alice", "Returned in expected condition; charger and case present")
	if err != nil {
		t.Fatalf("return loan: %v", err)
	}
	if returned.ActiveLoan != nil {
		t.Fatal("returned tool is still loaned")
	}
	if !strings.Contains(returned.Observations[len(returned.Observations)-1].Text, "return") {
		t.Fatal("return observation missing")
	}
}

func TestInSpaceOnlyToolCannotBeLoaned(t *testing.T) {
	app := NewDemoApp()
	_, err := app.CreateLoan("table-saw", "alice", time.Now().Add(24*time.Hour))
	if err == nil || !strings.Contains(err.Error(), "in-space") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestQualificationIsScopedToTheTool(t *testing.T) {
	app := NewDemoApp()
	state := app.State()
	if app.isQualifiedForTool("alice", state.Tools[0]) {
		t.Fatal("portable-power-tools qualification granted table saw access")
	}
	if !app.isQualifiedForTool("alice", state.Tools[1]) {
		t.Fatal("portable-power-tools qualification did not grant drill access")
	}
}

func TestPersistentAppReplaysEvidence(t *testing.T) {
	root := t.TempDir()
	app, err := NewPersistentDemoApp(root)
	if err != nil {
		t.Fatalf("create persistent app: %v", err)
	}
	if _, err := app.AddObservation("table-saw", "alice", "Fence alignment needs review", true, nil); err != nil {
		t.Fatalf("record durable observation: %v", err)
	}
	reloaded, err := NewPersistentDemoApp(root)
	if err != nil {
		t.Fatalf("reload persistent app: %v", err)
	}
	tool := reloaded.State().Tools[0]
	if !tool.SafetyHold || len(tool.Observations) != 1 {
		t.Fatalf("did not replay safety evidence: %+v", tool)
	}
}

func TestPersistentAppReplaysLargestAcceptedPhoto(t *testing.T) {
	root := t.TempDir()
	app, err := NewPersistentDemoApp(root)
	if err != nil {
		t.Fatalf("create persistent app: %v", err)
	}
	prefix := "data:image/png;base64,"
	photo := Photo{Name: "large.png", DataURL: prefix + strings.Repeat("a", maxPhotoDataURLBytes-len(prefix))}
	if _, err := app.AddObservation("table-saw", "alice", "Large photo evidence", false, []Photo{photo}); err != nil {
		t.Fatalf("record large photo: %v", err)
	}
	reloaded, err := NewPersistentDemoApp(root)
	if err != nil {
		t.Fatalf("replay large photo: %v", err)
	}
	if got := len(reloaded.State().Tools[0].Observations[0].Photos); got != 1 {
		t.Fatalf("replayed photos = %d, want 1", got)
	}
}

func TestFailedAppendDoesNotChangeLiveState(t *testing.T) {
	app := NewDemoApp()
	app.store = &Store{path: t.TempDir()}
	if _, err := app.AddObservation("table-saw", "alice", "Guard is loose", true, nil); err == nil {
		t.Fatal("recorded observation despite an unwritable event path")
	}
	tool := app.State().Tools[0]
	if tool.SafetyHold || len(tool.Observations) != 0 {
		t.Fatalf("failed append changed live state: %+v", tool)
	}
}

func TestPersistentAppFailsClosedOnMalformedEvidence(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "events.jsonl")
	if err := os.WriteFile(path, []byte(`{"type":"observation"`), 0o600); err != nil {
		t.Fatalf("write malformed evidence: %v", err)
	}
	if _, err := NewPersistentDemoApp(root); err == nil {
		t.Fatal("started with malformed evidence")
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read preserved evidence: %v", err)
	}
	if string(contents) != `{"type":"observation"` {
		t.Fatalf("malformed evidence changed to %q", contents)
	}
}

func TestLoanReplayPreservesAcceptedAreaPolicy(t *testing.T) {
	root := t.TempDir()
	app, err := NewPersistentDemoApp(root)
	if err != nil {
		t.Fatalf("create persistent app: %v", err)
	}
	app.state.Areas[0].PolicyVersion = "woodworking-lending/v2"
	app.state.Areas[0].Policy = "Return the drill with its charger and case."
	if _, err := app.CreateLoan("cordless-drill", "alice", time.Now().Add(24*time.Hour)); err != nil {
		t.Fatalf("create loan: %v", err)
	}
	reloaded, err := NewPersistentDemoApp(root)
	if err != nil {
		t.Fatalf("replay loan: %v", err)
	}
	loan := reloaded.State().Tools[1].ActiveLoan
	if loan == nil {
		t.Fatal("replayed loan is missing")
	}
	if loan.PolicyVersion != "woodworking-lending/v2" || loan.Policy != "Return the drill with its charger and case." {
		t.Fatalf("replayed policy = %q / %q", loan.PolicyVersion, loan.Policy)
	}
}

func TestPersistentAppReplaysLegacyLoanWithIncompleteTerms(t *testing.T) {
	root := t.TempDir()
	store, err := NewStore(root)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	dueAt := time.Now().Add(24 * time.Hour).UTC()
	if err := store.Append(Event{Type: "loan", ToolID: "cordless-drill", ActorID: "alice", DueAt: dueAt, CreatedAt: time.Now().UTC()}); err != nil {
		t.Fatalf("append legacy loan: %v", err)
	}
	app, err := NewPersistentDemoApp(root)
	if err != nil {
		t.Fatalf("replay legacy loan: %v", err)
	}
	loan := app.State().Tools[1].ActiveLoan
	if loan == nil || loan.TermsComplete {
		t.Fatalf("legacy loan terms completeness = %+v", loan)
	}
	if loan.PolicyVersion != "" || loan.Policy != "" || !loan.DueAt.Equal(dueAt) {
		t.Fatalf("legacy loan evidence = %+v", loan)
	}
}
