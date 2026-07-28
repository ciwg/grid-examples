package grid

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const RelayBatchFormat = "moks-relay-batch-v1"

type ImplementationClaim struct {
	PackageID      string `json:"package_id"`
	ProtocolPCID   string `json:"protocol_pcid"`
	Role           string `json:"role"`
	Summary        string `json:"summary,omitempty"`
	PackageVersion string `json:"package_version"`
}

type RouteRegistration struct {
	PackageID      string   `json:"package_id"`
	PackageVersion string   `json:"package_version"`
	ProtocolPCID   string   `json:"protocol_pcid"`
	Role           string   `json:"role"`
	Summary        string   `json:"summary,omitempty"`
	Families       []string `json:"families,omitempty"`
}

type RecordProof struct {
	Digest string `json:"digest"`
}

type ClaimProof struct {
	SignerPeerID string `json:"signer_peer_id"`
	PublicKey    string `json:"public_key"`
	Signature    string `json:"signature"`
}

type ClaimAttestation struct {
	ClaimIndex   int    `json:"claim_index"`
	SignerPeerID string `json:"signer_peer_id"`
	PublicKey    string `json:"public_key"`
	Signature    string `json:"signature"`
}

type RecordSignature struct {
	SignerPeerID string `json:"signer_peer_id"`
	PublicKey    string `json:"public_key"`
	Signature    string `json:"signature"`
}

type Batch struct {
	Format               string                `json:"format"`
	Implementation       string                `json:"implementation"`
	ExportedAt           string                `json:"exported_at"`
	ImplementationClaims []ImplementationClaim `json:"implementation_claims,omitempty"`
	Routes               []RouteRegistration   `json:"routes,omitempty"`
	ClaimProofs          []ClaimProof          `json:"claim_proofs,omitempty"`
	ClaimAttestations    []ClaimAttestation    `json:"claim_attestations,omitempty"`
	Records              []json.RawMessage     `json:"records"`
	RecordProofs         []RecordProof         `json:"record_proofs,omitempty"`
	RecordSignatures     []RecordSignature     `json:"record_signatures,omitempty"`
	Signature            string                `json:"signature,omitempty"`
}

