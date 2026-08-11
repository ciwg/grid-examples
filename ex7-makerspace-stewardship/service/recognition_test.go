package service

import (
	"crypto/ed25519"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestRecognitionPolicyRequiresMatchingKeyAndLabel(t *testing.T) {
	record, private := testRecord(t)
	_, raw, err := record.Sign(private)
	if err != nil {
		t.Fatalf("sign record: %v", err)
	}
	parsed, err := ParseRecord(raw)
	if err != nil {
		t.Fatalf("parse record: %v", err)
	}
	policy := NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": parsed.PublicKey})
	if !policy.recognizes(parsed) {
		t.Fatal("recognized key and label were rejected")
	}
	parsed.Signer = "carol"
	if policy.recognizes(parsed) {
		t.Fatal("recognized key was accepted under a different label")
	}
}

func TestLoadRecognitionPolicyRejectsInsecureAndMalformedFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "recognition.json")
	if err := os.WriteFile(path, []byte(`{"version":1,"keys":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRecognitionPolicy(path, false); err == nil {
		t.Fatal("accepted insecure recognition policy")
	}
	if err := os.WriteFile(path, []byte(`{"version":1,"keys":[{"label":"alice","ed25519_public_key_base64":"not-base64"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRecognitionPolicy(path, false); err == nil {
		t.Fatal("accepted malformed recognition policy")
	}
}

func TestLoadRecognitionPolicyLoadsPublicKeyOnly(t *testing.T) {
	_, private := testRecord(t)
	public := private.Public().(ed25519.PublicKey)
	path := filepath.Join(t.TempDir(), "recognition.json")
	data := []byte(`{"version":1,"keys":[{"label":"alice","ed25519_public_key_base64":"` + base64.StdEncoding.EncodeToString(public) + `"}]}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	policy, err := LoadRecognitionPolicy(path, false)
	if err != nil {
		t.Fatal(err)
	}
	if policy.keyLabels[keyID(public)] != "alice" {
		t.Fatal("loaded key is not recognized as alice")
	}
}
