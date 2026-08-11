package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"os"
	"path/filepath"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocol"
)

// AgentID is a local identifier derived from a public signing key.
type AgentID string

// AgentKey retains private signing material only in the local embodiment that
// owns it. The server uses the public key supplied by an Enrollment instead.
type AgentKey struct {
	private ed25519.PrivateKey
	public  ed25519.PublicKey
}

// Enrollment is the server-storable public local-admission binding.
type Enrollment struct {
	AgentID   AgentID `json:"agent_id"`
	PublicKey []byte  `json:"public_key"`
	Role      string  `json:"role"`
}

// EnrollmentProof is detached proof of possession for a public enrollment claim.
type EnrollmentProof struct {
	Signature []byte `json:"signature"`
}

func NewAgentKey() (AgentKey, error) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return AgentKey{}, fmt.Errorf("generate agent key: %w", err)
	}
	return AgentKey{private: private, public: public}, nil
}

// LoadOrCreateAgentKey keeps CLI private material at an explicit operator path.
// Intent: Make CLI key custody local and durable without server persistence.
// Source: DI-muzal
func LoadOrCreateAgentKey(path string) (AgentKey, error) {
	if path == "" {
		return AgentKey{}, fmt.Errorf("agent key path is required")
	}
	private, err := os.ReadFile(path)
	if err == nil {
		if len(private) != ed25519.PrivateKeySize {
			return AgentKey{}, fmt.Errorf("invalid agent private key length")
		}
		key := ed25519.PrivateKey(private)
		return AgentKey{private: key, public: key.Public().(ed25519.PublicKey)}, nil
	}
	if !os.IsNotExist(err) {
		return AgentKey{}, fmt.Errorf("read agent key: %w", err)
	}
	key, err := NewAgentKey()
	if err != nil {
		return AgentKey{}, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return AgentKey{}, fmt.Errorf("create agent key directory: %w", err)
	}
	if err := os.WriteFile(path, key.private, 0o600); err != nil {
		return AgentKey{}, fmt.Errorf("write agent key: %w", err)
	}
	return key, nil
}

func AgentIDForPublicKey(publicKey []byte) AgentID {
	sum := sha256.Sum256(publicKey)
	return AgentID("agent-" + base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(sum[:]))
}

func (key AgentKey) AgentID() AgentID {
	return AgentIDForPublicKey(key.public)
}

func (key AgentKey) PublicKey() []byte {
	return append([]byte(nil), key.public...)
}

func (key AgentKey) Sign(message []byte) []byte {
	return ed25519.Sign(key.private, message)
}

func Verify(publicKey []byte, message []byte, signature []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(publicKey), message, signature)
}

// Intent: Keep server enrollment limited to a public key and local role label,
// never client private key material or a generalized authorization claim. Source: DI-muzal
func NewEnrollment(key AgentKey, role string) Enrollment {
	return Enrollment{AgentID: key.AgentID(), PublicKey: key.PublicKey(), Role: role}
}

func (key AgentKey) SignEnrollment(enrollment Enrollment) (EnrollmentProof, error) {
	claim, err := protocol.Marshal(enrollment)
	if err != nil {
		return EnrollmentProof{}, fmt.Errorf("marshal enrollment claim: %w", err)
	}
	return EnrollmentProof{Signature: key.Sign(claim)}, nil
}

// Intent: Verify an adapter's key-possession proof before persisting only its
// public local admission binding. Source: DI-nusop
func VerifyEnrollment(enrollment Enrollment, proof EnrollmentProof) error {
	if enrollment.AgentID != AgentIDForPublicKey(enrollment.PublicKey) {
		return fmt.Errorf("agent ID does not match public key")
	}
	claim, err := protocol.Marshal(enrollment)
	if err != nil {
		return fmt.Errorf("marshal enrollment claim: %w", err)
	}
	if !Verify(enrollment.PublicKey, claim, proof.Signature) {
		return fmt.Errorf("invalid enrollment proof")
	}
	return nil
}
