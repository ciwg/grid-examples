package records

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
)

// Envelope is the canonical Grid carriage for a package-defined durable record.
// Family payload semantics remain package-defined; this type owns only common
// evidence and carriage slots. Source: DI-sidoh
type Envelope struct {
	Family          string
	ProtocolPCID    string
	RecordID        string
	Signer          string
	Timestamp       string
	Payload         json.RawMessage
	AuthorKeyID     string
	AuthorPublicKey string
	AuthorSignature string
}

func (envelope Envelope) Validate() error {
	if strings.TrimSpace(envelope.Family) == "" {
		return errors.New("family is required")
	}
	if strings.TrimSpace(envelope.ProtocolPCID) == "" {
		return errors.New("protocol_pcid is required")
	}
	if _, err := cid.Decode(envelope.ProtocolPCID); err != nil {
		return fmt.Errorf("protocol_pcid must be a valid CID: %w", err)
	}
	if expected := PackageProtocolPCID(envelope.Family); expected != "" && envelope.ProtocolPCID != expected {
		return fmt.Errorf("package record protocol_pcid does not match family %s", envelope.Family)
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
	if len(envelope.Payload) == 0 || !json.Valid(envelope.Payload) {
		return errors.New("payload must be valid JSON")
	}
	canonical, err := CanonicalJSON(envelope.Payload)
	if err != nil {
		return err
	}
	if !bytes.Equal(canonical, envelope.Payload) {
		return errors.New("payload must be canonical JSON")
	}
	if envelope.HasAuthorSignature() && (strings.TrimSpace(envelope.AuthorKeyID) == "" || strings.TrimSpace(envelope.AuthorPublicKey) == "" || strings.TrimSpace(envelope.AuthorSignature) == "") {
		return errors.New("complete author signature fields are required")
	}
	return nil
}

func Parse(raw []byte) (Envelope, error) {
	gridEnvelope, err := DecodeGrid(raw)
	if err != nil {
		return Envelope{}, err
	}
	if len(gridEnvelope.Slots) != 8 {
		return Envelope{}, errors.New("grid envelope has invalid package record slots")
	}
	family, ok := gridEnvelope.Slots[0].(string)
	if !ok {
		return Envelope{}, errors.New("package record family must be text")
	}
	recordID, ok := gridEnvelope.Slots[1].(string)
	if !ok {
		return Envelope{}, errors.New("package record ID must be text")
	}
	signer, ok := gridEnvelope.Slots[2].(string)
	if !ok {
		return Envelope{}, errors.New("package record signer must be text")
	}
	timestamp, ok := gridEnvelope.Slots[3].(string)
	if !ok {
		return Envelope{}, errors.New("package record timestamp must be text")
	}
	payload, ok := gridEnvelope.Slots[4].([]byte)
	if !ok {
		return Envelope{}, errors.New("package record payload must be JSON bytes")
	}
	keyID, err := nullableText(gridEnvelope.Slots[5])
	if err != nil {
		return Envelope{}, fmt.Errorf("package record author key ID: %w", err)
	}
	publicKey, err := nullableText(gridEnvelope.Slots[6])
	if err != nil {
		return Envelope{}, fmt.Errorf("package record author public key: %w", err)
	}
	signature, err := nullableText(gridEnvelope.Slots[7])
	if err != nil {
		return Envelope{}, fmt.Errorf("package record author signature: %w", err)
	}
	envelope := Envelope{Family: family, ProtocolPCID: gridEnvelope.ProtocolPCID.String(), RecordID: recordID, Signer: signer, Timestamp: timestamp, Payload: append(json.RawMessage{}, payload...), AuthorKeyID: keyID, AuthorPublicKey: publicKey, AuthorSignature: signature}
	return envelope, envelope.Validate()
}

func MustMarshal(envelope Envelope) []byte {
	raw, err := marshal(envelope)
	if err != nil {
		panic(err)
	}
	return raw
}

func (envelope Envelope) HasAuthorSignature() bool {
	return strings.TrimSpace(envelope.AuthorKeyID) != "" || strings.TrimSpace(envelope.AuthorPublicKey) != "" || strings.TrimSpace(envelope.AuthorSignature) != ""
}

func (envelope Envelope) SigningBytes() ([]byte, error) {
	envelope.AuthorKeyID = ""
	envelope.AuthorPublicKey = ""
	envelope.AuthorSignature = ""
	return marshal(envelope)
}

func marshal(envelope Envelope) ([]byte, error) {
	if err := envelope.Validate(); err != nil {
		return nil, err
	}
	protocolPCID, err := cid.Decode(envelope.ProtocolPCID)
	if err != nil {
		return nil, err
	}
	return EncodeGrid(GridEnvelope{ProtocolPCID: protocolPCID, Slots: []any{envelope.Family, envelope.RecordID, envelope.Signer, envelope.Timestamp, []byte(envelope.Payload), nullable(envelope.AuthorKeyID), nullable(envelope.AuthorPublicKey), nullable(envelope.AuthorSignature)}})
}

// CanonicalJSON returns the deterministic JSON representation embedded in the
// canonical Grid package-record payload slot. Source: DI-sidoh
func CanonicalJSON(raw []byte) ([]byte, error) {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

func nullable(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func nullableText(value any) (string, error) {
	if value == nil {
		return "", nil
	}
	text, ok := value.(string)
	if !ok {
		return "", errors.New("must be text or null")
	}
	return text, nil
}
