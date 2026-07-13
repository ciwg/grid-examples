package awareness

import "testing"

func TestApplyPrefersNewerState(t *testing.T) {
	t.Parallel()
	index := Index{
		"alice": PeerState{Author: "alice", Lamport: 1, MessageCID: "cid-1"},
	}
	next, applied := Apply(index, Message{
		Kind:        "state",
		DocumentID:  "demo",
		Author:      "alice",
		DisplayName: "Alice",
		Color:       "#123456",
		Cursor:      9,
		Head:        9,
		Typing:      true,
		Lamport:     2,
	}, "cid-2")
	if !applied {
		t.Fatalf("expected awareness update to apply")
	}
	if next["alice"].Cursor != 9 {
		t.Fatalf("unexpected cursor %d", next["alice"].Cursor)
	}
}
