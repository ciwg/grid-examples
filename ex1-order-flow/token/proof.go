package token

import (
	"crypto/rand"
	"fmt"

	cose "github.com/veraison/go-cose"
)

func SignProof(role string, signable []byte) ([]byte, error) {
	signer, err := cose.NewSigner(cose.AlgorithmEdDSA, PrivateKey(role))
	if err != nil {
		return nil, fmt.Errorf("new proof signer: %w", err)
	}
	headers := cose.Headers{
		Protected: cose.ProtectedHeader{
			cose.HeaderLabelAlgorithm: cose.AlgorithmEdDSA,
			cose.HeaderLabelKeyID:     []byte(role),
		},
	}
	proofBytes, err := cose.Sign1(rand.Reader, signer, headers, signable, nil)
	if err != nil {
		return nil, fmt.Errorf("sign proof: %w", err)
	}
	return proofBytes, nil
}

func VerifyProof(expectedRole string, signable []byte, proofBytes []byte) error {
	var message cose.Sign1Message
	if err := message.UnmarshalCBOR(proofBytes); err != nil {
		return fmt.Errorf("decode proof: %w", err)
	}
	verifier, err := cose.NewVerifier(cose.AlgorithmEdDSA, PublicKey(expectedRole))
	if err != nil {
		return fmt.Errorf("new verifier: %w", err)
	}
	if err := message.Verify(nil, verifier); err != nil {
		return fmt.Errorf("verify proof: %w", err)
	}
	if string(message.Payload) != string(signable) {
		return fmt.Errorf("proof payload does not match signable bytes")
	}
	return nil
}
