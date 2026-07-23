package web

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"io"
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
	addBrowserMetaHandler(mux)
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(withMockBrowserBridge(MustRead("index.html")))
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
		"Review",
		`id="mode-review" class="mode-pill is-active"`,
		`class="workspace workspace-review is-active"`,
		"Primary Flow",
		"Review draft items",
		"Review Queue",
		"Draft queue",
		"Problem hotspots",
		"Known record search",
		"Draft procedures",
		"Advanced filters",
		"Author",
		"Collaboration settings",
		"Revision decisions",
		"Writing context",
		"Writing Surface",
		"Operate",
		"Operate From Current Record",
		"Log work",
		"Attach evidence",
		"Review record",
		"Log work from current record",
		"Review this item",
		"Create",
		"Browse Collections",
		"Receiving checklist",
		"Current Record",
		"Next Step",
		"Record run for this item",
		"Continue draft",
		"Approve this item",
		"Focus Writing",
		"Capture Review Decision",
		"Log Work Performed",
		"Revisions",
		"Ready to use",
		"Receiving lead",
		"Debug payload",
		"run-item-select",
		"approval-target-select",
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
		withMockBrowserBridge(MustRead("index.html")),
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

func TestHeadlessBrowserFailsClosedWhenDirectBridgeIsUnavailable(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	mux := http.NewServeMux()
	addBrowserMetaHandler(mux)
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
		"Direct browser embodiment unavailable",
		"This embodiment currently requires Chrome or Chromium with the ex5 browser extension installed.",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}
</script>
<script src="/app.js" type="module"></script>`),
		1,
	)

	mux := http.NewServeMux()
	addBrowserMetaHandler(mux)
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
		withMockBrowserBridge(MustRead("index.html")),
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
	addBrowserMetaHandler(mux)
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
		withMockBrowserBridge(MustRead("index.html")),
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
	addBrowserMetaHandler(mux)
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
	addBrowserMetaHandler(mux)
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(withMockBrowserBridge(MustRead("index.html")))
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
		withMockBrowserBridge(MustRead("index.html")),
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
	addBrowserMetaHandler(mux)
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
		withMockBrowserBridge(MustRead("index.html")),
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
	addBrowserMetaHandler(mux)
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
		withMockBrowserBridge(MustRead("index.html")),
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
	addBrowserMetaHandler(mux)
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
		"problem-focused",
		"at PLACE-0001",
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
		withMockBrowserBridge(MustRead("index.html")),
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
	addBrowserMetaHandler(mux)
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
	mux.HandleFunc("/api/items/RECV-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"RECV-0001","kind":"receiving_check","status":"approved","title":"Inspect inbound pallet","summary":"Receiving check","current_revision":1,"approvals":[{"decision":"approved","role":"reviewer","actor":"boss","revision":1,"notes":"Ready to use"}],"revisions":[{"number":1,"title":"Inspect inbound pallet","author":"alice","created_at":"2026-07-20T16:03:00Z"}],"responsibility_ids":["RESP-0001"],"related_runs":[{"id":"RUN-0001","revision":1,"outcome":"accepted_with_notes","created_at":"2026-07-20T17:10:00Z","kind":"receiving_check"}],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T16:03:00Z","actor":"alice","title":"Inspect inbound pallet"}]}`)
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
		"Problem hotspots",
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

