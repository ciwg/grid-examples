package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/accounting"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/agent"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cfg, err := agent.ConfigFromEnv("accounting")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-accounting: %v\n", err)
		os.Exit(1)
	}
	if err := accounting.Run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-accounting: %v\n", err)
		os.Exit(1)
	}
}
