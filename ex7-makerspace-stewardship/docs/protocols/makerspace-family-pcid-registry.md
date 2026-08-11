# Makerspace Family pCID Registry

Status: frozen mappings

This append-only registry maps each first-party Ex7 makerspace family to the
CIDv1 of its exact immutable specification bytes. A pCID identifies a shared
wire-level semantic contract; it is not an identifier for a tool, person,
workflow, photo, or individual promise record. Source: DI-bosur; DI-pihav.

Create a workflow by composing these family pCIDs when its durable meanings are
already covered. Add a new pCID only for a genuinely new interoperable durable
meaning: create a new immutable versioned spec, calculate its CIDv1 from exact
bytes, append its mapping here, and update claims and tests. Do not rewrite or
delete historical specs or mappings. Unknown pCID records remain exact bytes
until a local interpreter and policy can assess them.

The shared envelope specification has document CID
`bafkreid4ebb6ywvwvumetn6pddhyh2pw5uvbvrt6j7wdwv7v7eovb5wdce`; it defines
common carriage and author-evidence slots but is not a family pCID selector.

| Family | Immutable specification | Fixed CIDv1 pCID |
| --- | --- | --- |
| `makerspace.equipment.observation.v1` | `makerspace-families/makerspace-equipment-observation-v1.md` | `bafkreifhodcald6kzib36rzeji27hnqjdkeycibnkcigcsz7mzejz6obiy` |
| `makerspace.safety.disposition.v1` | `makerspace-families/makerspace-safety-disposition-v1.md` | `bafkreigt3p2l4uel7wmjr4kple7o55ymchlhh43gajjwsgaeifoogeztc4` |
| `makerspace.offsite.loan.v1` | `makerspace-families/makerspace-offsite-loan-v1.md` | `bafkreibgbky6hbswdimkngjii5zglgvfkolxjxuonbmnqwcdjt4v2qyliq` |
| `makerspace.offsite.return.v1` | `makerspace-families/makerspace-offsite-return-v1.md` | `bafkreifjgcfwkbwzrgmtj2wfhk3kpzmjbj3wjxid52jau5rrpthabr3ehi` |
