package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCASPutCIDRoundTrip(t *testing.T) {
	cas, err := OpenCAS(t.TempDir())
	if err != nil {
		t.Fatalf("open CAS: %v", err)
	}
	body := []byte("workflow event")
	objectCID, err := cas.PutCID(body)
	if err != nil {
		t.Fatalf("put CID: %v", err)
	}
	if objectCID.Version() != 1 {
		t.Errorf("CID version = %d, want 1", objectCID.Version())
	}
	if objectCID.Type() != 0x55 {
		t.Errorf("CID codec = %d, want raw", objectCID.Type())
	}
	stored, err := cas.GetCID(objectCID)
	if err != nil {
		t.Fatalf("get CID: %v", err)
	}
	if string(stored) != string(body) {
		t.Errorf("stored body = %q, want %q", stored, body)
	}
}

func TestCASGetCIDReadsLegacyObject(t *testing.T) {
	cas, err := OpenCAS(t.TempDir())
	if err != nil {
		t.Fatalf("open CAS: %v", err)
	}
	body := []byte("legacy workflow artifact")
	legacyID, err := cas.Put(body)
	if err != nil {
		t.Fatalf("put legacy object: %v", err)
	}
	if _, err := os.Stat(cas.pathFor(legacyID)); err != nil {
		t.Fatalf("stat legacy object: %v", err)
	}
	objectCID, err := cas.cidFor(body)
	if err != nil {
		t.Fatalf("derive CID: %v", err)
	}
	stored, err := cas.GetCID(objectCID)
	if err != nil {
		t.Fatalf("get legacy through CID: %v", err)
	}
	if string(stored) != string(body) {
		t.Errorf("stored body = %q, want %q", stored, body)
	}
}

func TestCASListCIDsNormalizesLegacyObjects(t *testing.T) {
	cas, err := OpenCAS(t.TempDir())
	if err != nil {
		t.Fatalf("open CAS: %v", err)
	}
	legacyBody := []byte("legacy")
	if _, err := cas.Put(legacyBody); err != nil {
		t.Fatalf("put legacy object: %v", err)
	}
	currentBody := []byte("current")
	currentCID, err := cas.PutCID(currentBody)
	if err != nil {
		t.Fatalf("put CID object: %v", err)
	}
	legacyCID, err := cas.cidFor(legacyBody)
	if err != nil {
		t.Fatalf("derive legacy CID: %v", err)
	}
	objectCIDs, err := cas.ListCIDs()
	if err != nil {
		t.Fatalf("list CIDs: %v", err)
	}
	if len(objectCIDs) != 2 {
		t.Fatalf("CID count = %d, want 2", len(objectCIDs))
	}
	seen := map[string]bool{}
	for _, objectCID := range objectCIDs {
		seen[objectCID.String()] = true
	}
	if !seen[currentCID.String()] || !seen[legacyCID.String()] {
		t.Errorf("listed CIDs = %v, want %s and %s", objectCIDs, currentCID, legacyCID)
	}
}

func TestCASPutCIDRepairsCorruptExistingObject(t *testing.T) {
	cas, err := OpenCAS(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("repairable")
	objectCID, err := cas.cidFor(body)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cas.pathFor(objectCID.String()), []byte("corrupt"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := cas.PutCID(body); err != nil {
		t.Fatal(err)
	}
	stored, err := cas.GetCID(objectCID)
	if err != nil {
		t.Fatal(err)
	}
	if string(stored) != string(body) {
		t.Fatalf("repaired body = %q, want %q", stored, body)
	}
}

func TestCASPutCIDRepairsCorruptCIDWithValidLegacyFallback(t *testing.T) {
	cas, err := OpenCAS(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("legacy repairable")
	if _, err := cas.Put(body); err != nil {
		t.Fatal(err)
	}
	objectCID, err := cas.cidFor(body)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cas.pathFor(objectCID.String()), []byte("corrupt"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := cas.PutCID(body); err != nil {
		t.Fatal(err)
	}
	stored, err := cas.GetCID(objectCID)
	if err != nil {
		t.Fatal(err)
	}
	if string(stored) != string(body) {
		t.Fatalf("repaired fallback body = %q, want %q", stored, body)
	}
}

func TestCASListCIDsIgnoresUnexpectedFiles(t *testing.T) {
	cas, err := OpenCAS(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	objectCID, err := cas.PutCID([]byte("retained"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cas.root, ".cas-interrupted"), []byte("temporary"), 0o644); err != nil {
		t.Fatal(err)
	}
	objectCIDs, err := cas.ListCIDs()
	if err != nil {
		t.Fatal(err)
	}
	if len(objectCIDs) != 1 || objectCIDs[0] != objectCID {
		t.Fatalf("listed CIDs = %v, want %s", objectCIDs, objectCID)
	}
}
