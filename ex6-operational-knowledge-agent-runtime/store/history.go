package store

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
)

type StoredEnvelope struct {
	Envelope records.Envelope
	Raw      []byte
}

type History struct {
	mu      sync.RWMutex
	file    *os.File
	entries []StoredEnvelope
	seenRaw map[string]struct{}
}

func OpenHistory(root string) (*History, error) {
	if err := osMkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	file, err := openAppendOnlyLog(filepath.Join(root, "history.base64l"))
	if err != nil {
		return nil, err
	}
	lines, err := readLines(file)
	if err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return nil, fmt.Errorf("read history: %w; close history: %v", err, closeErr)
		}
		return nil, err
	}
	history := &History{file: file, seenRaw: map[string]struct{}{}}
	for _, line := range lines {
		raw, err := base64.RawStdEncoding.DecodeString(string(line))
		if err != nil {
			if closeErr := file.Close(); closeErr != nil {
				return nil, fmt.Errorf("decode history record: %w; close history: %v", err, closeErr)
			}
			return nil, fmt.Errorf("decode history record: %w", err)
		}
		envelope, err := records.Parse(raw)
		if err != nil {
			if closeErr := file.Close(); closeErr != nil {
				return nil, fmt.Errorf("parse history: %w; close history: %v", err, closeErr)
			}
			return nil, fmt.Errorf("parse history: %w", err)
		}
		history.entries = append(history.entries, StoredEnvelope{Envelope: envelope, Raw: raw})
		history.seenRaw[string(raw)] = struct{}{}
	}
	return history, nil
}

// Intent: Keep exact bytes as the durable source of truth so unknown package
// families can survive storage and later relay even before local semantics are
// available. Source: DI-moksu
func (history *History) AppendRaw(raw []byte) (records.Envelope, bool, error) {
	envelope, err := records.Parse(raw)
	if err != nil {
		return records.Envelope{}, false, err
	}
	history.mu.Lock()
	defer history.mu.Unlock()
	if _, exists := history.seenRaw[string(raw)]; exists {
		return envelope, false, nil
	}
	if err := appendLine(history.file, []byte(base64.RawStdEncoding.EncodeToString(raw))); err != nil {
		return records.Envelope{}, false, err
	}
	copied := append([]byte{}, raw...)
	history.entries = append(history.entries, StoredEnvelope{Envelope: envelope, Raw: copied})
	history.seenRaw[string(copied)] = struct{}{}
	return envelope, true, nil
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
