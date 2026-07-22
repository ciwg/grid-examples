# TODO pafur - ex5 browser dominant review path

## Decision Intent Log

ID: DI-pafur
Date: 2026-07-22 22:15:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining ambiguity in the browser landing job as its own ex5 UI TODO.
Intent: Make the browser declare one unmistakable primary review path so operators know where to start and what to do next after landing.
Constraints: Preserve the existing review surfaces, search, problem review, and inspector detail; clarify the dominant path instead of removing alternative entry points.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-pafur-ex5-browser-dominant-review-path.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Make the browser landing experience clearly communicate the main review job and
the default next step sequence.

## Tasks

- [x] pafur.1 Review the current browser landing surfaces for competing “start here” signals.
- [x] pafur.2 Define the smallest changes that make one review path clearly dominant without reducing the other entry points.
- [x] pafur.3 Implement that stronger landing-path guidance in the browser shell and inspector.
