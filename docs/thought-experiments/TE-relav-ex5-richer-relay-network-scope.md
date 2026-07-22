# ex5 richer relay-network scope

TE ID: TE-relav
## Status
decided

## Decision under test

What the next relay-network expansion step for `ex5` should be now that the
repo already ships:

- eight frozen signed families
- origin-aware non-bootstrap peer exchange
- canonical create-envelope-CID durable IDs
- CAS-backed authoritative replay/export for the frozen families
- websocket-preferred live drafting inside the local embodiment adapter

The question is no longer whether `ex5` can exchange peer-visible history at
all. It is which relay-network shape most honestly advances the current runtime
beyond the local-adapter export/import layer without reopening already-locked
family and storage decisions.

## Assumptions

- `ex5` already has a stable signed-family layer for items, approvals,
  evidence, runs, places, resources, links, and responsibilities.
- The current local HTTP adapter is still the embodiment surface for browser,
  CLI, and Neovim.
- Live drafting is intentionally separate from durable family exchange.
- Mallory may observe relay traffic and attempt replay or duplication, but
  cannot forge signed family envelopes.
- Mixed-version peers are relevant because the relay layer will arrive after
  the existing local-adapter import/export contract.

## Alternatives

### Alternative A: relay mailbox for whole-bundle exchange

Add a relay-facing mailbox layer that stores and forwards the current
peer-exchange bundles largely as they already exist. Peers publish bundles to a
relay address, and other peers fetch, acknowledge, and import them through the
existing bundle validator.

### Alternative B: incremental relay feed over origin-aware signed records

Add a relay-facing feed protocol that carries origin-aware signed family
records incrementally instead of shipping whole bundles each time. CAS blobs
for evidence travel by CID reference with explicit fetch rules, and peers track
per-origin progress against the relay feed.

### Alternative C: broad relay session protocol that combines durable exchange,
live drafting, and embodiment transport

Use this next slice to introduce a larger relay session model that tries to
cover durable record exchange, live draft collaboration, and possibly future
non-HTTP embodiment contracts at once.

## Scenario analysis

### Scenario 1: normal operation between two peers

Alice and Bob both run current `ex5` instances with some diverged history.

Alternative A makes the first relay rollout easy to explain: Alice exports a
bundle, the relay stores it, and Bob imports it. That fits the existing
runtime, but it repeats whole-family logs even when only a few new records were
added. It also keeps progress tracking coarse, because the unit of delivery is
still “the bundle.”

Alternative B introduces per-origin incremental progress. Alice can publish the
new signed records and blob references since Bob’s last relay cursor, and Bob
can ingest only what is new. That is operationally stronger and more aligned
with the current origin-aware dedupe model, but it adds explicit relay cursor
and record-window obligations.

Alternative C gives the broadest theoretical future, but it makes the first
relay step harder to reason about because durable exchange and live drafting
have different timing, ordering, and failure semantics.

### Scenario 2: duplicate delivery and replay by a noisy relay

Mallory or an unreliable relay redelivers previously seen artifacts.

Alternative A can survive this because bundle import already dedupes by origin
tuple, but replay handling stays expensive because large already-seen bundles
must still be downloaded and revalidated before the importer can discard most
of them.

Alternative B matches the current replay identity model more directly. The feed
can redeliver signed records, but per-origin cursors plus existing
origin-tuple dedupe make the discard path cheaper and clearer.

Alternative C again mixes too many semantics. Replay of durable history and
replay of live-collaboration frames are different classes of problems, and the
first relay slice should not hide that.

### Scenario 3: evidence blobs and CAS carriage

Carol records new evidence with an attachment, and Dave imports it later
through the relay path.

Alternative A can keep using the current inline-bundle blob carriage, but that
means each relay-delivered bundle may carry the same blob bytes again even when
only cursors differ in the surrounding history.

Alternative B can use the current CID-keyed evidence model more cleanly by
letting the durable feed carry the signed evidence record plus a referenced
blob set. That fits the current CAS-backed evidence direction better.

Alternative C tempts the design toward “one stream for everything,” which
blurs durable blob carriage with interactive collaboration payloads before the
repo has a compelling need to unify them.

### Scenario 4: mixed-version peers

Ellen runs the current adapter-only peer-exchange build, while Frank runs the
new relay-aware build.

Alternative A is easier to bridge because the relay layer can still materialize
the existing bundle format that Ellen understands.

Alternative B is still survivable, but it creates a translation obligation:
the relay-aware node or relay service may need to materialize an old-style
bundle for peers that do not understand incremental feed cursors yet.

Alternative C is worst here because mixed-version nodes would have to disagree
not only on relay exchange, but on live transport and embodiment semantics too.

### Scenario 5: long-horizon evolution

The repo later wants richer relay behavior, better selective sync, or direct
non-HTTP embodiment contracts.

Alternative A is serviceable as a stepping stone, but it risks becoming dead
weight. Once selective sync or richer relay policy matters, whole-bundle
mailboxes will look too blunt.

Alternative B creates a cleaner long-term base for later relay policies,
partial catch-up, and direct peer progress tracking. It also composes with the
already-shipped origin-aware exchange model instead of bypassing it.

Alternative C claims too much too early. It may eventually be desirable, but
using it as the next step creates obligations across durable history, live
drafting, and embodiment contracts all at once.

### Scenario 6: operational complexity and testability

The next slice still needs to be implemented, tested, and explained in docs.

Alternative A is simplest to ship quickly and easiest to demo.

Alternative B is more complex, but it stays inside the durable exchange lane
and produces a cleaner final shape.

Alternative C creates the most implementation and testing load while solving
several different problems simultaneously.

## Conclusions

Rejected:

- Alternative C. It combines durable relay exchange with live collaboration and
  embodiment-contract changes too early.

Surviving:

- Alternative A: relay mailbox for whole-bundle exchange
- Alternative B: incremental relay feed over origin-aware signed records

Recommendation:

- Alternative B

Why:

- It is the clearest PromiseGrid expansion after the current shipped scope.
- It builds directly on the already-implemented origin-aware dedupe model.
- It gives a better long-term relay shape than whole-bundle mailboxes.
- It avoids reopening live drafting or non-HTTP embodiment questions in the
  same slice.

## Implications for TODOs and pending DIs

- TODO `115` should lock the relay-network expansion around either a bundle
  mailbox step (`A`) or an incremental origin-aware feed step (`B`).
- If `B` is chosen, the next implementation pass will likely need explicit
  relay cursor, relay envelope window, and CAS blob fetch/storage decisions.
- Future non-HTTP embodiment work should remain separate unless a later TE
  explicitly proves that the relay feed and embodiment contracts should merge.

## Decision status

Alternative B locked by `DI-pazek`: the first richer relay-network slice for
`ex5` is an incremental feed over origin-aware signed records with separate
blob transfer by CID instead of whole-bundle mailbox exchange or a broader
session rewrite.
