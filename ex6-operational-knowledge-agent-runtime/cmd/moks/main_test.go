package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
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

func TestWorkflowDemoInventoryReceipt(t *testing.T) {
	output, err := runCLI(t, repoRoot(t), "workflow", "demo", "inventory-receipt")
	if err != nil {
		t.Fatalf("run inventory receipt demo: %v", err)
	}
	for _, line := range []string{
		"[ok] captured inventory-receipt as ",
		"[ok] verified manifest and local dependencies",
		"[ok] activated local workflow availability",
		"[ok] extracted the exact retained artifact for inspection",
	} {
		if !strings.Contains(output, line) {
			t.Fatalf("demo output missing %q: %s", line, output)
		}
	}
}

func TestWorkflowVerifyReportsExecutionReadiness(t *testing.T) {
	workdir := t.TempDir()
	sourceDir := filepath.Join(repoRoot(t), "workflows", "inventory-receipt")
	if _, err := runCLI(t, workdir, "workflow", "capture", sourceDir, "inventory-receipt"); err != nil {
		t.Fatal(err)
	}
	output, err := runCLI(t, workdir, "workflow", "verify", "inventory-receipt")
	if err != nil {
		t.Fatal(err)
	}
	var inactive struct {
		Contract          string `json:"contract"`
		AdapterAvailable  bool   `json:"adapter_available"`
		SchemaCASReady    bool   `json:"schema_cas_ready"`
		EligibleToExecute bool   `json:"eligible_to_execute"`
	}
	if err := json.Unmarshal([]byte(output), &inactive); err != nil {
		t.Fatal(err)
	}
	if inactive.Contract != "canonical" || !inactive.AdapterAvailable || !inactive.SchemaCASReady || inactive.EligibleToExecute {
		t.Fatalf("inactive verification = %#v", inactive)
	}
	if _, err := runCLI(t, workdir, "workflow", "activate", "inventory-receipt"); err != nil {
		t.Fatal(err)
	}
	output, err = runCLI(t, workdir, "workflow", "verify", "inventory-receipt")
	if err != nil {
		t.Fatal(err)
	}
	var active struct {
		EligibleToExecute bool `json:"eligible_to_execute"`
	}
	if err := json.Unmarshal([]byte(output), &active); err != nil {
		t.Fatal(err)
	}
	if !active.EligibleToExecute {
		t.Fatalf("active verification = %s", output)
	}
}

func TestWorkflowOverviewSummarizesInboxAndNextImport(t *testing.T) {
	workdir := t.TempDir()
	source, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open source: %v", err)
	}
	defer func() {
		if closeErr := source.Close(); closeErr != nil {
			t.Errorf("close source: %v", closeErr)
		}
	}()
	target, err := kernel.Open(filepath.Join(workdir, ".moks"))
	if err != nil {
		t.Fatalf("open target: %v", err)
	}
	if err := target.AllowPeer(grid.AllowedPeer{
		PeerID:            source.LocalPeerID(),
		BatchURL:          "https://source.invalid/relay/batch",
		ImportURL:         "https://source.invalid/relay/import",
		PublicKey:         source.LocalPeerPublicKey(),
		AllowPush:         true,
		AttesterClass:     "peer",
		AttestationWeight: 1,
		Federation:        "independent",
	}); err != nil {
		t.Fatalf("allow source: %v", err)
	}
	artifact, err := source.PutCAS([]byte("overview inbox workflow artifact"))
	if err != nil {
		t.Fatalf("store source artifact: %v", err)
	}
	if err := source.ImportWorkflow(kernel.Workflow{ID: "source-overview", ArtifactCID: artifact}); err != nil {
		t.Fatalf("import source workflow: %v", err)
	}
	transfer, err := source.ExportWorkflowTransfer("source-overview")
	if err != nil {
		t.Fatalf("export source workflow: %v", err)
	}
	if err := target.ImportWorkflowTransferFromPeer(source.LocalPeerID(), transfer); err != nil {
		t.Fatalf("receive source workflow: %v", err)
	}
	if closeErr := target.Close(); closeErr != nil {
		t.Fatalf("close target: %v", closeErr)
	}
	overview, err := runCLI(t, workdir, "workflow", "overview")
	if err != nil {
		t.Fatalf("workflow overview: %v", err)
	}
	for _, expected := range []string{
		"WORKFLOW OVERVIEW",
		"[inbox] " + transfer.ArtifactCID + " — ready to import",
		"NEXT: moks workflow inbox import " + transfer.ArtifactCID + " <alias>",
	} {
		if !strings.Contains(overview, expected) {
			t.Fatalf("overview missing %q: %s", expected, overview)
		}
	}
	if _, err := runCLI(t, workdir, "workflow", "inbox", "import", transfer.ArtifactCID, "received-overview"); err != nil {
		t.Fatalf("import overview inbox: %v", err)
	}
	overview, err = runCLI(t, workdir, "workflow", "overview")
	if err != nil {
		t.Fatalf("overview after inbox import: %v", err)
	}
	if strings.Contains(overview, "[inbox] "+transfer.ArtifactCID) || !strings.Contains(overview, "Needs attention: 1") || !strings.Contains(overview, "NEXT: moks workflow verify received-overview") {
		t.Fatalf("overview retains imported inbox attention: %s", overview)
	}
}

func TestWorkflowOverviewSelectsFirstImportableInboxArtifact(t *testing.T) {
	workdir := t.TempDir()
	source, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open source: %v", err)
	}
	defer func() {
		if closeErr := source.Close(); closeErr != nil {
			t.Errorf("close source: %v", closeErr)
		}
	}()
	target, err := kernel.Open(filepath.Join(workdir, ".moks"))
	if err != nil {
		t.Fatalf("open target: %v", err)
	}
	if err := target.AllowPeer(grid.AllowedPeer{
		PeerID:            source.LocalPeerID(),
		BatchURL:          "https://source.invalid/relay/batch",
		ImportURL:         "https://source.invalid/relay/import",
		PublicKey:         source.LocalPeerPublicKey(),
		AllowPush:         true,
		AttesterClass:     "peer",
		AttestationWeight: 1,
		Federation:        "independent",
	}); err != nil {
		t.Fatalf("allow source: %v", err)
	}
	artifactCIDs := make([]string, 0, 2)
	for _, fixture := range []struct {
		alias string
		body  string
	}{
		{alias: "source-first", body: "first overview inbox workflow artifact"},
		{alias: "source-second", body: "second overview inbox workflow artifact"},
	} {
		artifactCID, err := source.PutCAS([]byte(fixture.body))
		if err != nil {
			t.Fatalf("store source artifact: %v", err)
		}
		if err := source.ImportWorkflow(kernel.Workflow{ID: fixture.alias, ArtifactCID: artifactCID}); err != nil {
			t.Fatalf("import source workflow: %v", err)
		}
		transfer, err := source.ExportWorkflowTransfer(fixture.alias)
		if err != nil {
			t.Fatalf("export source workflow: %v", err)
		}
		if err := target.ImportWorkflowTransferFromPeer(source.LocalPeerID(), transfer); err != nil {
			t.Fatalf("receive source workflow: %v", err)
		}
		artifactCIDs = append(artifactCIDs, transfer.ArtifactCID)
	}
	if closeErr := target.Close(); closeErr != nil {
		t.Fatalf("close target: %v", closeErr)
	}
	slices.Sort(artifactCIDs)
	overview, err := runCLI(t, workdir, "workflow", "overview")
	if err != nil {
		t.Fatalf("workflow overview: %v", err)
	}
	if !strings.Contains(overview, "NEXT: moks workflow inbox import "+artifactCIDs[0]+" <alias>") {
		t.Fatalf("overview did not select first importable inbox artifact: %s", overview)
	}
}

