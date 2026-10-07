package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
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

func TestLaunchFormSuppliesInteractiveDefaults(t *testing.T) {
	form := NewLaunchForm(Config{})

	if got, want := form.config, (Config{Relay: "http://127.0.0.1:7025", DocumentID: "demo", Name: "Charm User", Color: defaultColor}); got != want {
		t.Fatalf("launch form defaults = %#v, want %#v", got, want)
	}
	if form.form == nil {
		t.Fatal("launch form did not create a Huh form")
	}
}

func TestLaunchFormRetainsProvidedValues(t *testing.T) {
	config := Config{Relay: "http://relay.test", DocumentID: "shared", Name: "Alice", Color: "#22c55e", AccessToken: "secret"}
	form := NewLaunchForm(config)

	if got := form.config; got != config {
		t.Fatalf("launch form config = %#v, want %#v", got, config)
	}
}

func TestLaunchNeedsFormOnlyWhenInteractiveValuesAreMissing(t *testing.T) {
	complete := Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Alice", Color: "#22c55e"}
	if launchNeedsForm(complete) {
		t.Fatal("complete flag configuration unexpectedly opens launch form")
	}
	if !launchNeedsForm(Config{Relay: complete.Relay, DocumentID: complete.DocumentID, Name: complete.Name}) {
		t.Fatal("missing color did not open launch form")
	}
}

func TestRequiredLaunchValue(t *testing.T) {
	if err := requiredLaunchValue("  "); err == nil {
		t.Fatal("blank launch value was accepted")
	}
	if err := requiredLaunchValue("demo"); err != nil {
		t.Fatalf("non-blank launch value rejected: %v", err)
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

func TestOpenedDocumentStartsAtBeginningForEditing(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.applySidecar(sidecarEvent{Type: "opened", DocID: "demo", Content: "first line\nsecond line\nthird line"})
	state.applySidecar(sidecarEvent{Type: "changed", DocID: "demo", Content: "first line\nsecond line\nthird line\nremote update"})

	if got := state.editor.Line(); got != 0 {
		t.Fatalf("remote document update cursor line = %d, want 0", got)
	}
	state.editor, _ = state.editor.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("X")})
	if got := state.editor.Value(); got != "Xfirst line\nsecond line\nthird line\nremote update" {
		t.Fatalf("opened document edit = %q, want insertion at beginning", got)
	}
}

func TestMouseWheelNavigatesSharedDocument(t *testing.T) {
	input := &testWriteCloser{}
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, &sidecar{stdin: input, events: make(chan sidecarEvent)})
	state.editor.SetValue("first line\nsecond line\nthird line\nfourth line")
	if got := state.editor.Line(); got != 3 {
		t.Fatalf("initial editor line = %d, want 3", got)
	}

	updated, _ := state.Update(tea.MouseMsg{Type: tea.MouseWheelUp, Button: tea.MouseButtonWheelUp, Action: tea.MouseActionPress})
	if got := updated.(model).editor.Line(); got != 0 {
		t.Fatalf("wheel up editor line = %d, want 0", got)
	}
}

