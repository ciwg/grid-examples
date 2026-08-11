# moks.quarantine.event.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `case_id`, `transition`, `actor`, and `evidence_id`, optional strings `receiving_id`, `receipt_run_id`, `exception`, and `notes`. It records a quarantine event; local policy decides its consequences.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
