package transport

import (
	"reflect"
	"testing"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

func TestFilterRelayFeedEventsUsesOriginCursor(t *testing.T) {
	events := []records.Event{
		{OriginPeerID: "alice", OriginSequence: 1, AttachmentCID: "cid-a"},
		{OriginPeerID: "alice", OriginSequence: 2},
		{OriginPeerID: "bob", OriginSequence: 1},
	}
	gotEvents, gotWanted := FilterRelayFeedEvents(events, map[string]uint64{"alice": 1})
	if len(gotEvents) != 2 {
		t.Fatalf("expected 2 unseen events, got %d", len(gotEvents))
	}
	if gotEvents[0].OriginPeerID != "alice" || gotEvents[0].OriginSequence != 2 {
		t.Fatalf("unexpected first unseen event: %+v", gotEvents[0])
	}
	if !gotWanted["alice#2"] || !gotWanted["bob#1"] {
		t.Fatalf("unexpected wanted origin map: %+v", gotWanted)
	}
}

func TestRequiredBlobCIDsForEventsDedupesAndSorts(t *testing.T) {
	events := []records.Event{
		{Type: "evidence_added", AttachmentCID: "cid-b"},
		{Type: "evidence_added", AttachmentCID: "cid-a"},
		{Type: "evidence_added", AttachmentCID: "cid-b"},
		{Type: "run_recorded", AttachmentCID: "cid-z"},
	}
	got := RequiredBlobCIDsForEvents(events)
	want := []string{"cid-a", "cid-b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("required blob cids mismatch: got %v want %v", got, want)
	}
}

func TestFilterKnowledgeItemRecordsByOrigin(t *testing.T) {
	in := []records.SignedKnowledgeItemRecord{
		{OriginPeerID: "alice", OriginSequence: 1, Sequence: 7},
		{OriginPeerID: "alice", OriginSequence: 2, Sequence: 8},
	}
	got := FilterKnowledgeItemRecordsByOrigin(in, map[string]bool{"alice#2": true})
	if len(got) != 1 || got[0].OriginSequence != 2 {
		t.Fatalf("unexpected filtered records: %+v", got)
	}
}
