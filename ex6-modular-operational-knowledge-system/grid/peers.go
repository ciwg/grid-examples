package grid

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type AllowedPeer struct {
	PeerID    string `json:"peer_id"`
	BatchURL  string `json:"batch_url"`
	ImportURL string `json:"import_url"`
	AllowPull bool   `json:"allow_pull"`
	AllowPush bool   `json:"allow_push"`
}

type PeerConfig struct {
	LocalPeerID  string        `json:"local_peer_id"`
	AllowedPeers []AllowedPeer `json:"allowed_peers"`
}

type PeerStore struct {
	mu     sync.RWMutex
	path   string
	config PeerConfig
}

func OpenPeerStore(root string, runtimeRoot string) (*PeerStore, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(root, "peers.json")
	store := &PeerStore{
		path: path,
		config: PeerConfig{
			LocalPeerID: defaultLocalPeerID(runtimeRoot),
		},
	}
	body, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := store.persistLocked(); err != nil {
				return nil, err
			}
			return store, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(body, &store.config); err != nil {
		return nil, err
	}
	if strings.TrimSpace(store.config.LocalPeerID) == "" {
		store.config.LocalPeerID = defaultLocalPeerID(runtimeRoot)
	}
	slices.SortFunc(store.config.AllowedPeers, func(left, right AllowedPeer) int {
		return strings.Compare(left.PeerID, right.PeerID)
	})
	if err := store.validateLocked(); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *PeerStore) LocalPeerID() string {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.config.LocalPeerID
}

func (store *PeerStore) AllowedPeers() []AllowedPeer {
	store.mu.RLock()
	defer store.mu.RUnlock()
	out := make([]AllowedPeer, len(store.config.AllowedPeers))
	copy(out, store.config.AllowedPeers)
	return out
}

func (store *PeerStore) Lookup(peerID string) (AllowedPeer, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, peer := range store.config.AllowedPeers {
		if peer.PeerID == peerID {
			return peer, true
		}
	}
	return AllowedPeer{}, false
}

func (store *PeerStore) SetAllowedPeer(peer AllowedPeer) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := validateAllowedPeer(peer); err != nil {
		return err
	}
	replaced := false
	for index := range store.config.AllowedPeers {
		if store.config.AllowedPeers[index].PeerID == peer.PeerID {
			store.config.AllowedPeers[index] = peer
			replaced = true
			break
		}
	}
	if !replaced {
		store.config.AllowedPeers = append(store.config.AllowedPeers, peer)
	}
	slices.SortFunc(store.config.AllowedPeers, func(left, right AllowedPeer) int {
		return strings.Compare(left.PeerID, right.PeerID)
	})
	return store.persistLocked()
}

func (store *PeerStore) RemoveAllowedPeer(peerID string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	filtered := store.config.AllowedPeers[:0]
	found := false
	for _, peer := range store.config.AllowedPeers {
		if peer.PeerID == peerID {
			found = true
			continue
		}
		filtered = append(filtered, peer)
	}
	if !found {
		return fmt.Errorf("unknown peer: %s", peerID)
	}
	store.config.AllowedPeers = filtered
	return store.persistLocked()
}

func (store *PeerStore) AllowsPull(peerID string) bool {
	peer, ok := store.Lookup(peerID)
	return ok && peer.AllowPull
}

func (store *PeerStore) AllowsPush(peerID string) bool {
	peer, ok := store.Lookup(peerID)
	return ok && peer.AllowPush
}

func (store *PeerStore) validateLocked() error {
	if strings.TrimSpace(store.config.LocalPeerID) == "" {
		return errors.New("local_peer_id is required")
	}
	seen := map[string]struct{}{}
	for _, peer := range store.config.AllowedPeers {
		if err := validateAllowedPeer(peer); err != nil {
			return err
		}
		if _, exists := seen[peer.PeerID]; exists {
			return fmt.Errorf("duplicate peer_id: %s", peer.PeerID)
		}
		seen[peer.PeerID] = struct{}{}
	}
	return nil
}

func (store *PeerStore) persistLocked() error {
	if err := store.validateLocked(); err != nil {
		return err
	}
	body, err := json.MarshalIndent(store.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(store.path, body, 0o644)
}

func validateAllowedPeer(peer AllowedPeer) error {
	if strings.TrimSpace(peer.PeerID) == "" {
		return errors.New("peer_id is required")
	}
	if strings.TrimSpace(peer.BatchURL) == "" {
		return errors.New("batch_url is required")
	}
	if strings.TrimSpace(peer.ImportURL) == "" {
		return errors.New("import_url is required")
	}
	if !peer.AllowPull && !peer.AllowPush {
		return errors.New("peer must allow pull, push, or both")
	}
	return nil
}

func defaultLocalPeerID(runtimeRoot string) string {
	sum := sha256.Sum256([]byte(runtimeRoot))
	return "peer-" + hex.EncodeToString(sum[:6])
}
