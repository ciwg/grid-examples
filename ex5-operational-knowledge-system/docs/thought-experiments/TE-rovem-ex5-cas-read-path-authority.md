# Ex5 CAS Read-Path Authority

TE ID: `TE-rovem`
## Status
needs DF

## Decision under test

What the first authoritative CAS-backed read/replay step for `ex5` should be
now that signed family envelopes and copied evidence blobs already dual-write
into CAS.

Related TODO:

- `104` - `ex5-operational-knowledge-system/TODO/TODO-tavob-ex5-cas-read-path-adoption.md`

## Assumptions

- `ex5` already dual-writes all five signed family envelopes into CAS by
  envelope CID.
- `ex5` already dual-writes copied evidence blobs into CAS by blob CID.
- Startup replay still reads `events.jsonl` plus the five signed family JSONL
  logs as the active source of truth.
- Places, resources, and runs are not yet frozen PromiseGrid families and
  still depend on the compatibility event log for local replay.
- The user wants real progress toward a stricter fully-on-grid runtime, not a
  permanent “CAS as checksum cache only” posture.

## Alternatives

### Alternative A

Keep the current logs authoritative and use CAS only for verification,
portability, and recovery tooling.

This avoids replay migration now, but leaves the main runtime behavior
log-first.

### Alternative B

Make CAS authoritative for the five signed family envelopes while keeping the
compatibility event log active for the still-unfrozen runtime state.

Under this choice:

- envelope bytes in CAS become the authoritative durable record for the five
  frozen families
- family JSONL logs become manifests or indexes used to enumerate the envelope
  CIDs until a better manifest layer exists
- `events.jsonl` remains necessary for places, resources, runs, and other
  unfrozen projections

### Alternative C

Attempt a full CAS-authoritative runtime now, including replay for the
unfrozen event-only state.

This would try to eliminate log-first replay broadly rather than only for the
already-frozen families.

## Scope and systems affected

- `ex5-operational-knowledge-system/TODO/TODO-tavob-ex5-cas-read-path-adoption.md`
- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/service/persistence.go`
- `ex5-operational-knowledge-system/service/app.go`
- signed family record loaders and verification paths
- CAS/object enumeration logic
- PromiseGrid claims, architecture, and practical implementation docs

## Scenario analysis

### Scenario 1: normal restart after local writes

Alice creates items, approvals, evidence, links, and responsibilities, then
restarts the runtime.

Alternative A:

- behaves exactly like today
- keeps CAS present but non-authoritative
- does not reduce the gap between “CAS exists” and “CAS is real runtime state”

Alternative B:

- lets the five frozen families replay from authoritative CAS-backed envelope
  bytes
- keeps the still-unfrozen event-only state on the compatibility path
- improves the runtime where PromiseGrid-native coverage already exists

Alternative C:

- tries to move the whole runtime at once
- immediately runs into the fact that places/resources/runs are still not on
  frozen family contracts

Result:

- B is the strongest implementable next step.

### Scenario 2: family log drift or partial loss

Bob loses part of a signed family log but still has the CAS objects and the
remaining manifests.

Alternative A:

- cannot replay authoritatively from CAS
- still treats the drifted log as the primary runtime problem

Alternative B:

- can define recovery around the CAS-backed envelope bytes for the frozen
  families
- reduces the blast radius of family-log corruption

Alternative C:

- also helps here
- but expands the migration problem to unfrozen state that the repo has not
  modeled in CAS-native form yet

Result:

- B meaningfully improves recovery without overreaching.

### Scenario 3: evidence portability and later peer exchange

Carol wants peer-visible `knowledge-evidence` exchange later.

Alternative A:

- keeps CAS blobs useful, but still secondary
- leaves evidence portability dependent on a future bigger storage cutover

Alternative B:

- makes the signed envelope side of evidence more grid-native now
- gives the later evidence-carriage work a cleaner storage story because both
  metadata envelopes and blobs already have authoritative CAS identities

Alternative C:

- may be ideal eventually
- but entangles evidence portability with a full runtime storage rewrite

Result:

- B improves the later evidence step without forcing full CAS-only replay now.

### Scenario 4: mixed-version or partially migrated nodes

Dave runs the current log-first build; Ellen runs the first CAS-authoritative
family replay build.

Alternative A:

- avoids migration complexity
- but leaves the repo with no actual CAS read-path progress

Alternative B:

- still preserves the family logs as manifests/indexes, so mixed-version
  compatibility is manageable
- gives the newer node stronger semantics without deleting the older path

Alternative C:

- raises the mixed-version risk sharply because older logs and unfrozen state
  would need broader reinterpretation

Result:

- B is the safest meaningful migration step.

### Scenario 5: long-horizon “fully on the grid” evolution

Steve wants CAS eventually to matter as more than a sidecar cache.

Alternative A:

- is too conservative for that goal
- risks leaving CAS permanently adjunct instead of authoritative

Alternative B:

- advances authority for the families that already have frozen protocols and
  signed envelopes
- leaves a clear remaining backlog instead of claiming full completion

Alternative C:

- aims at the final target directly
- but the unfrozen parts of the runtime are not ready for it yet

Result:

- B best matches a real staged migration toward the intended end state.

## Conclusions

Rejected alternatives:

- Alternative A: too conservative; it leaves CAS real only on write, not on
  authoritative replay
- Alternative C: too broad for a runtime that still has unfrozen event-only
  state

Surviving alternative:

- Alternative B: make CAS authoritative for the five frozen family envelopes
  first, while keeping compatibility event replay for the still-unfrozen state

Implications and future work:

- `104` should focus first on authoritative CAS-backed replay for the five
  frozen families, not on immediate CAS-only replay for every runtime record
- the family JSONL files may survive temporarily as manifests/indexes even
  after authority moves to CAS
- `105` remains the place to decide how peer-visible evidence exchange carries
  blob bytes, but it becomes cleaner once the evidence envelope side already
  reads authoritatively from CAS

## Decision status

Needs DF on whether to lock Alternative B now as the first authoritative CAS
read-path slice.
