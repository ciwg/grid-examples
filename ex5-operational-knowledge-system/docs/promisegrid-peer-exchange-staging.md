# ex5 PromiseGrid peer-exchange staging

This note records the current shipped relay-visible peer exchange slice in
`ex5`. Source: `DI-guzab`; `DI-voruk`; `DI-vamok`; `DI-faruv`; `DI-ruzok`;
`DI-rumek`.

## First relay-visible slice

The first peer-visible slice is intentionally smaller than the current local
runtime surface. It now carries these signed families:

- `knowledge-item`
- `knowledge-approval`
- `knowledge-evidence`
- `operational-run`
- `knowledge-link`
- `knowledge-responsibility`

Those six families are already signed and replay-verified. The shipped
exchange exports and imports the whole family logs plus their compatibility
events over the current local HTTP adapter. Evidence blobs are included inline
by CID so another host can actually resolve imported evidence. Source:
`DI-guzab`; `DI-voruk`; `DI-vamok`; `DI-faruv`; `DI-ruzok`; `DI-rumek`.

## Still outside the current peer-visible slice

The runtime still keeps some references outside the current peer-visible set:

- places
- resources

Run records and typed links can preserve those references and report them as
unresolved during bootstrap import, but they are not yet first-class exchanged
families. Source: `DI-guzab`; `DI-vamok`; `DI-faruv`.

## Current import behavior

The shipped importer now accepts whole-family exchange into non-empty runtimes:

- it dedupes delivered history by `(origin_peer_id, origin_sequence)`
- it preserves whole-family signed history for items, approvals, evidence,
  runs, links, and responsibilities
- it assigns a fresh local compatibility `Sequence` when new peer history is
  accepted
- it reports unresolved place/resource run context and unresolved place/resource
  link endpoints explicitly instead of trimming those artifacts away
- it still rejects create-event ID collisions such as two independent peers
  both minting the same local-facing `RECV-*`, `RUN-*`, or `RESP-*` ID

This keeps the shipped exchange honest: it is no longer bootstrap-only, but it
also does not pretend the runtime has already solved cross-peer entity
namespace reconciliation. Source: `DI-ruzok`; `DI-rumek`.

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

- reconcile peer-visible entity namespaces across independent peers
- stronger CAS-backed read/replay authority for exchanged artifacts outside the
  frozen families
- peer-visible place/resource families or another durable answer for those
  references

Source: `DI-guzab`; `DI-tivor`; `DI-rumek`.
