# Browser Live-Open Acknowledgement

TE ID: TE-morav
## Status
decided

## Decision under test

What should count as a successful browser live-open acknowledgement before
`ex5-operational-knowledge-system` reports the browser live transport as
connected?

This TE corresponds to TODO `vazek.1` / `vazek.2` / `vazek.3`.

## Assumptions

- Alice uses the browser embodiment over the shipped Chrome/Chromium native
  messaging path.
- The browser already uses one direct contract family for one-shot RPCs and
  live traffic.
- Browser live-open currently posts a `live-open` request and marks
  `socketConnected = true` before any reply is observed.
- The direct local runtime currently replies to a healthy `live-open` with a
  `live-state` payload; there is no separate explicit `live-opened` message.
- Silent fallback to an older HTTP browser path is not allowed.
- One-shot browser RPC timeout bounds are handled separately under TODO `139`.

## Alternatives

### A. Keep optimistic local connected state

Mark the browser live lane connected immediately when the page posts the
`live-open` message.

### B. First successful runtime live message counts as acknowledgement

Do not mark the browser live lane connected until the page receives the first
successful live response for that request, typically `live-state`, and possibly
`live-conflict` if the runtime reports a stale base right away.

### C. Invent a synthetic explicit `live-opened` acknowledgement

Add a new bridge/runtime message type whose only job is to say that the live
lane is open, then mark the browser connected on that synthetic acknowledgement
before any state payload is required.

## Scenario analysis

### Scenario 1: healthy live-open on a normal local runtime

Alice opens an item draft and the extension, native host, and runtime are all
healthy.

#### A. Keep optimistic local connected state

What it makes easier:
- no extra state transition logic

What it makes harder:
- connected status can briefly claim success before the runtime has actually
  answered

New obligations:
- none

#### B. First successful runtime live message counts as acknowledgement

What it makes easier:
- connected means the runtime actually answered
- no extra protocol shape is invented beyond the existing live contract

What it makes harder:
- the page must hold a short “opening” state until the first live reply lands

New obligations:
- reconnect logic must distinguish “opening” from “connected”

#### C. Invent a synthetic explicit `live-opened` acknowledgement

What it makes easier:
- gives a crisp explicit event for the page to wait on

What it makes harder:
- introduces a new live message kind whose meaning is transport-local rather
  than runtime-stateful
- splits acknowledgement truth from the actual first state payload

New obligations:
- document and test a new message kind across page, content script, native host,
  and runtime

### Scenario 2: live-open send succeeds but the reply stalls or disappears

Mallory is not required; a dropped bridge reply or stalled host can prevent the
first live response from reaching the page after `live-open` is posted.

#### A. Keep optimistic local connected state

What it makes easier:
- nothing

What it makes harder:
- the page can report itself connected even though no live state was ever
  acknowledged
- heartbeats and live updates can start from a false connected assumption

New obligations:
- repair the false-positive state later

#### B. First successful runtime live message counts as acknowledgement

What it makes easier:
- keeps connected state honest
- avoids sending follow-up live updates until the runtime has really answered

What it makes harder:
- the page needs an explicit “not yet connected” stance during the stall

New obligations:
- surface failure/retry from the still-opening state

#### C. Invent a synthetic explicit `live-opened` acknowledgement

What it makes easier:
- can acknowledge the lane before a full state payload

What it makes harder:
- if the synthetic ack arrives but state never does, the page is still ahead of
  reality

### Scenario 3: reconnect after disconnect

Alice loses the live lane and the page schedules a reconnect.

#### A. Keep optimistic local connected state

What it makes easier:
- reconnect path stays simple

What it makes harder:
- each reconnect can repeat the same false-positive connected window

#### B. First successful runtime live message counts as acknowledgement

What it makes easier:
- reconnect truth matches initial-open truth
- page state changes only after the runtime really resumes sending live data

What it makes harder:
- reconnect code must reset back to opening/unconnected until the first reply

#### C. Invent a synthetic explicit `live-opened` acknowledgement

What it makes easier:
- reconnect state can flip quickly on a small ack

What it makes harder:
- again separates transport-local success from actual runtime state delivery

### Scenario 4: long-horizon PromiseGrid alignment

Steve wants the browser embodiment to state transport truth as directly and
minimally as the rest of the PromiseGrid-aligned ex5 runtime.

#### A. Keep optimistic local connected state

What it makes easier:
- minimal churn

What it makes harder:
- keeps a known truth gap between local send and real runtime acknowledgement

#### B. First successful runtime live message counts as acknowledgement

What it makes easier:
- ties connection truth to the first real runtime response
- avoids inventing new top-level live semantics just to acknowledge transport
- keeps the browser embodiment on the same direct contract family it already
  uses

What it makes harder:
- requires the page to model an opening state explicitly

#### C. Invent a synthetic explicit `live-opened` acknowledgement

What it makes easier:
- makes acknowledgement explicit

What it makes harder:
- adds a new contract message whose meaning is mostly adapter bookkeeping

### Scenario 5: testing and maintenance

Dave wants deterministic tests that show healthy open, stalled open, and
disconnect/reconnect behavior clearly.

#### A. Keep optimistic local connected state

What it makes easier:
- no new tests

What it makes harder:
- preserves the current false-positive path

#### B. First successful runtime live message counts as acknowledgement

What it makes easier:
- tests can assert exactly when connected becomes true
- the existing live-state path already provides the natural success signal

What it makes harder:
- page and extension coverage must observe the pre-ack opening state

#### C. Invent a synthetic explicit `live-opened` acknowledgement

What it makes easier:
- tests can assert on a dedicated ack event

What it makes harder:
- all layers need new fixtures for a message that exists only for
  acknowledgement bookkeeping

## Conclusions

Rejected:
- `140A` keep optimistic local connected state
- `140C` invent a synthetic explicit `live-opened` acknowledgement

Surviving:
- `140B` first successful runtime live message counts as acknowledgement

`140B` is the most PromiseGrid-aligned alternative because it ties browser
connection truth to the first real runtime reply, avoids inventing a new
browser-only live message family, and keeps the direct contract anchored in
runtime state rather than local optimistic assumptions.

## Implications for open TODOs and pending DIs

- TODO `140` should lock `140B` and then decide whether only `live-state` or
  both `live-state` and `live-conflict` count as successful acknowledgement.
- The browser live docs should explicitly distinguish “opening” from
  “connected” once the DF is locked.
- TODO `139` remains separate and should not be folded into live-open
  acknowledgement semantics.

Locked result:
- `140B.2` -> `live-state` and `live-conflict` both count as successful live-open acknowledgement
