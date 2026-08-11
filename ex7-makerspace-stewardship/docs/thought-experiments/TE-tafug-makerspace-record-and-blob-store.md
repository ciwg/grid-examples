# Makerspace Framed Record Store and Blob Boundary

TE ID: TE-tafug

## Status

superseded by TE-zajop / DI-rifib

## Decision under test

How Ex7 persists exact signed canonical records and photo blobs so a browser action is acknowledged only after its evidence is durable and replayable.

This TE narrows TODO `bubuz` after the frozen envelope and family specs.

## Assumptions and scope

- Durable records are exact canonical signed Grid bytes; the store does not re-encode them.
- An observation with a requested safety hold creates two linked records and must become visible together or not at all.
- Photos are content-addressed blobs. Mallory may leave interrupted writes, malformed frames, missing blobs, or untrusted exact record bytes.
- No database, replication, garbage collector, repair tool, or multi-process writer protocol is added here.

## Alternatives

### A. Blob-first writes plus append-only transaction frames

Use these paths under `<runtime-root>`:

```text
<runtime-root>/authors.json
<runtime-root>/records.frames
<runtime-root>/blobs/<cidv1-base32>
<runtime-root>/tmp/<opaque-temporary-name>
```

`records.frames` starts with `MSR1\n`. Each append has an unsigned 64-bit big-endian payload length followed by canonical CBOR array of byte strings. Each string is exactly one Grid record. Ordinary actions use one-record frames; safety holds use an ordered observation/disposition two-record frame.

For every photo, write decoded bytes under `tmp`, fsync and close, rename to its `blobs/<cid>` path, then fsync the blob directory. Reuse an existing CID only after recomputing and matching bytes. Once blobs are durable, append the complete frame, fsync and close `records.frames`, then apply its records to memory. A request succeeds only then. A crash before frame durability may leave an orphan blob but no projected evidence; a partial frame fails startup.

### B. One record per frame

Use the same blob process but append observation and safety disposition as separate frames.

### C. Mutable database or snapshots

Persist tool state in a transactional database or snapshots instead of exact append-only record frames.

## Scenario analysis

### Normal operation

A makes photo bytes durable before their records refer to them. A safety-hold request becomes one durable two-record operation. Carol's clearance and Alice's loan/return replay from exact record bytes. B can retain an observation while losing its requested hold. C hides primary evidence behind database-specific state.

### Failure and incomplete writes

If a blob write fails, A writes no frame and applies no state. If a frame append fails, the request fails; a later partial frame fails closed instead of projecting a prefix. Missing referenced blobs also fail closed. A crash after blob durability but before frame durability leaves only an unreferenced blob. B lacks action-level atomicity. C weakens direct exact-byte audit and unknown-family retention.

### Concurrent actors and mixed versions

The existing app mutex serializes A's frame appends. Unknown pCIDs remain exact bytes in frames without local semantics. New families do not change frame parsing. B lacks grouped actions; C requires schema compatibility to retain unknown durable families.

### Long horizon, trust, and scale

A keeps framing stable while pCID families evolve and lets workflows compose records without a storage migration. It preserves semantic author evidence separately from any future relay signature. Blob CIDs prove bytes, not authority or truth. A adds bounded framing and atomic file operations but no external service. B is smaller but partial; C is larger and less PromiseGrid-legible.

## Conclusions

B is rejected because a requested hold must not partially project. C is rejected because Ex7 needs exact append-only evidence and unknown-byte preservation. A survives and is recommended.

## Output to decision framing

- **A (recommended):** `MSR1\n` transaction frames, blob-first fsynced writes, then framed-record fsync, then projection, at the exact paths above.

## Decision status

superseded by TE-zajop / DI-rifib. Its transaction framing and blob-first
durability mechanics remain applicable only under one participant agent's
runtime root; its single-runtime author-key path is not valid for Ex7's
decentralized product architecture.

## Implications and future work

- `tmp/` is only for incomplete blob writes. Successful renames remove their source; stale entries need an explicit future cleanup decision.
- Frame and blob size limits require a later bounded-parser test decision.
- Legacy `events.jsonl` is never automatically imported.
