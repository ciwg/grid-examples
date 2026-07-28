package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/builtin"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
	contextpkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/context"
	inventorypkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/inventory"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/knowledge"
	linkspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/links"
	maintenancepkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/maintenance"
	procedurespkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/procedures"
	receivingpkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/receiving"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/runs"
	trainingpkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/training"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	root, err := defaultRuntimeRoot()
	if err != nil {
		return err
	}
	runtime, err := kernel.Open(root)
	if err != nil {
		return err
	}
	defer func() {
		_ = runtime.Close()
	}()
	if err := registerBuiltins(runtime); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("command is required")
	}
	switch {
	case matchesPrefix(args, "package", "list"):
		return packageList(runtime)
	case matchesPrefix(args, "package", "inspect"):
		if len(args) != 3 {
			return errors.New("usage: package inspect <package-id>")
		}
		return packageInspect(runtime, args[2])
	case matchesPrefix(args, "package", "install"):
		if len(args) != 3 {
			return errors.New("usage: package install <dir>")
		}
		manifest, err := runtime.InstallPackageDir(ctx, args[2])
		if err != nil {
			return err
		}
		fmt.Printf("installed %s\n", manifest.ID)
		return nil
	case matchesPrefix(args, "relay", "export"):
		if len(args) != 3 {
			return errors.New("usage: relay export <path>")
		}
		return relayExport(runtime, args[2])
	case matchesPrefix(args, "relay", "import"):
		if len(args) != 3 {
			return errors.New("usage: relay import <path>")
		}
		return relayImport(ctx, runtime, args[2])
	case matchesPrefix(args, "relay", "serve"):
		if len(args) != 3 {
			return errors.New("usage: relay serve <addr>")
		}
		return relayServe(ctx, runtime, args[2])
	case matchesPrefix(args, "relay", "pull"):
		if len(args) != 3 {
			return errors.New("usage: relay pull <peer-id>")
		}
		return relayPull(ctx, runtime, args[2])
	case matchesPrefix(args, "relay", "push"):
		if len(args) != 3 {
			return errors.New("usage: relay push <peer-id>")
		}
		return relayPush(ctx, runtime, args[2])
	case matchesPrefix(args, "relay", "peer", "local"):
		if len(args) != 4 {
			return errors.New("usage: relay peer local show")
		}
		if args[3] != "show" {
			return errors.New("usage: relay peer local show")
		}
		fmt.Println(runtime.LocalPeerID())
		return nil
	case matchesPrefix(args, "relay", "peer", "list"):
		return relayPeerList(runtime)
	case matchesPrefix(args, "relay", "peer", "allow"):
		if len(args) != 7 {
			return errors.New("usage: relay peer allow <peer-id> <batch-url> <import-url> <pull|no-pull> <push|no-push>")
		}
		return relayPeerAllow(runtime, args[3:])
	case matchesPrefix(args, "relay", "peer", "revoke"):
		if len(args) != 5 {
			return errors.New("usage: relay peer revoke <peer-id>")
		}
		return runtime.RevokePeer(args[4])
	default:
		output, err := runtime.RunCommand(ctx, args)
		if err != nil {
			return err
		}
		if strings.TrimSpace(output) != "" {
			fmt.Println(output)
		}
		return nil
	}
}

func defaultRuntimeRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".moks"), nil
}

func matchesPrefix(args []string, prefix ...string) bool {
	if len(args) < len(prefix) {
		return false
	}
	for index := range prefix {
		if args[index] != prefix[index] {
			return false
		}
	}
	return true
}

func packageList(runtime *kernel.Runtime) error {
	for _, manifest := range runtime.PackageManifests() {
		fmt.Printf("%s\t%s\n", manifest.ID, manifest.Version)
	}
	return nil
}

