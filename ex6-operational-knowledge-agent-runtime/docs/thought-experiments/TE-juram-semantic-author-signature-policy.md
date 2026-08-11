# Semantic Author Signature Admission Policy

TE ID: TE-juram

## Status

decided

## Decision under test

Whether Ex6 should admit unsigned canonical package records after the
unreleased clean re-creation, and what local policy governs semantic author
signatures separately from relay-carriage signatures.

## Assumptions and scope

- Ex6 package records are canonical Grid CBOR bytes selected by frozen family
  pCIDs; the symbolic JSON compatibility path is intentionally absent under
  TE-hofor and DI-sidoh.
- Local append currently signs an unsigned record with the runtime key.
- Import verifies a supplied author signature, but admits an unsigned record.
- Relay-carriage signatures authenticate the peer that carried bytes; they do
  not establish the semantic author of those bytes.
- Alice is a local author, Bob imports a relay batch, Carol installs an
  external package, and Mallory supplies altered, unsigned, or replayed bytes.
- This decision covers package-record admission only. Route-promise and
  workflow-lifecycle protocols retain their separately declared scopes.

## Alternatives

### A. Require a valid semantic author signature for every canonical package record

Local append signs records when needed. Import and append reject unsigned,
incomplete, malformed, or invalid author evidence. The current local runtime
key is the author identity mechanism until a richer author-identity protocol is
separately specified.

### B. Add explicit local signature-admission modes

Persist a local policy that selects `require`, `allow-unsigned`, or perhaps
family-scoped exceptions. The default may be strict or permissive.

### C. Retain the current optional behavior

Sign local records, verify signatures when present, but admit unsigned records
on import and append.

## Scenario analysis

### Normal local authoring

Under A, Alice's local record is always signed before history storage and
contains durable authorship evidence. Under B, the outcome depends on an
operator-selected policy and must be shown in diagnostics. Under C, local
authoring remains signed today, but the protocol permits an unsigned caller to
reach durable history whenever signing is disabled or bypassed.

### Relay import and altered carriage

Under A, Bob requires both the independently applicable relay policy and a
valid semantic author signature before durable admission. Mallory cannot turn
relay carriage into unsigned semantic evidence. Under B, Bob can intentionally
relax this boundary, but the policy itself becomes durable local evidence that
must be exported, documented, and tested. Under C, a correctly signed relay
batch may carry an unsigned semantic record, leaving Bob unable to distinguish
"Bob's peer carried this" from "an author made this statement."

### External package authors

Under A, Carol's package must emit author-signed canonical records or submit
them to a local signing boundary that truthfully identifies the local author.
This makes external author identity explicit but may require package tooling.
Under B, Carol's package can operate under a scoped exception, increasing the
policy surface and the chance of accidental broad trust. Under C, external
packages can silently emit unsigned durable statements despite the canonical
protocol re-creation.

### Failure, corruption, and incomplete writes

Under A, missing one of key ID, public key, or signature is a deterministic
admission failure; exact bytes can remain outside durable package history for
diagnosis. Under B, failures must also identify which policy mode applied.
Under C, an incomplete signature set is rejected, but a fully absent set is
accepted, creating an avoidable ambiguity between corruption and policy.

### Mixed versions and migration

Under A, Alice and Bob must use the current unreleased canonical contract; this
matches the clean-break premise of TE-hofor. Under B, an allow-unsigned mode
reintroduces an operational compatibility path even though no released record
format requires it. Under C, the old compatibility rationale remains embedded
in active admission behavior without an explicit version or expiry boundary.

### Long-horizon trust and evolution

Under A, a later richer identity or delegation model is a new pCID-defined
protocol or an explicit local-policy extension; old signatures remain durable
evidence under their stated semantics. Under B, policy history can be useful
but must be carefully versioned and audited. Under C, consumers must always
treat signature absence as an ambiguous normal case, weakening evidence review.

### Scale and operational complexity

A has the smallest steady-state implementation and clearest operator story:
canonical package records require valid author evidence; relay evidence remains
separate. B adds useful flexibility but requires a persisted policy model,
CLI, documentation, export semantics, and tests before it is trustworthy. C
has little immediate code cost but transfers recurring ambiguity to every
consumer and reviewer.

## Conclusions

Alternative C is rejected for the unreleased clean-recreation scope: it keeps
an undocumented compatibility admission path after the legacy record format was
removed. Alternative B survives only if Ex6 has a concrete current operational
need for scoped unsigned admission; no such need appears in the active TODO
set. Alternative A is the recommended current scope: require and verify a
complete semantic author signature for every canonical package record while
keeping relay-carriage signatures separate and local policy authoritative for
what a valid signature means.

## Decision status

Locked: Alternative A. DI-damab requires a valid semantic author signature for
every canonical package record. Alternative B is deferred because Ex6 has no
current demonstrated need for unsigned admission modes.

## Implications for TODOs and pending DIs

- Supersede DI-sovem's optional unsigned-record compatibility constraint if A
  is selected.
- Update `records/`, `grid/`, `kernel/`, installed-package fixtures,
  documentation, tests, and `TODO-sibok` to reflect the selected boundary.
- The final TODO audit should mark the optional-legacy signature behavior as
  resolved or explicitly deferred according to the DF result.
