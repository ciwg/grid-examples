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
		"Problem Review",
		"Record Inspector",
		"detail-timeline",
		"detail-review",
		"search-kind",
		"search-status",
		"search-outcome",
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
		"renderProblemReview",
		"getSearchFilters",
		"buildSearchParams",
		"/api/problem-review",
		"inspectRecord",
		"renderDetailTimeline",
		"renderDetailReview",
		"detailStats",
		"runSearch",
		"Related runs",
		"related_runs",
		"Receiving history",
		"Receiving context review",
		"Receiving review",
		"receivingRunEntries",
		"receivingEvidenceEntries",
		"receivingContextEntries",
		"Inventory count history",
		"Inventory discrepancy",
		"inventoryAuditEntries",
		"inventoryContextEntries",
		"formatEvidenceFacts",
		"Linked runs",
		"detail-json",
		"Search receiving here",
		"Search counts here",
		"Search problems here",
		"problems only",
		"participant_id",
	}
	for _, marker := range required {
		if !strings.Contains(app, marker) {
			t.Fatalf("embedded app missing %q", marker)
		}
	}
}
