# moks.training.completion.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `training_id`, `person`, and `decision`, plus optional string `notes`. It records a training completion statement; local policy decides its consequences.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
