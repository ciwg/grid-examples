package records

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Envelope is the ex6 durable carriage unit for package-defined record families.
type Envelope struct {
	Family          string          `json:"family"`
	ProtocolPCID    string          `json:"protocol_pcid"`
	RecordID        string          `json:"record_id"`
	Signer          string          `json:"signer"`
	Timestamp       string          `json:"timestamp"`
	Payload         json.RawMessage `json:"payload"`
	AuthorKeyID     string          `json:"author_key_id,omitempty"`
	AuthorPublicKey string          `json:"author_public_key,omitempty"`
	AuthorSignature string          `json:"author_signature,omitempty"`
}

func (envelope Envelope) Validate() error {
	if strings.TrimSpace(envelope.Family) == "" {
		return errors.New("family is required")
	}
	if strings.TrimSpace(envelope.ProtocolPCID) == "" {
		return errors.New("protocol_pcid is required")
	}
	if strings.TrimSpace(envelope.RecordID) == "" {
		return errors.New("record_id is required")
	}
	if strings.TrimSpace(envelope.Timestamp) == "" {
		return errors.New("timestamp is required")
	}
	if _, err := time.Parse(time.RFC3339, envelope.Timestamp); err != nil {
		return fmt.Errorf("timestamp must be RFC3339: %w", err)
	}
	if len(envelope.Payload) == 0 {
		return errors.New("payload is required")
	}
	if envelope.HasAuthorSignature() {
		if strings.TrimSpace(envelope.AuthorKeyID) == "" {
			return errors.New("author_key_id is required when author signature is present")
		}
		if strings.TrimSpace(envelope.AuthorPublicKey) == "" {
			return errors.New("author_public_key is required when author signature is present")
		}
		if strings.TrimSpace(envelope.AuthorSignature) == "" {
			return errors.New("author_signature is required when author signature is present")
		}
	}
	return nil
}

func Parse(raw []byte) (Envelope, error) {
	var envelope Envelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return Envelope{}, err
	}
	if err := envelope.Validate(); err != nil {
		return Envelope{}, err
	}
	return envelope, nil
}

func MustMarshal(envelope Envelope) []byte {
	body, err := json.Marshal(envelope)
	if err != nil {
		panic(err)
	}
	return body
}

func (envelope Envelope) HasAuthorSignature() bool {
	return strings.TrimSpace(envelope.AuthorKeyID) != "" ||
		strings.TrimSpace(envelope.AuthorPublicKey) != "" ||
		strings.TrimSpace(envelope.AuthorSignature) != ""
}

func (envelope Envelope) SigningBytes() ([]byte, error) {
	signable := struct {
		Family       string          `json:"family"`
		ProtocolPCID string          `json:"protocol_pcid"`
		RecordID     string          `json:"record_id"`
		Signer       string          `json:"signer"`
		Timestamp    string          `json:"timestamp"`
		Payload      json.RawMessage `json:"payload"`
	}{
		Family:       envelope.Family,
		ProtocolPCID: envelope.ProtocolPCID,
		RecordID:     envelope.RecordID,
		Signer:       envelope.Signer,
		Timestamp:    envelope.Timestamp,
		Payload:      envelope.Payload,
	}
	return json.Marshal(signable)
}
