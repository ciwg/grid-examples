package store

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
)

type CAS struct {
	root string
}

func OpenCAS(root string) (*CAS, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &CAS{root: root}, nil
}

// Intent: Keep large package payload bodies content-addressed and immutable so
// append-only records can point at durable bytes without rewriting history.
// Source: DI-moksu
func (cas *CAS) Put(body []byte) (string, error) {
	sum := sha256.Sum256(body)
	objectID := "sha256:" + hex.EncodeToString(sum[:])
	path := cas.pathFor(objectID)
	if _, err := os.Stat(path); err == nil {
		return objectID, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		return "", err
	}
	return objectID, nil
}

func (cas *CAS) Get(objectID string) ([]byte, error) {
	return os.ReadFile(cas.pathFor(objectID))
}

func (cas *CAS) pathFor(objectID string) string {
	return filepath.Join(cas.root, objectID)
}
