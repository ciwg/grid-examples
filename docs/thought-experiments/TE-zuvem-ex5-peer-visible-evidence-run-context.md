# Ex5 Peer-Visible Evidence Run Context

TE ID: `TE-zuvem`
## Status
needs DF

## Decision under test

How `ex5` should satisfy run context when making `knowledge-evidence`
peer-visible, given that evidence currently replays only against existing local
`run_recorded` events and runs are not yet a frozen peer-visible family.

Related TODO:

- `105` - `ex5-operational-knowledge-system/TODO/TODO-faruv-ex5-peer-visible-evidence-exchange.md`

## Assumptions

- `knowledge-evidence` is already a frozen signed family.
- The current peer-exchange bundle does not export or import `run_recorded`
  events.
- `applyEventLocked("evidence_added")` still requires the referenced run to
  exist locally.
- The user wants the next step to be as PromiseGrid-complete as possible rather
  than a quietly lossy approximation.

## Alternatives

### Alternative A

Let the first peer-visible evidence slice carry the required `run_recorded`
compatibility events alongside evidence, without first freezing a run family.

This keeps evidence exchange moving now, but it extends peer-visible behavior
through unfrozen compatibility events.

### Alternative B

Freeze and claim a PromiseGrid-native run family before peer-visible evidence
exchange lands.

This makes evidence exchange wait for a more durable run-context contract, but
it keeps the peer-visible model centered on frozen families rather than on
ad-hoc compatibility exceptions.

### Alternative C

Defer peer-visible evidence exchange until some other peer-visible mechanism
already carries run context, without deciding that mechanism yet.

## Scope and systems affected

- `ex5-operational-knowledge-system/TODO/TODO-faruv-ex5-peer-visible-evidence-exchange.md`
- peer-exchange bundle shape
- any future run-family or compatibility-event transport decisions
- PromiseGrid claims and peer-exchange staging docs

## Scenario analysis

### Scenario 1: bootstrap import of evidence into a fresh runtime

Alice exports one run plus its evidence to Bob.

Alternative A:

- can work immediately by shipping the necessary `run_recorded` compatibility
  events with the evidence bundle
- but it introduces a peer-visible dependency on an unfrozen event type

Alternative B:

- delays the feature until a proper run family exists
- but then evidence imports onto a stronger durable contract

Alternative C:

- delays implementation without clarifying the eventual mechanism

Result:

- A is fastest; B is stronger.

### Scenario 2: strict PromiseGrid framing

Steve wants the peer-visible surface to rely on frozen family contracts rather
than ad-hoc exceptions as much as possible.

Alternative A:

- leaks compatibility-event transport into the peer-visible model
- risks creating a special-case precedent that later has to be unwound

Alternative B:

- aligns better with the family-by-family migration style already used so far

Alternative C:

- avoids the exception too
- but gives no concrete next step

Result:

- B is the strongest PromiseGrid-complete direction.

### Scenario 3: implementation speed

The repo wants progress now without inventing a large new family if it can be
avoided.

Alternative A:

- is clearly the shortest path
- uses existing run events instead of waiting for a new frozen contract

Alternative B:

- is slower
- but cleaner for long-horizon semantics

Alternative C:

- is the slowest and least informative

Result:

- A wins on speed; B wins on durable model quality.

## Conclusions

Rejected alternative:

- Alternative C: too vague; it defers the feature without choosing a real path

Surviving alternatives:

- Alternative A: first evidence exchange includes required `run_recorded`
  compatibility context
- Alternative B: freeze a run family first, then make evidence peer-visible

## Decision status

Needs DF between Alternative A (faster compatibility-event carry-along) and
Alternative B (cleaner PromiseGrid-first run family before evidence exchange).
