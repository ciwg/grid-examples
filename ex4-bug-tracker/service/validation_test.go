package service

import (
	"path/filepath"
	"testing"
)

func TestValidateIssueIDAllowsFiveDigits(t *testing.T) {
	t.Parallel()
	if err := validateIssueID("BUG-10000"); err != nil {
		t.Fatalf("validateIssueID rejected five-digit id: %v", err)
	}
}

func TestCreateIssueBeyondFourDigits(t *testing.T) {
	t.Parallel()
	app, err := NewApp(filepath.Join(t.TempDir(), ".bug-tracker"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	app.nextIssueNumber = 9999
	issue, err := app.CreateIssue("reporter", "Scale rollover", "Verify five-digit issue IDs remain usable.", SeverityLow)
	if err != nil {
		t.Fatalf("create issue: %v", err)
	}
	if issue.ID != "BUG-10000" {
		t.Fatalf("issue id = %q, want BUG-10000", issue.ID)
	}
	if err := validateIssueID(issue.ID); err != nil {
		t.Fatalf("validateIssueID(%q): %v", issue.ID, err)
	}
}
