package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/identity"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocol"
)

type mutationCapabilityClaims struct {
	Issuer      string    `cbor:"issuer"`
	Audience    string    `cbor:"audience"`
	DocumentID  string    `cbor:"document_id"`
	ProtocolCID string    `cbor:"protocol_cid"`
	Action      string    `cbor:"action"`
	ExpiresAt   time.Time `cbor:"expires_at"`
	IssuedAt    time.Time `cbor:"issued_at"`
	TokenID     string    `cbor:"token_id"`
}

type mutationCapabilityToken struct {
	Claims mutationCapabilityClaims `cbor:"claims"`
	Proof  protocol.Proof           `cbor:"proof"`
}

func issueMutationCapability(signer *identity.Identity, audience string, documentID string, protocolCID string, action string, ttl time.Duration) (string, mutationCapabilityClaims, error) {
	now := time.Now().UTC()
	tokenIDBytes := make([]byte, 8)
	if _, err := rand.Read(tokenIDBytes); err != nil {
		return "", mutationCapabilityClaims{}, fmt.Errorf("read token id: %w", err)
	}
	claims := mutationCapabilityClaims{
		Issuer:      signer.KeyID(),
		Audience:    audience,
		DocumentID:  documentID,
		ProtocolCID: protocolCID,
		Action:      action,
		IssuedAt:    now,
		ExpiresAt:   now.Add(ttl),
		TokenID:     hex.EncodeToString(tokenIDBytes),
	}
	signable, err := marshalCapabilityClaims(claims)
	if err != nil {
		return "", mutationCapabilityClaims{}, err
	}
	proof, err := signer.SignProof(signable)
	if err != nil {
		return "", mutationCapabilityClaims{}, fmt.Errorf("sign capability: %w", err)
	}
	tokenBytes, err := marshalCapabilityToken(mutationCapabilityToken{
		Claims: claims,
		Proof:  proof,
	})
	if err != nil {
		return "", mutationCapabilityClaims{}, err
	}
	return base64.StdEncoding.EncodeToString(tokenBytes), claims, nil
}

func verifyMutationCapability(tokenBase64 string, expectedAudience string, expectedDocumentID string, expectedProtocolCID string, expectedAction string, now time.Time) (mutationCapabilityClaims, error) {
	tokenBytes, err := base64.StdEncoding.DecodeString(tokenBase64)
	if err != nil {
		return mutationCapabilityClaims{}, fmt.Errorf("decode capability: %w", err)
	}
	token, err := unmarshalCapabilityToken(tokenBytes)
	if err != nil {
		return mutationCapabilityClaims{}, err
	}
	signable, err := marshalCapabilityClaims(token.Claims)
	if err != nil {
		return mutationCapabilityClaims{}, err
	}
	if err := identity.VerifyProof(signable, token.Proof); err != nil {
		return mutationCapabilityClaims{}, fmt.Errorf("verify capability proof: %w", err)
	}
	switch {
	case token.Claims.Issuer == "":
		return mutationCapabilityClaims{}, fmt.Errorf("capability issuer missing")
	case token.Claims.Issuer != token.Proof.KeyID:
		return mutationCapabilityClaims{}, fmt.Errorf("capability issuer mismatch")
	case expectedAudience != "" && token.Claims.Audience != expectedAudience:
		return mutationCapabilityClaims{}, fmt.Errorf("capability audience %q", token.Claims.Audience)
	case token.Claims.DocumentID != expectedDocumentID:
		return mutationCapabilityClaims{}, fmt.Errorf("capability document %q", token.Claims.DocumentID)
	case token.Claims.ProtocolCID != expectedProtocolCID:
		return mutationCapabilityClaims{}, fmt.Errorf("capability protocol %q", token.Claims.ProtocolCID)
	case token.Claims.Action != expectedAction:
		return mutationCapabilityClaims{}, fmt.Errorf("capability action %q", token.Claims.Action)
	case now.UTC().After(token.Claims.ExpiresAt.UTC()):
		return mutationCapabilityClaims{}, fmt.Errorf("capability expired at %s", token.Claims.ExpiresAt.UTC())
	}
	return token.Claims, nil
}

func marshalCapabilityClaims(claims mutationCapabilityClaims) ([]byte, error) {
	mode, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return nil, fmt.Errorf("capability cbor mode: %w", err)
	}
	return mode.Marshal(claims)
}

func marshalCapabilityToken(token mutationCapabilityToken) ([]byte, error) {
	mode, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return nil, fmt.Errorf("capability token cbor mode: %w", err)
	}
	return mode.Marshal(token)
}

func unmarshalCapabilityToken(raw []byte) (mutationCapabilityToken, error) {
	var token mutationCapabilityToken
	if err := cbor.Unmarshal(raw, &token); err != nil {
		return mutationCapabilityToken{}, fmt.Errorf("decode capability token: %w", err)
	}
	return token, nil
}
