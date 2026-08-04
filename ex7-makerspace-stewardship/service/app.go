package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// App keeps the demo's shared equipment evidence in memory. Intent: make the
// append-only observations and current safety status easy to inspect in the
// first standalone example. Source: temporary DI pending mint-handle recovery.
type App struct {
	mu    sync.RWMutex
	state State
	store *Store
}

func NewPersistentDemoApp(root string) (*App, error) {
	app := NewDemoApp()
	store, err := NewStore(root)
	if err != nil {
		return nil, err
	}
	app.store = store
	events, err := store.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if err := app.applyEvent(event); err != nil {
			return nil, err
		}
	}
	return app, nil
}

func NewDemoApp() *App {
	return &App{state: State{
		Members:        []Member{{ID: "alice", Name: "Alice Nguyen", Initials: "A.N."}, {ID: "carol", Name: "Carol Davis", Initials: "C.D."}, {ID: "dave", Name: "Dave Patel", Initials: "D.P."}},
		Areas:          []Area{{ID: "woodworking", Name: "Woodworking"}, {ID: "fiber", Name: "Fiber Arts"}},
		Authorities:    []Authority{{MemberID: "carol", AreaID: "woodworking", Scopes: []string{"recognize qualifications", "assess tool condition", "clear safety holds", "publish area policy"}}},
		Qualifications: []Qualification{{MemberID: "alice", AreaID: "woodworking", IssuedBy: "carol", Status: "accepted"}},
		Tools: []Tool{
			{ID: "table-saw", Name: "Table saw", AreaID: "woodworking", Condition: "Available for qualified in-space use"},
			{ID: "cordless-drill", Name: "Cordless drill", AreaID: "woodworking", OffSiteLoan: true, Condition: "Available; tool only, no bits included"},
			{ID: "sewing-machine", Name: "Sewing machine", AreaID: "fiber", OffSiteLoan: true, Condition: "Available; tension adjustment is stiff"},
		},
	}}
}

func (a *App) State() State {
	a.mu.RLock()
	defer a.mu.RUnlock()
	state := a.state
	state.Tools = make([]Tool, len(a.state.Tools))
	for i, tool := range a.state.Tools {
		state.Tools[i] = tool
		state.Tools[i].Observations = append([]Observation(nil), tool.Observations...)
		if tool.ActiveLoan != nil {
			loan := *tool.ActiveLoan
			state.Tools[i].ActiveLoan = &loan
		}
	}
	return state
}

func (a *App) AddObservation(toolID, reporterID, text string, safetyHold bool, photos []Photo) (Tool, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if text == "" {
		return Tool{}, errors.New("observation text is required")
	}
	if !a.memberExists(reporterID) {
		return Tool{}, errors.New("unknown member")
	}
	if err := validatePhotos(photos); err != nil {
		return Tool{}, err
	}
	tool, err := a.tool(toolID)
	if err != nil {
		return Tool{}, err
	}
	observation := Observation{ID: fmt.Sprintf("obs-%d", time.Now().UnixNano()), ToolID: toolID, ReporterID: reporterID, Text: text, SafetyHold: safetyHold, Photos: photos, CreatedAt: time.Now().UTC()}
	tool.Observations = append(tool.Observations, observation)
	if safetyHold {
		tool.SafetyHold = true
		tool.Condition = "Safety hold: awaiting area review"
	}
	if err := a.record(Event{Type: "observation", ToolID: toolID, ActorID: reporterID, Text: text, SafetyHold: safetyHold, Photos: photos, CreatedAt: observation.CreatedAt}); err != nil {
		return Tool{}, err
	}
	return *tool, nil
}

func validatePhotos(photos []Photo) error {
	const maxPhotoBytes = 1024 * 1024
	if len(photos) > 3 {
		return errors.New("at most three photos may be attached to one observation")
	}
	for _, photo := range photos {
		if !strings.HasPrefix(photo.DataURL, "data:image/") {
			return errors.New("photo must be an image")
		}
		if len(photo.DataURL) > maxPhotoBytes*2 {
			return errors.New("photo is too large")
		}
	}
	return nil
}

func (a *App) ClearSafetyHold(toolID, stewardID, assessment string) (Tool, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	tool, err := a.tool(toolID)
	if err != nil {
		return Tool{}, err
	}
	if !a.hasScope(stewardID, tool.AreaID, "clear safety holds") {
		return Tool{}, errors.New("member is not recognized to clear this area's safety holds")
	}
	if assessment == "" {
		return Tool{}, errors.New("inspection assessment is required")
	}
	tool.SafetyHold = false
	tool.Condition = assessment
	createdAt := time.Now().UTC()
	tool.Observations = append(tool.Observations, Observation{ID: fmt.Sprintf("obs-%d", time.Now().UnixNano()), ToolID: toolID, ReporterID: stewardID, Text: "Safety hold cleared after inspection: " + assessment, CreatedAt: createdAt})
	if err := a.record(Event{Type: "clear-safety-hold", ToolID: toolID, ActorID: stewardID, Text: assessment, CreatedAt: createdAt}); err != nil {
		return Tool{}, err
	}
	return *tool, nil
}

