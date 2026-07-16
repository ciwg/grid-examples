# TODO sudur - bug tracker review findings

## Decision Intent Log

ID: DI-jofoj
Date: 2026-07-16 01:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the senior test review findings for `ex4-bug-tracker` as one focused follow-up TODO so behavior, UI, and test gaps can be fixed systematically.
Intent: Preserve the review outcome in repo-local planning instead of leaving the bugs only in chat.
Constraints: Keep the findings concrete and implementation-ready; do not collapse distinct behavior, UI, and resilience issues into one vague task.
Affects: `ex4-bug-tracker/TODO/TODO.md`, `ex4-bug-tracker/TODO/TODO-sudur-bug-tracker-review-findings.md`

ID: DI-gitam
Date: 2026-07-16 01:12:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Enforce workflow integrity in the service layer by requiring assignment before engineer progress transitions and by restricting assignment changes to workflow states where assignment makes sense.
Intent: Make the server enforce the single-active-assignee workflow that the product and docs already present.
Constraints: Engineers may only advance issues assigned to themselves; triage may not assign `New` or `Resolved` issues; keep the existing fixed status model and reopen behavior.
Affects: `ex4-bug-tracker/service/app.go`, `ex4-bug-tracker/service/app_test.go`

ID: DI-zumog
Date: 2026-07-16 01:12:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make the browser role-aware by hiding invalid create/status actions and by preventing the create panel from visually overlapping an open issue detail view.
Intent: Stop inviting users to click controls that the server will reject and keep the main page focused on one primary task at a time.
Constraints: Keep the browser thin; mirror the existing server workflow instead of inventing a second browser-only permission model.
Affects: `ex4-bug-tracker/web/app.js`, `ex4-bug-tracker/web/index.html`

ID: DI-fakuv
Date: 2026-07-16 01:12:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add explicit browser error handling around startup and refresh fetches so the app can surface API failures and recover without leaving the page half-initialized.
Intent: Improve operator confidence and resilience without widening the app into a more complex client runtime.
Constraints: Keep the handling small and local; prefer clear error banners and stable page state over retries or background recovery loops.
Affects: `ex4-bug-tracker/web/app.js`

## Goal

Close the first serious behavior and UI bugs discovered during review of
`ex4-bug-tracker`.

## Tasks

- [x] sudur.1 Enforce assignment ownership in issue workflow so only the active assignee can move an issue into `In Progress` or `Resolved`, and unassigned issues cannot be advanced by arbitrary engineers.
- [x] sudur.2 Add assignment lifecycle guards so triage cannot assign `New` or `Resolved` issues without the required workflow state changes.
- [x] sudur.3 Make the browser role-aware by hiding or constraining controls that the current identity is not allowed to use, including create and status actions.
- [x] sudur.4 Fix the “New Issue” browser state so opening the create panel does not leave a selected issue detail panel visible underneath it.
- [x] sudur.5 Add browser resilience for failed or unavailable API calls, especially during startup and queue/detail refresh, and verify the recovery path with tests where practical.

## Evidence

- Review findings were captured from a senior QA-style pass over `ex4-bug-tracker` behavior, UI, and code.
- The findings cover workflow integrity, role-aware UI behavior, and browser failure handling.
- Verification passed with `go test ./...`, `errcheck ./...`, and `node --check web/app.js`.
