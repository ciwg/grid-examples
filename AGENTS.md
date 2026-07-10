# Repository Guidelines

## Project Structure & Module Organization
- Keep packages at module root or under purpose-named top-level directories (`contexts/`, `state/`, etc.); avoid `internal/` and `pkg/`.
- Do not commit local state files (for example `.grok`, `.grok.lock`) or generated binaries.

## Build, Test, and Development Commands
- `go test ./...` runs the test suite.
- `gofmt -w .` (or `go fmt ./...`) formats Go code.

## Agent Instruction Architecture (Required)
- `AGENTS.md` is the canonical home for repo-wide protocol, workflow, and vocabulary rules. Role-specific files such as `AGENTS-codex.md` and `AGENTS-ppx.md` are overlays: they may add stricter role constraints, environment setup, and role-specific procedures, but they must not duplicate or relax canonical repo-wide rules. (DI-034-20260508-060134)
- When a rule applies to every agent or every repo artifact, move it here and replace role-file copies with pointers. When a rule applies only to one agent's runtime environment, identity, credentials, branch lifecycle, or private logging system, keep it in that role overlay. (DI-034-20260508-060134)

## Promise Action Minimalism (Required)
- Future PromiseGrid protocol, simulation, POC, scoring, generation, and guide work must not invent workflow-specific top-level action kinds by default. The default future-facing top-level semantic action is `promise`. (DI-mosoj)
- Treat observation as a promise that the promiser observed something from its local vantage. Treat refusal as absence of a promise, a promise not to do something, or a promise that the agent does not currently promise the requested behavior. (DI-mosoj)
- Treat repair, offer, counteroffer, acceptance, routing, introduction, redemption, transfer, storage, computation, TCP-link changes, authorization, dispatch, grant, registration, and enforcement as pCID-defined payload semantics, local trust/evidence interpretation, or implementation-local mechanics unless a scoped TE/DI proves a distinct wire-level role. (DI-mosoj)
- Before adding any new top-level action kind, stop and answer: why is this not just an agent's voluntary promise with pCID-defined payload meaning? If that answer is not explicit in a locked TE/DI, do not add the action kind. (DI-mosoj)

## Decision-First Specification and Compliance Protocol (Required)
- Decision-first means decisions must be locked before coding; it does not forbid pre-decision analysis such as required thought experiments.
- The agent must collect and lock user decisions before making any code edits for a task.
- Locked decisions must be recorded as Decision Intent Log entries in the relevant `TODO/TODO-<handle>-<slug>.md` file(s) with clear intent and rationale.
- The agent must ask decision questions up front in a single intake round whenever possible.
- Required decision categories are architecture, design/behavior, implementation approach, function naming, variable naming, and file/path decisions.
- The agent must ask these as multiple-choice questions whenever practical.
- When a thought experiment (TE) is required, the agent must complete the TE before asking final DF questions. TEs narrow alternatives; DF questions and answers lock the decision before implementation.
- Thought experiments (TEs) are analysis artifacts; Decision Intent (DI) entries are the separate records that capture the locked decision after DF is resolved.

## Thought Experiment Protocol (Required)
- Before locking any non-trivial decision that will require DF questions and answers, the agent must run a thought experiment (TE) if multiple plausible designs remain.
- A TE happens before final DF questions. Its purpose is to narrow the design space so DF questions and answers are informed by explicit scenario analysis.
- The agent must not collapse a TE into a short opinion or recommendation. The agent must explicitly model concrete scenarios and consequences.
- Each new TE must have a unique proquint handle in the format `TE-<handle>`, where `<handle>` is minted by `tools/mint-handle` from the global TODO/TE/DR/DI handle namespace. Existing pre-upgrade TE handles and prior aliases remain historical records.
- The TE doc filename must be `TE-<handle>-<slug>.md` and live under `docs/thought-experiments/`, for example: `docs/thought-experiments/TE-mumuv-naming-reconciliation.md`. The slug is informational and may be edited; the proquint handle is permanent.

### TE Intake Requirements
- Before locking decisions or asking final DF questions, the agent must identify:
  - the decision being tested,
  - the candidate alternatives,
  - the assumptions and threat/trust model,
  - the scope and systems affected.
- If the TE relates to an existing TODO, the agent must reference the TODO handle and subtask handle (for example, `fonuz.1`).

### TE Execution Requirements
- Each TE must evaluate the same decision across multiple concrete scenarios.
- Scenarios must include, when relevant:
  - normal operation,
  - failure/corruption/incomplete writes,
  - concurrent actors or mixed-version nodes,
  - long-horizon evolution and migration,
  - trust-boundary changes,
  - scale effects (storage, bandwidth, CPU, operational complexity).
