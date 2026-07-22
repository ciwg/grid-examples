# ex5 PromiseGrid peer-exchange staging

This note records the current shipped relay-visible peer exchange slice in
`ex5`. Source: `DI-guzab`; `DI-voruk`; `DI-vamok`; `DI-faruv`.

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
bootstrap exchange exports and imports the whole family logs plus their
compatibility events over the current local HTTP adapter. Evidence blobs are
included inline by CID so another host can actually resolve imported evidence.
Source: `DI-guzab`; `DI-voruk`; `DI-vamok`; `DI-faruv`.

## Still outside the current peer-visible slice

The runtime still keeps some references outside the current peer-visible set:

- places
- resources

Run records and typed links can preserve those references and report them as
unresolved during bootstrap import, but they are not yet first-class exchanged
families. Source: `DI-guzab`; `DI-vamok`; `DI-faruv`.

## Bootstrap import behavior

The shipped importer is bootstrap-only:

- it imports only into an empty runtime
- it preserves whole family history for items, approvals, evidence, runs,
  links, and responsibilities
- it reports unresolved place/resource run context and unresolved place/resource
  link endpoints explicitly instead of trimming those artifacts away

This keeps the first exchange slice PromiseGrid-complete at the family level
without pretending `ex5` already has a safe multi-peer merge contract. Source:
`DI-voruk`.

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

- non-bootstrap peer exchange into non-empty runtimes
- stronger CAS-backed read/replay authority for exchanged artifacts
- peer-visible place/resource families or another durable answer for those
  references

Source: `DI-guzab`; `DI-tivor`.
