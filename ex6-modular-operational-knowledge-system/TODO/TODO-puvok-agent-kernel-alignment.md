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
