package service

import (
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
	if loaned.ActiveLoan == nil || loaned.ActiveLoan.PolicyVersion == "" {
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
