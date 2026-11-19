package service_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex3.5-grid-editor-dtls/service"
	"github.com/computerscienceiscool/grid-examples/ex3.5-grid-editor-dtls/transport"
)

const emptyReplicaBase64 = "hW9Kg8HDZmEAdQEQUDnUuZsuTLOKK6EtAqSUwAF91ThR16b5XY1P61eTHXkwnJNTicqZ35V+jMImBQWmigYBAgMCEwIjBkACVgIHFQkhAiMCNAFCAlYCgAECfwB/AX8Bf8660dIGfwB/B38HY29udGVudH8AfwEBfwR/AH8AAA=="

func TestHeadlessBrowserFreshDocumentCompletesStartup(t *testing.T) {
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
		server.URL+"/?doc=fresh-empty-document",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome dump dom: %v\n%s", err, string(output))
	}
	dom := string(output)
	for _, marker := range []string{
		"browser sync: webrtc",
		`id="status" class="status-pill status-ready"`,
		`data-language="markdown"`,
	} {
		if !strings.Contains(dom, marker) {
			t.Fatalf("fresh document startup missing %q\n%s", marker, dom)
		}
	}
}

func TestHeadlessBrowserEditorEditPersistsDocument(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	const documentID = "browser-editor-proof"
	const editedText = "# Browser editor proof\n\nThis text came from CodeMirror."
	server := newBrowserProbeServer(t, app, browserProbeOptions{editorText: editedText})
	defer server.Close()

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
		server.URL+"/?doc="+documentID,
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome edit document: %v\n%s", err, string(output))
	}
	if !strings.Contains(string(output), "This text came from CodeMirror.") {
		t.Fatalf("rendered CodeMirror document missing edit\n%s", string(output))
	}

	state := app.State(documentID)
	persisted, err := base64.StdEncoding.DecodeString(state.TextBase64)
	if err != nil {
		t.Fatalf("decode persisted document: %v", err)
	}
	if string(persisted) != editedText {
		t.Fatalf("persisted document = %q, want %q", string(persisted), editedText)
	}
}

