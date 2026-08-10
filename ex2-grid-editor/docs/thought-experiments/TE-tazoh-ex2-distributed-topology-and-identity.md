# Ex2 distributed topology and identity

TE ID: TE-tazoh

## Status

decided

## Decision status

Locked by DI-nilas: Ex2 uses decentralized multi-relay peer collaboration;
relay keys identify app nodes and local observers; browser and Neovim are local
embodiments; and one-host execution distinguishes one node from a multi-relay
simulation. DI-lorud is superseded.

## Decision under test

What is Ex2's canonical deployment and identity model: how many relay nodes
exist in a normal collaboration, what each relay's signing key represents, how
browser and Neovim embodiments relate to relays, and what a single-host run is
allowed to claim?

## Assumptions and trust model

- Ex2 currently has a Go relay that owns a signing identity, CAS, append-only
  message log, verification, peer ingestion, and four local-draft pCID
  handlers.
- Browser and Neovim currently use local HTTP/sidecar plumbing and own local
  Automerge replicas; they do not presently own separate durable signing keys.
- The code supports a second relay configured with `--peer`, while the README's
  quick start demonstrates one local relay serving browser and Neovim.
- A relay can be copied, restarted, moved, compromised, partitioned, or run
  beside another relay on one host. A browser or Neovim process can be closed,
  restarted, or moved to another machine.
- The current PromiseGrid guide distinguishes signing-key continuity from
  presentation identity, keeps protocol contract separate from host process
  shape, and treats frozen specs as future contract references. Ex2's profiles
  remain local drafts.
- Alice, Bob, and Carol are cooperative collaborators. Mallory may compromise
  one relay, replay retained bytes, impersonate a display name, or exploit an
  unclear topology claim.

## Alternatives

### A. Canonical multi-relay collaboration; single-host runs simulate nodes

Each collaborating machine runs its own relay with its own signing key and
data root. Browser and Neovim are embodiments connected to their local relay.
Relays exchange signed profile traffic as peers. A single host may run several
separate relay processes/data roots to simulate those logical nodes for tests
or demos.

### B. Canonical one-relay collaboration; peer relays are optional extension

One relay is the normal shared collaboration service for all browser and
Neovim embodiments. Its signing key represents the app service. Multi-relay
peer exchange is an optional advanced topology outside the main identity model.

### C. Canonical centralized relay service

One designated relay is the authoritative collaboration host and assigns the
meaning of participant identities and document state to all embodiments.

## Scenario analysis

### Normal two-person collaboration

Alice edits in a browser and Bob edits in Neovim. Under A, Alice's and Bob's
local relays sign and retain their own traffic, then peer exchange carries
protocol-defined artifacts. Each embodiment remains local UI/replica plumbing.
Under B, both embodiments can use one relay efficiently, but Bob's durable app
identity is the same service key as Alice's unless a separate identity layer is
added. Under C, the central relay becomes the effective authority for identity
and document continuity.

### Browser and Neovim on one laptop

Alice runs browser and Neovim on one machine. Under A, they may share Alice's
one local relay because they are two embodiments of Alice's app node; this is a
real one-node collaboration case, not a simulation of two relay identities.
To simulate Alice and Bob, the host runs two distinct relays with distinct data
roots and signing keys. Under B, the same one-relay configuration is the
canonical collaboration topology. Under C, it resembles the central service
model and hides the distributed trust boundary.

### Network partition and restart

Alice's relay is offline while Bob continues editing. A lets both retain signed
traffic and later exchange it; CRDT replicas converge under their defined
protocol behavior, while each relay's evidence remains locally owned. B makes
offline continuation depend on whether the shared relay is still available.
C turns a central-service outage into an authority and availability failure.

### Identity, presentation, and compromised node

Mallory copies Bob's display name or color. A and B both treat it as
presentation data, not a signing identity. Under A, Bob's relay signing key
names his current app node; a compromised relay affects that node's local
evidence and emissions but does not grant authority over Alice's relay. Under
B, the shared service key can make Alice and Bob indistinguishable at the
PromiseGrid-facing layer. C makes the central key a stronger single point of
authority and compromise.

### Evidence and observer records

Under A, an observation should identify the relay key and local data root that
observed it; the physical host is operational context, not a protocol
identity. Under B, the one service key identifies the observer, but it cannot
distinguish collaborating human/app nodes without new client-side signing
identity. Under C, observer records look authoritative even though they remain
only one host's evidence.

### Long-horizon evolution and scale

A preserves independent relay histories, permits new embodiments, and makes
single-host multi-relay simulation explicit. It creates operational work for
peer configuration, synchronization, duplicate suppression, key rotation, and
partition recovery. B is simpler for a small shared service, but any later
multi-machine claim needs a new identity and evidence model. C is easiest to
operate initially but conflicts most strongly with the guide's avoidance of a
permanent authority and makes migration/forking harder.

## Conclusions

C is rejected: it conflicts with Ex2's non-canonical relay/CRDT design and
would turn a convenience service into a durable central authority. B remains
viable for a deliberately hosted shared-editor example, but it does not
support a strong multi-machine PromiseGrid collaboration claim without a new
client/participant identity layer.

A is recommended if Ex2's purpose is a distributed PromiseGrid example:
multiple relay nodes are the canonical peer topology; each relay has a signing
identity and local evidence store; browser/Neovim are local embodiments; and a
single-host setup can either be a genuine one-node session or a multi-relay
simulation depending on how many relay processes, data roots, and keys it
runs. The distinction must be documented precisely. This does **not** make a
relay key a human identity, nor does it decide key rotation, delegation, or
cross-relay trust policy beyond the current local-draft scope.

## Decisions still requiring DF

1. **Canonical topology:** adopt canonical multi-relay peer collaboration with
   single-host multi-relay simulation (recommended), canonical one-relay
   collaboration with optional peers, or centralized relay authority?
2. **Relay identity:** treat each relay signing key as the current app-node
   identity and local evidence observer (recommended), or defer app-node
   identity and treat the key as only a transport implementation detail?
3. **Embodiment role:** treat browser and Neovim as local embodiments of the
   relay/app node until a separately designed client-signing layer exists
   (recommended), or treat their participant fields as independent peer
   identities now?
4. **Single-host wording:** distinguish a genuine one-node session from a
   multi-relay simulation by the number of relay processes, data roots, and
   signing keys (recommended), or describe every single-host run simply as a
   distributed simulation?

## Implications for open work

- After DF, supersede or confirm DI-lorud with the correctly scoped topology
  decision.
- `sojot.2` may then finalize whether an observation record needs an observer
  relay key and how documentation describes its local ownership.
- If A is selected, later guide/testing work must show at least one
  two-relay exchange scenario without treating browser or Neovim display data
  as signing identity.
