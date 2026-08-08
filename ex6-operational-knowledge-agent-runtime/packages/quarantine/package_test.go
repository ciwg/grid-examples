package quarantinepkg

import (
	"context"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	contextpkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/context"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/knowledge"
	receivingpkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/receiving"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/runs"
)

func TestQuarantineLifecycleRetainsExplicitEvents(t *testing.T) {
	runtime, err := kernel.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := runtime.Close(); err != nil {
			t.Error(err)
		}
	})
	for _, builtin := range []kernel.BuiltinPackage{contextpkg.Package(), knowledgepkg.Package(), runspkg.Package(), receivingpkg.Package(), Package()} {
		if err := runtime.RegisterBuiltin(builtin); err != nil {
			t.Fatal(err)
		}
	}
	for _, command := range [][]string{
		{"context", "place", "create", "dock", "Dock", "Inbound"},
		{"receiving", "create", "receipt", "dock", "Inbound", "Pallet"},
		{"receiving", "record-receipt", "receipt", "receipt-run", "dock", "Alice", "failed", "seal-mismatch"},
		{"quarantine", "open", "case-1", "receipt", "receipt-run", "Alice", "inspection-1", "seal-mismatch", "hold"},
		{"quarantine", "release", "case-1", "release-1", "Bob", "review-1", "cleared"},
	} {
		if _, err := runtime.RunCommand(context.Background(), command); err != nil {
			t.Fatalf("run %q: %v", command, err)
		}
	}
	listing, err := runtime.RunCommand(context.Background(), []string{"quarantine", "list"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(listing, "case-1\treceipt\trelease\tevents=2") {
		t.Fatalf("listing = %s", listing)
	}
	if _, err := runtime.RunCommand(context.Background(), []string{"quarantine", "reject", "case-1", "reject-1", "Carol", "review-2"}); err == nil {
		t.Fatal("terminal quarantine case accepted a second resolution")
	}
}
