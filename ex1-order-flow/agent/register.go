package agent

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/token"
)

func (client *Client) RegisterProfiles(ctx context.Context, pcids []cid.Cid) error {
	register := protocol.KernelRegisterMessage{
		Role:         client.role,
		ReceivePCIDs: make([][]byte, 0, len(pcids)),
	}
	for _, pcid := range pcids {
		register.ReceivePCIDs = append(register.ReceivePCIDs, pcid.Bytes())
	}
	payloadBytes, err := protocol.Marshal(register)
	if err != nil {
		return fmt.Errorf("marshal register payload: %w", err)
	}
	envelope := protocol.NewEnvelope(protocol.KernelRegisterProfile.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return err
	}
	proofBytes, err := token.SignProof(client.role, signable)
	if err != nil {
		return err
	}
	envelope = protocol.NewEnvelope(protocol.KernelRegisterProfile.CID, payloadBytes, proofBytes)
	return client.Register(ctx, envelope)
}
