package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/kernel"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	address := os.Getenv("PG_KERNEL_LISTEN_ADDR")
	if address == "" {
		address = ":7000"
	}
	dataDir := os.Getenv("PG_DATA_DIR")
	if dataDir == "" {
		fmt.Fprintln(os.Stderr, "pg-order-kernel: PG_DATA_DIR is required")
		os.Exit(1)
	}
	server := &kernel.Server{Address: address, DataDir: dataDir}
	if err := server.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-kernel: %v\n", err)
		os.Exit(1)
	}
}
