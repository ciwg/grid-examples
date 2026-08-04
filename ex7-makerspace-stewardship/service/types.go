package service

import "time"

type Member struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Initials string `json:"initials"`
}

type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Authority struct {
	MemberID string   `json:"memberId"`
	AreaID   string   `json:"areaId"`
	Scopes   []string `json:"scopes"`
}

type Qualification struct {
	MemberID string `json:"memberId"`
	AreaID   string `json:"areaId"`
	IssuedBy string `json:"issuedBy"`
	Status   string `json:"status"`
}

type Observation struct {
	ID         string    `json:"id"`
	ToolID     string    `json:"toolId"`
	ReporterID string    `json:"reporterId"`
	Text       string    `json:"text"`
	SafetyHold bool      `json:"safetyHold"`
	Photos     []Photo   `json:"photos,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Photo struct {
	Name    string `json:"name"`
	DataURL string `json:"dataUrl"`
}

type Tool struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	AreaID       string        `json:"areaId"`
	OffSiteLoan  bool          `json:"offSiteLoan"`
	Condition    string        `json:"condition"`
	SafetyHold   bool          `json:"safetyHold"`
	Observations []Observation `json:"observations"`
	ActiveLoan   *Loan         `json:"activeLoan,omitempty"`
}

type Loan struct {
	MemberID      string    `json:"memberId"`
	DueAt         time.Time `json:"dueAt"`
	CreatedAt     time.Time `json:"createdAt"`
	PolicyVersion string    `json:"policyVersion"`
	Policy        string    `json:"policy"`
}

type State struct {
	Members        []Member        `json:"members"`
	Areas          []Area          `json:"areas"`
	Authorities    []Authority     `json:"authorities"`
	Qualifications []Qualification `json:"qualifications"`
	Tools          []Tool          `json:"tools"`
}
