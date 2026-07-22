package web

import (
	"bytes"
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
		writeJSON(writer, `{"responsibilities":1,"places":1,"resources":1,"procedures":1,"training_items":0,"maintenance_items":0,"receiving_items":0,"inventory_items":0,"procedure_runs":1,"training_runs":0,"maintenance_runs":0,"receiving_runs":0,"inventory_runs":0,"approvals":2,"evidence":1,"links":1}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":0,"place_groups":[],"resource_groups":[]}`)
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

func TestHeadlessBrowserSurvivesRestrictedParticipantIdentityStartup(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		MustRead("index.html"),
		[]byte(`<script src="/app.js" type="module"></script>`),
		[]byte(`<script>
Object.defineProperty(window, "localStorage", {
  configurable: true,
  get() {
    throw new Error("storage blocked");
  },
});
if (window.crypto) {
  try {
    Object.defineProperty(window.crypto, "randomUUID", {
      configurable: true,
      value: undefined,
    });
  } catch (error) {
    window.__random_uuid_patch_error = String(error);
  }
}
</script>
<script src="/app.js" type="module"></script>`),
		1,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(rootHTML)
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
		writeJSON(writer, `{"responsibilities":0,"places":0,"resources":0,"procedures":1,"training_items":0,"maintenance_items":0,"receiving_items":0,"inventory_items":0,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"receiving_runs":0,"inventory_runs":0,"approvals":0,"evidence":0,"links":0}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":0,"place_groups":[],"resource_groups":[]}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"places":[]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Startup checklist","summary":"Boot line","current_revision":1,"working_version":1}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"item_id":"ITEM-0001","title":"Startup checklist","status":"draft","body":"# Startup checklist","version":1,"current_revision":1,"participants":[]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Startup checklist","summary":"Boot line","current_revision":1,"approvals":[],"revisions":[{"number":1,"title":"Startup checklist","author":"alice","created_at":"2026-07-20T15:59:00Z"}],"responsibility_ids":[],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T15:59:00Z","actor":"alice","title":"Startup checklist"}]}`)
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
		"Startup checklist",
		"Live Draft Studio",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
	if strings.Contains(dom, "storage blocked") && strings.Contains(dom, "Error: storage blocked") {
		t.Fatalf("restricted-storage startup leaked blocking error into DOM\n%s", dom)
	}
}

func TestHeadlessBrowserHandlesSearchFailureInApp(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		MustRead("index.html"),
		[]byte("</body>"),
		[]byte(`<script>
const searchTimer = setInterval(() => {
  const form = document.getElementById("search-form");
  if (!form) {
    return;
  }
  form.q.value = "broken search";
  form.requestSubmit();
  clearInterval(searchTimer);
}, 200);
</script></body>`),
		1,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(rootHTML)
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
		writeJSON(writer, `{"responsibilities":0,"places":0,"resources":0,"procedures":0,"training_items":0,"maintenance_items":0,"receiving_items":0,"inventory_items":0,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"receiving_runs":0,"inventory_runs":0,"approvals":0,"evidence":0,"links":0}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":0,"place_groups":[],"resource_groups":[]}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"places":[]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[]}`)
	})
	mux.HandleFunc("/api/search", func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "search backend unavailable", http.StatusInternalServerError)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	userDataDir := filepath.Join(t.TempDir(), "chrome-profile")
	command := exec.Command(
		chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--virtual-time-budget=2000",
		"--user-data-dir="+userDataDir,
		"--dump-dom",
		server.URL+"/",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome dump dom: %v\n%s", err, string(output))
	}
	dom := string(output)
	if !strings.Contains(dom, "search backend unavailable") {
		t.Fatalf("rendered dom missing search failure text\n%s", dom)
	}
}

