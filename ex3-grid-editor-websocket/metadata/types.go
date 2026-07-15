package metadata

type Message struct {
	Kind          string   `cbor:"kind"`
	DocumentID    string   `cbor:"document_id"`
	Author        string   `cbor:"author"`
	ParticipantID string   `cbor:"participant_id"`
	Title         string   `cbor:"title,omitempty"`
	Description   string   `cbor:"description,omitempty"`
	Summary       string   `cbor:"summary,omitempty"`
	Tags          []string `cbor:"tags,omitempty"`
	Collections   []string `cbor:"collections,omitempty"`
	Favorite      bool     `cbor:"favorite"`
	Archived      bool     `cbor:"archived"`
	UpdatedAt     string   `cbor:"updated_at"`
	Lamport       uint64   `cbor:"lamport"`
	Embodiment    string   `cbor:"embodiment,omitempty"`
}

type Record struct {
	Offset        uint64   `json:"offset"`
	EnvelopeCID   string   `json:"envelope_cid"`
	DocumentID    string   `json:"document_id"`
	Author        string   `json:"author"`
	ParticipantID string   `json:"participant_id"`
	Title         string   `json:"title,omitempty"`
	Description   string   `json:"description,omitempty"`
	Summary       string   `json:"summary,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Collections   []string `json:"collections,omitempty"`
	Favorite      bool     `json:"favorite"`
	Archived      bool     `json:"archived"`
	UpdatedAt     string   `json:"updated_at"`
	Embodiment    string   `json:"embodiment,omitempty"`
	ReceivedAt    string   `json:"received_at"`
	Lamport       uint64   `json:"lamport"`
}
