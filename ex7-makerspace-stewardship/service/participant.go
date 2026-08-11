package service

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

const (
	rootHistoryPCID = "bafkreia7cn4srmmkxbwxk2hoezedjvuyokhypcsddjd4evx56lhtmsq3nm"
	deviceAuthPCID  = "bafkreifmbhgjwfmwbemkf4ogsg3gvuavjhttkitzf3muie3dhv5tdn4hq4"
	revocationPCID  = "bafkreify46v4jp3dvz3szem6vrukscplr3kshku7fhrzn4mc26scnnvqoi"
	recoveryPCID    = "bafkreicjdzlriq3nfasza5nmnflpycche63wn2n5kauq66jivlwhhomesy"
	peerCardPCID    = "bafkreicstci6idwm6d6dbt52ppqyjcapibskz27qmnfuyntg6zck72fa24"
	carriagePCID    = "bafkreihrlojt47erjc6uawkm47s7evppp23tk3ljlkl347ten4v3kb624i"
)

type rootHistoryPayload struct {
	RootKey              string   `json:"root_key"`
	PreviousRootRecordID string   `json:"previous_root_record_id,omitempty"`
	HistoryNote          string   `json:"history_note"`
	RecoverySet          []string `json:"recovery_set"`
}

type deviceAuthorizationPayload struct {
	RootRecordID string `json:"root_record_id"`
	DeviceKey    string `json:"device_key"`
	DeviceLabel  string `json:"device_label"`
	NotBefore    string `json:"not_before"`
	NotAfter     string `json:"not_after,omitempty"`
}

type revocationPayload struct {
	SubjectKeyID string `json:"subject_key_id"`
	SubjectKind  string `json:"subject_kind"`
	EffectiveAt  string `json:"effective_at"`
	Reason       string `json:"reason"`
}

type recoveryPayload struct {
	RootRecordID       string   `json:"root_record_id"`
	RecoveryID         string   `json:"recovery_id"`
	ReplacementRootKey string   `json:"replacement_root_key"`
	RecoverySet        []string `json:"recovery_set"`
}

type peerCardPayload struct {
	RootRecordID          string   `json:"root_record_id"`
	ActiveDeviceRecordIDs []string `json:"active_device_record_ids"`
	ContactHints          []string `json:"contact_hints"`
}

type carriagePayload struct {
	SenderCardRecordID string   `json:"sender_card_record_id"`
	Cursor             string   `json:"cursor"`
	Records            []string `json:"records"`
}

type rootState struct {
	key         ed25519.PublicKey
	lineage     string
	active      bool
	recoverySet []string
}

type recoveryState struct {
	replacement ed25519.PublicKey
	votes       map[string]struct{}
}

type deviceState struct {
	key       ed25519.PublicKey
	lineage   string
	notBefore time.Time
	notAfter  time.Time
}

// ParticipantHistory verifies the signed key history that makes a participant
// record author meaningful. Intent: a label, account, or recognition entry
// cannot replace root-linked author evidence. Source: DI-kasaz.
type ParticipantHistory struct {
	roots               map[string]rootState
	devices             map[string]deviceState
	deviceRecords       map[string]deviceState
	revocations         map[string]time.Time
	recoveries          map[string]recoveryState
	completedRecoveries map[string]recoveryState
	peerCards           map[string]struct{}
}

func NewParticipantHistory() *ParticipantHistory {
	return &ParticipantHistory{
		roots:               make(map[string]rootState),
		devices:             make(map[string]deviceState),
		deviceRecords:       make(map[string]deviceState),
		revocations:         make(map[string]time.Time),
		recoveries:          make(map[string]recoveryState),
		completedRecoveries: make(map[string]recoveryState),
		peerCards:           make(map[string]struct{}),
	}
}

func (h *ParticipantHistory) Clone() *ParticipantHistory {
	clone := NewParticipantHistory()
	for id, root := range h.roots {
		root.key = append(ed25519.PublicKey(nil), root.key...)
		root.recoverySet = append([]string(nil), root.recoverySet...)
		clone.roots[id] = root
	}
	for id, device := range h.devices {
		device.key = append(ed25519.PublicKey(nil), device.key...)
		clone.devices[id] = device
	}
	for id, device := range h.deviceRecords {
		device.key = append(ed25519.PublicKey(nil), device.key...)
		clone.deviceRecords[id] = device
	}
	for id, at := range h.revocations {
		clone.revocations[id] = at
	}
	for id, recovery := range h.recoveries {
		recovery.replacement = append(ed25519.PublicKey(nil), recovery.replacement...)
		recovery.votes = cloneVotes(recovery.votes)
		clone.recoveries[id] = recovery
	}
	for id, recovery := range h.completedRecoveries {
		recovery.replacement = append(ed25519.PublicKey(nil), recovery.replacement...)
		recovery.votes = cloneVotes(recovery.votes)
		clone.completedRecoveries[id] = recovery
	}
	for id := range h.peerCards {
		clone.peerCards[id] = struct{}{}
	}
	return clone
}

