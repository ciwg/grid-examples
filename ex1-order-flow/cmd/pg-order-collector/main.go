package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/collector"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	address := os.Getenv("PG_COLLECTOR_LISTEN_ADDR")
	if address == "" {
		address = ":9200"
	}
	dataDir := os.Getenv("PG_DATA_DIR")
	if dataDir == "" {
		fmt.Fprintln(os.Stderr, "pg-order-collector: PG_DATA_DIR is required")
		os.Exit(1)
	}
	service := &collector.Service{
		Address: address,
		DataDir: dataDir,
	}
	if err := service.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-collector: %v\n", err)
		os.Exit(1)
	}
}
