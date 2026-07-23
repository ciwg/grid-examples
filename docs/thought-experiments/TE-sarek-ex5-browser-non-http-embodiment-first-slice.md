# Browser Non-HTTP Embodiment First Slice

TE ID: TE-sarek
## Status
needs DF

## Decision under test

What should the first direct non-HTTP browser embodiment contract be for
`ex5-operational-knowledge-system`, now that CLI and Neovim already prefer the
local Unix-socket contract and the browser remains the main embodiment still
projecting through the local HTTP adapter?

This TE corresponds to TODO `rubek.1` / `rubek.2` / `rubek.3`.

## Assumptions

- Alice uses the browser as the main interactive review and authoring surface.
- Bob uses the CLI and Carol uses Neovim against the same local runtime.
- Dave operates the separate remote relay service, which is out of scope for
  the browser embodiment contract itself.
- The direct browser contract is still local to the same machine and runtime.
- Durable PromiseGrid families, relay feed, CAS, and the terminal local socket
  contract already exist and should not be destabilized by the first browser
  slice.
- The first browser slice should improve PromiseGrid alignment without forcing
  a simultaneous rewrite of every browser action and every fallback policy.

## Alternatives

### A. Browser-side local Unix-socket bridge over the existing browser app

Keep the browser UI, but replace its primary HTTP `fetch`/websocket traffic
with a browser-local bridge that speaks the same local Unix-socket runtime
contract used by terminal embodiments. The bridge can be exposed through a
browser-capable helper process or extension-local native messaging layer.

### B. Browser-side local runtime sidecar with a typed browser contract

Add a dedicated local browser sidecar that is not the HTTP adapter and not the
raw terminal socket contract. The browser speaks one typed local browser
contract to that sidecar; the sidecar then talks to the runtime through the
typed local socket/runtime layer.

### C. Keep browser HTTP but tighten it further

Do not move the browser off HTTP yet. Instead, refine the current adapter,
reduce fallback ambiguity, and continue to rely on `local_http` as the browser
primary contract.

## Scenario analysis

### Scenario 1: Normal local operation

Alice opens the browser, inspects items and runs, edits a live draft, approves
records, and uses search/pending/problem review throughout the day.

#### A. Browser-side local Unix-socket bridge

What it makes easier:
- keeps one shared direct runtime contract family across browser, CLI, and
  Neovim
- reduces the biggest remaining embodiment impurity in one meaningful step
- lets browser and terminal embodiments converge on the same typed runtime
  operation vocabulary over time

What it makes harder:
- browsers do not natively speak Unix sockets, so a helper/bridge is required
- deployment and startup become more complex than today's simple `http://127.0.0.1`
- browser packaging and local trust prompts become part of the embodiment story

New obligations:
- define how the browser locates and authenticates the local bridge
- define how live-draft streaming maps onto the direct contract
- define explicit browser compatibility behavior if the bridge is unavailable

#### B. Browser-side local runtime sidecar

What it makes easier:
- gives the browser a contract designed specifically for browser constraints
- can hide Unix-socket details and OS-specific bridge logic from the browser app
- may allow a more stable browser UX during migration

What it makes harder:
- creates a second non-HTTP contract instead of converging on the existing
  direct runtime contract
- introduces another layer whose semantics can drift from the runtime
- pushes more meaning into a browser-specific layer, which is less
  PromiseGrid-clean

New obligations:
- define and maintain a new browser-specific local contract
- test and document the bridge-to-runtime mapping permanently
- keep sidecar evolution aligned with the terminal contract

#### C. Keep browser HTTP

What it makes easier:
- no deployment disruption
- preserves the currently working browser path
- avoids browser/OS helper complexity

What it makes harder:
- leaves the biggest remaining embodiment impurity in place
- keeps browser semantics tied to adapter-shaped routes
- delays convergence toward a more runtime-native embodiment story

New obligations:
- continue explaining why browser is special
- accept that `ex5` remains meaningfully below the strongest PromiseGrid bar

### Scenario 2: Browser startup when local support is missing or misconfigured

Alice launches the browser, but the local helper or bridge is missing,
misconfigured, or stale.

#### A. Browser-side local Unix-socket bridge

What it makes easier:
- failure is explicit at the embodiment boundary
- browser can clearly report “direct local embodiment unavailable”

What it makes harder:
- first-run and failure UX need careful design
- a silent fallback to HTTP would weaken the new contract immediately

New obligations:
- decide whether browser compatibility HTTP remains explicit opt-in or is
  removed for the first slice
- surface actionable operator guidance when the bridge is unavailable

#### B. Browser-side local runtime sidecar

What it makes easier:
- a browser-specific sidecar can potentially hide more startup issues

