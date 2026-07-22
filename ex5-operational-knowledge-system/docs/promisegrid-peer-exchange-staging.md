# ex5 PromiseGrid peer-exchange staging

This note records the current shipped relay-visible peer exchange slice in
`ex5`. Source: `DI-guzab`; `DI-voruk`; `DI-vamok`; `DI-faruv`; `DI-ruzok`;
`DI-rumek`; `DI-loruk`; `DI-pivul`.

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
exchange exports and imports the whole family logs plus their compatibility
events over the current local HTTP adapter. Evidence blobs are included inline
by CID so another host can actually resolve imported evidence. Source:
`DI-guzab`; `DI-voruk`; `DI-vamok`; `DI-faruv`; `DI-ruzok`; `DI-rumek`;
`DI-pivul`.

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

The first staged relay-visible work should:

- expose signed envelope export/import boundaries for the current portable
  families
- keep current local family logs and compatibility event logs during the first
  exchange pass
- carry evidence blobs inline by CID during the bootstrap phase
- avoid tightening browser, CLI, or Neovim beyond the current local HTTP
  adapter in the same slice

Source: `DI-guzab`.

## Follow-on backlog

- stronger CAS-backed read/replay authority for exchanged artifacts outside the
  frozen families

Source: `DI-guzab`; `DI-tivor`; `DI-rumek`; `DI-pivul`.
