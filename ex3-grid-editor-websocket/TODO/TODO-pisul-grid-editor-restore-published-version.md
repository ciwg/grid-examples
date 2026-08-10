# TODO pisul - grid-editor restore published version

Goal: Define and implement a user-directed workflow that makes a selected
earlier published document version the current working document without
rewriting existing history.

- [ ] pisul.1 Run a thought experiment for provenance, conflict handling,
  authority boundaries, and the distinction between restore and import.
- [ ] pisul.2 Lock the user-visible restore semantics, artifact references,
  behavior names, and storage paths through Decision Framing and a DI.
- [ ] pisul.3 Implement the selected workflow with append-only evidence and
  regression coverage.
- [ ] pisul.4 Document the workflow, its trust boundary, and verification.

Status: deferred. No restore behavior is authorized by this TODO alone. A
future TE and Decision Framing must decide whether a selected published version
creates a new working change, how concurrent edits are handled, and which
actor may request the operation.
