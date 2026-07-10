package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/agent"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/warehouse"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cfg, err := agent.ConfigFromEnv("warehouse")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-warehouse: %v\n", err)
		os.Exit(1)
	}
	if err := warehouse.Run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-warehouse: %v\n", err)
		os.Exit(1)
	}
}
