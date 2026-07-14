package publish

type Message struct {
	Kind              string `cbor:"kind"`
	DocumentID        string `cbor:"document_id"`
	Author            string `cbor:"author"`
	ParticipantID     string `cbor:"participant_id"`
	SourceKind        string `cbor:"source_kind"`
	SourceVersionID   string `cbor:"source_version_id,omitempty"`
	SourceVersionName string `cbor:"source_version_name,omitempty"`
	Title             string `cbor:"title"`
	Summary           string `cbor:"summary"`
	TextCID           string `cbor:"text_cid"`
	ReplicaCID        string `cbor:"replica_cid"`
	PublishedAt       string `cbor:"published_at"`
	Lamport           uint64 `cbor:"lamport"`
	Embodiment        string `cbor:"embodiment,omitempty"`
}

type Record struct {
	Offset            uint64 `json:"offset"`
	EnvelopeCID       string `json:"envelope_cid"`
	DocumentID        string `json:"document_id"`
	Author            string `json:"author"`
	ParticipantID     string `json:"participant_id"`
	SourceKind        string `json:"source_kind"`
	SourceVersionID   string `json:"source_version_id,omitempty"`
	SourceVersionName string `json:"source_version_name,omitempty"`
	Title             string `json:"title"`
	Summary           string `json:"summary"`
	TextCID           string `json:"text_cid"`
	ReplicaCID        string `json:"replica_cid"`
	PublishedAt       string `json:"published_at"`
	Embodiment        string `json:"embodiment,omitempty"`
	ReceivedAt        string `json:"received_at"`
}

type Resolved struct {
	Record        Record `json:"record"`
	TextBase64    string `json:"text_base64"`
	ReplicaBase64 string `json:"replica_base64"`
}
