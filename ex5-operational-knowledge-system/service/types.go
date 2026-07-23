package service

const (
	KnowledgeKindProcedure   = "procedure"
	KnowledgeKindTraining    = "training"
	KnowledgeKindMaintenance = "maintenance"
	KnowledgeKindReceiving   = "receiving_check"
	KnowledgeKindInventory   = "inventory_audit"
)

const (
	RunKindProcedure   = "procedure"
	RunKindTraining    = "training"
	RunKindMaintenance = "maintenance"
	RunKindReceiving   = "receiving_check"
	RunKindInventory   = "inventory_audit"
)

const (
	DecisionApproved = "approved"
	DecisionRejected = "rejected"
	DecisionNoted    = "noted"
)

const (
	ItemStatusDraft      = "draft"
	ItemStatusApproved   = "approved"
	ItemStatusSuperseded = "superseded"
)

type Meta struct {
	DataRoot                    string                         `json:"data_root"`
	LocalPeerID                 string                         `json:"local_peer_id"`
	KnowledgeKinds              []string                       `json:"knowledge_kinds"`
	RunKinds                    []string                       `json:"run_kinds"`
	ApprovalDecisions           []string                       `json:"approval_decisions"`
	ItemStatuses                []string                       `json:"item_statuses"`
	KnowledgeItemPCID           string                         `json:"knowledge_item_pcid"`
	KnowledgeApprovalPCID       string                         `json:"knowledge_approval_pcid"`
	KnowledgeEvidencePCID       string                         `json:"knowledge_evidence_pcid"`
	KnowledgeLinkPCID           string                         `json:"knowledge_link_pcid"`
	KnowledgeResponsibilityPCID string                         `json:"knowledge_responsibility_pcid"`
	OperationalRunPCID          string                         `json:"operational_run_pcid"`
	OperationalPlacePCID        string                         `json:"operational_place_pcid"`
	OperationalResourcePCID     string                         `json:"operational_resource_pcid"`
	PeerExchangeFormat          string                         `json:"peer_exchange_format"`
	PeerExchangeFamilies        []string                       `json:"peer_exchange_families"`
	RelayFeedFormat             string                         `json:"relay_feed_format"`
	RelayFeedFamilies           []string                       `json:"relay_feed_families"`
	CASObjectsEnabled           bool                           `json:"cas_objects_enabled"`
	CASAttachmentBlobsEnabled   bool                           `json:"cas_attachment_blobs_enabled"`
	CASDraftBodiesEnabled       bool                           `json:"cas_draft_bodies_enabled"`
	RelayBlobTransferEnabled    bool                           `json:"relay_blob_transfer_enabled"`
	LiveDraftWebSocketEnabled   bool                           `json:"live_draft_websocket_enabled"`
	LocalUnixSocketEnabled      bool                           `json:"local_unix_socket_enabled"`
	LocalUnixSocketPath         string                         `json:"local_unix_socket_path"`
	Embodiments                 map[string]EmbodimentTransport `json:"embodiments"`
}

type EmbodimentTransport struct {
	PrimaryAdapter      string   `json:"primary_adapter"`
	LiveDraftTransport  string   `json:"live_draft_transport,omitempty"`
	FallbackTransports  []string `json:"fallback_transports,omitempty"`
	CompatibilityMode   string   `json:"compatibility_mode"`
	LocalUnixSocketPath string   `json:"local_unix_socket_path,omitempty"`
}

type RelayMeta struct {
	DataRoot                   string   `json:"data_root"`
	ServiceName                string   `json:"service_name"`
	RoutePrefix                string   `json:"route_prefix"`
	RelayFeedFormat            string   `json:"relay_feed_format"`
	RelayFeedFamilies          []string `json:"relay_feed_families"`
	RelayBlobTransferEnabled   bool     `json:"relay_blob_transfer_enabled"`
	PublishRequiresStagedBlobs bool     `json:"publish_requires_staged_blobs"`
}

type RelayFeedRequest struct {
	KnownOrigins map[string]uint64 `json:"known_origins,omitempty"`
}

