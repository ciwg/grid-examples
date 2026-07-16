package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	root           string
	logPath        string
	attachmentsDir string
	mu             sync.Mutex
}

func OpenStore(root string) (*Store, []IssueEvent, error) {
	store := &Store{
		root:           root,
		logPath:        filepath.Join(root, "events.jsonl"),
		attachmentsDir: filepath.Join(root, "attachments"),
	}
	if err := os.MkdirAll(store.attachmentsDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("mkdir attachments: %w", err)
	}
	events, err := store.loadEvents()
	if err != nil {
		return nil, nil, err
	}
	return store, events, nil
}

func (store *Store) loadEvents() ([]IssueEvent, error) {
	file, err := os.Open(store.logPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open event log: %w", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var events []IssueEvent
	for scanner.Scan() {
		var event IssueEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, fmt.Errorf("decode event log line: %w", err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return nil, fmt.Errorf("scan event log: %w (close event log: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("scan event log: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close event log: %w", err)
	}
	return events, nil
}

func (store *Store) Append(event IssueEvent) (err error) {
	record, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	file, err := os.OpenFile(store.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open event log: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close event log: %w", closeErr)
		}
	}()
	if _, err := file.Write(append(record, '\n')); err != nil {
		return fmt.Errorf("append event log: %w", err)
	}
	return nil
}

func (store *Store) WriteAttachment(relativePath string, contents []byte) error {
	fullPath := filepath.Join(store.root, relativePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("mkdir attachment path: %w", err)
	}
	if err := os.WriteFile(fullPath, contents, 0o644); err != nil {
		return fmt.Errorf("write attachment: %w", err)
	}
	return nil
}

func (store *Store) RemoveAttachment(relativePath string) error {
	fullPath := filepath.Join(store.root, relativePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove attachment: %w", err)
	}
	return nil
}

func (store *Store) ReadAttachment(relativePath string) ([]byte, error) {
	fullPath := filepath.Join(store.root, relativePath)
	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read attachment: %w", err)
	}
	return bytes, nil
}
