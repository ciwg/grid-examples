// Command grid-tui provides Charm's terminal embodiment for Ex3 collaboration.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/tui"
)

func main() {
	relay := flag.String("relay", "http://127.0.0.1:7025", "grid relay base URL")
	documentID := flag.String("doc", "demo", "document ID to open")
	name := flag.String("name", "Charm User", "collaboration display name")
	color := flag.String("color", "#8b5cf6", "collaboration color (#RRGGBB)")
	flag.Parse()
	// Intent: Keep the public command stable while the richer terminal
	// workspace remains a reusable local presentation package. Source: DI-mutoh.
	if err := tui.Run(tui.Config{Relay: *relay, DocumentID: *documentID, Name: *name, Color: *color}); err != nil {
		fmt.Fprintln(os.Stderr, "Grid TUI:", err)
	}
}
