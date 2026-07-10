package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/agent"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/intake"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: pg-order-intake FIXTURE.json")
		os.Exit(1)
	}
	cfg, err := agent.ConfigFromEnv("intake")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-intake: %v\n", err)
		os.Exit(1)
	}
	result, err := intake.Run(context.Background(), cfg, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-intake: %v\n", err)
		os.Exit(1)
	}
	resultBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-intake: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(resultBytes))
}
