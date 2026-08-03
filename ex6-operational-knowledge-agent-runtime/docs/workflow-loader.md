# Workflow Loader: The Basket

The workflow loader is the runtime-owned, workflow-agnostic basket. It accepts
any valid workflow directory, turns it into one immutable CAS artifact, and
records the local lifecycle decision for that artifact. It does not execute a
workflow or create a new top-level PromiseGrid action kind.

Source: DI-lovek; DI-bavuk.

## The loading flow

    workflow directory
          ↓ validate
    deterministic tar archive
          ↓ CAS
    artifact CID
          ↓ import
    local workflow catalog
          ↓ optional activation
    available local workflow

An allowed peer can additionally transfer an exact artifact plus signed
lifecycle evidence through `moks workflow relay push <alias> <peer-id>`. The
receiver verifies the configured peer identity and signature, retains the
artifact in its own CAS, and retains the sender's lifecycle bytes in
`<runtime root>/workflow-evidence`, with local authenticated sender identities in
`<runtime root>/state/workflow-receipts.json`. It does not import or activate
the workflow: that availability decision remains local. The receiver uses the
CAS-derived inbox commands to inspect one retained artifact and only then makes
an explicit local import decision.

Source: DI-novuk; DI-jifuk; DI-rufir.

The artifact CID is the identity of the exact archived directory. The local
alias is an operator-facing name only. Deactivation or revocation removes local
availability without deleting the artifact or its lifecycle history.

## Retained adapter contracts

The loader retains old workflow artifacts as immutable evidence; upgrading the
runtime does not rewrite their manifests. The shipped trusted adapters preserve
their exact v1 input/output pCID pairs, so an already active historical artifact
continues to run under its declared contract. New captures carry canonical
embedded adapter schemas and reject retired pCIDs. There is no implicit
translation between the two contracts. Source: DI-lumek.

## Local basket state

The moks CLI uses .moks/ beneath its current working directory as its local
runtime root. It holds CAS objects, lifecycle events, the disposable projection
cache, local peer state, and history. It is intentionally ignored by Git:
workflow source directories are committed; the locally loaded basket state is
not.

Inspect the current local basket with:

    go run ./cmd/moks workflow list
    go run ./cmd/moks workflow inspect <alias-or-cid>
    go run ./cmd/moks workflow inbox list
    go run ./cmd/moks workflow inbox inspect <artifact-cid>

## What exists now

- CAS stores immutable bytes by CID.
- The runtime records imported, active, deactivated, and revoked lifecycle
  events as pCID-selected CBOR envelopes.
- Runtime startup rebuilds the disposable local workflow projection from CAS.
- A caller can import an already-present CAS artifact by CID.
- The workflow command family captures any directory with workflow.json,
  imports by alias/CID, lists, inspects, activates, deactivates, and revokes.
- Capture uses deterministic tar bytes and rejects unsafe filesystem entries.

## What a complete loader still needs

### 1. Workflow manifest and validation

Validate workflow.json: workflow ID, version, summary, required
packages, protocol pCIDs, declared inputs/outputs, and policy-file references.
Validation must reject missing required files, unsupported schema versions, and
unmet package/protocol dependencies.

### 2. Deterministic capture

Capture uses the approved deterministic tar format:

- sorted relative paths;
- fixed archive timestamps, ownership, and modes;
- no absolute paths, path traversal, device files, or symlinks;
- explicit inclusion/exclusion rules;
- size and file-count limits.

The same valid directory must always produce the same bytes and artifact CID.

### 3. Catalog and CLI

Expose the basket through a workflow command family:

    moks workflow capture <directory>
    moks workflow import <alias> <artifact-cid>
    moks workflow inbox list
    moks workflow inbox inspect <artifact-cid>
    moks workflow inbox import <artifact-cid> <alias>
    moks workflow list
    moks workflow inspect <alias-or-cid>
    moks workflow verify <alias-or-cid>
    moks workflow status <alias-or-cid>
    moks workflow relay push <alias> <peer-id>
    moks workflow demo <workflow-id>

The demo prints a short operator checklist followed by the final JSON status.
    moks workflow extract <alias-or-cid> <destination>
    moks workflow activate <alias>
    moks workflow deactivate <alias>
    moks workflow revoke <alias>

Capture stores and imports one directory. Import is also available for an
artifact CID obtained from another approved source.

### 4. Safe inspection and retrieval

Inspection must show the artifact CID, local alias, lifecycle state, accepted
event CID, manifest summary, and dependency validation result. Retrieval or
unpacking must use the same path-safety rules as capture and never write outside
an explicitly chosen destination.

### 5. Trust and execution boundary

Loading proves only that local CAS contains a structurally valid artifact.
Activation is a local availability decision. Neither loading nor activation
must silently grant route or worker-execution authority. An installed worker
can run only after its active package's self-checked `workflow_adapters`
declaration exactly matches the active artifact's adapter name and pCIDs. The
runtime sends it CBOR through the Docker confinement boundary and validates its
returned output before applying any proposed writes. Receive-promise
registration and broader trust admission remain later explicit steps. Source:
`DI-fofuh`; `TE-dovek`.

### 6. Tests

The loader needs deterministic-CID tests, invalid-archive tests, dependency
validation tests, restart/catalog rebuild tests, lifecycle withdrawal tests,
and CLI end-to-end tests.

## Demonstrable first slice

The first useful basket demonstration is:

1. Capture a procedure-execution workflow directory.
2. Display its CID and validated manifest.
3. Import it under a local alias.
4. Activate it.
5. Restart the runtime and show the same active artifact in workflow list.
6. Deactivate it and show that its CAS artifact and history remain inspectable.

That demonstrates generic workflow loading without pretending the runtime
already executes workflows.
