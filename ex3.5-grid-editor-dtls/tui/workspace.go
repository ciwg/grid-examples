// Package tui implements Grid's terminal collaboration workspace.
package tui

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf16"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	charmLog "github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

// Intent: Keep Charm's terminal presentation separate from Ex3's relay and
// Automerge contracts. Source: DI-mutoh.
const defaultColor = "#8b5cf6"

var (
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	mutedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	panelStyle      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(0, 1)
	menuStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	activeMenuStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("86")).Padding(0, 1)
	// Intent: Keep the active terminal user's editing position unmistakable on
	// light and dark terminals; peer colors remain reserved for remote awareness.
	// Source: DI-holoz.
	localCursorStyle = lipgloss.NewStyle().Background(lipgloss.Color("#facc15")).Foreground(lipgloss.Color("0")).Bold(true)
)

// Config is the public launcher contract for the terminal embodiment.
type Config struct{ Relay, DocumentID, Name, Color, AccessToken string }

// LaunchForm gathers a human operator's local presentation choices before the
// terminal embodiment joins a shared Ex3 document.
type LaunchForm struct {
	config Config
	form   *huh.Form
}

// NewLaunchForm creates the interactive launch step. Explicit CLI values are
// retained so a partially specified command only asks for what remains.
//
// Intent: Make interactive Grid TUI startup discoverable without changing the
// existing Automerge, WebSocket, or relay contracts. Source: DI-sidob.
func NewLaunchForm(config Config) *LaunchForm {
	config = launchDefaults(config)
	launch := &LaunchForm{config: config}
	launch.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Display name").
				Description("The name collaborators see in the peer legend.").
				Value(&launch.config.Name).
				Validate(requiredLaunchValue),
			huh.NewSelect[string]().
				Title("Peer color").
				Description("The color used for your cursor and presence.").
				Options(
					huh.NewOption("Purple", defaultColor),
					huh.NewOption("Blue", "#3b82f6"),
					huh.NewOption("Green", "#22c55e"),
					huh.NewOption("Orange", "#f97316"),
					huh.NewOption("Pink", "#ec4899"),
				).
				Value(&launch.config.Color),
			huh.NewInput().
				Title("Relay URL").
				Description("The Ex3 relay that carries collaboration updates.").
				Value(&launch.config.Relay).
				Validate(requiredLaunchValue),
			huh.NewInput().
				Title("Document ID").
				Description("The collaborative document to open.").
				Value(&launch.config.DocumentID).
				Validate(requiredLaunchValue),
		).
			Title("Grid TUI — Join collaboration session").
			Description("Choose your local identity, then join the shared document."),
	)
	return launch
}

// Run displays the form and returns the submitted launch configuration.
func (launch *LaunchForm) Run() (Config, error) {
	if err := launch.form.Run(); err != nil {
		return Config{}, err
	}
	return launch.config, nil
}

func launchDefaults(config Config) Config {
	if config.Relay == "" {
		config.Relay = "http://127.0.0.1:7025"
	}
	if config.DocumentID == "" {
		config.DocumentID = "demo"
	}
	if config.Name == "" {
		config.Name = "Charm User"
	}
	if config.Color == "" {
		config.Color = defaultColor
	}
	return config
}

func launchNeedsForm(config Config) bool {
	return config.Relay == "" || config.DocumentID == "" || config.Name == "" || config.Color == ""
}

func requiredLaunchValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("required")
	}
	return nil
}

