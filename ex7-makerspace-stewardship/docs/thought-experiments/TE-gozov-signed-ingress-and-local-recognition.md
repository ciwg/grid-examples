# Signed Ingress and Local Recognition for the First Ex7 Record Path

TE ID: TE-gozov

## Status

decided

## Decision under test

How the first live Ex7 record/projection path accepts semantic author evidence
without making a browser, account, shared HTTP process, relay, or self-declared
display label into the author.

This TE narrows TODO `bubuz`, slice 2, after TE-biban / DI-tohak. It is not a
decision about key generation, key custody, account login, device delegation,
key recovery, or byte carriage.

## Assumptions and scope

- The record store can now append and replay exact canonical signed records,
  but it deliberately does not assign their family semantics.
- `Record.Signer` is presentation/bootstrap text. A valid signature over the
  same record does not by itself prove that a signer label is an accepted
  makerspace participant or steward.
- Ex7's first protocol set consists of the frozen observation, safety
  disposition, off-site loan, and off-site return families. Their projection
  must operate only on known family pCIDs and validated payloads.
- A browser may submit a proposed action or exact bytes, but it cannot sign as
  Alice merely because it submits `"alice"` or has an account session.
- Alice and Carol may use independently held signing material outside this
  slice. The first Ex7 runtime need only verify exact signed bytes and apply a
  stated local recognition policy; it does not create or retain their private
  keys.
- Mallory may submit malformed bytes, a valid unknown-pCID record, a valid
  record from an unrecognized key, or a valid record whose signer text claims
  to be Carol.

## Alternatives

### A. External signed ingress with injected local recognition policy

The Ex7 agent accepts exact record bytes only after canonical parse and
signature verification. It retains valid bytes for unknown pCIDs or
unrecognized keys, but projects a known-family record only when an injected
local policy recognizes the full public-key fingerprint for the relevant local
participant/role context. The policy is ordinary local bootstrap/configuration
for this first embodiment, not a portable membership protocol.

The HTTP/UI surface may submit exact signed bytes to this ingress. Existing
browser workflow routes become unsigned draft/request behavior until a later
embodiment can obtain a participant signature. No route receives a private key
or turns an account/session into a signature.

### B. Process-generated signatures for selected fixture labels

The process retains private keys for Alice, Carol, and Dave and signs the
records that its browser routes request.

### C. Signature-valid records project by signer label

Any valid record with `Signer: "carol"` receives the local steward effect
because the label matches the fixture authority list.

## Scenario analysis

### Normal operation

Under A, Alice signs a loan record using signing material outside the local UI
process. The agent verifies the exact bytes, recognizes Alice's public-key
fingerprint in its local policy, stores the bytes, and projects the loan. Carol
signs an inspection/clearance record; the agent separately recognizes Carol's
key for the configured steward context before projecting a clear.

B makes the process the real signer even while its UI says Alice or Carol.
C lets Mallory create a fresh key, write `Signer: "carol"`, and obtain a
clearance effect. Neither preserves the author-evidence boundary.

### Failure, corruption, and incomplete input

Under A, malformed or invalid-signature bytes never enter the record store.
Valid unknown-family or unrecognized-key bytes may be retained intact but do
not alter the known makerspace view. If the recognition policy is unavailable
or malformed, startup/admission fails rather than accepting an unlabeled
fallback. The framed store still fails closed on incomplete durable frames.

B can continue signing while hiding the wrong author. C continues projecting
an attacker-controlled label. Both create misleading durable meaning, not
merely a recoverable operational error.

### Concurrent participants and mixed versions

Alice's agent can retain a valid record from a new family without projecting
it. Carol's agent can recognize a different locally configured set of keys or
roles and reach a different local view from the same bytes. A new binary may
add a known family while old binaries retain its exact records. No byte
carrier must decide which agent's recognition policy is correct.

B makes every participant dependent on the process's signing arrangements.
C has no key-to-role evidence boundary across participants.

### Long-horizon evolution

A keeps the first record conversion small. A later pCID-defined continuity,
revocation, membership, or delegation protocol may replace or supplement the
injected local policy without rewriting existing record bytes. A future account
or browser embodiment may request an external signature, but its existence is
not assumed by this path.

B requires a later migration away from process-held keys. C requires a later
retrofit from labels to keys and leaves historical projected effects difficult
to assess.

### Trust and scale

A makes the local policy explicit: a public-key fingerprint is recognized in a
specific makerspace context, and signature validity is separate from whether a
projection should occur. The policy can remain small for the first product
embodiment and does not claim universal membership, reputation, or identity
recovery. The operational cost is a clear configuration surface and tests for
recognized, unrecognized, and role-mismatched keys.

B gives an implementation component impersonation power. C treats mutable
presentation text as authority. Their smaller apparent setup cost is paid as a
false evidence claim.

## Conclusions

B is rejected because it would attribute process-made signatures to people.
C is rejected because a signer label is not cryptographic author or role
evidence.

A survives and is recommended. It is the smallest path that lets the framed
store become the live durable history without deciding how private keys are
created, where they live, how a browser reaches them, or how bytes travel
between agents.

## Output to decision framing

**A — accept externally signed exact records; retain valid-but-unrecognized or
unknown-family bytes; project only known-family records whose full public-key
fingerprints satisfy an injected local recognition policy. Browser/account
routes remain drafts or exact-byte ingress, never signing routes.**

## Decision status

locked: Alternative A, external signed ingress with injected local recognition
policy, by DI-piruf in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`.

## Implications and future work

- The new runtime needs a testable local recognition-policy interface and
  known-family payload validators.
- The browser workflow must visibly distinguish an unsigned request/draft from
  a stored signed record.
- Any persistent policy format, key-continuity protocol, account bootstrap, or
  relay/feed design needs its own decision rather than being inferred here.
