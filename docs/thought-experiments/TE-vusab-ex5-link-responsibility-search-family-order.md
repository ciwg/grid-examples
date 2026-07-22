# Ex5 Link, Responsibility, and Search-Metadata Family Order

TE ID: `TE-vusab`
## Status
decided

## Decision under test

How the next PromiseGrid migration step after `knowledge-evidence` should be
shaped:

- whether `knowledge-link`, `knowledge-responsibility`, and
  `knowledge-search-metadata` should all be implemented in one grouped slice, or
- whether `knowledge-link` and `knowledge-responsibility` should be implemented
  now while `knowledge-search-metadata` is explicitly deferred to a later TE.

Related TODOs:

- `097` - `ex5-operational-knowledge-system/TODO/TODO-votek-ex5-fourth-frozen-protocol-family.md`
- `098` - `ex5-operational-knowledge-system/TODO/TODO-sarib-ex5-fifth-frozen-protocol-family.md`
- `099` - `ex5-operational-knowledge-system/TODO/TODO-lomuk-ex5-sixth-frozen-protocol-family.md`

## Assumptions

- `ex5` already has three frozen PromiseGrid families:
  `knowledge-item`, `knowledge-approval`, and `knowledge-evidence`.
- `knowledge-link` currently maps cleanly to the existing `link_added` event.
- `knowledge-responsibility` currently maps cleanly to the existing
  `responsibility_created` event.
- `knowledge-search-metadata` is not currently represented by one dedicated
  durable event family. Search behavior mostly projects metadata already carried
  by item, run, responsibility, place, and resource state.
- Browser, CLI, and Neovim should remain on the current local HTTP adapter
  during these next family freezes.

## Alternatives

### Alternative A

Implement `knowledge-link` and `knowledge-responsibility` now as the next two
clean family freezes, and defer `knowledge-search-metadata` until a later TE
defines a real durable boundary for it.

### Alternative B

Force all three families into one grouped implementation pass now by inventing
or extracting a dedicated `knowledge-search-metadata` durable family alongside
the link and responsibility slices.

## Scope and systems affected

- `protocols/knowledge-link.md`
- `protocols/knowledge-responsibility.md`
- `protocols/knowledge-search-metadata.md`
- `protocols/profiles.go`
- `service/app.go`
- `service/persistence.go`
- `service/types.go`
- new link / responsibility signed-envelope helpers
- service replay tests and tamper tests
- PromiseGrid claims, changelog, and migration docs

## Scenario analysis

### Scenario 1: normal migration from the current runtime

Alice wants the next PromiseGrid step to land with minimal behavior drift for
browser, CLI, and Neovim.

Alternative A:

- uses two event families that already exist cleanly in the current runtime
- requires no new search-specific write path
- keeps the migration family-by-family and close to the shipped behavior
- leaves search metadata on the current projection model for one more round

Alternative B:

- lands more apparent PromiseGrid coverage in one pass
- but forces a search-metadata durable boundary that the current runtime does
  not actually expose as one event family
- risks inventing a synthetic family only to satisfy ordering pressure

Result:

- A is more faithful to the current runtime.

### Scenario 2: replay verification and tamper detection

Bob restarts after links and responsibilities were created across many
workflows.

Alternative A:

- can verify `link_added` and `responsibility_created` directly against replayed
  events
- has one obvious append-only log per family
- keeps verification rules simple and local

Alternative B:

- still verifies links and responsibilities cleanly
- but must also define what counts as a search-metadata durable artifact:
  projected labels only, tags, status, context filters, or all of them
- turns replay verification into a mixed “event family plus projection snapshot”
  problem

Result:

- A creates fewer new verification obligations.

### Scenario 3: mixed-version nodes and migration compatibility

Carol runs a slightly older ex5 node while Dave runs a newer one with the next
family slices.

Alternative A:

- introduces two narrow new signed logs without changing search semantics
- keeps compatibility bridges straightforward
- preserves the rule that search still comes from projected current state

Alternative B:

- introduces a new search-metadata artifact whose exact scope older nodes do not
  already understand
- risks disagreement over whether search should read family logs, projections,
  or both

Result:

- A is cleaner for staged mixed-version migration.

### Scenario 4: long-horizon evolution

Ellen wants `ex5` fully on-grid, but without freezing the wrong search boundary
too early.

Alternative A:

- lets link and responsibility become PromiseGrid-native now
- buys time to decide whether search metadata should remain derived from other
  families or become its own durable family later
- avoids locking a poor long-term search contract just to keep numerical order

Alternative B:

- accelerates the count of frozen families
- but may lock a search family whose semantics are still derivative and not
  independently meaningful
- could create later cleanup debt if search metadata is better modeled as a
  secondary projection over already-frozen families

Result:

- A preserves optionality without stalling real progress.

### Scenario 5: scale and operational complexity

Frank expects larger operational histories, richer link graphs, and more search
filters over time.

Alternative A:

- adds two new logs with obvious storage and verification costs
- avoids adding a search-specific durable stream before it is justified
- keeps the number of newly authoritative data paths small

Alternative B:

- adds more storage and indexing obligations immediately
- makes search metadata authoritative before the product has proven what must be
  durable versus what can remain projected

Result:

- A is lighter operationally and easier to reason about.

## Conclusions

Rejected:

- Alternative B: forcing `knowledge-search-metadata` into the same grouped pass
  would create a search-specific durable family before the current runtime has a
  clean durable boundary for it.

Surviving:

- Alternative A: implement `knowledge-link` and `knowledge-responsibility` now,
  and defer the `knowledge-search-metadata` boundary to a later TE.

## Implications for TODOs and pending DIs

- `097` and `098` can be grouped into one implementation batch if the user
  locks Alternative A.
- `099` should remain open, but its implementation should wait for a dedicated
  search-boundary TE.
- the next DF question should be whether to lock Alternative A and proceed with
  grouped implementation of `knowledge-link` and `knowledge-responsibility`
  now.

## Decision status

Alternative A locked by follow-on DIs `DI-votek`, `DI-sarib`, and `DI-lomuk`:
implement `knowledge-link` and `knowledge-responsibility` now, and leave
`knowledge-search-metadata` for a later dedicated boundary TE.
