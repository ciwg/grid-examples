# bug-tracker architecture

`ex4-bug-tracker` keeps the shared workflow contract small and durable.

## Topology

```text
Browser UI ----------------------\
                                  \
CLI -------------------------------> local HTTP server -> event projection + CAS
                                  /
Attachment object carriage --------/
```

The server owns:

- issue ID allocation
- service-scoped public enrollment and proof validation
- canonical pCID-selected envelope construction and verification
- accepted-envelope CAS and separate rejected-ingress observations
- append-only issue event persistence
- current issue projection
- attachment-object carriage and signed reference projection
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

## PromiseGrid scope boundary

The browser and CLI are independent signing embodiments behind local adapters
of one HTTP service. Each owns its private key; the service stores only an
explicit, service-scoped public enrollment binding. The service prepares exact
canonical signing bytes, verifies returned proofs, and accepts only final
`grid([42(pCID), payload, proof])` envelopes. `events.jsonl` remains local
projection history—not independently shared evidence or a general identity,
delegation, role-continuity, or authorization system. Source: `DI-muzal`;
`DI-kolaf`; `DI-ninul`.

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

Attachment bytes first enter local CAS as opaque objects. They do not change an
issue by themselves. A signed `issue-attachment-reference` promise names the
CAS CID and gives the object issue-specific meaning; only then does the local
projection append `attachment_added`. Downloads resolve through that accepted
reference. Source: `DI-ninul`; `DI-kolaf`.

## Runtime layout

The default runtime root is `.bug-tracker/` and currently contains:

- `events.jsonl`
  - append-only issue event log
- `cas/<CID>`
  - accepted envelopes and attachment objects
- `accepted-promises.jsonl`, `observations.jsonl`, `agent-bindings.jsonl`
  - local acceptance, rejection, and public enrollment records

This keeps the first example inspectable on disk and easy to reset locally.
Source: `DI-dajak`; `DI-nunit`.
