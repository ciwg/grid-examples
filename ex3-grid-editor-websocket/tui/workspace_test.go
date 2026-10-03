package tui

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type testWriteCloser struct{ bytes.Buffer }

func (writer *testWriteCloser) Close() error { return nil }

func TestSidecarConnectMessageIncludesRemoteAccessToken(t *testing.T) {
	message := sidecarConnectMessage("http://relay.example", "bootstrap-token", "charm-a", "Charm A", "#8b5cf6")

	if got, want := message["access_token"], any("bootstrap-token"); got != want {
		t.Fatalf("access token mismatch: got %q want %q", got, want)
	}
	if got, want := message["embodiment"], any("charm"); got != want {
		t.Fatalf("embodiment mismatch: got %q want %q", got, want)
	}
}

func TestEditorAcceptsTextAndSendsSharedDocumentUpdate(t *testing.T) {
	input := &testWriteCloser{}
	client := &sidecar{stdin: input, events: make(chan sidecarEvent)}
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, client)

	updated, _ := state.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	edited := updated.(model)
	if got := edited.editor.Value(); got != "a" {
		t.Fatalf("editor value = %q, want typed text", got)
	}
	lines := bytes.Split(bytes.TrimSpace(input.Bytes()), []byte("\n"))
	if len(lines) < 1 {
		t.Fatal("editor did not send a sidecar message")
	}
	var message map[string]any
	if err := json.Unmarshal(lines[0], &message); err != nil {
		t.Fatalf("decode sidecar message: %v", err)
	}
	if got, want := message["type"], any("set_text"); got != want {
		t.Fatalf("sidecar message type = %q, want %q", got, want)
	}
	if got, want := message["content"], any("a"); got != want {
		t.Fatalf("sidecar content = %q, want %q", got, want)
	}
}

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

func TestMenuKeyboardNavigationOpensDocumentMenu(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	updated, _ := state.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d"), Alt: true})
	menu := updated.(model)
	if !menu.menuOpen || menu.activeMenu != 0 {
		t.Fatalf("Alt+D did not open Document menu: %#v", menu)
	}
	updated, _ = menu.Update(tea.KeyMsg{Type: tea.KeyDown})
	if got := updated.(model).activeItem; got != 1 {
		t.Fatalf("down menu item = %d, want 1", got)
	}
}

func TestMouseClickUsesRenderedMenuRowAndWidths(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	viewMenu := state.menuAt(lipgloss.Width(menuStyle.Render("Document")) + 2)
	updated, _ := state.Update(tea.MouseMsg{X: lipgloss.Width(menuStyle.Render("Document")) + 2, Y: menuRow, Button: tea.MouseButtonLeft, Action: tea.MouseActionPress})
	clicked := updated.(model)
	if !clicked.menuOpen || clicked.activeMenu != viewMenu || clicked.activeMenu != 1 {
		t.Fatalf("mouse menu click = open:%t menu:%d, want Edit", clicked.menuOpen, clicked.activeMenu)
	}
}

func TestViewMenuChangesVisibleMode(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.activeMenu, state.activeItem, state.menuOpen = 2, 1, true
	state.activate()
	if state.mode != "preview" {
		t.Fatalf("preview menu mode = %q, want preview", state.mode)
	}
	state.activeItem, state.menuOpen = 2, true
	state.activate()
	if state.mode != "split" {
		t.Fatalf("split menu mode = %q, want split", state.mode)
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

func TestExportRequiresSecondConfirmation(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.openDialog("Export Markdown", "Absolute output path", "export:md", true)
	state.dialog.input.SetValue("/tmp/grid-tui-test.md")
	updated, _ := state.updateDialog(tea.KeyMsg{Type: tea.KeyEnter})
	confirmed := updated.(model)
	if confirmed.dialog == nil || confirmed.dialog.confirm {
		t.Fatal("first export Enter did not move to explicit confirmation")
	}
}

func TestRelayTraceMenuUsesExistingEx3HTTPPath(t *testing.T) {
	requestPath := ""
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestPath = request.URL.RequestURI()
		writer.Header().Set("Content-Type", "application/json")
		if _, err := writer.Write([]byte(`{"entries":[]}`)); err != nil {
			t.Errorf("write trace response: %v", err)
		}
	}))
	defer server.Close()
	state := newModel(Config{Relay: server.URL, DocumentID: "shared doc", Name: "Charm User", Color: defaultColor}, nil)
	message := state.traceCmd()()
	result, ok := message.(actionResultMsg)
	if !ok || result.status != "Relay trace refreshed" {
		t.Fatalf("trace result = %#v", message)
	}
	if requestPath != "/api/local/documents/shared%20doc/trace?limit=12" {
		t.Fatalf("trace path = %q", requestPath)
	}
}

func TestGlowMissingReportsInstallLocation(t *testing.T) {
	t.Setenv("PATH", "")
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	message := state.glowCmd("README.md")()
	result, ok := message.(actionResultMsg)
	if !ok || !strings.Contains(result.status, "https://charm.land/glow") {
		t.Fatalf("Glow missing result = %#v", message)
	}
}

func TestHarmonicaPresencePulseChangesTypingMarker(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.peers = []peer{{Name: "Alice", Color: "#ff0000", Typing: true}}
	for range 20 {
		updated, _ := state.Update(presenceTickMsg{})
		state = updated.(model)
	}
	if state.presencePulse <= 0.5 {
		t.Fatalf("Harmonica pulse = %f, want visible typing pulse", state.presencePulse)
	}
	if !strings.Contains(state.sidebar(), "● Alice") {
		t.Fatalf("typing sidebar did not render pulse: %q", state.sidebar())
	}
}
