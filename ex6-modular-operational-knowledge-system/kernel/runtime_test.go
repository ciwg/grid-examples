package kernel_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/builtin"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
	pkgmeta "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages"
	contextpkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/context"
	inventorypkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/inventory"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/knowledge"
	linkspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/links"
	maintenancepkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/maintenance"
	procedurespkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/procedures"
	receivingpkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/receiving"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/runs"
	trainingpkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/training"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/records"
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

func TestTrainingPackageCommands(t *testing.T) {
	runtime := newRuntime(t)
	if _, err := runtime.RunCommand(context.Background(), []string{"training", "create", "train-1", "ForkliftBasics", "Forklift-safety-basics"}); err != nil {
		t.Fatalf("create training: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"training", "record-session", "train-1", "run-3", "alice", "bob", "passed", "completed-lab"}); err != nil {
		t.Fatalf("record training session: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"training", "certify", "train-1", "comp-1", "alice", "certified", "ready-for-shift"}); err != nil {
		t.Fatalf("certify training: %v", err)
	}
	listing, err := runtime.RunCommand(context.Background(), []string{"training", "list"})
	if err != nil {
		t.Fatalf("list training: %v", err)
	}
	if !strings.Contains(listing, "train-1\tForkliftBasics\tsessions=1\tcompletions=1") {
		t.Fatalf("unexpected training listing: %s", listing)
	}
	inspect, err := runtime.RunCommand(context.Background(), []string{"training", "inspect", "train-1"})
	if err != nil {
		t.Fatalf("inspect training: %v", err)
	}
	if !strings.Contains(inspect, "title: ForkliftBasics") || !strings.Contains(inspect, "run-3:alice->bob") || !strings.Contains(inspect, "alice:certified") {
		t.Fatalf("unexpected training inspect: %s", inspect)
	}
}

func TestMaintenancePackageCommands(t *testing.T) {
	runtime := newRuntime(t)
	if _, err := runtime.RunCommand(context.Background(), []string{"context", "resource", "create", "res-2", "Mixer", "Paint-mixer", "place-1"}); err != nil {
		t.Fatalf("create maintenance resource: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"maintenance", "create", "maint-1", "res-2", "MixerCheck", "Quarterly-mixer-check"}); err != nil {
		t.Fatalf("create maintenance item: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"maintenance", "record-service", "maint-1", "run-4", "res-2", "carol", "completed", "lubed-bearings"}); err != nil {
		t.Fatalf("record maintenance service: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"maintenance", "record-finding", "maint-1", "find-1", "res-2", "serviceable", "ready-for-next-shift"}); err != nil {
		t.Fatalf("record maintenance finding: %v", err)
	}
	listing, err := runtime.RunCommand(context.Background(), []string{"maintenance", "list"})
	if err != nil {
		t.Fatalf("list maintenance: %v", err)
	}
	if !strings.Contains(listing, "maint-1\tres-2\tMixerCheck\tservices=1\tfindings=1") {
		t.Fatalf("unexpected maintenance listing: %s", listing)
	}
	inspect, err := runtime.RunCommand(context.Background(), []string{"maintenance", "inspect", "maint-1"})
	if err != nil {
		t.Fatalf("inspect maintenance: %v", err)
	}
	if !strings.Contains(inspect, "resource_id: res-2") || !strings.Contains(inspect, "run-4:carol@res-2") || !strings.Contains(inspect, "find-1:serviceable") {
		t.Fatalf("unexpected maintenance inspect: %s", inspect)
	}
}

func TestReceivingPackageCommands(t *testing.T) {
	runtime := newRuntime(t)
	if _, err := runtime.RunCommand(context.Background(), []string{"context", "place", "create", "place-2", "Dock", "Inbound-dock"}); err != nil {
		t.Fatalf("create receiving place: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"context", "resource", "create", "res-3", "PalletA", "Pallet-load", "place-2"}); err != nil {
		t.Fatalf("create receiving resource: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"receiving", "create", "recv-1", "place-2", "InboundPallet", "Receive-pallet-a"}); err != nil {
		t.Fatalf("create receiving item: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"receiving", "record-receipt", "recv-1", "run-5", "place-2", "dave", "accepted", "count-matched"}); err != nil {
		t.Fatalf("record receipt: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"receiving", "record-disposition", "recv-1", "disp-1", "stock", "res-3", "moved-to-storage"}); err != nil {
		t.Fatalf("record disposition: %v", err)
	}
	listing, err := runtime.RunCommand(context.Background(), []string{"receiving", "list"})
	if err != nil {
		t.Fatalf("list receiving: %v", err)
	}
	if !strings.Contains(listing, "recv-1\tplace-2\tInboundPallet\treceipts=1\tdispositions=1") {
		t.Fatalf("unexpected receiving listing: %s", listing)
	}
	inspect, err := runtime.RunCommand(context.Background(), []string{"receiving", "inspect", "recv-1"})
	if err != nil {
		t.Fatalf("inspect receiving: %v", err)
	}
	if !strings.Contains(inspect, "place_id: place-2") || !strings.Contains(inspect, "run-5:dave@place-2") || !strings.Contains(inspect, "disp-1:stock@res-3") {
		t.Fatalf("unexpected receiving inspect: %s", inspect)
	}
}

func TestInventoryPackageCommands(t *testing.T) {
	runtime := newRuntime(t)
	if _, err := runtime.RunCommand(context.Background(), []string{"context", "place", "create", "place-3", "Warehouse", "Main-storage"}); err != nil {
		t.Fatalf("create inventory place: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"context", "resource", "create", "res-4", "BoltBin", "Bin-of-bolts", "place-3"}); err != nil {
		t.Fatalf("create inventory resource: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"inventory", "create", "inv-1", "place-3", "BoltAudit", "Cycle-count-bolt-bin"}); err != nil {
		t.Fatalf("create inventory item: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"inventory", "record-count", "inv-1", "run-6", "place-3", "ellen", "42", "matched", "count-confirmed"}); err != nil {
		t.Fatalf("record inventory count: %v", err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"inventory", "record-reconcile", "inv-1", "rec-1", "balanced", "res-4", "no-adjustment-needed"}); err != nil {
		t.Fatalf("record inventory reconciliation: %v", err)
	}
	listing, err := runtime.RunCommand(context.Background(), []string{"inventory", "list"})
	if err != nil {
		t.Fatalf("list inventory: %v", err)
	}
	if !strings.Contains(listing, "inv-1\tplace-3\tBoltAudit\tcounts=1\treconciles=1") {
		t.Fatalf("unexpected inventory listing: %s", listing)
	}
	inspect, err := runtime.RunCommand(context.Background(), []string{"inventory", "inspect", "inv-1"})
	if err != nil {
		t.Fatalf("inspect inventory: %v", err)
	}
	if !strings.Contains(inspect, "place_id: place-3") || !strings.Contains(inspect, "run-6:ellen=42@place-3") || !strings.Contains(inspect, "rec-1:balanced@res-4") {
		t.Fatalf("unexpected inventory inspect: %s", inspect)
	}
}

func TestInstalledPackageManifestSelfCheck(t *testing.T) {
	runtime := newRuntime(t)
	packageDir := helperPackageDir(t, false)
	manifest, err := runtime.InstallPackageDir(context.Background(), packageDir)
	if err != nil {
		t.Fatalf("install package: %v", err)
	}
	if manifest.ID != "helper-agent" {
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

func TestInstalledPackageFamilyRequiresValidatorRoute(t *testing.T) {
	runtime := newRuntime(t)
	packageDir := helperPackageDir(t, false)
	manifestPath := filepath.Join(packageDir, "moks-package.json")
	body, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(body, &manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	manifest["claims"] = []map[string]any{
		{"protocol_pcid": "pcid:helper.echo.v1", "role": "reader", "summary": "Validates helper echo envelopes."},
	}
	updated, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	if err := os.WriteFile(manifestPath, updated, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	scriptPath := filepath.Join(packageDir, "helper-agent.sh")
	scriptBody, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("read script: %v", err)
	}
	updatedScript := strings.ReplaceAll(string(scriptBody), `"role":"family-validator"`, `"role":"reader"`)
	if err := os.WriteFile(scriptPath, []byte(updatedScript), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	if _, err := runtime.InstallPackageDir(context.Background(), packageDir); err == nil || !strings.Contains(err.Error(), "family helper.echo.v1 requires family-validator claim") {
		t.Fatalf("expected family-validator route rejection, got %v", err)
	}
}

func TestUnknownFamilyStoredAndLaterInterpreted(t *testing.T) {
	unknownRaw := []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"peer-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	runtimeA := newRuntime(t)
	if _, err := runtimeA.AppendRecord(context.Background(), unknownRaw); err != nil {
		t.Fatalf("append unknown raw: %v", err)
	}
	exported, err := runtimeA.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	exportedEnvelope, err := records.Parse(exported.Records[0])
	if err != nil {
		t.Fatalf("parse exported envelope: %v", err)
	}
	if exportedEnvelope.Family != "helper.echo.v1" {
		t.Fatalf("unexpected family: %s", exportedEnvelope.Family)
	}
	if string(exportedEnvelope.Payload) != `{"message":"hello"}` {
		t.Fatalf("expected payload preserved, got %s", string(exportedEnvelope.Payload))
	}
	if exportedEnvelope.AuthorSignature == "" {
		t.Fatal("expected authored unknown record to be signed")
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

func TestAppendRecordAddsSemanticAuthorSignature(t *testing.T) {
	runtime := newRuntime(t)
	raw := []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	envelope, err := runtime.AppendRecord(context.Background(), raw)
	if err != nil {
		t.Fatalf("append record: %v", err)
	}
	if envelope.AuthorKeyID != runtime.LocalPeerID() {
		t.Fatalf("unexpected author key id: %s", envelope.AuthorKeyID)
	}
	if envelope.AuthorPublicKey != runtime.LocalPeerPublicKey() {
		t.Fatalf("unexpected author public key: %s", envelope.AuthorPublicKey)
	}
	if envelope.AuthorSignature == "" {
		t.Fatal("expected author signature to be added")
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

func TestBatchMetadataValidation(t *testing.T) {
	runtime := newRuntime(t)
	err := runtime.ImportBatch(context.Background(), grid.Batch{
		Format:         grid.RelayBatchFormat,
		Implementation: "",
		ExportedAt:     "2026-07-28T00:00:00Z",
		Records:        []json.RawMessage{json.RawMessage(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"peer-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)},
	})
	if err == nil {
		t.Fatal("expected missing implementation error")
	}
}

func TestImportBatchIsIdempotentForExactBytes(t *testing.T) {
	raw := json.RawMessage(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"peer-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	runtime := newRuntime(t)
	batch := grid.Batch{
		Format:         grid.RelayBatchFormat,
		Implementation: "peer-a",
		ExportedAt:     "2026-07-28T00:00:00Z",
		ImplementationClaims: []grid.ImplementationClaim{
			{PackageID: "helper-agent", PackageVersion: "0.1.0", ProtocolPCID: "pcid:helper.echo.v1", Role: "family-validator"},
		},
		Records:      []json.RawMessage{raw},
		RecordProofs: grid.ProofsForRecords([]json.RawMessage{raw}),
	}
	if err := runtime.ImportBatch(context.Background(), batch); err != nil {
		t.Fatalf("first import: %v", err)
	}
	if err := runtime.ImportBatch(context.Background(), batch); err != nil {
		t.Fatalf("second import: %v", err)
	}
	if len(runtime.History()) != 1 {
		t.Fatalf("expected exact-byte dedupe after repeated import, got %d", len(runtime.History()))
	}
}

func TestExportBatchIncludesImplementationClaims(t *testing.T) {
	runtime := newRuntime(t)
	batch, err := runtime.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	if batch.Implementation != runtime.LocalPeerID() {
		t.Fatalf("unexpected implementation: %s", batch.Implementation)
	}
	if len(batch.RecordProofs) != len(batch.Records) {
		t.Fatalf("expected record proofs to match records, got %d proofs for %d records", len(batch.RecordProofs), len(batch.Records))
	}
	if len(batch.ClaimProofs) != len(batch.ImplementationClaims) {
		t.Fatalf("expected claim proofs to match claims, got %d proofs for %d claims", len(batch.ClaimProofs), len(batch.ImplementationClaims))
	}
	if len(batch.Routes) == 0 {
		t.Fatal("expected exported routes")
	}
	if len(batch.RecordSignatures) != len(batch.Records) {
		t.Fatalf("expected record signatures to match records, got %d signatures for %d records", len(batch.RecordSignatures), len(batch.Records))
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
	foundContextRoute := false
	for _, route := range batch.Routes {
		if route.ProtocolPCID == contextpkg.PlaceProtocol && route.Role == "family-validator" && route.PackageID == "context" {
			foundContextRoute = true
			break
		}
	}
	if !foundContextRoute {
		t.Fatalf("expected exported context route, got %#v", batch.Routes)
	}
}

func TestProtocolRoutesIncludeBuiltinClaims(t *testing.T) {
	runtime := newRuntime(t)
	routes := runtime.ProtocolRoutes()
	if len(routes) == 0 {
		t.Fatal("expected protocol routes")
	}
	foundContext := false
	foundOps := false
	for _, route := range routes {
		if route.ProtocolPCID == contextpkg.PlaceProtocol && route.Role == "family-validator" && route.PackageID == "context" {
			foundContext = true
			if route.RouteType != "direct" {
				t.Fatalf("unexpected context route type: %s", route.RouteType)
			}
			if !slices.Equal(route.Families, []string{"moks.context.place.v1"}) {
				t.Fatalf("unexpected context route families: %#v", route.Families)
			}
		}
		if route.ProtocolPCID == builtin.OpsFamilyProtocol && route.Role == "family-validator" && route.PackageID == "ops-note" {
			foundOps = true
			if route.RouteType != "direct" {
				t.Fatalf("unexpected ops route type: %s", route.RouteType)
			}
			if !slices.Equal(route.Families, []string{"moks.ops.note.v1"}) {
				t.Fatalf("unexpected ops route families: %#v", route.Families)
			}
		}
	}
	if !foundContext || !foundOps {
		t.Fatalf("expected builtin protocol routes, got %#v", routes)
	}
}

func TestProtocolRoutesCanModelParserHop(t *testing.T) {
	runtimeRoot := filepath.Join(t.TempDir(), ".moks")
	runtime, err := kernel.Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer func() {
		_ = runtime.Close()
	}()
	parserPackage := kernel.BuiltinPackage{
		Manifest: pkgmeta.Manifest{
			ID:      "parser-agent",
			Version: "0.1.0",
			Claims: []pkgmeta.ImplementationClaim{
				{
					ProtocolPCID:   "pcid:moks.raw.v1",
					Role:           "parser",
					RouteType:      "parser",
					EmitsProtocols: []string{"pcid:moks.parsed.v1"},
					Summary:        "Parses raw envelopes into parsed ones.",
				},
				{
					ProtocolPCID: "pcid:moks.parsed.v1",
					Role:         "handler",
					Summary:      "Handles parsed envelopes directly.",
				},
			},
		},
		Commands:   map[string]kernel.BuiltinCommand{},
		Validators: map[string]kernel.BuiltinValidator{},
	}
	if err := runtime.RegisterBuiltin(parserPackage); err != nil {
		t.Fatalf("register parser package: %v", err)
	}
	routes := runtime.ProtocolRoutes()
	foundParser := false
	foundHandler := false
	for _, route := range routes {
		if route.PackageID == "parser-agent" && route.ProtocolPCID == "pcid:moks.raw.v1" {
			foundParser = true
			if route.RouteType != "parser" {
				t.Fatalf("unexpected parser route type: %s", route.RouteType)
			}
			if !slices.Equal(route.EmitsProtocols, []string{"pcid:moks.parsed.v1"}) {
				t.Fatalf("unexpected parser emits protocols: %#v", route.EmitsProtocols)
			}
		}
		if route.PackageID == "parser-agent" && route.ProtocolPCID == "pcid:moks.parsed.v1" {
			foundHandler = true
			if route.RouteType != "direct" {
				t.Fatalf("unexpected handler route type: %s", route.RouteType)
			}
		}
	}
	if !foundParser || !foundHandler {
		t.Fatalf("expected parser hop routes, got %#v", routes)
	}
	filtered := runtime.ProtocolRoutesForProtocol("pcid:moks.raw.v1")
	if len(filtered) != 1 {
		t.Fatalf("expected one raw protocol route, got %#v", filtered)
	}
	if filtered[0].RouteType != "parser" {
		t.Fatalf("expected parser route type, got %#v", filtered)
	}
	plan := runtime.ProtocolRoutePlan("pcid:moks.raw.v1")
	if plan.Preferred == nil {
		t.Fatalf("expected preferred parser plan, got %#v", plan)
	}
	if plan.Preferred.Route.RouteType != "parser" {
		t.Fatalf("expected parser preferred route, got %#v", plan.Preferred)
	}
	if len(plan.Preferred.Next) != 1 || plan.Preferred.Next[0].Preferred == nil {
		t.Fatalf("expected downstream preferred route, got %#v", plan.Preferred)
	}
	if plan.Preferred.Next[0].Preferred.Route.ProtocolPCID != "pcid:moks.parsed.v1" {
		t.Fatalf("unexpected downstream preferred route: %#v", plan.Preferred.Next[0].Preferred)
	}
	if len(plan.Preferred.Explanation.Downstream) != 1 {
		t.Fatalf("expected one downstream explanation, got %#v", plan.Preferred.Explanation)
	}
	if !plan.Preferred.Explanation.Downstream[0].Executable {
		t.Fatalf("expected executable downstream explanation, got %#v", plan.Preferred.Explanation.Downstream[0])
	}
	if plan.Preferred.Explanation.Downstream[0].PreferredRoute == "" {
		t.Fatalf("expected downstream preferred route summary, got %#v", plan.Preferred.Explanation.Downstream[0])
	}
	tracePlan := runtime.ProtocolRoutePlanTrace("pcid:moks.raw.v1")
	if tracePlan.Explanation == nil || len(tracePlan.Explanation.Trace) == 0 {
		t.Fatalf("expected route plan trace, got %#v", tracePlan.Explanation)
	}
	if tracePlan.Explanation.TraceSummary == nil || tracePlan.Explanation.TraceSummary.TotalSteps == 0 {
		t.Fatalf("expected route plan trace summary, got %#v", tracePlan.Explanation)
	}
	if tracePlan.Explanation.TraceSummary.Scope != "root" || tracePlan.Explanation.TraceSummary.ProtocolPCID != "pcid:moks.raw.v1" {
		t.Fatalf("expected root trace summary scope metadata, got %#v", tracePlan.Explanation.TraceSummary)
	}
	if tracePlan.Explanation.TraceSummary.HopPath != "root" {
		t.Fatalf("expected root trace summary hop path, got %#v", tracePlan.Explanation.TraceSummary)
	}
	if tracePlan.Explanation.TraceSummary.HopSummary != "root" {
		t.Fatalf("expected root trace summary hop summary, got %#v", tracePlan.Explanation.TraceSummary)
	}
	if tracePlan.Explanation.TraceSummary.HopDepth != 0 || tracePlan.Explanation.TraceSummary.HopIndex != 0 {
		t.Fatalf("expected root trace summary depth/index, got %#v", tracePlan.Explanation.TraceSummary)
	}
	if len(tracePlan.Explanation.DownstreamTraceSummaries) != 1 {
		t.Fatalf("expected one top-level downstream trace summary, got %#v", tracePlan.Explanation.DownstreamTraceSummaries)
	}
	if tracePlan.Explanation.DownstreamTraceSummaries[0].Scope != "downstream" || tracePlan.Explanation.DownstreamTraceSummaries[0].ProtocolPCID != "pcid:moks.parsed.v1" {
		t.Fatalf("expected top-level downstream trace summary scope metadata, got %#v", tracePlan.Explanation.DownstreamTraceSummaries[0])
	}
	if tracePlan.Explanation.DownstreamTraceSummaries[0].HopPath != "root > parser-agent:parser:parser#1 > pcid:moks.parsed.v1#1" {
		t.Fatalf("expected top-level downstream trace summary hop path, got %#v", tracePlan.Explanation.DownstreamTraceSummaries[0])
	}
	if tracePlan.Explanation.DownstreamTraceSummaries[0].HopSummary != "parser-agent:parser:parser [1] -> pcid:moks.parsed.v1 [1]" {
		t.Fatalf("expected top-level downstream trace summary hop summary, got %#v", tracePlan.Explanation.DownstreamTraceSummaries[0])
	}
	if tracePlan.Explanation.DownstreamTraceSummaries[0].HopDepth != 1 || tracePlan.Explanation.DownstreamTraceSummaries[0].HopIndex != 1 {
		t.Fatalf("expected top-level downstream trace summary depth/index, got %#v", tracePlan.Explanation.DownstreamTraceSummaries[0])
	}
	focusedTrace := runtime.ProtocolRoutePlanTraceFocused("pcid:moks.raw.v1", kernel.RoutePlanTraceFilter{
		Kind:   "downstream",
		Target: "pcid:moks.parsed.v1",
	})
	if focusedTrace.Explanation == nil || len(focusedTrace.Explanation.Trace) == 0 {
		t.Fatalf("expected focused downstream trace, got %#v", focusedTrace.Explanation)
	}
	if focusedTrace.Explanation.TraceSummary == nil || focusedTrace.Explanation.TraceSummary.HiddenSteps < 0 {
		t.Fatalf("expected focused trace summary, got %#v", focusedTrace.Explanation)
	}
	if len(focusedTrace.Explanation.DownstreamTraceSummaries) != 1 || focusedTrace.Explanation.DownstreamTraceSummaries[0].ProtocolPCID != "pcid:moks.parsed.v1" {
		t.Fatalf("expected focused downstream trace summary list, got %#v", focusedTrace.Explanation.DownstreamTraceSummaries)
	}
	if len(tracePlan.Preferred.Explanation.Downstream) != 1 || tracePlan.Preferred.Explanation.Downstream[0].TraceSummary == nil {
		t.Fatalf("expected downstream trace summary in traced plan, got %#v", tracePlan.Preferred.Explanation)
	}
	if tracePlan.Preferred.Explanation.Downstream[0].TraceSummary.Scope != "downstream" || tracePlan.Preferred.Explanation.Downstream[0].TraceSummary.ProtocolPCID != "pcid:moks.parsed.v1" {
		t.Fatalf("expected downstream trace summary scope metadata, got %#v", tracePlan.Preferred.Explanation.Downstream[0].TraceSummary)
	}
	if tracePlan.Preferred.Explanation.Downstream[0].TraceSummary.HopPath != "root > parser-agent:parser:parser#1 > pcid:moks.parsed.v1#1" {
		t.Fatalf("expected downstream trace summary hop path, got %#v", tracePlan.Preferred.Explanation.Downstream[0].TraceSummary)
	}
	if tracePlan.Preferred.Explanation.Downstream[0].TraceSummary.HopSummary != "parser-agent:parser:parser [1] -> pcid:moks.parsed.v1 [1]" {
		t.Fatalf("expected downstream trace summary hop summary, got %#v", tracePlan.Preferred.Explanation.Downstream[0].TraceSummary)
	}
	if tracePlan.Preferred.Explanation.Downstream[0].TraceSummary.HopDepth != 1 || tracePlan.Preferred.Explanation.Downstream[0].TraceSummary.HopIndex != 1 {
		t.Fatalf("expected downstream trace summary depth/index, got %#v", tracePlan.Preferred.Explanation.Downstream[0].TraceSummary)
	}
}

func TestProtocolRoutePlanTraceKeepsRepeatedDownstreamProtocolsDistinct(t *testing.T) {
	runtimeRoot := filepath.Join(t.TempDir(), ".moks")
	runtime, err := kernel.Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer func() {
		_ = runtime.Close()
	}()
	alphaParser := kernel.BuiltinPackage{
		Manifest: pkgmeta.Manifest{
			ID:      "alpha-parser",
			Version: "0.1.0",
			Claims: []pkgmeta.ImplementationClaim{
				{
					ProtocolPCID:   "pcid:moks.multi.v1",
					Role:           "parser",
					RouteType:      "parser",
					EmitsProtocols: []string{"pcid:moks.shared.v1"},
					Summary:        "Parses multi envelopes via alpha path.",
				},
			},
		},
		Commands:   map[string]kernel.BuiltinCommand{},
		Validators: map[string]kernel.BuiltinValidator{},
	}
	betaParser := kernel.BuiltinPackage{
		Manifest: pkgmeta.Manifest{
			ID:      "beta-parser",
			Version: "0.1.0",
			Claims: []pkgmeta.ImplementationClaim{
				{
					ProtocolPCID:   "pcid:moks.multi.v1",
					Role:           "parser",
					RouteType:      "parser",
					EmitsProtocols: []string{"pcid:moks.shared.v1"},
					Summary:        "Parses multi envelopes via beta path.",
				},
				{
					ProtocolPCID: "pcid:moks.shared.v1",
					Role:         "handler",
					Summary:      "Handles shared parsed envelopes.",
				},
			},
		},
		Commands:   map[string]kernel.BuiltinCommand{},
		Validators: map[string]kernel.BuiltinValidator{},
	}
	if err := runtime.RegisterBuiltin(alphaParser); err != nil {
		t.Fatalf("register alpha parser: %v", err)
	}
	if err := runtime.RegisterBuiltin(betaParser); err != nil {
		t.Fatalf("register beta parser: %v", err)
	}
	tracePlan := runtime.ProtocolRoutePlanTrace("pcid:moks.multi.v1")
	if tracePlan.Explanation == nil {
		t.Fatalf("expected trace explanation, got %#v", tracePlan)
	}
	if len(tracePlan.Explanation.DownstreamTraceSummaries) != 2 {
		t.Fatalf("expected two downstream trace summaries, got %#v", tracePlan.Explanation.DownstreamTraceSummaries)
	}
	first := tracePlan.Explanation.DownstreamTraceSummaries[0]
	second := tracePlan.Explanation.DownstreamTraceSummaries[1]
	if first.ProtocolPCID != "pcid:moks.shared.v1" || second.ProtocolPCID != "pcid:moks.shared.v1" {
		t.Fatalf("expected repeated shared downstream protocol, got %#v", tracePlan.Explanation.DownstreamTraceSummaries)
	}
	if first.HopPath == second.HopPath {
		t.Fatalf("expected distinct hop paths for repeated downstream protocol, got %#v", tracePlan.Explanation.DownstreamTraceSummaries)
	}
	if first.HopSummary == second.HopSummary {
		t.Fatalf("expected distinct hop summaries for repeated downstream protocol, got %#v", tracePlan.Explanation.DownstreamTraceSummaries)
	}
	if first.HopDepth != 1 || second.HopDepth != 1 || first.HopIndex != 1 || second.HopIndex != 2 {
		t.Fatalf("expected deterministic downstream hop depth/index metadata, got %#v", tracePlan.Explanation.DownstreamTraceSummaries)
	}
}

func TestProtocolRoutePlanPrefersDirectExecutableRoute(t *testing.T) {
	runtimeRoot := filepath.Join(t.TempDir(), ".moks")
	runtime, err := kernel.Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer func() {
		_ = runtime.Close()
	}()
	hybridPackage := kernel.BuiltinPackage{
		Manifest: pkgmeta.Manifest{
			ID:      "hybrid-agent",
			Version: "0.1.0",
			Claims: []pkgmeta.ImplementationClaim{
				{
					ProtocolPCID: "pcid:moks.hybrid.v1",
					Role:         "handler",
					Summary:      "Handles hybrid envelopes directly.",
				},
				{
					ProtocolPCID:   "pcid:moks.hybrid.v1",
					Role:           "parser",
					RouteType:      "parser",
					EmitsProtocols: []string{"pcid:moks.hybrid.parsed.v1"},
					Summary:        "Parses hybrid envelopes.",
				},
				{
					ProtocolPCID: "pcid:moks.hybrid.parsed.v1",
					Role:         "handler",
					Summary:      "Handles parsed hybrid envelopes.",
				},
			},
		},
		Commands:   map[string]kernel.BuiltinCommand{},
		Validators: map[string]kernel.BuiltinValidator{},
	}
	if err := runtime.RegisterBuiltin(hybridPackage); err != nil {
		t.Fatalf("register hybrid package: %v", err)
	}
	plan := runtime.ProtocolRoutePlan("pcid:moks.hybrid.v1")
	if plan.Preferred == nil {
		t.Fatalf("expected preferred route, got %#v", plan)
	}
	if plan.Preferred.Route.RouteType != "direct" {
		t.Fatalf("expected direct route to win, got %#v", plan.Preferred)
	}
	if plan.Preferred.Route.Role != "handler" {
		t.Fatalf("expected direct handler to win, got %#v", plan.Preferred)
	}
	if len(plan.Candidates) < 2 {
		t.Fatalf("expected direct and parser candidates, got %#v", plan.Candidates)
	}
}

func TestProtocolRoutePlanPolicyCanPreferParserOverDirect(t *testing.T) {
	runtimeRoot := filepath.Join(t.TempDir(), ".moks")
	runtime, err := kernel.Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer func() {
		_ = runtime.Close()
	}()
	hybridPackage := kernel.BuiltinPackage{
		Manifest: pkgmeta.Manifest{
			ID:      "hybrid-agent",
			Version: "0.1.0",
			Claims: []pkgmeta.ImplementationClaim{
				{ProtocolPCID: "pcid:moks.hybrid.v1", Role: "handler", Summary: "Handles hybrid envelopes directly."},
				{ProtocolPCID: "pcid:moks.hybrid.v1", Role: "parser", RouteType: "parser", EmitsProtocols: []string{"pcid:moks.hybrid.parsed.v1"}, Summary: "Parses hybrid envelopes."},
				{ProtocolPCID: "pcid:moks.hybrid.parsed.v1", Role: "handler", Summary: "Handles parsed hybrid envelopes."},
			},
		},
		Commands:   map[string]kernel.BuiltinCommand{},
		Validators: map[string]kernel.BuiltinValidator{},
	}
	if err := runtime.RegisterBuiltin(hybridPackage); err != nil {
		t.Fatalf("register hybrid package: %v", err)
	}
	if err := runtime.SetRoutePlanPolicy(grid.RoutePlanPolicy{
		PreferRouteTypes: []string{"parser"},
		AvoidRouteTypes:  []string{"direct"},
	}); err != nil {
		t.Fatalf("set route plan policy: %v", err)
	}
	plan := runtime.ProtocolRoutePlan("pcid:moks.hybrid.v1")
	if plan.Preferred == nil {
		t.Fatalf("expected preferred route, got %#v", plan)
	}
	if plan.Preferred.Route.RouteType != "parser" {
		t.Fatalf("expected parser route to win under policy, got %#v", plan.Preferred)
	}
}

func TestProtocolRoutePlanPolicyCanOverrideOneProtocolOnly(t *testing.T) {
	runtimeRoot := filepath.Join(t.TempDir(), ".moks")
	runtime, err := kernel.Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer func() {
		_ = runtime.Close()
	}()
	hybridPackage := kernel.BuiltinPackage{
		Manifest: pkgmeta.Manifest{
			ID:      "hybrid-agent",
			Version: "0.1.0",
			Claims: []pkgmeta.ImplementationClaim{
				{ProtocolPCID: "pcid:moks.hybrid.v1", Role: "handler", Summary: "Handles hybrid envelopes directly."},
				{ProtocolPCID: "pcid:moks.hybrid.v1", Role: "parser", RouteType: "parser", EmitsProtocols: []string{"pcid:moks.hybrid.parsed.v1"}, Summary: "Parses hybrid envelopes."},
				{ProtocolPCID: "pcid:moks.hybrid.parsed.v1", Role: "handler", Summary: "Handles parsed hybrid envelopes."},
				{ProtocolPCID: "pcid:moks.other.v1", Role: "handler", Summary: "Handles other envelopes directly."},
				{ProtocolPCID: "pcid:moks.other.v1", Role: "parser", RouteType: "parser", EmitsProtocols: []string{"pcid:moks.other.parsed.v1"}, Summary: "Parses other envelopes."},
				{ProtocolPCID: "pcid:moks.other.parsed.v1", Role: "handler", Summary: "Handles parsed other envelopes."},
			},
		},
		Commands:   map[string]kernel.BuiltinCommand{},
		Validators: map[string]kernel.BuiltinValidator{},
	}
	if err := runtime.RegisterBuiltin(hybridPackage); err != nil {
		t.Fatalf("register hybrid package: %v", err)
	}
	if err := runtime.SetRoutePlanPolicy(grid.RoutePlanPolicy{
		PreferRouteTypes: []string{"direct"},
	}); err != nil {
		t.Fatalf("set global route plan policy: %v", err)
	}
	if err := runtime.SetProtocolRoutePlanPolicy("pcid:moks.hybrid.v1", grid.RoutePlanPolicy{
		PreferRouteTypes: []string{"parser"},
		AvoidRouteTypes:  []string{"direct"},
	}); err != nil {
		t.Fatalf("set protocol route plan policy: %v", err)
	}
	hybridPlan := runtime.ProtocolRoutePlan("pcid:moks.hybrid.v1")
	if hybridPlan.Preferred == nil || hybridPlan.Preferred.Route.RouteType != "parser" {
		t.Fatalf("expected parser route to win for hybrid protocol, got %#v", hybridPlan.Preferred)
	}
	otherPlan := runtime.ProtocolRoutePlan("pcid:moks.other.v1")
	if otherPlan.Preferred == nil || otherPlan.Preferred.Route.RouteType != "direct" {
		t.Fatalf("expected direct route to keep winning for other protocol, got %#v", otherPlan.Preferred)
	}
	effective := runtime.EffectiveRoutePlanPolicy("pcid:moks.hybrid.v1")
	if !slices.Equal(effective.PreferRouteTypes, []string{"parser"}) {
		t.Fatalf("unexpected effective prefer route types: %#v", effective.PreferRouteTypes)
	}
	if !slices.Equal(effective.AvoidRouteTypes, []string{"direct"}) {
		t.Fatalf("unexpected effective avoid route types: %#v", effective.AvoidRouteTypes)
	}
}

func TestProtocolRoutePlanPolicyCanOverrideByRoleWithinOneProtocol(t *testing.T) {
	runtimeRoot := filepath.Join(t.TempDir(), ".moks")
	runtime, err := kernel.Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer func() {
		_ = runtime.Close()
	}()
	rolePackage := kernel.BuiltinPackage{
		Manifest: pkgmeta.Manifest{
			ID:      "role-agent",
			Version: "0.1.0",
			Claims: []pkgmeta.ImplementationClaim{
				{ProtocolPCID: "pcid:moks.roleful.v1", Role: "handler", Summary: "Handles roleful envelopes."},
				{ProtocolPCID: "pcid:moks.roleful.v1", Role: "domain-behavior", Summary: "Applies domain behavior for roleful envelopes."},
			},
		},
		Commands:   map[string]kernel.BuiltinCommand{},
		Validators: map[string]kernel.BuiltinValidator{},
	}
	if err := runtime.RegisterBuiltin(rolePackage); err != nil {
		t.Fatalf("register role package: %v", err)
	}
	plan := runtime.ProtocolRoutePlan("pcid:moks.roleful.v1")
	if plan.Preferred == nil || plan.Preferred.Route.Role != "handler" {
		t.Fatalf("expected handler to win by default, got %#v", plan.Preferred)
	}
	if err := runtime.SetProtocolRoleRoutePlanPolicy("pcid:moks.roleful.v1", "domain-behavior", grid.RoutePlanPolicy{
		PreferRoles: []string{"domain-behavior"},
	}); err != nil {
		t.Fatalf("set role route plan policy for domain-behavior: %v", err)
	}
	if err := runtime.SetProtocolRoleRoutePlanPolicy("pcid:moks.roleful.v1", "handler", grid.RoutePlanPolicy{
		AvoidRoles: []string{"handler"},
	}); err != nil {
		t.Fatalf("set role route plan policy for handler: %v", err)
	}
	plan = runtime.ProtocolRoutePlan("pcid:moks.roleful.v1")
	if plan.Preferred == nil || plan.Preferred.Route.Role != "domain-behavior" {
		t.Fatalf("expected domain-behavior to win under role policy, got %#v", plan.Preferred)
	}
	effective := runtime.EffectiveRoutePlanPolicyForRole("pcid:moks.roleful.v1", "domain-behavior")
	if !slices.Equal(effective.PreferRoles, []string{"domain-behavior"}) {
		t.Fatalf("unexpected effective prefer roles: %#v", effective.PreferRoles)
	}
	if len(plan.Explanation.Winner) == 0 {
		t.Fatalf("expected plan winner explanation, got %#v", plan.Explanation)
	}
	if len(plan.Preferred.Explanation.Notes) == 0 {
		t.Fatalf("expected candidate explanation notes, got %#v", plan.Preferred.Explanation)
	}
	if !plan.Preferred.Explanation.PreferredByPolicy {
		t.Fatalf("expected preferred route to be marked preferred by policy, got %#v", plan.Preferred.Explanation)
	}
	if len(plan.Explanation.Comparisons) == 0 {
		t.Fatalf("expected pairwise comparisons, got %#v", plan.Explanation)
	}
}

func TestImportBatchRejectsRecordProofMismatch(t *testing.T) {
	raw := json.RawMessage(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"peer-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	runtime := newRuntime(t)
	batch := grid.Batch{
		Format:         grid.RelayBatchFormat,
		Implementation: "peer-a",
		ExportedAt:     "2026-07-28T00:00:00Z",
		Records:        []json.RawMessage{raw},
		RecordProofs:   []grid.RecordProof{{Digest: "sha256:deadbeef"}},
	}
	if err := runtime.ImportBatch(context.Background(), batch); err == nil {
		t.Fatal("expected record proof mismatch rejection")
	}
	if len(runtime.History()) != 0 {
		t.Fatalf("expected no history on proof mismatch, got %d", len(runtime.History()))
	}
}

func TestImportBatchRejectsRecordSignatureMismatch(t *testing.T) {
	raw := json.RawMessage(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"peer-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	runtime := newRuntime(t)
	batch := grid.Batch{
		Format:         grid.RelayBatchFormat,
		Implementation: "peer-a",
		ExportedAt:     "2026-07-28T00:00:00Z",
		Records:        []json.RawMessage{raw},
		RecordProofs:   grid.ProofsForRecords([]json.RawMessage{raw}),
		RecordSignatures: []grid.RecordSignature{{
			SignerPeerID: "peer-deadbeef",
			PublicKey:    "deadbeef",
			Signature:    "deadbeef",
		}},
	}
	if err := runtime.ImportBatch(context.Background(), batch); err == nil {
		t.Fatal("expected record signature mismatch rejection")
	}
	if len(runtime.History()) != 0 {
		t.Fatalf("expected no history on signature mismatch, got %d", len(runtime.History()))
	}
}

func TestImportBatchRejectsClaimProofMismatch(t *testing.T) {
	raw := json.RawMessage(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"peer-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	runtime := newRuntime(t)
	batch := grid.Batch{
		Format:         grid.RelayBatchFormat,
		Implementation: "peer-a",
		ExportedAt:     "2026-07-28T00:00:00Z",
		ImplementationClaims: []grid.ImplementationClaim{
			{PackageID: "helper-agent", PackageVersion: "0.1.0", ProtocolPCID: "pcid:helper.echo.v1", Role: "family-validator"},
		},
		ClaimProofs: []grid.ClaimProof{{
			SignerPeerID: "peer-deadbeef",
			PublicKey:    "deadbeef",
			Signature:    "deadbeef",
		}},
		Records:      []json.RawMessage{raw},
		RecordProofs: grid.ProofsForRecords([]json.RawMessage{raw}),
	}
	if err := runtime.ImportBatch(context.Background(), batch); err == nil {
		t.Fatal("expected claim proof mismatch rejection")
	}
	if len(runtime.History()) != 0 {
		t.Fatalf("expected no history on claim proof mismatch, got %d", len(runtime.History()))
	}
}

func TestImportBatchRejectsRouteClaimMismatch(t *testing.T) {
	raw := json.RawMessage(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"peer-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	runtime := newRuntime(t)
	batch := grid.Batch{
		Format:         grid.RelayBatchFormat,
		Implementation: "peer-a",
		ExportedAt:     "2026-07-28T00:00:00Z",
		ImplementationClaims: []grid.ImplementationClaim{
			{PackageID: "helper-agent", PackageVersion: "0.1.0", ProtocolPCID: "pcid:helper.echo.v1", Role: "family-validator"},
		},
		Routes: []grid.RouteRegistration{
			{PackageID: "helper-agent", PackageVersion: "0.1.0", ProtocolPCID: "pcid:helper.echo.v1", Role: "reader", Families: []string{"helper.echo.v1"}},
		},
		Records:      []json.RawMessage{raw},
		RecordProofs: grid.ProofsForRecords([]json.RawMessage{raw}),
	}
	if err := runtime.ImportBatch(context.Background(), batch); err == nil || !strings.Contains(err.Error(), "route registration missing matching claim") {
		t.Fatalf("expected route claim mismatch rejection, got %v", err)
	}
}

func TestImportBatchRejectsRouteParserMetadataMismatch(t *testing.T) {
	raw := json.RawMessage(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"peer-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	runtime := newRuntime(t)
	batch := grid.Batch{
		Format:         grid.RelayBatchFormat,
		Implementation: "peer-a",
		ExportedAt:     "2026-07-28T00:00:00Z",
		ImplementationClaims: []grid.ImplementationClaim{
			{PackageID: "parser-agent", PackageVersion: "0.1.0", ProtocolPCID: "pcid:helper.echo.v1", Role: "parser", RouteType: "parser", EmitsProtocols: []string{"pcid:helper.parsed.v1"}},
		},
		Routes: []grid.RouteRegistration{
			{PackageID: "parser-agent", PackageVersion: "0.1.0", ProtocolPCID: "pcid:helper.echo.v1", Role: "parser", RouteType: "transform", EmitsProtocols: []string{"pcid:helper.parsed.v1"}},
		},
		Records:      []json.RawMessage{raw},
		RecordProofs: grid.ProofsForRecords([]json.RawMessage{raw}),
	}
	if err := runtime.ImportBatch(context.Background(), batch); err == nil || !strings.Contains(err.Error(), "route registration route_type mismatch") {
		t.Fatalf("expected route parser metadata mismatch rejection, got %v", err)
	}
}

func TestAttestBatchClaimsAddsThirdPartyAttestations(t *testing.T) {
	exporter := newRuntime(t)
	attester := newRuntime(t)
	if _, err := exporter.AppendRecord(context.Background(), []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)); err != nil {
		t.Fatalf("append exporter record: %v", err)
	}
	batch, err := exporter.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	attested, err := attester.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims: %v", err)
	}
	if len(attested.ClaimAttestations) != len(attested.ImplementationClaims) {
		t.Fatalf("expected claim attestations to match claims, got %d for %d", len(attested.ClaimAttestations), len(attested.ImplementationClaims))
	}
	if attested.ClaimAttestations[0].SignerPeerID != attester.LocalPeerID() {
		t.Fatalf("unexpected attester peer id: %s", attested.ClaimAttestations[0].SignerPeerID)
	}
}

func TestImportBatchRejectsBadThirdPartyClaimAttestation(t *testing.T) {
	exporter := newRuntime(t)
	if _, err := exporter.AppendRecord(context.Background(), []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)); err != nil {
		t.Fatalf("append exporter record: %v", err)
	}
	batch, err := exporter.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	attester := newRuntime(t)
	batch, err = attester.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims: %v", err)
	}
	batch.ClaimAttestations[0].SignerPeerID = batch.Implementation
	runtime := newRuntime(t)
	if err := runtime.ImportBatch(context.Background(), batch); err == nil {
		t.Fatal("expected third-party claim attestation rejection")
	}
}

func TestImportBatchAcceptsClaimAttestationQuorum(t *testing.T) {
	exporter := newRuntime(t)
	if _, err := exporter.AppendRecord(context.Background(), []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-2","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)); err != nil {
		t.Fatalf("append exporter record: %v", err)
	}
	batch, err := exporter.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	claim := findClaimByProtocol(t, batch, contextpkg.PlaceProtocol)
	attester := newRuntime(t)
	batch, err = attester.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims: %v", err)
	}
	importer := newRuntime(t)
	if err := importer.AllowPeer(grid.AllowedPeer{
		PeerID:    attester.LocalPeerID(),
		BatchURL:  "http://attester.example/relay/batch",
		ImportURL: "http://attester.example/relay/import",
		PublicKey: attester.LocalPeerPublicKey(),
	}); err != nil {
		t.Fatalf("allow attester peer: %v", err)
	}
	if err := importer.SetClaimPolicy(grid.ClaimTrustPolicy{
		ProtocolPCID: claim.ProtocolPCID,
		Role:         claim.Role,
		MinAttesters: 1,
	}); err != nil {
		t.Fatalf("set claim policy: %v", err)
	}
	if err := importer.ImportBatch(context.Background(), batch); err != nil {
		t.Fatalf("import batch with satisfied quorum: %v", err)
	}
}

func TestImportBatchRejectsMissingClaimAttestationQuorum(t *testing.T) {
	exporter := newRuntime(t)
	if _, err := exporter.AppendRecord(context.Background(), []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-3","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)); err != nil {
		t.Fatalf("append exporter record: %v", err)
	}
	batch, err := exporter.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	claim := findClaimByProtocol(t, batch, contextpkg.PlaceProtocol)
	attester := newRuntime(t)
	batch, err = attester.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims: %v", err)
	}
	importer := newRuntime(t)
	if err := importer.AllowPeer(grid.AllowedPeer{
		PeerID:    attester.LocalPeerID(),
		BatchURL:  "http://attester.example/relay/batch",
		ImportURL: "http://attester.example/relay/import",
		PublicKey: attester.LocalPeerPublicKey(),
	}); err != nil {
		t.Fatalf("allow attester peer: %v", err)
	}
	if err := importer.SetClaimPolicy(grid.ClaimTrustPolicy{
		ProtocolPCID: claim.ProtocolPCID,
		Role:         claim.Role,
		MinAttesters: 2,
	}); err != nil {
		t.Fatalf("set claim policy: %v", err)
	}
	if err := importer.ImportBatch(context.Background(), batch); err == nil {
		t.Fatal("expected claim quorum rejection")
	}
}

func TestImportBatchAcceptsWeightedClaimTrust(t *testing.T) {
	exporter := newRuntime(t)
	if _, err := exporter.AppendRecord(context.Background(), []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-4","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)); err != nil {
		t.Fatalf("append exporter record: %v", err)
	}
	batch, err := exporter.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	claim := findClaimByProtocol(t, batch, contextpkg.PlaceProtocol)
	attesterA := newRuntime(t)
	attesterB := newRuntime(t)
	batch, err = attesterA.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims A: %v", err)
	}
	batch, err = attesterB.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims B: %v", err)
	}
	importer := newRuntime(t)
	for _, peer := range []struct {
		id     string
		pub    string
		class  string
		weight int
	}{
		{attesterA.LocalPeerID(), attesterA.LocalPeerPublicKey(), "auditor", 2},
		{attesterB.LocalPeerID(), attesterB.LocalPeerPublicKey(), "supplier", 1},
	} {
		if err := importer.AllowPeer(grid.AllowedPeer{
			PeerID:            peer.id,
			BatchURL:          "http://example/relay/batch",
			ImportURL:         "http://example/relay/import",
			PublicKey:         peer.pub,
			AttesterClass:     peer.class,
			AttestationWeight: peer.weight,
		}); err != nil {
			t.Fatalf("allow attester peer %s: %v", peer.id, err)
		}
	}
	if err := importer.SetClaimPolicy(grid.ClaimTrustPolicy{
		ProtocolPCID:   claim.ProtocolPCID,
		Role:           claim.Role,
		MinAttesters:   1,
		MinTrustWeight: 2,
		AllowedClasses: []string{"auditor"},
	}); err != nil {
		t.Fatalf("set weighted claim policy: %v", err)
	}
	if err := importer.ImportBatch(context.Background(), batch); err != nil {
		t.Fatalf("import batch with weighted trust: %v", err)
	}
}

func TestImportBatchRejectsWrongAttesterClass(t *testing.T) {
	exporter := newRuntime(t)
	if _, err := exporter.AppendRecord(context.Background(), []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-5","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)); err != nil {
		t.Fatalf("append exporter record: %v", err)
	}
	batch, err := exporter.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	claim := findClaimByProtocol(t, batch, contextpkg.PlaceProtocol)
	attester := newRuntime(t)
	batch, err = attester.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims: %v", err)
	}
	importer := newRuntime(t)
	if err := importer.AllowPeer(grid.AllowedPeer{
		PeerID:            attester.LocalPeerID(),
		BatchURL:          "http://example/relay/batch",
		ImportURL:         "http://example/relay/import",
		PublicKey:         attester.LocalPeerPublicKey(),
		AttesterClass:     "supplier",
		AttestationWeight: 3,
	}); err != nil {
		t.Fatalf("allow attester peer: %v", err)
	}
	if err := importer.SetClaimPolicy(grid.ClaimTrustPolicy{
		ProtocolPCID:   claim.ProtocolPCID,
		Role:           claim.Role,
		MinAttesters:   1,
		MinTrustWeight: 2,
		AllowedClasses: []string{"auditor"},
	}); err != nil {
		t.Fatalf("set weighted claim policy: %v", err)
	}
	if err := importer.ImportBatch(context.Background(), batch); err == nil {
		t.Fatal("expected claim trust class rejection")
	}
}

func TestImportBatchAcceptsFederatedClaimTrust(t *testing.T) {
	exporter := newRuntime(t)
	if _, err := exporter.AppendRecord(context.Background(), []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-6","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)); err != nil {
		t.Fatalf("append exporter record: %v", err)
	}
	batch, err := exporter.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	claim := findClaimByProtocol(t, batch, contextpkg.PlaceProtocol)
	attesterA := newRuntime(t)
	attesterB := newRuntime(t)
	batch, err = attesterA.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims A: %v", err)
	}
	batch, err = attesterB.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims B: %v", err)
	}
	importer := newRuntime(t)
	for _, peer := range []struct {
		id         string
		pub        string
		class      string
		weight     int
		federation string
	}{
		{attesterA.LocalPeerID(), attesterA.LocalPeerPublicKey(), "auditor", 2, "fed-a"},
		{attesterB.LocalPeerID(), attesterB.LocalPeerPublicKey(), "auditor", 2, "fed-b"},
	} {
		if err := importer.AllowPeer(grid.AllowedPeer{
			PeerID:            peer.id,
			BatchURL:          "http://example/relay/batch",
			ImportURL:         "http://example/relay/import",
			PublicKey:         peer.pub,
			AttesterClass:     peer.class,
			AttestationWeight: peer.weight,
			Federation:        peer.federation,
		}); err != nil {
			t.Fatalf("allow federated peer %s: %v", peer.id, err)
		}
	}
	if err := importer.SetClaimPolicy(grid.ClaimTrustPolicy{
		ProtocolPCID:       claim.ProtocolPCID,
		Role:               claim.Role,
		MinAttesters:       2,
		MinTrustWeight:     3,
		MinFederations:     2,
		AllowedClasses:     []string{"auditor"},
		AllowedFederations: []string{"fed-a", "fed-b"},
	}); err != nil {
		t.Fatalf("set federated claim policy: %v", err)
	}
	if err := importer.ImportBatch(context.Background(), batch); err != nil {
		t.Fatalf("import batch with federated trust: %v", err)
	}
}

func TestImportBatchRejectsSingleFederationSpread(t *testing.T) {
	exporter := newRuntime(t)
	if _, err := exporter.AppendRecord(context.Background(), []byte(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-7","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)); err != nil {
		t.Fatalf("append exporter record: %v", err)
	}
	batch, err := exporter.ExportBatch()
	if err != nil {
		t.Fatalf("export batch: %v", err)
	}
	claim := findClaimByProtocol(t, batch, contextpkg.PlaceProtocol)
	attesterA := newRuntime(t)
	attesterB := newRuntime(t)
	batch, err = attesterA.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims A: %v", err)
	}
	batch, err = attesterB.AttestBatchClaims(batch)
	if err != nil {
		t.Fatalf("attest batch claims B: %v", err)
	}
	importer := newRuntime(t)
	for _, peer := range []struct {
		id     string
		pub    string
		class  string
		weight int
	}{
		{attesterA.LocalPeerID(), attesterA.LocalPeerPublicKey(), "auditor", 2},
		{attesterB.LocalPeerID(), attesterB.LocalPeerPublicKey(), "auditor", 2},
	} {
		if err := importer.AllowPeer(grid.AllowedPeer{
			PeerID:            peer.id,
			BatchURL:          "http://example/relay/batch",
			ImportURL:         "http://example/relay/import",
			PublicKey:         peer.pub,
			AttesterClass:     peer.class,
			AttestationWeight: peer.weight,
			Federation:        "fed-a",
		}); err != nil {
			t.Fatalf("allow federated peer %s: %v", peer.id, err)
		}
	}
	if err := importer.SetClaimPolicy(grid.ClaimTrustPolicy{
		ProtocolPCID:       claim.ProtocolPCID,
		Role:               claim.Role,
		MinAttesters:       2,
		MinTrustWeight:     3,
		MinFederations:     2,
		AllowedClasses:     []string{"auditor"},
		AllowedFederations: []string{"fed-a", "fed-b"},
	}); err != nil {
		t.Fatalf("set federated claim policy: %v", err)
	}
	if err := importer.ImportBatch(context.Background(), batch); err == nil {
		t.Fatal("expected federation spread rejection")
	}
}

func TestImportBatchAcceptsLegacyUnsignedAuthorRecord(t *testing.T) {
	raw := json.RawMessage(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"}}`)
	runtime := newRuntime(t)
	batch := grid.Batch{
		Format:         grid.RelayBatchFormat,
		Implementation: "peer-a",
		ExportedAt:     "2026-07-28T00:00:00Z",
		Records:        []json.RawMessage{raw},
		RecordProofs:   grid.ProofsForRecords([]json.RawMessage{raw}),
	}
	if err := runtime.ImportBatch(context.Background(), batch); err != nil {
		t.Fatalf("import legacy unsigned author record: %v", err)
	}
	if len(runtime.History()) != 1 {
		t.Fatalf("expected imported history, got %d", len(runtime.History()))
	}
}

func TestImportBatchRejectsBadSemanticAuthorSignature(t *testing.T) {
	raw := json.RawMessage(`{"family":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1","record_id":"u-1","signer":"author-a","timestamp":"2026-07-28T00:00:00Z","payload":{"message":"hello"},"author_key_id":"peer-deadbeef","author_public_key":"deadbeef","author_signature":"deadbeef"}`)
	runtime := newRuntime(t)
	batch := grid.Batch{
		Format:         grid.RelayBatchFormat,
		Implementation: "peer-a",
		ExportedAt:     "2026-07-28T00:00:00Z",
		Records:        []json.RawMessage{raw},
		RecordProofs:   grid.ProofsForRecords([]json.RawMessage{raw}),
	}
	if err := runtime.ImportBatch(context.Background(), batch); err == nil {
		t.Fatal("expected semantic author signature rejection")
	}
	if len(runtime.History()) != 0 {
		t.Fatalf("expected no history on bad semantic author signature, got %d", len(runtime.History()))
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
	if err := runtime.RegisterBuiltin(inventorypkg.Package()); err != nil {
		t.Fatalf("register inventory package: %v", err)
	}
	if err := runtime.RegisterBuiltin(runspkg.Package()); err != nil {
		t.Fatalf("register runs package: %v", err)
	}
	if err := runtime.RegisterBuiltin(linkspkg.Package()); err != nil {
		t.Fatalf("register links package: %v", err)
	}
	if err := runtime.RegisterBuiltin(maintenancepkg.Package()); err != nil {
		t.Fatalf("register maintenance package: %v", err)
	}
	if err := runtime.RegisterBuiltin(receivingpkg.Package()); err != nil {
		t.Fatalf("register receiving package: %v", err)
	}
	if err := runtime.RegisterBuiltin(procedurespkg.Package()); err != nil {
		t.Fatalf("register procedures package: %v", err)
	}
	if err := runtime.RegisterBuiltin(trainingpkg.Package()); err != nil {
		t.Fatalf("register training package: %v", err)
	}
	if err := runtime.RegisterBuiltin(builtin.OpsPackage()); err != nil {
		t.Fatalf("register builtin: %v", err)
	}
	return runtime
}

func findClaimByProtocol(t *testing.T, batch grid.Batch, protocolPCID string) grid.ImplementationClaim {
	t.Helper()
	for _, claim := range batch.ImplementationClaims {
		if claim.ProtocolPCID == protocolPCID {
			return claim
		}
	}
	t.Fatalf("missing implementation claim for %s", protocolPCID)
	return grid.ImplementationClaim{}
}

func helperPackageDir(t *testing.T, mismatch bool) string {
	t.Helper()
	dir := t.TempDir()
	executable := filepath.Join(dir, "helper-agent.sh")
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
		"id":          "helper-agent",
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
	executable := filepath.Join(dir, "writer-agent.sh")
	script := `#!/bin/sh
set -eu
case "$1" in
  describe)
    cat <<'EOF'
{"id":"writer-agent","version":"0.1.0","description":"Test writer package","commands":[{"path":["writer","create"],"summary":"Create a writer record"}],"families":[{"name":"writer.note.v1","protocol_pcid":"pcid:writer.note.v1"}],"claims":[{"protocol_pcid":"pcid:writer.note.v1","role":"family-validator","summary":"Validates writer note envelopes."}]}
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
{"output":"created $3","cas":[{"alias":"body1","body":"payload for $3"}],"records":[{"family":"writer.note.v1","protocol_pcid":"pcid:writer.note.v1","record_id":"$3","signer":"writer-agent","timestamp":"2026-07-28T00:00:00Z","payload":{"title":"Writer","body_ref":"\$cas:body1"}}]}
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
		"id":          "writer-agent",
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
