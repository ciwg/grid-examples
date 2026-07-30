package kernel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ipfs/go-cid"
)

func TestWorkflowLifecycleReplaysAfterRestart(t *testing.T) {
	runtimeRoot := t.TempDir()
	runtime, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open initial runtime: %v", err)
	}

	artifactCID, err := runtime.PutCAS([]byte("workflow artifact"))
	if err != nil {
		t.Fatalf("store workflow artifact: %v", err)
	}
	if err := runtime.ImportWorkflow(Workflow{ID: "inventory-receipt", ArtifactCID: artifactCID}); err != nil {
		t.Fatalf("import workflow: %v", err)
	}
	canonicalArtifactCID := runtime.Workflows()[0].ArtifactCID
	if err := runtime.ActivateWorkflow("inventory-receipt"); err != nil {
		t.Fatalf("activate workflow: %v", err)
	}
	if err := runtime.DeactivateWorkflow("inventory-receipt"); err != nil {
		t.Fatalf("deactivate workflow: %v", err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatalf("close initial runtime: %v", err)
	}

	reopened, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("reopen runtime: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := reopened.Close(); closeErr != nil {
			t.Errorf("close reopened runtime: %v", closeErr)
		}
	})

	workflows := reopened.Workflows()
	if len(workflows) != 1 {
		t.Fatalf("workflow count = %d, want 1", len(workflows))
	}
	if workflows[0].ID != "inventory-receipt" {
		t.Errorf("workflow ID = %q, want inventory-receipt", workflows[0].ID)
	}
	if workflows[0].ArtifactCID != canonicalArtifactCID {
		t.Errorf("artifact CID = %q, want %q", workflows[0].ArtifactCID, canonicalArtifactCID)
	}
	if workflows[0].State != WorkflowDeactivated {
		t.Errorf("workflow state = %q, want %q", workflows[0].State, WorkflowDeactivated)
	}
}

func TestWorkflowLifecycleRebuildsDeletedCache(t *testing.T) {
	runtimeRoot := t.TempDir()
	runtime, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	artifactCID, err := runtime.PutCAS([]byte("workflow artifact"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	if err := runtime.ImportWorkflow(Workflow{ID: "inventory-receipt", ArtifactCID: artifactCID}); err != nil {
		t.Fatalf("import workflow: %v", err)
	}
	if err := runtime.ActivateWorkflow("inventory-receipt"); err != nil {
		t.Fatalf("activate workflow: %v", err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatalf("close runtime: %v", err)
	}
	cachePath := filepath.Join(runtimeRoot, "state", "workflow-lifecycle-cache.json")
	if err := os.Remove(cachePath); err != nil {
		t.Fatalf("remove cache: %v", err)
	}
	reopened, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("reopen runtime: %v", err)
	}
	defer func() {
		if err := reopened.Close(); err != nil {
			t.Errorf("close reopened runtime: %v", err)
		}
	}()
	if workflows := reopened.Workflows(); len(workflows) != 1 || workflows[0].State != WorkflowActive {
		t.Fatalf("rebuilt workflows = %#v, want active workflow", workflows)
	}
	if _, err := os.Stat(cachePath); err != nil {
		t.Fatalf("recreated cache: %v", err)
	}
}

func TestWorkflowLifecycleReplayRejectsMissingArtifact(t *testing.T) {
	runtimeRoot := t.TempDir()
	runtime, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	artifactCID, err := cid.Decode(testWorkflowArtifactCID)
	if err != nil {
		t.Fatalf("decode missing artifact CID: %v", err)
	}
	raw, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{
		State: WorkflowImported, WorkflowAlias: "missing-artifact", ArtifactCID: artifactCID,
	})
	if err != nil {
		t.Fatalf("encode lifecycle event: %v", err)
	}
	if _, err := runtime.cas.PutCID(raw); err != nil {
		t.Fatalf("store lifecycle event: %v", err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatalf("close runtime: %v", err)
	}
	reopened, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("reopen runtime: %v", err)
	}
	defer func() {
		if err := reopened.Close(); err != nil {
			t.Errorf("close reopened runtime: %v", err)
		}
	}()
	if workflows := reopened.Workflows(); len(workflows) != 0 {
		t.Fatalf("workflows = %#v, want no missing-artifact import", workflows)
	}
}

func TestWorkflowLifecycleCacheRecordsHeadEventCID(t *testing.T) {
	runtimeRoot := t.TempDir()
	runtime, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	artifactCID, err := runtime.PutCAS([]byte("workflow artifact"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	if err := runtime.ImportWorkflow(Workflow{ID: "inventory-receipt", ArtifactCID: artifactCID}); err != nil {
		t.Fatalf("import workflow: %v", err)
	}
	cache, err := os.ReadFile(filepath.Join(runtimeRoot, "state", "workflow-lifecycle-cache.json"))
	if err != nil {
		t.Fatalf("read cache: %v", err)
	}
	var entries []workflowLifecycleCacheEntry
	if err := json.Unmarshal(cache, &entries); err != nil {
		t.Fatalf("decode cache: %v", err)
	}
	if len(entries) != 1 || entries[0].EventCID == "" {
		t.Fatalf("cache entries = %#v, want one entry with event CID", entries)
	}
}

func TestWorkflowLifecycleReplaySkipsCorruptCASObject(t *testing.T) {
	runtimeRoot := t.TempDir()
	runtime, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	artifactCID, err := runtime.PutCAS([]byte("workflow artifact"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	if err := runtime.ImportWorkflow(Workflow{ID: "inventory-receipt", ArtifactCID: artifactCID}); err != nil {
		t.Fatalf("import workflow: %v", err)
	}
	corruptCID, err := runtime.cas.PutCID([]byte("will be corrupted"))
	if err != nil {
		t.Fatalf("store corruptible object: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runtimeRoot, "cas", corruptCID.String()), []byte("corrupt"), 0o644); err != nil {
		t.Fatalf("corrupt CAS object: %v", err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatalf("close runtime: %v", err)
	}
	reopened, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("reopen with corrupt object: %v", err)
	}
	defer func() {
		if err := reopened.Close(); err != nil {
			t.Errorf("close reopened runtime: %v", err)
		}
	}()
	if workflows := reopened.Workflows(); len(workflows) != 1 || workflows[0].ID != "inventory-receipt" {
		t.Fatalf("workflows = %#v, want retained valid workflow", workflows)
	}
}

func TestWorkflowLifecycleKeepsDurableEventWhenCacheWriteFails(t *testing.T) {
	runtimeRoot := t.TempDir()
	runtime, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	artifactCID, err := runtime.PutCAS([]byte("workflow artifact"))
	if err != nil {
		t.Fatalf("store artifact: %v", err)
	}
	runtime.workflows.cachePath = t.TempDir()
	if err := runtime.ImportWorkflow(Workflow{ID: "inventory-receipt", ArtifactCID: artifactCID}); err != nil {
		t.Fatalf("import with failed cache write: %v", err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatalf("close runtime: %v", err)
	}
	reopened, err := Open(runtimeRoot)
	if err != nil {
		t.Fatalf("reopen after failed cache write: %v", err)
	}
	defer func() {
		if err := reopened.Close(); err != nil {
			t.Errorf("close reopened runtime: %v", err)
		}
	}()
	if workflows := reopened.Workflows(); len(workflows) != 1 || workflows[0].State != WorkflowImported {
		t.Fatalf("workflows = %#v, want durable imported workflow", workflows)
	}
}