func TestWorkflowOverviewReportsEmptyRuntimeDeterministically(t *testing.T) {
	workdir := t.TempDir()
	first, err := runCLI(t, workdir, "workflow", "overview")
	if err != nil {
		t.Fatalf("first workflow overview: %v", err)
	}
	second, err := runCLI(t, workdir, "workflow", "overview")
	if err != nil {
		t.Fatalf("second workflow overview: %v", err)
	}
	if first != second {
		t.Fatalf("overview is not deterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	for _, expected := range []string{
		"Ready: 0 active workflows",
		"Needs attention: 0",
		"Recent activity: none",
		"NEXT: no action required",
	} {
		if !strings.Contains(first, expected) {
			t.Fatalf("overview missing %q: %s", expected, first)
		}
	}
}

func TestWorkflowOverviewListsReadyWorkflow(t *testing.T) {
	workdir := t.TempDir()
	sourceDir := filepath.Join(repoRoot(t), "workflows", "inventory-receipt")
	if _, err := runCLI(t, workdir, "workflow", "capture", sourceDir, "inventory-receipt"); err != nil {
		t.Fatalf("capture workflow: %v", err)
	}
	if _, err := runCLI(t, workdir, "workflow", "activate", "inventory-receipt"); err != nil {
		t.Fatalf("activate workflow: %v", err)
	}
	overview, err := runCLI(t, workdir, "workflow", "overview")
	if err != nil {
		t.Fatalf("workflow overview: %v", err)
	}
	if !strings.Contains(overview, "Ready: 1 active workflows") || !strings.Contains(overview, "[ready] inventory-receipt") {
		t.Fatalf("ready overview = %s", overview)
	}
}

func TestWorkflowOverviewPrioritizesSchemaReadinessBeforeActivation(t *testing.T) {
	workdir := t.TempDir()
	workflowDir := filepath.Join(t.TempDir(), "workflow")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatalf("make workflow directory: %v", err)
	}
	manifest := `{
  "id":"overview-schema",
  "version":"1.0.0",
  "summary":"Overview schema prerequisite test.",
  "required_packages":[],
  "required_protocols":[],
  "adapter":"inventory-receipt",
  "input_pcid":"bafkreie3xn5cs7in24a5aenl7kpyaa22e346wr4tcqm4evxgcn2v55yvne",
  "output_pcid":"bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne"
}`
	if err := os.WriteFile(filepath.Join(workflowDir, "workflow.json"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write workflow manifest: %v", err)
	}
	if _, err := runCLI(t, workdir, "workflow", "capture", workflowDir, "overview-schema"); err != nil {
		t.Fatalf("capture workflow: %v", err)
	}
	overview, err := runCLI(t, workdir, "workflow", "overview")
	if err != nil {
		t.Fatalf("workflow overview: %v", err)
	}
	if !strings.Contains(overview, "[workflow] overview-schema — canonical workflow schemas are not ready in CAS") || !strings.Contains(overview, "NEXT: moks workflow verify overview-schema") {
		t.Fatalf("schema-blocked overview = %s", overview)
	}
}

func TestWorkflowOverviewReportsDependencyFailureBeforeGenericAdapterBlocker(t *testing.T) {
	workdir := t.TempDir()
	workflowDir := filepath.Join(t.TempDir(), "workflow")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatalf("make workflow directory: %v", err)
	}
	manifest := `{
  "id":"overview-dependency",
  "version":"1.0.0",
  "summary":"Overview dependency prerequisite test.",
  "required_packages":["missing-package"],
  "required_protocols":[],
  "adapter":"missing-adapter",
  "input_pcid":"bafkreie3xn5cs7in24a5aenl7kpyaa22e346wr4tcqm4evxgcn2v55yvne",
  "output_pcid":"bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne"
}`
	if err := os.WriteFile(filepath.Join(workflowDir, "workflow.json"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write workflow manifest: %v", err)
	}
	if _, err := runCLI(t, workdir, "workflow", "capture", workflowDir, "overview-dependency"); err != nil {
		t.Fatalf("capture workflow: %v", err)
	}
	overview, err := runCLI(t, workdir, "workflow", "overview")
	if err != nil {
		t.Fatalf("workflow overview: %v", err)
	}
	if !strings.Contains(overview, "[workflow] overview-dependency — required package is not active: missing-package") || !strings.Contains(overview, "NEXT: moks workflow verify overview-dependency") {
		t.Fatalf("dependency-blocked overview = %s", overview)
	}
}

func TestWorkflowInboxCommandsInspectAndImportReceivedArtifact(t *testing.T) {
	workdir := t.TempDir()
	source, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open source: %v", err)
	}
	defer func() {
		if closeErr := source.Close(); closeErr != nil {
			t.Errorf("close source: %v", closeErr)
		}
	}()
	target, err := kernel.Open(filepath.Join(workdir, ".moks"))
	if err != nil {
		t.Fatalf("open target: %v", err)
	}
	if err := target.AllowPeer(grid.AllowedPeer{
		PeerID:            source.LocalPeerID(),
		BatchURL:          "https://source.invalid/relay/batch",
		ImportURL:         "https://source.invalid/relay/import",
		PublicKey:         source.LocalPeerPublicKey(),
		AllowPush:         true,
		AttesterClass:     "peer",
		AttestationWeight: 1,
		Federation:        "independent",
	}); err != nil {
		t.Fatalf("allow source: %v", err)
	}
	artifact, err := source.PutCAS([]byte("CLI inbox workflow artifact"))
	if err != nil {
		t.Fatalf("store source artifact: %v", err)
	}
	if err := source.ImportWorkflow(kernel.Workflow{ID: "source-workflow", ArtifactCID: artifact}); err != nil {
		t.Fatalf("import source workflow: %v", err)
	}
	transfer, err := source.ExportWorkflowTransfer("source-workflow")
	if err != nil {
		t.Fatalf("export source workflow: %v", err)
	}
	if err := target.ImportWorkflowTransferFromPeer(source.LocalPeerID(), transfer); err != nil {
		t.Fatalf("receive source workflow: %v", err)
	}
	if closeErr := target.Close(); closeErr != nil {
		t.Fatalf("close target: %v", closeErr)
	}

	list, err := runCLI(t, workdir, "workflow", "inbox", "list")
	if err != nil {
		t.Fatalf("list inbox: %v", err)
	}
	var entries []struct {
		ArtifactCID   string `json:"artifact_cid"`
		ReadyToImport bool   `json:"ready_to_import"`
	}
	if err := json.Unmarshal([]byte(list), &entries); err != nil {
		t.Fatalf("decode inbox list: %v", err)
	}
	if len(entries) != 1 || entries[0].ArtifactCID != transfer.ArtifactCID || !entries[0].ReadyToImport {
		t.Fatalf("inbox list = %s", list)
	}
	inspect, err := runCLI(t, workdir, "workflow", "inbox", "inspect", transfer.ArtifactCID)
	if err != nil || !strings.Contains(inspect, source.LocalPeerID()) {
		t.Fatalf("inspect inbox = %q, %v", inspect, err)
	}
	if _, err := runCLI(t, workdir, "workflow", "inbox", "import", transfer.ArtifactCID, "received-workflow"); err != nil {
		t.Fatalf("import inbox: %v", err)
	}
	workflows, err := runCLI(t, workdir, "workflow", "list")
	if err != nil || !strings.Contains(workflows, "received-workflow") {
		t.Fatalf("workflow list after inbox import = %q, %v", workflows, err)
	}
}

func TestRegistryCommandsPersistExactHostPolicy(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "registry", "allow", "REGISTRY.example:5000"); err != nil {
		t.Fatal(err)
	}
	output, err := runCLI(t, workdir, "registry", "list")
	if err != nil || output != "registry.example:5000" {
		t.Fatalf("registry list = %q, %v", output, err)
	}
	if _, err := runCLI(t, workdir, "registry", "remove", "registry.example:5000"); err != nil {
		t.Fatal(err)
	}
	output, err = runCLI(t, workdir, "registry", "list")
	if err != nil || output != "" {
		t.Fatalf("registry list after remove = %q, %v", output, err)
	}
}

func TestProcedureExecutionAdapterDockerEndToEnd(t *testing.T) {
	if os.Getenv("MOKS_DOCKER_INTEGRATION") != "1" {
		t.Skip("set MOKS_DOCKER_INTEGRATION=1 after building the pinned procedure-execution adapter image")
	}
	workdir := t.TempDir()
	root := repoRoot(t)
	packageDir := filepath.Join(root, "examples", "procedure-execution-adapter")
	workflowDir := filepath.Join(root, "workflows", "procedure-execution")
	if _, err := runCLI(t, workdir, "package", "install", packageDir); err != nil {
		t.Fatalf("install procedure adapter package: %v", err)
	}
	if _, err := runCLI(t, workdir, "procedures", "create", "proc-1", "DockCheck", "dock-intake"); err != nil {
		t.Fatalf("seed procedure: %v", err)
	}
	if _, err := runCLI(t, workdir, "workflow", "capture", workflowDir, "procedure-execution"); err != nil {
		t.Fatalf("capture workflow: %v", err)
	}
	if _, err := runCLI(t, workdir, "workflow", "activate", "procedure-execution"); err != nil {
		t.Fatalf("activate workflow: %v", err)
	}
	if _, err := runCLI(t, workdir, "workflow", "run", "start", "procedure-execution", "procedure_id", "proc-1", "run_id", "run-1", "actor", "alice", "outcome", "completed", "notes", "followed"); err != nil {
		t.Fatalf("run Docker adapter: %v", err)
	}
	procedure, err := runCLI(t, workdir, "procedures", "inspect", "proc-1")
	if err != nil {
		t.Fatalf("inspect procedure: %v", err)
	}
	if !strings.Contains(procedure, "uses: run-1") {
		t.Fatalf("procedure did not retain adapter use: %s", procedure)
	}
	runOutput, err := runCLI(t, workdir, "runs", "inspect", "run-1")
	if err != nil {
		t.Fatalf("inspect run: %v", err)
	}
	for _, expected := range []string{"item_id: proc-1", "actor: alice", "outcome: completed", "notes: followed"} {
		if !strings.Contains(runOutput, expected) {
			t.Fatalf("run inspect output missing %q: %s", expected, runOutput)
		}
	}
}

