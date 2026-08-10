# Ex4 local-history and PromiseGrid-evidence scope

TE ID: TE-mofol

## Status

decided

## Decision status

Locked by DI-nibuh: publish `CHANGELOG.md` scope/non-claims with
architecture/README navigation; classify `events.jsonl` as durable local
application history; and state that signed pCID-selected issue-promise
artifacts are a separately specified future layer. This TE is `kakon.1` and
does not change event storage, issue behavior, identity handling, or network
behavior.

## Decision under test

How should Ex4 document the boundary between its current append-only
`events.jsonl` application history and the bounded signed, pCID-selected
issue-promise/evidence layer selected by `DI-gisor` but not yet implemented?

## Assumptions and trust model

- Today, one local HTTP server validates built-in roles and writes
  `events.jsonl`; browser and CLI clients use that server. This is durable
  application history, not independently verifiable shared evidence.
- The selected future direction is a bounded issue-promise layer. Before its
  own TE and DF, Ex4 has no selected pCID documents, signing identity,
  accepted-artifact format, rejected-artifact format, or cross-tracker path.
- Alice operates one tracker. Bob uses its browser or CLI. Carol may later use
  a compatible implementation. Mallory may submit malformed or unauthorized
  input to Alice's HTTP service.
- Documentation must not promote local role checks to universal authorization,
  a server projection to another actor's intent, or an event record to a
  PromiseGrid promise merely because both are append-only.

## Alternatives

### A. Explicit two-layer scope and non-claim publication

Create a concise implementation-scope/non-claim document and extend the
architecture/README navigation. Name `events.jsonl` durable local application
history. State that future accepted pCID-selected signed promise artifacts, if
implemented, will be a separate layer with its own evidence boundary. State no
cross-tracker interoperability, agent identity, delegation, or role continuity
claim until that layer is specified and implemented.

### B. Reclassify current events as PromiseGrid evidence

Describe the existing append-only event log as promise/evidence now, with the
server's built-in roles and transition checks treated as the protocol.

### C. Hide the boundary until the promise layer ships

Leave the current broad durable-workflow wording and add no explicit
implementation scope or non-claims until a later implementation slice exists.

## Scenario analysis

### S1 — Alice inspects a normal local issue timeline

A lets Alice see the log as durable local history and understand that it
projects the server's current issue view. B lets a local server record be read
as a peer-visible promise without a selected pCID or signature. C leaves the
reader unable to distinguish durable workflow mechanics from protocol claims.

### S2 — Bob uses browser and CLI against Alice's service

Under A, browser and CLI are local adapters of one application workflow; their
shared display is not decentralized corroboration. B incorrectly makes a
common server response look like two independent agents agreeing. C omits the
boundary, so readers infer more interoperability than exists.

### S3 — Mallory sends malformed or unauthorized input

A can say the server's rejection and any local diagnostic are facts from
Alice's service vantage only. B risks presenting an HTTP denial as a universal
statement about Mallory or a protocol-level refusal. C offers no durable
interpretation rule for the reader.

### S4 — Carol runs another tracker or a mixed-version implementation

A says Carol's tracker has no promised compatibility until future pCIDs match
and exchange rules are implemented. B promises compatibility from similar
event names and role labels without a wire contract. C postpones a needed
non-claim and makes accidental compatibility assumptions more likely.

### S5 — Future signed issue promises and migration

A preserves `events.jsonl` as historical application data and provides a clean
new layer for accepted signed artifacts without rewriting old records. B makes
migration ambiguous because old unsigned server events were named as protocol
evidence. C may be mechanically reversible later, but creates misleading
reader expectations in the meantime.

### S6 — Scale and operational complexity

A adds concise documentation and focused tests only. B appears cheap but
creates a large future correction cost in trust and provenance claims. C is
cheap today but leaves the alignment work incomplete and hard to audit.

## Conclusions

B is rejected: append-only storage does not itself make an event a signed
promise or shared evidence. C is rejected: Ex4's current central boundary must
be explicit before a PromiseGrid implementation phase begins.

A survives and is recommended. It is the PromiseGrid-aligned way to describe
the current system honestly while reserving pCID-selected signed artifacts for
the separately locked `kakon.2` implementation design.

## Decisions still requiring DF

1. **Scope publication:** create a concise `CHANGELOG.md` implementation scope
   and non-claim declaration, with architecture/README links (recommended),
   place all scope content only in `docs/architecture.md`, or defer publication
   until code exists?
2. **Current event classification:** call `events.jsonl` durable local
   application history (recommended), relay-local PromiseGrid evidence, or a
   signed promise ledger?
3. **Future-layer wording:** state that pCID-selected signed promise artifacts
   are a separately specified future layer with no current interoperability or
   identity/delegation claim (recommended), imply they are already active, or
   omit the future boundary?
4. **Reader navigation:** link the scope declaration and architecture from
   README (recommended), link only architecture, or keep scope unlinked?

## Implications for open TODOs and pending DIs

- `kakon.1` can publish the selected scope/non-claim documents after DF.
- `kakon.2` must not reuse `events.jsonl` as the accepted promise artifact
  merely for convenience; its TE must select pCIDs, signatures, and evidence
  rules.
- `kakon.3` must preserve the distinction between local application history,
  local diagnostics, and any later accepted promise artifacts.
