# ex5 next PromiseGrid transport substrate

TE ID: TE-sorav
## Status
decided

## Decision under test

What the next reusable PromiseGrid substrate slice should be after the durable
record core moved into `promisegrid/records/`.

This TE corresponds to TODO `ravem.1` / `ravem.2` / `ravem.3`.

## Assumptions

- `ex5` already ships one real reusable substrate boundary under
  `promisegrid/records/`.
- The next extraction should be proven by shipped behavior, not by speculative
  framework design.
- Browser, CLI, Neovim, peer exchange, and relay feed all already depend on a
  real direct-contract and relay-visible transport layer.
- `service/` should keep ownership of ex5-specific persistence, projections,
  workflow composition, and user-facing semantics.

## Alternatives

### A. Stop after `promisegrid/records/`

Keep the record substrate as the only reusable layer for now.

### B. Extract the wire transport substrate next

Extract the pure peer-exchange and relay-feed wire shapes plus origin-aware
filtering helpers into a new `promisegrid/transport/` package built on top of
the already-reusable record types.

### C. Extract local socket/browser/native-host framing first

Treat the embodiment contract framing as the next reusable layer before the
relay/feed wire types.

## Scenario analysis

### Scenario 1: normal peer and relay exchange

Alice exports a peer bundle, Bob imports it, and Carol publishes incremental
relay batches.

- A keeps the system working, but the relay/feed wire contract remains owned
  only as `ex5` app code.
- B extracts the already-shipped transport truth that multiple ex5 surfaces
  now depend on, without generalizing ex5 search/review/workflow semantics.
- C extracts an embodiment-facing layer first, even though the peer/relay wire
  contract is the cleaner shared substrate slice today.

### Scenario 2: mixed-version maintenance

Dave changes relay filtering or origin-key logic later.

- A keeps those mechanics coupled to the ex5 app layer.
- B makes the relay/feed wire contract testable as a substrate package in its
  own right.
- C helps direct embodiments, but leaves the network-carried transport truth
  less explicit than it should be.

### Scenario 3: PromiseGrid boundary honesty

Ellen asks what is reusable substrate and what is still ex5-specific.

- A says only record durability is reusable so far.
- B says durable records plus transport wire shapes are reusable, while
  projections and workflows remain ex5-specific.
- C says embodiment framing is reusable before the relay/feed contract itself,
  which is a weaker foundation ordering.

## Conclusions

Rejected:

- Alternative A: too conservative now that relay/feed carriage is a real,
  shipped, reusable contract.
- Alternative C: extracts a less foundational layer before the network-visible
  wire contract that already has stronger reuse evidence.

Surviving:

- Alternative B: extract the peer-exchange and relay-feed wire layer next.

## Implications for TODOs and pending DIs

- TODO `141` is locked to Alternative `B`.
- The new reusable slice should live under `promisegrid/transport/`.
- `service/` should consume that package directly instead of continuing to own
  duplicate transport wire structs where the boundary is now clear.
