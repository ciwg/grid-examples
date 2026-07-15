package crdt

type Message struct {
	Kind          string `cbor:"kind"`
	DocumentID    string `cbor:"document_id"`
	Author        string `cbor:"author"`
	ParticipantID string `cbor:"participant_id"`
	RecipientID   string `cbor:"recipient_id,omitempty"`
	ChangeBytes   []byte `cbor:"change_bytes"`
	Lamport       uint64 `cbor:"lamport"`
	Embodiment    string `cbor:"embodiment,omitempty"`
}

type SyncRecord struct {
	Offset        uint64 `json:"offset"`
	EnvelopeCID   string `json:"envelope_cid"`
	ParticipantID string `json:"participant_id"`
	RecipientID   string `json:"recipient_id,omitempty"`
	Author        string `json:"author"`
	MessageBase64 string `json:"message_base64"`
	Embodiment    string `json:"embodiment,omitempty"`
	ReceivedAt    string `json:"received_at"`
}
