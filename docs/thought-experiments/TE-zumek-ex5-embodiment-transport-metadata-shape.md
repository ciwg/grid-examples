# ex5 embodiment transport metadata shape

TE ID: TE-zumek
## Status
decided

## Decision under test

How `ex5` should publish embodiment transport metadata now that the shipped
runtime has different primary and compatibility transport semantics for:

- browser
- CLI
- Neovim

The current `Meta` shape exposes:

- `primary_embodiment_adapter`
- `terminal_embodiment_adapter`
- `browser_live_draft_transport`
- `neovim_live_draft_transport`
- `local_unix_socket_*`

That is better than the older single global transport field, but it still
compresses the embodiment story and leaves “primary vs fallback vs explicit
compatibility opt-in” only partly machine-readable.

## Assumptions

- The underlying shipped transport behavior does not change in this TODO.
- Browser stays on the local HTTP adapter.
- CLI stays on the local Unix socket by default and only uses HTTP through
  explicit `-socket=off`.
- Neovim stays socket-first with websocket and HTTP fallback.
- The goal is to make the metadata contract more honest before changing more
  behavior.

## Alternatives

### Alternative A: extend the current flat field set

Keep the current `Meta` layout and add a few more top-level fields, such as:

- `cli_transport`
- `cli_http_compatibility_mode`
- `browser_http_compatibility_mode`
- `neovim_http_compatibility_mode`

### Alternative B: add one embodiment transport table

Keep existing compatibility fields only if needed, but introduce one explicit
structured table keyed by embodiment, where each embodiment declares:

- primary adapter
- live-draft primary transport
- fallback transports
- whether compatibility transport is implicit fallback or explicit opt-in

### Alternative C: replace the current fields entirely with a new structured
transport table

Remove the existing flat fields and expose only the structured embodiment
transport table.

## Scenario analysis

### Scenario 1: current human-facing docs and tooling

Alice reads `/api/meta` through docs and tests.

Alternative A is easiest to patch into current docs, but it keeps the contract
fragmented across many top-level fields.

Alternative B gives Alice a clearer machine-readable summary without forcing
all current readers to migrate at once.

Alternative C is the cleanest final shape, but it breaks the current metadata
surface immediately and forces a broader migration than this TODO probably
needs.

### Scenario 2: embodiment-specific transport truth

Bob wants to know:

- what browser primary transport is
- what CLI primary transport is
- what Neovim primary transport is
- whether fallback is implicit or explicit

Alternative A answers these questions only indirectly through multiple fields.

Alternative B answers them directly in one embodiment table while still
allowing the older fields to coexist during transition.

Alternative C also answers them directly, but at a higher migration cost.

### Scenario 3: PromiseGrid alignment

Carol wants the runtime to describe embodiment contracts explicitly instead of
leaving them partly in prose.

Alternative A improves the situation, but still feels adapter-era and ad hoc.

Alternative B is the clearest staged PromiseGrid move: embodiment contracts
become explicit objects without forcing an abrupt contract deletion.

Alternative C is the most pure end state, but it is more disruptive than this
safe-first TODO needs.

### Scenario 4: long-horizon evolution

Dave later adds a browser non-HTTP embodiment or another terminal mode.

Alternative A scales poorly because more top-level fields accumulate.

Alternative B scales well because each embodiment record can evolve without
exploding the top-level schema.

Alternative C also scales well, but only after paying the larger migration
cost now.

### Scenario 5: test and doc alignment cost

Ellen wants this TODO to stay a safe refinement.

Alternative A is the lowest code churn, but it leaves the contract less clear.

Alternative B has moderate churn but the best clarity-per-change ratio.

Alternative C has the highest churn because every existing metadata consumer
must switch at once.

## Conclusions

Rejected:

- Alternative A. It improves clarity, but not enough.

Surviving:

- Alternative B: add one embodiment transport table alongside the current
  fields during transition
- Alternative C: replace the current fields entirely with a new structured
  transport table

## Implications for TODOs and pending DIs

- TODO `129` is locked to Alternative `C`.
- `/api/meta` should expose one top-level `embodiments` table keyed by
  `browser`, `cli`, and `neovim`.
- Tests and docs should align to that replacement metadata shape in the same
  batch.
