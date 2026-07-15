package cas

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocol"
)

type Store struct {
	root string
}

func Open(root string) (*Store, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir cas root: %w", err)
	}
	return &Store{root: root}, nil
}

// Intent: Persist exact signed envelope bytes under a stable content-addressed
// path so replay, verification, and relay forwarding all share the same byte
// identity without creating a canonical document host. Source: DI-ramuv
func (store *Store) Put(bytes []byte) (string, error) {
	cidValue, err := protocol.CIDForBytes(bytes)
	if err != nil {
		return "", fmt.Errorf("cid for bytes: %w", err)
	}
	address := cidValue.String()
	path := store.pathFor(address)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("mkdir cid dir: %w", err)
	}
	if _, err := os.Stat(path); err == nil {
		return address, nil
	}
	if err := os.WriteFile(path, bytes, 0o644); err != nil {
		return "", fmt.Errorf("write cid bytes: %w", err)
	}
	return address, nil
}

func (store *Store) Get(address string) ([]byte, error) {
	bytes, err := os.ReadFile(store.pathFor(address))
	if err != nil {
		return nil, fmt.Errorf("read cid bytes: %w", err)
	}
	cidValue, err := protocol.CIDForBytes(bytes)
	if err != nil {
		return nil, fmt.Errorf("recompute cid: %w", err)
	}
	if cidValue.String() != address {
		return nil, fmt.Errorf("cid mismatch: got %s want %s", cidValue, address)
	}
	return bytes, nil
}

func (store *Store) pathFor(address string) string {
	prefix := address
	if len(prefix) > 2 {
		prefix = prefix[:2]
	}
	return filepath.Join(store.root, prefix, sanitizeAddress(address))
}

func sanitizeAddress(address string) string {
	return strings.ReplaceAll(address, "/", "_")
}
