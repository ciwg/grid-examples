# TODO fozoz — Ex3 PromiseGrid alignment

## Decision Intent Log

### DI-figak

- ID: DI-figak
- Date: 2026-08-10 11:44:27 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Align Ex3 in documentation-first order: publish source-derived local-draft pCID and provisional remote-admission scope/non-claims; audit remote capability, WebSocket-carriage, and relay-local rejected-ingress evidence boundaries; add focused regression coverage while retaining existing cross-embodiment and multi-relay tests; create an exercise-local testing guide; then finish a reader-facing README pass that keeps private/incognito manual verification explicitly open in TODO tamuk.
- Intent: Make Ex3's current decentralized multi-relay demonstration auditable and guide-aligned without overstating its bootstrap-token/capability mechanism as a universal PromiseGrid auth API or converting automated browser hardening into a completed private-browser environmental guarantee.
- Constraints: WebSocket remains carriage rather than protocol meaning; preserve pCID-selected profile semantics; retain TODO tamuk until real manual private/incognito verification occurs; do not redesign generalized identity, delegation, or authorization during this alignment plan.
- Affects: `ex3-grid-editor-websocket/{README.md,CHANGELOG.md,docs,protocols,service,web,cmd,TODO}`, `TODO/handle-namespace.tsv`, `ex3-grid-editor-websocket/docs/thought-experiments/TE-hovij-ex3-promisegrid-alignment-plan.md`

### DI-hadil

- ID: DI-hadil
- Date: 2026-08-10 12:57:29 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Publish Ex3's four exact source-derived local-draft pCIDs through a `CHANGELOG.md` scope and non-claim declaration plus a primary inventory in `docs/architecture.md`, linked from the README; classify bootstrap-secret and short-lived mutation capabilities as relay-local implementation admission mechanics; state that the relay key is current app-node identity while capability audiences and embodiment fields are not general person identity, delegation, or role continuity; and disclose that TODO tamuk's manual private/incognito verification remains open despite automated hardening.
- Intent: Make the current Ex3 contract auditable without promoting WebSocket carriage or relay-local remote admission into a fifth public profile, a frozen PromiseGrid auth API, or a completed cross-environment browser guarantee.
- Constraints: Derive published pCIDs from exact current protocol-source bytes; do not alter profile documents or remote-admission behavior in this slice; do not create a normal frozen-spec implementation-promise claim; preserve TODO tamuk as open; add source-derived regression coverage in `fozoz.3`.
- Affects: `ex3-grid-editor-websocket/{CHANGELOG.md,README.md,docs/architecture.md,protocols/*.md}`, `ex3-grid-editor-websocket/docs/thought-experiments/TE-hujos-ex3-local-draft-remote-admission-scope.md`, `ex3-grid-editor-websocket/TODO/TODO-fozoz-promisegrid-alignment.md`, `TODO/handle-namespace.tsv`

### DI-pazis

- ID: DI-pazis
- Date: 2026-08-10 13:04:37 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: For bounded rejected peer-envelope ingress, retain exact bytes in CAS and append one relay-local observation per receipt while excluding observations from accepted log, replay, and peer feed; complete pCID support and payload validation before accepted CAS/log/replay mutation; retain remote capability denials only as non-secret local diagnostics; and return/close WebSocket framing or JSON failures without durable raw-frame evidence.
- Intent: Preserve useful relay-local evidence and strict rebuildable accepted replay without exposing bearer/bootstrap material, inventing WebSocket protocol semantics, or turning local observations into global validity or intent claims.
- Constraints: Apply ingress-size bounds before raw retention; do not retain capability tokens or bootstrap secrets; do not relay rejected bytes; classify a valid but unsupported pCID as this relay's `no_supported_handler`, not a global invalidity claim; keep capability diagnostics separate from peer-envelope observations.
- Affects: `ex3-grid-editor-websocket/{cas,service,store}/*.go`, `ex3-grid-editor-websocket/{cas,service,store}/*_test.go`, `ex3-grid-editor-websocket/docs/{architecture.md,testing.md}`, `ex3-grid-editor-websocket/docs/thought-experiments/TE-sozol-ex3-remote-admission-and-ingress-evidence-policy.md`, `ex3-grid-editor-websocket/TODO/TODO-fozoz-promisegrid-alignment.md`, `TODO/handle-namespace.tsv`

### DI-darif

- ID: DI-darif
- Date: 2026-08-10 13:08:32 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Implement rejected peer-envelope records as `store.Observation` and `store.ObservationLog` in `observations.jsonl`; implement non-secret remote capability denial records as `store.AdmissionDiagnostic` and `store.AdmissionDiagnosticLog` in separate `admission-diagnostics.jsonl`; and refactor `ingestEnvelopeLocked` to fully classify envelope support and payload before choosing either accepted CAS/log/replay persistence or the observation path.
- Intent: Name records for what the relay actually knows, keep admission events distinct from peer-envelope observations, and structurally prevent rejected inputs from entering accepted replay.
- Constraints: Preserve DI-pazis's raw-retention, secret-exclusion, and WebSocket rules; do not represent a sender as an authenticated peer before accepted envelope verification; do not add a top-level protocol action; supersedes no prior DI.
- Affects: `ex3-grid-editor-websocket/{service,store}/*.go`, `ex3-grid-editor-websocket/{service,store}/*_test.go`, `ex3-grid-editor-websocket/TODO/TODO-fozoz-promisegrid-alignment.md`, `TODO/handle-namespace.tsv`

