package nvim

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNeovimPluginIncludesLiveDraftCommands(t *testing.T) {
	root := filepath.Join("lua", "oks", "init.lua")
	body, err := os.ReadFile(root)
	if err != nil {
		t.Fatalf("read plugin: %v", err)
	}
	required := []string{
		"OksOpen",
		"OksRefresh",
		"OksPush",
		"OksInfo",
		"OksInspect",
		"OksInspectRun",
		"OksInspectEntity",
		"OksSearch",
		"OksPending",
		"OksApproveItem",
		"OksApproveRun",
		"OksSupersedeItem",
		"OksClose",
		"/api/items/",
		"/api/runs/",
		"/approvals",
		"/api/search?q=",
		"/api/places/",
		"/api/resources/",
		"/api/responsibilities/",
		"/live",
		"## Links",
		"linked_item_ids",
		"related_runs",
		"## Revisions",
		"## Approvals",
		"## Evidence",
		"## Filters",
		"## Draft items",
		"## Unreviewed runs",
		"## Problem runs",
		"inspect: :OksInspectRun",
		"revision = 0",
		"current_revision",
		"item supersede failed",
		"participant_id",
		"BufWriteCmd",
		"live_draft_winid",
		"nvim_list_wins",
		"nvim_win_get_buf",
		"wipe_buffer",
	}
	text := string(body)
	for _, marker := range required {
		if !strings.Contains(text, marker) {
			t.Fatalf("plugin missing %q", marker)
		}
	}
}

func TestNeovimLauncherOpensItemThroughRuntimePath(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "scripts", "oks-nvim"))
	if err != nil {
		t.Fatalf("read launcher: %v", err)
	}
	required := []string{
		"runtimepath+=",
		"OksOpen",
		"OKS_BASE_URL",
		"OKS_DISPLAY_NAME",
		"OKS_COLOR",
	}
	text := string(body)
	for _, marker := range required {
		if !strings.Contains(text, marker) {
			t.Fatalf("launcher missing %q", marker)
		}
	}
}
