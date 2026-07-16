package service_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/service"
)

func TestCreateIssueDefaults(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)

	issue, err := app.CreateIssue("reporter", "Upload crash", "Uploading a log file crashes the page.", service.SeverityHigh)
	if err != nil {
		t.Fatalf("create issue: %v", err)
	}
	if issue.ID != "BUG-0001" {
		t.Fatalf("issue id = %q, want BUG-0001", issue.ID)
	}
	if issue.Status != service.StatusNew {
		t.Fatalf("status = %q, want %q", issue.Status, service.StatusNew)
	}
	if issue.Team != "CORE" {
		t.Fatalf("team = %q, want CORE", issue.Team)
	}
	if len(issue.Timeline) != 1 || issue.Timeline[0].Type != "created" {
		t.Fatalf("timeline = %#v, want single created event", issue.Timeline)
	}
}

func TestIssueWorkflowAndReopenClearsAssignee(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	issue := mustCreateIssue(t, app)

	if _, err := app.ChangeStatus("triage", issue.ID, service.StatusTriaged); err != nil {
		t.Fatalf("triage status: %v", err)
	}
	if _, err := app.AssignIssue("triage", issue.ID, "engineer"); err != nil {
		t.Fatalf("assign issue: %v", err)
	}
	if _, err := app.ChangeStatus("engineer", issue.ID, service.StatusInProgress); err != nil {
		t.Fatalf("start issue: %v", err)
	}
	resolved, err := app.ChangeStatus("engineer", issue.ID, service.StatusResolved)
	if err != nil {
		t.Fatalf("resolve issue: %v", err)
	}
	if resolved.Assignee != "engineer" {
		t.Fatalf("resolved assignee = %q, want engineer", resolved.Assignee)
	}
	reopened, err := app.ChangeStatus("reporter", issue.ID, service.StatusTriaged)
	if err != nil {
		t.Fatalf("reopen issue: %v", err)
	}
	if reopened.Assignee != "" {
		t.Fatalf("reopened assignee = %q, want empty", reopened.Assignee)
	}
	if reopened.Status != service.StatusTriaged {
		t.Fatalf("reopened status = %q, want %q", reopened.Status, service.StatusTriaged)
	}
}

func TestAttachmentRoundTrip(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	issue := mustCreateIssue(t, app)

	updated, err := app.AddAttachment("reporter", issue.ID, "trace.log", "text/plain", []byte("stack trace"))
	if err != nil {
		t.Fatalf("add attachment: %v", err)
	}
	event := updated.Timeline[len(updated.Timeline)-1]
	if event.Type != "attachment_added" {
		t.Fatalf("last event type = %q, want attachment_added", event.Type)
	}
	download, err := app.DownloadAttachment(issue.ID, event.AttachmentID)
	if err != nil {
		t.Fatalf("download attachment: %v", err)
	}
	if download.Name != "trace.log" {
		t.Fatalf("download name = %q, want trace.log", download.Name)
	}
	if !bytes.Equal(download.Bytes, []byte("stack trace")) {
		t.Fatalf("download bytes = %q, want stack trace", string(download.Bytes))
	}
}

func TestListIssuesFilters(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	first := mustCreateIssue(t, app)
	second, err := app.CreateIssue("reporter", "Login glitch", "The login form loses focus.", service.SeverityLow)
	if err != nil {
		t.Fatalf("create second issue: %v", err)
	}
	if _, err := app.ChangeStatus("triage", first.ID, service.StatusTriaged); err != nil {
		t.Fatalf("triage first issue: %v", err)
	}
	if _, err := app.AssignIssue("triage", first.ID, "engineer"); err != nil {
		t.Fatalf("assign first issue: %v", err)
	}
	assigned, err := app.ListIssues("", "engineer")
	if err != nil {
		t.Fatalf("list assigned issues: %v", err)
	}
	if len(assigned) != 1 || assigned[0].ID != first.ID {
		t.Fatalf("assigned issues = %#v, want only %s", assigned, first.ID)
	}
	newIssues, err := app.ListIssues(service.StatusNew, "")
	if err != nil {
		t.Fatalf("list new issues: %v", err)
	}
	if len(newIssues) != 1 || newIssues[0].ID != second.ID {
		t.Fatalf("new issues = %#v, want only %s", newIssues, second.ID)
	}
}

