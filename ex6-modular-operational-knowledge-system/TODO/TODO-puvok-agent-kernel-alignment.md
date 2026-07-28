# TODO puvok - agent kernel alignment

## Decision Intent Log

ID: DI-puvok
Date: 2026-07-28 23:35:00
Status: active
Decision: Record the boss-note alignment that the ex6 runtime should eventually be understood as node services / agent roles, while packages/apps should eventually be understood as agents.
Intent: Keep current ex6 work aligned with stronger architecture guidance without throwing away the runtime-centered direction that already fits PromiseGrid-oriented modular growth.
Constraints: Treat current ex6 runtime/package embodiment as intermediate; preserve runtime-centered reasoning; do not yet force a full rewrite from manifests into pure startup-promise routing.
Affects: `docs/agent-kernel-alignment.md`, README links, and future kernel/app boundary work inside ex6.

ID: DI-rutom
Date: 2026-07-28 10:15:43
Status: active
Decision: Make protocol routing explicit as a kernel service by deriving route registrations from package claims, requiring `family-validator` routes for registered families, and exposing the route table from the CLI.
Intent: Move ex6 one concrete step toward the routing-agent model from the boss notes without throwing away the current manifest/self-check package boundary.
Constraints: Keep the current package manifest format; treat claim-derived routes as the intermediate embodiment of startup promises; do not introduce parser-agent hops yet.
Affects: `kernel/` route registration and activation behavior, `cmd/moks/` route inspection, runtime/package tests, and ex6 docs that describe the kernel/app boundary.

ID: DI-ruvot
Date: 2026-07-28 10:23:13
Status: active
Decision: Export the claim-derived route model in relay batch metadata, keep all declared routing roles visible there, and validate that exported routes are consistent with exported implementation claims.
Intent: Make routing roles visible across runtimes so ex6 does not hide the routing model inside one local process while it is moving toward a more agent-shaped kernel.
Constraints: Reuse the current batch/signature surface where possible; do not invent parser-agent metadata yet; keep route metadata derivative of package claims rather than a second independent declaration source.
Affects: `grid/` batch types and validation, `kernel/` batch export and route translation, runtime tests, and ex6 docs describing routing roles.

ID: DI-lafek
Date: 2026-07-28 10:30:40
Status: active
Decision: Let package claims model explicit `direct`, `parser`, and `transform` route types, require parser/transform routes to declare their emitted protocols, and carry that metadata through the route table and relay exports.
Intent: Make multi-hop routing visible in ex6 so the current kernel can describe “parser first, then app handler” style flows without pretending every protocol route is a direct terminal handler.
Constraints: Keep the claim-derived route model as the single source of truth; default missing route type to `direct`; do not yet implement runtime parser execution or parser-specific trust policy.
Affects: `packages/manifest.go`, `kernel/routes.go`, `kernel/runtime.go`, relay batch metadata, CLI route inspection, tests, and routing docs.

ID: DI-fotav
Date: 2026-07-28 10:33:55
Status: active
Decision: Add a real routing query surface that filters routes by input `protocol_pcid`, expose it as `moks route inspect <protocol-pcid>`, and return the matching direct/parser/transform routes in machine-readable JSON.
Intent: Let route consumers ask what handlers or hops exist for one protocol without scraping the whole route table.
Constraints: Query the existing claim-derived route table only; do not execute parser/transform hops yet; keep missing-protocol queries deterministic and non-fatal.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, CLI/runtime tests, and routing docs.

ID: DI-pabut
Date: 2026-07-28 10:38:28
Status: active
Decision: Add a route-plan selector for one input `protocol_pcid`, prefer executable direct paths over parser/transform hops, and expose the preferred plan plus ordered candidates as `moks route plan <protocol-pcid>`.
Intent: Move ex6 from route discovery toward actual kernel routing choice while keeping the current route table as the single source of truth.
Constraints: Choose plans from the current claim-derived route table only; treat parser/transform paths as executable only when their emitted protocols also resolve to preferred downstream plans; do not execute any plan yet.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, CLI/runtime tests, and routing docs.

