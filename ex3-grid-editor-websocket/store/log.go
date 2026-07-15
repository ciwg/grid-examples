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

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocol"
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
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("close created log: %w", err)
		}
		return log, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open log: %w", err)
	}
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
		if closeErr := file.Close(); closeErr != nil {
			return nil, fmt.Errorf("scan log: %w (close file: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("scan log: %w", err)
	}
	// Intent: Make the copied log loader and appender report close failures
	// explicitly so the verification-only cleanup stays local to ex3. Source:
	// DI-rokod
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close log: %w", err)
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
	// Intent: Keep the copied append path explicit about file-close failures so
	// `errcheck` passes without changing ex2. Source: DI-rokod
	line, err := json.Marshal(entry)
	if err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return Entry{}, fmt.Errorf("marshal entry: %w (close file: %v)", err, closeErr)
		}
		return Entry{}, fmt.Errorf("marshal entry: %w", err)
	}
	if _, err := file.Write(append(line, '\n')); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return Entry{}, fmt.Errorf("append entry: %w (close file: %v)", err, closeErr)
		}
		return Entry{}, fmt.Errorf("append entry: %w", err)
	}
	if err := file.Close(); err != nil {
		return Entry{}, fmt.Errorf("close append log: %w", err)
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
