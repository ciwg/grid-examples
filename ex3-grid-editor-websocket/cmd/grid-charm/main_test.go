package main

import (
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/service"
)

func TestModelFiltersRelayObservedEntries(t *testing.T) {
	model := newModel("http://relay.test", "demo")
	entries := []service.TraceEntry{{Protocol: "live-document", Kind: "change"}, {Protocol: "live-awareness", Kind: "update"}}
	model.familyIndex = 1
	model.setEntries(entries)
	if got := len(model.list.Items()); got != 1 {
		t.Fatalf("family filter kept %d entries, want 1", got)
	}
	model.familyIndex = 0
	model.kindIndex = 2
	model.setEntries(entries)
	if got := len(model.list.Items()); got != 1 {
		t.Fatalf("kind filter kept %d entries, want 1", got)
	}
}

func TestShortID(t *testing.T) {
	if got := shortID("abcdefghijklmnopqrstuvwxyz"); got != "abcdefghijklmnopqr…" {
		t.Fatalf("shortID() = %q", got)
	}
}
