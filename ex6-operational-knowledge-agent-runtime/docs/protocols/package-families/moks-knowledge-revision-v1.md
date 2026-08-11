# moks.knowledge.revision.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty `item_id`, positive integer `revision`, non-empty `title`, and non-empty `body_ref`. It records a revision of a knowledge item; it grants no authority.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
