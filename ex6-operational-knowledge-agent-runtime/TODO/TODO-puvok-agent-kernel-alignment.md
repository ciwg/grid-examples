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

ID: DI-nuvom
Date: 2026-08-07 22:24:19 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Replace package-claim-derived live route authority with paired, explicit app-acceptance and routing-role delivery promises stored as durable local records. Package claims remain bootstrap hints only and never become active promises automatically.
Intent: Keep the shipped ex6 example aligned with PromiseGrid's voluntary local-promise model instead of inferring a running agent's availability from package installation.
Constraints: The first slice remains local, deterministic, non-executing, and non-networked. Conditions, public names, and exact touched paths require further DF before code.
Affects: `docs/thought-experiments/TE-ravuk-agent-route-registration.md`, route-planning implementation, durable local state, tests, CLI, and documentation.

ID: DI-komaz
Date: 2026-08-07 22:40:36 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Give each first-slice route promise an explicit opaque local `agent_id`; it is neither a package ID nor a durable PromiseGrid identity. Any package association remains bootstrap or implementation metadata outside the promise's agent identity.
Intent: Keep live voluntary promises attributable to app agents without falsely treating a package label as identity or prematurely claiming signing-key continuity.
Constraints: This local non-networked slice has no signing or cross-node identity claim. A later signed protocol must use a new pCID if it changes the record's meaning or encoding.
Affects: `docs/protocols/route-promises.md`, route-promise record/replay code, route planning, tests, and CLI.
Supersedes: DI-nuvom

ID: DI-butam
Date: 2026-08-07 22:42:57 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Require an explicit local append-only `AgentBinding` that maps each opaque `agent_id` to installed implementation package and route metadata. A binding is local adapter metadata, not a promise and not an identity claim.
Intent: Let route planning connect voluntarily published app and routing-role promises to the current implementation while keeping package installation incapable of creating a live promise by inference.
Constraints: The planner must require a valid binding plus enabled receive and delivery promises. Bindings are retained in local CAS and rebuilt from it; no network exchange, signing, or automatic binding from package claims occurs in this slice.
Affects: `docs/protocols/route-promises.md`, route-promise record/replay code, runtime setup, route planning, tests, and CLI.

ID: DI-zolil
Date: 2026-08-07 22:48:30 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Expose first-slice local route evidence through `moks route bind <agent-id> <package-id> <true|false>`, `moks route promise receive <agent-id> <pcid> <true|false>`, and `moks route promise deliver <router-id> <recipient-id> <pcid> <true|false>`.
Intent: Make bindings and promises explicit operator actions that create retained local evidence instead of letting package installation silently create live route authority.
Constraints: Commands create only local CAS records; they do not communicate with peers, authenticate durable identity, or execute a route.
Affects: `cmd/moks/main.go`, CLI tests, route-promise documentation, and operator-facing route planning.

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

ID: DI-tusek
Date: 2026-07-29 10:31:47
Status: active
Decision: Report invalid branch-query filter values, starting with malformed depth filters, as ignored filters in scope inspection diagnostics.
Intent: Let operators see when a query value was accepted syntactically by the CLI but ignored semantically by grouped-branch filtering.
Constraints: Keep invalid-filter reporting read-only; preserve current grouped branch behavior; do not turn malformed depth filters into command errors in this slice.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-vusem
Date: 2026-07-29 10:34:06
Status: active
Decision: Report invalid sort values as ignored filters and explain that grouped-branch ordering fell back to label ordering.
Intent: Let operators distinguish an invalid `sort` input from a deliberate label-sort query.
Constraints: Keep invalid-sort reporting read-only; preserve current grouped branch behavior; do not turn invalid sort values into command errors in this slice.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-zusev
Date: 2026-07-29 10:38:02
Status: active
Decision: Report degenerate whitespace-only label and summary filters as ignored filters in scope inspection diagnostics.
Intent: Let operators see when a text filter was syntactically present in the query but trimmed away before grouped-branch matching.
Constraints: Keep degenerate-text reporting read-only; preserve current grouped branch behavior; do not turn whitespace-only text filters into command errors in this slice.
Affects: `kernel/routes.go`, tests, and routing docs.

