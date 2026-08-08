## Title

ex5 minimal workflow substrate evidence

## TE ID

TE-vunek

## Status

needs DF

## Decision under test

Whether `ex5-operational-knowledge-system` now has enough evidence to extract a
minimal reusable PromiseGrid workflow substrate beyond the existing
`promisegrid/records`, `promisegrid/transport`, and `promisegrid/store`
packages, or whether workflow composition should stay inside `service/`.

## Assumptions

- `ex5` now has explicit reusable substrate slices for durable record truth,
  peer and relay transport wire truth, and append-only log/CAS/hydration
  persistence.
- The remaining higher layers in `service/` include entity-specific projections
  and workflows for places, resources, responsibilities, items, runs,
  evidence, approvals, links, search, pending review, and problem review.
- The user wants the most PromiseGrid-aligned outcome, but also wants the repo
  to avoid speculative abstraction and keep working behavior stable.
- This TE is about whether a new workflow substrate is justified by actual
  shipped evidence, not about whether workflow extraction could be imagined in
  the future.

## Alternatives

### Alternative A

Do not extract any workflow substrate now.

- Keep workflow composition in `ex5-operational-knowledge-system/service/`
- Treat the current typed `operation` families as embodiment contracts for the
  `ex5` application, not proof of a reusable generic workflow layer
- Tighten docs so the current substrate boundary is stated honestly

### Alternative B

Extract a very narrow workflow substrate now.

- Pull out only the smallest plausibly shared workflow/query composition slice
- Likely candidates would be review-queue assembly, projection-backed search
  composition, or typed operation grouping
- Leave all entity-specific create/operate semantics in `service/`

### Alternative C

Extract a broad workflow substrate now.

- Move create, review, approval, search, and projection orchestration into a
  reusable PromiseGrid workflow layer
- Treat `service/` primarily as one adapter/user-surface host over that layer

## Scenario analysis

### Scenario 1: normal ex5 application use

Alice uses browser, CLI, and Neovim to create places and resources, author
items, run work, attach evidence, approve outcomes, and inspect pending review.

- Under A, those flows stay together inside the operational-knowledge app that
  names them. The substrate remains focused on records, transport, and
  persistence.
- Under B, the repo starts to split out a narrow orchestration layer, but the
  candidate workflows are still described in ex5 terms such as pending review,
  problem review, and operational search.
- Under C, the app's workflow semantics are generalized immediately, even
  though their current meaning is still tied to this one application.

What gets easier:

- A: clear ownership of app-specific flows
- B: an earlier start on possible workflow reuse
- C: the strongest appearance of a generalized runtime stack

What gets harder:

- A: no immediate workflow reuse package exists
- B: the repo has to maintain a new boundary that may not reflect real reuse
- C: the substrate likely absorbs one app's semantics prematurely

### Scenario 2: failure, corruption, and replay

Bob replays a damaged runtime, Carol rehydrates records from CAS, and Dave
checks whether workflow repair semantics are substrate or app responsibilities.

- Under A, durable repair stays in `promisegrid/store` and
  `promisegrid/records`, while workflow recomputation remains an app concern.
- Under B, a narrow workflow layer may need to define what counts as generic
  recomputation versus ex5-specific projection policy.
- Under C, broad workflow extraction creates larger obligations around generic
  repair semantics that the repo has not yet proven across multiple apps.

What gets easier:

- A: durable substrate duties stay crisp
- B: some repeated orchestration code could move out if it is truly generic
- C: one broad story about recovery and workflows

What gets harder:

- A: none of the app workflow recomposition is reusable yet
- B: the boundary between substrate repair and app projection policy becomes
  harder to police
- C: failure semantics risk becoming framework-shaped around ex5

### Scenario 3: mixed-version nodes and staged evolution

Ellen wants one more year of incremental evolution without a speculative
framework wave.

- Under A, the repo can keep extracting only the layers with clear reuse proof.
- Under B, the repo starts a new package that still depends on ex5-shaped
  meanings, increasing migration and compatibility burden.
- Under C, staged evolution becomes harder because the app/runtime line shifts
  broadly in one wave.

What gets easier:

- A: evidence-first extraction continues
- B: some future workflow reuse might start sooner
- C: a potential big-picture end state if the abstraction is right

What gets harder:

