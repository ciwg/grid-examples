package service

import records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"

func effectiveOriginPeerID(event OperationalEvent, localPeerID string) string {
	return records.EffectiveOriginPeerID(records.Event(event), localPeerID)
}

func effectiveOriginSequence(event OperationalEvent) uint64 {
	return records.EffectiveOriginSequence(records.Event(event))
}
