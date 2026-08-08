# Ex5 Run Family Boundary

TE ID: `TE-vamok`
## Status
decided

## Decision under test

What the first frozen PromiseGrid-native run family for `ex5` should contain so
peer-visible evidence can attach to a durable run-context contract without
collapsing other already-frozen families back into one large record.

Related TODO:

- `108` - `ex5-operational-knowledge-system/TODO/TODO-zaruv-ex5-sixth-frozen-protocol-family-run.md`

## Assumptions

- `ex5` already has frozen families for `knowledge-item`,
  `knowledge-approval`, `knowledge-evidence`, `knowledge-link`, and
  `knowledge-responsibility`.
- The current local `run_recorded` compatibility event includes:
  - run ID
  - item ID and revision
  - actor
  - outcome and notes
  - place/resource/responsibility references
  - optional machine/location strings
- Evidence, approvals, and typed links are already modeled as separate
  families or compatibility projections anchored to runs.
- The user prefers the more PromiseGrid-complete path rather than a shortcut
  that carries unfrozen compatibility events forever.

## Alternatives

### Alternative A

Freeze a minimal run family around just the execution anchor:

- run identity
- item identity
- revision number
- actor
- timestamp
- outcome
- notes

All place/resource/responsibility context stays outside the first frozen run
family.

### Alternative B

Freeze a richer run family around the full current durable run context:

- run identity
- item identity
- revision number
- actor
- timestamp
- outcome
- notes
- place ID
- resource IDs
- responsibility IDs
- machine
- location

Evidence, approvals, and links remain separate families anchored to the run.

### Alternative C

Freeze a maximal run family that also subsumes evidence, approvals, and links
inside the run contract.

This would make a run one large aggregate instead of keeping those already
separate durable families independent.

## Scope and systems affected

- `ex5-operational-knowledge-system/TODO/TODO-zaruv-ex5-sixth-frozen-protocol-family-run.md`
- `ex5-operational-knowledge-system/TODO/TODO.md`
- new run protocol doc under `ex5-operational-knowledge-system/protocols/`
- `ex5-operational-knowledge-system/protocols/profiles.go`
- `ex5-operational-knowledge-system/service/app.go`
- new signed run envelope builder/verification code
- peer-exchange and evidence follow-on work under TODO `105`

## Scenario analysis

### Scenario 1: evidence needs a durable run anchor

Alice records a receiving run, then attaches signed evidence to it. Bob later
imports both the run and the evidence.

Alternative A:

- gives evidence a durable run anchor
- but omits place/resource/responsibility context that is often part of why the
  evidence matters operationally
- would force later peer-visible consumers to reconstruct context from other
  non-run sources or compatibility history

Alternative B:

- gives evidence a durable run anchor plus the full current operational run
  context
- makes the run family immediately useful as the thing evidence points at

Alternative C:

- also gives evidence an anchor
- but collapses already-separated evidence into the run contract and fights the
  family-by-family migration shape already used in `ex5`

Result:

- B fits the current evidence dependency best without over-aggregating.

### Scenario 2: run review and approval semantics

Carol reviews a run, approves it, and later inspects why the run happened in a
given place with a given resource set.

Alternative A:

- preserves run approvals as a separate family
- but weakens the run family as an independently portable operational record

Alternative B:

- keeps approvals separate while still making the run itself rich enough to
  stand alone in peer exchange and inspection

Alternative C:

- removes the clean family boundary between run and approval
- would force later changes in approval semantics to mutate the run family

Result:

- B keeps the family layering cleaner.

### Scenario 3: typed links to runs

Dave links a responsibility or item to a run and later shares that graph.

Alternative A:

- makes linked runs portable only in a skeletal way
- the linked run would exist, but much of its operational meaning would still
  live elsewhere

Alternative B:

- gives linked runs the richer context already expected by current `RunRecord`
  projections

Alternative C:

- again over-aggregates and undermines the independent link/evidence/approval
  families

Result:

- B best matches the current graph semantics.

### Scenario 4: long-horizon fully-on-grid evolution

Steve wants `ex5` eventually to share real operational runs across peers, not
just local browser/CLI views.

Alternative A:

- is more conservative and easier to freeze
- but may require a second run-family expansion later, which is exactly the
  kind of stepping stone the current PromiseGrid direction is trying to avoid

Alternative B:

- freezes the run family closer to the current real operational meaning of a
  run
- gives later evidence, approval, and link exchange a stronger base record

Alternative C:

- goes too far by re-absorbing already-separated durable families

Result:

- B is the strongest long-horizon choice.

### Scenario 5: implementation risk

Ellen wants the next slice to stay tractable.

Alternative A:

- has the smallest initial payload
- but likely leaves too much follow-on expansion pressure

Alternative B:

- is bigger, but still structurally aligned with the existing `RunRecord`
  projection and current `run_recorded` event payload
- can reuse more of the current compatibility-event semantics directly

Alternative C:

- is largest and most disruptive because it would cut across evidence,
  approval, and link families already frozen separately

Result:

- B is the best balance of tractability and completeness.

## Conclusions

Rejected alternatives:

- Alternative C: too aggregate-heavy; it collapses already-separated families
- Alternative A: plausible, but too skeletal for the role runs already play in
  evidence, approval, and link context

Surviving alternative:

- Alternative B: freeze a rich run family that captures the full current run
  context, while keeping evidence, approvals, and links as separate families

Unresolved questions that still require user choice:

- none; naming and richer scope were locked by `DI-vamok`

## Decision status

Locked by `DI-vamok`:

1. family name: `operational-run`
2. lock Alternative B as the first run-family scope
3. carry current `place_id`, `resource_ids`, and `responsibility_ids` as-is in
   the first frozen family, with later peer-stable tightening deferred to TODO
   `107`
