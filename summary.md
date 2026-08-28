# Grid Exercises: Summary of Programs

This is a practical map of the seven examples in this repository. Each section
gives the first command to try, the important technology, and the actual Grid
boundary. It is a personal refresher, not a repository artifact.

## Ex1 — Order flow

- **Folder:** `~/lab/cswg/grid-examples/ex1-order-flow/`
- **What it does:** An order-fulfilment flow among collector, kernel, seller,
  warehouse, accounting, carrier, and intake roles. The demo includes normal,
  warehouse-refusal, and carrier-timeout scenarios.
- **Technology:** Go; Docker Compose; Bash. No browser, Node, or Neovim.
- **Grid shape:** Roles exchange signed messages through a kernel and retain
  their own raw bytes and local observations. The five profiles are local
  drafts; a timeout or refusal record is not proof of another actor's intent.
- **First run:** `bash docker/run-demo.sh`. It uses
  `/tmp/grid-examples-ex1-data` by default and leaves evidence there.

## Ex2 — Grid Editor

- **Folder:** `~/lab/cswg/grid-examples/ex2-grid-editor/`
- **What it does:** Collaborative document editing in browser and Neovim,
  including CRDT text, presence, metadata, and publish/import flows.
- **Technology:** Go relay; JavaScript; Automerge; CodeMirror; optional Node,
  npm, Neovim, and Docker Compose.
- **Grid shape:** The relay signs, verifies, stores, and exchanges
  pCID-selected envelopes while editors retain their own CRDT replicas. One
  relay normally represents one logical node.
- **First run:** `go run ./cmd/grid-relay --listen 127.0.0.1:7015 --data-root .grid-editor/alice`, then open
  `http://127.0.0.1:7015/?doc=demo`.

## Ex3 — WebSocket Grid Editor

- **Folder:** `~/lab/cswg/grid-examples/ex3-grid-editor-websocket/`
- **What it does:** The Grid Editor shape with WebSocket live carriage,
  remote-access bootstrap, and browser startup/recovery hardening.
- **Technology:** Go relay; JavaScript; Automerge; CodeMirror; optional Node,
  npm, Docker Compose, and Neovim sidecar.
- **Grid shape:** It retains pCID-selected document, awareness, metadata, and
  publishing profiles; WebSocket is carriage, not a new profile. The bootstrap
  token only mints short-lived relay-local capabilities.
- **First run:** `go run ./cmd/grid-relay --listen 127.0.0.1:7025 --data-root .grid-editor/alice`, then open
  `http://127.0.0.1:7025/?doc=demo`.

## Ex4 — Bug Tracker

- **Folder:** `~/lab/cswg/grid-examples/ex4-bug-tracker/`
- **What it does:** Browser-first bug tracker with issue timelines,
  assignments, fixed status transitions, attachments, and engineer CLI.
- **Technology:** Go HTTP server and CLI; browser WebCrypto; JSONL; CAS.
- **Grid shape:** Browser and CLI produce pCID-selected signed issue-report,
  lifecycle, and attachment-reference promises. Enrollment and acceptance are
  service-local policy, not general identity or authority.
- **First run:** `bash scripts/run-demo.sh`, then open `http://127.0.0.1:7035/`.

## Ex5 — Operational Knowledge System

- **Folder:** `~/lab/cswg/grid-examples/ex5-operational-knowledge-system/`
- **What it does:** Operational knowledge, review, receiving, inventory,
  training, maintenance, evidence, and approval application with browser, CLI,
  and early Neovim surfaces.
- **Technology:** Go 1.24.13; append-only storage; CAS; browser
  extension/native host; optional Neovim.
- **Grid shape:** Records bounded signed envelopes, peer/relay exchange, and
  operational projections. Its trust and authority claims are deliberately
  local and explicitly scoped.
- **First run:** `./scripts/load-sample-data.sh /tmp/ex5-newcomer-runtime` then
  `go run ./cmd/operational-knowledge -data-root /tmp/ex5-newcomer-runtime`;
  open `http://127.0.0.1:7045/`.

## Ex6 — Operational Knowledge Agent Runtime (OKAR)

- **Folder:** `~/lab/cswg/grid-examples/ex6-operational-knowledge-agent-runtime/`
- **What it does:** Standalone runtime for built-in and installed packages,
  workflows, relay carriage, package activation, and CLI operation.
- **Technology:** Go `moks` CLI; canonical CBOR Grid records; immutable pCIDs;
  append-only history; CAS; optional Docker worker integration.
- **Grid shape:** Frozen family specs name durable records; unknown family bytes
  are retained without inferred meaning. Semantic author evidence, relay
  carriage, and local policy remain distinct.
- **First run:** `./examples/builtin-quickstart.sh` for a throwaway workspace.

## Ex7 — Makerspace Stewardship

- **Folder:** `~/lab/cswg/grid-examples/ex7-makerspace-stewardship/`
- **What it does:** Stores and assesses makerspace records: participant roots
  and devices, recovery, revocation, peer cards, equipment observations, and
  local role-policy projection.
- **Technology:** Go server; browser display/ingress; canonical record frames;
  local storage; Chrome DevTools for the two-agent proof.
- **Grid shape:** Accepts externally signed pCID-selected records and projects
  only recognized, authorized records. The browser does not sign; unrecognized
  records remain retained evidence without changing the view.
- **First run:** `go run ./cmd/makerspace-stewardship`, then open
  `http://127.0.0.1:7037/`. Read the operator/testing guides before using
  recognition policy or the browser proof.
