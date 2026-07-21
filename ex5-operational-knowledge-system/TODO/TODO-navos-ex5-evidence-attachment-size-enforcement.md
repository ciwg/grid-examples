# TODO navos - ex5 evidence attachment size enforcement

## Decision Intent Log

ID: DI-navos
Date: 2026-07-21 13:25:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track evidence attachment size enforcement as its own ex5 bug fix TODO.
Intent: Prevent oversized uploads from being silently accepted or truncated when the runtime intends to enforce a maximum attachment size.
Constraints: Keep the fix local to ex5 evidence upload behavior; add regression tests and update the HTTP/docs contract in the same pass.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-navos-ex5-evidence-attachment-size-enforcement.md`, `ex5-operational-knowledge-system/service/server.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/README.md`

## Goal

Make evidence uploads reject oversized attachments explicitly instead of reading
past the intended limit and saving partial bytes as if the upload succeeded.

## Review finding

`handleEvidence` caps the request body and reads the uploaded file through
`io.LimitReader(file, maxEvidenceAttachmentBytes+1)`, but it never rejects the
`+1` overflow case before calling `AddEvidence`. That means an attachment just
over the intended limit can be accepted and stored instead of failing clearly.

## Tasks

- [x] navos.1 Reject attachments whose decoded file body exceeds `maxEvidenceAttachmentBytes` instead of storing truncated or overflow bytes.
- [x] navos.2 Add HTTP regression coverage for just-over-limit uploads and document the enforced attachment size contract.

## Status

- done
- derived from the 2026-07-21 deep ex5 review
