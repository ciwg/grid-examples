# TODO sifeg - ex5 CLI search query encoding

## Decision Intent Log

ID: DI-sifeg
Date: 2026-07-21 16:12:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track CLI search query encoding as its own ex5 review TODO.
Intent: Keep the CLI search embodiment trustworthy for normal operator queries that include spaces or reserved URL characters.
Constraints: Stay inside `ex5`; preserve the current CLI command shape; fix URL construction and add targeted regression coverage.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-sifeg-ex5-cli-search-query-encoding.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`

## Goal

Make `oks-cli search ...` encode the query string correctly before it hits the
HTTP adapter.

## Review finding

The CLI currently concatenates the raw search string onto `?q=`. Queries with
spaces or reserved URL characters can therefore fail or change meaning before
they reach the service.

## Tasks

- [x] sifeg.1 Encode CLI search query parameters correctly.
- [x] sifeg.2 Add regression coverage for queries with spaces and reserved characters.
- [x] sifeg.3 Update CLI docs/examples if the user-facing guidance needs tightening.

## Status

- done
- derived from the 2026-07-21 deep ex5 review
