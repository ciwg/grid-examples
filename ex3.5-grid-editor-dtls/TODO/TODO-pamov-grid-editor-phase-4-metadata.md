# TODO pamov - grid-editor phase 4 document metadata

## Decision Intent Log

ID: DI-loruk
Date: 2026-07-13 18:35:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement the next Phase 4 slice as relay-signed document metadata instead of permissions or restore, covering description, summary, tags, collections, favorites, archive state, and document search over relay-known metadata.
Intent: Move the most immediately useful document-management features into a PromiseGrid-shaped backend slice without waiting for the more sensitive permissions or restore semantics.
Constraints: Live CRDT editing remains unchanged; publish/import remains separate; title and other older Phase 2 browser-local workflow fields may still exist locally, but Phase 4 metadata must have a relay-owned signed protocol and latest-state semantics.
Affects: `ex2-grid-editor/protocols`, `ex2-grid-editor/service`, `ex2-grid-editor/web`, `ex2-grid-editor/docs`

ID: DI-sukip
Date: 2026-07-13 18:35:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: The metadata slice will use a new `protocols/document-metadata.md` spec and a relay-signed latest-state model per document, while browser search/collection views query relay metadata rather than inventing a separate authoritative registry.
Intent: Keep document metadata durable and relay-verifiable, but avoid building a full authoritative document-management server before owner/admin and restore decisions are separately tested.
Constraints: Metadata updates are current-time actions; relay peer feeds should carry metadata envelopes; search is limited to relay-known documents and metadata in this slice; browser-local saved versions and comments stay outside the metadata protocol.
Affects: `ex2-grid-editor/protocols/document-metadata.md`, `ex2-grid-editor/service`, `ex2-grid-editor/web`, `ex2-grid-editor/docs/architecture.md`

Goal: Add the first PromiseGrid-native document metadata slice with relay-signed metadata state, relay-backed document search, and browser surfaces for description, tags, collections, favorites, and archive.

- [x] pamov.1 Add the document-metadata protocol doc and relay metadata types.
- [x] pamov.2 Add relay metadata update/get/search endpoints and peer ingestion rules.
- [x] pamov.3 Add browser metadata editing, search, and document-management UI.
- [x] pamov.4 Add tests and docs for the metadata slice.
