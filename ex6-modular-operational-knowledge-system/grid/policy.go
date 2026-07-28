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
	ProtocolPCID       string   `json:"protocol_pcid"`
	Role               string   `json:"role"`
	MinAttesters       int      `json:"min_attesters"`
	MinTrustWeight     int      `json:"min_trust_weight,omitempty"`
	MinFederations     int      `json:"min_federations,omitempty"`
	AllowedAttesters   []string `json:"allowed_attesters,omitempty"`
	AllowedClasses     []string `json:"allowed_classes,omitempty"`
	AllowedFederations []string `json:"allowed_federations,omitempty"`
}

type RoutePlanPolicy struct {
	PreferRouteTypes []string `json:"prefer_route_types,omitempty"`
	AvoidRouteTypes  []string `json:"avoid_route_types,omitempty"`
	PreferRoles      []string `json:"prefer_roles,omitempty"`
	AvoidRoles       []string `json:"avoid_roles,omitempty"`
}

type ProtocolRoutePlanPolicy struct {
	ProtocolPCID string `json:"protocol_pcid"`
	RoutePlanPolicy
}

type ProtocolRoleRoutePlanPolicy struct {
	ProtocolPCID string `json:"protocol_pcid"`
	Role         string `json:"role"`
	RoutePlanPolicy
}

type TrustPolicy struct {
	ClaimPolicies                 []ClaimTrustPolicy            `json:"claim_policies"`
	RoutePlanPolicy               RoutePlanPolicy               `json:"route_plan_policy,omitempty"`
	ProtocolRoutePlanPolicies     []ProtocolRoutePlanPolicy     `json:"protocol_route_plan_policies,omitempty"`
	ProtocolRoleRoutePlanPolicies []ProtocolRoleRoutePlanPolicy `json:"protocol_role_route_plan_policies,omitempty"`
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

func (store *PolicyStore) RoutePlanPolicy() RoutePlanPolicy {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return cloneRoutePlanPolicy(store.policy.RoutePlanPolicy)
}

func (store *PolicyStore) ProtocolRoutePlanPolicies() []ProtocolRoutePlanPolicy {
	store.mu.RLock()
	defer store.mu.RUnlock()
	out := make([]ProtocolRoutePlanPolicy, 0, len(store.policy.ProtocolRoutePlanPolicies))
	for _, policy := range store.policy.ProtocolRoutePlanPolicies {
		out = append(out, cloneProtocolRoutePlanPolicy(policy))
	}
	return out
}

func (store *PolicyStore) ProtocolRoleRoutePlanPolicies() []ProtocolRoleRoutePlanPolicy {
	store.mu.RLock()
	defer store.mu.RUnlock()
	out := make([]ProtocolRoleRoutePlanPolicy, 0, len(store.policy.ProtocolRoleRoutePlanPolicies))
	for _, policy := range store.policy.ProtocolRoleRoutePlanPolicies {
		out = append(out, cloneProtocolRoleRoutePlanPolicy(policy))
	}
	return out
}

// Intent: Let route planning keep one global default policy while still
// adapting route preferences for specific input protocols when local operators
// know a given pCID needs different routing behavior. Source: DI-posek
func (store *PolicyStore) EffectiveRoutePlanPolicy(protocolPCID string) RoutePlanPolicy {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return effectiveRoutePlanPolicyLocked(store.policy.RoutePlanPolicy, store.policy.ProtocolRoutePlanPolicies, protocolPCID)
}

// Intent: Let operators steer candidate selection within one protocol by
// layering exact-role overrides on top of the global and per-protocol route
// policy instead of forcing all roles for that protocol to share one policy.
// Source: DI-rivuk
func (store *PolicyStore) EffectiveRoutePlanPolicyForRole(protocolPCID string, role string) RoutePlanPolicy {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return effectiveRoutePlanPolicyForRoleLocked(
		store.policy.RoutePlanPolicy,
		store.policy.ProtocolRoutePlanPolicies,
		store.policy.ProtocolRoleRoutePlanPolicies,
		protocolPCID,
		role,
	)
}

func (store *PolicyStore) SetRoutePlanPolicy(policy RoutePlanPolicy) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	policy = normalizeRoutePlanPolicy(policy)
	if err := validateRoutePlanPolicy(policy); err != nil {
		return err
	}
	store.policy.RoutePlanPolicy = policy
	return store.persistLocked()
}

