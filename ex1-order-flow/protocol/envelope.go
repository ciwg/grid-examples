package protocol

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
)

type Envelope struct {
	PCID         cid.Cid
	PayloadBytes []byte
	ProofBytes   []byte
}

func NewEnvelope(pcid cid.Cid, payloadBytes []byte, proofBytes []byte) Envelope {
	return Envelope{
		PCID:         pcid,
		PayloadBytes: append([]byte(nil), payloadBytes...),
		ProofBytes:   append([]byte(nil), proofBytes...),
	}
}

// Intent: Preserve slot 1 as the payload's raw CBOR item so the signed wire
// shape remains `grid([42(pCID), payload, proof])` rather than wrapping the
// payload bytes inside a CBOR byte string. Source: DI-vurad
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
			envelope.ProofBytes,
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
		return Envelope{}, fmt.Errorf("decode outer slots: %w", err)
	}
	if len(slots) != 3 {
		return Envelope{}, fmt.Errorf("unexpected envelope slot count")
	}
	var pcidTag cbor.RawTag
	if err := Unmarshal(slots[0], &pcidTag); err != nil {
		return Envelope{}, fmt.Errorf("decode pcid tag: %w", err)
	}
	if pcidTag.Number != cidTag {
		return Envelope{}, fmt.Errorf("unexpected pcid tag %d", pcidTag.Number)
	}
	var pcidBytes []byte
	if err := Unmarshal(pcidTag.Content, &pcidBytes); err != nil {
		return Envelope{}, fmt.Errorf("decode pcid bytes: %w", err)
	}
	pcid, err := cid.Cast(pcidBytes)
	if err != nil {
		return Envelope{}, fmt.Errorf("cast pcid: %w", err)
	}
	payloadBytes := append([]byte(nil), slots[1]...)
	var proofBytes []byte
	if err := Unmarshal(slots[2], &proofBytes); err != nil {
		return Envelope{}, fmt.Errorf("decode proof bytes: %w", err)
	}
	return NewEnvelope(pcid, payloadBytes, proofBytes), nil
}

func rawPCIDTag(pcid cid.Cid) cbor.RawTag {
	return cbor.RawTag{
		Number:  cidTag,
		Content: cbor.RawMessage(MustMarshal(pcid.Bytes())),
	}
}
