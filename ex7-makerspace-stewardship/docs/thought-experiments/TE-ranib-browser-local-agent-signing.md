# Browser-to-Local-Agent Signing Embodiment

TE ID: TE-ranib

## Status

superseded by TE-folok / DI-janup

## Decision under test

How a participant browser signs Ex7 promises and submits exact bytes to that participant's own local agent without making a remote service the keyholder.

## Assumptions

- Every participant runs an independent local Ex7 agent.
- A relay/feed carries bytes but never signs or authorizes for an agent.
- Alice must be able to inspect and approve the exact pCID-selected promise before her key signs it.

## Alternatives

### A. Browser-held WebCrypto key with local-agent byte submission

The browser creates/imports Alice's non-extractable Ed25519 key in browser secure storage. It constructs the canonical signing view, shows Alice the pCID, payload, and recipient/context, then signs locally. It submits the exact signed Grid bytes to Alice's local agent over a loopback-only authenticated session. The agent verifies and stores bytes but cannot sign as Alice.

### B. Local-agent key with browser approval

The local agent stores Alice's private key and signs after browser approval.

### C. Remote web account signing

A hosted service holds the key or receives an exportable key.

## Scenario analysis

In normal operation A makes the browser embodiment the keyholder and the local agent a verifier/store. If the agent is offline, Alice can retain a pending signed record locally and submit it later; no other node can manufacture it. B makes the agent the keyholder and reintroduces impersonation risk. C centralizes identity and availability.

If Mallory serves a malicious UI, A still requires the visible signing prompt to name exact pCID/payload, but browser-origin compromise remains a local endpoint risk. B expands the trusted computing base to the agent's key store. C gives Mallory or the host a global target.

For long-horizon evolution, A can add alternative embodiments that produce the same signed bytes. The local agent and relay do not need a browser-specific authority rule. B and C bind identity to one implementation host.

## Conclusions

B and C are rejected. A is recommended: browser-held non-extractable keys, explicit signing view, and loopback-only exact-byte submission to the participant's own agent.

## Decision status

needs DF
