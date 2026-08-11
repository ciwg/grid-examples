# Initial Makerspace Family Boundaries and Frozen-Spec Paths

TE ID: TE-fivod

## Status

decided

## Decision under test

Which initial durable makerspace evidence families Ex7 should freeze, and where
their immutable specifications and central pCID registry should live, before
any pCID is calculated.

This TE narrows TODO `bubuz`, slice 1. It follows the clean canonical
re-creation locked by TE-malap / DI-nohos.

## Assumptions and scope

- The durable boundary uses a shared canonical Grid record envelope with
  semantic author evidence. Individual pCIDs define makerspace payload meaning;
  they do not identify individual tools, members, browser actions, or workflow
  instances.
- Existing startup data for members, areas, recognized steward roles,
  qualifications, and tool catalogue is local bootstrap/configuration in this
  first slice. It is not silently elevated into a durable interoperable family
  merely because the browser displays it.
- The currently demonstrated durable evidence is: condition observation,
  safety-hold disposition, off-site loan acceptance with policy snapshot, and
  return condition. The chosen set must represent those meanings without
  top-level action verbs beyond `promise`.
- Alice can report a condition and accept a loan; Carol can make a steward
  safety disposition that Alice's local policy recognizes. Mallory can send an
  unknown pCID, a record with a known pCID but invalid payload, or an otherwise
  valid record that the local policy does not treat as authoritative.
- Frozen specs are immutable. Adding a future workflow that only composes these
  meanings must not add a pCID or require an Ex7 rebuild.

## Alternatives

### A. Four focused evidence families plus one shared record envelope

Freeze a shared `makerspace-record-v1` envelope specification and four
independent payload-family specifications:

1. `makerspace-equipment-observation-v1` — a promiser's condition observation
   about one tool, including optional photo/blob references.
2. `makerspace-safety-disposition-v1` — a promiser's asserted safety
   disposition for one tool: place a hold or clear a hold with inspection
   evidence. Local policy decides whether a particular signer is recognized to
   clear a hold.
3. `makerspace-offsite-loan-v1` — a borrower's voluntary return commitment for
   one loanable tool, including due time and the exact accepted area-policy
   snapshot.
4. `makerspace-offsite-return-v1` — a borrower's observation/claim that one
   specific loaned tool was returned, including condition evidence.

Place the immutable files under
`docs/protocols/makerspace-families/` and the append-only mapping in
`docs/protocols/makerspace-family-pcid-registry.md`:

```text
docs/protocols/makerspace-record-v1.md
docs/protocols/makerspace-families/makerspace-equipment-observation-v1.md
docs/protocols/makerspace-families/makerspace-safety-disposition-v1.md
docs/protocols/makerspace-families/makerspace-offsite-loan-v1.md
docs/protocols/makerspace-families/makerspace-offsite-return-v1.md
docs/protocols/makerspace-family-pcid-registry.md
```

### B. One broad stewardship-evidence family

Freeze one `makerspace-stewardship-evidence-v1` pCID with a `kind` field for
observation, hold, clearance, loan, and return. Put one spec and one registry
entry under `docs/protocols/`.

### C. Seven or more families, including bootstrap governance

Freeze the four evidence families plus separate durable families for member
identity, area policy, steward recognition, qualification, and tool catalogue
in this first slice.

## Scenario analysis

### Normal operation

Alice observes a loose guard, Carol clears the resulting hold after inspection,
Alice accepts drill-loan terms, and later records a return. A gives each of
those claims a precise pCID-selected payload while retaining a shared envelope
and one projection path. The UI can show a chronological evidence story without
pretending its REST route names are protocol semantics.

B is simpler initially, but every implementation must parse a growing `kind`
union before it knows which signers, fields, and validation rules apply. C adds
families for data Ex7 currently treats as local setup, even though the demo has
no durable governance lifecycle for those facts.

### Failure, corruption, and incomplete writes

Under A, the shared envelope and framed store have one canonical validation and
durability path; each family validates only its own payload. A partial record
or malformed known-family payload fails closed before projection. The existing
loan policy snapshot remains inside the loan pCID where its evidentiary meaning
belongs.

B centralizes every payload validation into an ever-growing switch. C has more
cross-family ordering rules before a tool can be interpreted at all, increasing
the chance that partial bootstrap history makes the demo unusable.

### Concurrent actors and mixed-version nodes

Alice's observation and Carol's safety disposition can arrive in either order;
the projection applies local, pCID-defined assessment rules and does not claim
that either record is a global command. A newer node can retain an unknown
fifth family as exact bytes without deciding its semantics. A node that knows
the four families can continue operating on them.

B requires a new `kind` interpretation for every extension, making an unknown
variant an ambiguous half-known record. C makes ordinary new nodes depend on
the full bootstrap-governance protocol set before they can even host the
existing equipment evidence.

### Long-horizon evolution and workflows

A lets a future calibration workflow compose observation and safety-disposition
records. It needs no new pCID unless it introduces a genuinely new durable
meaning. A future interoperable recognition or qualification lifecycle may add
its own family and spec without changing existing mappings. The registry stays
append-only, so old pCIDs continue to name their exact frozen bytes.

B turns every new meaning into a version decision for the broad family. C
prematurely freezes bootstrap models that are likely to change as the makerspace
gains real membership and policy processes.

### Trust boundary and authority

A separates evidence from assessment. Carol's safety-disposition record is
signed evidence of Carol's promise/claim; Ex7's local recognition configuration
decides whether it clears the displayed hold. Alice's loan is evidence of her
return commitment, not an authorization grant. A future governance family can
make recognition portable only when its own trust, key-continuity, and revocation
rules are ready.

B makes those different assessments share one permissive field vocabulary. C
risks suggesting that static demo members and qualifications are durable global
authority records before their protocol rules exist.

### Scale and operational complexity

A creates five short, reviewable immutable specs and one small registry. It is
large enough for distinct validation and future composition, but small enough
to re-create the current runtime coherently. B has fewer files but a broader
and more brittle contract. C increases specification, test, projection, and
key-policy work without supporting an additional current user flow.

## Conclusions

B is rejected: a broad `kind` union weakens per-pCID payload and assessment
boundaries. C is rejected for the first slice: it mistakes local bootstrap
configuration for an already-designed interoperable governance protocol.

A survives and is recommended. Freeze the shared envelope plus the four focused
evidence-family specs at the exact paths listed above. Keep members, areas,
qualifications, steward recognition, and tool catalogue explicitly local
bootstrap data until a separate TE defines their durable interoperability and
trust rules.

## Output to decision framing

The surviving DF choice is:

- **A (recommended):** four focused makerspace evidence families, one shared
  record-envelope spec, and the exact `docs/protocols/makerspace-*` paths listed
  in this TE.

If selected, the next DI will lock the family names, paths, and bootstrap
boundary; the following implementation slice will create the immutable spec
files, calculate their CIDv1 pCIDs, and add the checked-in registry and tests.

## Decision status

locked: Alternative A, four focused evidence families and the exact frozen-spec
paths, by DI-bosur in `../../TODO/TODO-bubuz-canonical-makerspace-records.md`.

## Implications and future work

- `docs/protocols/` is a new Ex7 directory and will become the canonical
  location for frozen protocol specs and their registry.
- Signature key storage and the exact shared-envelope fields remain a later
  decision under TODO `bubuz` after family boundaries are locked.
- A new workflow or user-created workflow normally composes existing family
  pCIDs; it earns a new pCID only for a new interoperable durable meaning.
