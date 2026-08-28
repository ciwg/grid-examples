# TODO haruv — Add tailored exercise Makefiles

## Decision Intent Log

### DI-jigap

- ID: DI-jigap
- Date: 2026-08-28 09:45:02 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Add one Makefile to each exercise folder with targets named for that exercise's checked-in run, demo, proof, build, or verification workflow rather than a uniform target set.
- Intent: Let a returning user discover and invoke the documented workflow from the exercise root without concealing distinct process, container, or evidence requirements.
- Constraints: Each recipe must delegate to existing commands or scripts; do not add a target for Ex3's destructive reset script; targets that start servers remain foreground processes; runtime locations and cleanup remain controlled by existing commands and scripts.
- Affects: `ex1-order-flow/Makefile`, `ex2-grid-editor/Makefile`, `ex3-grid-editor-websocket/Makefile`, `ex4-bug-tracker/Makefile`, `ex5-operational-knowledge-system/Makefile`, `ex6-operational-knowledge-agent-runtime/Makefile`, and `ex7-makerspace-stewardship/Makefile`.

## Tasks

- [x] haruv.1 Add exercise-specific target menus.
- [x] haruv.2 Verify target discovery and safe local checks.
