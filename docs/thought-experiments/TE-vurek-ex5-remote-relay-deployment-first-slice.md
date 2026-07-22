# ex5 remote relay deployment first slice

TE ID: TE-vurek
## Status
decided

## Decision under test

What the first dedicated remote relay deployment slice for
`ex5-operational-knowledge-system` should be now that the repo already ships:

- eight frozen signed families
- origin-aware ongoing peer exchange into non-empty runtimes
- incremental relay-feed export/import over the local adapter
- separate CID-addressed blob transfer for ongoing evidence carriage
- CAS-backed authoritative replay/export for the frozen families

The question is no longer whether `ex5` can exchange PromiseGrid-shaped
history. The question is which remote deployment shape most honestly extends
that relay model beyond the current local-adapter exchange layer without
collapsing back into a central app server or reusing embodiment routes as the
remote transport contract.

## Assumptions

- `ex5` nodes already know how to export/import signed relay-feed batches and
  stage/fetch blobs by CID.
- Browser, CLI, and Neovim still use the current local HTTP adapter for their
  embodiment surface.
- Live drafting remains a separate collaboration lane and is not part of this
  remote relay deployment decision.
- Mallory may observe or replay traffic through the relay, but cannot forge
  signed family envelopes.
- Mixed-version nodes matter because the repo already ships bundle-based
  bootstrap exchange and local-adapter relay feed.

## Alternatives

### Alternative A: remote proxy of the current local relay adapter

Deploy a remote service that largely republishes the current local relay-feed
and blob routes for one runtime instance. Remote peers push and pull through
that deployed adapter, and the relay behaves much like a network-reachable
front door for one existing ex5 node.

### Alternative B: dedicated neutral relay service for feed plus blob carriage

Deploy a separate relay service whose job is only:

- store origin-aware incremental relay-feed batches
- track peer-visible feed progress or acknowledgements
- store and serve CID-addressed blob objects
- never become the authoritative business runtime for browser/CLI/Neovim

The relay is a transport/storage intermediary for signed artifacts, not the
main ex5 application.

### Alternative C: central hosted ex5 server as the relay

Deploy one central hosted ex5 instance and make remote peers act more like
clients of that hosted runtime. Durable history, transport, and possibly later
embodiment behavior converge into one hosted server role.

## Scenario analysis

### Scenario 1: normal multi-peer operation

Alice, Bob, and Carol each operate separate ex5 runtimes at different sites.

Alternative A is easy to explain because it looks like “host the thing that
already exists.” But the deployment target is still one runtime-shaped
adapter, so it blurs relay responsibility with local app responsibility. The
remote front door becomes coupled to one node’s adapter semantics.

Alternative B lets each site remain its own runtime while the relay only moves
signed records and blobs. That is closer to the current PromiseGrid direction:
the relay exists to carry artifacts, not to own the application state.

Alternative C simplifies the network story by giving everyone one hosted
server, but it weakens the peer/runtime separation. The “relay” becomes the
main app server, which is a different product shape.

### Scenario 2: replay, duplication, and noisy transport

Mallory or an unreliable remote relay redelivers the same feed windows or blob
references repeatedly.

Alternative A can survive because runtime import already dedupes by
`(origin_peer_id, origin_sequence)`, but replay handling is tied to a hosted
app-shaped surface rather than to a narrower relay role.

Alternative B matches the current origin-aware relay feed more directly. The
relay can be noisy, but each node’s import path remains clearly about durable
artifact dedupe and blob presence, not about replaying higher-level app
operations.

Alternative C can also dedupe, but because the hosted server is now the main
system of coordination, replay handling becomes bound to the hosted runtime’s
application lifecycle rather than a cleaner transport/storage layer.

### Scenario 3: evidence blobs and CAS carriage

Dave records evidence with a large attachment, and Ellen must catch up through
the remote path later.

Alternative A can reuse the current feed-plus-blob routes almost verbatim, but
it still makes the deployed runtime adapter the remote contract. Blob
availability becomes one more concern of the hosted app node.

