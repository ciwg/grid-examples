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
		// Intent: Move the copied example off ex2's default loopback ports so the
		// two examples can run side by side on one machine. Source: DI-vatub
		listen            = flag.String("listen", "127.0.0.1:7025", "listen address")
		dataRoot          = flag.String("data-root", ".grid-editor", "local runtime data root")
		remoteAccessToken = flag.String("remote-access-token", "", "optional bootstrap token for remote mutation sessions")
		peers             peerFlags
	)
	flag.Var(&peers, "peer", "peer base URL to poll for signed messages (repeatable)")
	flag.Parse()

	app, err := service.NewApp(*dataRoot, service.AppOptions{RemoteAccessToken: *remoteAccessToken})
	if err != nil {
		log.Fatalf("new relay: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.StartPeerPolling(ctx, peers, 2*time.Second)

	server := service.NewServer(app)
	log.Printf("grid-relay listening on %s", *listen)
	log.Printf("relay author=%s", app.Meta().LocalID)
	log.Printf("live-document pCID=%s", app.Meta().DocumentPCID)
	log.Printf("live-awareness pCID=%s", app.Meta().AwarenessPCID)
	log.Printf("remote mutation bootstrap enabled=%t", app.RemoteAccessEnabled())
	log.Fatal(http.ListenAndServe(*listen, server.Handler()))
}
