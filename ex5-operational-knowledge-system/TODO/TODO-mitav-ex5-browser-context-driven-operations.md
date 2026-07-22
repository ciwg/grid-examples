# TODO mitav - ex5 browser context-driven operations

## Decision Intent Log

ID: DI-mitav
Date: 2026-07-22 21:58:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track context-driven browser operation flows as their own ex5 UI TODO.
Intent: Make browser run recording, evidence capture, and approvals feel more like actions taken from the current record context and less like separate generic forms.
Constraints: Preserve the existing generic forms and manual overrides; add stronger context-derived actions and defaults instead of removing the broader form surfaces.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-mitav-ex5-browser-context-driven-operations.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Make high-frequency browser operational actions start from current context
rather than from blank generic forms.

## Tasks

- [x] mitav.1 Review the current browser record inspector and operation forms for context handoff gaps.
- [x] mitav.2 Define the smallest context-driven run, evidence, and approval action improvements.
- [x] mitav.3 Implement those browser action paths while preserving generic fallbacks.
