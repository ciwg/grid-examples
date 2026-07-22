package protocols

import (
	"fmt"

	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

const (
	gridTag = uint64(0x67726964)
	cidTag  = uint64(42)
)

// Intent: Compute the first frozen ex5 profile and envelope identities from
// raw bytes so runtime behavior is selected by exact shipped spec content.
// Source: DI-mibor
func CIDForBytes(data []byte) (cid.Cid, error) {
	sum, err := mh.Sum(data, mh.SHA2_256, -1)
	if err != nil {
		return cid.Undef, fmt.Errorf("multihash bytes: %w", err)
	}
	return cid.NewCidV1(cid.Raw, sum), nil
}
