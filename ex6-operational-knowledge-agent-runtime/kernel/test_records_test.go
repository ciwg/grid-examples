package kernel

import (
	"encoding/json"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
)

// canonicalTestRecord builds exact canonical Grid bytes for tests that exercise
// durable package evidence. External fixtures pass their frozen pCID explicitly
// rather than deriving a protocol identity from a family label. Source: DI-sidoh, DI-solan.
func canonicalTestRecord(t *testing.T, protocolPCID, family, recordID, signer, timestamp string, payload any) []byte {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := records.CanonicalJSON(body)
	if err != nil {
		t.Fatal(err)
	}
	envelope := records.Envelope{Family: family, ProtocolPCID: protocolPCID, RecordID: recordID, Signer: signer, Timestamp: timestamp, Payload: canonical}
	if err := envelope.Validate(); err != nil {
		t.Fatal(err)
	}
	return records.MustMarshal(envelope)
}
