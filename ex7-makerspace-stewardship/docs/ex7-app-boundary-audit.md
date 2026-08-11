# Ex7 PromiseGrid App-Boundary Audit

## Status

Historical audit baseline; not an implementation-promise claim and not a
protocol freeze. This records the condition of Ex7 before its source-grounded
redesign. Its table is intentionally preserved as baseline evidence; evaluate
the current runtime through `docs/promisegrid-implementation-claims.md` and
`docs/testing.md`. Source: DI-tohak.

## Governing interpretation

Ex7 is a PromiseGrid **app**. It may define makerspace-specific protocol
contracts, but it must not present an app-local host, account, key custody
choice, browser, relay, storage layout, or runtime process as a universal
PromiseGrid requirement. A shipped implementation claim must name the exact
frozen spec doc-CID and scope it actually implements. Source: PromiseGrid
Development Guide, App Devs and Kernel Devs sections; `TE-nibar`, `TE-zukug`,
and `TE-liviv` in the guide's cited upstream sources.

## Classification

| Surface | Baseline evidence | Classification | Audit result at time of audit |
| --- | --- | --- | --- |
| Equipment observation, safety disposition, off-site loan, and off-site return documents | `docs/protocols/makerspace-families/` and the makerspace pCID registry | Ex7 protocol-design artifacts | They describe plausible app-specific durable meanings, but the running application does not emit or consume their Grid records. They are not currently implemented protocol claims. |
| `makerspace-record-v1` | `docs/protocols/makerspace-record-v1.md` | Ex7 protocol-design artifact | It defines an Ex7 record profile. It must not be described as a final universal Grid envelope; upstream envelope/signature details remain provisional. |
| Participant continuity, revocation, and peer-card documents | uncommitted `docs/protocols/participant-identity/` | Draft design artifacts requiring reset | They were drafted as if they completed a universal identity boundary. The guide instead makes signing-key continuity a protocol concern while leaving embodiment and host shape explicit. These documents need a source-grounded redesign before any freeze or implementation claim. |
| `Record`, `ParseRecord`, and record tests | `service/record.go`, `service/record_test.go` | Unintegrated local implementation experiment | The code can sign/parse a custom Ex7 record profile, but no running action calls it. The test proves only that helper's local behavior, not a makerspace protocol implementation. |
| `AppendRecords` / `records.frames` | `service/store.go` | Unintegrated local implementation experiment | The function is unused. The running application still appends/replays JSONL events. It is not evidence of framed-record durability or Grid replay. |
| Browser UI and JSON HTTP endpoints | `web/`, `service/server.go` | Host-local UI/bridge behavior | Requests submit caller-selected member IDs and JSON fields to one local process. They are not signed Grid ingress, account authority, or a peer protocol. |
| Historical `events.jsonl` and in-memory projection | prior local-demo baseline | Superseded local-demo evidence | The JSONL append/replay path was removed by the DI-tohak / DI-piruf record conversion. Historical thought experiments retain its former scope as evidence; the running store now uses exact `records.frames` and signed-record replay. |
| Members, qualifications, authorities, areas, and tools | `service/app.go`, `service/types.go` | Local bootstrap/presentation/policy | These built-in fixtures are neither portable membership nor durable role continuity. They are appropriate only when labelled as local demo policy. |
| Relay/feed | repository search and runtime entrypoint | Deferred / absent | No relay or peer carriage is implemented. No relay claim may be made. |
| Existing makerspace accounts | repository search and UI | Absent | The current demo has selected fixture IDs, not accounts. A future account UI is host policy, never author evidence. |
| Current verification | `service/*_test.go`, existing guides | Local-demo evidence | Tests cover workflow rules, JSON parsing, replay, and isolated record helper behavior. They do not prove inter-peer interoperability, spec conformance, signing-key continuity, relay carriage, or account/author separation. |

## Immediate redesign constraints

1. Do not claim the local demo implements any makerspace pCID until its running
   paths create, validate, retain, and project the exact bytes required by that
   named Ex7 spec.
2. Do not claim a universal Grid identity, envelope, relay, bootstrap, runtime,
   or storage standard. Where upstream remains provisional, Ex7 must state its
   app-specific assumptions and its host-dependent choices.
3. Treat friendly names and fixture IDs as presentation/local policy, not
   signing-key continuity or authority.
4. Keep protocol-owned durable evidence distinct from local projections, caches,
   UI drafts, and bridge/API mechanics.
5. Before a spec is called frozen, provide its assumptions, known issues, and
   open questions; publish its doc-CID through the spec-side freeze history;
   then add only the implementation-side conformance claim the code can prove.
6. Preserve unknown pCID bytes without semantic acceptance only when Ex7
   actually implements that behavior and tests it.

## Redesign order

The next design artifact must map the makerspace workflows to explicit Ex7
protocol contracts and state which concerns stay host-local. Only after that
mapping is locked may Ex7 redesign its code, tests, and implementation promise
claims.
