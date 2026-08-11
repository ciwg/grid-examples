# Ex6 Semantic Author-Signature Admission Policy

Status: current local runtime policy

## Rule

Ex6 requires every canonical package record admitted through local append or
relay import to carry a complete, valid semantic author signature. A record
with missing, incomplete, malformed, or invalid author evidence is refused and
does not enter durable package history. Source: DI-damab.

## Boundary

This is Ex6's local admission policy. It does not alter the frozen pCID or
wire-level meaning of any package family. Each family pCID continues to select
its immutable family specification and the shared canonical carriage slots.

The signature is durable evidence that an author signed the canonical record
content. It does not automatically grant authority, approval, role membership,
workflow execution, or trust. Those remain separate local-policy decisions.

## Relay distinction

Relay-carriage signatures authenticate the peer that carried exact bytes.
Semantic author signatures authenticate the record's signed content. Ex6 checks
both where their respective import policy applies; neither substitutes for the
other.

## Local signing

When a local append receives an unsigned record from a local command boundary,
the runtime signs it with its current local author key before admission. An
imported unsigned record is not locally re-authored: it is refused, because
adding a new signature would not establish who made the original semantic
statement.

## Future runtimes

Another runtime may adopt a different explicitly documented local admission
policy while using the same frozen family pCIDs. A change to record encoding,
signature meaning, or payload semantics requires a new protocol specification
and pCID; a local trust-policy change alone does not. Source: DI-damab.
