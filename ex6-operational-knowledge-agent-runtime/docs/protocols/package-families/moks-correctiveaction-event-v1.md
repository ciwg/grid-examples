# moks.correctiveaction.event.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `action_id`, `quarantine_case_id`, `actor`, `evidence_id`, and `summary`, plus optional string `notes`. It records a corrective-action event; local policy decides its consequences.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