func (h *ParticipantHistory) Apply(record Record) error {
	switch record.Protocol {
	case rootHistoryPCID:
		return h.applyRoot(record)
	case deviceAuthPCID:
		return h.applyDeviceAuthorization(record)
	case revocationPCID:
		return h.applyRevocation(record)
	case recoveryPCID:
		return h.applyRecovery(record)
	case peerCardPCID:
		return h.applyPeerCard(record)
	default:
		return nil
	}
}

func (h *ParticipantHistory) applyRoot(record Record) error {
	if _, exists := h.roots[record.ID]; exists {
		return nil
	}
	var payload rootHistoryPayload
	if err := decodePayload(record.Payload, &payload); err != nil {
		return fmt.Errorf("decode root-history payload: %w", err)
	}
	rootKey, err := decodePublicKey(payload.RootKey)
	if err != nil {
		return fmt.Errorf("root_key: %w", err)
	}
	if payload.HistoryNote == "" || len(payload.RecoverySet) != 3 || !distinctPublicKeys(payload.RecoverySet) {
		return errors.New("invalid root-history payload")
	}
	for _, encoded := range payload.RecoverySet {
		if _, err := decodePublicKey(encoded); err != nil {
			return fmt.Errorf("recovery_set: %w", err)
		}
	}
	lineage := record.ID
	if payload.PreviousRootRecordID == "" {
		if !bytes.Equal(record.PublicKey, rootKey) {
			return errors.New("root establishment must be self-signed")
		}
	} else {
		previous, ok := h.roots[payload.PreviousRootRecordID]
		if !ok || !previous.active {
			return errors.New("root continuation has no active predecessor")
		}
		recovered, recoveryComplete := h.completedRecoveries[payload.PreviousRootRecordID]
		ordinaryContinuation := bytes.Equal(record.PublicKey, previous.key) && !h.revokedAt(record.KeyID, record.CreatedAt)
		recoveryContinuation := recoveryComplete && bytes.Equal(record.PublicKey, rootKey) && bytes.Equal(rootKey, recovered.replacement)
		if !ordinaryContinuation && !recoveryContinuation {
			return errors.New("root continuation is not signed by predecessor or completed recovery key")
		}
		previous.active = false
		h.roots[payload.PreviousRootRecordID] = previous
		lineage = previous.lineage
	}
	h.roots[record.ID] = rootState{key: rootKey, lineage: lineage, active: true, recoverySet: append([]string(nil), payload.RecoverySet...)}
	return nil
}

func (h *ParticipantHistory) applyDeviceAuthorization(record Record) error {
	var payload deviceAuthorizationPayload
	if err := decodePayload(record.Payload, &payload); err != nil {
		return fmt.Errorf("decode device-authorization payload: %w", err)
	}
	root, ok := h.roots[payload.RootRecordID]
	if !ok || !root.active || !bytes.Equal(record.PublicKey, root.key) || h.revokedAt(record.KeyID, record.CreatedAt) {
		return errors.New("device authorization is not signed by an active root")
	}
	key, err := decodePublicKey(payload.DeviceKey)
	if err != nil {
		return fmt.Errorf("device_key: %w", err)
	}
	notBefore, err := time.Parse(time.RFC3339, payload.NotBefore)
	if err != nil || payload.DeviceLabel == "" {
		return errors.New("invalid device authorization")
	}
	var notAfter time.Time
	if payload.NotAfter != "" {
		notAfter, err = time.Parse(time.RFC3339, payload.NotAfter)
		if err != nil || !notAfter.After(notBefore) {
			return errors.New("invalid device authorization expiry")
		}
	}
	device := deviceState{key: key, lineage: root.lineage, notBefore: notBefore, notAfter: notAfter}
	h.devices[keyID(key)] = device
	h.deviceRecords[record.ID] = device
	return nil
}