What it makes harder:
- failures may become harder to diagnose because the sidecar can fail
  independently of the runtime
- more moving parts means more states to test

New obligations:
- define sidecar discovery, health, and restart behavior

#### C. Keep browser HTTP

What it makes easier:
- startup behavior remains simple

What it makes harder:
- does not resolve the PromiseGrid concern at all

### Scenario 3: Mixed-version migration

Bob and Carol already use the terminal socket contract. Alice upgrades to the
new browser embodiment while some installations still run the older browser
HTTP path.

#### A. Browser-side local Unix-socket bridge

What it makes easier:
- the underlying runtime contract can stay shared with terminal embodiments
- mixed browser versions can coexist if the HTTP adapter remains as explicit
  compatibility during rollout

What it makes harder:
- requires a transition period where both browser-local bridge and HTTP browser
  adapter exist

New obligations:
- publish capability metadata that says which browser embodiment contracts are
  available
- decide how long compatibility browser HTTP remains supported

#### B. Browser-side local runtime sidecar

What it makes easier:
- a sidecar can version independently from the runtime

What it makes harder:
- version skew grows because browser, sidecar, and runtime can all drift

New obligations:
- maintain compatibility matrix across three moving layers

#### C. Keep browser HTTP

What it makes easier:
- no migration

What it makes harder:
- no actual improvement

### Scenario 4: Long-horizon PromiseGrid evolution

Steve wants `ex5` to move from “strongly aligned example” toward the most
PromiseGrid-pure embodiment story practical in this repo.

#### A. Browser-side local Unix-socket bridge

What it makes easier:
- converges all local embodiments toward one runtime-native contract family
- keeps the adapter distinction clearer: HTTP becomes compatibility/browser
  legacy, not the core browser definition
- pairs naturally with the typed runtime-operation work from TODO `131`

What it makes harder:
- the browser still needs some helper technology because browsers do not speak
  Unix sockets directly

New obligations:
- define the browser bridge as a thin carriage layer, not a second semantic
  runtime

#### B. Browser-side local runtime sidecar

What it makes easier:
- may feel more product-like in the short term

What it makes harder:
- moves the repo toward a browser-special substrate instead of one shared local
  runtime contract

New obligations:
- justify why a browser-specific contract is not just a compatibility layer

#### C. Keep browser HTTP

What it makes easier:
- nothing new

What it makes harder:
- caps PromiseGrid alignment below the strongest available shape

### Scenario 5: Scale and operational complexity

The browser is used heavily for review/search/live drafting, and the repo must
remain testable and maintainable.

#### A. Browser-side local Unix-socket bridge

What it makes easier:
- fewer semantic layers long-term if the bridge stays thin
- lets the runtime contract be tested once and reused

What it makes harder:
- browser integration testing becomes more complex initially

New obligations:
- add browser smoke coverage that proves the bridge uses the typed runtime
  contract rather than HTTP adapter routes

#### B. Browser-side local runtime sidecar

What it makes easier:
- can tailor performance or buffering to browser needs

What it makes harder:
- more code, more tests, more drift risk

New obligations:
- maintain dedicated sidecar behavior and test harnesses

#### C. Keep browser HTTP

What it makes easier:
- no immediate new test harnesses

What it makes harder:
- leaves the strongest alignment improvement undone

## Rejected alternatives

- **C. Keep browser HTTP but tighten it further**
  - rejected because it does not actually satisfy TODO `132` and leaves the
    main remaining embodiment impurity intact

## Surviving alternatives

- **A. Browser-side local Unix-socket bridge over the existing browser app**
- **B. Browser-side local runtime sidecar with a typed browser contract**

## Conclusions

The strongest PromiseGrid-aligned direction is **A**.

Why:
- it converges the browser toward the same direct local runtime contract family
  already used by terminal embodiments
- it minimizes new semantic layers
- it builds naturally on TODO `131`, where the terminal socket contract has
  already started moving from route-shaped forwarding to typed runtime
  operations

Alternative **B** survives only as the more browser-comfortable staging option,
not the more PromiseGrid-pure one. It is viable if deployment constraints or
browser helper realities force an intermediate layer, but it creates another
contract that the repo would then have to keep aligned forever.

## Implications for open TODOs and pending DIs

- TODO `132` should lock between surviving alternatives **A** and **B**
- If `132A` is chosen, the next DF questions should cover:
  - the concrete browser bridge technology
  - whether browser compatibility HTTP remains during rollout
  - how browser live drafting maps onto the typed runtime-operation layer
- TODO `133` should remain deferred until `132` clarifies whether `ex5`
  extracts a broader reusable browser/runtime substrate or only a thin browser
  bridge
