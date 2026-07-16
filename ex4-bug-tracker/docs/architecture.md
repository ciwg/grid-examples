# bug-tracker architecture

`ex4-bug-tracker` keeps the shared workflow contract small and durable.

## Topology

```text
Browser UI ----------------------\
                                  \
CLI -------------------------------> local HTTP server -> append-only event log
                                  /
Attachment upload/download -------/
```

The server owns:

- issue ID allocation
- built-in identity and role validation
- append-only issue event persistence
- current issue projection
- attachment copy and download behavior
- the queue and detail HTTP surface

The browser UI owns:

- queue and detail presentation
- local form state
- current issue selection
- the identity picker

The CLI owns:

- engineer-focused issue commands
- text output for assigned work
- the same HTTP requests the browser relies on

This keeps the durable workflow state in one place while allowing multiple
embodiments over the same shared model. Source: `DI-dajak`; `DI-nunit`;
`DI-ninuf`.

## Shared data model

Every issue carries:

- `id`
- `title`
- `description`
- `severity`
- `status`
- `reporter`
- `assignee`
- `team`
- `created_at`
- `updated_at`

Each issue also has a timeline of append-only events. V1 event types are:

- `created`
- `commented`
- `assigned`
- `status_changed`
- `attachment_added`

The server projects current queue/detail state from those events instead of
mutating a canonical issue row in place. Source: `DI-nunit`.

## Workflow model

V1 uses a fixed status flow:

- `New`
- `Triaged`
- `In Progress`
- `Resolved`

Allowed transitions are intentionally narrow:

- triage: `New -> Triaged`
- engineer: `Triaged -> In Progress`
- engineer: `In Progress -> Resolved`
- reporter or triage: `Resolved -> Triaged`

Reopen clears the active assignee while preserving the full prior history.
Source: `DI-ninuf`; `DI-gofub`.

## Identity and team model

Built-in identities are:

- `reporter`
- `triage`
- `engineer`

The first slice uses those fixed roles instead of a broader auth or user
management system. Every issue also stores `team=CORE`, but the team field is
hidden in the current UI and CLI. That gives the storage model a seam for later
multi-team work without complicating the first usable release. Source:
`DI-ninuf`; `DI-gofub`.

## Attachment model

Attachments are not external references. The server copies uploaded files into
the local runtime root under `.bug-tracker/attachments/`, then records an
attachment event in the issue timeline. Downloads resolve back through the
server, not directly to original host paths. Source: `DI-nunit`.

## Runtime layout

The default runtime root is `.bug-tracker/` and currently contains:

- `events.jsonl`
  - append-only issue event log
- `attachments/`
  - copied uploaded files organized under per-issue paths

This keeps the first example inspectable on disk and easy to reset locally.
Source: `DI-dajak`; `DI-nunit`.
