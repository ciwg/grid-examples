## Title

ex5 websocket collaboration transport scope

## TE ID

TE-tivok

## Status

decided

## Decision under test

How `ex5-operational-knowledge-system` should close TODO `005` and move from
HTTP-polled live drafting to a more PromiseGrid-aligned websocket collaboration
transport without reopening the whole runtime design.

## Related TODO

- `005` `ex5-operational-knowledge-system/TODO/TODO-masad-ex5-websocket-collaboration-transport.md`

## Assumptions

- `ex5` already has one shared local runtime, one local HTTP adapter, eight
  frozen signed families, origin-aware peer exchange, CAS-backed durable family
  replay, and CAS-backed shared draft body reload.
- The current missing gap is transport for shared live drafting and presence,
  not durable-history semantics.
- Websocket is carriage, not protocol meaning. Durable PromiseGrid semantics
  stay in the frozen family records and existing peer-exchange layer.
- The current browser and Neovim live-draft paths both hit
  `GET/POST /api/items/{id}/live`.
- CLI is not part of the live drafting path.
- This TE is about local embodiment transport for collaborative drafting, not a
  full multi-relay network design.

## Threat / trust model

- Alice and Bob are cooperative operators editing the same item through browser
  and Neovim embodiments.
- Carol upgrades one client while Bob stays on an older polling client.
- Mallory is not assumed to have remote network access in this slice; the
  transport remains within the current local-adapter trust envelope unless a
  later TE expands it.
- The runtime must preserve optimistic conflict detection and must not let
  websocket carriage silently bypass durable revision boundaries.

## Alternatives

### Alternative A

Browser-only websocket transport.

- Browser switches live drafting and presence to websocket.
- Neovim keeps polling and posting over the existing HTTP live-draft route.
- HTTP stays primary for Neovim and as browser fallback.

### Alternative B

Shared websocket live transport for browser and Neovim.

- Browser and Neovim both prefer websocket for shared draft state and presence.
- Existing HTTP live routes remain as compatibility fallback and for simple
  non-streaming state fetches.
- Durable revisions, approvals, runs, and peer-exchange stay on their current
  routes.

### Alternative C

Broader relay/network-first websocket push before local embodiment transport
cleanup.

- Treat TODO `005` as the start of richer relay-visible network behavior.
- Expand beyond local live drafting into a larger websocket/relay slice now.

## Scenario analysis

### S1. Normal local collaboration: Alice in browser, Bob in Neovim

Alternative A:

- Alice gets lower-latency live drafting and presence.
- Bob still sees state on the polling cadence.
- Mixed transport is workable, but the shared-draft experience stays visibly
  uneven across embodiments.
- PromiseGrid alignment improves only partially because the two active authoring
  embodiments still do not share one live transport story.

Alternative B:

- Alice and Bob share one live transport story for the same working draft.
- Presence, body updates, and conflict visibility converge faster and more
  symmetrically.
- HTTP remains available for fallback and for the broader adapter surface.
- This is the cleanest step if the repo wants live collaboration to stop being
  “browser websocket someday” and become a real shared capability.

Alternative C:

- Could lead to a bigger end-state, but it mixes local embodiment transport with
  broader relay/network concerns before the local live path is simplified.
- Scope grows sharply and risks delaying a clear ex5 improvement.

### S2. Failure and restart: websocket drops or one client cannot upgrade

Alternative A:

- Browser needs a polling fallback anyway.
- Neovim remains untouched and therefore stable, but the codebase keeps two
  first-class live transport stories for longer.

Alternative B:

- Both browser and Neovim can attempt websocket first, then fall back to the
  existing HTTP live route when websocket is unavailable.
- This keeps the recovery story uniform and avoids making websocket availability
  embodiment-specific.
- The runtime already has optimistic live-draft conflict handling, so fallback
  still lands on the same shared draft semantics.

Alternative C:

- Failure handling becomes entangled with relay/network behavior, not just local
  transport fallback.
- More chances to leave the repo half-migrated.

### S3. Mixed-version rollout: Carol upgrades, Bob does not

Alternative A:

- Browser can upgrade independently.
- Neovim remains compatible because nothing changes there.
- This is the lowest-risk rollout, but it entrenches asymmetry.

Alternative B:

- Mixed-version operation remains possible if websocket is additive and the HTTP
  live route stays intact.
- Older browser or Neovim clients can keep polling.
- New clients can prefer websocket and still talk to the same live draft model.
- This adds more implementation work than A, but it does not require a flag day.

Alternative C:

- Mixed-version semantics become harder because relay/network expectations would
  also be moving.

### S4. Long-horizon PromiseGrid alignment

Alternative A:

- Leaves `ex5` with one embodiment on websocket and one on polling for the same
  shared-draft capability.
- That still reads like a stepping stone.

Alternative B:

- Most directly closes the current alignment gap identified in review:
  collaborative drafting gets a real websocket transport while durable
  PromiseGrid semantics stay where they already belong.
- Keeps websocket as carriage, not meaning.
- Fits the current architecture note that HTTP is the sole embodiment adapter,
  because websocket can still be an adapter path under the same local server.

Alternative C:

- Could eventually be more ambitious, but it is not the cleanest next move from
  the current repo state because it reopens unresolved network scope.

### S5. Operational complexity and code ownership

Alternative A:

- Smallest diff.
- Lowest testing burden.
- But it creates long-lived split transport ownership between browser and
  Neovim.

Alternative B:

- Moderate diff.
- Requires shared websocket server support and client transport code in both
  browser and Neovim.
- The complexity is still bounded because the live-draft payloads already
  exist; only carriage changes.

Alternative C:

- Highest complexity.
- Requires design decisions that are outside TODO `005`’s current phrasing and
  beyond the review finding we are trying to close.

## Conclusions

Rejected:

- Alternative C. It is too broad for the current gap and would mix local
  collaboration transport with a larger relay-network expansion before the repo
  proves the narrower websocket slice cleanly.

Surviving:

- Alternative A: browser-only websocket transport
- Alternative B: shared websocket transport for browser and Neovim, with HTTP
  fallback preserved

Recommendation:

- Alternative B

Why:

- It is the clearest PromiseGrid-aligned answer to the current review finding.
- It improves both active authoring embodiments instead of only the browser.
- It keeps websocket as carriage over the existing live-draft semantics instead
  of pretending the runtime needs a new durable protocol family for drafts.
- It can still preserve HTTP fallback for mixed-version and failure recovery.

## Decision status

Locked to Alternative B by `DI-noruv`.

## Implications for open TODOs and pending DIs

- Locking Alternative B would let TODO `005` become an implementation task
  instead of an open direction placeholder.
- If Alternative B is chosen, no new durable PromiseGrid family is required for
  shared drafts in this slice; draft semantics remain local shared working
  state.
- A later TE would still be needed before claiming richer relay-network
  behavior beyond the local-adapter live transport.
