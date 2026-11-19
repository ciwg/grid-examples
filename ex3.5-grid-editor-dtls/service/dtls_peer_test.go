package service_test

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"testing"

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
