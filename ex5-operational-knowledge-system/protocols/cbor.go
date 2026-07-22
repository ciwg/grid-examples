package protocols

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

var (
	// Intent: Freeze the first ex5 PromiseGrid-native family on canonical CBOR
	// so the same knowledge-item payload bytes drive pCID selection, signing,
	// and replay verification. Source: DI-mibor
	canonicalEncMode = mustCanonicalEncMode()
	strictDecMode    = mustStrictDecMode()
)

func mustCanonicalEncMode() cbor.EncMode {
	mode, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		panic(fmt.Errorf("canonical enc mode: %w", err))
	}
	return mode
}

func mustStrictDecMode() cbor.DecMode {
	options := cbor.DecOptions{
		DupMapKey: cbor.DupMapKeyEnforcedAPF,
	}
	mode, err := options.DecMode()
	if err != nil {
		panic(fmt.Errorf("strict dec mode: %w", err))
	}
	return mode
}

func Marshal(value any) ([]byte, error) {
	return canonicalEncMode.Marshal(value)
}

func MustMarshal(value any) []byte {
	body, err := Marshal(value)
	if err != nil {
		panic(err)
	}
	return body
}

func Unmarshal(data []byte, value any) error {
	return strictDecMode.Unmarshal(data, value)
}
