package service

import records "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/records"

type RuntimeIdentity = records.RuntimeIdentity

// Intent: Keep service-level runtime startup on the extracted PromiseGrid
// record identity without changing the existing store and app ownership model.
// Source: DI-ragiv
func LoadOrCreateRuntimeIdentity(path string) (*RuntimeIdentity, error) {
	return records.LoadOrCreateRuntimeIdentity(path)
}

func VerifyRuntimeProof(signable []byte, proofBytes []byte) error {
	return records.VerifyRuntimeProof(signable, proofBytes)
}
