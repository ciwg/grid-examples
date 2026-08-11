// Command makerspace-record-fixture creates deterministic test-only signed
// ingress material. It is not a participant key-management tool. Source:
// DI-likoh.
package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex7-makerspace-stewardship/service"
)

const observationPCID = "bafkreifhodcald6kzib36rzeji27hnqjdkeycibnkcigcsz7mzejz6obiy"
const safetyPCID = "bafkreigt3p2l4uel7wmjr4kple7o55ymchlhh43gajjwsgaeifoogeztc4"

func main() {
	runtimeRoot := flag.String("runtime-root", "", "agent runtime root to receive recognition.json")
	label := flag.String("label", "alice", "fixture signer label")
	seedText := flag.String("seed", "ex7-operator-proof", "deterministic test-only signer seed")
	toolID := flag.String("tool", "table-saw", "observed tool ID")
	observation := flag.String("observation", "Guard is loose", "observation text")
	kind := flag.String("kind", "observation", "fixture kind: observation or safety-hold or safety-clear")
	writePolicy := flag.Bool("write-recognition", true, "write this public key to recognition.json")
	flag.Parse()
	if *runtimeRoot == "" || *label == "" || *toolID == "" || *observation == "" {
		fmt.Fprintln(os.Stderr, "runtime-root, label, tool, and observation are required")
		os.Exit(2)
	}
	seed := sha256.Sum256([]byte(*seedText))
	private := ed25519.NewKeyFromSeed(seed[:])
	public := private.Public().(ed25519.PublicKey)
	digest := sha256.Sum256(public)
	keyID := "ed25519:" + hex.EncodeToString(digest[:])
	protocol := observationPCID
	payloadValue := map[string]string{"observation": *observation, "tool_id": *toolID}
	if *kind == "safety-hold" || *kind == "safety-clear" {
		protocol = safetyPCID
		disposition := "hold"
		if *kind == "safety-clear" {
			disposition = "clear"
		}
		payloadValue = map[string]string{"assessment": *observation, "disposition": disposition, "tool_id": *toolID}
	}
	payload, err := json.Marshal(payloadValue)
	if err != nil {
		panic(err)
	}
	record := service.Record{Protocol: protocol, ID: "fixture-" + *kind, Signer: *label, CreatedAt: time.Now().UTC().Format(time.RFC3339), Payload: payload, KeyID: keyID, PublicKey: public}
	_, raw, err := record.Sign(private)
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(*runtimeRoot, 0o750); err != nil {
		panic(err)
	}
	if *writePolicy {
		policy, err := json.Marshal(map[string]any{"version": 1, "keys": []map[string]string{{"label": *label, "ed25519_public_key_base64": base64.StdEncoding.EncodeToString(public)}}})
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(filepath.Join(*runtimeRoot, "recognition.json"), policy, 0o600); err != nil {
			panic(err)
		}
	}
	fmt.Println(base64.StdEncoding.EncodeToString(raw))
}
