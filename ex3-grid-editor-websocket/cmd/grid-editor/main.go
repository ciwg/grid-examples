package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/service"
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
		// Intent: Keep the copied example's one-process launcher on ex3's own
		// default loopback port so it does not collide with ex2. Source: DI-vatub
		listen   = flag.String("listen", "127.0.0.1:7025", "listen address")
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
