# Ex5 Embodiment Contract Tightening Beyond HTTP

TE ID: `TE-ramut`
## Status
decided

## Decision under test

Whether the next embodiment-contract step for `ex5` should keep browser, CLI,
and Neovim routing through the local HTTP adapter, or introduce a more direct
runtime-facing contract now that peer exchange, CAS authority, and canonical
identity are much stronger.

Related TODO:

- `112` - `ex5-operational-knowledge-system/TODO/TODO-lurog-ex5-embodiment-contract-tightening-beyond-http.md`

## Assumptions

- `ex5` now ships eight PromiseGrid-native families plus origin-aware peer
  exchange and canonical create-envelope-CID durable IDs.
- The local HTTP API is still the only shipped embodiment adapter for browser,
  CLI, and Neovim.
- Browser and Neovim both reuse the shared live-draft surface, and CLI uses
  the same search/detail routes for operational review.
- The user wants the most PromiseGrid-complete practical next step, but does
  not want gratuitous stepping stones.
- `111` may still strengthen remaining local storage authority first.

## Threat and trust model

- Alice uses the browser, Bob uses the CLI, and Carol uses Neovim against the
  same runtime.
- Dave runs a mixed-version embodiment during migration.
- Mallory can tamper with local adapter requests if the boundary is poorly
  specified or duplicated across too many competing surfaces.
- The runtime should not create two conflicting embodiment contracts that drift
  apart semantically.

## Alternatives

### Alternative A

Keep HTTP as the sole embodiment adapter contract for now, and tighten only by
making the HTTP surface more explicitly describe the richer runtime beneath it.

Under this model:

- browser, CLI, and Neovim continue to route through `/api/*`
- capability metadata and docs remain the embodiment-tightening mechanism
- no direct runtime socket, file, or library contract is added yet

### Alternative B

Introduce one new local runtime-facing contract for terminal embodiments while
keeping the browser on HTTP.

Under this model:

- CLI and/or Neovim can speak to a narrower runtime-native boundary
- browser remains on HTTP for practical reasons
- ex5 ends up with two embodiment surfaces that must stay semantically aligned

### Alternative C

Start moving all embodiments away from HTTP toward a more direct runtime-facing
contract now.

Under this model:

- browser, CLI, and Neovim all begin migrating off the current HTTP adapter
- HTTP stops being the primary delivery surface
- the runtime contract becomes more explicit and less adapter-shaped

## Scope and systems affected

- `docs/thought-experiments/TE-ramut-ex5-embodiment-contract-tightening-beyond-http.md`
- `ex5-operational-knowledge-system/TODO/TODO-lurog-ex5-embodiment-contract-tightening-beyond-http.md`
- `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`
- `ex5-operational-knowledge-system/docs/http-api-guide.md`
- CLI and Neovim command implementations if a new surface is introduced
- browser transport behavior if HTTP stops being primary

## Scenario analysis

### Scenario 1: ordinary multi-embodiment use on one host

Alice writes in the browser, Bob reviews in the CLI, and Carol snapshots a
draft in Neovim.

Alternative A:

- all three keep using one adapter contract
- runtime semantics improve underneath them without splitting the embodiment
  surface
- operator documentation remains simpler

Alternative B:

- terminal embodiments gain a more direct contract
- browser still uses HTTP
- ex5 now has two embodiment surfaces to keep coherent

Alternative C:

- every embodiment changes at once
- this is the strongest break from the old model, but also the highest-risk
  migration with the broadest implementation blast radius

Result:

- A is the safest way to preserve one operational surface.
- B and C only help if the current adapter is actively blocking runtime work.

### Scenario 2: mixed-version migration

Dave upgrades the CLI first while Ellen keeps using an older browser build.

Alternative A:

- mixed-version migration is straightforward because the HTTP surface remains
  the common denominator

Alternative B:

- terminal and browser migrations can diverge
- compatibility testing burden increases because two embodiment surfaces must
  stay aligned over time

Alternative C:

- migration becomes all-or-nothing across the whole product

Result:

- A is clearly easiest to migrate.

### Scenario 3: PromiseGrid completeness pressure

Steve wants `ex5` to feel more honestly grid-native, not merely adapter-local.

Alternative A:

- the runtime can become more grid-native without pretending HTTP itself is the
  peer contract
- docs and capability metadata can remain explicit that HTTP is only the local
  embodiment adapter
- the durable/runtime semantics do the real PromiseGrid work underneath

Alternative B:

- terminal embodiments may feel more direct
- but the browser still anchors the product around HTTP
- this gives only partial embodiment tightening

Alternative C:

- this is the strongest embodiment-level break with HTTP
- but it only pays off if there is a clearly superior runtime-facing contract
  ready now

Result:

- A is still compatible with a fully on-grid runtime underneath.
- C is stronger only if the repo is ready for a large embodiment rewrite now.

### Scenario 4: trust boundary clarity

Mallory tries to exploit semantic drift between two embodiment surfaces.

Alternative A:

- one adapter contract means one translation layer to defend

Alternative B:

- CLI/Neovim and browser can drift if semantics are not maintained carefully

Alternative C:

- the new runtime-facing contract may be cleaner long-term
- but during migration the risk of drift is highest because all embodiment code
  paths are changing

Result:

- A minimizes semantic split risk.

### Scenario 5: long-horizon maintenance

Carol maintains `ex5` after the PromiseGrid runtime is richer than it was when
the HTTP surface was first introduced.

Alternative A:

- maintenance remains centered on one adapter plus one runtime
- if a future direct contract becomes necessary, it can be introduced against a
  more settled runtime/storage model

Alternative B:

- maintenance now spans two embodiment contracts indefinitely unless the HTTP
  surface is later retired

Alternative C:

- maintenance eventually converges on a more direct runtime surface
- but only after a disruptive migration with high near-term cost

Result:

- A is the strongest choice unless a concrete adapter limitation is already
  blocking current functionality.

## Conclusions

Rejected alternatives:

- Alternative B: it creates two embodiment contracts without delivering a fully
  unified replacement.
- Alternative C: it is too broad for the current migration stage and would
  force a large rewrite before the remaining local-only runtime state is fully
  settled.

Surviving alternative:

- Alternative A: keep HTTP as the sole embodiment adapter for now, and tighten
  the contract by making the richer runtime semantics more explicit beneath it.

## Implications and future work

- If Alternative A is chosen, TODO 112 should probably become a doc/capability
  tightening pass rather than a transport rewrite.
- If later work discovers a concrete embodiment blocked by HTTP, a new TE
  should test that narrower limitation instead of reopening the whole
  embodiment question generically.
