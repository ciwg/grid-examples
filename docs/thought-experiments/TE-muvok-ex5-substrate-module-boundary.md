## Title

ex5 substrate module boundary

## TE ID

TE-muvok

## Status

needs DF

## Decision under test

Whether the reusable PromiseGrid substrate currently living under
`ex5-operational-knowledge-system/promisegrid/` should remain inside the `ex5`
module boundary for now or graduate to a separate package/module boundary.

## Assumptions

- The repo now has three explicit reusable substrate slices inside `ex5`:
  `promisegrid/records/`, `promisegrid/transport/`, and `promisegrid/store/`.
- Those packages currently have one real consumer set: the `ex5`
  operational-knowledge application under `service/` plus its shipped
  embodiments and relay binary.
- `145A` is now locked: workflow composition remains inside `service/`, so the
  substrate surface is still intentionally incomplete at the application layer.
- The user wants the most PromiseGrid-aligned path, but not speculative
  packaging churn or a module split that outruns the currently proven
  substrate.

## Alternatives

### Alternative A

Keep `promisegrid/*` inside the `ex5` module boundary for now.

- Preserve the current `go.mod` and import paths
- Treat `promisegrid/*` as a semantic substrate boundary without claiming a
  separate product/module boundary yet
- Revisit the split only after more than one consumer or a stronger stability
  line exists

### Alternative B

Create a separate reusable PromiseGrid module boundary now.

- Move the current `promisegrid/*` packages into a standalone module or
  similarly strong packaging boundary
- Rewire `ex5` to consume that externalized substrate
- Use packaging to reinforce the semantic split already documented

### Alternative C

Perform a broader repo-level PromiseGrid packaging reframe now.

- Promote the current substrate slices into a more expansive shared runtime
  area
- Start shaping a multi-app PromiseGrid layout immediately
- Treat `ex5` as only one app already sitting on that broader package line

## Scenario analysis

### Scenario 1: normal development on the current shipped system

Alice continues evolving browser, CLI, Neovim, relay, and storage behavior in
the current `ex5` module.

- Under A, the semantic substrate line is already clear while development
  remains low-friction inside one module.
- Under B, every substrate change becomes a packaging and import-boundary
  change as well as a behavior change.
- Under C, the repo takes on a broader layout shift before a second real app or
  second substrate consumer exists.

What gets easier:

- A: incremental substrate growth without packaging churn
- B: stronger apparent separation between reusable and app code
- C: strongest immediate multi-package framing

What gets harder:

- A: the packaging line remains less formal than the semantic line
- B: current iteration speed slows because packaging churn joins every
  substrate change
- C: architectural scope expands sharply beyond what the repo currently proves

### Scenario 2: failure, compatibility, and release burden

Bob fixes a replay bug in `promisegrid/records`, Carol changes a relay-feed
shape in `promisegrid/transport`, and Dave updates CAS hydration in
`promisegrid/store`.

- Under A, those changes stay inside one versioned module while the substrate
  surface is still settling.
- Under B, the repo must decide whether substrate and `ex5` version
  independently, even though there is still only one shipped application using
  the substrate.
- Under C, broader packaging multiplies compatibility and release obligations
  before the substrate boundary is mature enough to justify them.

What gets easier:

- A: one release line while substrate is still maturing
- B: future independent release potential
- C: immediate large-scale reuse story

What gets harder:

- A: no independent substrate release artifact exists yet
- B: version skew becomes a new operational burden immediately
- C: compatibility work expands faster than reuse evidence

### Scenario 3: second-application thought experiment

Ellen imagines a future second PromiseGrid application with different entities
and workflows.

- Under A, that future app can still import `promisegrid/*` once the need is
  real, and the repo can split packaging at that point using actual consumer
  evidence.
- Under B, the split happens in anticipation of that app rather than because
  the app already exists.
- Under C, the repo effectively commits to a framework-first packaging shape
  before the second app proves what should really be shared.

What gets easier:

