package artifact

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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
	exactCID, err := protocol.CIDForBytes(envelopeBytes)
	if err != nil {
		return "", fmt.Errorf("cid for envelope: %w", err)
	}
	exactCIDText = exactCID.String()
	messagePath := filepath.Join(store.root, "message-cas", exactCIDText+".cbor")
	if writeErr := os.WriteFile(messagePath, envelopeBytes, 0o644); writeErr != nil {
		return "", fmt.Errorf("write local artifact: %w", writeErr)
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
