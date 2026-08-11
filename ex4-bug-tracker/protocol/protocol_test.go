package protocol_test

import (
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/identity"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocol"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocols"
)

func TestEnvelopeRoundTripAndProof(t *testing.T) {
	key, err := identity.NewAgentKey()
	if err != nil {
		t.Fatalf("new key: %v", err)
	}
	pcid, err := protocol.CIDForBytes(protocols.MustRead(protocols.IssueReportSpec))
	if err != nil {
		t.Fatalf("profile CID: %v", err)
	}
	payload, err := protocol.Marshal(protocol.IssueReport{AgentID: string(key.AgentID()), IssuedAt: "2026-08-10T00:00:00Z", Team: "CORE", Title: "example", Description: "body", Severity: "Low"})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	envelope := protocol.NewEnvelope(pcid, payload, protocol.Proof{Algorithm: "Ed25519", AgentID: string(key.AgentID()), PublicKey: key.PublicKey()})
	signable, err := envelope.SignableBytes()
	if err != nil {
		t.Fatalf("signable bytes: %v", err)
	}
	envelope.Proof.Signature = key.Sign(signable)
	wire, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("wire bytes: %v", err)
	}
	decoded, err := protocol.ParseEnvelope(wire)
	if err != nil {
		t.Fatalf("parse envelope: %v", err)
	}
	decodedSignable, err := decoded.SignableBytes()
	if err != nil {
		t.Fatalf("decoded signable bytes: %v", err)
	}
	if decoded.PCID.String() != pcid.String() {
		t.Fatalf("pCID mismatch: got %s want %s", decoded.PCID, pcid)
	}
	if !identity.Verify(decoded.Proof.PublicKey, decodedSignable, decoded.Proof.Signature) {
		t.Fatal("proof did not verify")
	}
}

func TestProfilesHaveDistinctPCIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, name := range []string{protocols.IssueReportSpec, protocols.IssueLifecycleUpdateSpec, protocols.IssueAttachmentReferenceSpec} {
		pcid, err := protocol.CIDForBytes(protocols.MustRead(name))
		if err != nil {
			t.Fatalf("%s pCID: %v", name, err)
		}
		if seen[pcid.String()] {
			t.Fatalf("duplicate pCID for %s", name)
		}
		seen[pcid.String()] = true
	}
}
