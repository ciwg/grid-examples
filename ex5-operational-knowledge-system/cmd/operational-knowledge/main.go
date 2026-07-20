package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/service"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:7045", "listen address")
	dataRoot := flag.String("data-root", ".operational-knowledge-system", "runtime data root")
	flag.Parse()

	// Intent: Start one local ex5 runtime that owns durable history and serves
	// equal browser and CLI embodiments over the same state root. Source:
	// DI-radok; DI-zuvob
	app, err := service.NewApp(*dataRoot)
	if err != nil {
		log.Fatalf("new app: %v", err)
	}
	log.Printf("operational-knowledge listening on http://%s", *listen)
	if err := http.ListenAndServe(*listen, service.NewServer(app).Handler()); err != nil {
		log.Fatal(err)
	}
}
