package store_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/store"
)

func TestObservationLogAppendsRecords(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "observations.jsonl")
	log, err := store.OpenObservationLog(path)
	if err != nil {
		t.Fatalf("open observation log: %v", err)
	}
	if err := log.Append(store.Observation{Kind: "malformed_input", ObservedAt: "2026-08-10T10:00:00Z", ObserverKeyID: "relay-a", RawCID: "raw-a"}); err != nil {
		t.Fatalf("append malformed observation: %v", err)
	}
	if err := log.Append(store.Observation{Kind: "no_supported_handler", ObservedAt: "2026-08-10T10:00:01Z", ObserverKeyID: "relay-a", RawCID: "raw-b", ObservedPCID: "pcid-b"}); err != nil {
		t.Fatalf("append unsupported observation: %v", err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open observation file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("close observation file: %v", closeErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	var records []store.Observation
	for scanner.Scan() {
		var record store.Observation
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			t.Fatalf("decode observation: %v", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan observation file: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("observation count = %d, want 2", len(records))
	}
	if records[0].Kind != "malformed_input" || records[1].Kind != "no_supported_handler" {
		t.Fatalf("observation kinds = %#v", records)
	}
}
