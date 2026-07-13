package awareness

type PeerState struct {
	Author        string `json:"author"`
	ParticipantID string `json:"participant_id"`
	DisplayName   string `json:"display_name"`
	Color         string `json:"color"`
	Cursor        int    `json:"cursor"`
	Head          int    `json:"head"`
	Typing        bool   `json:"typing"`
	Lamport       uint64 `json:"lamport"`
	Embodiment    string `json:"embodiment"`
	MessageCID    string `json:"message_cid"`
}

type Index map[string]PeerState

func Apply(index Index, message Message, messageCID string) (Index, bool) {
	if index == nil {
		index = Index{}
	}
	key := message.ParticipantID
	if key == "" {
		key = message.Author
	}
	current, ok := index[key]
	if ok && !wins(current.Lamport, current.Author, current.MessageCID, message.Lamport, message.Author, messageCID) {
		return index, false
	}
	index[key] = PeerState{
		Author:        message.Author,
		ParticipantID: key,
		DisplayName:   message.DisplayName,
		Color:         message.Color,
		Cursor:        message.Cursor,
		Head:          message.Head,
		Typing:        message.Typing,
		Lamport:       message.Lamport,
		Embodiment:    message.Embodiment,
		MessageCID:    messageCID,
	}
	return index, true
}

// Intent: Give awareness the same deterministic latest-state selection rule as
// document updates while still keeping awareness a separate protocol family and
// projection. Source: DI-tofug; DI-jilin
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