- The agent must compare alternatives under the same assumptions instead of switching assumptions mid-analysis.
- The agent must state what each alternative makes easier, what it makes harder, and what new obligations it creates.

### TE Authoring Conventions
- Named actors follow the cryptography-literature alphabetical convention. Use Alice, Bob, Carol, Dave, Ellen, Frank, and so on for cooperative actors; use Mallory for adversaries; name Steve explicitly only when his repo-owner role is load-bearing in the scenario. Use this convention in TEs, scenario analyses, tabletop simulations, DR/DI prose, and worked examples in specs or docs when named actors are useful. Do not invent ad-hoc names when the convention fits. (DI-034-20260508-060134)

### TE Output to DF
- After the TE, the agent must identify:
  - rejected alternatives,
  - surviving alternatives,
  - unresolved questions that still require user choice,
  - any new naming/path/runtime decisions exposed by the TE.
- Final DF questions must be framed from the surviving alternatives identified by the TE. The agent must not ask broad DF questions that ignore TE results.

### TE Artifacts
- For each completed TE, the agent must write a verbatim copy of the thought experiment into a standalone file under `docs/thought-experiments/`.
- The doc filename must begin with the TE ID and then use a descriptive suffix.
- The doc must stand on its own and include:
  - title,
  - TE ID,
  - decision under test,
  - assumptions,
  - alternatives,
  - scenario analysis,
  - conclusions,
  - implications for the repo's open TODOs and pending DIs.

### TE Decision Rules
- A TE does not by itself lock a decision.
- After the TE, the agent must either:
  - ask the user to choose among the surviving alternatives, or
  - recommend one surviving alternative and clearly state why the others were rejected.
- After user choice is resolved, the agent must record the locked result via the existing DI process before implementation.
- If a TE exposes a new ambiguity, dependency, or naming/path decision, the agent must stop and resolve that before implementation.

### TE Final Handoff Requirements
- In the final response for TE work, the agent must include:
  - which TE was completed,
  - the TE ID,
  - the doc path under `docs/thought-experiments/`,
  - the surviving alternatives,
  - the recommended conclusion or the exact DF question that remains for the user.
- Hard gate: for decisions that require a TE, work is incomplete until the TE doc exists and the resulting decision status is explicit (`needs DF`, `locked`, or `deferred`).

### TE Editing Policy (Required)

Once a TE is filed in `docs/thought-experiments/`, edits to it follow a categorized policy. The policy is locked in `DI-020-20260502-213103` (categorized editing regimes), `DI-020-20260502-213104` (uniform applicability across all TE corpora; this rule applies wherever TEs are stored, not just under `docs/thought-experiments/`), and `DI-020-20260502-213105` (holistic reading by default for substantive questions, single-TE reading allowed for obviously mechanical ones). The Cat-1 clause of `DI-020-20260502-213103` was superseded on 2026-05-02 by `DI-020-20260502-232651` (Cat-1a / Cat-1b split). Four further Cat-3 navigational refinements are appended to TE-dabol (`docs/thought-experiments/TE-dabol-te-editing-policy-and-holistic-corpus.md`) and to TE-vudaf (`docs/thought-experiments/TE-vudaf-editing-policy-tabletop.md`); the agent must read both TEs and their `## Refinements` sections before performing any TE edit.

The seven categories are:

