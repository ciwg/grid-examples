package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

type RuntimeIdentity struct {
	privateKey ed25519.PrivateKey
}

// Intent: Keep the first ex5 PromiseGrid-native signing identity in one
// runtime-owned seed file so knowledge-item envelopes have stable local
// authorship and replay verification. Source: DI-mibor
func LoadOrCreateRuntimeIdentity(path string) (*RuntimeIdentity, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir identity dir: %w", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		seed := make([]byte, ed25519.SeedSize)
		if _, err := rand.Read(seed); err != nil {
			return nil, fmt.Errorf("read seed: %w", err)
		}
		if err := os.WriteFile(path, seed, 0o600); err != nil {
			return nil, fmt.Errorf("write seed: %w", err)
		}
	}
	seed, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read seed: %w", err)
	}
	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("unexpected seed size %d", len(seed))
	}
	return &RuntimeIdentity{privateKey: ed25519.NewKeyFromSeed(seed)}, nil
}

func (identity *RuntimeIdentity) PublicKey() ed25519.PublicKey {
	return identity.privateKey.Public().(ed25519.PublicKey)
}

func (identity *RuntimeIdentity) KeyID() string {
	return hex.EncodeToString(identity.PublicKey())
}

func (identity *RuntimeIdentity) PeerID() string {
	return identity.KeyID()
}

func (identity *RuntimeIdentity) SignProof(signable []byte) ([]byte, error) {
	proof := struct {
		Algorithm string `cbor:"algorithm"`
		KeyID     string `cbor:"key_id"`
		PublicKey []byte `cbor:"public_key"`
		Signature []byte `cbor:"signature"`
	}{
		Algorithm: "Ed25519",
		KeyID:     identity.KeyID(),
		PublicKey: append([]byte(nil), identity.PublicKey()...),
		Signature: ed25519.Sign(identity.privateKey, signable),
	}
	return protocols.Marshal(proof)
}

func VerifyRuntimeProof(signable []byte, proofBytes []byte) error {
	var proof struct {
		Algorithm string `cbor:"algorithm"`
		KeyID     string `cbor:"key_id"`
		PublicKey []byte `cbor:"public_key"`
		Signature []byte `cbor:"signature"`
	}
	if err := protocols.Unmarshal(proofBytes, &proof); err != nil {
		return fmt.Errorf("decode proof: %w", err)
	}
	if proof.Algorithm != "Ed25519" {
		return fmt.Errorf("unsupported proof algorithm %q", proof.Algorithm)
	}
	if len(proof.PublicKey) != ed25519.PublicKeySize {
		return fmt.Errorf("unexpected public key size %d", len(proof.PublicKey))
	}
	if proof.KeyID != hex.EncodeToString(proof.PublicKey) {
		return fmt.Errorf("proof key id mismatch")
	}
	if len(proof.Signature) != ed25519.SignatureSize {
		return fmt.Errorf("unexpected signature size %d", len(proof.Signature))
	}
	if !ed25519.Verify(ed25519.PublicKey(proof.PublicKey), signable, proof.Signature) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}
