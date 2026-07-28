package kernel_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/builtin"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
)

func TestBuiltinPackageCommandAndCAS(t *testing.T) {
	runtime := newRuntime(t)
	output, err := runtime.RunCommand(context.Background(), []string{"ops", "note", "add", "note-1", "Daily", "checklist complete"})
	if err != nil {
		t.Fatalf("add note: %v", err)
	}
	if !strings.Contains(output, "stored note-1") {
		t.Fatalf("unexpected output: %s", output)
	}
	listing, err := runtime.RunCommand(context.Background(), []string{"ops", "note", "list"})
	if err != nil {
		t.Fatalf("list notes: %v", err)
	}
	if !strings.Contains(listing, "note-1\tDaily\tsha256:") {
		t.Fatalf("unexpected listing: %s", listing)
	}
}

func TestInstalledPackageManifestSelfCheck(t *testing.T) {
	runtime := newRuntime(t)
	packageDir := helperPackageDir(t, false)
	manifest, err := runtime.InstallPackageDir(context.Background(), packageDir)
	if err != nil {
		t.Fatalf("install package: %v", err)
	}
	if manifest.ID != "helper-egg" {
		t.Fatalf("unexpected package id: %s", manifest.ID)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"helper", "echo", "hello"}); err != nil {
		t.Fatalf("run installed package command: %v", err)
	}
}

func TestInstalledPackageManifestMismatchRejected(t *testing.T) {
	runtime := newRuntime(t)
	packageDir := helperPackageDir(t, true)
	if _, err := runtime.InstallPackageDir(context.Background(), packageDir); err == nil {
		t.Fatal("expected manifest self-check mismatch")
	}
}

func TestUnknownFamilyStoredAndLaterInterpreted(t *testing.T) {
	unknownRaw := []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"peer-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	runtimeA := newRuntime(t)
	if _, err := runtimeA.AppendRecord(context.Background(), unknownRaw); err != nil {
		t.Fatalf("append unknown raw: %v", err)
	}
	exported := runtimeA.ExportBatch()
	if got := string(exported.Records[0]); got != string(unknownRaw) {
		t.Fatalf("expected exact bytes preserved, got %s", got)
	}
	runtimeB := newRuntime(t)
	if err := runtimeB.ImportBatch(context.Background(), exported); err != nil {
		t.Fatalf("import batch: %v", err)
	}
	packageDir := helperPackageDir(t, false)
	if _, err := runtimeB.InstallPackageDir(context.Background(), packageDir); err != nil {
		t.Fatalf("install helper package: %v", err)
	}
	if len(runtimeB.History()) != 1 {
		t.Fatalf("expected imported history to remain available, got %d", len(runtimeB.History()))
	}
	if runtimeB.History()[0].Envelope.Family != "helper.echo.v1" {
		t.Fatalf("unexpected family: %s", runtimeB.History()[0].Envelope.Family)
	}
}

func TestHistorySurvivesRestart(t *testing.T) {
	root := t.TempDir()
	runtime, err := kernel.Open(root)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	if err := runtime.RegisterBuiltin(builtin.OpsPackage()); err != nil {
		t.Fatalf("register builtin: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"ops", "note", "add", "note-2", "Shift", "handoff ready"}); err != nil {
		t.Fatalf("add note: %v", err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatalf("close runtime: %v", err)
	}
	reopened, err := kernel.Open(root)
	if err != nil {
		t.Fatalf("reopen runtime: %v", err)
	}
	defer func() {
		_ = reopened.Close()
	}()
	if err := reopened.RegisterBuiltin(builtin.OpsPackage()); err != nil {
		t.Fatalf("register builtin on reopen: %v", err)
	}
	if len(reopened.History()) != 1 {
		t.Fatalf("expected one entry after reopen, got %d", len(reopened.History()))
	}
}

func TestBatchFormatValidation(t *testing.T) {
	runtime := newRuntime(t)
	err := runtime.ImportBatch(context.Background(), grid.Batch{Format: "wrong"})
	if err == nil {
		t.Fatal("expected wrong batch format error")
	}
}

func newRuntime(t *testing.T) *kernel.Runtime {
	t.Helper()
	runtime, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	t.Cleanup(func() {
		_ = runtime.Close()
	})
	if err := runtime.RegisterBuiltin(builtin.OpsPackage()); err != nil {
		t.Fatalf("register builtin: %v", err)
	}
	return runtime
}

func helperPackageDir(t *testing.T, mismatch bool) string {
	t.Helper()
	dir := t.TempDir()
	executable := filepath.Join(dir, "helper-egg.sh")
	script := `#!/bin/sh
set -eu
case "$1" in
  describe)
    cat <<'EOF'
{"id":"helper-egg","version":"0.1.0","description":"Test helper package","commands":[{"path":["helper","echo"],"summary":"Echo a string"}],"families":[{"name":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1"}]}
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
    if [ "${3-}" = "" ]; then
      echo "missing echo value" >&2
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
	manifest := map[string]any{
		"id":          "helper-egg",
		"version":     "0.1.0",
		"description": "Test helper package",
		"executable":  executable,
		"commands": []map[string]any{
			{"path": []string{"helper", "echo"}, "summary": "Echo a string"},
		},
		"families": []map[string]any{
			{"name": "helper.echo.v1", "protocol_pcid": "pcid:helper.echo.v1"},
		},
	}
	if mismatch {
		manifest["commands"] = []map[string]any{
			{"path": []string{"helper", "wrong"}, "summary": "Wrong command"},
		}
	}
	body, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "moks-package.json"), body, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	return dir
}
