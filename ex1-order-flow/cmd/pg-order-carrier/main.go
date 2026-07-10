package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/agent"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/carrier"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cfg, err := agent.ConfigFromEnv("carrier")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-carrier: %v\n", err)
		os.Exit(1)
	}
	if err := carrier.Run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-carrier: %v\n", err)
		os.Exit(1)
	}
}