func TestWorkflowVerifyReportsMissingDependencyReadiness(t *testing.T) {
	workdir := t.TempDir()
	sourceDir := filepath.Join(t.TempDir(), "missing-dependency")
	if err := os.Mkdir(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"id":"missing-dependency","version":"1","summary":"missing dependency","required_packages":["missing"],"required_protocols":[]}`
	if err := os.WriteFile(filepath.Join(sourceDir, "workflow.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runCLI(t, workdir, "workflow", "capture", sourceDir, "missing-dependency"); err != nil {
		t.Fatal(err)
	}
	output, err := runCLI(t, workdir, "workflow", "verify", "missing-dependency")
	if err != nil {
		t.Fatal(err)
	}
	var verification struct {
		EligibleToExecute bool   `json:"eligible_to_execute"`
		Reason            string `json:"reason"`
	}
	if err := json.Unmarshal([]byte(output), &verification); err != nil {
		t.Fatal(err)
	}
	if verification.EligibleToExecute || verification.Reason != "required package is not active: missing" {
		t.Fatalf("missing dependency verification = %#v", verification)
	}
}

func TestWorkflowDemoMaintenanceRound(t *testing.T) {
	output, err := runCLI(t, repoRoot(t), "workflow", "demo", "maintenance-round")
	if err != nil {
		t.Fatalf("run maintenance round demo: %v", err)
	}
	for _, line := range []string{
		"[ok] captured maintenance-round as ",
		"[ok] verified manifest and local dependencies",
		"[ok] activated local workflow availability",
		"[ok] extracted the exact retained artifact for inspection",
	} {
		if !strings.Contains(output, line) {
			t.Fatalf("demo output missing %q: %s", line, output)
		}
	}
}

func TestWorkflowDemoReceivingCheck(t *testing.T) {
	output, err := runCLI(t, repoRoot(t), "workflow", "demo", "receiving-check")
	if err != nil {
		t.Fatalf("run receiving check demo: %v", err)
	}
	for _, line := range []string{
		"[ok] captured receiving-check as ",
		"[ok] verified manifest and local dependencies",
		"[ok] activated local workflow availability",
		"[ok] extracted the exact retained artifact for inspection",
	} {
		if !strings.Contains(output, line) {
			t.Fatalf("demo output missing %q: %s", line, output)
		}
	}
}

func TestWorkflowDemoTrainingQualification(t *testing.T) {
	output, err := runCLI(t, repoRoot(t), "workflow", "demo", "training-qualification")
	if err != nil {
		t.Fatalf("run training qualification demo: %v", err)
	}
	for _, line := range []string{
		"[ok] captured training-qualification as ",
		"[ok] verified manifest and local dependencies",
		"[ok] activated local workflow availability",
		"[ok] extracted the exact retained artifact for inspection",
	} {
		if !strings.Contains(output, line) {
			t.Fatalf("demo output missing %q: %s", line, output)
		}
	}
}

func TestWorkflowDemoInventoryDiscrepancyReview(t *testing.T) {
	output, err := runCLI(t, repoRoot(t), "workflow", "demo", "inventory-discrepancy-review")
	if err != nil {
		t.Fatalf("run inventory discrepancy review demo: %v", err)
	}
	for _, line := range []string{
		"[ok] captured inventory-discrepancy-review as ",
		"[ok] verified manifest and local dependencies",
		"[ok] activated local workflow availability",
		"[ok] extracted the exact retained artifact for inspection",
	} {
		if !strings.Contains(output, line) {
			t.Fatalf("demo output missing %q: %s", line, output)
		}
	}
}

func TestWorkflowDemoKnowledgeReview(t *testing.T) {
	output, err := runCLI(t, repoRoot(t), "workflow", "demo", "knowledge-review")
	if err != nil {
		t.Fatalf("run knowledge review demo: %v", err)
	}
	for _, line := range []string{
		"[ok] captured knowledge-review as ",
		"[ok] verified manifest and local dependencies",
		"[ok] activated local workflow availability",
		"[ok] extracted the exact retained artifact for inspection",
	} {
		if !strings.Contains(output, line) {
			t.Fatalf("demo output missing %q: %s", line, output)
		}
	}
}

func TestMultiWorkflowScenarioUsesSharedMainProgramRuntime(t *testing.T) {
	// Intent: Prove that separately loaded workflow artifacts execute their
	// declared built-in package commands over one shared runtime. Source: DI-lumek
	workdir := t.TempDir()
	for _, workflowID := range []string{
		"procedure-execution",
		"inventory-receipt",
		"maintenance-round",
		"receiving-check",
		"training-qualification",
		"inventory-discrepancy-review",
		"knowledge-review",
	} {
		sourceDir := filepath.Join(repoRoot(t), "workflows", workflowID)
		if _, err := runCLI(t, workdir, "workflow", "capture", sourceDir, workflowID); err != nil {
			t.Fatalf("capture %s: %v", workflowID, err)
		}
		if _, err := runCLI(t, workdir, "workflow", "activate", workflowID); err != nil {
			t.Fatalf("activate %s: %v", workflowID, err)
		}
	}
	commands := [][]string{
		{"context", "place", "create", "dock", "Dock", "Inbound"},
		{"context", "resource", "create", "scale", "Scale", "Bench-scale", "dock"},
		{"procedures", "create", "dock-procedure", "Dock-procedure", "Inspect-inbound-goods"},
		{"receiving", "create", "receipt", "dock", "Inbound-receipt", "Pallet-inspection"},
		{"receiving", "record-receipt", "receipt", "receive-run", "dock", "Alice", "accepted", "sealed"},
		{"receiving", "record-disposition", "receipt", "receipt-disposition", "accepted", "scale", "accepted"},
		{"inventory", "create", "stock", "dock", "Inbound-stock", "Count-after-receipt"},
		{"inventory", "record-count", "stock", "count-run", "dock", "Bob", "8", "counted", "counted"},
		{"inventory", "record-reconcile", "stock", "reconcile", "investigate", "scale", "variance"},
		{"maintenance", "create", "scale-check", "scale", "Scale-check", "Inspect-scale"},
		{"maintenance", "record-service", "scale-check", "maintenance-run", "scale", "Carol", "completed", "calibrated"},
		{"maintenance", "record-finding", "scale-check", "scale-finding", "scale", "accepted", "stable"},
		{"training", "create", "dock-training", "Dock-training", "Receiving-training"},
		{"training", "record-session", "dock-training", "training-run", "Dave", "Ellen", "completed", "demonstrated"},
		{"training", "certify", "dock-training", "training-certification", "Dave", "certified", "approved"},
		{"knowledge", "item", "create", "dock-guide", "procedure", "Dock-guide", "Receiving-guide"},
		{"knowledge", "revision", "snapshot", "dock-guide", "dock-guide-revision", "1", "Dock-guide", "Inspect-before-acceptance"},
		{"knowledge", "item", "approve", "dock-guide", "dock-guide-approval", "approved"},
	}
	for _, command := range commands {
		if _, err := runCLI(t, workdir, command...); err != nil {
			t.Fatalf("run %q: %v", command, err)
		}
	}
	workflows, err := runCLI(t, workdir, "workflow", "list")
	if err != nil {
		t.Fatalf("list workflows: %v", err)
	}
	if strings.Count(workflows, `"state": "active"`) != 7 {
		t.Fatalf("active workflow list = %s", workflows)
	}
	workflowRuns := [][]string{
		{"workflow", "run", "start", "procedure-execution", "procedure_id", "dock-procedure", "run_id", "workflow-procedure-run", "actor", "Alice", "outcome", "accepted", "notes", "sealed"},
		{"workflow", "run", "start", "inventory-receipt", "inventory_id", "stock", "run_id", "workflow-inventory-run", "place_id", "dock", "counter", "Bob", "quantity", "8", "outcome", "counted", "notes", "counted"},
		{"workflow", "run", "start", "maintenance-round", "maintenance_id", "scale-check", "run_id", "workflow-maintenance-run", "resource_id", "scale", "performer", "Carol", "outcome", "completed", "notes", "calibrated"},
		{"workflow", "run", "start", "receiving-check", "receiving_id", "receipt", "run_id", "workflow-receiving-run", "place_id", "dock", "receiver", "Alice", "outcome", "accepted", "notes", "sealed"},
		{"workflow", "run", "start", "training-qualification", "training_id", "dock-training", "run_id", "workflow-training-run", "trainee", "Dave", "instructor", "Ellen", "outcome", "completed", "notes", "demonstrated"},
		{"workflow", "run", "start", "inventory-discrepancy-review", "inventory_id", "stock", "event_id", "workflow-reconcile", "decision", "investigate", "resource_id", "scale", "notes", "variance"},
		{"workflow", "run", "start", "knowledge-review", "item_id", "dock-guide", "event_id", "workflow-guide-approval", "notes", "approved"},
	}
	for _, command := range workflowRuns {
		output, err := runCLI(t, workdir, command...)
		if err != nil || !strings.Contains(output, `"state": "completed"`) {
			t.Fatalf("run %q: %v\n%s", command, err, output)
		}
	}
	for _, inspect := range [][]string{
		{"receiving", "inspect", "receipt"},
		{"inventory", "inspect", "stock"},
		{"maintenance", "inspect", "scale-check"},
		{"training", "inspect", "dock-training"},
		{"knowledge", "item", "inspect", "dock-guide"},
	} {
		if _, err := runCLI(t, workdir, inspect...); err != nil {
			t.Fatalf("inspect %q: %v", inspect, err)
		}
	}
}

func TestWorkflowRelayEndpointTransfersArtifactWithoutActivatingIt(t *testing.T) {
	source, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open source: %v", err)
	}
	defer func() {
		if closeErr := source.Close(); closeErr != nil {
			t.Errorf("close source: %v", closeErr)
		}
	}()
	target, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open target: %v", err)
	}
	defer func() {
		if closeErr := target.Close(); closeErr != nil {
			t.Errorf("close target: %v", closeErr)
		}
	}()
	if err := target.AllowPeer(grid.AllowedPeer{
		PeerID:            source.LocalPeerID(),
		BatchURL:          "https://source.invalid/relay/batch",
		ImportURL:         "https://source.invalid/relay/import",
		PublicKey:         source.LocalPeerPublicKey(),
		AllowPush:         true,
		AttesterClass:     "peer",
		AttestationWeight: 1,
		Federation:        "independent",
	}); err != nil {
		t.Fatalf("allow source: %v", err)
	}
	artifactCID, err := source.PutCAS([]byte("workflow relay endpoint artifact"))
	if err != nil {
		t.Fatalf("put source artifact: %v", err)
	}
	if err := source.ImportWorkflow(kernel.Workflow{ID: "endpoint-handoff", ArtifactCID: artifactCID}); err != nil {
		t.Fatalf("import source workflow: %v", err)
	}
	sourceWorkflow := source.Workflows()[0]
	server := httptest.NewServer(relayHandler(context.Background(), target))
	defer server.Close()
	if err := source.AllowPeer(grid.AllowedPeer{
		PeerID:            target.LocalPeerID(),
		BatchURL:          server.URL + "/relay/batch",
		ImportURL:         server.URL + "/relay/import",
		WorkflowImportURL: server.URL + "/relay/workflow/import",
		PublicKey:         target.LocalPeerPublicKey(),
		AllowPush:         true,
		AttesterClass:     "peer",
		AttestationWeight: 1,
		Federation:        "independent",
	}); err != nil {
		t.Fatalf("allow target: %v", err)
	}
	if err := workflowRelayPush(context.Background(), source, "endpoint-handoff", target.LocalPeerID()); err != nil {
		t.Fatalf("push workflow: %v", err)
	}
	if workflows := target.Workflows(); len(workflows) != 0 {
		t.Fatalf("workflow endpoint changed target lifecycle: %#v", workflows)
	}
	if _, err := target.GetCAS(sourceWorkflow.ArtifactCID); err != nil {
		t.Fatalf("target missing transferred artifact: %v", err)
	}
}

// Intent: Prove that two independently rooted, deployed relay processes
// exchange exact workflow evidence without letting receipt change the
// receiver's workflow lifecycle. Source: DI-novuk
func TestDeployedTwoNodeWorkflowRelay(t *testing.T) {
	if os.Getenv("MOKS_DEPLOYED_RELAY_INTEGRATION") != "1" {
		t.Skip("set MOKS_DEPLOYED_RELAY_INTEGRATION=1 to build and run two local moks relay processes")
	}
	root := repoRoot(t)
	binary := filepath.Join(t.TempDir(), "moks")
	build := exec.Command("go", "build", "-o", binary, "./cmd/moks")
	build.Dir = root
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build moks binary: %v: %s", err, output)
	}
	aliceDir := t.TempDir()
	bobDir := t.TempDir()
	workflowDir := filepath.Join(root, "workflows", "procedure-execution")
	capture := runDeployedMoks(t, binary, aliceDir, "workflow", "capture", workflowDir, "procedure-execution")
	var captured []kernel.Workflow
	if err := json.Unmarshal([]byte(capture), &captured); err != nil || len(captured) != 1 {
		t.Fatalf("decode captured workflow %q: %v", capture, err)
	}
	if _, err := runDeployedMoksResult(binary, aliceDir, "workflow", "activate", "procedure-execution"); err != nil {
		t.Fatalf("activate Alice workflow: %v", err)
	}
	aliceRuntime, err := kernel.Open(filepath.Join(aliceDir, ".moks"))
	if err != nil {
		t.Fatalf("open Alice runtime for transfer inspection: %v", err)
	}
	transfer, err := aliceRuntime.ExportWorkflowTransfer("procedure-execution")
	if closeErr := aliceRuntime.Close(); closeErr != nil {
		t.Fatalf("close Alice transfer inspection runtime: %v", closeErr)
	}
	if err != nil {
		t.Fatalf("export Alice workflow transfer: %v", err)
	}
	bobAddress := reserveDeployedRelayAddress(t)
	aliceCard := deployedRelayPeerCard(t, aliceDir, "127.0.0.1:1")
	bobCard := deployedRelayPeerCard(t, bobDir, bobAddress)
	allowDeployedRelayPeer(t, binary, aliceDir, bobCard, "no-pull", "push")
	allowDeployedRelayPeer(t, binary, bobDir, aliceCard, "no-pull", "push")
	bob := startDeployedRelayNode(t, binary, bobDir, bobAddress)
	push := runDeployedMoks(t, binary, aliceDir, "workflow", "relay", "push", "procedure-execution", bob.card.PeerID)
	if !strings.Contains(push, "workflow relay pushed procedure-execution to "+bob.card.PeerID) {
		t.Fatalf("unexpected relay push output: %s", push)
	}
	bob.stop(t)
	bobWorkflows := runDeployedMoks(t, binary, bobDir, "workflow", "list")
	if strings.TrimSpace(bobWorkflows) != "[]" {
		t.Fatalf("workflow receipt changed Bob lifecycle: %s", bobWorkflows)
	}
	if captured[0].ArtifactCID == "" {
		t.Fatal("capture did not produce an artifact CID")
	}
	bobRuntime, err := kernel.Open(filepath.Join(bobDir, ".moks"))
	if err != nil {
		t.Fatalf("open Bob runtime for transfer inspection: %v", err)
	}
	artifact, err := bobRuntime.GetCAS(captured[0].ArtifactCID)
	if closeErr := bobRuntime.Close(); closeErr != nil {
		t.Fatalf("close Bob transfer inspection runtime: %v", closeErr)
	}
	if err != nil {
		t.Fatalf("read Bob transferred artifact: %v", err)
	}
	if !bytes.Equal(artifact, transfer.Artifact) {
		t.Fatal("Bob transferred artifact bytes differ from Alice export")
	}
	evidenceStore, err := store.OpenCAS(filepath.Join(bobDir, ".moks", "workflow-evidence"))
	if err != nil {
		t.Fatalf("open Bob workflow evidence CAS: %v", err)
	}
	evidenceCIDs, err := evidenceStore.ListCIDs()
	if err != nil {
		t.Fatalf("list Bob workflow evidence: %v", err)
	}
	if len(evidenceCIDs) != 1 {
		t.Fatalf("Bob workflow evidence count = %d, want 1", len(evidenceCIDs))
	}
	evidence, err := evidenceStore.GetCID(evidenceCIDs[0])
	if err != nil {
		t.Fatalf("read Bob workflow evidence: %v", err)
	}
	if !bytes.Equal(evidence, transfer.LifecycleEvent) {
		t.Fatal("Bob lifecycle evidence bytes differ from Alice export")
	}
	// Intent: Exercise Bob's explicit portable-image acquisition only when the
	// caller supplies both Docker permission and a durable registry package.
	// Source: DI-zivut
	if os.Getenv("MOKS_PORTABLE_REGISTRY_INTEGRATION") != "1" || os.Getenv("MOKS_DOCKER_INTEGRATION") != "1" {
		return
	}
	packageDir := os.Getenv("MOKS_PORTABLE_ADAPTER_PACKAGE")
	if packageDir == "" {
		t.Skip("set MOKS_PORTABLE_ADAPTER_PACKAGE to a package pinned to the supplied portable registry digest")
	}
	registryHost := os.Getenv("MOKS_PORTABLE_REGISTRY_HOST")
	if registryHost == "" {
		t.Skip("set MOKS_PORTABLE_REGISTRY_HOST to the exact registry host named by the supplied package")
	}
	if _, err := runDeployedMoksResult(binary, bobDir, "package", "install", packageDir); err != nil {
		t.Fatalf("install Bob procedure adapter package: %v", err)
	}
	if _, err := runDeployedMoksResult(binary, bobDir, "registry", "allow", registryHost); err != nil {
		t.Fatalf("allow Bob procedure adapter registry: %v", err)
	}
	if _, err := runDeployedMoksResult(binary, bobDir, "workflow", "import", "procedure-execution", captured[0].ArtifactCID); err != nil {
		t.Fatalf("import Bob transferred workflow: %v", err)
	}
	if _, err := runDeployedMoksResult(binary, bobDir, "workflow", "activate", "procedure-execution"); err != nil {
		t.Fatalf("activate Bob workflow: %v", err)
	}
	if _, err := runDeployedMoksResult(binary, bobDir, "workflow", "image", "pull", "procedure-execution"); err != nil {
		t.Fatalf("pull Bob adapter image: %v", err)
	}
	if _, err := runDeployedMoksResult(binary, bobDir, "procedures", "create", "proc-1", "DockCheck", "dock-intake"); err != nil {
		t.Fatalf("seed Bob procedure: %v", err)
	}
	if _, err := runDeployedMoksResult(binary, bobDir, "workflow", "run", "start", "procedure-execution", "procedure_id", "proc-1", "run_id", "run-1", "actor", "bob", "outcome", "completed", "notes", "received"); err != nil {
		t.Fatalf("run Bob pulled adapter: %v", err)
	}
	runOutput := runDeployedMoks(t, binary, bobDir, "runs", "inspect", "run-1")
	if !strings.Contains(runOutput, "actor: bob") || !strings.Contains(runOutput, "notes: received") {
		t.Fatalf("Bob run did not retain adapter result: %s", runOutput)
	}
}

type deployedRelayNode struct {
	card    grid.PeerCard
	command *exec.Cmd
}

func reserveDeployedRelayAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve relay address: %v", err)
	}
	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("release relay address: %v", err)
	}
	return address
}

func deployedRelayPeerCard(t *testing.T, workdir string, address string) grid.PeerCard {
	t.Helper()
	runtime, err := kernel.Open(filepath.Join(workdir, ".moks"))
	if err != nil {
		t.Fatalf("open relay node runtime: %v", err)
	}
	defer func() {
		if closeErr := runtime.Close(); closeErr != nil {
			t.Errorf("close relay node runtime: %v", closeErr)
		}
	}()
	baseURL := "http://" + address
	return grid.PeerCard{
		PeerID:            runtime.LocalPeerID(),
		PublicKey:         runtime.LocalPeerPublicKey(),
		BatchURL:          baseURL + "/relay/batch",
		ImportURL:         baseURL + "/relay/import",
		WorkflowImportURL: baseURL + "/relay/workflow/import",
		DiscoverURL:       baseURL + "/relay/peer-card",
	}
}

func startDeployedRelayNode(t *testing.T, binary string, workdir string, address string) deployedRelayNode {
	t.Helper()
	command := exec.Command(binary, "relay", "serve", address)
	command.Dir = workdir
	if err := command.Start(); err != nil {
		t.Fatalf("start relay node: %v", err)
	}
	node := deployedRelayNode{command: command}
	t.Cleanup(func() { node.stop(t) })
	cardURL := "http://" + address + "/relay/peer-card"
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		response, err := http.Get(cardURL)
		if err == nil {
			if response.StatusCode == http.StatusOK {
				var card grid.PeerCard
				decodeErr := json.NewDecoder(response.Body).Decode(&card)
				if closeErr := response.Body.Close(); closeErr != nil {
					t.Fatalf("close relay peer card response: %v", closeErr)
				}
				if decodeErr == nil {
					node.card = card
					return node
				}
			} else if closeErr := response.Body.Close(); closeErr != nil {
				t.Fatalf("close unsuccessful relay peer card response: %v", closeErr)
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("relay node did not serve peer card at %s", cardURL)
	return deployedRelayNode{}
}

func (node deployedRelayNode) stop(t *testing.T) {
	t.Helper()
	if node.command == nil || (node.command.ProcessState != nil && node.command.ProcessState.Exited()) {
		return
	}
	if err := node.command.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		t.Errorf("stop relay node: %v", err)
	}
	if err := node.command.Wait(); err != nil && node.command.ProcessState == nil {
		t.Errorf("wait for relay node: %v", err)
	}
}

func allowDeployedRelayPeer(t *testing.T, binary string, workdir string, peer grid.PeerCard, pull string, push string) {
	t.Helper()
	if _, err := runDeployedMoksResult(binary, workdir, "relay", "peer", "allow", peer.PeerID, peer.BatchURL, peer.ImportURL, peer.PublicKey, pull, push); err != nil {
		t.Fatalf("allow peer %s: %v", peer.PeerID, err)
	}
}

func runDeployedMoks(t *testing.T, binary string, workdir string, args ...string) string {
	t.Helper()
	output, err := runDeployedMoksResult(binary, workdir, args...)
	if err != nil {
		t.Fatalf("run moks %q: %v", args, err)
	}
	return output
}

func runDeployedMoksResult(binary string, workdir string, args ...string) (string, error) {
	command := exec.Command(binary, args...)
	command.Dir = workdir
	output, err := command.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func TestInstalledWriterAgentExample(t *testing.T) {
	workdir := t.TempDir()
	exampleDir := filepath.Join(repoRoot(t), "examples", "writer-agent")
	if output, err := runCLI(t, workdir, "package", "install", exampleDir); err != nil {
		t.Fatalf("install writer agent: %v", err)
	} else if !strings.Contains(output, "installed writer-agent") {
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

func TestRouteListShowsProtocolRoutes(t *testing.T) {
	workdir := t.TempDir()
	output, err := runCLI(t, workdir, "route", "list")
	if err != nil {
		t.Fatalf("route list: %v", err)
	}
	if !strings.Contains(output, "pcid:moks.context.place.v1\tfamily-validator\tdirect\tcontext\t0.1.0\tmoks.context.place.v1\t") {
		t.Fatalf("route list missing context place validator: %s", output)
	}
	if !strings.Contains(output, "pcid:moks.ops.note.v1\tfamily-validator\tdirect\tops-note\t0.1.0\tmoks.ops.note.v1\t") {
		t.Fatalf("route list missing ops note validator: %s", output)
	}
}

func TestRouteInspectShowsProtocolRoutesForPCID(t *testing.T) {
	workdir := t.TempDir()
	output, err := runCLI(t, workdir, "route", "inspect", "pcid:moks.context.place.v1")
	if err != nil {
		t.Fatalf("route inspect: %v", err)
	}
	if !strings.Contains(output, `"protocol_pcid": "pcid:moks.context.place.v1"`) {
		t.Fatalf("route inspect missing protocol: %s", output)
	}
	if !strings.Contains(output, `"role": "family-validator"`) {
		t.Fatalf("route inspect missing role: %s", output)
	}
	if !strings.Contains(output, `"route_type": "direct"`) {
		t.Fatalf("route inspect missing route type: %s", output)
	}
}

func TestRoutePlanShowsPreferredRouteForPCID(t *testing.T) {
	workdir := t.TempDir()
	publishContextRoutePromises(t, workdir)
	output, err := runCLI(t, workdir, "route", "plan", "pcid:moks.context.place.v1")
	if err != nil {
		t.Fatalf("route plan: %v", err)
	}
	if !strings.Contains(output, `"protocol_pcid": "pcid:moks.context.place.v1"`) {
		t.Fatalf("route plan missing protocol: %s", output)
	}
	if !strings.Contains(output, `"preferred"`) {
		t.Fatalf("route plan missing preferred plan: %s", output)
	}
	if !strings.Contains(output, `"executable": true`) {
		t.Fatalf("route plan missing executable preferred route: %s", output)
	}
	if !strings.Contains(output, `"explanation"`) || !strings.Contains(output, `"winner"`) {
		t.Fatalf("route plan missing explanation: %s", output)
	}
	if strings.Contains(output, `"comparisons"`) && !strings.Contains(output, `"decision_path"`) {
		t.Fatalf("route plan comparison detail is incomplete: %s", output)
	}
}

func TestRoutePlanTraceShowsPlannerSteps(t *testing.T) {
	workdir := t.TempDir()
	publishContextRoutePromises(t, workdir)
	output, err := runCLI(t, workdir, "route", "plan", "pcid:moks.context.place.v1", "trace")
	if err != nil {
		t.Fatalf("route plan trace: %v", err)
	}
	if !strings.Contains(output, `"trace"`) || !strings.Contains(output, `"plan-start"`) {
		t.Fatalf("route plan trace missing trace steps: %s", output)
	}
	if !strings.Contains(output, `"preferred"`) {
		t.Fatalf("route plan trace missing preferred route: %s", output)
	}
	if !strings.Contains(output, `"trace_summary"`) || !strings.Contains(output, `"total_steps"`) {
		t.Fatalf("route plan trace missing trace summary: %s", output)
	}
	if !strings.Contains(output, `"scope": "root"`) || !strings.Contains(output, `"protocol_pcid": "pcid:moks.context.place.v1"`) {
		t.Fatalf("route plan trace missing scope metadata: %s", output)
	}
	if !strings.Contains(output, `"hop_path": "root"`) {
		t.Fatalf("route plan trace missing root hop path: %s", output)
	}
	if !strings.Contains(output, `"hop_summary": "root"`) {
		t.Fatalf("route plan trace missing root hop summary: %s", output)
	}
	if !strings.Contains(output, `"hop_depth": 0`) || !strings.Contains(output, `"hop_index": 0`) {
		t.Fatalf("route plan trace missing root hop depth/index: %s", output)
	}
}

func TestRoutePlanTraceCanFocusOnCandidate(t *testing.T) {
	workdir := t.TempDir()
	output, err := runCLI(t, workdir, "route", "plan", "pcid:moks.context.place.v1", "trace", "candidate", "context:family-validator:direct")
	if err != nil {
		t.Fatalf("route plan trace candidate focus: %v", err)
	}
	if !strings.Contains(output, `"trace"`) || !strings.Contains(output, `context:family-validator:direct`) {
		t.Fatalf("focused candidate trace missing target: %s", output)
	}
	if !strings.Contains(output, `"hidden_steps"`) || !strings.Contains(output, `"filter"`) {
		t.Fatalf("focused candidate trace missing filter summary: %s", output)
	}
}

func TestRoutePlanTraceCanFocusOnDepth(t *testing.T) {
	workdir := t.TempDir()
	output, err := runCLI(t, workdir, "route", "plan", "pcid:moks.context.place.v1", "trace", "depth", "0")
	if err != nil {
		t.Fatalf("route plan trace depth focus: %v", err)
	}
	if !strings.Contains(output, `"filter"`) || !strings.Contains(output, `"target": "0"`) {
		t.Fatalf("focused depth trace missing filter target: %s", output)
	}
}

func TestRoutePlanTraceCanCombineFilters(t *testing.T) {
	workdir := t.TempDir()
	output, err := runCLI(t, workdir, "route", "plan", "pcid:moks.context.place.v1", "trace", "depth", "0", "candidate", "context:family-validator:direct")
	if err != nil {
		t.Fatalf("route plan trace combined filters: %v", err)
	}
	if !strings.Contains(output, `"clauses"`) || !strings.Contains(output, `"kind": "depth"`) || !strings.Contains(output, `"kind": "candidate"`) {
		t.Fatalf("combined trace filter missing clauses: %s", output)
	}
}

func TestRoutePlanTraceCanUseNamedScope(t *testing.T) {
	workdir := t.TempDir()
	output, err := runCLI(t, workdir, "route", "plan", "pcid:moks.context.place.v1", "trace", "scope", "direct-hops")
	if err != nil {
		t.Fatalf("route plan trace named scope: %v", err)
	}
	if !strings.Contains(output, `"kind": "scope"`) || !strings.Contains(output, `"target": "direct-hops"`) {
		t.Fatalf("named scope trace missing scope clause: %s", output)
	}
}

func TestRoutePlanTraceCanUseLocalScopeAlias(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "route", "scope", "set", "root-only", "depth", "0"); err != nil {
		t.Fatalf("route scope set: %v", err)
	}
	listOutput, err := runCLI(t, workdir, "route", "scope", "list")
	if err != nil {
		t.Fatalf("route scope list: %v", err)
	}
	if !strings.Contains(listOutput, `"name": "root-only"`) {
		t.Fatalf("route scope list missing local alias: %s", listOutput)
	}
	output, err := runCLI(t, workdir, "route", "plan", "pcid:moks.context.place.v1", "trace", "scope", "root-only")
	if err != nil {
		t.Fatalf("route plan trace local scope alias: %v", err)
	}
	if !strings.Contains(output, `"target": "root-only"`) || !strings.Contains(output, `"hop_depth": 0`) {
		t.Fatalf("local scope alias trace missing alias filter or root depth: %s", output)
	}
	if _, err := runCLI(t, workdir, "route", "scope", "remove", "root-only"); err != nil {
		t.Fatalf("route scope remove: %v", err)
	}
}

func TestRouteScopeInspectShowsExpandedAliasClauses(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "route", "scope", "set", "branch-base", "scope", "direct-hops", "candidate", "context:family-validator:direct"); err != nil {
		t.Fatalf("route scope set branch-base: %v", err)
	}
	if _, err := runCLI(t, workdir, "route", "scope", "set", "branch-expanded", "scope", "branch-base", "downstream", "pcid:moks.context.place.v1"); err != nil {
		t.Fatalf("route scope set branch-expanded: %v", err)
	}
	output, err := runCLI(t, workdir, "route", "scope", "inspect", "branch-expanded")
	if err != nil {
		t.Fatalf("route scope inspect: %v", err)
	}
	if !strings.Contains(output, `"name": "branch-expanded"`) || !strings.Contains(output, `"raw_clauses"`) || !strings.Contains(output, `"expanded_clauses"`) {
		t.Fatalf("route scope inspect missing inspection fields: %s", output)
	}
	if !strings.Contains(output, `"kind": "scope"`) || !strings.Contains(output, `"target": "branch-base"`) {
		t.Fatalf("route scope inspect missing raw alias composition: %s", output)
	}
	if !strings.Contains(output, `"kind": "depth"`) || !strings.Contains(output, `"target": "1"`) {
		t.Fatalf("route scope inspect missing expanded built-in clause: %s", output)
	}
	if !strings.Contains(output, `"expanded_details"`) || !strings.Contains(output, `"provenance"`) {
		t.Fatalf("route scope inspect missing expanded provenance: %s", output)
	}
	if !strings.Contains(output, `"branch-expanded"`) || !strings.Contains(output, `"branch-base"`) || !strings.Contains(output, `"direct-hops"`) {
		t.Fatalf("route scope inspect missing provenance chain: %s", output)
	}
	if !strings.Contains(output, `"groups"`) || !strings.Contains(output, `"branch"`) {
		t.Fatalf("route scope inspect missing grouped branches: %s", output)
	}
	if !strings.Contains(output, `"label": "branch-1"`) || !strings.Contains(output, `"summary": "branch-expanded`) || !strings.Contains(output, `branch-base`) || !strings.Contains(output, `direct-hops`) {
		t.Fatalf("route scope inspect missing branch label or summary: %s", output)
	}
}

