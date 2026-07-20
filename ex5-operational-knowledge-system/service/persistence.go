package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	root      string
	events    *os.File
	eventPath string
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
	return &Store{root: root, events: eventsFile, eventPath: eventPath}, events, nil
}

func readEvents(file *os.File) ([]OperationalEvent, error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer file.Seek(0, os.SEEK_END)
	scanner := bufio.NewScanner(file)
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

func (store *Store) SaveAttachment(entityID string, filename string, data []byte) (string, int64, error) {
	dir := filepath.Join(store.root, "attachments", entityID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", 0, err
	}
	target := filepath.Join(dir, filepath.Base(filename))
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return "", 0, err
	}
	return target, int64(len(data)), nil
}

func (store *Store) Close() error {
	if store == nil || store.events == nil {
		return nil
	}
	return store.events.Close()
}
