# Ex2 PromiseGrid alignment plan

TE ID: TE-vodik

## Status

decided

## Decision status

Locked by DI-busiz: use the documentation-first sequence, publish a
local-draft scope declaration, decide evidence policy separately before runtime
changes, and add the testing guide with regression coverage.

## Decision under test

How should the existing Ex2 Grid Editor example be brought into current
PromiseGrid guide alignment without rewriting mechanics that already satisfy
the guide or presenting its repo-local protocol drafts as frozen upstream
contracts?

## Assumptions and trust model

- Ex2 currently derives four pCIDs from exact repo-local protocol documents:
  `live-document`, `live-awareness`, `document-metadata`, and
  `publish-document`.
- Its relay signs, verifies, stores, and forwards exact envelopes; browser and
  Neovim are local embodiments with CRDT replicas, not independent durable
  peer identities. The relay does not own canonical merged text.
- Ex2 passed `go vet ./...` and `go test ./...` on 2026-08-10, including its
  browser/Neovim interoperability test.
- Frozen specs by pCID and explicit implementation-promise claims are the
  future authoritative form. Ex2's local drafts remain provisional
  orientation, not interoperability claims.
- Alice runs a normal editor session; Bob operates a second relay; Mallory may
  send malformed or unsupported envelopes, copy stale artifacts, or mistake a
  relay record for shared proof.
- The current guide permits `grid([42(pCID), payload, proof])` as a local
  profile when the pCID-owned spec defines the slots. Ex2's slot-1 raw CBOR
  item is therefore not a gap by itself.

## Alternatives

### A. Documentation-first, then evidence policy and regression coverage

First make Ex2's four local contracts, pCIDs, scope, embodiment boundary, and
non-claims visible. Next decide unsupported/malformed-input evidence policy,
then add focused regression coverage and a testing guide. Finish with a
guide-facing documentation pass.

### B. Mechanic-first hardening

Immediately modify relay retention, rejection, or identity behavior, then
document whatever implementation results.

### C. Documentation-only declaration

Publish an inventory and scope declaration but defer unsupported-pCID policy
and regression coverage indefinitely.

## Scenario analysis

### Normal local collaboration

Alice opens the browser while Bob uses Neovim through the same relay. A makes
the current claim legible: the relay is the signing app identity, the two UIs
are embodiments, and local CRDT replicas converge from signed relay traffic.
B may improve internals but leaves readers unsure what actually claims to be
implemented. C improves orientation but cannot show that published claims stay
true as specs or code evolve.

### Second relay and a changed protocol draft

Bob runs a second relay and one local protocol document changes, producing a
new pCID. A first publishes an inventory derived from exact source bytes and
explicitly calls the profiles local drafts; later tests detect stale inventory
or scope records. B risks changing runtime behavior without an auditable
publication boundary. C describes the distinction but does not protect it from
drift.

### Unsupported or malformed envelope

Mallory sends an envelope that parses but selects an unsupported pCID, or sends
bytes that fail parsing or proof verification. The current relay stores some
exact bytes before rejecting unknown pCIDs, but the meaning, retention, and
replay policy are not documented or fully tested. A treats that as a separate
observer-evidence decision: retain exact local bytes where selected, record
only what the relay observed, and avoid claiming global pCID invalidity or
another agent's intent. B risks silently choosing a trust policy during a code
patch. C leaves operators unable to distinguish local rejection from durable
evidence.

### Long-horizon identity and embodiment change

Alice later switches between browser and Neovim, while Bob rotates a relay
key. A records the current boundary honestly: the relay key is the current
app-side signing identity; embodiment labels and participant IDs are local
presentation/routing fields, not signing-key continuity claims. It also leaves
key rotation for a separately scoped future decision. B can accidentally turn
UI labels into authority. C risks readers inferring that a browser-local name
is a portable peer identity.

### Scale and maintenance

A divides work into small, reviewable tasks: source-derived pCID consistency,
scope declaration, evidence policy, tests, testing guide, and final reader
walkthrough. B creates a broader hardening patch with unclear provenance. C is
cheap initially but causes documentation and runtime claims to drift.

## Conclusions

B is rejected because the review found no urgent broken core mechanic and
because changing rejection or identity behavior before defining its evidence
meaning would be backwards. C is rejected because current PromiseGrid guidance
requires implementation claims to be explicit and because an untested scope
record will drift.

A survives and is recommended. The proposed Ex2 alignment plan is:

1. Publish a source-derived four-profile pCID inventory and a local-draft
   implementation-scope declaration with explicit non-claims.
2. Decide the relay's malformed/unsupported-envelope evidence policy before
   changing its retention or rejection behavior.
3. Add regression coverage for source-derived inventory consistency, signed
   rejection/evidence behavior selected in step 2, replay/CAS retention, and
   existing cross-embodiment interoperability.
4. Add `docs/testing.md`, link it from the README, and state what each test
   layer proves.
5. Complete a final guide-facing README pass explaining local-draft scope,
   relay-versus-embodiment identity, reproducible operation, and artifact
   inspection.

## Decisions still requiring DF

1. **Alignment order:** adopt A's documentation-first sequence
   (recommended), harden mechanics first, or documentation-only declaration?
2. **Scope declaration:** publish a `CHANGELOG.md` local-draft
   implementation-scope record naming all four current pCIDs (recommended),
   or wait for frozen upstream specs?
3. **Evidence policy:** run a dedicated TE and DF for malformed and
   unsupported-envelope retention before code changes (recommended), or treat
   current relay behavior as sufficient without a decision record?
4. **Testing guide:** add the required Ex2 `docs/testing.md` and README link
   in the regression-coverage slice (recommended), or defer it to the final
   guide pass?

## Implications for open work

- A new Ex2 alignment TODO should be filed after these DFs are locked.
- Each step that changes protocol wording, persistence, identity, or runtime
  behavior must have its own TE and DI before implementation.
- Existing phase TODOs remain historical feature records; this alignment plan
  should not reopen completed feature work merely to rename terminology.