func TestRouteScopeInspectShowsSkippedBranches(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "route", "scope", "set", "cycle-a", "scope", "cycle-b", "depth", "1"); err != nil {
		t.Fatalf("route scope set cycle-a: %v", err)
	}
	if _, err := runCLI(t, workdir, "route", "scope", "set", "cycle-b", "scope", "cycle-a"); err != nil {
		t.Fatalf("route scope set cycle-b: %v", err)
	}
	output, err := runCLI(t, workdir, "route", "scope", "inspect", "cycle-a")
	if err != nil {
		t.Fatalf("route scope inspect cycle-a: %v", err)
	}
	if !strings.Contains(output, `"skipped"`) || !strings.Contains(output, `"reason": "cycle"`) {
		t.Fatalf("route scope inspect missing cycle skip diagnostics: %s", output)
	}
	if !strings.Contains(output, `"branch": [`) || !strings.Contains(output, `"cycle-a"`) || !strings.Contains(output, `"cycle-b"`) {
		t.Fatalf("route scope inspect missing cycle branch attachment: %s", output)
	}
	if _, err := runCLI(t, workdir, "route", "scope", "set", "dangling", "scope", "missing-scope", "candidate", "context:family-validator:direct"); err != nil {
		t.Fatalf("route scope set dangling: %v", err)
	}
	output, err = runCLI(t, workdir, "route", "scope", "inspect", "dangling")
	if err != nil {
		t.Fatalf("route scope inspect dangling: %v", err)
	}
	if !strings.Contains(output, `"reason": "unknown-scope"`) || !strings.Contains(output, `"scope": "missing-scope"`) {
		t.Fatalf("route scope inspect missing unknown-scope diagnostics: %s", output)
	}
	if !strings.Contains(output, `"groups"`) || !strings.Contains(output, `"dangling"`) {
		t.Fatalf("route scope inspect missing grouped skip attachment: %s", output)
	}
}

