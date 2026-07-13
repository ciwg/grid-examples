package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var relayURL = flag.String("relay", "http://127.0.0.1:7001", "grid relay base URL")
	flag.Parse()
	// Intent: Make the missing Go-side CRDT replica explicit in this slice so
	// the repo exposes the locked command path without falsely claiming that the
	// Neovim sidecar is already feature-complete. Source: DI-lumek
	_, _ = fmt.Fprintf(os.Stderr, "grid-nvim-sidecar is a transitional scaffold in this slice; use the relay at %s with the current Neovim compatibility client while the real Go CRDT sidecar lands.\n", *relayURL)
	os.Exit(1)
}
