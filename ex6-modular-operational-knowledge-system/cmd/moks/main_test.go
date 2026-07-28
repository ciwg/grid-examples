package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
)

func TestBuiltinQuickstartFlow(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "context", "place", "create", "place-1", "Receiving", "Inbound-area"); err != nil {
		t.Fatalf("create place: %v", err)
	}
	if _, err := runCLI(t, workdir, "context", "resource", "create", "res-1", "Scale", "Bench-scale", "place-1"); err != nil {
		t.Fatalf("create resource: %v", err)
	}
	if _, err := runCLI(t, workdir, "procedures", "create", "proc-1", "DockCheck", "Check-the-dock"); err != nil {
		t.Fatalf("create procedure: %v", err)
	}
	if _, err := runCLI(t, workdir, "procedures", "record-use", "proc-1", "run-1", "alice", "ok", "followed-v1"); err != nil {
		t.Fatalf("record procedure use: %v", err)
	}
	if _, err := runCLI(t, workdir, "runs", "evidence", "add", "run-1", "ev-1", "photo", "kind=image,shift=night", "blob", "payload"); err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	if _, err := runCLI(t, workdir, "runs", "approve", "run-1", "ap-1", "accepted", "looks-good"); err != nil {
		t.Fatalf("approve run: %v", err)
	}
	inspectProcedure, err := runCLI(t, workdir, "procedures", "inspect", "proc-1")
	if err != nil {
		t.Fatalf("inspect procedure: %v", err)
	}
	if !strings.Contains(inspectProcedure, "title: DockCheck") || !strings.Contains(inspectProcedure, "uses: run-1") {
		t.Fatalf("unexpected procedure inspect output: %s", inspectProcedure)
	}
	inspectRun, err := runCLI(t, workdir, "runs", "inspect", "run-1")
	if err != nil {
		t.Fatalf("inspect run: %v", err)
	}
	if !strings.Contains(inspectRun, "ev-1:photo") || !strings.Contains(inspectRun, "ap-1:accepted") {
		t.Fatalf("unexpected run inspect output: %s", inspectRun)
	}
	exportPath := filepath.Join(workdir, "relay.json")
	if _, err := runCLI(t, workdir, "relay", "export", exportPath); err != nil {
		t.Fatalf("export relay: %v", err)
	}
	if _, err := os.Stat(exportPath); err != nil {
		t.Fatalf("stat relay export: %v", err)
	}
}

func TestInstalledWriterEggExample(t *testing.T) {
	workdir := t.TempDir()
	exampleDir := filepath.Join(repoRoot(t), "examples", "writer-egg")
	if output, err := runCLI(t, workdir, "package", "install", exampleDir); err != nil {
		t.Fatalf("install writer egg: %v", err)
	} else if !strings.Contains(output, "installed writer-egg") {
		t.Fatalf("unexpected install output: %s", output)
	}
	output, err := runCLI(t, workdir, "writer", "create", "writer-1")
	if err != nil {
		t.Fatalf("writer create: %v", err)
	}
	if !strings.Contains(output, "created writer-1") {
		t.Fatalf("unexpected writer output: %s", output)
	}
	exportPath := filepath.Join(workdir, "writer-relay.json")
	if _, err := runCLI(t, workdir, "relay", "export", exportPath); err != nil {
		t.Fatalf("export writer relay: %v", err)
	}
	body, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("read writer relay export: %v", err)
	}
	if !bytes.Contains(body, []byte("writer.note.v1")) {
		t.Fatalf("writer relay export missing record family: %s", string(body))
	}
}