func TestRouteScopeInspectCanOrderAndFilterGroups(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "route", "scope", "set", "branch-base", "scope", "direct-hops", "candidate", "context:family-validator:direct"); err != nil {
		t.Fatalf("route scope set branch-base: %v", err)
	}
	if _, err := runCLI(t, workdir, "route", "scope", "set", "branch-expanded", "scope", "branch-base", "downstream", "pcid:moks.context.place.v1"); err != nil {
		t.Fatalf("route scope set branch-expanded: %v", err)
	}
	output, err := runCLI(t, workdir, "route", "scope", "inspect", "branch-expanded", "sort", "summary", "summary", "branch-base")
	if err != nil {
		t.Fatalf("route scope inspect with query: %v", err)
	}
	if !strings.Contains(output, `"active_query"`) || !strings.Contains(output, `"sort_by": "summary"`) || !strings.Contains(output, `"summary_filter": "branch-base"`) {
		t.Fatalf("route scope inspect missing active query echo: %s", output)
	}
	if !strings.Contains(output, `"query_summary"`) || !strings.Contains(output, `"matched_groups": 2`) || !strings.Contains(output, `"hidden_groups": 1`) || !strings.Contains(output, `"ordering": "summary"`) {
		t.Fatalf("route scope inspect missing query summary: %s", output)
	}
	if !strings.Contains(output, `"groups"`) || !strings.Contains(output, `"branch-base"`) {
		t.Fatalf("route scope inspect missing filtered grouped branches: %s", output)
	}
	if strings.Contains(output, `"summary": "branch-expanded"`) {
		t.Fatalf("route scope inspect summary filter kept unexpected root-only branch: %s", output)
	}
	depthOutput, err := runCLI(t, workdir, "route", "scope", "inspect", "branch-expanded", "depth", "3")
	if err != nil {
		t.Fatalf("route scope inspect with depth filter: %v", err)
	}
	if !strings.Contains(depthOutput, `"active_query"`) || !strings.Contains(depthOutput, `"depth_filter": "3"`) {
		t.Fatalf("route scope inspect missing depth query echo: %s", depthOutput)
	}
	if !strings.Contains(depthOutput, `"query_summary"`) || !strings.Contains(depthOutput, `"matched_groups": 1`) || !strings.Contains(depthOutput, `"hidden_groups": 2`) || !strings.Contains(depthOutput, `"ordering": "label"`) {
		t.Fatalf("route scope inspect missing depth query summary: %s", depthOutput)
	}
	if !strings.Contains(depthOutput, `"query_diagnostics"`) || !strings.Contains(depthOutput, `"default_ordering_applied": true`) || !strings.Contains(depthOutput, `"default_ordering_reason": "no sort provided; defaulted to label ordering"`) {
		t.Fatalf("route scope inspect missing depth query diagnostics: %s", depthOutput)
	}
	if !strings.Contains(depthOutput, `"depth": 3`) || strings.Contains(depthOutput, `"depth": 2`) {
		t.Fatalf("route scope inspect depth filter produced unexpected groups: %s", depthOutput)
	}
	zeroOutput, err := runCLI(t, workdir, "route", "scope", "inspect", "branch-expanded", "depth", "9")
	if err != nil {
		t.Fatalf("route scope inspect with zero-match filter: %v", err)
	}
	if !strings.Contains(zeroOutput, `"matched_groups": 0`) || !strings.Contains(zeroOutput, `"hidden_groups": 3`) || !strings.Contains(zeroOutput, `"zero_matches": true`) || !strings.Contains(zeroOutput, `"zero_match_reason": "no grouped branches matched filters: depth=9"`) {
		t.Fatalf("route scope inspect missing zero-match diagnostics: %s", zeroOutput)
	}
	invalidOutput, err := runCLI(t, workdir, "route", "scope", "inspect", "branch-expanded", "depth", "abc")
	if err != nil {
		t.Fatalf("route scope inspect with invalid depth filter: %v", err)
	}
	if !strings.Contains(invalidOutput, `"matched_groups": 3`) || !strings.Contains(invalidOutput, `"hidden_groups": 0`) || !strings.Contains(invalidOutput, `"ignored_filters": [`) || !strings.Contains(invalidOutput, `"depth filter \"abc\" ignored: expected \u003cn\u003e or \u003cn+\u003e"`) {
		t.Fatalf("route scope inspect missing invalid-depth diagnostics: %s", invalidOutput)
	}
	invalidSortOutput, err := runCLI(t, workdir, "route", "scope", "inspect", "branch-expanded", "sort", "weird")
	if err != nil {
		t.Fatalf("route scope inspect with invalid sort: %v", err)
	}
	if !strings.Contains(invalidSortOutput, `"matched_groups": 3`) || !strings.Contains(invalidSortOutput, `"hidden_groups": 0`) || !strings.Contains(invalidSortOutput, `"ordering": "label"`) || !strings.Contains(invalidSortOutput, `"default_ordering_reason": "invalid sort ignored; defaulted to label ordering"`) || !strings.Contains(invalidSortOutput, `"sort \"weird\" ignored: expected depth, label, or summary"`) {
		t.Fatalf("route scope inspect missing invalid-sort diagnostics: %s", invalidSortOutput)
	}
	invalidTextOutput, err := runCLI(t, workdir, "route", "scope", "inspect", "branch-expanded", "label", "   ", "summary", "\t")
	if err != nil {
		t.Fatalf("route scope inspect with invalid text filters: %v", err)
	}
	if !strings.Contains(invalidTextOutput, `"matched_groups": 3`) || !strings.Contains(invalidTextOutput, `"hidden_groups": 0`) || !strings.Contains(invalidTextOutput, `"label_filter": "   "`) || !strings.Contains(invalidTextOutput, `"summary_filter": "\t"`) || !strings.Contains(invalidTextOutput, `"label filter \"   \" ignored: expected non-whitespace text"`) || !strings.Contains(invalidTextOutput, `"summary filter \"\\t\" ignored: expected non-whitespace text"`) {
		t.Fatalf("route scope inspect missing invalid-text diagnostics: %s", invalidTextOutput)
	}
}