- **Cat-1a (current-pointer paths).** A path reference that names the current location of a file. Mechanical sweep in place; no top-of-file note required.
- **Cat-1b (historical-quotation paths).** A path reference that quotes an earlier corpus state — inside a markdown blockquote, attributed to another TE ("TE-N states ..."), in past tense ("TE-magup used the path ..."), inside a `## Refinements` section, supersedence note, or `Decision status` line. Left untouched; rewriting would falsify the historical record. Per-match classification with five heuristics (quotation context; Refinements / supersedence framing; past tense; default Cat-1a; when-in-doubt-Cat-1b). Sweep tools may emit matches with surrounding context for human review but must not auto-rewrite.
- **Cat-2 (vocabulary updates).** A rename of a term whose meaning is unchanged (typo fixes, terminology consolidation). In place, with a top-of-file note pointing at the driving TE or TODO. The note must enumerate by ID every DI that lives in the affected TE, paired with an explicit promise that the rewrite preserves each DI's meaning. A TE without DIs gets a one-line `no DIs in this file` note. Form: `Cat-2 vocabulary update per <driving TE or TODO>: '<old term>' -> '<new term>'. The following DIs in this file are unchanged in meaning: DI-XXX-..., DI-YYY-..., DI-ZZZ-... .` Mandatory pre-step: grep the entire corpus for the old term inside quotation contexts (markdown blockquotes; fenced code blocks presented as citations; single/double-quoted phrases attributed to another TE via `TE-N states`, `TE-N reads`, `originally said`, `as of TE-N`, `the corpus showed`); each match is classified Cat-2 (sweep) or Cat-2-historical (leave) per the same heuristics as Cat-1a/Cat-1b.
- **Cat-3 (navigational forward pointers).** Append a dated entry to the TE's `## Refinements` section (created if absent, placed after `## Decision status`) describing where the affected reader should now look. The TE body above is unchanged. No DI is filed for a Cat-3 entry. Procedural tightenings of an existing category's how-to are Cat-3.
- **Cat-4 (resolved-implication forward pointers).** Same shape as Cat-3, used when an item from the TE's `Implications and future work` list has resolved (a TODO filed; a DR opened; a downstream TE landed). Append-only; no body edit.
- **Cat-5 / Cat-6 / Cat-7 (substantive supersedence).** A material change to a locked DI's meaning, scope, or applicability requires a new TE that supersedes the affected one. The new TE carries its own DFs and DIs; the older TE's `## Decision status` is updated to `superseded by TE-<id>` and its top-of-file `## Status` field is updated to `superseded by TE-<id> / DI-<id>`. The older TE's body is otherwise untouched.

Every TE in the corpus carries a top-of-file `## Status` field placed immediately after the TE ID line. Canonical values: `needs DF`, `decided`, `decided, refined`, `superseded by TE-<id> / DI-<id>`, `withdrawn`. Legacy values preserved during retrofit: `stub`, `open`, `recommended for immediate adoption`, `locked for the <protocol>`. New TEs prefer canonical values; the field is updated by Cat-1a sweep when the TE's state changes.

The `## Refinements` section is the single append-only home for Cat-3 / Cat-4 entries on a TE. Entries are dated (`### YYYY-MM-DD — <title>`) and ordered chronologically. The body of the TE above the `## Refinements` section is treated as historical evidence: a Cat-1a path-rename or Cat-2 vocabulary sweep on the body is permitted under its category rules; a Cat-3 / Cat-4 forward-pointer is appended to `## Refinements` rather than rewriting the body; a Cat-5 / Cat-6 / Cat-7 substantive change is filed as a new superseding TE rather than as an edit. The four Cat-3 Refinements on TE-dabol are themselves examples of this shape: each one tightens a category's procedure without changing the locked policy, and is filed as a Refinement entry rather than as a new DI.

Reading default: holistic. When deciding whether an edit is mechanical or substantive, when interpreting a single TE's claims, or when reasoning about whether a refinement is Cat-3 or Cat-5–7, the agent must read the corpus holistically (the affected TE plus the corpus's editing-policy chain: TE-dabol, TE-vudaf, and any other TEs they cite or that cite them). Single-TE reading is reserved for obviously mechanical questions (a single typo; a path that has demonstrably moved; a Status field retrofit) and only after the holistic read has confirmed the question is mechanical. When in doubt, read holistically.

Applicability: this policy applies uniformly to every TE corpus in this repository, regardless of which protocol or harness directory it lives in. Per-protocol corpora may add stricter rules but may not relax these rules.

### Naming Decisions (Required)
- The agent must not invent function names or variable names that are not already covered by locked naming decisions.
- If naming is not covered, the agent must stop and ask multiple-choice naming options before continuing.

### File/Path Decisions (Required)
- Path approvals are mandatory for all touched paths:
  - repo-changed files (create/rename/move/delete),
  - runtime touched paths (read/write/delete), including input files, output files, DB files, caches, fixtures, and temporary test files.
- The agent must ask path approvals one path at a time via multiple-choice questions.
- Path-question order must be dependency order.
- Each path question must include: action, exact path (or approved dynamic pattern ID), purpose, class (`prod-code | prod-data | test | temp`), and lifecycle intent.
- Temporary test paths require explicit approval and an explicit cleanup plan before handoff.
- Dynamic/runtime-generated paths must be approved by pattern, with:
  - allowed root bounds,
  - allowed actions,
  - concrete examples.
- The agent must ask one multiple-choice approval per dynamic path pattern.
- If any unapproved runtime path appears, the agent must stop and ask before continuing.

