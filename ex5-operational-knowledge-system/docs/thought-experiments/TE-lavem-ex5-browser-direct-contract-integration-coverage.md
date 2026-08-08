# Browser Direct-Contract Integration Coverage

TE ID: TE-lavem
## Status
decided

## Decision under test

What the next strongest deterministic browser direct-contract integration
coverage should be now that `ex5` ships a Chrome/Chromium MV3 extension,
content-script page bridge, native-messaging browser host, and local runtime
socket, but the main browser smoke tests still rely on `withMockBrowserBridge`
instead of the real extension scripts.

This TE corresponds to TODO `pavur.1` / `pavur.2` / `pavur.3`.

## Assumptions

- Alice loads `web/app.js` in Chrome or Chromium through the shipped extension.
- Bob maintains the Go native host and wants deterministic tests that can run
  in CI without machine-local browser registration.
- Carol reads the test/docs boundary to understand which browser layers are
  real and which are still simulated.
- The current test stack already covers:
  - manifest and native-host manifest shape in `chrome-extension/assets_test.go`
  - browser-host framing and typed forwarding in `service/browser_host_test.go`
  - browser UI behavior with a synthetic in-page bridge in
    `web/browser_smoke_test.go`
- The remaining gap is the real extension/native-host boundary between
  `content.js`, `background.js`, and the native host contract, especially for
  fail-closed startup and direct-contract RPC/live forwarding.

## Alternatives

### A. Keep the current split

Leave the current coverage as-is:

- asset tests for manifest shape
- Go tests for browser host framing
- headless UI smoke tests with `withMockBrowserBridge(...)`

### B. Add deterministic extension/native-host contract tests

Add explicit deterministic tests around the real shipped boundary by loading
the actual `background.js` and `content.js` logic under a stubbed Chrome API
environment, while pairing that with the existing Go browser-host tests.

This would cover:

- handshake/readiness failure and success through the real content/background
  message flow
- one-shot RPC forwarding and error propagation
- live-port forwarding/disconnect semantics
- explicit documentation of what is still mocked in the browser smoke tests

### C. Add a full Chrome/native-host end-to-end harness

Drive real Chrome/Chromium with the shipped extension loaded and a registered
native-messaging host, then run browser tests through the actual extension
packaging and native-host registration path.

## Scenario analysis

### Scenario 1: normal browser startup

Alice launches the browser embodiment on a machine with the extension
installed and the native host available.

#### A. Keep the current split

What it makes easier:
- no new harness work

What it makes harder:
- startup coverage still depends on the mock bridge rather than the shipped
  content/background path
- the real handshake chain stays under-tested

New obligations:
- accept that readiness regressions can slip past the current smoke suite

#### B. Add deterministic extension/native-host contract tests

What it makes easier:
- verifies the real bridge message choreography without requiring external
  browser installation state
- catches content/background regressions close to where they occur

What it makes harder:
- requires a small script-level test harness for the Chrome APIs

New obligations:
- clearly document which browser behaviors are still UI-smoke mocked

#### C. Add a full Chrome/native-host end-to-end harness

What it makes easier:
- closest test to the exact shipped operator path

What it makes harder:
- difficult registration/setup lifecycle
- much more CI and local-environment sensitivity
- more nondeterministic failure modes unrelated to runtime semantics

### Scenario 2: missing or misregistered native host

Bob breaks native-host registration or the host disappears from the local
system.

#### A. Keep the current split

What it makes easier:
- no new complexity

What it makes harder:
- the current browser smoke suite still cannot prove the failure propagates
  across the real extension boundary

#### B. Add deterministic extension/native-host contract tests

What it makes easier:
- can directly inject `chrome.runtime.lastError`, missing responses, or port
  disconnects and assert fail-closed behavior
- keeps those failure tests precise and repeatable

What it makes harder:
- requires deliberate mock objects for the Chrome runtime APIs

#### C. Add a full Chrome/native-host end-to-end harness

What it makes easier:
- can test real installation failures

What it makes harder:
- the distinction between product regression and host-registration environment
  failure becomes noisy

### Scenario 3: live-draft transport integrity

Carol needs confidence that browser live traffic still rides the direct
native-messaging path and fails honestly on disconnect.

#### A. Keep the current split

What it makes easier:
- no additional work

What it makes harder:
- live-open/live-update/live-close remain mostly trusted by structure rather
  than directly exercised at the extension boundary

#### B. Add deterministic extension/native-host contract tests

What it makes easier:
- directly covers `chrome.runtime.connect`, background port bridging, native
  disconnect propagation, and page error delivery

What it makes harder:
- modest test harness work

#### C. Add a full Chrome/native-host end-to-end harness

What it makes easier:
- most realistic live-path test

What it makes harder:
- expensive and brittle for a first pass

### Scenario 4: PromiseGrid alignment

Steve wants the test story to validate the real direct browser embodiment
without turning coverage into a deployment-lab project.

#### A. Keep the current split

What it makes easier:
- no new harness

What it makes harder:
- leaves the most PromiseGrid-important browser boundary still represented
  mainly by a synthetic page bridge

#### B. Add deterministic extension/native-host contract tests

What it makes easier:
- tests the real embodiment boundary where PromiseGrid alignment matters:
  page bridge -> extension -> native host -> runtime socket
- preserves deterministic, reviewable semantics
- keeps HTTP-era mock behavior from overstating what is truly covered

What it makes harder:
- requires a small JS harness for the Chrome APIs

#### C. Add a full Chrome/native-host end-to-end harness

What it makes easier:
- maximum operational realism

What it makes harder:
- over-couples PromiseGrid contract validation to environment setup and browser
  packaging details

### Scenario 5: maintenance cost

Dave wants a meaningful improvement that the repo can keep passing.

#### A. Keep the current split

What it makes easier:
- lowest immediate cost

What it makes harder:
- leaves the known review finding unresolved

#### B. Add deterministic extension/native-host contract tests

What it makes easier:
- good signal-to-maintenance ratio
- scales to later browser direct-contract follow-ons

What it makes harder:
- one new JS test surface to maintain

#### C. Add a full Chrome/native-host end-to-end harness

What it makes easier:
- broadest realism if maintained well

What it makes harder:
- highest long-term maintenance and environment burden

## Conclusions

Rejected:

- Alternative A: too weak; it leaves the real browser boundary under-tested.
- Alternative C: too heavy for the next deterministic coverage step.

Surviving:

- Alternative B: add deterministic extension/native-host contract tests

Alternative B is the most PromiseGrid-aligned surviving path because it tests
the actual shipped browser embodiment boundary directly, while keeping the
coverage deterministic and reviewable instead of depending on fragile external
browser-registration setup.

## Implications for open TODOs and pending DIs

- TODO `137` should lock a deterministic extension/native-host contract test
  wave, not a no-op and not a full deployment-lab harness.
- Locked result: `137B.2`, meaning the first deterministic contract test wave
  covers startup/readiness, one-shot RPC forwarding, and live-port
  forwarding/disconnect behavior together.
