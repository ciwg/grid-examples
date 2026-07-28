package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/packages"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/records"
	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/store"
)

type BuiltinCommand func(context.Context, *Runtime, []string) (string, error)
type BuiltinValidator func(records.Envelope) error

type BuiltinPackage struct {
	Manifest   packages.Manifest
	Commands   map[string]BuiltinCommand
	Validators map[string]BuiltinValidator
}

type activePackage struct {
	manifest    packages.Manifest
	builtin     bool
	commands    map[string]BuiltinCommand
	validators  map[string]BuiltinValidator
	external    packages.Runner
	packageRoot string
}

type registeredFamily struct {
	owner        *activePackage
	protocolPCID string
}

type Runtime struct {
	root         string
	packagesRoot string
	history      *store.History
	cas          *store.CAS
	packages     map[string]*activePackage
	commands     map[string]*activePackage
	families     map[string]registeredFamily
}

func Open(root string) (*Runtime, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	history, err := store.OpenHistory(filepath.Join(root, "state"))
	if err != nil {
		return nil, err
	}
	casStore, err := store.OpenCAS(filepath.Join(root, "cas"))
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	runtime := &Runtime{
		root:         root,
		packagesRoot: filepath.Join(root, "packages"),
		history:      history,
		cas:          casStore,
		packages:     map[string]*activePackage{},
		commands:     map[string]*activePackage{},
		families:     map[string]registeredFamily{},
	}
	if err := os.MkdirAll(runtime.packagesRoot, 0o755); err != nil {
		_ = history.Close()
		return nil, err
	}
	return runtime, nil
}

func (runtime *Runtime) Close() error {
	return runtime.history.Close()
}

func (runtime *Runtime) RegisterBuiltin(pkg BuiltinPackage) error {
	if err := pkg.Manifest.Validate(); err != nil {
		return err
	}
	registered := &activePackage{
		manifest:   pkg.Manifest,
		builtin:    true,
		commands:   pkg.Commands,
		validators: pkg.Validators,
	}
	return runtime.activatePackage(registered)
}

func (runtime *Runtime) ActivateInstalled(ctx context.Context, manifestPath string) error {
	manifest, err := packages.LoadManifest(manifestPath)
	if err != nil {
		return err
	}
	if strings.TrimSpace(manifest.Executable) == "" {
		return fmt.Errorf("package %s executable is required for installed packages", manifest.ID)
	}
	described, err := (packages.Runner{Executable: manifest.Executable}).Describe(ctx)
	if err != nil {
		return err
	}
	described.Executable = manifest.Executable
	if !manifest.Equal(described) {
		return fmt.Errorf("manifest self-check mismatch for package %s", manifest.ID)
	}
	return runtime.activatePackage(&activePackage{
		manifest:    manifest,
		external:    packages.Runner{Executable: manifest.Executable},
		packageRoot: filepath.Dir(manifestPath),
	})
}

func (runtime *Runtime) activatePackage(pkg *activePackage) error {
	if _, exists := runtime.packages[pkg.manifest.ID]; exists {
		return fmt.Errorf("package already active: %s", pkg.manifest.ID)
	}
	for _, command := range pkg.manifest.Commands {
		key := command.Key()
		if _, exists := runtime.commands[key]; exists {
			return fmt.Errorf("command already registered: %s", key)
		}
	}
	for _, family := range pkg.manifest.Families {
		if _, exists := runtime.families[family.Name]; exists {
			return fmt.Errorf("family already registered: %s", family.Name)
		}
	}
	runtime.packages[pkg.manifest.ID] = pkg
	for _, command := range pkg.manifest.Commands {
		runtime.commands[command.Key()] = pkg
	}
	for _, family := range pkg.manifest.Families {
		runtime.families[family.Name] = registeredFamily{
			owner:        pkg,
			protocolPCID: family.ProtocolPCID,
		}
	}
	return nil
}

func (runtime *Runtime) PackageManifests() []packages.Manifest {
	manifests := make([]packages.Manifest, 0, len(runtime.packages))
	for _, pkg := range runtime.packages {
		manifests = append(manifests, pkg.manifest)
	}
	packages.SortManifests(manifests)
	return manifests
}

func (runtime *Runtime) PackageManifest(id string) (packages.Manifest, bool) {
	pkg, ok := runtime.packages[id]
	if !ok {
		return packages.Manifest{}, false
	}
	return pkg.manifest, true
}

func (runtime *Runtime) PutCAS(body []byte) (string, error) {
	return runtime.cas.Put(body)
}

