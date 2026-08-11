# Participant Identity, Authorized Devices, and Cross-Device Approval

TE ID: TE-bunot

## Status

superseded by TE-folok / DI-janup

## Decision under test

How Ex7 identifies a participant across personal devices and shared makerspace browsers without central accounts or browser-profile identity.

This supersedes TE-ranib's browser-held-key assumption.

## Assumptions

- Alice may borrow from her laptop and return using a makerspace kiosk.
- Identity must survive browser replacement, device loss, and new devices.
- A kiosk, relay, and UI never receive Alice's private identity key.

## Alternatives

### A. Participant root identity with signed device authorization

Alice has a durable root identity key and append-only signed key-continuity chain. Alice authorizes device keys using frozen device-authorization records. Devices sign ordinary makerspace records and include identity/device evidence. Kiosks display requests and use QR or short-lived pairing to send them to Alice's approved phone, hardware signer, or agent. Alice's device signs; the kiosk receives only exact record bytes.

### B. One copied private key

Alice exports and imports one private key into every browser/device.

### C. Central login/account recovery

A hosted service maps Alice to devices and recovers keys.

## Scenario analysis

A lets Alice use any authorized device without creating a new author identity. Her laptop signs a loan and her phone signs a kiosk return; the kiosk never holds a key. If Alice loses a phone, another authorized/root recovery path signs a device-revocation record. Agents retain old evidence but stop treating that device as active under local policy. If all recovery keys are lost, continuity is not fabricated; a later witness/recovery protocol is required.

B makes every copied device a full compromise point and makes revocation coarse. C creates the central identity authority Ex7 rejects. Under partitions, A retains exact authorization/revocation evidence and reconciles later through non-authoritative carriage. Conflicting chains remain evidence, not registry edits.

## Conclusions

B and C are rejected. A is recommended: participant root identity, signed authorized-device records, signed device revocation/continuity evidence, and cross-device approval where shared browsers only display/carry requests.

## Output to decision framing

**A is recommended.** It requires frozen pCID families for participant key continuity, device authorization, and device revocation before Ex7 signs makerspace records. Browser WebCrypto may be one authorized-device embodiment, but is not participant identity.

## Decision status

needs DF