type RelayFeedBatch struct {
	Format                         string                                `json:"format"`
	ExportedAt                     string                                `json:"exported_at"`
	Implementation                 string                                `json:"implementation"`
	ExportingPeerID                string                                `json:"exporting_peer_id"`
	KnowledgeItemPCID              string                                `json:"knowledge_item_pcid"`
	KnowledgeApprovalPCID          string                                `json:"knowledge_approval_pcid"`
	KnowledgeEvidencePCID          string                                `json:"knowledge_evidence_pcid"`
	KnowledgeLinkPCID              string                                `json:"knowledge_link_pcid"`
	KnowledgeResponsibilityPCID    string                                `json:"knowledge_responsibility_pcid"`
	OperationalRunPCID             string                                `json:"operational_run_pcid"`
	OperationalPlacePCID           string                                `json:"operational_place_pcid"`
	OperationalResourcePCID        string                                `json:"operational_resource_pcid"`
	Events                         []OperationalEvent                    `json:"events"`
	KnowledgeItemRecords           []SignedKnowledgeItemRecord           `json:"knowledge_item_records"`
	KnowledgeApprovalRecords       []SignedKnowledgeApprovalRecord       `json:"knowledge_approval_records"`
	KnowledgeEvidenceRecords       []SignedKnowledgeEvidenceRecord       `json:"knowledge_evidence_records"`
	OperationalRunRecords          []SignedOperationalRunRecord          `json:"operational_run_records"`
	OperationalPlaceRecords        []SignedOperationalPlaceRecord        `json:"operational_place_records"`
	OperationalResourceRecords     []SignedOperationalResourceRecord     `json:"operational_resource_records"`
	KnowledgeLinkRecords           []SignedKnowledgeLinkRecord           `json:"knowledge_link_records"`
	KnowledgeResponsibilityRecords []SignedKnowledgeResponsibilityRecord `json:"knowledge_responsibility_records"`
	RequiredBlobCIDs               []string                              `json:"required_blob_cids,omitempty"`
}

type RelayFeedImportResult struct {
	ImportedEvents               int                       `json:"imported_events"`
	ImportedKnowledgeItems       int                       `json:"imported_knowledge_items"`
	ImportedKnowledgeApprovals   int                       `json:"imported_knowledge_approvals"`
	ImportedKnowledgeEvidence    int                       `json:"imported_knowledge_evidence"`
	ImportedOperationalRuns      int                       `json:"imported_operational_runs"`
	ImportedOperationalPlaces    int                       `json:"imported_operational_places"`
	ImportedOperationalResources int                       `json:"imported_operational_resources"`
	ImportedKnowledgeLinks       int                       `json:"imported_knowledge_links"`
	ImportedResponsibilities     int                       `json:"imported_responsibilities"`
	MissingBlobCIDs              []string                  `json:"missing_blob_cids,omitempty"`
	UnresolvedReferences         []PeerExchangeImportIssue `json:"unresolved_references"`
}

type RelayPublishResult struct {
	PublishedEvents               int      `json:"published_events"`
	PublishedKnowledgeItems       int      `json:"published_knowledge_items"`
	PublishedKnowledgeApprovals   int      `json:"published_knowledge_approvals"`
	PublishedKnowledgeEvidence    int      `json:"published_knowledge_evidence"`
	PublishedOperationalRuns      int      `json:"published_operational_runs"`
	PublishedOperationalPlaces    int      `json:"published_operational_places"`
	PublishedOperationalResources int      `json:"published_operational_resources"`
	PublishedKnowledgeLinks       int      `json:"published_knowledge_links"`
	PublishedResponsibilities     int      `json:"published_responsibilities"`
	MissingBlobCIDs               []string `json:"missing_blob_cids,omitempty"`
}

type PeerExchangeBundle struct {
	Format                         string                                `json:"format"`
	ExportedAt                     string                                `json:"exported_at"`
	Implementation                 string                                `json:"implementation"`
	ExportingPeerID                string                                `json:"exporting_peer_id"`
	KnowledgeItemPCID              string                                `json:"knowledge_item_pcid"`
	KnowledgeApprovalPCID          string                                `json:"knowledge_approval_pcid"`
	KnowledgeEvidencePCID          string                                `json:"knowledge_evidence_pcid"`
	KnowledgeLinkPCID              string                                `json:"knowledge_link_pcid"`
	KnowledgeResponsibilityPCID    string                                `json:"knowledge_responsibility_pcid"`
	OperationalRunPCID             string                                `json:"operational_run_pcid"`
	OperationalPlacePCID           string                                `json:"operational_place_pcid"`
	OperationalResourcePCID        string                                `json:"operational_resource_pcid"`
	Events                         []OperationalEvent                    `json:"events"`
	KnowledgeItemRecords           []SignedKnowledgeItemRecord           `json:"knowledge_item_records"`
	KnowledgeApprovalRecords       []SignedKnowledgeApprovalRecord       `json:"knowledge_approval_records"`
	KnowledgeEvidenceRecords       []SignedKnowledgeEvidenceRecord       `json:"knowledge_evidence_records"`
	OperationalRunRecords          []SignedOperationalRunRecord          `json:"operational_run_records"`
	OperationalPlaceRecords        []SignedOperationalPlaceRecord        `json:"operational_place_records"`
	OperationalResourceRecords     []SignedOperationalResourceRecord     `json:"operational_resource_records"`
	KnowledgeLinkRecords           []SignedKnowledgeLinkRecord           `json:"knowledge_link_records"`
	KnowledgeResponsibilityRecords []SignedKnowledgeResponsibilityRecord `json:"knowledge_responsibility_records"`
	CASBlobObjects                 map[string]string                     `json:"cas_blob_objects,omitempty"`
}

