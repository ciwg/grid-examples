# PromiseGrid implementation claims

Ex7 implements the following frozen document identities in its scoped local
participant-agent embodiment. `makerspace-record-v1` supplies the exact
canonical signed-record profile; each other pCID selects one independently
versioned meaning. Source: DI-tohak; DI-piruf; DI-likoh; DI-kasaz; DI-sisad.

| Contract | Frozen document CID / pCID | Implemented boundary |
| --- | --- | --- |
| Makerspace Record v1 | `bafkreid4ebb6ywvwvumetn6pddhyh2pw5uvbvrt6j7wdwv7v7eovb5wdce` | Canonical parsing, signature verification, framed exact-byte retention, and replay. |
| Equipment Observation v1 | `bafkreifhodcald6kzib36rzeji27hnqjdkeycibnkcigcsz7mzejz6obiy` | Payload validation and locally recognized projection. |
| Safety Disposition v1 | `bafkreigt3p2l4uel7wmjr4kple7o55ymchlhh43gajjwsgaeifoogeztc4` | Payload validation and local hold/clear assessment; clearing additionally requires the local steward role. |
| Off-Site Loan v1 | `bafkreibgbky6hbswdimkngjii5zglgvfkolxjxuonbmnqwcdjt4v2qyliq` | Payload validation and local loan projection. |
| Off-Site Return v1 | `bafkreifjgcfwkbwzrgmtj2wfhk3kpzmjbj3wjxid52jau5rrpthabr3ehi` | Payload validation and local linked-return projection. |
| Participant Root/History v1 | `bafkreia7cn4srmmkxbwxk2hoezedjvuyokhypcsddjd4evx56lhtmsq3nm` | Root continuity and declared recovery-set history. |
| Participant Device Authorization v1 | `bafkreifmbhgjwfmwbemkf4ogsg3gvuavjhttkitzf3muie3dhv5tdn4hq4` | Root-authorized device keys for ordinary semantic author evidence. |
| Participant Key Revocation v1 | `bafkreify46v4jp3dvz3szem6vrukscplr3kshku7fhrzn4mc26scnnvqoi` | Participant-scoped key revocation. |
| Participant Threshold Recovery v1 | `bafkreicjdzlriq3nfasza5nmnflpycche63wn2n5kauq66jivlwhhomesy` | Two matching witnesses from the declared three-key recovery set before replacement-root continuation. |
| Participant Terminal Approval v1 | `bafkreidztcisyvexrlia4eos7wko27e4rqt7ivbmikjbkr5tzbbts3rcd4` | Unsigned terminal request, local approval, and one-time exact-record return. |
| Participant Peer Card v1 | `bafkreicstci6idwm6d6dbt52ppqyjcapibskz27qmnfuyntg6zck72fa24` | Signed root-linked contact/card evidence. |
| Exact Record Carriage v1 | `bafkreihrlojt47erjc6uawkm47s7evppp23tk3ljlkl347ten4v3kb624i` | Opaque exact-byte wrapper; enclosed records are independently validated and locally assessed. |

The runtime applies participant history in record order before treating a root
or device key as an author of a makerspace record. `recognition.json` is local
role assessment only. Unknown pCIDs and unrecognized keys are retained without
known semantics; a carriage signature is not substituted for enclosed author
evidence.

This claim excludes browser signing, accounts as author evidence, a running
relay/discovery service, blob retrieval, portable governance, a global trust
registry, and a PromiseGrid-wide envelope or runtime standard.

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
