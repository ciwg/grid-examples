package service_test

import (
	"context"
	"encoding/base64"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex3.5-grid-editor-dtls/service"
	"github.com/computerscienceiscool/grid-examples/ex3.5-grid-editor-dtls/transport"
)

func TestDirectDTLSPeerReplicationPersistsSignedEnvelope(t *testing.T) {
	appA, err := service.NewApp(filepath.Join(t.TempDir(), "node-a"))
	if err != nil {
		t.Fatalf("new node a: %v", err)
	}
	appB, err := service.NewApp(filepath.Join(t.TempDir(), "node-b"))
	if err != nil {
		t.Fatalf("new node b: %v", err)
	}
	serverIdentity, _, _, err := transport.NewDTLSIdentity()
	if err != nil {
		t.Fatalf("server identity: %v", err)
	}
	clientIdentity, _, _, err := transport.NewDTLSIdentity()
	if err != nil {
		t.Fatalf("client identity: %v", err)
	}
	listener, err := transport.DTLSListen("127.0.0.1:0", serverIdentity)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = transport.DTLSServe(ctx, listener, func(payload []byte) ([]byte, error) { return []byte("accepted"), appB.IngestRawBase64(string(payload)) })
	}()
	record, err := appA.PostSync("demo", "alice", "bob", base64.StdEncoding.EncodeToString([]byte{1, 2, 3}), "browser", "", "")
	if err != nil {
		t.Fatalf("post sync: %v", err)
	}
	messages, _ := appA.PeerMessagesSince(record.Offset, 1)
	receipt, err := transport.DTLSSend(listener.Addr().String(), clientIdentity, serverIdentity.DER, serverIdentity.Fingerprint, []byte(messages[0]))
	if err != nil || string(receipt) != "accepted" {
		t.Fatalf("send receipt=%q err=%v", receipt, err)
	}
	feed := appB.SyncFeed("demo", 0, 1)
	if len(feed.Messages) != 1 {
		t.Fatalf("received records = %d, want 1", len(feed.Messages))
	}
	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}
}

func TestDirectDTLSRetriesAfterLossAndReceiverRestart(t *testing.T) {
	rootA := filepath.Join(t.TempDir(), "node-a")
	rootB := filepath.Join(t.TempDir(), "node-b")
	appA, err := service.NewApp(rootA)
	if err != nil {
		t.Fatalf("new node a: %v", err)
	}
	peerIdentity, err := transport.LoadOrCreateDTLSIdentity(filepath.Join(rootB, "dtls", "certificate.pem"), filepath.Join(rootB, "dtls", "private-key.pem"))
	if err != nil {
		t.Fatalf("prepare node b identity: %v", err)
	}
	endpoint := freeUDPAddress(t)
	ctxA, cancelA := context.WithCancel(context.Background())
	defer cancelA()
	_, err = appA.StartDTLSPeerDelivery(ctxA, service.DTLSPeerConfig{ListenEndpoint: "127.0.0.1:0", Peers: []service.DTLSPeerInvite{{Endpoint: endpoint, CertificateBase64: base64.StdEncoding.EncodeToString(peerIdentity.DER), CertificateFingerprint: peerIdentity.Fingerprint}}})
	if err != nil {
		t.Fatalf("start node a: %v", err)
	}
	if _, err := appA.PostSync("demo", "alice", "bob", base64.StdEncoding.EncodeToString([]byte{1}), "browser", "", ""); err != nil {
		t.Fatalf("post while peer absent: %v", err)
	}
	appB, err := service.NewApp(rootB)
	if err != nil {
		t.Fatalf("new node b: %v", err)
	}
	ctxB, cancelB := context.WithCancel(context.Background())
	if _, err := appB.StartDTLSPeerDelivery(ctxB, service.DTLSPeerConfig{ListenEndpoint: endpoint}); err != nil {
		t.Fatalf("start node b: %v", err)
	}
	waitForSyncCount(t, appB, 1)
	cancelB()
	if _, err := appA.PostSync("demo", "alice", "bob", base64.StdEncoding.EncodeToString([]byte{2}), "browser", "", ""); err != nil {
		t.Fatalf("post during restart: %v", err)
	}
	appB, err = service.NewApp(rootB)
	if err != nil {
		t.Fatalf("restart node b: %v", err)
	}
	ctxB, cancelB = context.WithCancel(context.Background())
	defer cancelB()
	if _, err := appB.StartDTLSPeerDelivery(ctxB, service.DTLSPeerConfig{ListenEndpoint: endpoint}); err != nil {
		t.Fatalf("restart dtls node b: %v", err)
	}
	waitForSyncCount(t, appB, 2)
}

func freeUDPAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	if err != nil {
		t.Fatalf("reserve udp port: %v", err)
	}
	address := listener.LocalAddr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("release udp port: %v", err)
	}
	return address
}

func waitForSyncCount(t *testing.T, app interface {
	SyncFeed(string, uint64, int) service.SyncFeed
}, want int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if len(app.SyncFeed("demo", 0, 16).Messages) == want {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %d synchronized messages", want)
}