func (store *PolicyStore) SetProtocolRoutePlanPolicy(protocolPCID string, policy RoutePlanPolicy) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	protocolPCID = strings.TrimSpace(protocolPCID)
	if protocolPCID == "" {
		return errors.New("protocol_pcid is required")
	}
	policy = normalizeRoutePlanPolicy(policy)
	if err := validateRoutePlanPolicy(policy); err != nil {
		return err
	}
	replaced := false
	for index := range store.policy.ProtocolRoutePlanPolicies {
		if store.policy.ProtocolRoutePlanPolicies[index].ProtocolPCID == protocolPCID {
			store.policy.ProtocolRoutePlanPolicies[index] = ProtocolRoutePlanPolicy{
				ProtocolPCID:    protocolPCID,
				RoutePlanPolicy: policy,
			}
			replaced = true
			break
		}
	}
	if !replaced {
		store.policy.ProtocolRoutePlanPolicies = append(store.policy.ProtocolRoutePlanPolicies, ProtocolRoutePlanPolicy{
			ProtocolPCID:    protocolPCID,
			RoutePlanPolicy: policy,
		})
	}
	sortProtocolRoutePlanPolicies(store.policy.ProtocolRoutePlanPolicies)
	return store.persistLocked()
}

func (store *PolicyStore) SetProtocolRoleRoutePlanPolicy(protocolPCID string, role string, policy RoutePlanPolicy) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	protocolPCID = strings.TrimSpace(protocolPCID)
	role = strings.TrimSpace(role)
	if protocolPCID == "" {
		return errors.New("protocol_pcid is required")
	}
	if role == "" {
		return errors.New("role is required")
	}
	policy = normalizeRoutePlanPolicy(policy)
	if err := validateRoutePlanPolicy(policy); err != nil {
		return err
	}
	replaced := false
	for index := range store.policy.ProtocolRoleRoutePlanPolicies {
		current := store.policy.ProtocolRoleRoutePlanPolicies[index]
		if current.ProtocolPCID == protocolPCID && current.Role == role {
			store.policy.ProtocolRoleRoutePlanPolicies[index] = ProtocolRoleRoutePlanPolicy{
				ProtocolPCID:    protocolPCID,
				Role:            role,
				RoutePlanPolicy: policy,
			}
			replaced = true
			break
		}
	}
	if !replaced {
		store.policy.ProtocolRoleRoutePlanPolicies = append(store.policy.ProtocolRoleRoutePlanPolicies, ProtocolRoleRoutePlanPolicy{
			ProtocolPCID:    protocolPCID,
			Role:            role,
			RoutePlanPolicy: policy,
		})
	}
	sortProtocolRoleRoutePlanPolicies(store.policy.ProtocolRoleRoutePlanPolicies)
	return store.persistLocked()
}

func (store *PolicyStore) RemoveProtocolRoutePlanPolicy(protocolPCID string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	protocolPCID = strings.TrimSpace(protocolPCID)
	filtered := store.policy.ProtocolRoutePlanPolicies[:0]
	found := false
	for _, policy := range store.policy.ProtocolRoutePlanPolicies {
		if policy.ProtocolPCID == protocolPCID {
			found = true
			continue
		}
		filtered = append(filtered, policy)
	}
	if !found {
		return fmt.Errorf("unknown protocol route plan policy: %s", protocolPCID)
	}
	store.policy.ProtocolRoutePlanPolicies = filtered
	return store.persistLocked()
}

