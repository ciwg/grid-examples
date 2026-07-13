package protocols

import "embed"

const (
	LiveDocumentSpec  = "live-document.md"
	LiveAwarenessSpec = "live-awareness.md"
)

// Intent: pCIDs must derive from the exact local draft spec bytes so the app
// can keep protocol identity content-addressed even before upstream specs
// exist. Source: DI-tofug

//go:embed *.md
var files embed.FS

func MustRead(name string) []byte {
	bytes, err := files.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return bytes
}