type PeerExchangeImportIssue struct {
	RecordType string `json:"record_type"`
	RecordID   string `json:"record_id"`
	Reason     string `json:"reason"`
}

type PeerExchangeImportResult struct {
	ImportedEvents               int                       `json:"imported_events"`
	ImportedKnowledgeItems       int                       `json:"imported_knowledge_items"`
	ImportedKnowledgeApprovals   int                       `json:"imported_knowledge_approvals"`
	ImportedKnowledgeEvidence    int                       `json:"imported_knowledge_evidence"`
	ImportedOperationalRuns      int                       `json:"imported_operational_runs"`
	ImportedOperationalPlaces    int                       `json:"imported_operational_places"`
	ImportedOperationalResources int                       `json:"imported_operational_resources"`
	ImportedKnowledgeLinks       int                       `json:"imported_knowledge_links"`
	ImportedResponsibilities     int                       `json:"imported_responsibilities"`
	UnresolvedReferences         []PeerExchangeImportIssue `json:"unresolved_references"`
}

type Place struct {
	ID            string             `json:"id"`
	AliasID       string             `json:"alias_id,omitempty"`
	Kind          string             `json:"kind"`
	Name          string             `json:"name"`
	Summary       string             `json:"summary"`
	ParentID      string             `json:"parent_id"`
	Tags          []string           `json:"tags"`
	CreatedAt     string             `json:"created_at"`
	UpdatedAt     string             `json:"updated_at"`
	ChildPlaceIDs []string           `json:"child_place_ids"`
	ResourceIDs   []string           `json:"resource_ids"`
	RelatedRuns   []RunRecord        `json:"related_runs"`
	Links         []Link             `json:"links"`
	Timeline      []OperationalEvent `json:"timeline"`
}

type Resource struct {
	ID          string             `json:"id"`
	AliasID     string             `json:"alias_id,omitempty"`
	Kind        string             `json:"kind"`
	Name        string             `json:"name"`
	Summary     string             `json:"summary"`
	PlaceID     string             `json:"place_id"`
	Tags        []string           `json:"tags"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	RelatedRuns []RunRecord        `json:"related_runs"`
	Links       []Link             `json:"links"`
	Timeline    []OperationalEvent `json:"timeline"`
}

type Responsibility struct {
	ID             string             `json:"id"`
	AliasID        string             `json:"alias_id,omitempty"`
	Title          string             `json:"title"`
	Summary        string             `json:"summary"`
	Team           string             `json:"team"`
	Tags           []string           `json:"tags"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
	LinkedItemIDs  []string           `json:"linked_item_ids"`
	LinkedRunIDs   []string           `json:"linked_run_ids"`
	RelatedRuns    []RunRecord        `json:"related_runs"`
	LinkedRoleKeys []string           `json:"linked_role_keys"`
	Links          []Link             `json:"links"`
	Timeline       []OperationalEvent `json:"timeline"`
}

type KnowledgeItem struct {
	ID                string              `json:"id"`
	AliasID           string              `json:"alias_id,omitempty"`
	Kind              string              `json:"kind"`
	Status            string              `json:"status"`
	Title             string              `json:"title"`
	Summary           string              `json:"summary"`
	Tags              []string            `json:"tags"`
	ResponsibilityIDs []string            `json:"responsibility_ids"`
	CreatedAt         string              `json:"created_at"`
	UpdatedAt         string              `json:"updated_at"`
	CurrentRevision   int                 `json:"current_revision"`
	WorkingBody       string              `json:"working_body"`
	WorkingVersion    int                 `json:"working_version"`
	WorkingUpdatedAt  string              `json:"working_updated_at"`
	Revisions         []KnowledgeRevision `json:"revisions"`
	RelatedRuns       []RunRecord         `json:"related_runs"`
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
	AliasID           string             `json:"alias_id,omitempty"`
	Kind              string             `json:"kind"`
	ItemID            string             `json:"item_id"`
	ItemKind          string             `json:"item_kind"`
	Revision          int                `json:"revision"`
	Actor             string             `json:"actor"`
	Outcome           string             `json:"outcome"`
	Notes             string             `json:"notes"`
	PlaceID           string             `json:"place_id"`
	ResourceIDs       []string           `json:"resource_ids"`
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
	AliasID        string            `json:"alias_id,omitempty"`
	Summary        string            `json:"summary"`
	Facts          map[string]string `json:"facts"`
	AttachmentName string            `json:"attachment_name"`
	AttachmentPath string            `json:"attachment_path"`
	AttachmentCID  string            `json:"attachment_cid"`
	AttachmentSize int64             `json:"attachment_size"`
	Actor          string            `json:"actor"`
	CreatedAt      string            `json:"created_at"`
}

