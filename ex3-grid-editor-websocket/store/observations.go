package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Observation is one relay-local fact about a rejected bounded peer envelope.
// It is deliberately not an accepted protocol message or a statement about a
// sender's identity or intent. Source: DI-pazis; DI-darif.
type Observation struct {
	Kind          string `json:"kind"`
	ObservedAt    string `json:"observed_at"`
	ObserverKeyID string `json:"observer_key_id"`
	RawCID        string `json:"raw_cid,omitempty"`
	ObservedPCID  string `json:"observed_pcid,omitempty"`
	Reason        string `json:"reason,omitempty"`
}

// ObservationLog appends exception evidence separately from the accepted log,
// so restarting a relay never interprets a rejected envelope as replay input.
// Source: DI-pazis; DI-darif.
type ObservationLog struct {
	path string
	mu   sync.Mutex
}

func OpenObservationLog(path string) (*ObservationLog, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir observation log dir: %w", err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open observation log: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close observation log: %w", err)
	}
	return &ObservationLog{path: path}, nil
}

func (log *ObservationLog) Append(observation Observation) error {
	log.mu.Lock()
	defer log.mu.Unlock()

	if observation.ObservedAt == "" {
		observation.ObservedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	line, err := json.Marshal(observation)
	if err != nil {
		return fmt.Errorf("marshal observation: %w", err)
	}
	file, err := os.OpenFile(log.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open observation append: %w", err)
	}
	if _, err := file.Write(append(line, '\n')); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return fmt.Errorf("append observation: %w (close file: %v)", err, closeErr)
		}
		return fmt.Errorf("append observation: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close observation append: %w", err)
	}
	return nil
}

func ReadObservations(path string) (observations []Observation, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open observation log: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close observation log: %w", closeErr)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var observation Observation
		if err := json.Unmarshal([]byte(line), &observation); err != nil {
			return nil, fmt.Errorf("decode observation: %w", err)
		}
		observations = append(observations, observation)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan observations: %w", err)
	}
	return observations, nil
}
