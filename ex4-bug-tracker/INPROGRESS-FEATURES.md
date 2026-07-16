# INPROGRESS FEATURES

Current status for `ex4-bug-tracker`.

## Current slice

Status: implemented slice complete

- browser queue view
- browser issue detail view
- new issue flow
- comment flow
- assignment flow
- status change flow
- reopen flow
- attachment upload and download
- engineer CLI
- append-only runtime storage
- seeded demo launch path

Goal:

- make the first bug-tracker slice usable on its own
- keep the workflow durable, inspectable, and easy to extend

## Confirmed

1. Human-friendly issue IDs
2. Fixed workflow statuses
3. Reporter / triage / engineer identities
4. Single active assignee
5. Reopen from `Resolved` to `Triaged`
6. Queue filtering by status
7. Queue filtering by assignee
8. Merged issue timeline
9. Real uploaded file attachments
10. CLI list assigned issues
11. CLI show issue detail
12. CLI comment on issue
13. CLI move issue to `In Progress`
14. CLI move issue to `Resolved`
15. Hidden built-in team field for later multi-team work
16. Seeded demo runtime with representative issue states

## Likely next steps

- add richer queue sorting and filtering
- add more than one engineer identity
- show a read-only team badge while keeping team fixed
- add triage CLI commands
- improve timeline formatting and summaries

## Not in this slice

- websocket or live collaboration transport
- rich text or editor-centric behavior
- user signup or full auth system
- notifications
- labels, watchers, or custom workflows
- visible multi-team routing