func (store *PolicyStore) RemoveProtocolRoleRoutePlanPolicy(protocolPCID string, role string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	protocolPCID = strings.TrimSpace(protocolPCID)
	role = strings.TrimSpace(role)
	filtered := store.policy.ProtocolRoleRoutePlanPolicies[:0]
	found := false
	for _, policy := range store.policy.ProtocolRoleRoutePlanPolicies {
		if policy.ProtocolPCID == protocolPCID && policy.Role == role {
			found = true
			continue
		}
		filtered = append(filtered, policy)
	}
	if !found {
		return fmt.Errorf("unknown protocol role route plan policy: %s %s", protocolPCID, role)
	}
	store.policy.ProtocolRoleRoutePlanPolicies = filtered
	return store.persistLocked()
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
// quorums, weights, and federation spread against runtime-owned policy instead
// of treating any valid countersignature as automatically sufficient. Source: DI-rumek
func (store *PolicyStore) VerifyClaimPolicies(batch Batch, knownPeers []AllowedPeer) error {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for index, claim := range batch.ImplementationClaims {
		policy, ok := store.matchClaimPolicyLocked(claim)
		if !ok || (policy.MinAttesters == 0 && policy.MinTrustWeight == 0 && policy.MinFederations == 0) {
			continue
		}
		matched := map[string]struct{}{}
		federations := map[string]struct{}{}
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
			if len(policy.AllowedFederations) > 0 && !slices.Contains(policy.AllowedFederations, peer.Federation) {
				continue
			}
			if _, exists := matched[attestation.SignerPeerID]; exists {
				continue
			}
			matched[attestation.SignerPeerID] = struct{}{}
			federations[peer.Federation] = struct{}{}
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
		if len(federations) < policy.MinFederations {
			return fmt.Errorf(
				"claim attestation federation spread failed for %s role %s: need %d, got %d",
				claim.ProtocolPCID,
				claim.Role,
				policy.MinFederations,
				len(federations),
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
	store.policy.RoutePlanPolicy = normalizeRoutePlanPolicy(store.policy.RoutePlanPolicy)
	if err := validateRoutePlanPolicy(store.policy.RoutePlanPolicy); err != nil {
		return err
	}
	routeSeen := map[string]struct{}{}
	for index, policy := range store.policy.ProtocolRoutePlanPolicies {
		policy = normalizeProtocolRoutePlanPolicy(policy)
		if err := validateProtocolRoutePlanPolicy(policy); err != nil {
			return err
		}
		if _, exists := routeSeen[policy.ProtocolPCID]; exists {
			return fmt.Errorf("duplicate protocol route plan policy: %s", policy.ProtocolPCID)
		}
		routeSeen[policy.ProtocolPCID] = struct{}{}
		store.policy.ProtocolRoutePlanPolicies[index] = policy
	}
	roleSeen := map[string]struct{}{}
	for index, policy := range store.policy.ProtocolRoleRoutePlanPolicies {
		policy = normalizeProtocolRoleRoutePlanPolicy(policy)
		if err := validateProtocolRoleRoutePlanPolicy(policy); err != nil {
			return err
		}
		key := policy.ProtocolPCID + "\x00" + policy.Role
		if _, exists := roleSeen[key]; exists {
			return fmt.Errorf("duplicate protocol role route plan policy: %s %s", policy.ProtocolPCID, policy.Role)
		}
		roleSeen[key] = struct{}{}
		store.policy.ProtocolRoleRoutePlanPolicies[index] = policy
	}
	sortClaimPolicies(store.policy.ClaimPolicies)
	sortProtocolRoutePlanPolicies(store.policy.ProtocolRoutePlanPolicies)
	sortProtocolRoleRoutePlanPolicies(store.policy.ProtocolRoleRoutePlanPolicies)
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
	if policy.MinFederations < 0 {
		return errors.New("min_federations must be zero or greater")
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
	seenFederations := map[string]struct{}{}
	for _, federation := range policy.AllowedFederations {
		if strings.TrimSpace(federation) == "" {
			return errors.New("allowed_federations cannot contain blanks")
		}
		if _, exists := seenFederations[federation]; exists {
			return fmt.Errorf("duplicate allowed_federation: %s", federation)
		}
		seenFederations[federation] = struct{}{}
	}
	if len(policy.AllowedAttesters) > 0 && policy.MinAttesters > len(policy.AllowedAttesters) {
		return errors.New("min_attesters cannot exceed allowed_attesters length")
	}
	if len(policy.AllowedFederations) > 0 && policy.MinFederations > len(policy.AllowedFederations) {
		return errors.New("min_federations cannot exceed allowed_federations length")
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
	for index := range policy.AllowedFederations {
		policy.AllowedFederations[index] = strings.TrimSpace(policy.AllowedFederations[index])
	}
	slices.Sort(policy.AllowedFederations)
	policy.AllowedFederations = slices.Compact(policy.AllowedFederations)
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

func validateRoutePlanPolicy(policy RoutePlanPolicy) error {
	if err := validateDistinctNonBlank(policy.PreferRouteTypes, "prefer_route_types"); err != nil {
		return err
	}
	if err := validateDistinctNonBlank(policy.AvoidRouteTypes, "avoid_route_types"); err != nil {
		return err
	}
	if err := validateDistinctNonBlank(policy.PreferRoles, "prefer_roles"); err != nil {
		return err
	}
	if err := validateDistinctNonBlank(policy.AvoidRoles, "avoid_roles"); err != nil {
		return err
	}
	return nil
}

func normalizeRoutePlanPolicy(policy RoutePlanPolicy) RoutePlanPolicy {
	policy.PreferRouteTypes = normalizeStringList(policy.PreferRouteTypes)
	policy.AvoidRouteTypes = normalizeStringList(policy.AvoidRouteTypes)
	policy.PreferRoles = normalizeStringList(policy.PreferRoles)
	policy.AvoidRoles = normalizeStringList(policy.AvoidRoles)
	return policy
}

func cloneRoutePlanPolicy(policy RoutePlanPolicy) RoutePlanPolicy {
	return RoutePlanPolicy{
		PreferRouteTypes: append([]string{}, policy.PreferRouteTypes...),
		AvoidRouteTypes:  append([]string{}, policy.AvoidRouteTypes...),
		PreferRoles:      append([]string{}, policy.PreferRoles...),
		AvoidRoles:       append([]string{}, policy.AvoidRoles...),
	}
}

func cloneProtocolRoutePlanPolicy(policy ProtocolRoutePlanPolicy) ProtocolRoutePlanPolicy {
	return ProtocolRoutePlanPolicy{
		ProtocolPCID:    policy.ProtocolPCID,
		RoutePlanPolicy: cloneRoutePlanPolicy(policy.RoutePlanPolicy),
	}
}

func cloneProtocolRoleRoutePlanPolicy(policy ProtocolRoleRoutePlanPolicy) ProtocolRoleRoutePlanPolicy {
	return ProtocolRoleRoutePlanPolicy{
		ProtocolPCID:    policy.ProtocolPCID,
		Role:            policy.Role,
		RoutePlanPolicy: cloneRoutePlanPolicy(policy.RoutePlanPolicy),
	}
}

func normalizeProtocolRoutePlanPolicy(policy ProtocolRoutePlanPolicy) ProtocolRoutePlanPolicy {
	policy.ProtocolPCID = strings.TrimSpace(policy.ProtocolPCID)
	policy.RoutePlanPolicy = normalizeRoutePlanPolicy(policy.RoutePlanPolicy)
	return policy
}

func normalizeProtocolRoleRoutePlanPolicy(policy ProtocolRoleRoutePlanPolicy) ProtocolRoleRoutePlanPolicy {
	policy.ProtocolPCID = strings.TrimSpace(policy.ProtocolPCID)
	policy.Role = strings.TrimSpace(policy.Role)
	policy.RoutePlanPolicy = normalizeRoutePlanPolicy(policy.RoutePlanPolicy)
	return policy
}

func validateProtocolRoutePlanPolicy(policy ProtocolRoutePlanPolicy) error {
	if policy.ProtocolPCID == "" {
		return errors.New("protocol_pcid is required")
	}
	return validateRoutePlanPolicy(policy.RoutePlanPolicy)
}

func validateProtocolRoleRoutePlanPolicy(policy ProtocolRoleRoutePlanPolicy) error {
	if policy.ProtocolPCID == "" {
		return errors.New("protocol_pcid is required")
	}
	if policy.Role == "" {
		return errors.New("role is required")
	}
	return validateRoutePlanPolicy(policy.RoutePlanPolicy)
}

func effectiveRoutePlanPolicyLocked(global RoutePlanPolicy, overrides []ProtocolRoutePlanPolicy, protocolPCID string) RoutePlanPolicy {
	effective := cloneRoutePlanPolicy(global)
	protocolPCID = strings.TrimSpace(protocolPCID)
	if protocolPCID == "" {
		return effective
	}
	for _, override := range overrides {
		if override.ProtocolPCID != protocolPCID {
			continue
		}
		if len(override.PreferRouteTypes) > 0 {
			effective.PreferRouteTypes = append([]string{}, override.PreferRouteTypes...)
		}
		if len(override.AvoidRouteTypes) > 0 {
			effective.AvoidRouteTypes = append([]string{}, override.AvoidRouteTypes...)
		}
		if len(override.PreferRoles) > 0 {
			effective.PreferRoles = append([]string{}, override.PreferRoles...)
		}
		if len(override.AvoidRoles) > 0 {
			effective.AvoidRoles = append([]string{}, override.AvoidRoles...)
		}
		break
	}
	return effective
}

func effectiveRoutePlanPolicyForRoleLocked(
	global RoutePlanPolicy,
	protocolOverrides []ProtocolRoutePlanPolicy,
	roleOverrides []ProtocolRoleRoutePlanPolicy,
	protocolPCID string,
	role string,
) RoutePlanPolicy {
	effective := effectiveRoutePlanPolicyLocked(global, protocolOverrides, protocolPCID)
	role = strings.TrimSpace(role)
	if role == "" {
		return effective
	}
	for _, override := range roleOverrides {
		if override.ProtocolPCID != protocolPCID || override.Role != role {
			continue
		}
		if len(override.PreferRouteTypes) > 0 {
			effective.PreferRouteTypes = append([]string{}, override.PreferRouteTypes...)
		}
		if len(override.AvoidRouteTypes) > 0 {
			effective.AvoidRouteTypes = append([]string{}, override.AvoidRouteTypes...)
		}
		if len(override.PreferRoles) > 0 {
			effective.PreferRoles = append([]string{}, override.PreferRoles...)
		}
		if len(override.AvoidRoles) > 0 {
			effective.AvoidRoles = append([]string{}, override.AvoidRoles...)
		}
		break
	}
	return effective
}

func sortProtocolRoutePlanPolicies(policies []ProtocolRoutePlanPolicy) {
	slices.SortFunc(policies, func(left, right ProtocolRoutePlanPolicy) int {
		return strings.Compare(left.ProtocolPCID, right.ProtocolPCID)
	})
}

func sortProtocolRoleRoutePlanPolicies(policies []ProtocolRoleRoutePlanPolicy) {
	slices.SortFunc(policies, func(left, right ProtocolRoleRoutePlanPolicy) int {
		if left.ProtocolPCID != right.ProtocolPCID {
			return strings.Compare(left.ProtocolPCID, right.ProtocolPCID)
		}
		return strings.Compare(left.Role, right.Role)
	})
}

func normalizeStringList(values []string) []string {
	out := append([]string{}, values...)
	for index := range out {
		out[index] = strings.TrimSpace(out[index])
	}
	slices.Sort(out)
	return slices.Compact(out)
}

func validateDistinctNonBlank(values []string, field string) error {
	seen := map[string]struct{}{}
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s cannot contain blanks", field)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("duplicate %s entry: %s", field, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}