func packageInspect(runtime *kernel.Runtime, id string) error {
	manifest, ok := runtime.PackageManifest(id)
	if !ok {
		return fmt.Errorf("unknown package: %s", id)
	}
	body, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func relayExport(runtime *kernel.Runtime, path string) error {
	body, err := json.MarshalIndent(runtime.ExportBatch(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}

func relayImport(ctx context.Context, runtime *kernel.Runtime, path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var batch grid.Batch
	if err := json.Unmarshal(body, &batch); err != nil {
		return err
	}
	return runtime.ImportBatch(ctx, batch)
}

func relayServe(ctx context.Context, runtime *kernel.Runtime, addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: relayHandler(ctx, runtime),
	}
	fmt.Printf("relay serving on %s\n", addr)
	return server.ListenAndServe()
}

func relayHandler(ctx context.Context, runtime *kernel.Runtime) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /relay/batch", func(writer http.ResponseWriter, _ *http.Request) {
		body, err := json.MarshalIndent(runtime.ExportBatch(), "", "  ")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("X-Moks-Peer-ID", runtime.LocalPeerID())
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(body)
	})
	mux.HandleFunc("POST /relay/import", func(writer http.ResponseWriter, request *http.Request) {
		peerID := strings.TrimSpace(request.Header.Get("X-Moks-Peer-ID"))
		if peerID == "" {
			http.Error(writer, "missing X-Moks-Peer-ID header", http.StatusForbidden)
			return
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		var batch grid.Batch
		if err := json.Unmarshal(body, &batch); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err := runtime.ImportBatchFromPeer(ctx, peerID, batch, "push"); err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"status":"ok"}`))
	})
	return mux
}

func relayPull(ctx context.Context, runtime *kernel.Runtime, peerID string) error {
	peer, ok := runtime.LookupPeer(peerID)
	if !ok {
		return fmt.Errorf("peer not allowed: %s", peerID)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, peer.BatchURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Moks-Peer-ID", runtime.LocalPeerID())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("relay pull failed: %s: %s", response.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var batch grid.Batch
	if err := json.Unmarshal(body, &batch); err != nil {
		return err
	}
	return runtime.ImportBatchFromPeer(ctx, peerID, batch, "pull")
}

func relayPush(ctx context.Context, runtime *kernel.Runtime, peerID string) error {
	peer, ok := runtime.LookupPeer(peerID)
	if !ok {
		return fmt.Errorf("peer not allowed: %s", peerID)
	}
	body, err := json.Marshal(runtime.ExportBatch())
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, peer.ImportURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Moks-Peer-ID", runtime.LocalPeerID())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(response.Body)
		return fmt.Errorf("relay push failed: %s: %s", response.Status, strings.TrimSpace(string(responseBody)))
	}
	return nil
}

func registerBuiltins(runtime *kernel.Runtime) error {
	if err := runtime.RegisterBuiltin(contextpkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(knowledgepkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(inventorypkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(runspkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(linkspkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(maintenancepkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(receivingpkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(procedurespkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(trainingpkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(builtin.OpsPackage()); err != nil {
		return err
	}
	return nil
}

func relayPeerList(runtime *kernel.Runtime) error {
	for _, peer := range runtime.AllowedPeers() {
		fmt.Printf("%s\tpull=%t\tpush=%t\tbatch=%s\timport=%s\n", peer.PeerID, peer.AllowPull, peer.AllowPush, peer.BatchURL, peer.ImportURL)
	}
	return nil
}

func relayPeerAllow(runtime *kernel.Runtime, args []string) error {
	allowPull := args[3] == "pull"
	allowPush := args[4] == "push"
	if args[3] != "pull" && args[3] != "no-pull" {
		return errors.New("usage: relay peer allow <peer-id> <batch-url> <import-url> <pull|no-pull> <push|no-push>")
	}
	if args[4] != "push" && args[4] != "no-push" {
		return errors.New("usage: relay peer allow <peer-id> <batch-url> <import-url> <pull|no-pull> <push|no-push>")
	}
	return runtime.AllowPeer(grid.AllowedPeer{
		PeerID:    args[0],
		BatchURL:  args[1],
		ImportURL: args[2],
		AllowPull: allowPull,
		AllowPush: allowPush,
	})
}
