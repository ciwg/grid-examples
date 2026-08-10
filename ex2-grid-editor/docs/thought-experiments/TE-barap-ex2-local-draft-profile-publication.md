# Ex2 local-draft profile publication

TE ID: TE-barap

## Status

decided

## Decision status

Locked by DI-ralit: publish a local-draft `CHANGELOG.md` scope declaration and
architecture inventory, distinguish relay identity from embodiment fields, and
list the four explicit non-claim boundaries.

## Decision under test

How should Ex2 publish the current derived pCIDs and implementation scope of
its four repo-local protocol profiles so readers can audit what the example
actually implements without confusing those drafts with frozen PromiseGrid
contracts?

## Assumptions and trust model

- `live-document`, `live-awareness`, `document-metadata`, and
  `publish-document` are embedded local Markdown specs; each pCID derives from
  its exact bytes at runtime.
- The current PromiseGrid guide treats frozen specs by pCID as the eventual
  contract reference. Draft specs and executable examples are provisional
  orientation until an explicit freeze and implementation claim exists.
- Ex2's relay is the current signing app identity; browser and Neovim are
  local embodiments. Their labels and participant IDs do not establish
  independent signing-key continuity.
- Alice is an operator deciding whether Ex2 matches her environment. Bob runs
  a second relay with a later checkout. Mallory cites a pCID from an old or
  altered checkout as if it proved current compatibility.
- The existing README, architecture guide, and protocol docs already explain
  mechanisms. They do not publish a derived-pCID inventory plus a bounded
  implementation-scope declaration in one auditable place.

## Alternatives

### A. `CHANGELOG.md` local-draft scope declaration plus README/design inventory

Create an Ex2 `CHANGELOG.md` that names all four current pCIDs, implemented
components, local-draft status, and explicit non-claims. Add the same
source-derived inventory to the design/README entry path and protect it later
with a regression test.

### B. A normal implementation-promise claim using the four draft pCIDs

Use the guide's future-facing frozen-spec claim form now, with a
`partially-implements` status against the local profiles.

### C. Protocol docs only

Leave each profile's draft status in its own protocol document and do not add
an implementation-scope declaration or reader-facing inventory.

## Scenario analysis

### Normal local browser and Neovim session

Alice runs browser and Neovim against one relay. A lets her see that both are
embodiments of the one signing relay app, which four local profiles it uses,
and which parts are only local implementation scope. B makes the same session
sound like a frozen-spec claim it cannot support. C requires her to infer the
actual composition from source and implementation paths.

### Second relay with a changed profile document

Bob changes `live-document.md`, which produces a new pCID, then tries peer
exchange with Alice's checkout. A makes the changing identity visible in the
inventory and makes a mismatch an explicit local-draft compatibility question.
B falsely suggests a stable implementation promise. C leaves no single
published record that tells either operator which exact draft each relay used.

### Stale artifacts and mixed versions

Mallory presents a recorded envelope selecting an old pCID. A allows an
operator to compare it to the current inventory while preserving the principle
that old bytes remain evidence rather than automatically current traffic. B
invites a false claim that old and new drafts share a frozen contract. C
provides no reader-facing comparison point.

### Long-horizon evolution and eventual freezing

An upstream profile later freezes with its own doc-CID. A can add a later,
separate formal implementation-promise entry while preserving the historical
local-draft declaration. B must unwind an overclaim. C lacks an explicit
baseline from which to state what changed.

### Trust and identity boundary

Alice sees a browser display name and a relay key ID. A distinguishes the
relay signing identity from embodiment-local presentation fields and does not
mistake either for a general peer-identity or key-rotation system. B and C
leave that distinction less visible, increasing the chance that friendly UI
labels are treated as protocol authority.

## Conclusions

C is rejected because the development guide requires implementation scope to
be explicit rather than inferred from a repository path, branch, or demo. B
is rejected because the four documents are local drafts, not frozen upstream
specs, so the normal implementation-promise claim form would overstate their
status.

A survives and is recommended: publish a clearly labeled local-draft scope
declaration in `CHANGELOG.md`, add an exact source-derived profile inventory to
the guide-facing documentation, distinguish relay signing identity from local
embodiment labels, and enumerate what Ex2 does not claim. Later regression
coverage must derive the pCIDs from the source documents and compare the
published values.

## Decisions still requiring DF

1. **Publication form:** use the local-draft `CHANGELOG.md` scope declaration
   plus reader-facing inventory (recommended), use a formal implementation
   promise against drafts, or rely on protocol docs only?
2. **Inventory location:** place the primary inventory in
   `docs/architecture.md` and link it from the README (recommended), or put
   the complete inventory directly in the README?
3. **Scope boundary:** explicitly state that the relay signing key is the
   current app identity while browser/Neovim fields are local embodiment or
   presentation data (recommended), or omit identity-boundary wording?
4. **Non-claims:** explicitly exclude frozen-spec conformance, independent-peer
   interoperability, general key-rotation/role-continuity policy, and a
   portable runtime/storage contract (recommended), or list only frozen-spec
   conformance?

## Implications for open work

- `sojot.1` can proceed after these DFs are locked and recorded in
  `TODO-sojot`.
- Publishing source-derived pCIDs requires a later regression check; this TE
  does not authorize profile-document edits that would change them.
- `sojot.2` remains the required place to decide malformed and
  unsupported-envelope evidence policy.
