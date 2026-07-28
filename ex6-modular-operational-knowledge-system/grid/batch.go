package grid

import "encoding/json"

const RelayBatchFormat = "moks-relay-batch-v1"

type ImplementationClaim struct {
	PackageID      string `json:"package_id"`
	ProtocolPCID   string `json:"protocol_pcid"`
	Role           string `json:"role"`
	Summary        string `json:"summary,omitempty"`
	PackageVersion string `json:"package_version"`
}

type Batch struct {
	Format               string                `json:"format"`
	Implementation       string                `json:"implementation"`
	ExportedAt           string                `json:"exported_at"`
	ImplementationClaims []ImplementationClaim `json:"implementation_claims,omitempty"`
	Records              []json.RawMessage     `json:"records"`
}
