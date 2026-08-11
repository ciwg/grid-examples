# Ex7 source-grounded recreation roadmap

## Status

Planning guide for DI-tohak. This is not an implementation claim.

## Historical starting point

At filing, Ex7 demonstrated local makerspace workflow behavior through a JSON
HTTP process and JSONL replay. The recreation began with named contracts and
their exact byte identities, not with a generic runtime, account, relay,
browser, or identity product. Source: DI-tohak.

## Current implementation status

The live embodiment now uses canonical signed records in `records.frames`,
projects the four makerspace families only after participant-history and local
policy assessment, validates 2-of-3 recovery and peer-card-linked exact-byte
carriage, and supports terminal-to-agent local approval. The Chrome proof is
`scripts/run-two-agent-browser-proof.sh`; it creates a disposable two-agent
session and records its evidence root. The current implementation claim and
explicit non-claims are maintained in `docs/promisegrid-implementation-claims.md`.
Source: DI-tohak; DI-piruf; DI-likoh; DI-kasaz; DI-sisad; DI-hibok; DI-fuzar.

## Ordered slices

### 1. Verify protocol artifacts

Recalculate and test every listed makerspace pCID against the exact immutable
specification bytes. Verify the record-profile document identity that a later
implementation claim will cite. If a value does not match, correct the
registry/claim relationship before writing runtime behavior.

### 2. Put the four contracts on the live record path

Replace `Event{Type: ...}` persistence with exact Ex7 record validation,
durable append, replay, and family-specific projection. An observation plus a
safety hold becomes the distinct linked meanings specified by the two family
contracts. A return links to its exact loan record. The UI/API is only a local
request path; it is not the durable protocol format.

This slice must prove canonical bytes, malformed and partial-write failure,
duplicate handling, unknown-pCID retention without interpretation, and replay
from durable bytes. It must not make a browser or local HTTP process the author.

### 3. Lock the signing embodiment

Before the UI emits semantic author evidence, complete a separate
source-grounded decision for the participant signing/ingress embodiment. It
must say which key may sign, what continuity evidence exists, and what the UI
may do when the signer is unavailable. Existing account sessions remain local
bootstrap or recognition inputs, not signatures.

### 4. Add byte carriage only when specified

If Ex7 needs relay, feed, or direct exchange, define its named carriage
contract and test it after local verified record ownership exists. Carriage may
retain and move exact bytes; it does not decide makerspace semantics or local
trust.

### 5. Publish evidence, not aspiration

After the live path passes conformance tests, add `docs/testing.md`,
implementation claims, `CHANGELOG.md`, architecture documentation, and an
opt-in end-to-end proof. Each implementation claim names exact frozen document
identities, scope, supported embodiment, dependencies, unsupported cases, and
deferrals.

## Explicit deferrals

No current Ex7 claim covers browser signing, account-backed authoring, a
running relay/discovery service, blob availability, portable
membership/governance, a global trust registry, or a PromiseGrid-wide
envelope/runtime standard. Each requires its own contract and evidence if Ex7
later needs it.
