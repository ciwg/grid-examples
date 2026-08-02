package packages

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRunnerDescribeValidateAndRun(t *testing.T) {
	executable := filepath.Join(t.TempDir(), "helper-agent.sh")
	script := `#!/bin/sh
set -eu
case "$1" in
  describe)
    cat <<'EOF'
{"id":"helper-agent","version":"0.1.0","description":"Test helper package","commands":[{"path":["helper","echo"],"summary":"Echo a string"}],"families":[{"name":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1"}],"claims":[{"protocol_pcid":"pcid:helper.echo.v1","role":"family-validator","summary":"Validates helper echo envelopes."}]}
EOF
    ;;
  validate)
    body="$(cat)"
    case "$body" in
      *'"family":"helper.echo.v1"'*) exit 0 ;;
      *) echo "wrong family" >&2; exit 1 ;;
    esac
    ;;
  run)
    if [ "$2" != "helper echo" ]; then
      echo "unknown helper command" >&2
      exit 1
    fi
    printf '%s\n' "$3"
    ;;
  *)
    echo "unknown helper verb" >&2
    exit 1
    ;;
esac
`
	if err := os.WriteFile(executable, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper script: %v", err)
	}
	runner := Runner{Executable: executable}
	ctx := context.Background()
	manifest, err := runner.Describe(ctx)
	if err != nil {
		t.Fatalf("describe: %v", err)
	}
	if manifest.ID != "helper-agent" {
		t.Fatalf("unexpected id: %s", manifest.ID)
	}
	raw := []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"one","signer":"helper","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	if err := runner.ValidateEnvelope(ctx, raw); err != nil {
		t.Fatalf("validate: %v", err)
	}
	output, err := runner.RunCommand(ctx, "helper echo", []string{"hello"})
	if err != nil {
		t.Fatalf("run command: %v", err)
	}
	if output.Output != "hello" {
		t.Fatalf("unexpected output: %s", output.Output)
	}
}

func TestRunnerRunCommandStructuredResult(t *testing.T) {
	executable := filepath.Join(t.TempDir(), "helper-agent.sh")
	script := `#!/bin/sh
set -eu
case "$1" in
  describe)
    printf '{}\n'
    ;;
  validate)
    exit 0
    ;;
  run)
    cat <<'EOF'
{"output":"created","cas":[{"alias":"body1","body":"hello body"}],"records":[{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"one","signer":"helper","timestamp":"2026-07-28T00:00:00Z","payload":{"body_ref":"$cas:body1"}}]}
EOF
    ;;
  *)
    exit 1
    ;;
esac
`
	if err := os.WriteFile(executable, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper script: %v", err)
	}
	runner := Runner{Executable: executable}
	result, err := runner.RunCommand(context.Background(), "helper create", nil)
	if err != nil {
		t.Fatalf("run structured command: %v", err)
	}
	if result.Output != "created" || len(result.CAS) != 1 || len(result.Records) != 1 {
		t.Fatalf("unexpected structured result: %#v", result)
	}
	var record map[string]any
	if err := json.Unmarshal(result.Records[0], &record); err != nil {
		t.Fatalf("unmarshal record: %v", err)
	}
}

func TestManifestValidateRequiresEmitsForParserClaims(t *testing.T) {
	manifest := Manifest{
		ID:      "parser-agent",
		Version: "0.1.0",
		Claims: []ImplementationClaim{
			{ProtocolPCID: "pcid:raw.example.v1", Role: "parser", RouteType: "parser"},
		},
	}
	if err := manifest.Validate(); err == nil {
		t.Fatal("expected parser claim validation failure")
	}
}

func TestManifestValidateWorkflowAdapters(t *testing.T) {
	manifest := Manifest{
		ID:      "procedure-execution-adapter",
		Version: "0.1.0",
		WorkflowAdapters: []WorkflowAdapter{{
			Name:       "procedure-execution",
			Image:      "example/procedure-execution:1",
			Command:    []string{"worker"},
			InputPCID:  "bafkreiahdp34nto2rnnqde26jw3xnkd6xnlalnr72sug3w7tjb3bhhoj4q",
			OutputPCID: "bafkreifmttp5fwt3yvxvkb7ni6kwg3j3arl7mbjsyzszf7s7crxrncch24",
			CPUs:       "0.5",
			Memory:     "128m",
			PIDsLimit:  64,
			Timeout:    "30s",
		}},
	}
	if err := manifest.Validate(); err != nil {
		t.Fatalf("validate workflow adapter: %v", err)
	}
	changed := manifest
	changed.WorkflowAdapters = append([]WorkflowAdapter{}, manifest.WorkflowAdapters...)
	changed.WorkflowAdapters[0].OutputPCID = "bafkreih6yllp2v7e5bmerznebzmohniezsv64hpqe2m33h6jclq6rfzqdu"
	if manifest.Equal(changed) {
		t.Fatal("manifest equality ignored workflow adapter output contract")
	}
	changed.WorkflowAdapters[0].PIDsLimit = 0
	if err := changed.Validate(); err == nil {
		t.Fatal("accepted workflow adapter without PID limit")
	}
}
