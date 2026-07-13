package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

//go:embed helper.bundle.cjs
var helperBundle []byte

//go:embed automerge_wasm_bg.wasm
var helperWASM []byte

func main() {
	var relayURL = flag.String("relay", "http://127.0.0.1:7001", "grid relay base URL")
	flag.Parse()

	nodePath, err := exec.LookPath("node")
	if err != nil {
		fmt.Fprintf(os.Stderr, "grid-nvim-sidecar requires node in PATH: %v\n", err)
		os.Exit(1)
	}
	helperPath, cleanup, err := writeHelper()
	if err != nil {
		fmt.Fprintf(os.Stderr, "grid-nvim-sidecar helper setup failed: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	// Intent: Keep the locked Go command path while hiding the existing
	// Node-based Automerge engine behind it, so Neovim gets a real CRDT sidecar
	// now without promoting the helper implementation detail into the public
	// surface. Source: DI-sulod
	command := exec.Command(nodePath, helperPath, "--relay", *relayURL)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		fmt.Fprintf(os.Stderr, "grid-nvim-sidecar helper failed: %v\n", err)
		os.Exit(1)
	}
}

func writeHelper() (string, func(), error) {
	dir, err := os.MkdirTemp("", "grid-nvim-sidecar-*")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(dir, "helper.bundle.cjs")
	if err := os.WriteFile(path, helperBundle, 0o700); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, err
	}
	wasmPath := filepath.Join(dir, "automerge_wasm_bg.wasm")
	if err := os.WriteFile(wasmPath, helperWASM, 0o600); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, err
	}
	return path, func() {
		_ = os.RemoveAll(dir)
	}, nil
}
