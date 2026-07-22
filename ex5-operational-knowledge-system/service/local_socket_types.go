package service

type LocalEmbodimentRequest struct {
	Type          string            `json:"type"`
	Method        string            `json:"method,omitempty"`
	Path          string            `json:"path,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Body          string            `json:"body,omitempty"`
	BodyBase64    string            `json:"body_base64,omitempty"`
	ItemID        string            `json:"item_id,omitempty"`
	ParticipantID string            `json:"participant_id,omitempty"`
	DisplayName   string            `json:"display_name,omitempty"`
	Color         string            `json:"color,omitempty"`
	Cursor        int               `json:"cursor,omitempty"`
	Head          int               `json:"head,omitempty"`
	Typing        bool              `json:"typing,omitempty"`
	BaseVersion   int               `json:"base_version,omitempty"`
	UpdateBody    bool              `json:"update_body,omitempty"`
}

type LocalEmbodimentResponse struct {
	Type    string         `json:"type"`
	Status  int            `json:"status,omitempty"`
	Headers map[string]any `json:"headers,omitempty"`
	Body    string         `json:"body,omitempty"`
	State   LiveItemState  `json:"state"`
	Message string         `json:"message,omitempty"`
}
