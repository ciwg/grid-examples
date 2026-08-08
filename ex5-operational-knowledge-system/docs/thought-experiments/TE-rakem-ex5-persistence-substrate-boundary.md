# ex5 persistence substrate boundary

TE ID: TE-rakem
## Status
decided

## Decision under test

What the next reusable PromiseGrid persistence substrate should be after
`promisegrid/records/` and `promisegrid/transport/` already exist.

This TE corresponds to TODO `mufek.1` / `mufek.2` / `mufek.3`.

## Assumptions

- `ex5` already has reusable durable record truth under
  `promisegrid/records/`.
- `ex5` already has reusable peer-exchange and relay-feed wire truth under
  `promisegrid/transport/`.
- The next extraction should be driven by shipped reuse evidence, not by a
  desire to make `promisegrid/` look larger.
- Runtime drafts, attachment materialization, projection rebuilds, and browser
  or terminal embodiment behavior are still example-local unless this TE shows
  otherwise.
- The relay binary is important reuse evidence because it already reuses much
  of the runtime's log and CAS persistence behavior without sharing the same
  operator workflows.

## Alternatives

### A. Keep all persistence in `service/`

Leave the whole `Store` type and its helper methods inside the example app for
now.

### B. Extract a minimal persistence substrate next

Extract only the reusable append-only event/family-log and CAS object
mechanics, plus authoritative frozen-envelope hydration, into a new
PromiseGrid substrate package. Keep attachment paths, draft manifests, and
other ex5-specific storage policy in `service/`.

### C. Extract the whole `Store` boundary next

Move nearly all of `service/persistence.go` into substrate, including draft
manifests, attachment materialization, local path layout, and other current
storage policy.

## Scenario analysis

### Scenario 1: normal runtime startup and replay

Alice starts one local runtime and expects durable history plus frozen-family
records to reload cleanly.

- A keeps startup behavior stable, but the reusable part of persistence stays
  hidden inside `service/`.
- B extracts the shipped reusable core that both startup and relay already
  depend on: append-only logs, CAS object storage, and authoritative envelope
  hydration.
- C also keeps startup working, but it promotes ex5's current draft and
  attachment layout into substrate before there is evidence another app wants
  those exact policies.

### Scenario 2: relay durability without operator workflows

Bob runs the dedicated relay and wants durable signed records and blobs, but
not local drafts, attachment file paths, or browser-oriented runtime state.

- A leaves runtime and relay sharing persistence behavior mostly by copy inside
  `service/`.
- B extracts the exact overlap the relay already proves: family logs, events,
  CAS objects, and authoritative frozen-envelope hydration.
- C forces relay and runtime to share more local operator policy than the
  relay actually needs.

### Scenario 3: imported evidence and local operator compatibility

Carol imports evidence from another node and needs it materialized into the
local compatibility attachment tree for current ex5 operator flows.

- A keeps that compatibility behavior local, which is safe but leaves the
  underlying portable CAS layer less reusable than it could be.
- B keeps attachment rematerialization local while still extracting the shared
  blob store underneath it.
- C incorrectly treats one current ex5 compatibility layout as part of the
  general PromiseGrid substrate.

### Scenario 4: shared drafts and embodiment recovery

Dave restarts the runtime and expects browser draft bodies to reload from CAS
through the local draft manifests.

- A keeps this entirely local and coupled to the app.
- B still keeps this local, because the current draft manifest format is tied
  to ex5 collaborative editing behavior rather than already-proven general
  PromiseGrid storage semantics.
- C over-generalizes a still-example-specific draft recovery policy.

### Scenario 5: future second app on the same substrate

Ellen later wants another PromiseGrid example app to reuse signed-record and
transport substrate pieces, but with a different local cache or draft policy.

- A offers no reusable persistence slice beyond records and transport.
- B offers a clean base: logs, CAS objects, and authoritative frozen-envelope
  hydration can be reused while the new app chooses its own draft and
  attachment policy.
- C makes the second app inherit ex5-local storage conventions or fork them
  back out again.

### Scenario 6: long-horizon migration and mixed versions

Frank upgrades one node while older runtimes still have manifest-only envelope
copies or old draft files.

- A preserves the current migration behavior, but the strongest reusable
  migration mechanics remain app-owned.
- B extracts the already-proven migration rule that belongs to durable record
  truth: CAS is authoritative for frozen envelopes, with one-time manifest
  backfill where needed.
- C ties that durable migration rule to unrelated draft and attachment
  conventions, making future substrate evolution harder.

## Conclusions

Rejected:

- Alternative A: too conservative now that runtime and relay already share a
  real persistence core in practice.
- Alternative C: too broad because it promotes ex5-specific storage policy
  into substrate before reuse evidence exists.

Surviving:

- Alternative B: extract a minimal persistence substrate around append-only
  logs, CAS objects, and authoritative frozen-envelope hydration while leaving
  drafts, attachment rematerialization, and local file-layout policy in
  `service/`.

## Implications for TODOs and pending DIs

- TODO `144` should lock to Alternative `B` if the goal is the strongest
  PromiseGrid-aligned next step.
- The remaining DF is naming and scoping the first persistence substrate
  package.
- TODO `145` should stay separate, because workflow substrate evidence should
  still be judged after persistence is cleaner.
- TODO `146` should stay last, because module packaging should follow proven
  substrate slices instead of leading them.
