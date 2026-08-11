# Terminal-to-Agent Cross-Origin Transport

TE ID: TE-jusur

## Status

decided

## Decision under test

How a shared terminal browser transports an unsigned approval request to a
participant-agent URL and polls the response without treating web origin,
account, CORS, or URL possession as author evidence.

## Alternatives

### A. Explicit user-supplied HTTPS or loopback target with narrow CORS

The terminal displays an explicit participant-agent URL. It permits only
`https` or loopback `http` URLs, sends the existing unsigned request, and uses
the returned one-time token only for polling. The participant agent allows the
terminal origin only for request creation and token-gated polling; local
approval remains loopback-only.

### B. Wildcard CORS

Allow every browser origin to create, poll, and approve requests.

### C. Account-routed request proxy

Use a makerspace account service to forward, authenticate, or approve the
request.

## Scenario analysis

Alice enters her own reachable agent URL at Bob's terminal. Under A, Bob's
origin may create an unsigned draft and poll its token, but cannot approve it.
Alice approves only through her local loopback page; the returned record is
verified through signed history. Mallory may send malformed target URLs or
replay a token, but gets neither a signing key nor approval authority.

B lets arbitrary pages manufacture request load and attempt approval. C makes
the account service a transport and authoring dependency. A preserves a
participant-controlled transport choice while keeping durable evidence in the
existing signed record contracts.

## Conclusion

Alternative A is selected by DI-fuzar. CORS is transport permission only; it
is never identity, signature, recognition, or approval evidence.
