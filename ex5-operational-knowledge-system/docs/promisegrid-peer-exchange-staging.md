# ex5 PromiseGrid peer-exchange staging

This note records what `ex5` will treat as the first relay-visible peer
exchange slice after the five local signed families were frozen. Source:
`DI-guzab`.

## First relay-visible slice

The first peer-visible slice is intentionally smaller than the current local
runtime surface. It carries only the attachment-free signed families:

- `knowledge-item`
- `knowledge-approval`
- `knowledge-link`
- `knowledge-responsibility`

Those four families are already signed, replay-verified, and portable without
depending on local attachment paths. Source: `DI-guzab`.

## Deferred from the first slice

`knowledge-evidence` is not in the first relay-visible slice.

The current evidence family signs durable metadata plus attachment references,
but attachment bytes still live on the local copied-file path. Until `ex5`
settles CAS-backed storage or another portable blob-carriage rule, peer-visible
evidence exchange would overstate what another host can actually resolve.
Source: `DI-ribof`; `DI-guzab`.

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