type peer struct {
	ID         string `json:"participant_id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Anchor     int    `json:"anchor"`
	Head       int    `json:"head"`
	Typing     bool   `json:"typing"`
	Embodiment string `json:"embodiment"`
}
type sidecarEvent struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	Peers     []peer `json:"peers"`
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
	Replica   string `json:"replica_base64"`
	DocID     string `json:"doc_id"`
}
type ioWriteCloser interface {
	Write([]byte) (int, error)
	Close() error
}
type ioReader interface{ Read([]byte) (int, error) }

type sidecar struct {
	command *exec.Cmd
	stdin   ioWriteCloser
	events  chan sidecarEvent
	mu      sync.Mutex
	logger  *charmLog.Logger
}

func sidecarConnectMessage(relay, accessToken, participantID, name, color string) map[string]any {
	return map[string]any{
		"type":           "connect",
		"relay_url":      relay,
		"access_token":   accessToken,
		"participant_id": participantID,
		"display_name":   name,
		"color":          color,
		"embodiment":     "charm",
	}
}

func startSidecar(relay, accessToken, participantID, name, color string, logger *charmLog.Logger) (*sidecar, error) {
	command := exec.Command("go", "run", "./cmd/grid-nvim-sidecar", "--relay", relay)
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("sidecar stdin: %w", err)
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("sidecar stdout: %w", err)
	}
	command.Stderr = os.Stderr
	if err := command.Start(); err != nil {
		return nil, fmt.Errorf("start Automerge sidecar: %w", err)
	}
	client := &sidecar{command: command, stdin: stdin, events: make(chan sidecarEvent, 32), logger: logger}
	go client.read(stdout)
	if err := client.send(sidecarConnectMessage(relay, accessToken, participantID, name, color)); err != nil {
		return nil, err
	}
	return client, nil
}
func (client *sidecar) read(output ioReader) {
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		var event sidecarEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			client.logger.Warn("discarding malformed sidecar output", "error", err)
			continue
		}
		client.events <- event
	}
	if err := scanner.Err(); err != nil {
		client.logger.Warn("sidecar output ended", "error", err)
	}
	close(client.events)
}
func (client *sidecar) send(message map[string]any) error {
	encoded, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("encode sidecar message: %w", err)
	}
	client.mu.Lock()
	defer client.mu.Unlock()
	if _, err := client.stdin.Write(append(encoded, '\n')); err != nil {
		return fmt.Errorf("write sidecar message: %w", err)
	}
	return nil
}
func (client *sidecar) stop() {
	if err := client.send(map[string]any{"type": "close"}); err != nil {
		client.logger.Warn("close sidecar", "error", err)
	}
	if err := client.stdin.Close(); err != nil {
		client.logger.Warn("close sidecar stdin", "error", err)
	}
	if err := client.command.Wait(); err != nil {
		client.logger.Warn("sidecar exited", "error", err)
	}
}

type sidecarMsg struct{ event sidecarEvent }
type actionResultMsg struct {
	status, content, panelKind, panelAction string
}
type dialog struct {
	title, prompt, action string
	input                 textinput.Model
	confirm               bool
}

type model struct {
	relay, documentID, participantID, name, color string
	editor                                        textarea.Model
	client                                        *sidecar
	peers                                         []peer
	activity                                      []string
	status                                        string
	width, height, activeMenu, activeItem         int
	menuOpen                                      bool
	mode                                          string
	dialog                                        *dialog
	panel                                         *workspacePanel
	spinner                                       spinner.Model
	busy                                          bool
	replica                                       string
}

const menuRow = 2

var menus = []struct {
	name  string
	items []string
}{
	{"Document", []string{"Open shared document", "New blank document", "Duplicate document", "Import file"}},
	{"Edit", []string{"Find", "Replace", "Go to line", "Bold current line", "Italic current line"}},
	{"View", []string{"Editor", "Markdown preview", "Split editor + preview"}},
	{"Collaborate", []string{"Change profile", "Peer legend", "Activity"}},
	{"Relay", []string{"Refresh relay trace", "Save relay snapshot", "Published versions", "Catalog search"}},
	{"Publish", []string{"Export Markdown", "Export HTML", "Export plain text", "Export Automerge", "Export audit report", "Publish exchange", "Import exchange"}},
	{"Help", []string{"Keyboard shortcuts", "Glow: Ex3 README", "Glow: Ex3 testing guide", "Glow: protocols"}},
}

func newModel(config Config, client *sidecar) model {
	editor := textarea.New()
	editor.Placeholder = "Waiting for the shared document…"
	editor.ShowLineNumbers = true
	editor.Prompt = ""
	editor.Focus()
	spin := spinner.New()
	spin.Style = menuStyle
	return model{relay: config.Relay, documentID: config.DocumentID, participantID: fmt.Sprintf("charm-%d", time.Now().UnixNano()), name: config.Name, color: config.Color, editor: editor, client: client, status: "connecting to relay", mode: "editor", spinner: spin}
}
func (m model) Init() tea.Cmd { return waitForSidecar(m.client.events) }
func waitForSidecar(events <-chan sidecarEvent) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-events
		if !ok {
			return sidecarMsg{sidecarEvent{Type: "error", Message: "sidecar stopped"}}
		}
		return sidecarMsg{event}
	}
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := message.(tea.WindowSizeMsg); ok {
		m.width, m.height = size.Width, size.Height
		m.editor.SetWidth(max(20, size.Width*2/3-8))
		m.editor.SetHeight(max(5, size.Height-10))
		if m.panel != nil {
			return m, m.panel.update(size)
		}
		return m, nil
	}
	if m.panel != nil {
		switch message.(type) {
		case sidecarMsg, actionResultMsg, spinner.TickMsg:
			// These messages update shared workspace state before the local panel.
		default:
			return m.updatePanel(message)
		}
	}
	switch msg := message.(type) {
	case sidecarMsg:
		m.applySidecar(msg.event)
		return m, waitForSidecar(m.client.events)
	case actionResultMsg:
		m.busy = false
		m.status = msg.status
		if msg.content != "" {
			m.activity = append(m.activity, msg.content)
		}
		if msg.panelKind == panelResults {
			m.panel = newResultsPanel(msg.status, msg.panelAction, msg.content, m.width, m.height)
			m.editor.Blur()
		}
		return m, nil
	case spinner.TickMsg:
		if !m.busy {
			return m, nil
		}
		var command tea.Cmd
		m.spinner, command = m.spinner.Update(msg)
		return m, command
	case tea.MouseMsg:
		if msg.Type == tea.MouseWheelUp || msg.Type == tea.MouseWheelDown {
			// Intent: Bubble Tea's textarea does not own wheel scrolling, so
			// translate the wheel into editor navigation and keep the document
			// viewport reachable without leaving the terminal workspace.
			// Source: DI-mutoh.
			keyType := tea.KeyUp
			if msg.Type == tea.MouseWheelDown {
				keyType = tea.KeyDown
			}
			commands := make([]tea.Cmd, 0, 3)
			for range 3 {
				var command tea.Cmd
				m.editor, command = m.editor.Update(tea.KeyMsg{Type: keyType})
				commands = append(commands, command)
			}
			m.publishCursor(false)
			return m, tea.Batch(commands...)
		}
		if msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress {
			if msg.Y == menuRow {
				m.activeMenu = m.menuAt(msg.X)
				m.menuOpen = true
				return m, nil
			}
			if !m.menuOpen && m.mode != "preview" {
				m.placeCursor(msg.X, msg.Y)
				m.publishCursor(false)
				return m, nil
			}
		}
	case tea.KeyMsg:
		if m.dialog != nil {
			return m.updateDialog(msg)
		}
		key := msg.String()
		if key == "ctrl+c" || key == "q" && !m.editor.Focused() {
			return m, tea.Quit
		}
		if key == "esc" {
			m.menuOpen = false
			m.editor.Focus()
			return m, nil
		}
		if key == "alt+d" {
			m.activeMenu = 0
			m.menuOpen = true
			return m, nil
		}
		if key == "alt+e" {
			m.activeMenu = 1
			m.menuOpen = true
			return m, nil
		}
		if key == "alt+v" {
			m.activeMenu = 2
			m.menuOpen = true
			return m, nil
		}
		if key == "alt+c" {
			m.activeMenu = 3
			m.menuOpen = true
			return m, nil
		}
		if key == "alt+r" {
			m.activeMenu = 4
			m.menuOpen = true
			return m, nil
		}
		if key == "alt+p" {
			m.activeMenu = 5
			m.menuOpen = true
			return m, nil
		}
		if key == "alt+h" {
			m.activeMenu = 6
			m.menuOpen = true
			return m, nil
		}
		if m.menuOpen {
			switch key {
			case "left":
				m.activeMenu = (m.activeMenu + len(menus) - 1) % len(menus)
				m.activeItem = 0
			case "right":
				m.activeMenu = (m.activeMenu + 1) % len(menus)
				m.activeItem = 0
			case "up":
				m.activeItem = (m.activeItem + len(menus[m.activeMenu].items) - 1) % len(menus[m.activeMenu].items)
			case "down":
				m.activeItem = (m.activeItem + 1) % len(menus[m.activeMenu].items)
			case "enter":
				return m, m.activate()
			}
			return m, nil
		}
		if key == "ctrl+p" || key == "f2" {
			m.mode = "preview"
			m.status = "Markdown preview"
			return m, nil
		}
	}
	if m.mode == "preview" {
		return m, nil
	}
	before := m.editor.Value()
	var command tea.Cmd
	m.editor, command = m.editor.Update(message)
	if after := m.editor.Value(); after != before {
		if err := m.client.send(map[string]any{"type": "set_text", "content": after}); err != nil {
			m.status = err.Error()
		}
	}
	m.publishCursor(false)
	return m, command
}

// updatePanel gives a focused Bubbles surface first access to input, then
// returns focus to the shared editor on close or selection. This keeps local
// navigation independent of Ex3's shared document protocol. Source: DI-vujub.
func (m model) updatePanel(message tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := message.(tea.KeyMsg); ok {
		if key.String() == "esc" {
			m.panel = nil
			m.editor.Focus()
			return m, nil
		}
		if m.panel.hasPicker {
			if selected, path := m.panel.picker.DidSelectFile(key); selected {
				m.panel = nil
				m.editor.Focus()
				return m, m.busyCmd(m.importCmd(path))
			}
		}
		if m.panel.hasList && key.String() == "enter" {
			value := m.panel.selectedValue()
			action := m.panel.action
			if value != "" && action == "open" {
				m.panel = nil
				m.editor.Focus()
				return m, m.busyCmd(m.openCmd(value, false))
			}
		}
	}
	return m, m.panel.update(message)
}

func (m *model) busyCmd(command tea.Cmd) tea.Cmd {
	if command == nil {
		return nil
	}
	m.busy = true
	return tea.Batch(m.spinner.Tick, command)
}

// menuAt maps an actual terminal column to the rendered menu label. The
// header occupies two rows, so click handling must not assume the menu begins
// at row zero or that all labels have the same width. Source: DI-mutoh.
func (m model) menuAt(column int) int {
	start := 0
	for index, menu := range menus {
		style := menuStyle
		if m.menuOpen && index == m.activeMenu {
			style = activeMenuStyle
		}
		width := lipgloss.Width(style.Render(menu.name))
		if column >= start && column < start+width {
			return index
		}
		start += width + 1
	}
	return min(len(menus)-1, max(0, column))
}

func (m *model) activate() tea.Cmd {
	action := menus[m.activeMenu].items[m.activeItem]
	m.menuOpen = false
	m.editor.Focus()
	switch action {
	case "Editor":
		m.mode = "editor"
		m.status = "Editing shared document"
	case "Markdown preview":
		m.mode = "preview"
		m.status = "Markdown preview"
	case "Split editor + preview":
		m.mode = "split"
		m.status = "Split editor and preview"
	case "Bold current line":
		m.wrapCurrentLine("**", "**")
	case "Italic current line":
		m.wrapCurrentLine("_", "_")
	case "Peer legend":
		m.status = fmt.Sprintf("%d remote peer(s) in legend", len(m.peers))
	case "Activity":
		m.panel = newActivityPanel(m.activity, m.width, m.height)
		m.editor.Blur()
	case "Keyboard shortcuts":
		m.panel = newHelpPanel(m.width)
		m.editor.Blur()
	case "Refresh relay trace":
		return m.busyCmd(m.traceCmd())
	case "Save relay snapshot":
		return m.busyCmd(m.snapshotCmd())
	case "Published versions":
		return m.busyCmd(m.publishedCmd())
	case "Catalog search":
		m.openDialog("Catalog search", "Search relay metadata", "catalog", false)
	case "Open shared document":
		return m.busyCmd(m.catalogCmd("", "open"))
	case "New blank document":
		m.openDialog("New blank document", "New document ID", "new", false)
	case "Duplicate document":
		m.openDialog("Duplicate document", "New document ID", "duplicate", false)
	case "Import file":
		var command tea.Cmd
		m.panel, command = newFilePickerPanel(m.width, m.height)
		m.editor.Blur()
		return command
	case "Find":
		m.openDialog("Find", "Text to find", "find", false)
	case "Replace":
		m.openDialog("Replace", "Find|replace (for example old|new)", "replace", false)
	case "Go to line":
		m.openDialog("Go to line", "Line number", "goto", false)
	case "Change profile":
		m.openDialog("Collaboration profile", "name|#RRGGBB", "profile", false)
	case "Export Markdown":
		m.openDialog("Export Markdown", "Absolute output path", "export:md", true)
	case "Export HTML":
		m.openDialog("Export HTML", "Absolute output path", "export:html", true)
	case "Export plain text":
		m.openDialog("Export text", "Absolute output path", "export:txt", true)
	case "Export Automerge":
		m.openDialog("Export Automerge", "Absolute output path", "export:automerge", true)
	case "Export audit report":
		m.openDialog("Export audit report", "Absolute output path", "export:audit", true)
	case "Publish exchange":
		m.openDialog("Publish exchange", "Title", "publish", false)
	case "Import exchange":
		m.openDialog("Import exchange", "Published exchange URL", "exchange", false)
	case "Glow: Ex3 README":
		return m.glowCmd("README.md")
	case "Glow: Ex3 testing guide":
		return m.glowCmd("docs/testing.md")
	case "Glow: protocols":
		return m.glowCmd("protocols")
	}
	return nil
}

func (m *model) openDialog(title, prompt, action string, confirm bool) {
	input := textinput.New()
	input.Placeholder = prompt
	input.Focus()
	m.dialog = &dialog{title: title, prompt: prompt, action: action, input: input, confirm: confirm}
	m.editor.Blur()
}
func (m model) updateDialog(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.String() == "esc" {
		m.dialog = nil
		m.editor.Focus()
		return m, nil
	}
	if key.String() == "enter" {
		if m.dialog.confirm {
			m.dialog.confirm = false
			m.status = "Confirm write: press Enter again to write " + m.dialog.input.Value()
			return m, nil
		}
		d := m.dialog
		m.dialog = nil
		m.editor.Focus()
		return m, m.runDialog(d)
	}
	var cmd tea.Cmd
	m.dialog.input, cmd = m.dialog.input.Update(key)
	return m, cmd
}
func (m *model) runDialog(d *dialog) tea.Cmd {
	value := strings.TrimSpace(d.input.Value())
	if value == "" {
		m.status = "No value entered"
		return nil
	}
	switch d.action {
	case "open", "new":
		return m.busyCmd(m.openCmd(value, d.action == "new"))
	case "duplicate":
		return m.busyCmd(m.duplicateCmd(value))
	case "find":
		m.find(value)
		return nil
	case "replace":
		m.replace(value)
		return nil
	case "goto":
		m.gotoLine(value)
		return nil
	case "profile":
		m.setProfile(value)
		return nil
	case "import":
		return m.busyCmd(m.importCmd(value))
	case "catalog":
		return m.busyCmd(m.catalogCmd(value, ""))
	case "publish":
		return m.busyCmd(m.publishCmd(value))
	case "exchange":
		return m.busyCmd(m.exchangeCmd(value))
	default:
		if strings.HasPrefix(d.action, "export:") {
			return m.busyCmd(m.exportCmd(strings.TrimPrefix(d.action, "export:"), value))
		}
	}
	return nil
}

func (m *model) applySidecar(event sidecarEvent) {
	switch event.Type {
	case "opened", "changed":
		if event.Content != m.editor.Value() {
			previousLine := m.editor.Line()
			m.editor.SetValue(event.Content)
			// Intent: Sidecar changes replace the textarea value, whose default
			// insertion behavior places the cursor at the document end. Restore
			// the local line so remote collaboration never strands a user there.
			// Source: DI-mutoh.
			for m.editor.Line() > previousLine {
				m.editor.CursorUp()
			}
			m.editor.CursorStart()
		}
		if event.DocID != "" {
			m.documentID = event.DocID
		}
		m.status = "live document synchronized"
	case "awareness":
		m.peers = event.Peers
	case "relay_status":
		if event.Connected {
			m.status = "relay connected"
		} else {
			m.status = "relay disconnected"
		}
	case "state":
		m.replica = event.Replica
	case "error":
		m.status = "sidecar: " + event.Message
	case "info":
		m.activity = append(m.activity, event.Message)
	}
}
func (m *model) publishCursor(typing bool) {
	offset := editorUTF16Offset(m.editor)
	if err := m.client.send(map[string]any{"type": "set_cursor", "anchor": offset, "head": offset, "typing": typing}); err != nil {
		m.status = err.Error()
	}
}
func (m *model) placeCursor(x, y int) {
	const editorTop = 3
	const textColumn = 8
	if y < editorTop || x < textColumn {
		return
	}
	line := max(0, m.editor.Line()-max(5, m.height-10)/2) + y - editorTop
	line = min(line, m.editor.LineCount()-1)
	for m.editor.Line() < line {
		m.editor.CursorDown()
	}
	for m.editor.Line() > line {
		m.editor.CursorUp()
	}
	m.editor.SetCursor(max(0, x-textColumn))
}
func (m *model) wrapCurrentLine(before, after string) {
	lines := strings.Split(m.editor.Value(), "\n")
	line := m.editor.Line()
	if line < 0 || line >= len(lines) {
		return
	}
	lines[line] = before + lines[line] + after
	m.editor.SetValue(strings.Join(lines, "\n"))
	if err := m.client.send(map[string]any{"type": "set_text", "content": m.editor.Value()}); err != nil {
		m.status = err.Error()
	} else {
		m.status = "Formatted current line"
	}
}
func (m *model) find(value string) {
	index := strings.Index(m.editor.Value(), value)
	if index < 0 {
		m.status = "Not found: " + value
		return
	}
	m.gotoOffset(index)
	m.status = "Found: " + value
}
func (m *model) replace(value string) {
	parts := strings.SplitN(value, "|", 2)
	if len(parts) != 2 {
		m.status = "Replace needs find|replace"
		return
	}
	text := strings.Replace(m.editor.Value(), parts[0], parts[1], 1)
	if text == m.editor.Value() {
		m.status = "Not found: " + parts[0]
		return
	}
	m.editor.SetValue(text)
	if err := m.client.send(map[string]any{"type": "set_text", "content": text}); err != nil {
		m.status = err.Error()
	} else {
		m.status = "Replaced first match"
	}
}
func (m *model) gotoLine(value string) {
	var line int
	if _, err := fmt.Sscanf(value, "%d", &line); err != nil || line < 1 || line > m.editor.LineCount() {
		m.status = "Line out of range"
		return
	}
	for m.editor.Line() < line-1 {
		m.editor.CursorDown()
	}
	for m.editor.Line() > line-1 {
		m.editor.CursorUp()
	}
	m.status = fmt.Sprintf("Line %d", line)
}
func (m *model) gotoOffset(offset int) {
	before := m.editor.Value()[:offset]
	m.gotoLine(fmt.Sprintf("%d", strings.Count(before, "\n")+1))
	m.editor.SetCursor(len([]rune(before[strings.LastIndex(before, "\n")+1:])))
}
func (m *model) setProfile(value string) {
	parts := strings.SplitN(value, "|", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || !validColor(strings.TrimSpace(parts[1])) {
		m.status = "Profile needs name|#RRGGBB"
		return
	}
	m.name, m.color = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	if err := m.client.send(map[string]any{"type": "set_user", "display_name": m.name, "color": m.color}); err != nil {
		m.status = err.Error()
	} else {
		m.status = "Collaboration profile updated"
	}
}

func (m *model) openCmd(documentID string, blank bool) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.send(map[string]any{"type": "open", "doc_id": documentID}); err != nil {
			return actionResultMsg{status: err.Error()}
		}
		if blank {
			if err := m.client.send(map[string]any{"type": "set_text", "content": ""}); err != nil {
				return actionResultMsg{status: err.Error()}
			}
		}
		m.documentID = documentID
		return actionResultMsg{status: "Opening shared document " + documentID}
	}
}
func (m *model) duplicateCmd(documentID string) tea.Cmd {
	content := m.editor.Value()
	return func() tea.Msg {
		if err := m.client.send(map[string]any{"type": "open", "doc_id": documentID}); err != nil {
			return actionResultMsg{status: err.Error()}
		}
		if err := m.client.send(map[string]any{"type": "set_text", "content": content}); err != nil {
			return actionResultMsg{status: err.Error()}
		}
		m.documentID = documentID
		return actionResultMsg{status: "Duplicating to " + documentID}
	}
}
func (m *model) importCmd(path string) tea.Cmd {
	return func() tea.Msg {
		data, err := os.ReadFile(path)
		if err != nil {
			return actionResultMsg{status: "Import: " + err.Error()}
		}
		next := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if next == "" {
			next = "imported-document"
		}
		if err := m.client.send(map[string]any{"type": "open", "doc_id": next}); err != nil {
			return actionResultMsg{status: err.Error()}
		}
		if err := m.client.send(map[string]any{"type": "set_text", "content": string(data)}); err != nil {
			return actionResultMsg{status: err.Error()}
		}
		m.documentID = next
		return actionResultMsg{status: "Imported " + path + " into " + next}
	}
}

func (m *model) traceCmd() tea.Cmd {
	return m.getJSON("/api/local/documents/"+url.PathEscape(m.documentID)+"/trace?limit=12", func(body []byte) actionResultMsg {
		return actionResultMsg{status: "Relay trace refreshed", content: string(body)}
	})
}
func (m *model) publishedCmd() tea.Cmd {
	return m.getJSON("/api/local/documents/"+url.PathEscape(m.documentID)+"/published", func(body []byte) actionResultMsg {
		return actionResultMsg{status: "Published versions", content: string(body), panelKind: panelResults}
	})
}
func (m *model) catalogCmd(query, action string) tea.Cmd {
	return m.getJSON("/api/local/metadata/search?q="+url.QueryEscape(query), func(body []byte) actionResultMsg {
		return actionResultMsg{status: "Relay document catalog", content: string(body), panelKind: panelResults, panelAction: action}
	})
}
func (m *model) getJSON(path string, done func([]byte) actionResultMsg) tea.Cmd {
	relay := m.relay
	return func() tea.Msg {
		response, err := http.Get(strings.TrimRight(relay, "/") + path)
		if err != nil {
			return actionResultMsg{status: err.Error()}
		}
		body, err := readResponse(response)
		if err != nil {
			return actionResultMsg{status: err.Error()}
		}
		if response.StatusCode/100 != 2 {
			return actionResultMsg{status: fmt.Sprintf("Relay HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(body)))}
		}
		return done(body)
	}
}
func (m *model) snapshotCmd() tea.Cmd {
	return m.withReplica(func(replica string) tea.Cmd {
		return m.postJSON("/api/local/documents/"+url.PathEscape(m.documentID)+"/snapshot", map[string]any{"participant_id": m.participantID, "text_base64": base64.StdEncoding.EncodeToString([]byte(m.editor.Value())), "replica_base64": replica}, "Relay snapshot saved")
	})
}
func (m *model) publishCmd(title string) tea.Cmd {
	return m.withReplica(func(replica string) tea.Cmd {
		return m.postJSON("/api/local/documents/"+url.PathEscape(m.documentID)+"/publish", map[string]any{"participant_id": m.participantID, "source_kind": "current", "title": title, "summary": summarize(m.editor.Value()), "text_base64": base64.StdEncoding.EncodeToString([]byte(m.editor.Value())), "replica_base64": replica, "embodiment": "charm"}, "Signed exchange published")
	})
}
func (m *model) withReplica(next func(string) tea.Cmd) tea.Cmd {
	if m.replica != "" {
		return next(m.replica)
	}
	return func() tea.Msg {
		if err := m.client.send(map[string]any{"type": "get_state"}); err != nil {
			return actionResultMsg{status: err.Error()}
		}
		return actionResultMsg{status: "Requesting Automerge state; choose the action again when ready"}
	}
}
func (m *model) postJSON(path string, payload any, success string) tea.Cmd {
	relay := m.relay
	return func() tea.Msg {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return actionResultMsg{status: err.Error()}
		}
		response, err := http.Post(strings.TrimRight(relay, "/")+path, "application/json", bytes.NewReader(encoded))
		if err != nil {
			return actionResultMsg{status: err.Error()}
		}
		body, err := readResponse(response)
		if err != nil {
			return actionResultMsg{status: err.Error()}
		}
		if response.StatusCode/100 != 2 {
			return actionResultMsg{status: fmt.Sprintf("Relay HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(body)))}
		}
		return actionResultMsg{status: success, content: string(body)}
	}
}
func (m *model) exchangeCmd(raw string) tea.Cmd {
	return func() tea.Msg {
		u, err := url.Parse(raw)
		if err != nil || u.Scheme == "" {
			return actionResultMsg{status: "Invalid published exchange URL"}
		}
		response, err := http.Get(raw)
		if err != nil {
			return actionResultMsg{status: err.Error()}
		}
		body, err := readResponse(response)
		if err != nil {
			return actionResultMsg{status: "Read exchange: " + err.Error()}
		}
		var payload struct {
			Record struct {
				Title   string `json:"title"`
				TextCID string `json:"text_cid"`
			} `json:"record"`
			TextBase64 string `json:"text_base64"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			return actionResultMsg{status: "Read exchange: " + err.Error()}
		}
		if payload.TextBase64 == "" {
			return actionResultMsg{status: "Exchange has no importable text"}
		}
		text, err := base64.StdEncoding.DecodeString(payload.TextBase64)
		if err != nil {
			return actionResultMsg{status: "Decode exchange: " + err.Error()}
		}
		documentID := "imported-" + fmt.Sprint(time.Now().Unix())
		if err := m.client.send(map[string]any{"type": "open", "doc_id": documentID}); err != nil {
			return actionResultMsg{status: err.Error()}
		}
		if err := m.client.send(map[string]any{"type": "set_text", "content": string(text)}); err != nil {
			return actionResultMsg{status: err.Error()}
		}
		m.documentID = documentID
		return actionResultMsg{status: "Imported exchange as " + documentID}
	}
}

