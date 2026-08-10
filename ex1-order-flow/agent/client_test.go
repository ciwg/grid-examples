package agent

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/artifact"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
)

func TestReceiveRecordsMalformedInput(t *testing.T) {
	client, peer, root := testClient(t)
	defer closeTestPeer(t, peer)
	go func() {
		if err := writeFrame(context.Background(), peer, []byte{0xff}); err != nil {
			t.Errorf("write malformed frame: %v", err)
		}
	}()
	if _, _, err := client.Receive(context.Background()); err == nil {
		t.Fatal("Receive succeeded for malformed input")
	}
	assertObservationKind(t, root, "malformed_input")
}

func TestReceiveTypedRecordsInvalidProof(t *testing.T) {
	client, peer, root := testClient(t)
	defer closeTestPeer(t, peer)
	payload, err := protocol.Marshal(map[string]string{"kind": "submit"})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	envelope, err := BuildSignedEnvelope("seller", protocol.OrderProfile, payload)
	if err != nil {
		t.Fatalf("build envelope: %v", err)
	}
	envelope.ProofBytes[len(envelope.ProofBytes)-1] ^= 0x01
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}
	go func() {
		if err := writeFrame(context.Background(), peer, envelopeBytes); err != nil {
			t.Errorf("write invalid-proof frame: %v", err)
		}
	}()
	var received map[string]string
	if _, _, err := ReceiveTyped(context.Background(), client, "seller", protocol.OrderProfile, &received); err == nil {
		t.Fatal("ReceiveTyped succeeded for invalid proof")
	}
	assertObservationKind(t, root, "invalid_proof")
}

func TestReceiveTypedRecordsUnexpectedPCID(t *testing.T) {
	client, peer, root := testClient(t)
	defer closeTestPeer(t, peer)
	payload, err := protocol.Marshal(map[string]string{"kind": "request"})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	envelope, err := BuildSignedEnvelope("seller", protocol.PickPackProfile, payload)
	if err != nil {
		t.Fatalf("build envelope: %v", err)
	}
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}
	go func() {
		if err := writeFrame(context.Background(), peer, envelopeBytes); err != nil {
			t.Errorf("write unexpected-pcid frame: %v", err)
		}
	}()
	var received map[string]string
	if _, _, err := ReceiveTyped(context.Background(), client, "seller", protocol.OrderProfile, &received); err == nil {
		t.Fatal("ReceiveTyped succeeded for unexpected pCID")
	}
	assertObservationKind(t, root, "unexpected_pcid")
}

func testClient(t *testing.T) (*Client, net.Conn, string) {
	t.Helper()
	root := t.TempDir()
	store, err := artifact.NewStore("alice", root)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	clientConn, peer := net.Pipe()
	t.Cleanup(func() {
		if closeErr := clientConn.Close(); closeErr != nil {
			t.Errorf("close client connection: %v", closeErr)
		}
	})
	return &Client{role: "alice", conn: clientConn, store: store}, peer, root
}

func closeTestPeer(t *testing.T, peer net.Conn) {
	t.Helper()
	if err := peer.Close(); err != nil {
		t.Errorf("close test peer: %v", err)
	}
}

func assertObservationKind(t *testing.T, root string, kind string) {
	t.Helper()
	observations, err := os.ReadFile(filepath.Join(root, "observations.jsonl"))
	if err != nil {
		t.Fatalf("read observations: %v", err)
	}
	if !strings.Contains(string(observations), `"kind":"`+kind+`"`) {
		t.Fatalf("observations = %s, want %q", observations, kind)
	}
}
