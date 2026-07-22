# TODO bavum - ex5 browser mode separation

## Decision Intent Log

ID: DI-bavum
Date: 2026-07-22 21:55:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track stronger browser mode separation as its own ex5 UI TODO.
Intent: Reduce the remaining “one oversized surface” feeling by making review, authoring, operation, creation, and browsing feel more behaviorally distinct without removing any current functionality.
Constraints: Preserve the current single-page browser shell and keep every current browser workflow reachable; focus on interaction structure and mode clarity rather than backend changes.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-bavum-ex5-browser-mode-separation.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Make the browser feel like a set of clear workflow modes instead of one large
multi-purpose surface.

## Tasks

- [x] bavum.1 Review the current section transitions and interaction boundaries.
- [x] bavum.2 Define stronger mode-specific affordances for Review, Author, Operate, Create, and Browse.
- [x] bavum.3 Implement the smallest browser changes that make the modes feel behaviorally distinct.
