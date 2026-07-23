package records

type Event struct {
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

type SignedKnowledgeItemRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	ItemID         string `json:"item_id"`
	EventType      string `json:"event_type"`
	Revision       int    `json:"revision"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type SignedKnowledgeApprovalRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	ApprovalID     string `json:"approval_id"`
	TargetType     string `json:"target_type"`
	TargetID       string `json:"target_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type SignedKnowledgeEvidenceRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	EvidenceID     string `json:"evidence_id"`
	RunID          string `json:"run_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type SignedOperationalRunRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	RunID          string `json:"run_id"`
	ItemID         string `json:"item_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type SignedOperationalPlaceRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	PlaceID        string `json:"place_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type SignedOperationalResourceRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	ResourceID     string `json:"resource_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type SignedKnowledgeLinkRecord struct {
	Sequence       uint64 `json:"sequence"`
	OriginPeerID   string `json:"origin_peer_id"`
	OriginSequence uint64 `json:"origin_sequence"`
	LinkID         string `json:"link_id"`
	PCID           string `json:"pcid"`
	EnvelopeCID    string `json:"envelope_cid"`
	EnvelopeBase64 string `json:"envelope_base64"`
	RecordedAt     string `json:"recorded_at"`
	Implementation string `json:"implementation"`
}

type SignedKnowledgeResponsibilityRecord struct {
	Sequence         uint64 `json:"sequence"`
	OriginPeerID     string `json:"origin_peer_id"`
	OriginSequence   uint64 `json:"origin_sequence"`
	ResponsibilityID string `json:"responsibility_id"`
	PCID             string `json:"pcid"`
	EnvelopeCID      string `json:"envelope_cid"`
	EnvelopeBase64   string `json:"envelope_base64"`
	RecordedAt       string `json:"recorded_at"`
	Implementation   string `json:"implementation"`
}

type Signer interface {
	PeerID() string
	SignProof([]byte) ([]byte, error)
}
