# TODO pisul - grid-editor restore published version

## Decision Intent Log

### DI-hihok

- ID: DI-hihok
- Date: 2026-08-14 19:29:00 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Restore only from an exact resolved `publish-document` manifest CID. An admitted participant's restore creates a new, append-only, CRDT-merged working state and a dedicated repo-local-draft `restore-published-version` promise; it never rewrites live history, treats the historical publication as current authority, or claims byte-identical output in the presence of concurrent edits. Existing mutation-capability admission authorizes requests, and success requires durable coupling of the restore promise and its live CRDT change.
- Intent: Preserve immutable publication provenance while allowing the existing collaborative document to continue through an explicit current-time, pCID-defined promise instead of centralized replacement or import-by-another-name.
- Constraints: Import remains a new-document workflow; restore source is a verified manifest and its referenced CAS bytes; the UI must disclose merge-derived output and source provenance; no owner/admin/delegation/recognized-role-continuity policy is introduced; bounded failure evidence is required when validation or durable coupling fails; the exact implementation names and paths remain subject to the required follow-on DF before code edits.
- Affects: `ex3-grid-editor-websocket/{protocols,service,web,docs,TODO}`, `TODO/handle-namespace.tsv`, `ex3-grid-editor-websocket/docs/thought-experiments/TE-bazav-ex3-restore-published-version.md`

### DI-tibum

- ID: DI-tibum
- Date: 2026-08-14 19:30:08 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Encode the selected publish-manifest provenance and the exact CRDT change that continues the existing document in one dedicated `restore-published-version` promise artifact, rather than append a normal live edit and a separately linked restore record.
- Intent: Make every successful restore replayable as one atomic, pCID-defined append-only fact, with no crash window in which later readers must guess whether a live change and a restoration claim belong together.
- Constraints: The artifact must name the exact source manifest and contain or canonically identify the exact applied CRDT change; replay must reconstruct the same live-document effect from that one accepted artifact; no authoritative snapshot replacement, hidden history rewrite, or browser-local provenance record is permitted; exact field/function/path names remain subject to the required follow-on DF before code edits.
- Affects: `ex3-grid-editor-websocket/{protocols,service,web,docs,TODO}`, `TODO/handle-namespace.tsv`, `ex3-grid-editor-websocket/docs/thought-experiments/TE-bazav-ex3-restore-published-version.md`

### DI-jomij

