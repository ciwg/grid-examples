package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/service"
)

const emptyReplicaBase64 = "hW9Kg8HDZmEAdQEQUDnUuZsuTLOKK6EtAqSUwAF91ThR16b5XY1P61eTHXkwnJNTicqZ35V+jMImBQWmigYBAgMCEwIjBkACVgIHFQkhAiMCNAFCAlYCgAECfwB/AX8Bf8660dIGfwB/B38HY29udGVudH8AfwEBfwR/AH8AAA=="

func TestHeadlessBrowserLateJoinRendersSharedText(t *testing.T) {
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
		`"sync_websocket_ready":true`,
		`"recovery_sync_reads":0`,
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

func TestHeadlessBrowserRecoversFromBlankSnapshotState(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := newBrowserProbeServer(t, app, browserProbeOptions{poisonBlankSnapshotState: true})
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
		`"sync_websocket_ready":true`,
		`"recovery_sync_reads":1`,
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

type browserProbeOptions struct {
	poisonBlankSnapshotState bool
}

func newBrowserProbeServer(t *testing.T, app *service.App, options ...browserProbeOptions) *httptest.Server {
	t.Helper()

	indexPath := filepath.Join(repoRoot(t), "web", "index.html")
	indexHTML, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index html: %v", err)
	}
	probeOptions := browserProbeOptions{}
	if len(options) > 0 {
		probeOptions = options[0]
	}
	requireRecovery := "false"
	if probeOptions.poisonBlankSnapshotState {
		requireRecovery = "true"
	}
	probeHTML := bytes.Replace(indexHTML, []byte("</body>"), []byte(`<div id="startup-probe" hidden></div>
<script>
// Intent: Observe local browser carriage events separately from the rendered
// transport label so this headless proof does not treat UI timing as protocol
// evidence. Source: DI-gofut; DI-raron
const transportEvidence = { sync_websocket_ready: false, recovery_sync_reads: 0 };
const nativeFetch = window.fetch.bind(window);
window.fetch = async (input, init) => {
  const requestURL = new URL(input instanceof Request ? input.url : String(input), window.location.href);
  const method = init?.method || (input instanceof Request ? input.method : "GET");
  if (method === "GET" && requestURL.pathname.endsWith("/sync")) {
    transportEvidence.recovery_sync_reads += 1;
  }
  return nativeFetch(input, init);
};
const NativeWebSocket = window.WebSocket;
window.WebSocket = function (...args) {
  const socket = new NativeWebSocket(...args);
  const socketURL = new URL(args[0], window.location.href);
  if (socketURL.pathname.endsWith("/sync-socket")) {
    socket.addEventListener("message", (event) => {
      try {
        if (JSON.parse(event.data).type === "sync-ready") {
          transportEvidence.sync_websocket_ready = true;
        }
      } catch (_error) {
        // Ignore non-JSON websocket frames; Ex3's sync-ready evidence is JSON.
      }
    });
  }
  return socket;
};
window.WebSocket.prototype = NativeWebSocket.prototype;
Object.setPrototypeOf(window.WebSocket, NativeWebSocket);
const startupProbeTimer = setInterval(() => {
  const content = document.querySelector(".cm-content");
  const target = document.getElementById("startup-probe");
  if (!content || !target) {
    return;
  }
  const text = content.textContent || "";
  target.textContent = JSON.stringify({ text, transport_evidence: transportEvidence });
  if (text.includes("Shared text should appear here.") && transportEvidence.sync_websocket_ready && (!`+requireRecovery+` || transportEvidence.recovery_sync_reads > 0)) {
    clearInterval(startupProbeTimer);
  }
}, 50);
</script></body>`), 1)

	handler := service.NewServer(app).Handler()
	server := &httptest.Server{
		Listener: listenTCP4OrSkip(t),
		Config: &http.Server{
			Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.URL.Path == "/" {
					writer.Header().Set("Content-Type", "text/html; charset=utf-8")
					_, _ = writer.Write(probeHTML)
					return
				}
				if probeOptions.poisonBlankSnapshotState && request.URL.Path == "/api/local/documents/demo/state" {
					state := app.State("demo")
					state.ReplicaBase64 = emptyReplicaBase64
					state.TextBase64 = ""
					state.SnapshotPresent = true
					state.SnapshotOffset = state.NextOffset + 99
					writer.Header().Set("Content-Type", "application/json")
					if err := json.NewEncoder(writer).Encode(state); err != nil {
						t.Fatalf("encode poisoned state: %v", err)
					}
					return
				}
				handler.ServeHTTP(writer, request)
			}),
		},
	}
	server.Start()
	return server
}
