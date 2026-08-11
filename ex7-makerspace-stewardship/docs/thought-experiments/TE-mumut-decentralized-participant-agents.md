# Decentralized Participant Agents and Non-Authoritative Carriage

TE ID: TE-mumut

## Status

needs DF

## Decision under test

How Ex7 becomes a decentralized Grid product: who holds author keys, where records are signed and projected, and how records move without a central server becoming author or authority.

This supersedes the single-runtime key ownership premise of TE-zadam / DI-simus.

## Assumptions

- Alice, Carol, and every participant are independent agents, not labels owned by one server process.
- Each agent owns its private signing key; no relay or shared HTTP service may sign a participant record or decide its makerspace meaning.
- Exact signed pCID-selected records are common evidence. Every agent keeps its own durable history and computes local projection/trust policy.
- A relay/feed may carry exact bytes but has no authority to alter, accept, revoke, or reinterpret them.

## Alternatives

### A. Independent participant agents with participant-held keys

Each participant runs an Ex7 agent and holds their private key. Their browser signs locally under explicit user control, sends exact signed bytes to that participant's local agent for storage, and exchanges bytes with other agents through a non-authoritative relay/feed. Each agent verifies signatures and applies its own recognition policy.

### B. Shared server with per-person keys

One shared service stores all keys and signs after receiving browser requests.

### C. Shared server as record proxy

Browsers sign but submit through one shared server that is the required durable-history and current-state source.

## Scenario analysis

### Normal operation

Under A, Alice signs her observation and loan promise; Carol signs her inspection and safety disposition. Their agents retain exact bytes and may exchange them through a relay. No component tells every other agent what state must be. B makes the server the key holder. C avoids server signing but centralizes availability and history.

### Failure and partition

If a relay fails, A's agents retain evidence and later exchange it. If Alice is offline, Carol may retain Carol's own record but cannot manufacture Alice's promise. B and C make the shared service a continuity dependency.

### Trust, authority, and revocation

A makes signatures portable while assessment stays local: an agent may recognize Carol's active key for a clearance, reject it, or retain it without projection. Key continuity and revocation are agent/protocol concerns, not central-account properties. B gives the service impersonation power. C gives it de facto gatekeeping power.

### Long horizon and scale

A supports participants, forks, and alternative relays without rewriting records. The relay scales as carriage infrastructure because it carries bytes rather than identity or state. B and C accumulate central governance power inconsistent with the product goal.

## Conclusions

B and C are rejected. A is the only architecture consistent with no central server: independently signing participant agents, locally retained projections, and non-authoritative byte carriage.

## Output to decision framing

**A is recommended:** participant-held keys; one local Ex7 agent/runtime per participant; direct local signing; per-agent durable store/projection; relay/feed only for exact record carriage. A separate follow-on decision must select the browser-to-local-agent embodiment and key-continuity/revocation protocol before code is written.

## Decision status

needs DF