func TestHeadlessBrowserHandlesPlaceCreateFailureInApp(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		MustRead("index.html"),
		[]byte("</body>"),
		[]byte(`<script>
const placeTimer = setInterval(() => {
  const form = document.getElementById("place-form");
  if (!form) {
    return;
  }
  form.actor.value = "alice";
  form.kind.value = "area";
  form.name.value = "Receiving";
  form.summary.value = "Inbound area";
  form.requestSubmit();
  clearInterval(placeTimer);
}, 200);
</script></body>`),
		1,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(rootHTML)
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
		writeJSON(writer, `{"responsibilities":0,"places":0,"resources":0,"procedures":0,"training_items":0,"maintenance_items":0,"receiving_items":0,"inventory_items":0,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"receiving_runs":0,"inventory_runs":0,"approvals":0,"evidence":0,"links":0}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":0,"place_groups":[],"resource_groups":[]}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			http.Error(writer, "place validation failed", http.StatusBadRequest)
			return
		}
		writeJSON(writer, `{"places":[]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[]}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	userDataDir := filepath.Join(t.TempDir(), "chrome-profile")
	command := exec.Command(
		chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--virtual-time-budget=2000",
		"--user-data-dir="+userDataDir,
		"--dump-dom",
		server.URL+"/",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome dump dom: %v\n%s", err, string(output))
	}
	dom := string(output)
	if !strings.Contains(dom, "place validation failed") {
		t.Fatalf("rendered dom missing place failure text\n%s", dom)
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
		writeJSON(writer, `{"responsibilities":1,"places":1,"resources":1,"procedures":0,"training_items":0,"maintenance_items":0,"receiving_items":0,"inventory_items":1,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"receiving_runs":0,"inventory_runs":1,"approvals":1,"evidence":1,"links":0}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":0,"place_groups":[],"resource_groups":[]}`)
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
		"Inventory count history",
		"RUN-0001",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func TestHeadlessBrowserRendersReceivingCheckReview(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		MustRead("index.html"),
		[]byte("</body>"),
		[]byte(`<script>
const runClickTimer = setInterval(() => {
  const runButton = document.querySelector("#run-list button");
  if (!runButton) {
    return;
  }
  runButton.click();
  clearInterval(runClickTimer);
}, 200);
</script></body>`),
		1,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(rootHTML)
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
		writeJSON(writer, `{"responsibilities":1,"places":1,"resources":1,"procedures":0,"training_items":0,"maintenance_items":0,"receiving_items":1,"inventory_items":0,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"receiving_runs":1,"inventory_runs":0,"approvals":1,"evidence":1,"links":0}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":1,"place_groups":[{"group_type":"place","group_id":"PLACE-0001","kind":"dock","name":"Dock A","problem_count":1,"receiving_problems":1,"inventory_problems":0,"highlights":["outcome: accepted_with_notes"],"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"}]}],"resource_groups":[]}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"places":[{"id":"PLACE-0001","kind":"dock","name":"Dock A","summary":"Inbound receiving dock","parent_id":"","child_place_ids":[],"resource_ids":["RES-0001"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z"}],"timeline":[{"type":"place_created","timestamp":"2026-07-20T17:00:00Z","actor":"alice","summary":"Inbound receiving dock"}]}]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[{"id":"RES-0001","kind":"container","name":"Inbound pallet","summary":"Pallet staged for receipt","place_id":"PLACE-0001","related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z"}],"tags":["inbound"],"links":[],"timeline":[{"type":"resource_created","timestamp":"2026-07-20T17:01:00Z","actor":"alice","summary":"Pallet staged for receipt"}]}]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[{"id":"RESP-0001","title":"Receiving lead","summary":"Owns intake checks","team":"OPS","linked_item_ids":["RECV-0001"],"linked_run_ids":["RUN-0001"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z"}],"linked_role_keys":["reviewer"],"timeline":[{"type":"responsibility_created","timestamp":"2026-07-20T17:02:00Z","actor":"alice","summary":"Owns intake checks"}]}]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[{"id":"RECV-0001","kind":"receiving_check","status":"approved","title":"Inspect inbound pallet","summary":"Receiving check for inbound pallet","current_revision":1,"working_version":2}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"}]}`)
	})
	mux.HandleFunc("/api/runs/RUN-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"responsibility_ids":["RESP-0001"],"notes":"Outer wrap torn","evidence":[{"id":"EVID-0001","summary":"Receiving inspection","facts":{"supplier":"Acme Parts","packing_slip":"PS-1234","received_units":"18","expected_units":"20","variance":"-2","condition":"wrap torn"},"actor":"bob","created_at":"2026-07-20T17:11:00Z"}],"approvals":[{"decision":"approved","role":"reviewer","actor":"boss","notes":"Reviewed at dock","created_at":"2026-07-20T17:12:00Z"}],"timeline":[{"type":"run_recorded","timestamp":"2026-07-20T17:10:00Z","actor":"bob","outcome":"accepted_with_notes","notes":"Outer wrap torn"},{"type":"evidence_added","timestamp":"2026-07-20T17:11:00Z","actor":"bob","summary":"Receiving inspection"}]}`)
	})
	mux.HandleFunc("/api/items/RECV-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"item_id":"RECV-0001","title":"Inspect inbound pallet","status":"approved","body":"# Inspect inbound pallet","version":2,"current_revision":1,"participants":[]}`)
	})
	mux.HandleFunc("/api/items/RECV-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"RECV-0001","kind":"receiving_check","status":"approved","title":"Inspect inbound pallet","summary":"Receiving check for inbound pallet","current_revision":1,"approvals":[{"decision":"approved","role":"reviewer","actor":"boss","revision":1,"notes":"Ready for intake use"}],"revisions":[{"number":1,"title":"Inspect inbound pallet","author":"alice","created_at":"2026-07-20T16:59:00Z"}],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z"}],"responsibility_ids":["RESP-0001"],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T16:59:00Z","actor":"alice","title":"Inspect inbound pallet"}]}`)
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
		"Inspect inbound pallet",
		"Receiving review",
		"RUN-0001",
		"supplier: Acme Parts",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func TestHeadlessBrowserRendersContextReviewFacts(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		MustRead("index.html"),
		[]byte("</body>"),
		[]byte(`<script>
const placeClickTimer = setInterval(() => {
  const placeButton = document.querySelector("#place-list button");
  if (!placeButton) {
    return;
  }
  placeButton.click();
  clearInterval(placeClickTimer);
}, 200);
</script></body>`),
		1,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(rootHTML)
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
		writeJSON(writer, `{"responsibilities":1,"places":1,"resources":1,"procedures":0,"training_items":0,"maintenance_items":0,"receiving_items":1,"inventory_items":1,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"receiving_runs":1,"inventory_runs":1,"approvals":1,"evidence":2,"links":0}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":2,"place_groups":[{"group_type":"place","group_id":"PLACE-0001","kind":"area","name":"Receiving","problem_count":2,"receiving_problems":1,"inventory_problems":1,"highlights":["condition: wrap torn","discrepancy: -2"],"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"},{"id":"RUN-0002","kind":"inventory_audit","item_id":"INV-0001","revision":1,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Counted receiving bin"}]}],"resource_groups":[]}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"places":[{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area","parent_id":"","child_place_ids":[],"resource_ids":["RES-0001"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z","evidence":[{"id":"EVID-0001","summary":"Receiving inspection","facts":{"supplier":"Acme Parts","variance":"-2","condition":"wrap torn"},"actor":"bob","created_at":"2026-07-20T17:11:00Z"}]},{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z","evidence":[{"id":"EVID-0002","summary":"Cycle count","facts":{"expected_count":"12","actual_count":"10","discrepancy":"-2"},"actor":"bob","created_at":"2026-07-20T18:11:00Z"}]}],"timeline":[{"type":"place_created","timestamp":"2026-07-20T16:00:00Z","actor":"alice","summary":"Inbound inspection area"}]}]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[{"id":"RES-0001","kind":"container","name":"RJ45 Bin","summary":"Connector bin","place_id":"PLACE-0001","related_runs":[{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z","evidence":[{"id":"EVID-0002","summary":"Cycle count","facts":{"expected_count":"12","actual_count":"10","discrepancy":"-2"},"actor":"bob","created_at":"2026-07-20T18:11:00Z"}]}],"tags":["parts"],"links":[],"timeline":[{"type":"resource_created","timestamp":"2026-07-20T16:01:00Z","actor":"alice","summary":"Connector bin"}]}]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[{"id":"RESP-0001","title":"Receiving lead","summary":"Owns intake checks","team":"OPS","linked_item_ids":["RECV-0001","INV-0001"],"linked_run_ids":["RUN-0001","RUN-0002"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z","evidence":[{"id":"EVID-0001","summary":"Receiving inspection","facts":{"supplier":"Acme Parts","variance":"-2","condition":"wrap torn"},"actor":"bob","created_at":"2026-07-20T17:11:00Z"}]},{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z","evidence":[{"id":"EVID-0002","summary":"Cycle count","facts":{"expected_count":"12","actual_count":"10","discrepancy":"-2"},"actor":"bob","created_at":"2026-07-20T18:11:00Z"}]}],"linked_role_keys":["reviewer"],"timeline":[{"type":"responsibility_created","timestamp":"2026-07-20T16:02:00Z","actor":"alice","summary":"Owns intake checks"}]}]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[{"id":"RECV-0001","kind":"receiving_check","status":"approved","title":"Inspect inbound pallet","summary":"Receiving check","current_revision":1,"working_version":2},{"id":"INV-0001","kind":"inventory_audit","status":"approved","title":"Count receiving bin","summary":"Cycle count","current_revision":1,"working_version":2}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"},{"id":"RUN-0002","kind":"inventory_audit","item_id":"INV-0001","revision":1,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Counted receiving bin"}]}`)
	})
	mux.HandleFunc("/api/items/RECV-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"item_id":"RECV-0001","title":"Inspect inbound pallet","status":"approved","body":"# Inspect inbound pallet","version":2,"current_revision":1,"participants":[]}`)
	})
	mux.HandleFunc("/api/items/RECV-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"RECV-0001","kind":"receiving_check","status":"approved","title":"Inspect inbound pallet","summary":"Receiving check","current_revision":1,"approvals":[],"revisions":[{"number":1,"title":"Inspect inbound pallet","author":"alice","created_at":"2026-07-20T16:59:00Z"}],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z"}],"responsibility_ids":["RESP-0001"],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T16:59:00Z","actor":"alice","title":"Inspect inbound pallet"}]}`)
	})
	mux.HandleFunc("/api/items/INV-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"item_id":"INV-0001","title":"Count receiving bin","status":"approved","body":"# Count receiving bin","version":2,"current_revision":1,"participants":[]}`)
	})
	mux.HandleFunc("/api/items/INV-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"INV-0001","kind":"inventory_audit","status":"approved","title":"Count receiving bin","summary":"Cycle count","current_revision":1,"approvals":[],"revisions":[{"number":1,"title":"Count receiving bin","author":"alice","created_at":"2026-07-20T16:59:00Z"}],"related_runs":[{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z"}],"responsibility_ids":["RESP-0001"],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T16:59:00Z","actor":"alice","title":"Count receiving bin"}]}`)
	})
	mux.HandleFunc("/api/places/PLACE-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area","parent_id":"","child_place_ids":[],"resource_ids":["RES-0001"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z","evidence":[{"id":"EVID-0001","summary":"Receiving inspection","facts":{"supplier":"Acme Parts","variance":"-2","condition":"wrap torn"},"actor":"bob","created_at":"2026-07-20T17:11:00Z"}]},{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z","evidence":[{"id":"EVID-0002","summary":"Cycle count","facts":{"expected_count":"12","actual_count":"10","discrepancy":"-2"},"actor":"bob","created_at":"2026-07-20T18:11:00Z"}]}],"timeline":[{"type":"place_created","timestamp":"2026-07-20T16:00:00Z","actor":"alice","summary":"Inbound inspection area"}]}`)
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
		"Receiving context review",
		"supplier: Acme Parts",
		"Inventory count history",
		"expected_count: 12",
		"discrepancy: -2",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func TestHeadlessBrowserSupportsContextHistoryDrilldownFilters(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		MustRead("index.html"),
		[]byte("</body>"),
		[]byte(`<script>
