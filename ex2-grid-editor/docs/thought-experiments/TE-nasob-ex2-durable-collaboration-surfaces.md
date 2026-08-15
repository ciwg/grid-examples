# Ex2 durable collaboration surfaces

TE ID: TE-nasob

## Status

decided

## Decision under test

Where should Ex2 present document activity, comments, saved-version history,
`last viewed`, and `last edited` after live-awareness peers age out of the
roster?

This resolves TODO `mivor.5`.

## Assumptions and trust model

- `live-awareness` is ephemeral, pCID-selected current presence; it is not
  membership, durable activity, or proof of another participant's intent.
- Ex2 is a repo-local-draft example. Browser workflow/review state is local to
  that browser unless a separate profile and evidence design says otherwise.
- The existing `DocumentRegistry` already stores per-document timestamps,
  activity, comments, and saved versions in browser-local preferences.
- Alice may view or edit a document while Bob's last awareness observation
  ages out. Mallory may delay or stop presence traffic. Neither condition
  creates a durable claim about another participant.

## Alternatives

### A. Separate browser-local workflow and review surfaces

Keep document timestamps, activity, comments, and saved versions in the
existing `DocumentRegistry`; render them in the document metadata, activity,
comment, and version surfaces. Keep the peer list limited to current local
presence interpretation.

### B. Retain aged peers in the main roster as historical presence

Show departed peers beside live/stale/offline peers and reuse their last
awareness state as activity history.

### C. Make the relay authoritatively persist collaboration history

Have each relay turn awareness, comments, and document timing into a shared
durable activity/membership feed.

## Scenario analysis

### Normal operation

Alice opens and edits a document. A records her browser-local view/edit events
and renders current peers separately. B makes the roster ambiguous: it no
longer answers who is here now. C expands a current example relay into an
unlocked shared-history authority.

### Failure and delayed traffic

Bob's last awareness update becomes old while Alice continues working. A lets
Bob disappear from the live roster while Alice's local document history stays
available. B mistakes absence for a historical relation. C risks presenting a
relay's incomplete or delayed receipt as a global account of activity.

### Concurrent relays and long-horizon evolution

Different relays see different current awareness observations. A keeps each
embodiment's presence local while leaving room for a future pCID-defined,
evidence-backed collaboration-history design. B couples history to volatile
roster timing. C would require separate decisions for signatures, replication,
retention, privacy, authority, and cross-node convergence.

## Conclusion

Alternative A is the PromiseGrid-aligned design and is already implemented.
`DocumentRegistry` owns local workflow/review metadata; `main.js` renders
separate activity, comment, version, and timestamp surfaces. The live peer
roster remains limited to the most recent accepted awareness observations.

No new pCID, relay persistence, identity claim, or shared evidence contract is
introduced. Any future portable collaboration-history profile requires a new
TE and DI.

## Decision status

Locked by the existing local-registry and review-surface decisions
DI-dovoz, DI-nuvif, DI-safor, DI-lapek, and the live-presence boundary
DI-mivor. This TE records that those existing decisions satisfy `mivor.5`.