func (h *ParticipantHistory) applyRevocation(record Record) error {
	var payload revocationPayload
	if err := decodePayload(record.Payload, &payload); err != nil {
		return fmt.Errorf("decode key-revocation payload: %w", err)
	}
	if payload.SubjectKeyID == "" || (payload.SubjectKind != "root" && payload.SubjectKind != "device") || payload.Reason == "" {
		return errors.New("invalid key revocation")
	}
	effectiveAt, err := time.Parse(time.RFC3339, payload.EffectiveAt)
	if err != nil {
		return fmt.Errorf("revocation effective_at: %w", err)
	}
	signerRoot, ok := h.activeRoot(record.PublicKey, record.CreatedAt)
	if !ok {
		return errors.New("revocation is not signed by an active root")
	}
	if payload.SubjectKind == "root" {
		subject, ok := h.rootByKeyID(payload.SubjectKeyID)
		if !ok || subject.lineage != signerRoot.lineage {
			return errors.New("revocation root subject is outside signer history")
		}
	}
	if payload.SubjectKind == "device" {
		subject, ok := h.devices[payload.SubjectKeyID]
		if !ok || subject.lineage != signerRoot.lineage {
			return errors.New("revocation device subject is outside signer history")
		}
	}
	if current, exists := h.revocations[payload.SubjectKeyID]; !exists || effectiveAt.Before(current) {
		h.revocations[payload.SubjectKeyID] = effectiveAt
	}
	return nil
}

func (h *ParticipantHistory) applyRecovery(record Record) error {
	var payload recoveryPayload
	if err := decodePayload(record.Payload, &payload); err != nil {
		return fmt.Errorf("decode threshold-recovery payload: %w", err)
	}
	root, ok := h.roots[payload.RootRecordID]
	if !ok || payload.RecoveryID == "" || len(payload.RecoverySet) != 3 || !distinctPublicKeys(payload.RecoverySet) || !sameStrings(payload.RecoverySet, root.recoverySet) {
		return errors.New("invalid threshold-recovery payload")
	}
	replacement, err := decodePublicKey(payload.ReplacementRootKey)
	if err != nil {
		return fmt.Errorf("replacement_root_key: %w", err)
	}
	if !containsPublicKey(payload.RecoverySet, record.PublicKey) {
		return errors.New("recovery witness is not declared by root history")
	}
	stateID := payload.RootRecordID + "\x00" + payload.RecoveryID
	state, exists := h.recoveries[stateID]
	if !exists {
		state = recoveryState{replacement: replacement, votes: make(map[string]struct{})}
	} else if len(state.replacement) == 0 || !bytes.Equal(state.replacement, replacement) {
		// Intent: Preserve conflicting recovery promises as evidence without
		// allowing either replacement to activate. Source: DI-sisad.
		state.replacement = nil
		state.votes = nil
		h.recoveries[stateID] = state
		return nil
	}
	state.votes[record.KeyID] = struct{}{}
	h.recoveries[stateID] = state
	if len(state.votes) >= 2 {
		if completed, exists := h.completedRecoveries[payload.RootRecordID]; exists && !bytes.Equal(completed.replacement, replacement) {
			// Intent: Keep both durable witness sets while refusing to choose a
			// replacement root from conflicting completed recoveries. Source: DI-sisad.
			completed.replacement = nil
			h.completedRecoveries[payload.RootRecordID] = completed
			return nil
		}
		h.completedRecoveries[payload.RootRecordID] = state
	}
	return nil
}

func (h *ParticipantHistory) applyPeerCard(record Record) error {
	var payload peerCardPayload
	if err := decodePayload(record.Payload, &payload); err != nil {
		return fmt.Errorf("decode peer-card payload: %w", err)
	}
	root, ok := h.roots[payload.RootRecordID]
	if !ok {
		return errors.New("peer card references unknown root history")
	}
	if !h.Authorizes(record) {
		return errors.New("peer card signer is not an active root or device")
	}
	if signerRoot, ok := h.activeRoot(record.PublicKey, record.CreatedAt); ok {
		if signerRoot.lineage != root.lineage {
			return errors.New("peer card root is outside signer history")
		}
	} else {
		device, ok := h.devices[record.KeyID]
		if !ok || device.lineage != root.lineage {
			return errors.New("peer card device is outside signer history")
		}
	}
	for _, id := range payload.ActiveDeviceRecordIDs {
		device, ok := h.deviceRecords[id]
		if !ok || device.lineage != root.lineage {
			return errors.New("peer card references device outside root history")
		}
	}
	for _, hint := range payload.ContactHints {
		if hint == "" {
			return errors.New("peer card has empty contact hint")
		}
	}
	h.peerCards[record.ID] = struct{}{}
	return nil
}

