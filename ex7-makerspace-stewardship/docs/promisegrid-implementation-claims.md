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
