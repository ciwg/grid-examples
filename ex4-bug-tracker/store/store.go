package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/identity"
)

type AcceptedPromise struct {
	ArtifactCID string `json:"artifact_cid"`
	PCID        string `json:"pcid"`
	AgentID     string `json:"agent_id"`
	AcceptedAt  string `json:"accepted_at"`
}

type Observation struct {
	Reason      string `json:"reason"`
	ArtifactCID string `json:"artifact_cid,omitempty"`
	ObservedAt  string `json:"observed_at"`
}

type AcceptedPromiseLog struct {
	path string
	mu   sync.Mutex
}
type ObservationLog struct {
	path string
	mu   sync.Mutex
}
type AgentBindingLog struct {
	path string
	mu   sync.Mutex
}

// ArtifactStore owns the three distinct append-only local log streams.
type ArtifactStore struct {
	AcceptedPromises *AcceptedPromiseLog
	Observations     *ObservationLog
	AgentBindings    *AgentBindingLog
}

func Open(root string) (*ArtifactStore, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("create store root: %w", err)
	}
	return &ArtifactStore{
		AcceptedPromises: &AcceptedPromiseLog{path: filepath.Join(root, "accepted-promises.jsonl")},
		Observations:     &ObservationLog{path: filepath.Join(root, "observations.jsonl")},
		AgentBindings:    &AgentBindingLog{path: filepath.Join(root, "agent-bindings.jsonl")},
	}, nil
}

func (log *AcceptedPromiseLog) Append(record AcceptedPromise) error {
	log.mu.Lock()
	defer log.mu.Unlock()
	return appendJSON(log.path, record)
}
func (log *ObservationLog) Append(record Observation) error {
	log.mu.Lock()
	defer log.mu.Unlock()
	return appendJSON(log.path, record)
}
func (log *AgentBindingLog) Append(record identity.Enrollment) error {
	log.mu.Lock()
	defer log.mu.Unlock()
	return appendJSON(log.path, record)
}

func (log *AgentBindingLog) Load() (records []identity.Enrollment, err error) {
	file, err := os.Open(log.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open enrollment log: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close enrollment log: %w", closeErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var record identity.Enrollment
		if decodeErr := json.Unmarshal(scanner.Bytes(), &record); decodeErr != nil {
			return nil, fmt.Errorf("decode enrollment log: %w", decodeErr)
		}
		records = append(records, record)
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return nil, fmt.Errorf("scan enrollment log: %w", scanErr)
	}
	return records, nil
}

func appendJSON(path string, value any) (err error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal log record: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open log: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close log: %w", closeErr)
		}
	}()
	if _, err := file.Write(append(bytes, '\n')); err != nil {
		return fmt.Errorf("append log: %w", err)
	}
	return nil
}
