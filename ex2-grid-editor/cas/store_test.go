package cas_test

import (
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/cas"
)

func TestStoreRoundTrip(t *testing.T) {
	t.Parallel()
	storeValue, err := cas.Open(filepath.Join(t.TempDir(), "cas"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	address, err := storeValue.Put([]byte("signed envelope bytes"))
	if err != nil {
		t.Fatalf("put bytes: %v", err)
	}
	got, err := storeValue.Get(address)
	if err != nil {
		t.Fatalf("get bytes: %v", err)
	}
	if string(got) != "signed envelope bytes" {
		t.Fatalf("content mismatch: got %q", got)
	}
}
