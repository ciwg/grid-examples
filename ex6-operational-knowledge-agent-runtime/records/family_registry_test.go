package records

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

// TestPackageFamilyRegistryMatchesFrozenSpecifications keeps compiled pCIDs
// bound to the exact immutable spec bytes and their published registry values.
// Source: DI-solan.
func TestPackageFamilyRegistryMatchesFrozenSpecifications(t *testing.T) {
	registryPath := filepath.Join("..", "docs", "protocols", "package-family-pcid-registry.md")
	registryBody, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatalf("read human registry: %v", err)
	}
	specificationDirectory := filepath.Join("..", "docs", "protocols", "package-families")
	specifications, err := filepath.Glob(filepath.Join(specificationDirectory, "*.md"))
	if err != nil {
		t.Fatalf("list frozen specifications: %v", err)
	}
	if len(specifications) != len(packageFamilyPCIDs) {
		t.Fatalf("frozen specification count = %d, built-in registry count = %d", len(specifications), len(packageFamilyPCIDs))
	}
	for family, expectedPCID := range packageFamilyPCIDs {
		specificationName := strings.ReplaceAll(family, ".", "-") + ".md"
		specificationPath := filepath.Join(specificationDirectory, specificationName)
		specification, err := os.ReadFile(specificationPath)
		if err != nil {
			t.Fatalf("read frozen specification for %s: %v", family, err)
		}
		digest, err := multihash.Sum(specification, multihash.SHA2_256, -1)
		if err != nil {
			t.Fatalf("hash frozen specification for %s: %v", family, err)
		}
		actualPCID := cid.NewCidV1(cid.Raw, digest).String()
		if actualPCID != expectedPCID {
			t.Fatalf("pCID for %s = %s, want %s", family, actualPCID, expectedPCID)
		}
		expectedRow := "| `" + family + "` | [" + specificationName + "](package-families/" + specificationName + ") | `" + expectedPCID + "` |"
		if !strings.Contains(string(registryBody), expectedRow) {
			t.Fatalf("human registry missing published mapping for %s", family)
		}
	}
}

// TestExternalFamilyRequiresAnExplicitValidPCID keeps the external-package
// boundary explicit before a record reaches the Grid encoder. Source: DI-solan.
func TestExternalFamilyRequiresAnExplicitValidPCID(t *testing.T) {
	envelope := Envelope{
		Family:       "writer.note.v1",
		ProtocolPCID: "bafkreigwh6qript7zma7gu6fgxixmno2eglo3v2bhwpqr3dg5utiyagmca",
		RecordID:     "writer-1",
		Signer:       "writer-agent",
		Timestamp:    "2026-08-11T00:00:00Z",
		Payload:      []byte(`{"note":"hello","writer_id":"writer-1"}`),
	}
	if err := envelope.Validate(); err != nil {
		t.Fatalf("validate explicit external pCID: %v", err)
	}
	envelope.ProtocolPCID = ""
	if err := envelope.Validate(); err == nil {
		t.Fatal("expected missing external pCID rejection")
	}
	envelope.ProtocolPCID = "not-a-cid"
	if err := envelope.Validate(); err == nil {
		t.Fatal("expected malformed external pCID rejection")
	}
}
