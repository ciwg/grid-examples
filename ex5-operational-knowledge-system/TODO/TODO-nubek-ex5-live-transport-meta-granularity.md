# TODO nubek - ex5 live transport meta granularity

## Decision Intent Log

ID: DI-torak
Date: 2026-07-22 17:15:03 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Replace the singular `live_draft_preferred_transport` field in `/api/meta` with embodiment-specific `browser_live_draft_transport` and `neovim_live_draft_transport` fields.
Intent: Make the capability contract reflect the shipped embodiment split directly instead of compressing it into one misleading global answer.
Constraints: Keep `/api/meta` as the capability surface and remove the obsolete singular field rather than carrying compatibility baggage.
Affects: `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `docs/thought-experiments/TE-rovek-ex5-live-transport-meta-granularity.md`
Supersedes: DI-latik

ID: DI-latik
Date: 2026-07-22 19:12:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the over-compressed `/api/meta` live transport contract as a dedicated follow-on.
Intent: Replace the single global `live_draft_preferred_transport` idea with capability metadata that reflects the shipped embodiment split more honestly.
Constraints: Keep `/api/meta` as the capability surface; focus on contract clarity, not on adding new transports.
Affects: `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nubek-ex5-live-transport-meta-granularity.md`

## Goal

Make ex5 live transport capability metadata embodiment-aware instead of
compressing browser and terminal behavior into one global preferred-transport
field.

## Tasks

- [x] nubek.1 Define the embodiment-aware capability shape for live transport metadata.
- [x] nubek.2 Implement the refined `Meta` contract and align tests.
- [x] nubek.3 Update the HTTP adapter docs to describe the refined capability fields accurately.

## Status

- closed
- resolved by replacing the singular live-transport field with embodiment-specific metadata
