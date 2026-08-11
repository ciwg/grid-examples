package service

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func testRecord(t *testing.T) (Record, ed25519.PrivateKey) {
	t.Helper()
	seed := sha256.Sum256([]byte("ex7 deterministic record signer"))
	private := ed25519.NewKeyFromSeed(seed[:])
	public := private.Public().(ed25519.PublicKey)
	digest := sha256.Sum256(public)
	return Record{Protocol: "bafkreifhodcald6kzib36rzeji27hnqjdkeycibnkcigcsz7mzejz6obiy", ID: "rec-1", Signer: "alice", CreatedAt: "2026-08-11T18:00:00Z", Payload: []byte(`{"observation":"guard loose","tool_id":"table-saw"}`), KeyID: "ed25519:" + hex.EncodeToString(digest[:]), PublicKey: public}, private
}

func TestRecordCanonicalSignAndParse(t *testing.T) {
	record, private := testRecord(t)
	_, raw, err := record.Sign(private)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseRecord(raw)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Protocol != record.Protocol || parsed.Signer != "alice" {
		t.Fatalf("parsed = %#v", parsed)
	}
}

func TestRecordRejectsTamperedBytes(t *testing.T) {
	record, private := testRecord(t)
	_, raw, err := record.Sign(private)
	if err != nil {
		t.Fatal(err)
	}
	raw[len(raw)-1] ^= 1
	if _, err := ParseRecord(raw); err == nil {
		t.Fatal("accepted tampered record")
	}
}
