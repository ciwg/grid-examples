package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Event struct {
	Type       string    `json:"type"`
	ToolID     string    `json:"toolId"`
	ActorID    string    `json:"actorId"`
	Text       string    `json:"text"`
	SafetyHold bool      `json:"safetyHold,omitempty"`
	Photos     []Photo   `json:"photos,omitempty"`
	DueAt      time.Time `json:"dueAt,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Store struct {
	path string
}

func NewStore(root string) (*Store, error) {
	if root == "" {
		return nil, errors.New("runtime root is required")
	}
	if err := os.MkdirAll(root, 0o750); err != nil {
		return nil, err
	}
	return &Store{path: filepath.Join(root, "events.jsonl")}, nil
}

func (s *Store) Append(event Event) error {
	file, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(event); err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	return file.Close()
}

func (s *Store) ReadAll() ([]Event, error) {
	file, err := os.Open(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var events []Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		closeErr := file.Close()
		if closeErr != nil {
			return nil, errors.Join(err, closeErr)
		}
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return events, nil
}
