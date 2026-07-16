package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var issueIDPattern = regexp.MustCompile(`^BUG-\d{4}$`)

func validateIdentity(id string) (Identity, error) {
	identity, ok := IdentityByID(id)
	if !ok {
		return Identity{}, fmt.Errorf("unknown user %q", id)
	}
	return identity, nil
}

func validateIssueID(issueID string) error {
	if !issueIDPattern.MatchString(issueID) {
		return fmt.Errorf("invalid issue id %q", issueID)
	}
	return nil
}

func validateTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("title is required")
	}
	if len(title) > 120 {
		return errors.New("title is too long")
	}
	return nil
}

func validateDescription(description string) error {
	if strings.TrimSpace(description) == "" {
		return errors.New("description is required")
	}
	if len(description) > 16000 {
		return errors.New("description is too long")
	}
	return nil
}

func validateComment(comment string) error {
	if strings.TrimSpace(comment) == "" {
		return errors.New("comment is required")
	}
	if len(comment) > 8000 {
		return errors.New("comment is too long")
	}
	return nil
}

func validateSeverity(severity string) error {
	for _, candidate := range severities {
		if candidate == severity {
			return nil
		}
	}
	return fmt.Errorf("invalid severity %q", severity)
}

func validateStatus(status string) error {
	for _, candidate := range statuses {
		if candidate == status {
			return nil
		}
	}
	return fmt.Errorf("invalid status %q", status)
}

func validateAssignee(assignee string) error {
	if assignee == "" {
		return nil
	}
	identity, err := validateIdentity(assignee)
	if err != nil {
		return err
	}
	if identity.Role != RoleEngineer {
		return fmt.Errorf("assignee %q is not an engineer", assignee)
	}
	return nil
}

func sanitizeAttachmentName(name string) string {
	name = strings.TrimSpace(filepath.Base(name))
	name = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '.', r == '-', r == '_':
			return r
		default:
			return '-'
		}
	}, name)
	name = strings.Trim(name, "-")
	if name == "" {
		return "attachment.bin"
	}
	return name
}
