# Workflow overview operator flow

TE ID: TE-nakum

## Status

decided

## Decision under test

How `moks workflow overview` should present EX6 workflow state to an operator,
a team lead, or an executive observer without requiring them to assemble
multiple CID-oriented commands.

## Assumptions and trust model

- Workflow receipt, local import, activation, and execution remain separate
  decisions.
- The overview is read-only and derives its information from existing runtime
  state; it must not pull images, import artifacts, or activate workflows.
- An operator needs exact next commands, while a boss needs an immediately
  legible account of current state and blockers.
- The command runs locally and must remain useful when inbox evidence or an
  adapter dependency is incomplete.

## Alternatives

1. **Raw JSON aggregation.** Emit the existing workflow, inbox, and readiness
   structures together as one JSON document.
2. **Human-first summary with an optional JSON mode.** Print a compact,
   deterministic status screen with an explicit `NEXT:` line; support
   `moks workflow overview --json` for automation.
3. **Interactive terminal dashboard.** Add live refresh, selection, and
   navigation inside a terminal UI.

## Scenario analysis

### Normal operation

Alice runs the overview before a shift. Option 1 contains all facts but makes
her interpret nested JSON. Option 2 reports active-and-ready workflows first,
then received artifacts waiting for an explicit import, and finishes with the
next safe command. Option 3 adds navigation but hides the simple one-command
briefing behind an interactive interface.

### Failure, corruption, and incomplete state

Bob has a received artifact whose evidence or adapter image is unavailable.
Option 1 carries the raw reason but does not prioritize it. Option 2 puts it in
`NEEDS ATTENTION`, states the blocker, and names the non-destructive next
command. Option 3 can visually highlight it, but requires more terminal and
input-handling failure paths.

### Concurrent actors and mixed-version nodes

Carol receives artifacts while Dave runs the overview. A read-only snapshot is
safe under all options. Option 2 can state that it is a snapshot while keeping
ordering deterministic. Older nodes simply lack the command; no wire contract
or lifecycle meaning changes.

### Long-horizon scale

At many workflows, raw JSON remains machine-friendly but is not a team status
surface. A text summary can cap each section and show counts plus the most
urgent entries. A TUI creates an ongoing maintenance commitment before EX6 has
proved a stable operator vocabulary.

## Recommended output and operator flow

Option 2 survives. The default command should print a concise fixed-order
screen:

```text
WORKFLOW OVERVIEW

Ready: 2 active workflows
  [ready] procedure-execution
  [ready] maintenance-round

Needs attention: 1
  [inbox] bafk... from peer-alice — ready to import

Recent activity: none

NEXT: moks workflow inbox import bafk... receiving-check
```

The overview derives four sections:

1. **Ready** — active workflows whose readiness reports are eligible to execute.
2. **Needs attention** — inactive or blocked workflows and all inbox entries,
   with their readiness reason.
3. **Recent activity** — the most recent local workflow-run summaries, when
   present.
4. **NEXT** — exactly one safe, deterministic operator command. It prefers the
   first CID-sorted importable inbox artifact, then an image pull or diagnostic
   verification of a dependency blocker, then activation of a verified imported
   workflow; otherwise it says that no action is required.

The first slice is human-facing only. It has no color, terminal control
sequences, network calls, or mutations, so it can be pasted into a meeting note
or run on a minimal host.

## Decision status

Locked by DI-sotad. EX6 first ships a human-first overview only; a future JSON
representation requires a separate decision.

## Implications for open work

- Add deterministic overview rendering tests for ready, inbox, blocked, and
  empty states.
- Keep all underlying lifecycle and readiness commands authoritative; overview
  is an operator projection only.
- Do not add interactive terminal behavior until a later operator-use review.
