# Ex1 local profile contract expansion

TE ID: TE-faviz

## Status

decided

## Decision status

Locked by DI-garis: retain the current reusable-until-expiry bearer-token
behavior; locally reject malformed or unauthorized input without requiring a
signed reply; and expand the five existing local draft profile documents in
place using a common contract template.

## Decision under test

How should Ex1 expand its five terse local profile documents into auditable
PromiseGrid-facing draft contracts without falsely presenting them as frozen
upstream specifications or changing behavior that belongs to later work? This
TE is the analysis prerequisite for `lubav.2`.

## Assumptions and trust model

- The PromiseGrid development guide is the alignment standard: a selected pCID
  owns payload and handler semantics, while frozen-spec conformance claims are
  explicit and bounded.
- Ex1's `order`, `pick_pack`, `accounting`, `shipment`, and `kernel_register`
  profiles are local draft/example profiles derived from `specdocs/`, not
  asserted upstream frozen specifications.
- Alice runs intake, Bob seller, Carol warehouse, Dave accounting, and Ellen
  carrier. The kernel forwards exact bytes by pCID and does not assess business
  promises.
- Current capability tokens are signed bearer tokens containing issuer,
  audience, pCID, kind, action, issued-at, expiry, and token ID. They are
  reusable until expiry; no redemption or deduplication ledger exists.
- This documentation task must not silently add replay prevention, a new
  evidence protocol, a registry, new wire actions, or an upstream-spec claim.

## Alternatives

### A. Complete local-draft contract per existing pCID, using one authoring template

Each existing profile gains its own status, payload grammar, valid kinds,
signing scope, capability semantics, validation/failure behavior, parent and
evidence links, and durable/transient statement. The template is documentation
only; no generic wire protocol is introduced.

### B. Shared capability/evidence protocol referenced by the business profiles

Existing profiles would define only their domain fields, while a new shared
profile owns tokens, evidence, and common failures.

### C. Keep the short profile headers and rely on Go code as the contract

The Markdown continues to name only the outer envelope; readers infer all
behavior from structs and services.

## Scenario analysis

### Normal fulfillment

Alice sends an order submit to Bob; Bob validates it, requests work from Carol,
Dave, and Ellen, then sends Alice a final. Every payload is selected by pCID,
signed, and linked to the prior work. A makes valid kinds, token audience, and
evidence links readable from the selected profile. B adds an unimplemented
protocol boundary. C makes the workflow run but leaves its wire meaning in
code rather than a contract.

### Invalid or expired authority

Mallory sends Bob malformed bytes or a token issued to another audience. Bob
must not perform business work. Current code validates pCID, proof, and token
claims before work; after a valid request, a defined failure path can produce a
signed final. A can state local rejection versus signed response precisely. B
does not improve it without behavior changes. C leaves silence, rejection, and
signed refusal ambiguous.

### Duplicate delivery and mixed-version peers

Carol may receive a repeated valid request, and a peer may select a pCID for a
changed profile with the same human-facing name. Current Ex1 has bearer tokens
reusable until expiry and no exactly-once promise. A documents that limit and
requires an unexpected pCID to be rejected. B requires a new redemption ledger
and compatibility protocol. C risks implying that token IDs prevent replay.

### Timeout and incomplete work

Ellen may remain silent until Bob's order deadline expires. Bob may emit a
signed `failed` carrier final. If Bob is silent past Alice's deadline, Alice
has only a local timeout observation, not proof of Bob's intent. A documents
the current profile response while deferring the durable observer record to
`lubav.4`. B improperly folds an evidence-design decision into this slice. C
hides the guide-required distinction.

### Long horizon, trust boundaries, and scale

Editing a local spec changes its pCID. Old bytes stay interpretable under the
old profile; new bytes must not be claimed compatible by name alone. Ex1's
deterministic role keys are fixture identities, not a production continuity
system. A makes status and evolution visible while retaining a small kernel.
B adds shared coordination overhead. C makes independent audit and forking
harder as the example grows.

## Conclusions

C is rejected because code is not an auditable protocol contract. B is rejected
for this slice because it changes the protocol surface and requires new
behavior decisions. A survives and is recommended: write one complete,
explicitly local-draft contract for each existing pCID using a common template.

The recommended scope documents existing signed bearer-token behavior, local
rejection of malformed/unauthorized input, and signed results/finals for
defined valid paths. It does not add redemption/replay guarantees and defers
durable observer timeout evidence to `lubav.4`.

## Decisions still requiring DF

1. **Capability reuse:** document the current reusable-until-expiry bearer
   tokens (recommended), or change Ex1 now to one-time redemption with a
   durable redemption ledger?
2. **Failure response:** specify local rejection for malformed/unauthorized
   input and signed results/finals only for valid profile paths (recommended),
   or require a signed refusal for every rejected input?
3. **Profile presentation:** expand the five existing `specdocs/*.md` files in
   place with a common template (recommended), or create separate contract
   documents and leave the current files as summaries?

## Implications for open work

- `lubav.2` can proceed after the three DF choices are locked.
- `lubav.3` remains responsible for any implementation-promise publication.
- `lubav.4` remains responsible for durable observer evidence of timeout,
  malformed input, refusal, and unknown-pCID receipt.
