package protocol_test

import (
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/identity"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocol"
)

func TestEnvelopeRoundTrip(t *testing.T) {
	t.Parallel()
	identityValue, err := identity.LoadOrCreate(filepath.Join(t.TempDir(), "seed"))
	if err != nil {
		t.Fatalf("load identity: %v", err)
	}
	pcid, err := protocol.CIDForBytes([]byte("spec bytes"))
	if err != nil {
		t.Fatalf("pcid: %v", err)
	}
	payload := protocol.MustMarshal(map[string]any{"kind": "replace", "content": "hello"})
	envelope := protocol.NewEnvelope(pcid, payload, protocol.Proof{})
	signable, err := envelope.SignableBytes()
	if err != nil {
		t.Fatalf("signable: %v", err)
	}
	proof, err := identityValue.SignProof(signable)
	if err != nil {
		t.Fatalf("sign proof: %v", err)
	}
	envelope.Proof = proof
	wire, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("wire bytes: %v", err)
	}
	decoded, err := protocol.ParseEnvelope(wire)
	if err != nil {
		t.Fatalf("parse envelope: %v", err)
	}
	signable, err = decoded.SignableBytes()
	if err != nil {
		t.Fatalf("decoded signable: %v", err)
	}
	if err := identity.VerifyProof(signable, decoded.Proof); err != nil {
		t.Fatalf("verify proof: %v", err)
	}
	if decoded.PCID.String() != pcid.String() {
		t.Fatalf("pcid mismatch: got %s want %s", decoded.PCID, pcid)
	}
}
