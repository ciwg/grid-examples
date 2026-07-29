package kernel

import "testing"

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
	if workflows[0].ArtifactCID != artifactCID {
		t.Errorf("artifact CID = %q, want %q", workflows[0].ArtifactCID, artifactCID)
	}
	if workflows[0].State != WorkflowDeactivated {
		t.Errorf("workflow state = %q, want %q", workflows[0].State, WorkflowDeactivated)
	}
}
