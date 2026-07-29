# Grid-Compliant Workflow Lifecycle

TE ID: TE-gavuk

## Status

needs DF

## Decision under test

Determine how OKAR imports, activates, deactivates, replaces, and retains
workflow artifacts without conflating artifact evidence, route eligibility, or
agent execution authority.

## Handle note

`tools/mint-handle` was unavailable in this checkout. `gavuk` is a manual
proquint-style handle, matching the existing OKAR local exception.

## Assumptions and trust model

- A workflow is a content-addressed artifact plus protocol-scoped implementation
  promises, not a new PromiseGrid top-level action kind.
- Import provenance and durable history must remain legible after local removal.
- The kernel hosts local lifecycle/resource roles; app agents retain workflow and
  protocol judgment. Source: PromiseGrid Development Guide, Kernel Devs.
- Docker confinement is a local execution mechanism; import never grants worker
  execution authority.
- Mallory can supply malformed, unsigned, revoked, or incompatible workflow
  artifacts; Alice may deactivate a workflow without deleting its evidence.

## Alternatives

### A. Import immediately activates a workflow

Copy a package into the runtime package root and derive routes immediately.

### B. Imported artifact registry with explicit activation state

Store imported workflow identity/provenance separately from local activation.
Only active entries may expose routes or be selected for worker execution.

### C. Delete package artifacts on removal

Treat removal as deleting installed files and all local references.

### D. Deactivate or revoke eligibility while retaining evidence

Keep the artifact CID, provenance, and historical records; remove only local
route and execution eligibility. Replacements are new additive artifacts.

## Scenario analysis

### Normal import and activation

Alice imports a signed workflow artifact. A activates it before an operator can
assess its claims. B records import first, then requires explicit local
activation. B separates evidence from authority and is compatible with a later
agent receive-promise registration. C/D do not change import behavior.

### Malformed or untrusted artifact

Mallory supplies an artifact with an unknown or unsupported pCID. A risks
making local availability look like protocol support. B retains a rejected or
inactive artifact only if provenance validation permits storage, and never
exposes a route or worker. This preserves evidence without making a promise.

### Deactivation and removal

Alice deactivates a workflow after a vulnerability. C destroys the evidence
needed to explain prior records and breaks historical interpretation. D removes
local route/worker eligibility while retaining the content-addressed artifact
and history. A/B determine activation; D is the only survivable removal model.

### Mixed versions and replacement

Carol imports version two while version one has historical records. A makes
replacement implicit and obscures which implementation was active. B+D retain
both artifact identities, activate one or both explicitly, and let routing
policy select among active compatible promises without rewriting pCID history.

### Long-horizon scale and forks

Many artifacts require rebuildable local indexes rather than canonical mutable
state. B+D allow a local registry projection to be rebuilt from retained CAS
artifacts and local lifecycle events. They also permit forks to coexist as
distinct content-addressed workflow artifacts.

## Conclusions

Reject A and C. The surviving design is B+D: import stores immutable artifact
evidence; explicit activation grants local route/worker eligibility; deactivate
or revoke withdraws that eligibility without deleting artifacts or history;
replacement is additive by CID.

## Decisions still requiring DF

1. Registry persistence: durable local lifecycle records or rebuildable state
   projection from append-only local events.
2. Import source: CAS object only, local package directory plus CAS capture, or
   both.
3. Deactivation terms: one state `inactive`, or separate `deactivated` and
   `revoked` states.
4. Public API names for import, activation, deactivation, and listing.
5. Exact code, state, test, and documentation paths.

## Decision status

needs DF

## Refinements

None.