### DI-lubij

- ID: DI-lubij
- Date: 2026-08-10 13:13:42 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Name the no-side-effect ingress preflight helper `validateEnvelopeLocked`, the local rejected-envelope recorder `recordRejectedEnvelopeLocked`, and the relay evidence fields `observations` and `admissionDiagnostics`.
- Intent: Keep names centered on validation and relay-local observation rather than prematurely identifying a sender as a peer or leaking file-format mechanics into behavior vocabulary.
- Constraints: These helpers must preserve DI-pazis's accepted/rejected separation and DI-darif's record names and paths; do not add a new protocol action or sender-identity claim.
- Affects: `ex3-grid-editor-websocket/service/app.go`, `ex3-grid-editor-websocket/TODO/TODO-fozoz-promisegrid-alignment.md`, `TODO/handle-namespace.tsv`

### DI-lozut

- ID: DI-lozut
- Date: 2026-08-10 13:22:21 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Name the relay-local, non-secret remote-admission recording method `RecordAdmissionDiagnostic(transport, reason string)`.
- Intent: Record what the relay observed at its admission boundary without characterizing a requester, capability, or remote participant as globally failed or invalid.
- Constraints: The method must never accept or persist bearer capabilities, bootstrap secrets, or raw WebSocket frames; it writes only to the DI-darif admission-diagnostic stream and remains outside accepted replay.
- Affects: `ex3-grid-editor-websocket/{service/app.go,service/server.go,service/live_socket.go}`, `ex3-grid-editor-websocket/TODO/TODO-fozoz-promisegrid-alignment.md`, `TODO/handle-namespace.tsv`

### DI-dilav

- ID: DI-dilav
- Date: 2026-08-10 13:39:50 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Add focused regression coverage that derives each Ex3 local-draft pCID from exact `protocols/*.md` source bytes and compares both `docs/architecture.md` and `CHANGELOG.md`; directly prove raw CAS retention, one observation per rejected receipt, accepted replay/feed exclusion, and non-secret admission diagnostics; and retain existing browser/Neovim, peer-relay, WebSocket, headless late-join, browser-JavaScript, and sidecar suites as decentralized coverage without duplicating topology tests.
- Intent: Make Ex3's published and relay-local claims deterministic and auditable while preserving existing cross-embodiment evidence as the correct level for decentralized behavior.
- Constraints: Do not create a duplicate topology suite; do not replace exact-source comparisons with runtime constants; do not treat automated tests as closure of TODO tamuk's manual private-browser requirement; use `fozoz.4` to document the resulting test layers.
- Affects: `ex3-grid-editor-websocket/{protocol,service,store}/*_test.go`, `ex3-grid-editor-websocket/{CHANGELOG.md,docs/architecture.md,docs/testing.md}`, `ex3-grid-editor-websocket/docs/thought-experiments/TE-nadut-ex3-published-claim-regression-coverage.md`, `ex3-grid-editor-websocket/TODO/TODO-fozoz-promisegrid-alignment.md`, `TODO/handle-namespace.tsv`

### DI-hosit

- ID: DI-hosit
- Date: 2026-08-10 13:52:23 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Create one concise Ex3 `docs/testing.md` linked from README; document grouped Go (`go vet ./...`, `go test ./...`, `errcheck ./...`), browser (`npm test`, `npm run build`), and sidecar (`npm run build`) commands; distinguish focused pCID/evidence/admission tests from existing interoperability, headless, browser-JavaScript, sidecar, and WebSocket suites; and explain both `t.TempDir()` isolation and TODO tamuk's still-open manual private-window verification.
- Intent: Give Ex3 contributors reproducible verification and precise confidence boundaries without treating any local suite, temporary relay root, or automated private-browser hardening as global trust proof or closure of the outstanding manual condition.
- Constraints: Keep detailed testing semantics out of the main README; do not add or change runtime behavior; retain TODO tamuk as open; preserve architecture and scope documents as their respective sources of truth.
- Affects: `ex3-grid-editor-websocket/{README.md,docs/testing.md}`, `ex3-grid-editor-websocket/TODO/TODO-fozoz-promisegrid-alignment.md`, `ex3-grid-editor-websocket/docs/thought-experiments/TE-dohuf-ex3-testing-guide.md`, `TODO/handle-namespace.tsv`

## Alignment plan

- [x] fozoz.1 Publish source-derived local-draft pCID inventory and provisional
  remote-admission scope declaration with explicit non-claims.
- [x] fozoz.2 Run a TE and DF for remote capability, WebSocket-carriage, and
  relay-local rejected-ingress evidence policy before behavior changes.
- [x] fozoz.3 Add focused regression coverage for published pCIDs, remote
  admission/evidence boundaries, and existing decentralized interoperability.
- [x] fozoz.4 Add `docs/testing.md`, link it from the README, and document Go,
  browser JavaScript, sidecar, and cross-node verification layers.
- [ ] fozoz.5 Complete the final README guide pass for local-draft scope,
  provisional remote admission, reproducible multi-relay operation, evidence
  inspection, and the still-open manual private-browser condition.
