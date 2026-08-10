package kernel

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ipfs/go-cid"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/artifact"
)

func TestNoRegisteredRecipientObservation(t *testing.T) {
	root := t.TempDir()
	store, err := artifact.NewStore("kernel", root)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	server := &Server{subscribers: map[string]map[*clientConn]bool{}, store: store}
	if server.broadcast(nil, cid.Undef, []byte("raw")) {
		t.Fatal("broadcast reported a recipient")
	}
	server.recordObservation("no_registered_recipient", "raw-cid", cid.Undef.String(), "", "no registered recipient")
	raw, err := os.ReadFile(filepath.Join(root, "observations.jsonl"))
	if err != nil {
		t.Fatalf("read observations: %v", err)
	}
	if !strings.Contains(string(raw), `"kind":"no_registered_recipient"`) {
		t.Fatalf("observation = %s", raw)
	}
}
