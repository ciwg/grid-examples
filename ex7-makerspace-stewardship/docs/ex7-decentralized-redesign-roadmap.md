# Ex7 decentralized redesign roadmap

## Status

Planned implementation roadmap. This document sequences the locked decisions
in DI-lazim, DI-janup, and DI-sinov; it does not create a central service,
author keyring, or shared authoritative state.

## Target boundary

Ex7 is a set of independent participant agents. Each agent owns its own
signing identity, durable exact-record history, blob store, and local
projection. A participant's browser or makerspace account helps them find and
use their agent; neither is cryptographic author evidence. A relay carries
exact bytes between agents but does not sign, validate makerspace meaning,
choose trust, or compute a global current state. Source: DI-lazim, DI-janup,
DI-sinov.

The current JSONL-backed local demo is pre-recreation code. It must not be
presented as this target architecture while the slices below remain open.

## Slice 1 — Participant identity and agent bootstrap

**Goal.** Give every participant agent a durable root signing identity without
making a browser profile, makerspace account, relay, or shared HTTP process the
identity owner.

**Deliverables.**

- Define and freeze the participant root-identity, key-continuity, and
  revocation family specifications before code relies on them.
- Implement local agent initialization and durable root-key storage under an
  agent-owned runtime root.
- Implement an explicit public bootstrap/peer-card representation that lets a
  UI or another agent discover how to contact an agent without possessing its
  private key.
- Treat account login only as local UI bootstrap and recognition policy.

**Done when.** Alice and Carol can each initialize separate agents; each agent
can restart with its own identity; public identity material verifies without a
website; and no process can produce an Alice signature from Alice's account
session alone.

**Guardrail.** Device authorization, recovery, and browser-held keys are not
silently substituted for participant identity. If needed, they require their
own frozen family protocols and evidence. Source: DI-janup.

## Slice 2 — Signed exact-record ingress

**Goal.** Make an agent accept, preserve, and assess exact signed PromiseGrid
records from its participant or peers.

**Deliverables.**

- Finalize the shared canonical Grid envelope against the frozen makerspace
  family pCIDs and the participant-key policy from Slice 1.
- Implement deterministic encode, sign, parse, and verification operations.
- Add an ingress boundary that accepts exact record bytes, preserves unknown
  pCIDs and locally untrusted but well-framed records, and only projects
  evidence recognized by that agent's local policy.
- Retire the implication that browser JSON or an HTTP handler is durable
  author evidence.

**Done when.** Alice's agent signs an off-site-loan record; Carol's agent can
verify and retain those exact bytes; tampering is rejected; and an unknown pCID
remains retrievable byte-for-byte without gaining known makerspace semantics.

**Guardrail.** Relay-carriage authentication, if later needed, is separate
from the semantic author signature. Source: DI-bosur, DI-lazim, DI-sinov.

## Slice 3 — Per-agent local records, blobs, and projection

**Goal.** Replace legacy event persistence with exact per-agent history and a
derived local view.

**Deliverables.**

- Implement the MSR1 framed-record and CID-addressed blob store per agent,
  with atomic write order, fsync, bounded frame/blob sizes, and fail-closed
  replay.
- Replay exact records into a local projection only after durable append and
  local evidence assessment.
- Convert the current App actions from legacy events into draft/request and
  projection behavior; no App or server component signs for an absent
  participant.
- Add corruption, partial-frame, duplicate, unknown-family, and replay tests.

**Done when.** A restarted agent reproduces its projection solely from its
framed exact bytes and blobs; incomplete/corrupt frames fail closed; and the
legacy `events.jsonl` path is no longer a supported durable format. Source:
DI-nohos, DI-rifib, DI-sazir, DI-sinov.

**Guardrail.** Storage is agent-local. A runtime root is never a shared
makerspace authority or a place where an account service keeps participant
private keys.

## Slice 4 — Non-authoritative relay carriage

**Goal.** Allow agents to exchange exact evidence without creating a central
record or state authority.

**Deliverables.**

- Define and freeze a relay/feed carriage specification, including batch,
  cursor, duplicate, retention, and transport-integrity semantics.
- Implement relay append/pull of opaque exact record bytes (and any explicitly
  specified blob carriage), with no family-specific projection.
- Implement agent push/pull and idempotent local retention.
- Add two-agent tests across a relay outage, duplication, malformed input, and
  an unknown pCID.

**Done when.** Alice can publish a record after local durable acceptance;
Carol can later retrieve the same bytes through a relay and independently
choose whether to project them; relay outage delays carriage but never changes
authorship or local history. Source: DI-lazim, DI-sinov.

**Guardrail.** The relay cannot sign for Alice, declare a key trusted or
revoked, reject an unknown family because it is unknown, or publish a global
makerspace state.

## Slice 5 — Browser and account embodiment

**Goal.** Restore a usable makerspace experience without conflating a login
with a signature.

**Deliverables.**

- Implement browser-to-agent discovery, pairing, and request flow.
- Let an existing makerspace account supply only UI session, discovery, and
  local policy inputs.
- Render durable signed records and locally derived state distinctly from
  unsigned drafts or unavailable-agent requests.
- Provide a kiosk path that requests a reachable participant agent signature
  and otherwise stores or exports an unsigned draft.

**Done when.** Alice can use a makerspace kiosk account to request a loan from
Alice's own reachable agent; when that agent is unavailable, the UI visibly
cannot claim a signed Alice promise. Source: DI-janup.

**Guardrail.** No browser private key is the participant identity by default,
and no account token is converted into an author signature.

## Slice 6 — End-to-end proof and aligned documentation

**Goal.** Prove the real embodiment and make all public claims match it.

**Deliverables.**

- Build deterministic unit coverage for codec, key policy, ingress, framed
  replay, blobs, projection, relay, and account/UI separation.
- Build an opt-in multi-agent end-to-end proof: Alice's agent signs a loan,
  a relay carries exact bytes, Carol's agent verifies and projects according
  to Carol's local policy, and the browser distinguishes a signed result from
  a draft.
- Update README, architecture, testing guide, implementation claims,
  CHANGELOG, protocol registry, and TODO/deferral inventory together.
- Run `go test ./...`, `go vet ./...`, `errcheck ./...`, protocol-CID/label
  consistency scans, and the opt-in E2E proof.

**Done when.** The commands and stored evidence demonstrate the full path with
no central signer or authoritative server claim, and every deferred path is
explicit.

## Ordering and stop conditions

The slices are deliberately ordered. Slice 2 cannot be trusted before Slice 1
defines whose signature it verifies; Slice 3 cannot replace event persistence
before Slice 2 supplies exact verified bytes; Slice 4 cannot be added before
agents own durable records; and a browser must not be treated as an author
before the participant-agent boundary exists.

If a slice exposes a new protocol family, key-continuity rule, runtime path,
or trust boundary, stop that slice for its required thought experiment,
decision framing, frozen specification, and pCID. Do not bridge the gap with a
server-held key, account-derived signature, or undocumented compatibility
path.