### Decision Lock and Stop Rule
- The agent must produce a Decision Lock summary with decision IDs before code edits begin.
- The agent must not proceed if any required decision is missing, ambiguous, or conflicting.
- The agent must stop and ask immediately if a new decision need appears during implementation.
- The agent must not assume defaults for locked categories unless the user explicitly approves defaults.

### Compliance Ownership (Agent)
- The agent must treat user decisions as authoritative and implement to those decisions.
- The agent must run a compliance self-review before finalizing and must fix all non-compliance before handoff.
- Hard gate: work is incomplete until compliance is PASS, or the user explicitly approves an exception.
- The user should not need to manually inspect diffs to determine compliance.

### Required final handoff artifacts
- `Decision Compliance: PASS/FAIL`
- Decision Matrix mapping each locked decision ID to implementation evidence.
- Inline diff annotations in the form `path:line -> decision_id -> rationale`.
- Runtime Path Touch Matrix listing each approved runtime path/pattern, action used, and where it is implemented/validated.
- `Exceptions:` listing only user-approved deviations.
- Every non-trivial behavior change must include intent provenance per existing DI requirements.

## Coding Style & Naming Conventions
- Use object-oriented design with structs and methods; avoid large functions and global state.
- Follow generally accepted object oriented design patterns.
- Keep Go code `gofmt`-clean; package names should be short and lower-case.
- Prefer focused edits over broad refactors unless required.
- Add and maintain explanatory comments for non-obvious logic.
- Use `git mv` for file moves/renames to preserve history.

## Error Handling Policy (Required)
- Never use `|| true` in scripts, templates, or make recipes. Always inspect
  command exit codes explicitly with `if/else` branches and handle each outcome.
- For non-fatal cleanup/diagnostics steps, record command status (exit code and
  logs) explicitly; do not fail silently.
- In Go code, never ignore errors with `_ = ...`; handle, propagate, or report
  errors explicitly.
- Run `errcheck ./...` and keep it passing for Go changes.

## DR/DI Source-of-Truth Protocol (Required)
- In this repo, DR and DI logs are the primary source of truth for decisions and open questions.
- Documents and code are outputs of that process and must link back to DR/DI records.
- Person identity in DR/DI records must use full email with label format: `user@example.com (FirstName)`.
- In DRs, `Asked by` and person-valued `Waiting on` fields must use that format.
- In DIs, `Author` must use that format.
- A DI `Author` is the decision-maker, not merely the recorder. If an agent records a DI for Steve's chat decision, Steve is the `Author`; an agent may be `Author` only when Steve explicitly delegates that decision authority to that agent. (DI-034-20260508-060134)
- A settled statement in docs (or critical logic in code comments) must cite at least one DI ID.
- An unresolved question or uncertainty must cite at least one DR ID.
- If an unresolved question has no DR yet, create a DR before finalizing the change.
- During TODO-kugod migration, apply these rules incrementally as sections/files are brought under DR/DI tracking.
- Intent: New coordination artifacts use a single proquint handle namespace so TODO, TE, DR, and DI references do not depend on timestamp or integer allocation. Source: DI-nisam


## Comment Preservation Protocol (Required)
- Never remove existing code comments unless they are replaced in the same patch by equal-or-better explanatory comments near the same logic.
- When rewriting or refactoring code, port old explanatory intent first, then improve wording.
- If a touched non-trivial code block has no comments, add explanatory comments.
- Do not treat shorter comments as better unless they preserve all important intent.
- For any non-trivial behavior change, include a behavior-level comment with:
  - `Intent:` a short, clear rationale (a sentence or a few; no hard cap if more is needed for clarity).
  - `Source:` a DI ID in the format `DI-<handle>`.
  - `<handle>` is minted by `tools/mint-handle` and is globally unique across TODO, TE, DR, and DI owners.
  - Optional: TODO file/section reference for faster lookup.
- If a comment must be dropped with no replacement, stop and ask the user before proceeding.
- Before editing a file, review existing comments in that file.
- Maintain a `## Decision Intent Log` at the top of relevant `TODO/TODO-<handle>-<slug>.md` files.
- Treat DI logs as append-only history. Do not rewrite or delete prior entries.
- When intent evolves, add a new DI entry and set `Supersedes: <old-di-id>`.
- DI entries must include:
  - `ID: DI-<handle>`
  - `Date: YYYY-MM-DD HH:MM:SS`
  - `Status: active|superseded`
  - `Decision:`
  - `Intent:`
  - `Constraints:`
  - `Affects:`
  - `Supersedes:` (optional)
