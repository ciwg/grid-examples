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

type RecordProof struct {
	Digest string `json:"digest"`
}

type Batch struct {
	Format               string                `json:"format"`
	Implementation       string                `json:"implementation"`
	ExportedAt           string                `json:"exported_at"`
	ImplementationClaims []ImplementationClaim `json:"implementation_claims,omitempty"`
	Records              []json.RawMessage     `json:"records"`
	RecordProofs         []RecordProof         `json:"record_proofs,omitempty"`
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
		Records              []json.RawMessage     `json:"records"`
		RecordProofs         []RecordProof         `json:"record_proofs,omitempty"`
	}{
		Format:               batch.Format,
		Implementation:       batch.Implementation,
		ExportedAt:           batch.ExportedAt,
		ImplementationClaims: batch.ImplementationClaims,
		Records:              batch.Records,
		RecordProofs:         batch.RecordProofs,
	}
	return json.Marshal(signable)
}

func ProofForRecord(raw json.RawMessage) RecordProof {
	canonical := canonicalRecordBytes(raw)
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

func canonicalRecordBytes(raw json.RawMessage) []byte {
	var compact bytes.Buffer
	if err := json.Compact(&compact, raw); err != nil {
		return append([]byte{}, raw...)
	}
	return compact.Bytes()
}
