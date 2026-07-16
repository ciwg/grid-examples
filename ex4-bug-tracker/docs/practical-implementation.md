# Practical implementation notes

This is a practical guide for how `ex4-bug-tracker` is currently shaped and
what the next changes should preserve.

## Keep the app durable-first

The useful thing about this example is not live collaboration. It is the
durable issue trail.

- issue creation is an event
- assignment is an event
- status changes are events
- comments are events
- attachments are events

The queue and detail views are projections over that event history. Preserve
that property when extending the app. Source: `DI-nunit`.

## Keep the browser shell thin

The browser should stay a straightforward working surface over the local server:

- load metadata and identity options
- list issues
- fetch one issue
- post form-driven actions
- render the merged timeline

Do not let front-end-only local state become the source of truth for issue
workflow. Source: `DI-dajak`; `DI-ninuf`.

## Keep the CLI narrow

The CLI is intentionally smaller than the browser:

- list assigned issues
- show one issue
- comment
- start work
- resolve work

That keeps the second embodiment useful without having to reproduce the full
administrative surface. If the CLI grows later, keep its scope legible instead
of chasing parity for its own sake. Source: `DI-ninuf`.

## Status changes should stay explicit

The fixed status model is a feature of this first slice, not a missing piece.

- `New`
- `Triaged`
- `In Progress`
- `Resolved`

Reopen is also explicit: it moves the issue back to `Triaged` and clears the
assignee. Do not silently infer reopen from a comment or an assignment change.
Source: `DI-ninuf`; `DI-gofub`.

## Team is intentionally hidden

The current implementation stores `team=CORE` on every issue, but does not
show team controls in the browser or CLI.

That is deliberate:

- product behavior is single-team today
- storage is multi-team-ready
- the UI does not yet expose a concept the product is not using

If later work adds visible team routing, it should build on that stored field
instead of replacing it. Source: `DI-gofub`.

## Attachments should remain app-managed copies

Attachment behavior should keep these properties:

- the server accepts an uploaded file
- the file is copied into `.bug-tracker/attachments/`
- the issue timeline records the upload
- downloads flow back through the server

That avoids hidden dependence on outside host paths and keeps the issue history
portable within the local runtime root. Source: `DI-nunit`.

## Good near-term extension directions

Reasonable next steps that preserve the current shape:

- richer seeded demo cases and walkthrough notes
- better queue filters and sorting
- richer engineer and triage CLI commands
- per-issue severity or assignee badges in more places
- visible but still fixed team badges
- a second engineer identity for clearer assignment examples

Less good next steps for this example:

- websocket/live typing work
- rich text editor behavior
- a large user-management/auth system
- generic plugin architecture before the core workflow hardens
