# Coordinated Finished Ex7 Contract Set

TE ID: TE-zizum

## Status

decided

## Decision under test

The exact coordinated contracts required to implement the finished TE-bijad
participant-agent product without later architecture substitutions.

## Proposed complete contract set

1. **Participant root/history.** An offline-capable root key authorizes device
   keys; an active key or configured recovery threshold may append replacement
   and revocation evidence. Every transition is signed, exact, durable, and
   locally assessed.
2. **Devices and recovery.** Ordinary makerspace records are signed by active
   device keys. A participant predeclares an `m-of-n` recovery set; recovery
   activates a replacement only with the threshold's exact signed evidence.
   Conflicting histories remain evidence for local assessment.
3. **Personal and terminal embodiment.** A personal device signs. A shared
   terminal creates an unsigned request containing exact proposed payload
   bytes; it receives a signed result from a reachable authorized device or
   remains visibly unsigned. It stores no participant private key.
4. **Account boundary.** A makerspace account only opens the terminal UI and
   supplies local facility context. It cannot authorize a device, alter key
   history, sign, recover, or project a record.
5. **Discovery and carriage.** A signed peer card advertises participant root,
   active-device history reference, and optional contact hints. Direct and
   relay carriage use a frozen pCID contract for opaque exact records, cursor,
   duplicate, retention, and outage behavior. Carriers do not interpret or
   author makerspace semantics.
6. **Local state.** Each agent retains records/blobs, verifies known contracts,
   preserves unknown bytes, and derives only a local projection and trust view.
   `recognition.json` becomes bootstrap input, never a substitute for signed
   key history.
7. **Finished proof.** Alice signs from two devices; uses a terminal; loses and
   revokes one device; recovers through threshold evidence; Carol performs a
   steward action; Mallory attacks labels/replay/unknown pCIDs; a relay outage,
   duplicate delivery, restart, and independent projections are proven.

## Alternatives

### A. Freeze and implement this entire coordinated set

All seven contracts receive immutable specs, pCIDs/doc-CIDs, runtime support,
and end-to-end evidence before Ex7 is called complete.

### B. Implement only selected portions

Leave recovery, terminal use, discovery, or carriage as assumed future work.

## Conclusion

B is rejected because it recreates the incomplete-product boundary TE-bijad
was created to end. A is the only finished PromiseGrid-aligned Ex7 scope.

## Output to decision framing

**A — freeze and implement the whole coordinated contract set.** The remaining
DF must select exact `m`/`n`, recovery-key roles, payload fields, pCID paths,
runtime paths, and acceptance vectors as one lock.

## Decision status

locked: Alternative A and the coordinated parameters below, by DI-girup in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`.

## Locked parameters

- Recovery is participant-selected 2-of-3 recovery keys; active device/root
  keys may revoke or replace, while recovery needs two distinct witnesses.
- Frozen family paths are `docs/protocols/participant/` for root history,
  device authorization, revocation, and recovery; `peer/` for peer card; and
  `carriage/` for direct/relay feed. Each gets a registry entry and CID check.
- Per-agent runtime paths are `<agent-root>/identity/`, `recognition.json`,
  `records.frames`, `blobs/`, and `requests/`; terminal drafts are bounded,
  unsigned request files under `requests/` and never record evidence.
- Acceptance proves two-device signing, terminal approval/draft fallback,
  revocation, 2-of-3 recovery, peer-card verification, relay outage/duplicate
  handling, unknown preservation, Mallory containment, restart, and distinct
  Alice/Carol local projections.
