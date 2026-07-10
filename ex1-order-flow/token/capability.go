package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	cose "github.com/veraison/go-cose"
)

const (
	claimIssuer     = int64(1)
	claimSubject    = int64(2)
	claimAudience   = int64(3)
	claimExpiration = int64(4)
	claimIssuedAt   = int64(6)
	claimTokenID    = int64(7)

	privateProtocol = int64(-70001)
	privateKind     = int64(-70002)
	privateAction   = int64(-70003)
)

type CapabilityClaims struct {
	Issuer       string
	Subject      string
	Audience     string
	ProtocolCID  string
	Kind         string
	Action       string
	ExpiresAt    time.Time
	IssuedAt     time.Time
	TokenID      string
}

func IssueCapability(role string, audience string, protocolCID string, kind string, action string, ttl time.Duration) ([]byte, CapabilityClaims, error) {
	now := time.Now().UTC()
	tokenIDBytes := make([]byte, 8)
	if _, err := rand.Read(tokenIDBytes); err != nil {
		return nil, CapabilityClaims{}, fmt.Errorf("rand token id: %w", err)
	}
	claims := CapabilityClaims{
		Issuer:      role,
		Subject:     "capability_token",
		Audience:    audience,
		ProtocolCID: protocolCID,
		Kind:        kind,
		Action:      action,
		IssuedAt:    now,
		ExpiresAt:   now.Add(ttl),
		TokenID:     hex.EncodeToString(tokenIDBytes),
	}
	payload, err := marshalClaims(claims)
	if err != nil {
		return nil, CapabilityClaims{}, err
	}
	signer, err := cose.NewSigner(cose.AlgorithmEdDSA, PrivateKey(role))
	if err != nil {
		return nil, CapabilityClaims{}, fmt.Errorf("new token signer: %w", err)
	}
	headers := cose.Headers{
		Protected: cose.ProtectedHeader{
			cose.HeaderLabelAlgorithm: cose.AlgorithmEdDSA,
			cose.HeaderLabelKeyID:     []byte(role),
		},
	}
	tokenBytes, err := cose.Sign1(rand.Reader, signer, headers, payload, nil)
	if err != nil {
		return nil, CapabilityClaims{}, fmt.Errorf("sign capability token: %w", err)
	}
	return tokenBytes, claims, nil
}

func VerifyCapability(tokenBytes []byte, expectedIssuer string, expectedAudience string, expectedProtocolCID string, expectedKind string, expectedAction string, now time.Time) (CapabilityClaims, error) {
	var message cose.Sign1Message
	if err := message.UnmarshalCBOR(tokenBytes); err != nil {
		return CapabilityClaims{}, fmt.Errorf("decode token: %w", err)
	}
	verifier, err := cose.NewVerifier(cose.AlgorithmEdDSA, PublicKey(expectedIssuer))
	if err != nil {
		return CapabilityClaims{}, fmt.Errorf("new token verifier: %w", err)
	}
	if err := message.Verify(nil, verifier); err != nil {
		return CapabilityClaims{}, fmt.Errorf("verify token signature: %w", err)
	}
	claims, err := unmarshalClaims(message.Payload)
	if err != nil {
		return CapabilityClaims{}, err
	}
	switch {
	case claims.Issuer != expectedIssuer:
		return CapabilityClaims{}, fmt.Errorf("capability issuer %q", claims.Issuer)
	case claims.Audience != expectedAudience:
		return CapabilityClaims{}, fmt.Errorf("capability audience %q", claims.Audience)
	case claims.ProtocolCID != expectedProtocolCID:
		return CapabilityClaims{}, fmt.Errorf("capability protocol %q", claims.ProtocolCID)
	case claims.Kind != expectedKind:
		return CapabilityClaims{}, fmt.Errorf("capability kind %q", claims.Kind)
	case claims.Action != expectedAction:
		return CapabilityClaims{}, fmt.Errorf("capability action %q", claims.Action)
	case claims.Subject != "capability_token":
		return CapabilityClaims{}, fmt.Errorf("capability subject %q", claims.Subject)
	case now.UTC().After(claims.ExpiresAt.UTC()):
		return CapabilityClaims{}, fmt.Errorf("capability token expired at %s", claims.ExpiresAt.UTC())
	}
	return claims, nil
}

func marshalClaims(claims CapabilityClaims) ([]byte, error) {
	terms := map[any]any{
		claimIssuer:     claims.Issuer,
		claimSubject:    claims.Subject,
		claimAudience:   claims.Audience,
		claimExpiration: claims.ExpiresAt.Unix(),
		claimIssuedAt:   claims.IssuedAt.Unix(),
		claimTokenID:    []byte(claims.TokenID),
		privateProtocol: claims.ProtocolCID,
		privateKind:     claims.Kind,
		privateAction:   claims.Action,
	}
	mode, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return nil, fmt.Errorf("cbor enc mode: %w", err)
	}
	return mode.Marshal(terms)
}

func unmarshalClaims(payload []byte) (CapabilityClaims, error) {
	var terms map[any]any
	if err := cbor.Unmarshal(payload, &terms); err != nil {
		return CapabilityClaims{}, fmt.Errorf("decode cwt claims: %w", err)
	}
	claims := CapabilityClaims{
		Issuer:      stringTerm(cwtValue(terms, claimIssuer)),
		Subject:     stringTerm(cwtValue(terms, claimSubject)),
		Audience:    stringTerm(cwtValue(terms, claimAudience)),
		ProtocolCID: stringTerm(cwtValue(terms, privateProtocol)),
		Kind:        stringTerm(cwtValue(terms, privateKind)),
		Action:      stringTerm(cwtValue(terms, privateAction)),
		TokenID:     bytesTextTerm(cwtValue(terms, claimTokenID)),
	}
	claims.IssuedAt = time.Unix(int64Term(cwtValue(terms, claimIssuedAt)), 0).UTC()
	claims.ExpiresAt = time.Unix(int64Term(cwtValue(terms, claimExpiration)), 0).UTC()
	return claims, nil
}

func cwtValue(terms map[any]any, label int64) any {
	for _, key := range []any{label, uint64(label), int(label)} {
		if value, ok := terms[key]; ok {
			return value
		}
	}
	return nil
}

func stringTerm(value any) string {
	text, _ := value.(string)
	return text
}

func bytesTextTerm(value any) string {
	switch raw := value.(type) {
	case []byte:
		return string(raw)
	case string:
		return raw
	default:
		return ""
	}
}

func int64Term(value any) int64 {
	switch number := value.(type) {
	case int64:
		return number
	case uint64:
		return int64(number)
	case uint32:
		return int64(number)
	case uint16:
		return int64(number)
	case uint8:
		return int64(number)
	case int:
		return int64(number)
	default:
		return 0
	}
}
