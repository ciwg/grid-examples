# ex5 PromiseGrid peer-exchange staging

This note records the first shipped relay-visible peer exchange slice in
`ex5`. Source: `DI-guzab`; `DI-voruk`.

## First relay-visible slice

The first peer-visible slice is intentionally smaller than the current local
runtime surface. It carries only the attachment-free signed families:

- `knowledge-item`
- `knowledge-approval`
- `knowledge-link`
- `knowledge-responsibility`

Those four families are already signed, replay-verified, and portable without
depending on local attachment paths. The shipped bootstrap exchange exports and
imports the whole family logs plus their compatibility events over the current
local HTTP adapter. Source: `DI-guzab`; `DI-voruk`.

## Deferred from the first slice

`knowledge-evidence` is not in the first relay-visible slice.

The current evidence family signs durable metadata plus attachment references,
but attachment bytes still live on the local copied-file path. Until `ex5`
settles CAS-backed storage or another portable blob-carriage rule, peer-visible
evidence exchange would overstate what another host can actually resolve.
Source: `DI-ribof`; `DI-guzab`.

## Bootstrap import behavior

The shipped importer is bootstrap-only:

- it imports only into an empty runtime
- it preserves whole approval and link family history
- it reports unresolved run approvals and unresolved place/resource/run link
  endpoints explicitly instead of trimming those artifacts away

This keeps the first exchange slice PromiseGrid-complete at the family level
without pretending `ex5` already has a safe multi-peer merge contract. Source:
`DI-voruk`.

## Staged runtime/storage shape

The first staged relay-visible work should:

- expose signed envelope export/import boundaries for the four portable
  families
- keep current local family logs and compatibility event logs during the first
  exchange pass
- leave attachment carriage and evidence portability to the later storage
  decision
- avoid tightening browser, CLI, or Neovim beyond the current local HTTP
  adapter in the same slice

Source: `DI-guzab`.

## Follow-on backlog

- `101`: CAS-backed storage and portable evidence/blob carriage
- `102`: embodiment contract tightening after the peer/storage layers settle

Source: `DI-guzab`.
