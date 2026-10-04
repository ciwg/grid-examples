// Command grid-tui provides Charm's terminal embodiment for Ex3 collaboration.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/tui"
)

func main() {
	relay := flag.String("relay", "", "grid relay base URL (opens the launch form when omitted)")
	documentID := flag.String("doc", "", "document ID to open (opens the launch form when omitted)")
	name := flag.String("name", "", "collaboration display name (opens the launch form when omitted)")
	color := flag.String("color", "", "collaboration color (#RRGGBB; opens the launch form when omitted)")
	accessToken := flag.String("access-token", "", "optional relay bootstrap token for remote collaboration")
	flag.Parse()
	// Intent: Keep the public command stable while the richer terminal
	// workspace remains a reusable local presentation package. Source: DI-mutoh.
	if err := tui.Run(tui.Config{Relay: *relay, DocumentID: *documentID, Name: *name, Color: *color, AccessToken: *accessToken}); err != nil {
		fmt.Fprintln(os.Stderr, "Grid TUI:", err)
	}
}