- After editing, run a comment-delta audit on each touched code file using: `git diff -U0 -- <file> | rg -n '^-\\s*//|^-\\s*/\\*|^\\+\\s*//|^\\+\\s*/\\*'`.
- Resolve all removed-comment lines before finalizing unless explicit user approval was given.
- In the final response, include:
  - `Comment audit: PASS/FAIL`, with file list.
  - `Intent provenance audit: PASS/FAIL`, listing files with behavior changes and DI sources.
- Hard gate: behavior-changing work is incomplete unless comments preserve intent and include DI provenance.
- Do not remove comments or documentation; update them if outdated or incorrect.

### Comment + DI Examples
- Comment format example:
  - `// Intent: Keep context resolution stable across workspace scans to avoid target drift between plan and run. Source: DI-vapoj`
- Decision Intent Log entry template (for TODO files):
  - `ID: DI-<handle>`
  - `Date: YYYY-MM-DD HH:MM:SS`
  - `Status: active`
  - `Decision: <what was decided>`
  - `Intent: <short clear rationale>`
  - `Constraints: <hard limits, dependencies, assumptions>`
  - `Affects: <paths, modules, commands, docs>`
  - `Supersedes: <old DI ID, optional>`

# DR Records

The DR/ directory stores Decision Request (DR) records for coordination work.

Rules:
- One DR per file.
- DR files are append-only event logs.
- Keep TODO files as snapshots; link TODOs to DR files for open questions.
- Person identity format: `user@example.com (FirstName)`.

Recommended file naming:
- `DR-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`.

Required DR fields:
- `DR-ID`
- `Date`
- `Asked by` (person identity format above)
- `State` (`open | decided | blocked | implemented | closed`)
- `Question`
- `Why this blocks progress`
- `Affects` (repos/files/components)
- `Unblocks` (TODO IDs/tasks)
- `Waiting on` (person identity format above, or DI ID)
- `Decision` (filled when decided)
- `Linked DI`
- `Related commits`
- `Last updated`

Reference pattern:
- From TODO files: `../DR/<filename>.md`


## Testing Guidelines
- Use Go's standard `testing` package with deterministic tests.
- Avoid network calls in tests unless explicitly required and documented.
- When changing `plan/run` behavior, add coverage for both command paths when possible.

## Commit & Pull Request Guidelines
- Treat a line containing only `commit` as: add and commit all changes with an AGENTS-compliant message.
- Use short, imperative, capitalized commit subjects.
- Summarize changes per file in commit bodies.
- Stage files explicitly (avoid `git add .` / `git add -A`).
- Do not force-push repo branches unless a scoped DI/DR explicitly authorizes the exception. Role overlays may add stricter no-force-push rules for their branches or may document a narrow private-remote exception, but the repo-wide default is no history rewrites. (DI-034-20260508-060134)
- Do not commit local state files, generated binaries, credentials, tokens, signing keys, or other secrets. (DI-034-20260508-060134)

## Glossary
- **TE**: Thought Experiment. Analysis doc under `docs/thought-experiments/TE-<handle>-<slug>.md` or an approved per-protocol TE corpus. The handle is a proquint minted by `tools/mint-handle`. (DI-034-20260508-060134)
- **DR**: Decision Request. Open question or decision-tracking record under `DR/DR-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`. (DI-nisam; DI-034-20260508-060134)
- **DI**: Decision Intent. Locked decision record inside a `## Decision Intent Log` in a protocol TODO file. New DI ID format is `DI-<handle>`, where `<handle>` is minted by `tools/mint-handle`. (DI-nisam; DI-034-20260508-060134)
- **DF**: Decision Framing. The multiple-choice intake round used to lock a decision after any required TE narrows the alternatives. (DI-034-20260508-060134)
- **TODO**: Task-tracking file under `TODO/TODO-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`. The master cross-listed index is `TODO/TODO.md`; per-protocol queues live at `TODO/TODO.md`. (DI-nisam; DI-034-20260508-060134)
- **twig**: Short kebab-case task name used in branch names, usually as `<user>/<twig>` such as `ppx/<twig>` or `stevegt/<twig>`. (DI-034-20260508-060134)
- **pCID**: Protocol CID. A pCID is the content hash of a spec document that defines a wire protocol; it is analogous to a TCP/UDP port number with no central registry because the spec hash is the port number. A pCID is not the hash of a particular message, payload, or promise body. (DI-009-20260429-173359; DI-034-20260508-060134)
