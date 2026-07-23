# Ex5 Implementation Promise Claims

This file is the canonical B-side implementation-promise surface for
`ex5-operational-knowledge-system`. It publishes the exact frozen family spec
doc-CIDs the shipped implementation claims, and it also records which shipped
components implement or delegate which PromiseGrid-facing surfaces. Source:
`DI-jovek`; `DI-nubor`.

## Frozen Family Claims

### knowledge-item

claim:           partially-implements
spec:            bafkreih5sww7apgdurye4n6el4nsdgwgbcawochqe42t3g3rrkjjehccgy
scope:           local-runtime signed envelopes for `knowledge_item_created`, `revision_added`, `knowledge_item_status_changed`, and `knowledge_item_superseded`, surfaced through the shared runtime plus the direct browser, CLI, and Neovim embodiment contracts
breaking-change: false
notes:           `knowledge-item` is a shipped frozen family. Browser now reaches it through the Chrome/Chromium native-messaging embodiment, while CLI and Neovim prefer the direct local Unix-socket contract and keep HTTP only as explicit compatibility transport.

### knowledge-approval

claim:           partially-implements
spec:            bafkreibdxszozhrp335vpi5v4hxaoqwp46q2xd22n4iwbniabhu6nfykau
scope:           local-runtime signed envelopes for `approval_recorded`, covering both knowledge-item and run approvals through the shared runtime plus the direct browser, CLI, and Neovim embodiment contracts
breaking-change: false
notes:           `knowledge-approval` is a shipped frozen family for named-role review outcomes. Lifecycle status transitions that result from approvals remain part of `knowledge-item`.

### knowledge-evidence

claim:           partially-implements
spec:            bafkreidyre7waqivwh7ef5hb35rlzogpp3lbt4sdprygoi2ii47vaxh7h4
scope:           local-runtime signed envelopes for `evidence_added`, covering structured evidence metadata plus attachment references, with copied evidence blobs dual-written into CID-addressed CAS storage
breaking-change: false
notes:           `knowledge-evidence` is a shipped frozen family. Attachment bytes remain outside the family payload itself even though the runtime now also stages copied blobs by CID.

### knowledge-link

claim:           partially-implements
spec:            bafkreihl643tk2lawdfvyuexfrd3gtx3hrksdcrze66vvqubfbjxga3xui
scope:           local-runtime signed envelopes for `link_added`, covering typed operational links across items, runs, responsibilities, places, and resources
breaking-change: false
notes:           `knowledge-link` is a shipped frozen family. The runtime validates typed links before durable publication and then projects them through browser, CLI, and Neovim inspectors.

### knowledge-responsibility

claim:           partially-implements
spec:            bafkreihtw2i5j7au7xxuetrp2hunanl6rzyaiffg3ibuboqv46jlj56jfe
scope:           local-runtime signed envelopes for `responsibility_created`, covering first-class durable responsibilities through the shared runtime plus the direct browser, CLI, and Neovim embodiment contracts
breaking-change: false
notes:           `knowledge-responsibility` is a shipped frozen family. Responsibility drilldowns, review flows, and typed-link projections remain app/runtime behavior on top of the frozen records.

### operational-run

claim:           partially-implements
spec:            bafkreicn5t2ghs2d6b3olhhzayh2cle2n67ltu5lzcvxe4sp5x67frs5du
scope:           local-runtime signed envelopes for `run_recorded`, covering durable performed operational execution records through the shared runtime plus the direct browser, CLI, and Neovim embodiment contracts
breaking-change: false
notes:           `operational-run` is a shipped frozen family. Evidence, approvals, and typed links remain separate frozen families that anchor to the run rather than fields inside the run payload itself.

### operational-place

claim:           partially-implements
spec:            bafkreic3cc7tlg4gbbktkipkdzn2gjptibrmyszl75s4rjkn4umfglgtmu
scope:           local-runtime signed envelopes for first-class operational place records through the shared runtime plus the direct browser, CLI, and Neovim embodiment contracts
breaking-change: false
notes:           `operational-place` is a shipped frozen family. Place hierarchy, run drilldowns, and grouped problem projections remain app/runtime behavior on top of the frozen records.

