# Ex5 Knowledge-Approval Family Boundary

TE ID: `TE-tipav`
## Status
decided

## Decision under test

What the second frozen ex5 PromiseGrid family should own when `knowledge-approval`
is frozen: both knowledge-item and run approvals in one family, item approvals
only, or separate item-vs-run approval families.

Related TODO:

- `095` - `ex5-operational-knowledge-system/TODO/TODO-vosul-ex5-second-frozen-protocol-family.md`

## Assumptions

- `ex5` already has one frozen PromiseGrid family: `knowledge-item`.
- The current runtime models both item approvals and run approvals through the
  same `approval_recorded` event type and the same `Approval` struct.
- Knowledge-item lifecycle status changes already remain a separate item-family
  concern.
- The PromiseGrid dev guide favors staged, spec-first, claim-driven migration
  by narrow durable family.
- Browser, CLI, and Neovim should remain on the current local HTTP adapter
  during this migration slice.

## Alternatives

### Alternative A

Freeze `knowledge-approval` as one family that covers both knowledge-item and
run approvals.

### Alternative B

Freeze `knowledge-approval` now, but limit it to knowledge-item approvals and
leave run approvals on the bridge layer.

### Alternative C

Split the approval work into separate families: one for knowledge-item
approvals and one for run approvals.

## Scope and systems affected

- `protocols/knowledge-approval.md`
- ex5 runtime approval storage and replay
- approval creation paths in browser, CLI, and Neovim through the shared HTTP
  adapter
- implementation claims and PromiseGrid boundary docs
- any later trust-bearing interpretation of approval artifacts

## Scenario analysis

### Scenario 1: normal operator approval workflow

Alice approves a knowledge-item revision from the browser. Bob later approves a
performed run from CLI or Neovim.

Alternative A:

- matches the current runtime model directly
- keeps one approval artifact shape regardless of target type
- lets `target_type` and `target_id` carry the boundary between item and run
  approvals without needing separate protocol families

Alternative B:

- freezes only part of the existing approval behavior
- leaves run approvals in a limbo state where they are still durable and
  trust-bearing, but not yet part of the frozen family
- creates an odd product story because users experience both as “approvals”

Alternative C:

- forces a distinction that the current product does not naturally expose
- creates extra protocol surface for little immediate gain
- adds migration and documentation overhead

Result:

- A is the cleanest fit to the shipped behavior.
- B under-freezes a clearly related trust-bearing artifact.
- C over-splits the contract.

### Scenario 2: failure, corruption, or incomplete writes

Carol restarts the runtime after one approval message is corrupted or only
partially written.

Alternative A:

- adds one approval envelope log and one verifier path
- keeps replay responsibilities centralized for all approvals

Alternative B:

- still needs a verifier path, but only for one subset of approvals
- leaves the rest of the trust-bearing approval history outside the new checks

Alternative C:

- creates multiple approval replay/verification paths
- increases the chance of inconsistent handling between item and run review
  records

Result:

- A keeps the failure surface narrowest while still covering the full current
  approval concept.

### Scenario 3: mixed-version and staged migration

Dave runs a newer ex5 build with frozen approval envelopes while Ellen still
uses older browser or Neovim flows through the same local HTTP adapter.

Alternative A:

- works well with adapter-preserving migration because all approval writes flow
  through one signed family behind the same server behavior

Alternative B:

- creates different migration behavior for item approvals versus run approvals
- increases the chance of docs and tests drifting across the two cases

Alternative C:

- makes mixed-version claims heavier because two approval-family contracts must
  be described and migrated independently

Result:

- A keeps staged migration simplest.

### Scenario 4: long-horizon evolution

Frank wants later families like evidence and links to follow the same shape of
“one durable concept per frozen family.”

Alternative A:

- keeps approval as one durable concept: named-role review outcomes
- leaves target-specific semantics in the payload rather than the top-level
  family split

Alternative B:

- risks leaving run approvals permanently second-class

Alternative C:

- encourages over-fragmentation where each target type becomes its own family

Result:

- A best matches the PromiseGrid “freeze the durable concept” rule.

### Scenario 5: trust-boundary clarity

Mallory audits which artifacts in ex5 represent trust-bearing human review.

Alternative A:

- gives one clear answer: approval artifacts live in `knowledge-approval`

Alternative B:

- requires an arbitrary explanation for why run approvals are not yet in the
  same family

Alternative C:

- requires explaining why target type, not durable concept, drives the family
  boundary

Result:

- A is the clearest trust-boundary story.

## Conclusions

Rejected alternatives:

- Alternative B: rejects part of the already-shipped approval concept without a
  good durable-boundary reason.
- Alternative C: splits the family by target type rather than by durable
  concept, adding unnecessary protocol surface.

Surviving alternative:

- Alternative A: one `knowledge-approval` family covering both knowledge-item
  and run approvals.

Recommended conclusion:

- Freeze `knowledge-approval` as one family.
- Keep lifecycle status changes outside it; those remain in `knowledge-item`.
- Implement one signed-envelope runtime slice for `approval_recorded`
  artifacts, with `target_type` and `target_id` carried in the payload.

## Implications for open TODOs and pending DIs

- `095` should lock the second family as `knowledge-approval` and implement the
  signed approval artifact log.
- The implementation claims doc should move from “one frozen family” to “two
  frozen families.”
- A later follow-on should likely make `knowledge-evidence` the third family
  unless new TE work shows a different trust-bearing dependency order.
