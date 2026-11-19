package awareness

import (
	"testing"
	"time"
)

func TestApplyPrefersNewerState(t *testing.T) {
	t.Parallel()
	index := Index{
		"alice-browser": PeerState{Author: "alice", ParticipantID: "alice-browser", Lamport: 1, MessageCID: "cid-1"},
	}
	next, applied := Apply(index, Message{
		Kind:          "state",
		DocumentID:    "demo",
		Author:        "alice",
		ParticipantID: "alice-browser",
		DisplayName:   "Alice",
		Color:         "#123456",
		Cursor:        9,
		Head:          9,
		Typing:        true,
		Lamport:       2,
	}, "cid-2", time.Unix(10, 0))
	if !applied {
		t.Fatalf("expected awareness update to apply")
	}
	if next["alice-browser"].Cursor != 9 {
		t.Fatalf("unexpected cursor %d", next["alice-browser"].Cursor)
	}
}

func TestApplyRejectsStaleState(t *testing.T) {
	t.Parallel()
	index := Index{
		"alice-browser": {
			Author:        "alice",
			ParticipantID: "alice-browser",
			DisplayName:   "Alice",
			Color:         "#123456",
			Cursor:        9,
			Head:          9,
			Typing:        true,
			Lamport:       4,
			MessageCID:    "cid-4",
		},
	}
	next, applied := Apply(index, Message{
		Kind:          "state",
		DocumentID:    "demo",
		Author:        "alice",
		ParticipantID: "alice-browser",
		DisplayName:   "Alice old",
		Color:         "#654321",
		Cursor:        2,
		Head:          2,
		Typing:        false,
		Lamport:       3,
	}, "cid-3", time.Unix(11, 0))
	if applied {
		t.Fatalf("expected stale awareness update to be rejected")
	}
	if next["alice-browser"].Cursor != 9 {
		t.Fatalf("unexpected cursor after stale update %d", next["alice-browser"].Cursor)
	}
}

func TestApplyFallsBackToAuthorWhenParticipantMissing(t *testing.T) {
	t.Parallel()
	next, applied := Apply(nil, Message{
		Kind:        "state",
		DocumentID:  "demo",
		Author:      "alice",
		DisplayName: "Alice",
		Color:       "#abcdef",
		Cursor:      5,
		Head:        7,
		Typing:      true,
		Lamport:     1,
	}, "cid-1", time.Unix(12, 0))
	if !applied {
		t.Fatalf("expected awareness update to apply")
	}
	state, ok := next["alice"]
	if !ok {
		t.Fatalf("expected author fallback key")
	}
	if state.Head != 7 {
		t.Fatalf("unexpected head %d", state.Head)
	}
}

func TestApplyBreaksLamportTiesDeterministically(t *testing.T) {
	t.Parallel()
	index := Index{
		"shared": {
			Author:        "alice",
			ParticipantID: "shared",
			DisplayName:   "Alice",
			Color:         "#111111",
			Cursor:        1,
			Head:          1,
			Typing:        false,
			Lamport:       2,
			MessageCID:    "cid-a",
		},
	}
	next, applied := Apply(index, Message{
		Kind:          "state",
		DocumentID:    "demo",
		Author:        "bob",
		ParticipantID: "shared",
		DisplayName:   "Bob",
		Color:         "#222222",
		Cursor:        8,
		Head:          8,
		Typing:        true,
		Lamport:       2,
	}, "cid-b", time.Unix(13, 0))
	if !applied {
		t.Fatalf("expected equal-lamport update with higher author tie-break to apply")
	}
	if next["shared"].Author != "bob" {
		t.Fatalf("unexpected winner %q", next["shared"].Author)
	}
}
