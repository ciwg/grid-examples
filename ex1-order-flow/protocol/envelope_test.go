package protocol

import (
	"testing"

	"github.com/fxamacker/cbor/v2"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/token"
)

func TestEnvelopeEncodesPayloadAsRawCBORItem(t *testing.T) {
	payloadBytes := MustMarshal(map[string]any{
		"kind":               "submit",
		"customer_order_ref": "demo-001",
	})
	envelope := NewEnvelope(OrderProfile.CID, payloadBytes, []byte("proof"))
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("Bytes(): %v", err)
	}

	var outer cbor.RawTag
	if err := Unmarshal(envelopeBytes, &outer); err != nil {
		t.Fatalf("Unmarshal outer: %v", err)
	}
	if outer.Number != gridTag {
		t.Fatalf("outer.Number = %d, want %d", outer.Number, gridTag)
	}

	var slots []cbor.RawMessage
	if err := Unmarshal(outer.Content, &slots); err != nil {
		t.Fatalf("Unmarshal slots: %v", err)
	}
	if len(slots) != 3 {
		t.Fatalf("len(slots) = %d, want 3", len(slots))
	}
	if len(slots[1]) == 0 {
		t.Fatal("slot 1 is empty")
	}
	if slots[1][0]>>5 != 5 {
		t.Fatalf("slot 1 major type = %d, want map item not bstr", slots[1][0]>>5)
	}
	if got := string(slots[1]); got != string(payloadBytes) {
		t.Fatalf("slot 1 bytes changed during encoding")
	}
}

func TestParseEnvelopePreservesPayloadBytesAndProofVerification(t *testing.T) {
	payloadBytes := MustMarshal(map[string]any{
		"kind":               "request",
		"customer_order_ref": "demo-002",
		"amount_cents":       uint64(1999),
	})
	envelope := NewEnvelope(AccountingProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		t.Fatalf("SignableBytes(): %v", err)
	}
	proofBytes, err := token.SignProof("seller", signable)
	if err != nil {
		t.Fatalf("SignProof(): %v", err)
	}
	envelope = NewEnvelope(AccountingProfile.CID, payloadBytes, proofBytes)

	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("Bytes(): %v", err)
	}
	parsed, err := ParseEnvelope(envelopeBytes)
	if err != nil {
		t.Fatalf("ParseEnvelope(): %v", err)
	}
	if parsed.PCID.String() != AccountingProfile.CID.String() {
		t.Fatalf("parsed.PCID = %s, want %s", parsed.PCID, AccountingProfile.CID)
	}
	if got := string(parsed.PayloadBytes); got != string(payloadBytes) {
		t.Fatalf("payload bytes changed during round trip")
	}
	if err := token.VerifyProof("seller", mustSignableBytes(t, parsed), parsed.ProofBytes); err != nil {
		t.Fatalf("VerifyProof(): %v", err)
	}
}

func mustSignableBytes(t *testing.T, envelope Envelope) []byte {
	t.Helper()
	signable, err := envelope.SignableBytes()
	if err != nil {
		t.Fatalf("SignableBytes(): %v", err)
	}
	return signable
}
