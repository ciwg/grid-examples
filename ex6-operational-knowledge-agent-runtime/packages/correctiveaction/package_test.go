package correctiveaction

import (
	"context"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	contextpkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/context"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/knowledge"
	quarantinepkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/quarantine"
	receivingpkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/receiving"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/runs"
)

func TestOpenActionRequiresRejectedQuarantineCase(t *testing.T) {
	runtime, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := runtime.Close(); err != nil {
			t.Error(err)
		}
	})
	for _, p := range []kernel.BuiltinPackage{contextpkg.Package(), knowledgepkg.Package(), runspkg.Package(), receivingpkg.Package(), quarantinepkg.Package(), Package()} {
		if err := runtime.RegisterBuiltin(p); err != nil {
			t.Fatal(err)
		}
	}
	for _, command := range [][]string{{"context", "place", "create", "dock", "Dock", "Inbound"}, {"receiving", "create", "receipt", "dock", "Receipt", "Inbound"}, {"receiving", "record-receipt", "receipt", "run", "dock", "Alice", "failed", "mismatch"}, {"quarantine", "open", "case", "receipt", "run", "Alice", "inspection", "mismatch"}, {"quarantine", "reject", "case", "reject-event", "Bob", "review", "reject"}, {"correctiveaction", "open", "action", "case", "Carol", "review", "Correct supplier label"}} {
		if _, err := runtime.RunCommand(context.Background(), command); err != nil {
			t.Fatalf("%q: %v", command, err)
		}
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"quarantine", "open", "open-case", "receipt", "run", "Alice", "inspection-2", "mismatch"}); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"correctiveaction", "open", "bad", "open-case", "Carol", "review", "must reject first"}); err == nil {
		t.Fatal("open quarantine case accepted for corrective action")
	}
}
