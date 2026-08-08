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
	Type       string  `json:"type"`
	ToolID     string  `json:"toolId"`
	ActorID    string  `json:"actorId"`
	Text       string  `json:"text"`
	SafetyHold bool    `json:"safetyHold,omitempty"`
	Photos     []Photo `json:"photos,omitempty"`
	// Intent: Preserve the exact loan terms accepted at checkout so replay never
	// substitutes a later area policy. Source: DI-patag.
	Loan *Loan `json:"loan,omitempty"`
	// Intent: Decode pre-snapshot loan events without rewriting their evidence.
	// Source: DI-malih.
	DueAt     time.Time `json:"dueAt,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

const (
	maxPhotoDataURLBytes = 2 * 1024 * 1024
	maxEventBytes        = 8 * 1024 * 1024
)

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
	closeWithError := func(operationErr error) error {
		if closeErr := file.Close(); closeErr != nil {
			return errors.Join(operationErr, closeErr)
		}
		return operationErr
	}
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(event); err != nil {
		return closeWithError(err)
	}
	// Intent: Do not acknowledge evidence until its bytes reach stable storage.
	// Source: DI-dapod.
	if err := file.Sync(); err != nil {
		return closeWithError(err)
	}
	return closeWithError(nil)
}

func (s *Store) ReadAll() (events []Event, returnErr error) {
	file, err := os.Open(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			returnErr = errors.Join(returnErr, err)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventBytes)
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
		return nil, err
	}
	return events, nil
}
