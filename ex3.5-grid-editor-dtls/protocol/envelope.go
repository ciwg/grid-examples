package protocol

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
)

type Proof struct {
	Algorithm string `cbor:"alg"`
	KeyID     string `cbor:"kid"`
	PublicKey []byte `cbor:"pub"`
	Signature []byte `cbor:"sig"`
}

type Envelope struct {
	PCID         cid.Cid
	PayloadBytes []byte
	Proof        Proof
}

func NewEnvelope(pcid cid.Cid, payloadBytes []byte, proof Proof) Envelope {
	return Envelope{
		PCID:         pcid,
		PayloadBytes: append([]byte(nil), payloadBytes...),
		Proof:        proof,
	}
}

// Intent: Preserve slot 1 as the payload's raw CBOR item so the signed wire
// shape remains `grid([42(pCID), payload, proof])` rather than a byte-string
// wrapper. Source: DI-tofug
func (envelope Envelope) SignableBytes() ([]byte, error) {
	signable := []any{
		rawPCIDTag(envelope.PCID),
		cbor.RawMessage(envelope.PayloadBytes),
	}
	return Marshal(signable)
}

func (envelope Envelope) Bytes() ([]byte, error) {
	outer := cbor.RawTag{
		Number: gridTag,
		Content: MustMarshal([]any{
			rawPCIDTag(envelope.PCID),
			cbor.RawMessage(envelope.PayloadBytes),
			envelope.Proof,
		}),
	}
	return Marshal(outer)
}

func ParseEnvelope(envelopeBytes []byte) (Envelope, error) {
	var outer cbor.RawTag
	if err := Unmarshal(envelopeBytes, &outer); err != nil {
		return Envelope{}, fmt.Errorf("decode outer envelope: %w", err)
	}
	if outer.Number != gridTag {
		return Envelope{}, fmt.Errorf("unexpected outer tag %d", outer.Number)
	}
	var slots []cbor.RawMessage
	if err := Unmarshal(outer.Content, &slots); err != nil {
		return Envelope{}, fmt.Errorf("decode envelope slots: %w", err)
	}
	if len(slots) != 3 {
		return Envelope{}, fmt.Errorf("unexpected envelope slot count %d", len(slots))
	}
	var pcidTag cbor.RawTag
	if err := Unmarshal(slots[0], &pcidTag); err != nil {
		return Envelope{}, fmt.Errorf("decode pCID tag: %w", err)
	}
	if pcidTag.Number != cidTag {
		return Envelope{}, fmt.Errorf("unexpected pCID tag %d", pcidTag.Number)
	}
	var pcidBytes []byte
	if err := Unmarshal(pcidTag.Content, &pcidBytes); err != nil {
		return Envelope{}, fmt.Errorf("decode pCID bytes: %w", err)
	}
	pcid, err := cid.Cast(pcidBytes)
	if err != nil {
		return Envelope{}, fmt.Errorf("cast pCID bytes: %w", err)
	}
	var proof Proof
	if err := Unmarshal(slots[2], &proof); err != nil {
		return Envelope{}, fmt.Errorf("decode proof: %w", err)
	}
	return NewEnvelope(pcid, slots[1], proof), nil
}

func rawPCIDTag(pcid cid.Cid) cbor.RawTag {
	return cbor.RawTag{
		Number:  cidTag,
		Content: cbor.RawMessage(MustMarshal(pcid.Bytes())),
	}
}
