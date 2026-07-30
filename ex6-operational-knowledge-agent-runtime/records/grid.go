package records

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
)

const (
	GridTagNumber  uint64 = 1735551332
	ProtocolCIDTag uint64 = 42
)

// GridEnvelope is the canonical CBOR carriage for pCID-selected messages.
// Intent: Preserve exact PromiseGrid envelope bytes while keeping every slot
// after the pCID under the selected protocol's control. Source: DI-bavuk
type GridEnvelope struct {
	ProtocolPCID cid.Cid
	Slots        []any
}

func EncodeGrid(envelope GridEnvelope) ([]byte, error) {
	if err := validateProtocolCID(envelope.ProtocolPCID); err != nil {
		return nil, err
	}
	values := make([]any, 0, len(envelope.Slots)+1)
	values = append(values, cbor.Tag{Number: ProtocolCIDTag, Content: envelope.ProtocolPCID.Bytes()})
	values = append(values, envelope.Slots...)
	encoded, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return nil, err
	}
	return encoded.Marshal(cbor.Tag{Number: GridTagNumber, Content: values})
}

func DecodeGrid(raw []byte) (GridEnvelope, error) {
	var wrapped cbor.Tag
	if err := cbor.Unmarshal(raw, &wrapped); err != nil {
		return GridEnvelope{}, err
	}
	if wrapped.Number != GridTagNumber {
		return GridEnvelope{}, fmt.Errorf("grid tag = %d, want %d", wrapped.Number, GridTagNumber)
	}
	values, ok := wrapped.Content.([]any)
	if !ok || len(values) == 0 {
		return GridEnvelope{}, errors.New("grid content must be a non-empty array")
	}
	selector, ok := values[0].(cbor.Tag)
	if !ok || selector.Number != ProtocolCIDTag {
		return GridEnvelope{}, errors.New("grid selector must be tag 42")
	}
	protocolBytes, ok := selector.Content.([]byte)
	if !ok {
		return GridEnvelope{}, errors.New("grid selector must contain CID bytes")
	}
	protocolPCID, err := cid.Cast(protocolBytes)
	if err != nil {
		return GridEnvelope{}, fmt.Errorf("grid selector CID: %w", err)
	}
	if err := validateProtocolCID(protocolPCID); err != nil {
		return GridEnvelope{}, err
	}
	envelope := GridEnvelope{ProtocolPCID: protocolPCID, Slots: values[1:]}
	canonical, err := EncodeGrid(envelope)
	if err != nil {
		return GridEnvelope{}, err
	}
	if !bytes.Equal(raw, canonical) {
		return GridEnvelope{}, errors.New("grid envelope is not canonical CBOR")
	}
	return envelope, nil
}

func validateProtocolCID(protocolPCID cid.Cid) error {
	if !protocolPCID.Defined() {
		return errors.New("protocol pCID is required")
	}
	if protocolPCID.Version() != 1 {
		return errors.New("protocol pCID must be CIDv1")
	}
	if _, err := cid.Cast(protocolPCID.Bytes()); err != nil {
		return fmt.Errorf("protocol pCID: %w", err)
	}
	return nil
}
