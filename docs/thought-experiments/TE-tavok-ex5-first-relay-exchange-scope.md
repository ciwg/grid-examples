# Ex5 First Relay-Visible Exchange Scope

TE ID: `TE-tavok`
## Status
decided

## Decision under test

What the first relay-visible PromiseGrid peer-exchange slice for `ex5` should
carry now that five durable families are frozen locally.

Related TODO:

- `100` - `ex5-operational-knowledge-system/TODO/TODO-fativ-ex5-relay-visible-peer-exchange.md`

## Assumptions

- `ex5` now has five frozen PromiseGrid-native local runtime families:
  `knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
  `knowledge-link`, and `knowledge-responsibility`.
- Those families are already signed, replay-verified, and persisted in local
  append-only family logs.
- `knowledge-evidence` currently carries durable metadata plus attachment
  references, but attachment bytes stay on the local copied-file path.
- CAS-backed envelope storage is still an open backlog item under `101`.
- Browser, CLI, and Neovim still bind to the current local HTTP adapter, and
  that embodiment question remains a later backlog item under `102`.

## Alternatives

### Alternative A

Make the first relay-visible exchange slice carry all five frozen families
immediately, including `knowledge-evidence`.

This would maximize early peer-visible coverage but would expose attachment
references whose current form is only locally meaningful.

### Alternative B

Make the first relay-visible exchange slice carry only the attachment-free
families:

- `knowledge-item`
- `knowledge-approval`
- `knowledge-link`
- `knowledge-responsibility`

Defer peer-visible `knowledge-evidence` exchange until CAS-backed storage or a
portable attachment-carriage rule exists.

### Alternative C

Defer all relay-visible peer exchange until CAS-backed storage exists, then
introduce exchange only after the storage layer is redesigned.

## Scope and systems affected

- `ex5-operational-knowledge-system/TODO/TODO-fativ-ex5-relay-visible-peer-exchange.md`
- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`
- `ex5-operational-knowledge-system/docs/architecture.md`
- `ex5-operational-knowledge-system/docs/practical-implementation.md`
- `ex5-operational-knowledge-system/README.md`
- new peer-exchange staging documentation
- later runtime/storage work for `100`, `101`, and `102`

## Scenario analysis

### Scenario 1: normal peer-visible item and approval sharing

Alice and Bob want to exchange a procedure revision, its approval decision, a
typed link, and the responsibility record that owns the work.

Alternative A:

- carries all of those artifacts
- also carries evidence records, even though evidence attachment references are
  not yet portable across hosts
- makes the first exchange slice look broader than the runtime semantics
  actually support

Alternative B:

- carries the item, approval, link, and responsibility artifacts cleanly
- avoids introducing peer-visible attachment references that another host
  cannot resolve
- gives `ex5` one honest first peer-visible slice over already-portable
  artifacts

Alternative C:

- defers even these clearly portable artifacts
- keeps peer exchange blocked on a storage redesign that is not required for
  item/approval/link/responsibility semantics

Result:

- B provides the cleanest first exchange without overstating portability.

### Scenario 2: evidence records with local attachment paths

Carol records evidence with a copied attachment. Dave later receives that
evidence artifact through relay-visible exchange.

Alternative A:

- exposes evidence metadata plus a path reference that only makes sense on
  Carol's host
- creates an implicit obligation either to follow a broken path or to invent
  peer-side rules that have not been specified yet

Alternative B:

- does not exchange peer-visible evidence artifacts yet
- keeps local evidence behavior honest while preserving the later CAS decision
  space

Alternative C:

- also avoids the bad attachment-reference leak
- but blocks the exchange of artifact families that are already portable

Result:

- B avoids a misleading cross-host evidence contract.

### Scenario 3: mixed-version nodes

Ellen runs the current local-only `ex5`; Frank runs the first relay-visible
peer-exchange build.

Alternative A:

- requires broader mixed-version handling immediately, including evidence
  artifacts with unresolved attachment semantics
- increases the chance that older nodes accept peer-visible artifacts they
  cannot represent honestly

Alternative B:

- narrows the first compatibility surface to four portable families
- reduces mixed-version obligations while still proving the relay-visible shape

Alternative C:

- has the smallest compatibility risk
- but leaves the peer-exchange backlog untested longer than necessary

Result:

- B is the best staged compatibility step.

### Scenario 4: trust-boundary changes

Mallory can observe relayed envelopes but cannot forge signatures.

Alternative A:

- increases visible surface area immediately, including evidence metadata that
  may reveal operational context before the repo has settled the right portable
  storage story

Alternative B:

- exposes the four attachment-free durable families first
- keeps the first trust-boundary expansion constrained to artifact types the
  runtime already models most clearly

Alternative C:

- avoids immediate trust-boundary growth
- but also delays validating the relay-visible trust model at all

Result:

- B gives the smallest useful trust-boundary expansion.

### Scenario 5: long-horizon migration

Steve wants `ex5` fully on-grid without repainting the whole runtime twice.

Alternative A:

- risks locking in awkward evidence portability semantics before CAS or blob
  exchange is designed
- could force later search/storage/attachment migration cleanup

Alternative B:

- lets peer exchange start over portable signed families now
- leaves `knowledge-evidence` ready to join later once CAS-backed or otherwise
  portable attachment carriage is specified

Alternative C:

- postpones learning from real peer-visible exchange
- increases the chance that storage-first work is designed in a vacuum

Result:

- B gives the best staged migration path.

## Conclusions

Rejected alternatives:

- Alternative A: too broad for the current evidence attachment semantics
- Alternative C: too conservative; it blocks peer-visible progress on already
  portable signed families

Surviving alternative:

- Alternative B: first relay-visible exchange should carry only the current
  attachment-free families and defer peer-visible `knowledge-evidence`

Implications and future work:

- `100` should close with a staged peer-exchange plan centered on
  `knowledge-item`, `knowledge-approval`, `knowledge-link`, and
  `knowledge-responsibility`
- `101` remains the place to decide CAS-backed storage and any portable
  evidence/blob carriage
- `102` remains later because embodiment contract tightening still depends on
  the peer/storage layers being more concrete

## Decision status

Alternative B locked by `DI-guzab`: the first relay-visible `ex5`
peer-exchange slice carries `knowledge-item`, `knowledge-approval`,
`knowledge-link`, and `knowledge-responsibility`, while peer-visible
`knowledge-evidence` waits for the later storage/carriage decision.
