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
		"Record Inspector",
		"detail-timeline",
		"detail-review",
		"search-kind",
		"search-status",
		"Runs",
		"receiving_check",
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
		"renderSearchResults",
		"getSearchFilters",
		"buildSearchParams",
		"inspectRecord",
		"renderDetailTimeline",
		"renderDetailReview",
		"detailStats",
		"Related runs",
		"related_runs",
		"Receiving history",
		"Receiving review",
		"receivingRunEntries",
		"receivingEvidenceEntries",
		"Inventory audit history",
		"Inventory discrepancy",
		"inventoryAuditEntries",
		"formatEvidenceFacts",
		"Linked runs",
		"detail-json",
		"participant_id",
	}
	for _, marker := range required {
		if !strings.Contains(app, marker) {
			t.Fatalf("embedded app missing %q", marker)
		}
	}
}
