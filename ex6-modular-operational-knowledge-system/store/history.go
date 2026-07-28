package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/records"
)

type StoredEnvelope struct {
	Envelope records.Envelope
	Raw      []byte
}

type History struct {
	mu      sync.RWMutex
	file    *os.File
	entries []StoredEnvelope
}

func OpenHistory(root string) (*History, error) {
	if err := osMkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	file, err := openAppendOnlyLog(filepath.Join(root, "history.jsonl"))
	if err != nil {
		return nil, err
	}
	lines, err := readLines(file)
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	history := &History{file: file}
	for _, line := range lines {
		envelope, err := records.Parse(line)
		if err != nil {
			_ = file.Close()
			return nil, fmt.Errorf("parse history: %w", err)
		}
		history.entries = append(history.entries, StoredEnvelope{Envelope: envelope, Raw: line})
	}
	return history, nil
}

// Intent: Keep exact bytes as the durable source of truth so unknown package
// families can survive storage and later relay even before local semantics are
// available. Source: DI-moksu
func (history *History) AppendRaw(raw []byte) (records.Envelope, error) {
	envelope, err := records.Parse(raw)
	if err != nil {
		return records.Envelope{}, err
	}
	history.mu.Lock()
	defer history.mu.Unlock()
	if err := appendLine(history.file, raw); err != nil {
		return records.Envelope{}, err
	}
	copied := append([]byte{}, raw...)
	history.entries = append(history.entries, StoredEnvelope{Envelope: envelope, Raw: copied})
	return envelope, nil
}

func (history *History) Entries() []StoredEnvelope {
	history.mu.RLock()
	defer history.mu.RUnlock()
	out := make([]StoredEnvelope, len(history.entries))
	for index, entry := range history.entries {
		out[index] = StoredEnvelope{
			Envelope: entry.Envelope,
			Raw:      append([]byte{}, entry.Raw...),
		}
	}
	return out
}

func (history *History) Close() error {
	history.mu.Lock()
	defer history.mu.Unlock()
	return history.file.Close()
}
