# TODO kakon - Ex4 PromiseGrid alignment

## Decision Intent Log

### DI-gisor

- ID: DI-gisor
- Date: 2026-08-10 14:46:29 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Align Ex4 by adding a bounded pCID-defined, signed issue-promise layer while retaining the current browser and CLI workflow as local adapters/projections; describe `events.jsonl` as durable local application history; describe built-in identities and role checks as local application policy; and add an exercise-local `docs/testing.md` linked from README.
- Intent: Make Ex4 a genuine PromiseGrid example without falsely treating its current central HTTP server, local role enforcement, or unsigned event log as decentralized shared proof or a universal authorization system.
- Constraints: Preserve the existing usable workflow while its adapter/projection boundary is made explicit; do not add workflow-specific top-level protocol actions; do not claim generalized identity, delegation, revocation, cross-tracker recognition, or federation; require a separate TE and DF before selecting issue-promise pCIDs, artifact/evidence rules, adapter behavior, or remote exchange.
- Affects: `ex4-bug-tracker/{README.md,docs,TODO,protocols,service,store,web,cmd}`, `TODO/handle-namespace.tsv`, `ex4-bug-tracker/docs/thought-experiments/TE-fugos-ex4-promisegrid-alignment-plan.md`

### DI-nibuh

- ID: DI-nibuh
- Date: 2026-08-10 14:54:48 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Publish a concise `CHANGELOG.md` implementation-scope and non-claim declaration, link it and the architecture from README, classify `events.jsonl` as durable local application history, and state that pCID-selected signed issue-promise artifacts are a separately specified future layer with no current interoperability, agent identity, delegation, or role-continuity claim.
- Intent: Give Ex4 an honest PromiseGrid boundary before issue-promise behavior exists, so durable local server events cannot be mistaken for signed promises, shared evidence, or a general authorization system.
- Constraints: Do not modify issue behavior, persistence, role enforcement, event records, or transport in this documentation slice; do not present future artifacts as active; preserve `events.jsonl` as historical local application data; require `kakon.2` TE/DF before selecting pCIDs, signing, artifact acceptance, rejection, adapter projection, or remote exchange.
- Affects: `ex4-bug-tracker/{CHANGELOG.md,README.md,docs/architecture.md,TODO}`, `TODO/handle-namespace.tsv`, `ex4-bug-tracker/docs/thought-experiments/TE-mofol-ex4-local-history-evidence-scope.md`

## Alignment plan

- [x] kakon.1 Publish the current local-workflow scope, non-claims, and the
  boundary between durable application history and PromiseGrid evidence.
- [ ] kakon.2 Run a TE and DF for the bounded issue-promise profiles, their
  pCID-selected meaning, accepted/rejected artifacts, and adapter projection.
- [ ] kakon.3 Implement the locked issue-promise layer with local durable
  records, bounded rejection handling, and regression coverage.
- [ ] kakon.4 Add `docs/testing.md`, link it from README, and document local
  workflow checks separately from promise/evidence and interoperability checks.
- [ ] kakon.5 Complete the final README and architecture alignment pass.

Status: active. The existing application workflow remains usable, but no
PromiseGrid protocol behavior is authorized until `kakon.2` completes its TE
and Decision Framing.

`kakon.1` completed with `CHANGELOG.md`, `docs/architecture.md`, and README
scope/navigation updates. Source: `DI-nibuh`.
