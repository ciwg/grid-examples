# Grid Exercises

This repository is a set of working application exercises for learning how to
build on the PromiseGrid: begin with a pCID-selected contract, keep durable
evidence distinct from local projections and workflow policy, preserve exact
bytes where meaning is not known, and make transport and user-interface
embodiments serve the contract rather than redefine it. The exercises are
useful as concrete reference implementations because they pair runnable
applications with specifications, decision records, tests, and explicit limits
on what each application does **not** claim.

Each exercise below links to its own README, which is the authoritative guide
for setup, architecture, verification, implementation claims, and deferred
work. They demonstrate different application domains and maturity points; they
are not interchangeable protocol specifications or a single product.

## Exercise index

### [Ex1 — Order Flow](ex1-order-flow/README.md)

Ex1 is a multi-agent order-fulfillment example: independent seller,
warehouse, accounting, carrier, intake, collector, and kernel roles exchange
signed messages and retain role-local evidence. Its Docker demo makes the
message graph and exceptional refusal/timeout observations inspectable.

It implements five local-draft profiles. It does not claim frozen-spec
conformance, independent-peer interoperability, or that a local observation
settles another participant's intent. Source: `DI-josir`; `DI-motiv`.

### [Ex2 — Grid Editor](ex2-grid-editor/README.md)

Ex2 is a shared-document editor embodied in both a browser and Neovim. It
uses signed, pCID-addressed envelopes, content-addressed persistence, local
CRDT replicas for convergence, and distinct document, awareness, metadata,
and publish families.

Its four profiles are repo-local drafts, and its browser and Neovim paths are
local embodiments rather than separate cryptographic identities. Relay-local
observations are evidence at that relay, not shared proof. Source: `DI-bafar`;
`DI-nilas`; `DI-todav`.

### [Ex3 — Grid Editor with WebSocket Carriage](ex3-grid-editor-websocket/README.md)

Ex3 develops the shared-editor example with WebSocket-preferred live document
and awareness carriage, while retaining explicit pCID-selected meanings and
HTTP metadata/publish surfaces. It includes browser, Neovim, relay, and
headless recovery evidence for late joins and stale snapshots.

WebSocket is carriage, not a fifth protocol family; the four profiles remain
repo-local drafts. Its remote mutation bootstrap and capability arrangement is
application-local admission policy, not a general identity or authorization
API. Source: `DI-vipat`; `DI-bitus`; `DI-gofut`.

### [Ex4 — Bug Tracker](ex4-bug-tracker/README.md)

Ex4 is a browser-first bug tracker with an engineer CLI. Issue reports,
lifecycle updates, and attachment references are signed, pCID-selected
promises; accepted bytes are retained in CAS and the visible issue state is a
projection of durable history.

Its `events.jsonl` workflow history and enrollment/role policy are local to
this application. It does not claim global identity, delegation, federation,
or authorization, and its first release uses a single built-in team. Source:
`DI-kolaf`; `DI-ninul`; `DI-gofub`.

### [Ex5 — Operational Knowledge System](ex5-operational-knowledge-system/README.md)

Ex5 is a durable operational-memory application for procedures, training,
maintenance, inventory checks, evidence, approvals, and later retrieval. It
has browser, CLI, and optional Neovim embodiments over a local runtime, and
extracts reusable record, transport, and storage substrate under
`promisegrid/`.

This is a completed scoped runtime slice, not a generalized ERP or a complete
PromiseGrid product boundary. Chrome native messaging is verified; Chromium is
explicitly deferred, and richer Neovim work remains future product work.
Source: `DI-lavek`; `DI-rasok`; `DI-punek`.

### [Ex6 — Operational Knowledge Agent Runtime](ex6-operational-knowledge-agent-runtime/README.md)

Ex6, OKAR, is a standalone runtime for first-party and installed operational
knowledge packages. Its 27 built-in families have immutable specifications and
fixed CIDv1 pCIDs; durable package evidence uses canonical Grid CBOR, and
workflows compose family contracts rather than becoming pCIDs themselves.

Unknown-family records are retained as exact bytes without inferred meaning.
Author signatures, relay carriage, route availability, approvals, and workflow
execution remain separate local-policy questions; Ex6 does not claim global
identity, universal authority, consensus, or trust from installation alone.
Source: `DI-sidoh`; `DI-jusij`; `DI-solan`.

### [Ex7 — Makerspace Stewardship](ex7-makerspace-stewardship/README.md)

Ex7 is a makerspace record-agent example. It validates externally signed,
pCID-selected makerspace records against participant root/device history,
recovery, peer cards, exact-byte carriage, and a local recognition policy,
then derives a local view only from records it recognizes.

The browser submits already signed bytes; a selected name, session, or account
is not author evidence. Browser signing, account identity, a running relay,
and portable governance are intentionally outside the present implementation.
Source: `DI-tohak`; `DI-piruf`; `DI-kasaz`; `DI-sisad`.

## Reading and verification

Start with the exercise closest to the kind of application you are building,
then read its protocol/specification inventory, architecture, implementation
claims, and testing guide before reusing a pattern. The current cross-exercise
verification commands and evidence paths live in each exercise's testing guide;
completed evidence and explicit deferrals do not become universal conformance
claims.
