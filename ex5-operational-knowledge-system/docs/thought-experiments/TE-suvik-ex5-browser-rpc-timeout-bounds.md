# Browser RPC Timeout Bounds

TE ID: TE-suvik
## Status
decided

## Decision under test

How should `ex5-operational-knowledge-system` bound browser one-shot
direct-contract RPC waits so lost extension/native-host replies do not leave
the browser UI hanging indefinitely?

This TE corresponds to TODO `talem.1` / `talem.2` / `talem.3`.

## Assumptions

- Alice uses the browser embodiment over the shipped Chrome/Chromium native
  messaging path.
- Bob and Carol continue using the CLI and Neovim direct contracts over the
  same local runtime.
- The browser is already fail-closed for unsupported browsers and for missing
  direct-contract readiness.
- Browser one-shot RPCs currently create a pending promise in `web/app.js` and
  depend on a later `rpc-response` or `error` message from the content script.
- The current implementation has no timeout or cleanup if the reply never
  returns.
- Silent fallback to the older HTTP browser path is not allowed.
- Live-draft connection acknowledgement is a separate question under TODO
  `140`; this TE is only about one-shot RPCs.

## Alternatives

### A. Keep one-shot RPC waits unbounded

Leave the current promise lifecycle unchanged. Browser UI actions wait until a
reply eventually returns or the page is reloaded.

### B. Page-owned per-RPC timeout with cleanup

Each one-shot browser RPC gets a bounded timer in `web/app.js`. If no reply
arrives before the bound, the pending request is removed, the promise rejects
with a direct-contract timeout error, and later stray replies are ignored.

### C. Extension-owned timeout only

Keep page promises unbounded, but add a timeout in the content script or
extension worker so the extension sends back an explicit error when the native
host path stalls.

## Scenario analysis

### Scenario 1: normal local operation on a healthy runtime

Alice opens the dashboard, searches, and writes a new item while the extension,
native host, and local runtime are all healthy.

#### A. Keep one-shot RPC waits unbounded

What it makes easier:
- no new timeout tuning or cleanup logic

What it makes harder:
- keeps browser action completion dependent on perfect reply delivery
- offers no bounded failure semantics if the local bridge regresses later

New obligations:
- none

#### B. Page-owned per-RPC timeout with cleanup

What it makes easier:
- the UI has a clear completion or failure bound for every one-shot action
- timeout ownership sits next to the pending promise map that actually needs
  cleanup

What it makes harder:
- requires choosing a bound that is strict enough to fail closed without
  spuriously rejecting healthy local work

New obligations:
- document the timeout as part of the direct browser contract
- test healthy completion versus timeout cleanup

#### C. Extension-owned timeout only

What it makes easier:
- keeps timeout logic closer to the native-messaging hop

What it makes harder:
- the page still has no self-owned cleanup if the content script or extension
  never answers at all
- the browser UI remains dependent on another layer to remember to respond

New obligations:
- coordinate timeout rules across page and extension layers later if the page
  still needs its own bound

### Scenario 2: dropped or blackholed reply

Mallory does not need to be malicious here; a bug, dropped message, or stalled
extension/native-host callback loses the reply after the browser posts the
request.

#### A. Keep one-shot RPC waits unbounded

What it makes easier:
- nothing

What it makes harder:
- the UI action hangs forever
- pending request state accumulates until reload

New obligations:
- operators must guess whether to retry, reload, or wait longer

#### B. Page-owned per-RPC timeout with cleanup

What it makes easier:
- the browser fails closed within a known bound
- pending request state is reclaimed deterministically
- late stray responses become harmless because the page has already forgotten
  the request

What it makes harder:
- timeout errors become visible more often when the bridge is unhealthy

New obligations:
- make the timeout error actionable and clearly tied to the direct browser
  embodiment

#### C. Extension-owned timeout only

What it makes easier:
- can recover some native-host stalls if the extension is still healthy

What it makes harder:
- cannot protect the page from a content-script or extension reply that never
  reaches it
- does not clean up the page's pending promise unless a later error is emitted

New obligations:
- add more page-level protection later anyway

### Scenario 3: concurrent browser RPCs

Alice opens multiple views quickly, causing several outstanding browser bridge
 requests at once.

#### A. Keep one-shot RPC waits unbounded

What it makes easier:
- no timer bookkeeping

What it makes harder:
- one lost reply leaves one pending entry stuck indefinitely
- concurrency makes it harder to notice which action is actually hung

#### B. Page-owned per-RPC timeout with cleanup

What it makes easier:
- each request owns its own lifecycle and deadline
- timeout cleanup stays keyed by request ID in the same structure that resolves
  replies today

What it makes harder:
- requires per-request timers instead of a single global watchdog

New obligations:
- make sure late replies for timed-out requests are ignored safely

#### C. Extension-owned timeout only

What it makes easier:
- timeout handling stays centralized in one extension layer

What it makes harder:
- still leaves the page's request table exposed if the extension path itself
  wedges

### Scenario 4: long-horizon PromiseGrid alignment

Steve wants the browser direct contract to be as explicit and bounded as the
terminal direct contracts, without silently demoting to compatibility lanes.

#### A. Keep one-shot RPC waits unbounded

What it makes easier:
- minimal code churn

What it makes harder:
- leaves the browser contract less honest than the surrounding fail-closed
  design
- makes success/failure semantics depend on eventual message luck instead of a
  declared contract

#### B. Page-owned per-RPC timeout with cleanup

What it makes easier:
- keeps the browser embodiment explicit: every one-shot contract either
  succeeds or fails within a known bound
- places timeout ownership at the layer that owns user-visible action promises
- avoids reviving the old HTTP adapter as a hidden escape hatch

What it makes harder:
- requires choosing and documenting one explicit bound

#### C. Extension-owned timeout only

What it makes easier:
- makes the extension more protective on the native-host hop

What it makes harder:
- keeps the browser page as a partially passive client instead of owning its
  own direct-contract truth

### Scenario 5: testing and maintenance

Dave wants deterministic tests that exercise healthy reply, explicit error, and
timeout cleanup behavior.

#### A. Keep one-shot RPC waits unbounded

What it makes easier:
- no new tests

What it makes harder:
- the known hang path remains untested and unbounded

#### B. Page-owned per-RPC timeout with cleanup

What it makes easier:
- smoke tests can simulate no reply and assert a bounded browser error
- script-level tests can verify that late replies are ignored after timeout

What it makes harder:
- tests need a controllable timeout horizon

#### C. Extension-owned timeout only

What it makes easier:
- extension tests can simulate native-host stalls directly

What it makes harder:
- still leaves a page-level gap in deterministic coverage

## Conclusions

Rejected:
- `139A` keep one-shot RPC waits unbounded
- `139C` extension-owned timeout only

Surviving:
- `139B` page-owned per-RPC timeout with cleanup

`139B` is the most PromiseGrid-aligned alternative because it keeps timeout
truth at the page layer that owns the user-visible promise, fails closed
without demoting to HTTP, and turns lost replies into bounded explicit errors
instead of indefinite hangs.

## Implications for open TODOs and pending DIs

- TODO `139` should lock `139B` and then choose one explicit timeout bound for
  one-shot browser RPCs.
- TODO `140` remains separate; live-open acknowledgement should not be folded
  into the one-shot RPC timeout rule.
- The browser direct contract docs should state the one-shot timeout bound once
  DF is locked.

Locked result:
- `139B.2` -> 1000ms page-owned per-RPC timeout with cleanup
