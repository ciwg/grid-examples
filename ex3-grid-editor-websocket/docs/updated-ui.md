# Updated UI

- Page: `http://127.0.0.1:7026/?doc=demo&access_token=ex3-demo-access`
- Screenshot: `screenshots/ex3-updated-ui.png`
- Captured: `2026-07-19 12:43 America/Los_Angeles`

![Ex3 updated UI](screenshots/ex3-updated-ui.png)

## Visible State In The Captured Screenshot

The screenshot shows the default `paper` theme with the `demo` document open and the page loaded near the top of the layout.

Important note:
- this screenshot was captured **before** the later sidebar reorder that moved
  `PromiseGrid Flow` above `Document`
- the current code now places `PromiseGrid Flow` first in the left sidebar,
  then `Relay`, then `Document`
- that means the screenshot is still useful for the overall shell look, but
  the exact top-to-bottom sidebar order is now slightly out of date

Visible status in the screenshot:

- Header brand: `grid-editor`
- Subtitle: `PromiseGrid collaborative editor example`
- Current document ID: `demo`
- Current title: `Document demo`
- Peer count pill: `0 peers`
- Share link line: `http://127.0.0.1:7026/?doc=demo&access_token=ex3-demo-access`
- Top-right connection state: `connecting...`
- Auto-save state: `auto-save idle`
- Message counter: `messages: -`
- Replica indicator: `local replica: -`
- Quick start banner is visible
- One visible peer badge row entry: `Browser User  you`
- Main editor pane is visually empty in this screenshot

## Top-Level Layout

The page is split into two main columns:

1. A left sidebar of stacked cards for PromiseGrid flow, relay details, document controls, identity, workspace, metadata, peers, and review/history.
2. A right editor shell containing the toolbar, session status, quick-start banner, document tabs, presence badges, editor pane, preview pane, and modal overlays.

The captured screenshot only shows the upper portion of the sidebar. The lower cards still exist on the page but are below the first viewport cut.

## Sidebar Sections

### 1. Sidebar Header

Purpose:
- brand the demo as `grid-editor`
- identify it as a PromiseGrid collaborative editor example
- expose a `Help` button

Visible in screenshot:
- yes

### 2. PromiseGrid Flow Card

Purpose:
- show the live data-flow story for the current document
- display the transport mode in use
- show relay-observed message traffic for the current document
- let the user click a message for decoded inspection

Layout note:
- this card is now at the top of the left sidebar
- it is not the top of the whole page, because the page header still sits
  above the cards

Elements:
- transport pill
- flow diagram
- trace caption
- `message-trace` list

Flow diagram labels:
- `Browser`
- `signed grid message`
- `Relay`
- `peer feed`
- `Peer relay`
- `websocket fanout`
- `Other editor`

Visible in screenshot:
- no, this screenshot predates the sidebar reorder

### 3. Relay Card

Purpose:
- expose the local relay/author identity and the pCID identifiers used by the demo

Readouts:
- `Author`
- `live-document pCID`
- `live-awareness pCID`
- `document-metadata pCID`
- `publish-document pCID`
- `demo profile` pill

Visible in screenshot:
- no, below the first viewport cut

### 4. Document Card

Purpose:
- choose which shared document is open
- rename the current document
- create or duplicate a shared document
- copy/share/paste the current tokenized link
- show simple document timestamps

Controls:
- `Document ID` text box
- `Open` button
- `Title` text box
- `New Shared Doc`
- `Duplicate`
- `Paste Link`
- `Copy Link`
- `Email Link`

Readouts:
- `Current link`
- `Created`
- `Last viewed`
- `Last edited`
- `Last exported`
- peer count pill in the card header

Visible in screenshot:
- yes

Current screenshot values:
- document id: `demo`
- title: `Document demo`
- peer count: `0 peers`
- last edited: `-`
- last exported: `-`

### 5. You Card

Purpose:
- define the local participant identity that the browser advertises to the relay and other editors

Controls:
- `Display name`
- `Color`
- `Settings` button

Readouts:
- color preview swatch
- color preview name/value
- participant ID

Visible in screenshot:
- partially visible

Current screenshot values:
- display name: `Browser User`
- color picker: blue value corresponding to the current participant color

### 6. Workspace Card

Purpose:
- show open tabs, recent docs, and templates
- provide quick doc/template generation actions

Subsections:
- `Open tabs`
- `Recent docs`
- `Templates`

Buttons:
- `Generate Demo Doc`
- `Sample Doc`

Visible in screenshot:
- no, below the first viewport cut

### 7. Metadata Card

Purpose:
- edit relay-backed descriptive metadata for the document

Fields:
- `Description`
- `Summary`
- `Tags`
- `Collections`
- `Favorite`
- `Archived`
- `Save Metadata`

Visible in screenshot:
- no

### 8. Peers Card

Purpose:
- list live participants and their presence state

Readouts:
- live presence legend
- peer list

Visible in screenshot:
- no

### 9. Review Card

Purpose:
- collect non-editor support views for inspection and presentation

Subsections:
- `Outline`
- `Saved versions`
- `Recent participants`
- `Activity`
- `Published exchanges`
- `Catalog search`

