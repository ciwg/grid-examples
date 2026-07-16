package service

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type App struct {
	dataRoot        string
	store           *Store
	mu              sync.Mutex
	issues          map[string]*Issue
	nextIssueNumber int
	nextSequence    uint64
}

func NewApp(dataRoot string) (*App, error) {
	store, events, err := OpenStore(dataRoot)
	if err != nil {
		return nil, err
	}
	app := &App{
		dataRoot: dataRoot,
		store:    store,
		issues:   map[string]*Issue{},
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
		DataRoot:   app.dataRoot,
		Statuses:   append([]string(nil), statuses...),
		Severities: append([]string(nil), severities...),
		Identities: append([]Identity(nil), identities...),
		Team:       defaultTeam,
	}
}

func (app *App) ListIssues(status string, assignee string) ([]IssueSummary, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	if status != "" {
		if err := validateStatus(status); err != nil {
			return nil, err
		}
	}
	if err := validateAssignee(assignee); err != nil {
		return nil, err
	}
	result := []IssueSummary{}
	for _, issue := range app.issues {
		if status != "" && issue.Status != status {
			continue
		}
		if assignee != "" && issue.Assignee != assignee {
			continue
		}
		result = append(result, issue.Summary())
	}
	sort.Slice(result, func(i int, j int) bool {
		if result[i].UpdatedAt == result[j].UpdatedAt {
			return result[i].ID > result[j].ID
		}
		return result[i].UpdatedAt > result[j].UpdatedAt
	})
	return result, nil
}

func (app *App) GetIssue(issueID string) (Issue, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	issue, ok := app.issues[issueID]
	if !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	return cloneIssue(issue), nil
}

