# Ex1 local-draft implementation publication

TE ID: TE-birut

## Status

decided

## Decision status

Locked by DI-josir: publish a local-draft implementation-scope declaration in
`ex1-order-flow/CHANGELOG.md`, explicitly separate from a frozen-spec
implementation-promise claim.

## Decision under test

How should Ex1 publish its implementation scope while its five profile pCIDs
are explicitly local draft/example contracts rather than frozen upstream
PromiseGrid specification doc-CIDs? This TE is the prerequisite for `lubav.3`.

## Assumptions and trust model

- The development guide reserves formal implementation-promise claims for
  explicit frozen spec doc-CIDs and expects them in an implementation-side
  `CHANGELOG.md`.
- Ex1 currently implements five named local draft profiles and must not claim
  upstream interoperability or frozen-spec conformance.
- A reader needs to distinguish protocol-facing behavior from host-local Docker
  topology, deterministic fixture keys, storage layout, collector operations,
  and other example mechanics.
- The publication must remain useful across profile-pCID changes and later
  migration to frozen upstream specs.

## Alternatives

### A. Publish a `CHANGELOG.md` local-draft scope declaration, explicitly not a frozen-spec conformance claim

The record lists the five local profile pCIDs, their implemented components,
the draft/non-interoperability boundary, and host-local behavior that is not
claimed. It reserves the formal guide claim fields for a future frozen-spec
entry.

### B. Publish a normal `partially-implements` implementation-promise entry for the local pCIDs

The entry would reuse the guide's formal claim form even though the listed
documents are local drafts.

### C. Publish nothing until upstream specs freeze

Ex1 would defer all scope publication until a frozen upstream spec is available.

## Scenario analysis

### Current demo adoption

Alice wants to reproduce the local order-flow demo and determine which pCID
profiles its handlers speak. A gives Alice that information and prevents her
from treating local fixtures as a portable standard. B makes the claim format
look stronger than its evidence. C leaves scope implicit in source layout.

### Independent or mixed-version peer

Bob has an older Ex1 checkout with different profile bytes. A says that a
changed local spec has a new pCID and no compatibility is implied; Bob must
assess the named local bytes. B may imply a stable conformance relationship
that does not exist. C provides no visible migration boundary.

### Long-horizon graduation

Carol later freezes or adopts an upstream spec. A permits a subsequent normal
`implements` or `partially-implements` entry naming that exact frozen doc-CID
and clearly separates the prior local-draft history. B muddies which entry was
actually an upstream claim. C loses useful provenance for the example stage.

### Trust and operational boundary

Mallory can copy Ex1's public deterministic fixture keys, Docker setup, or
collector layout; none proves recognized production-role continuity. A lists
those as deliberately unclaimed host-local mechanics. B risks readers treating
the packaging as protocol authority. C does not explain the distinction.

## Conclusions

B is rejected because it would use a frozen-spec claim format for documents
that Ex1 itself calls local drafts. C is rejected because the guide requires
implementation scope to be legible rather than inferred from a branch or demo.

A survives and is recommended. Add `ex1-order-flow/CHANGELOG.md` with a clearly
labelled local-draft implementation-scope declaration. It must name the five
current pCIDs and components, state that it is not a frozen-upstream-spec
conformance claim, and list explicitly unclaimed host-local behavior. A later
frozen-spec adoption receives a separate formal implementation-promise entry.

## Decision still requiring DF

Should Ex1 adopt Alternative A, a local-draft scope declaration in
`ex1-order-flow/CHANGELOG.md` (recommended), or defer publication entirely
until it implements an explicit frozen upstream spec (Alternative C)?

## Implications for open work

- `lubav.3` can proceed after the publication choice is locked.
- No new profile wire semantics, handler behavior, or pCIDs are introduced.
- `lubav.7` should link the resulting publication record from guide-facing
  documentation after the remaining alignment work is complete.