func (batch Batch) Validate() error {
	if batch.Format != RelayBatchFormat {
		return fmt.Errorf("unsupported batch format: %s", batch.Format)
	}
	if strings.TrimSpace(batch.Implementation) == "" {
		return errors.New("batch implementation is required")
	}
	if strings.TrimSpace(batch.ExportedAt) == "" {
		return errors.New("batch exported_at is required")
	}
	if _, err := time.Parse(time.RFC3339, batch.ExportedAt); err != nil {
		return fmt.Errorf("batch exported_at must be RFC3339: %w", err)
	}
	seenClaims := map[string]struct{}{}
	for _, claim := range batch.ImplementationClaims {
		if strings.TrimSpace(claim.PackageID) == "" {
			return errors.New("claim package_id is required")
		}
		if strings.TrimSpace(claim.PackageVersion) == "" {
			return errors.New("claim package_version is required")
		}
		if strings.TrimSpace(claim.ProtocolPCID) == "" {
			return errors.New("claim protocol_pcid is required")
		}
		if strings.TrimSpace(claim.Role) == "" {
			return errors.New("claim role is required")
		}
		key := claim.PackageID + "\x00" + claim.PackageVersion + "\x00" + claim.ProtocolPCID + "\x00" + claim.Role
		if _, exists := seenClaims[key]; exists {
			return fmt.Errorf("duplicate implementation claim: %s", key)
		}
		seenClaims[key] = struct{}{}
	}
	seenRoutes := map[string]struct{}{}
	// Intent: Keep exported route metadata derivative of exported claims so
	// runtimes can observe routing roles across the grid without creating a
	// second independent declaration surface. Source: DI-ruvot
	for _, route := range batch.Routes {
		if strings.TrimSpace(route.PackageID) == "" {
			return errors.New("route package_id is required")
		}
		if strings.TrimSpace(route.PackageVersion) == "" {
			return errors.New("route package_version is required")
		}
		if strings.TrimSpace(route.ProtocolPCID) == "" {
			return errors.New("route protocol_pcid is required")
		}
		if strings.TrimSpace(route.Role) == "" {
			return errors.New("route role is required")
		}
		key := route.PackageID + "\x00" + route.PackageVersion + "\x00" + route.ProtocolPCID + "\x00" + route.Role
		if _, exists := seenRoutes[key]; exists {
			return fmt.Errorf("duplicate route registration: %s", key)
		}
		if _, exists := seenClaims[key]; !exists {
			return fmt.Errorf("route registration missing matching claim: %s", key)
		}
		familyNames := map[string]struct{}{}
		for _, family := range route.Families {
			if strings.TrimSpace(family) == "" {
				return fmt.Errorf("route family is required for %s", key)
			}
			if _, exists := familyNames[family]; exists {
				return fmt.Errorf("duplicate route family %q for %s", family, key)
			}
			familyNames[family] = struct{}{}
		}
		seenRoutes[key] = struct{}{}
	}
	if len(batch.ClaimProofs) > 0 && len(batch.ClaimProofs) != len(batch.ImplementationClaims) {
		return errors.New("claim_proofs must match implementation_claims length")
	}
	for _, proof := range batch.ClaimProofs {
		if strings.TrimSpace(proof.SignerPeerID) == "" {
			return errors.New("claim proof signer_peer_id is required")
		}
		if strings.TrimSpace(proof.PublicKey) == "" {
			return errors.New("claim proof public_key is required")
		}
		if strings.TrimSpace(proof.Signature) == "" {
			return errors.New("claim proof signature is required")
		}
	}
	for _, attestation := range batch.ClaimAttestations {
		if attestation.ClaimIndex < 0 || attestation.ClaimIndex >= len(batch.ImplementationClaims) {
			return fmt.Errorf("claim attestation index out of range: %d", attestation.ClaimIndex)
		}
		if strings.TrimSpace(attestation.SignerPeerID) == "" {
			return errors.New("claim attestation signer_peer_id is required")
		}
		if strings.TrimSpace(attestation.PublicKey) == "" {
			return errors.New("claim attestation public_key is required")
		}
		if strings.TrimSpace(attestation.Signature) == "" {
			return errors.New("claim attestation signature is required")
		}
	}
	if len(batch.Records) == 0 {
		return errors.New("batch records are required")
	}
	seenRecords := map[string]struct{}{}
	for _, record := range batch.Records {
		if len(record) == 0 {
			return errors.New("batch record must not be empty")
		}
		key := string(record)
		if _, exists := seenRecords[key]; exists {
			return errors.New("duplicate raw record in batch")
		}
		seenRecords[key] = struct{}{}
	}
	if len(batch.RecordProofs) > 0 && len(batch.RecordProofs) != len(batch.Records) {
		return errors.New("record_proofs must match records length")
	}
	for _, proof := range batch.RecordProofs {
		if strings.TrimSpace(proof.Digest) == "" {
			return errors.New("record proof digest is required")
		}
	}
	if len(batch.RecordSignatures) > 0 && len(batch.RecordSignatures) != len(batch.Records) {
		return errors.New("record_signatures must match records length")
	}
	for _, signature := range batch.RecordSignatures {
		if strings.TrimSpace(signature.SignerPeerID) == "" {
			return errors.New("record signature signer_peer_id is required")
		}
		if strings.TrimSpace(signature.PublicKey) == "" {
			return errors.New("record signature public_key is required")
		}
		if strings.TrimSpace(signature.Signature) == "" {
			return errors.New("record signature signature is required")
		}
	}
	return nil
}

