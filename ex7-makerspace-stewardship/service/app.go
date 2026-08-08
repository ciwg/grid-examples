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
		Members: []Member{{ID: "alice", Name: "Alice Nguyen", Initials: "A.N."}, {ID: "carol", Name: "Carol Davis", Initials: "C.D."}, {ID: "dave", Name: "Dave Patel", Initials: "D.P."}},
		Areas: []Area{
			{ID: "woodworking", Name: "Woodworking", PolicyVersion: "v1", Policy: "Qualified members may use tools in the space. Portable tools may be loaned when their own terms allow it. A safety hold prevents self-service use until cleared after inspection.", DelegatedBy: "Makerspace governance"},
			{ID: "fiber", Name: "Fiber Arts", PolicyVersion: "v1", Policy: "Members use equipment in the space according to current area guidance. Portable tools may have separate loan terms.", DelegatedBy: "Makerspace governance"},
		},
		Authorities:    []Authority{{MemberID: "carol", AreaID: "woodworking", Scopes: []string{"recognize qualifications", "assess tool condition", "clear safety holds", "publish area policy"}, RecognizedBy: "Makerspace governance", ReviewAt: "2027-01-01"}},
		Qualifications: []Qualification{{MemberID: "alice", AreaID: "woodworking", Scope: "portable-power-tools", IssuedBy: "carol", Status: "accepted"}},
		Tools: []Tool{
			{ID: "table-saw", Name: "Table saw", AreaID: "woodworking", RequiredQualification: "table-saw-safety", Condition: "Available for qualified in-space use"},
			{ID: "cordless-drill", Name: "Cordless drill", AreaID: "woodworking", RequiredQualification: "portable-power-tools", OffSiteLoan: true, Condition: "Available; tool only, no bits included"},
			{ID: "sewing-machine", Name: "Sewing machine", AreaID: "fiber", RequiredQualification: "fiber-equipment", OffSiteLoan: true, Condition: "Available; tension adjustment is stiff"},
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
	createdAt := time.Now().UTC()
	event := Event{Type: "observation", ToolID: toolID, ActorID: reporterID, Text: text, SafetyHold: safetyHold, Photos: photos, CreatedAt: createdAt}
	if err := a.record(event); err != nil {
		return Tool{}, err
	}
	// Intent: State only reflects evidence that has reached stable storage.
	// Source: DI-pending-mint-ex7-001.
	if err := a.applyEvent(event); err != nil {
		return Tool{}, err
	}
	return *tool, nil
}

func validatePhotos(photos []Photo) error {
	if len(photos) > 3 {
		return errors.New("at most three photos may be attached to one observation")
	}
	for _, photo := range photos {
		if !strings.HasPrefix(photo.DataURL, "data:image/") {
			return errors.New("photo must be an image")
		}
		if len(photo.DataURL) > maxPhotoDataURLBytes {
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
	createdAt := time.Now().UTC()
	event := Event{Type: "clear-safety-hold", ToolID: toolID, ActorID: stewardID, Text: assessment, CreatedAt: createdAt}
	if err := a.record(event); err != nil {
		return Tool{}, err
	}
	// Intent: State only reflects evidence that has reached stable storage.
	// Source: DI-pending-mint-ex7-001.
	if err := a.applyEvent(event); err != nil {
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
	if !a.isQualifiedForTool(memberID, *tool) {
		return Tool{}, errors.New("member does not hold this tool's required qualification")
	}
	if !dueAt.After(time.Now()) {
		return Tool{}, errors.New("return deadline must be in the future")
	}
	loan := &Loan{
		MemberID:      memberID,
		DueAt:         dueAt.UTC(),
		CreatedAt:     time.Now().UTC(),
		TermsComplete: true,
	}
	for _, candidate := range a.state.Areas {
		if candidate.ID == tool.AreaID {
			loan.PolicyVersion = candidate.PolicyVersion
			loan.Policy = candidate.Policy
			break
		}
	}
	if loan.PolicyVersion == "" || loan.Policy == "" {
		return Tool{}, errors.New("tool area has no loan policy")
	}
	createdAt := time.Now().UTC()
	loan.CreatedAt = createdAt
	event := Event{Type: "loan", ToolID: toolID, ActorID: memberID, Loan: loan, CreatedAt: createdAt}
	if err := a.record(event); err != nil {
		return Tool{}, err
	}
	// Intent: State only reflects evidence that has reached stable storage.
	// Source: DI-pending-mint-ex7-001.
	if err := a.applyEvent(event); err != nil {
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
	createdAt := time.Now().UTC()
	event := Event{Type: "return", ToolID: toolID, ActorID: memberID, Text: condition, CreatedAt: createdAt}
	if err := a.record(event); err != nil {
		return Tool{}, err
	}
	// Intent: State only reflects evidence that has reached stable storage.
	// Source: DI-pending-mint-ex7-001.
	if err := a.applyEvent(event); err != nil {
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
		if event.Loan == nil {
			if event.DueAt.IsZero() {
				return errors.New("legacy loan event lacks a return deadline")
			}
			// Intent: Retain known legacy loan facts while exposing unavailable terms
			// instead of inferring a current policy. Source: DI-pending-mint-ex7-004.
			tool.ActiveLoan = &Loan{MemberID: event.ActorID, DueAt: event.DueAt, CreatedAt: event.CreatedAt, TermsComplete: false}
			tool.Observations = append(tool.Observations, Observation{ID: fmt.Sprintf("replay-%d", len(tool.Observations)+1), ToolID: event.ToolID, ReporterID: event.ActorID, Text: "Off-site loan replayed with accepted terms unavailable. Return promised by " + event.DueAt.Format(time.RFC822), CreatedAt: event.CreatedAt})
			break
		}
		if event.Loan.MemberID != event.ActorID || !event.Loan.TermsComplete || event.Loan.PolicyVersion == "" || event.Loan.Policy == "" {
			return errors.New("loan event lacks a complete accepted policy snapshot")
		}
		loan := *event.Loan
		tool.ActiveLoan = &loan
		tool.Observations = append(tool.Observations, Observation{ID: fmt.Sprintf("replay-%d", len(tool.Observations)+1), ToolID: event.ToolID, ReporterID: event.ActorID, Text: "Off-site loan accepted. Return promised by " + loan.DueAt.Format(time.RFC822), CreatedAt: event.CreatedAt})
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
func (a *App) isQualifiedForTool(memberID string, tool Tool) bool {
	for _, qualification := range a.state.Qualifications {
		if qualification.MemberID == memberID && qualification.AreaID == tool.AreaID && qualification.Scope == tool.RequiredQualification && qualification.Status == "accepted" {
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
