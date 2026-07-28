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
