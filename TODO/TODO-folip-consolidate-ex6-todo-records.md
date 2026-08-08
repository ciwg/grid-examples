# TODO folip — Consolidate duplicate ex6 TODO records

## Decision Intent Log

### DI-vajol

- ID: DI-vajol
- Date: 2026-08-07 21:56:47 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Retain the richer ex6 `sibok` and `nupad` TODO records as their canonical files; merge the short records' unique task/scope content into them; then remove the two redundant short files.
- Intent: Restore one authoritative TODO file per handle without losing decision history, active work, or Git history.
- Constraints: Preserve the two open `sibok` workflow-hardening tasks; use Git removal only after their content is present in the canonical file; do not alter runtime code or historical decisions.
- Affects: the four ex6 `sibok`/`nupad` TODO files, `TODO/TODO.md`, this TODO record, and `TODO/handle-namespace.tsv`.

## Tasks

- [x] folip.1 Merge the short records' unique content into the canonical records.
- [x] folip.2 Remove the two redundant short records and verify no dangling path references.
- [x] folip.3 Commit and push the consolidation.
