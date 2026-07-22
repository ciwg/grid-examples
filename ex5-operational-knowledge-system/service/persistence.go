package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const maxEventLineBytes = 1 << 20

type Store struct {
	root                  string
	events                *os.File
	knowledgeItemMessages *os.File
	eventPath             string
	knowledgeItemPath     string
	draftPath             string
	identity              *RuntimeIdentity
}

type PersistedDraft struct {
	Body      string `json:"body"`
	Version   int    `json:"version"`
	UpdatedAt string `json:"updated_at"`
}

// Intent: Keep durable operational truth in an ex5-local append-only log plus
// copied attachments so the example can preserve history independently of any
// browser or CLI session state. Source: DI-radok; DI-zuvob
func OpenStore(root string) (*Store, []OperationalEvent, []SignedKnowledgeItemRecord, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, nil, nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "attachments"), 0o755); err != nil {
		return nil, nil, nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "drafts"), 0o755); err != nil {
		return nil, nil, nil, err
	}
	identity, err := LoadOrCreateRuntimeIdentity(filepath.Join(root, "identity", "knowledge-item-ed25519.seed"))
	if err != nil {
		return nil, nil, nil, err
	}
	eventPath := filepath.Join(root, "events.jsonl")
	eventsFile, err := os.OpenFile(eventPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, err
	}
	events, err := readEvents(eventsFile)
	if err != nil {
		return nil, nil, nil, errors.Join(err, eventsFile.Close())
	}
	knowledgeItemPath := filepath.Join(root, "knowledge-item-messages.jsonl")
	knowledgeItemFile, err := os.OpenFile(knowledgeItemPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, errors.Join(err, eventsFile.Close())
	}
	knowledgeItemRecords, err := readSignedKnowledgeItemRecords(knowledgeItemFile)
	if err != nil {
		return nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close())
	}
	if _, err := eventsFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close())
	}
	if _, err := knowledgeItemFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close())
	}
	return &Store{
		root:                  root,
		events:                eventsFile,
		knowledgeItemMessages: knowledgeItemFile,
		eventPath:             eventPath,
		knowledgeItemPath:     knowledgeItemPath,
		draftPath:             filepath.Join(root, "drafts"),
		identity:              identity,
	}, events, knowledgeItemRecords, nil
}

func readEvents(file *os.File) (events []OperationalEvent, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	// Intent: Replay the full stored event log within the server's current
	// request-size envelope so durable large revisions do not become
	// unreadable after restart. Source: DI-busor
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	events = []OperationalEvent{}
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

func readSignedKnowledgeItemRecords(file *os.File) (records []SignedKnowledgeItemRecord, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	records = []SignedKnowledgeItemRecord{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record SignedKnowledgeItemRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode knowledge-item record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeItemRecord(record SignedKnowledgeItemRecord) error {
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := store.knowledgeItemMessages.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.knowledgeItemMessages.Sync()
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
		return "", 0, errors.Join(err, tempFile.Close())
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
	if store == nil {
		return nil
	}
	if store.events == nil && store.knowledgeItemMessages == nil {
		return nil
	}
	var err error
	if store.events != nil {
		err = errors.Join(err, store.events.Close())
	}
	if store.knowledgeItemMessages != nil {
		err = errors.Join(err, store.knowledgeItemMessages.Close())
	}
	return err
}
