package protocol

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

var (
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
	bytes, err := Marshal(value)
	if err != nil {
		panic(err)
	}
	return bytes
}

func Unmarshal(data []byte, value any) error {
	return strictDecMode.Unmarshal(data, value)
}
