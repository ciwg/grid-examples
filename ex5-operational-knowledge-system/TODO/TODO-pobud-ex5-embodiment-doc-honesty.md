# TODO pobud - ex5 embodiment doc honesty

## Decision Intent Log

ID: DI-pobud
Date: 2026-07-21 16:13:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the browser/CLI embodiment-equality doc overstatement as its own ex5 review TODO.
Intent: Describe the current ex5 embodiment depth honestly instead of claiming browser and CLI parity that the implementation does not provide.
Constraints: Stay inside `ex5`; this is primarily a docs-truthfulness pass unless the user later decides to expand the CLI surface instead.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-pobud-ex5-embodiment-doc-honesty.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Bring the ex5 docs into line with the actual embodiment depth of the browser,
CLI, and Neovim surfaces.

## Review finding

The docs still describe browser and CLI as equal first-class embodiments even
though the CLI is materially narrower: it lacks evidence upload, typed-link
creation, live-draft access, structured/problem search filters, responsibility
detail browsing, and contextual drilldowns.

## Tasks

- [x] pobud.1 Audit ex5 docs for embodiment-equality language that overstates CLI capability.
- [x] pobud.2 Rewrite those sections to describe the shared runtime truthfully without implying feature parity.
- [x] pobud.3 Keep any remaining “equal” language only where it is genuinely true.

## Status

- done
- derived from the 2026-07-21 deep ex5 review