### operational-resource

claim:           partially-implements
spec:            bafkreihtahpvdzmtjr4ouf5oy3ixv7anfccn7mdda7yr7vz5khkv7pe3k4
scope:           local-runtime signed envelopes for first-class operational resource records through the shared runtime plus the direct browser, CLI, and Neovim embodiment contracts
breaking-change: false
notes:           `operational-resource` is a shipped frozen family. Resource hierarchy, run drilldowns, and grouped problem projections remain app/runtime behavior on top of the frozen records.

## Component Claims

### local runtime (`service/*` plus `cmd/operational-knowledge`)

claim:           partially-implements
spec:            bafkreih5sww7apgdurye4n6el4nsdgwgbcawochqe42t3g3rrkjjehccgy, bafkreibdxszozhrp335vpi5v4hxaoqwp46q2xd22n4iwbniabhu6nfykau, bafkreidyre7waqivwh7ef5hb35rlzogpp3lbt4sdprygoi2ii47vaxh7h4, bafkreihl643tk2lawdfvyuexfrd3gtx3hrksdcrze66vvqubfbjxga3xui, bafkreihtw2i5j7au7xxuetrp2hunanl6rzyaiffg3ibuboqv46jlj56jfe, bafkreicn5t2ghs2d6b3olhhzayh2cle2n67ltu5lzcvxe4sp5x67frs5du, bafkreic3cc7tlg4gbbktkipkdzn2gjptibrmyszl75s4rjkn4umfglgtmu, bafkreihtahpvdzmtjr4ouf5oy3ixv7anfccn7mdda7yr7vz5khkv7pe3k4
scope:           one shared local runtime implementing append-only event history, signed-envelope durability for the eight frozen families, CAS-backed replay/export, peer exchange, relay-feed export/import, projections, direct local embodiment contracts, and explicit HTTP compatibility transport
breaking-change: false
notes:           This is the main PromiseGrid implementation surface in ex5. Browser, CLI, Neovim, the browser native host, and the relay binary all delegate frozen-family semantics to this runtime rather than reimplementing those semantics independently.

### browser embodiment (`web/*`, `chrome-extension/*`, and `cmd/operational-browser-host`)

claim:           partially-implements
spec:            bafkreih5sww7apgdurye4n6el4nsdgwgbcawochqe42t3g3rrkjjehccgy, bafkreibdxszozhrp335vpi5v4hxaoqwp46q2xd22n4iwbniabhu6nfykau, bafkreidyre7waqivwh7ef5hb35rlzogpp3lbt4sdprygoi2ii47vaxh7h4, bafkreihl643tk2lawdfvyuexfrd3gtx3hrksdcrze66vvqubfbjxga3xui, bafkreihtw2i5j7au7xxuetrp2hunanl6rzyaiffg3ibuboqv46jlj56jfe, bafkreicn5t2ghs2d6b3olhhzayh2cle2n67ltu5lzcvxe4sp5x67frs5du, bafkreic3cc7tlg4gbbktkipkdzn2gjptibrmyszl75s4rjkn4umfglgtmu, bafkreihtahpvdzmtjr4ouf5oy3ixv7anfccn7mdda7yr7vz5khkv7pe3k4
scope:           Chrome/Chromium Manifest V3 browser embodiment over native messaging, with typed runtime operations for day-to-day read, create, operate, and live-draft flows and no silent fallback to the older browser HTTP path
breaking-change: false
notes:           This embodiment depends on the shipped extension plus native host and delegates durable family semantics to the local runtime. Unsupported browsers are outside the current implementation promise.

### CLI embodiment (`cmd/oks-cli`)