func TestRoutePolicySetAndShow(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "route", "policy", "set", "parser", "direct", "parser", "-"); err != nil {
		t.Fatalf("route policy set: %v", err)
	}
	output, err := runCLI(t, workdir, "route", "policy", "show")
	if err != nil {
		t.Fatalf("route policy show: %v", err)
	}
	if !strings.Contains(output, `"prefer_route_types": [`) || !strings.Contains(output, `"parser"`) {
		t.Fatalf("route policy show missing prefer_route_types: %s", output)
	}
	if !strings.Contains(output, `"avoid_route_types": [`) || !strings.Contains(output, `"direct"`) {
		t.Fatalf("route policy show missing avoid_route_types: %s", output)
	}
}

func TestRoutePolicySetForProtocolAndShowEffective(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "route", "policy", "set", "direct", "-", "-", "-"); err != nil {
		t.Fatalf("route policy set: %v", err)
	}
	if _, err := runCLI(t, workdir, "route", "policy", "set-for", "pcid:moks.context.place.v1", "parser", "direct", "-", "-"); err != nil {
		t.Fatalf("route policy set-for: %v", err)
	}
	output, err := runCLI(t, workdir, "route", "policy", "show")
	if err != nil {
		t.Fatalf("route policy show: %v", err)
	}
	if !strings.Contains(output, `"protocol_pcid": "pcid:moks.context.place.v1"`) {
		t.Fatalf("route policy show missing protocol override: %s", output)
	}
	effective, err := runCLI(t, workdir, "route", "policy", "show", "pcid:moks.context.place.v1")
	if err != nil {
		t.Fatalf("route policy show for protocol: %v", err)
	}
	if !strings.Contains(effective, `"protocol_pcid": "pcid:moks.context.place.v1"`) {
		t.Fatalf("route policy effective output missing protocol: %s", effective)
	}
	if !strings.Contains(effective, `"global"`) || !strings.Contains(effective, `"effective"`) {
		t.Fatalf("route policy effective output missing global/effective blocks: %s", effective)
	}
	if !strings.Contains(effective, `"prefer_route_types": [`) || !strings.Contains(effective, `"parser"`) {
		t.Fatalf("route policy effective output missing parser override: %s", effective)
	}
	if !strings.Contains(effective, `"avoid_route_types": [`) || !strings.Contains(effective, `"direct"`) {
		t.Fatalf("route policy effective output missing direct avoidance: %s", effective)
	}
	if _, err := runCLI(t, workdir, "route", "policy", "remove", "pcid:moks.context.place.v1"); err != nil {
		t.Fatalf("route policy remove: %v", err)
	}
	afterRemove, err := runCLI(t, workdir, "route", "policy", "show")
	if err != nil {
		t.Fatalf("route policy show after remove: %v", err)
	}
	if strings.Contains(afterRemove, `"protocol_pcid": "pcid:moks.context.place.v1"`) {
		t.Fatalf("route policy override still present after remove: %s", afterRemove)
	}
}

