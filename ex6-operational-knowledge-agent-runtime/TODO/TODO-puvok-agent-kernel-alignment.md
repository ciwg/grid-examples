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

ID: DI-dovak
Date: 2026-07-28 11:08:00
Status: active
Decision: Add focused trace filters for one candidate path or one downstream protocol so `route plan ... trace` can return only the relevant subset of planner steps on larger route sets.
Intent: Keep trace output readable without losing the ability to inspect the exact planner sequence for the one path or downstream protocol an operator cares about.
Constraints: Reuse the existing trace data rather than recomputing a second planner; keep the unfocused trace unchanged; use explicit filter kinds instead of free-form search.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-buvok
Date: 2026-07-28 11:09:13
Status: active
Decision: Add trace summary counts to route-plan trace output so filtered traces report the total planner steps, kept steps, dropped steps, and active filter.
Intent: Let operators immediately see how much of the full planner trace was filtered away.
Constraints: Keep unfocused trace summaries truthful too; derive counts from the already-recorded trace; do not change planner execution.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-zafek
Date: 2026-07-28 11:10:55
Status: active
Decision: Add explicit trace scope metadata to route-plan summaries so each summary states which `protocol_pcid` it describes and whether it belongs to the root plan or a downstream hop.
Intent: Let operators distinguish the root protocol trace summary from nested downstream trace summaries when parser and transform routes produce multi-hop plans.
Constraints: Reuse the existing trace data and downstream-plan explanation shape; do not change planner execution or focused-trace filtering.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-rukav
Date: 2026-07-28 11:13:40
Status: active
Decision: Add top-level downstream trace summaries to traced route-plan output so nested protocol hops expose their own scoped counts directly in the trace payload.
Intent: Let operators see separate downstream-hop trace summaries without having to drill into candidate explanation blocks.
Constraints: Reuse recursive traced subplans; keep the existing root trace summary; do not change planner execution or remove downstream summaries from candidate explanations.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-vatuk
Date: 2026-07-28 11:19:10
Status: active
Decision: Add stable hop-path labels to route trace summaries and use them to distinguish repeated downstream hops that share the same `protocol_pcid`.
Intent: Let operators tell same-protocol downstream hops apart in traced route-plan output.
Constraints: Keep the current protocol and scope metadata; preserve deterministic ordering; do not change planner execution.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-lupav
Date: 2026-07-28 11:24:05
Status: active
Decision: Add a short human-readable `hop_summary` string alongside `hop_path` in route trace summaries.
Intent: Let operators scan multi-hop trace output faster without parsing the full structured path label.
Constraints: Derive the summary from the same deterministic route metadata as `hop_path`; keep `hop_path` as the exact identity field.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-sovak
Date: 2026-07-28 11:28:40
Status: active
Decision: Add `hop_depth` and `hop_index` metadata to route trace summaries so downstream hops can be sorted and filtered by distance from the root protocol.
Intent: Make multi-hop trace output easier to analyze programmatically and operationally.
Constraints: Keep the metadata deterministic; use depth `0` and index `0` for the root summary; preserve existing hop identity fields.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-vobek
Date: 2026-07-28 11:33:20
Status: active
Decision: Add a `depth` trace filter mode that accepts exact depth like `1` and inclusive lower-bound depth like `2+`.
Intent: Let operators inspect only direct downstream hops, only deeper hops, or the full trace without changing the planner itself.
Constraints: Reuse the current trace filtering surface; keep existing `candidate` and `downstream` modes; treat malformed depth filters as no-op fallback to the full trace.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-zumok
Date: 2026-07-28 11:39:30
Status: active
Decision: Allow combined trace filters by treating trace filters as an ordered list of `kind target` clauses combined with logical AND.
Intent: Let operators narrow traces across multiple dimensions in one request, such as depth plus downstream protocol.
Constraints: Preserve the current single-clause forms; keep clause order accepted but semantically commutative; malformed clauses fall back to the unfiltered trace.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-zamuk
Date: 2026-07-28 11:44:10
Status: active
Decision: Add named trace scopes as clause presets, with `direct-hops`, `deep-hops`, and `one-branch:<candidate-id>` mapped onto existing clause logic.
Intent: Let operators request common routing views without spelling every clause manually.
Constraints: Keep presets as syntactic sugar over the current clause engine; preserve combined filters; treat unknown scopes as no-op fallback to the full trace.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-bemok
Date: 2026-07-29 09:18:00
Status: active
Decision: Add runtime-owned local trace scope aliases, persist them beside the existing route policy state, expose them through a `route scope` CLI family, and resolve them as additional `scope` targets during trace filtering.
Intent: Let operators define reusable local trace views without editing code or being limited to the built-in named scopes.
Constraints: Keep built-in scopes working unchanged; treat local aliases as additive operator configuration; preserve unknown-scope no-op fallback; keep alias resolution deterministic and non-executing.
Affects: `grid/policy.go`, `kernel/runtime.go`, `kernel/routes.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-rusek
Date: 2026-07-29 09:34:00
Status: active
Decision: Add trace scope introspection that can show both the raw stored clause list and the fully expanded clause list for a built-in or local scope name.
Intent: Let operators see exactly what a scope resolves to, especially when local aliases compose other aliases or built-in scopes.
Constraints: Keep expansion deterministic and non-executing; preserve built-in scope names as reserved; expose introspection through the existing `route scope` CLI family.
Affects: `kernel/routes.go`, `kernel/runtime.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-fusek
Date: 2026-07-29 09:48:00
Status: active
Decision: Extend trace scope inspection to report skipped expansion branches, including alias cycles and unresolved scope references, alongside the clauses that were successfully expanded.
Intent: Let operators see not just the final expanded clause list, but also why some scope branches were dropped during expansion.
Constraints: Keep filtering behavior unchanged; only inspection gains explicit skipped-branch diagnostics; keep cycle handling deterministic and non-recursive.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-zusek
Date: 2026-07-29 10:02:00
Status: active
Decision: Extend trace scope inspection so each expanded clause includes provenance showing which built-in scope or local alias chain produced it.
Intent: Let operators trace an expanded clause back through alias composition instead of only seeing the final flattened clause list.
Constraints: Keep route filtering behavior unchanged; add provenance only to inspection output; preserve deterministic alias expansion order.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-vusek
Date: 2026-07-29 10:14:00
Status: active
Decision: Add grouped scope inspection views that collect expanded clauses by provenance branch, while keeping the existing flat expanded list.
Intent: Let operators inspect composed scopes branch-by-branch instead of reconstructing branch structure from a flat clause list.
Constraints: Keep grouping deterministic; do not change expansion order or filtering behavior; treat grouping as additional inspection output only.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-busek
Date: 2026-07-29 10:26:00
Status: active
Decision: Add a short deterministic label and a human-readable summary string for each grouped provenance branch in scope inspection output.
Intent: Let operators scan grouped branches quickly without parsing the raw provenance path array first.
Constraints: Keep labels deterministic from grouped-branch order; preserve the raw provenance branch array; do not change grouping or filtering behavior.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-yusek
Date: 2026-07-29 10:41:00
Status: active
Decision: Attach skipped scope-expansion diagnostics to the grouped provenance branch they came from, while keeping the flat skipped list.
Intent: Let operators see which alias branch dropped a scope reference or hit a cycle instead of inferring that relationship from the flat skipped list alone.
Constraints: Keep the existing flat skipped list; keep branch attachment deterministic from expansion provenance; do not change filtering behavior.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-lusek
Date: 2026-07-29 10:56:00
Status: active
Decision: Add branch ordering and filtering helpers to scope inspection so grouped branches can be sorted or filtered by depth, label, or summary.
Intent: Let operators narrow and reorder grouped branch output without losing the existing full inspection surfaces.
Constraints: Keep raw and expanded clause lists unchanged; only the grouped branch view and its attached skips are reordered or filtered; preserve deterministic output.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-musek
Date: 2026-07-29 11:08:00
Status: active
Decision: Echo the active branch query in scope inspection output whenever grouped branches were sorted or filtered by query parameters.
Intent: Let operators see exactly which branch-query parameters produced the visible grouped branch view.
Constraints: Keep the query echo read-only; preserve current grouped branch behavior; omit the field when no branch query is active.
Affects: `kernel/routes.go`, `cmd/moks/main.go`, tests, and routing docs.

ID: DI-nusek
Date: 2026-07-29 10:27:57
Status: active
Decision: Add a short branch-query summary to scope inspection output that reports total groups, matched groups, hidden groups, and the effective ordering used for the grouped branch view.
Intent: Let operators see the size and ordering impact of an active branch query without reconstructing it from the grouped branch list by hand.
Constraints: Keep the summary read-only; preserve current grouped branch behavior; omit the summary when no branch query is active.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-pusek
Date: 2026-07-29 10:27:57
Status: active
Decision: Add branch-query diagnostics to scope inspection output for default label-order fallback and zero-match grouped-branch results.
Intent: Let operators tell whether an empty or unexpectedly ordered grouped-branch view came from query normalization rather than from missing route-scope data.
Constraints: Keep diagnostics read-only; preserve current grouped branch behavior; omit diagnostics when there is nothing notable to explain.
Affects: `kernel/routes.go`, tests, and routing docs.
