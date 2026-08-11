# Participant Agents with Account Bootstrap Only

TE ID: TE-folok

## Status

decided

## Decision

Ex7 participants are independently owned agents with their own durable signing identity and local evidence store. Existing makerspace website accounts are only UI login, discovery/bootstrap, and local makerspace policy. An account session is not cryptographic author evidence and cannot sign or create a participant promise.

## Scenarios

Alice may sign in to a makerspace kiosk with her existing account, which lets the UI find or request a connection to Alice's own agent. The kiosk displays an exact reservation/loan request and Alice's reachable agent signs it. If Alice's agent is unavailable, the kiosk may retain an unsigned draft but cannot claim an Alice promise. Carol's agent signs Carol's inspection separately.

The account service may fail, be replaced, or disagree without changing any signed record's author. A relay carries bytes but cannot authenticate as Alice, sign for Alice, or decide state. Each agent locally decides how account/bootstrap information informs UI and recognition policy.

## Conclusion

This is Ex7's smallest honest decentralized boundary. Device authorization and recovery are future protocols, not core assumptions.

## Decision status

locked by DI-janup.
