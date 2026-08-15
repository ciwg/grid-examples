# Ex2 presence-lifecycle refresh ownership

TE ID: TE-zubup

## Status

decided

## Decision under test

Where should Ex2 recompute the age of a peer's last accepted
`live-awareness` observation so the browser and Neovim roster moves from
`live` to `stale` to `offline` and finally disappears without new traffic?

This TE completes TODO `mivor`, especially `mivor.1` through `mivor.4`.

## Assumptions and trust model

- Ex2 is a decentralized, multi-relay example. Each relay verifies and retains
  its own signed profile traffic; no relay is a global membership authority.
- `live-awareness` is a repo-local-draft, pCID-selected, latest-state presence
  profile. It communicates an observer's current presentation state, not
  durable membership, a promise of liveness, or authority over another peer.
- The locked lifecycle is: live for 0–1 minute, stale for 1–5 minutes,
  offline for 5–15 minutes, and removed after 15 minutes in the normal
  profile. The existing demo profile keeps longer windows for recordings.
- Browser and Neovim already receive `last_seen_at` with the relay's accepted
  awareness projection. Existing code can calculate the states, but currently
  recalculates only when another event causes rendering.
- Mallory may stop sending traffic, delay a valid message, operate an older
  implementation, or make a relay unavailable. No local timer may turn that
  absence into a claim about Mallory's intent.

## Alternatives

### A. Embodiment-local periodic refresh

The browser and Neovim each retain the latest accepted awareness projection
and periodically recompute only their local rendering state. The relay keeps
its current latest-state projection and does not delete or publish expiry
messages. The refresh cadence is bounded and local; it exists solely to make
the already-locked display windows visible without fresh traffic.

### B. Relay-side expiry and removal

Each relay runs an expiry process that removes peers from its awareness
projection and publishes the resulting smaller roster to connected clients.

### C. Hybrid relay expiry plus embodiment-local refresh

Relays remove stale awareness from their projection, while clients also use a
timer to age any last received projection between relay responses.

## Scenario analysis

### Normal operation

Alice and Bob exchange awareness while editing the same document. Bob closes
his editor. Under A, Alice's browser and Carol's Neovim each move Bob through
the locked display states using their own clocks; neither emits protocol
traffic. Under B, their view changes only when their chosen relay runs expiry
and sends a new projection. C produces the same visible result but duplicates
the authority for deciding when a record is no longer shown.

A makes the existing observer-local meaning explicit. B makes a relay's timer
look like a stronger membership decision than it is. C creates two clocks and
two sources of display behavior without adding protocol value.

### Relay outage, delayed traffic, and incomplete writes

Alice's relay becomes unreachable after it supplied Bob's last accepted
awareness state. Under A, Carol's local view still ages that received fact;
it does not assert Alice left or revoked anything. Under B, the UI freezes
until the relay returns. Under C, the local timer avoids the freeze, which
shows that the relay expiry half is unnecessary for this display purpose.

If a relay crashes before retaining a new update, all alternatives retain only
the last accepted observation. A states that limit honestly. B and C risk
making a relay-local deletion look like a durable protocol event even though
no signed expiry fact exists.

### Concurrent relays and mixed versions

Alice's relay sees Bob at 12:00:00 while Carol's relay sees Bob at 12:00:30.
Under A, each embodiment renders the age of the observation it actually has;
temporary disagreement is expected and visible only as local presentation.
B can make relay differences look like contradictory peer membership claims.
C preserves that disagreement while adding timer duplication.

An older embodiment that does not run a refresh continues to show its last
view. A does not make that older rendering a protocol incompatibility: the
wire record and its pCID remain unchanged. B requires all relays to change
for a UI lifecycle improvement; C inherits both upgrade obligations.

### Long-horizon evolution and scale

At small scale, A costs one bounded local timer per open embodiment and no
extra network messages. At larger scale, it still works on an in-memory latest
projection and does not create persistence churn. B adds expiry jobs, relay
state mutation, reconnection semantics, and timing consistency questions for
every relay. C inherits both costs.

If future work introduces durable activity, comments, versions, `last viewed`,
or `last edited`, A keeps those records separate from the live roster. B and C
make it easier to accidentally reuse an expiry projection as historical truth.

## Conclusions

Alternative A is the PromiseGrid-aligned survivor. It treats presence aging as
an embodiment-local interpretation of the last local observation, preserves
the signed `live-awareness` contract unchanged, and does not elevate a relay
into membership authority.

Alternatives B and C are rejected for this slice. They add relay-side
membership-like behavior and duplicated clocks without improving the
pCID-defined awareness contract or evidence boundary.

## Required DF questions

1. Should both embodiments refresh locally on a bounded periodic timer and
   immediately after awareness events (A), or should an alternate design be
   reconsidered?
2. Should the timer use the next lifecycle boundary for its next wake-up
   (recommended) or a fixed short polling interval?
3. Should deterministic tests inject a clock/scheduler seam (recommended) or
   test only the pure classification function with fixed timestamps?
4. Which permanent source/test/doc paths should own the browser, Neovim, and
   regression coverage changes?

## Decision status

Locked by DI-dizut and DI-dazin: each embodiment refreshes its own current
awareness presentation at exact local lifecycle boundaries. The browser owns a
pure, deterministically tested scheduler module; Neovim uses the same local
next-boundary rule. Relays retain and distribute their latest accepted
awareness observations without expiry or membership authority.

## Implications and future work

If locked, this work changes only local presentation refresh and deterministic
coverage. It does not add a new protocol action, pCID, durable event, relay
expiry mechanism, or historical-collaboration surface. `mivor.5` remains a
separate future design task.
