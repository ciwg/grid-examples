# ex5 PromiseGrid peer-exchange staging

This note records the current shipped relay-visible peer exchange slice in
`ex5`. Source: `DI-guzab`; `DI-voruk`; `DI-vamok`; `DI-faruv`; `DI-ruzok`;
`DI-rumek`; `DI-loruk`; `DI-pivul`; `DI-pazek`.

## First relay-visible slice

The first peer-visible slice is intentionally smaller than the current local
runtime surface. It now carries these signed families:

- `knowledge-item`
- `knowledge-approval`
- `knowledge-evidence`
- `operational-run`
- `operational-place`
- `operational-resource`
- `knowledge-link`
- `knowledge-responsibility`

Those eight families are already signed and replay-verified. The shipped
exchange now has two layers over the current local HTTP adapter:

- whole-family bundle export/import for compatibility and bootstrap
- incremental relay-feed export/import over origin-aware signed records

The bootstrap bundle still carries inline evidence blobs by CID, but the
incremental relay feed names required evidence blobs only by CID and leaves the
raw blob transfer to separate relay blob routes. Source: `DI-guzab`;
`DI-voruk`; `DI-vamok`; `DI-faruv`; `DI-ruzok`; `DI-rumek`; `DI-pivul`;
`DI-pazek`.

## Still outside the current peer-visible slice

Place and resource references are no longer outside the current peer-visible
set. The remaining follow-on work is narrower: the still-unfrozen
runtime/projection state beyond the eight signed families remains local
compatibility state for now. Source: `DI-pivul`; `DI-tivor`.

## Current import behavior

The shipped importer now accepts whole-family exchange into non-empty runtimes:

- it dedupes delivered history by `(origin_peer_id, origin_sequence)`
- it preserves whole-family signed history for items, approvals, evidence,
  runs, places, resources, links, and responsibilities
- it assigns a fresh local compatibility `Sequence` when new peer history is
  accepted
- peer-visible entities now use the create-envelope CID as the durable ID and
  preserve the older short ID only as an alias

This keeps the shipped exchange honest: it is no longer bootstrap-only, and it
no longer depends on local counter-minted IDs being globally unique across
peers. Source: `DI-ruzok`; `DI-rumek`; `DI-loruk`; `DI-pivul`.

## Staged runtime/storage shape

The shipped staged relay-visible work now:

- exposes signed envelope export/import boundaries for the current portable
  families
- keeps current local family logs and compatibility event logs during the
  relay pass
- keeps inline evidence blobs only in the bootstrap bundle path
- exposes incremental relay feed plus separate CID blob transfer routes
- avoids tightening browser, CLI, or Neovim beyond the current local HTTP
  adapter in the same slice

Source: `DI-guzab`; `DI-pazek`.

## Follow-on backlog

- stronger CAS-backed read/replay authority for exchanged artifacts outside the
  frozen families

Source: `DI-guzab`; `DI-tivor`; `DI-rumek`; `DI-pivul`; `DI-pazek`.