- ID: DI-jomij
- Date: 2026-08-14 19:31:15 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Create the permanent repo-local-draft restore protocol source at `ex3-grid-editor-websocket/protocols/restore-published-version.md`.
- Intent: Give the dedicated restore promise an exact, content-addressed source document rather than silently extending the distinct publish handoff profile.
- Constraints: The file defines a repo-local draft only; its derived pCID must be surfaced alongside Ex3's current profile inventory; it must retain the existing `grid([42(pCID), payload, proof])` envelope boundary and the DI-hihok/DI-tibum append-only constraints.
- Affects: `ex3-grid-editor-websocket/protocols/restore-published-version.md`, `ex3-grid-editor-websocket/protocols/docs.go`, `ex3-grid-editor-websocket/{service,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-bahit

- ID: DI-bahit
- Date: 2026-08-14 19:31:39 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Create `ex3-grid-editor-websocket/restore/types.go` as the permanent production-code home for the dedicated restore promise payload and its replayable record projection.
- Intent: Keep pCID-owned restore semantics separate from publish handoff types, matching Ex3's established protocol-package boundary.
- Constraints: The package does not define authority, import, or snapshot-replacement behavior; it represents only the DI-hihok/DI-tibum atomic restore artifact and must be replayable from the append-only message log.
- Affects: `ex3-grid-editor-websocket/restore/types.go`, `ex3-grid-editor-websocket/{service,protocols,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-lihud

- ID: DI-lihud
- Date: 2026-08-14 19:32:04 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Name the service operation `service.App.RestorePublishedVersion` and expose it only through `POST /api/local/documents/{document_id}/restore-published-version`.
- Intent: Make the operation visibly distinct from import and ordinary live editing while preserving the exact selected published-version provenance in its name.
- Constraints: The endpoint accepts the one atomic restore artifact selected by DI-tibum; it remains behind existing mutation-capability admission; it must not expose snapshot replacement, a generic restore endpoint, or an owner/admin control plane.
- Affects: `ex3-grid-editor-websocket/{service/app.go,service/server.go,service/*_test.go,web/src/main.js,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-lukis

- ID: DI-lukis
- Date: 2026-08-14 19:32:22 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Extend the existing browser sources `web/src/automerge-relay.js` and `web/src/main.js`, then rebuild the existing served `web/app.js` bundle; do not create a parallel restore browser module.
- Intent: Reuse the same browser-local Automerge machinery that produces ordinary live-document changes, so the atomic restore artifact contains exact compatible CRDT bytes rather than a second browser representation.
- Constraints: The browser retains local construction of the CRDT change but submits one DI-tibum atomic artifact to the DI-lihud endpoint; generated `web/app.js` remains a checked-in build artifact of those source changes; no browser-local restore provenance store is introduced.
- Affects: `ex3-grid-editor-websocket/{web/src/automerge-relay.js,web/src/main.js,web/app.js,web/*_test.go,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-nogat

- ID: DI-nogat
- Date: 2026-08-14 19:32:42 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Name the restore payload fields `source_manifest_cid` and `live_change_base64`; name the browser helpers `buildReplacementChange` and `restorePublishedVersion`.
- Intent: Keep immutable source provenance, exact CRDT carriage, browser-side construction, and the user-directed operation distinct and legible across the pCID-defined artifact, server, and UI.
- Constraints: `source_manifest_cid` identifies only a resolved publish manifest; `live_change_base64` carries only the exact CRDT change the server validates and replays; these names do not create a generic source/change or snapshot-replacement API.
- Affects: `ex3-grid-editor-websocket/{restore/types.go,service/app.go,service/server.go,web/src/automerge-relay.js,web/src/main.js,web/app.js,*_test.go,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-sabud

- ID: DI-sabud
- Date: 2026-08-14 19:33:08 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/service/app.go` as the permanent production-code home for `App.RestorePublishedVersion`, restore pCID registration, replay dispatch, and atomic append handling.
- Intent: Keep accepted protocol state, exact-byte persistence, and replay semantics in Ex3's existing canonical application owner rather than split restore behavior into a second service-state path.
- Constraints: The method must enforce DI-hihok through DI-nogat, retain the same append-only/CAS model as other profiles, and leave import, websocket carriage, and snapshot replacement outside the operation.
- Affects: `ex3-grid-editor-websocket/service/app.go`, `ex3-grid-editor-websocket/{restore,protocols,service/*_test.go,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-vugim

- ID: DI-vugim
- Date: 2026-08-14 19:33:22 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/service/server.go` as the permanent production-code home for the capability-gated `POST /api/local/documents/{document_id}/restore-published-version` adapter.
- Intent: Keep HTTP parsing, request bounds, local/capability admission, and response handling at Ex3's established service adapter boundary while leaving restore semantics in `App.RestorePublishedVersion`.
- Constraints: The adapter must accept only the DI-tibum atomic artifact, reuse existing mutation-capability admission, and must not add a generic restore, remote authority, or snapshot-replacement endpoint.
- Affects: `ex3-grid-editor-websocket/service/server.go`, `ex3-grid-editor-websocket/{service/app.go,service/server_test.go,web/src/main.js,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-dumah

- ID: DI-dumah
- Date: 2026-08-14 19:33:34 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/protocols/docs.go` to register `restore-published-version.md` as an embedded protocol source for pCID derivation.
- Intent: Ensure the restore profile's identity is derived from its exact checked-in source bytes through Ex3's existing content-addressed registry rather than a second registry mechanism.
- Constraints: No pCID is handwritten; the runtime must expose the derived value with its other profile identities; changes to the source document intentionally change the local-draft profile identity.
- Affects: `ex3-grid-editor-websocket/{protocols/docs.go,protocols/restore-published-version.md,service/app.go,protocol/*_test.go,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-pazad

- ID: DI-pazad
- Date: 2026-08-14 19:33:47 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/web/src/main.js` as the permanent browser-source home for the explicit `restorePublishedVersion` user action in the existing published-version review surface.
- Intent: Keep restoration visibly user-directed and adjacent to the immutable published artifact it cites, rather than hiding a state-changing action behind a separate browser subsystem.
- Constraints: The UI must identify the selected source manifest and disclose merge-derived results; it must submit only the DI-tibum atomic artifact through the DI-lihud endpoint; import remains a distinct new-document action.
- Affects: `ex3-grid-editor-websocket/{web/src/main.js,web/app.js,web/src/automerge-relay.js,service/server.go,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-fijab

- ID: DI-fijab
- Date: 2026-08-14 19:34:24 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/service/app_test.go` for deterministic restore source validation, atomic artifact replay, restart recovery, and concurrent-merge provenance coverage.
- Intent: Prove the DI-hihok/DI-tibum semantic and durability guarantees at Ex3's existing application-state test boundary.
- Constraints: Tests must distinguish a restored-derived working state from byte-identical historical replacement; they must not claim general authority or interoperability beyond Ex3's repo-local-draft profiles.
- Affects: `ex3-grid-editor-websocket/{service/app_test.go,service/app.go,restore/types.go,protocols/restore-published-version.md,docs/testing.md,TODO}`, `TODO/handle-namespace.tsv`

### DI-fofal

- ID: DI-fofal
- Date: 2026-08-14 19:35:29 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/service/server_test.go` for restore endpoint capability denial, malformed manifest/change rejection, and successful atomic request carriage coverage.
- Intent: Prove the HTTP adapter preserves the narrow DI-vugim admission and exact-artifact boundary rather than becoming an alternate authority or restore mechanism.
- Constraints: Tests must cover existing local/capability admission, must reject malformed inputs before state change, and must not rely on external network services.
- Affects: `ex3-grid-editor-websocket/{service/server_test.go,service/server.go,service/app.go,docs/testing.md,TODO}`, `TODO/handle-namespace.tsv`

### DI-sipuj

- ID: DI-sipuj
- Date: 2026-08-14 19:35:56 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/web/src/automerge-relay.test.mjs` to prove `buildReplacementChange` creates the exact CRDT change carried by the atomic restore artifact.
- Intent: Keep browser-side restore construction deterministic and directly tied to the bytes the service validates and replays.
- Constraints: The test covers only embodiment-local CRDT construction; server acceptance, provenance, and capability behavior remain covered by the approved Go test paths.
- Affects: `ex3-grid-editor-websocket/{web/src/automerge-relay.test.mjs,web/src/automerge-relay.js,web/src/main.js,web/app.js,docs/testing.md,TODO}`, `TODO/handle-namespace.tsv`

### DI-posam

- ID: DI-posam
- Date: 2026-08-14 19:36:13 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/docs/testing.md` to document restore verification, manifest/CRDT provenance assertions, and the repo-local-draft trust boundary.
- Intent: Give contributors a reproducible account of what restore tests establish without presenting a local relay's evidence as centralized authority or frozen interoperability proof.
- Constraints: The guide must distinguish publication, import, and restore; it must state local temporary-root/evidence limits; it must link test layers to DI-hihok through DI-tibum.
- Affects: `ex3-grid-editor-websocket/{docs/testing.md,service/app_test.go,service/server_test.go,web/src/automerge-relay.test.mjs,README.md,TODO}`, `TODO/handle-namespace.tsv`

### DI-zapud

- ID: DI-zapud
- Date: 2026-08-14 19:36:29 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/README.md` to explain that restore continues the existing document through a provenance-bearing CRDT-merged change, whereas import creates a new local document.
- Intent: Keep the primary reader path honest about restore's append-only, non-authoritative semantics and prevent users from treating it as a destructive rollback or a rename of import.
- Constraints: README wording must preserve the repo-local-draft/non-frozen scope and point detailed verification to `docs/testing.md`; it must not imply owner/admin authority or byte-identical output under concurrent edits.
- Affects: `ex3-grid-editor-websocket/{README.md,docs/testing.md,docs/architecture.md,web/src/main.js,TODO}`, `TODO/handle-namespace.tsv`

### DI-gadun

- ID: DI-gadun
- Date: 2026-08-14 19:36:56 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/docs/architecture.md` to add the restore profile's derived pCID inventory entry and one-artifact replay boundary.
- Intent: Keep the primary technical architecture record aligned with the exact local-draft profile source and make the restore durability model inspectable without relying on source-code inference.
- Constraints: The document must identify the profile as repo-local-draft, preserve the existing non-claims, and distinguish immutable publish provenance from the new current-time restore promise.
- Affects: `ex3-grid-editor-websocket/{docs/architecture.md,protocols/restore-published-version.md,protocol/profile_inventory_test.go,README.md,docs/testing.md,TODO}`, `TODO/handle-namespace.tsv`

### DI-puman

- ID: DI-puman
- Date: 2026-08-14 19:37:09 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/protocol/profile_inventory_test.go` to derive the restore profile pCID from its source bytes and require its publication in both the architecture inventory and implementation scope.
- Intent: Make profile-publication drift fail deterministically at the reader-facing claim boundary rather than leaving the new restore pCID undocumented or hand-copied.
- Constraints: The test must derive—not hard-code—the pCID; it extends only Ex3's repo-local-draft inventory and does not assert frozen-spec conformance or external interoperability.
- Affects: `ex3-grid-editor-websocket/{protocol/profile_inventory_test.go,protocols/restore-published-version.md,docs/architecture.md,CHANGELOG.md,TODO}`, `TODO/handle-namespace.tsv`

### DI-mitus

- ID: DI-mitus
- Date: 2026-08-14 19:37:35 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/CHANGELOG.md` to publish the source-derived repo-local-draft restore pCID and its explicit non-claims.
- Intent: Keep Ex3's implementation promise honest and auditable by disclosing exactly what the new profile claims before readers or tests treat it as a stable interoperability commitment.
- Constraints: The declaration must state that the profile is repo-local-draft, excludes frozen-spec conformance and generalized authority, and preserves the append-only, capability-gated semantics locked by DI-hihok through DI-tibum.
- Affects: `ex3-grid-editor-websocket/{CHANGELOG.md,protocol/profile_inventory_test.go,protocols/restore-published-version.md,docs/architecture.md,README.md,TODO}`, `TODO/handle-namespace.tsv`

### DI-nihiz

- ID: DI-nihiz
- Date: 2026-08-14 19:39:23 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Issue a dedicated `restore-published-version` mutation capability through Ex3's existing session/bootstrap mechanism and require that capability for restore requests.
- Intent: Scope admission to the exact new pCID-defined promise while preserving the established capability mechanism and avoiding any owner/admin/recognized-role authority claim.
- Constraints: The restore capability is a short-lived relay-local admission mechanism, not a person identity, delegation, or authorization standard; restore remains capability-gated and does not reuse an unrelated live-document capability.
- Affects: `ex3-grid-editor-websocket/{service/bootstrap.go,service/app.go,service/server.go,service/*_test.go,web/src/main.js,docs,TODO}`, `TODO/handle-namespace.tsv`

### DI-dosid

- ID: DI-dosid
- Date: 2026-08-14 19:39:53 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Modify `ex3-grid-editor-websocket/service/bootstrap.go` to issue the dedicated restore pCID capability through the existing session response.
- Intent: Reuse Ex3's narrow, short-lived capability mechanism while giving the new pCID-defined restore promise an exact audience rather than borrowing another profile's admission token.
- Constraints: No new authority registry, role claim, or alternate bootstrap channel; the capability remains relay-local provisional admission and must be covered by existing-session regression paths.
- Affects: `ex3-grid-editor-websocket/{service/bootstrap.go,service/*_test.go,service/app.go,service/server.go,web/src/main.js,docs,TODO}`, `TODO/handle-namespace.tsv`

Goal: Define and implement a user-directed workflow that makes a selected
earlier published document version the current working document without
rewriting existing history.

- [x] pisul.1 Run a thought experiment for provenance, conflict handling,
  authority boundaries, and the distinction between restore and import in
  `docs/thought-experiments/TE-bazav-ex3-restore-published-version.md`.
- [x] pisul.2 Lock the user-visible restore semantics, artifact references,
  behavior names, and storage paths through Decision Framing and DIs
  `DI-hihok`, `DI-tibum`, `DI-jomij`, `DI-bahit`, `DI-lihud`, `DI-lukis`,
  `DI-nogat`, `DI-sabud`, `DI-vugim`, `DI-dumah`, `DI-pazad`, `DI-fijab`,
  `DI-fofal`, `DI-sipuj`, `DI-posam`, `DI-zapud`, `DI-gadun`, `DI-puman`,
  and `DI-mitus`.
- [x] pisul.3 Implement the selected atomic restore workflow with append-only
  evidence and regression coverage.
- [x] pisul.4 Document the workflow, its trust boundary, and verification.

Status: complete. TE-bazav and the Decision Intent Log lock and the approved
implementation verifies an append-only, capability-gated atomic restore
promise.
