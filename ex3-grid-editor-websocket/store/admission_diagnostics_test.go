package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdmissionDiagnosticLogAppendsNonSecretRecord(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "admission-diagnostics.jsonl")
	log, err := OpenAdmissionDiagnosticLog(path)
	if err != nil {
		t.Fatalf("open admission diagnostic log: %v", err)
	}
	if err := log.Append(AdmissionDiagnostic{
		Kind:          "admission_denied",
		ObserverKeyID: "relay-key",
		Transport:     "websocket",
		Reason:        "capability expired",
	}); err != nil {
		t.Fatalf("append admission diagnostic: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read admission diagnostic: %v", err)
	}
	for _, forbidden := range []string{"capability\":", "access_token", "bootstrap"} {
		if strings.Contains(string(raw), forbidden) {
			t.Fatalf("diagnostic retained forbidden material %q: %s", forbidden, raw)
		}
	}
}