Catalog search controls:
- metadata search query
- include archived toggle
- `Search Catalog`
- results list

Visible in screenshot:
- no

## Editor Shell Sections

### 10. Toolbar

Purpose:
- expose the main editing and demo actions

Buttons visible in the screenshot:
- `Search`
- `Bold`
- `Italic`
- `Underline`
- `Preview`
- `Split View`
- `Import`
- `Export / Exchange`
- `Snapshot`
- `Bookmark`
- `Comment`
- `Save Version`
- `Summary`
- `Focus`
- `Inspect`

Visible in screenshot:
- yes

### 11. Toolbar Status Cluster

Purpose:
- expose current connection and local editor state

Readouts:
- `connecting...`
- `auto-save idle`
- `messages: -`
- `local replica: -`

Visible in screenshot:
- yes

### 12. Quick Start Banner

Purpose:
- orient a new user to the expected demo flow

Text:
- `Open a shared doc, try Preview, then use Export / Exchange or Snapshot when you want a stable handoff or demo artifact.`

Buttons:
- `Settings`
- `Dismiss`

Visible in screenshot:
- yes

### 13. Document Tab Bar

Purpose:
- show currently open document tabs

Visible tab in screenshot:
- `Document demo`
- secondary slug text: `demo`

Visible in screenshot:
- yes

### 14. Editor Presence Row

Purpose:
- show local and remote participant badges above the editor

Visible badge in screenshot:
- `Browser User`
- `you`

Visible in screenshot:
- yes

### 15. Editor Pane

Purpose:
- host the live collaborative editor itself

Visible in screenshot:
- yes

Current screenshot state:
- empty white editor region with no visible document text

### 16. Preview Pane

Purpose:
- show rendered preview of the same document

State:
- hidden by default
- shown when `Preview` or `Split View` is used

Visible in screenshot:
- no

## Hidden Overlays And Panels

These are part of the current UI even though they are not visible in the captured screenshot.

### 17. Settings Panel

Contains:
- theme selector
- line numbers toggle
- font size range
- dyslexia-friendly spacing toggle
- presence profile selector
- shortcut bindings

### 18. Help Panel

Contains:
- keyboard shortcut/help grid

### 19. Search Panel

Contains:
- find
- replace
- case sensitive toggle
- regex toggle
- `Find Next`
- `Replace All`
- `Go To Line`

### 20. Export / Exchange Panel

Contains:
- export buttons for Markdown, HTML, Plain Text, and Automerge
- copy buttons
- publish/import exchange buttons
- audit report export

### 21. Comments Panel

Contains:
- selected text
- comment body
- save/resolve actions
- comment list

### 22. Document Summary Panel

Contains:
- generated summary text
- `Read Aloud`
- `Voice Input`

### 23. PromiseGrid Inspector Panel

Purpose:
- show expanded debug/inspection output, including clicked message details

Why it exists:
- it lets the demo explain a real selected PromiseGrid message instead of
  asking viewers to trust that “something decentralized” is happening
- it shows the current browser transport state and the selected relay-observed
  message in one place

What the inspector payload means:
- `documentID`
  - the active shared document
- `browser_transport`
  - the transport modes currently used by the browser page
- `browser_transport.sync`
  - live document sync transport
- `browser_transport.awareness`
  - live presence transport
- `browser_transport.relay_path`
  - explicit proof that the browser is going through the relay path
- `selected_message`
  - the clicked trace entry from the PromiseGrid Flow panel

Important `selected_message` fields:
- `offset`
  - relay log position of the message
- `envelope_cid`
  - CID of the signed outer envelope
- `protocol`
  - protocol family such as `live-awareness`
- `pcid`
  - protocol CID selecting the message semantics
- `kind`
  - message kind inside that protocol family
- `document_id`
  - document this message belongs to
- `participant_id`
  - participant that emitted the message
- `author`
  - signing key id
- `embodiment`
  - browser or Neovim-style source embodiment
- `lamport`
  - Lamport clock value carried on the message
- `received_at`
  - local relay receive timestamp
- `summary`
  - short human-readable explanation of the message
- `envelope_base64`
  - raw signed envelope bytes
- `payload_base64`
  - raw payload bytes
- `proof_algorithm`
  - signing algorithm
- `proof_key_id`
  - proof key id
- `decoded_payload`
  - decoded protocol payload fields for human reading

### 24. Hidden File Import Control

Purpose:
- support importing `.md`, `.txt`, `.html`, `.json`, and images

### 25. Toast Stack

Purpose:
- show transient status/error notifications

## Demo-Relevant Reading Of The Current UI

The current page is designed around three parallel stories:

1. **Editing story**
   - toolbar
   - editor pane
   - preview pane
   - comments
   - summary

2. **Collaboration story**
   - document card
   - peers card
   - presence row
   - relay card

3. **Presentation / explanation story**
   - PromiseGrid Flow card
   - PromiseGrid Inspector
   - Review card
   - export/exchange actions

The screenshot captured only the top of that experience, which is why the visible page mostly emphasizes the editor shell and the `Document` / `You` cards rather than the lower PromiseGrid-specific explainer cards.
