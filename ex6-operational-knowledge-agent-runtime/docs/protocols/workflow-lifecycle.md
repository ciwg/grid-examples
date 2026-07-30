# OKAR Workflow Lifecycle Protocol

Status: implementation-local profile; pCID assignment pending

## Scope

This document defines the pCID-selected CBOR message shape for retained local
workflow lifecycle events. It implements DI-bavuk. It is an OKAR profile of the
current PromiseGrid envelope candidate, not a claim that the outer envelope is
a frozen global PromiseGrid API.

The protocol records a runtime resource owner's local decision about whether a
workflow artifact is imported, active, deactivated, or revoked. It does not
grant execution authority, make a promise for an app agent, or introduce a new
top-level PromiseGrid action kind.

## pCID assignment

The pCID is the CIDv1 raw-sha2-256 CID of this exact specification document.
The pCID is intentionally not embedded in this document: embedding its own CID
would change the bytes being identified. Before implementation, the finalized
bytes must be hashed, the resulting CIDv1 base32 text must be recorded in a
CID-named protocol-spec symlink, and the exact binary CID bytes must be
hardcoded in the implementation.

Until that step is complete, no runtime may treat this draft as an accepted
lifecycle protocol. Source: DI-bavuk.

## Envelope

Each event is the exact canonical CBOR bytes of:

```cddl
grid( ; CBOR tag 1735551332 / 0x67726964
  [ selector: 42(pCID),
    operation: uint,
    workflow_alias: tstr,
    artifact_cid: bstr,
    parents: [* bstr],
  ]
)
```

`selector` is CBOR tag 42 around the binary CID bytes for this specification's
pCID. `artifact_cid` and every `parents` entry are binary CID bytes, never
printable CID text or raw digest bytes. The envelope's full exact byte sequence
is the CAS object body. Its CID is the lifecycle event CID.

The current PromiseGrid envelope rule is provisional. Any future frozen
replacement has a different pCID and must be implemented as a separate
protocol, not by mutating this protocol's meaning.

## Slot meanings

| Slot | Type | Meaning |
| --- | --- | --- |
| 0 | `42(pCID)` | This protocol specification's binary CID selector. |
| 1 | `uint` | Lifecycle operation defined below. |
| 2 | `tstr` | Non-empty local workflow alias. It is a local operator label, not a global object identity. |
| 3 | `bstr` | Binary CID of the immutable workflow artifact held in local CAS. |
| 4 | `[bstr]` | Parent lifecycle event CIDs for the same workflow artifact. |

## Operations

| Value | Name | Required predecessor |
| --- | --- | --- |
| 0 | `import` | No parent. |
| 1 | `activate` | `import` or `deactivate`. |
| 2 | `deactivate` | `activate`. |
| 3 | `revoke` | `import`, `activate`, or `deactivate`. |

`revoke` is terminal for the artifact timeline. A replacement is a new imported
workflow artifact CID, not reactivation of a revoked artifact.

## Canonical encoding rules

- The outer value MUST use CBOR tag `1735551332`; the inner vector MUST be a
  definite-length array of exactly five items.
- The selector MUST be tag 42 around a byte string containing one complete CID.
- The operation MUST use the shortest canonical unsigned integer encoding.
- The alias MUST be valid non-empty UTF-8 after trimming surrounding Unicode
  whitespace. It MUST NOT be used as a registry key outside the local runtime.
- The artifact CID and each parent CID MUST be complete binary CID bytes. A
  raw SHA-256 digest, printable base32 CID, or `sha256:<hex>` string is invalid
  on the wire.
- The parent vector MUST be a definite-length array. `import` has exactly zero
  parents. Every other operation has exactly one parent until a future pCID
  defines explicit lifecycle reconciliation.
- No extra envelope slots, payload fields, indefinite-length containers, or
  alternate CBOR encodings are valid for this pCID.

## Local acceptance and replay rules

1. The runtime reads exact event bytes from its local CAS; it does not rebuild
   an event from a JSON representation before validation.
2. It verifies canonical CBOR, the outer `grid()` tag, selector pCID, array
   arity, CID encodings, and the operation-specific parent count.
3. For `import`, it verifies the artifact CID resolves in the local CAS. For a
   later operation, it verifies the one parent event resolves locally, validates
   under this same pCID, names the same artifact CID, and has a permitted prior
   operation.
4. It derives one accepted head per artifact CID. Multiple competing valid
   children of one parent remain a local conflict; no projection may silently
   select one. A future pCID may define a reconciliation operation.
5. Unknown pCIDs, invalid arities, invalid parents, unknown artifacts, and
   conflicting heads produce explicit local non-commitments. They never create
   route eligibility or worker execution authority.
6. The local workflow alias is a projection label. Events are indexed by
   artifact CID; an alias collision is a local operator error and does not
   merge artifact timelines.

## Projection cache

A runtime MAY maintain a local cache of accepted artifact heads and derived
workflow state. The cache is not authoritative: deleting it, corrupting it, or
requesting a rebuild causes the runtime to reconstruct it solely from retained
CAS event bytes. A cache entry records at least artifact CID, accepted event
CID, local alias, and derived lifecycle state.

JSONL MAY be emitted for diagnostics or cache export, but it MUST NOT be read as
an authority for event acceptance, lifecycle transitions, or replay.

## Storage and exchange

The runtime stores the original canonical envelope bytes under the event CID in
its local CAS. It retains only the artifacts and event DAG portions it chooses
to store. A node does not assume another node has the same complete CAS.

If an event is exchanged with another agent, its exact stored CBOR bytes travel
over an approved transport as the same `grid()` message. JSON is not a
PromiseGrid inter-agent message encoding. Local cache files and Docker mounts
must not become hidden inter-agent communication channels.

## Error behavior

The implementation returns a local error that identifies the failed validation
class without asserting a global fact. It does not acknowledge an event that it
did not accept, does not mutate an accepted head after failure, and does not
delete retained artifact or event bytes while resolving malformed input.

## Examples

The following is structural notation, not JSON or a wire encoding:

```text
grid([
  42(<workflow-lifecycle-pcid-bytes>),
  0,
  "inventory-receipt",
  <workflow-artifact-cid-bytes>,
  []
])
```

The next event for the same artifact may activate it:

```text
