package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex7-makerspace-stewardship/service"
)

func main() {
	root := flag.String("root", "", "empty parent directory for Alice and Bob agent roots")
	passphraseStdin := flag.Bool("passphrase-stdin", false, "read bootstrap identity passphrase from standard input")
	flag.Parse()
	if *root == "" || !*passphraseStdin {
		fmt.Fprintln(os.Stderr, "root and -passphrase-stdin are required")
		os.Exit(2)
	}
	passphrase, err := io.ReadAll(io.LimitReader(os.Stdin, 4097))
	if err != nil {
		panic(err)
	}
	passphrase = bytes.TrimSpace(passphrase)
	if len(passphrase) == 0 {
		fmt.Fprintln(os.Stderr, "passphrase is required")
		os.Exit(2)
	}
	if err := os.MkdirAll(*root, 0o700); err != nil {
		panic(err)
	}
	alice, err := service.CreateParticipantIdentity(filepath.Join(*root, "alice"), "alice", passphrase)
	if err != nil {
		panic(err)
	}
	bob, err := service.CreateParticipantIdentity(filepath.Join(*root, "bob"), "bob", passphrase)
	if err != nil {
		panic(err)
	}
	for _, agent := range []struct {
		name     string
		identity *service.ParticipantIdentity
	}{{"alice", alice}, {"bob", bob}} {
		public := agent.identity.Private.Public().(ed25519.PublicKey)
		recognition, _ := json.Marshal(map[string]any{"version": 1, "keys": []map[string]string{{"label": "alice", "ed25519_public_key_base64": base64.StdEncoding.EncodeToString(alice.Private.Public().(ed25519.PublicKey))}, {"label": "bob", "ed25519_public_key_base64": base64.StdEncoding.EncodeToString(bob.Private.Public().(ed25519.PublicKey))}}})
		if err := os.WriteFile(filepath.Join(*root, agent.name, "recognition.json"), recognition, 0o600); err != nil {
			panic(err)
		}
		rootKey := mustPrivate()
		recovery := []string{base64.StdEncoding.EncodeToString(mustPrivate().Public().(ed25519.PublicKey)), base64.StdEncoding.EncodeToString(mustPrivate().Public().(ed25519.PublicKey)), base64.StdEncoding.EncodeToString(mustPrivate().Public().(ed25519.PublicKey))}
		now := time.Now().UTC().Format(time.RFC3339)
		rootRecord := service.Record{Protocol: "bafkreia7cn4srmmkxbwxk2hoezedjvuyokhypcsddjd4evx56lhtmsq3nm", ID: agent.name + "-root-1", Signer: agent.name + " root", CreatedAt: now, Payload: []byte(fmt.Sprintf(`{"root_key":"%s","history_note":"bootstrap","recovery_set":["%s","%s","%s"]}`, base64.StdEncoding.EncodeToString(rootKey.Public().(ed25519.PublicKey)), recovery[0], recovery[1], recovery[2])), KeyID: keyID(rootKey.Public().(ed25519.PublicKey)), PublicKey: rootKey.Public().(ed25519.PublicKey)}
		_, rootRaw, err := rootRecord.Sign(rootKey)
		if err != nil {
			panic(err)
		}
		deviceRecord := service.Record{Protocol: "bafkreifmbhgjwfmwbemkf4ogsg3gvuavjhttkitzf3muie3dhv5tdn4hq4", ID: agent.name + "-device-1", Signer: agent.name + " root", CreatedAt: now, Payload: []byte(fmt.Sprintf(`{"root_record_id":"%s","device_key":"%s","device_label":"%s device","not_before":"%s"}`, rootRecord.ID, base64.StdEncoding.EncodeToString(public), agent.name, now)), KeyID: keyID(rootKey.Public().(ed25519.PublicKey)), PublicKey: rootKey.Public().(ed25519.PublicKey)}
		_, deviceRaw, err := deviceRecord.Sign(rootKey)
		if err != nil {
			panic(err)
		}
		app, err := service.NewPersistentParticipantApp(filepath.Join(*root, agent.name), service.NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": alice.Private.Public().(ed25519.PublicKey), "bob": bob.Private.Public().(ed25519.PublicKey)}), agent.identity)
		if err != nil {
			panic(err)
		}
		if err := app.IngestRecords([][]byte{rootRaw, deviceRaw}); err != nil {
			panic(err)
		}
		if agent.name == "alice" {
			peerCard := service.Record{Protocol: "bafkreicstci6idwm6d6dbt52ppqyjcapibskz27qmnfuyntg6zck72fa24", ID: "alice-peer-card-1", Signer: "alice", CreatedAt: now, Payload: []byte(fmt.Sprintf(`{"root_record_id":"%s","active_device_record_ids":["%s"],"contact_hints":[]}`, rootRecord.ID, deviceRecord.ID)), KeyID: keyID(public), PublicKey: public}
			_, peerCardRaw, peerCardErr := peerCard.Sign(agent.identity.Private)
			if peerCardErr != nil {
				panic(peerCardErr)
			}
			carrierSeed := sha256.Sum256([]byte("ex7 deterministic bootstrap carrier"))
			carrierKey := ed25519.NewKeyFromSeed(carrierSeed[:])
			carrierPublic := carrierKey.Public().(ed25519.PublicKey)
			carriage := service.Record{Protocol: "bafkreihrlojt47erjc6uawkm47s7evppp23tk3ljlkl347ten4v3kb624i", ID: "alice-bootstrap-carriage-1", Signer: "bootstrap carrier", CreatedAt: now, Payload: []byte(fmt.Sprintf(`{"sender_card_record_id":"%s","cursor":"bootstrap-1","records":["%s","%s","%s"]}`, peerCard.ID, base64.StdEncoding.EncodeToString(rootRaw), base64.StdEncoding.EncodeToString(deviceRaw), base64.StdEncoding.EncodeToString(peerCardRaw))), KeyID: keyID(carrierPublic), PublicKey: carrierPublic}
			_, carriageRaw, carriageErr := carriage.Sign(carrierKey)
			if carriageErr != nil {
				panic(carriageErr)
			}
			// Intent: Preserve exact public bootstrap evidence for later peer-card
			// and carriage delivery; this file never contains private-key material.
			// Source: DI-fuzar.
			artifact, marshalErr := json.Marshal(struct {
				Records []string `json:"records"`
			}{Records: []string{base64.StdEncoding.EncodeToString(rootRaw), base64.StdEncoding.EncodeToString(deviceRaw), base64.StdEncoding.EncodeToString(peerCardRaw), base64.StdEncoding.EncodeToString(carriageRaw)}})
			if marshalErr != nil {
				panic(marshalErr)
			}
			if writeErr := os.WriteFile(filepath.Join(*root, "alice", "bootstrap-records.json"), artifact, 0o600); writeErr != nil {
				panic(writeErr)
			}
		}
	}
	artifactBytes, artifactReadErr := os.ReadFile(filepath.Join(*root, "alice", "bootstrap-records.json"))
	if artifactReadErr != nil {
		panic(artifactReadErr)
	}
	var artifact struct {
		Records []string `json:"records"`
	}
	if artifactDecodeErr := json.Unmarshal(artifactBytes, &artifact); artifactDecodeErr != nil || len(artifact.Records) != 4 {
		panic("invalid Alice bootstrap artifact")
	}
	aliceHistory := make([][]byte, 3)
	for index := range aliceHistory {
		decoded, decodeErr := base64.StdEncoding.DecodeString(artifact.Records[index])
		if decodeErr != nil {
			panic(decodeErr)
		}
		aliceHistory[index] = decoded
	}
	carriageRaw, carriageDecodeErr := base64.StdEncoding.DecodeString(artifact.Records[3])
	if carriageDecodeErr != nil {
		panic(carriageDecodeErr)
	}
	bobApp, bobOpenErr := service.NewPersistentParticipantApp(filepath.Join(*root, "bob"), service.NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": alice.Private.Public().(ed25519.PublicKey), "bob": bob.Private.Public().(ed25519.PublicKey)}), bob)
	if bobOpenErr != nil {
		panic(bobOpenErr)
	}
	if bobHistoryErr := bobApp.IngestRecords(aliceHistory); bobHistoryErr != nil {
		panic(bobHistoryErr)
	}
	if bobCarriageErr := bobApp.IngestRecords([][]byte{carriageRaw}); bobCarriageErr != nil {
		panic(bobCarriageErr)
	}
	fmt.Printf("bootstrap complete: %s\n", *root)
}

func mustPrivate() ed25519.PrivateKey {
	seed := make([]byte, ed25519.SeedSize)
	if _, err := rand.Read(seed); err != nil {
		panic(err)
	}
	return ed25519.NewKeyFromSeed(seed)
}
func keyID(public ed25519.PublicKey) string {
	return "ed25519:" + fmt.Sprintf("%x", serviceKeyDigest(public))
}
func serviceKeyDigest(public ed25519.PublicKey) [32]byte { return sha256.Sum256(public) }
