# Ex5 Terminal Capability Matrix

This matrix is a current-state capability view for `ex5`. It is meant to act as
the quick spreadsheet-style summary of what Browser, CLI, and Neovim can do
today. Source: `DI-movar`.

For a newcomer path that uses the checked-in sample runtime and walks through
the receiving, inventory, training, and maintenance storylines step by step,
start with the [User Guide](./user-guide.md). Source: `DI-rubav`.

## Capability Matrix

| Capability | Browser | CLI | Neovim |
| --- | --- | --- | --- |
| Create responsibilities, places, resources | Yes | Yes | No |
| Create knowledge items | Yes | Yes | No |
| Edit live draft body | Yes | No | Yes |
| Refresh live draft | Yes | No | Yes |
| Snapshot revision | Yes | Yes | Yes |
| Inspect item detail | Yes | Yes | Yes |
| Inspect run detail | Yes | Yes | Yes |
| Inspect place/resource/responsibility detail | Yes | Yes | Yes |
| Typed-link browsing | Yes | Via detail output | Yes |
| Typed-link creation | No dedicated create form | Yes | No |
| Record runs | Yes | Yes | No |
| Upload evidence | Yes | Yes | No |
| Approve item | Yes | Yes | Yes |
| Approve run | Yes | Yes | Yes |
| Supersede item | Yes | Yes | Yes |
| Free-text search | Yes | Yes | Yes |
| Structured search filters | Yes | Yes | Yes |
| Pending-review queue | Yes | Yes | Yes |
| Grouped problem review | Yes | Yes | Yes |
| Context drilldowns from place/resource/responsibility | Yes | Yes | Partial via inspectors |

The matrix reflects current shipped behavior, not planned parity. CLI and
Neovim are already strong terminal surfaces, but they are intentionally staged
and still not perfectly symmetric with the browser. Source: `DI-fudok`;
`DI-ravum`; `DI-salup`; `DI-lorav`; `DI-jabup`; `DI-vogar`.

CLI transport note: the CLI now treats the direct local Unix socket as the
required terminal contract by default and fails closed when that socket is
unavailable. Use `-socket=off` only when you explicitly want the HTTP
compatibility path. Neovim still keeps HTTP as a compatibility fallback behind
its socket-first behavior. Source: `DI-zorav`.

## Terminal-First Summary

### CLI is strongest for

- one-shot create and mutate commands
- shell-only durable revision snapshots
- evidence upload
- typed-link creation
- shell-friendly review queues
- direct drilldown summaries

Source: `DI-zanub`; `DI-vuteg`; `DI-ravum`; `DI-salup`; `DI-muvok`.

### Neovim is strongest for

- live draft editing
- durable revision snapshots from the current live draft
- staying inside one editor session
- inspecting items, runs, and linked entities
- structured search and grouped problem review without leaving the editor
- pending-review browsing
- item/run approval and item supersede actions

Source: `DI-fudok`; `DI-givot`; `DI-lorav`; `DI-vamor`; `DI-bafor`;
`DI-pudor`.

### Browser is strongest for

- widest surface coverage
- shared visual review panels
- grouped hotspot review
- richer contextual navigation

Source: `DI-pogul`; `DI-zemok`; `DI-vemur`; `DI-ruvot`.

## Review Queue Note

CLI and Neovim pending-review views both depend on the shared `/api/search`
projection carrying an explicit `approvals` array for each run. Omitted
`approvals` is treated as contract drift, not as genuine unreviewed work.
Source: `DI-davur`.
