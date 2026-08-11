# moks.maintenance.service.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `maintenance_id`, `run_id`, `resource_id`, and `performer`. It records a maintenance service execution; local policy decides its consequences.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
