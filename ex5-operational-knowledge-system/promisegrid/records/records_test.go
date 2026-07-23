package records

import (
	"path/filepath"
	"testing"
)

func TestBuildAndVerifySignedKnowledgeItemRecord(t *testing.T) {
	identity, err := LoadOrCreateRuntimeIdentity(filepath.Join(t.TempDir(), "identity.seed"))
	if err != nil {
		t.Fatalf("load identity: %v", err)
	}
	event := Event{
		Sequence:          7,
		Timestamp:         "2026-07-22T19:38:06Z",
		EntityID:          "ITEM-0007",
		Type:              "knowledge_item_created",
		Actor:             "alice",
		Kind:              "procedure",
		Status:            "draft",
		Title:             "Title",
		Summary:           "Summary",
		Body:              "Body",
		Tags:              []string{"tag"},
		ResponsibilityIDs: []string{"RESP-1"},
		Revision:          1,
	}
	record, ok, err := BuildSignedKnowledgeItemRecord(identity, event)
	if err != nil {
		t.Fatalf("build record: %v", err)
	}
	if !ok {
		t.Fatalf("expected record")
	}
	if err := VerifySignedKnowledgeItemRecords([]Event{NormalizeEvent(event, identity.PeerID())}, []SignedKnowledgeItemRecord{record}); err != nil {
		t.Fatalf("verify record: %v", err)
	}
}

func TestDecoratePeerVisibleEventCanonicalIDs(t *testing.T) {
	events := []Event{{
		Sequence:       3,
		OriginPeerID:   "peer-a",
		OriginSequence: 3,
		EntityID:       "ITEM-0003",
		Type:           "knowledge_item_created",
	}}
	records := []SignedKnowledgeItemRecord{{
		Sequence:       3,
		OriginPeerID:   "peer-a",
		OriginSequence: 3,
		ItemID:         "ITEM-0003",
		EventType:      "knowledge_item_created",
		EnvelopeCID:    "bafy-test-cid",
	}}
	decorated := DecoratePeerVisibleEventCanonicalIDs(events, records, nil, nil, nil, nil, nil, nil, nil)
	if got := decorated[0].CanonicalID; got != "bafy-test-cid" {
		t.Fatalf("canonical id = %q", got)
	}
	if got := decorated[0].DisplayID; got != "ITEM-0003" {
		t.Fatalf("display id = %q", got)
	}
}
