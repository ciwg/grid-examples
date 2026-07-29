# Agent Execution Sandbox Boundary

TE ID: TE-safuk

## Status

superseded by TE-dovek / DI-dovek

## Decision under test

Determine the mandatory isolation boundary that OKAR must enforce before it executes independently authored agents, parser agents, or transform agents in response to routed protocol messages.

## Non-negotiable architectural requirement

OKAR must not treat package installation, route registration, trust metadata, or route planning as a substitute for execution isolation. Before an agent is allowed to execute, the runtime must enforce a sandbox boundary that constrains filesystem access, network access, process capabilities, resource consumption, and message authority. This is a required part of the agent-runtime architecture, not a later hardening exercise.

## Assumptions and trust model

- Apps/packages are agents that can be independently authored, upgraded, and loaded from content-addressed storage.
- Routing, parsing, and transforming are agent roles. A routing plan describes intended delivery; it does not grant execution authority.
- Alice may run a cooperative local agent. Mallory may supply a malformed, compromised, or deliberately hostile agent package or payload.
- Direct agent-to-agent messages remain allowed, but they must not bypass the receiver's execution boundary or the node's policy enforcement.
- Current OKAR route planning is a control-plane prototype. It does not yet execute agents, parsers, or transforms.

## Alternatives

### A. In-process Go plugin or shared-runtime execution

Load agent code directly into the OKAR process and rely on package validation, trust policy, and cooperative APIs.

### B. OS-process sandbox per agent invocation or supervised agent worker

Run each agent in a separate operating-system process with an explicit capability manifest, constrained filesystem and network views, resource limits, a bounded message interface, and a supervisor-owned lifecycle.

### C. WebAssembly sandbox as the required first execution environment

Compile or adapt all agents to WASM, expose only host functions granted by OKAR, and execute through a WASM runtime with resource limits and capability-scoped host calls.

### D. Container or microVM isolation as the required first execution environment

Run every agent inside a container or microVM and use the node runtime as the supervisor and message broker.

## Scenario analysis

### Normal local operation

Alice installs a training agent that accepts a known pCID. Under A, Alice's code can directly access the runtime's memory and process authority. B gives the agent a narrow input/output protocol while the supervisor retains state and routing authority. C has similarly strong API-level confinement but requires a WASM tool chain and host ABI from the outset. D provides broad isolation but makes simple local agent startup operationally heavy.

### Hostile package or parser payload

Mallory supplies a parser agent that attempts to read credentials, scan the local network, fork processes, or exhaust memory while processing a routed message. A cannot defensibly contain these actions. B can deny filesystem and network access, bound CPU/memory/time, and terminate the worker. C can deny host functions and meter execution, but unsafe host ABI design can still reintroduce authority. D also contains the workload, but image and orchestration policy become part of the attack surface.

### Parser-first and transform chains

Carol's parser receives an envelope pCID, emits a downstream pCID, and Alice's application handles that output. The parser must receive only its input message and declared output channel; it must not receive Alice's durable store or router credentials. B naturally models this with supervised workers and message pipes. C naturally models it with capability-scoped host calls. A relies on convention. D is workable but imposes a deployment boundary at every hop.

### Mixed-version nodes and gradual adoption

Alice may run an existing native package while Carol publishes a newer agent format. B can introduce one stable worker protocol and admit native executables without changing the language/runtime they were authored in. C requires every agent to be WASM-ready before execution. D requires image or VM packaging conventions. A is easy initially but makes later migration difficult because old agents have already depended on ambient runtime authority.

### Resource pressure and scale

Thousands of short messages make one fresh process per message expensive. B can use supervised, pooled workers with per-invocation limits and restart policies. C may be more efficient for short-lived work but only after the host ABI is stable. D consumes the most operational resources per agent. A is cheapest but shares failure and compromise across the whole node.

### Long-horizon evolution and trust changes

Trust in Alice's author signature can change without changing the fact that her agent remains confined. B, C, and D preserve that separation: trust decides whether to schedule an agent; sandboxing limits what it can do if scheduled. A encourages a dangerous shortcut where a trusted signature is treated as permission to receive ambient process authority.

## Conclusions

Alternative A is rejected. It cannot provide a credible boundary for independently authored agents and would make later confinement a breaking redesign.

Alternatives B, C, and D survive. B is recommended for OKAR's first executable agent slice because it works with the existing executable-package direction, supports gradual migration, and establishes supervisor-owned boundaries without requiring a universal language or heavyweight deployment substrate. The process boundary must be capability-enforced, not merely a child-process convention.

C is a likely additional execution backend once OKAR defines a stable host ABI. D may be appropriate for higher-risk or multi-tenant deployments, but should not be required to prove the first local execution path.

## Decisions still requiring DF

1. Should the first sandboxed backend be a supervised OS-process worker, with WASM and container/microVM backends deferred? Recommended.
2. Should workers be one-shot per invocation or long-lived supervised processes?
3. Which capabilities are denied by default and which initial capabilities, if any, may an agent request: filesystem, network, durable-store access, clock, randomness, subprocess creation, and outbound direct messaging?
4. What public names should identify the execution boundary, worker contract, and capability manifest?
5. How should the runtime report sandbox refusal, timeout, resource exhaustion, or worker crash without inventing new top-level PromiseGrid action kinds?

## Implications for open TODOs and pending DIs

- TODO puvok needs a DI that establishes the execution sandbox as a prerequisite for route-plan execution and agent lifecycle work.
- The registration design in TE-ravuk may proceed as control-plane work, but it must not claim to make registered agents executable.
- Routing, parser, and transform plans remain descriptive until the selected sandbox backend and supervisor contract are implemented.

## Decision status

superseded by TE-dovek / DI-dovek

## Refinements

### 2026-07-29 — Docker backend supersedence

Superseded by TE-dovek: Docker is the selected first execution backend.
