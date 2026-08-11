# moks.runs.approval.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `run_id` and `decision`, and optional string `notes`. It records an approval statement; local policy determines whether it is accepted or effective.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
