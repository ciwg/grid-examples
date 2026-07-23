package chromeextension

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type manifest struct {
	ManifestVersion int      `json:"manifest_version"`
	Permissions     []string `json:"permissions"`
	Background      struct {
		ServiceWorker string `json:"service_worker"`
	} `json:"background"`
	ContentScripts []struct {
		Matches []string `json:"matches"`
		JS      []string `json:"js"`
	} `json:"content_scripts"`
	Key string `json:"key"`
}

func TestChromeExtensionManifestPinsTheDirectBrowserEmbodiment(t *testing.T) {
	root := repoChromeExtensionRoot(t)
	payload, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest manifest
	if err := json.Unmarshal(payload, &manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	if manifest.ManifestVersion != 3 {
		t.Fatalf("expected mv3 manifest: %+v", manifest)
	}
	if manifest.Background.ServiceWorker != "background.js" {
		t.Fatalf("unexpected background worker: %+v", manifest)
	}
	if len(manifest.ContentScripts) != 1 || len(manifest.ContentScripts[0].JS) != 1 || manifest.ContentScripts[0].JS[0] != "content.js" {
		t.Fatalf("unexpected content script config: %+v", manifest.ContentScripts)
	}
	if !containsString(manifest.Permissions, "nativeMessaging") {
		t.Fatalf("expected nativeMessaging permission: %+v", manifest.Permissions)
	}
	if manifest.Key == "" {
		t.Fatal("expected pinned extension key")
	}
}

func TestChromeExtensionNativeHostManifestPinsTheAllowedOrigin(t *testing.T) {
	root := repoChromeExtensionRoot(t)
	payload, err := os.ReadFile(filepath.Join(root, "native-host", "operational_browser_host.json"))
	if err != nil {
		t.Fatalf("read native host manifest: %v", err)
	}
	text := string(payload)
	if !strings.Contains(text, `"operational_browser_host"`) {
		t.Fatalf("missing native host name: %s", text)
	}
	if !strings.Contains(text, `"chrome-extension://miagfmaampfgjkojhccdilogehbjijpe/"`) {
		t.Fatalf("missing fixed allowed origin: %s", text)
	}
	if !strings.Contains(text, `"__BROWSER_HOST_PATH__"`) {
		t.Fatalf("missing host path placeholder: %s", text)
	}
}

func repoChromeExtensionRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return wd
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
