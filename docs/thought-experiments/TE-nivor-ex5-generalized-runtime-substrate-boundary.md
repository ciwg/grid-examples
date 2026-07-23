## Title

ex5 generalized runtime substrate boundary

## TE ID

TE-nivor

## Status

needs DF

## Decision under test

Whether `ex5-operational-knowledge-system` should open a real extraction wave
for a more generalized PromiseGrid runtime substrate, and if so, which boundary
should remain `ex5`-specific instead of being pulled into a reusable layer.

## Assumptions

- `ex5` is no longer a small local-only demo. It now ships signed frozen
  families, CAS-backed replay, direct local embodiment contracts, a dedicated
  remote relay, and a Chrome/Chromium native-messaging browser embodiment.
- The repo currently mixes two kinds of things inside `ex5`:
  - reusable PromiseGrid-shaped mechanics such as identity, origin tracking,
    envelope replay, relay-feed carriage, local direct-contract plumbing, and
    CAS-backed storage
  - operational-knowledge-specific projections and workflows such as items,
    runs, places, resources, responsibilities, approvals, search views, and
    review queues
- The user wants the most PromiseGrid-aligned path, but also wants the system
  to keep working, stay tested, and avoid speculative abstraction churn.
- This TE is about architectural boundary and extraction scope, not about
  changing the already shipped durable semantics.

## Alternatives

### Alternative A

Keep `ex5` as the bounded example application and do not open a real extraction
wave now.

- Continue to document the reusable-looking pieces as implementation details
  inside `ex5`
- Improve comments and docs when needed, but do not split code into a broader
  substrate
- Accept that `ex5` remains both example app and host for substrate-like logic

### Alternative B

Open a thin generalized runtime substrate wave now.

- Extract the directly reusable PromiseGrid-shaped mechanics first:
  - identity and origin ordering
  - frozen-family envelope replay and verification
  - CAS object and blob carriage
  - peer-exchange and relay-feed transport logic
  - direct local contract framing and transport helpers
- Keep `ex5`-specific operational domain logic in the example app:
  - knowledge-item, run, place, resource, responsibility projections
  - search, pending review, and problem review projections
  - browser/CLI/Neovim user-facing workflow composition
- Treat adapters and embodiments as consumers of the substrate, not as the
  substrate itself

### Alternative C

Open a broad generalized runtime framework wave now.

- Extract not only the lower-level PromiseGrid mechanics but also the higher
  read-model, operation, and embodiment layers into a generalized runtime
- Aim for a reusable runtime/product substrate that could host multiple
  applications beyond `ex5`
- Move quickly toward a larger multi-package architecture in one wave

## Scenario analysis

### Scenario 1: normal operation in the current shipped system

Alice runs one local runtime, edits from browser and Neovim, reviews via CLI,
and exchanges records through relay.

- Alternative A keeps the working system untouched, but the repo continues to
  present one example app as both the application and the implicit substrate.
- Alternative B separates the reusable PromiseGrid mechanics from the
  operational app logic without changing the meaning of the shipped workflows.
- Alternative C tries to redraw the entire runtime/app line at once, which
  creates more moving parts than the current working system requires.

What gets easier:

- A: zero architectural churn
- B: clearer line between substrate mechanics and example-app semantics
- C: maximum future reuse if the split is right

What gets harder:

- A: substrate reuse remains implicit and tangled with `ex5`
- B: targeted extraction work and package-boundary design
- C: significantly broader migration burden

### Scenario 2: failure, corruption, and compatibility burden

Bob hits a corrupted envelope log, Carol replays older runtime state, and Dave
uses a slightly older embodiment.

- Under A, the recovery logic still works, but substrate-like concerns stay
  embedded in example-app files and cannot be reasoned about independently.
- Under B, the failure-handling code for replay, origin tracking, CAS, relay,
  and local contract framing can become reusable and testable on their own,
  while the app projections remain `ex5`-specific.
- Under C, the recovery semantics might become cleaner eventually, but only
  after a much larger extraction lands without introducing regressions.

What gets easier:

- A: no migration risk
- B: reusable mechanics gain tighter isolated tests and clearer ownership
- C: eventual full-stack uniformity

What gets harder:

- A: bug fixes in reusable mechanics still land only as `ex5` internals
- B: some cross-package seams have to be made explicit
- C: compatibility and regression surface expands sharply

### Scenario 3: mixed-version evolution and staged adoption

Ellen wants to evolve the repo over time without forcing one giant rewrite.

- Alternative A delays the split entirely, which avoids risk now but makes the
  future extraction bigger because more behavior will accumulate inside `ex5`.
