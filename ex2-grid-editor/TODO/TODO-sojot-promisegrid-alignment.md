# TODO sojot — Ex2 PromiseGrid alignment

## Decision Intent Log

### DI-busiz

- ID: DI-busiz
- Date: 2026-08-10 09:53:03 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Align Ex2 in documentation-first order: publish the four source-derived local-draft pCIDs and implementation scope; decide malformed and unsupported-envelope evidence policy in a dedicated TE before behavior changes; add focused regression coverage and the required testing guide; then complete a guide-facing documentation pass.
- Intent: Preserve Ex2's already-working signed, pCID-selected, non-canonical CRDT relay mechanics while making its provisional contract, evidence boundaries, and implementation claims explicit and auditable under current PromiseGrid guidance.
- Constraints: Do not represent local drafts as frozen upstream specs or interoperability claims; do not change relay retention, rejection, signing identity, or embodiment semantics without a scoped TE and DI; use the exercise-local `docs/testing.md` guide; preserve existing phase TODO history.
- Affects: `ex2-grid-editor/{README.md,CHANGELOG.md,docs,protocols,service,TODO}`, `TODO/handle-namespace.tsv`, `ex2-grid-editor/docs/thought-experiments/TE-vodik-ex2-promisegrid-alignment-plan.md`

### DI-ralit

- ID: DI-ralit
- Date: 2026-08-10 09:58:56 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Publish Ex2's four current repo-local draft pCIDs through a `CHANGELOG.md` scope declaration and a primary inventory in `docs/architecture.md` linked from the README; state that the relay signing key is the current app identity while browser and Neovim fields are embodiment or presentation data; and explicitly exclude frozen-spec conformance, independent-peer interoperability, general key-rotation/role-continuity policy, and a portable runtime/storage contract.
- Intent: Make the existing pCID-selected editor contract auditable and guide-aligned without turning local draft specs, embodiment labels, or the Docker/runtime layout into stronger PromiseGrid claims than Ex2 can support.
- Constraints: Derive published pCIDs from exact current protocol source bytes; do not alter those source docs in this slice; do not create a normal frozen-spec implementation-promise claim; add regression coverage for publication consistency in `sojot.3`.
- Affects: `ex2-grid-editor/{CHANGELOG.md,README.md,docs/architecture.md,protocols/*.md}`, `ex2-grid-editor/docs/thought-experiments/TE-barap-ex2-local-draft-profile-publication.md`, `ex2-grid-editor/TODO/TODO-sojot-promisegrid-alignment.md`, `TODO/handle-namespace.tsv`

### DI-todav

- ID: DI-todav
- Date: 2026-08-10 10:07:28 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Retain every bounded malformed, invalid-proof, or unsupported-profile inbound byte sequence in relay CAS and append one separate local `observations.jsonl` record per receipt; use `no_supported_handler` for a valid signed envelope that this relay cannot handle; and keep exception observations separate from accepted-message log replay.
- Intent: Preserve exact relay-local evidence for later assessment while keeping the accepted `message-log.jsonl` and its replay projections limited to supported, verified envelopes, without claiming global pCID invalidity or another participant's intent.
- Constraints: Apply existing ingress-size bounds before retention; do not relay rejected bytes; do not use `unknown_pcid`; retain only local observer facts; store observations beside the relay's existing `cas/` and `message-log.jsonl`; do not make observations replay input.
- Affects: `ex2-grid-editor/{cas,service,store}/*.go`, `ex2-grid-editor/{cas,service,store}/*_test.go`, `ex2-grid-editor/docs/{architecture.md,testing.md}`, `ex2-grid-editor/TODO/TODO-sojot-promisegrid-alignment.md`, `ex2-grid-editor/docs/thought-experiments/TE-bulaj-ex2-relay-exception-evidence-policy.md`, `TODO/handle-namespace.tsv`

### DI-lorud

- ID: DI-lorud
- Date: 2026-08-10 10:25:52 -0700
- Status: superseded
- Author: jj@thesalleys.com (JJ)
- Decision: Treat Ex2 as a multi-machine, multi-relay application architecture; its single-machine browser/Neovim and multi-relay local runs are simulations of those logical relay nodes, not a different one-node product model.
- Intent: Keep Ex2's topology and evidence reasoning aligned with distributed PromiseGrid operation while preserving a convenient local development and demonstration mode.
- Constraints: A shared host does not collapse independent relay nodes; each relay keeps its own data root and signing identity; later documentation must call single-host execution a local simulation rather than a separate architecture.
- Affects: `ex2-grid-editor/{README.md,docs/architecture.md,docs/testing.md,service,store}`, `ex2-grid-editor/TODO/TODO-sojot-promisegrid-alignment.md`, `TODO/handle-namespace.tsv`

### DI-nilas

- ID: DI-nilas
- Date: 2026-08-10 10:35:49 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Ex2's canonical decentralized topology is multi-relay peer collaboration: each relay signing key is the current app-node identity and local evidence observer; browser and Neovim are local embodiments until a separate client-signing design exists; and a single host is either a genuine one-node session (one relay/key/data root) or a distributed simulation (multiple relay processes, keys, and data roots).
- Intent: Preserve decentralized PromiseGrid collaboration without turning a shared relay, UI participant field, host name, or display label into a permanent authority or a false cryptographic identity.
- Constraints: Do not represent the relay key as a human identity; do not define key rotation, delegation, or cross-relay trust policy in this decision; document two-relay exchange in later guide/testing work; supersedes DI-lorud.
- Affects: `ex2-grid-editor/{README.md,docs/architecture.md,docs/testing.md,service,store}`, `ex2-grid-editor/docs/thought-experiments/TE-tazoh-ex2-distributed-topology-and-identity.md`, `ex2-grid-editor/TODO/TODO-sojot-promisegrid-alignment.md`, `TODO/handle-namespace.tsv`

## Alignment plan

- [x] sojot.1 Publish a four-profile source-derived pCID inventory and
  local-draft implementation-scope declaration with explicit non-claims.
- [x] sojot.2 Run a TE and DF for malformed and unsupported-envelope retention,
  local observations, and replay behavior before changing relay semantics.
- [ ] sojot.3 Add regression coverage for published pCID inventory, selected
  evidence behavior, raw/CAS replay, and cross-embodiment interoperability.
- [ ] sojot.4 Add `docs/testing.md`, link it from the README, and document the
  verification commands, test layers, and artifact assertions.
- [ ] sojot.5 Complete the final README guide pass for local-draft scope,
  relay-versus-embodiment identity, reproducible operation, and artifact
  inspection.
