package service

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	maxDocumentIDLen    = 128
	maxParticipantIDLen = 128
	maxEmbodimentLen    = 32
	maxDisplayNameLen   = 80
	maxChangeBytesLen   = 1 << 20
	defaultFeedLimit    = 256
	maxFeedLimit        = 512
)

var (
	participantIDPattern = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)
	embodimentPattern    = regexp.MustCompile(`^[A-Za-z0-9._-]*$`)
	colorPattern         = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
)

// Intent: Keep local relay inputs bounded and structurally predictable so one
// malformed client request cannot collapse participant identity, explode feed
// size, or push arbitrarily large payloads into the signed relay path. Source:
// DI-rabod
func validateDocumentID(documentID string) error {
	if documentID == "" {
		return fmt.Errorf("document id is required")
	}
	if utf8.RuneCountInString(documentID) > maxDocumentIDLen {
		return fmt.Errorf("document id too long")
	}
	if strings.ContainsAny(documentID, `/\`) {
		return fmt.Errorf("document id must not contain path separators")
	}
	for _, value := range documentID {
		if value < 0x20 || value == 0x7f {
			return fmt.Errorf("document id must not contain control characters")
		}
	}
	return nil
}

func validateParticipantID(participantID string) error {
	if participantID == "" {
		return fmt.Errorf("participant id is required")
	}
	if len(participantID) > maxParticipantIDLen {
		return fmt.Errorf("participant id too long")
	}
	if !participantIDPattern.MatchString(participantID) {
		return fmt.Errorf("participant id contains unsupported characters")
	}
	return nil
}

func validateRecipientID(recipientID string) error {
	if recipientID == "" {
		return nil
	}
	if len(recipientID) > maxParticipantIDLen {
		return fmt.Errorf("recipient id too long")
	}
	if !participantIDPattern.MatchString(recipientID) {
		return fmt.Errorf("recipient id contains unsupported characters")
	}
	return nil
}

func validateEmbodiment(embodiment string) error {
	if len(embodiment) > maxEmbodimentLen {
		return fmt.Errorf("embodiment too long")
	}
	if !embodimentPattern.MatchString(embodiment) {
		return fmt.Errorf("embodiment contains unsupported characters")
	}
	return nil
}

func validateDisplayName(displayName string) error {
	if displayName == "" {
		return fmt.Errorf("display name is required")
	}
	if utf8.RuneCountInString(displayName) > maxDisplayNameLen {
		return fmt.Errorf("display name too long")
	}
	for _, value := range displayName {
		if value < 0x20 || value == 0x7f {
			return fmt.Errorf("display name must not contain control characters")
		}
	}
	return nil
}

func validateColor(color string) error {
	if !colorPattern.MatchString(color) {
		return fmt.Errorf("color must be a #RRGGBB value")
	}
	return nil
}

func validateCursorValue(name string, value int) error {
	if value < 0 {
		return fmt.Errorf("%s must be non-negative", name)
	}
	return nil
}

func validateChangeBytes(messageBytes []byte) error {
	if len(messageBytes) == 0 {
		return fmt.Errorf("change bytes are required")
	}
	if len(messageBytes) > maxChangeBytesLen {
		return fmt.Errorf("change bytes too large")
	}
	return nil
}

func clampFeedLimit(limit int) int {
	if limit <= 0 {
		return defaultFeedLimit
	}
	if limit > maxFeedLimit {
		return maxFeedLimit
	}
	return limit
}
