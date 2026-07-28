# EX6 Current State

This document is intentionally blunt.

`ex6` is **not done**.

It currently contains the runtime foundation, one small built-in example
package, and the first nine real first-party built-in packages. It does not
yet contain the full operational knowledge product as a set of ex6 packages.
Source: `DI-moksu`; `DI-lupok`; `DI-lorup`; `DI-vakod`; `DI-pamuk`;
`DI-figar`; `DI-tusav`; `DI-sivuk`; `DI-ramek`; `DI-zibek`; `DI-lavom`;
`DI-zafek`.

## Implemented

- `moks package list`
- `moks package inspect <package-id>`
- `moks package install <dir>`
- `moks route list`
- `moks route inspect <protocol-pcid>`
- `moks route plan <protocol-pcid>`
- `moks route plan <protocol-pcid> trace`
- `moks route plan <protocol-pcid> trace candidate <package-id:role:route-type>`
- `moks route plan <protocol-pcid> trace downstream <protocol-pcid>`
- route-plan trace summaries now distinguish root and downstream protocol scope
- `moks route policy show`
- `moks route policy set <prefer-route-types|-> <avoid-route-types|-> <prefer-roles|-> <avoid-roles|->`
- `moks route policy show <protocol-pcid>`
- `moks route policy set-for <protocol-pcid> <prefer-route-types|-> <avoid-route-types|-> <prefer-roles|-> <avoid-roles|->`
- `moks route policy remove <protocol-pcid>`
- `moks route policy show <protocol-pcid> <role>`
- `moks route policy set-for-role <protocol-pcid> <role> <prefer-route-types|-> <avoid-route-types|-> <prefer-roles|-> <avoid-roles|->`
- `moks route policy remove-role <protocol-pcid> <role>`
- `moks relay export <path>`
- `moks relay import <path>`
- `moks relay serve <addr>`
- `moks relay pull <peer-id>`
- `moks relay push <peer-id>`
- `moks relay peer local show`
- `moks relay peer discover <peer-card-url>`
- `moks relay peer discover <peer-card-url> seed`
- `moks relay peer allow <peer-id> <batch-url> <import-url> <public-key> <pull|no-pull> <push|no-push>`
- `moks relay peer promote <peer-id> <pull|push|both>`
- `moks relay peer classify <peer-id> <class> <weight>`
- `moks relay peer federate <peer-id> <federation>`
- `moks relay peer list`
- `moks relay peer revoke <peer-id>`
- `moks relay policy claim list`
- `moks relay policy claim set <protocol-pcid> <role|*> <min-attesters> <any|peer-id,peer-id>`
- `moks relay policy claim set-weighted <protocol-pcid> <role|*> <min-attesters> <min-weight> <any|peer-id,peer-id> <any|class,class>`
- `moks relay policy claim set-federated <protocol-pcid> <role|*> <min-attesters> <min-weight> <min-federations> <any|peer-id,peer-id> <any|class,class> <any|federation,federation>`
- `moks relay policy claim remove <protocol-pcid> <role|*>`
- package manifest validation
- installed-package self-check
- explicit route registration derived from package claims
- exported route metadata alongside relay claims
- explicit direct/parser/transform route typing
- runtime-owned route-planning preference policy
- per-`protocol_pcid` route-planning policy overrides over the global planner defaults
- per-`protocol_pcid + role` route-planning policy overrides inside one protocol
- route-plan introspection that explains why the preferred route won
- pairwise route-plan comparison detail across the full candidate set
- explicit downstream-plan explanation summaries for parser/transform hops
- route-plan trace mode with exact planner decision sequence
- focused trace filters for one candidate path or one downstream protocol
- trace summary counts showing total, shown, and hidden steps
- append-only durable history
- CAS object storage
- one built-in example package: `ops-note`
- one real first-party built-in package: `context`
- one real first-party built-in package: `knowledge`
- one real first-party built-in package: `runs`
- one real first-party built-in package: `links`
- one real first-party built-in package: `procedures`
- one real first-party built-in package: `training`
- one real first-party built-in package: `maintenance`
- one real first-party built-in package: `receiving`
- one real first-party built-in package: `inventory`
- first-party package source directories for the current 9-package spine
- a starter package template for outside authors
- runtime-mediated durable writes for installed executable packages
- runnable built-in and installed-package examples
- exact-byte relay import dedupe and relay batch metadata validation
- live HTTP peer exchange on top of the current relay batch shape
- peer allow rules and peer-identity matching on live relay import
- runtime-owned local relay keypairs and signed live relay batches
- peer public-key verification on live relay pull/push/import
- relay peer-card discovery without automatic trust grants
- optional seeded peer registration that still starts as `no-pull` and `no-push`
- peer-policy promotion shortcuts that reuse stored metadata
- per-record digest proofs on relay batches
- per-record relay-carriage signatures by the exporting peer
- claim-level proofs for advertised implementation claims
- semantic author-level signatures on durable records
- third-party attestation support for implementation claims
- runtime-owned attestation policy and quorum for implementation claims
- weighted attestation trust and attester classes
- federated trust semantics with distinct federation spread