- A: future split can be designed around real second-consumer needs
- B: a second app would find an already separated module
- C: a second app could slot into a large prebuilt structure

What gets harder:

- A: the second app would still need the split when it arrives
- B: the current module line may freeze around one app's present needs
- C: later apps may inherit abstractions chosen too early

### Scenario 4: PromiseGrid honesty versus packaging symbolism

Frank asks whether PromiseGrid alignment requires a separate module today.

- Under A, the repo answers honestly: semantic substrate is real, but separate
  packaging is not yet justified by consumer or release evidence.
- Under B, the repo uses packaging to signal seriousness, but risks implying a
  more stable or independently consumable substrate than the code currently
  proves.
- Under C, the repo uses packaging to announce a broader runtime vision before
  the remaining example-specific surfaces have truly settled.

What gets easier:

- A: semantics stay ahead of branding
- B: the separation becomes more visible in the filesystem and import paths
- C: the repo looks closest to a generalized PromiseGrid product

What gets harder:

- A: some readers may want a stronger packaging line than the code yet needs
- B: packaging symbolism may outrun substance
- C: the packaging story may overstate current generality

### Scenario 5: current substrate incompleteness after `145A`

Grace looks at the state after locking `145A`: workflows remain intentionally
inside `service/`.

- Under A, the packaging reflects that the substrate is real but still
  intentionally partial.
- Under B, the packaging suggests the extracted pieces are now ready to stand
  alone even though the boundary above them is intentionally unresolved.
- Under C, the repo broadens packaging right after deciding not to broaden
  workflow substrate, which pulls in opposite directions.

What gets easier:

- A: the package line stays consistent with the intentionally partial substrate
- B: the lower-level substrate is isolated sooner
- C: broad reframe now instead of later

What gets harder:

- A: no standalone module yet
- B: the standalone module may invite premature expansion pressure
- C: mismatch between intentionally partial substrate and aggressively broad
  packaging

### Scenario 6: maintenance and import-path stability

Heidi looks at the current imports. `service/*` already depends on
`github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/...`.

- Under A, import paths stay stable while the extracted slices continue to
  settle.
- Under B, all current substrate imports churn now, even though only `ex5`
  consumes them.
- Under C, import and directory churn grows further because the repo would
  also be choosing a broader shared layout.

What gets easier:

- A: minimal churn
- B: potential future independent consumption path
- C: strongest repo-wide package separation

What gets harder:

- A: no early independent import root
- B: immediate refactor cost without second-consumer proof
- C: bigger refactor and longer-lived migration risk

## Conclusions

Rejected:

- Alternative C: too broad and too speculative. It asks packaging to lead the
  architecture before the repo has more than one real application or a broader
  settled substrate surface.

Surviving:

- Alternative A: keep the reusable substrate inside the `ex5` module boundary
  for now
- Alternative B: split the current substrate into a separate module boundary
  now

Recommended:

- Alternative A is the strongest PromiseGrid-aligned surviving choice.

Why A is more aligned than B:

- PromiseGrid alignment here is evidence-first, not symbolism-first.
- The semantic substrate boundary is already explicit and honest.
- There is still one real consumer set, one release line, and an intentionally
  partial substrate after `145A`.
- A separate module becomes more justified when there is real independent
  consumption, a second application, or a need for independent versioning and
  release cadence.

Alternative B remains viable later if future evidence appears, especially:

- a second application consuming `promisegrid/*`
- independent versioning pressure between substrate and `ex5`
- a stabilized substrate surface that no longer changes in lockstep with ex5's
  remaining app/runtime evolution

## Implications for open TODOs and pending DIs

- `TODO-tolav-ex5-substrate-module-boundary.md` should either lock Alternative
  A and close as an intentional stay-nested decision, or continue only if the
  user wants to force an early packaging split despite the current evidence.
- `DI-tolav` remains active until the user resolves the post-TE DF.
- A future module-boundary wave should be justified by consumer, versioning, or
  release evidence, not only by aesthetics or symbolic separation.
