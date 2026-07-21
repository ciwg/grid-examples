package service_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/service"
)

func TestHeadlessBrowserLateJoinRendersSharedText(t *testing.T) {
	t.Parallel()

	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := newBrowserProbeServer(t, app)
	defer server.Close()

	repo := repoRoot(t)
	browserA := startJSONProcess(t, repo, "node", filepath.Join(repo, "service", "testdata", "browser-harness.mjs"))
	defer browserA.Close()
	browserA.WaitForType(t, "ready")
	browserA.Send(t, map[string]any{
		"type":           "connect",
		"relay_url":      server.URL,
		"participant_id": "browser-a",
		"doc_id":         "demo",
		"display_name":   "Browser A",
		"color":          "#1d6fd6",
	})
	browserA.WaitForType(t, "opened")

	const shared = "# Incognito Probe\n\nShared text should appear here."
	browserA.Send(t, map[string]any{
		"type":    "set_text",
		"content": shared,
	})
	browserA.WaitForMessage(t, func(message map[string]any) bool {
		return stringField(t, message, "type") == "local_change" && stringField(t, message, "content") == shared
	})

	userDataDir := filepath.Join(t.TempDir(), "chrome-profile")
	command := exec.Command(
		chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--incognito",
		"--virtual-time-budget=6000",
		"--user-data-dir="+userDataDir,
		"--dump-dom",
		server.URL+"/?doc=demo",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome dump dom: %v\n%s", err, string(output))
	}
	dom := string(output)
	required := []string{
		"browser sync: websocket",
		"# Incognito Probe",
		"Shared text should appear here.",
		`id="startup-probe"`,
	}
	for _, marker := range required {
		if !strings.Contains(dom, marker) {
			t.Fatalf("rendered dom missing %q\n%s", marker, dom)
		}
	}
}

func newBrowserProbeServer(t *testing.T, app *service.App) *httptest.Server {
	t.Helper()

	indexPath := filepath.Join(repoRoot(t), "web", "index.html")
	indexHTML, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index html: %v", err)
	}
	probeHTML := bytes.Replace(indexHTML, []byte("</body>"), []byte(`<div id="startup-probe" hidden></div>
<script>
const startupProbeTimer = setInterval(() => {
  const content = document.querySelector(".cm-content");
  const target = document.getElementById("startup-probe");
  if (!content || !target) {
    return;
  }
  target.textContent = content.textContent || "";
  if (target.textContent.includes("Shared text should appear here.")) {
    clearInterval(startupProbeTimer);
  }
}, 50);
</script></body>`), 1)

	handler := service.NewServer(app).Handler()
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/" {
			writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = writer.Write(probeHTML)
			return
		}
		handler.ServeHTTP(writer, request)
	}))
}
