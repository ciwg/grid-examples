package document

import "testing"

func TestApplyPrefersHigherLamport(t *testing.T) {
	t.Parallel()
	current := State{Lamport: 1, Author: "a", MessageCID: "cid-1"}
	message, err := NewMessage("demo", "hello", 2, "b", "test", "")
	if err != nil {
		t.Fatalf("new message: %v", err)
	}
	next, applied, err := Apply(current, message, "cid-2")
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if !applied {
		t.Fatalf("expected message to apply")
	}
	if next.Lamport != 2 {
		t.Fatalf("unexpected lamport %d", next.Lamport)
	}
}
