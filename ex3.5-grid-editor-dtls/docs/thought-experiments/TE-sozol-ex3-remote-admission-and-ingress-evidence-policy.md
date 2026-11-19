# Ex3 remote-admission and ingress evidence policy

TE ID: TE-sozol

## Status

decided

## Decision status

Locked by DI-pazis: use separate bounded peer-ingress observations, validate
before accepted persistence, retain only non-secret capability diagnostics, and
keep WebSocket failure evidence transport-local. This TE does not yet alter
ingress, WebSocket, or capability behavior.

## Decision under test

How should an Ex3 relay distinguish accepted signed envelopes from rejected
peer ingress, remote-admission denials, and WebSocket-carriage failures so
local evidence is durable where useful without turning invalid input, bearer
material, or transport noise into accepted PromiseGrid replay?

## Assumptions and trust model

- Ex3's accepted `message-log.jsonl` rebuilds relay projections on restart and
  feeds peer exchange. It must contain only envelopes whose pCID and payload
  the relay supports after verification. Source: DI-figak.
- Current ingress parses, verifies, and persists a signed envelope before all
  pCID-specific payload handling completes. A valid proof with an unsupported
  pCID, or a supported pCID with an invalid payload, can therefore be written
  before the final rejection. This is the boundary under test.
- Remote capabilities are bearer material and a bootstrap secret is long-lived
  local operator configuration. Neither must be recorded in raw evidence or
  logs. Source: DI-hadil; DI-povip.
- WebSocket framing and JSON are carriage-local. Broken upgrade/auth/frame
  traffic must not create a profile meaning or a new top-level PromiseGrid
  action. Source: DI-figak.
- Alice operates a relay, Bob sends a valid supported envelope, Carol sends a
  valid but unsupported profile, Dave sends malformed bytes, and Mallory sends
  expired/mismatched capabilities or malformed WebSocket frames.

## Alternatives

### A. Separate bounded peer-ingress observations from accepted replay

For bounded peer-envelope ingress, retain exact rejected bytes in CAS and
append one relay-local observation per receipt after classification; keep those
records out of `message-log.jsonl`, replay, and peer feed. Record remote
capability denial and WebSocket framing/auth failures only as non-secret local
operational diagnostics, not raw token/frame evidence.

### B. Keep all parsable signed envelopes in the accepted log

Continue persisting a verified envelope before pCID-specific support and
payload validation, even when projection/replay later fails.

### C. Reject everything without durable local evidence

Return errors for all unsupported or malformed ingress and retain neither raw
bytes nor observations.

## Scenario analysis

### Normal supported peer envelope

Under A, Bob's verified and fully supported envelope enters CAS, accepted log,
replay, and peer exchange exactly once. B has the same normal path. C has the
same normal path but offers no distinction when later failures occur.

### Carol sends a valid unsupported pCID

A retains Carol's bounded exact bytes in CAS and records that this relay has
`no_supported_handler`; it does not claim Carol is wrong globally and does not
make the envelope replay input. B lets a locally unsupported envelope poison
the accepted log and potentially break restart reconstruction. C prevents log
poisoning but loses useful local evidence for an upgrade or compatibility
diagnosis.

### Dave sends malformed or invalid-proof bytes

A applies the existing byte bound, retains the exact bounded bytes, records a
local classification such as malformed input or invalid proof, and excludes
them from accepted paths. B cannot safely create a meaningful accepted entry.
C makes later local assessment impossible.

### Mallory presents an expired or mismatched capability

A returns the denial to the requesting remote client and may retain a
non-secret operational diagnostic such as reason, time, transport, and relay
observer. It never retains the bearer token or bootstrap secret. B gives no
clean place for authorization evidence and risks mixing it with peer messages.
C loses rate/abuse diagnostics but avoids secret retention.

### Mallory sends a malformed WebSocket frame or JSON request

A reports an error or closes the socket according to the transport path,
without storing raw frame bytes. Framing failures are not envelopes and may be
unbounded/hostile. B has no safe accepted-log representation. C is similar on
transport but lacks the bounded peer-envelope evidence distinction.

### Restart, concurrent ingress, and long-horizon evolution

A gives the accepted log a strict replay contract while a separate append-only
observation stream preserves local exception facts. B makes restart behavior
depend on inputs the relay already rejected. C discards compatibility and
incident evidence. A adds storage/retention obligations, so observation
records need bounded input and must remain explicitly relay-local.

## Conclusions

B is rejected because accepted replay must not include an envelope the relay
cannot fully support and project. C is rejected because bounded peer-ingress
evidence is useful for compatibility and incident assessment.

A survives and is recommended: use an Ex2-style separate observation stream
for bounded peer-envelope ingress; ensure accepted CAS/log/replay mutation
happens only after full support and payload validation; retain no bearer token,
bootstrap secret, or raw WebSocket frame; and treat all observations as local
relay facts rather than global validity or intent claims.

## Decisions still requiring DF

1. **Rejected peer-envelope persistence:** retain bounded raw peer-envelope
   bytes in CAS plus one local observation per receipt, excluded from accepted
   log/replay/feed (recommended), keep them in the accepted log, or discard
   them entirely?
2. **Acceptance point:** validate pCID support and payload before accepted
   CAS/log/replay mutation (recommended), retain current persist-before-final
   validation, or defer accepted persistence until asynchronous processing?
3. **Remote capability denials:** record only non-secret local diagnostics
   (reason/time/transport/observer) and never bearer or bootstrap material
   (recommended), retain raw authorization material, or retain nothing?
4. **WebSocket failures:** return client errors/close the connection without
   durable raw-frame evidence (recommended), retain raw frames as observations,
   or treat frames as a new peer-envelope protocol?

## Implications for open work

- `fozoz.2` can implement a separate observation store and reorder accepted
  ingress after DF is locked.
- `fozoz.3` should test raw retention, one observation per receipt, accepted
  replay/feed exclusion, non-secret capability diagnostics, and WebSocket
  handling boundaries.
- The scope declaration and README must describe the resulting local evidence
  boundary without exposing secret material or overclaiming global authority.
