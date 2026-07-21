package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const maxEventLineBytes = 1 << 20

type Store struct {
	root      string
	events    *os.File
	eventPath string
	draftPath string
}

type PersistedDraft struct {
	Body      string `json:"body"`
	Version   int    `json:"version"`
	UpdatedAt string `json:"updated_at"`
}

// Intent: Keep durable operational truth in an ex5-local append-only log plus
// copied attachments so the example can preserve history independently of any
// browser or CLI session state. Source: DI-radok; DI-zuvob
func OpenStore(root string) (*Store, []OperationalEvent, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "attachments"), 0o755); err != nil {
		return nil, nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "drafts"), 0o755); err != nil {
		return nil, nil, err
	}
	eventPath := filepath.Join(root, "events.jsonl")
	eventsFile, err := os.OpenFile(eventPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, err
	}
	events, err := readEvents(eventsFile)
	if err != nil {
		eventsFile.Close()
		return nil, nil, err
	}
	if _, err := eventsFile.Seek(0, os.SEEK_END); err != nil {
		eventsFile.Close()
		return nil, nil, err
	}
	return &Store{root: root, events: eventsFile, eventPath: eventPath, draftPath: filepath.Join(root, "drafts")}, events, nil
}

func readEvents(file *os.File) ([]OperationalEvent, error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer file.Seek(0, os.SEEK_END)
	scanner := bufio.NewScanner(file)
	// Intent: Replay the full stored event log within the server's current
	// request-size envelope so durable large revisions do not become
	// unreadable after restart. Source: DI-busor
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	events := []OperationalEvent{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var event OperationalEvent
		if err := json.Unmarshal(line, &event); err != nil {
			return nil, fmt.Errorf("decode event: %w", err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (store *Store) AppendEvent(event OperationalEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if _, err := store.events.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.events.Sync()
}

// Intent: Preserve evidence attachment history by storing each uploaded file at
// a unique immutable path instead of overwriting earlier evidence bytes when a
// later upload reuses the same filename. Source: DI-busor
func (store *Store) SaveAttachment(entityID string, filename string, data []byte) (string, int64, error) {
	dir := filepath.Join(store.root, "attachments", entityID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", 0, err
	}
	base := filepath.Base(strings.TrimSpace(filename))
	if base == "" || base == "." {
		base = "attachment.bin"
	}
	tempFile, err := os.CreateTemp(dir, "evidence-*-"+base)
	if err != nil {
		return "", 0, err
	}
	target := tempFile.Name()
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return "", 0, err
	}
	if err := tempFile.Close(); err != nil {
		return "", 0, err
	}
	return target, int64(len(data)), nil
}

// Intent: Restore the shared browser working bodies on startup without mixing
// them into the append-only event log, so durable revision history stays
// explicit while collaborative drafting can resume after a restart. Source:
// DI-lusov; DI-zoruk
func (store *Store) LoadDrafts() (map[string]PersistedDraft, error) {
	entries, err := os.ReadDir(store.draftPath)
	if err != nil {
		return nil, err
	}
	out := map[string]PersistedDraft{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		body, err := os.ReadFile(filepath.Join(store.draftPath, entry.Name()))
		if err != nil {
			return nil, err
		}
		var draft PersistedDraft
		if err := json.Unmarshal(body, &draft); err != nil {
			return nil, fmt.Errorf("decode draft %s: %w", entry.Name(), err)
		}
		out[strings.TrimSuffix(entry.Name(), ".json")] = draft
	}
	return out, nil
}

// Intent: Persist the current shared working body separately from durable
// revision snapshots so browser collaboration can converge on one draft
// without rewriting historical revision events. Source: DI-lusov; DI-zoruk
func (store *Store) SaveDraft(entityID string, draft PersistedDraft) error {
	body, err := json.Marshal(draft)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(store.draftPath, entityID+".json"), body, 0o644)
}

func (store *Store) Close() error {
	if store == nil || store.events == nil {
		return nil
	}
	return store.events.Close()
}
