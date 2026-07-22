package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	socketServer := service.NewLocalEmbodimentServer(app, service.EmbodimentSocketPath(*dataRoot))
	go func() {
		// Intent: Start the direct Unix-socket embodiment contract alongside the
		// browser HTTP adapter so terminal embodiments can leave HTTP without
		// splitting runtime ownership or durable state. Source: DI-favel
		if err := socketServer.ListenAndServe(); err != nil {
			log.Fatalf("local embodiment socket: %v", err)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		if err := socketServer.Close(); err != nil {
			log.Printf("close local embodiment socket: %v", err)
		}
		os.Exit(0)
	}()
	log.Printf("operational-knowledge listening on http://%s", *listen)
	if err := http.ListenAndServe(*listen, service.NewServer(app).Handler()); err != nil {
		log.Fatal(err)
	}
}
