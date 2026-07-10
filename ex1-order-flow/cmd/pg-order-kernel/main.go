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
	server := &kernel.Server{Address: address}
	if err := server.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-kernel: %v\n", err)
		os.Exit(1)
	}
}
