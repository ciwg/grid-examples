package main

import (
	"fmt"
	"os"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/service"
)

func main() {
	host := service.NewBrowserHost()
	if err := host.ServeSession(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