type Approval struct {
	ID         string `json:"id"`
	AliasID    string `json:"alias_id,omitempty"`
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
	AliasID   string `json:"alias_id,omitempty"`
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
	Places           int `json:"places"`
	Resources        int `json:"resources"`
	Procedures       int `json:"procedures"`
	TrainingItems    int `json:"training_items"`
	MaintenanceItems int `json:"maintenance_items"`
	ReceivingItems   int `json:"receiving_items"`
	InventoryItems   int `json:"inventory_items"`
	ProcedureRuns    int `json:"procedure_runs"`
	TrainingRuns     int `json:"training_runs"`
	MaintenanceRuns  int `json:"maintenance_runs"`
	ReceivingRuns    int `json:"receiving_runs"`
	InventoryRuns    int `json:"inventory_runs"`
	Approvals        int `json:"approvals"`
	Evidence         int `json:"evidence"`
	Links            int `json:"links"`
}

type ProblemReview struct {
	ProblemRuns    int                  `json:"problem_runs"`
	PlaceGroups    []ProblemReviewGroup `json:"place_groups"`
	ResourceGroups []ProblemReviewGroup `json:"resource_groups"`
}

type ProblemReviewGroup struct {
	GroupType         string      `json:"group_type"`
	GroupID           string      `json:"group_id"`
	Kind              string      `json:"kind"`
	Name              string      `json:"name"`
	ProblemCount      int         `json:"problem_count"`
	ReceivingProblems int         `json:"receiving_problems"`
	InventoryProblems int         `json:"inventory_problems"`
	HighlightExamples []string    `json:"highlights"`
	Runs              []RunRecord `json:"runs"`
}

// Intent: Keep one shared search shape for browser forms, drilldown actions,
// and HTTP query parsing so context history filters stay aligned. Source:
// DI-vafuk
type SearchOptions struct {
	Query            string `json:"query"`
	Kind             string `json:"kind"`
	Status           string `json:"status"`
	Outcome          string `json:"outcome"`
	PlaceID          string `json:"place_id"`
	ResourceID       string `json:"resource_id"`
	ResponsibilityID string `json:"responsibility_id"`
	Problem          bool   `json:"problem"`
}

type LivePresence struct {
	ParticipantID string `json:"participant_id"`
	DisplayName   string `json:"display_name"`
	Color         string `json:"color"`
	Cursor        int    `json:"cursor"`
	Head          int    `json:"head"`
	Typing        bool   `json:"typing"`
	LastSeenAt    string `json:"last_seen_at"`
}

type LiveItemState struct {
	ItemID          string         `json:"item_id"`
	Title           string         `json:"title"`
	Status          string         `json:"status"`
	Body            string         `json:"body"`
	Version         int            `json:"version"`
	CurrentRevision int            `json:"current_revision"`
	Participants    []LivePresence `json:"participants"`
}

type OperationalEvent struct {
	Sequence          uint64            `json:"sequence"`
	OriginPeerID      string            `json:"origin_peer_id"`
	OriginSequence    uint64            `json:"origin_sequence"`
	Timestamp         string            `json:"timestamp"`
	EntityType        string            `json:"entity_type"`
	EntityID          string            `json:"entity_id"`
	DisplayID         string            `json:"display_id,omitempty"`
	CanonicalID       string            `json:"canonical_id,omitempty"`
	Type              string            `json:"type"`
	Actor             string            `json:"actor"`
	Name              string            `json:"name"`
	Title             string            `json:"title"`
	Summary           string            `json:"summary"`
	Body              string            `json:"body"`
	Kind              string            `json:"kind"`
	Status            string            `json:"status"`
	Tags              []string          `json:"tags"`
	Team              string            `json:"team"`
	ParentID          string            `json:"parent_id"`
	PlaceID           string            `json:"place_id"`
	ResourceIDs       []string          `json:"resource_ids"`
	ResponsibilityIDs []string          `json:"responsibility_ids"`
	RoleKeys          []string          `json:"role_keys"`
	Revision          int               `json:"revision"`
	Outcome           string            `json:"outcome"`
	Notes             string            `json:"notes"`
	Machine           string            `json:"machine"`
	Location          string            `json:"location"`
	AttachmentName    string            `json:"attachment_name"`
	AttachmentPath    string            `json:"attachment_path"`
	AttachmentCID     string            `json:"attachment_cid"`
	AttachmentSize    int64             `json:"attachment_size"`
	EvidenceID        string            `json:"evidence_id"`
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
