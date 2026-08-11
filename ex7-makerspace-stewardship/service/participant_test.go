package service

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"testing"
	"time"
)

func TestParticipantHistoryRequiresRootAuthorizedDeviceAndHonorsRevocation(t *testing.T) {
	bootstrap, devicePublic := testParticipantBootstrap(t, "alice")
	app, err := NewPersistentRecordApp(t.TempDir(), NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": devicePublic}))
	if err != nil {
		t.Fatalf("create participant app: %v", err)
	}
	if err := app.IngestRecords(bootstrap); err != nil {
		t.Fatalf("ingest root and device history: %v", err)
	}

	observation, _ := signedTestRecord(t, observationPCID, "alice-observation", "alice", `{"observation":"Guard is secure","tool_id":"table-saw"}`)
	if err := app.IngestRecords([][]byte{observation}); err != nil {
		t.Fatalf("ingest authorized device record: %v", err)
	}

	root := testPrivate("alice root")
	rootPublic := root.Public().(ed25519.PublicKey)
	revocation := Record{
		Protocol:  revocationPCID,
		ID:        "alice-device-revocation",
		Signer:    "alice root",
		CreatedAt: "2026-08-11T22:20:00Z",
		Payload: []byte(fmt.Sprintf(`{"subject_key_id":"%s","subject_kind":"device","effective_at":"2026-08-11T22:20:00Z","reason":"device lost"}`,
			keyID(devicePublic))),
		KeyID:     keyID(rootPublic),
		PublicKey: rootPublic,
	}
	_, revocationRaw, err := revocation.Sign(root)
	if err != nil {
		t.Fatalf("sign revocation: %v", err)
	}
	if err := app.IngestRecords([][]byte{revocationRaw}); err != nil {
		t.Fatalf("ingest revocation: %v", err)
	}

	device := testPrivate("alice")
	after := Record{Protocol: observationPCID, ID: "alice-after-revocation", Signer: "alice", CreatedAt: "2026-08-11T22:21:00Z", Payload: []byte(`{"observation":"Should not project","tool_id":"table-saw"}`), KeyID: keyID(device.Public().(ed25519.PublicKey)), PublicKey: device.Public().(ed25519.PublicKey)}
	_, afterRevocation, err := after.Sign(device)
	if err != nil {
		t.Fatalf("sign post-revocation record: %v", err)
	}
	if err := app.IngestRecords([][]byte{afterRevocation}); err == nil {
		t.Fatal("accepted record from revoked device")
	}
}

func TestParticipantSignerSignsOnlyWithActiveHistoryLinkedKey(t *testing.T) {
	bootstrap, _ := testParticipantBootstrap(t, "alice")
	history := NewParticipantHistory()
	for _, raw := range bootstrap {
		record, err := ParseRecord(raw)
		if err != nil {
			t.Fatalf("parse bootstrap: %v", err)
		}
		if err := history.Apply(record); err != nil {
			t.Fatalf("apply bootstrap: %v", err)
		}
	}
	signer, err := NewParticipantSigner(testPrivate("alice"), "alice", history)
	if err != nil {
		t.Fatalf("create device signer: %v", err)
	}
	signed, raw, err := signer.Sign(Record{Protocol: observationPCID, ID: "signed-by-device", CreatedAt: "2026-08-11T22:10:00Z", Payload: []byte(`{"observation":"Guard is secure","tool_id":"table-saw"}`)})
	if err != nil {
		t.Fatalf("sign active device record: %v", err)
	}
	if signed.KeyID == "" || string(raw) == "" {
		t.Fatal("signer did not produce exact signed record bytes")
	}
	if _, err := ParseRecord(raw); err != nil {
		t.Fatalf("parse signed record: %v", err)
	}

	unknown, err := NewParticipantSigner(testPrivate("mallory"), "mallory", history)
	if err != nil {
		t.Fatalf("create unknown signer: %v", err)
	}
	if _, _, err := unknown.Sign(Record{Protocol: observationPCID, ID: "mallory-record", CreatedAt: "2026-08-11T22:10:00Z", Payload: []byte(`{"observation":"No authorization","tool_id":"table-saw"}`)}); err == nil {
		t.Fatal("unknown device signed a participant record")
	}
}

func TestRootHistoryRejectsUnannouncedRecoverySet(t *testing.T) {
	root := testPrivate("alice root")
	public := root.Public().(ed25519.PublicKey)
	recovery := base64.StdEncoding.EncodeToString(testPrivate("recovery").Public().(ed25519.PublicKey))
	record := Record{Protocol: rootHistoryPCID, ID: "bad-root", Signer: "alice", CreatedAt: "2026-08-11T22:00:00Z", Payload: []byte(fmt.Sprintf(`{"root_key":"%s","history_note":"bad","recovery_set":["%s","%s","%s"]}`, base64.StdEncoding.EncodeToString(public), recovery, recovery, recovery)), KeyID: keyID(public), PublicKey: public}
	_, raw, err := record.Sign(root)
	if err != nil {
		t.Fatalf("sign malformed root history: %v", err)
	}
	parsed, err := ParseRecord(raw)
	if err != nil {
		t.Fatalf("parse signed root history: %v", err)
	}
	if err := NewParticipantHistory().Apply(parsed); err == nil {
		t.Fatal("accepted duplicate recovery keys")
	}
}

func TestParticipantHistoryRejectsCrossParticipantRevocation(t *testing.T) {
	aliceBootstrap, _ := testParticipantBootstrap(t, "alice")
	bobBootstrap, bobDevice := testParticipantBootstrap(t, "bob")
	history := NewParticipantHistory()
	for _, raw := range append(aliceBootstrap, bobBootstrap...) {
		record, err := ParseRecord(raw)
		if err != nil {
			t.Fatalf("parse bootstrap: %v", err)
		}
		if err := history.Apply(record); err != nil {
			t.Fatalf("apply bootstrap: %v", err)
		}
	}
	aliceRoot := testPrivate("alice root")
	alicePublic := aliceRoot.Public().(ed25519.PublicKey)
	revocation := Record{Protocol: revocationPCID, ID: "alice-revokes-bob", Signer: "alice root", CreatedAt: "2026-08-11T22:20:00Z", Payload: []byte(fmt.Sprintf(`{"subject_key_id":"%s","subject_kind":"device","effective_at":"2026-08-11T22:20:00Z","reason":"not alice's key"}`, keyID(bobDevice))), KeyID: keyID(alicePublic), PublicKey: alicePublic}
	_, raw, err := revocation.Sign(aliceRoot)
	if err != nil {
		t.Fatalf("sign cross-participant revocation: %v", err)
	}
	parsed, err := ParseRecord(raw)
	if err != nil {
		t.Fatalf("parse cross-participant revocation: %v", err)
	}
	if err := history.Apply(parsed); err == nil {
		t.Fatal("accepted cross-participant revocation")
	}
}

func TestParticipantHistoryActivatesReplacementRootOnlyAfterTwoRecoveryWitnesses(t *testing.T) {
	bootstrap, _ := testParticipantBootstrap(t, "alice")
	history := NewParticipantHistory()
	for _, raw := range bootstrap {
		record, err := ParseRecord(raw)
		if err != nil {
			t.Fatalf("parse bootstrap: %v", err)
		}
		if err := history.Apply(record); err != nil {
			t.Fatalf("apply bootstrap: %v", err)
		}
	}
	replacement := testPrivate("alice replacement root")
	replacementPublic := replacement.Public().(ed25519.PublicKey)
	recoverySet := []string{
		base64.StdEncoding.EncodeToString(testPrivate("alice recovery 1").Public().(ed25519.PublicKey)),
		base64.StdEncoding.EncodeToString(testPrivate("alice recovery 2").Public().(ed25519.PublicKey)),
		base64.StdEncoding.EncodeToString(testPrivate("alice recovery 3").Public().(ed25519.PublicKey)),
	}
	payload := []byte(fmt.Sprintf(`{"root_record_id":"alice-root-1","recovery_id":"lost-laptop","replacement_root_key":"%s","recovery_set":["%s","%s","%s"]}`,
		base64.StdEncoding.EncodeToString(replacementPublic), recoverySet[0], recoverySet[1], recoverySet[2]))
	for index, witness := range []ed25519.PrivateKey{testPrivate("alice recovery 1"), testPrivate("alice recovery 2")} {
		public := witness.Public().(ed25519.PublicKey)
		record := Record{Protocol: recoveryPCID, ID: fmt.Sprintf("alice-recovery-%d", index), Signer: fmt.Sprintf("alice recovery %d", index+1), CreatedAt: "2026-08-11T22:20:00Z", Payload: payload, KeyID: keyID(public), PublicKey: public}
		_, raw, err := record.Sign(witness)
		if err != nil {
			t.Fatalf("sign recovery witness: %v", err)
		}
		parsed, err := ParseRecord(raw)
		if err != nil {
			t.Fatalf("parse recovery witness: %v", err)
		}
		if err := history.Apply(parsed); err != nil {
			t.Fatalf("apply recovery witness: %v", err)
		}
		if index == 0 {
			continuation := Record{Protocol: rootHistoryPCID, ID: "alice-root-2", Signer: "alice replacement root", CreatedAt: "2026-08-11T22:21:00Z", Payload: []byte(fmt.Sprintf(`{"root_key":"%s","previous_root_record_id":"alice-root-1","history_note":"recovered","recovery_set":["%s","%s","%s"]}`, base64.StdEncoding.EncodeToString(replacementPublic), recoverySet[0], recoverySet[1], recoverySet[2])), KeyID: keyID(replacementPublic), PublicKey: replacementPublic}
			_, raw, err := continuation.Sign(replacement)
			if err != nil {
				t.Fatalf("sign premature continuation: %v", err)
			}
			parsed, err := ParseRecord(raw)
			if err != nil {
				t.Fatalf("parse premature continuation: %v", err)
			}
			if err := history.Apply(parsed); err == nil {
				t.Fatal("one recovery witness activated replacement root")
			}
		}
	}
	continuation := Record{Protocol: rootHistoryPCID, ID: "alice-root-2", Signer: "alice replacement root", CreatedAt: "2026-08-11T22:22:00Z", Payload: []byte(fmt.Sprintf(`{"root_key":"%s","previous_root_record_id":"alice-root-1","history_note":"recovered","recovery_set":["%s","%s","%s"]}`, base64.StdEncoding.EncodeToString(replacementPublic), recoverySet[0], recoverySet[1], recoverySet[2])), KeyID: keyID(replacementPublic), PublicKey: replacementPublic}
	_, raw, err := continuation.Sign(replacement)
	if err != nil {
		t.Fatalf("sign recovery continuation: %v", err)
	}
	parsed, err := ParseRecord(raw)
	if err != nil {
		t.Fatalf("parse recovery continuation: %v", err)
	}
	if err := history.Apply(parsed); err != nil {
		t.Fatalf("two witnesses did not activate replacement root: %v", err)
	}
}

func TestPeerCardLinkedCarriageRetainsAndProjectsExactRecord(t *testing.T) {
	bootstrap, alicePublic := testParticipantBootstrap(t, "alice")
	device := testPrivate("alice")
	card := Record{Protocol: peerCardPCID, ID: "alice-card-1", Signer: "alice", CreatedAt: "2026-08-11T22:02:00Z", Payload: []byte(`{"root_record_id":"alice-root-1","active_device_record_ids":["alice-device-1"],"contact_hints":[]}`), KeyID: keyID(alicePublic), PublicKey: alicePublic}
	_, cardRaw, err := card.Sign(device)
	if err != nil {
		t.Fatalf("sign peer card: %v", err)
	}
	observation, _ := signedTestRecord(t, observationPCID, "alice-carried-observation", "alice", `{"observation":"Carried exactly","tool_id":"table-saw"}`)
	carrier := testPrivate("carrier")
	carrierPublic := carrier.Public().(ed25519.PublicKey)
	carriage := Record{Protocol: carriagePCID, ID: "carrier-batch-1", Signer: "carrier", CreatedAt: "2026-08-11T22:10:00Z", Payload: []byte(fmt.Sprintf(`{"sender_card_record_id":"alice-card-1","cursor":"1","records":["%s"]}`, base64.StdEncoding.EncodeToString(observation))), KeyID: keyID(carrierPublic), PublicKey: carrierPublic}
	_, carriageRaw, err := carriage.Sign(carrier)
	if err != nil {
		t.Fatalf("sign carriage: %v", err)
	}
	app, err := NewPersistentRecordApp(t.TempDir(), NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": alicePublic}))
	if err != nil {
		t.Fatalf("create participant app: %v", err)
	}
	if err := app.IngestRecords(append(append(bootstrap, cardRaw), carriageRaw)); err != nil {
		t.Fatalf("ingest peer card and carriage: %v", err)
	}
	if got := len(app.State().Tools[0].Observations); got != 1 {
		t.Fatalf("carried observation projection = %d", got)
	}
	frames, err := app.store.ReadRecordFrames()
	if err != nil {
		t.Fatalf("read carried exact bytes: %v", err)
	}
	if len(frames) != 1 || len(frames[0]) != 5 || !bytes.Equal(frames[0][4], observation) {
		t.Fatal("carriage did not retain enclosed exact bytes")
	}
}

func TestTerminalApprovalCreatesDurableSignedRecordOnlyAfterApproval(t *testing.T) {
	root := t.TempDir()
	identity := &ParticipantIdentity{Label: "alice", Private: testPrivate("alice")}
	app, err := NewPersistentParticipantApp(root, NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": identity.Private.Public().(ed25519.PublicKey)}), identity)
	if err != nil {
		t.Fatalf("create participant app: %v", err)
	}
	bootstrap, _ := testParticipantBootstrap(t, "alice")
	if err := app.IngestRecords(bootstrap); err != nil {
		t.Fatalf("ingest participant history: %v", err)
	}
	now := time.Date(2026, 8, 11, 22, 10, 0, 0, time.UTC)
	request, err := app.CreateApprovalRequest(observationPCID, []byte(`{"observation":"Terminal request","tool_id":"table-saw"}`), now)
	if err != nil {
		t.Fatalf("create terminal request: %v", err)
	}
	if got := len(app.State().Tools[0].Observations); got != 0 {
		t.Fatalf("unsigned request projected %d observations", got)
	}
	if _, err := app.ApprovalRequest(request.RequestID, "wrong", now); err == nil {
		t.Fatal("accepted wrong polling token")
	}
	approved, err := app.ApproveRequest(request.RequestID, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("approve terminal request: %v", err)
	}
	if approved.State != "approved" || approved.SignedRecordBase64 == "" {
		t.Fatal("approval did not return exact signed record")
	}
	if got := len(app.State().Tools[0].Observations); got != 1 {
		t.Fatalf("approved request projected %d observations", got)
	}
}

func TestTerminalApprovalExpiresRejectsReplayAndRejectsUnknownTarget(t *testing.T) {
	root := t.TempDir()
	identity := &ParticipantIdentity{Label: "alice", Private: testPrivate("alice")}
	app, err := NewPersistentParticipantApp(root, NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": identity.Private.Public().(ed25519.PublicKey)}), identity)
	if err != nil {
		t.Fatal(err)
	}
	bootstrap, _ := testParticipantBootstrap(t, "alice")
	if err := app.IngestRecords(bootstrap); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 11, 22, 10, 0, 0, time.UTC)
	if _, err := app.CreateApprovalRequest("unknown-pcid", []byte(`{}`), now); err == nil {
		t.Fatal("accepted hostile unknown target")
	}
	expired, err := app.CreateApprovalRequest(observationPCID, []byte(`{"observation":"expired","tool_id":"table-saw"}`), now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.ApprovalRequest(expired.RequestID, expired.ApprovalToken, now.Add(11*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if _, err := app.ApproveRequest(expired.RequestID, now.Add(11*time.Minute)); err == nil {
		t.Fatal("approved expired request")
	}
	request, err := app.CreateApprovalRequest(observationPCID, []byte(`{"observation":"once","tool_id":"table-saw"}`), now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.ApproveRequest(request.RequestID, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if _, err := app.ApproveRequest(request.RequestID, now.Add(2*time.Minute)); err == nil {
		t.Fatal("replayed approval request")
	}
}
