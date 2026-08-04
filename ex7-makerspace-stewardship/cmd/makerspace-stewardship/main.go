package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/computerscienceiscool/grid-examples/ex7-makerspace-stewardship/service"
)

func main() {
	runtimeRoot := flag.String("runtime-root", ".makerspace-stewardship", "directory for append-only local evidence")
	flag.Parse()
	app, err := service.NewPersistentDemoApp(*runtimeRoot)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Makerspace Stewardship listening on http://127.0.0.1:7037")
	log.Fatal(http.ListenAndServe("127.0.0.1:7037", app.Handler()))
}
