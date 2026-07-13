package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/service"
)

type peerFlags []string

func (flags *peerFlags) String() string {
	return ""
}

func (flags *peerFlags) Set(value string) error {
	*flags = append(*flags, value)
	return nil
}

func main() {
	var (
		listen   = flag.String("listen", "127.0.0.1:7015", "listen address")
		dataRoot = flag.String("data-root", ".grid-editor", "local runtime data root")
		peers    peerFlags
	)
	flag.Var(&peers, "peer", "peer base URL to poll for signed messages (repeatable)")
	flag.Parse()

	app, err := service.NewApp(*dataRoot)
	if err != nil {
		log.Fatalf("new app: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.StartPeerPolling(ctx, peers, 2*time.Second)

	server := service.NewServer(app)
	log.Printf("grid-editor listening on %s", *listen)
	log.Printf("local author=%s", app.Meta().LocalID)
	log.Printf("live-document pCID=%s", app.Meta().DocumentPCID)
	log.Printf("live-awareness pCID=%s", app.Meta().AwarenessPCID)
	log.Fatal(http.ListenAndServe(*listen, server.Handler()))
}
