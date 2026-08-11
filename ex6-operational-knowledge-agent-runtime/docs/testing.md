# Ex6 Testing and Evidence Guide

This guide describes the verification evidence for Ex6's current PromiseGrid scope. It covers canonical package records, frozen pCID identity, local runtime behavior, relay carriage, workflow adapters, and package boundaries. Source: DI-sidoh, DI-jusij, DI-solan.

## Verification commands

Run these commands from the Ex6 directory:

```sh
go test ./...
go vet ./...
errcheck ./...
```

The two integration proofs are opt-in because they need the local Docker
daemon or a second runtime process. They run from the same current checkout:

```sh
MOKS_DOCKER_INTEGRATION=1 go test ./cmd/moks -run '^TestProcedureExecutionAdapterDockerEndToEnd$' -count=1
MOKS_DEPLOYED_RELAY_INTEGRATION=1 go test ./cmd/moks -run '^TestDeployedTwoNodeWorkflowRelay$' -count=1
```

The Docker proof builds the procedure-execution worker from the checkout,
creates a temporary package descriptor pinned to that build's immutable image
ID, and removes that exact test image at cleanup. It therefore does not use
the checked-in development image ID as test input and cannot silently pass
against a stale local image. The deployed-relay proof starts two temporary
local runtime nodes and exercises their signed relay exchange. Source:
DI-fofuh; DI-bidam.

## What each layer proves

| Layer | Evidence | What it proves |
| --- | --- | --- |
| `records` | `TestPackageFamilyRegistryMatchesFrozenSpecifications` | Every built-in pCID literal equals the CIDv1 of its exact immutable family-spec bytes; every mapping is published in the human registry. |
| `records` | canonical Grid envelope tests | Durable package records use canonical CBOR Grid carriage, canonical JSON payload bytes, RFC3339 timestamps, and correctly shaped author evidence slots. |
| `packages` | external runner tests | Installed packages describe claims, validate canonical records, and return base64 `[][]byte` records; external families must name an explicit frozen pCID. |
| `kernel` | runtime, history, import, and workflow tests | Known families use local validators; unknown-family bytes are retained without inferred semantics; Ex6 rejects unsigned semantic records at append/import; relay carriage remains separate; route promises and workflow operations remain local-policy decisions. |
| `grid` | batch and peer tests | Relay batches carry exact record bytes and relay signatures independently of semantic author evidence. |
| `cmd/moks` | CLI tests | Package install, routes, relay export/import, workflow actions, and the checked-in writer example expose the same runtime contract to operators. |
| opt-in Docker E2E | `TestProcedureExecutionAdapterDockerEndToEnd` | A freshly built, Docker-confined worker receives exact CBOR, satisfies explicit local route evidence, completes a workflow, and yields runtime-signed durable run and procedure-use records. |
| opt-in two-node relay E2E | `TestDeployedTwoNodeWorkflowRelay` | Two separate runtime roots exchange signed relay carriage and preserve the verified workflow evidence without granting automatic execution authority. |
| `go vet` | static analysis | Go code passes standard vet checks. |
| `errcheck` | error-path analysis | Go error returns are not silently dropped. |

## Frozen-spec evidence

The built-in family set is documented in [the pCID registry](protocols/package-family-pcid-registry.md). The runtime contains the corresponding checked-in values, but it does not read Markdown files in production. `TestPackageFamilyRegistryMatchesFrozenSpecifications` recomputes every CIDv1 from the linked immutable spec bytes and fails for a missing, extra, or mismatched mapping.

An operator or package author adding a workflow normally reuses existing family pCIDs and does not recompile Ex6. A new interoperable family must supply its own immutable spec and explicit pCID; a third-party package owns that spec outside the built-in registry. Unknown family bytes remain retained until a local package provides an interpreter and local policy permits its use.

## Local evidence locations

The runtime root selected by the CLI contains local operational evidence, including canonical history bytes, CAS data, relay state, package installations, workflow artifacts, and workflow receipt metadata. These are implementation-local storage locations, not global authority claims. The tests use temporary runtime roots and inspect behavior through the public runtime and CLI surfaces rather than relying on persistent test fixtures.

## Scope and non-claims

Passing these commands does not claim a global identity system, universal authority, automatic workflow execution, global consensus, or automatic trust in a received package or workflow. Ex6 requires semantic author evidence for local package-record admission, but signatures remain evidence rather than authority; relay signatures prove carriage separately. Source: DI-sidoh, DI-jusij, DI-solan, DI-damab.
