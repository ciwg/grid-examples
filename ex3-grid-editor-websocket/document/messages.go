package document

import "github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocol"

type Message struct {
	Kind          string `cbor:"kind"`
	DocumentID    string `cbor:"document_id"`
	Content       string `cbor:"content"`
	ContentCID    string `cbor:"content_cid"`
	Lamport       uint64 `cbor:"lamport"`
	Author        string `cbor:"author"`
	ParticipantID string `cbor:"participant_id,omitempty"`
	Embodiment    string `cbor:"embodiment,omitempty"`
	PreviousCID   string `cbor:"previous_cid,omitempty"`
}

func NewMessage(documentID string, content string, lamport uint64, author string, participantID string, embodiment string, previousCID string) (Message, error) {
	contentCID, err := protocol.CIDForBytes([]byte(content))
	if err != nil {
		return Message{}, err
	}
	return Message{
		Kind:          "replace",
		DocumentID:    documentID,
		Content:       content,
		ContentCID:    contentCID.String(),
		Lamport:       lamport,
		Author:        author,
		ParticipantID: participantID,
		Embodiment:    embodiment,
		PreviousCID:   previousCID,
	}, nil
}
