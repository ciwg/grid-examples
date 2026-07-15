package document

import (
	"fmt"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocol"
)

type State struct {
	DocumentID string `json:"document_id"`
	Content    string `json:"content"`
	ContentCID string `json:"content_cid"`
	MessageCID string `json:"message_cid"`
	Lamport    uint64 `json:"lamport"`
	Author     string `json:"author"`
}

func Apply(current State, message Message, messageCID string) (State, bool, error) {
	if message.Kind != "replace" {
		return current, false, fmt.Errorf("unexpected document message kind %q", message.Kind)
	}
	contentCID, err := protocol.CIDForBytes([]byte(message.Content))
	if err != nil {
		return current, false, fmt.Errorf("content cid: %w", err)
	}
	if contentCID.String() != message.ContentCID {
		return current, false, fmt.Errorf("content cid mismatch")
	}
	if !wins(current.Lamport, current.Author, current.MessageCID, message.Lamport, message.Author, messageCID) {
		return current, false, nil
	}
	return State{
		DocumentID: message.DocumentID,
		Content:    message.Content,
		ContentCID: message.ContentCID,
		MessageCID: messageCID,
		Lamport:    message.Lamport,
		Author:     message.Author,
	}, true, nil
}

// Intent: Keep document convergence deterministic across replay and mixed-host
// delivery by using an explicit `(lamport, author, message_cid)` ordering rule
// instead of adapter-local timing behavior. Source: DI-jilin
func wins(currentLamport uint64, currentAuthor string, currentMessageCID string, nextLamport uint64, nextAuthor string, nextMessageCID string) bool {
	if currentMessageCID == "" {
		return true
	}
	if nextLamport != currentLamport {
		return nextLamport > currentLamport
	}
	if nextAuthor != currentAuthor {
		return nextAuthor > currentAuthor
	}
	return nextMessageCID > currentMessageCID
}
