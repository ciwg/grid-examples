package web

import (
	"strings"
	"testing"
)

func TestEmbeddedIndexIncludesOperationalWorkflowSections(t *testing.T) {
	index := string(MustRead("index.html"))
	required := []string{
		"Create Place",
		"Create Resource",
		"Live Draft Studio",
		"inventory_audit",
		"Record Run",
	}
	for _, marker := range required {
		if !strings.Contains(index, marker) {
			t.Fatalf("embedded index missing %q", marker)
		}
	}
}

func TestEmbeddedAppIncludesLiveDraftWorkflowHooks(t *testing.T) {
	app := string(MustRead("app.js"))
	required := []string{
		"editor-item-id",
		"/api/items/${editorState.itemID}/live",
		"editor-snapshot",
		"renderPlaces",
		"participant_id",
	}
	for _, marker := range required {
		if !strings.Contains(app, marker) {
			t.Fatalf("embedded app missing %q", marker)
		}
	}
}