- Alternative B supports staged migration. The substrate can grow under the
  current binaries and tests while `ex5` keeps shipping the same workflows.
- Alternative C pulls many boundaries at once and makes staged adoption harder
  because there are more package moves and more partial states to support.

What gets easier:

- A: immediate stability
- B: additive extraction with current binaries preserved
- C: big-picture end state if everything lands cleanly

What gets harder:

- A: later extraction cost increases
- B: temporary mixed package ownership during transition
- C: rollout coordination and review burden

### Scenario 4: trust-boundary and PromiseGrid clarity

Frank asks what in the repo is “the example application” versus “the reusable
PromiseGrid runtime substrate.”

- Under A, the answer stays muddy. The docs can say `ex5` is an example, but
  the code still mixes durable substrate mechanics with application semantics.
- Under B, the answer becomes explicit: reusable substrate handles identity,
  envelope durability, relay/feed/blob carriage, and direct-contract plumbing,
  while `ex5` remains the operational knowledge application that projects and
  names those records.
- Under C, the answer might become even broader, but only by defining a more
  general runtime/app framework than the repo has yet proven necessary.

What gets easier:

- A: no new packaging decisions
- B: honest substrate/app distinction
- C: strongest eventual general-runtime story

What gets harder:

- A: PromiseGrid alignment stays partially rhetorical
- B: extraction boundary has to be chosen carefully
- C: higher risk of speculative abstraction

### Scenario 5: scale and maintenance over time

Grace imagines one more example app or another protocol slice appearing in the
repo later.

- Alternative A means the next app either duplicates substrate-like mechanics
  or reuses them informally from `ex5`, both of which are weak long-term
  outcomes.
- Alternative B gives the repo a reusable lower layer without pretending the
  application-level review/search/workflow semantics are generic already.
- Alternative C optimizes for a possible multi-app future immediately, but at
  the cost of building generalized surfaces that may still be premature.

What gets easier:

- A: nothing new to maintain today
- B: later apps can reuse real substrate pieces
- C: later apps might fit quickly if the abstraction proves right

What gets harder:

- A: future duplication or informal coupling
- B: the substrate API has to stay intentionally narrow
- C: maintenance of broader abstractions before second-app proof exists

### Scenario 6: embodiment and relay evolution after `132`

Heidi looks at the current state after browser direct contract, terminal direct
contract, and dedicated relay.

- Alternative A leaves those new direct-contract and relay pieces inside `ex5`
  as app-owned implementation details, even though they now look more like
  reusable substrate capabilities.
- Alternative B treats those pieces as the first real evidence that a reusable
  substrate exists, but resists extracting higher-level app workflows that are
  still clearly `ex5`-specific.
- Alternative C pulls embodiment workflow semantics and review/search
  operations into a generalized layer too early, before a second application
  proves those surfaces are generic.

What gets easier:

- A: no further restructuring
- B: direct-contract and relay logic can become reusable without overstating
  the generality of review/search/app projections
- C: single large architectural story

What gets harder:

- A: direct-contract and relay work remain harder to reuse
- B: package split requires careful boundary discipline
- C: risk of building framework abstractions around one app's semantics

## Conclusions

Rejected:

- Alternative C: too broad and too speculative. It tries to generalize not
  only PromiseGrid-shaped substrate mechanics but also application-shaped
  workflow and projection semantics that are not yet proven reusable.

Surviving:

- Alternative A: keep `ex5` as a bounded example and defer extraction
- Alternative B: open a thin generalized substrate wave around reusable
  PromiseGrid mechanics while leaving operational-knowledge projections and
  workflows inside `ex5`

Alternative B is the strongest PromiseGrid-aligned surviving path. It draws the
substrate boundary where the repo now has real evidence of reuse:

- identity and origin ordering
- signed-envelope replay/verification
- CAS-backed object and blob handling
- peer-exchange and relay-feed carriage
- direct local contract framing and embodiment-bridge helpers

It does **not** pretend that `ex5` search, review queues, approval workflows,
or operational projections are already generic substrate.

Alternative A remains viable if the user decides the repo should stop at an
honest example boundary for now. It is the safer deferral choice, but less
PromiseGrid-aligned than B because the reusable mechanics would continue to
live only as `ex5` internals.

## Implications for open TODOs and pending DIs

- TODO `133` should now narrow to a real decision between deferring extraction
  (`133A`) and opening a thin reusable substrate wave (`133B`).
- If `133B` is chosen, the next follow-on should be a new TODO that names the
  exact initial substrate package boundary and migration order.
- If `133A` is chosen, `133` should close as an intentional deferral with docs
  that keep the boundary honest.
