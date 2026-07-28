package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/builtin"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/kernel"
	contextpkg "github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages/context"
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
	if err := runtime.RegisterBuiltin(contextpkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(knowledgepkg.Package()); err != nil {
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
