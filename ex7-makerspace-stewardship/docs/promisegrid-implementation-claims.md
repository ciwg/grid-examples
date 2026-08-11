# PromiseGrid implementation claims

Ex7 implements its four frozen makerspace record families and the scoped
`makerspace-record-v1` profile for exact signed ingress, framed local
retention, and local projection under `recognition.json`. For the implemented
root-history, device-authorization, revocation, threshold-recovery, and
peer-card pCIDs, it verifies signed participant history in record order before
treating a root or device key as an author of a makerspace record. Two matching
recovery witnesses may authorize a replacement-root continuation; conflicting
witness evidence is retained without choosing either replacement. A peer-card
linked exact-byte carriage wrapper retains its opaque enclosed records and
then validates each enclosed record independently. `recognition.json` is only
local role assessment. Unknown pCIDs and unrecognized keys are retained
without known semantics. This claim excludes browser signing, accounts as
author evidence, a running relay/discovery service, blob retrieval, and
portable governance. Source: DI-tohak; DI-piruf; DI-likoh; DI-kasaz; DI-sisad.

The implemented browser embodiment is limited to unsigned request transport:
Bob's terminal page names an explicit loopback target for Alice's participant
agent, Alice's local page approves or declines the request, and Bob polls a
one-time return location for the exact signed record. Alice's authorized device
key supplies author evidence only after local approval. The browser origin,
account state, target URL, and approval token do not supply author evidence.
`scripts/run-two-agent-browser-proof.sh` is the repeatable Chrome evidence
path; it prints a disposable `/tmp/ex7-browser-proof.XXXXXX` evidence root
containing `bob-final.png`, the approval response, and process logs. Source:
DI-fuzar; DI-hibok; DI-kasaz.
