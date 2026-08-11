package service

import (
	"os"
	"testing"
)

func TestStoreReplaysExactValidatedRecordFrames(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	record, private := testRecord(t)
	_, raw, err := record.Sign(private)
	if err != nil {
		t.Fatalf("sign record: %v", err)
	}
	if err := store.AppendRecords([][]byte{raw}); err != nil {
		t.Fatalf("append record: %v", err)
	}
	frames, err := store.ReadRecordFrames()
	if err != nil {
		t.Fatalf("replay frames: %v", err)
	}
	if len(frames) != 1 || len(frames[0]) != 1 || string(frames[0][0]) != string(raw) {
		t.Fatalf("replayed frames = %#v", frames)
	}
}

func TestStoreRejectsInvalidRecordBeforeWrite(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	if err := store.AppendRecords([][]byte{[]byte("not a Grid record")}); err == nil {
		t.Fatal("accepted invalid record")
	}
	if _, err := os.Stat(store.framePath); !os.IsNotExist(err) {
		t.Fatalf("record store exists after rejected append: %v", err)
	}
}

func TestStoreFailsClosedOnPartialRecordFrame(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	if err := os.WriteFile(store.framePath, []byte("MSR1\n\x00\x00\x00"), 0o600); err != nil {
		t.Fatalf("write partial frame: %v", err)
	}
	if _, err := store.ReadRecordFrames(); err == nil {
		t.Fatal("replayed partial record frame")
	}
}