func TestRelayHandlerExportsAndImportsBatch(t *testing.T) {
	source := newRuntimeForCLI(t)
	if _, err := source.RunCommand(context.Background(), []string{"context", "place", "create", "place-1", "Receiving", "Inbound-area"}); err != nil {
		t.Fatalf("seed source runtime: %v", err)
	}
	if err := source.AllowPeer(grid.AllowedPeer{
		PeerID:    "peer-target",
		BatchURL:  "http://peer-target/relay/batch",
		ImportURL: "http://peer-target/relay/import",
		PublicKey: source.LocalPeerPublicKey(),
		AllowPull: true,
		AllowPush: true,
	}); err != nil {
		t.Fatalf("allow source peer: %v", err)
	}
	server := httptest.NewServer(relayHandler(context.Background(), source))
	defer server.Close()

	response, err := http.Get(server.URL + "/relay/batch")
	if err != nil {
		t.Fatalf("get relay batch: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("unexpected relay batch status: %s", response.Status)
	}
	var batch grid.Batch
	if err := json.NewDecoder(response.Body).Decode(&batch); err != nil {
		t.Fatalf("decode relay batch: %v", err)
	}
	if len(batch.Records) == 0 {
		t.Fatal("expected exported relay records")
	}

	target := newRuntimeForCLI(t)
	if err := target.AllowPeer(grid.AllowedPeer{
		PeerID:    source.LocalPeerID(),
		BatchURL:  server.URL + "/relay/batch",
		ImportURL: server.URL + "/relay/import",
		PublicKey: source.LocalPeerPublicKey(),
		AllowPull: true,
		AllowPush: true,
	}); err != nil {
		t.Fatalf("allow target peer: %v", err)
	}
	importBody, err := json.Marshal(batch)
	if err != nil {
		t.Fatalf("marshal import batch: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/relay/import", bytes.NewReader(importBody))
	request.Header.Set("X-Moks-Peer-ID", source.LocalPeerID())
	recorder := httptest.NewRecorder()
	relayHandler(context.Background(), target).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected import status: %d body=%s", recorder.Code, recorder.Body.String())
	}
	if len(target.History()) == 0 {
		t.Fatal("expected imported history")
	}
}

func TestRelayPeerCardAndDiscover(t *testing.T) {
	source := newRuntimeForCLI(t)
	server := httptest.NewServer(relayHandler(context.Background(), source))
	defer server.Close()

	response, err := http.Get(server.URL + "/relay/peer-card")
	if err != nil {
		t.Fatalf("get peer card: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("unexpected peer card status: %s", response.Status)
	}
	var card grid.PeerCard
	if err := json.NewDecoder(response.Body).Decode(&card); err != nil {
		t.Fatalf("decode peer card: %v", err)
	}
	if err := card.Validate(); err != nil {
		t.Fatalf("validate peer card: %v", err)
	}
	if card.PeerID != source.LocalPeerID() {
		t.Fatalf("unexpected peer id: %s", card.PeerID)
	}
	if card.PublicKey != source.LocalPeerPublicKey() {
		t.Fatalf("unexpected public key: %s", card.PublicKey)
	}

	workdir := t.TempDir()
	output, err := runCLI(t, workdir, "relay", "peer", "discover", server.URL+"/relay/peer-card")
	if err != nil {
		t.Fatalf("discover peer: %v", err)
	}
	if !strings.Contains(output, "peer_id: "+source.LocalPeerID()) {
		t.Fatalf("discover output missing peer id: %s", output)
	}
	if !strings.Contains(output, "allow_command: moks relay peer allow "+source.LocalPeerID()) {
		t.Fatalf("discover output missing allow command: %s", output)
	}
}

func TestRelayPullImportsFromPeer(t *testing.T) {
	source := newRuntimeForCLI(t)
	if _, err := source.RunCommand(context.Background(), []string{"context", "place", "create", "place-1", "Receiving", "Inbound-area"}); err != nil {
		t.Fatalf("seed source runtime: %v", err)
	}
	server := httptest.NewServer(relayHandler(context.Background(), source))
	defer server.Close()

	target := newRuntimeForCLI(t)
	if err := target.AllowPeer(grid.AllowedPeer{
		PeerID:    source.LocalPeerID(),
		BatchURL:  server.URL + "/relay/batch",
		ImportURL: server.URL + "/relay/import",
		PublicKey: source.LocalPeerPublicKey(),
		AllowPull: true,
		AllowPush: false,
	}); err != nil {
		t.Fatalf("allow target peer: %v", err)
	}
	if err := relayPull(context.Background(), target, source.LocalPeerID()); err != nil {
		t.Fatalf("relay pull: %v", err)
	}
	if len(target.History()) != 1 {
		t.Fatalf("expected one imported record, got %d", len(target.History()))
	}
}

func TestRelayPushPostsToPeer(t *testing.T) {
	source := newRuntimeForCLI(t)
	if _, err := source.RunCommand(context.Background(), []string{"context", "place", "create", "place-1", "Receiving", "Inbound-area"}); err != nil {
		t.Fatalf("seed source runtime: %v", err)
	}
	target := newRuntimeForCLI(t)
	server := httptest.NewServer(relayHandler(context.Background(), target))
	defer server.Close()
	if err := target.AllowPeer(grid.AllowedPeer{
		PeerID:    source.LocalPeerID(),
		BatchURL:  server.URL + "/relay/batch",
		ImportURL: server.URL + "/relay/import",
		PublicKey: source.LocalPeerPublicKey(),
		AllowPull: true,
		AllowPush: true,
	}); err != nil {
		t.Fatalf("allow target peer: %v", err)
	}
	if err := source.AllowPeer(grid.AllowedPeer{
		PeerID:    target.LocalPeerID(),
		BatchURL:  server.URL + "/relay/batch",
		ImportURL: server.URL + "/relay/import",
		PublicKey: target.LocalPeerPublicKey(),
		AllowPull: true,
		AllowPush: true,
	}); err != nil {
		t.Fatalf("allow source peer: %v", err)
	}
	if err := relayPush(context.Background(), source, target.LocalPeerID()); err != nil {
		t.Fatalf("relay push: %v", err)
	}
	if len(target.History()) != 1 {
		t.Fatalf("expected pushed record on peer, got %d", len(target.History()))
	}
}

func TestRelayImportRejectsUnknownPeer(t *testing.T) {
	target := newRuntimeForCLI(t)
	batch, err := target.SignedExportBatch()
	if err != nil {
		t.Fatalf("sign batch: %v", err)
	}
	body, err := json.Marshal(batch)
	if err != nil {
		t.Fatalf("marshal batch: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/relay/import", bytes.NewReader(body))
	request.Header.Set("X-Moks-Peer-ID", "unknown-peer")
	recorder := httptest.NewRecorder()
	relayHandler(context.Background(), target).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden for unknown peer, got %d", recorder.Code)
	}
}

func TestRelayPullRejectsUnexpectedPeerIdentity(t *testing.T) {
	source := newRuntimeForCLI(t)
	if _, err := source.RunCommand(context.Background(), []string{"context", "place", "create", "place-1", "Receiving", "Inbound-area"}); err != nil {
		t.Fatalf("seed source runtime: %v", err)
	}
	server := httptest.NewServer(relayHandler(context.Background(), source))
	defer server.Close()

	target := newRuntimeForCLI(t)
	if err := target.AllowPeer(grid.AllowedPeer{
		PeerID:    "wrong-peer",
		BatchURL:  server.URL + "/relay/batch",
		ImportURL: server.URL + "/relay/import",
		PublicKey: source.LocalPeerPublicKey(),
		AllowPull: true,
		AllowPush: false,
	}); err != nil {
		t.Fatalf("allow wrong peer: %v", err)
	}
	if err := relayPull(context.Background(), target, "wrong-peer"); err == nil {
		t.Fatal("expected peer identity mismatch")
	}
}

func TestRelayPullRejectsBadSignature(t *testing.T) {
	source := newRuntimeForCLI(t)
	if _, err := source.RunCommand(context.Background(), []string{"context", "place", "create", "place-1", "Receiving", "Inbound-area"}); err != nil {
		t.Fatalf("seed source runtime: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		batch, err := source.SignedExportBatch()
		if err != nil {
			t.Fatalf("sign batch: %v", err)
		}
		batch.Signature = "00"
		if err := json.NewEncoder(writer).Encode(batch); err != nil {
			t.Fatalf("encode batch: %v", err)
		}
	}))
	defer server.Close()

	target := newRuntimeForCLI(t)
	if err := target.AllowPeer(grid.AllowedPeer{
		PeerID:    source.LocalPeerID(),
		BatchURL:  server.URL,
		ImportURL: server.URL + "/relay/import",
		PublicKey: source.LocalPeerPublicKey(),
		AllowPull: true,
		AllowPush: false,
	}); err != nil {
		t.Fatalf("allow target peer: %v", err)
	}
	if err := relayPull(context.Background(), target, source.LocalPeerID()); err == nil {
		t.Fatal("expected signature verification failure")
	}
}

func runCLI(t *testing.T, workdir string, args ...string) (string, error) {
	t.Helper()
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	if err := os.Chdir(workdir); err != nil {
		t.Fatalf("chdir to workdir: %v", err)
	}
	os.Stdout = writer
	runErr := run(context.Background(), args)
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	os.Stdout = originalStdout
	if err := os.Chdir(originalWD); err != nil {
		t.Fatalf("restore cwd: %v", err)
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close reader: %v", err)
	}
	return strings.TrimSpace(string(body)), runErr
}

func repoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

func newRuntimeForCLI(t *testing.T) *kernel.Runtime {
	t.Helper()
	runtime, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	t.Cleanup(func() {
		_ = runtime.Close()
	})
	registerAllBuiltinsForTest(t, runtime)
	return runtime
}

func registerAllBuiltinsForTest(t *testing.T, runtime *kernel.Runtime) {
	t.Helper()
	if err := registerBuiltins(runtime); err != nil {
		t.Fatalf("register builtins: %v", err)
	}
}
