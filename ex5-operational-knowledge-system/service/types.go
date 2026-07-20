package service

const (
	KnowledgeKindProcedure   = "procedure"
	KnowledgeKindTraining    = "training"
	KnowledgeKindMaintenance = "maintenance"
)

const (
	RunKindProcedure   = "procedure"
	RunKindTraining    = "training"
	RunKindMaintenance = "maintenance"
)

const (
	DecisionApproved = "approved"
	DecisionRejected = "rejected"
	DecisionNoted    = "noted"
)

type Meta struct {
	DataRoot          string   `json:"data_root"`
	KnowledgeKinds    []string `json:"knowledge_kinds"`
	RunKinds          []string `json:"run_kinds"`
	ApprovalDecisions []string `json:"approval_decisions"`
}

type Responsibility struct {
	ID             string             `json:"id"`
	Title          string             `json:"title"`
	Summary        string             `json:"summary"`
	Team           string             `json:"team"`
	Tags           []string           `json:"tags"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
	LinkedItemIDs  []string           `json:"linked_item_ids"`
	LinkedRunIDs   []string           `json:"linked_run_ids"`
	LinkedRoleKeys []string           `json:"linked_role_keys"`
	Timeline       []OperationalEvent `json:"timeline"`
}

type KnowledgeItem struct {
	ID                string              `json:"id"`
	Kind              string              `json:"kind"`
	Title             string              `json:"title"`
	Summary           string              `json:"summary"`
	Tags              []string            `json:"tags"`
	ResponsibilityIDs []string            `json:"responsibility_ids"`
	CreatedAt         string              `json:"created_at"`
	UpdatedAt         string              `json:"updated_at"`
	CurrentRevision   int                 `json:"current_revision"`
	Revisions         []KnowledgeRevision `json:"revisions"`
	Approvals         []Approval          `json:"approvals"`
	Links             []Link              `json:"links"`
	Timeline          []OperationalEvent  `json:"timeline"`
}

type KnowledgeRevision struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	Body      string   `json:"body"`
	Tags      []string `json:"tags"`
	Author    string   `json:"author"`
	CreatedAt string   `json:"created_at"`
}

type RunRecord struct {
	ID                string             `json:"id"`
	Kind              string             `json:"kind"`
	ItemID            string             `json:"item_id"`
	ItemKind          string             `json:"item_kind"`
	Revision          int                `json:"revision"`
	Actor             string             `json:"actor"`
	Outcome           string             `json:"outcome"`
	Notes             string             `json:"notes"`
	Machine           string             `json:"machine"`
	Location          string             `json:"location"`
	ResponsibilityIDs []string           `json:"responsibility_ids"`
	CreatedAt         string             `json:"created_at"`
	UpdatedAt         string             `json:"updated_at"`
	Evidence          []Evidence         `json:"evidence"`
	Approvals         []Approval         `json:"approvals"`
	Links             []Link             `json:"links"`
	Timeline          []OperationalEvent `json:"timeline"`
}

type Evidence struct {
	ID             string            `json:"id"`
	Summary        string            `json:"summary"`
	Facts          map[string]string `json:"facts"`
	AttachmentName string            `json:"attachment_name"`
	AttachmentPath string            `json:"attachment_path"`
	AttachmentSize int64             `json:"attachment_size"`
	Actor          string            `json:"actor"`
	CreatedAt      string            `json:"created_at"`
}

type Approval struct {
	ID         string `json:"id"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Revision   int    `json:"revision"`
	RunID      string `json:"run_id"`
	Role       string `json:"role"`
	Decision   string `json:"decision"`
	Actor      string `json:"actor"`
	Notes      string `json:"notes"`
	CreatedAt  string `json:"created_at"`
}

type Link struct {
	ID        string `json:"id"`
	FromType  string `json:"from_type"`
	FromID    string `json:"from_id"`
	ToType    string `json:"to_type"`
	ToID      string `json:"to_id"`
	Relation  string `json:"relation"`
	Notes     string `json:"notes"`
	Actor     string `json:"actor"`
	CreatedAt string `json:"created_at"`
}

type Dashboard struct {
	Responsibilities int `json:"responsibilities"`
	Procedures       int `json:"procedures"`
	TrainingItems    int `json:"training_items"`
	MaintenanceItems int `json:"maintenance_items"`
	ProcedureRuns    int `json:"procedure_runs"`
	TrainingRuns     int `json:"training_runs"`
	MaintenanceRuns  int `json:"maintenance_runs"`
	Approvals        int `json:"approvals"`
	Evidence         int `json:"evidence"`
	Links            int `json:"links"`
}

type OperationalEvent struct {
	Sequence          uint64            `json:"sequence"`
	Timestamp         string            `json:"timestamp"`
	EntityType        string            `json:"entity_type"`
	EntityID          string            `json:"entity_id"`
	Type              string            `json:"type"`
	Actor             string            `json:"actor"`
	Title             string            `json:"title"`
	Summary           string            `json:"summary"`
	Body              string            `json:"body"`
	Kind              string            `json:"kind"`
	Tags              []string          `json:"tags"`
	Team              string            `json:"team"`
	ResponsibilityIDs []string          `json:"responsibility_ids"`
	RoleKeys          []string          `json:"role_keys"`
	Revision          int               `json:"revision"`
	Outcome           string            `json:"outcome"`
	Notes             string            `json:"notes"`
	Machine           string            `json:"machine"`
	Location          string            `json:"location"`
	AttachmentName    string            `json:"attachment_name"`
	AttachmentPath    string            `json:"attachment_path"`
	AttachmentSize    int64             `json:"attachment_size"`
	Facts             map[string]string `json:"facts"`
	TargetType        string            `json:"target_type"`
	TargetID          string            `json:"target_id"`
	RunID             string            `json:"run_id"`
	Decision          string            `json:"decision"`
	Role              string            `json:"role"`
	FromType          string            `json:"from_type"`
	FromID            string            `json:"from_id"`
	ToType            string            `json:"to_type"`
	ToID              string            `json:"to_id"`
	Relation          string            `json:"relation"`
}
