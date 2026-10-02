// Command grid-charm is the Charm terminal embodiment for the Ex3 relay.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/service"
)

const refreshInterval = time.Second

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	panelStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(0, 1)
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

type traceItem struct{ entry service.TraceEntry }

func (item traceItem) Title() string {
	return fmt.Sprintf("%-17s %-14s %s", item.entry.Protocol, item.entry.Kind, item.entry.Summary)
}
func (item traceItem) Description() string {
	return fmt.Sprintf("%s · %s", item.entry.ReceivedAt, shortID(item.entry.PCID))
}
func (item traceItem) FilterValue() string {
	return strings.Join([]string{item.entry.Protocol, item.entry.Kind, item.entry.Summary, item.entry.PCID}, " ")
}

type traceLoadedMsg struct {
	feed service.TraceFeed
	err  error
}

type tickMsg time.Time

type model struct {
	relay       string
	documentID  string
	list        list.Model
	document    textinput.Model
	familyIndex int
	kindIndex   int
	status      string
	width       int
	height      int
}

var families = []string{"all", "live-document", "live-awareness", "document-metadata", "publish-document"}
var kinds = []string{"all", "change", "update", "publish", "restore"}

func newModel(relay string, documentID string) model {
	input := textinput.New()
	input.Placeholder = "document ID"
	input.SetValue(documentID)
	input.Width = 28
	entries := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	entries.Title = "Relay-observed envelopes"
	entries.SetShowStatusBar(false)
	entries.SetShowHelp(false)
	entries.SetFilteringEnabled(true)
	return model{relay: strings.TrimRight(relay, "/"), documentID: documentID, list: entries, document: input, status: "connecting"}
}

func (model model) Init() tea.Cmd {
	return tea.Batch(fetchTrace(model.relay, model.documentID), nextTick())
}

func (model model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.width, model.height = message.Width, message.Height
		model.resize()
	case traceLoadedMsg:
		if message.err != nil {
			model.status = "relay error: " + message.err.Error()
			return model, nil
		}
		model.status = fmt.Sprintf("%d accepted relay-observed envelope(s)", len(message.feed.Entries))
		return model, model.setEntries(message.feed.Entries)
	case tickMsg:
		return model, tea.Batch(fetchTrace(model.relay, model.documentID), nextTick())
	case tea.KeyMsg:
		if model.document.Focused() {
			switch message.String() {
			case "enter":
				candidate := strings.TrimSpace(model.document.Value())
				if candidate != "" {
					model.documentID = candidate
				}
				model.document.Blur()
				return model, fetchTrace(model.relay, model.documentID)
			case "esc":
				model.document.SetValue(model.documentID)
				model.document.Blur()
				return model, nil
			}
			var command tea.Cmd
			model.document, command = model.document.Update(message)
			return model, command
		}
		switch message.String() {
		case "q", "ctrl+c":
			return model, tea.Quit
		case "d":
			model.document.Focus()
			return model, nil
		case "f":
			model.familyIndex = (model.familyIndex + 1) % len(families)
			return model, fetchTrace(model.relay, model.documentID)
		case "k":
			model.kindIndex = (model.kindIndex + 1) % len(kinds)
			return model, fetchTrace(model.relay, model.documentID)
		case "r":
			return model, fetchTrace(model.relay, model.documentID)
		}
	}
	var command tea.Cmd
	model.list, command = model.list.Update(message)
	return model, command
}

func (model *model) resize() {
	if model.width == 0 || model.height == 0 {
		return
	}
	listWidth := model.width / 2
	model.list.SetSize(listWidth-2, model.height-8)
}

func (model *model) setEntries(entries []service.TraceEntry) tea.Cmd {
	items := make([]list.Item, 0, len(entries))
	for _, entry := range entries {
		if model.matches(entry) {
			items = append(items, traceItem{entry: entry})
		}
	}
	return model.list.SetItems(items)
}

