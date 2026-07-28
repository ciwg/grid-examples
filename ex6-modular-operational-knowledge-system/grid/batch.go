package grid

import "encoding/json"

const RelayBatchFormat = "moks-relay-batch-v1"

type Batch struct {
	Format  string            `json:"format"`
	Records []json.RawMessage `json:"records"`
}
