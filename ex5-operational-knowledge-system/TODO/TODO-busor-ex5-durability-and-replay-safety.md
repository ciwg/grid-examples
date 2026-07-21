# TODO busor - ex5 durability and replay safety

## Decision Intent Log

ID: DI-busor
Date: 2026-07-21 12:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Split the review findings so durability and replay hazards are tracked as their own ex5 fix TODO.
Intent: Make the highest-risk persistence issues visible and independently executable.
Constraints: Focus on evidence attachment immutability and event-log replay robustness; keep docs and tests in the same implementation pass.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vurab-ex5-review-followups.md`, `ex5-operational-knowledge-system/TODO/TODO-busor-ex5-durability-and-replay-safety.md`

## Goal

Fix the durability problems that can corrupt historical evidence or make the
runtime fail to replay stored events after restart.

## Tasks

- [x] busor.1 Make evidence attachment storage immutable so later uploads cannot overwrite older evidence bytes through the same attachment path.
- [x] busor.2 Raise or replace the event-log replay scanner limit and add restart coverage for large item or revision bodies.

## Status

- done
- derived from the 2026-07-21 extensive ex5 review
