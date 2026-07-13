# TODO dafek - grid-editor hardening review

## Decision Intent Log

ID: DI-povuz
Date: 2026-07-12 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the first post-CRDT review pass as a dedicated hardening queue covering relay security, sync robustness, and embodiment usability gaps.
Intent: Keep the highest-risk follow-up work visible in one place after the relay/browser/Neovim slices landed, instead of letting security and operator-facing regressions stay implicit in chat history.
Constraints: This queue records findings only; it does not by itself lock the implementation strategy for auth, capability tokens, or transport changes that may need later TE/DI work.
Affects: `ex2-grid-editor/service`, `ex2-grid-editor/identity`, `ex2-grid-editor/web`, `ex2-grid-editor/nvim`, `ex2-grid-editor/scripts`, `ex2-grid-editor/TODO/TODO.md`

Goal: Close the highest-value security, correctness, and usability gaps found in the July 12, 2026 review pass.

- [ ] dafek.1 Lock down the relay signing surface so arbitrary reachable HTTP clients cannot make the local relay sign and publish document or awareness messages as the local author.
  Evidence: `service/server.go`, `service/app.go`
- [ ] dafek.2 Bind decoded payload `author` values to the verified proof key ID during relay ingestion so peers cannot sign with one key while claiming another author identity in the payload.
  Evidence: `service/app.go`, `identity/store.go`
- [ ] dafek.3 Add request validation and resource bounds for `document_id`, `participant_id`, `message_base64`, and feed pagination so malformed or oversized requests cannot collapse participants or force unbounded memory/response growth.
  Evidence: `service/server.go`, `service/app.go`, `crdt`, `awareness`
- [ ] dafek.4 Surface relay sync failures clearly in the browser and peer-polling paths instead of silently dropping failed POSTs or peer-fetch errors.
  Evidence: `web/src/automerge-relay.js`, `service/app.go`
- [ ] dafek.5 Unify the default relay URL and connection docs across the launcher, Neovim plugin, sidecar, and README so manual and scripted startup paths stop disagreeing about `7001` vs `7015`.
  Evidence: `nvim/lua/grid_editor/init.lua`, `cmd/grid-nvim-sidecar/main.go`, `scripts/grid-editor-nvim`, `README.md`
