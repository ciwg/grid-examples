# Frozen Family pCID Registry

TE ID: TE-zonal

## Status

decided

## Decision under test

How Ex6 assigns immutable pCIDs to package-record families while allowing people to add, retire, and compose workflows over time.

## Scope and affected systems

This decision applies to the active package families, their manifests, route promises, relay metadata, canonical record fixtures, implementation claims, and package-author documentation. It does not assign a pCID to a workflow instance, package installation, executable, or durable record.

## Assumptions and trust model

- A pCID is the CIDv1 of an immutable protocol specification, not a symbolic label or a hash of a runtime record.
- Durable records keep their exact canonical Grid bytes and declare one family pCID.
- Alice can add a workflow or package locally. Bob can later receive the corresponding records. Mallory can carry opaque bytes but cannot change the meaning of an existing pCID.
- Local policy decides which known-family records to interpret or accept. Relay-carriage evidence remains separate from semantic author evidence.

## Alternatives

### A. One immutable specification per family, with an append-only central registry

Each record family has a standalone specification file. The registry maps its symbolic family name to that file and to the CIDv1 calculated from the file's exact frozen bytes. New semantics add a new specification and registry entry. Existing entries are never rewritten or removed.

### B. One shared specification containing all families

One document defines every active family and produces one collective CID. A change to any family changes the shared document and therefore its CID.

### C. Assign a pCID to every workflow or package

Each workflow, egg, or installed package receives a distinct pCID regardless of whether its record semantics are already covered by a family specification.

## Scenario analysis

### Normal operation: Alice adds a workflow using existing semantics

Under A, Alice's new workflow references the existing procedure, run, evidence, and approval family pCIDs. No protocol change occurs, and Bob can interpret the resulting records using the existing specifications. Under B, no document change is required, but unrelated families still share a revision boundary. Under C, Alice creates needless protocol identities for familiar semantics and Bob must discover package-specific identifiers before understanding a routine workflow.

### New interoperable semantics

If Alice adds a workflow that introduces a new durable record meaning, A adds one versioned standalone specification and one registry entry. Existing family pCIDs remain stable. B changes the shared specification and forces every consumer to reason about whether its collective identifier changed for an unrelated family. C can identify the workflow but fails to distinguish an implementation/package identity from the shared wire meaning Bob must understand.

### Failure, corruption, and incomplete writes

With A, a registry entry is valid only when its spec file and calculated CID agree; an incomplete new entry cannot alter old mappings. With B, a partial edit threatens one document serving all families. With C, independently created package/workflow identifiers give no stable way to validate whether two records actually share semantics.

### Concurrent actors and mixed-version nodes

With A, Alice may publish a record using a newly appended pCID while Bob preserves its exact bytes as unknown until he installs that specification. Existing pCIDs keep their meanings. B couples adoption of one family revision to the whole set. C multiplies discovery and compatibility state for each workflow/package rather than for each semantic contract.

### Long-horizon evolution and retirement

With A, changed semantics receive a new versioned family specification and pCID; retirement is documented as deprecation, never deletion or rewriting. Historical records remain interpretable. B makes unrelated families churn whenever one evolves. C leaves historical readers dependent on the availability and identity of arbitrary application packages.

### Trust boundaries

A preserves the boundary between a public protocol meaning, local acceptance policy, semantic author signatures, and relay signatures. B does not improve that boundary. C incorrectly encourages treating an installed package or workflow identity as protocol authority.

### Scale and operational complexity

A introduces one small file per active family and an append-only registry, but keeps updates narrow and reviewable. B has fewer files but creates a high-contention, high-blast-radius document. C creates unbounded identifiers as ordinary users create workflows and makes routing, fixtures, and documentation needlessly complex.

## Conclusion

Alternative A survives. Alternatives B and C are rejected. Ex6 will maintain one immutable spec file per active record family and an append-only central registry of fixed CIDv1 pCIDs. Workflows normally compose existing family pCIDs. A new pCID is required only for a new shared wire-level semantic contract. Retired workflows or families are deprecated; their frozen specifications and mappings remain available for historical interpretation.

## Decision status

locked by user confirmation on 2026-08-11; DI-jusij records the implementation decision.

## Implications and future work

- Create the frozen family-specification set and registry before wiring runtime constants.
- Make manifests, routes, fixtures, tests, and relay metadata consume registry values.
- Document the workflow-composition rule and the distinction between pCIDs, package identities, workflow definitions, and record IDs.
