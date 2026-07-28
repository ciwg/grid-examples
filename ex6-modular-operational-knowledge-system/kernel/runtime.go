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
	peers        *grid.PeerStore
	policies     *grid.PolicyStore
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
	peerStore, err := grid.OpenPeerStore(filepath.Join(root, "state"), root)
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	policyStore, err := grid.OpenPolicyStore(filepath.Join(root, "state"))
	if err != nil {
		_ = history.Close()
		return nil, err
	}
	runtime := &Runtime{
		root:         root,
		packagesRoot: filepath.Join(root, "packages"),
		history:      history,
		cas:          casStore,
		peers:        peerStore,
		policies:     policyStore,
		packages:     map[string]*activePackage{},
		commands:     map[string]*activePackage{},
		families:     map[string]registeredFamily{},
	}
	if err := os.MkdirAll(runtime.packagesRoot, 0o755); err != nil {
		_ = history.Close()
		return nil, err
	}
	if err := runtime.activateInstalledFromRoot(context.Background()); err != nil {
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

func (runtime *Runtime) LocalPeerID() string {
	return runtime.peers.LocalPeerID()
}

func (runtime *Runtime) LocalPeerPublicKey() string {
	return runtime.peers.LocalPublicKey()
}

func (runtime *Runtime) AllowedPeers() []grid.AllowedPeer {
	return runtime.peers.AllowedPeers()
}

func (runtime *Runtime) AllowPeer(peer grid.AllowedPeer) error {
	return runtime.peers.SetAllowedPeer(peer)
}

func (runtime *Runtime) RevokePeer(peerID string) error {
	return runtime.peers.RemoveAllowedPeer(peerID)
}

func (runtime *Runtime) LookupPeer(peerID string) (grid.AllowedPeer, bool) {
	return runtime.peers.Lookup(peerID)
}

func (runtime *Runtime) ClaimPolicies() []grid.ClaimTrustPolicy {
	return runtime.policies.ClaimPolicies()
}

func (runtime *Runtime) SetClaimPolicy(policy grid.ClaimTrustPolicy) error {
	for _, peerID := range policy.AllowedAttesters {
		if _, ok := runtime.LookupPeer(peerID); !ok {
			return fmt.Errorf("unknown attester peer: %s", peerID)
		}
	}
	return runtime.policies.SetClaimPolicy(policy)
}

func (runtime *Runtime) RemoveClaimPolicy(protocolPCID string, role string) error {
	return runtime.policies.RemoveClaimPolicy(protocolPCID, role)
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
	return runtime.appendRecord(ctx, raw, true)
}

func (runtime *Runtime) appendRecord(ctx context.Context, raw []byte, signLocal bool) (records.Envelope, error) {
	envelope, err := records.Parse(raw)
	if err != nil {
		return records.Envelope{}, err
	}
	if signLocal && !envelope.HasAuthorSignature() {
		// Intent: Sign locally authored durable records once at creation time so
		// later relay trust can distinguish semantic authoring from carriage.
		// Source: DI-sovem
		envelope, err = runtime.peers.SignAuthorEnvelope(envelope)
		if err != nil {
			return records.Envelope{}, err
		}
		raw = records.MustMarshal(envelope)
	}
	// Intent: Verify embedded semantic author signatures before durable storage
	// so bad author proofs are rejected even outside relay exchange.
	// Source: DI-sovem
	if err := runtime.peers.VerifyAuthorEnvelope(envelope); err != nil {
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
	envelope, _, err = runtime.history.AppendRaw(raw)
	return envelope, err
}

func (runtime *Runtime) History() []store.StoredEnvelope {
	return runtime.history.Entries()
}

func (runtime *Runtime) ExportBatch() (grid.Batch, error) {
	entries := runtime.history.Entries()
	rawRecords := make([]json.RawMessage, 0, len(entries))
	for _, entry := range entries {
		rawRecords = append(rawRecords, append(json.RawMessage{}, entry.Raw...))
	}
	claims := runtime.ImplementationClaims()
	claimProofs, err := runtime.peers.SignClaims(claims)
	if err != nil {
		return grid.Batch{}, err
	}
	recordSignatures, err := runtime.peers.SignRecords(rawRecords)
	if err != nil {
		return grid.Batch{}, err
	}
	return grid.Batch{
		Format:               grid.RelayBatchFormat,
		Implementation:       runtime.LocalPeerID(),
		ExportedAt:           time.Now().UTC().Format(time.RFC3339),
		ImplementationClaims: claims,
		ClaimProofs:          claimProofs,
		Records:              rawRecords,
		RecordProofs:         grid.ProofsForRecords(rawRecords),
		RecordSignatures:     recordSignatures,
	}, nil
}

func (runtime *Runtime) SignedExportBatch() (grid.Batch, error) {
	batch, err := runtime.ExportBatch()
	if err != nil {
		return grid.Batch{}, err
	}
	return runtime.peers.SignBatch(batch)
}

func (runtime *Runtime) AttestBatchClaims(batch grid.Batch) (grid.Batch, error) {
	if batch.Implementation == runtime.LocalPeerID() {
		return grid.Batch{}, errors.New("third-party claim attestation requires a different peer")
	}
	attestations, err := runtime.peers.AttestClaims(batch.ImplementationClaims)
	if err != nil {
		return grid.Batch{}, err
	}
	batch.ClaimAttestations = append(batch.ClaimAttestations[:0], attestations...)
	return batch, nil
}

func (runtime *Runtime) ImportBatch(ctx context.Context, batch grid.Batch) error {
	// Intent: Treat the current relay shell as idempotent exact-byte carriage so
	// repeated imports stop re-appending identical durable records while malformed
	// batch metadata is rejected before touching local history. Source: DI-sibok
	if err := batch.Validate(); err != nil {
		return err
	}
	// Intent: Verify per-record digests before durable import so receivers can
	// reject tampered relay contents even when they do not yet understand the
	// record family semantics. Source: DI-zumep
	if err := batch.VerifyClaimProofs(); err != nil {
		return err
	}
	if err := batch.VerifyClaimAttestations(); err != nil {
		return err
	}
	// Intent: Verify exported implementation claims against the exporting peer's
	// key material before treating those claims as trustworthy batch metadata.
	// Source: DI-luzef
	if err := runtime.peers.VerifyClaimProofs(batch); err != nil {
		return err
	}
	// Intent: Verify outside countersigners for implementation claims so import
	// can distinguish exporter self-claims from third-party attestation.
	// Source: DI-fogem
	if err := runtime.peers.VerifyClaimAttestations(batch); err != nil {
		return err
	}
	// Intent: Require local claim-attestation quorum only when local runtime
	// policy says a claim needs it, so trust remains explicit and operator-owned.
	// Source: DI-movek
	if err := runtime.policies.VerifyClaimPolicies(batch, runtime.AllowedPeers()); err != nil {
		return err
	}
	if err := batch.VerifyRecordProofs(); err != nil {
		return err
	}
	// Intent: Verify relay-carriage signatures per record so import does not rely
	// solely on the enclosing batch signature for transport trust. Source: DI-ravud
	if err := runtime.peers.VerifyRecordSignatures(batch); err != nil {
		return err
	}
	for _, raw := range batch.Records {
		if _, err := runtime.appendRecord(ctx, raw, false); err != nil {
			return err
		}
	}
	return nil
}

// Intent: Keep live multi-peer exchange behind explicit allow rules and reject
// batches whose claimed peer identity or signature does not match the allowed
// peer entry.
// Source: DI-zotem
func (runtime *Runtime) ImportBatchFromPeer(ctx context.Context, peerID string, batch grid.Batch, direction string) error {
	if strings.TrimSpace(peerID) == "" {
		return errors.New("peer id is required")
	}
	peer, ok := runtime.LookupPeer(peerID)
	if !ok {
		return fmt.Errorf("peer not allowed: %s", peerID)
	}
	switch direction {
	case "pull":
		if !peer.AllowPull {
			return fmt.Errorf("peer %s is not allowed for pull", peerID)
		}
	case "push":
		if !peer.AllowPush {
			return fmt.Errorf("peer %s is not allowed for push", peerID)
		}
	default:
		return fmt.Errorf("unknown peer direction: %s", direction)
	}
	if batch.Implementation != peerID {
		return fmt.Errorf("batch implementation %s does not match peer %s", batch.Implementation, peerID)
	}
	if err := runtime.peers.VerifyPeerBatch(peerID, batch); err != nil {
		return err
	}
	return runtime.ImportBatch(ctx, batch)
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
		result, err := owner.external.RunCommand(ctx, key, args[width:])
		if err != nil {
			return "", err
		}
		return runtime.applyExternalCommandResult(ctx, result)
	}
	return "", fmt.Errorf("unknown command: %s", strings.Join(args, " "))
}

// Intent: Keep durable writes basket-mediated even for installed eggs so
// external packages can extend the system without bypassing runtime-owned CAS
// and append-only history. Source: DI-rovum
func (runtime *Runtime) applyExternalCommandResult(ctx context.Context, result packages.CommandResult) (string, error) {
	replacements := map[string]string{}
	for _, write := range result.CAS {
		if strings.TrimSpace(write.Alias) == "" {
			return "", errors.New("cas alias is required")
		}
		objectID, err := runtime.PutCAS([]byte(write.Body))
		if err != nil {
			return "", err
		}
		replacements["$cas:"+write.Alias] = objectID
	}
	for _, raw := range result.Records {
		replaced, err := replaceCASAliases(raw, replacements)
		if err != nil {
			return "", err
		}
		if _, err := runtime.AppendRecord(ctx, replaced); err != nil {
			return "", err
		}
	}
	return result.Output, nil
}

func replaceCASAliases(raw []byte, replacements map[string]string) ([]byte, error) {
	if len(replacements) == 0 {
		return raw, nil
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	value = replaceAliasesRecursive(value, replacements)
	return json.Marshal(value)
}

func replaceAliasesRecursive(value any, replacements map[string]string) any {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			typed[key] = replaceAliasesRecursive(child, replacements)
		}
		return typed
	case []any:
		for index, child := range typed {
			typed[index] = replaceAliasesRecursive(child, replacements)
		}
		return typed
	case string:
		if replacement, ok := replacements[typed]; ok {
			return replacement
		}
		return typed
	default:
		return value
	}
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

// Intent: Re-activate installed eggs from the runtime-owned package root on
// startup so installation survives later CLI invocations and process restarts.
// Source: DI-rovum
func (runtime *Runtime) activateInstalledFromRoot(ctx context.Context) error {
	entries, err := os.ReadDir(runtime.packagesRoot)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		manifestPath := filepath.Join(runtime.packagesRoot, entry.Name(), "moks-package.json")
		if _, err := os.Stat(manifestPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
		if err := runtime.ActivateInstalled(ctx, manifestPath); err != nil {
			return err
		}
	}
	return nil
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
