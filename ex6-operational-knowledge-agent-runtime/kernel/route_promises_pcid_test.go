package kernel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
)

func TestRoutePromiseProtocolPCIDMatchesSpecification(t *testing.T) {
	cas, err := store.OpenCAS(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join("..", "docs", "protocols", "route-promises.md"))
	if err != nil {
		t.Fatal(err)
	}
	actual, err := cas.PutCID(raw)
	if err != nil {
		t.Fatal(err)
	}
	if actual.String() != RoutePromisesProtocolPCID {
		t.Fatalf("route promise pCID = %s, want %s", actual, RoutePromisesProtocolPCID)
	}
}
