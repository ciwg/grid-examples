package transport

import (
	"sort"
	"strconv"
	"strings"

	records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"
)

// Intent: Keep origin-aware relay and peer-exchange filtering in the reusable
// transport substrate so ex5 app code stops owning identical wire-selection
// mechanics. Source: DI-vurem
func OriginEventKey(peerID string, originSequence uint64) string {
	return strings.TrimSpace(peerID) + "#" + strconv.FormatUint(originSequence, 10)
}

// Intent: Reuse one record-origin key shape across relay and peer-exchange
// filtering so the transport substrate selects signed records by durable
// origin identity instead of local sequence. Source: DI-vurem
func RecordOriginKey(peerID string, originSequence uint64, sequence uint64) string {
	return records.RecordOriginKey(peerID, originSequence, sequence)
}

// Intent: Emit only relay-visible unseen events above the caller's origin
// cursor map so relay carriage stays incremental and origin-aware. Source:
// DI-vurem
func FilterRelayFeedEvents(events []records.Event, known map[string]uint64) ([]records.Event, map[string]bool) {
	out := make([]records.Event, 0, len(events))
	wanted := map[string]bool{}
	for _, event := range events {
		seen := known[event.OriginPeerID]
		if event.OriginSequence <= seen {
			continue
		}
		out = append(out, event)
		wanted[OriginEventKey(event.OriginPeerID, event.OriginSequence)] = true
	}
	return out, wanted
}

// Intent: Keep evidence-blob dependency selection inside the reusable
// transport substrate so both local peer exchange and remote relay batch logic
// name the same required CID set. Source: DI-vurem
func RequiredBlobCIDsForEvents(events []records.Event) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, event := range events {
		if event.Type != "evidence_added" || strings.TrimSpace(event.AttachmentCID) == "" {
			continue
		}
		cid := strings.TrimSpace(event.AttachmentCID)
		if seen[cid] {
			continue
		}
		seen[cid] = true
		out = append(out, cid)
	}
	sort.Strings(out)
	return out
}

func FilterKnowledgeItemRecordsByOrigin(in []records.SignedKnowledgeItemRecord, wanted map[string]bool) []records.SignedKnowledgeItemRecord {
	out := []records.SignedKnowledgeItemRecord{}
	for _, record := range in {
		if wanted[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func FilterKnowledgeApprovalRecordsByOrigin(in []records.SignedKnowledgeApprovalRecord, wanted map[string]bool) []records.SignedKnowledgeApprovalRecord {
	out := []records.SignedKnowledgeApprovalRecord{}
	for _, record := range in {
		if wanted[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func FilterKnowledgeEvidenceRecordsByOrigin(in []records.SignedKnowledgeEvidenceRecord, wanted map[string]bool) []records.SignedKnowledgeEvidenceRecord {
	out := []records.SignedKnowledgeEvidenceRecord{}
	for _, record := range in {
		if wanted[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func FilterOperationalRunRecordsByOrigin(in []records.SignedOperationalRunRecord, wanted map[string]bool) []records.SignedOperationalRunRecord {
	out := []records.SignedOperationalRunRecord{}
	for _, record := range in {
		if wanted[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func FilterOperationalPlaceRecordsByOrigin(in []records.SignedOperationalPlaceRecord, wanted map[string]bool) []records.SignedOperationalPlaceRecord {
	out := []records.SignedOperationalPlaceRecord{}
	for _, record := range in {
		if wanted[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func FilterOperationalResourceRecordsByOrigin(in []records.SignedOperationalResourceRecord, wanted map[string]bool) []records.SignedOperationalResourceRecord {
	out := []records.SignedOperationalResourceRecord{}
	for _, record := range in {
		if wanted[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func FilterKnowledgeLinkRecordsByOrigin(in []records.SignedKnowledgeLinkRecord, wanted map[string]bool) []records.SignedKnowledgeLinkRecord {
	out := []records.SignedKnowledgeLinkRecord{}
	for _, record := range in {
		if wanted[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}

func FilterKnowledgeResponsibilityRecordsByOrigin(in []records.SignedKnowledgeResponsibilityRecord, wanted map[string]bool) []records.SignedKnowledgeResponsibilityRecord {
	out := []records.SignedKnowledgeResponsibilityRecord{}
	for _, record := range in {
		if wanted[RecordOriginKey(record.OriginPeerID, record.OriginSequence, record.Sequence)] {
			out = append(out, record)
		}
	}
	return out
}
