package store

import (
	"path/filepath"
	"testing"
)

func TestObservationLogAppendsRecords(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "observations.jsonl")
	log, err := OpenObservationLog(path)
	if err != nil {
		t.Fatalf("open observation log: %v", err)
	}
	if err := log.Append(Observation{
		Kind:          "no_supported_handler",
		ObserverKeyID: "relay-key",
		RawCID:        "bafkraw",
		ObservedPCID:  "bafkpcid",
		Reason:        "unknown pCID bafkpcid",
	}); err != nil {
		t.Fatalf("append observation: %v", err)
	}
	observations, err := ReadObservations(path)
	if err != nil {
		t.Fatalf("read observations: %v", err)
	}
	if len(observations) != 1 {
		t.Fatalf("observation count = %d, want 1", len(observations))
	}
	if observations[0].Kind != "no_supported_handler" || observations[0].RawCID != "bafkraw" {
		t.Fatalf("observation = %#v", observations[0])
	}
}
