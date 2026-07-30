# CAS-Backed Workflow Lifecycle Events

TE ID: TE-kavuk

## Status

needs DF

## Decision under test

Determine how OKAR represents local workflow import, activation, deactivation,
and revocation so that lifecycle history is content-addressed, protocol-defined,
and replayable without treating a mutable JSONL journal as authoritative state.

This extends TODO `puvok`, DI-lovek, and TE-gavuk. `tools/mint-handle` was not
available in this checkout, so `kavuk` is a manually selected unused proquint
handle for this TE.

## Assumptions and trust model

- A workflow artifact is immutable CAS content. Import does not grant route,
  worker, or execution authority.
- Each lifecycle decision is local to the runtime resource owner. It is not a
  promise made on behalf of an app agent and is not a new top-level PromiseGrid
  action kind.
- When lifecycle events cross an agent boundary, they are PromiseGrid messages:
  exact CBOR bytes in a `grid()` envelope selected by a pCID. The pCID belongs
  to a frozen RFC-like lifecycle protocol specification.
- CAS is authoritative for retained artifacts and lifecycle events. Mutable
  files may accelerate lookup but must be disposable and rebuildable from CAS.
- Alice operates a node and locally decides availability. Bob supplies a valid
  workflow artifact. Mallory may supply malformed events, unknown pCIDs,
  competing histories, or events that Alice does not choose to retain.
- The current `workflow-lifecycle.jsonl` file is a temporary local journal. It
  cannot remain the sole durable authority under PromiseGrid Good Practices.

## Alternatives

### A. CAS-backed `grid()` lifecycle events with per-workflow parent links

Store the exact CBOR `grid()` envelope in CAS. Its pCID identifies a lifecycle
protocol specification. The pCID-defined payload carries the workflow artifact
CID, the lifecycle operation, and zero or more parent event CIDs. Each workflow
artifact has a local event timeline; later events reference the immediately
prior accepted event for that artifact. A local index may cache current heads
and state, but replay reconstructs it from retained CAS event bytes.

### B. CAS-backed `grid()` lifecycle events with one node-wide parent chain

Store the same kind of exact envelopes in CAS, but every local lifecycle event
links to the prior lifecycle event accepted by the node. The node-wide chain is
the replay source; a workflow projection is derived by filtering it.

### C. pCID-labelled DAG-CBOR records without `grid()` envelopes

Store a structured lifecycle record in CAS with a pCID field, but do not encode
it as a `grid()` envelope. A later adapter would translate it before exchange.

### D. Continue JSONL as the authoritative lifecycle journal

Keep lifecycle events in JSONL and optionally archive their content in CAS.

## Scenario analysis

### Normal import and withdrawal

Alice imports Bob's workflow artifact, then later deactivates it. A records two
immutable, pCID-selected envelopes in CAS and derives the current state from
the latter event. The artifact and both decisions remain independently
inspectable. B provides the same retention, but unrelated workflow activity
becomes part of Alice's chain. C retains bytes but has no PromiseGrid-native
wire representation. D is easy locally but makes a mutable file authoritative.

### Corruption and incomplete writes

If a cache write is interrupted, A and B rebuild from previously committed CAS
objects and retain only complete event objects. A malformed object, unknown
pCID, unsupported arity, or broken parent reference produces an explicit local
non-commitment rather than an invented acknowledgement. C has similar CAS
durability but weaker protocol validation. D can truncate the only record of a
state transition and cannot distinguish a partial line from an authoritative
event without extra recovery rules.

### Concurrent actors and mixed versions

Alice may retain her own lifecycle decisions while Bob exports a newer workflow
artifact and Carol uses a different lifecycle-protocol pCID. A permits separate
per-artifact timelines and does not force unrelated conflicts. B serializes
everything through one node timeline, increasing contention and making partial
replication less useful. C requires an adapter agreement before any agent can
interpret the record. D has no pCID to select a parser or define compatibility.

### Trust-boundary changes

Mallory can present a validly encoded event without being trusted by Alice. A
and B let Alice retain, reject, or locally rank events using pCID-defined
semantics and local trust rules; storage does not convert an event into worker
eligibility. C risks application-specific reinterpretation before a protocol
is frozen. D conflates local file possession with durable lifecycle authority.

### Long-horizon evolution and scale

A supports one or more parallel DAG timelines, selective replication, and
per-workflow replay. A node can retain only the workflow histories it chose to
store. B provides simple whole-node auditing but makes replay, export, and
compaction proportional to unrelated activity. C delays the needed protocol
work. D keeps an opaque growing journal and makes cache migration risky.

### Local operational cost

A requires a lifecycle protocol specification, CBOR/grid encoding, CID parent
validation, and a projection/index builder. B needs the same mechanisms but a
simpler linear traversal. C is cheaper now but creates a second non-native
format. D is cheapest now but violates the source-of-truth requirement and
would require a disruptive migration later.

## Conclusions

Alternative D is rejected: JSONL must not remain authoritative. Alternative C
is rejected: pCID metadata without a `grid()` envelope is not a PromiseGrid
message shape and would require a translation layer before agent exchange.

A and B survive. A is recommended because workflow lifecycle is naturally
artifact-scoped, avoids node-wide coupling, supports partial CAS retention, and
matches the guide's parallel-timeline model. A creates obligations to define a
frozen lifecycle protocol, preserve exact CBOR envelope bytes, validate parent
links, and rebuild any mutable index from CAS.

## Decisions still requiring DF

1. Choose the surviving event topology: A per-workflow parent links or B one
   node-wide chain.
2. Choose the event identity exposed by public APIs: retain the current local
   workflow ID as an alias alongside artifact CID, or make artifact CID the
   sole lifecycle key.
3. Choose the first cache posture: no cache (scan retained CAS), or a
   rebuildable local head/projection cache with an explicit reset/rebuild path.
4. Approve the new protocol-spec path and the production code/test paths after
   the selected model determines the exact files.

## Implications for open TODOs and pending DIs

- TODO `puvok` must replace the JSONL-authoritative registry with CAS-backed
  lifecycle events before treating workflow lifecycle as grid-aligned.
- TODO `sibok` must preserve workflow artifact CIDs and lifecycle parent CIDs
  across deactivation, revocation, replacement, and restart.
- DI-lovek remains directionally valid but needs a superseding DI that makes
  CAS lifecycle events authoritative and JSONL only a disposable projection or
  diagnostic export.
