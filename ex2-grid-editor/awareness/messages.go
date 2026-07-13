package awareness

type Message struct {
	Kind          string `cbor:"kind"`
	DocumentID    string `cbor:"document_id"`
	Author        string `cbor:"author"`
	ParticipantID string `cbor:"participant_id"`
	DisplayName   string `cbor:"display_name"`
	Color         string `cbor:"color"`
	Cursor        int    `cbor:"cursor"`
	Head          int    `cbor:"head"`
	Typing        bool   `cbor:"typing"`
	Lamport       uint64 `cbor:"lamport"`
	Embodiment    string `cbor:"embodiment,omitempty"`
}
