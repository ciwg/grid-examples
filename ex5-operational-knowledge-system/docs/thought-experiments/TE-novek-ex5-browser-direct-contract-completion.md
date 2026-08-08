# Browser Direct-Contract Completion

TE ID: TE-novek
## Status
decided

## Decision under test

What the next browser direct-contract completion slice should be now that
browser writes and some reads already use typed runtime operations, but the
browser still falls back to generic `type:"request"` forwarding or plain HTTP
bootstrap for dashboard/catalog refresh, structured search, and live-state
bootstrap.

This TE corresponds to TODO `lunav.1` / `lunav.2` / `lunav.3`.

## Assumptions

- Alice uses the browser as a first-class embodiment over the shipped
  Chrome/Chromium native-messaging bridge.
- Bob wants the browser embodiment to carry runtime intent directly instead of
  re-entering the runtime as route-shaped HTTP requests wherever that mapping
  is already well understood.
- Carol still needs the current browser UI to load from the local HTTP shell
  and fetch `/api/meta` before the direct embodiment is ready.
- The current shipped direct browser contract already covers:
  - readiness with `runtime_ready`
  - inspect/search/problem-review read operations
  - main create/operate mutation operations
  - live drafting over the native-messaging path
- The remaining browser fallback surfaces are:
  - dashboard and collection refresh reads
  - structured search via `runSearch(...)`
  - live-state bootstrap via `GET /api/items/{id}/live`
  - any unmapped request path that still reaches `bridgeRPC({type:"request", ...})`

## Alternatives

### A. Leave the remaining browser request fallback in place

Accept the current split:

- typed operations for current inspect/search/problem-review and mutation paths
- generic `type:"request"` fallback for the rest

### B. Finish the remaining browser read slice next

Move the remaining browser read/bootstrap/query surfaces onto typed browser
operations, while leaving the local HTTP shell and `/api/meta` bootstrap alone.

That would cover:

- dashboard
- catalog lists for places/resources/responsibilities/items/runs
- structured search
- live draft bootstrap state

### C. Eliminate generic browser request forwarding entirely in one wave

Move both the remaining reads and any remaining generic request lane together
so the browser no longer emits `type:"request"` at all outside bootstrap HTTP.

## Scenario analysis

### Scenario 1: normal browser use

Alice loads the browser, reviews queues, opens a draft, searches, and edits.

#### A. Leave the remaining fallback in place

What it makes easier:
- no additional implementation risk

What it makes harder:
- the browser direct embodiment remains partly route-shaped
- PromiseGrid claims stay ahead of the browser’s actual semantic boundary

#### B. Finish the remaining browser read slice next

What it makes easier:
- the browser embodiment becomes much more runtime-shaped in daily use
- dashboard/search/catalog/live bootstrap become direct operations too
- keeps the remaining generic request lane very small and easier to reason about

What it makes harder:
- requires a broader typed read surface and matching tests

#### C. Eliminate generic browser request forwarding entirely in one wave

What it makes easier:
- strongest direct-contract purity

What it makes harder:
- wider migration surface at once
- higher risk of missing a less-obvious fallback path

### Scenario 2: PromiseGrid alignment

Steve evaluates whether the browser embodiment names runtime intent directly
instead of carrying adapter semantics inside a new transport.

#### A. Leave the remaining fallback in place

What it makes easier:
- preserves current stability

What it makes harder:
- leaves a meaningful amount of browser traffic still route-shaped

#### B. Finish the remaining browser read slice next

What it makes easier:
- removes the most visible remaining semantic impurity in day-to-day browser
  use
- keeps HTTP shell/bootstrap limited to where it is actually needed
- gives the browser one clearer direct-contract story without forcing a
  maximal one-wave rewrite

What it makes harder:
- still leaves a small generic request escape hatch for truly unmapped paths

#### C. Eliminate generic browser request forwarding entirely in one wave

What it makes easier:
- maximal embodiment purity

What it makes harder:
- raises the chance of an unnecessarily broad breakage wave compared with the
  concrete remaining browser surfaces already identified

### Scenario 3: maintenance cost

Bob wants the next move to be durable and testable.

#### A. Leave the remaining fallback in place

What it makes easier:
- no new work

What it makes harder:
- future reviews will keep rediscovering the same remaining seam

#### B. Finish the remaining browser read slice next

What it makes easier:
- one bounded follow-on closes the most obvious remaining browser gap
- tests can enumerate a concrete set of typed read operations

What it makes harder:
- requires some additional operation naming and projection plumbing

#### C. Eliminate generic browser request forwarding entirely in one wave

What it makes easier:
- one-time closure if done perfectly

What it makes harder:
- harder to scope, review, and validate deterministically

### Scenario 4: browser shell/bootstrap boundary

Carol still needs the browser shell and initial meta/bootstrap to come through
the local HTTP surface before direct embodiment readiness is proven.

#### A. Leave the remaining fallback in place

What it makes easier:
- no distinction work

What it makes harder:
- shell/bootstrap and runtime semantics remain less cleanly separated

#### B. Finish the remaining browser read slice next

What it makes easier:
- keeps the necessary HTTP shell/bootstrap boundary explicit
- moves runtime reads off route-shaped forwarding without pretending the browser
  shell itself must disappear

What it makes harder:
- requires careful wording so HTTP shell/bootstrap and typed runtime reads do
  not get conflated

#### C. Eliminate generic browser request forwarding entirely in one wave

What it makes easier:
- stronger eventual separation

What it makes harder:
- easy to overreach and accidentally couple shell/bootstrap cleanup to runtime
  semantic cleanup

## Conclusions

Rejected:

- Alternative A: too weak; it leaves the known browser semantic seam in place.
- Alternative C: stronger in theory, but too broad for the concrete remaining
  surfaces already identified.

Surviving:

- Alternative B: finish the remaining browser read slice next

Alternative B is the most PromiseGrid-aligned surviving path because it moves
the browser’s remaining day-to-day runtime reads onto typed operations while
keeping the necessary HTTP shell/bootstrap boundary explicit and bounded.

## Implications for open TODOs and pending DIs

- TODO `138` should lock a bounded typed-read completion wave, not defer the
  seam and not attempt a maximal browser rewrite in one step.
- Locked result: `138B.2`, meaning dashboard, catalog refresh, structured
  search, and live-state bootstrap move together in the first patch.
