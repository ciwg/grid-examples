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
		"OksClose",
		"/api/items/",
		"/live",
		"related_runs",
		"## Revisions",
		"## Approvals",
		"participant_id",
		"BufWriteCmd",
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
