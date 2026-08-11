# moks.receiving.disposition.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `receiving_id` and `decision`, optional string `resource_id`, and optional string `notes`. It records a receiving disposition statement; local policy decides its consequences.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
