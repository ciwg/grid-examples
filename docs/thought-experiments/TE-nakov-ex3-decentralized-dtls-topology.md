# Ex3.5 decentralized DTLS Grid Editor topology

**TE ID:** TE-nakov  
**Status:** needs DF

## Decision under test

How should `ex3.5-grid-editor-dtls` provide a working collaborative editor
with no central relay or service while retaining Ex3's browser and Neovim
embodiments and using direct DTLS/UDP peer links?

## Scope and assumptions

- Each participant runs a local Grid Editor node that retains its own identity,
  signed objects, CRDT state, and peer configuration.
- Alice and Bob collaborate through direct DTLS/UDP links. A browser reaches
  its node through a local bridge because browser JavaScript has no raw-DTLS
  socket API.
- Mallory may drop, reorder, replay, inject, or flood datagrams, advertise
  stale peer information, or operate an incompatible node.
- The existing pCID-selected document and awareness payload meanings remain
  intact. DTLS is a carriage/binding choice and does not create new PromiseGrid
  top-level actions.
- The bundled PromiseGrid Development Guide 0.1.0 treats rendezvous and
  ingress as host- and spec-dependent; it does not prescribe a universal peer
  discovery method. This is current guide guidance, not a frozen rule.

## Alternatives

### A. Explicit direct peer addresses

Each node is started with one or more peer DTLS endpoints, such as
`--peer-dtls=alice.example:7026`. Peers exchange signed protocol objects
directly. Address exchange happens out of band (a command-line flag, copied
invite, or local configuration), with no hosted discovery service.

### B. Local-network discovery

Nodes discover each other on a LAN with mDNS or equivalent local broadcast,
then establish direct DTLS sessions.

### C. Decentralized overlay discovery

Nodes participate in a DHT, gossip mesh, or referral network to discover
addresses, then establish direct DTLS sessions.

### D. Shared relay/service

A central relay accepts browser and sidecar traffic and distributes it to
participants.

## Scenario analysis

### Normal remote collaboration

Alice gives Bob her explicit endpoint and DTLS certificate fingerprint. A
creates a direct session and both nodes exchange signed document changes; it
has no central service dependency. B works only when both are on a suitable
LAN. C can find Alice without a prior address exchange, but is a separate
distributed-systems project. D is operationally simple but violates the
no-central-service requirement.

### Packet loss, offline peers, and incomplete writes

Alice edits while Bob is offline. In A, Alice persists the signed change
locally and retries it when the configured peer reconnects; Bob deduplicates
by signed-object identity. B and C need the same durable outbound queue, plus
discovery-state repair. D centralizes the queue and therefore changes the
ownership model. Awareness remains newest-state, not durable queue work.

### Concurrent writers and mixed versions

Alice and Bob edit concurrently while Carol runs an older node. A uses the
existing CRDT convergence rules and versioned protocol/payload checks on each
direct link; an incompatible Carol is refused without blocking Alice and Bob.
B adds local multicast version/discovery compatibility. C adds overlay
versioning, referral validation, and routing compatibility. D hides some
differences behind the relay, but at the cost of a central boundary.

### Trust boundary and adversarial peers

Mallory advertises a forged endpoint. A makes the trusted peer address and
certificate fingerprint explicit, then validates DTLS plus the existing signed
envelope on every received object. B adds unauthenticated LAN advertisements
that must be treated as hints only. C adds routing/referral claims, Sybil
pressure, cache expiry, and local trust policy. D makes the relay an
availability and traffic-analysis focal point.

### Scale, NAT, and long-horizon evolution

A is straightforward and demonstrably decentralized, but requires reachable
addresses or user-managed port forwarding/VPN and is practical for a small
known peer set. B improves same-LAN usability but does not solve Internet NAT.
C can improve discovery but still cannot bypass NAT without additional
connectivity machinery and carries the greatest CPU, bandwidth, storage, and
operational complexity. D scales deployment coordination but is rejected by
the stated topology requirement.

## Comparison

| Alternative | Makes easier | Makes harder / obligation | Result |
| --- | --- | --- | --- |
| A | Direct, inspectable, no-service collaboration | Manual endpoint exchange and reachable UDP | Survives; recommended first implementation |
| B | Same-LAN discovery | LAN-only discovery, advertisement trust | Deferred extension |
| C | Wider discovery | Overlay, NAT, Sybil, and routing design | Deferred extension |
| D | Centralized availability/queueing | Violates no-central-service requirement | Rejected |

## Conclusions

Alternative A is the smallest architecture that is both genuinely decentralized
and runnable now: every node owns persistence and connects directly to known
peers by DTLS. It makes rendezvous an explicit, user-controlled input instead
of quietly replacing decentralization with a discovery service.

B and C remain legitimate future extensions, but neither is needed to prove a
working direct-peer editor. The implementation must document its reachability
assumption clearly: two hosts need mutually reachable UDP endpoints (or an
operator-provided private network), and no central component is involved.

## Output to DF

- **Rejected:** D, because it contradicts the no-central-service requirement.
- **Surviving:** A now; B and C as later extensions.
- **Recommended decision:** A — direct configured peer endpoints and durable
  per-node outbound recovery queues.
- **New decisions exposed:** replacement/supersedence of the prior relay-based
  DI, the peer-invite format, DTLS fingerprint trust record, direct-peer
  configuration paths, bridge ownership, and multi-node proof harness.

## Implications for existing decisions

The relay-based architecture captured in DI-ditaz and its dependent relay
transport decisions cannot govern this topology if A is selected. They require
explicit supersedence in the relevant TODO log before implementation. The
existing `TE-sohoz` remains historical analysis of a relay-oriented browser
constraint; it does not settle this direct-peer topology.

## Sources

- Bundled PromiseGrid Development Guide 0.1.0 — App Devs: bootstrap and
  rendezvous are host- and spec-dependent; no universal discovery method is
  currently frozen.
- [RFC 9147: DTLS 1.3](https://datatracker.ietf.org/doc/html/rfc9147) — DTLS
  protects UDP datagrams without supplying reliable, ordered application data.
