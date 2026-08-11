package service

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fxamacker/cbor/v2"
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
	maxFrameBytes        = 1 * 1024 * 1024
)

type Store struct {
	root      string
	framePath string
}

func NewStore(root string) (*Store, error) {
	if root == "" {
		return nil, errors.New("runtime root is required")
	}
	if err := os.MkdirAll(root, 0o750); err != nil {
		return nil, err
	}
	return &Store{root: root, framePath: filepath.Join(root, "records.frames")}, nil
}

// AppendRecords durably appends exact canonical record bytes as one projection
// transaction. Intent: preserve participant evidence without re-encoding it.
// Source: DI-rifib; DI-sinov.
func (s *Store) AppendRecords(records [][]byte) (err error) {
	if len(records) == 0 {
		return errors.New("record frame is empty")
	}
	for _, record := range records {
		if _, err := ParseRecord(record); err != nil {
			return fmt.Errorf("validate record: %w", err)
		}
	}
	payload, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return err
	}
	body, err := payload.Marshal(records)
	if err != nil {
		return err
	}
	if len(body) > maxFrameBytes {
		return errors.New("record frame is too large")
	}
	file, err := os.OpenFile(s.framePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	if info, err := file.Stat(); err != nil {
		return err
	} else if info.Size() == 0 {
		if _, err := file.Write([]byte("MSR1\n")); err != nil {
			return err
		}
	}
	var size [8]byte
	binary.BigEndian.PutUint64(size[:], uint64(len(body)))
	if _, err := file.Write(size[:]); err != nil {
		return err
	}
	if _, err := file.Write(body); err != nil {
		return err
	}
	return file.Sync()
}

// ReadRecordFrames replays only complete canonical frames containing exact
// valid record bytes. Intent: preserve durable evidence without assigning
// family semantics at the storage boundary. Source: DI-tohak.
func (s *Store) ReadRecordFrames() (frames [][][]byte, returnErr error) {
	file, err := os.Open(s.framePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			returnErr = errors.Join(returnErr, closeErr)
		}
	}()

	header := make([]byte, len("MSR1\n"))
	if _, err := io.ReadFull(file, header); err != nil {
		return nil, fmt.Errorf("read record-store header: %w", err)
	}
	if string(header) != "MSR1\n" {
		return nil, errors.New("invalid record-store header")
	}
	encoder, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return nil, err
	}
	for {
		var size [8]byte
		read, err := io.ReadFull(file, size[:])
		if errors.Is(err, io.EOF) && read == 0 {
			return frames, nil
		}
		if err != nil {
			return nil, fmt.Errorf("read record-frame length: %w", err)
		}
		length := binary.BigEndian.Uint64(size[:])
		if length == 0 || length > maxFrameBytes {
			return nil, errors.New("invalid record-frame length")
		}
		body := make([]byte, int(length))
		if _, err := io.ReadFull(file, body); err != nil {
			return nil, fmt.Errorf("read record-frame body: %w", err)
		}
		var records [][]byte
		if err := cbor.Unmarshal(body, &records); err != nil {
			return nil, fmt.Errorf("decode record frame: %w", err)
		}
		canonical, err := encoder.Marshal(records)
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(body, canonical) {
			return nil, errors.New("record frame is not canonical")
		}
		if len(records) == 0 {
			return nil, errors.New("record frame is empty")
		}
		for _, record := range records {
			if _, err := ParseRecord(record); err != nil {
				return nil, fmt.Errorf("validate stored record: %w", err)
			}
		}
		frames = append(frames, records)
	}
}
