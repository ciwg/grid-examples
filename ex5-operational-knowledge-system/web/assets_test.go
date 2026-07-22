package web

import (
	"strings"
	"testing"
)

func TestEmbeddedIndexIncludesOperationalWorkflowSections(t *testing.T) {
	index := string(MustRead("index.html"))
	required := []string{
		"workspace-status",
		"Review",
		"Author",
		"Operate",
		"Create",
		"Browse Collections",
		"Create Place",
		"Create Resource",
		"Live Draft Studio",
		"Problem Review",
		"Primary Flow",
		"focus-drafts",
		"Current Record",
		"detail-primary",
		"search-clear",
		"search-debug",
		"run-item-select",
		"evidence-run-select",
		"approval-target-select",
		"resource-place-select",
		"detail-timeline",
		"detail-review",
		"search-kind",
		"search-status",
		"search-outcome",
		"Runs",
		"receiving_check",
		"inventory_audit",
		"Log Work Performed",
		"Capture Review Decision",
		"Attach Evidence",
		"Draft Item",
		"Procedure / Checklist",
		"Review This Record",
		"Save Run",
		"Save Evidence",
		"Save Decision",
	}
	for _, marker := range required {
		if !strings.Contains(index, marker) {
			t.Fatalf("embedded index missing %q", marker)
		}
	}
	reviewIndex := strings.Index(index, `workspace workspace-review`)
	createIndex := strings.Index(index, `workspace workspace-create`)
	if reviewIndex == -1 || createIndex == -1 || reviewIndex > createIndex {
		t.Fatalf("embedded index does not keep review ahead of create")
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
		"refreshActionCatalog",
		"renderApprovalTargetOptions",
		"applyContextDefaults",
		"startRunFromContext",
		"startEvidenceFromRun",
		"startApprovalFromContext",
		"renderDetailPrimary",
		"makePrimaryActionButton",
		"workspace-status",
		"setWorkspaceStatus",
		"clearSearch",
		"makeActionButton",
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
		"Record run for this item",
		"Add evidence to this run",
		"Approve this run",
		"Review hotspots",
		"Next Step",
		"Continue draft",
		"problems only",
		"participant_id",
		"safeParticipantStorage",
		"createMemoryStorage",
		"createParticipantID",
		"randomUUID",
		"runHandled",
	}
	for _, marker := range required {
		if !strings.Contains(app, marker) {
			t.Fatalf("embedded app missing %q", marker)
		}
	}
}
