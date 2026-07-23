package service

type LocalEmbodimentRequest struct {
	Type                 string            `json:"type"`
	Operation            string            `json:"operation,omitempty"`
	Method               string            `json:"method,omitempty"`
	Path                 string            `json:"path,omitempty"`
	Headers              map[string]string `json:"headers,omitempty"`
	Body                 string            `json:"body,omitempty"`
	BodyBase64           string            `json:"body_base64,omitempty"`
	Actor                string            `json:"actor,omitempty"`
	Kind                 string            `json:"kind,omitempty"`
	Name                 string            `json:"name,omitempty"`
	Title                string            `json:"title,omitempty"`
	Summary              string            `json:"summary,omitempty"`
	Notes                string            `json:"notes,omitempty"`
	Tags                 []string          `json:"tags,omitempty"`
	ParentID             string            `json:"parent_id,omitempty"`
	PlaceID              string            `json:"place_id,omitempty"`
	ResourceIDs          []string          `json:"resource_ids,omitempty"`
	Machine              string            `json:"machine,omitempty"`
	Location             string            `json:"location,omitempty"`
	Facts                map[string]string `json:"facts,omitempty"`
	Role                 string            `json:"role,omitempty"`
	Decision             string            `json:"decision,omitempty"`
	RoleKeys             []string          `json:"role_keys,omitempty"`
	RunID                string            `json:"run_id,omitempty"`
	EntityType           string            `json:"entity_type,omitempty"`
	EntityID             string            `json:"entity_id,omitempty"`
	SearchOptions        *SearchOptions    `json:"search_options,omitempty"`
	ItemID               string            `json:"item_id,omitempty"`
	Revision             int               `json:"revision,omitempty"`
	Outcome              string            `json:"outcome,omitempty"`
	ResponsibilityIDs    []string          `json:"responsibility_ids,omitempty"`
	AttachmentName       string            `json:"attachment_name,omitempty"`
	AttachmentBodyBase64 string            `json:"attachment_body_base64,omitempty"`
	ParticipantID        string            `json:"participant_id,omitempty"`
	DisplayName          string            `json:"display_name,omitempty"`
	Color                string            `json:"color,omitempty"`
	Cursor               int               `json:"cursor,omitempty"`
	Head                 int               `json:"head,omitempty"`
	Typing               bool              `json:"typing,omitempty"`
	BaseVersion          int               `json:"base_version,omitempty"`
	UpdateBody           bool              `json:"update_body,omitempty"`
}

type LocalEmbodimentResponse struct {
	Type    string         `json:"type"`
	Status  int            `json:"status,omitempty"`
	Headers map[string]any `json:"headers,omitempty"`
	Body    string         `json:"body,omitempty"`
	State   LiveItemState  `json:"state"`
	Message string         `json:"message,omitempty"`
}
