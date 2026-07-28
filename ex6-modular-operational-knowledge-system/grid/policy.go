package grid

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type ClaimTrustPolicy struct {
	ProtocolPCID     string   `json:"protocol_pcid"`
	Role             string   `json:"role"`
	MinAttesters     int      `json:"min_attesters"`
	MinTrustWeight   int      `json:"min_trust_weight,omitempty"`
	AllowedAttesters []string `json:"allowed_attesters,omitempty"`
	AllowedClasses   []string `json:"allowed_classes,omitempty"`
}

type TrustPolicy struct {
	ClaimPolicies []ClaimTrustPolicy `json:"claim_policies"`
}

type PolicyStore struct {
	mu     sync.RWMutex
	path   string
	policy TrustPolicy
}

func OpenPolicyStore(root string) (*PolicyStore, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	store := &PolicyStore{
		path:   filepath.Join(root, "attestation-policy.json"),
		policy: TrustPolicy{},
	}
	body, err := os.ReadFile(store.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := store.persistLocked(); err != nil {
				return nil, err
			}
			return store, nil
		}
		return nil, err
	}
	if err := unmarshalJSON(body, &store.policy); err != nil {
		return nil, err
	}
	if err := store.validateLocked(); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *PolicyStore) ClaimPolicies() []ClaimTrustPolicy {
	store.mu.RLock()
	defer store.mu.RUnlock()
	out := make([]ClaimTrustPolicy, len(store.policy.ClaimPolicies))
	copy(out, store.policy.ClaimPolicies)
	return out
}

func (store *PolicyStore) SetClaimPolicy(policy ClaimTrustPolicy) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	policy = normalizeClaimTrustPolicy(policy)
	if err := validateClaimTrustPolicy(policy); err != nil {
		return err
	}
	replaced := false
	for index := range store.policy.ClaimPolicies {
		current := store.policy.ClaimPolicies[index]
		if current.ProtocolPCID == policy.ProtocolPCID && current.Role == policy.Role {
			store.policy.ClaimPolicies[index] = policy
			replaced = true
			break
		}
	}
	if !replaced {
		store.policy.ClaimPolicies = append(store.policy.ClaimPolicies, policy)
	}
	sortClaimPolicies(store.policy.ClaimPolicies)
	return store.persistLocked()
}

func (store *PolicyStore) RemoveClaimPolicy(protocolPCID string, role string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	protocolPCID = strings.TrimSpace(protocolPCID)
	role = normalizeClaimRole(role)
	filtered := store.policy.ClaimPolicies[:0]
	found := false
	for _, policy := range store.policy.ClaimPolicies {
		if policy.ProtocolPCID == protocolPCID && policy.Role == role {
			found = true
			continue
		}
		filtered = append(filtered, policy)
	}
	if !found {
		return fmt.Errorf("unknown claim policy: %s %s", protocolPCID, role)
	}
	store.policy.ClaimPolicies = filtered
	return store.persistLocked()
}

// Intent: Keep claim-attestation trust local and explicit by evaluating import
// quorums against runtime-owned policy instead of treating any valid
// countersignature as automatically sufficient. Source: DI-movek
func (store *PolicyStore) VerifyClaimPolicies(batch Batch, knownPeers []AllowedPeer) error {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for index, claim := range batch.ImplementationClaims {
		policy, ok := store.matchClaimPolicyLocked(claim)
		if !ok || (policy.MinAttesters == 0 && policy.MinTrustWeight == 0) {
			continue
		}
		matched := map[string]struct{}{}
		weight := 0
		for _, attestation := range batch.ClaimAttestations {
			if attestation.ClaimIndex != index {
				continue
			}
			peer, ok := findKnownPeer(knownPeers, attestation.SignerPeerID)
			if !ok {
				continue
			}
			if len(policy.AllowedAttesters) > 0 && !slices.Contains(policy.AllowedAttesters, attestation.SignerPeerID) {
				continue
			}
			if len(policy.AllowedClasses) > 0 && !slices.Contains(policy.AllowedClasses, peer.AttesterClass) {
				continue
			}
			if _, exists := matched[attestation.SignerPeerID]; exists {
				continue
			}
			matched[attestation.SignerPeerID] = struct{}{}
			weight += peer.AttestationWeight
		}
		if len(matched) < policy.MinAttesters {
			return fmt.Errorf(
				"claim attestation quorum failed for %s role %s: need %d, got %d",
				claim.ProtocolPCID,
				claim.Role,
				policy.MinAttesters,
				len(matched),
			)
		}
		if weight < policy.MinTrustWeight {
			return fmt.Errorf(
				"claim attestation trust weight failed for %s role %s: need %d, got %d",
				claim.ProtocolPCID,
				claim.Role,
				policy.MinTrustWeight,
				weight,
			)
		}
	}
	return nil
}

