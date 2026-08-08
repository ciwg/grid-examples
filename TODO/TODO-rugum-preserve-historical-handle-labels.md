# TODO rugum — Preserve historical handle labels

## Decision Intent Log

### DI-jutok

- ID: DI-jutok
- Date: 2026-08-07 21:42:11 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Preserve historical manual and noncanonical identifiers as stable references; mint handles only for new coordination records unless JJ explicitly requests a scoped migration.
- Intent: Avoid broad, non-functional reference churn while ensuring all future records use the repository-wide allocator.
- Constraints: Do not reinterpret this as permission to retain new temporary IDs; ex7's temporary IDs have already been migrated.
- Affects: `AGENTS.md`, `TODO/handle-namespace.tsv`, `TODO/TODO.md`, and this TODO record.

## Tasks

- [x] rugum.1 Record the historical-label preservation rule in `AGENTS.md`.
- [x] rugum.2 Commit and push the rule with an explicit decision summary in the commit message.
