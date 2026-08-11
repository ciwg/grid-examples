package service

import (
	"crypto/ed25519"
	"crypto/sha256"
	"testing"
)

func signedTestRecord(t *testing.T, protocol, id, signer, payload string) ([]byte, ed25519.PublicKey) {
	t.Helper()
	seed := sha256.Sum256([]byte("ex7 record signer: " + signer))
	private := ed25519.NewKeyFromSeed(seed[:])
	public := private.Public().(ed25519.PublicKey)
	record := Record{Protocol: protocol, ID: id, Signer: signer, CreatedAt: "2026-08-11T22:10:00Z", Payload: []byte(payload), KeyID: keyID(public), PublicKey: public}
	_, raw, err := record.Sign(private)
	if err != nil {
		t.Fatalf("sign record: %v", err)
	}
	return raw, public
}

func TestPersistentRecordAppReplaysObservationAndSafetyDisposition(t *testing.T) {
	observation, aliceKey := signedTestRecord(t, observationPCID, "obs-1", "alice", `{"observation":"Guard is loose","tool_id":"table-saw"}`)
	hold, carolKey := signedTestRecord(t, safetyPCID, "hold-1", "carol", `{"assessment":"Guard is loose","basis_record_id":"obs-1","disposition":"hold","tool_id":"table-saw"}`)
	policy := NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": aliceKey, "carol": carolKey})
	root := t.TempDir()
	app, err := NewPersistentRecordApp(root, policy)
	if err != nil {
		t.Fatalf("create persistent record app: %v", err)
	}
	if err := app.IngestRecords([][]byte{observation, hold}); err != nil {
		t.Fatalf("ingest evidence frame: %v", err)
	}
	reloaded, err := NewPersistentRecordApp(root, policy)
	if err != nil {
		t.Fatalf("replay record app: %v", err)
	}
	tool := reloaded.State().Tools[0]
	if !tool.SafetyHold || len(tool.Observations) != 1 {
		t.Fatalf("replayed tool = %+v", tool)
	}
}

func TestPersistentRecordAppRetainsUnrecognizedEvidenceWithoutProjection(t *testing.T) {
	record, _ := signedTestRecord(t, observationPCID, "obs-1", "mallory", `{"observation":"Untrusted claim","tool_id":"table-saw"}`)
	root := t.TempDir()
	app, err := NewPersistentRecordApp(root, RecognitionPolicy{})
	if err != nil {
		t.Fatalf("create persistent record app: %v", err)
	}
	if err := app.IngestRecords([][]byte{record}); err != nil {
		t.Fatalf("retain unrecognized record: %v", err)
	}
	if got := len(app.State().Tools[0].Observations); got != 0 {
		t.Fatalf("unrecognized record changed projection: %d observations", got)
	}
	frames, err := app.store.ReadRecordFrames()
	if err != nil {
		t.Fatalf("read retained evidence: %v", err)
	}
	if len(frames) != 1 || string(frames[0][0]) != string(record) {
		t.Fatal("unrecognized record was not retained exactly")
	}
}

func TestPersistentRecordAppReplaysUnknownPCIDWithoutProjection(t *testing.T) {
	unknownPCID := rawCIDv1([]byte("ex7 unknown family test specification"))
	record, aliceKey := signedTestRecord(t, unknownPCID, "unknown-1", "alice", `{}`)
	root := t.TempDir()
	policy := NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": aliceKey})
	app, err := NewPersistentRecordApp(root, policy)
	if err != nil {
		t.Fatalf("create persistent record app: %v", err)
	}
	if err := app.IngestRecords([][]byte{record}); err != nil {
		t.Fatalf("retain unknown family: %v", err)
	}
	reloaded, err := NewPersistentRecordApp(root, policy)
	if err != nil {
		t.Fatalf("replay unknown family: %v", err)
	}
	if got := len(reloaded.State().Tools[0].Observations); got != 0 {
		t.Fatalf("unknown pCID changed projection: %d observations", got)
	}
	frames, err := reloaded.store.ReadRecordFrames()
	if err != nil || len(frames) != 1 || string(frames[0][0]) != string(record) {
		t.Fatalf("unknown pCID replay = %#v, %v", frames, err)
	}
}

func TestPersistentRecordAppDoesNotApplyRoleMismatchedRecognizedKey(t *testing.T) {
	hold, aliceKey := signedTestRecord(t, safetyPCID, "hold-1", "alice", `{"assessment":"Guard is loose","disposition":"hold","tool_id":"table-saw"}`)
	clear, _ := signedTestRecord(t, safetyPCID, "clear-1", "alice", `{"assessment":"Looks safe","disposition":"clear","tool_id":"table-saw"}`)
	app, err := NewPersistentRecordApp(t.TempDir(), NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": aliceKey}))
	if err != nil {
		t.Fatalf("create persistent record app: %v", err)
	}
	if err := app.IngestRecords([][]byte{hold, clear}); err != nil {
		t.Fatalf("ingest role-mismatched records: %v", err)
	}
	if !app.State().Tools[0].SafetyHold {
		t.Fatal("recognized non-steward key cleared a safety hold")
	}
}

func TestPersistentRecordAppRejectsMalformedKnownPayloadBeforeWrite(t *testing.T) {
	record, aliceKey := signedTestRecord(t, observationPCID, "obs-1", "alice", `{"observation":"Missing tool"}`)
	app, err := NewPersistentRecordApp(t.TempDir(), NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": aliceKey}))
	if err != nil {
		t.Fatalf("create persistent record app: %v", err)
	}
	if err := app.IngestRecords([][]byte{record}); err == nil {
		t.Fatal("accepted malformed known-family payload")
	}
	if _, err := app.store.ReadRecordFrames(); err != nil {
		t.Fatalf("read empty record history: %v", err)
	}
}

func TestPersistentRecordAppProjectsLoanAndLinkedReturn(t *testing.T) {
	loan, aliceKey := signedTestRecord(t, loanPCID, "loan-1", "alice", `{"borrower_id":"alice","due_at":"2030-01-02T15:04:05Z","policy":"Return with charger","policy_version":"v1","tool_id":"cordless-drill"}`)
	returned, _ := signedTestRecord(t, returnPCID, "return-1", "alice", `{"condition":"Returned with charger","loan_record_id":"loan-1","tool_id":"cordless-drill"}`)
	app, err := NewPersistentRecordApp(t.TempDir(), NewRecognitionPolicy(map[string]ed25519.PublicKey{"alice": aliceKey}))
	if err != nil {
		t.Fatalf("create persistent record app: %v", err)
	}
	if err := app.IngestRecords([][]byte{loan}); err != nil {
		t.Fatalf("ingest loan: %v", err)
	}
	if app.State().Tools[1].ActiveLoan == nil {
		t.Fatal("loan was not projected")
	}
	if err := app.IngestRecords([][]byte{returned}); err != nil {
		t.Fatalf("ingest return: %v", err)
	}
	if app.State().Tools[1].ActiveLoan != nil {
		t.Fatal("linked return did not clear active loan")
	}
}
