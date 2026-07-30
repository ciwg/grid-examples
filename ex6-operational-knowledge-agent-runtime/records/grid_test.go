package records

import (
	"testing"

	"github.com/ipfs/go-cid"
)

const testProtocolCID = "bafkreihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku"

func TestGridRoundTrip(t *testing.T) {
	protocolPCID, err := cid.Decode(testProtocolCID)
	if err != nil {
		t.Fatalf("decode protocol CID: %v", err)
	}
	raw, err := EncodeGrid(GridEnvelope{
		ProtocolPCID: protocolPCID,
		Slots:        []any{uint64(7), "payload"},
	})
	if err != nil {
		t.Fatalf("encode grid: %v", err)
	}
	decoded, err := DecodeGrid(raw)
	if err != nil {
		t.Fatalf("decode grid: %v", err)
	}
	if decoded.ProtocolPCID != protocolPCID {
		t.Errorf("protocol CID = %s, want %s", decoded.ProtocolPCID, protocolPCID)
	}
	if len(decoded.Slots) != 2 {
		t.Fatalf("slot count = %d, want 2", len(decoded.Slots))
	}
}

func TestDecodeGridRejectsNonCanonicalEncoding(t *testing.T) {
	protocolPCID, err := cid.Decode(testProtocolCID)
	if err != nil {
		t.Fatalf("decode protocol CID: %v", err)
	}
	nonCanonical := append([]byte{0xda, 0x67, 0x72, 0x69, 0x64, 0x82, 0xd8, 0x2a, 0x58, 0x24}, protocolPCID.Bytes()...)
	nonCanonical = append(nonCanonical, 0x18, 0x07)
	if _, err := DecodeGrid(nonCanonical); err == nil {
		t.Fatal("decode noncanonical grid succeeded")
	}
}

func TestEncodeGridRejectsCIDv0(t *testing.T) {
	legacyCID, err := cid.Decode("QmYwAPJzv5CZsnAzt8auVTLF3HhWQH1d9q5JZd4z2Qb9N1")
	if err != nil {
		t.Fatalf("decode legacy CID: %v", err)
	}
	if _, err := EncodeGrid(GridEnvelope{ProtocolPCID: legacyCID}); err == nil {
		t.Fatal("encode CIDv0 protocol selector succeeded")
	}
}
