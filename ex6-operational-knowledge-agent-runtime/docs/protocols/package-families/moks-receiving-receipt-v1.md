# moks.receiving.receipt.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `receiving_id`, `run_id`, `place_id`, and `receiver`. It records a receipt execution; local policy decides its consequences.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
