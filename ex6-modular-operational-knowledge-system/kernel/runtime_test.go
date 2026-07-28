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
	contextpkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/context"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/knowledge"
	linkspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/links"
	procedurespkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/procedures"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/runs"
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

func TestContextPackageCommands(t *testing.T) {
	runtime := newRuntime(t)
	if _, err := runtime.RunCommand(context.Background(), []string{"context", "place", "create", "place-1", "Receiving", "Inbound-area"}); err != nil {
		t.Fatalf("create place: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"context", "responsibility", "create", "resp-1", "Reviewer", "Checks-runs"}); err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"context", "resource", "create", "res-1", "Scale", "Bench-scale", "place-1"}); err != nil {
		t.Fatalf("create resource: %v", err)
	}
	listing, err := runtime.RunCommand(context.Background(), []string{"context", "resource", "list"})
	if err != nil {
		t.Fatalf("list resource: %v", err)
	}
	if !strings.Contains(listing, "res-1\tScale\tBench-scale\tplace-1") {
		t.Fatalf("unexpected resource listing: %s", listing)
	}
	inspect, err := runtime.RunCommand(context.Background(), []string{"context", "place", "inspect", "place-1"})
	if err != nil {
		t.Fatalf("inspect place: %v", err)
	}
	if !strings.Contains(inspect, "name: Receiving") {
		t.Fatalf("unexpected place inspect: %s", inspect)
	}
}

func TestKnowledgePackageCommands(t *testing.T) {
	runtime := newRuntime(t)
	if _, err := runtime.RunCommand(context.Background(), []string{"knowledge", "item", "create", "item-1", "procedure", "DockCheck", "Dock-intake-check"}); err != nil {
		t.Fatalf("create item: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"knowledge", "revision", "snapshot", "item-1", "rev-1", "1", "DockCheck-v1", "Check", "the", "dock"}); err != nil {
		t.Fatalf("snapshot revision: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"knowledge", "item", "approve", "item-1", "life-1", "approved"}); err != nil {
		t.Fatalf("approve item: %v", err)
	}
	listing, err := runtime.RunCommand(context.Background(), []string{"knowledge", "item", "list"})
	if err != nil {
		t.Fatalf("list items: %v", err)
	}
	if !strings.Contains(listing, "item-1\tprocedure\tDockCheck\tapproved\trev=1") {
		t.Fatalf("unexpected item listing: %s", listing)
	}
	inspect, err := runtime.RunCommand(context.Background(), []string{"knowledge", "item", "inspect", "item-1"})
	if err != nil {
		t.Fatalf("inspect item: %v", err)
	}
	if !strings.Contains(inspect, "revision: 1") || !strings.Contains(inspect, "status: approved") {
		t.Fatalf("unexpected item inspect: %s", inspect)
	}
}

func TestRunsPackageCommands(t *testing.T) {
	runtime := newRuntime(t)
	if _, err := runtime.RunCommand(context.Background(), []string{"runs", "record", "run-1", "item-1", "alice", "ok", "dock-complete"}); err != nil {
		t.Fatalf("record run: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"runs", "evidence", "add", "run-1", "ev-1", "photo", "kind=image,shift=night", "blob", "payload"}); err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"runs", "approve", "run-1", "ap-1", "accepted", "looks-good"}); err != nil {
		t.Fatalf("approve run: %v", err)
	}
	inspect, err := runtime.RunCommand(context.Background(), []string{"runs", "inspect", "run-1"})
	if err != nil {
		t.Fatalf("inspect run: %v", err)
	}
	if !strings.Contains(inspect, "actor: alice") || !strings.Contains(inspect, "ev-1:photo") || !strings.Contains(inspect, "ap-1:accepted") {
		t.Fatalf("unexpected run inspect: %s", inspect)
	}
}

func TestLinksPackageCommands(t *testing.T) {
	runtime := newRuntime(t)
	if _, err := runtime.RunCommand(context.Background(), []string{"links", "create", "link-1", "item", "item-1", "run", "run-1", "used-by"}); err != nil {
		t.Fatalf("create link: %v", err)
	}
	inspect, err := runtime.RunCommand(context.Background(), []string{"links", "inspect", "link-1"})
	if err != nil {
		t.Fatalf("inspect link: %v", err)
	}
	if !strings.Contains(inspect, "relation: used-by") || !strings.Contains(inspect, "from: item item-1") {
		t.Fatalf("unexpected link inspect: %s", inspect)
	}
}

func TestProceduresPackageCommands(t *testing.T) {
	runtime := newRuntime(t)
	if _, err := runtime.RunCommand(context.Background(), []string{"procedures", "create", "proc-1", "DockCheck", "Check-the-dock"}); err != nil {
		t.Fatalf("create procedure: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"procedures", "record-use", "proc-1", "run-2", "bob", "ok", "followed-v1"}); err != nil {
		t.Fatalf("record procedure use: %v", err)
	}
	inspect, err := runtime.RunCommand(context.Background(), []string{"procedures", "inspect", "proc-1"})
	if err != nil {
		t.Fatalf("inspect procedure: %v", err)
	}
	if !strings.Contains(inspect, "title: DockCheck") || !strings.Contains(inspect, "uses: run-2") {
		t.Fatalf("unexpected procedure inspect: %s", inspect)
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

func TestInstalledPackageCanMutateDurableStateThroughBasket(t *testing.T) {
	runtime := newRuntime(t)
	packageDir := helperWriterPackageDir(t)
	if _, err := runtime.InstallPackageDir(context.Background(), packageDir); err != nil {
		t.Fatalf("install writer package: %v", err)
	}
	output, err := runtime.RunCommand(context.Background(), []string{"writer", "create", "w-1"})
	if err != nil {
		t.Fatalf("run writer package: %v", err)
	}
	if output != "created w-1" {
		t.Fatalf("unexpected output: %s", output)
	}
	if len(runtime.History()) == 0 {
		t.Fatal("expected durable record from installed package")
	}
	found := false
	for _, entry := range runtime.History() {
		if entry.Envelope.Family == "writer.note.v1" && entry.Envelope.RecordID == "w-1" {
			found = true
			if !strings.Contains(string(entry.Envelope.Payload), "sha256:") {
				t.Fatalf("expected CAS reference in payload: %s", string(entry.Envelope.Payload))
			}
		}
	}
	if !found {
		t.Fatal("expected writer.note.v1 record")
	}
}

func TestInstalledPackageSurvivesRestart(t *testing.T) {
	root := t.TempDir()
	runtime, err := kernel.Open(root)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	if _, err := runtime.InstallPackageDir(context.Background(), helperWriterPackageDir(t)); err != nil {
		t.Fatalf("install writer package: %v", err)
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
	output, err := reopened.RunCommand(context.Background(), []string{"writer", "create", "w-2"})
	if err != nil {
		t.Fatalf("run writer command after reopen: %v", err)
	}
	if output != "created w-2" {
		t.Fatalf("unexpected output after reopen: %s", output)
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

func TestKnownFamilyProtocolMismatchRejected(t *testing.T) {
	runtime := newRuntime(t)
	raw := []byte(`{"family":"moks.ops.note.v1","protocol_pcid":"pcid:wrong","record_id":"bad-1","signer":"ops-note","timestamp":"2026-07-28T00:00:00Z","payload":{"title":"x","body_ref":"sha256:abc"}}`)
	if _, err := runtime.AppendRecord(context.Background(), raw); err == nil {
		t.Fatal("expected protocol mismatch rejection")
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

func TestExportBatchIncludesImplementationClaims(t *testing.T) {
	runtime := newRuntime(t)
	batch := runtime.ExportBatch()
	if batch.Implementation != "moks-ex6" {
		t.Fatalf("unexpected implementation: %s", batch.Implementation)
	}
	if len(batch.ImplementationClaims) < 2 {
		t.Fatal("expected implementation claims")
	}
	foundOps := false
	foundContext := false
	foundKnowledge := false
	foundRuns := false
	foundLinks := false
	foundProcedures := false
	for _, claim := range batch.ImplementationClaims {
		if claim.ProtocolPCID == builtin.OpsFamilyProtocol {
			foundOps = true
		}
		if claim.ProtocolPCID == contextpkg.PlaceProtocol {
			foundContext = true
		}
		if claim.ProtocolPCID == knowledgepkg.ItemProtocol {
			foundKnowledge = true
		}
		if claim.ProtocolPCID == runspkg.RunProtocol {
			foundRuns = true
		}
		if claim.ProtocolPCID == linkspkg.TypedProtocol {
			foundLinks = true
		}
		if claim.ProtocolPCID == procedurespkg.ItemProtocol {
			foundProcedures = true
		}
	}
	if !foundOps || !foundContext || !foundKnowledge || !foundRuns || !foundLinks || !foundProcedures {
		t.Fatalf("expected ops and context claims, got %#v", batch.ImplementationClaims)
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
	if err := runtime.RegisterBuiltin(contextpkg.Package()); err != nil {
		t.Fatalf("register context package: %v", err)
	}
	if err := runtime.RegisterBuiltin(knowledgepkg.Package()); err != nil {
		t.Fatalf("register knowledge package: %v", err)
	}
	if err := runtime.RegisterBuiltin(runspkg.Package()); err != nil {
		t.Fatalf("register runs package: %v", err)
	}
	if err := runtime.RegisterBuiltin(linkspkg.Package()); err != nil {
		t.Fatalf("register links package: %v", err)
	}
	if err := runtime.RegisterBuiltin(procedurespkg.Package()); err != nil {
		t.Fatalf("register procedures package: %v", err)
	}
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
		"claims": []map[string]any{
			{"protocol_pcid": "pcid:helper.echo.v1", "role": "family-validator", "summary": "Validates helper echo envelopes."},
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

func helperWriterPackageDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	executable := filepath.Join(dir, "writer-egg.sh")
	script := `#!/bin/sh
set -eu
case "$1" in
  describe)
    cat <<'EOF'
{"id":"writer-egg","version":"0.1.0","description":"Test writer package","commands":[{"path":["writer","create"],"summary":"Create a writer record"}],"families":[{"name":"writer.note.v1","protocol_pcid":"pcid:writer.note.v1"}],"claims":[{"protocol_pcid":"pcid:writer.note.v1","role":"family-validator","summary":"Validates writer note envelopes."}]}
EOF
    ;;
  validate)
    body="$(cat)"
    case "$body" in
      *'"family":"writer.note.v1"'*) exit 0 ;;
      *) echo "wrong family" >&2; exit 1 ;;
    esac
    ;;
  run)
    if [ "$2" != "writer create" ]; then
      echo "unknown helper command" >&2
      exit 1
    fi
    cat <<EOF
{"output":"created $3","cas":[{"alias":"body1","body":"payload for $3"}],"records":[{"family":"writer.note.v1","protocol_pcid":"pcid:writer.note.v1","record_id":"$3","signer":"writer-egg","timestamp":"2026-07-28T00:00:00Z","payload":{"title":"Writer","body_ref":"\$cas:body1"}}]}
EOF
    ;;
  *)
    echo "unknown helper verb" >&2
    exit 1
    ;;
esac
`
	if err := os.WriteFile(executable, []byte(script), 0o755); err != nil {
		t.Fatalf("write writer script: %v", err)
	}
	manifest := map[string]any{
		"id":          "writer-egg",
		"version":     "0.1.0",
		"description": "Test writer package",
		"executable":  executable,
		"commands": []map[string]any{
			{"path": []string{"writer", "create"}, "summary": "Create a writer record"},
		},
		"families": []map[string]any{
			{"name": "writer.note.v1", "protocol_pcid": "pcid:writer.note.v1"},
		},
		"claims": []map[string]any{
			{"protocol_pcid": "pcid:writer.note.v1", "role": "family-validator", "summary": "Validates writer note envelopes."},
		},
	}
	body, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("marshal writer manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "moks-package.json"), body, 0o644); err != nil {
		t.Fatalf("write writer manifest: %v", err)
	}
	return dir
}
