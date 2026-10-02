// Command grid-tui provides Charm's terminal embodiment for Ex3 collaboration.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode/utf16"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	charmLog "github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

const (
	defaultColor = "#8b5cf6"
	defaultName  = "Charm User"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	panelStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(0, 1)
	selectionDim = lipgloss.NewStyle().Background(lipgloss.Color("60"))
	// Intent: Keep the active terminal user's editing position unmistakable on
	// light and dark terminals; peer colors remain reserved for remote awareness.
	// Source: DI-holoz
	localCursorStyle = lipgloss.NewStyle().Background(lipgloss.Color("#facc15")).Foreground(lipgloss.Color("0")).Bold(true)
)

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
}

type sidecar struct {
	command *exec.Cmd
	stdin   ioWriteCloser
	events  chan sidecarEvent
	mu      sync.Mutex
	logger  *charmLog.Logger
}

type ioWriteCloser interface {
	Write([]byte) (int, error)
	Close() error
}

func startSidecar(relay, participantID, name, color string, logger *charmLog.Logger) (*sidecar, error) {
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
	if err := client.send(map[string]any{
		"type": "connect", "relay_url": relay, "participant_id": participantID,
		"display_name": name, "color": color, "embodiment": "charm",
	}); err != nil {
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

type ioReader interface{ Read([]byte) (int, error) }

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

type model struct {
	relay         string
	documentID    string
	participantID string
	name          string
	color         string
	editor        textarea.Model
	client        *sidecar
	peers         []peer
	activity      []string
	status        string
	width         int
	height        int
	preview       bool
}

func newModel(relay, documentID, name, color string, client *sidecar) model {
	editor := textarea.New()
	editor.Placeholder = "Waiting for the shared document…"
	editor.ShowLineNumbers = true
	editor.Prompt = ""
	editor.Focus()
	return model{
		relay: relay, documentID: documentID, name: name, color: color,
		participantID: fmt.Sprintf("charm-%d", time.Now().UnixNano()), editor: editor, client: client,
		status: "connecting to relay",
	}
}

func (model model) Init() tea.Cmd { return waitForSidecar(model.client.events) }

func waitForSidecar(events <-chan sidecarEvent) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-events
		if !ok {
			return sidecarMsg{event: sidecarEvent{Type: "error", Message: "sidecar stopped"}}
		}
		return sidecarMsg{event: event}
	}
}

func (model model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.width, model.height = message.Width, message.Height
		model.editor.SetWidth(max(20, message.Width*2/3-8))
		model.editor.SetHeight(max(5, message.Height-8))
	case sidecarMsg:
		model.applySidecar(message.event)
		return model, waitForSidecar(model.client.events)
	case tea.KeyMsg:
		if message.String() == "ctrl+c" || message.String() == "q" && !model.editor.Focused() {
			return model, tea.Quit
		}
		if message.String() == "ctrl+p" || message.String() == "alt+p" || message.String() == "f2" {
			model.preview = !model.preview
			if model.preview {
				model.status = "Markdown preview"
			} else {
				model.status = "Editing shared document"
			}
			return model, nil
		}
	case tea.MouseMsg:
		if message.Button == tea.MouseButtonLeft && message.Action == tea.MouseActionPress && !model.preview {
			model.placeCursor(message.X, message.Y)
			model.publishCursor(false)
			return model, nil
		}
	}
	before := model.editor.Value()
	var command tea.Cmd
	model.editor, command = model.editor.Update(message)
	if after := model.editor.Value(); after != before {
		if err := model.client.send(map[string]any{"type": "set_text", "content": after}); err != nil {
			model.status = err.Error()
		}
	}
	model.publishCursor(false)
	return model, command
}

func (model *model) placeCursor(x, y int) {
	const editorTop = 3
	const textColumn = 8
	if y < editorTop || x < textColumn {
		return
	}
	line := max(0, model.editor.Line()-max(5, model.height-8)/2) + y - editorTop
	line = min(line, model.editor.LineCount()-1)
	for model.editor.Line() < line {
		model.editor.CursorDown()
	}
	for model.editor.Line() > line {
		model.editor.CursorUp()
	}
	model.editor.SetCursor(max(0, x-textColumn))
}

func (model *model) applySidecar(event sidecarEvent) {
	switch event.Type {
	case "opened", "changed":
		if event.Content != model.editor.Value() {
			model.editor.SetValue(event.Content)
		}
		model.status = "live document synchronized"
	case "awareness":
		model.peers = event.Peers
	case "relay_status":
		if event.Connected {
			model.status = "relay connected"
		} else {
			model.status = "relay disconnected"
		}
	case "error":
		model.status = "sidecar: " + event.Message
	case "info":
		model.activity = append(model.activity, event.Message)
	}
}

func (model *model) publishCursor(typing bool) {
	offset := editorUTF16Offset(model.editor)
	if err := model.client.send(map[string]any{"type": "set_cursor", "anchor": offset, "head": offset, "typing": typing}); err != nil {
		model.status = err.Error()
	}
}

