# Charm.land UI Exploration for the Grid Examples

**Status:** proposal only — no exercise behavior, protocol, runtime, or UI has
been changed by this document.

**Purpose:** identify where Charm.land’s terminal tooling could make the
examples more approachable, demonstrable, and pleasant to operate without
mistaking a presentation improvement for a PromiseGrid or trust-model change.

## Executive summary

Charm is most compelling here as a **terminal-embodiment toolkit**, not as a
replacement for the browser surfaces. Its Go-native tools line up unusually
well with this repository: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
for stateful terminal applications, [Bubbles](https://github.com/charmbracelet/bubbles)
and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for components and
layout, [Huh](https://github.com/charmbracelet/huh) for guided forms,
[Gum](https://github.com/charmbracelet/gum) for interactive shell/demo steps,
and [Glow](https://github.com/charmbracelet/glow) for readable Markdown in a
terminal. Charm describes this collection as open-source tooling for building
terminal applications and forms. [Charm’s library overview](https://charm.land/)

The best early candidates are Ex5 and Ex6, where the primary operators already
work in a terminal. Ex1 is the best low-risk demo win. Ex2–Ex4 and Ex7 should
keep their browser-first embodiments; Charm should add a complementary
operator/inspection surface or a better demo flow, not duplicate their UI.

## Full Charm tool scan

This review covers the eight libraries and seven applications currently
presented by Charm. The library catalogue describes Bubble Tea, Huh, Lip
Gloss, Wish, Bubbles, Glamour, Log, and Harmonica; the applications catalogue
adds Pop, Mods, Wishlist, VHS, Soft Serve, Glow, and Skate.
[Charm libraries](https://charm.land/libs/) and
[Charm applications](https://charm.land/apps/)

| Tool | What it is | Fit for this repository | Recommendation |
| --- | --- | --- | --- |
| **Bubble Tea** | Stateful Go TUI framework. | Excellent for Ex1, Ex4, Ex5, Ex6, and Ex7; useful for Ex2/Ex3 inspector tools. | Primary interactive foundation. |
| **Bubbles** | Reusable Bubble Tea components (lists, text inputs, viewports, spinners, etc.). | Excellent companion to Bubble Tea. | Use it instead of reimplementing standard terminal controls. |
| **Lip Gloss** | Terminal style and layout toolkit. | Excellent companion to every new TUI. | Establish a small, accessible shared visual vocabulary after the first TUI proves out. |
| **Huh** | Terminal forms and prompts. | Strong for guided demo setup and deliberate write actions. | Use for explicit, reviewable forms; never as evidence of identity or authorization. |
| **Gum** | Interactive shell-script prompts and helpers. | Excellent for existing demo scripts and fixture selection. | Best low-effort improvement for Ex1, Ex3, Ex4, and Ex7 demos. |
| **VHS** | Code-defined terminal recordings/GIFs. | Excellent for READMEs, how-tos, regression-visible demos, and eventual proof. | Add only after a stable terminal flow exists; prioritize Ex1 and Ex6 recordings. |
| **Glow** | Markdown reading CLI/TUI. | Strong for all docs-heavy exercises. | Use as an optional reader for guides, protocol explanations, run reports, and item bodies. |
| **Glamour** | Go Markdown-rendering library. | Useful when Markdown must be rendered *inside* a Go TUI rather than launched in Glow. | Use selectively inside Ex5/Ex6 detail panes; do not add both it and Glow without a concrete use. |
| **Log** | Go logging library. | Potentially useful for new TUI diagnostics, but not itself a UI upgrade. | Evaluate only in the implementation task; preserve existing logging/evidence behavior. |
| **Harmonica** | Physics-based terminal animation library. | Low fit for evidence and operations tools. | Avoid initially: motion can distract from dense, auditable information. |
| **Wish** | Library for SSH applications. | Conditional fit only for a future remote terminal embodiment. | Do not introduce now; it would add a new network and authentication boundary. |
| **Pop** | Terminal email application. | Low direct fit. | Do not couple it to workflow notifications without a separately approved notification/e-mail design. |
| **Mods** | AI CLI. | Outside the requested UI scope and can introduce model/data handling questions. | Do not integrate into exercises now. |
| **Wishlist** | SSH directory/bastion tool. | No current fit. | Treat as an operator’s personal tool, not an exercise dependency. |
| **Soft Serve** | Self-hosted SSH Git server. | No direct fit for the exercises’ product UI. | Out of scope. |
| **Skate** | Personal key-value store. | Poor fit: the examples already make explicit durable-storage and evidence claims. | Do not use as application storage or evidence persistence. |

The practical shortlist is therefore **Bubble Tea + Bubbles + Lip Gloss** for
interactive screens, **Huh** for intentional forms, **Gum** for scripts,
**VHS** for demonstrations, and **Glow/Glamour** for reading. The remaining
tools are either implementation support or introduce a boundary this proposal
does not need.

## How to read the recommendations

| Label | Meaning |
| --- | --- |
| **Fit** | How naturally the exercise’s current embodiment and workflows suit a Charm addition. |
| **Difficulty** | Relative implementation effort: S (days), M (about a week), L (multi-week), XL (larger cross-cutting effort). |
| **Priority** | Suggested sequencing for a future, separately approved implementation effort. |

Every recommendation below is deliberately bounded to local presentation,
operator guidance, or use of already available local APIs. None should alter a
pCID, a signed-envelope shape, a trust decision, or who can author a promise.

## Shared design principles

1. **Keep the scriptable CLI.** A TUI is an additional embodiment, not a
   replacement for commands, JSON/text output, automation, or tests.
2. **Read first; make mutations explicit.** Use a focused detail view before
   an action, and show the selected record, target, and consequences before a
   form is submitted.
3. **Do not make a UI look more authoritative than it is.** Relay-local
   observations, retained unknown records, unsigned requests, and local policy
   must retain their existing labels and boundaries.
4. **Favor progressive disclosure.** A compact queue, a focused record view,
   and a visible key map are more useful than reproducing every command on one
   screen.
5. **Use the browser where it is already the better tool.** Rich shared text
   editing and complex browser workflows should remain in their existing
   browser/Neovim embodiments.

## Per-exercise opportunities

### Ex1 — Order flow

**Current shape:** Docker-driven multi-agent order fulfillment demo with
seller, warehouse, accounting, carrier, kernel, collector, and analyzer. It
already exposes happy, refusal, and timeout fixtures, then leaves local
evidence artifacts for inspection. Review basis: [README](../ex1-order-flow/README.md).

| Proposal | How it helps | Difficulty | Priority |
| --- | --- | --- | --- |
| Add a `Gum`-guided demo launcher that selects a fixture, explains the roles, confirms the disposable runtime root, and opens the next inspection step. | Makes the first demo self-explanatory without hiding the existing `bash docker/run-demo.sh` path. | S | 1 |
| Add a read-only Bubble Tea “run storyboard” over the existing analyzer and collected artifacts. | Presents the order path, each role’s response, and the exact distinction between signed refusal and locally observed timeout in one navigable view. | M | 2 |
| Render the final analyzer report and design notes with Glow. | Gives a polished terminal handoff with no new application behavior. | S | 3 |

**Guardrails:** do not turn the storyboard into a claim that the collector is
the global authority or that an observation proves another role’s intent. The
current JSONL/CAS artifacts remain the inspectable evidence.

### Ex2 — Grid Editor (relay/polling embodiment)

**Current shape:** browser and Neovim collaborative editor over a local Go
relay; the browser already has rich document, presence, metadata, review,
publish, and diagnostics surfaces. Review basis:
[README](../ex2-grid-editor/README.md) and
[browser shell](../ex2-grid-editor/web/index.html).

| Proposal | How it helps | Difficulty | Priority |
| --- | --- | --- | --- |
| Add a small `gridctl` Bubble Tea relay/document inspector: document list, peer presence, relay identity, pCID labels, publication list, and local observations. | Gives an operator a fast, keyboard-first health and evidence view without competing with CodeMirror or Neovim editing. | M | 2 |
| Use Huh for a guided “open or create document” and relay-connect workflow. | Reduces command-flag discovery for demos while retaining explicit document and relay selection. | S–M | 3 |
| Use Glow for protocol/design/help pages invoked from the terminal workflow. | Makes the protocol boundary legible at the point of operation. | S | 4 |

**Not recommended:** rebuilding the collaborative editor itself in Bubble Tea.
The existing browser/Neovim surfaces are better suited to rich text editing,
cursor awareness, comments, and split Markdown preview.

### Ex3 — Grid Editor (WebSocket embodiment)

**Current shape:** the WebSocket variant of Grid Editor, with a browser-side
PromiseGrid flow card and live relay-observed message trace in addition to the
document editor. Review basis:
[README](../ex3-grid-editor-websocket/README.md) and
[browser shell](../ex3-grid-editor-websocket/web/index.html).

| Proposal | How it helps | Difficulty | Priority |
| --- | --- | --- | --- |
| Add a Bubble Tea live transport inspector with a compact event list and a drill-in pane for decoded fields, pCID, envelope identity, and local relay observation. | Makes the WebSocket/peer-feed/fanout story dramatically easier to demonstrate in a terminal without cluttering the editor. | M–L | 2 |
| Add focused filters for document ID, transport stage, message family, and error/observation state. | Lets a presenter isolate one causal path instead of scrolling through raw traffic. | M | 3 |
| Add a Gum-based two-relay simulation chooser around the existing scripts. | Makes it easy to select a safe, reproducible topology before launching it. | S | 4 |

**Guardrails:** raw bytes and decoded payloads should be clearly labelled as
relay-observed local evidence. Never imply that a WebSocket frame is anonymous
or that tracing changes envelope admission.

### Ex4 — Bug Tracker

**Current shape:** browser-first issue queue and detail/timeline UI with a
small engineer CLI for assigned work. Review basis:
[README](../ex4-bug-tracker/README.md) and
[browser shell](../ex4-bug-tracker/web/index.html).

| Proposal | How it helps | Difficulty | Priority |
| --- | --- | --- | --- |
| Add a Bubble Tea engineer workbench: assigned-issue list, filters, issue detail, append-only timeline, and keyboard handoff to the existing command action. | Turns the existing CLI into a genuinely useful daily triage surface while preserving the browser for fuller issue management. | M | 2 |
| Use Huh for new issue, assignment, state transition, and comment forms. | Provides validation and a clear review step before a durable action. | M | 3 |
| Add a Gum-driven seeded-demo chooser and a Glow issue/timeline reader. | Makes the demo and a read-only incident handoff much easier to show. | S | 4 |

**Guardrails:** the form must continue to use the selected built-in local
identity and existing service checks; styling must not suggest global identity
or shared proof.

### Ex5 — Operational Knowledge System

**Current shape:** mature browser, CLI, and Neovim embodiments over the same
local operational runtime. The browser is deliberately review-first and
includes queues, context drilldowns, live drafts, run/evidence/approval flows,
and search. Review basis:
[README](../ex5-operational-knowledge-system/README.md) and
[browser shell](../ex5-operational-knowledge-system/web/index.html).

| Proposal | How it helps | Difficulty | Priority |
| --- | --- | --- | --- |
| Add an `oks` Bubble Tea review workbench over the existing direct local socket: draft queue, problem hotspots, search results, and record/timeline drilldown. | This is the strongest Charm opportunity: it makes the existing review-first model feel immediate and coherent for keyboard-centric operators. | L | 1 |
| Add staged Huh flows for “log work,” “attach evidence,” and “review record.” | Mirrors the browser’s deliberate progressive disclosure in the terminal and reduces accidental context loss. | M–L | 2 |
| Use Lip Gloss to establish a shared terminal visual language for status, record kinds, warnings, context breadcrumbs, and queue density. | Makes a complex operational model scannable without changing the runtime. | M | 2 |
| Add Glow-powered reading for item bodies, procedures, and evidence summaries. | Gives long-form operational knowledge a much stronger terminal reading experience. | S–M | 3 |

**Implementation shape:** retain `oks-cli` for scripting and direct commands;
add a separate interactive command or a clearly optional TUI mode. Use the
same direct local socket contract rather than introducing a browser-only or
second service path.

**Guardrails:** record IDs, revision/approval boundaries, attachment
immutability, and the distinction between a live draft and a durable revision
must remain visible. A pleasant review screen must not make local policy look
like universal authority.

### Ex6 — Operational Knowledge Agent Runtime (OKAR)

**Current shape:** terminal-first `moks` runtime with package activation,
route planning, peer/trust controls, workflow artifacts, durable history, and
27 frozen package-record families. It has no browser or Neovim embodiment.
Review basis: [README](../ex6-operational-knowledge-agent-runtime/README.md).

| Proposal | How it helps | Difficulty | Priority |
| --- | --- | --- | --- |
| Add a Bubble Tea `moks overview` workbench for active packages, routes, peer-policy state, workflow readiness, inbox attention, and the next safe read-only command. | Gives the most complex terminal-first example a discoverable, calm home screen without sacrificing its explicit commands. | L | 1 |
| Add route-plan explorer with candidate comparison, downstream-hop tree, trace filters, and detail pane. | Converts a dense JSON/text analysis task into an excellent teaching and operations interface. | L | 2 |
| Use Huh for explicit peer-policy promotion, route-policy changes, and workflow import choices. | Makes consequential local-policy choices reviewable before the existing command executes. | M | 2 |
| Use Glow for package manifests, workflow definitions, route explanations, and implementation-claim documents. | Improves terminal documentation consumption with almost no semantic risk. | S | 3 |

**Critical guardrails:** preserve every noninteractive `moks` command and its
machine-readable output. Never put an action behind a default confirmation,
and never describe discovery, a signature, installation, or a peer card as a
trust grant.

### Ex7 — Makerspace Stewardship

**Current shape:** browser projection for equipment, eligibility, policies,
recognized authority, and evidence. The browser can submit already signed
record bytes and create an unsigned request; it must not pretend that a
selected member or browser account signs a promise. Review basis:
[README](../ex7-makerspace-stewardship/README.md) and
[browser shell](../ex7-makerspace-stewardship/web/index.html).

| Proposal | How it helps | Difficulty | Priority |
| --- | --- | --- | --- |
| Add a read-only Bubble Tea stewardship console for equipment state, member eligibility, recognized authority, retained/recognized record counts, and evidence timeline. | Gives operators a fast audit and demo surface while honoring the existing signed-record boundary. | M | 2 |
| Add a terminal signed-record inspection and submission assistant that validates/presents bytes before calling the current ingress path. | Makes opaque base64 record handling less error-prone and more explainable. | M | 3 |
| Use Gum/Glow for the two-agent proof walkthrough and evidence report. | Makes the trust-boundary demonstration easier to reproduce and present. | S | 3 |

**Critical guardrails:** do not create a “sign as member” terminal control.
Any approval UI must remain an unsigned request routed to the actual
participant agent; retained unknown records must remain evidence, not
projection input.

## Recommended delivery sequence

| Phase | Scope | Why it comes first |
| --- | --- | --- |
| **1 — demonstrate** | Ex1 Gum/Glow demo flow; Ex6 Glow readers; Ex7 proof walkthrough; VHS recordings once each flow is stable. | Small changes with immediate presentation value and minimal runtime coupling. |
| **2 — prove the interaction model** | Ex6 `moks overview` and one route-plan explorer slice, or Ex5 review workbench. | Tests whether a Charm TUI improves dense operational work before spreading a framework across modules. |
| **3 — extend selectively** | Ex4 engineer workbench; Ex2 relay inspector; Ex3 transport inspector; Ex7 stewardship console. | Reuse established conventions only where terminal work is a real complement to the browser. |
| **4 — form polish** | Huh confirmation/forms and shared Lip Gloss conventions where the selected TUI has proven useful. | Avoids building a large visual system before the navigation and operator tasks earn it. |

## Decisions to make before any implementation

1. **Choose the first proving ground:** Ex6 offers the clearest terminal-first
   fit; Ex5 offers the most visible operator workflow; Ex1 offers the fastest
   demo payoff.
2. **Choose an embodiment boundary:** separate `*-tui` commands are safest for
   script compatibility; adding an interactive mode to an existing command is
   more convenient but needs careful flag/output design.
3. **Choose a support floor:** terminal color capability, mouse support,
   terminal resize behavior, accessibility expectations, and noninteractive
   fallback should be explicitly defined.
4. **Define verification:** each TUI slice should have deterministic model/
   update tests, a noninteractive command-path regression test, and a
   repeatable terminal recording or screenshot proof for the demo surface.

## Bottom line

Charm can make these examples feel as intentional and inviting in a terminal
as their ideas deserve. The highest-value move is not “convert everything to a
TUI.” It is to add a focused, beautiful terminal embodiment where the
operator already works: Ex6 first for runtime/route/trust visibility, Ex5
next for operational review, and Ex1 for the most immediately delightful demo.