func TestHeadlessBrowserSearchesByRecordID(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		withMockBrowserBridge(MustRead("index.html")),
		[]byte("</body>"),
		[]byte(`<script>
const searchByIDTimer = setInterval(() => {
  const form = document.getElementById("search-form");
  if (!form) {
    return;
  }
  form.q.value = "ITEM-0001";
  form.requestSubmit();
  clearInterval(searchByIDTimer);
}, 200);
</script></body>`),
		1,
	)

	var requestedQuery string
	mux := http.NewServeMux()
	addBrowserMetaHandler(mux)
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
		writeJSON(writer, `{"responsibilities":1,"places":1,"resources":1,"procedures":1,"training_items":0,"maintenance_items":0,"receiving_items":0,"inventory_items":0,"procedure_runs":0,"training_runs":0,"maintenance_runs":0,"receiving_runs":0,"inventory_runs":0,"approvals":1,"evidence":0,"links":0}`)
	})
	mux.HandleFunc("/api/problem-review", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"problem_runs":0,"place_groups":[],"resource_groups":[]}`)
	})
	mux.HandleFunc("/api/places", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"places":[{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound inspection area","parent_id":"","child_place_ids":[],"resource_ids":[],"timeline":[]}]}`)
	})
	mux.HandleFunc("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"resources":[{"id":"RES-0001","kind":"container","name":"RJ45 Bin","summary":"Connector bin","place_id":"PLACE-0001","tags":[],"links":[],"timeline":[]}]}`)
	})
	mux.HandleFunc("/api/responsibilities", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"responsibilities":[{"id":"RESP-0001","title":"Receiving lead","summary":"Owns receiving checks","team":"OPS","linked_item_ids":["ITEM-0001"],"linked_run_ids":[],"linked_role_keys":["reviewer"],"timeline":[]}]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Receiving checklist","summary":"Procedure draft","current_revision":2,"working_version":2}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Receiving checklist","summary":"Procedure draft","current_revision":2,"approvals":[],"revisions":[{"number":1,"title":"Receiving checklist","author":"alice","created_at":"2026-07-20T15:59:00Z"},{"number":2,"title":"Receiving checklist","author":"alice","created_at":"2026-07-20T16:03:00Z"}],"responsibility_ids":["RESP-0001"],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T15:59:00Z","actor":"alice","title":"Receiving checklist"}]}`)
	})
	mux.HandleFunc("/api/search", func(writer http.ResponseWriter, request *http.Request) {
		requestedQuery = request.URL.Query().Get("q")
		writeJSON(writer, `{"filters":{"query":"ITEM-0001"},"places":[],"resources":[],"responsibilities":[],"items":[{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Receiving checklist","summary":"Procedure draft","current_revision":2,"working_version":2}],"runs":[]}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	userDataDir := filepath.Join(t.TempDir(), "chrome-profile")
	command := exec.Command(
		chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--virtual-time-budget=3000",
		"--user-data-dir="+userDataDir,
		"--dump-dom",
		server.URL+"/",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome dump dom: %v\n%s", err, string(output))
	}
	if requestedQuery != "ITEM-0001" {
		t.Fatalf("expected browser search query ITEM-0001, got %q", requestedQuery)
	}
	dom := string(output)
	required := []string{
		`Showing results for searching for "ITEM-0001".`,
		"items (1)",
		"ITEM-0001",
		"Receiving checklist",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func TestHeadlessBrowserLoadsAuthoringOnlyAfterExplicitAuthorMode(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		withMockBrowserBridge(MustRead("index.html")),
		[]byte("</body>"),
		[]byte(`<script>
const authorModeTimer = setInterval(() => {
  const detail = document.getElementById("detail-meta");
  const button = document.getElementById("mode-author");
  if (!detail || !button) {
    return;
  }
  if (!detail.textContent.includes("ITEM-0001")) {
    return;
  }
  fetch("/test/activate-author", { method: "POST" }).finally(() => button.click());
  clearInterval(authorModeTimer);
}, 250);
</script></body>`),
		1,
	)

	preAuthorLiveRequests := 0
	postAuthorLiveRequests := 0
	authorActivated := false
	mux := http.NewServeMux()
	addBrowserMetaHandler(mux)
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
		writeJSON(writer, `{"items":[{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Startup checklist","summary":"Boot line","current_revision":1,"working_version":2}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Startup checklist","summary":"Boot line","current_revision":1,"approvals":[],"revisions":[{"number":1,"title":"Startup checklist","author":"alice","created_at":"2026-07-20T15:59:00Z"}],"responsibility_ids":[],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T15:59:00Z","actor":"alice","title":"Startup checklist"}]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		if authorActivated {
			postAuthorLiveRequests++
		} else {
			preAuthorLiveRequests++
		}
		writeJSON(writer, `{"item_id":"ITEM-0001","title":"Startup checklist","status":"draft","body":"# Startup checklist","version":2,"current_revision":1,"participants":[{"participant_id":"browser-a","display_name":"Alice","color":"#0c6d62","cursor":3,"head":3,"typing":true}]}`)
	})
	mux.HandleFunc("/test/activate-author", func(writer http.ResponseWriter, request *http.Request) {
		authorActivated = true
		writer.WriteHeader(http.StatusNoContent)
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
	if preAuthorLiveRequests != 0 {
		t.Fatalf("expected no live-draft loads before explicit Author activation, got %d", preAuthorLiveRequests)
	}
	if postAuthorLiveRequests == 0 {
		t.Fatalf("expected live-draft loads after explicit Author activation")
	}
	dom := string(output)
	required := []string{
		`id="mode-author" class="mode-pill is-active"`,
		"Startup checklist",
		"Alice",
		"Live Version",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func TestHeadlessBrowserRecordsRunFromCurrentItemContext(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		withMockBrowserBridge(MustRead("index.html")),
		[]byte("</body>"),
		[]byte(`<script>
const runSubmitTimer = setInterval(() => {
  const launch = document.getElementById("operate-run-current");
  const form = document.getElementById("run-form");
  if (!launch || !form || launch.disabled) {
    return;
  }
  launch.click();
  form.outcome.value = "completed";
  form.notes.value = "Logged from item context";
  form.requestSubmit();
  clearInterval(runSubmitTimer);
}, 250);
</script></body>`),
		1,
	)

	var postedBody string
	var postedRuns int
	mux := http.NewServeMux()
	addBrowserMetaHandler(mux)
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
		writeJSON(writer, `{"responsibilities":1,"places":0,"resources":0,"procedures":1,"training_items":0,"maintenance_items":0,"receiving_items":0,"inventory_items":0,"procedure_runs":1,"training_runs":0,"maintenance_runs":0,"receiving_runs":0,"inventory_runs":0,"approvals":0,"evidence":0,"links":0}`)
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
		writeJSON(writer, `{"responsibilities":[{"id":"RESP-0001","title":"Receiving lead","summary":"Owns checks","team":"OPS","linked_item_ids":["ITEM-0001"],"linked_run_ids":[],"linked_role_keys":["reviewer"],"timeline":[]}]}`)
	})
	mux.HandleFunc("/api/items", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"items":[{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Receiving checklist","summary":"Procedure draft","current_revision":2,"working_version":2}]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Receiving checklist","summary":"Procedure draft","current_revision":2,"approvals":[],"revisions":[{"number":1,"title":"Receiving checklist","author":"alice","created_at":"2026-07-20T15:59:00Z"},{"number":2,"title":"Receiving checklist","author":"alice","created_at":"2026-07-20T16:03:00Z"}],"responsibility_ids":["RESP-0001"],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T15:59:00Z","actor":"alice","title":"Receiving checklist"}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			body := new(bytes.Buffer)
			_, _ = body.ReadFrom(request.Body)
			postedBody = body.String()
			postedRuns = 2
			writeJSON(writer, `{"id":"RUN-0002","kind":"procedure","item_id":"ITEM-0001","revision":2,"outcome":"completed","notes":"Logged from item context"}`)
			return
		}
		if postedRuns == 0 {
			writeJSON(writer, `{"runs":[{"id":"RUN-0001","kind":"procedure","item_id":"ITEM-0001","revision":2,"outcome":"completed","place_id":"","resource_ids":[],"notes":"Previous run"}]}`)
			return
		}
		writeJSON(writer, `{"runs":[{"id":"RUN-0001","kind":"procedure","item_id":"ITEM-0001","revision":2,"outcome":"completed","place_id":"","resource_ids":[],"notes":"Previous run"},{"id":"RUN-0002","kind":"procedure","item_id":"ITEM-0001","revision":2,"outcome":"completed","place_id":"","resource_ids":[],"notes":"Logged from item context"}]}`)
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
	if !strings.Contains(postedBody, `"item_id":"ITEM-0001"`) || !strings.Contains(postedBody, `"revision":2`) {
		t.Fatalf("run submit did not preserve current-item defaults: %s", postedBody)
	}
	dom := string(output)
	required := []string{
		"RUN-0002",
		"Logged from item context",
		"Saved via /api/runs",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func TestHeadlessBrowserSnapshotsRevisionFromAuthorMode(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootHTML := bytes.Replace(
		withMockBrowserBridge(MustRead("index.html")),
		[]byte("</body>"),
		[]byte(`<script>
const snapshotTimer = setInterval(() => {
  const author = document.getElementById("mode-author");
  const meta = document.getElementById("editor-meta");
  const snapshot = document.getElementById("editor-snapshot");
  if (!author || !meta || !snapshot) {
    return;
  }
  if (!document.getElementById("editor-item-id").value) {
    author.click();
    return;
  }
  if (!meta.textContent.includes("live v")) {
    return;
  }
  if (!window.__oksSnapshotReadyAt) {
    window.__oksSnapshotReadyAt = Date.now() + 1000;
    return;
  }
  if (Date.now() < window.__oksSnapshotReadyAt) {
    return;
  }
  snapshot.click();
  clearInterval(snapshotTimer);
}, 250);
</script></body>`),
		1,
	)

	var livePostBody string
	var revisionPostBody string
	liveSocketBodies := make(chan string, 4)
	mux := http.NewServeMux()
	addBrowserMetaHandler(mux)
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
		writeJSON(writer, `{"items":[{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Startup checklist","summary":"Boot line","current_revision":2,"working_version":3}]}`)
	})
	mux.HandleFunc("/api/runs", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"runs":[]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			http.NotFound(writer, request)
			return
		}
		writeJSON(writer, `{"id":"ITEM-0001","kind":"procedure","status":"draft","title":"Startup checklist","summary":"Boot line","tags":["startup"],"current_revision":2,"approvals":[],"revisions":[{"number":1,"title":"Startup checklist","author":"alice","created_at":"2026-07-20T15:59:00Z"},{"number":2,"title":"Startup checklist","author":"alice","created_at":"2026-07-20T16:03:00Z"}],"responsibility_ids":[],"timeline":[{"type":"knowledge_item_created","timestamp":"2026-07-20T15:59:00Z","actor":"alice","title":"Startup checklist"}]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001/live", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			body := new(bytes.Buffer)
			_, _ = body.ReadFrom(request.Body)
			if strings.Contains(body.String(), `"update_body":true`) {
				livePostBody = body.String()
			}
			writeJSON(writer, `{"item_id":"ITEM-0001","title":"Startup checklist","status":"draft","body":"# Startup checklist","version":4,"current_revision":2,"participants":[]}`)
			return
		}
		writeJSON(writer, `{"item_id":"ITEM-0001","title":"Startup checklist","status":"draft","body":"# Startup checklist","version":3,"current_revision":2,"participants":[]}`)
	})
	mux.HandleFunc("/api/items/ITEM-0001/live/socket", func(writer http.ResponseWriter, request *http.Request) {
		socket, err := acceptTestWebSocket(writer, request)
		if err != nil {
			t.Fatalf("accept test websocket: %v", err)
		}
		defer func() {
			_ = socket.Close()
		}()
		if err := socket.WriteJSON(`{"type":"live-state","state":{"item_id":"ITEM-0001","title":"Startup checklist","status":"draft","body":"# Startup checklist","version":3,"current_revision":2,"participants":[]}}`); err != nil {
			t.Fatalf("write test websocket state: %v", err)
		}
		for {
			payload, err := socket.ReadJSON()
			if err != nil {
				return
			}
			if strings.Contains(payload, `"update_body":true`) {
				select {
				case liveSocketBodies <- payload:
				default:
				}
				return
			}
		}
	})
	mux.HandleFunc("/api/items/ITEM-0001/revisions", func(writer http.ResponseWriter, request *http.Request) {
		body := new(bytes.Buffer)
		_, _ = body.ReadFrom(request.Body)
		revisionPostBody = body.String()
		writeJSON(writer, `{"id":"ITEM-0001","current_revision":3}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	userDataDir := filepath.Join(t.TempDir(), "chrome-profile")
	command := exec.Command(
		chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--virtual-time-budget=6000",
		"--user-data-dir="+userDataDir,
		"--dump-dom",
		server.URL+"/",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome dump dom: %v\n%s", err, string(output))
	}
	if livePostBody != "" && !strings.Contains(livePostBody, `"# Startup checklist"`) {
		t.Fatalf("snapshot bridge did not carry the expected live draft body: %s", livePostBody)
	}
	if !strings.Contains(revisionPostBody, `"# Startup checklist"`) || !strings.Contains(revisionPostBody, `"tags":["startup"]`) {
		t.Fatalf("snapshot did not post the expected durable revision payload: %s", revisionPostBody)
	}
	dom := string(output)
	required := []string{
		"Snapshot created as revision 3",
		"Startup checklist",
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

type testWebSocketConn struct {
	conn  http.Hijacker
	raw   interface{ Close() error }
	read  func() (string, error)
	write func(string) error
}

func (socket *testWebSocketConn) Close() error {
	return socket.raw.Close()
}

func (socket *testWebSocketConn) ReadJSON() (string, error) {
	return socket.read()
}

func (socket *testWebSocketConn) WriteJSON(payload string) error {
	return socket.write(payload)
}

func acceptTestWebSocket(writer http.ResponseWriter, request *http.Request) (*testWebSocketConn, error) {
	hijacker, ok := writer.(http.Hijacker)
	if !ok {
		return nil, http.ErrNotSupported
	}
	conn, buffer, err := hijacker.Hijack()
	if err != nil {
		return nil, err
	}
	key := request.Header.Get("Sec-WebSocket-Key")
	sum := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	accept := base64.StdEncoding.EncodeToString(sum[:])
	response := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
	if _, err := buffer.WriteString(response); err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := buffer.Flush(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &testWebSocketConn{
		raw: conn,
		read: func() (string, error) {
			return readTestWebSocketFrame(conn)
		},
		write: func(payload string) error {
			return writeTestWebSocketFrame(conn, payload)
		},
	}, nil
}

func readTestWebSocketFrame(conn interface {
	Read([]byte) (int, error)
}) (string, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		return "", err
	}
	length := int(header[1] & 0x7f)
	if length == 126 {
		extended := make([]byte, 2)
		if _, err := io.ReadFull(conn, extended); err != nil {
			return "", err
		}
		length = int(binary.BigEndian.Uint16(extended))
	} else if length == 127 {
		extended := make([]byte, 8)
		if _, err := io.ReadFull(conn, extended); err != nil {
			return "", err
		}
		length = int(binary.BigEndian.Uint64(extended))
	}
	masked := header[1]&0x80 != 0
	maskKey := make([]byte, 4)
	if masked {
		if _, err := io.ReadFull(conn, maskKey); err != nil {
			return "", err
		}
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return "", err
	}
	if masked {
		for i := range payload {
			payload[i] ^= maskKey[i%4]
		}
	}
	return string(payload), nil
}

func writeTestWebSocketFrame(conn interface {
	Write([]byte) (int, error)
}, payload string) error {
	body := []byte(payload)
	frame := []byte{0x81}
	switch {
	case len(body) < 126:
		frame = append(frame, byte(len(body)))
	case len(body) <= 0xffff:
		frame = append(frame, 126, 0, 0)
		binary.BigEndian.PutUint16(frame[len(frame)-2:], uint16(len(body)))
	default:
		frame = append(frame, 127, 0, 0, 0, 0, 0, 0, 0, 0)
		binary.BigEndian.PutUint64(frame[len(frame)-8:], uint64(len(body)))
	}
	frame = append(frame, body...)
	_, err := conn.Write(frame)
	return err
}

func withMockBrowserBridge(html []byte) []byte {
	// Intent: Keep the browser smoke suite focused on page-level UI behavior
	// while stronger deterministic extension/native-host contract tests cover
	// the real shipped bridge boundary separately. Source: DI-vasem
	mock := `<script>
window.addEventListener("message", async (event) => {
  if (event.source !== window || !event.data || event.data.__oks_bridge !== true || event.data.direction !== "page->bridge") {
    return;
  }
  const message = event.data;
  const reply = (payload) => window.postMessage({
    __oks_bridge: true,
    direction: "bridge->page",
    request_id: message.request_id,
    ...payload,
  }, window.location.origin);
  if (message.kind === "handshake") {
    const response = await fetch("/api/meta").then((metaResponse) => metaResponse.json());
    const ok = !!(response && response.local_unix_socket_path && response.embodiments && response.embodiments.browser && response.embodiments.browser.primary_adapter === "chrome_native_messaging");
    reply({ kind: "handshake", ok });
    return;
  }
  if (message.kind === "rpc") {
    const request = message.request || {};
    let path = request.path;
    let method = request.method || "GET";
    let headers = request.headers || {};
    let body = request.body || undefined;
    if (!path && request.type === "operation") {
      if (request.operation === "inspect_item") {
        path = "/api/items/" + request.item_id;
      } else if (request.operation === "inspect_run") {
        path = "/api/runs/" + request.run_id;
      } else if (request.operation === "inspect_entity") {
        const roots = {
          place: "places",
          resource: "resources",
          responsibility: "responsibilities",
        };
        path = "/api/" + roots[request.entity_type] + "/" + request.entity_id;
      } else if (request.operation === "search") {
        const params = new URLSearchParams();
        const options = request.search_options || {};
        for (const [key, value] of Object.entries(options)) {
          if (value !== "" && value !== false) {
            params.set(key, String(value));
          }
        }
        path = "/api/search?" + params.toString();
      } else if (request.operation === "problem_review") {
        path = "/api/problem-review";
      } else if (request.operation === "create_place") {
        path = "/api/places";
        method = "POST";
        headers = { "Content-Type": "application/json" };
        body = JSON.stringify({
          actor: request.actor,
          kind: request.kind,
          name: request.name,
          summary: request.summary,
          parent_id: request.parent_id,
          tags: request.tags || [],
        });
      } else if (request.operation === "create_resource") {
        path = "/api/resources";
        method = "POST";
        headers = { "Content-Type": "application/json" };
        body = JSON.stringify({
          actor: request.actor,
          kind: request.kind,
          name: request.name,
          summary: request.summary,
          place_id: request.place_id,
          tags: request.tags || [],
        });
      } else if (request.operation === "create_responsibility") {
        path = "/api/responsibilities";
        method = "POST";
        headers = { "Content-Type": "application/json" };
        body = JSON.stringify({
          actor: request.actor,
          title: request.title,
          summary: request.summary,
          role_keys: request.role_keys || [],
          tags: request.tags || [],
        });
      } else if (request.operation === "create_item") {
        path = "/api/items";
        method = "POST";
        headers = { "Content-Type": "application/json" };
        body = JSON.stringify({
          actor: request.actor,
          kind: request.kind,
          title: request.title,
          summary: request.summary,
          body: request.body,
          tags: request.tags || [],
          responsibility_ids: request.responsibility_ids || [],
        });
      } else if (request.operation === "record_run") {
        path = "/api/runs";
        method = "POST";
        headers = { "Content-Type": "application/json" };
        body = JSON.stringify({
          actor: request.actor,
          kind: request.kind,
          item_id: request.item_id,
          revision: request.revision,
          outcome: request.outcome,
          notes: request.notes,
          machine: request.machine,
          location: request.location,
          place_id: request.place_id,
          resource_ids: request.resource_ids || [],
          responsibility_ids: request.responsibility_ids || [],
        });
      } else if (request.operation === "record_item_approval") {
        path = "/api/items/" + request.item_id + "/approvals";
        method = "POST";
        headers = { "Content-Type": "application/json" };
        body = JSON.stringify({
          actor: request.actor,
          revision: request.revision,
          role: request.role,
          decision: request.decision,
          notes: request.notes,
        });
      } else if (request.operation === "record_run_approval") {
        path = "/api/runs/" + request.run_id + "/approvals";
        method = "POST";
        headers = { "Content-Type": "application/json" };
        body = JSON.stringify({
          actor: request.actor,
          role: request.role,
          decision: request.decision,
          notes: request.notes,
        });
      } else if (request.operation === "add_revision") {
        path = "/api/items/" + request.item_id + "/revisions";
        method = "POST";
        headers = { "Content-Type": "application/json" };
        body = JSON.stringify({
          actor: request.actor,
          title: request.title,
          summary: request.summary,
          body: request.body,
          tags: request.tags || [],
        });
      } else if (request.operation === "supersede_item") {
        path = "/api/items/" + request.item_id + "/supersede";
        method = "POST";
        headers = { "Content-Type": "application/json" };
        body = JSON.stringify({
          actor: request.actor,
          notes: request.notes,
        });
      } else if (request.operation === "add_evidence") {
        path = "/api/runs/" + request.run_id + "/evidence";
        method = "POST";
        const payload = new FormData();
        payload.set("actor", request.actor || "");
        payload.set("summary", request.summary || "");
        payload.set("facts_json", JSON.stringify(request.facts || {}));
        if (request.attachment_name && request.attachment_body_base64) {
          const binary = atob(request.attachment_body_base64);
          const bytes = new Uint8Array(binary.length);
          for (let index = 0; index < binary.length; index += 1) {
            bytes[index] = binary.charCodeAt(index);
          }
          payload.set("attachment", new Blob([bytes]), request.attachment_name);
        }
        body = payload;
        headers = {};
      }
    }
    const response = await fetch(path, {
      method,
      headers,
      body,
    });
    const text = await response.text();
    reply({
      kind: "rpc-response",
      response: {
        type: "response",
        status: response.status,
        headers: { content_type: response.headers.get("Content-Type") || "" },
        body: text,
      },
    });
    return;
  }
  if (message.kind === "live-open") {
    const state = await fetch("/api/items/" + message.request.item_id + "/live").then((response) => response.text());
    reply({
      kind: "live-message",
      response: {
        type: "live-state",
        state: JSON.parse(state),
      },
    });
    return;
  }
  if (message.kind === "live-update") {
    const payload = message.request || {};
    const response = await fetch("/api/items/" + (message.request.item_id || "") + "/live", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        participant_id: payload.participant_id,
        display_name: payload.display_name,
        color: payload.color,
        cursor: payload.cursor,
        head: payload.head,
        typing: payload.typing,
        base_version: payload.base_version,
        update_body: payload.update_body,
        body: payload.body,
      }),
    });
    const text = await response.text();
    const parsed = JSON.parse(text);
    reply({
      kind: "live-message",
      response: response.status === 409 ? { type: "live-conflict", state: parsed.state } : { type: "live-state", state: parsed },
    });
  }
});
</script>
<script src="/app.js" type="module"></script>`
	return bytes.Replace(html, []byte(`<script src="/app.js" type="module"></script>`), []byte(mock), 1)
}

func addBrowserMetaHandler(mux *http.ServeMux) {
	mux.HandleFunc("/api/meta", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, `{"local_unix_socket_path":"/tmp/embodiment.sock","embodiments":{"browser":{"primary_adapter":"chrome_native_messaging","live_draft_transport":"native_messaging","compatibility_mode":"chrome_or_chromium_required"}}}`)
	})
}

func writeJSON(writer http.ResponseWriter, body string) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = writer.Write([]byte(body))
}