func (runtime *Runtime) GetCAS(objectID string) ([]byte, error) {
	return runtime.cas.Get(objectID)
}

// Intent: Validate known families through their owning egg while still
// preserving unknown-family carriage as durable exact bytes for the grid.
// Source: DI-lupok
func (runtime *Runtime) AppendRecord(ctx context.Context, raw []byte) (records.Envelope, error) {
	envelope, err := records.Parse(raw)
	if err != nil {
		return records.Envelope{}, err
	}
	if registered, exists := runtime.families[envelope.Family]; exists {
		if envelope.ProtocolPCID != registered.protocolPCID {
			return records.Envelope{}, fmt.Errorf("family %s expects protocol_pcid %s, got %s", envelope.Family, registered.protocolPCID, envelope.ProtocolPCID)
		}
		owner := registered.owner
		if owner.builtin {
			if validator := owner.validators[envelope.Family]; validator != nil {
				if err := validator(envelope); err != nil {
					return records.Envelope{}, err
				}
			}
		} else {
			if err := owner.external.ValidateEnvelope(ctx, raw); err != nil {
				return records.Envelope{}, err
			}
		}
	}
	return runtime.history.AppendRaw(raw)
}

func (runtime *Runtime) History() []store.StoredEnvelope {
	return runtime.history.Entries()
}

func (runtime *Runtime) ExportBatch() grid.Batch {
	entries := runtime.history.Entries()
	rawRecords := make([]json.RawMessage, 0, len(entries))
	for _, entry := range entries {
		rawRecords = append(rawRecords, append(json.RawMessage{}, entry.Raw...))
	}
	claims := runtime.ImplementationClaims()
	return grid.Batch{
		Format:               grid.RelayBatchFormat,
		Implementation:       "moks-ex6",
		ExportedAt:           time.Now().UTC().Format(time.RFC3339),
		ImplementationClaims: claims,
		Records:              rawRecords,
	}
}

func (runtime *Runtime) ImportBatch(ctx context.Context, batch grid.Batch) error {
	if batch.Format != grid.RelayBatchFormat {
		return fmt.Errorf("unsupported batch format: %s", batch.Format)
	}
	for _, raw := range batch.Records {
		if _, err := runtime.AppendRecord(ctx, raw); err != nil {
			return err
		}
	}
	return nil
}

// Intent: Publish explicit per-package protocol claims so relay exports say
// what the active eggs believe they implement instead of implying that through
// local layout or family names alone. Source: DI-lupok
func (runtime *Runtime) ImplementationClaims() []grid.ImplementationClaim {
	claims := []grid.ImplementationClaim{}
	for _, pkg := range runtime.PackageManifests() {
		for _, claim := range pkg.Claims {
			claims = append(claims, grid.ImplementationClaim{
				PackageID:      pkg.ID,
				PackageVersion: pkg.Version,
				ProtocolPCID:   claim.ProtocolPCID,
				Role:           claim.Role,
				Summary:        claim.Summary,
			})
		}
	}
	return claims
}

func (runtime *Runtime) RunCommand(ctx context.Context, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("command is required")
	}
	for width := len(args); width > 0; width-- {
		key := strings.Join(args[:width], " ")
		owner := runtime.commands[key]
		if owner == nil {
			continue
		}
		if owner.builtin {
			handler := owner.commands[key]
			if handler == nil {
				return "", fmt.Errorf("missing builtin handler for %s", key)
			}
			return handler(ctx, runtime, args[width:])
		}
		return owner.external.RunCommand(ctx, key, args[width:])
	}
	return "", fmt.Errorf("unknown command: %s", strings.Join(args, " "))
}

func (runtime *Runtime) InstallPackageDir(ctx context.Context, sourceDir string) (packages.Manifest, error) {
	manifest, err := packages.LoadManifest(filepath.Join(sourceDir, "moks-package.json"))
	if err != nil {
		return packages.Manifest{}, err
	}
	destination := filepath.Join(runtime.packagesRoot, manifest.ID)
	if err := copyDirectory(sourceDir, destination); err != nil {
		return packages.Manifest{}, err
	}
	if err := runtime.ActivateInstalled(ctx, filepath.Join(destination, "moks-package.json")); err != nil {
		return packages.Manifest{}, err
	}
	return manifest, nil
}

func NewEnvelope(family string, protocolPCID string, recordID string, signer string, payload any) (records.Envelope, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return records.Envelope{}, err
	}
	envelope := records.Envelope{
		Family:       family,
		ProtocolPCID: protocolPCID,
		RecordID:     recordID,
		Signer:       signer,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Payload:      body,
	}
	return envelope, envelope.Validate()
}
