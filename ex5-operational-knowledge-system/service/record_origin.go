package service

import records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"

func normalizeKnowledgeItemRecordOrigins(in []SignedKnowledgeItemRecord, events []OperationalEvent) []SignedKnowledgeItemRecord {
	recordSlice := make([]records.SignedKnowledgeItemRecord, len(in))
	eventSlice := make([]records.Event, len(events))
	for i, record := range in {
		recordSlice[i] = records.SignedKnowledgeItemRecord(record)
	}
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	out := records.NormalizeKnowledgeItemRecordOrigins(recordSlice, eventSlice)
	converted := make([]SignedKnowledgeItemRecord, len(out))
	for i, record := range out {
		converted[i] = SignedKnowledgeItemRecord(record)
	}
	return converted
}

func normalizeKnowledgeApprovalRecordOrigins(in []SignedKnowledgeApprovalRecord, events []OperationalEvent) []SignedKnowledgeApprovalRecord {
	recordSlice := make([]records.SignedKnowledgeApprovalRecord, len(in))
	eventSlice := make([]records.Event, len(events))
	for i, record := range in {
		recordSlice[i] = records.SignedKnowledgeApprovalRecord(record)
	}
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	out := records.NormalizeKnowledgeApprovalRecordOrigins(recordSlice, eventSlice)
	converted := make([]SignedKnowledgeApprovalRecord, len(out))
	for i, record := range out {
		converted[i] = SignedKnowledgeApprovalRecord(record)
	}
	return converted
}

func normalizeKnowledgeEvidenceRecordOrigins(in []SignedKnowledgeEvidenceRecord, events []OperationalEvent) []SignedKnowledgeEvidenceRecord {
	recordSlice := make([]records.SignedKnowledgeEvidenceRecord, len(in))
	eventSlice := make([]records.Event, len(events))
	for i, record := range in {
		recordSlice[i] = records.SignedKnowledgeEvidenceRecord(record)
	}
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	out := records.NormalizeKnowledgeEvidenceRecordOrigins(recordSlice, eventSlice)
	converted := make([]SignedKnowledgeEvidenceRecord, len(out))
	for i, record := range out {
		converted[i] = SignedKnowledgeEvidenceRecord(record)
	}
	return converted
}

func normalizeOperationalRunRecordOrigins(in []SignedOperationalRunRecord, events []OperationalEvent) []SignedOperationalRunRecord {
	recordSlice := make([]records.SignedOperationalRunRecord, len(in))
	eventSlice := make([]records.Event, len(events))
	for i, record := range in {
		recordSlice[i] = records.SignedOperationalRunRecord(record)
	}
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	out := records.NormalizeOperationalRunRecordOrigins(recordSlice, eventSlice)
	converted := make([]SignedOperationalRunRecord, len(out))
	for i, record := range out {
		converted[i] = SignedOperationalRunRecord(record)
	}
	return converted
}

func normalizeOperationalPlaceRecordOrigins(in []SignedOperationalPlaceRecord, events []OperationalEvent) []SignedOperationalPlaceRecord {
	recordSlice := make([]records.SignedOperationalPlaceRecord, len(in))
	eventSlice := make([]records.Event, len(events))
	for i, record := range in {
		recordSlice[i] = records.SignedOperationalPlaceRecord(record)
	}
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	out := records.NormalizeOperationalPlaceRecordOrigins(recordSlice, eventSlice)
	converted := make([]SignedOperationalPlaceRecord, len(out))
	for i, record := range out {
		converted[i] = SignedOperationalPlaceRecord(record)
	}
	return converted
}

func normalizeOperationalResourceRecordOrigins(in []SignedOperationalResourceRecord, events []OperationalEvent) []SignedOperationalResourceRecord {
	recordSlice := make([]records.SignedOperationalResourceRecord, len(in))
	eventSlice := make([]records.Event, len(events))
	for i, record := range in {
		recordSlice[i] = records.SignedOperationalResourceRecord(record)
	}
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	out := records.NormalizeOperationalResourceRecordOrigins(recordSlice, eventSlice)
	converted := make([]SignedOperationalResourceRecord, len(out))
	for i, record := range out {
		converted[i] = SignedOperationalResourceRecord(record)
	}
	return converted
}

func normalizeKnowledgeLinkRecordOrigins(in []SignedKnowledgeLinkRecord, events []OperationalEvent) []SignedKnowledgeLinkRecord {
	recordSlice := make([]records.SignedKnowledgeLinkRecord, len(in))
	eventSlice := make([]records.Event, len(events))
	for i, record := range in {
		recordSlice[i] = records.SignedKnowledgeLinkRecord(record)
	}
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	out := records.NormalizeKnowledgeLinkRecordOrigins(recordSlice, eventSlice)
	converted := make([]SignedKnowledgeLinkRecord, len(out))
	for i, record := range out {
		converted[i] = SignedKnowledgeLinkRecord(record)
	}
	return converted
}

func normalizeKnowledgeResponsibilityRecordOrigins(in []SignedKnowledgeResponsibilityRecord, events []OperationalEvent) []SignedKnowledgeResponsibilityRecord {
	recordSlice := make([]records.SignedKnowledgeResponsibilityRecord, len(in))
	eventSlice := make([]records.Event, len(events))
	for i, record := range in {
		recordSlice[i] = records.SignedKnowledgeResponsibilityRecord(record)
	}
	for i, event := range events {
		eventSlice[i] = records.Event(event)
	}
	out := records.NormalizeKnowledgeResponsibilityRecordOrigins(recordSlice, eventSlice)
	converted := make([]SignedKnowledgeResponsibilityRecord, len(out))
	for i, record := range out {
		converted[i] = SignedKnowledgeResponsibilityRecord(record)
	}
	return converted
}
