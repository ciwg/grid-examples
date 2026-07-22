# Ex5 Implementation Promise Claims

claim:           partially-implements
spec:            bafkreih5sww7apgdurye4n6el4nsdgwgbcawochqe42t3g3rrkjjehccgy
scope:           local-runtime signed envelopes for `knowledge_item_created`, `revision_added`, `knowledge_item_status_changed`, and `knowledge_item_superseded`, with browser/CLI/Neovim still using the local HTTP adapter and the remaining ex5 families still bridged
breaking-change: false
notes:           This is the first PromiseGrid-native runtime slice in ex5. It freezes the `knowledge-item` family first and leaves approval, evidence, link, responsibility, and search-metadata families for later staged migration.