func TestHeadlessBrowserEditorsConvergeThroughDirectDTLS(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootA := filepath.Join(t.TempDir(), "node-a")
	rootB := filepath.Join(t.TempDir(), "node-b")
	identityA, err := transport.LoadOrCreateDTLSIdentity(filepath.Join(rootA, "dtls", "certificate.pem"), filepath.Join(rootA, "dtls", "private-key.pem"))
	if err != nil {
		t.Fatalf("load node a identity: %v", err)
	}
	identityB, err := transport.LoadOrCreateDTLSIdentity(filepath.Join(rootB, "dtls", "certificate.pem"), filepath.Join(rootB, "dtls", "private-key.pem"))
	if err != nil {
		t.Fatalf("load node b identity: %v", err)
	}
	appA, err := service.NewApp(rootA)
	if err != nil {
		t.Fatalf("new node a: %v", err)
	}
	appB, err := service.NewApp(rootB)
	if err != nil {
		t.Fatalf("new node b: %v", err)
	}

	endpointA := freeUDPAddress(t)
	endpointB := freeUDPAddress(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := appB.StartDTLSPeerDelivery(ctx, service.DTLSPeerConfig{
		ListenEndpoint: endpointB,
		Peers: []service.DTLSPeerInvite{{
			Endpoint:               endpointA,
			CertificateBase64:      base64.StdEncoding.EncodeToString(identityA.DER),
			CertificateFingerprint: identityA.Fingerprint,
		}},
	}); err != nil {
		t.Fatalf("start node b dtls: %v", err)
	}
	if _, err := appA.StartDTLSPeerDelivery(ctx, service.DTLSPeerConfig{
		ListenEndpoint: endpointA,
		Peers: []service.DTLSPeerInvite{{
			Endpoint:               endpointB,
			CertificateBase64:      base64.StdEncoding.EncodeToString(identityB.DER),
			CertificateFingerprint: identityB.Fingerprint,
		}},
	}); err != nil {
		t.Fatalf("start node a dtls: %v", err)
	}

	const documentID = "real-browser-direct-dtls"
	const textFromA = "# Browser A\n\nThis edit crossed direct DTLS."
	serverA := newBrowserProbeServer(t, appA, browserProbeOptions{editorText: textFromA})
	defer serverA.Close()
	if dom := runHeadlessBrowser(t, chromePath, serverA.URL, documentID); !strings.Contains(dom, "This edit crossed direct DTLS.") {
		t.Fatalf("browser a did not render its CodeMirror edit\n%s", dom)
	}
	waitForPersistedBrowserText(t, appA, documentID, textFromA)
	waitForBrowserSyncMessageCount(t, appB, documentID, 1)

	serverB := newBrowserProbeServer(t, appB)
	defer serverB.Close()
	if dom := runHeadlessBrowser(t, chromePath, serverB.URL, documentID); !strings.Contains(dom, "This edit crossed direct DTLS.") {
		t.Fatalf("browser b did not render node a's direct-DTLS edit\n%s", dom)
	}

	const textFromB = "# Browser B\n\nThis edit returned through direct DTLS."
	serverBEdit := newBrowserProbeServer(t, appB, browserProbeOptions{editorText: textFromB})
	defer serverBEdit.Close()
	if dom := runHeadlessBrowser(t, chromePath, serverBEdit.URL, documentID); !strings.Contains(dom, "This edit returned through direct DTLS.") {
		t.Fatalf("browser b did not render its CodeMirror edit\n%s", dom)
	}
	waitForPersistedBrowserText(t, appB, documentID, textFromB)
	waitForBrowserSyncMessageCount(t, appA, documentID, 2)

	serverAFinal := newBrowserProbeServer(t, appA)
	defer serverAFinal.Close()
	if dom := runHeadlessBrowser(t, chromePath, serverAFinal.URL, documentID); !strings.Contains(dom, "This edit returned through direct DTLS.") {
		t.Fatalf("browser a did not render node b's direct-DTLS edit\n%s", dom)
	}
}

func TestHeadlessBrowserAndNeovimConvergeThroughDirectDTLS(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootA := filepath.Join(t.TempDir(), "node-a")
	rootB := filepath.Join(t.TempDir(), "node-b")
	identityA, err := transport.LoadOrCreateDTLSIdentity(filepath.Join(rootA, "dtls", "certificate.pem"), filepath.Join(rootA, "dtls", "private-key.pem"))
	if err != nil {
		t.Fatalf("load node a identity: %v", err)
	}
	identityB, err := transport.LoadOrCreateDTLSIdentity(filepath.Join(rootB, "dtls", "certificate.pem"), filepath.Join(rootB, "dtls", "private-key.pem"))
	if err != nil {
		t.Fatalf("load node b identity: %v", err)
	}
	appA, err := service.NewApp(rootA)
	if err != nil {
		t.Fatalf("new node a: %v", err)
	}
	appB, err := service.NewApp(rootB)
	if err != nil {
		t.Fatalf("new node b: %v", err)
	}
	endpointA := freeUDPAddress(t)
	endpointB := freeUDPAddress(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := appB.StartDTLSPeerDelivery(ctx, service.DTLSPeerConfig{ListenEndpoint: endpointB, Peers: []service.DTLSPeerInvite{{Endpoint: endpointA, CertificateBase64: base64.StdEncoding.EncodeToString(identityA.DER), CertificateFingerprint: identityA.Fingerprint}}}); err != nil {
		t.Fatalf("start node b dtls: %v", err)
	}
	if _, err := appA.StartDTLSPeerDelivery(ctx, service.DTLSPeerConfig{ListenEndpoint: endpointA, Peers: []service.DTLSPeerInvite{{Endpoint: endpointB, CertificateBase64: base64.StdEncoding.EncodeToString(identityB.DER), CertificateFingerprint: identityB.Fingerprint}}}); err != nil {
		t.Fatalf("start node a dtls: %v", err)
	}
	serverA := newBrowserProbeServer(t, appA, browserProbeOptions{editorText: "browser to neovim over direct dtls"})
	defer serverA.Close()
	serverB := httptest.NewUnstartedServer(service.NewServer(appB).Handler())
	serverB.Listener = listenTCP4OrSkip(t)
	serverB.Start()
	defer serverB.Close()

	repo := repoRoot(t)
	sidecar := startJSONProcess(t, repo, "node", filepath.Join(repo, "cmd", "grid-nvim-sidecar", "helper.bundle.cjs"), "--relay", serverB.URL)
	defer sidecar.Close()
	sidecar.WaitForType(t, "info")
	sidecar.Send(t, map[string]any{"type": "connect", "relay_url": serverB.URL, "participant_id": "nvim-b", "display_name": "Neovim B", "color": "#d66f1d"})
	sidecar.WaitForType(t, "connected")
	sidecar.Send(t, map[string]any{"type": "open", "doc_id": "real-browser-neovim-dtls"})
	sidecar.WaitForType(t, "opened")

	runHeadlessBrowser(t, chromePath, serverA.URL, "real-browser-neovim-dtls")
	waitForBrowserSyncMessageCount(t, appB, "real-browser-neovim-dtls", 1)
	sidecar.WaitForMessage(t, func(message map[string]any) bool {
		return stringField(t, message, "type") == "changed" && stringField(t, message, "content") == "browser to neovim over direct dtls"
	})

	const textFromNeovim = "neovim to browser over direct dtls"
	sidecar.Send(t, map[string]any{"type": "set_text", "content": textFromNeovim})
	waitForBrowserSyncMessageCount(t, appB, "real-browser-neovim-dtls", 2)
	waitForBrowserSyncMessageCount(t, appA, "real-browser-neovim-dtls", 2)
	serverAFinal := newBrowserProbeServer(t, appA)
	defer serverAFinal.Close()
	if dom := runHeadlessBrowser(t, chromePath, serverAFinal.URL, "real-browser-neovim-dtls"); !strings.Contains(dom, textFromNeovim) {
		t.Fatalf("browser did not render Neovim's direct-DTLS edit\n%s", dom)
	}
}

func TestHeadlessBrowserDirectDTLSRetriesAfterReceiverRestart(t *testing.T) {
	chromePath, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not available")
	}

	rootA := filepath.Join(t.TempDir(), "node-a")
	rootB := filepath.Join(t.TempDir(), "node-b")
	identityB, err := transport.LoadOrCreateDTLSIdentity(filepath.Join(rootB, "dtls", "certificate.pem"), filepath.Join(rootB, "dtls", "private-key.pem"))
	if err != nil {
		t.Fatalf("load node b identity: %v", err)
	}
	appA, err := service.NewApp(rootA)
	if err != nil {
		t.Fatalf("new node a: %v", err)
	}
	endpointB := freeUDPAddress(t)
	lossProxy := newPacketLossProxy(t, endpointB, 1)
	defer func() {
		if closeErr := lossProxy.Close(); closeErr != nil {
			t.Errorf("close packet-loss proxy: %v", closeErr)
		}
	}()
	ctxA, cancelA := context.WithCancel(context.Background())
	defer cancelA()
	if _, err := appA.StartDTLSPeerDelivery(ctxA, service.DTLSPeerConfig{ListenEndpoint: "127.0.0.1:0", Peers: []service.DTLSPeerInvite{{Endpoint: lossProxy.Endpoint(), CertificateBase64: base64.StdEncoding.EncodeToString(identityB.DER), CertificateFingerprint: identityB.Fingerprint}}}); err != nil {
		t.Fatalf("start node a dtls: %v", err)
	}

	const documentID = "browser-dtls-recovery"
	const beforeRestart = "browser edit persisted while receiver was unavailable"
	serverAFirst := newBrowserProbeServer(t, appA, browserProbeOptions{editorText: beforeRestart})
	runHeadlessBrowser(t, chromePath, serverAFirst.URL, documentID)
	serverAFirst.Close()
	waitForPersistedBrowserText(t, appA, documentID, beforeRestart)

	appB, err := service.NewApp(rootB)
	if err != nil {
		t.Fatalf("new node b: %v", err)
	}
	ctxB, cancelB := context.WithCancel(context.Background())
	if _, err := appB.StartDTLSPeerDelivery(ctxB, service.DTLSPeerConfig{ListenEndpoint: endpointB}); err != nil {
		t.Fatalf("start node b dtls: %v", err)
	}
	waitForBrowserSyncMessageCount(t, appB, documentID, 1)
	if lossProxy.DroppedPackets() == 0 {
		t.Fatal("packet-loss proxy did not drop a direct-DTLS datagram")
	}
	serverBFirst := newBrowserProbeServer(t, appB)
	if dom := runHeadlessBrowser(t, chromePath, serverBFirst.URL, documentID); !strings.Contains(dom, beforeRestart) {
		t.Fatalf("browser b did not render retry after receiver became available\n%s", dom)
	}
	serverBFirst.Close()
	cancelB()

	const duringRestart = "browser edit persisted while receiver restarted"
	serverASecond := newBrowserProbeServer(t, appA, browserProbeOptions{editorText: duringRestart})
	runHeadlessBrowser(t, chromePath, serverASecond.URL, documentID)
	serverASecond.Close()
	waitForPersistedBrowserText(t, appA, documentID, duringRestart)

	appB, err = service.NewApp(rootB)
	if err != nil {
		t.Fatalf("restart node b: %v", err)
	}
	ctxB, cancelB = context.WithCancel(context.Background())
	defer cancelB()
	if _, err := appB.StartDTLSPeerDelivery(ctxB, service.DTLSPeerConfig{ListenEndpoint: endpointB}); err != nil {
		t.Fatalf("restart node b dtls: %v", err)
	}
	waitForBrowserSyncMessageCount(t, appB, documentID, 2)
	serverBFinal := newBrowserProbeServer(t, appB)
	defer serverBFinal.Close()
	if dom := runHeadlessBrowser(t, chromePath, serverBFinal.URL, documentID); !strings.Contains(dom, duringRestart) {
		t.Fatalf("browser b did not render retry after receiver restart\n%s", dom)
	}
}

func runHeadlessBrowser(t *testing.T, chromePath, serverURL, documentID string) string {
	t.Helper()
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
		serverURL+"/?doc="+documentID,
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("chrome document editor: %v\n%s", err, string(output))
	}
	return string(output)
}

