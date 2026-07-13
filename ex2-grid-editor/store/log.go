package store

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/protocol"
)

type Entry struct {
	Offset      uint64 `json:"offset"`
	EnvelopeCID string `json:"envelope_cid"`
	PCID        string `json:"pcid"`
	RawBase64   string `json:"raw_base64"`
	ReceivedAt  string `json:"received_at"`
}

type Log struct {
	path    string
	entries []Entry
	mu      sync.Mutex
}

// Intent: Keep the local service history append-only so projections can be
// rebuilt from exact observed message bytes rather than mutable summary state.
// Source: DI-jilin
func Open(path string) (*Log, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir log dir: %w", err)
	}
	log := &Log{path: path}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_CREATE, 0o644)
		if err != nil {
			return nil, fmt.Errorf("create log: %w", err)
		}
		_ = file.Close()
		return log, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open log: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("decode log entry: %w", err)
		}
		log.entries = append(log.entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan log: %w", err)
	}
	return log, nil
}

func (log *Log) Append(envelopeBytes []byte, pcid string) (Entry, error) {
	log.mu.Lock()
	defer log.mu.Unlock()
	envelopeCID, err := protocol.CIDForBytes(envelopeBytes)
	if err != nil {
		return Entry{}, fmt.Errorf("envelope cid: %w", err)
	}
	entry := Entry{
		Offset:      uint64(len(log.entries)),
		EnvelopeCID: envelopeCID.String(),
		PCID:        pcid,
		RawBase64:   base64.StdEncoding.EncodeToString(envelopeBytes),
		ReceivedAt:  time.Now().Format(time.RFC3339Nano),
	}
	file, err := os.OpenFile(log.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return Entry{}, fmt.Errorf("open append log: %w", err)
	}
	defer file.Close()
	line, err := json.Marshal(entry)
	if err != nil {
		return Entry{}, fmt.Errorf("marshal entry: %w", err)
	}
	if _, err := file.Write(append(line, '\n')); err != nil {
		return Entry{}, fmt.Errorf("append entry: %w", err)
	}
	log.entries = append(log.entries, entry)
	return entry, nil
}

func (log *Log) All() []Entry {
	log.mu.Lock()
	defer log.mu.Unlock()
	return append([]Entry(nil), log.entries...)
}

func (log *Log) EntriesSince(offset uint64) []Entry {
	log.mu.Lock()
	defer log.mu.Unlock()
	if offset >= uint64(len(log.entries)) {
		return nil
	}
	return append([]Entry(nil), log.entries[offset:]...)
}

func (log *Log) NextOffset() uint64 {
	log.mu.Lock()
	defer log.mu.Unlock()
	return uint64(len(log.entries))
}
