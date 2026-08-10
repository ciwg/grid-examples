# Ex4 PromiseGrid alignment plan

TE ID: TE-fugos

## Status

decided

## Decision status

Locked by DI-gisor: use Alternative B, describe `events.jsonl` as durable local
application history, describe built-in identities and role checks as local
application policy, and add an exercise-local testing guide. This TE does not
change Ex4's workflow, storage, identity model, API, or browser/CLI behavior.

## Decision under test

How should the Ex4 Bug Tracker example align with the PromiseGrid development
guide while preserving its useful durable issue workflow and avoiding a false
claim that its current central HTTP server, built-in roles, or event log already
constitute a decentralized PromiseGrid protocol?

## Assumptions and trust model

- The current server owns issue-ID allocation, built-in identity/role checks,
  transition validation, attachment copying, event persistence, and projection.
  Its browser and CLI are two local embodiments of that same server API.
- Current `events.jsonl` is durable local history. It is not signed,
  pCID-selected, independently replayable peer evidence, or shared proof of an
  actor's intent.
- Alice may operate a tracker for one team; Bob may use the browser or CLI;
  Carol may later operate a second tracker; Mallory may submit malformed,
  replayed, or unauthorized HTTP requests.
- PromiseGrid alignment must keep top-level distributed semantics centered on
  voluntary promises and pCID-defined payload meaning. It must not introduce
  new workflow-specific top-level action kinds merely to mirror issue states.
- Ex4 may remain a useful teaching example even if some implementation-local
  mechanics stay local, provided its documentation states the boundary plainly.

## Alternatives

### A. Documentation-first local-workflow alignment

Keep Ex4's current local server-owned issue workflow. Publish an explicit scope
and non-claim statement, distinguish local event history from shared evidence,
document the centralized identity/role boundary, add an exercise testing guide,
and add focused regressions for the documented local durability and rejection
boundaries. Do not add pCIDs, signed envelopes, or remote peer exchange in this
alignment pass.

### B. Bounded PromiseGrid issue-promise layer

Keep the existing application UX but define one or more local-draft pCID specs
for signed issue-related promises and preserve the existing HTTP workflow as a
local adapter/projection. Publish durable accepted and rejected evidence rules,
and scope any cross-node exchange narrowly. The server may validate its local
policy but may not present that policy as universal proof of another actor's
intent.

### C. Full decentralized bug-tracker redesign

Replace built-in identities, central transition enforcement, and the current
event model with independently operating agents, decentralized authorization,
cross-tracker synchronization, and generalized role/delegation semantics.

## Scenario analysis

### S1 — Alice runs one local tracker

Under A, Alice retains the simple browser/CLI workflow and can inspect the
append-only log, while the docs stop calling its local events shared workflow
proof. Under B, Alice gains explicit pCID-selected promise semantics, but must
also understand which changes are accepted protocol artifacts versus local UI
or policy projection. Under C, Alice faces a much larger operational system
before the existing example can teach issue tracking.

### S2 — Bob submits a normal issue update

Under A, the server's built-in role rule is honestly an application-local
authorization check. Under B, Bob can issue a signed promise whose pCID defines
what a report, triage, assignment, or resolution claim means; the server's
chosen projection remains local interpretation rather than a global fact. C
requires a full identity and delegation model before a routine update can be
explained.

### S3 — Mallory submits malformed or unauthorized input

A can deterministically reject the request and document that the resulting
server log is local diagnostic evidence. B additionally requires a clear
accepted/rejected artifact boundary, bounded raw retention, and no conversion
of a local rejection into a universal statement about Mallory. C multiplies
attack surface through authority, synchronization, revocation, and federation
without first proving the basic issue-promise contract.

### S4 — Carol operates a second tracker or a mixed-version node

A makes no interoperability claim: Carol's tracker is another local
application. B can compare pCIDs before exchanging compatible signed promise
artifacts and treats an unsupported pCID as a local lack of support, not global
invalidity. C aspires to federation but must first settle discovery, identity,
authorization, conflict, and evidence semantics.

### S5 — Long-horizon evolution and restore

A preserves current event history but offers no portable issue artifact or
cross-node restore semantics. B creates an explicit seam: later restore,
assignment, and lifecycle policy can refer to selected promise artifacts
without rewriting history, while remaining future work. C risks freezing a
premature universal workflow, role, and permissions model into an example.

### S6 — Scale and operational cost

A has the lowest storage, CPU, bandwidth, and teaching cost. B adds signing,
CAS/evidence handling, pCID documents, verification, and migration obligations
but stays bounded if introduced in slices. C introduces continuing distributed
systems cost disproportionate to this example's current purpose.

## Conclusions

C is rejected for this alignment pass: it would invent generalized identity,
delegation, authorization, and federation policy before Ex4 has an auditable
promise contract.

A and B survive. A is the smallest honest alignment and should be the first
phase whenever the aim is documentation, testing, and accurate boundaries. B
is the more PromiseGrid-native destination, but it should be a separately
scoped implementation phase after its pCID, promise, authority, artifact, and
adapter decisions are explicitly locked.

## Decisions still requiring DF

1. **Alignment scope:** start with A's local-workflow documentation/testing
   alignment, start directly with B's bounded issue-promise layer, or defer
   Ex4 alignment?
2. **Current event wording:** call `events.jsonl` durable local application
   history (recommended), call it relay-local observation/evidence, or present
   it as shared workflow proof?
3. **Identity/role wording:** describe built-in identities and transition rules
   as local application policy (recommended), as temporary PromiseGrid agent
   identities, or as a general authorization model?
4. **Testing guide:** add an Ex4 `docs/testing.md` linked from README
   (recommended), place testing details only in README, or defer the guide?

## Implications for open TODOs and pending DIs

- A new Ex4 alignment TODO should record the selected plan before any code or
  documentation changes.
- A local-workflow alignment must not claim pCIDs, signed envelopes,
  decentralized evidence, or cross-tracker interoperability that Ex4 does not
  implement.
- A future issue-promise layer requires its own TE and DF for profile meaning,
  artifact storage, acceptance/rejection, local projection, and any remote
  exchange.