func TestRoutePolicySetForRoleAndShowEffective(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "route", "policy", "set-for-role", "pcid:moks.context.place.v1", "family-validator", "-", "-", "family-validator", "-"); err != nil {
		t.Fatalf("route policy set-for-role: %v", err)
	}
	output, err := runCLI(t, workdir, "route", "policy", "show")
	if err != nil {
		t.Fatalf("route policy show: %v", err)
	}
	if !strings.Contains(output, `"protocol_roles": [`) || !strings.Contains(output, `"role": "family-validator"`) {
		t.Fatalf("route policy show missing protocol role override: %s", output)
	}
	effective, err := runCLI(t, workdir, "route", "policy", "show", "pcid:moks.context.place.v1", "family-validator")
	if err != nil {
		t.Fatalf("route policy show for protocol role: %v", err)
	}
	if !strings.Contains(effective, `"role": "family-validator"`) {
		t.Fatalf("route policy effective output missing role: %s", effective)
	}
	if !strings.Contains(effective, `"protocol"`) || !strings.Contains(effective, `"effective"`) {
		t.Fatalf("route policy effective output missing protocol/effective blocks: %s", effective)
	}
	if !strings.Contains(effective, `"prefer_roles": [`) || !strings.Contains(effective, `"family-validator"`) {
		t.Fatalf("route policy effective output missing family-validator preference: %s", effective)
	}
	if _, err := runCLI(t, workdir, "route", "policy", "remove-role", "pcid:moks.context.place.v1", "family-validator"); err != nil {
		t.Fatalf("route policy remove-role: %v", err)
	}
	afterRemove, err := runCLI(t, workdir, "route", "policy", "show")
	if err != nil {
		t.Fatalf("route policy show after remove-role: %v", err)
	}
	if strings.Contains(afterRemove, `"role": "family-validator"`) {
		t.Fatalf("route policy role override still present after remove-role: %s", afterRemove)
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
	if !strings.Contains(output, "seeded_untrusted: false") {
		t.Fatalf("discover output missing unseeded state: %s", output)
	}
	if !strings.Contains(output, "allow_command: moks relay peer allow "+source.LocalPeerID()) {
		t.Fatalf("discover output missing allow command: %s", output)
	}
}

