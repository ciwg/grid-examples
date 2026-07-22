# Ex5 PromiseGrid Wire Slice Decision

TE ID: `TE-lafiz`
## Status
decided

## Decision under test

Whether `ex5-operational-knowledge-system` should open a real PromiseGrid
runtime/wire implementation slice now, or whether it should first stay on the
current local-runtime layer and freeze one narrow protocol contract before any
signed-envelope or relay-visible work begins.

Related TODO:

- `093` - `ex5-operational-knowledge-system/TODO/TODO-ragup-ex5-promisegrid-wire-slice-decision.md`

## Assumptions

- `ex5` already ships inside the PromiseGrid example set and should follow the
  PromiseGrid dev guide rather than inventing its own app-contract rules.
- The current shipped ex5 runtime already has meaningful durable behavior:
  append-only events, local projections, browser/CLI/Neovim embodiments, and a
  local HTTP adapter.
- The current shipped ex5 runtime does not yet implement frozen `pCID`-selected
  runtime behavior, signed `grid([42(pCID), ...])` envelopes, or relay-visible
  peer exchange.
- The current ex5 spec prose is still reference/spec guidance rather than a
  frozen spec with a published `pCID` and implementation promise claim.
- The PromiseGrid dev guide currently says app work should be spec-first:
  choose one explicit protocol spec, use its `pCID` when frozen, and publish
  implementation promise claims rather than treating branch paths or local
  adapters as the contract.

## Alternatives

### Alternative A

Do not open a real wire/runtime slice yet. Keep ex5 on the current local
runtime layer and instead make the next PromiseGrid step a narrow protocol
freeze/claim slice.

### Alternative B

Open a narrow real PromiseGrid slice now:

- freeze one ex5 protocol family immediately
- thread its `pCID` through runtime behavior
- start signed envelope creation/verification locally
- keep relay exchange out of scope for the first wire slice

### Alternative C

Open the larger PromiseGrid runtime slice now:

- freeze protocol identity
- implement signed envelopes
- add relay-visible peer exchange or transport-facing exchange machinery in the
  same program

## Scope and systems affected

- `ex5` technical docs and implementation claims
- ex5 runtime storage model
- ex5 event history representation
- local HTTP adapter semantics
- browser/CLI/Neovim embodiment contract boundaries
- future transport, signing, verification, and conformance claims

## Scenario analysis

### Scenario 1: normal operator workflow on the shipped local runtime

Alice runs the local ex5 server, drafts or revises a procedure, records runs,
captures evidence, approves work, and reviews later history from browser, CLI,
and Neovim.

Alternative A:

- keeps the current working system stable
- keeps local HTTP and append-only event semantics as the real shipped
  contract for the embodiments
- lets the next PromiseGrid step start from one deliberately chosen protocol
  family instead of trying to reinterpret all existing runtime behavior at once

Alternative B:

- immediately introduces a second representation boundary inside a system that
  is currently coherent around one local event/runtime model
- can work if the first slice is genuinely narrow
- creates migration work because some current durable records will need a clear
  mapping to a newly frozen protocol contract

Alternative C:

- adds the most capability at once
- also creates the most moving parts at once: signed envelopes, conformance
  semantics, possible exchange boundaries, and durable-format migration
- makes routine ex5 feature work subordinate to a transport/runtime rewrite

Result:

- A is easiest on the currently working operator flow.
- B is viable only if the protocol to freeze is already clear.
- C is too large for the current ex5 state.

### Scenario 2: failure, corruption, or incomplete writes

Bob restarts ex5 after partial writes, large evidence, or replay-sensitive
durability events. The current runtime already has recent durability hardening.

Alternative A:

- preserves the current replay model and recent durability fixes
- avoids immediately introducing a second durable format or signature-verifier
  failure mode

Alternative B:

- adds a new durable contract boundary and therefore new replay obligations
- requires precise answers about canonical bytes, what exactly is signed, how
  old local events map to the new wire-facing representation, and what happens
  when verification fails

Alternative C:

- adds those obligations plus transport or relay-facing partial-delivery and
  exchange-failure cases

Result:

- A keeps the durability surface smallest.
- B can be safe, but only after the exact protocol bytes are frozen.
- C compounds durability and transport failure modes too early.

### Scenario 3: concurrent actors and mixed-version nodes

