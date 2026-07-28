package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
