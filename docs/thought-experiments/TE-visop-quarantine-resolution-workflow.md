# Quarantine resolution workflow

TE ID: TE-visop

## Status

decided

## Decision under test

How Ex6 should express the authorized terminal resolution of an open
quarantine case: one typed workflow, two transition-specific workflows, or no
workflow artifact at all.

## Assumptions

- TE-maluz and DI-hogid established that receiving-exception only opens a
  case; release and reject are separate, explicit transitions.
- The quarantine package already retains append-only `open`, `release`, and
  `reject` events and refuses a second terminal transition.
- The runtime persists workflow input/output and dispatches only an active,
  pCID-matched artifact through its selected adapter.
- Alice may inspect and propose a release, Bob may reject a case, and Carol
  may review the same evidence. Mallory may supply an untrusted artifact or
  misleading evidence but cannot make the runtime append a transition without
  the local activation and route-evidence requirements.

## Alternatives

1. **One quarantine-resolution workflow.** One typed input carries `case_id`,
   `event_id`, `actor`, `evidence_id`, `decision`, and `notes`, where
   `decision` is exactly `release` or `reject`.
2. **Two transition-specific workflows.** Separate `quarantine-release` and
   `quarantine-reject` artifacts each encode one terminal transition.
3. **Direct package commands only.** Operators call `quarantine release` or
   `quarantine reject` without a workflow artifact.

## Scenario analysis

### Normal authorized review

Alice reviews an open case, determines that the inspection evidence clears the
material, and records a release. Bob reviews a different case and determines
that the material must be rejected.

Alternative 1 gives both actions one stable resolution contract while keeping
the chosen terminal promise explicit in the typed `decision` field. The
resulting event remains either `release` or `reject`; the workflow does not
turn the two meanings into an inferred status.

Alternative 2 gives maximal artifact-level separation but duplicates schema,
policy, adapter, and lifecycle handling for transitions that share the same
case/evidence shape. It adds no additional durable distinction because the
package event already records the transition.

Alternative 3 records the package event but loses the workflow's retained
input, activation decision, pCID contract, and route-promise evidence.

### Interrupted or repeated resolution

Carol starts a resolution and the runtime stops before the adapter applies the
event. Later she retries; meanwhile Bob attempts a competing reject.

Alternative 1 retains a typed run that distinguishes an incomplete attempted
resolution from a completed terminal event. The package's rebuildable case view
accepts the first terminal event and rejects the later competing one.

Alternative 2 has the same package-level protection, but creates two possible
artifact contracts for the same concurrency race. Alternative 3 leaves only a
command error and no retained workflow evidence for the attempted action.

### Mixed versions and long horizon

Alice retains an older workflow while a later policy requires stronger evidence
for resolution. A future corrective-action workflow may consume the retained
resolution event but is not part of this slice.

Alternative 1 versions one general resolution pCID when its evidence contract
changes. Old artifacts remain inspectable and executable only under their
explicit local activation. The package event preserves the concrete terminal
meaning across versions.

Alternative 2 doubles the migration and policy surface. Alternative 3 offers
no artifact contract to version or retain.

### Trust boundary and scale

Mallory offers an artifact that claims it can release cases automatically.
Across many sites, local operators use different evidence sources and review
policies.

Alternative 1 permits each site to activate only a trusted resolution artifact
and make explicit receive/delivery promises for its pCID. The workflow carries
an evidence reference but does not treat a package claim or imported artifact
as authority. A local policy can later constrain accepted decision values or
evidence without moving quarantine semantics into the runtime.

Alternative 2 preserves that property with more artifacts to govern.
Alternative 3 weakens the visible trust boundary by bypassing the workflow
lifecycle altogether.

## Conclusions

Alternative 3 is rejected because it bypasses the workflow evidence and
activation model.

Alternative 2 survives for a future case where release and reject genuinely
require different package sets, evidence schemas, or authorities. Those
differences are not established in the current scope.

Alternative 1 is recommended. One pCID-defined resolution workflow keeps the
terminal decision explicit, reuses the package's durable event distinction,
and avoids duplicating an otherwise identical artifact contract.

## Output to decision framing

The surviving recommendation is one `quarantine-resolution` workflow with a
typed `decision` field restricted to `release` or `reject`.

Remaining decisions:

1. Whether the first workflow requires the case to link to a specific prior
   opening-event ID in addition to the package's open-state validation.
2. Whether evidence requires one reference or separate review and authority
   references.
3. Exact source paths, adapter naming, schemas, output fields, test paths,
   and runtime temporary-path pattern.

## Implications

- A new TODO and DI set is required before implementation.
- The existing quarantine package remains the sole owner of terminal-event
  validation; the workflow adapter selects a package command only after its
  typed decision is validated.

## Decision status

Locked by DI-nufav.