Carol and Dave use different ex5 builds or different embodiments at the same
time.

Alternative A:

- keeps concurrency inside the current local runtime model
- mixed-version behavior remains a local deployment question, not yet a
  protocol-conformance question

Alternative B:

- mixed-version semantics suddenly matter at the protocol layer
- requires an auditable implementation promise claim and exact understanding of
  which frozen spec the build claims to implement
- is only responsible if the frozen spec exists first

Alternative C:

- turns mixed-version behavior into both a conformance question and an exchange
  interoperability question

Result:

- The dev guide supports B only after explicit spec freeze and implementation
  claims exist.
- C is premature for ex5.

### Scenario 4: long-horizon evolution and migration

Ellen wants ex5 to evolve over time without locking early mistakes into the
wrong wire contract.

Alternative A:

- delays wire-level commitment until one narrow contract is worth freezing
- keeps room to refine protocol-family boundaries using the already working
  local runtime

Alternative B:

- can create a healthy long-horizon path if the first frozen contract is truly
  narrow and valuable
- becomes expensive if the first frozen contract is too broad or maps poorly to
  the current ex5 event model

Alternative C:

- risks freezing too much too early
- raises the chance that later refinement becomes breaking migration instead of
  a cleanly isolated next protocol

Result:

- A is the safest evolutionary choice right now.
- B becomes attractive after a narrow first protocol family is explicitly
  chosen and frozen.

### Scenario 5: trust-boundary changes

Frank wants to move ex5 from a single local runtime toward independently
assessable peers.

Alternative A:

- does not claim more trust-boundary machinery than the runtime currently
  implements
- keeps documentation honest

Alternative B:

- is the first alternative that truly starts the PromiseGrid trust-boundary
  shift
- only makes sense if the signed bytes and protocol meaning are already fixed

Alternative C:

- expands trust-boundary handling the fastest
- also demands the most unsettled design work all at once

Result:

- The trust-boundary shift should begin spec-first, not transport-first.
- That favors A now, then B later.

### Scenario 6: scale, storage, bandwidth, and operational complexity

Mallory is not an attacker here; the pressure is simple growth in data size,
operator count, and operational expectations.

Alternative A:

- keeps complexity proportional to the current local-runtime workload
- avoids adding signature storage, envelope indexing, and exchange machinery
  before they are needed for a frozen protocol contract

Alternative B:

- adds complexity in a bounded way if the first slice is narrow
- still creates permanent new obligations around envelope storage and
  conformance publication

Alternative C:

- adds the highest storage and operational complexity immediately

Result:

- A remains the best immediate choice.
- B is the next acceptable choice only after a narrow frozen protocol is ready.

## Conclusions

### Rejected alternatives

- **Alternative C** is rejected. It is too broad for the current ex5 state and
  collapses protocol freeze, signed-envelope runtime behavior, and
  relay/transport work into one step.
- **Alternative B, right now** is rejected as the immediate next slice. The
  PromiseGrid dev guide says app work should be spec-first: pick the exact
  protocol, freeze it, then claim and implement it. Ex5 has protocol-family
  framing, but not yet one frozen ex5 protocol contract that is clearly ready
  to drive runtime behavior.

### Surviving alternative

- **Alternative A** survives: do not open the wire/runtime slice yet. The next
  PromiseGrid-aligned step should be a narrow protocol-freeze and
  implementation-claim slice.

## Implications for open TODOs and pending DIs

- `093` should close with the decision that ex5 should **not** open a real
  wire/runtime slice yet.
- A new follow-on TODO should track the narrower prerequisite: choose and
  freeze one first ex5 protocol family that is worth turning into a real
  PromiseGrid contract.
- `005` websocket transport stays deferred; this TE does not reopen that work.
- `016` Neovim follow-ons are unaffected by this decision.
- `DI-ragup` remains the source for the tracked decision question.

## Recommended conclusion

Do **not** begin the real ex5 wire/runtime PromiseGrid slice yet.

Instead:

1. pick one narrow ex5 protocol family
2. freeze that protocol contract explicitly
3. publish a real implementation promise claim against that frozen spec
4. only then open the first signed-envelope runtime slice

That sequence matches the PromiseGrid dev guide much better than starting with
runtime/wire machinery before ex5 has one clearly frozen protocol contract to
implement.
