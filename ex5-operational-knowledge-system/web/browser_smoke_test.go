package web

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHeadlessBrowserRendersOperationalWorkflow(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(MustRead("index.html"))
	})
	mux.HandleFunc("/app.js", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		_, _ = writer.Write(MustRead("app.js"))
	})
	mux.HandleFunc("/style.css", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = writer.Write(MustRead("style.css"))
	})
	mux.HandleFunc("/api/dashboard", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":1,"places":1,"resources":1,"procedures":1,"training_items":0,"maintenance_items":0,"inventory_items":0,"procedure_runs":1,"training_runs":0,"maintenance_runs":0,"inventory_runs":0,"approvals":2,"evidence":1,"links":1}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"places":[{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area","parent_id":"","child_place_ids":[],"resource_ids":["RES-0001"],"timeline":[{"type":"place_created","timestamp":"2026-07-20T16:00:00Z","actor":"alice","summary":"Inbound inspection area"}]}]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[{"id":"RES-0001","kind":"container","name":"RJ45 Bin","summary":"Connector bin","place_id":"PLACE-0001","tags":["parts"],"links":[],"timeline":[{"type":"resource_created","timestamp":"2026-07-20T16:01:00Z","actor":"alice","summary":"Connector bin"}]}]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[{"id":"RESP-0001","title":"Receiving lead","summary":"Owns receiving checks","team":"OPS","linked_item_ids":["ITEM-0001"],"linked_run_ids":["RUN-0001"],"linked_role_keys":["reviewer"],"timeline":[{"type":"responsibility_created","timestamp":"2026-07-20T16:02:00Z","actor":"alice","summary":"Owns receiving checks"}]}]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[{"id":"ITEM-0001","kind":"procedure","status":"approved","title":"Receiving checklist","summary":"Procedure draft","current_revision":2,"working_version":3}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[{"id":"RUN-0001","kind":"procedure","item_id":"ITEM-0001","revision":2,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Count completed"}]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"item_id":"ITEM-0001","title":"Receiving checklist","status":"approved","body":"# Receiving checklist","version":3,"current_revision":2,"participants":[{"participant_id":"browser-a","display_name":"Alice","color":"#1d6fd6","cursor":4,"head":4,"typing":true}]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"ITEM-0001","kind":"procedure","status":"approved","title":"Receiving checklist","summary":"Procedure draft","current_revision":2,"approvals":[{"decision":"approved","role":"reviewer","actor":"boss","revision":2,"notes":"Ready to use"}],"revisions":[{"number":1,"title":"Receiving checklist","author":"alice","created_at":"2026-07-20T15:59:00Z"},{"number":2,"title":"Receiving checklist","author":"alice","created_at":"2026-07-20T16:03:00Z"}],"responsibility_ids":["RESP-0001"],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T15:59:00Z","actor":"alice","title":"Receiving checklist"},{"type":"approval_recorded","timestamp":"2026-07-20T16:04:00Z","actor":"boss","decision":"approved","revision":2,"notes":"Ready to use"}]}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	userDataDir := filepath.Join(t.TempDir(), "chrome-profile")
	command := exec.Command(
		chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--virtual-time-budget=4000",
		"--user-data-dir="+userDataDir,
		"--dump-dom",
		server.URL+"/",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome dump dom: %v\n%s", err, string(output))
	}
	dom := string(output)
	required := []string{
		"Operational Knowledge System",
		"Receiving checklist",
		"Record Inspector",
		"Revisions",
		"Ready to use",
		"Receiving lead",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func TestHeadlessBrowserRendersInventoryAuditHistory(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(MustRead("index.html"))
	})
	mux.HandleFunc("/app.js", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		_, _ = writer.Write(MustRead("app.js"))
	})
	mux.HandleFunc("/style.css", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = writer.Write(MustRead("style.css"))
	})
	mux.HandleFunc("/api/dashboard", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":1,"places":1,"resources":1,"procedures":0,"training_items":0,"maintenance_items":0,"inventory_items":1,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"inventory_runs":1,"approvals":1,"evidence":1,"links":0}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"places":[{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area","parent_id":"","child_place_ids":[],"resource_ids":["RES-0001"],"related_runs":[{"id":"RUN-0001","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T16:10:00Z"}],"timeline":[{"type":"place_created","timestamp":"2026-07-20T16:00:00Z","actor":"alice","summary":"Inbound inspection area"}]}]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[{"id":"RES-0001","kind":"container","name":"RJ45 Bin","summary":"Connector bin","place_id":"PLACE-0001","related_runs":[{"id":"RUN-0001","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T16:10:00Z"}],"tags":["parts"],"links":[],"timeline":[{"type":"resource_created","timestamp":"2026-07-20T16:01:00Z","actor":"alice","summary":"Connector bin"}]}]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[{"id":"RESP-0001","title":"Receiving lead","summary":"Owns receiving checks","team":"OPS","linked_item_ids":["INV-0001"],"linked_run_ids":["RUN-0001"],"related_runs":[{"id":"RUN-0001","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T16:10:00Z"}],"linked_role_keys":["reviewer"],"timeline":[{"type":"responsibility_created","timestamp":"2026-07-20T16:02:00Z","actor":"alice","summary":"Owns receiving checks"}]}]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[{"id":"INV-0001","kind":"inventory_audit","status":"approved","title":"Count RJ45 bin","summary":"Cycle count for receiving bin","current_revision":1,"working_version":2}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[{"id":"RUN-0001","kind":"inventory_audit","item_id":"INV-0001","revision":1,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Counted receiving bin"}]}`)
	})
	mux.HandleFunc("/api/items/INV-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"item_id":"INV-0001","title":"Count RJ45 bin","status":"approved","body":"# Count RJ45 bin","version":2,"current_revision":1,"participants":[]}`)
	})
	mux.HandleFunc("/api/items/INV-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"INV-0001","kind":"inventory_audit","status":"approved","title":"Count RJ45 bin","summary":"Cycle count for receiving bin","current_revision":1,"approvals":[{"decision":"approved","role":"reviewer","actor":"boss","revision":1,"notes":"Ready to use"}],"revisions":[{"number":1,"title":"Count RJ45 bin","author":"alice","created_at":"2026-07-20T15:59:00Z"}],"related_runs":[{"id":"RUN-0001","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T16:10:00Z"}],"responsibility_ids":["RESP-0001"],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T15:59:00Z","actor":"alice","title":"Count RJ45 bin"}]}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	userDataDir := filepath.Join(t.TempDir(), "chrome-profile")
	command := exec.Command(
		chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--virtual-time-budget=4000",
		"--user-data-dir="+userDataDir,
		"--dump-dom",
		server.URL+"/",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome dump dom: %v\n%s", err, string(output))
	}
	dom := string(output)
	required := []string{
		"Count RJ45 bin",
		"Inventory audit history",
		"RUN-0001",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func writeJSON(writer http.ResponseWriter, body string) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = writer.Write([]byte(body))
}
