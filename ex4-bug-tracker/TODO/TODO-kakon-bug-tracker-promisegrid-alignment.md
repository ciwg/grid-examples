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

### DI-tosoj

- ID: DI-tosoj
- Date: 2026-08-10 15:23:07 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Name the reusable envelope/CID types `protocol.Envelope`, `protocol.Proof`, and `protocol.CIDForBytes`; name local key/binding types `identity.AgentKey`, `identity.AgentID`, and `identity.Enrollment`; name durable stores `store.ArtifactStore`, `store.AcceptedPromiseLog`, `store.ObservationLog`, and `store.AgentBindingLog`; and name profile payloads `protocol.IssueReport`, `protocol.IssueLifecycleUpdate`, and `protocol.IssueAttachmentReference`.
- Intent: Reuse the guide-consistent Ex3 vocabulary for the common envelope while keeping Ex4's local identity, durable artifact, and pCID-owned payload boundaries legible.
- Constraints: These names carry no generalized agent identity, recognized role, authorization, or transport claim; keep `service.SubmitPromise` and `service.EnrollAgent` as previously locked; do not rename historical `IssueEvent` local-history records.
- Affects: `ex4-bug-tracker/{protocol,identity,cas,store,service,TODO}`, `TODO/handle-namespace.tsv`

### DI-nusop

- ID: DI-nusop
- Date: 2026-08-10 15:32:15 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Make enrollment an adapter-submitted canonical-CBOR `identity.Enrollment` claim accompanied by a detached `identity.EnrollmentProof`; `service.EnrollAgent` recomputes `agent_id` from the submitted public key and verifies the proof over the claim before appending the public local binding.
- Intent: Prove local private-key possession at enrollment without server private-key access, server-signing on a client's behalf, or an implied general authorization claim.
- Constraints: The proof is local bootstrap/admission evidence only; no silent enrollment, pre-enrolled private key, recognized-role continuity, delegation, revocation, or cross-machine trust is introduced; invalid enrollment is recorded only as a bounded local observation/diagnostic.
- Affects: `ex4-bug-tracker/{identity,store,service,cmd,web,TODO}`, `TODO/handle-namespace.tsv`

### DI-tofuf

- ID: DI-tofuf
- Date: 2026-08-10 15:34:24 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Carry both `POST /api/promises` signed envelopes and `POST /api/agents/enroll` canonical enrollment-claim/proof requests as raw CBOR request bodies rather than JSON/base64 wrappers.
- Intent: Preserve exact bytes at the signed protocol boundary so adapter HTTP carriage does not become an alternate unsigned representation.
- Constraints: The endpoints remain local adapters only; content type must identify CBOR; the server retains exact accepted envelope bytes in CAS and does not alter a client-signed envelope before verification.
- Affects: `ex4-bug-tracker/{service,cmd,web,TODO}`, `TODO/handle-namespace.tsv`

### DI-basok

- ID: DI-basok
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Limit raw-CBOR enrollment and signed-promise adapter requests to 64 KiB.
- Intent: Bound local adapter ingress while keeping attachment bytes outside promise requests in CAS.
- Constraints: Applies only to the metadata-only profile requests; attachment objects retain their separate existing size policy.

### DI-zubot

- ID: DI-zubot
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Persist each browser profile's non-extractable WebCrypto signing key in IndexedDB.
- Intent: Keep browser private keys embodiment-local and make a private/incognito profile naturally a distinct local agent.
- Constraints: Do not serialize private keys into localStorage or send them to the server; a fresh profile requires separate local enrollment.

### DI-tigid

- ID: DI-tigid
- Status: superseded
- Author: jj@thesalleys.com (JJ)
- Decision: Use a local signer bridge: Go builds canonical signable bytes, the browser's non-extractable IndexedDB WebCrypto key signs those bytes, and Go assembles and verifies the returned exact envelope before the browser submits it unchanged.
- Intent: Preserve browser private-key custody while maintaining one verified canonical-CBOR envelope implementation.
- Constraints: Local-only bridge; server never receives a private key or signs on a client's behalf; no browser CBOR encoder; exact returned envelope bytes must verify before submission.

### DI-pusip