func TestRelayPeerDiscoverSeedCreatesUntrustedPeerEntry(t *testing.T) {
	source := newRuntimeForCLI(t)
	server := httptest.NewServer(relayHandler(context.Background(), source))
	defer server.Close()

	workdir := t.TempDir()
	output, err := runCLI(t, workdir, "relay", "peer", "discover", server.URL+"/relay/peer-card", "seed")
	if err != nil {
		t.Fatalf("discover and seed peer: %v", err)
	}
	if !strings.Contains(output, "seeded_untrusted: true") {
		t.Fatalf("discover output missing seeded state: %s", output)
	}

	runtime, err := kernel.Open(filepath.Join(workdir, ".moks"))
	if err != nil {
		t.Fatalf("open seeded runtime: %v", err)
	}
	defer func() {
		_ = runtime.Close()
	}()
	peer, ok := runtime.LookupPeer(source.LocalPeerID())
	if !ok {
		t.Fatal("expected seeded peer entry")
	}
	if peer.AllowPull || peer.AllowPush {
		t.Fatalf("expected seeded peer to remain untrusted, got pull=%t push=%t", peer.AllowPull, peer.AllowPush)
	}
	if peer.PublicKey != source.LocalPeerPublicKey() {
		t.Fatalf("unexpected seeded public key: %s", peer.PublicKey)
	}
}

func TestRelayPeerPromoteUsesSeededMetadata(t *testing.T) {
	source := newRuntimeForCLI(t)
	server := httptest.NewServer(relayHandler(context.Background(), source))
	defer server.Close()

	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "relay", "peer", "discover", server.URL+"/relay/peer-card", "seed"); err != nil {
		t.Fatalf("discover and seed peer: %v", err)
	}
	output, err := runCLI(t, workdir, "relay", "peer", "promote", source.LocalPeerID(), "both")
	if err != nil {
		t.Fatalf("promote seeded peer: %v", err)
	}
	if !strings.Contains(output, "promoted "+source.LocalPeerID()+" pull=true push=true") {
		t.Fatalf("unexpected promote output: %s", output)
	}

	runtime, err := kernel.Open(filepath.Join(workdir, ".moks"))
	if err != nil {
		t.Fatalf("open promoted runtime: %v", err)
	}
	defer func() {
		_ = runtime.Close()
	}()
	peer, ok := runtime.LookupPeer(source.LocalPeerID())
	if !ok {
		t.Fatal("expected promoted peer entry")
	}
	if !peer.AllowPull || !peer.AllowPush {
		t.Fatalf("expected promoted peer to allow both, got pull=%t push=%t", peer.AllowPull, peer.AllowPush)
	}
	if peer.BatchURL != server.URL+"/relay/batch" || peer.ImportURL != server.URL+"/relay/import" {
		t.Fatalf("expected stored peer metadata to be reused, got batch=%s import=%s", peer.BatchURL, peer.ImportURL)
	}
}

func TestRelayPolicyClaimSetAndList(t *testing.T) {
	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "relay", "policy", "claim", "set", "pcid:test.echo.v1", "*", "1", "any"); err != nil {
		t.Fatalf("set claim policy: %v", err)
	}
	output, err := runCLI(t, workdir, "relay", "policy", "claim", "list")
	if err != nil {
		t.Fatalf("list claim policies: %v", err)
	}
	if !strings.Contains(output, "pcid:test.echo.v1\t*\t1\t0\t0\tany-known-peer\tany-class\tany-federation") {
		t.Fatalf("unexpected claim policy output: %s", output)
	}
}

func TestRelayPeerClassifyAndWeightedPolicy(t *testing.T) {
	source := newRuntimeForCLI(t)
	server := httptest.NewServer(relayHandler(context.Background(), source))
	defer server.Close()

	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "relay", "peer", "discover", server.URL+"/relay/peer-card", "seed"); err != nil {
		t.Fatalf("discover and seed peer: %v", err)
	}
	if _, err := runCLI(t, workdir, "relay", "peer", "classify", source.LocalPeerID(), "auditor", "3"); err != nil {
		t.Fatalf("classify peer: %v", err)
	}
	if _, err := runCLI(t, workdir, "relay", "policy", "claim", "set-weighted", "pcid:test.echo.v1", "*", "1", "3", "any", "auditor"); err != nil {
		t.Fatalf("set weighted claim policy: %v", err)
	}
	output, err := runCLI(t, workdir, "relay", "peer", "list")
	if err != nil {
		t.Fatalf("list peers: %v", err)
	}
	if !strings.Contains(output, "\tclass=auditor\tweight=3\t") {
		t.Fatalf("unexpected peer list output: %s", output)
	}
	policies, err := runCLI(t, workdir, "relay", "policy", "claim", "list")
	if err != nil {
		t.Fatalf("list weighted claim policies: %v", err)
	}
	if !strings.Contains(policies, "pcid:test.echo.v1\t*\t1\t3\t0\tany-known-peer\tauditor\tany-federation") {
		t.Fatalf("unexpected weighted policy output: %s", policies)
	}
}

func TestRelayPeerFederateAndFederatedPolicy(t *testing.T) {
	source := newRuntimeForCLI(t)
	server := httptest.NewServer(relayHandler(context.Background(), source))
	defer server.Close()

	workdir := t.TempDir()
	if _, err := runCLI(t, workdir, "relay", "peer", "discover", server.URL+"/relay/peer-card", "seed"); err != nil {
		t.Fatalf("discover and seed peer: %v", err)
	}
	if _, err := runCLI(t, workdir, "relay", "peer", "classify", source.LocalPeerID(), "auditor", "3"); err != nil {
		t.Fatalf("classify peer: %v", err)
	}
	if _, err := runCLI(t, workdir, "relay", "peer", "federate", source.LocalPeerID(), "fed-a"); err != nil {
		t.Fatalf("federate peer: %v", err)
	}
	if _, err := runCLI(t, workdir, "relay", "policy", "claim", "set-federated", "pcid:test.echo.v1", "*", "1", "3", "1", "any", "auditor", "fed-a"); err != nil {
		t.Fatalf("set federated claim policy: %v", err)
	}
	output, err := runCLI(t, workdir, "relay", "peer", "list")
	if err != nil {
		t.Fatalf("list peers: %v", err)
	}
	if !strings.Contains(output, "\tfederation=fed-a\t") {
		t.Fatalf("unexpected peer federation output: %s", output)
	}
	policies, err := runCLI(t, workdir, "relay", "policy", "claim", "list")
	if err != nil {
		t.Fatalf("list federated claim policies: %v", err)
	}
	if !strings.Contains(policies, "pcid:test.echo.v1\t*\t1\t3\t1\tany-known-peer\tauditor\tfed-a") {
		t.Fatalf("unexpected federated policy output: %s", policies)
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

func publishContextRoutePromises(t *testing.T, workdir string) {
	t.Helper()
	commands := [][]string{
		{"route", "bind", "context-app", "context", "true"},
		{"route", "promise", "receive", "context-app", "pcid:moks.context.place.v1", "true"},
		{"route", "promise", "deliver", "local-router", "context-app", "pcid:moks.context.place.v1", "true"},
	}
	for _, command := range commands {
		if _, err := runCLI(t, workdir, command...); err != nil {
			t.Fatalf("publish route evidence %q: %v", command, err)
		}
	}
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
