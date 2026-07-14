# TODO ravup - grid-editor phase 4 publish and import exchange

## Decision Intent Log

ID: DI-tavul
Date: 2026-07-13 16:45:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement the first Phase 4 slice as relay-signed publish plus import/exchange, not publish-only, so Grid Editor can both emit and ingest durable document handoff artifacts.
Intent: Cover the full outbound and inbound movement of document data in one slice instead of shipping a one-way export that does not prove the round trip.
Constraints: The live CRDT relay remains unchanged; publish/import is a separate current-time action; permissions, owner/admin, and restore semantics stay out of scope in this slice.
Affects: `ex2-grid-editor/service`, `ex2-grid-editor/protocols`, `ex2-grid-editor/web`, `ex2-grid-editor/docs`

ID: DI-gosaf
Date: 2026-07-13 16:45:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Publish may target either the current document state or a named saved version, and the publish/import protocol will live in a new `protocols/publish-document.md` document.
Intent: Support both ad hoc sharing and stable version publication while keeping publish/import semantics separate from `live-document`.
Constraints: Saved versions remain local browser review objects in this slice; the relay publishes a signed manifest that references CAS-backed bytes; import/exchange may materialize a new local document from the published artifact without changing past live history.
Affects: `ex2-grid-editor/protocols/publish-document.md`, `ex2-grid-editor/service`, `ex2-grid-editor/web`, `ex2-grid-editor/README.md`

Goal: Add the first PromiseGrid-native publish/import slice with relay-signed manifests, CAS-backed exchange objects, and browser publish/import actions.

- [x] ravup.1 Add the publish-document protocol doc and relay manifest types.
- [x] ravup.2 Add relay publish/import endpoints, CAS-backed artifact resolution, and replay indexing.
- [x] ravup.3 Add browser publish/import UI and current-or-saved-version selection.
- [x] ravup.4 Add tests and docs for the publish/import exchange slice.
