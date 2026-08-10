package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Observation is one relay-local record of bounded bytes the relay assessed.
// Intent: Keep malformed, invalid-proof, and unsupported-handler evidence
// separate from accepted messages and their replay projections. Source: DI-todav; DI-nilas
type Observation struct {
	Kind          string `json:"kind"`
	ObservedAt    string `json:"observed_at"`
	ObserverKeyID string `json:"observer_key_id"`
	RawCID        string `json:"raw_cid"`
	ObservedPCID  string `json:"observed_pcid,omitempty"`
	Reason        string `json:"reason,omitempty"`
}

type ObservationLog struct {
	path string
	mu   sync.Mutex
}

func OpenObservationLog(path string) (*ObservationLog, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir observation log dir: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("create observation log: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close observation log: %w", err)
	}
	return &ObservationLog{path: path}, nil
}

func (log *ObservationLog) Append(record Observation) (err error) {
	encoded, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal observation: %w", err)
	}
	log.mu.Lock()
	defer log.mu.Unlock()
	file, err := os.OpenFile(log.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open observation log: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close observation log: %w", closeErr)
		}
	}()
	if _, err := file.Write(append(encoded, '\n')); err != nil {
		return fmt.Errorf("append observation: %w", err)
	}
	return nil
}
