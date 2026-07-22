# Ex5 CAS Storage Migration Order

TE ID: `TE-nadok`
## Status
decided

## Decision under test

How `ex5` should introduce CAS-backed storage for signed PromiseGrid families
and copied evidence blobs without breaking the current local durable runtime.

Related TODO:

- `101` - `ex5-operational-knowledge-system/TODO/TODO-nobek-ex5-cas-envelope-storage.md`

## Assumptions

- `ex5` currently persists one compatibility `events.jsonl`, five signed family
  logs, a runtime-local `attachments/` tree, and a `drafts/` tree.
- The first relay-visible peer-exchange slice is now staged for the four
  attachment-free families, while peer-visible `knowledge-evidence` is deferred
  because attachment bytes are not yet portable.
- Browser, CLI, and Neovim still use the local HTTP adapter, and embodiment
  tightening remains backlog `102`.

## Alternatives

### Alternative A

Add CAS-backed storage as an additive sidecar first:

- dual-write signed family envelopes into CAS
- dual-write copied evidence blobs into CAS
- keep `events.jsonl`, current family logs, and the current attachment tree as
  compatibility/manifests during migration

### Alternative B

Replace the current family logs and copied attachment tree with CAS as the new
  immediate source of truth in one step.

### Alternative C

Add CAS only for copied evidence blobs first, while leaving signed family
envelopes entirely on the current family logs.

## Scope and systems affected

- `ex5-operational-knowledge-system/TODO/TODO-nobek-ex5-cas-envelope-storage.md`
- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`
- `ex5-operational-knowledge-system/docs/architecture.md`
- `ex5-operational-knowledge-system/docs/practical-implementation.md`
- `ex5-operational-knowledge-system/README.md`
- new CAS staging documentation
- later runtime/storage implementation work for the remaining PromiseGrid
  backlog

## Scenario analysis

### Scenario 1: normal local write path

Alice creates a new item revision, approval, link, responsibility record, and
evidence upload.

Alternative A:

- preserves the current successful local write path
- adds CAS copies for both signed envelopes and copied blobs
- lets the repo verify CAS behavior without forcing every current read path to
  switch immediately

Alternative B:

- forces a full storage cutover at once
- increases the chance of breaking startup replay, current tests, and local
  operators during migration

Alternative C:

- improves blob portability
- but leaves envelopes outside CAS, which weakens the PromiseGrid storage story
  and keeps peer-visible family exchange partially log-bound

Result:

- A gives the safest staged migration while still moving both envelopes and
  blobs into CAS.

### Scenario 2: restart and replay after mixed old/new writes

Bob restarts after some writes predate CAS and some are dual-written.

Alternative A:

- can continue to replay from the current logs while validating CAS population
- gives a clean migration window where old and new records coexist honestly

Alternative B:

- must solve old-to-new migration and new-only startup semantics immediately
- creates avoidable cutover risk

Alternative C:

- avoids full cutover risk
- but still leaves family envelopes outside the new storage regime

Result:

- A best handles mixed historical state.

### Scenario 3: peer-visible evidence portability

Carol wants later peer-visible `knowledge-evidence` exchange.

Alternative A:

- creates the missing portable storage target for both envelope objects and
  copied evidence blobs
- leaves the current local attachment tree available until carriage semantics
  are proven

Alternative B:

- also creates a portable storage target
- but couples it to a high-risk immediate storage rewrite

Alternative C:

- solves blob portability
- but still leaves evidence envelopes outside CAS, making the storage model
  asymmetrical

Result:

- A is the best foundation for later evidence exchange.

### Scenario 4: long-horizon operational complexity

Dave needs another engineer to understand the migration in six months.

Alternative A:

- has a clear staged story:
  1. dual-write
  2. validate
  3. move read paths later
- keeps rollback understandable

Alternative B:

- compresses every obligation into one storage rewrite
- makes rollback and blame assignment harder

Alternative C:

- creates an awkward split where blobs are CAS-native but signed family
  envelopes are not

Result:

- A creates the clearest long-horizon migration contract.

### Scenario 5: trust and integrity

Mallory can inspect files on disk but cannot forge signatures.

Alternative A:

- adds content-addressable integrity hooks for both envelopes and blobs without
  removing the already-verified signed family logs too early

Alternative B:

- eventually gives a cleaner storage story
- but expands migration risk and debugging difficulty during the first pass

Alternative C:

- improves blob integrity only
- keeps envelope-address integrity stuck on the old path

Result:

- A gives the best integrity improvement per unit of migration risk.

## Conclusions

Rejected alternatives:

- Alternative B: too much cutover risk for the first CAS step
- Alternative C: too narrow; it leaves signed envelopes outside the CAS story

Surviving alternative:

- Alternative A: introduce CAS as an additive sidecar for both signed family
  envelopes and copied evidence blobs while keeping current logs and copied
  attachment paths during migration

Implications and future work:

- `101` should close with an additive CAS staging plan
- the next implementation slice can dual-write envelopes and blobs to CAS
  without changing embodiment contracts yet
- `102` still remains later because embodiments do not need a direct runtime
  contract change during the first CAS pass

## Decision status

Alternative A locked by `DI-ribek`: add CAS-backed storage as an additive
sidecar for signed family envelopes and copied evidence blobs, keep the current
family logs and attachment tree as compatibility/manifests during migration,
and defer any embodiment-contract tightening to later work.
