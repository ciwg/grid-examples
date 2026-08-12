# Participant Root Identity and Agent Bootstrap

TE ID: TE-mamop

## Status

superseded by TE-zizum / DI-girup

## Decision under test

For TODO `giman`, Slice 1: select the first Ex7 participant identity model,
the authority that signs ordinary makerspace promises, and the public bootstrap
artifact that lets an account UI or another agent find an agent without making
either one an identity authority.

## Assumptions and trust model

- Alice and Carol are independent Ex7 participant agents. Each holds its own
  private material; no makerspace website, browser session, relay, or shared
  service may possess or use it.
- Existing makerspace accounts can choose a UI experience and supply local
  policy inputs, but an account credential is not author evidence.
- Durable continuity and revocation evidence must survive a website or relay
  failure and be assessable by each receiving agent.
- A peer card can carry public contact hints. It is advisory discovery data,
  not proof of membership, authorization, current key status, or makerspace
  state.
- Device authorization, delegated hardware signers, social recovery, and
  multi-party recovery remain out of scope unless this decision necessarily
  creates them.

## Alternatives

### A. Participant root key signs ordinary promises directly

Agent initialization creates one Ed25519 participant root key. That key signs
ordinary makerspace promises as well as the continuity and revocation promises
for its identity. A public peer card names the root public key and optional
agent contact hints. Rotation is a promise signed by the currently active root
key that names the replacement root key; revocation is a promise from an active
root key that names a compromised or retired key and a reason. A later agent or
device delegation protocol may authorize narrower keys, but it is not required
for this first embodiment.

### B. Root key authorizes a separate agent operational key

Initialization creates a root key and a distinct agent operational key. The
root signs an authorization promise for the operational key; the operational
key signs ordinary makerspace promises. Continuity and revocation must define
both root replacement and operational-key withdrawal. The peer card publishes
the root and current operational key.

### C. Account or relay-managed participant key

A makerspace account service or relay holds, selects, or recovers the key that
signs ordinary promises, while the agent acts as a client or cache.

## Scenario analysis

### Normal operation

Under A, Alice initializes one agent, publishes a peer card containing only
her public root key and a contact hint, and signs a loan promise locally. Carol
can verify that same public key directly from the record and assess its trust
under Carol's local policy. The peer card helps Carol find Alice's agent but
does not turn the contact hint into authority.

B reduces routine use of a root key but immediately adds another protocol
family, validity rules, nested revocation, and a choice of whether an
operational key represents an agent, device, browser, or hardware token. That
is the delegated-embodiment design intentionally deferred by DI-janup.

C makes familiar account login convenient but gives the account or relay the
ability to impersonate Alice or block continuity, contradicting DI-lazim.

### Key compromise, rotation, and revocation

With A, Alice still holding the active root key signs a continuity promise to a
new root key and may sign a revocation promise for the old key. Receiving
agents retain both exact records and locally decide whether the chain satisfies
their policy. If Alice loses every active key, Ex7 does not invent recovery;
that requires a later witness/recovery protocol. A compromised active root can
make harmful promises until other agents learn and apply revocation evidence,
which is visible rather than hidden by a registry.

B can revoke only the operational key while preserving the root, but must
define which authorization wins during partitions, whether expired authority
invalidates a record made earlier, and how an agent migration differs from a
device migration. It has better compartmentalization only after those
additional obligations are designed and tested.

C may make rotation feel instant, but that is a centralized assertion that
other participants cannot independently verify without trusting the service.

### Offline, partitioned, and mixed-version operation

Under A, Alice's agent continues to sign and retain promises offline. A relay
or account outage does not affect authorship. Carol can cache the peer card and
later receive continuity/revocation promises as exact bytes. An older agent
that does not recognize the continuity pCID retains it as unknown bytes rather
than silently treating a new key as trusted.

B multiplies mixed-version requirements because every verifier must understand
authorization as well as continuity. C makes an outage a key-discovery and
possibly signing dependency.

### Long-horizon evolution and scale

A establishes a stable participant public-key anchor and minimal peer-card
format. Future delegated signers can be introduced as additional pCID-defined
promises without changing the meaning of records already signed by the root.
B pre-commits Ex7 to delegation semantics before there is a real delegated
embodiment to test. C centralizes the relationship that PromiseGrid records are
meant to preserve independently.

## Protocol and artifact consequences

Alternative A needs three immutable specifications before code depends on
them:

1. `participant-key-continuity-v1`: a `promise` whose payload identifies the
   prior and replacement Ed25519 public keys and effective intent; the active
   prior key supplies semantic author evidence.
2. `participant-key-revocation-v1`: a `promise` whose payload identifies a
   key, a bounded revocation reason, and the authoring active root key.
3. `participant-peer-card-v1`: a public bootstrap/discovery encoding containing
   a root public key, stable public fingerprint, protocol version, and optional
   advisory contact hints. It is not a durable makerspace promise and carries
   no account credential, private key, trust decision, or global identifier.

Each specification receives a frozen CIDv1 pCID. Continuity and revocation
records use the existing Grid envelope with top-level semantic action
`promise`; the peer card is a separate public protocol artifact, never an
author-evidence substitute.

## Conclusions

C is rejected because it is a central keyholder. B is deferred, not rejected
forever: it is the correct direction only when Ex7 has a concrete delegated
agent/device embodiment and can define its authority lifecycle. A is the
smallest complete participant-agent identity boundary now. It keeps authorship
with Alice's agent, supports verifiable continuity, and leaves future
delegation additive.

## Output to decision framing

Surviving choice: **A — root key signs ordinary promises directly.**

Recommended locks, if accepted:

- Freeze the three specifications at
  `docs/protocols/participant-identity/` and publish their CIDv1 pCIDs in a
  central registry alongside the makerspace registry.
- Store the private root key only inside the local participant-agent runtime
  root; determine the exact runtime path, file format, permissions, and
  initialization command in the implementation decision framing.
- Make the peer card public and advisory; account login may use it for local
  discovery but may not rewrite it, sign with it, or treat it as membership.

## Decision status

superseded by TE-zizum / DI-girup. DI-basun locked Alternative A, participant
root key signs ordinary promises directly, in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`; the later coordinated
participant-root, authorized-device, and threshold-recovery contract set
supersedes that root-only embodiment.
