# Ex4 implementation scope

## 2026-08-10 — Local workflow and future issue-promise boundary

This is a scope declaration for the Ex4 Bug Tracker example, not a formal
PromiseGrid implementation-promise claim against a frozen upstream
specification doc-CID. Source: `DI-nibuh`.

### Current local workflow

Ex4 currently implements one local HTTP-server workflow. Its browser and CLI
are local adapters of that service. The service allocates issue IDs, applies
built-in identity and role checks, validates the fixed issue-status flow,
copies attachments, and projects the current issue view.

`events.jsonl` is durable local application history for that server workflow.
It records the server's accepted issue events and supports local projection and
inspection. It is not a signed PromiseGrid promise ledger, independently
verifiable shared evidence, or proof of another actor's intent. Source:
`DI-nibuh`.

### Planned bounded issue-promise layer

Ex4's selected alignment direction is a later bounded layer of signed,
pCID-selected issue-promise artifacts. That layer is not implemented yet. Its
separate `kakon.2` thought experiment and Decision Framing must select the
profile meaning, signing identity, accepted/rejected artifact rules, adapter
projection, and any remote exchange before behavior changes. Source:
`DI-gisor`; `DI-nibuh`.

### Explicit non-claims

Ex4 does not currently claim:

- conformance to a frozen upstream PromiseGrid spec doc-CID;
- signed, pCID-selected issue promises or independently verifiable shared
  issue evidence;
- interoperability with another tracker based only on matching issue fields,
  HTTP endpoints, event names, or role labels;
- that built-in identities, role checks, or fixed status transitions establish
  general agent identity, delegation, revocation, role continuity, or a
  portable authorization system; or
- that a local server event, projection, rejection, or attachment proves
  another participant's intent.

When Ex4 implements a selected frozen upstream specification, a later entry
will use the guide's formal implementation-promise fields (`claim`, `spec`,
`scope`, `breaking-change`, and `notes`) and name that exact frozen spec
doc-CID.
