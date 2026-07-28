# TODO sibok - grid hardening beyond relay batch shell

## Decision Intent Log

ID: DI-sibok
Date: 2026-07-28 18:20:00
Status: active
Decision: Harden the current relay path by deduplicating exact-byte records during append/import and validating relay batch metadata before accepting a batch.
Intent: Keep repeated relay exchange from inflating append-only history with identical records while preserving unknown-family exact-byte carriage and making malformed batches fail early.
Constraints: Keep the current batch wire shape; do not invent a speculative final PromiseGrid transport; dedupe exact bytes only, not semantic near-matches.
Affects: `store/history.go`, `kernel/runtime.go`, `grid/batch.go`, runtime tests, and ex6 current-state docs.

ID: DI-nasek
Date: 2026-07-28 18:35:00
Status: active
Decision: Add a minimal live peer exchange layer with HTTP `relay serve`, `relay pull`, and `relay push` commands that carry the existing validated relay batch format.
Intent: Move ex6 beyond file-only batch exchange while keeping the current protocol narrow, testable, and compatible with the hardened batch layer.
Constraints: Reuse the current batch wire shape; no speculative peer discovery or trust protocol yet; keep endpoints explicit and local-runtime owned.
Affects: `cmd/moks`, relay CLI tests, and ex6 current-state docs.

ID: DI-rupem
Date: 2026-07-28 18:50:00
Status: active
Decision: Add runtime-owned peer allow rules, stable local peer identity, and batch-to-peer identity matching for live relay pull/push/import.
Intent: Make multi-peer exchange safer by requiring explicit peer registration and by rejecting live imports whose declared batch identity does not match the allowed peer entry.
Constraints: Keep trust policy explicit and local; no automatic peer discovery; no cryptographic signing layer yet.
Affects: `grid/peers.go`, `kernel/runtime.go`, `cmd/moks`, relay tests, and ex6 current-state docs.

ID: DI-zotem
Date: 2026-07-28 19:05:00
Status: active
Decision: Add signed relay batches with runtime-owned local key material and required peer public keys for live peer exchange.
Intent: Make pull/push/import trust stronger than peer-ID matching alone by letting runtimes prove which allowed peer exported a live batch.
Constraints: Sign the current relay batch as a whole instead of inventing a new wire format; keep trust rooted in explicit local peer registration; leave broader record-level proofs and discovery for later slices.
Affects: `grid/peers.go`, `grid/batch.go`, `kernel/runtime.go`, `cmd/moks`, relay tests, and ex6 current-state docs.

ID: DI-vemut
Date: 2026-07-28 19:20:00
Status: active
Decision: Add explicit peer discovery through a relay peer-card endpoint and a CLI discovery command that fetches peer metadata without auto-granting trust.
Intent: Let operators learn a peer's identity, public key, and relay endpoints from the grid itself while keeping allow/pull/push policy as a separate explicit step.
Constraints: Discovery must not silently enable pull or push; reuse the current live relay surface; keep peer registration human-auditable by printing the exact allow command to run next.
Affects: `grid/peers.go`, `cmd/moks`, relay tests, and ex6 current-state docs.

ID: DI-kasud
Date: 2026-07-28 19:35:00
Status: active
Decision: Let peer discovery optionally seed a local peer entry with `no-pull` and `no-push`, while keeping plain discovery read-only by default.
Intent: Remove manual transcription from the operator flow without collapsing discovery into trust or enabling exchange permissions implicitly.
Constraints: Seeding must remain explicit; seeded peers stay untrusted for exchange until a later `relay peer allow ... pull|push` command changes policy; discovery output must state that boundary clearly.
Affects: `cmd/moks`, relay tests, README, and ex6 current-state docs.

ID: DI-lutep
Date: 2026-07-28 20:00:00
Status: active
Decision: Add peer-policy promotion shortcuts that reuse stored peer metadata and only change pull/push policy.
Intent: Finish the explicit trust workflow by letting operators promote a seeded peer without retyping batch URLs, import URLs, or public keys.
Constraints: Promotion must stay explicit and local; it may only operate on already registered peers; it must not fetch fresh metadata during promotion.
Affects: `cmd/moks`, relay tests, README, and ex6 current-state docs.

## Goal

Make the current relay shell safer and less noisy under repeated exchange.

## Scope

- exact-byte record dedupe
- relay batch metadata validation
- repeated import idempotence tests
- current-state hardening note

## Why

Without this pass, repeated imports can append the same raw record forever and
malformed relay batches are not rejected early enough.
