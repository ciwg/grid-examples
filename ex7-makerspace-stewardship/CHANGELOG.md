# Changelog

## Unreleased — Ex7 participant-agent embodiment

Implements the following frozen Ex7 protocol documents in the scoped local
agent embodiment:

| Contract | Frozen document CID / pCID | Claim |
| --- | --- | --- |
| Makerspace Record v1 | `bafkreid4ebb6ywvwvumetn6pddhyh2pw5uvbvrt6j7wdwv7v7eovb5wdce` | Implements exact canonical signed-record parsing, validation, framed retention, and replay. |
| Equipment Observation v1 | `bafkreifhodcald6kzib36rzeji27hnqjdkeycibnkcigcsz7mzejz6obiy` | Implements payload validation and locally recognized projection. |
| Safety Disposition v1 | `bafkreigt3p2l4uel7wmjr4kple7o55ymchlhh43gajjwsgaeifoogeztc4` | Implements payload validation and local hold/clear assessment; a recognized key still requires the configured steward role to clear. |
| Off-Site Loan v1 | `bafkreibgbky6hbswdimkngjii5zglgvfkolxjxuonbmnqwcdjt4v2qyliq` | Implements payload validation and locally recognized loan projection. |
| Off-Site Return v1 | `bafkreifjgcfwkbwzrgmtj2wfhk3kpzmjbj3wjxid52jau5rrpthabr3ehi` | Implements payload validation and linked locally recognized return projection. |
| Participant Root/History v1 | `bafkreia7cn4srmmkxbwxk2hoezedjvuyokhypcsddjd4evx56lhtmsq3nm` | Implements signed root continuity and recovery-set history. |
| Participant Device Authorization v1 | `bafkreifmbhgjwfmwbemkf4ogsg3gvuavjhttkitzf3muie3dhv5tdn4hq4` | Implements root-authorized device signing for semantic author evidence. |
| Participant Key Revocation v1 | `bafkreify46v4jp3dvz3szem6vrukscplr3kshku7fhrzn4mc26scnnvqoi` | Implements participant-scoped device/root revocation assessment. |
| Participant Threshold Recovery v1 | `bafkreicjdzlriq3nfasza5nmnflpycche63wn2n5kauq66jivlwhhomesy` | Implements declared 2-of-3 recovery-witness verification before root continuation. |
| Participant Terminal Approval v1 | `bafkreidztcisyvexrlia4eos7wko27e4rqt7ivbmikjbkr5tzbbts3rcd4` | Implements unsigned terminal request, local approval, and one-time exact-record return. |
| Participant Peer Card v1 | `bafkreicstci6idwm6d6dbt52ppqyjcapibskz27qmnfuyntg6zck72fa24` | Implements signed root-linked peer-card assessment. |
| Exact Record Carriage v1 | `bafkreihrlojt47erjc6uawkm47s7evppp23tk3ljlkl347ten4v3kb624i` | Implements opaque exact-byte carriage with independent enclosed-record validation. |

This embodiment uses agent-local `recognition.json` policy only for local role
assessment. It retains unknown pCIDs and unrecognized keys without assigning
makerspace semantics. The browser proof transports unsigned terminal requests;
the participant agent returns the author-signed record only after local
approval. It does not claim browser signing, account authoring, a running
relay/discovery service, blob retrieval, portable governance, a global trust
registry, or a PromiseGrid-wide envelope/runtime. Source: DI-tohak; DI-piruf;
DI-likoh; DI-kasaz; DI-sisad; DI-hibok; DI-fuzar.
