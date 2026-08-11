# Ex6 PromiseGrid Implementation Claims

Status: current verified scope

Ex6 is an unreleased local operational-knowledge runtime that demonstrates PromiseGrid-aligned durable evidence, protocol selection, local policy, package extension, relay carriage, and workflow composition. It is not a generalized ERP, global identity system, consensus network, or universal authority framework. Source: DI-sidoh, DI-jusij, DI-solan.

## Implemented evidence boundary

Package-defined durable records are canonical Grid CBOR bytes using the shared [package-record carriage](protocols/package-record-v1.md). The envelope carries family name, record ID, signer label, RFC3339 timestamp, canonical JSON payload bytes, and optional semantic author evidence. A complete author signature signs the same canonical envelope with author slots empty. It is evidence interpreted under local policy, not an automatic authorization grant.

The runtime keeps semantic author evidence separate from relay-carriage signatures. Relay batches carry exact record bytes and can be signed by the carrying peer without changing who authored a semantic statement. Unknown family records are preserved as exact bytes; the runtime does not invent a validator or semantics for them.

## Frozen built-in protocols

Ex6 implements 27 built-in package-record families. Each has one immutable family specification and a fixed CIDv1 pCID in the [package-family registry](protocols/package-family-pcid-registry.md). The checked-in Go registry is verified against those exact files by the test suite. Built-in manifests, validator routes, relay metadata, record fixtures, and package implementations consume those fixed values.

The frozen set covers the context, knowledge, links, procedures, runs, inventory, receiving, quarantine, corrective-action, maintenance, training, and operational-note families listed in the registry. The family specification defines payload meaning; the shared carriage specification defines canonical record encoding and evidence slots.

## Workflow and package extension

A workflow is a local composition of package capabilities and existing family pCIDs. Creating, editing, enabling, disabling, or retiring a workflow does not assign a pCID by default and does not require recompiling Ex6.

A new interoperable durable-record meaning requires a new versioned immutable specification and its CIDv1 pCID. Third-party packages own their own family specifications and present explicit pCIDs in their manifests and output records. Installing a package may provide a local validator/interpreter, but does not create global authority or automatically activate a route promise. The runtime preserves unknown pCID bytes until local code and policy choose to interpret them.

## Embodiments

- **CLI and local runtime:** package commands, local durable storage, route promises, relay import/export, and workflow lifecycle/run operations use one local runtime root.
- **Installed packages:** external executables describe their claims, validate records, and return canonical record bytes wrapped as base64 `[][]byte` JSON transport.
- **Relay:** peer and batch behavior is a separate carriage mechanism, not the application runtime and not proof of semantic authority.
- **Workflow adapters:** installed adapters may propose canonical records only through their declared contract and local validation/policy checks; receiving, importing, or activating a workflow does not by itself execute it.

## Explicit non-claims and deferrals

Ex6 does not claim global agent identity, a universal trust root, automatic authority from a signature, automatic authority from package installation, automatic workflow activation/execution, consensus, distributed transaction semantics, or support for arbitrary unknown records beyond exact-byte retention. Workflow role effects, approvals, route selection, signature acceptance, and adapter execution remain explicit local-policy decisions. A semantic or carriage change requires a new versioned protocol specification and pCID.

## Verification

See [testing.md](testing.md) for the commands and evidence each test layer supplies. The current verification gate is `go test ./...`, `go vet ./...`, and `errcheck ./...`.
