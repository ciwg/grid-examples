package service

import "strings"

func effectiveOriginPeerID(event OperationalEvent, localPeerID string) string {
	if strings.TrimSpace(event.OriginPeerID) != "" {
		return event.OriginPeerID
	}
	return localPeerID
}

func effectiveOriginSequence(event OperationalEvent) uint64 {
	if event.OriginSequence != 0 {
		return event.OriginSequence
	}
	return event.Sequence
}