func (a *App) CreateLoan(toolID, memberID string, dueAt time.Time) (Tool, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	tool, err := a.tool(toolID)
	if err != nil {
		return Tool{}, err
	}
	if !tool.OffSiteLoan {
		return Tool{}, errors.New("this tool is for in-space use only")
	}
	if tool.SafetyHold {
		return Tool{}, errors.New("this tool has a safety hold")
	}
	if tool.ActiveLoan != nil {
		return Tool{}, errors.New("this tool is already loaned")
	}
	if !a.isQualified(memberID, tool.AreaID) {
		return Tool{}, errors.New("member is not qualified for this area")
	}
	if !dueAt.After(time.Now()) {
		return Tool{}, errors.New("return deadline must be in the future")
	}
	tool.ActiveLoan = &Loan{
		MemberID:      memberID,
		DueAt:         dueAt.UTC(),
		CreatedAt:     time.Now().UTC(),
		PolicyVersion: "woodworking-off-site-lending/v1",
		Policy:        "Terms accepted at checkout; later policy changes do not rewrite this loan.",
	}
	createdAt := time.Now().UTC()
	tool.Observations = append(tool.Observations, Observation{
		ID:         fmt.Sprintf("obs-%d", time.Now().UnixNano()),
		ToolID:     toolID,
		ReporterID: memberID,
		Text:       "Off-site loan accepted. Return promised by " + dueAt.UTC().Format(time.RFC822),
		CreatedAt:  createdAt,
	})
	if err := a.record(Event{Type: "loan", ToolID: toolID, ActorID: memberID, DueAt: dueAt.UTC(), CreatedAt: createdAt}); err != nil {
		return Tool{}, err
	}
	return *tool, nil
}

func (a *App) ReturnLoan(toolID, memberID, condition string) (Tool, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	tool, err := a.tool(toolID)
	if err != nil {
		return Tool{}, err
	}
	if tool.ActiveLoan == nil {
		return Tool{}, errors.New("this tool is not currently loaned")
	}
	if tool.ActiveLoan.MemberID != memberID {
		return Tool{}, errors.New("only the current borrower can record this return")
	}
	if condition == "" {
		return Tool{}, errors.New("return condition is required")
	}
	tool.ActiveLoan = nil
	tool.Condition = condition
	createdAt := time.Now().UTC()
	tool.Observations = append(tool.Observations, Observation{
		ID:         fmt.Sprintf("obs-%d", time.Now().UnixNano()),
		ToolID:     toolID,
		ReporterID: memberID,
		Text:       "Off-site return recorded: " + condition,
		CreatedAt:  createdAt,
	})
	if err := a.record(Event{Type: "return", ToolID: toolID, ActorID: memberID, Text: condition, CreatedAt: createdAt}); err != nil {
		return Tool{}, err
	}
	return *tool, nil
}

func (a *App) record(event Event) error {
	if a.store == nil {
		return nil
	}
	return a.store.Append(event)
}

func (a *App) applyEvent(event Event) error {
	tool, err := a.tool(event.ToolID)
	if err != nil {
		return err
	}
	switch event.Type {
	case "observation":
		tool.Observations = append(tool.Observations, Observation{ID: fmt.Sprintf("replay-%d", len(tool.Observations)+1), ToolID: event.ToolID, ReporterID: event.ActorID, Text: event.Text, SafetyHold: event.SafetyHold, Photos: event.Photos, CreatedAt: event.CreatedAt})
		if event.SafetyHold {
			tool.SafetyHold = true
			tool.Condition = "Safety hold: awaiting area review"
		}
	case "clear-safety-hold":
		tool.SafetyHold = false
		tool.Condition = event.Text
		tool.Observations = append(tool.Observations, Observation{ID: fmt.Sprintf("replay-%d", len(tool.Observations)+1), ToolID: event.ToolID, ReporterID: event.ActorID, Text: "Safety hold cleared after inspection: " + event.Text, CreatedAt: event.CreatedAt})
	case "loan":
		tool.ActiveLoan = &Loan{MemberID: event.ActorID, DueAt: event.DueAt, CreatedAt: event.CreatedAt, PolicyVersion: "woodworking-off-site-lending/v1", Policy: "Terms accepted at checkout; later policy changes do not rewrite this loan."}
		tool.Observations = append(tool.Observations, Observation{ID: fmt.Sprintf("replay-%d", len(tool.Observations)+1), ToolID: event.ToolID, ReporterID: event.ActorID, Text: "Off-site loan accepted. Return promised by " + event.DueAt.Format(time.RFC822), CreatedAt: event.CreatedAt})
	case "return":
		tool.ActiveLoan = nil
		tool.Condition = event.Text
		tool.Observations = append(tool.Observations, Observation{ID: fmt.Sprintf("replay-%d", len(tool.Observations)+1), ToolID: event.ToolID, ReporterID: event.ActorID, Text: "Off-site return recorded: " + event.Text, CreatedAt: event.CreatedAt})
	default:
		return fmt.Errorf("unknown stored event type %q", event.Type)
	}
	return nil
}

func (a *App) tool(id string) (*Tool, error) {
	for i := range a.state.Tools {
		if a.state.Tools[i].ID == id {
			return &a.state.Tools[i], nil
		}
	}
	return nil, errors.New("unknown tool")
}
func (a *App) memberExists(id string) bool {
	for _, member := range a.state.Members {
		if member.ID == id {
			return true
		}
	}
	return false
}
func (a *App) isQualified(memberID, areaID string) bool {
	for _, qualification := range a.state.Qualifications {
		if qualification.MemberID == memberID && qualification.AreaID == areaID && qualification.Status == "accepted" {
			return true
		}
	}
	return false
}
func (a *App) hasScope(memberID, areaID, scope string) bool {
	for _, authority := range a.state.Authorities {
		if authority.MemberID == memberID && authority.AreaID == areaID {
			for _, value := range authority.Scopes {
				if value == scope {
					return true
				}
			}
		}
	}
	return false
}
