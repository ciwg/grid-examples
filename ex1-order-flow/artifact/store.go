package artifact

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
)

type Store struct {
	role    string
	root    string
	logPath string
	mu      sync.Mutex
}

type MessageRecord struct {
	Direction string   `json:"direction"`
	ExactCID  string   `json:"exact_cid"`
	PCID      string   `json:"pcid"`
	Parents   []string `json:"parents"`
}

// ObservationRecord is local durable evidence of what one Ex1 component saw.
// Intent: Keep timeout and validation observations distinct from another
// agent's signed promise or refusal. Source: DI-vihoz; DI-riguz; DI-purum
type ObservationRecord struct {
	Kind         string    `json:"kind"`
	ObservedAt   time.Time `json:"observed_at"`
	ObserverRole string    `json:"observer_role"`
	RawCID       string    `json:"raw_cid,omitempty"`
	ExpectedCID  string    `json:"expected_cid,omitempty"`
	ObservedPCID string    `json:"observed_pcid,omitempty"`
	Reason       string    `json:"reason,omitempty"`
}

func NewStore(role string, root string) (*Store, error) {
	store := &Store{
		role:    role,
		root:    root,
		logPath: filepath.Join(root, "messages.jsonl"),
	}
	if err := os.MkdirAll(filepath.Join(root, "message-cas"), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir store: %w", err)
	}
	return store, nil
}

func (store *Store) SaveEnvelope(direction string, envelopeBytes []byte, parents []string, pcidText string) (exactCIDText string, err error) {
	exactCIDText, err = store.SaveRawBytes(envelopeBytes)
	if err != nil {
		return "", err
	}
	record := MessageRecord{
		Direction: direction,
		ExactCID:  exactCIDText,
		PCID:      pcidText,
		Parents:   append([]string(nil), parents...),
	}
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("marshal message record: %w", err)
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	file, err := os.OpenFile(store.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("open message log: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if _, err := file.Write(append(recordBytes, '\n')); err != nil {
		return "", fmt.Errorf("append message log: %w", err)
	}
	return exactCIDText, nil
}

func (store *Store) SaveRawBytes(raw []byte) (string, error) {
	exactCID, err := protocol.CIDForBytes(raw)
	if err != nil {
		return "", fmt.Errorf("cid for raw bytes: %w", err)
	}
	exactCIDText := exactCID.String()
	path := filepath.Join(store.root, "message-cas", exactCIDText+".cbor")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return "", fmt.Errorf("write raw artifact: %w", err)
	}
	return exactCIDText, nil
}

func (store *Store) AppendObservation(record ObservationRecord) (err error) {
	if record.ObservedAt.IsZero() {
		record.ObservedAt = time.Now().UTC()
	}
	if record.ObserverRole == "" {
		record.ObserverRole = store.role
	}
	encoded, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal observation: %w", err)
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	file, err := os.OpenFile(filepath.Join(store.root, "observations.jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open observation log: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if _, err := file.Write(append(encoded, '\n')); err != nil {
		return fmt.Errorf("append observation: %w", err)
	}
	return nil
}

type CollectorArtifactRecord struct {
	SourceRole     string   `json:"source_role"`
	ExactCID       string   `json:"exact_cid"`
	PCID           string   `json:"pcid"`
	ParentCIDs     []string `json:"parent_cids"`
	EnvelopeBase64 string   `json:"envelope_base64"`
}

func EnvelopeBase64(envelopeBytes []byte) string {
	return base64.StdEncoding.EncodeToString(envelopeBytes)
}
