# Ex5 Implementation Promise Claims

## knowledge-item

claim:           partially-implements
spec:            bafkreih5sww7apgdurye4n6el4nsdgwgbcawochqe42t3g3rrkjjehccgy
scope:           local-runtime signed envelopes for `knowledge_item_created`, `revision_added`, `knowledge_item_status_changed`, and `knowledge_item_superseded`, with browser/CLI/Neovim still using the local HTTP adapter and the remaining ex5 families still bridged
breaking-change: false
notes:           This is the first PromiseGrid-native runtime slice in ex5. It freezes the `knowledge-item` family first and leaves approval, evidence, link, responsibility, and search-metadata families for later staged migration.

## knowledge-approval

claim:           partially-implements
spec:            bafkreibdxszozhrp335vpi5v4hxaoqwp46q2xd22n4iwbniabhu6nfykau
scope:           local-runtime signed envelopes for `approval_recorded`, covering both knowledge-item and run approvals while browser/CLI/Neovim still use the local HTTP adapter and the remaining ex5 families still bridged
breaking-change: false
notes:           This is the second PromiseGrid-native runtime slice in ex5. It freezes `knowledge-approval` as one durable family for named-role review outcomes while keeping lifecycle status changes in `knowledge-item`.

## knowledge-evidence

claim:           partially-implements
spec:            bafkreidyre7waqivwh7ef5hb35rlzogpp3lbt4sdprygoi2ii47vaxh7h4
scope:           local-runtime signed envelopes for `evidence_added`, covering evidence metadata plus attachment name/path/size references while browser/CLI/Neovim still use the local HTTP adapter and copied attachment bytes remain on current runtime storage
breaking-change: false
notes:           This is the third PromiseGrid-native runtime slice in ex5. It freezes `knowledge-evidence` around structured evidence plus attachment references while leaving attachment bytes outside the family.

## knowledge-link

claim:           partially-implements
spec:            bafkreihl643tk2lawdfvyuexfrd3gtx3hrksdcrze66vvqubfbjxga3xui
scope:           local-runtime signed envelopes for `link_added`, covering typed operational links while browser/CLI/Neovim still use the local HTTP adapter and the remaining ex5 families still bridged
breaking-change: false
notes:           This is the fourth PromiseGrid-native runtime slice in ex5. It freezes `knowledge-link` around stable typed link artifacts without changing the embodiment adapter contract.

## knowledge-responsibility

claim:           partially-implements
spec:            bafkreihtw2i5j7au7xxuetrp2hunanl6rzyaiffg3ibuboqv46jlj56jfe
scope:           local-runtime signed envelopes for `responsibility_created`, covering first-class durable responsibilities while browser/CLI/Neovim still use the local HTTP adapter and the remaining ex5 families still bridged
breaking-change: false
notes:           This is the fifth PromiseGrid-native runtime slice in ex5. It freezes `knowledge-responsibility` around stable role-bearing operational duty records while keeping links, runs, and search projections on the current adapter path.
