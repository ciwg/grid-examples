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

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
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