ID: DI-lovek
Date: 2026-07-29 13:59:57
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Model workflows as immutable imported artifacts with append-only local lifecycle events and a rebuildable registry projection; support both local package-directory capture into CAS and direct CAS artifact import; distinguish `deactivated` from `revoked`; expose `ImportWorkflow`, `ActivateWorkflow`, `DeactivateWorkflow`, `RevokeWorkflow`, and `Workflows`.
Intent: Preserve grid-visible artifact provenance and historical interpretation while making route and Docker-worker eligibility an explicit local lifecycle decision rather than an automatic consequence of import.
Constraints: Import must not grant execution authority; deactivation/revocation must never delete CAS artifacts or durable record history; replacement is additive by artifact CID; lifecycle state is local runtime mechanics and not a new top-level PromiseGrid action kind.
Affects: `kernel/workflows.go`, runtime lifecycle persistence and startup, package/CAS import translation, Docker-worker eligibility, tests, and operator documentation.

ID: DI-bavuk
Date: 2026-07-29 14:38:44
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Store lifecycle events as exact pCID-selected CBOR `grid()` envelopes in CAS, linked by parent CID per workflow artifact; keep workflow IDs as local aliases; derive a disposable local projection cache; and define the lifecycle payload as a fixed-shape CBOR array.
Intent: Make CAS event history authoritative while preserving selective local retention, artifact-scoped replay, exact wire bytes, and explicit local lifecycle decisions.
Constraints: The pCID protocol specification is frozen before implementation; lifecycle mechanics do not introduce a new top-level PromiseGrid action kind; JSONL is not authoritative and may exist only as a rebuildable projection or diagnostic export; unknown pCIDs, invalid arities, and invalid parents produce explicit local non-commitments.
Affects: lifecycle protocol specification, `kernel/workflows.go`, CAS event storage and replay, local cache rebuild, tests, and operator documentation.
Supersedes: DI-lovek (lifecycle-persistence clause only; manual handle because `tools/mint-handle` was unavailable)

ID: DI-fofuh
Date: 2026-08-02 12:49:04 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Declare executable workflow adapters in the active package's `moks-package.json`; use `procedure-execution-adapter` as the first adapter package; and run every such adapter only through the Docker-confined worker boundary.
Intent: Keep package self-description, installation, activation self-check, and workflow-adapter authority in one inspectable contract while preventing an executable workflow from inheriting host or runtime authority.
Constraints: The adapter declaration identifies its adapter name, Docker image/command, and exact input/output pCIDs; the runtime must require the active artifact manifest and active package declaration to match exactly; workers receive exact CBOR on stdin and may return only typed CBOR plus proposed CAS/record writes; the runtime validates output pCID and applies durable writes itself; no direct host-process fallback, runtime-root mount, CAS/history/peer-key mount, ambient secrets, network, or Docker socket is permitted.
Affects: `moks-package.json` manifest contract, package activation self-check, Docker worker dispatch, workflow execution validation, package-author documentation, workflow documentation, and future tests.

ID: DI-harib
Date: 2026-08-02 13:41:27 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Distribute portable workflow-adapter images through digest-pinned OCI registry references, and permit a runtime to acquire an image only from an operator-configured registry allow-list.
Intent: Make a transferred workflow independently executable on a receiving node without allowing an active package manifest to choose arbitrary network destinations or expanding workflow relay into an image-content transport.
Constraints: Image references must remain registry-qualified immutable digests; workflow relay continues to carry artifact and lifecycle evidence only; a failed or disallowed pull leaves the adapter unavailable and does not alter lifecycle state; Docker worker confinement remains unchanged.
Affects: adapter package manifests, local registry policy persistence, image availability/readiness reporting, Docker pull verification, workflow verification, operator documentation, and future two-node tests.