func (batch Batch) SigningBytes() ([]byte, error) {
	// Intent: Sign a stable view of the current relay batch without inventing a
	// new wire format so live peers can verify who exported the batch.
	// Source: DI-zotem
	signable := struct {
		Format               string                `json:"format"`
		Implementation       string                `json:"implementation"`
		ExportedAt           string                `json:"exported_at"`
		ImplementationClaims []ImplementationClaim `json:"implementation_claims,omitempty"`
		Routes               []RouteRegistration   `json:"routes,omitempty"`
		ClaimProofs          []ClaimProof          `json:"claim_proofs,omitempty"`
		ClaimAttestations    []ClaimAttestation    `json:"claim_attestations,omitempty"`
		Records              []json.RawMessage     `json:"records"`
		RecordProofs         []RecordProof         `json:"record_proofs,omitempty"`
		RecordSignatures     []RecordSignature     `json:"record_signatures,omitempty"`
	}{
		Format:               batch.Format,
		Implementation:       batch.Implementation,
		ExportedAt:           batch.ExportedAt,
		ImplementationClaims: batch.ImplementationClaims,
		Routes:               batch.Routes,
		ClaimProofs:          batch.ClaimProofs,
		ClaimAttestations:    batch.ClaimAttestations,
		Records:              batch.Records,
		RecordProofs:         batch.RecordProofs,
		RecordSignatures:     batch.RecordSignatures,
	}
	return json.Marshal(signable)
}

func ClaimSigningBytes(claim ImplementationClaim) []byte {
	body, err := json.Marshal(claim)
	if err != nil {
		return []byte(claim.PackageID + "\x00" + claim.PackageVersion + "\x00" + claim.ProtocolPCID + "\x00" + claim.Role + "\x00" + claim.Summary)
	}
	return body
}

func ClaimAttestationSigningBytes(claimIndex int, claim ImplementationClaim) []byte {
	signable := struct {
		ClaimIndex int                 `json:"claim_index"`
		Claim      ImplementationClaim `json:"claim"`
	}{
		ClaimIndex: claimIndex,
		Claim:      claim,
	}
	body, err := json.Marshal(signable)
	if err != nil {
		return append([]byte(fmt.Sprintf("%d\x00", claimIndex)), ClaimSigningBytes(claim)...)
	}
	return body
}

func ProofForRecord(raw json.RawMessage) RecordProof {
	canonical := RecordSigningBytes(raw)
	sum := sha256.Sum256(canonical)
	return RecordProof{Digest: "sha256:" + hex.EncodeToString(sum[:])}
}

func ProofsForRecords(records []json.RawMessage) []RecordProof {
	proofs := make([]RecordProof, 0, len(records))
	for _, record := range records {
		proofs = append(proofs, ProofForRecord(record))
	}
	return proofs
}

func (batch Batch) VerifyRecordProofs() error {
	// Intent: Carry per-record digests inside relay batches so receivers can
	// detect tampering or accidental mutation before durable import.
	// Source: DI-zumep
	if len(batch.RecordProofs) == 0 {
		return nil
	}
	if len(batch.RecordProofs) != len(batch.Records) {
		return errors.New("record_proofs must match records length")
	}
	for index, record := range batch.Records {
		expected := ProofForRecord(record)
		if batch.RecordProofs[index].Digest != expected.Digest {
			return fmt.Errorf("record proof mismatch at index %d", index)
		}
	}
	return nil
}

func (batch Batch) VerifyClaimProofs() error {
	// Intent: Carry per-claim proofs so receivers can verify that the exporting
	// peer actually signed the implementation claims it is advertising.
	// Source: DI-luzef
	if len(batch.ClaimProofs) == 0 {
		return nil
	}
	if len(batch.ClaimProofs) != len(batch.ImplementationClaims) {
		return errors.New("claim_proofs must match implementation_claims length")
	}
	return nil
}

func (batch Batch) VerifyClaimAttestations() error {
	// Intent: Carry third-party attestations over specific implementation claims
	// so receivers can distinguish exporter self-claims from outside approval.
	// Source: DI-fogem
	for _, attestation := range batch.ClaimAttestations {
		if attestation.ClaimIndex < 0 || attestation.ClaimIndex >= len(batch.ImplementationClaims) {
			return fmt.Errorf("claim attestation index out of range: %d", attestation.ClaimIndex)
		}
	}
	return nil
}

func RecordSigningBytes(raw json.RawMessage) []byte {
	var compact bytes.Buffer
	if err := json.Compact(&compact, raw); err != nil {
		return append([]byte{}, raw...)
	}
	return compact.Bytes()
}
