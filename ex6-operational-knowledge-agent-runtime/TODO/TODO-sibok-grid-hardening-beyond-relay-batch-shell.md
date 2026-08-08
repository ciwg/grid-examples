# TODO sibok - grid hardening beyond relay batch shell

## Decision Intent Log

ID: DI-novuk
Date: 2026-07-29 18:30:00
Status: active
Decision: Exchange workflow artifacts through a dedicated peer-authenticated endpoint carrying one exact artifact and lifecycle-envelope bundle, rather than through JSON relay batches or shared storage.
Intent: Prove PromiseGrid-style independent-node workflow exchange while keeping routine relay records and binary CAS artifacts separate.
Constraints: Each receiver verifies and stores bytes in its own CAS; received lifecycle evidence never automatically imports or activates a workflow; top-level PromiseGrid semantics remain `promise`.
Affects: `grid/peers.go`, workflow relay endpoint and transfer codec, multi-node simulation tests, peer-card metadata, workflow loader docs, and `docs/thought-experiments/TE-novuk-workflow-relay-endpoint.md`.

ID: DI-jifuk
Date: 2026-08-03 07:22:07 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Expose received workflow artifacts through a read-only CAS-derived inbox and import one only through `moks workflow inbox import <artifact-cid> <alias>`.
Intent: Make retained remote artifacts discoverable while keeping local import, activation, and execution as independent decisions made by the receiving operator.
Constraints: Scan existing artifact and workflow-evidence CAS stores; do not add a receipt ledger or top-level PromiseGrid action; group evidence by artifact CID; refuse inbox convenience import when matching valid evidence is absent.
Affects: workflow relay receiver view, workflow CLI, inbox tests, operator documentation, and `docs/thought-experiments/TE-rasih-workflow-receipt-inbox.md`.

ID: DI-rufir
Date: 2026-08-03 07:28:53 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Retain authenticated workflow-transfer sender identity in a local receipt-metadata sidecar keyed by evidence CID and artifact CID.
Intent: Make receipt inbox provenance auditable without rewriting signed evidence or treating receipt as a local workflow lifecycle event.
Constraints: Sidecar data is local projection metadata only; exact evidence bytes remain in `workflow-evidence`; incomplete metadata/evidence pairs are visible but cannot use inbox convenience import.
Affects: workflow relay import, receipt inbox projection, tests, and `docs/thought-experiments/TE-kitol-receipt-metadata-sidecar.md`.

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

ID: DI-zumep
Date: 2026-07-28 20:15:00
Status: active
Decision: Add per-record digest proofs to relay batches and verify those proofs during import before appending durable history.
Intent: Strengthen trust beyond whole-batch signatures by letting receivers detect tampering or mutation at the record level, including for unknown-family carriage.
Constraints: Keep the current batch format additive for now; export proofs on new batches and verify them when present; do not yet invent a full claim-proof or record-signature wire layer.
Affects: `grid/batch.go`, `kernel/runtime.go`, runtime tests, CLI relay tests, README, and ex6 current-state docs.

ID: DI-ravud
Date: 2026-07-28 20:35:00
Status: active
Decision: Add relay-carriage record signatures by the exporting peer and verify them during import.
Intent: Move durable trust deeper than batch-level signatures by letting receivers validate each carried record against the exporting peer's stable identity and key material.
Constraints: These signatures attest relay carriage, not semantic authorship; keep the layer additive to the current batch format; leave claim-level proofs and author-level record signatures for later.
Affects: `grid/batch.go`, `grid/peers.go`, `kernel/runtime.go`, runtime tests, README, and ex6 current-state docs.

ID: DI-luzef
Date: 2026-07-28 20:50:00
Status: active
Decision: Add claim-level proofs by having the exporting peer sign each implementation claim it advertises in a relay batch.
Intent: Make implementation claims trustworthy batch metadata instead of unsigned declarations, while keeping the trust root explicit and local to the exporting peer identity.
Constraints: These proofs attest what the exporting peer claims to implement, not a global truth about the package; keep the layer additive to the current batch format; leave third-party attestation and richer claim semantics for later.
Affects: `grid/batch.go`, `grid/peers.go`, `kernel/runtime.go`, runtime tests, README, and ex6 current-state docs.

ID: DI-sovem
Date: 2026-07-28 21:10:00
Status: active
Decision: Add semantic author-level signatures to durable record envelopes, signing local records at authoring time and verifying embedded author signatures during import and append.
Intent: Move trust from "the exporting peer carried this record" to "the named record content was signed when authored," while remaining backward-compatible with older unsigned records.
Constraints: Keep the signature layer additive and optional for legacy records; semantic author signatures are separate from relay-carriage signatures; use the local runtime key as the current author-signing root until richer author identity exists.
Affects: `records/envelope.go`, `grid/peers.go`, `kernel/runtime.go`, runtime tests, README, and ex6 current-state docs.

