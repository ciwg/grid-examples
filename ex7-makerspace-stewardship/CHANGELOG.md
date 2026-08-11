# Changelog

## Unreleased — Ex7 signed-record ingress embodiment

Implements the following frozen Ex7 protocol documents in the scoped local
agent embodiment:

| Contract | Frozen document CID / pCID | Claim |
| --- | --- | --- |
| Makerspace Record v1 | `bafkreid4ebb6ywvwvumetn6pddhyh2pw5uvbvrt6j7wdwv7v7eovb5wdce` | Implements exact canonical signed-record parsing, validation, framed retention, and replay. |
| Equipment Observation v1 | `bafkreifhodcald6kzib36rzeji27hnqjdkeycibnkcigcsz7mzejz6obiy` | Implements payload validation and locally recognized projection. |
| Safety Disposition v1 | `bafkreigt3p2l4uel7wmjr4kple7o55ymchlhh43gajjwsgaeifoogeztc4` | Implements payload validation and local hold/clear assessment; a recognized key still requires the configured steward role to clear. |
| Off-Site Loan v1 | `bafkreibgbky6hbswdimkngjii5zglgvfkolxjxuonbmnqwcdjt4v2qyliq` | Implements payload validation and locally recognized loan projection. |
| Off-Site Return v1 | `bafkreifjgcfwkbwzrgmtj2wfhk3kpzmjbj3wjxid52jau5rrpthabr3ehi` | Implements payload validation and linked locally recognized return projection. |

This embodiment accepts externally signed bytes and uses agent-local
`recognition.json` policy. It retains unknown pCIDs and unrecognized keys
without assigning makerspace semantics. It does not claim browser signing,
account authoring, key continuity/recovery, relay carriage, blob retrieval, or
portable governance. Source: DI-tohak; DI-piruf; DI-likoh.
