# Ex5 Embodiment Contract Tightening Timing

TE ID: `TE-lavok`
## Status
decided

## Decision under test

When browser, CLI, and Neovim should stop being described primarily through the
local HTTP adapter and instead bind more directly to the shipped PromiseGrid
runtime contract.

Related TODO:

- `102` - `ex5-operational-knowledge-system/TODO/TODO-zidor-ex5-embodiment-contract-tightening.md`

## Assumptions

- `ex5` has five frozen local signed families.
- The first relay-visible exchange slice is staged for the four
  attachment-free families.
- CAS-backed storage is staged as an additive sidecar for signed envelopes and
  copied evidence blobs.
- The current embodiments are stable over the local HTTP adapter and already
  share one runtime.

## Alternatives

### Alternative A

Tighten the embodiment contract now: redefine browser, CLI, and Neovim as if
they already bind directly to the PromiseGrid runtime contract instead of the
local HTTP adapter.

### Alternative B

Keep embodiments on the local HTTP adapter until the first relay-visible
exchange layer and additive CAS layer actually exist in the runtime, then
tighten the embodiment contract in a later implementation slice.

### Alternative C

Keep the local HTTP adapter as the primary embodiment contract indefinitely and
do not plan any stronger direct-runtime language.

## Scope and systems affected

- `ex5-operational-knowledge-system/TODO/TODO-zidor-ex5-embodiment-contract-tightening.md`
- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`
- `ex5-operational-knowledge-system/docs/architecture.md`
- `ex5-operational-knowledge-system/docs/practical-implementation.md`
- `ex5-operational-knowledge-system/README.md`
- new embodiment-staging documentation

## Scenario analysis

### Scenario 1: normal operator use today

Alice uses the browser, Bob uses the CLI, and Carol uses Neovim.

Alternative A:

- describes the embodiments as more direct-runtime than the shipped code
  actually is
- risks doc dishonesty and confusing migration language

Alternative B:

- keeps today’s description honest
- still preserves a clear later point when the wording should tighten

Alternative C:

- stays honest today
- but gives up on stating a later on-grid embodiment target clearly

Result:

- B is the best honest-current, staged-future description.

### Scenario 2: first relay-visible exchange ships

Dave implements the first peer-visible exchange layer over the four
attachment-free families.

Alternative A:

- would already have forced embodiment wording ahead of runtime reality

Alternative B:

- now has a concrete runtime event that justifies starting embodiment
  tightening work

Alternative C:

- still leaves embodiments described only through HTTP even after the runtime
  grows beyond that local-only story

Result:

- B ties tightening to a concrete runtime milestone.

### Scenario 3: additive CAS lands

Ellen adds CAS-backed storage for signed envelopes and copied evidence blobs.

Alternative A:

- has already front-loaded embodiment wording before storage reality caught up

Alternative B:

- can start tightening once both peer exchange and additive CAS exist, because
  the runtime contract is now materially richer than the adapter alone

Alternative C:

- keeps useful local-adapter language
- but understates the runtime contract once peer/storage layers are real

Result:

- B best matches the actual migration sequence.

### Scenario 4: long-horizon maintainability

Frank reads the docs six months later.

Alternative A:

- makes it harder to tell which parts are implemented versus merely intended

Alternative B:

- gives a crisp sequence:
  1. freeze families
  2. stage peer exchange
  3. stage CAS
  4. tighten embodiment language after those runtime layers exist

Alternative C:

- keeps immediate honesty
- but blurs the longer-term goal of getting embodiments more directly on-grid

Result:

- B creates the clearest migration narrative.

## Conclusions

Rejected alternatives:

- Alternative A: too early; it would overstate shipped embodiment/runtime
  coupling
- Alternative C: too static; it would fail to describe the intended later
  embodiment tightening point

Surviving alternative:

- Alternative B: keep embodiments on the local HTTP adapter until the first
  relay-visible exchange and additive CAS layers are implemented, then tighten
  the embodiment contract in a later slice

Implications and future work:

- `102` should close as a timing/staging decision
- the current docs should keep the local HTTP adapter description
- later implementation work can open a narrower embodiment-migration slice once
  peer/storage code exists

## Decision status

Alternative B locked by `DI-vabek`: keep browser, CLI, and Neovim described
through the local HTTP adapter until the first relay-visible exchange and
additive CAS layers are actually implemented, then tighten embodiment/runtime
language in a later slice.