func TestSidebarDoesNotPushEditorBelowTerminal(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.width = 80
	state.height = 24
	for index := range 67 {
		state.peers = append(state.peers, peer{ID: fmt.Sprintf("peer-%d", index), Name: fmt.Sprintf("Peer %d", index), Color: "#8b5cf6", Embodiment: "browser"})
	}

	if got, limit := lipgloss.Height(state.sidebar()), max(5, state.height-10); got > limit {
		t.Fatalf("sidebar height = %d, want no more than %d", got, limit)
	}
	if !strings.Contains(state.sidebar(), "… 65 more peers") {
		t.Fatalf("sidebar did not summarize hidden peers: %q", state.sidebar())
	}
	if got := lipgloss.Height(state.View()); got > state.height {
		t.Fatalf("full workspace height = %d, want no more than terminal height %d", got, state.height)
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

func TestCatalogSearchUsesExistingMetadataAPI(t *testing.T) {
	requestPath := ""
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestPath = request.URL.RequestURI()
		writer.Header().Set("Content-Type", "application/json")
		if _, err := writer.Write([]byte(`{"results":[]}`)); err != nil {
			t.Errorf("write catalog response: %v", err)
		}
	}))
	defer server.Close()
	state := newModel(Config{Relay: server.URL, DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	message := state.catalogCmd("road map", "open")()
	result, ok := message.(actionResultMsg)
	if !ok || result.panelKind != panelResults || result.panelAction != "open" {
		t.Fatalf("catalog result = %#v", message)
	}
	if requestPath != "/api/local/metadata/search?q=road+map" {
		t.Fatalf("catalog path = %q", requestPath)
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

func TestTypingPresenceUsesStableMarker(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.peers = []peer{{Name: "Alice", Color: "#ff0000", Typing: true}}
	if !strings.Contains(state.sidebar(), "● Alice") {
		t.Fatalf("typing sidebar did not render stable marker: %q", state.sidebar())
	}
}

func TestWorkspaceInitWaitsForSidecarWithoutAnimation(t *testing.T) {
	events := make(chan sidecarEvent)
	close(events)
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, &sidecar{events: events})

	message := state.Init()()
	if _, ok := message.(sidecarMsg); !ok {
		t.Fatalf("workspace startup message = %T, want sidecarMsg without animation batch", message)
	}
}

func TestHelpMenuOpensBubblesKeyBindingPanel(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.width, state.height = 100, 30
	state.activeMenu, state.activeItem, state.menuOpen = 6, 0, true
	state.activate()

	if state.panel == nil || !state.panel.hasHelp {
		t.Fatalf("Help menu panel = %#v, want Bubbles help panel", state.panel)
	}
	if got := state.panel.view(); !strings.Contains(got, "Markdown preview") {
		t.Fatalf("help panel did not render Bubbles key bindings: %q", got)
	}
}

func TestDocumentAndRelayResultsUseBubblesListAndTable(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.width, state.height = 100, 30
	updated, _ := state.Update(actionResultMsg{status: "Relay document catalog", content: `{"documents":[{"document_id":"alpha","title":"Alpha notes"}]}`, panelKind: panelResults, panelAction: "open"})
	result := updated.(model)

	if result.panel == nil || !result.panel.hasList || !result.panel.hasTable {
		t.Fatalf("catalog panel = %#v, want Bubbles list and table", result.panel)
	}
	if got := result.panel.selectedValue(); got != "alpha" {
		t.Fatalf("selected document = %q, want alpha", got)
	}
}

func TestMetadataSearchResultsPopulateDocumentChooser(t *testing.T) {
	items, rows := resultRows(`{"query":"","results":[{"document_id":"alpha","title":"Alpha notes"}]}`)
	if len(items) != 1 || items[0].value != "alpha" {
		t.Fatalf("metadata items = %#v, want alpha document", items)
	}
	if len(rows) != 1 || rows[0][0] != "alpha" {
		t.Fatalf("metadata rows = %#v, want alpha document", rows)
	}
}

func TestImportMenuOpensBubblesFilePicker(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.width, state.height = 100, 30
	state.activeMenu, state.activeItem, state.menuOpen = 0, 3, true
	state.activate()

	if state.panel == nil || !state.panel.hasPicker {
		t.Fatalf("import panel = %#v, want Bubbles file picker", state.panel)
	}
	if state.dialog != nil {
		t.Fatal("import retained the absolute-path text dialog")
	}
}

func TestActivityMenuUsesScrollableBubblesViewport(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.width, state.height = 100, 30
	state.activity = []string{"first relay event", "second relay event"}
	state.activeMenu, state.activeItem, state.menuOpen = 3, 2, true
	state.activate()

	if state.panel == nil || !state.panel.hasViewport {
		t.Fatalf("activity panel = %#v, want Bubbles viewport", state.panel)
	}
	if got := state.panel.view(); !strings.Contains(got, "second relay event") {
		t.Fatalf("activity panel did not render local history: %q", got)
	}
}

func TestActivityPanelOwnsMouseWheelInsteadOfSharedEditor(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.width, state.height = 100, 20
	state.editor.SetValue("one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten")
	state.activity = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	state.panel = newActivityPanel(state.activity, state.width, state.height)
	beforeLine := state.editor.Line()

	updated, _ := state.Update(tea.MouseMsg{Type: tea.MouseWheelDown, Button: tea.MouseButtonWheelDown, Action: tea.MouseActionPress})
	result := updated.(model)
	if got := result.editor.Line(); got != beforeLine {
		t.Fatalf("panel wheel changed editor line = %d, want %d", got, beforeLine)
	}
}

func TestBusyCommandStartsSpinnerOnlyForActiveWork(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	if command := state.busyCmd(nil); command != nil || state.busy {
		t.Fatal("nil work unexpectedly started the spinner")
	}
	command := state.busyCmd(func() tea.Msg { return actionResultMsg{status: "done"} })
	if command == nil || !state.busy {
		t.Fatal("active work did not start the Bubbles spinner")
	}
}

func TestRunDialogWrapsExportInBusySpinner(t *testing.T) {
	state := newModel(Config{Relay: "http://relay.test", DocumentID: "demo", Name: "Charm User", Color: defaultColor}, nil)
	state.runDialog(&dialog{action: "export:txt", input: textinput.New()})
	// A blank export path intentionally does not begin a write or spinner.
	if state.busy {
		t.Fatal("blank export unexpectedly started spinner")
	}
	dialogInput := textinput.New()
	dialogInput.SetValue("/tmp/grid-tui-export.txt")
	state.runDialog(&dialog{action: "export:txt", input: dialogInput})
	if !state.busy {
		t.Fatal("export work did not start spinner")
	}
}