claim:           partially-implements
spec:            bafkreih5sww7apgdurye4n6el4nsdgwgbcawochqe42t3g3rrkjjehccgy, bafkreibdxszozhrp335vpi5v4hxaoqwp46q2xd22n4iwbniabhu6nfykau, bafkreidyre7waqivwh7ef5hb35rlzogpp3lbt4sdprygoi2ii47vaxh7h4, bafkreihl643tk2lawdfvyuexfrd3gtx3hrksdcrze66vvqubfbjxga3xui, bafkreihtw2i5j7au7xxuetrp2hunanl6rzyaiffg3ibuboqv46jlj56jfe, bafkreicn5t2ghs2d6b3olhhzayh2cle2n67ltu5lzcvxe4sp5x67frs5du, bafkreic3cc7tlg4gbbktkipkdzn2gjptibrmyszl75s4rjkn4umfglgtmu, bafkreihtahpvdzmtjr4ouf5oy3ixv7anfccn7mdda7yr7vz5khkv7pe3k4
scope:           direct local Unix-socket embodiment for inspect, search, review, and mutation flows over the shared runtime, with HTTP available only as explicit compatibility opt-in through `-socket=off`
breaking-change: false
notes:           The CLI does not claim independent durable-family semantics. It is a terminal embodiment over the shared runtime contract and now fails closed by default when the preferred direct socket contract is unavailable.

### Neovim embodiment (`nvim/*` and `scripts/oks-nvim`)

claim:           partially-implements
spec:            bafkreih5sww7apgdurye4n6el4nsdgwgbcawochqe42t3g3rrkjjehccgy, bafkreibdxszozhrp335vpi5v4hxaoqwp46q2xd22n4iwbniabhu6nfykau, bafkreidyre7waqivwh7ef5hb35rlzogpp3lbt4sdprygoi2ii47vaxh7h4, bafkreihl643tk2lawdfvyuexfrd3gtx3hrksdcrze66vvqubfbjxga3xui, bafkreihtw2i5j7au7xxuetrp2hunanl6rzyaiffg3ibuboqv46jlj56jfe, bafkreicn5t2ghs2d6b3olhhzayh2cle2n67ltu5lzcvxe4sp5x67frs5du, bafkreic3cc7tlg4gbbktkipkdzn2gjptibrmyszl75s4rjkn4umfglgtmu, bafkreihtahpvdzmtjr4ouf5oy3ixv7anfccn7mdda7yr7vz5khkv7pe3k4
scope:           direct local Unix-socket embodiment for live drafting, inspect/search/review operations, and selected approval/supersede actions over the shared runtime, with websocket/HTTP available only through explicit compatibility opt-in
breaking-change: false
notes:           Neovim is a direct terminal embodiment over the shared runtime contract, not an independent runtime. Compatibility mode is explicit because silent cross-adapter fallback would weaken the embodiment contract honesty.

### remote relay (`cmd/operational-relay`)

claim:           extends
spec:            bafkreih5sww7apgdurye4n6el4nsdgwgbcawochqe42t3g3rrkjjehccgy, bafkreibdxszozhrp335vpi5v4hxaoqwp46q2xd22n4iwbniabhu6nfykau, bafkreidyre7waqivwh7ef5hb35rlzogpp3lbt4sdprygoi2ii47vaxh7h4, bafkreihl643tk2lawdfvyuexfrd3gtx3hrksdcrze66vvqubfbjxga3xui, bafkreihtw2i5j7au7xxuetrp2hunanl6rzyaiffg3ibuboqv46jlj56jfe, bafkreicn5t2ghs2d6b3olhhzayh2cle2n67ltu5lzcvxe4sp5x67frs5du, bafkreic3cc7tlg4gbbktkipkdzn2gjptibrmyszl75s4rjkn4umfglgtmu, bafkreihtahpvdzmtjr4ouf5oy3ixv7anfccn7mdda7yr7vz5khkv7pe3k4
scope:           origin-aware relay-feed and CID-addressed blob carriage for the eight frozen families under `/relay/v1`, without claiming to replace the main local runtime or browser/terminal embodiment contracts
breaking-change: false
notes:           The relay extends ex5's deployable PromiseGrid surface by hosting relay/feed/blob transport for frozen-family artifacts. It delegates family semantics to the frozen specs and the runtime's signed-envelope discipline rather than becoming the main app runtime itself.