func (store *PolicyStore) validateLocked() error {
	seen := map[string]struct{}{}
	for _, policy := range store.policy.ClaimPolicies {
		if err := validateClaimTrustPolicy(policy); err != nil {
			return err
		}
		key := policy.ProtocolPCID + "\x00" + policy.Role
		if _, exists := seen[key]; exists {
			return fmt.Errorf("duplicate claim policy: %s %s", policy.ProtocolPCID, policy.Role)
		}
		seen[key] = struct{}{}
	}
	sortClaimPolicies(store.policy.ClaimPolicies)
	return nil
}

func (store *PolicyStore) persistLocked() error {
	if err := store.validateLocked(); err != nil {
		return err
	}
	body, err := marshalIndentJSON(store.policy)
	if err != nil {
		return err
	}
	return os.WriteFile(store.path, body, 0o644)
}

func (store *PolicyStore) matchClaimPolicyLocked(claim ImplementationClaim) (ClaimTrustPolicy, bool) {
	for _, policy := range store.policy.ClaimPolicies {
		if policy.ProtocolPCID == claim.ProtocolPCID && policy.Role == claim.Role {
			return policy, true
		}
	}
	for _, policy := range store.policy.ClaimPolicies {
		if policy.ProtocolPCID == claim.ProtocolPCID && policy.Role == "*" {
			return policy, true
		}
	}
	return ClaimTrustPolicy{}, false
}

func validateClaimTrustPolicy(policy ClaimTrustPolicy) error {
	if strings.TrimSpace(policy.ProtocolPCID) == "" {
		return errors.New("protocol_pcid is required")
	}
	if normalizeClaimRole(policy.Role) == "" {
		return errors.New("role is required")
	}
	if policy.MinAttesters < 0 {
		return errors.New("min_attesters must be zero or greater")
	}
	if policy.MinTrustWeight < 0 {
		return errors.New("min_trust_weight must be zero or greater")
	}
	seen := map[string]struct{}{}
	for _, peerID := range policy.AllowedAttesters {
		if strings.TrimSpace(peerID) == "" {
			return errors.New("allowed_attesters cannot contain blanks")
		}
		if _, exists := seen[peerID]; exists {
			return fmt.Errorf("duplicate allowed_attester: %s", peerID)
		}
		seen[peerID] = struct{}{}
	}
	seenClasses := map[string]struct{}{}
	for _, class := range policy.AllowedClasses {
		if strings.TrimSpace(class) == "" {
			return errors.New("allowed_classes cannot contain blanks")
		}
		if _, exists := seenClasses[class]; exists {
			return fmt.Errorf("duplicate allowed_class: %s", class)
		}
		seenClasses[class] = struct{}{}
	}
	if len(policy.AllowedAttesters) > 0 && policy.MinAttesters > len(policy.AllowedAttesters) {
		return errors.New("min_attesters cannot exceed allowed_attesters length")
	}
	return nil
}

func normalizeClaimTrustPolicy(policy ClaimTrustPolicy) ClaimTrustPolicy {
	policy.ProtocolPCID = strings.TrimSpace(policy.ProtocolPCID)
	policy.Role = normalizeClaimRole(policy.Role)
	for index := range policy.AllowedAttesters {
		policy.AllowedAttesters[index] = strings.TrimSpace(policy.AllowedAttesters[index])
	}
	slices.Sort(policy.AllowedAttesters)
	policy.AllowedAttesters = slices.Compact(policy.AllowedAttesters)
	for index := range policy.AllowedClasses {
		policy.AllowedClasses[index] = strings.TrimSpace(policy.AllowedClasses[index])
	}
	slices.Sort(policy.AllowedClasses)
	policy.AllowedClasses = slices.Compact(policy.AllowedClasses)
	return policy
}

func normalizeClaimRole(role string) string {
	role = strings.TrimSpace(role)
	if role == "" {
		return "*"
	}
	return role
}

func sortClaimPolicies(policies []ClaimTrustPolicy) {
	slices.SortFunc(policies, func(left, right ClaimTrustPolicy) int {
		if left.ProtocolPCID != right.ProtocolPCID {
			return strings.Compare(left.ProtocolPCID, right.ProtocolPCID)
		}
		return strings.Compare(left.Role, right.Role)
	})
}

func findKnownPeer(peers []AllowedPeer, peerID string) (AllowedPeer, bool) {
	for _, peer := range peers {
		if peer.PeerID == peerID {
			return peer, true
		}
	}
	return AllowedPeer{}, false
}