ID: DI-hapak
Date: 2026-08-02 13:58:08 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Persist the portable-image registry policy as an exact canonical host[:port] runtime allow-list and expose it through `moks registry allow|list|remove`.
Intent: Keep image acquisition understandable and auditable as one local network-trust decision, without broad hostname patterns or package-lifecycle-coupled grants.
Constraints: Reject schemes, paths, wildcard characters, empty hosts, and non-canonical duplicates; do not add credentials, wildcard matching, per-package exceptions, or image transfer in this slice.
Affects: runtime policy persistence, registry CLI commands, OCI image-reference parsing, workflow readiness output, tests, and operator documentation.

ID: DI-zivut
Date: 2026-08-02 14:07:40 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Acquire a portable adapter image only through the explicit `moks workflow image pull <alias>` command.
Intent: Keep workflow verification read-only and prevent lifecycle/start commands from hiding registry network activity or Docker state changes.
Constraints: Resolve the image only from a verified artifact and matching installed package declaration; require the local registry allow-list and exact digest verification; pulling neither activates a workflow nor starts a run; unavailable images keep execution fail-closed.
Affects: workflow CLI, image availability/readiness state, Docker pull integration, tests, and operator documentation.

ID: DI-kojab
Date: 2026-08-07 22:34:58 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: For the first route-promise slice, define the record contract in `docs/protocols/route-promises.md`; use `ReceivePromise` and `DeliveryPromise` as durable local, append-only CAS records; support enabled-only conditions; rebuild any local index from retained record bytes.
Intent: Make ex6's routing model conform to the PromiseGrid Development Guide without adding unstated clock semantics or treating package claims as active agent promises.
Constraints: Package claims remain bootstrap hints only. The first slice remains local, deterministic, non-executing, and non-networked. Exact implementation and test paths must be approved before code edits.
Affects: `docs/protocols/route-promises.md`, `docs/thought-experiments/TE-ravuk-agent-route-registration.md`, route-planning implementation, durable local state, tests, CLI, and documentation.
Supersedes: DI-nuvom

ID: DI-bidam
Date: 2026-08-08 10:39:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Before invoking a Docker workflow worker, require active workflow state plus enabled local AgentBinding, ReceivePromise, and DeliveryPromise evidence for the adapter's owning package and workflow input pCID. If evidence is missing, disabled, malformed, or conflicted, return a pre-dispatch local refusal without creating a workflow-run event.
Intent: Make Docker dispatch obey the same explicit voluntary route evidence as planning, while keeping a non-dispatched worker from being misrepresented as a failed run or broken promise.
Constraints: Retain existing Docker confinement and post-acceptance run lifecycle events; do not claim network delivery, durable identity, or distributed locking; use only `kernel/runtime.go` and `kernel/workflow_runs_test.go` in this slice.
Affects: `kernel/runtime.go`, `kernel/workflow_runs_test.go`, and `docs/thought-experiments/TE-niliv-workflow-dispatch-route-evidence.md`.

ID: DI-guraj
Date: 2026-08-08 10:42:32 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Name the runtime's private adapter-name-to-active-package projection `workflowAdapterPackages`; keep the pre-dispatch eligibility check inline in `StartWorkflowRun`.
Intent: Make the stored relation legible without broadening the first dispatch-gate slice with a new public API or helper abstraction.
Constraints: The map is runtime-local derived state populated only during active-package registration; it does not create a binding or promise.
Affects: `kernel/runtime.go` and `kernel/workflow_runs_test.go`.

## Alignment Implementation Queue

- [x] Persist append-only local workflow lifecycle events and rebuild the local registry at startup. (DI-lovek; TE-gavuk; cdc0621)
- [x] Split workflow import from explicit activation; route and worker eligibility require active state. (DI-lovek; TE-gavuk)
- [x] Add separate deactivation and revocation withdrawal paths without deleting CAS or durable history. (DI-lovek; TE-gavuk)
- [x] Model pCID-scoped app receive promises and routing-role delivery promises before route-plan execution. (TE-ravuk; DI-kojab; DI-komaz; DI-butam; DI-zolil; 51b51cf)
- [x] Bind Docker worker dispatch only to active registered receive promises and record local lifecycle events. (TE-dovek; TE-niliv; DI-bidam; DI-guraj)
- [x] Extend active package manifests with Docker-confined workflow adapter declarations and bind matching active artifacts to those declarations. (TE-dovek; DI-fofuh; d86cef2)
