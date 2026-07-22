# Ex5 CAS Authority Beyond Frozen Families

TE ID: `TE-vusor`
## Status
decided

## Decision under test

What remaining ex5 runtime state should move onto authoritative CAS-backed
replay/read behavior next, now that the eight frozen PromiseGrid families
already replay and export authoritatively from CAS.

Related TODO:

- `111` - `ex5-operational-knowledge-system/TODO/TODO-tasem-ex5-cas-authority-beyond-frozen-families.md`

## Assumptions

- `ex5` already ships eight frozen signed families:
  `knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
  `knowledge-link`, `knowledge-responsibility`, `operational-run`,
  `operational-place`, and `operational-resource`.
- Those eight families already dual-write signed envelopes into CAS by
  envelope CID and replay/export authoritatively from CAS.
- Evidence blobs already dual-write into CAS by blob CID.
- Search metadata is already settled as derived projection state, not a
  durable family.
- Live presence is intentionally ephemeral in-memory state, not durable
  runtime state.
- Shared live drafts are still persisted separately under `drafts/*.json` and
  loaded directly from the filesystem instead of from CAS.
- The user wants the most PromiseGrid-complete practical next step, not a
  cosmetic storage shuffle.

## Threat and trust model

- Alice and Bob are honest operators restarting the same local node at
  different times and expecting shared draft continuity.
- Carol imports peer-visible durable history and later resumes authoring
  locally against that imported item set.
- Dave runs an older runtime that still writes `drafts/*.json` directly.
- Mallory can tamper with local draft files, delete manifest files, or leave a
  partially written draft update on disk.
- The runtime should not confuse shared draft state with durable historical
  revision state, and should not claim stronger peer visibility than it has.

## Alternatives

### Alternative A

Make shared live drafts the next authoritative CAS-backed local runtime state,
while keeping them explicitly outside the frozen peer-visible family set.

Under this model:

- each saved draft body is stored in CAS by content CID
- a small local manifest keyed by item ID points to the latest draft CID plus
  its version/updated-at metadata
- startup loads draft state authoritatively from CAS using that manifest
- older direct `drafts/*.json` files may be backfilled into the manifest/CAS
  path during migration
- presence remains in-memory only
- search metadata remains derived projection state

### Alternative B

Promote shared live drafts into a new durable PromiseGrid family now and make
that new family authoritative from CAS.

Under this model:

- drafts stop being merely local embodiment state
- a new frozen draft family or draft-event family is defined
- CAS authority comes from signed draft-family artifacts rather than from a
  local manifest over draft bodies
- later peer-visible draft exchange becomes easier

### Alternative C

Keep authoritative CAS behavior limited to the eight frozen families and leave
all remaining runtime state on the current direct local paths for now.

Under this model:

- `drafts/*.json` remains the direct source of truth for shared draft state
- CAS authority remains a family-only feature
- no broader storage unification is attempted yet

## Scope and systems affected

- `docs/thought-experiments/TE-vusor-ex5-cas-authority-beyond-frozen-families.md`
- `ex5-operational-knowledge-system/TODO/TODO-tasem-ex5-cas-authority-beyond-frozen-families.md`
- `ex5-operational-knowledge-system/service/persistence.go`
- `ex5-operational-knowledge-system/service/app.go`
- `ex5-operational-knowledge-system/service/types.go` if capability metadata changes
- live-draft tests in `service/app_test.go`
- PromiseGrid claims and CAS staging docs

## Scenario analysis

### Scenario 1: normal restart with an active shared draft

Alice edits a live draft in the browser, leaves without snapshotting a durable
revision, then restarts the node later.

Alternative A:

- the current draft body can be restored from a CAS-backed draft body plus a
  small local manifest pointer
- the runtime preserves the existing product claim that shared drafts are
  resumable local working state, not durable historical revisions
- storage semantics become more uniform without changing peer-visible meaning

Alternative B:

- the draft is restored through a new signed family model
- this is stronger if drafts are meant to become part of the durable grid
  contract
- it also changes the meaning of shared drafts substantially and pulls a new
  family into scope

Alternative C:

- restart behavior stays simple and already works
- no additional CAS authority is gained for the last meaningful persisted
  local-only state

Result:

- A strengthens the storage model while preserving the current draft boundary.
- B is stronger only if ex5 is ready to treat drafts as a durable family now.

### Scenario 2: partial write or torn local manifest

Bob saves a draft and the process crashes between writing local pointers and
writing local compatibility state.

Alternative A:

- the runtime can validate that the manifest points at an existing CAS object
- the authoritative draft bytes remain content-addressed even if a compatibility
  JSON file is missing or stale
- the manifest layer still introduces one extra pointer that must be updated
  carefully

Alternative B:

- a signed draft-family append model would give the strongest crash story
- but it would also require deciding draft-family lifecycle, replay, and
  snapshot relationships now

Alternative C:

- a torn `drafts/*.json` write can still corrupt or lose the current draft
  directly
- no content-addressed recovery layer exists for draft bodies

Result:

- A improves crash integrity materially without forcing a new family.
- B is theoretically strongest, but with much more scope and new obligations.

### Scenario 3: mixed-version migration

Dave upgrades first while Ellen still runs a build that only understands the
plain `drafts/*.json` path.

Alternative A:

- new runtimes can backfill older draft files into CAS plus manifest pointers
- older runtimes can keep writing the same local draft file shape during a
  transition if compatibility output remains
- migration cost stays local to one persistence subsystem

Alternative B:

- mixed-version compatibility becomes much harder because older nodes would not
  understand the new draft family at all
- the migration immediately becomes a protocol rollout instead of a storage
  rollout

Alternative C:

- mixed-version operation is simplest because nothing changes
- the broader CAS authority backlog remains open

Result:

- A is the cleanest migration path.

### Scenario 4: long-horizon PromiseGrid evolution

Steve later decides whether shared live drafts should ever become peer-visible
or signed PromiseGrid artifacts.

Alternative A:

- keeps that future choice open
- draft bytes are already content-addressed, which helps later peer-visible
  exchange if the product ever goes there
- it does not prematurely claim that drafts are already part of the durable
  peer contract

Alternative B:

- decides that question now by promoting drafts into a durable family
- makes later peer-visible authoring easier
- also commits ex5 immediately to semantics it has deliberately avoided so far:
  draft history, draft conflict meaning, and draft-family trust boundaries

Alternative C:

- leaves both the storage and protocol questions unresolved

Result:

- A is the strongest staged move unless the product is explicitly ready to
  freeze draft semantics now.

### Scenario 5: scale and operational complexity

Carol edits many large drafts over time while the runtime keeps accumulating
durable artifacts.

Alternative A:

- CAS naturally deduplicates identical draft bodies
- local manifests stay small
- operational complexity is limited to one more manifest/pointer subsystem

Alternative B:

- a draft family would likely accumulate more signed artifacts and replay logic
- this is appropriate only if those draft states truly belong in the durable
  operational graph

Alternative C:

- operational complexity remains low
- deduplication and content-addressability for draft bodies never arrives

Result:

- A provides most of the storage benefits with much lower protocol complexity.

## Conclusions

Rejected alternative:

- Alternative C: it leaves the only meaningful remaining persisted local-only
  state outside authoritative CAS behavior.

Surviving alternatives:

- Alternative A: make shared live drafts authoritative from CAS through a
  manifest-plus-CAS local storage model, while keeping them outside the frozen
  peer-visible family set.
- Alternative B: promote drafts into a new frozen durable family now.

## Implications and future work

- If Alternative A is chosen, TODO 111 should narrow to draft manifest/CAS
  authority and compatibility backfill.
- If Alternative B is chosen, a new protocol-family TODO should likely be
  filed before implementation, because the scope is no longer merely a storage
  change.
- Presence should remain explicitly out of scope for both alternatives.
- Search metadata should remain derived projection state in either outcome.
