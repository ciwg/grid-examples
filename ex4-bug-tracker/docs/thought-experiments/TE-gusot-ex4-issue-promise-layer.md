# Ex4 bounded issue-promise layer

TE ID: TE-gusot

## Status

decided

## Decision status

Locked by DI-ninul: use three signed pCID-selected issue-promise profiles,
distinct local adapter signing keys, append-only accepted artifacts with local
rebuildable projection and bounded rejected evidence, and no remote exchange.
This TE is `kakon.2` and does not change issue behavior, event storage,
identity handling, or network behavior.

## Decision under test

What bounded pCID-selected signed issue-promise layer should Ex4 add so browser and CLI can submit voluntary issue-related promises, while the HTTP server remains a local adapter/projection and `events.jsonl` remains distinct local history?

## Assumptions and trust model

- Current `IssueEvent` records are unsigned local server history. Fixed identities and role checks are not cryptographic agent identity or universal authorization. Source: `DI-nibuh`.
- The only top-level distributed semantic action is `promise`. Reporting, commenting, assigning, status change, and attachments are pCID-defined payload meanings or local mechanics, not top-level actions. Source: `DI-mosoj`.
- Alice operates Ex4; Bob and Carol submit through browser or CLI; Mallory can alter, replay, omit, or send an unsupported artifact.
- The first layer makes no cross-tracker federation, generalized delegation, revocation, recognized-role, or remote-exchange claim.

## Alternatives

### A. One broad issue-workflow profile

One pCID-selected `issue-workflow` promise profile has a tagged payload for report, comment, assignment, status, and attachment references. Each local agent signs its own envelope; the server validates and locally projects it.

### B. Three semantic promise profiles

Separate pCID-selected profiles cover issue report, issue lifecycle update (comment, assignment, or status), and attachment reference. Each local agent signs its own envelope. The server stores accepted artifacts append-only and content-addressably, uses a rebuildable local projection/index, and records rejected input only as bounded local observation or diagnostics. Browser and CLI remain local adapters; no remote exchange occurs in the first slice.

### C. One profile per current event name

Independent pCIDs mirror `created`, `commented`, `assigned`, `status_changed`, and `attachment_added`.

### D. Server-signed or unsigned event promotion

Reuse `events.jsonl` as the promise log, unsigned or server-signed on users' behalf.

## Scenario analysis

### S1 — Bob reports a new issue

A is compact but every consumer must understand a large tagged union. B gives reporting a stable pCID-defined meaning. C adds a pCID for every local storage event. D turns a server record into an apparent promise without Bob's signature.

### S2 — Carol comments, assigns, or changes status

A centralizes updates but couples simple comments to a broad schema. B groups updates to an existing issue under one explicit lifecycle meaning while retaining a payload discriminator. C makes evolution and mixed-version support expensive. D lets server role checks masquerade as Carol's promise.

### S3 — Mallory sends malformed, replayed, or unsupported input

A and B can verify pCID and signature before local projection; B can call an unsupported profile local non-support. C expands the unsupported surface. D cannot distinguish a user-originated promise from arbitrary HTTP data. A rejection remains local evidence, not a universal statement about Mallory.

### S4 — Alice restarts and rebuilds state

A and B retain accepted artifacts and rebuild a local projection. B's separate report, update, and attachment meanings permit selective projection and migration. C requires many branches. D confuses old local history with new accepted artifacts.

### S5 — Carol runs a future compatible tracker

A compares one pCID but couples unrelated lifecycle evolution. B compares exactly the report, update, and attachment meanings a peer supports; absence is local non-support, not global invalidity. C adds coordination points. D supplies no portable signed contract. This slice still does not exchange remotely.

### S6 — Identity, authority, and scale evolution

A and B keep a later key-binding change separate from built-in role labels. B best separates signer statement, server policy, accepted artifact, rejected observation, and historic local log. C freezes UI taxonomy as protocol taxonomy. D attributes promises to the server or no signer. B adds three specs and modest verification/storage cost; C adds needless test and migration cost.

## Conclusions

D is rejected: neither unsigned records nor a server signature on behalf of a user establishes that user's voluntary promise. C is rejected: local event names are implementation details, not automatically peer protocols.

A and B survive. B is recommended because report creation, lifecycle updates, and attachment references have distinct stable meanings without producing a pCID per UI or storage event. It maintains clear provenance among signer promise, local policy, accepted artifact, rejected evidence, and existing history. A remains viable only when one broad schema is deliberately preferred over compatibility and evolution boundaries.

## Decisions still requiring DF

1. **Profile granularity:** use B's report/lifecycle-update/attachment pCIDs (recommended), one broad issue-workflow pCID, or one pCID per local event?
2. **Promiser identity:** give each built-in local identity a distinct local signing key explicitly bound to that adapter identity (recommended), have the server sign on behalf of users, or keep artifacts unsigned?
3. **Accepted/rejected storage:** use append-only content-addressed accepted artifacts plus a rebuildable local projection/index and bounded rejected observations/diagnostics (recommended), append everything to `events.jsonl`, or store only projected mutable rows?
4. **First-slice adapter/exchange:** have browser and CLI submit signed promises to the local server adapter without remote exchange (recommended), let the server synthesize promises from commands, or add cross-tracker exchange now?

## Implications for open TODOs and pending DIs

- `kakon.2` must lock all four choices before profile names, paths, schemas, helpers, or implementation details are selected.
- `kakon.3` must preserve `events.jsonl` as historical local application data or use an explicitly selected compatibility migration; it must not silently relabel old events as signed promises.
- `kakon.4` must test accepted signed artifacts, local rejected evidence, server projection, and browser/CLI adapters distinctly.
