package protocol

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProfileInventoryMatchesPublishedDocuments(t *testing.T) {
	design, err := os.ReadFile(filepath.Join("..", "docs", "design.md"))
	if err != nil {
		t.Fatalf("read design inventory: %v", err)
	}
	scope, err := os.ReadFile(filepath.Join("..", "CHANGELOG.md"))
	if err != nil {
		t.Fatalf("read implementation scope: %v", err)
	}
	profiles := []struct {
		profile  Profile
		specName string
	}{
		{OrderProfile, "order.md"},
		{PickPackProfile, "pick-pack.md"},
		{AccountingProfile, "accounting.md"},
		{ShipmentProfile, "shipment.md"},
		{KernelRegisterProfile, "kernel-register.md"},
	}
	for _, test := range profiles {
		test := test
		t.Run(test.profile.Name, func(t *testing.T) {
			spec, err := os.ReadFile(filepath.Join("..", "specdocs", test.specName))
			if err != nil {
				t.Fatalf("read %s: %v", test.specName, err)
			}
			derived, err := CIDForBytes(spec)
			if err != nil {
				t.Fatalf("derive pCID: %v", err)
			}
			if derived.String() != test.profile.CID.String() {
				t.Fatalf("profile pCID = %s, derived from %s = %s", test.profile.CID, test.specName, derived)
			}
			if !strings.Contains(string(design), derived.String()) {
				t.Fatalf("design inventory does not publish %s", derived)
			}
			if !strings.Contains(string(scope), derived.String()) {
				t.Fatalf("implementation scope does not publish %s", derived)
			}
		})
	}
}