- A: a future workflow substrate may still need work later
- B: temporary mixed ownership with uncertain boundaries
- C: larger churn and more partial states to support

### Scenario 4: trust-boundary and PromiseGrid clarity

Frank asks which parts of the repo are truly PromiseGrid substrate and which
parts are the operational-knowledge application.

- Under A, the answer stays honest: substrate owns durable records, transport,
  and persistence; `service/` owns ex5 projections and workflows.
- Under B, the answer gets muddier unless the extracted slice is obviously
  generic across more than one application.
- Under C, the answer sounds more generalized, but only by assuming workflows
  that are still application-shaped are reusable substrate.

What gets easier:

- A: honest substrate/app boundary
- B: a possible stepping stone toward future workflow reuse
- C: a more ambitious general-runtime narrative

What gets harder:

- A: the repo stops short of claiming a workflow substrate today
- B: readers may over-read a tiny extracted layer as proof of broad workflow
  generality
- C: PromiseGrid substrate risks being defined by ex5 user flows

### Scenario 5: second-application thought experiment

Grace imagines a future second PromiseGrid application that does not use ex5's
operational entities or review vocabulary.

- Under A, that future app can still reuse records, transport, and persistence
  without inheriting ex5 workflow semantics.
- Under B, the second app might inherit a narrow workflow layer whose actual
  API was named around ex5 review/search concepts.
- Under C, the second app almost certainly inherits a workflow substrate that
  was generalized from one example before cross-app evidence existed.

What gets easier:

- A: future apps only adopt proven substrate
- B: there may be a starting point for limited orchestration reuse
- C: future apps get a broad stack immediately

What gets harder:

- A: no workflow layer is ready-made
- B: later apps may need to unlearn ex5 vocabulary hidden in a narrow layer
- C: later apps may be constrained by abstractions that fit ex5 better than
  them

### Scenario 6: scale and maintenance of current operation families

Heidi looks at the current typed `operation` families for browser, CLI, and
Neovim and asks whether they already imply a generic workflow substrate.

- Under A, those operation names remain embodiment contracts for the ex5 app.
  They prove direct contract structure, not generic workflow semantics.
- Under B, the repo would likely start extracting around those operation names,
  even though many still map directly to ex5 entities and review flows.
- Under C, the repo would treat those operation families as the beginning of a
  fully generic orchestration layer, which overstates what they currently show.

What gets easier:

- A: operation families stay honest about what they currently are
- B: some operation-grouping code could be shared
- C: one uniform workflow layer story

What gets harder:

- A: no immediate workflow package
- B: the extracted layer may mostly wrap ex5 vocabulary
- C: genericity is claimed before it is proven

## Conclusions

Rejected:

- Alternative C: too broad and too speculative. The shipped workflows are still
  clearly application-shaped, and broad extraction would freeze ex5 semantics
  into PromiseGrid substrate prematurely.

Surviving:

- Alternative A: keep workflow composition in `service/` for now and treat the
  current substrate boundary as honest and complete enough
- Alternative B: extract a very narrow workflow layer now if, and only if, a
  slice can be named without ex5-specific entity or review vocabulary

Recommended:

- Alternative A is the strongest PromiseGrid-aligned surviving choice.

Why A is more aligned than B:

- PromiseGrid alignment here is evidence-first, not abstraction-first.
- The repo has real reuse proof for record, transport, and persistence
  substrate, but not yet for workflow composition.
- The current typed operation families are direct embodiment contracts for
  `ex5`, not proof that review, search, approval, or operational projections
  are app-agnostic substrate.

Alternative B remains conceptually possible, but only after a future TE can
show one workflow slice that is:

- named without ex5-specific entity or review vocabulary
- useful across more than one application shape
- narrower than the current `service/app.go` orchestration surface

## Implications for open TODOs and pending DIs

- `TODO-rasuv-ex5-minimal-workflow-substrate-evidence.md` should either lock
  Alternative A and close as an intentional no-extraction decision, or capture
  a sharply defined follow-on proof target if the user prefers to keep
  exploring Alternative B.
- `TODO-tolav-ex5-substrate-module-boundary.md` should remain deferred behind
  this decision. Module-boundary work is less justified if workflow substrate
  evidence is still intentionally absent.
- `DI-rasuv` remains active until the user locks the post-TE decision.