- ID: DI-pusip
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Supersede DI-tigid's same-machine/local-only bridge constraint. The prepare/finalize signer bridge is network-reachable through the normal service transport, while each browser profile retains its own non-extractable private key and the service applies only service-scoped enrollment and acceptance policy.
- Intent: Support decentralized multi-machine browser embodiments without moving client private keys to the service or turning local acceptance into global authority.
- Constraints: Exact signable and finalized envelope bytes remain Go-canonical and are verified before acceptance; network reachability does not add cross-tracker federation, global identity, delegation, revocation, or recognized-role continuity.
- Supersedes: DI-tigid

### DI-gugah

- ID: DI-gugah
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Expose network-reachable raw-CBOR `POST /api/promises/prepare` and `POST /api/promises/finalize` endpoints. Prepare validates service-scoped enrollment and profile payload then returns Go-canonical signable bytes; finalize verifies the browser proof, assembles the exact envelope, verifies it again, and returns its unchanged bytes for `POST /api/promises` submission.
- Intent: Support multi-machine browser embodiments while preserving browser private-key custody and one canonical Go envelope implementation.
- Constraints: Apply the 64 KiB bound; no server private key or signing on behalf of a browser; service enrollment/acceptance remains service-scoped, not global; no cross-tracker federation is implied.

### DI-mofab

- ID: DI-mofab
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Name the local signer-bridge service methods `PreparePromise` and `FinalizePromise`.
- Intent: Describe canonical-byte construction and proof verification as adapter mechanics while keeping `promise` as the only protocol-level semantic action.
- Constraints: These names do not introduce new wire actions or change the pCID-owned payload meanings.

### DI-fadog

- ID: DI-fadog
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Name the raw-CBOR prepare request `PromiseDraft` and the finalize request `PromiseProof`.
- Intent: Keep the network adapter contract explicit without introducing a new protocol action or obscuring the signer-owned proof boundary.
- Constraints: These are local bridge request types; pCID-selected payload semantics and final envelope proof remain protocol-owned.

### DI-rutul

- ID: DI-rutul
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Every pCID crossing the signer bridge uses CBOR tag 42 with binary CID bytes; printable CID text is limited to documentation, logs, and diagnostics.
- Intent: Keep prepare/finalize selector bytes identical to the final Ex4 grid envelope.
- Constraints: Supersedes any under-specified text-CID bridge interpretation; no bridge request may use a text CID as its protocol selector.

### DI-kolaf

- ID: DI-kolaf
- Date: 2026-08-10 16:20:00 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Supersede raw-CBOR carriage for the prepare/finalize adapter requests with bounded JSON requests that carry only a local profile name and profile fields. The service resolves the profile's embedded source to its pCID, creates canonical tag-42 signable bytes, verifies a client proof, and returns the unchanged final raw-CBOR envelope for `/api/promises`.
- Intent: Remove the unsafe second browser CBOR implementation while retaining browser/CLI private-key custody, multi-machine reachability, and the exact PromiseGrid wire boundary.
- Constraints: JSON is local adapter carriage only, never a promise artifact; neither bridge request carries a pCID; final signed envelopes remain raw canonical CBOR and `grid([42(pCID), payload, proof])`; drafts are bounded, short-lived, and service-scoped; no private key reaches the service.
- Affects: `ex4-bug-tracker/{service,identity,cmd,web,docs,TODO}`
- Supersedes: DI-tofuf (prepare/finalize request carriage only); DI-gugah (prepare/finalize request carriage only)

## Alignment plan

- [x] kakon.1 Publish the current local-workflow scope, non-claims, and the
  boundary between durable application history and PromiseGrid evidence.
- [x] kakon.2 Run a TE and DF for the bounded issue-promise profiles, their
  pCID-selected meaning, accepted/rejected artifacts, and adapter projection.
- [x] kakon.3 Implement the locked issue-promise layer with local durable
  records, bounded rejection handling, and regression coverage.
- [ ] kakon.4 Add `docs/testing.md`, link it from README, and document local
  workflow checks separately from promise/evidence and interoperability checks.
- [ ] kakon.5 Complete the final README and architecture alignment pass.

Status: active. `kakon.3` implements the bounded local-draft promise layer;
the remaining documentation and final alignment checks are tracked below.

`kakon.1` completed with `CHANGELOG.md`, `docs/architecture.md`, and README
scope/navigation updates. Source: `DI-nibuh`.
