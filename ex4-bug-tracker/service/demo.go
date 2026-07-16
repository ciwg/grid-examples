package service

// Intent: Keep the demo path honest by seeding a small issue set through the
// same append-only app methods the browser and CLI use, rather than inventing a
// parallel demo-only storage format or UI. Source: DI-zogof
func (app *App) SeedDemoIfEmpty() (bool, error) {
	issues, err := app.ListIssues("", "")
	if err != nil {
		return false, err
	}
	if len(issues) > 0 {
		return false, nil
	}
	if err := app.seedDemo(); err != nil {
		return false, err
	}
	return true, nil
}

func (app *App) seedDemo() error {
	issueA, err := app.CreateIssue("reporter", "Safari upload crash", "Uploading a 4 MB log file in Safari closes the page without any visible error.", SeverityHigh)
	if err != nil {
		return err
	}
	if _, err := app.AddComment("reporter", issueA.ID, "I can reproduce this every time with the latest support bundle from QA."); err != nil {
		return err
	}
	if _, err := app.AddAttachment("reporter", issueA.ID, "safari-console.txt", "text/plain", []byte("TypeError: undefined is not an object\nat upload.js:118")); err != nil {
		return err
	}

	issueB, err := app.CreateIssue("reporter", "CSV import trims account names", "Imported CSV rows lose trailing spaces and collapse some quoted account names.", SeverityMedium)
	if err != nil {
		return err
	}
	if _, err := app.ChangeStatus("triage", issueB.ID, StatusTriaged); err != nil {
		return err
	}
	if _, err := app.AssignIssue("triage", issueB.ID, "engineer"); err != nil {
		return err
	}
	if _, err := app.ChangeStatus("engineer", issueB.ID, StatusInProgress); err != nil {
		return err
	}
	if _, err := app.AddComment("engineer", issueB.ID, "Working on a parser fix. The bug looks isolated to the normalization pass."); err != nil {
		return err
	}

	issueC, err := app.CreateIssue("reporter", "Notification badge stays red after refresh", "The header badge still shows unread alerts after I reload the page.", SeverityLow)
	if err != nil {
		return err
	}
	if _, err := app.ChangeStatus("triage", issueC.ID, StatusTriaged); err != nil {
		return err
	}
	if _, err := app.AssignIssue("triage", issueC.ID, "engineer"); err != nil {
		return err
	}
	if _, err := app.ChangeStatus("engineer", issueC.ID, StatusInProgress); err != nil {
		return err
	}
	if _, err := app.AddComment("engineer", issueC.ID, "Fix landed locally. Verifying cache invalidation now."); err != nil {
		return err
	}
	if _, err := app.ChangeStatus("engineer", issueC.ID, StatusResolved); err != nil {
		return err
	}

	issueD, err := app.CreateIssue("reporter", "Password reset email shows the old branding", "The footer still uses the winter campaign logo instead of the current product branding.", SeverityLow)
	if err != nil {
		return err
	}
	if _, err := app.ChangeStatus("triage", issueD.ID, StatusTriaged); err != nil {
		return err
	}
	if _, err := app.AssignIssue("triage", issueD.ID, "engineer"); err != nil {
		return err
	}
	if _, err := app.ChangeStatus("engineer", issueD.ID, StatusInProgress); err != nil {
		return err
	}
	if _, err := app.ChangeStatus("engineer", issueD.ID, StatusResolved); err != nil {
		return err
	}
	if _, err := app.AddComment("reporter", issueD.ID, "Still seeing the old logo in staging after today's deploy."); err != nil {
		return err
	}
	if _, err := app.ChangeStatus("reporter", issueD.ID, StatusTriaged); err != nil {
		return err
	}

	return nil
}
