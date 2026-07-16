package service

import "strings"

const (
	StatusNew        = "New"
	StatusTriaged    = "Triaged"
	StatusInProgress = "In Progress"
	StatusResolved   = "Resolved"

	SeverityLow      = "Low"
	SeverityMedium   = "Medium"
	SeverityHigh     = "High"
	SeverityCritical = "Critical"

	RoleReporter = "reporter"
	RoleTriage   = "triage"
	RoleEngineer = "engineer"

	defaultTeam = "CORE"
)

var (
	statuses   = []string{StatusNew, StatusTriaged, StatusInProgress, StatusResolved}
	severities = []string{SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical}
	identities = []Identity{
		{ID: "reporter", DisplayName: "Reporter", Role: RoleReporter},
		{ID: "triage", DisplayName: "Triage", Role: RoleTriage},
		{ID: "engineer", DisplayName: "Engineer", Role: RoleEngineer},
	}
)

type Identity struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
}

type Meta struct {
	DataRoot   string     `json:"data_root"`
	Statuses   []string   `json:"statuses"`
	Severities []string   `json:"severities"`
	Identities []Identity `json:"identities"`
	Team       string     `json:"team"`
}

type IssueSummary struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
	Reporter    string `json:"reporter"`
	Assignee    string `json:"assignee"`
	Team        string `json:"team"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Issue struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Severity    string       `json:"severity"`
	Status      string       `json:"status"`
	Reporter    string       `json:"reporter"`
	Assignee    string       `json:"assignee"`
	Team        string       `json:"team"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
	Timeline    []IssueEvent `json:"timeline"`
}

type IssueEvent struct {
	Sequence              uint64 `json:"sequence"`
	IssueID               string `json:"issue_id"`
	Type                  string `json:"type"`
	Actor                 string `json:"actor"`
	Timestamp             string `json:"timestamp"`
	Title                 string `json:"title,omitempty"`
	Description           string `json:"description,omitempty"`
	Severity              string `json:"severity,omitempty"`
	Status                string `json:"status,omitempty"`
	PreviousStatus        string `json:"previous_status,omitempty"`
	Assignee              string `json:"assignee,omitempty"`
	PreviousAssignee      string `json:"previous_assignee,omitempty"`
	Comment               string `json:"comment,omitempty"`
	Team                  string `json:"team,omitempty"`
	AttachmentID          string `json:"attachment_id,omitempty"`
	AttachmentName        string `json:"attachment_name,omitempty"`
	AttachmentContentType string `json:"attachment_content_type,omitempty"`
	AttachmentSize        int64  `json:"attachment_size,omitempty"`
	AttachmentPath        string `json:"-"`
}

type AttachmentDownload struct {
	Name        string
	ContentType string
	Bytes       []byte
}

func cloneIssue(issue *Issue) Issue {
	if issue == nil {
		return Issue{}
	}
	clone := *issue
	clone.Timeline = append([]IssueEvent(nil), issue.Timeline...)
	return clone
}

func (issue *Issue) Summary() IssueSummary {
	return IssueSummary{
		ID:          issue.ID,
		Title:       issue.Title,
		Description: issue.Description,
		Severity:    issue.Severity,
		Status:      issue.Status,
		Reporter:    issue.Reporter,
		Assignee:    issue.Assignee,
		Team:        issue.Team,
		CreatedAt:   issue.CreatedAt,
		UpdatedAt:   issue.UpdatedAt,
	}
}

func IdentityByID(id string) (Identity, bool) {
	for _, identity := range identities {
		if identity.ID == id {
			return identity, true
		}
	}
	return Identity{}, false
}

func engineerIDs() []string {
	result := []string{}
	for _, identity := range identities {
		if identity.Role == RoleEngineer {
			result = append(result, identity.ID)
		}
	}
	return result
}

func normalizeLineEndings(value string) string {
	return strings.ReplaceAll(value, "\r\n", "\n")
}
