# Canonical Makerspace Record Re-Creation

TE ID: TE-malap

## Status

decided

## Decision under test

How Ex7 should replace its unreleased, local JSONL `Event{Type: ...}` evidence
model with canonical PromiseGrid records while preserving the makerspace
demonstration's useful behavior: observations, safety holds, steward
inspections, off-site loan terms, and returns.

## Assumptions and scope

- Ex7 is unreleased. Existing `.makerspace-stewardship/events.jsonl` files are
  development artifacts, not a public compatibility commitment.
- Alice may observe an unsafe table saw and make a voluntary loan promise for
  an eligible portable tool. Carol may make a protocol-defined inspection or
  clearance promise under the Woodworking role that Alice locally recognizes.
  Mallory may provide malformed, unsigned, unknown, or stale record bytes.
- Durable makerspace evidence must be canonical CBOR `grid(...)` bytes selected
  by frozen pCIDs. A browser JSON request remains a local ingress convenience,
  not the durable or relay format.
- Every active makerspace family needs one immutable specification whose pCID is
  the CIDv1 of that exact specification. The service's current-state display is
  a local derived projection, not an authoritative mutable record.
- Ex7 will follow the current local admission policy used by the aligned
  examples: semantic author evidence is verified before durable admission;
  relay-carriage signatures, when relay is added, remain distinct evidence.
- This TE does not choose a universal identity system, global authority,
  consensus, automatic policy enforcement, or a browser embodiment.

## Alternatives

### A. Clean canonical re-creation

Create a new canonical record store for Ex7. Define frozen specifications for
the active makerspace evidence families; append exact Grid bytes; derive the
browser state from those records. Use only the top-level PromiseGrid semantic
action `promise`; each family pCID defines whether its payload means an
observation, safety hold, inspection, loan commitment, or return observation.
The old JSONL file is not replayed by the new runtime. It remains an explicitly
legacy development artifact that may be archived or inspected outside normal
admission.

### B. Compatibility-first dual reader

Keep JSONL as a supported durable format while adding canonical records for new
writes. Replay both formats into one projection and retain the legacy `Type`
switch as a first-class semantic path.

### C. Automatic JSONL conversion on startup

On first startup, parse historical JSON events, synthesize canonical records,
and have the current runtime sign or otherwise admit them as if they were new
authors' evidence.

### D. Keep the local JSONL model

Document its local-only limit and defer canonical records, pCIDs, author
evidence, and relay compatibility indefinitely.

## Scenario analysis

### Normal operation

Alice records an observation that the table-saw guard is loose, Carol records
an inspection that clears a hold, and Alice accepts a cordless-drill loan with
a due date and policy snapshot. Under A, each browser request creates one
exact, pCID-selected canonical record. The local projection explains the
current tool state by reading those records, while the frozen family specs tell
another implementation what each payload means. The UI can still use friendly
JSON without making it the protocol.

B makes the happy path appear smooth, but a reader must continue to understand
two unrelated durable semantic systems. D remains easy to run but makes the
PromiseGrid label a presentation claim. C gives old events new-looking bytes,
but cannot truthfully prove that Alice or Carol authored the transformed record.

### Failure, corruption, and incomplete writes

If power fails during Alice's write, A can retain the existing append-before-
projection, fsync, and fail-closed discipline, now over framed canonical bytes.
A malformed or incomplete frame remains evidence of a failed write rather than
an instruction to derive state from a selected prefix.

B must distinguish corruption in two parsers and decide ordering between JSON
and CBOR sources. C adds a worse failure window: a partly converted log mixes
old history with synthetic new history. D detects malformed JSON but cannot
make the stored bytes portable or pCID-defined.

### Concurrent actors and mixed-version nodes

With A, Alice's and Carol's records have independent author evidence and exact
bytes. A newer local projection may preserve a record for an unknown pCID
without assigning it known makerspace meaning. A node that does not implement a
family can carry it exactly and show it as unknown rather than discard it.

B requires every version to retain the legacy event vocabulary forever. C
changes evidence meaning across an upgrade boundary and creates different
results when two runtimes convert the same JSON in different local contexts.
D has no credible mixed-node format at all.

### Long-horizon evolution and migration

Years later, a new makerspace may add a maintenance-calibration workflow. A
adds a new frozen family spec and pCID; existing records and their meanings do
not change. A workflow composes those pCIDs rather than receiving a new
application-wide protocol identity. Forks can select a different spec without
rewriting the old evidence.

B accumulates permanent compatibility code and blurs which of two records is
normative. C overwrites the boundary between historical local claims and new
canonical evidence. D leaves every future workflow tied to application source
and JSON field conventions.

### Trust-boundary changes

Under A, semantic author signatures establish who made a durable makerspace
promise under the configured local policy. The projection may decide whether
Carol's inspection is locally recognized to clear a hold, but it does not turn
that local policy decision into global authority. Later relay import verifies
carriage separately and does not locally re-author imported records.

B keeps unsigned legacy events as equally meaningful alongside signed records.
C is unacceptable because runtime signing would make the runtime appear to be
the historical author. D explicitly remains unauthenticated local state and
cannot make a PromiseGrid author-evidence claim.

### Scale and operational complexity

A adds a codec, frozen-spec registry, signature handling, framed append log,
and projection tests. That complexity is concentrated at one shared durable
boundary and gives every later family the same path. Photo blobs can be moved
to content-addressed objects rather than repeated data URLs in durable records.

B keeps all of A's new machinery plus dual-format handling. C adds conversion,
rollback, and provenance ambiguity. D is smallest now but defers every required
alignment capability and makes the eventual re-creation larger.

## Conclusions

D is rejected because it leaves the central PromiseGrid evidence boundary
absent. C is rejected because it fabricates canonical-looking historical author
evidence. B is rejected because Ex7 is unreleased and would preserve the
legacy JSON event model as an indefinite second protocol.

A survives and is recommended: cleanly re-create Ex7 around canonical
pCID-selected promise records. Preserve the user-facing makerspace story and
the current fail-closed durability rule, but do not carry the legacy JSONL
semantic path forward. Treat pre-re-creation JSONL only as explicit archived
development evidence, never as automatically admitted canonical history.

## Output to decision framing

The remaining decision is whether to lock Alternative A as the Ex7 migration
strategy. If locked, the first implementation slice will:

1. create one new Ex7 alignment TODO with an append-only DI log;
2. define the initial frozen makerspace family-spec set and central fixed pCID
   registry;
3. replace JSONL durable events with framed canonical Grid bytes and a derived
   projection; and
4. retain JSON only at the browser ingress/egress boundary, with no automatic
   legacy replay or re-authoring.

## Decision status

locked: Alternative A, clean canonical re-creation, by DI-nohos in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`.

## Implications and future work

- The existing local-demo documentation must be superseded where it calls
  `events.jsonl` the durable evidence source.
- A separate decision will define the initial family names, frozen spec paths,
  signature-key storage, and exact record-store framing after this migration
  strategy is locked.
- Browser interaction evidence and relay embodiment are follow-on slices, not
  prerequisites for choosing the canonical durable boundary.
