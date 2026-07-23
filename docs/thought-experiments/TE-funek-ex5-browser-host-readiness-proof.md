# Browser Host Readiness Proof

TE ID: TE-funek
## Status
decided

## Decision under test

How should `ex5-operational-knowledge-system` prove that the shipped
Chrome/Chromium direct browser embodiment is actually ready before the browser
marks itself usable?

This TE corresponds to TODO `fovek.1` / `fovek.2` / `fovek.3`.

## Assumptions

- Alice uses the browser as the main review and authoring surface.
- Bob uses the CLI and Carol uses Neovim against the same local runtime.
- The browser embodiment is already locked to Chrome/Chromium Manifest V3 plus
  a native-messaging host and must remain fail-closed when that direct contract
  is unavailable.
- The page content script, extension worker, native host, and local
  `embodiment.sock` path are distinct layers that can fail independently.
- The current browser startup only proves `/api/meta` plus page-bridge
  presence; it does not yet prove that the native host is reachable.
- Silent fallback to the older HTTP browser path is not allowed.

## Alternatives

### A. Page-bridge handshake only

Treat the current handshake as sufficient. If the content script answers, the
browser marks itself ready and discovers native-host failure only when the
first real RPC or live request runs.

### B. One-shot native-host readiness RPC at startup

Add a small startup probe that crosses the page bridge, extension worker, and
native host, then performs one cheap round-trip against the local runtime
socket before the browser marks itself ready.

### C. Full live-open readiness proof at startup

Require the browser to open a live-draft session at startup as the readiness
proof, so browser readiness is not declared until the most demanding direct
contract lane succeeds.

## Scenario analysis

### Scenario 1: normal startup on a healthy local runtime

Alice launches Chrome with the extension and native host installed correctly.

#### A. Page-bridge handshake only

What it makes easier:
- startup remains fast and simple
- no extra native-host probe message is needed

What it makes harder:
- readiness is overstated because only the content script is proven
- operators can still hit a delayed failure on the first real action

New obligations:
- keep explaining that “ready” only means the page bridge is present

#### B. One-shot native-host readiness RPC at startup

What it makes easier:
- startup truth matches the actual direct browser contract
- operators fail early, before they start work
- keeps the readiness probe smaller than a full live-draft session

What it makes harder:
- one more startup round-trip is required
- the system needs one agreed cheap probe operation

New obligations:
- define which runtime message counts as the readiness probe
- document the exact failure state when the native host is missing

#### C. Full live-open readiness proof at startup

What it makes easier:
- proves the strongest lane directly, including live traffic

What it makes harder:
- startup now depends on a heavier workflow-specific action
- browser readiness becomes tied to item/live-draft state, not just embodiment
  availability

New obligations:
- define which item, participant identity, and close behavior the startup proof
  uses

### Scenario 2: extension installed but native host missing or misregistered

Alice has the extension, but the native host path or registration is broken.

#### A. Page-bridge handshake only

What it makes easier:
- nothing new

What it makes harder:
- the browser claims readiness even though the direct contract is broken
- the first real browser action becomes the discovery point

New obligations:
- surface better delayed errors later

#### B. One-shot native-host readiness RPC at startup

What it makes easier:
- failure is detected immediately
- the browser can show one precise “direct browser embodiment unavailable”
  state before any workflow starts

What it makes harder:
- startup errors become more visible, which can feel harsher

New obligations:
- make the message actionable and specific

#### C. Full live-open readiness proof at startup

What it makes easier:
- also catches the failure immediately

What it makes harder:
- uses a much larger probe than needed just to discover that the native host is
  missing

### Scenario 3: runtime socket missing or stale

Bob or Carol may still have stale local configuration, or Alice starts Chrome
before the runtime is up.

#### A. Page-bridge handshake only

What it makes easier:
- startup stays independent of the runtime

What it makes harder:
- readiness no longer means the embodied runtime is actually reachable

#### B. One-shot native-host readiness RPC at startup

What it makes easier:
- the browser proves the full embodiment chain, including the local runtime
  socket
- the error surface is immediate and honest

What it makes harder:
- startup now depends on runtime availability, not just extension presence

#### C. Full live-open readiness proof at startup

What it makes easier:
- also proves runtime availability

What it makes harder:
- again pulls item/live state into what should be a transport-availability
  decision

### Scenario 4: long-horizon PromiseGrid alignment

Steve wants the browser embodiment story to be as explicit and honest as the
terminal direct-contract story.

#### A. Page-bridge handshake only

What it makes easier:
- minimal code churn

What it makes harder:
- keeps a truth gap between advertised readiness and actual embodiment
  availability

#### B. One-shot native-host readiness RPC at startup

What it makes easier:
- keeps readiness tied to the actual embodiment boundary
- mirrors the repo’s broader direction of explicit primary-contract truth
- avoids overfitting readiness to one workflow

What it makes harder:
- requires one more explicit typed probe across the direct contract

#### C. Full live-open readiness proof at startup

What it makes easier:
- proves the strongest lane

What it makes harder:
- mixes embodiment availability with collaboration semantics
- makes the readiness concept less general and less reusable

### Scenario 5: test coverage and maintenance

Dave wants deterministic tests that clearly show which layer failed.

#### A. Page-bridge handshake only

What it makes easier:
- no new coverage required

What it makes harder:
- the current tests still overstate what startup proves

#### B. One-shot native-host readiness RPC at startup

What it makes easier:
- allows focused tests for healthy host, missing host, and stale runtime socket
- keeps the readiness test smaller than the live-draft stack

What it makes harder:
- requires new fixtures at the content-script and native-host boundary

#### C. Full live-open readiness proof at startup

What it makes easier:
- one path covers both readiness and live semantics

What it makes harder:
- test setup becomes heavier and more fragile than the readiness question
  requires

## Conclusions

Rejected:

- Alternative A: it leaves the current truth gap in place; the browser can
  still report readiness while the native host is unavailable.
- Alternative C: it over-couples readiness to live-draft workflow behavior
  instead of proving the browser embodiment boundary cleanly.

Surviving:

- Alternative B: one-shot native-host readiness RPC at startup

Alternative B is the most PromiseGrid-aligned path. It proves the real browser
embodiment chain without reintroducing HTTP fallback and without making
workflow-specific live sessions the definition of readiness.

## Implications for open TODOs and pending DIs

- TODO `134` should lock one explicit readiness probe that crosses page bridge,
  extension worker, native host, and local runtime socket.
- The implementation still needs a DF choice on the exact probe shape so the
  probe remains typed and cheap instead of degenerating back into generic route
  tunneling.

The locked DF result is:

- `134B`
- `134B.1`
- one-shot readiness proof through typed `operation: "runtime_ready"`