Alternative B keeps the strongest separation. Feed records travel as signed
history windows; blobs travel as CID-addressed objects through the relay’s blob
store. That aligns with the current CAS direction and keeps the relay focused
on transport/storage.

Alternative C can certainly host blobs too, but now evidence carriage is
centralized inside the hosted ex5 server role rather than kept as a reusable
artifact relay layer.

### Scenario 4: mixed-version compatibility

Frank still runs the older bundle-bootstrap build while Grace runs the current
incremental relay-feed build.

Alternative A is good at short-term compatibility because the hosted adapter
can continue exposing the older bundle routes next to the newer relay routes.
But that compatibility story depends on the hosted app-shaped surface staying
in the middle.

Alternative B is still compatible, but it creates an explicit translation
obligation: either the relay or a relay-aware peer may need to materialize a
bootstrap bundle for older peers while newer peers use incremental feed
windows. That is more work, but it preserves the cleaner long-term split.

Alternative C is worst here because older peers would now be compared against a
centralized hosted-runtime model that changes not only transport placement but
also the product shape.

### Scenario 5: long-horizon evolution

The repo later wants stronger remote relay policy, store-and-forward
deployment, or partial sync for selected peers.

Alternative A is a workable bridge, but the hosted adapter shape risks
becoming dead weight. Later transport/storage evolution would still be coupled
to one app-facing surface.

Alternative B creates the strongest base for later relay evolution. Selective
feed replay, blob caching, acknowledgement policy, and remote durability are
all easier to reason about when the relay is already a dedicated feed/blob
intermediary rather than a hosted local adapter.

Alternative C commits too early to a different end state: central hosted ex5
as the main coordination model.

### Scenario 6: operator and embodiment boundaries

Heidi uses the browser; Ivan uses CLI; Judy uses Neovim.

Alternative A risks letting the remote relay deployment shape leak upward into
the embodiment story because the hosted adapter looks like a network-visible
version of the same surface the embodiments already use locally.

Alternative B keeps the boundary cleaner. The local adapter remains the
embodiment surface, while remote relay deployment is a separate transport lane
for signed durable history and blobs.

Alternative C again mixes concerns: the hosted server begins to look like both
the embodiment server and the durable relay layer.

### Scenario 7: implementation and testability

The repo needs a first slice that can be built and verified without reopening
all of ex5.

Alternative A is probably the fastest first deployment to demo.

Alternative B is more work because it needs a dedicated relay role and
translation decisions, but it stays cleaner architecturally and aligns better
with the existing origin-aware feed model.

Alternative C is the heaviest because it changes the deployment role, product
shape, and future embodiment expectations at once.

## Conclusions

Rejected:

- Alternative C. It changes the remote deployment question into a hosted-app
  centralization question too early.

Surviving:

- Alternative A: remote proxy of the current local relay adapter
- Alternative B: dedicated neutral relay service for feed plus blob carriage

Recommendation:

- Alternative B

Why:

- It is the clearest PromiseGrid-shaped remote deployment path.
- It preserves the separation between local embodiment adapter and remote relay
  transport.
- It composes cleanly with the already-shipped origin-aware incremental feed
  and CID-addressed blob model.
- It avoids turning the next step into “host the whole app remotely” before
  that product choice is actually desired.

## Implications for TODOs and pending DIs

- TODO `116` should lock either:
  - `A`: deploy the current local relay adapter remotely as the first slice, or
  - `B`: introduce a dedicated feed/blob relay service as the first remote
    deployment slice.
- If `B` is chosen, the next DF will likely need one more scope lock about
  whether the first relay service stores feed windows ephemerally or durably.
- Future non-HTTP embodiment work under TODO `117` should remain separate
  unless a later TE proves that remote relay deployment and embodiment
  contracts should merge.

## Decision status

Locked and implemented:

- `116B`
- durable relay store
- separate `operational-relay` binary
- relay-only `/relay/v1` route surface
- per-origin append-only relay history with cursor-map pulls
- staged blobs required before evidence-bearing publish