const placeClickTimer = setInterval(() => {
  const placeButton = document.querySelector("#place-list button");
  if (!placeButton) {
    return;
  }
  placeButton.click();
  clearInterval(placeClickTimer);
}, 200);
const searchClickTimer = setInterval(() => {
  const buttons = Array.from(document.querySelectorAll("#detail-actions button"));
  const target = buttons.find((button) => button.textContent.includes("Search problems here"));
  if (!target) {
    return;
  }
  target.click();
  clearInterval(searchClickTimer);
}, 700);
</script></body>`),
		1,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(rootHTML)
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
		writeJSON(writer, `{"responsibilities":1,"places":1,"resources":1,"procedures":0,"training_items":0,"maintenance_items":0,"receiving_items":1,"inventory_items":1,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"receiving_runs":1,"inventory_runs":1,"approvals":1,"evidence":2,"links":0}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":2,"place_groups":[{"group_type":"place","group_id":"PLACE-0001","kind":"area","name":"Receiving","problem_count":2,"receiving_problems":1,"inventory_problems":1,"highlights":["condition: wrap torn","discrepancy: -2","outcome: accepted_with_notes"],"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"},{"id":"RUN-0002","kind":"inventory_audit","item_id":"INV-0001","revision":1,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Counted receiving bin"}]}],"resource_groups":[]}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"places":[{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area","parent_id":"","child_place_ids":[],"resource_ids":["RES-0001"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z","evidence":[{"id":"EVID-0001","summary":"Receiving inspection","facts":{"supplier":"Acme Parts","variance":"-2","condition":"wrap torn"},"actor":"bob","created_at":"2026-07-20T17:11:00Z"}]},{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z","evidence":[{"id":"EVID-0002","summary":"Cycle count","facts":{"expected_count":"12","actual_count":"10","discrepancy":"-2"},"actor":"bob","created_at":"2026-07-20T18:11:00Z"}]}],"timeline":[{"type":"place_created","timestamp":"2026-07-20T16:00:00Z","actor":"alice","summary":"Inbound inspection area"}]}]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[{"id":"RES-0001","kind":"container","name":"RJ45 Bin","summary":"Connector bin","place_id":"PLACE-0001","related_runs":[{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z","evidence":[{"id":"EVID-0002","summary":"Cycle count","facts":{"expected_count":"12","actual_count":"10","discrepancy":"-2"},"actor":"bob","created_at":"2026-07-20T18:11:00Z"}]}],"tags":["parts"],"links":[],"timeline":[{"type":"resource_created","timestamp":"2026-07-20T16:01:00Z","actor":"alice","summary":"Connector bin"}]}]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[{"id":"RESP-0001","title":"Receiving lead","summary":"Owns intake checks","team":"OPS","linked_item_ids":["RECV-0001","INV-0001"],"linked_run_ids":["RUN-0001","RUN-0002"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z","evidence":[{"id":"EVID-0001","summary":"Receiving inspection","facts":{"supplier":"Acme Parts","variance":"-2","condition":"wrap torn"},"actor":"bob","created_at":"2026-07-20T17:11:00Z"}]},{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z","evidence":[{"id":"EVID-0002","summary":"Cycle count","facts":{"expected_count":"12","actual_count":"10","discrepancy":"-2"},"actor":"bob","created_at":"2026-07-20T18:11:00Z"}]}],"linked_role_keys":["reviewer"],"timeline":[{"type":"responsibility_created","timestamp":"2026-07-20T16:02:00Z","actor":"alice","summary":"Owns intake checks"}]}]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[{"id":"RECV-0001","kind":"receiving_check","status":"approved","title":"Inspect inbound pallet","summary":"Receiving check","current_revision":1,"working_version":2},{"id":"INV-0001","kind":"inventory_audit","status":"approved","title":"Count receiving bin","summary":"Cycle count","current_revision":1,"working_version":2}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"},{"id":"RUN-0002","kind":"inventory_audit","item_id":"INV-0001","revision":1,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Counted receiving bin"}]}`)
	})
	mux.HandleFunc("/api/search", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"filters":{"query":"","kind":"","status":"","outcome":"","place_id":"PLACE-0001","resource_id":"","responsibility_id":"","problem":true},"places":[{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area"}],"resources":[],"responsibilities":[],"items":[],"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"},{"id":"RUN-0002","kind":"inventory_audit","item_id":"INV-0001","revision":1,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Counted receiving bin"}]}`)
	})
	mux.HandleFunc("/api/items/RECV-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"item_id":"RECV-0001","title":"Inspect inbound pallet","status":"approved","body":"# Inspect inbound pallet","version":2,"current_revision":1,"participants":[]}`)
	})
	mux.HandleFunc("/api/items/INV-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"item_id":"INV-0001","title":"Count receiving bin","status":"approved","body":"# Count receiving bin","version":2,"current_revision":1,"participants":[]}`)
	})
	mux.HandleFunc("/api/places/PLACE-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area","parent_id":"","child_place_ids":[],"resource_ids":["RES-0001"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z","evidence":[{"id":"EVID-0001","summary":"Receiving inspection","facts":{"supplier":"Acme Parts","variance":"-2","condition":"wrap torn"},"actor":"bob","created_at":"2026-07-20T17:11:00Z"}]},{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z","evidence":[{"id":"EVID-0002","summary":"Cycle count","facts":{"expected_count":"12","actual_count":"10","discrepancy":"-2"},"actor":"bob","created_at":"2026-07-20T18:11:00Z"}]}],"timeline":[{"type":"place_created","timestamp":"2026-07-20T16:00:00Z","actor":"alice","summary":"Inbound inspection area"}]}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	userDataDir := filepath.Join(t.TempDir(), "chrome-profile")
	command := exec.Command(
		chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--virtual-time-budget=5000",
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
		"Search problems here",
		"problems only",
		"place: PLACE-0001",
		"RUN-0001",
		"RUN-0002",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func TestHeadlessBrowserRendersGroupedProblemReview(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		MustRead("index.html"),
		[]byte("</body>"),
		[]byte(`<script>
const problemClickTimer = setInterval(() => {
  const problemButton = document.querySelector("#problem-review button");
  if (!problemButton) {
    return;
  }
  problemButton.click();
  clearInterval(problemClickTimer);
}, 250);
</script></body>`),
		1,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(rootHTML)
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
		writeJSON(writer, `{"responsibilities":1,"places":1,"resources":1,"procedures":0,"training_items":0,"maintenance_items":0,"receiving_items":1,"inventory_items":1,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"receiving_runs":1,"inventory_runs":1,"approvals":1,"evidence":2,"links":0}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":2,"place_groups":[{"group_type":"place","group_id":"PLACE-0001","kind":"area","name":"Receiving","problem_count":2,"receiving_problems":1,"inventory_problems":1,"highlights":["condition: wrap torn","discrepancy: -2","outcome: accepted_with_notes"],"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"},{"id":"RUN-0002","kind":"inventory_audit","item_id":"INV-0001","revision":1,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Counted receiving bin"}]}],"resource_groups":[{"group_type":"resource","group_id":"RES-0001","kind":"container","name":"RJ45 Bin","problem_count":2,"receiving_problems":1,"inventory_problems":1,"highlights":["condition: wrap torn","discrepancy: -2"],"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"},{"id":"RUN-0002","kind":"inventory_audit","item_id":"INV-0001","revision":1,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Counted receiving bin"}]}]}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"places":[{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area","parent_id":"","child_place_ids":[],"resource_ids":["RES-0001"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z"},{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z"}],"timeline":[{"type":"place_created","timestamp":"2026-07-20T16:00:00Z","actor":"alice","summary":"Inbound inspection area"}]}]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[{"id":"RES-0001","kind":"container","name":"RJ45 Bin","summary":"Connector bin","place_id":"PLACE-0001","related_runs":[{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z"}],"tags":["parts"],"links":[],"timeline":[{"type":"resource_created","timestamp":"2026-07-20T16:01:00Z","actor":"alice","summary":"Connector bin"}]}]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[{"id":"RESP-0001","title":"Receiving lead","summary":"Owns intake checks","team":"OPS","linked_item_ids":["RECV-0001","INV-0001"],"linked_run_ids":["RUN-0001","RUN-0002"],"related_runs":[],"linked_role_keys":["reviewer"],"timeline":[{"type":"responsibility_created","timestamp":"2026-07-20T16:02:00Z","actor":"alice","summary":"Owns intake checks"}]}]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[{"id":"RECV-0001","kind":"receiving_check","status":"approved","title":"Inspect inbound pallet","summary":"Receiving check","current_revision":1,"working_version":2},{"id":"INV-0001","kind":"inventory_audit","status":"approved","title":"Count receiving bin","summary":"Cycle count","current_revision":1,"working_version":2}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"RECV-0001","revision":1,"outcome":"accepted_with_notes","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Outer wrap torn"},{"id":"RUN-0002","kind":"inventory_audit","item_id":"INV-0001","revision":1,"outcome":"completed","place_id":"PLACE-0001","resource_ids":["RES-0001"],"notes":"Counted receiving bin"}]}`)
	})
	mux.HandleFunc("/api/items/RECV-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"item_id":"RECV-0001","title":"Inspect inbound pallet","status":"approved","body":"# Inspect inbound pallet","version":2,"current_revision":1,"participants":[]}`)
	})
	mux.HandleFunc("/api/places/PLACE-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area","parent_id":"","child_place_ids":[],"resource_ids":["RES-0001"],"related_runs":[{"id":"RUN-0001","kind":"receiving_check","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z","evidence":[{"id":"EVID-0001","summary":"Receiving inspection","facts":{"supplier":"Acme Parts","variance":"-2","condition":"wrap torn"},"actor":"bob","created_at":"2026-07-20T17:11:00Z"}]},{"id":"RUN-0002","kind":"inventory_audit","revision":1,"outcome":"completed","created_at":"2026-07-20T18:10:00Z","evidence":[{"id":"EVID-0002","summary":"Cycle count","facts":{"expected_count":"12","actual_count":"10","discrepancy":"-2"},"actor":"bob","created_at":"2026-07-20T18:11:00Z"}]}],"timeline":[{"type":"place_created","timestamp":"2026-07-20T16:00:00Z","actor":"alice","summary":"Inbound inspection area"}]}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	userDataDir := filepath.Join(t.TempDir(), "chrome-profile")
	command := exec.Command(
		chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--virtual-time-budget=5000",
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
		"Problem Review",
		"Places with repeated problems",
		"2 problems",
		"condition: wrap torn",
		"discrepancy: -2",
		"PLACE-0001",
		"Receiving context review",
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
