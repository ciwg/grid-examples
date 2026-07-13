package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/protocol"
)

type Identity struct {
	privateKey ed25519.PrivateKey
}

// Intent: Keep the durable signing identity in one local service-owned place so
// both embodiments share the same author identity and replay story. Source:
// DI-jilin
func LoadOrCreate(path string) (*Identity, error) {
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
	return &Identity{privateKey: ed25519.NewKeyFromSeed(seed)}, nil
}

func (identity *Identity) PublicKey() ed25519.PublicKey {
	return identity.privateKey.Public().(ed25519.PublicKey)
}

func (identity *Identity) KeyID() string {
	return hex.EncodeToString(identity.PublicKey())
}

func (identity *Identity) SignProof(signable []byte) (protocol.Proof, error) {
	signature := ed25519.Sign(identity.privateKey, signable)
	return protocol.Proof{
		Algorithm: "Ed25519",
		KeyID:     identity.KeyID(),
		PublicKey: append([]byte(nil), identity.PublicKey()...),
		Signature: append([]byte(nil), signature...),
	}, nil
}

func VerifyProof(signable []byte, proof protocol.Proof) error {
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