func waitForPersistedBrowserText(t *testing.T, app *service.App, documentID, want string) {
	t.Helper()
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		state := app.State(documentID)
		text, err := base64.StdEncoding.DecodeString(state.TextBase64)
		if err != nil {
			t.Fatalf("decode %s document state: %v", documentID, err)
		}
		if string(text) == want {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	state := app.State(documentID)
	text, err := base64.StdEncoding.DecodeString(state.TextBase64)
	if err != nil {
		t.Fatalf("decode final %s document state: %v", documentID, err)
	}
	t.Fatalf("persisted document = %q, want %q", string(text), want)
}

func waitForBrowserSyncMessageCount(t *testing.T, app *service.App, documentID string, want int) {
	t.Helper()
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		if len(app.SyncFeed(documentID, 0, 64).Messages) >= want {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	got := len(app.SyncFeed(documentID, 0, 64).Messages)
	t.Fatalf("direct-DTLS sync messages = %d, want at least %d", got, want)
}

type packetLossProxy struct {
	listener  *net.UDPConn
	upstream  *net.UDPAddr
	dropFirst int

	mu      sync.Mutex
	dropped int
	flows   map[string]*packetLossProxyFlow
}

type packetLossProxyFlow struct {
	client   *net.UDPAddr
	upstream *net.UDPConn
}

func newPacketLossProxy(t *testing.T, upstream string, dropFirst int) *packetLossProxy {
	t.Helper()
	upstreamAddress, err := net.ResolveUDPAddr("udp", upstream)
	if err != nil {
		t.Fatalf("resolve packet-loss upstream: %v", err)
	}
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	if err != nil {
		t.Fatalf("listen packet-loss proxy: %v", err)
	}
	proxy := &packetLossProxy{
		listener:  listener,
		upstream:  upstreamAddress,
		dropFirst: dropFirst,
		flows:     make(map[string]*packetLossProxyFlow),
	}
	go proxy.forward()
	return proxy
}

func (proxy *packetLossProxy) Endpoint() string {
	return proxy.listener.LocalAddr().String()
}

func (proxy *packetLossProxy) DroppedPackets() int {
	proxy.mu.Lock()
	defer proxy.mu.Unlock()
	return proxy.dropped
}

func (proxy *packetLossProxy) Close() error {
	proxy.mu.Lock()
	flows := make([]*packetLossProxyFlow, 0, len(proxy.flows))
	for _, flow := range proxy.flows {
		flows = append(flows, flow)
	}
	proxy.mu.Unlock()
	for _, flow := range flows {
		if err := flow.upstream.Close(); err != nil {
			return err
		}
	}
	return proxy.listener.Close()
}

func (proxy *packetLossProxy) forward() {
	buffer := make([]byte, 64*1024)
	for {
		count, client, err := proxy.listener.ReadFromUDP(buffer)
		if err != nil {
			return
		}
		packet := append([]byte(nil), buffer[:count]...)
		proxy.mu.Lock()
		if proxy.dropped < proxy.dropFirst {
			proxy.dropped++
			proxy.mu.Unlock()
			continue
		}
		flow, flowErr := proxy.flowLocked(client)
		proxy.mu.Unlock()
		if flowErr != nil {
			return
		}
		if _, writeErr := flow.upstream.Write(packet); writeErr != nil {
			return
		}
	}
}

func (proxy *packetLossProxy) flowLocked(client *net.UDPAddr) (*packetLossProxyFlow, error) {
	key := client.String()
	if flow := proxy.flows[key]; flow != nil {
		return flow, nil
	}
	upstream, err := net.DialUDP("udp", nil, proxy.upstream)
	if err != nil {
		return nil, err
	}
	flow := &packetLossProxyFlow{client: client, upstream: upstream}
	proxy.flows[key] = flow
	go proxy.returnPackets(flow)
	return flow, nil
}

func (proxy *packetLossProxy) returnPackets(flow *packetLossProxyFlow) {
	buffer := make([]byte, 64*1024)
	for {
		count, err := flow.upstream.Read(buffer)
		if err != nil {
			return
		}
		if _, writeErr := proxy.listener.WriteToUDP(buffer[:count], flow.client); writeErr != nil {
			return
		}
	}
}

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
		"browser sync: webrtc",
		`"webrtc_offer_posted":true`,
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
		"browser sync: webrtc",
		`"webrtc_offer_posted":true`,
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
	editorText               string
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
	editorText, err := json.Marshal(probeOptions.editorText)
	if err != nil {
		t.Fatalf("marshal browser edit text: %v", err)
	}
	probeHTML := bytes.Replace(indexHTML, []byte("</body>"), []byte(`<div id="startup-probe" hidden></div>
<script>
// Intent: Observe local browser carriage events separately from the rendered
// transport label so this headless proof does not treat UI timing as protocol
// evidence. Source: DI-gofut; DI-raron
const transportEvidence = { webrtc_offer_posted: false, recovery_sync_reads: 0 };
const browserEditorProofText = `+string(editorText)+`;
let browserEditorProofApplied = false;
const nativeFetch = window.fetch.bind(window);
window.fetch = async (input, init) => {
  const requestURL = new URL(input instanceof Request ? input.url : String(input), window.location.href);
  const method = init?.method || (input instanceof Request ? input.method : "GET");
  if (method === "GET" && requestURL.pathname.endsWith("/sync")) {
    transportEvidence.recovery_sync_reads += 1;
  }
  if (method === "POST" && requestURL.pathname.endsWith("/webrtc-offer")) {
    transportEvidence.webrtc_offer_posted = true;
  }
  return nativeFetch(input, init);
};
const startupProbeTimer = setInterval(() => {
  const content = document.querySelector(".cm-content");
  const target = document.getElementById("startup-probe");
  if (!content || !target) {
    return;
  }
  const text = content.textContent || "";
  if (browserEditorProofText && !browserEditorProofApplied && document.getElementById("status")?.classList.contains("status-ready")) {
    // Intent: Exercise the actual CodeMirror browser surface so this test
    // proves an editor-originated change reaches the persistent document path.
    // Source: DI-zegov; DI-gafit
    content.focus();
    document.execCommand("selectAll", false, null);
    document.execCommand("insertText", false, browserEditorProofText);
    browserEditorProofApplied = true;
  }
  target.textContent = JSON.stringify({ text, transport_evidence: transportEvidence });
  if (text.includes("Shared text should appear here.") && transportEvidence.webrtc_offer_posted && (!`+requireRecovery+` || transportEvidence.recovery_sync_reads > 0)) {
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
