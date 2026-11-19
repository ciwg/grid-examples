package web

import "embed"

//go:embed index.html app.js style.css
var files embed.FS

func MustRead(name string) []byte {
	bytes, err := files.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return bytes
}
