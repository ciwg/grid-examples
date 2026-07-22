# Ex5 Peer-Visible Entity Namespace Reconciliation

TE ID: `TE-loruk`
## Status
decided

## Decision under test

How `ex5` should reconcile peer-visible durable entity identity across
independent runtimes once those runtimes have already minted overlapping
local-facing IDs such as `ITEM-0001`, `RECV-0001`, `RUN-0001`, `RESP-0001`,
and `LINK-0001`.

Related TODO:

- `109` - `ex5-operational-knowledge-system/TODO/TODO-zavok-ex5-peer-visible-entity-namespace-reconciliation.md`

## Assumptions

- `ex5` already ships origin-aware event identity via
  `(origin_peer_id, origin_sequence)`.
- `ex5` already ships six frozen signed families and non-bootstrap peer import
  for unseen origin tuples.
- The current runtime still uses locally minted entity IDs as durable keys in
  events, records, and projections.
- The current importer honestly rejects colliding create-event IDs rather than
  silently merging them.
- Browser, CLI, and Neovim surfaces currently show those local-facing IDs
  directly.
- The goal is the most PromiseGrid-complete durable identity model that still
  has a practical migration path from the current runtime.

## Threat and trust model

- Alice and Bob are honest peers with independent runtimes that may mint the
  same local-facing IDs.
- Carol is an honest mixed-version peer exchanging history during migration.
- Dave imports the same or overlapping history multiple times.
- Mallory may replay bundles, reorder valid deliveries, or attempt to exploit
  ambiguous entity identity to misbind approvals, links, evidence, or runs.
- The runtime must not silently attach imported records to the wrong existing
  entity because two peers happened to mint the same human-facing ID.

## Alternatives

### Alternative A

Replace the current durable entity IDs with globally stable canonical IDs, and
treat the existing `ITEM-*`, `RUN-*`, `RESP-*`, and similar values as
presentation aliases only.

Under this model:

- each entity gets a new canonical durable ID, likely derived from:
  - origin peer identity
  - origin create-event order
- events, frozen-family payloads, projections, and exchange all use the
  canonical durable ID
- local-facing IDs remain as display aliases or short labels

### Alternative B

Keep the current local-facing IDs as the embodiment display surface, but add a
separate canonical peer-stable entity key internally and on the wire.

Under this model:

- each entity keeps its current local-facing ID for browser/CLI/Neovim display
- each entity also gains a canonical peer-stable key, likely derived from:
  - create-event `origin_peer_id`
  - create-event `origin_sequence`
- durable references between records use the canonical key
- adapters can still show short local-facing IDs as aliases

### Alternative C

Keep the current entity IDs as the only durable IDs and resolve collisions with
import-side rewrite or translation tables only.

Under this model:

- imported peers that collide on `ITEM-*`, `RUN-*`, `RESP-*`, and similar IDs
  are remapped locally during import
- the runtime maintains translation state to preserve equivalence across peers
- no new canonical entity key is introduced into the durable family model

## Scope and systems affected

- `docs/thought-experiments/TE-loruk-ex5-peer-visible-entity-namespace-reconciliation.md`
- `ex5-operational-knowledge-system/TODO/TODO-zavok-ex5-peer-visible-entity-namespace-reconciliation.md`
- `ex5-operational-knowledge-system/service/types.go`
- `ex5-operational-knowledge-system/service/app.go`
- `ex5-operational-knowledge-system/service/peer_exchange.go`
- frozen-family envelope helpers and payload schemas
- search, detail routes, and embodiment displays that surface entity IDs
- replay/import tests
- PromiseGrid claims, architecture docs, and UI guide wording where IDs are described

## Scenario analysis

### Scenario 1: two independent peers both create `ITEM-0001`

Alice and Bob each create one receiving-check item before any exchange. Both
now have a valid local `ITEM-0001`.

Alternative A:

- solves the collision directly because the durable identity stops being
  `ITEM-0001`
- imported history can coexist without translation ambiguity
- existing embodiment surfaces must learn to show aliases or short IDs instead
  of relying on the canonical durable ID directly

Alternative B:

- also solves the collision directly because the canonical peer-stable key is
  separate from the displayed local-facing ID
- lets embodiments keep their current operator-facing short IDs longer
- still requires the runtime to move durable references onto the canonical key

Alternative C:

- keeps the wrong thing durable
- import must rewrite or translate colliding IDs and then keep that mapping
  forever
