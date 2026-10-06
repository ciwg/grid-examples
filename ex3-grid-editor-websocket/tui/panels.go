package tui

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	panelHelp     = "help"
	panelActivity = "activity"
	panelFiles    = "files"
	panelResults  = "results"
)

// workspacePanel keeps a Bubbles surface local to Grid TUI. It deliberately
// contains no document or awareness state: the panel is an implementation
// convenience, not collaboration protocol state. Source: DI-vujub.
type workspacePanel struct {
	title       string
	kind        string
	action      string
	list        list.Model
	table       table.Model
	picker      filepicker.Model
	viewport    viewport.Model
	help        help.Model
	hasHelp     bool
	hasList     bool
	hasTable    bool
	hasPicker   bool
	hasViewport bool
}

func newHelpPanel(width int) *workspacePanel {
	return &workspacePanel{title: "Keyboard shortcuts", kind: panelHelp, help: newWorkspaceHelp(width), hasHelp: true}
}

type resultItem struct {
	title, detail, value string
}

type workspaceKeyMap struct {
	menus, preview, close, quit key.Binding
}

func newWorkspaceKeyMap() workspaceKeyMap {
	return workspaceKeyMap{
		menus:   key.NewBinding(key.WithKeys("alt+d", "alt+e", "alt+v", "alt+c", "alt+r", "alt+p", "alt+h"), key.WithHelp("Alt+key", "open menu")),
		preview: key.NewBinding(key.WithKeys("ctrl+p", "f2"), key.WithHelp("Ctrl+P", "Markdown preview")),
		close:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("Esc", "close panel")),
		quit:    key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("Ctrl+C", "quit")),
	}
}

func (keys workspaceKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keys.menus, keys.preview, keys.close, keys.quit}
}

func (keys workspaceKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{keys.menus, keys.preview}, {keys.close, keys.quit}}
}

func newWorkspaceHelp(width int) help.Model {
	helpView := help.New()
	helpView.Width = max(30, width-8)
	helpView.ShowAll = true
	return helpView
}

func (item resultItem) Title() string       { return item.title }
func (item resultItem) Description() string { return item.detail }
func (item resultItem) FilterValue() string { return item.title + " " + item.detail }

func newActivityPanel(entries []string, width, height int) *workspacePanel {
	view := viewport.New(max(28, width-8), max(5, height-10))
	content := "No local relay activity yet."
	if len(entries) > 0 {
		content = strings.Join(entries, "\n\n")
	}
	view.SetContent(content)
	return &workspacePanel{title: "Relay activity", kind: panelActivity, viewport: view, hasViewport: true}
}

func newFilePickerPanel(width, height int) (*workspacePanel, tea.Cmd) {
	picker := filepicker.New()
	picker.SetHeight(max(5, height-10))
	return &workspacePanel{title: "Import a file", kind: panelFiles, picker: picker, hasPicker: true}, picker.Init()
}

func newResultsPanel(title, action, content string, width, height int) *workspacePanel {
	items, rows := resultRows(content)
	listItems := make([]list.Item, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, item)
	}
	selector := list.New(listItems, list.NewDefaultDelegate(), max(28, width/2-6), max(8, height-12))
	selector.Title = "Choose a document"
	selector.SetShowStatusBar(false)
	selector.SetShowHelp(false)
	metadata := table.New(
		table.WithColumns([]table.Column{{Title: "ID", Width: max(12, width/4-4)}, {Title: "Details", Width: max(20, width/2-12)}}),
		table.WithRows(rows),
		table.WithWidth(max(34, width/2-4)),
		table.WithHeight(max(8, height-12)),
	)
	return &workspacePanel{title: title, kind: panelResults, action: action, list: selector, table: metadata, hasList: true, hasTable: true}
}

func resultRows(content string) ([]resultItem, []table.Row) {
	var payload any
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		item := resultItem{title: "Relay response", detail: strings.TrimSpace(content), value: strings.TrimSpace(content)}
		return []resultItem{item}, []table.Row{{item.title, item.detail}}
	}
	records := resultRecords(payload)
	if len(records) == 0 {
		return []resultItem{{title: "No records", detail: "The relay returned no selectable records."}}, nil
	}
	items := make([]resultItem, 0, len(records))
	rows := make([]table.Row, 0, len(records))
	for index, record := range records {
		identifier := firstRecordValue(record, "document_id", "doc_id", "id", "title", "cid")
		if identifier == "" {
			identifier = fmt.Sprintf("record %d", index+1)
		}
		detail := recordDetail(record, identifier)
		items = append(items, resultItem{title: identifier, detail: detail, value: identifier})
		rows = append(rows, table.Row{identifier, detail})
	}
	return items, rows
}

func resultRecords(payload any) []map[string]any {
	switch value := payload.(type) {
	case []any:
		records := make([]map[string]any, 0, len(value))
		for _, item := range value {
			if record, ok := item.(map[string]any); ok {
				records = append(records, record)
			}
		}
		return records
	case map[string]any:
		for _, key := range []string{"documents", "entries", "records", "versions", "items"} {
			if nested, ok := value[key]; ok {
				if records := resultRecords(nested); len(records) > 0 {
					return records
				}
			}
		}
		return []map[string]any{value}
	default:
		return nil
	}
}

func firstRecordValue(record map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := record[key]; ok {
			return fmt.Sprint(value)
		}
	}
	return ""
}

func recordDetail(record map[string]any, identifier string) string {
	keys := make([]string, 0, len(record))
	for key := range record {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value := fmt.Sprint(record[key])
		if value != identifier && value != "" {
			parts = append(parts, key+": "+value)
		}
	}
	if len(parts) == 0 {
		return "Relay metadata"
	}
	return strings.Join(parts, " · ")
}

func (panel *workspacePanel) update(message tea.Msg) tea.Cmd {
	var command tea.Cmd
	switch {
	case panel.hasPicker:
		panel.picker, command = panel.picker.Update(message)
	case panel.hasViewport:
		panel.viewport, command = panel.viewport.Update(message)
	case panel.hasList:
		panel.list, command = panel.list.Update(message)
		if panel.hasTable {
			panel.table, _ = panel.table.Update(message)
		}
	}
	return command
}

func (panel workspacePanel) view() string {
	var body string
	switch {
	case panel.hasHelp:
		body = panel.help.View(newWorkspaceKeyMap()) + "\n" + mutedStyle.Render("Esc closes")
	case panel.hasPicker:
		body = panel.picker.View() + "\n" + mutedStyle.Render("Enter selects a file · Esc closes")
	case panel.hasViewport:
		body = panel.viewport.View() + "\n" + mutedStyle.Render("↑/↓ scroll · Esc closes")
	case panel.hasList:
		body = lipgloss.JoinHorizontal(lipgloss.Top, panel.list.View(), "  ", panel.table.View())
		body += "\n" + mutedStyle.Render("Type to filter · Enter chooses · Esc closes")
	}
	return panelStyle.Render(titleStyle.Render(panel.title) + "\n" + body)
}

func (panel workspacePanel) selectedValue() string {
	item, ok := panel.list.SelectedItem().(resultItem)
	if !ok {
		return ""
	}
	return item.value
}
