package analyzer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
)

type Summary struct {
	MessageCount    int      `json:"message_count"`
	UniquePCIDs     []string `json:"unique_pcids"`
	ArtifactCount   int      `json:"artifact_count"`
	HasIntakeResult bool     `json:"has_intake_result"`
}

func Analyze(root string) (summary Summary, err error) {
	dagPath := filepath.Join(root, "collector", "message-dag.jsonl")
	file, err := os.Open(dagPath)
	if err != nil {
		return Summary{}, fmt.Errorf("open dag file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	seenPCIDs := map[string]bool{}
	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
		var record struct {
			PCID     string `json:"pcid"`
			ExactCID string `json:"exact_cid"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return Summary{}, fmt.Errorf("decode dag record: %w", err)
		}
		if _, err := protocol.ProfileByCIDText(record.PCID); err != nil {
			return Summary{}, fmt.Errorf("unknown dag pcid %q: %w", record.PCID, err)
		}
		seenPCIDs[record.PCID] = true
		artifactPath := filepath.Join(root, "collector", "message-cas", record.ExactCID+".cbor")
		envelopeBytes, err := os.ReadFile(artifactPath)
		if err != nil {
			return Summary{}, fmt.Errorf("read artifact %s: %w", artifactPath, err)
		}
		envelope, err := protocol.ParseEnvelope(envelopeBytes)
		if err != nil {
			return Summary{}, fmt.Errorf("parse artifact envelope: %w", err)
		}
		if len(envelope.ProofBytes) == 0 {
			return Summary{}, fmt.Errorf("artifact %s has no proof bytes", record.ExactCID)
		}
		var payload map[string]any
		if err := protocol.Unmarshal(envelope.PayloadBytes, &payload); err != nil {
			return Summary{}, fmt.Errorf("decode payload map: %w", err)
		}
		if _, exists := payload["protocol"]; exists {
			return Summary{}, fmt.Errorf("payload repeated protocol name in artifact %s", record.ExactCID)
		}
	}
	if err := scanner.Err(); err != nil {
		return Summary{}, fmt.Errorf("scan dag file: %w", err)
	}
	artifacts, err := filepath.Glob(filepath.Join(root, "collector", "message-cas", "*.cbor"))
	if err != nil {
		return Summary{}, fmt.Errorf("glob artifacts: %w", err)
	}
	unique := make([]string, 0, len(seenPCIDs))
	for pcidText := range seenPCIDs {
		unique = append(unique, pcidText)
	}
	_, intakeErr := os.Stat(filepath.Join(root, "intake", "result.json"))
	summary = Summary{
		MessageCount:    count,
		UniquePCIDs:     unique,
		ArtifactCount:   len(artifacts),
		HasIntakeResult: intakeErr == nil,
	}
	if len(unique) > 4 {
		return summary, fmt.Errorf("unique business pcids = %d, want <= 4", len(unique))
	}
	if summary.ArtifactCount == 0 {
		return summary, fmt.Errorf("no retained message artifacts")
	}
	if !summary.HasIntakeResult {
		return summary, fmt.Errorf("missing intake result")
	}
	return summary, nil
}
