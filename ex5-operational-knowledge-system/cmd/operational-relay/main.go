package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/service"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:7046", "listen address")
	dataRoot := flag.String("data-root", ".operational-relay", "relay data root")
	flag.Parse()

	// Intent: Run the remote relay as its own deployable durable service instead
	// of as another mode of the local embodiment/runtime server. Source:
	// DI-rovik
	relay, err := service.NewRelay(*dataRoot)
	if err != nil {
		log.Fatalf("new relay: %v", err)
	}
	defer func() {
		if err := relay.Close(); err != nil {
			log.Printf("close relay: %v", err)
		}
	}()

	log.Printf("operational-relay listening on http://%s%s", *listen, "/relay/v1")
	if err := http.ListenAndServe(*listen, service.NewRelayServer(relay).Handler()); err != nil {
		log.Fatal(err)
	}
}
