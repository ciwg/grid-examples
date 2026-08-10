package artifact

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestStoreRetainsRawBytesAndAppendsObservation(t *testing.T) {
	root := t.TempDir()
	store, err := NewStore("alice", root)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	raw := []byte{0x82, 0x01, 0x02}
	rawCID, err := store.SaveRawBytes(raw)
	if err != nil {
		t.Fatalf("save raw bytes: %v", err)
	}
	if err := store.AppendObservation(ObservationRecord{Kind: "malformed_input", RawCID: rawCID}); err != nil {
		t.Fatalf("append first observation: %v", err)
	}
	if err := store.AppendObservation(ObservationRecord{Kind: "invalid_proof", RawCID: rawCID}); err != nil {
		t.Fatalf("append second observation: %v", err)
	}
	if _, err := NewStore("alice", root); err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	replayed, err := os.ReadFile(filepath.Join(root, "message-cas", rawCID+".cbor"))
	if err != nil {
		t.Fatalf("read retained raw bytes: %v", err)
	}
	if string(replayed) != string(raw) {
		t.Fatalf("retained raw bytes = %x, want %x", replayed, raw)
	}
	logFile, err := os.Open(filepath.Join(root, "observations.jsonl"))
	if err != nil {
		t.Fatalf("open observation log: %v", err)
	}
	defer func() {
		if closeErr := logFile.Close(); closeErr != nil {
			t.Errorf("close observation log: %v", closeErr)
		}
	}()
	decoder := json.NewDecoder(logFile)
	var observations []ObservationRecord
	for {
		var observation ObservationRecord
		err := decoder.Decode(&observation)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("decode observation: %v", err)
		}
		observations = append(observations, observation)
	}
	if len(observations) != 2 {
		t.Fatalf("observation count = %d, want 2", len(observations))
	}
	for _, observation := range observations {
		if observation.ObserverRole != "alice" || observation.RawCID != rawCID {
			t.Fatalf("observation = %#v", observation)
		}
	}
}
