package grid

import (
	"crypto/ed25519"
	"crypto/rand"
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

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/records"
)

type AllowedPeer struct {
	PeerID    string `json:"peer_id"`
	BatchURL  string `json:"batch_url"`
	ImportURL string `json:"import_url"`
	PublicKey string `json:"public_key"`
	AllowPull bool   `json:"allow_pull"`
	AllowPush bool   `json:"allow_push"`
}

type PeerConfig struct {
	LocalPeerID     string        `json:"local_peer_id"`
	LocalPublicKey  string        `json:"local_public_key"`
	LocalPrivateKey string        `json:"local_private_key"`
	AllowedPeers    []AllowedPeer `json:"allowed_peers"`
}

type PeerCard struct {
	PeerID      string `json:"peer_id"`
	PublicKey   string `json:"public_key"`
	BatchURL    string `json:"batch_url"`
	ImportURL   string `json:"import_url"`
	DiscoverURL string `json:"discover_url"`
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
		path:   path,
		config: PeerConfig{},
	}
	body, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := store.ensureLocalIdentityLocked(runtimeRoot); err != nil {
				return nil, err
			}
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
	if err := store.ensureLocalIdentityLocked(runtimeRoot); err != nil {
		return nil, err
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

func (store *PeerStore) LocalPublicKey() string {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.config.LocalPublicKey
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

func (store *PeerStore) SignBatch(batch Batch) (Batch, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	privateKey, err := hex.DecodeString(store.config.LocalPrivateKey)
	if err != nil {
		return Batch{}, err
	}
	signingBytes, err := batch.SigningBytes()
	if err != nil {
		return Batch{}, err
	}
	batch.Signature = hex.EncodeToString(ed25519.Sign(ed25519.PrivateKey(privateKey), signingBytes))
	return batch, nil
}

func (store *PeerStore) SignRecords(records []json.RawMessage) ([]RecordSignature, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	privateKey, err := hex.DecodeString(store.config.LocalPrivateKey)
	if err != nil {
		return nil, err
	}
	signatures := make([]RecordSignature, 0, len(records))
	for _, record := range records {
		signingBytes := RecordSigningBytes(record)
		signatures = append(signatures, RecordSignature{
			SignerPeerID: store.config.LocalPeerID,
			PublicKey:    store.config.LocalPublicKey,
			Signature:    hex.EncodeToString(ed25519.Sign(ed25519.PrivateKey(privateKey), signingBytes)),
		})
	}
	return signatures, nil
}

func (store *PeerStore) SignClaims(claims []ImplementationClaim) ([]ClaimProof, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	privateKey, err := hex.DecodeString(store.config.LocalPrivateKey)
	if err != nil {
		return nil, err
	}
	proofs := make([]ClaimProof, 0, len(claims))
	for _, claim := range claims {
		signingBytes := ClaimSigningBytes(claim)
		proofs = append(proofs, ClaimProof{
			SignerPeerID: store.config.LocalPeerID,
			PublicKey:    store.config.LocalPublicKey,
			Signature:    hex.EncodeToString(ed25519.Sign(ed25519.PrivateKey(privateKey), signingBytes)),
		})
	}
	return proofs, nil
}

func (store *PeerStore) AttestClaims(claims []ImplementationClaim) ([]ClaimAttestation, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	privateKey, err := hex.DecodeString(store.config.LocalPrivateKey)
	if err != nil {
		return nil, err
	}
	attestations := make([]ClaimAttestation, 0, len(claims))
	for index, claim := range claims {
		signingBytes := ClaimAttestationSigningBytes(index, claim)
		attestations = append(attestations, ClaimAttestation{
			ClaimIndex:   index,
			SignerPeerID: store.config.LocalPeerID,
			PublicKey:    store.config.LocalPublicKey,
			Signature:    hex.EncodeToString(ed25519.Sign(ed25519.PrivateKey(privateKey), signingBytes)),
		})
	}
	return attestations, nil
}

func (store *PeerStore) SignAuthorEnvelope(envelope records.Envelope) (records.Envelope, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	privateKey, err := hex.DecodeString(store.config.LocalPrivateKey)
	if err != nil {
		return records.Envelope{}, err
	}
	signingBytes, err := envelope.SigningBytes()
	if err != nil {
		return records.Envelope{}, err
	}
	envelope.AuthorKeyID = store.config.LocalPeerID
	envelope.AuthorPublicKey = store.config.LocalPublicKey
	envelope.AuthorSignature = hex.EncodeToString(ed25519.Sign(ed25519.PrivateKey(privateKey), signingBytes))
	return envelope, nil
}

func (store *PeerStore) VerifyPeerBatch(peerID string, batch Batch) error {
	// Intent: Keep live relay trust rooted in explicit peer registration by
	// verifying each batch against the allowed peer's configured public key.
	// Source: DI-zotem
	peer, ok := store.Lookup(peerID)
	if !ok {
		return fmt.Errorf("peer not allowed: %s", peerID)
	}
	if strings.TrimSpace(peer.PublicKey) == "" {
		return fmt.Errorf("peer %s public key is required", peerID)
	}
	if strings.TrimSpace(batch.Signature) == "" {
		return errors.New("batch signature is required")
	}
	publicKey, err := hex.DecodeString(peer.PublicKey)
	if err != nil {
		return err
	}
	signature, err := hex.DecodeString(batch.Signature)
	if err != nil {
		return err
	}
	signingBytes, err := batch.SigningBytes()
	if err != nil {
		return err
	}
	if !ed25519.Verify(ed25519.PublicKey(publicKey), signingBytes, signature) {
		return errors.New("batch signature verification failed")
	}
	return nil
}

func (store *PeerStore) VerifyRecordSignatures(batch Batch) error {
	// Intent: Bind each carried record to the exporting peer's key material so
	// receivers can verify relay-carriage signatures per record, not just once
	// for the enclosing batch.
	// Source: DI-ravud
	if len(batch.RecordSignatures) == 0 {
		return nil
	}
	if len(batch.RecordSignatures) != len(batch.Records) {
		return errors.New("record_signatures must match records length")
	}
	for index, signature := range batch.RecordSignatures {
		if peerIDFromPublicKey(signature.PublicKey) != signature.SignerPeerID {
			return fmt.Errorf("record signature signer mismatch at index %d", index)
		}
		publicKey, err := hex.DecodeString(signature.PublicKey)
		if err != nil {
			return err
		}
		sigBytes, err := hex.DecodeString(signature.Signature)
		if err != nil {
			return err
		}
		if !ed25519.Verify(ed25519.PublicKey(publicKey), RecordSigningBytes(batch.Records[index]), sigBytes) {
			return fmt.Errorf("record signature verification failed at index %d", index)
		}
	}
	return nil
}

func (store *PeerStore) VerifyClaimProofs(batch Batch) error {
	// Intent: Verify each advertised implementation claim against the exporting
	// peer's key material so claim metadata is not just unsigned batch content.
	// Source: DI-luzef
	if len(batch.ClaimProofs) == 0 {
		return nil
	}
	if len(batch.ClaimProofs) != len(batch.ImplementationClaims) {
		return errors.New("claim_proofs must match implementation_claims length")
	}
	for index, proof := range batch.ClaimProofs {
		if peerIDFromPublicKey(proof.PublicKey) != proof.SignerPeerID {
			return fmt.Errorf("claim proof signer mismatch at index %d", index)
		}
		publicKey, err := hex.DecodeString(proof.PublicKey)
		if err != nil {
			return err
		}
		sigBytes, err := hex.DecodeString(proof.Signature)
		if err != nil {
			return err
		}
		if !ed25519.Verify(ed25519.PublicKey(publicKey), ClaimSigningBytes(batch.ImplementationClaims[index]), sigBytes) {
			return fmt.Errorf("claim proof verification failed at index %d", index)
		}
	}
	return nil
}

func (store *PeerStore) VerifyAuthorEnvelope(envelope records.Envelope) error {
	// Intent: Verify semantic author signatures on durable records so receivers
	// can trust that the named record content was signed at authoring time, not
	// only during later relay carriage.
	// Source: DI-sovem
	if !envelope.HasAuthorSignature() {
		return nil
	}
	if peerIDFromPublicKey(envelope.AuthorPublicKey) != envelope.AuthorKeyID {
		return errors.New("author key id does not match author public key")
	}
	publicKey, err := hex.DecodeString(envelope.AuthorPublicKey)
	if err != nil {
		return err
	}
	signature, err := hex.DecodeString(envelope.AuthorSignature)
	if err != nil {
		return err
	}
	signingBytes, err := envelope.SigningBytes()
	if err != nil {
		return err
	}
	if !ed25519.Verify(ed25519.PublicKey(publicKey), signingBytes, signature) {
		return errors.New("author signature verification failed")
	}
	return nil
}

func (store *PeerStore) VerifyClaimAttestations(batch Batch) error {
	// Intent: Verify third-party claim attestations separately from exporter
	// self-claims so trust can include outside countersigners.
	// Source: DI-fogem
	for _, attestation := range batch.ClaimAttestations {
		if attestation.SignerPeerID == batch.Implementation {
			return fmt.Errorf("claim attestation must be third-party at index %d", attestation.ClaimIndex)
		}
		if peerIDFromPublicKey(attestation.PublicKey) != attestation.SignerPeerID {
			return fmt.Errorf("claim attestation signer mismatch at index %d", attestation.ClaimIndex)
		}
		publicKey, err := hex.DecodeString(attestation.PublicKey)
		if err != nil {
			return err
		}
		sigBytes, err := hex.DecodeString(attestation.Signature)
		if err != nil {
			return err
		}
		claim := batch.ImplementationClaims[attestation.ClaimIndex]
		if !ed25519.Verify(ed25519.PublicKey(publicKey), ClaimAttestationSigningBytes(attestation.ClaimIndex, claim), sigBytes) {
			return fmt.Errorf("claim attestation verification failed at index %d", attestation.ClaimIndex)
		}
	}
	return nil
}

func (store *PeerStore) validateLocked() error {
	if strings.TrimSpace(store.config.LocalPeerID) == "" {
		return errors.New("local_peer_id is required")
	}
	if strings.TrimSpace(store.config.LocalPublicKey) == "" {
		return errors.New("local_public_key is required")
	}
	if strings.TrimSpace(store.config.LocalPrivateKey) == "" {
		return errors.New("local_private_key is required")
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
	// Intent: Allow explicit local registration of untrusted peers so discovery
	// can seed metadata without silently granting pull or push rights.
	// Source: DI-kasud
	if strings.TrimSpace(peer.PeerID) == "" {
		return errors.New("peer_id is required")
	}
	if strings.TrimSpace(peer.BatchURL) == "" {
		return errors.New("batch_url is required")
	}
	if strings.TrimSpace(peer.ImportURL) == "" {
		return errors.New("import_url is required")
	}
	if strings.TrimSpace(peer.PublicKey) == "" {
		return errors.New("public_key is required")
	}
	return nil
}

func (card PeerCard) Validate() error {
	if strings.TrimSpace(card.PeerID) == "" {
		return errors.New("peer_id is required")
	}
	if strings.TrimSpace(card.PublicKey) == "" {
		return errors.New("public_key is required")
	}
	if strings.TrimSpace(card.BatchURL) == "" {
		return errors.New("batch_url is required")
	}
	if strings.TrimSpace(card.ImportURL) == "" {
		return errors.New("import_url is required")
	}
	if strings.TrimSpace(card.DiscoverURL) == "" {
		return errors.New("discover_url is required")
	}
	return nil
}

func (store *PeerStore) ensureLocalIdentityLocked(runtimeRoot string) error {
	if strings.TrimSpace(store.config.LocalPublicKey) != "" && strings.TrimSpace(store.config.LocalPrivateKey) != "" {
		if strings.TrimSpace(store.config.LocalPeerID) == "" {
			store.config.LocalPeerID = peerIDFromPublicKey(store.config.LocalPublicKey)
		}
		return nil
	}
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	store.config.LocalPublicKey = hex.EncodeToString(publicKey)
	store.config.LocalPrivateKey = hex.EncodeToString(privateKey)
	if strings.TrimSpace(store.config.LocalPeerID) == "" {
		store.config.LocalPeerID = peerIDFromPublicKey(store.config.LocalPublicKey)
	}
	if strings.TrimSpace(store.config.LocalPeerID) == "" {
		store.config.LocalPeerID = defaultLocalPeerID(runtimeRoot)
	}
	return nil
}

func defaultLocalPeerID(runtimeRoot string) string {
	sum := sha256.Sum256([]byte("runtime:" + runtimeRoot))
	return "peer-" + hex.EncodeToString(sum[:6])
}

func peerIDFromPublicKey(publicKey string) string {
	sum := sha256.Sum256([]byte(publicKey))
	return "peer-" + hex.EncodeToString(sum[:6])
}
