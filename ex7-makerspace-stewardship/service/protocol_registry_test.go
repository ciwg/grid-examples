package service

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

// TestFrozenMakerspaceFamilyPCIDs keeps the registry tied to the immutable
// specification bytes instead of treating pCID strings as hand-maintained
// labels. Source: DI-tohak.
func TestFrozenMakerspaceFamilyPCIDs(t *testing.T) {
	t.Parallel()

	families := []struct {
		name string
		path string
		pcid string
	}{
		{
			name: "equipment observation",
			path: filepath.Join("..", "docs", "protocols", "makerspace-families", "makerspace-equipment-observation-v1.md"),
			pcid: "bafkreifhodcald6kzib36rzeji27hnqjdkeycibnkcigcsz7mzejz6obiy",
		},
		{
			name: "safety disposition",
			path: filepath.Join("..", "docs", "protocols", "makerspace-families", "makerspace-safety-disposition-v1.md"),
			pcid: "bafkreigt3p2l4uel7wmjr4kple7o55ymchlhh43gajjwsgaeifoogeztc4",
		},
		{
			name: "off-site loan",
			path: filepath.Join("..", "docs", "protocols", "makerspace-families", "makerspace-offsite-loan-v1.md"),
			pcid: "bafkreibgbky6hbswdimkngjii5zglgvfkolxjxuonbmnqwcdjt4v2qyliq",
		},
		{
			name: "off-site return",
			path: filepath.Join("..", "docs", "protocols", "makerspace-families", "makerspace-offsite-return-v1.md"),
			pcid: "bafkreifjgcfwkbwzrgmtj2wfhk3kpzmjbj3wjxid52jau5rrpthabr3ehi",
		},
	}

	for _, family := range families {
		family := family
		t.Run(family.name, func(t *testing.T) {
			t.Parallel()
			bytes, err := os.ReadFile(family.path)
			if err != nil {
				t.Fatalf("read frozen specification: %v", err)
			}
			got := rawCIDv1(bytes)
			if got != family.pcid {
				t.Fatalf("pCID = %s, want %s", got, family.pcid)
			}
		})
	}
}

// TestMakerspaceRecordProfileDocumentCID preserves the exact document identity
// that a later Ex7 implementation claim must cite. Source: DI-tohak.
func TestMakerspaceRecordProfileDocumentCID(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "docs", "protocols", "makerspace-record-v1.md")
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read frozen record profile: %v", err)
	}
	const want = "bafkreid4ebb6ywvwvumetn6pddhyh2pw5uvbvrt6j7wdwv7v7eovb5wdce"
	if got := rawCIDv1(bytes); got != want {
		t.Fatalf("record-profile document CID = %s, want %s", got, want)
	}
}

// TestFrozenParticipantAgentPCIDs binds every finished-product protocol family
// and its human-readable registry entry to immutable specification bytes.
// Source: DI-sisad.
func TestFrozenParticipantAgentPCIDs(t *testing.T) {
	t.Parallel()

	families := []struct {
		name     string
		specPath string
		registry string
		pcid     string
	}{
		{"root history", filepath.Join("..", "docs", "protocols", "participant", "participant-root-history-v1.md"), filepath.Join("..", "docs", "protocols", "participant-pcid-registry.md"), "bafkreia7cn4srmmkxbwxk2hoezedjvuyokhypcsddjd4evx56lhtmsq3nm"},
		{"device authorization", filepath.Join("..", "docs", "protocols", "participant", "participant-device-authorization-v1.md"), filepath.Join("..", "docs", "protocols", "participant-pcid-registry.md"), "bafkreifmbhgjwfmwbemkf4ogsg3gvuavjhttkitzf3muie3dhv5tdn4hq4"},
		{"key revocation", filepath.Join("..", "docs", "protocols", "participant", "participant-key-revocation-v1.md"), filepath.Join("..", "docs", "protocols", "participant-pcid-registry.md"), "bafkreify46v4jp3dvz3szem6vrukscplr3kshku7fhrzn4mc26scnnvqoi"},
		{"threshold recovery", filepath.Join("..", "docs", "protocols", "participant", "participant-threshold-recovery-v1.md"), filepath.Join("..", "docs", "protocols", "participant-pcid-registry.md"), "bafkreicjdzlriq3nfasza5nmnflpycche63wn2n5kauq66jivlwhhomesy"},
		{"terminal approval", filepath.Join("..", "docs", "protocols", "participant", "participant-terminal-approval-v1.md"), filepath.Join("..", "docs", "protocols", "participant-pcid-registry.md"), "bafkreidztcisyvexrlia4eos7wko27e4rqt7ivbmikjbkr5tzbbts3rcd4"},
		{"peer card", filepath.Join("..", "docs", "protocols", "peer", "participant-peer-card-v1.md"), filepath.Join("..", "docs", "protocols", "peer-pcid-registry.md"), "bafkreicstci6idwm6d6dbt52ppqyjcapibskz27qmnfuyntg6zck72fa24"},
		{"exact record carriage", filepath.Join("..", "docs", "protocols", "carriage", "exact-record-carriage-v1.md"), filepath.Join("..", "docs", "protocols", "carriage-pcid-registry.md"), "bafkreihrlojt47erjc6uawkm47s7evppp23tk3ljlkl347ten4v3kb624i"},
	}

	for _, family := range families {
		family := family
		t.Run(family.name, func(t *testing.T) {
			t.Parallel()
			bytes, err := os.ReadFile(family.specPath)
			if err != nil {
				t.Fatalf("read frozen specification: %v", err)
			}
			if got := rawCIDv1(bytes); got != family.pcid {
				t.Fatalf("pCID = %s, want %s", got, family.pcid)
			}
			registry, err := os.ReadFile(family.registry)
			if err != nil {
				t.Fatalf("read pCID registry: %v", err)
			}
			if !strings.Contains(string(registry), family.pcid) {
				t.Fatalf("registry %s does not declare %s", family.registry, family.pcid)
			}
		})
	}
}

func rawCIDv1(bytes []byte) string {
	digest := sha256.Sum256(bytes)
	multihash, err := multihash.Encode(digest[:], multihash.SHA2_256)
	if err != nil {
		panic(err)
	}
	return cid.NewCidV1(cid.Raw, multihash).String()
}