func readResponse(response *http.Response) ([]byte, error) {
	body, readErr := io.ReadAll(response.Body)
	closeErr := response.Body.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	return body, nil
}
func (m *model) exportCmd(format, path string) tea.Cmd {
	text, replica, doc := m.editor.Value(), m.replica, m.documentID
	return func() tea.Msg {
		var data []byte
		switch format {
		case "html":
			data = []byte("<!doctype html><html><head><meta charset=\"utf-8\"><title>" + html.EscapeString(doc) + "</title></head><body><pre>" + html.EscapeString(text) + "</pre></body></html>")
		case "automerge":
			if replica == "" {
				return actionResultMsg{status: "Automerge state not ready; select export again"}
			}
			decoded, err := base64.StdEncoding.DecodeString(replica)
			if err != nil {
				return actionResultMsg{status: err.Error()}
			}
			data = decoded
		case "audit":
			data = []byte("# Grid TUI audit report\n\nDocument: " + doc + "\nRelay: " + m.relay + "\nParticipant: " + m.participantID + "\n\n## Relay activity\n" + strings.Join(m.activity, "\n"))
		default:
			data = []byte(text)
		}
		if err := os.WriteFile(path, data, 0o600); err != nil {
			return actionResultMsg{status: "Export: " + err.Error()}
		}
		return actionResultMsg{status: "Exported " + format + " to " + path}
	}
}
func (m *model) glowCmd(target string) tea.Cmd {
	return func() tea.Msg {
		glowPath, err := exec.LookPath("glow")
		if err != nil {
			return actionResultMsg{status: "Glow is not installed; install it from https://charm.land/glow"}
		}
		command := exec.Command(glowPath, target)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			return actionResultMsg{status: "Glow: " + err.Error()}
		}
		return actionResultMsg{status: "Closed Glow reader"}
	}
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading Grid TUI…"
	}
	header := titleStyle.Render("Grid TUI") + "  " + mutedStyle.Render("Charm collaboration workspace") + "\n" + mutedStyle.Render("document "+m.documentID+"  relay "+m.relay+"  you "+m.name)
	var labels []string
	for index, menu := range menus {
		style := menuStyle
		if m.menuOpen && index == m.activeMenu {
			style = activeMenuStyle
		}
		labels = append(labels, style.Render(menu.name))
	}
	bar := strings.Join(labels, " ")
	main := m.mainView()
	sidebar := panelStyle.Width(max(24, m.width-m.width*2/3-2)).Render(m.sidebar())
	status := m.status
	if m.busy {
		status = m.spinner.View() + " " + status
	}
	footer := mutedStyle.Render(status + "  •  Alt+key menu  •  arrows/Enter choose  •  Esc close  •  Ctrl+C quit")
	if m.panel != nil {
		return lipgloss.JoinVertical(lipgloss.Left, header, bar, m.panel.view(), footer)
	}
	screen := lipgloss.JoinVertical(lipgloss.Left, header, bar, lipgloss.JoinHorizontal(lipgloss.Top, main, sidebar), footer)
	if m.menuOpen {
		screen += "\n" + m.menuView()
	}
	if m.dialog != nil {
		screen += "\n" + panelStyle.Width(min(72, m.width-4)).Render(titleStyle.Render(m.dialog.title)+"\n"+mutedStyle.Render(m.dialog.prompt)+"\n"+m.dialog.input.View()+"\n"+mutedStyle.Render("Enter continues · Esc cancels"))
	}
	return screen
}
func (m model) mainView() string {
	width := max(30, m.width*2/3-2)
	if m.mode == "preview" {
		return panelStyle.Width(width).Render(titleStyle.Render("Markdown preview") + "\n\n" + renderMarkdown(m.editor.Value(), width-6))
	}
	editor := renderEditor(m.editor.Value(), m.peers, m.editor.Line(), editorUTF16Offset(m.editor), m.height-10)
	if m.mode == "split" {
		return panelStyle.Width(width).Render(titleStyle.Render("Editor") + "\n" + editor + "\n\n" + titleStyle.Render("Preview") + "\n" + renderMarkdown(m.editor.Value(), width-6))
	}
	return panelStyle.Width(width).Render(editor)
}
func (m model) menuView() string {
	menu := menus[m.activeMenu]
	lines := []string{titleStyle.Render(menu.name)}
	for index, item := range menu.items {
		prefix := "  "
		if index == m.activeItem {
			prefix = "› "
			item = activeMenuStyle.Render(item)
		}
		lines = append(lines, prefix+item)
	}
	return panelStyle.Render(strings.Join(lines, "\n"))
}
func (m model) sidebar() string {
	lines := []string{titleStyle.Render("Live peers")}
	if len(m.peers) == 0 {
		lines = append(lines, mutedStyle.Render("No remote peers observed"))
	}
	// Intent: A relay can report many historical peers. Keep the collaboration
	// legend inside the terminal height so it cannot push the editor below the
	// visible screen. Source: DI-mutoh.
	peerLimit := max(1, (max(2, m.height-15)-5)/2)
	visiblePeers := min(len(m.peers), peerLimit)
	for _, peer := range m.peers[:visiblePeers] {
		state := "idle"
		if peer.Typing {
			state = "typing"
		}
		// Intent: Render remote typing as a stable marker. A 15-FPS terminal-wide
		// Harmonica redraw made the workspace visibly throb even while idle.
		// Source: DI-tubol. TODO-jufip.
		marker := "■"
		if peer.Typing {
			marker = "●"
		}
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color(peer.Color)).Render(marker)+" "+peer.Name+"\n  "+peer.Embodiment+" · "+state+" · "+lineColumn(m.editor.Value(), peer.Anchor))
	}
	if hiddenPeers := len(m.peers) - visiblePeers; hiddenPeers > 0 {
		lines = append(lines, mutedStyle.Render(fmt.Sprintf("… %d more peers", hiddenPeers)))
	}
	lines = append(lines, "", titleStyle.Render("Relay activity"))
	for _, entry := range last(m.activity, 2) {
		lines = append(lines, "• "+entry)
	}
	return strings.Join(lines, "\n")
}
func renderEditor(text string, peers []peer, localLine, localOffset, height int) string {
	lines := strings.Split(text, "\n")
	start := max(0, localLine-height/2)
	end := min(len(lines), start+height)
	var rendered []string
	for line := start; line < end; line++ {
		rendered = append(rendered, fmt.Sprintf("%4d  %s", line+1, renderLine(lines[line], peers, line, text, localOffset)))
	}
	return strings.Join(rendered, "\n")
}
func renderLine(line string, peers []peer, lineNumber int, text string, localOffset int) string {
	base := lineStartUTF16(text, lineNumber)
	var result strings.Builder
	for index, character := range []rune(line) {
		position := base + utf16Length(string([]rune(line)[:index]))
		if position == localOffset {
			result.WriteString(localCursorStyle.Render("▌"))
		}
		style := lipgloss.NewStyle()
		for _, peer := range peers {
			from, to := min(peer.Anchor, peer.Head), max(peer.Anchor, peer.Head)
			if position >= from && position < to {
				style = style.Background(lipgloss.Color(peer.Color))
			}
			if position == peer.Head {
				style = lipgloss.NewStyle().Background(lipgloss.Color(peer.Color)).Foreground(lipgloss.Color("0"))
			}
		}
		result.WriteString(style.Render(string(character)))
	}
	if localOffset == base+utf16Length(line) {
		result.WriteString(localCursorStyle.Render("▌"))
	}
	// Intent: Preserve an end-of-line remote caret, where there is no document
	// character to carry the peer's cursor color. Source: DI-mutoh.
	for _, peer := range peers {
		if peer.Head == base+utf16Length(line) {
			result.WriteString(lipgloss.NewStyle().Background(lipgloss.Color(peer.Color)).Render(" "))
		}
	}
	return result.String()
}
func editorUTF16Offset(editor textarea.Model) int {
	lines := strings.Split(editor.Value(), "\n")
	line := min(editor.Line(), len(lines)-1)
	var offset int
	for index := 0; index < line; index++ {
		offset += utf16Length(lines[index]) + 1
	}
	info := editor.LineInfo()
	runes := []rune(lines[line])
	offset += utf16Length(string(runes[:min(info.StartColumn+info.ColumnOffset, len(runes))]))
	return offset
}
func utf16Length(value string) int {
	var length int
	for _, character := range value {
		length += len(utf16.Encode([]rune{character}))
	}
	return length
}
func lineStartUTF16(text string, lineNumber int) int {
	lines := strings.Split(text, "\n")
	offset := 0
	for index := 0; index < min(lineNumber, len(lines)); index++ {
		offset += utf16Length(lines[index]) + 1
	}
	return offset
}
func lineColumn(text string, offset int) string {
	lines := strings.Split(text, "\n")
	consumed := 0
	for index, line := range lines {
		size := utf16Length(line)
		if offset <= consumed+size {
			return fmt.Sprintf("%d:%d", index+1, offset-consumed+1)
		}
		consumed += size + 1
	}
	return fmt.Sprintf("%d:1", len(lines))
}
func renderMarkdown(text string, width int) string {
	renderer, err := glamour.NewTermRenderer(glamour.WithWordWrap(width), glamour.WithColorProfile(termenv.TrueColor), glamour.WithStandardStyle("dark"))
	if err != nil {
		return text
	}
	result, err := renderer.Render(text)
	if err != nil {
		return text
	}
	return result
}
func summarize(text string) string {
	text = strings.Join(strings.Fields(text), " ")
	if len(text) > 180 {
		return text[:180]
	}
	return text
}
func last(values []string, count int) []string {
	if len(values) <= count {
		return values
	}
	return values[len(values)-count:]
}
func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}
func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}
func validColor(value string) bool {
	if len(value) != 7 || value[0] != '#' {
		return false
	}
	for _, character := range value[1:] {
		if !strings.ContainsRune("0123456789abcdefABCDEF", character) {
			return false
		}
	}
	return true
}

// Run starts the Charm terminal embodiment using the existing Ex3 sidecar.
func Run(config Config) error {
	if launchNeedsForm(config) {
		selected, err := NewLaunchForm(config).Run()
		if err != nil {
			return fmt.Errorf("Grid TUI launch form: %w", err)
		}
		config = selected
	}
	config = launchDefaults(config)
	lipgloss.SetColorProfile(termenv.TrueColor)
	logger := charmLog.New(os.Stderr)
	participantID := fmt.Sprintf("charm-%d", time.Now().UnixNano())
	client, err := startSidecar(config.Relay, config.AccessToken, participantID, config.Name, config.Color, logger)
	if err != nil {
		return err
	}
	defer client.stop()
	if err := client.send(map[string]any{"type": "open", "doc_id": config.DocumentID}); err != nil {
		return err
	}
	m := newModel(config, client)
	m.participantID = participantID
	// Intent: Enable click-only mouse reporting for menus while avoiding full
	// motion tracking, whose raw terminal escape bytes can enter shared text.
	// Source: DI-mutoh.
	_, err = tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	return err
}

var _ = errors.New
