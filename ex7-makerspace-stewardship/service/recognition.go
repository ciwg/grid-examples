package service

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// RecognitionPolicy binds locally recognized participant labels to full public
// key fingerprints. Intent: keep local role assessment separate from a
// self-declared record label or signature validity alone. Source: DI-piruf.
type RecognitionPolicy struct {
	keyLabels map[string]string
}

// NewRecognitionPolicy creates local bootstrap policy from public keys only;
// no participant private key is accepted or retained by the runtime.
func NewRecognitionPolicy(keys map[string]ed25519.PublicKey) RecognitionPolicy {
	keyLabels := make(map[string]string, len(keys))
	for label, key := range keys {
		keyLabels[keyID(key)] = label
	}
	return RecognitionPolicy{keyLabels: keyLabels}
}

func (p RecognitionPolicy) recognizes(record Record) bool {
	return p.keyLabels[record.KeyID] == record.Signer
}

type recognitionFile struct {
	Version int                `json:"version"`
	Keys    []recognitionEntry `json:"keys"`
}

type recognitionEntry struct {
	Label            string `json:"label"`
	Ed25519PublicKey string `json:"ed25519_public_key_base64"`
}

// LoadRecognitionPolicy reads public local recognition input only. Intent:
// make local assessment repeatable without accepting private keys or browser
// policy mutation. Source: DI-likoh.
func LoadRecognitionPolicy(path string, allowEmpty bool) (RecognitionPolicy, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) && allowEmpty {
		return RecognitionPolicy{}, nil
	}
	if err != nil {
		return RecognitionPolicy{}, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return RecognitionPolicy{}, err
	}
	if info.Mode().Perm() != 0o600 {
		return RecognitionPolicy{}, errors.New("recognition policy must have mode 0600")
	}
	var file recognitionFile
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&file); err != nil {
		return RecognitionPolicy{}, err
	}
	if file.Version != 1 {
		return RecognitionPolicy{}, errors.New("unsupported recognition policy version")
	}
	keys := make(map[string]ed25519.PublicKey, len(file.Keys))
	seenKeyIDs := make(map[string]struct{}, len(file.Keys))
	for _, entry := range file.Keys {
		if entry.Label == "" {
			return RecognitionPolicy{}, errors.New("recognition entry label is required")
		}
		if _, exists := keys[entry.Label]; exists {
			return RecognitionPolicy{}, fmt.Errorf("duplicate recognition label %q", entry.Label)
		}
		decoded, err := base64.StdEncoding.DecodeString(entry.Ed25519PublicKey)
		if err != nil || len(decoded) != ed25519.PublicKeySize {
			return RecognitionPolicy{}, fmt.Errorf("invalid public key for %q", entry.Label)
		}
		key := ed25519.PublicKey(decoded)
		id := keyID(key)
		if _, exists := seenKeyIDs[id]; exists {
			return RecognitionPolicy{}, fmt.Errorf("duplicate recognition public key for %q", entry.Label)
		}
		seenKeyIDs[id] = struct{}{}
		keys[entry.Label] = key
	}
	return NewRecognitionPolicy(keys), nil
}
