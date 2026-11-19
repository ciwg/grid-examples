# Ex3 headless transport-contract proof

TE ID: TE-zolik

## Status

decided

## Decision status

Alternative A was selected by jj@thesalleys.com (JJ) on 2026-08-11 and locked
by DI-gofut; DI-raron locks the test-only observation name.

## Decision under test

How should Ex3's headless browser proof establish its live-transport claim
when a browser may also need the explicitly bounded HTTP recovery read for a
blank or stale snapshot?

## Assumptions and trust model

- The live-document and live-awareness pCID-selected payload semantics remain
  unchanged by WebSocket, HTTP, test probes, or browser rendering. WebSocket
  is carriage, not a fifth protocol profile. Source: DI-vubih; DI-hadil.
- The browser normally uses WebSocket for live sync and awareness. When a
  browser opens blank despite relay history, it may make one bounded HTTP
  sync-feed recovery read. Source: DI-sulor; DI-vobek.
- Alice runs the isolated test relay and browser. Bob has already written a
  durable Automerge change. Mallory may mistake one rendered UI string, an
  HTTP fetch, or a passing local test for a statement about protocol meaning
  or global network behavior.
- The proof must show behavior without adding a browser-only protocol,
  changing signed record bytes, or treating test instrumentation as runtime
  authority.

## Alternatives

### A. Record transport events and recovery events separately

The headless page reports whether the sync WebSocket reached its ready state
and whether the test-injected blank snapshot caused the one allowed HTTP
recovery fetch. The test asserts those facts separately, then retains the
visible transport summary as secondary UI evidence.

### B. Assert only the rendered transport label

Continue requiring the page text `browser sync: websocket` after startup.

### C. Assert only restored document content

Remove the transport observation and prove only that Bob's text appears in
Alice's headless page.

## Scenario analysis

### Normal late join

With A, the proof distinguishes the successful WebSocket ready frame from the
absence of recovery. It proves the documented preferred carriage path and
detects an unexpected HTTP recovery call. B has less diagnostic value: the
label can be sampled before or after another asynchronous state transition.
C proves convergence but cannot establish the published live-transport claim.

### Injected stale or blank snapshot

With A, the proof requires both facts: WebSocket readiness remains observable
and exactly the bounded recovery read restores Bob's history from offset zero.
This matches the existing implementation boundary: recovery is a repair of
local startup state, not a replacement wire protocol. B incorrectly makes one
late UI label carry both claims. C can conceal an unbounded or accidental
polling path as long as text eventually appears.

### WebSocket failure or incomplete write

With A, a missing ready event, an unexpected recovery request, or a recovery
request with an incorrect offset yields a targeted failure artifact. The test
does not infer that a WebSocket frame is durable evidence; it only records
which local carriage mechanism was observed. B yields a label mismatch with
no cause. C can pass when fallback behavior masks a live-transport failure.

### Concurrent actors and mixed browser timing

With A, Alice's probe reports ordered local events while Bob's pre-existing
change remains the content assertion. The ordering explains headless timing
variation without weakening the requirement for a ready WebSocket. B is
sensitive to when the DOM serializer samples a status field. C gives no
carriage evidence at all.

### Future transport evolution

With A, a future, separately specified carriage option can add an explicit
event type and test branch without rewriting payload semantics or changing the
meaning of existing records. B encourages an ever-more-complex display string.
C makes transport changes invisible to the proof.

### Scale and maintenance

A adds a small test-only probe contract but keeps the runtime contract
observable and failure output actionable. B is shorter but creates recurring
timing ambiguity. C is smallest but insufficient for Ex3's current documented
WebSocket claim.

## Conclusions

B is rejected because the UI label is an asynchronous presentation of state,
not direct evidence of a WebSocket-ready event or a recovery fetch. C is
rejected because restored content alone does not verify the stated preferred
carriage behavior.

A survives and was selected: the headless proof will capture WebSocket
readiness and bounded HTTP recovery as separate local observations, retain the
DOM label as secondary UI evidence, and leave pCID-selected payload semantics
unchanged.

## Locked decisions requiring DI

1. Add a test-only browser probe for sync WebSocket readiness and recovery
   fetch observation.
2. Assert the events independently in the normal and poisoned-snapshot tests.
3. Retain the visible transport label assertion only as UI evidence.
4. Do not add fallback polling, change live-document payload semantics, or
   represent the probe as a network-wide truth claim.

## Implications for open work

- The implementation belongs with Ex3's completed WebSocket and browser
  startup evidence, and needs a new DI in the relevant Ex3 TODO before code
  changes.
- `docs/testing.md` should describe the strengthened headless evidence once
  the implementation passes.
