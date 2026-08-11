package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
)

const deviceIdentityFile = "device-key.json"

type encryptedDeviceIdentity struct {
	Version    int    `json:"version"`
	Salt       string `json:"argon2id_salt_base64"`
	Nonce      string `json:"aes_gcm_nonce_base64"`
	Ciphertext string `json:"ciphertext_base64"`
}

type deviceIdentityPlaintext struct {
	Label      string `json:"label"`
	PrivateKey string `json:"device_private_key_base64"`
}

// ParticipantIdentity is a participant-owned local signing embodiment. Intent:
// the encrypted private key is local custody only; signed history remains the
// public author-evidence boundary. Source: DI-hibok.
type ParticipantIdentity struct {
	Label   string
	Private ed25519.PrivateKey
}

func CreateParticipantIdentity(agentRoot, label string, passphrase []byte) (*ParticipantIdentity, error) {
	if label == "" || len(passphrase) == 0 {
		return nil, errors.New("identity label and passphrase are required")
	}
	identity := &ParticipantIdentity{Label: label, Private: ed25519.NewKeyFromSeed(randomIdentitySeed())}
	if err := identity.Save(agentRoot, passphrase); err != nil {
		return nil, err
	}
	return identity, nil
}

func LoadParticipantIdentity(agentRoot string, passphrase []byte) (*ParticipantIdentity, error) {
	if len(passphrase) == 0 {
		return nil, errors.New("identity passphrase is required")
	}
	path := filepath.Join(agentRoot, "identity", deviceIdentityFile)
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.Mode().Perm() != 0o600 {
		return nil, errors.New("identity file must have mode 0600")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var encrypted encryptedDeviceIdentity
	if err := json.Unmarshal(raw, &encrypted); err != nil || encrypted.Version != 1 {
		return nil, errors.New("invalid encrypted identity file")
	}
	salt, err := base64.StdEncoding.DecodeString(encrypted.Salt)
	if err != nil || len(salt) != 16 {
		return nil, errors.New("invalid identity salt")
	}
	nonce, err := base64.StdEncoding.DecodeString(encrypted.Nonce)
	if err != nil {
		return nil, errors.New("invalid identity nonce")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted.Ciphertext)
	if err != nil {
		return nil, errors.New("invalid identity ciphertext")
	}
	block, err := aes.NewCipher(identityKey(passphrase, salt))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil || len(nonce) != gcm.NonceSize() {
		return nil, errors.New("invalid identity encryption")
	}
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("identity passphrase did not decrypt key")
	}
	var decoded deviceIdentityPlaintext
	if err := json.Unmarshal(plain, &decoded); err != nil {
		return nil, err
	}
	private, err := base64.StdEncoding.DecodeString(decoded.PrivateKey)
	if err != nil || len(private) != ed25519.PrivateKeySize || decoded.Label == "" {
		return nil, errors.New("invalid identity plaintext")
	}
	return &ParticipantIdentity{Label: decoded.Label, Private: ed25519.PrivateKey(private)}, nil
}

func (i *ParticipantIdentity) Save(agentRoot string, passphrase []byte) error {
	if i == nil || i.Label == "" || len(i.Private) != ed25519.PrivateKeySize || len(passphrase) == 0 {
		return errors.New("invalid participant identity")
	}
	directory := filepath.Join(agentRoot, "identity")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return err
	}
	salt, nonce := make([]byte, 16), make([]byte, 12)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	plain, err := json.Marshal(deviceIdentityPlaintext{Label: i.Label, PrivateKey: base64.StdEncoding.EncodeToString(i.Private)})
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(identityKey(passphrase, salt))
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	if len(nonce) != gcm.NonceSize() {
		return fmt.Errorf("identity nonce size")
	}
	encrypted, err := json.Marshal(encryptedDeviceIdentity{Version: 1, Salt: base64.StdEncoding.EncodeToString(salt), Nonce: base64.StdEncoding.EncodeToString(nonce), Ciphertext: base64.StdEncoding.EncodeToString(gcm.Seal(nil, nonce, plain, nil))})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(directory, deviceIdentityFile), encrypted, 0o600)
}

func identityKey(passphrase, salt []byte) []byte {
	return argon2.IDKey(passphrase, salt, 3, 64*1024, 1, 32)
}
func randomIdentitySeed() []byte {
	seed := make([]byte, ed25519.SeedSize)
	if _, err := rand.Read(seed); err != nil {
		panic(err)
	}
	return seed
}
