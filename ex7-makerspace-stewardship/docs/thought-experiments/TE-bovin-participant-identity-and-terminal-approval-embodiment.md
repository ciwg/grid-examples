# Participant Identity and Terminal Approval Embodiment

TE ID: TE-bovin

## Status

decided

## Decision under test

How an Ex7 participant agent holds its own signing capability and how a shared
terminal obtains a signed record without holding a participant private key or
treating an account session as author evidence. This TE implements the
participant-agent and terminal-request terms locked by TE-zizum / DI-girup.

## Assumptions and affected scope

- The six frozen protocol families already define durable root, device,
  revocation, recovery, peer-card, and carriage records.
- A personal participant agent may sign only with a root-authorized key.
- A terminal draft is bounded unsigned local data under
  `<agent-root>/requests/`; it is never an evidence record and is deleted after
  expiry or successful signing.
- A makerspace account may open a terminal UI and provide local facility
  context, but cannot approve, sign, or modify participant history.
- The affected implementation paths are `<agent-root>/identity/`,
  `<agent-root>/requests/`, `service/`, `cmd/`, local HTTP endpoints, tests,
  and operator documentation.

## Alternatives

### A. Participant-owned encrypted local identity plus explicit approval exchange

Each personal agent holds an encrypted local device private key under its own
`identity/` directory, unlocked only for its local process. The durable history
contains its public authorization. A terminal writes an unsigned request with
the exact target pCID, canonical payload bytes, expiry, one-time approval
token, and return address. The participant's already-authorized personal agent
retrieves or receives the request, displays those exact values, signs a fresh
Grid record only after local approval, and returns the exact bytes using the
one-time token. The terminal stores only the unsigned request and returned
public evidence.

### B. Terminal-held browser key

The shared terminal browser creates or imports a participant key and signs the
request after account login.

### C. Account or relay approval service

The terminal asks a website or relay to hold a key or create an approval
signature for the account holder.

## Scenario analysis

### Normal operation

Alice uses her home agent and her personal phone agent, both authorized by her
root history. At a shared makerspace terminal she chooses a loan request. The
terminal shows an unsigned draft and sends its exact pCID/payload to Alice's
reachable personal agent. Alice sees the same bytes, approves, and her agent
returns an exact signed record. The terminal submits it to local ingress. Carol
performs a steward action through her own agent.

A preserves one participant history across embodiments. B creates a new
terminal-held authoring authority. C makes a service able to impersonate the
participant.

### Lost terminal, failure, and incomplete requests

If the terminal closes or the connection fails, A leaves only an expired
unsigned request. A later reply with the wrong token, different pCID, altered
payload, or expired deadline is rejected. Alice's agent may be offline; no
promise exists until her approved exact record is returned. B risks retained
browser key material. C turns service availability into an authoring
dependency.

### Concurrent and malicious actors

Mallory copies an account session, changes a label, replays a request, or
submits an unknown pCID. Under A, none gives Mallory a private key. The
participant agent verifies the request bounds, presents the exact fields, and
uses each approval token once. A received record remains independently
verifiable through participant history. B lets a shared endpoint accumulate
participant keys. C lets Mallory target account or relay approval machinery.

### Long-horizon evolution and scale

A keeps durable protocol bytes independent of local request transport. Agents
can later use direct local links, locally chosen rendezvous, or carriage
without changing author evidence. Encrypted private-key storage remains a
participant-agent embodiment, not a global identity service. B couples keys to
browser storage. C couples durable authorship to an operational service.

## Conclusion

B and C are rejected: neither preserves participant-owned author evidence at
a shared terminal. A survives.

## Output to decision framing

Alternative A is selected. DI-hibok locks encrypted `identity/device-key.json`,
bounded `requests/<request-id>.json` drafts, token-limited polling, and the
local approval API names.

## Decision status

locked: Alternative A by DI-hibok in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`.
