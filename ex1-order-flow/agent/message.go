package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/token"
)

func BuildSignedEnvelope(role string, pcid protocol.Profile, payloadBytes []byte) (protocol.Envelope, error) {
	envelope := protocol.NewEnvelope(pcid.CID, payloadBytes, nil)
	signable, err := envelope.SignableBytes()
	if err != nil {
		return protocol.Envelope{}, err
	}
	proofBytes, err := token.SignProof(role, signable)
	if err != nil {
		return protocol.Envelope{}, err
	}
	return protocol.NewEnvelope(pcid.CID, payloadBytes, proofBytes), nil
}

func VerifySignedEnvelope(expectedRole string, envelope protocol.Envelope) error {
	signable, err := envelope.SignableBytes()
	if err != nil {
		return err
	}
	return token.VerifyProof(expectedRole, signable, envelope.ProofBytes)
}

func IssueMessageCapability(issuer string, audience string, profile protocol.Profile, kind string) ([]byte, error) {
	tokenBytes, _, err := token.IssueCapability(issuer, audience, profile.CID.String(), kind, "send", 5*time.Minute)
	if err != nil {
		return nil, err
	}
	return tokenBytes, nil
}

func VerifyMessageCapability(tokenBytes []byte, expectedIssuer string, audience string, profile protocol.Profile, kind string) error {
	_, err := token.VerifyCapability(tokenBytes, expectedIssuer, audience, profile.CID.String(), kind, "send", time.Now().UTC())
	return err
}

func ReceiveTyped[T any](ctx context.Context, client *Client, expectedRole string, expectedPCID protocol.Profile, out *T) (protocol.Envelope, string, error) {
	envelope, exactCID, err := client.Receive(ctx)
	if err != nil {
		return protocol.Envelope{}, "", err
	}
	if envelope.PCID.String() != expectedPCID.CID.String() {
		return protocol.Envelope{}, "", fmt.Errorf("unexpected pCID %s", envelope.PCID.String())
	}
	if err := VerifySignedEnvelope(expectedRole, envelope); err != nil {
		return protocol.Envelope{}, "", err
	}
	if err := protocol.Unmarshal(envelope.PayloadBytes, out); err != nil {
		return protocol.Envelope{}, "", err
	}
	return envelope, exactCID, nil
}