ID: DI-fogem
Date: 2026-07-28 21:25:00
Status: active
Decision: Add third-party claim attestations so non-exporting peers can countersign specific implementation claims in a relay batch.
Intent: Extend trust beyond exporter self-claims by letting outside peers attest individual claims without becoming the exporter or record author.
Constraints: Attestations must be distinct from exporter claim proofs; they sign indexed claims, remain additive to the current batch format, and do not yet express policy weight or quorum.
Affects: `grid/batch.go`, `grid/peers.go`, `kernel/runtime.go`, runtime tests, README, and ex6 current-state docs.

ID: DI-movek
Date: 2026-07-28 22:10:00
Status: active
Decision: Add runtime-owned attestation policy and quorum for implementation claims, scoped by protocol pCID and role.
Intent: Make third-party claim attestations locally meaningful by letting each runtime decide which claims need countersigners and how many independent attesters are enough before import.
Constraints: Keep policy local to the importing runtime; use minimum-count quorum only for now; reuse known peer registrations as the current attester identity set; do not invent weight, federation, or global consensus semantics yet.
Affects: `grid/policy.go`, `kernel/runtime.go`, `cmd/moks`, runtime and CLI tests, README, and ex6 current-state docs.

ID: DI-ravok
Date: 2026-07-28 22:40:00
Status: active
Decision: Extend local claim-attestation policy with attester classes and trust weights, and let peer entries carry local class/weight metadata.
Intent: Move beyond simple minimum-count quorum so a runtime can distinguish stronger and weaker countersigners while still keeping trust decisions local and explicit.
Constraints: Keep class and weight local to the importing runtime; default discovered or allowed peers to class `peer` and weight `1`; do not claim federation, reputation markets, or global consensus semantics.
Affects: `grid/peers.go`, `grid/policy.go`, `kernel/runtime.go`, `cmd/moks`, runtime and CLI tests, README, and ex6 current-state docs.

ID: DI-rumek
Date: 2026-07-28 23:05:00
Status: active
Decision: Extend local claim-attestation policy with federation labels and minimum distinct federation spread.
Intent: Make trust stronger than raw local weights by letting an importing runtime require attestations from more than one federation, while still keeping federation semantics explicit and local.
Constraints: Federation labels remain local metadata on known peers; spread is counted by distinct federation labels among matched attesters; do not claim global federation discovery, federation PKI, or cross-runtime consensus semantics.
Affects: `grid/peers.go`, `grid/policy.go`, `kernel/runtime.go`, `cmd/moks`, runtime and CLI tests, README, and ex6 current-state docs.

ID: DI-sotad
Date: 2026-08-03 08:06:34 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Ship `moks workflow overview` first as a human-first, read-only operator screen whose `NEXT:` line names `moks workflow inbox import <cid> <alias>` when a received artifact is ready.
Intent: Give a boss and team one clear workflow status briefing without requiring CID-oriented command assembly or allowing overview to change runtime state.
Constraints: No JSON mode, terminal UI, network call, pull, import, activation, or execution in the first slice; preserve existing commands as authoritative detail views.
Affects: workflow CLI, overview renderer/tests, operator documentation, and `docs/thought-experiments/TE-nakum-workflow-overview-operator-flow.md`.

ID: DI-gihor
Date: 2026-08-03 08:11:58 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make the overview's Recent activity section order workflow-run heads by a durable UTC timestamp recorded in a new v2 run-event selector; retain and display v1 heads with time unavailable.
Intent: Give operators an honest recent-activity view without treating content-address order as chronology or invalidating retained CAS history.
Constraints: New events use the v2 selector and eight slots; replay accepts v1 seven-slot events; timestamps are local projection ordering evidence, not distributed clock truth; overview remains read-only and human-only.
Affects: `kernel/workflow_runs.go`, `kernel/runtime.go`, `cmd/moks/main.go`, run/CLI tests, `docs/protocols/workflow-run-v2.md`, and operator documentation.

ID: DI-rupit
Date: 2026-08-07 22:09:45 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep the current workflow-lifecycle pCID unchanged: a replacement remains a separately imported immutable artifact CID, while revocation is a local withdrawal of availability that retains artifact and lifecycle evidence.
Intent: Keep ex6's shipped example code consistent with the PromiseGrid Development Guide's local-promise model and the frozen workflow-lifecycle protocol.
Constraints: Do not add a replacement-link field, automatic replacement behavior, or a new lifecycle operation under the current pCID. Prove revocation persists across restart and prevents later activation or workflow execution.
Affects: `kernel/workflows_test.go`, `kernel/workflow_runs_test.go`, `docs/protocols/workflow-lifecycle.md`, and this TODO record.

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

## Workflow Artifact Hardening

- [x] Prove that revocation persists across restart, retains its artifact and lifecycle evidence, and prevents reactivation or workflow execution; under the current protocol, a replacement remains a separate imported artifact CID. (DI-rupit; TE-gavuk)
- [ ] Support local package-directory capture into CAS and direct CAS workflow-artifact import with validation. (DI-lovek; TE-gavuk)