func (model model) matches(entry service.TraceEntry) bool {
	if family := families[model.familyIndex]; family != "all" && entry.Protocol != family {
		return false
	}
	if kind := kinds[model.kindIndex]; kind != "all" && entry.Kind != kind {
		return false
	}
	return true
}

func (model model) View() string {
	if model.width == 0 {
		return "Loading Grid Charm…"
	}
	header := titleStyle.Render("Grid Charm") + "  " + labelStyle.Render("actual Bubble Tea relay inspector") + "\n" +
		labelStyle.Render("relay") + " " + model.relay + "   " + labelStyle.Render("document") + " " + model.documentID + "\n"
	filters := labelStyle.Render("filters") + "  [d] document " + model.document.View() + "  [f] family " + families[model.familyIndex] + "  [k] message " + kinds[model.kindIndex]
	detail := model.detail()
	body := lipgloss.JoinHorizontal(lipgloss.Top, panelStyle.Width(model.width/2-2).Render(model.list.View()), panelStyle.Width(model.width-model.width/2-2).Render(detail))
	footer := helpStyle.Render(model.status + "   •   [r] refresh   [/] search list   [q] quit")
	return lipgloss.JoinVertical(lipgloss.Left, header, filters, body, footer)
}

func (model model) detail() string {
	selected, ok := model.list.SelectedItem().(traceItem)
	if !ok {
		return "Decoded fields\n\nSelect a relay-observed envelope to inspect its protocol, pCID, identity, and payload."
	}
	payload, err := json.MarshalIndent(selected.entry.DecodedPayload, "", "  ")
	if err != nil {
		payload = []byte(err.Error())
	}
	return titleStyle.Render("Decoded fields") + "\n\n" +
		labelStyle.Render("protocol") + "  " + selected.entry.Protocol + "\n" +
		labelStyle.Render("pCID") + "      " + selected.entry.PCID + "\n" +
		labelStyle.Render("envelope") + "  " + selected.entry.EnvelopeCID + "\n" +
		labelStyle.Render("author") + "    " + selected.entry.Author + "\n\n" + string(payload)
}

func fetchTrace(relay string, documentID string) tea.Cmd {
	// Intent: Present the relay's accepted local observations without implying
	// that this terminal embodiment owns document truth or peer identity.
	// Source: DI-holoz
	return func() tea.Msg {
		endpoint := strings.TrimRight(relay, "/") + "/api/local/documents/" + url.PathEscape(documentID) + "/trace?limit=100"
		response, err := (&http.Client{Timeout: 2 * time.Second}).Get(endpoint)
		if err != nil {
			return traceLoadedMsg{err: err}
		}
		if response.StatusCode != http.StatusOK {
			if closeErr := response.Body.Close(); closeErr != nil {
				return traceLoadedMsg{err: closeErr}
			}
			return traceLoadedMsg{err: fmt.Errorf("%s", response.Status)}
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			if closeErr := response.Body.Close(); closeErr != nil {
				return traceLoadedMsg{err: closeErr}
			}
			return traceLoadedMsg{err: err}
		}
		if err := response.Body.Close(); err != nil {
			return traceLoadedMsg{err: err}
		}
		var feed service.TraceFeed
		if err := json.Unmarshal(body, &feed); err != nil {
			return traceLoadedMsg{err: err}
		}
		return traceLoadedMsg{feed: feed}
	}
}

func nextTick() tea.Cmd {
	return tea.Tick(refreshInterval, func(time.Time) tea.Msg { return tickMsg(time.Now()) })
}
func shortID(value string) string {
	if len(value) <= 18 {
		return value
	}
	return value[:18] + "…"
}

func main() {
	relay := flag.String("relay", "http://127.0.0.1:7025", "grid relay base URL")
	documentID := flag.String("doc", "demo", "document ID to inspect")
	flag.Parse()
	if _, err := tea.NewProgram(newModel(*relay, *documentID), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Grid Charm error:", err)
	}
}
