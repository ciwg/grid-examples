package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	observationPCID = "bafkreifhodcald6kzib36rzeji27hnqjdkeycibnkcigcsz7mzejz6obiy"
	safetyPCID      = "bafkreigt3p2l4uel7wmjr4kple7o55ymchlhh43gajjwsgaeifoogeztc4"
	loanPCID        = "bafkreibgbky6hbswdimkngjii5zglgvfkolxjxuonbmnqwcdjt4v2qyliq"
	returnPCID      = "bafkreifjgcfwkbwzrgmtj2wfhk3kpzmjbj3wjxid52jau5rrpthabr3ehi"
)

type observationPayload struct {
	ToolID      string           `json:"tool_id"`
	Observation string           `json:"observation"`
	Photos      []photoReference `json:"photos,omitempty"`
}

type photoReference struct {
	BlobCID   string `json:"blob_cid"`
	MediaType string `json:"media_type"`
	Name      string `json:"name"`
}

type safetyPayload struct {
	ToolID      string           `json:"tool_id"`
	Assessment  string           `json:"assessment"`
	Disposition string           `json:"disposition"`
	BasisRecord string           `json:"basis_record_id,omitempty"`
	Photos      []photoReference `json:"photos,omitempty"`
}

type loanPayload struct {
	ToolID        string `json:"tool_id"`
	BorrowerID    string `json:"borrower_id"`
	DueAt         string `json:"due_at"`
	PolicyVersion string `json:"policy_version"`
	Policy        string `json:"policy"`
}

type returnPayload struct {
	ToolID       string           `json:"tool_id"`
	LoanRecordID string           `json:"loan_record_id"`
	Condition    string           `json:"condition"`
	Photos       []photoReference `json:"photos,omitempty"`
}

// validateKnownRecord checks only semantics selected by a known Ex7 family.
// Intent: unknown records remain retainable while known malformed semantics do
// not become durable projectable evidence. Source: DI-tohak; DI-piruf.
func validateKnownRecord(record Record) error {
	switch record.Protocol {
	case observationPCID:
		var payload observationPayload
		if err := decodePayload(record.Payload, &payload); err != nil {
			return err
		}
		if payload.ToolID == "" || payload.Observation == "" {
			return errors.New("observation payload requires tool_id and observation")
		}
		if err := validatePhotoReferences(payload.Photos); err != nil {
			return err
		}
	case safetyPCID:
		var payload safetyPayload
		if err := decodePayload(record.Payload, &payload); err != nil {
			return err
		}
		if payload.ToolID == "" || payload.Assessment == "" || (payload.Disposition != "hold" && payload.Disposition != "clear") {
			return errors.New("invalid safety-disposition payload")
		}
		if err := validatePhotoReferences(payload.Photos); err != nil {
			return err
		}
	case loanPCID:
		var payload loanPayload
		if err := decodePayload(record.Payload, &payload); err != nil {
			return err
		}
		if payload.ToolID == "" || payload.BorrowerID == "" || payload.BorrowerID != record.Signer || payload.PolicyVersion == "" || payload.Policy == "" {
			return errors.New("invalid off-site-loan payload")
		}
		if _, err := time.Parse(time.RFC3339, payload.DueAt); err != nil {
			return fmt.Errorf("invalid loan due_at: %w", err)
		}
	case returnPCID:
		var payload returnPayload
		if err := decodePayload(record.Payload, &payload); err != nil {
			return err
		}
		if payload.ToolID == "" || payload.LoanRecordID == "" || payload.Condition == "" {
			return errors.New("invalid off-site-return payload")
		}
		if err := validatePhotoReferences(payload.Photos); err != nil {
			return err
		}
	}
	return nil
}

func validatePhotoReferences(photos []photoReference) error {
	for _, photo := range photos {
		if photo.BlobCID == "" || photo.MediaType == "" || photo.Name == "" {
			return errors.New("invalid photo reference")
		}
	}
	return nil
}

func decodePayload(raw []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("payload contains trailing JSON")
	}
	return nil
}