func (h *ParticipantHistory) ValidateCarriage(record Record) ([][]byte, error) {
	var payload carriagePayload
	if err := decodePayload(record.Payload, &payload); err != nil {
		return nil, fmt.Errorf("decode exact-record-carriage payload: %w", err)
	}
	if payload.SenderCardRecordID == "" || payload.Cursor == "" || len(payload.Records) == 0 {
		return nil, errors.New("invalid exact-record-carriage payload")
	}
	if _, ok := h.peerCards[payload.SenderCardRecordID]; !ok {
		return nil, errors.New("carriage references unknown peer card")
	}
	records := make([][]byte, 0, len(payload.Records))
	for _, encoded := range payload.Records {
		raw, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, errors.New("carriage record is not base64")
		}
		if _, err := ParseRecord(raw); err != nil {
			return nil, fmt.Errorf("carriage enclosed record: %w", err)
		}
		records = append(records, raw)
	}
	return records, nil
}

func (h *ParticipantHistory) Authorizes(record Record) bool {
	if h.revokedAt(record.KeyID, record.CreatedAt) {
		return false
	}
	if h.isActiveRoot(record.PublicKey, record.CreatedAt) {
		return true
	}
	device, ok := h.devices[record.KeyID]
	if !ok || !bytes.Equal(device.key, record.PublicKey) {
		return false
	}
	createdAt, err := time.Parse(time.RFC3339, record.CreatedAt)
	if err != nil || createdAt.Before(device.notBefore) {
		return false
	}
	return device.notAfter.IsZero() || createdAt.Before(device.notAfter)
}

func (h *ParticipantHistory) isActiveRoot(publicKey ed25519.PublicKey, createdAt string) bool {
	_, ok := h.activeRoot(publicKey, createdAt)
	return ok
}

func (h *ParticipantHistory) activeRoot(publicKey ed25519.PublicKey, createdAt string) (rootState, bool) {
	for _, root := range h.roots {
		if root.active && bytes.Equal(root.key, publicKey) && !h.revokedAt(keyID(publicKey), createdAt) {
			return root, true
		}
	}
	return rootState{}, false
}

func (h *ParticipantHistory) rootByKeyID(id string) (rootState, bool) {
	for _, root := range h.roots {
		if keyID(root.key) == id {
			return root, true
		}
	}
	return rootState{}, false
}

func (h *ParticipantHistory) revokedAt(key, createdAt string) bool {
	effectiveAt, exists := h.revocations[key]
	if !exists {
		return false
	}
	created, err := time.Parse(time.RFC3339, createdAt)
	return err != nil || !created.Before(effectiveAt)
}

func decodePublicKey(encoded string) (ed25519.PublicKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(decoded) != ed25519.PublicKeySize {
		return nil, errors.New("must be a base64 Ed25519 public key")
	}
	return ed25519.PublicKey(decoded), nil
}

func distinctPublicKeys(encoded []string) bool {
	seen := make(map[string]struct{}, len(encoded))
	for _, value := range encoded {
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}

func containsPublicKey(encoded []string, public ed25519.PublicKey) bool {
	for _, value := range encoded {
		key, err := decodePublicKey(value)
		if err == nil && bytes.Equal(key, public) {
			return true
		}
	}
	return false
}

func sameStrings(left, right []string) bool {
	return len(left) == len(right) && func() bool {
		for i := range left {
			if left[i] != right[i] {
				return false
			}
		}
		return true
	}()
}

func cloneVotes(votes map[string]struct{}) map[string]struct{} {
	clone := make(map[string]struct{}, len(votes))
	for id := range votes {
		clone[id] = struct{}{}
	}
	return clone
}

// ParticipantSigner creates exact records only with an active history-linked
// root or device key. Source: DI-kasaz.
type ParticipantSigner struct {
	private ed25519.PrivateKey
	label   string
	history *ParticipantHistory
}

func NewParticipantSigner(private ed25519.PrivateKey, label string, history *ParticipantHistory) (*ParticipantSigner, error) {
	if len(private) != ed25519.PrivateKeySize || label == "" || history == nil {
		return nil, errors.New("participant signer requires private key, label, and history")
	}
	return &ParticipantSigner{private: private, label: label, history: history}, nil
}

func (s *ParticipantSigner) Sign(record Record) (Record, []byte, error) {
	public := s.private.Public().(ed25519.PublicKey)
	record.Signer = s.label
	record.PublicKey = public
	record.KeyID = keyID(public)
	if record.CreatedAt == "" {
		return Record{}, nil, errors.New("record created_at is required")
	}
	if record.Protocol != rootHistoryPCID && !s.history.Authorizes(record) {
		return Record{}, nil, errors.New("signer is not an active participant root or device")
	}
	return record.Sign(s.private)
}
