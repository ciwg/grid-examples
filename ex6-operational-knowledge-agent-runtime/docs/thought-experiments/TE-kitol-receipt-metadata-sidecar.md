# Receipt metadata sidecar

TE ID: TE-kitol

## Status

decided

## Decision under test

How the receipt inbox retains the authenticated sending peer identity alongside
exact remote lifecycle evidence without rewriting that evidence or creating a
local workflow lifecycle event.

## Assumptions

- The relay endpoint authenticates `peerID` before accepting a transfer.
- `workflow-evidence` stores exact lifecycle bytes addressed by evidence CID.
- Lifecycle bytes do not contain the authenticated transport sender identity.

## Alternatives

1. Omit sender identity from the inbox.
2. Wrap or rewrite lifecycle evidence to add sender identity.
3. Store a local receipt-metadata sidecar keyed by evidence CID.

## Scenario analysis

Omitting identity makes a receipt list less auditable and cannot distinguish
two peers carrying the same artifact. Rewriting evidence changes the exact
bytes whose CID and signature Bob verified. A sidecar preserves evidence bytes,
records only Bob's local observation that a specific authenticated peer supplied
them, and can be rebuilt/validated against the evidence CAS.

If metadata write fails after evidence retention, the inbox reports an
incomplete receipt and refuses its convenience import; Bob may still retain the
evidence for diagnostics. Multiple peers create separate evidence-CID metadata
entries while grouping under one artifact CID. The sidecar creates no remote
authority and no new PromiseGrid wire action.

## Conclusion

Choose option 3: a local receipt-metadata sidecar maps evidence CID to the
authenticated peer ID and artifact CID. It is local projection data, never a
workflow lifecycle event, and never modifies signed evidence bytes.

## Decision status

Locked by DI-rufir.
