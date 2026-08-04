package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed index.html style.css app.js
var files embed.FS

func Handler() http.Handler {
	root, err := fs.Sub(files, ".")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(root))
}
