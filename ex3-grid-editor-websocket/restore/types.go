package restore

type Message struct {
	Kind              string `cbor:"kind"`
	DocumentID        string `cbor:"document_id"`
	Author            string `cbor:"author"`
	ParticipantID     string `cbor:"participant_id"`
	SourceManifestCID string `cbor:"source_manifest_cid"`
	LiveChangeBytes   []byte `cbor:"live_change_base64"`
	RestoredAt        string `cbor:"restored_at"`
	Lamport           uint64 `cbor:"lamport"`
	Embodiment        string `cbor:"embodiment,omitempty"`
}

type Record struct {
	Offset            uint64 `json:"offset"`
	EnvelopeCID       string `json:"envelope_cid"`
	DocumentID        string `json:"document_id"`
	ParticipantID     string `json:"participant_id"`
	Author            string `json:"author"`
	SourceManifestCID string `json:"source_manifest_cid"`
	MessageBase64     string `json:"message_base64"`
	Embodiment        string `json:"embodiment,omitempty"`
	ReceivedAt        string `json:"received_at"`
}
