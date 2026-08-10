package protocol

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProfileInventoryMatchesPublishedDocuments(t *testing.T) {
	// Intent: Derive pCIDs from source bytes so published inventory drift fails
	// at the reader-facing claim boundary. Source: DI-dilav.
	architecture, err := os.ReadFile(filepath.Join("..", "docs", "architecture.md"))
	if err != nil {
		t.Fatalf("read architecture inventory: %v", err)
	}
	scope, err := os.ReadFile(filepath.Join("..", "CHANGELOG.md"))
	if err != nil {
		t.Fatalf("read implementation scope: %v", err)
	}
	profiles := []string{
		"live-document.md",
		"live-awareness.md",
		"document-metadata.md",
		"publish-document.md",
	}
	for _, profile := range profiles {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			spec, err := os.ReadFile(filepath.Join("..", "protocols", profile))
			if err != nil {
				t.Fatalf("read source spec: %v", err)
			}
			pcid, err := CIDForBytes(spec)
			if err != nil {
				t.Fatalf("derive pCID: %v", err)
			}
			if !strings.Contains(string(architecture), pcid.String()) {
				t.Fatalf("architecture inventory does not publish %s", pcid)
			}
			if !strings.Contains(string(scope), pcid.String()) {
				t.Fatalf("implementation scope does not publish %s", pcid)
			}
		})
	}
}
