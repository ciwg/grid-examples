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

### DI-ninul

- ID: DI-ninul
- Date: 2026-08-10 15:08:40 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Implement three pCID-selected signed promise profiles for issue report, issue lifecycle update, and attachment reference; give each built-in local adapter identity a distinct local signing key; retain accepted artifacts append-only and content-addressably with a rebuildable local projection/index plus bounded local rejected observations/diagnostics; and have browser/CLI adapters submit signed promises only to the local server in the first slice, with no cross-tracker exchange.
- Intent: Make Ex4's first PromiseGrid layer provenance-preserving and evolvable without turning its current server, event log, roles, or HTTP commands into a false decentralized protocol or prematurely adding federation.
- Constraints: The only top-level distributed action is `promise`; report, comment, assignment, status, and attachment meanings are pCID-defined payload semantics. Preserve `events.jsonl` as historical local application data; do not server-sign on behalf of users; do not add generalized identity, delegation, revocation, or remote exchange; resolve profile, schema, path, function, and variable names through later Decision Framing before code edits.
- Affects: `ex4-bug-tracker/{protocols,cas,store,service,web,cmd,docs,TODO}`, `TODO/handle-namespace.tsv`, `ex4-bug-tracker/docs/thought-experiments/TE-gusot-ex4-issue-promise-layer.md`

### DI-muzal

- ID: DI-muzal
- Date: 2026-08-10 15:16:42 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Refine the `DI-ninul` signing-key choice so each browser profile generates and retains its own signing key, each CLI uses an explicit operator-selected `--agent-key` path, each enrolled public key receives an `agent_id` derived from that key, and first use performs explicit local enrollment with a public-key proof and selected local role; the server retains only public enrollment bindings.
- Intent: Ensure browser and CLI embodiments make their own signed promises while preserving the distinction between local `agent_id`, presentation/role metadata, and any future recognized role continuity.
- Constraints: The server must never retain or use a client private key or sign on behalf of a client; enrollment is only Ex4-local bootstrap/admission, not generalized identity, delegation, revocation, or cross-machine trust; a fresh/private browser profile is a fresh embodiment requiring separate enrollment; key rotation and recognized-role continuity require separate TE/DI work.
- Affects: `ex4-bug-tracker/{service,cmd,web,docs,TODO}`, `TODO/handle-namespace.tsv`, `ex4-bug-tracker/docs/thought-experiments/TE-nahup-ex4-local-adapter-signing-identity.md`

### DI-gonok

- ID: DI-gonok
- Date: 2026-08-10 15:19:59 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Reuse Ex3's provisional CBOR envelope shape `grid([42(pCID), payload, proof])`; create the `issue-report`, `issue-lifecycle-update`, and `issue-attachment-reference` profile documents; use `protocol/`, `identity/`, `cas/`, and `store/` packages with `service/` as the local adapter/projection; store accepted objects at `cas/<CID>` with `accepted-promises.jsonl`, `observations.jsonl`, and `agent-bindings.jsonl`; name the write endpoints `POST /api/promises` and `POST /api/agents/enroll` and their service operations `SubmitPromise` and `EnrollAgent`; and move browser/CLI writes to signed-promise submission while retaining read/projection endpoints only.
- Intent: Make Ex4's first signed promise slice structurally consistent with the guide and Ex3 while preserving exact artifact provenance, client key custody, local-only admission, and explicit provisional scope.
- Constraints: These are local-draft profiles, not frozen specs; `events.jsonl` remains historical local application history; the server never holds client private keys or signs for clients; no unsigned mutation endpoints, cross-tracker exchange, generalized authorization, delegation, revocation, or recognized-role continuity are added; every behavior-changing code path must cite this DI.
- Affects: `ex4-bug-tracker/{protocols,protocol,identity,cas,store,service,cmd,web,docs,TODO}`, `TODO/handle-namespace.tsv`

## Alignment plan

- [x] kakon.1 Publish the current local-workflow scope, non-claims, and the
  boundary between durable application history and PromiseGrid evidence.
- [x] kakon.2 Run a TE and DF for the bounded issue-promise profiles, their
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
