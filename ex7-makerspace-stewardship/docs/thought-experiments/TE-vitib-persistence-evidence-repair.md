# Persistence evidence repair for Makerspace Stewardship

**TE ID:** TE-vitib

## Status

decided

Identifier migration per TODO-tagup: `TE-pending-mint-ex7` -> `TE-vitib`. The following DIs in the related decision record are unchanged in meaning: DI-dapod, DI-patag, DI-sapun, DI-malih, DI-lasif.

## Decision under test

How `ex7-makerspace-stewardship` should make equipment evidence durable and replayable while fixing the review findings for oversized photo events, state changes that survive failed writes, and off-site loan terms that are not faithfully reconstructed.

The TE ID was minted through the repository-root allocator after the initial evidence repair; the decision content and scenario analysis are unchanged. Source: DI-dapod.

## Assumptions and scope

- Scope: the local append-only event log at `.makerspace-stewardship/events.jsonl`, the Go service replay path, and regression tests in `ex7-makerspace-stewardship`.
- Alice is a qualified member who reports an unsafe table saw and borrows an eligible portable tool. Carol is the recognized Woodworking steward. Mallory can corrupt a local log only if she already crosses the local filesystem trust boundary; the application must not silently manufacture evidence from corrupted bytes.
- The application remains a single-process local demo. This TE does not add identity authentication, multi-node replication, policy authoring, or a new top-level PromiseGrid action kind.
- A successful mutating HTTP response means the corresponding evidence is durably committed to the local log.

## Alternatives

### A. Event representation

1. Keep append-only events, including a complete copy of accepted loan policy terms in each loan event.
2. Keep append-only events but reconstruct loan terms from the then-current in-memory policy during replay.
3. Replace the event log with mutable state snapshots.

### B. Mutation order and durability

1. Encode, write, `fsync`, and close an event; apply it to in-memory state only after that succeeds.
2. Apply the mutation first, then append the event.
3. Append and close without `fsync` before reporting success.

### C. Corrupt or incomplete log handling

1. Fail closed: refuse startup, preserve the evidence file untouched, and report the error.
2. Recover the valid prefix while retaining the corrupt tail for repair.

## Scenario analysis

### Normal operation

Alice records a safety hold, Carol records an inspection that clears it, and Alice borrows then returns the cordless drill. Under A1/B1, each accepted action has one exact event and the returned state can be reconstructed after restart. A2 makes a currently visible loan appear to have accepted terms, but cannot prove those are the terms Alice accepted. A3 simplifies reads but loses the ordered evidence trail that explains why the current condition changed.

### Power failure, full disk, and incomplete writes

If the disk fills after Alice's observation, B2 leaves live memory showing a safety hold while the restart state has none. B1 returns an error and leaves both live and replayed state unchanged. `fsync` in B1 narrows the gap between an acknowledged response and a power loss; B3 cannot make the same claim. A partial final JSON line is possible after an interrupted write. C1 makes this evidence problem visible and prevents silently operating on an incomplete history. C2 improves availability but creates a repair policy, a user-visible degraded state, and ambiguity about the omitted final action.

### Concurrent actors and mixed versions

Carol clearing a hold and Alice creating a loan are serialized by the application's mutex. An event that carries the loan's complete policy snapshot remains intelligible to a newer process even if the current area policy changes. A2 instead assigns old loans whatever terms the new process happens to have compiled in. Older binaries will fail closed on an unknown event-schema change, which is preferable to silently changing historical terms; this repair does not introduce multi-process writer coordination.

### Long-horizon evolution and migration

Years later, the Fiber Arts area may revise its policy. A1 preserves the actual policy version and text that accompanied the sewing-machine loan, so a migration can keep it as immutable history. A2 cannot distinguish an old policy from a new one. A3 requires a separate audit-history system to recover the same provenance. The temporary TE filename is a documentation migration obligation only; it does not affect runtime evidence.

### Trust-boundary changes

Mallory's modification of the log is not prevented by this local demo's filesystem model. C1 ensures the service does not conceal malformed evidence or continue from an arbitrarily selected prefix. A1 captures the terms that the service presented but is not a cryptographic signature; adding signatures or authorization is outside this repair's scope.

### Scale effects

Photo data URLs can exceed `bufio.Scanner`'s default 64 KiB token limit. The reader must use a bounded mechanism sized consistently with the maximum accepted event, avoiding both restart failures for valid images and unbounded allocation for malicious files. Per-event `fsync` trades throughput for a simple, strong local acknowledgment boundary. That cost is appropriate for occasional volunteer makerspace records; batch/group durability would need its own decision.

## Conclusions

The user selected A1, B1, and C1:

- Persist a complete loan-policy snapshot in the loan event and replay it verbatim.
- Make evidence durable with write, `fsync`, and close before applying the corresponding state mutation or returning success.
- Fail closed on malformed or incomplete evidence logs, preserving the file for investigation.

Rejected: current-policy reconstruction loses accepted-term provenance; mutate-first writes create phantom live state; close-only writes weaken the recorded-evidence guarantee; prefix recovery hides an unresolved evidence gap; mutable snapshots lose ordered evidence.

## Decision status

locked: DI-dapod, DI-patag, and DI-sapun in `../../TODO/TODO-tagup-persistence-evidence-repair.md`.

## Implications and future work

- Add regression tests for a valid near-limit photo event, failed durable append, corrupt/incomplete log startup, and policy-snapshot replay.
- Keep the runtime event-size bound and HTTP photo-size validation aligned.
- A future multi-process or signed-evidence design requires a separate TE and DI.