Discovery rule:
- discovery is not trust
- plain discovery fetches metadata and prints next commands only
- seeded discovery writes a local peer entry, but it still cannot pull or push until explicitly allowed
- policy promotion changes trust locally without re-fetching peer metadata

Source: `DI-moksu`; `DI-lupok`; `DI-rovum`; `DI-nupad`; `DI-sibok`; `DI-nasek`; `DI-rupem`; `DI-zotem`; `DI-vemut`; `DI-kasud`; `DI-lutep`; `DI-zumep`; `DI-ravud`; `DI-luzef`; `DI-sovem`; `DI-fogem`; `DI-movek`; `DI-ravok`; `DI-rumek`; `DI-rutom`; `DI-ruvot`; `DI-lafek`; `DI-fotav`; `DI-pabut`; `DI-matek`; `DI-posek`; `DI-rivuk`; `DI-lavik`; `DI-fobek`; `DI-povak`; `DI-rusom`; `DI-dovak`; `DI-buvok`.

## Not Implemented

- browser embodiment
- Neovim embodiment
- ex5 review/search/operate parity
- richer author identity, cross-runtime federation semantics, and broader proof discipline
- automatic peer discovery with implicit trust and stronger trust beyond explicit allow rules
- scaffold command for new packages
- later domain packages beyond the current first-pass ex5 set

Source: `DI-moksu`; `DI-lupok`; `DI-zotem`; `DI-vemut`; `DI-kasud`; `DI-lutep`; `DI-zumep`; `DI-ravud`; `DI-luzef`; `DI-sovem`; `DI-fogem`; `DI-movek`; `DI-ravok`; `DI-rumek`.

## Why The Gap Exists

The current repo state reflects a deliberate kernel-first slice:

- first lock the runtime contract
- then build real packages against that contract

That first step happened.
The second step has now started with the built-in `context`, `knowledge`,
`runs`, `links`, `procedures`, `training`, `maintenance`, and `receiving`
packages plus `inventory`, but the rest is still mostly open.
Source: `DI-moksu`; `DI-lupok`; `DI-lorup`; `DI-vakod`; `DI-pamuk`;
`DI-figar`; `DI-tusav`; `DI-sivuk`; `DI-ramek`; `DI-zibek`; `DI-lavom`.

## Next Product Work

Best next product work:

1. make package creation easier than manual template copying
2. harden and deepen the current package set
3. harden relay/peer behavior beyond the current batch shell

Source: `DI-moksu`; `DI-rovum`; `DI-nupad`; `DI-sibok`; `DI-sivuk`; `DI-ramek`; `DI-zibek`; `DI-lavom`.
