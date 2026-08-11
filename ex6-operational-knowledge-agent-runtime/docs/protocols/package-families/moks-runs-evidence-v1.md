# moks.runs.evidence.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `run_id` and `summary`, optional object `facts` mapping strings to strings, and optional non-empty string `body_ref`. It records evidence about a run; it grants no authority.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