- makes future peer exchange and debugging harder because the durable ID still
  changes by environment

Result:

- A and B both solve the real collision; C adds ongoing translation debt.

### Scenario 2: imported approvals, links, and evidence must attach to the
correct entity

Carol imports a run approval, a typed link, and an evidence record that all
reference a colliding local-facing ID.

Alternative A:

- approval/link/evidence attachment becomes straightforward once all durable
  references use the canonical global ID
- embodiments can still show short aliases after resolution

Alternative B:

- also gives one canonical durable join key for all references
- keeps the operator-facing surfaces more stable because the old short IDs can
  remain visible as aliases

Alternative C:

- relies on the correctness of translation tables everywhere
- every route, every projection, and every imported reference must remember
  which rewritten local ID corresponds to which remote source
- a missed translation becomes a silent semantic bug

Result:

- A and B give one real join key; C keeps a fragile indirection layer.

### Scenario 3: mixed-version migration

Dave upgrades to the reconciled model while Ellen still runs the older local-ID
runtime.

Alternative A:

- is the cleanest end state
- but it is also the most disruptive migration because the old visible IDs stop
  being the durable identity immediately

Alternative B:

- is easier to stage
- older embodiments can still keep showing the short local-facing IDs while the
  runtime shifts durable references to canonical keys underneath
- gives the repo a compatibility bridge without keeping the wrong durable model

Alternative C:

- seems easiest at first because visible IDs do not change
- but it forces the runtime to maintain compatibility translation forever

Result:

- B gives the smoothest migration path.

### Scenario 4: long-horizon PromiseGrid completeness

Steve wants `ex5` to be fully on the grid rather than merely tolerant of some
imports.

Alternative A:

- is the most direct grid-native outcome
- canonical global IDs become the obvious durable identity for every entity
- but the embodiment impact is heavier

Alternative B:

- is nearly as strong at the durable runtime layer
- preserves a cleaner human-facing adapter story by separating:
  - canonical durable identity
  - operator-facing display alias

Alternative C:

- is the least PromiseGrid-complete because it keeps local ID minting as the
  durable identity and papers over collisions later

Result:

- A and B are both PromiseGrid-worthy; C is not.

### Scenario 5: debugging and operator cognition

Frank is using the browser and CLI while support staff diagnose cross-peer
history.

Alternative A:

- gives one stable durable ID everywhere
- but the IDs may become longer and less readable unless extra alias handling
  is added for operators

Alternative B:

- gives one stable durable key for the runtime and one short alias for humans
- makes logs, APIs, and peer exchange explicit without forcing every UI surface
  to become more verbose immediately

Alternative C:

- keeps short IDs visible
- but support staff must constantly reason about hidden translation tables

Result:

- B gives the best operator/runtime split.

### Scenario 6: implementation complexity and future debt

Mallory is not needed here; this is about repo maintenance cost.

Alternative A:

- requires the broadest one-time refactor
- but reduces long-term ambiguity once complete

Alternative B:

- requires a broad refactor too, but allows a staged embodiment transition
- gives the repo a durable canonical model without forcing one big visible-ID
  cutover

Alternative C:

- may appear cheaper immediately
- but imposes ongoing complexity in every import, lookup, projection, and
  embodiment contract

Result:

- B has the best implementation-to-future-debt ratio.

## Conclusions

Rejected alternative:

- Alternative C: import-side rewrite or translation without a canonical
  peer-stable entity key would keep the wrong durable identity model and create
  long-term ambiguity.

Surviving alternatives:

- Alternative A: canonical globally stable entity IDs replace current durable
  IDs and old short IDs become aliases
- Alternative B: canonical peer-stable entity keys become the durable runtime
  identity while current short IDs remain display aliases for embodiments

## Implications for TODOs and pending DIs

- `109` needs one locked answer on whether the runtime should expose the new
  canonical identity directly as the entity ID, or keep it separate from the
  displayed short ID.
- If Alternative A is chosen, embodiment routes and visible IDs likely change
  more aggressively in one pass.
- If Alternative B is chosen, the runtime still becomes PromiseGrid-complete at
  the durable identity layer, but browser/CLI/Neovim can transition more
  gradually.
- In either case, TODO `109` becomes a broad but coherent identity refactor
  across the six frozen families and their references.

## Decision status

Alternative A locked by `DI-loruk`: canonical durable entity IDs are derived
from the create-event envelope CID, and the old short IDs become aliases only.