func editorUTF16Offset(editor textarea.Model) int {
	lines := strings.Split(editor.Value(), "\n")
	line := min(editor.Line(), len(lines)-1)
	var offset int
	for index := 0; index < line; index++ {
		offset += utf16Length(lines[index]) + 1
	}
	info := editor.LineInfo()
	offset += utf16Length(string([]rune(lines[line])[:min(info.StartColumn+info.ColumnOffset, len([]rune(lines[line])))]))
	return offset
}

func utf16Length(value string) int {
	var length int
	for _, character := range value {
		length += len(utf16.Encode([]rune{character}))
	}
	return length
}

func (model model) View() string {
	if model.width == 0 {
		return "Loading Grid TUI…"
	}
	header := titleStyle.Render("Grid TUI") + "  " + mutedStyle.Render("Charm collaboration workspace") + "\n" +
		mutedStyle.Render("document "+model.documentID+"  relay "+model.relay+"  you "+model.name)
	var main string
	if model.preview {
		main = panelStyle.Width(model.width*2/3 - 2).Render(titleStyle.Render("Markdown preview") + "\n\n" + renderMarkdown(model.editor.Value(), model.width*2/3-6))
	} else {
		main = panelStyle.Width(model.width*2/3 - 2).Render(renderEditor(model.editor.Value(), model.peers, model.editor.Line(), editorUTF16Offset(model.editor), model.height-8))
	}
	sidebar := panelStyle.Width(model.width - model.width*2/3 - 2).Render(model.sidebar())
	footer := mutedStyle.Render(model.status + "   •   ctrl+p / F2 preview   •   click to move cursor   •   ctrl+c quit")
	return lipgloss.JoinVertical(lipgloss.Left, header, lipgloss.JoinHorizontal(lipgloss.Top, main, sidebar), footer)
}

func (model model) sidebar() string {
	var lines []string
	lines = append(lines, titleStyle.Render("Live peers"))
	if len(model.peers) == 0 {
		lines = append(lines, mutedStyle.Render("No remote peers observed"))
	}
	for _, peer := range model.peers {
		state := "idle"
		if peer.Typing {
			state = "typing"
		}
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color(peer.Color)).Render("■")+" "+peer.Name+"\n  "+peer.Embodiment+" · "+state+" · "+lineColumn(model.editor.Value(), peer.Anchor))
	}
	lines = append(lines, "", titleStyle.Render("Relay activity"))
	if len(model.activity) == 0 {
		lines = append(lines, mutedStyle.Render("Live sync and awareness are relay-observed."))
	}
	for _, entry := range last(model.activity, 6) {
		lines = append(lines, "• "+entry)
	}
	return strings.Join(lines, "\n")
}

func renderEditor(text string, peers []peer, localLine, localOffset int, height int) string {
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
				style = selectionDim.Background(lipgloss.Color(peer.Color))
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
	for _, peer := range peers {
		if peer.Head == base+utf16Length(line) {
			result.WriteString(lipgloss.NewStyle().Background(lipgloss.Color(peer.Color)).Render(" "))
		}
	}
	return result.String()
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

func profile() (string, string, error) {
	name, color := defaultName, defaultColor
	form := huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("Your display name").Value(&name).Validate(nonEmpty),
		huh.NewInput().Title("Your color (#RRGGBB)").Value(&color).Validate(validColor),
	))
	if err := form.Run(); err != nil {
		return "", "", err
	}
	return name, color, nil
}

func nonEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("a display name is required")
	}
	return nil
}
func validColor(value string) error {
	if len(value) != 7 || value[0] != '#' {
		return errors.New("use #RRGGBB")
	}
	for _, character := range value[1:] {
		if !strings.ContainsRune("0123456789abcdefABCDEF", character) {
			return errors.New("use #RRGGBB")
		}
	}
	return nil
}

func main() {
	lipgloss.SetColorProfile(termenv.TrueColor)
	relay := flag.String("relay", "http://127.0.0.1:7025", "grid relay base URL")
	documentID := flag.String("doc", "demo", "document ID to open")
	flag.Parse()
	name, color, err := profile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Grid TUI profile:", err)
		return
	}
	logger := charmLog.New(os.Stderr)
	participantID := fmt.Sprintf("charm-%d", time.Now().UnixNano())
	client, err := startSidecar(*relay, participantID, name, color, logger)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Grid TUI:", err)
		return
	}
	defer client.stop()
	if err := client.send(map[string]any{"type": "open", "doc_id": *documentID}); err != nil {
		fmt.Fprintln(os.Stderr, "Grid TUI:", err)
		return
	}
	model := newModel(*relay, *documentID, name, color, client)
	model.participantID = participantID
	if _, err := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		logger.Error("Grid TUI stopped", "error", err)
	}
}