ID: DI-matek
Date: 2026-07-28 10:40:49
Status: active
Decision: Make route-plan ordering policy-driven with runtime-owned global prefer/avoid lists for route types and roles, persist that policy beside the existing runtime policy state, and expose it through CLI show/set commands.
Intent: Replace one fixed built-in plan order with explicit operator-owned routing preference while keeping the planner deterministic.
Constraints: First slice is global, not per-protocol; prefer/avoid changes ordering but does not by itself make a non-executable route executable; keep the current route table as the only planning input.
Affects: `grid/policy.go`, `kernel/routes.go`, `kernel/runtime.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-posek
Date: 2026-07-28 10:47:56
Status: active
Decision: Keep the current global route-plan policy as the default, add per-`protocol_pcid` route-plan policy overrides, and have per-protocol overrides replace only the specific prefer/avoid fields they set while inheriting all other fields from the global policy.
Intent: Let the routing planner adapt to protocol-specific needs without forcing one global route preference order onto every protocol.
Constraints: Keep the claim-derived route table as the only planning input; preserve deterministic planning; treat empty override fields as inherit-from-global instead of clearing global defaults.
Affects: `grid/policy.go`, `kernel/routes.go`, `kernel/runtime.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-rivuk
Date: 2026-07-28 11:04:00
Status: active
Decision: Add role-scoped route-plan policy overrides keyed by `protocol_pcid + role`, apply them on top of the global and per-protocol planner policy, and let the planner evaluate each candidate route against the effective policy for that candidate role.
Intent: Let one protocol prefer one route role while still avoiding another without forcing every role on that protocol through one shared preference set.
Constraints: Keep route planning deterministic; preserve global and per-protocol inheritance; do not add route execution yet; treat role-scoped overrides as explicit exact-role policy, not wildcard matching.
Affects: `grid/policy.go`, `kernel/routes.go`, `kernel/runtime.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-lavik
Date: 2026-07-28 10:53:52
Status: active
Decision: Make `moks route plan <protocol-pcid>` include plan introspection that explains candidate executability, the active global/protocol/role policy layers for each route, and why the winning route outranked the next candidate.
Intent: Let operators inspect not just what route won, but why it won under the current layered planner policy.
Constraints: Keep the current `route plan` JSON surface as the primary output; do not add route execution; explanations must remain deterministic and derived from the current route table and planner policy only.
Affects: `grid/policy.go`, `kernel/routes.go`, `kernel/runtime.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-fobek
Date: 2026-07-28 10:58:15
Status: active
Decision: Extend route-plan introspection with pairwise comparison detail for every candidate pair, and expose the exact tie-break and policy reasons that made one candidate rank ahead of the other.
Intent: Let operators inspect the full ordering logic across the whole candidate set instead of only seeing the final winner explanation.
Constraints: Keep comparison output deterministic and derived only from the current planner rules; do not execute routes; keep the comparison explanation inside the existing `route plan` JSON surface.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-povak
Date: 2026-07-28 11:02:06
Status: active
Decision: Add explicit downstream-plan explanation summaries to parser and transform candidates so multi-hop route plans expose nested preferred routes and downstream winner reasons end to end.
Intent: Let operators understand not just that a parser/transform route depends on downstream protocols, but exactly how those downstream plans resolved.
Constraints: Reuse the existing recursive route-plan structure; do not execute routes; keep downstream explanations deterministic and derived from the nested route plans already built by the planner.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-rusom
Date: 2026-07-28 11:06:21
Status: active
Decision: Add an explicit `trace` mode for `moks route plan <protocol-pcid>` that records the planner's actual step-by-step decision sequence, including candidate discovery, pairwise comparisons, and preferred-route selection.
Intent: Let operators inspect the exact comparison order the kernel executed instead of reconstructing it from the final sorted output.
Constraints: Keep normal `route plan` output unchanged unless trace mode is requested; keep trace output deterministic; do not execute routes.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, tests, and routing docs.