func TestReporterCannotResolveIssue(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	issue := mustCreateIssue(t, app)
	if _, err := app.ChangeStatus("reporter", issue.ID, service.StatusResolved); err == nil {
		t.Fatalf("reporter resolved issue without error")
	}
}

func TestEngineerCannotStartUnassignedIssue(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	issue := mustCreateIssue(t, app)
	if _, err := app.ChangeStatus("triage", issue.ID, service.StatusTriaged); err != nil {
		t.Fatalf("triage issue: %v", err)
	}
	if _, err := app.ChangeStatus("engineer", issue.ID, service.StatusInProgress); err == nil {
		t.Fatalf("engineer started unassigned issue without error")
	}
}

func TestCannotAssignNewOrResolvedIssue(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	issue := mustCreateIssue(t, app)
	if _, err := app.AssignIssue("triage", issue.ID, "engineer"); err == nil {
		t.Fatalf("assigned new issue without error")
	}
	if _, err := app.ChangeStatus("triage", issue.ID, service.StatusTriaged); err != nil {
		t.Fatalf("triage issue: %v", err)
	}
	if _, err := app.AssignIssue("triage", issue.ID, "engineer"); err != nil {
		t.Fatalf("assign triaged issue: %v", err)
	}
	if _, err := app.ChangeStatus("engineer", issue.ID, service.StatusInProgress); err != nil {
		t.Fatalf("start issue: %v", err)
	}
	if _, err := app.ChangeStatus("engineer", issue.ID, service.StatusResolved); err != nil {
		t.Fatalf("resolve issue: %v", err)
	}
	if _, err := app.AssignIssue("triage", issue.ID, "engineer"); err == nil {
		t.Fatalf("assigned resolved issue without error")
	}
}

func TestSeedDemoIfEmpty(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)

	seeded, err := app.SeedDemoIfEmpty()
	if err != nil {
		t.Fatalf("seed demo: %v", err)
	}
	if !seeded {
		t.Fatalf("seeded = false, want true")
	}
	issues, err := app.ListIssues("", "")
	if err != nil {
		t.Fatalf("list issues: %v", err)
	}
	if len(issues) != 4 {
		t.Fatalf("issue count = %d, want 4", len(issues))
	}
	foundResolved := false
	foundReopened := false
	foundAttachment := false
	for _, summary := range issues {
		detail, err := app.GetIssue(summary.ID)
		if err != nil {
			t.Fatalf("get issue %s: %v", summary.ID, err)
		}
		if detail.Status == service.StatusResolved {
			foundResolved = true
		}
		for _, event := range detail.Timeline {
			if event.Type == "status_changed" && event.PreviousStatus == service.StatusResolved && event.Status == service.StatusTriaged {
				foundReopened = true
			}
			if event.Type == "attachment_added" {
				foundAttachment = true
			}
		}
	}
	if !foundResolved {
		t.Fatalf("expected one resolved demo issue")
	}
	if !foundReopened {
		t.Fatalf("expected one reopened demo issue")
	}
	if !foundAttachment {
		t.Fatalf("expected one demo attachment")
	}
}

func TestSeedDemoIfEmptyDoesNotDuplicate(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)

	seeded, err := app.SeedDemoIfEmpty()
	if err != nil {
		t.Fatalf("first seed: %v", err)
	}
	if !seeded {
		t.Fatalf("first seeded = false, want true")
	}
	seeded, err = app.SeedDemoIfEmpty()
	if err != nil {
		t.Fatalf("second seed: %v", err)
	}
	if seeded {
		t.Fatalf("second seeded = true, want false")
	}
	issues, err := app.ListIssues("", "")
	if err != nil {
		t.Fatalf("list issues: %v", err)
	}
	if len(issues) != 4 {
		t.Fatalf("issue count after second seed = %d, want 4", len(issues))
	}
}

func newTestApp(t *testing.T) *service.App {
	t.Helper()
	app, err := service.NewApp(filepath.Join(t.TempDir(), ".bug-tracker"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	return app
}

func mustCreateIssue(t *testing.T, app *service.App) service.Issue {
	t.Helper()
	issue, err := app.CreateIssue("reporter", "Upload crash", "Uploading a log file crashes the page.", service.SeverityHigh)
	if err != nil {
		t.Fatalf("create issue: %v", err)
	}
	return issue
}