// Intent: Keep the first bug-tracker slice durable-first by persisting issue
// creation as an append-only event and projecting queue/detail state from that
// history instead of mutating a canonical row in place. Source: DI-nunit
func (app *App) CreateIssue(actor string, title string, description string, severity string) (Issue, error) {
	identity, err := validateIdentity(actor)
	if err != nil {
		return Issue{}, err
	}
	if identity.Role != RoleReporter {
		return Issue{}, fmt.Errorf("%s cannot create issues", actor)
	}
	if err := validateTitle(title); err != nil {
		return Issue{}, err
	}
	if err := validateDescription(description); err != nil {
		return Issue{}, err
	}
	if err := validateSeverity(severity); err != nil {
		return Issue{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	issueID := fmt.Sprintf("BUG-%04d", app.nextIssueNumber+1)
	event := IssueEvent{
		IssueID:        issueID,
		Type:           "created",
		Actor:          actor,
		Title:          strings.TrimSpace(title),
		Description:    normalizeLineEndings(description),
		Severity:       severity,
		Status:         StatusNew,
		Team:           defaultTeam,
		PreviousStatus: "",
	}
	if _, err := app.appendEventLocked(event); err != nil {
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

func (app *App) AddComment(actor string, issueID string, comment string) (Issue, error) {
	if _, err := validateIdentity(actor); err != nil {
		return Issue{}, err
	}
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	if err := validateComment(comment); err != nil {
		return Issue{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	if _, ok := app.issues[issueID]; !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	if _, err := app.appendEventLocked(IssueEvent{
		IssueID: issueID,
		Type:    "commented",
		Actor:   actor,
		Comment: normalizeLineEndings(comment),
	}); err != nil {
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

func (app *App) AssignIssue(actor string, issueID string, assignee string) (Issue, error) {
	identity, err := validateIdentity(actor)
	if err != nil {
		return Issue{}, err
	}
	if identity.Role != RoleTriage {
		return Issue{}, fmt.Errorf("%s cannot assign issues", actor)
	}
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	if err := validateAssignee(assignee); err != nil {
		return Issue{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	issue, ok := app.issues[issueID]
	if !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	// Intent: Keep assignment meaningful in the fixed bug-tracker workflow by
	// only allowing triage to assign issues that are ready for or already in
	// active work, instead of attaching owners to untouched or already resolved
	// issues. Source: DI-gitam
	if issue.Status == StatusNew || issue.Status == StatusResolved {
		return Issue{}, fmt.Errorf("issue %q cannot be assigned while %s", issueID, issue.Status)
	}
	if _, err := app.appendEventLocked(IssueEvent{
		IssueID:          issueID,
		Type:             "assigned",
		Actor:            actor,
		Assignee:         assignee,
		PreviousAssignee: issue.Assignee,
	}); err != nil {
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

// Intent: Keep the workflow legible by allowing only the locked v1 transitions
// and recording reopen as a normal history event that clears active ownership
// instead of mutating away the earlier resolved state. Source: DI-ninuf;
// DI-gofub
func (app *App) ChangeStatus(actor string, issueID string, status string) (Issue, error) {
	identity, err := validateIdentity(actor)
	if err != nil {
		return Issue{}, err
	}
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	if err := validateStatus(status); err != nil {
		return Issue{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	issue, ok := app.issues[issueID]
	if !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	if issue.Status == status {
		return Issue{}, fmt.Errorf("issue %q is already %s", issueID, status)
	}
	if err := validateTransition(identity.Role, issue.Status, status); err != nil {
		return Issue{}, err
	}
	// Intent: Treat the assignee as the active workflow owner for engineer work
	// so only the assigned engineer can start or resolve implementation work.
	// Source: DI-gitam
	if identity.Role == RoleEngineer && (status == StatusInProgress || status == StatusResolved) {
		if issue.Assignee == "" {
			return Issue{}, fmt.Errorf("issue %q is unassigned", issueID)
		}
		if issue.Assignee != actor {
			return Issue{}, fmt.Errorf("issue %q is assigned to %s", issueID, issue.Assignee)
		}
	}
	event := IssueEvent{
		IssueID:        issueID,
		Type:           "status_changed",
		Actor:          actor,
		Status:         status,
		PreviousStatus: issue.Status,
		Assignee:       issue.Assignee,
	}
	if issue.Status == StatusResolved && status == StatusTriaged {
		event.PreviousAssignee = issue.Assignee
		event.Assignee = ""
	}
	if _, err := app.appendEventLocked(event); err != nil {
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

func validateTransition(role string, from string, to string) error {
	switch {
	case role == RoleTriage && from == StatusNew && to == StatusTriaged:
		return nil
	case role == RoleEngineer && from == StatusTriaged && to == StatusInProgress:
		return nil
	case role == RoleEngineer && from == StatusInProgress && to == StatusResolved:
		return nil
	case (role == RoleReporter || role == RoleTriage) && from == StatusResolved && to == StatusTriaged:
		return nil
	default:
		return fmt.Errorf("status change %q -> %q is not allowed for %s", from, to, role)
	}
}

func (app *App) AddAttachment(actor string, issueID string, filename string, contentType string, bytes []byte) (Issue, error) {
	if _, err := validateIdentity(actor); err != nil {
		return Issue{}, err
	}
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	if len(bytes) == 0 {
		return Issue{}, fmt.Errorf("attachment is empty")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	if _, ok := app.issues[issueID]; !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	nextSequence := app.nextSequence + 1
	attachmentID := fmt.Sprintf("ATT-%06d", nextSequence)
	safeName := sanitizeAttachmentName(filename)
	relativePath := filepath.Join("attachments", issueID, fmt.Sprintf("%06d-%s", nextSequence, safeName))
	if err := app.store.WriteAttachment(relativePath, bytes); err != nil {
		return Issue{}, err
	}
	event := IssueEvent{
		IssueID:               issueID,
		Type:                  "attachment_added",
		Actor:                 actor,
		AttachmentID:          attachmentID,
		AttachmentName:        safeName,
		AttachmentPath:        relativePath,
		AttachmentContentType: contentType,
		AttachmentSize:        int64(len(bytes)),
	}
	if _, err := app.appendEventLocked(event); err != nil {
		if removeErr := app.store.RemoveAttachment(relativePath); removeErr != nil {
			return Issue{}, fmt.Errorf("%v (attachment cleanup failed: %v)", err, removeErr)
		}
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

func (app *App) DownloadAttachment(issueID string, attachmentID string) (AttachmentDownload, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	if err := validateIssueID(issueID); err != nil {
		return AttachmentDownload{}, err
	}
	issue, ok := app.issues[issueID]
	if !ok {
		return AttachmentDownload{}, fmt.Errorf("issue %q not found", issueID)
	}
	for _, event := range issue.Timeline {
		if event.Type == "attachment_added" && event.AttachmentID == attachmentID {
			bytes, err := app.store.ReadAttachment(event.AttachmentPath)
			if err != nil {
				return AttachmentDownload{}, err
			}
			return AttachmentDownload{
				Name:        event.AttachmentName,
				ContentType: event.AttachmentContentType,
				Bytes:       bytes,
			}, nil
		}
	}
	return AttachmentDownload{}, fmt.Errorf("attachment %q not found", attachmentID)
}

func (app *App) appendEventLocked(event IssueEvent) (IssueEvent, error) {
	app.nextSequence++
	event.Sequence = app.nextSequence
	if event.Timestamp == "" {
		event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if err := app.store.Append(event); err != nil {
		app.nextSequence--
		return IssueEvent{}, err
	}
	if err := app.applyEventLocked(event); err != nil {
		return IssueEvent{}, err
	}
	return event, nil
}

func (app *App) applyEventLocked(event IssueEvent) error {
	if event.Sequence > app.nextSequence {
		app.nextSequence = event.Sequence
	}
	if number := issueNumber(event.IssueID); number > app.nextIssueNumber {
		app.nextIssueNumber = number
	}
	switch event.Type {
	case "created":
		issue := &Issue{
			ID:          event.IssueID,
			Title:       event.Title,
			Description: event.Description,
			Severity:    event.Severity,
			Status:      event.Status,
			Reporter:    event.Actor,
			Assignee:    "",
			Team:        defaultTeam,
			CreatedAt:   event.Timestamp,
			UpdatedAt:   event.Timestamp,
			Timeline:    []IssueEvent{event},
		}
		if event.Team != "" {
			issue.Team = event.Team
		}
		app.issues[event.IssueID] = issue
		return nil
	default:
		issue, ok := app.issues[event.IssueID]
		if !ok {
			return fmt.Errorf("issue %q not found for event %q", event.IssueID, event.Type)
		}
		switch event.Type {
		case "commented":
		case "assigned":
			issue.Assignee = event.Assignee
		case "status_changed":
			issue.Status = event.Status
			issue.Assignee = event.Assignee
		case "attachment_added":
		default:
			return fmt.Errorf("unknown event type %q", event.Type)
		}
		issue.UpdatedAt = event.Timestamp
		issue.Timeline = append(issue.Timeline, event)
		return nil
	}
}

func issueNumber(issueID string) int {
	if !strings.HasPrefix(issueID, "BUG-") {
		return 0
	}
	number, err := strconv.Atoi(strings.TrimPrefix(issueID, "BUG-"))
	if err != nil {
		return 0
	}
	return number
}
