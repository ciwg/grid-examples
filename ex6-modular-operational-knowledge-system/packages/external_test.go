package packages

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRunnerDescribeValidateAndRun(t *testing.T) {
	executable := filepath.Join(t.TempDir(), "helper-egg.sh")
	script := `#!/bin/sh
set -eu
case "$1" in
  describe)
    cat <<'EOF'
{"id":"helper-egg","version":"0.1.0","description":"Test helper package","commands":[{"path":["helper","echo"],"summary":"Echo a string"}],"families":[{"name":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1"}],"claims":[{"protocol_pcid":"pcid:helper.echo.v1","role":"family-validator","summary":"Validates helper echo envelopes."}]}
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
	if manifest.ID != "helper-egg" {
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
	if output != "hello" {
		t.Fatalf("unexpected output: %s", output)
	}
}
