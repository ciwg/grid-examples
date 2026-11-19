# Ex3 restore a published version

TE ID: TE-bazav
## Status

decided, refined

## Decision under test

How should Ex3 let a participant make a selected published document version the
current working document without rewriting live history, treating a historical
publish as authority, or conflating restore with import?

This TE addresses TODO `pisul`, especially `pisul.1`. It follows the existing
`publish-document` boundary: a publish manifest is a durable, signed,
repo-local-draft handoff artifact; it is not the live editing protocol and is
not a frozen upstream PromiseGrid specification.

## Assumptions and trust model

- `live-document` remains the CRDT editing surface; its log and snapshots are
  append-only local relay history.
- A resolved `publish-document` manifest names exact CAS-backed text and
  Automerge replica bytes. Its envelope CID is stable provenance for the
  historical selection.
- Alice, Bob, Carol, and Dave are cooperative participants. Mallory may send
  malformed, missing, stale, or misleading artifact references, or may act
  with only the capabilities the relay granted her.
- Ex3 currently has local/capability mutation admission, not a frozen owner,
  administrator, delegation, or recognized-role-continuity protocol.
- Import already materializes a **new** local document from a publish artifact.
  Restore must instead affect the selected existing live document.
- The top-level protocol semantic remains `promise`; a restore-specific meaning
  belongs to a pCID-defined payload and evidence record, not to a new global
  action kind. Source: `DI-mosoj`.

## Alternatives

### A — Treat restore as import under a different label

Resolve the manifest and create a new local document, then direct Alice to use
that document instead of the current one.

This has simple provenance and no effect on current collaboration, but it does
not satisfy the requested operation: the original working document is not made
current from the selected published version.

### B — Append a provenance-bearing restore change to the existing document

An admitted participant selects an exact publish-manifest CID. The relay
resolves and verifies its referenced bytes, then appends a new current-time
restore promise and a CRDT change derived from that selected version. The new
record names the source manifest CID, source text/replica CIDs, requesting
participant, and the observed live-document position. Existing history and the
manifest remain unchanged.

The result is a new working state, not a claim that the historical state is now
globally authoritative or that all concurrent edits disappeared.

### C — Authoritatively replace the current live state

The relay replaces its current snapshot/log with the selected version and
declares that version the sole current document, possibly limiting the action
to an owner or administrator.

This gives a familiar centralized-product undo experience, but it rewrites or
obscures history, discards or overrides concurrent work, and requires an
authority/role model Ex3 does not have.

## Scenario analysis

### Normal operation

Alice selects a publish manifest for yesterday's approved checklist while the
document is otherwise quiet.

- A produces a useful new copy but does not restore the shared working
  document.
- B produces a new, inspectable change in the existing document and preserves
  the exact historical source through its manifest CID. Bob can see both the
  resulting text and the provenance record.
- C appears simplest in the UI, but makes the server's declaration—not the
  append-only evidence—the source of truth.

### Concurrent actors and mixed-version nodes

Alice requests restore while Bob is editing a different paragraph, and Carol
joins from a node that has not yet seen the manifest.

- A avoids conflict only by avoiding shared restoration.
- B carries the restored content as a new CRDT change. Bob's concurrent change
  remains a concurrent CRDT input; the resulting text may not be byte-identical
  to the selected publication, so the UI must say it is a restore-derived
  working state and expose the manifest reference. Carol can fetch/assess the
  referenced manifest independently before treating the claim as meaningful.
- C either drops Bob's work or invents server precedence. Both make one relay
  an unnecessary central editor authority and make mixed-version recovery
  ambiguous.

### Corruption, incomplete writes, and unavailable artifacts

Alice selects an envelope CID whose manifest is missing, invalid, or points to
missing CAS bytes; a process fails between validation and append.

- A can fail safely, but still gives no restoration behavior.
- B validates the manifest signature and referenced bytes before creating any
  live change. A failed validation or incomplete attempt produces bounded local
  diagnostic/observation evidence and no partial restore record. A durable
  accepted restore record is sufficient on replay to identify exactly which
  source was applied.
- C risks leaving an ambiguous half-replaced snapshot unless it introduces a
  transactional authority layer beyond Ex3's current model.

### Trust-boundary change

Mallory can read published artifacts but lacks a mutation capability. Later,
an organization wants a reviewed restoration policy.

- A is harmless but does not meet the shared-workflow need.
- B reuses existing admission for the ability to request a current-time change;
  it does not claim that the requester owns the document or has a permanent
  role. A later pCID-defined review/acceptance protocol can add stricter local
  trust policy without falsifying this record's provenance.
- C demands a new owner/admin/delegation model immediately and would falsely
  imply recognized role continuity unless a separate protocol establishes it.

### Long-horizon evolution and forks

Years later, an implementation migrates storage or a community forks its
workflow rules.

- A preserves the publication but loses the intent to continue the original
  document.
- B retains immutable source references plus a new current-time record. A
  migrating or forked implementation can replay, reject, or reinterpret the
  restore profile by its pCID without erasing the underlying publication or
  prior live history.
- C couples continuity to the particular server's mutable current-state claim;
  it is difficult to audit whether a later snapshot reflects history, a policy
  override, or data loss.

### Scale and operational complexity

Large documents and many published versions make repeated full restorations
expensive.

- A duplicates documents and multiplies user-facing branches.
- B needs bounded manifest/CAS validation, a clear observed-position marker,
  and normal CRDT/storage limits, but reuses existing CAS and append-only
  history. It leaves compaction as a separate local implementation concern.
- C may seem efficient by replacing a snapshot, but shifts complexity into
  coordination, locks, authorization, rollback, and recovery.

## Conclusions

- Reject A: it is the already-shipped import behavior, not restore.
- Reject C: it centralizes authority, weakens provenance, and conflicts with
  append-only repairable collaboration.
- Keep B: a restore is an explicit, current-time, pCID-defined promise and
  CRDT-derived working change that references an immutable published manifest.

## Surviving design and new decisions exposed

The surviving direction is B, subject to DF. It exposes these decisions:

1. Should only a resolved publish-manifest CID be accepted as a restore source,
   or may browser-local saved versions also be restored directly?
2. Should a concurrent restore always create a CRDT-merged working state with
   an explicit non-byte-identical warning, or fail when the observed position
   has advanced?
3. Should existing local/capability admission authorize the request, or should
   restore wait for a new review/role protocol?
4. What pCID-defined profile name, payload fields, UI wording, function names,
   and storage paths make the new provenance record unambiguous?
5. Must the restore promise and live CRDT change be durably coupled before the
   UI reports success, and what bounded diagnostic evidence records a failure?

## Implications for TODOs and pending DIs

- `pisul.2` must lock the five decisions above through DF and record them in a
  DI before code changes.
- `pisul.3` must add deterministic tests for valid restore provenance,
  concurrent edits, unavailable/tampered artifacts, capability denial, and
  replay after restart.
- `pisul.4` must document the distinction among publication, import, and
  restore, including the local-draft/non-frozen scope.
- Existing TE-vafor remains historical publication provenance. Restore extends
  it by reference; it does not revise or reinterpret its locked publish/import
  boundary.

## Decision status

locked by DI-hihok

## Refinements

### 2026-08-14 — Atomic restore artifact

DI-tibum locks the durable-coupling detail exposed by this TE: one dedicated
restore promise artifact carries both the selected publish-manifest provenance
and the exact CRDT change. This avoids an unpaired two-record recovery state.
