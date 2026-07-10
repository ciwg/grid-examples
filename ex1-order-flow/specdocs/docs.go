package specdocs

import "embed"

// Intent: pCIDs must derive from exact local spec bytes so the implementation
// can keep protocol identity content-addressed without relying on hard-coded
// names alone. Source: DI-movab; DI-lihit

//go:embed *.md
var files embed.FS

func MustRead(name string) []byte {
	bytes, err := files.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return bytes
}
