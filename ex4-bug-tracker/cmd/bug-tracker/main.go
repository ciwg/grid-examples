package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/service"
)

func main() {
	var (
		// Intent: Keep the bug tracker self-contained and easy to run locally
		// with its own runtime root and default port instead of reusing any grid
		// editor launcher assumptions. Source: DI-dajak
		listen   = flag.String("listen", "127.0.0.1:7035", "listen address")
		dataRoot = flag.String("data-root", ".bug-tracker", "local runtime data root")
		seedDemo = flag.Bool("seed-demo", false, "seed a starter bug tracker dataset when the runtime root is empty")
	)
	flag.Parse()

	app, err := service.NewApp(*dataRoot)
	if err != nil {
		log.Fatalf("new app: %v", err)
	}
	if *seedDemo {
		seeded, err := app.SeedDemoIfEmpty()
		if err != nil {
			log.Fatalf("seed demo: %v", err)
		}
		log.Printf("demo seed applied=%t", seeded)
	}
	server := service.NewServer(app)
	log.Printf("bug-tracker listening on %s", *listen)
	log.Printf("runtime root=%s", app.Meta().DataRoot)
	log.Fatal(http.ListenAndServe(*listen, server.Handler()))
}
