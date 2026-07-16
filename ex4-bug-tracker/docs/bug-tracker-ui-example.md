# bug-tracker UI example

This page explains the browser UI for `ex4-bug-tracker`.

The browser is a thin working surface on top of the local Go server. The server
owns issue history, identity checks, assignment/status rules, and attachment
storage. The browser owns queue/detail presentation and form interactions.
Source: `DI-dajak`; `DI-nunit`; `DI-ninuf`.

## What this screen is

- Left side: queue filters and issue list.
- Right side: issue creation or issue detail.
- Bottom of detail: comment and attachment actions.

## Page regions

### 1. App header

- `Bug Tracker`
  - the app name
- `User`
  - identity picker for `reporter`, `triage`, and `engineer`
- `New Issue`
  - opens the issue creation panel

### 2. Queue filters

- `Status`
  - filters the queue by one fixed workflow status
- `Assignee`
  - filters the queue by engineer assignment

These filters are projections over current issue state, not separate stored
objects. Source: `DI-nunit`.

### 3. Queue list

Each issue row shows:

- issue ID
- title
- current status
- severity
- current assignee or `unassigned`

The queue sorts by most recently updated issue first.

### 4. Create issue panel

- `Title`
  - short issue summary
- `Severity`
  - one of `Low`, `Medium`, `High`, `Critical`
- `Description`
  - issue body from the reporter
- `Create Issue`
  - submits a new issue as `New`
- `Cancel`
  - closes the panel without saving

Only the `reporter` identity can create new issues in the current slice.
Source: `DI-ninuf`.

### 5. Issue detail header

- issue ID
- title
- current status badge

This is the main durable summary for one issue.

### 6. Summary grid

- `Severity`
- `Reporter`
- `Assignee`
- `Updated`

The hidden `team` field is not shown in the current UI even though the server
stores it on every issue. Source: `DI-gofub`.

### 7. Description block

Shows the original issue description exactly as stored on creation.

### 8. Assignment control

- `Assign`
  - choose the current engineer assignee
- `Update`
  - writes the assignment change

This control is only visible for the `triage` identity. Source: `DI-ninuf`.

### 9. Status control

- `Status`
  - choose a workflow state
- `Apply`
  - writes the status change

The server still enforces transition rules. The browser only exposes the form.
Source: `DI-ninuf`; `DI-gofub`.

### 10. Timeline

The timeline merges all durable issue activity into one chronological list:

- creation
- comments
- assignment changes
- status changes
- attachment uploads

This is the central history surface for the app. Source: `DI-nunit`.

### 11. Add comment

- text area for a new issue comment
- `Post Comment`
  - appends a comment event

Comments are stored as issue timeline events, not a separate chat system.
Source: `DI-nunit`.

### 12. Upload attachment

- file picker
- `Upload`
  - uploads a file into app-managed storage and appends an attachment event

Attachment links in the timeline download through the server. Source:
`DI-nunit`.

## Current browser caveats

- The UI assumes a small built-in identity set.
- Team is stored but intentionally hidden.
- There is no live presence or websocket collaboration layer.
- Timeline rendering is intentionally plain and inspectable.

## Seeded demo walkthrough

If you start the app with `bash scripts/run-demo.sh`, the seeded runtime gives
you four useful queue states immediately:

- `BUG-0001`
  - `New`
  - includes a starter attachment
- `BUG-0002`
  - `In Progress`
  - shows assignment plus engineer comments
- `BUG-0003`
  - `Resolved`
  - shows a finished issue
- `BUG-0004`
  - `Triaged`
  - shows a resolved issue that has been reopened and had its assignee cleared

That makes it easy to demonstrate the queue, timeline, attachments, and reopen
flow without manual setup. Source: `DI-zogof`.
