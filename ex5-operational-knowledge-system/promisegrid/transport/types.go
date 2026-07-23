package transport

import records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"

const (
	PeerExchangeBundleFormat = "ex5-peer-exchange-v5"
	RelayFeedFormat          = "ex5-relay-feed-v1"
)

type RelayFeedRequest struct {
	KnownOrigins map[string]uint64 `json:"known_origins,omitempty"`
}

type RelayFeedBatch struct {
	Format                         string                                        `json:"format"`
	ExportedAt                     string                                        `json:"exported_at"`
	Implementation                 string                                        `json:"implementation"`
	ExportingPeerID                string                                        `json:"exporting_peer_id"`
	KnowledgeItemPCID              string                                        `json:"knowledge_item_pcid"`
	KnowledgeApprovalPCID          string                                        `json:"knowledge_approval_pcid"`
	KnowledgeEvidencePCID          string                                        `json:"knowledge_evidence_pcid"`
	KnowledgeLinkPCID              string                                        `json:"knowledge_link_pcid"`
	KnowledgeResponsibilityPCID    string                                        `json:"knowledge_responsibility_pcid"`
	OperationalRunPCID             string                                        `json:"operational_run_pcid"`
	OperationalPlacePCID           string                                        `json:"operational_place_pcid"`
	OperationalResourcePCID        string                                        `json:"operational_resource_pcid"`
	Events                         []records.Event                               `json:"events"`
	KnowledgeItemRecords           []records.SignedKnowledgeItemRecord           `json:"knowledge_item_records"`
	KnowledgeApprovalRecords       []records.SignedKnowledgeApprovalRecord       `json:"knowledge_approval_records"`
	KnowledgeEvidenceRecords       []records.SignedKnowledgeEvidenceRecord       `json:"knowledge_evidence_records"`
	OperationalRunRecords          []records.SignedOperationalRunRecord          `json:"operational_run_records"`
	OperationalPlaceRecords        []records.SignedOperationalPlaceRecord        `json:"operational_place_records"`
	OperationalResourceRecords     []records.SignedOperationalResourceRecord     `json:"operational_resource_records"`
	KnowledgeLinkRecords           []records.SignedKnowledgeLinkRecord           `json:"knowledge_link_records"`
	KnowledgeResponsibilityRecords []records.SignedKnowledgeResponsibilityRecord `json:"knowledge_responsibility_records"`
	RequiredBlobCIDs               []string                                      `json:"required_blob_cids,omitempty"`
}

type PeerExchangeImportIssue struct {
	RecordType string `json:"record_type"`
	RecordID   string `json:"record_id"`
	Reason     string `json:"reason"`
}

type RelayFeedImportResult struct {
	ImportedEvents               int                       `json:"imported_events"`
	ImportedKnowledgeItems       int                       `json:"imported_knowledge_items"`
	ImportedKnowledgeApprovals   int                       `json:"imported_knowledge_approvals"`
	ImportedKnowledgeEvidence    int                       `json:"imported_knowledge_evidence"`
	ImportedOperationalRuns      int                       `json:"imported_operational_runs"`
	ImportedOperationalPlaces    int                       `json:"imported_operational_places"`
	ImportedOperationalResources int                       `json:"imported_operational_resources"`
	ImportedKnowledgeLinks       int                       `json:"imported_knowledge_links"`
	ImportedResponsibilities     int                       `json:"imported_responsibilities"`
	MissingBlobCIDs              []string                  `json:"missing_blob_cids,omitempty"`
	UnresolvedReferences         []PeerExchangeImportIssue `json:"unresolved_references"`
}

type RelayPublishResult struct {
	PublishedEvents               int      `json:"published_events"`
	PublishedKnowledgeItems       int      `json:"published_knowledge_items"`
	PublishedKnowledgeApprovals   int      `json:"published_knowledge_approvals"`
	PublishedKnowledgeEvidence    int      `json:"published_knowledge_evidence"`
	PublishedOperationalRuns      int      `json:"published_operational_runs"`
	PublishedOperationalPlaces    int      `json:"published_operational_places"`
	PublishedOperationalResources int      `json:"published_operational_resources"`
	PublishedKnowledgeLinks       int      `json:"published_knowledge_links"`
	PublishedResponsibilities     int      `json:"published_responsibilities"`
	MissingBlobCIDs               []string `json:"missing_blob_cids,omitempty"`
}

type PeerExchangeBundle struct {
	Format                         string                                        `json:"format"`
	ExportedAt                     string                                        `json:"exported_at"`
	Implementation                 string                                        `json:"implementation"`
	ExportingPeerID                string                                        `json:"exporting_peer_id"`
	KnowledgeItemPCID              string                                        `json:"knowledge_item_pcid"`
	KnowledgeApprovalPCID          string                                        `json:"knowledge_approval_pcid"`
	KnowledgeEvidencePCID          string                                        `json:"knowledge_evidence_pcid"`
	KnowledgeLinkPCID              string                                        `json:"knowledge_link_pcid"`
	KnowledgeResponsibilityPCID    string                                        `json:"knowledge_responsibility_pcid"`
	OperationalRunPCID             string                                        `json:"operational_run_pcid"`
	OperationalPlacePCID           string                                        `json:"operational_place_pcid"`
	OperationalResourcePCID        string                                        `json:"operational_resource_pcid"`
	Events                         []records.Event                               `json:"events"`
	KnowledgeItemRecords           []records.SignedKnowledgeItemRecord           `json:"knowledge_item_records"`
	KnowledgeApprovalRecords       []records.SignedKnowledgeApprovalRecord       `json:"knowledge_approval_records"`
	KnowledgeEvidenceRecords       []records.SignedKnowledgeEvidenceRecord       `json:"knowledge_evidence_records"`
	OperationalRunRecords          []records.SignedOperationalRunRecord          `json:"operational_run_records"`
	OperationalPlaceRecords        []records.SignedOperationalPlaceRecord        `json:"operational_place_records"`
	OperationalResourceRecords     []records.SignedOperationalResourceRecord     `json:"operational_resource_records"`
	KnowledgeLinkRecords           []records.SignedKnowledgeLinkRecord           `json:"knowledge_link_records"`
	KnowledgeResponsibilityRecords []records.SignedKnowledgeResponsibilityRecord `json:"knowledge_responsibility_records"`
	CASBlobObjects                 map[string]string                             `json:"cas_blob_objects,omitempty"`
}

type PeerExchangeImportResult struct {
	ImportedEvents               int                       `json:"imported_events"`
	ImportedKnowledgeItems       int                       `json:"imported_knowledge_items"`
	ImportedKnowledgeApprovals   int                       `json:"imported_knowledge_approvals"`
	ImportedKnowledgeEvidence    int                       `json:"imported_knowledge_evidence"`
	ImportedOperationalRuns      int                       `json:"imported_operational_runs"`
	ImportedOperationalPlaces    int                       `json:"imported_operational_places"`
	ImportedOperationalResources int                       `json:"imported_operational_resources"`
	ImportedKnowledgeLinks       int                       `json:"imported_knowledge_links"`
	ImportedResponsibilities     int                       `json:"imported_responsibilities"`
	UnresolvedReferences         []PeerExchangeImportIssue `json:"unresolved_references"`
}
