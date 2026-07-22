package service

func normalizeKnowledgeItemRecordOrigins(records []SignedKnowledgeItemRecord, events []OperationalEvent) []SignedKnowledgeItemRecord {
	out := make([]SignedKnowledgeItemRecord, 0, len(records))
	eventsBySequence := map[uint64]OperationalEvent{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func normalizeKnowledgeApprovalRecordOrigins(records []SignedKnowledgeApprovalRecord, events []OperationalEvent) []SignedKnowledgeApprovalRecord {
	out := make([]SignedKnowledgeApprovalRecord, 0, len(records))
	eventsBySequence := map[uint64]OperationalEvent{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func normalizeKnowledgeEvidenceRecordOrigins(records []SignedKnowledgeEvidenceRecord, events []OperationalEvent) []SignedKnowledgeEvidenceRecord {
	out := make([]SignedKnowledgeEvidenceRecord, 0, len(records))
	eventsBySequence := map[uint64]OperationalEvent{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func normalizeOperationalRunRecordOrigins(records []SignedOperationalRunRecord, events []OperationalEvent) []SignedOperationalRunRecord {
	out := make([]SignedOperationalRunRecord, 0, len(records))
	eventsBySequence := map[uint64]OperationalEvent{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func normalizeOperationalPlaceRecordOrigins(records []SignedOperationalPlaceRecord, events []OperationalEvent) []SignedOperationalPlaceRecord {
	out := make([]SignedOperationalPlaceRecord, 0, len(records))
	eventsBySequence := map[uint64]OperationalEvent{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func normalizeOperationalResourceRecordOrigins(records []SignedOperationalResourceRecord, events []OperationalEvent) []SignedOperationalResourceRecord {
	out := make([]SignedOperationalResourceRecord, 0, len(records))
	eventsBySequence := map[uint64]OperationalEvent{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func normalizeKnowledgeLinkRecordOrigins(records []SignedKnowledgeLinkRecord, events []OperationalEvent) []SignedKnowledgeLinkRecord {
	out := make([]SignedKnowledgeLinkRecord, 0, len(records))
	eventsBySequence := map[uint64]OperationalEvent{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}

func normalizeKnowledgeResponsibilityRecordOrigins(records []SignedKnowledgeResponsibilityRecord, events []OperationalEvent) []SignedKnowledgeResponsibilityRecord {
	out := make([]SignedKnowledgeResponsibilityRecord, 0, len(records))
	eventsBySequence := map[uint64]OperationalEvent{}
	for _, event := range events {
		eventsBySequence[event.Sequence] = event
	}
	for _, record := range records {
		if event, ok := eventsBySequence[record.Sequence]; ok {
			if record.OriginPeerID == "" {
				record.OriginPeerID = event.OriginPeerID
			}
			if record.OriginSequence == 0 {
				record.OriginSequence = event.OriginSequence
			}
		}
		out = append(out, record)
	}
	return out
}
