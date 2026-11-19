package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AdmissionDiagnostic records a relay-local remote-admission denial without
// retaining bearer capabilities, bootstrap secrets, or raw WebSocket frames.
// Source: DI-pazis; DI-darif.
type AdmissionDiagnostic struct {
	Kind          string `json:"kind"`
	ObservedAt    string `json:"observed_at"`
	ObserverKeyID string `json:"observer_key_id"`
	Transport     string `json:"transport"`
	Reason        string `json:"reason"`
}

// AdmissionDiagnosticLog keeps remote-admission diagnostics separate from
// rejected peer-envelope observations and accepted message replay. Source:
// DI-pazis; DI-darif.
type AdmissionDiagnosticLog struct {
	path string
	mu   sync.Mutex
}

func OpenAdmissionDiagnosticLog(path string) (*AdmissionDiagnosticLog, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir admission diagnostic log dir: %w", err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open admission diagnostic log: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close admission diagnostic log: %w", err)
	}
	return &AdmissionDiagnosticLog{path: path}, nil
}

func (log *AdmissionDiagnosticLog) Append(diagnostic AdmissionDiagnostic) error {
	log.mu.Lock()
	defer log.mu.Unlock()

	if diagnostic.ObservedAt == "" {
		diagnostic.ObservedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	line, err := json.Marshal(diagnostic)
	if err != nil {
		return fmt.Errorf("marshal admission diagnostic: %w", err)
	}
	file, err := os.OpenFile(log.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open admission diagnostic append: %w", err)
	}
	if _, err := file.Write(append(line, '\n')); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return fmt.Errorf("append admission diagnostic: %w (close file: %v)", err, closeErr)
		}
		return fmt.Errorf("append admission diagnostic: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close admission diagnostic append: %w", err)
	}
	return nil
}
