package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func TestUTF16LengthAndLineColumn(t *testing.T) {
	text := "A😀\nBeta"
	if got := utf16Length("A😀"); got != 3 {
		t.Fatalf("utf16Length() = %d, want 3", got)
	}
	if got := lineColumn(text, 3); got != "1:4" {
		t.Fatalf("lineColumn() = %q, want 1:4", got)
	}
	if got := lineColumn(text, 4); got != "2:1" {
		t.Fatalf("lineColumn() = %q, want 2:1", got)
	}
}

func TestPreviewShortcutsToggleTheVisibleMode(t *testing.T) {
	state := newModel("http://relay.test", "demo", "Charm User", "#8b5cf6", nil)
	updated, _ := state.Update(tea.KeyMsg{Type: tea.KeyCtrlP})
	if !updated.(model).preview {
		t.Fatal("ctrl+p did not enable Markdown preview")
	}
}

func TestEditorUTF16OffsetTracksBubbleTextareaCursor(t *testing.T) {
	editor := textarea.New()
	editor.SetValue("A😀\nBeta")
	editor.SetCursor(4)
	editor.CursorDown()
	editor.SetCursor(2)
	if got := editorUTF16Offset(editor); got != 6 {
		t.Fatalf("editorUTF16Offset() = %d, want 6", got)
	}
}

func TestRenderLineMarksRemoteCaretWithoutEmbeddingName(t *testing.T) {
	rendered := renderLine("hello", []peer{{Name: "Alice", Color: "#ff0000", Anchor: 5, Head: 5}}, 0, "hello", -1)
	if strings.Contains(rendered, "Alice") {
		t.Fatalf("peer name was embedded in editor output: %q", rendered)
	}
	if rendered == "hello" {
		t.Fatal("remote cursor did not alter rendered output")
	}
}
