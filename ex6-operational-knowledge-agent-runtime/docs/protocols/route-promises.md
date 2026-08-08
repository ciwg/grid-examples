# OKAR Route Promises Protocol

## Scope

This document defines the implementation-local, pCID-selected CBOR records that
let an OKAR runtime retain explicit local route evidence. It implements
DI-kojab, DI-komaz, and DI-butam. It is not a frozen global PromiseGrid API.

An app agent publishes a `ReceivePromise` for a pCID it promises to receive. A
routing-role agent separately publishes a `DeliveryPromise` for the same pCID
and named recipient. An `AgentBinding` maps an opaque local app label to the
installed package metadata that this runtime can use as an implementation
adapter. It is not a promise and is not an identity claim.

`agent_id` is a non-empty local app label. It is neither a package ID nor a
durable PromiseGrid identity; this unsigned, non-networked profile makes no
claim about signing-key continuity. Package claims are bootstrap hints only and
must never create a binding or promise automatically. Source: DI-komaz;
DI-butam.

## Envelope

Each record is exact canonical CBOR bytes of:

```cddl
grid(
  [ selector: 42(pCID),
    kind: uint,
    agent_id: tstr,
    protocol_pcid: tstr,
    recipient_agent_id: tstr / null,
    package_id: tstr / null,
    enabled: bool,
    parents: [* bstr],
  ]
)
```

The selector is the binary CID for this specification. The full canonical CBOR
sequence is retained under its own CID in local CAS.

| `kind` | Record | Required fields |
| --- | --- | --- |
| 0 | `AgentBinding` | `agent_id`, `package_id`, `enabled`; `protocol_pcid` empty; `recipient_agent_id` null |
| 1 | `ReceivePromise` | `agent_id`, `protocol_pcid`, `enabled`; `recipient_agent_id` and `package_id` null |
| 2 | `DeliveryPromise` | `agent_id`, `protocol_pcid`, `recipient_agent_id`, `enabled`; `package_id` null |

`parents` is empty for the first record for one logical key and contains exactly
one accepted prior record CID for later updates. The key is `agent_id` for a
binding, `agent_id + protocol_pcid` for a receive promise, and `agent_id +
recipient_agent_id + protocol_pcid` for a delivery promise. Competing heads
make that key unusable rather than allowing the projection to choose one.

## Local replay and planning

The runtime scans retained local CAS objects, accepts only valid records with a
complete predecessor chain, and rebuilds a disposable projection. It ignores
unknown selectors, malformed records, missing parents, and competing heads.

A package-derived route becomes executable only when all of these local records
are enabled and valid: a binding for its package, a receive promise by the bound
agent for the input pCID, and a delivery promise naming that recipient and pCID.
The first slice remains deterministic, local, non-executing, unsigned, and
non-networked. It does not treat package claims as promises and does not send
or accept any record across a transport. Source: DI-kojab; DI-butam.

The local CLI records this evidence explicitly with `moks route bind <agent-id>
<package-id> <true|false>`, `moks route promise receive <agent-id> <pcid>
<true|false>`, and `moks route promise deliver <router-id> <recipient-id>
<pcid> <true|false>`. These commands do not communicate with another agent or
execute a route. Source: DI-zolil.

## Evolution

Changing this record meaning or encoding requires a new pCID and a separate
protocol document. A future signed or networked version must not reinterpret
these local labels as durable identities. Source: DI-komaz.
