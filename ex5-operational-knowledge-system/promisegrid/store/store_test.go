package store

import (
	"os"
	"path/filepath"
	"testing"
)

type testEnvelopeRecord struct {
	Sequence       uint64
	EnvelopeCID    string
	EnvelopeBase64 string
}

func TestCASStoreRoundTrip(t *testing.T) {
	cas := NewCASStore(filepath.Join(t.TempDir(), "cas", "objects"))
	cid, err := cas.WriteObject([]byte("hello"))
	if err != nil {
		t.Fatalf("write object: %v", err)
	}
	body, err := cas.LoadObject(cid)
	if err != nil {
		t.Fatalf("load object: %v", err)
	}
	if string(body) != "hello" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestAuthoritativeEnvelopeBackfillsFromManifest(t *testing.T) {
	cas := NewCASStore(filepath.Join(t.TempDir(), "cas", "objects"))
	cid, err := cas.WriteObject([]byte("envelope"))
	if err != nil {
		t.Fatalf("seed object: %v", err)
	}
	casPath := cas.ObjectPath(cid)
	if err := os.Remove(casPath); err != nil {
		t.Fatalf("remove cas object: %v", err)
	}
	base64Envelope := "ZW52ZWxvcGU="
	authoritative, err := cas.AuthoritativeEnvelopeBase64(cid, base64Envelope)
	if err != nil {
		t.Fatalf("authoritative envelope: %v", err)
	}
	if authoritative != base64Envelope {
		t.Fatalf("unexpected authoritative envelope %q", authoritative)
	}
	reloaded, err := cas.LoadObject(cid)
	if err != nil {
		t.Fatalf("reload backfilled object: %v", err)
	}
	if string(reloaded) != "envelope" {
		t.Fatalf("unexpected backfilled body %q", reloaded)
	}
}

func TestHydrateAuthoritativeEnvelopes(t *testing.T) {
	cas := NewCASStore(filepath.Join(t.TempDir(), "cas", "objects"))
	cid, err := cas.WriteObject([]byte("authoritative"))
	if err != nil {
		t.Fatalf("seed object: %v", err)
	}
	records, err := HydrateAuthoritativeEnvelopes(
		cas,
		[]testEnvelopeRecord{{Sequence: 7, EnvelopeCID: cid, EnvelopeBase64: "dGFtcGVyZWQ="}},
		func(record testEnvelopeRecord) string { return record.EnvelopeCID },
		func(record testEnvelopeRecord) string { return record.EnvelopeBase64 },
		func(record testEnvelopeRecord) string { return "test envelope 7" },
		func(record *testEnvelopeRecord, base64Envelope string) { record.EnvelopeBase64 = base64Envelope },
	)
	if err != nil {
		t.Fatalf("hydrate envelopes: %v", err)
	}
	if got := records[0].EnvelopeBase64; got != "YXV0aG9yaXRhdGl2ZQ==" {
		t.Fatalf("unexpected hydrated envelope %q", got)
	}
}
