# Ex7 decentralized redesign roadmap

## Decision Intent Log

No new architecture decision is introduced by this roadmap. It sequences the
already locked participant-agent boundary in DI-lazim, DI-janup, and DI-sinov
into independently verifiable implementation slices.

## Purpose

Replace Ex7's pre-recreation single-process local demo boundary with a real
decentralized PromiseGrid embodiment. A participant agent, not a website,
browser session, or relay, authors a participant's durable promise. The full
roadmap is published at
[`docs/ex7-decentralized-redesign-roadmap.md`](../docs/ex7-decentralized-redesign-roadmap.md).

## Open slices

- [ ] 1. Participant identity and agent bootstrap.
- [ ] 2. Signed exact-record ingress.
- [ ] 3. Per-agent framed records, blobs, and local projection.
- [ ] 4. Non-authoritative relay carriage.
- [ ] 5. Browser and account embodiment.
- [ ] 6. End-to-end proof and aligned documentation.

## Completion condition

Ex7 is not described as PromiseGrid-complete until every slice has its stated
evidence, the account/relay non-authority boundary is documented, and the
remaining identity recovery or device-authorization work is explicit rather
than implied.
