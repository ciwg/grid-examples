package cas

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocol"
)

// Store retains immutable bytes at cas/<CID> under one Ex4 runtime root.
type Store struct {
	root string
}

func Open(root string) (*Store, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("create CAS root: %w", err)
	}
	return &Store{root: root}, nil
}

func (store *Store) Put(bytes []byte) (string, error) {
	cid, err := protocol.CIDForBytes(bytes)
	if err != nil {
		return "", fmt.Errorf("derive object CID: %w", err)
	}
	path := filepath.Join(store.root, cid.String())
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) != string(bytes) {
			return "", fmt.Errorf("CID collision at %s", cid)
		}
		return cid.String(), nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("read CAS object: %w", err)
	}
	if err := os.WriteFile(path, bytes, 0o644); err != nil {
		return "", fmt.Errorf("write CAS object: %w", err)
	}
	return cid.String(), nil
}

func (store *Store) Get(cid string) ([]byte, error) {
	bytes, err := os.ReadFile(filepath.Join(store.root, cid))
	if err != nil {
		return nil, fmt.Errorf("read CAS object: %w", err)
	}
	return bytes, nil
}
