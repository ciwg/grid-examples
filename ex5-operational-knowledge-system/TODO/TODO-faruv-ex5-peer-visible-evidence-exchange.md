# TODO faruv - make ex5 knowledge-evidence peer-visible and portable

## Goal

Add peer-visible `knowledge-evidence` exchange for ex5, including portable
blob carriage that another host can actually resolve.

## Why this exists

The evidence family is frozen and signed locally, but it is still excluded from
peer exchange because attachment bytes are not yet carried in a peer-portable
way.

## Tasks

- [x] faruv.1 Run the required TE for evidence blob carriage and remote
  resolvability.
- [ ] faruv.2 Lock how evidence metadata, attachment CIDs, and blob bytes move
  between peers.
- [ ] faruv.3 Extend peer-exchange bundle rules to cover `knowledge-evidence`.
- [ ] faruv.4 Implement evidence export/import over the settled portable blob
  carriage path.
- [ ] faruv.5 Add round-trip, missing-blob, and tamper coverage for peer
  evidence exchange.

## Status

- open
- `knowledge-evidence` is signed locally today but not yet peer-visible
- `TE-fubok` completed; first blob carriage is locked to self-contained CID-keyed
  evidence plus blobs
- blocked on a newly exposed run-context portability boundary because imported
  evidence still depends on runs that are outside the current peer-visible slice
- `TE-zuvem` filed; awaiting DF between compatibility run carry-along and a
  cleaner run-family-first path
