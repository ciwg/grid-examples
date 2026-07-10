package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/analyzer"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: pg-order-analyze DATA_ROOT")
		os.Exit(1)
	}
	summary, err := analyzer.Analyze(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-analyze: %v\n", err)
		os.Exit(1)
	}
	summaryBytes, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-order-analyze: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(summaryBytes))
}
