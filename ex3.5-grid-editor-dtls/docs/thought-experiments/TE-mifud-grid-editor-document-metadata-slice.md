# grid-editor document metadata slice

TE ID: TE-mifud
## Status
decided

## Decision under test

How should `grid-editor` add its first PromiseGrid-native document metadata
slice for description, tags, collections, favorites, archive, and search
without turning the relay into a full authoritative product server?

## Assumptions

- The live CRDT editing path remains `live-document`.
- Publish/import is already a separate durable relay path.
- Alice, Bob, and Carol are cooperative users; Mallory is the adversary.
- The next useful backend slice is document metadata, not permissions or
  restore.

## Alternatives

### Alternative 1

Keep document metadata browser-local in the existing Phase 2 registry.

Meaning:
- titles, descriptions, tags, favorites, collections, and archive state all
  live only in browser storage
- search works only over browser-local state
- relay remains unaware of metadata

### Alternative 2

Add a new relay-signed `document-metadata` protocol with latest-state
semantics per document.

Meaning:
- relay accepts local metadata updates
- relay signs metadata envelopes
- relay stores them durably and shares them with peers
- browser search queries relay-known metadata instead of inventing a second
  authoritative store

### Alternative 3

Build a full authoritative document-management registry now.

Meaning:
- relay owns canonical document metadata, search, collections, favorites,
  archive rules, and maybe future ownership semantics all at once
- browser and Neovim become thinner clients over that registry

## Scenario analysis

### S1 — Alice wants a stable description and tags visible on another relay

Alternative 1 fails because the metadata never leaves Alice's browser.

Alternative 2 succeeds because the metadata is a signed relay artifact that
other relays can ingest and serve back.

Alternative 3 also succeeds, but adds more server authority than needed.

### S2 — Bob wants to search documents he has learned through relay traffic

Alternative 1 only searches Bob's own browser-local cache.

Alternative 2 lets Bob's browser ask the relay for relay-known metadata search
results. That is enough for the first slice without inventing a global index.

Alternative 3 can support richer search, but it commits early to a larger
registry model that will collide with later permissions decisions.

### S3 — Alice archives a document

Alternative 1 makes archive state local and non-shareable.

Alternative 2 treats archive as current-time signed metadata. That keeps it
 durable and shareable, but still narrow.

Alternative 3 folds archive into a larger registry with more policy than this
slice needs.

### S4 — Mallory tampers with document labels or favorites

Alternative 1 gives the weakest provenance because the relay never signs the
metadata.

Alternative 2 preserves the existing signed-envelope evidence path.

Alternative 3 can do the same, but with a larger attack surface.

### S5 — Future permissions and restore

Alternative 1 becomes a dead end.

Alternative 2 composes well: permissions can later decide who may emit
metadata updates, and restore can reference metadata without replacing it.

Alternative 3 risks front-loading policy and ownership semantics before those
separate decisions are tested.

## Conclusions

- Reject Alternative 1 because browser-local metadata does not satisfy the
  PromiseGrid-native backend goal.
- Reject Alternative 3 because it bundles too much future policy and
  authority into the first metadata slice.
- Keep Alternative 2 as the surviving design.

## Surviving alternative

Alternative 2:
- add a relay-signed `document-metadata` protocol
- store latest metadata state durably
- let search operate over relay-known metadata
- keep titles and older Phase 2 browser-only fields as compatibility state
  until a later migration explicitly replaces them

## New decisions exposed

- whether metadata should use latest-state or append-only semantics
- whether search should be browser-local, relay-local, or global
- whether favorites/collections belong in the same metadata protocol

## Implications for TODOs and DIs

- add a new Phase 4 metadata TODO
- add a new `document-metadata` protocol spec
- document the split between relay-backed metadata and browser-local workflow
  metadata

## Decision status

locked by `DI-loruk` and `DI-sukip`
