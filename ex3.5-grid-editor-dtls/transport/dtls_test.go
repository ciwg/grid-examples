package transport

import (
	"context"
	"io"
	"path/filepath"
	"testing"
)

func TestDTLSDialExchangesBytesWithPinnedPeer(t *testing.T) {
	serverIdentity, _, _, err := NewDTLSIdentity()
	if err != nil {
		t.Fatalf("create server identity: %v", err)
	}
	clientIdentity, _, _, err := NewDTLSIdentity()
	if err != nil {
		t.Fatalf("create client identity: %v", err)
	}
	listener, err := DTLSListen("127.0.0.1:0", serverIdentity)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			t.Errorf("close listener: %v", closeErr)
		}
	}()

	serverResult := make(chan error, 1)
	go func() {
		connection, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverResult <- acceptErr
			return
		}
		defer connection.Close()
		payload := make([]byte, len("direct-peer"))
		if _, readErr := io.ReadFull(connection, payload); readErr != nil {
			serverResult <- readErr
			return
		}
		if string(payload) != "direct-peer" {
			serverResult <- io.ErrUnexpectedEOF
			return
		}
		_, writeErr := connection.Write([]byte("ack"))
		serverResult <- writeErr
	}()

	connection, err := DTLSDial(listener.Addr().String(), clientIdentity, serverIdentity.DER, serverIdentity.Fingerprint)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() {
		if closeErr := connection.Close(); closeErr != nil {
			t.Errorf("close client: %v", closeErr)
		}
	}()
	if _, err := connection.Write([]byte("direct-peer")); err != nil {
		t.Fatalf("write: %v", err)
	}
	ack := make([]byte, len("ack"))
	if _, err := io.ReadFull(connection, ack); err != nil {
		t.Fatalf("read ack: %v", err)
	}
	if string(ack) != "ack" {
		t.Fatalf("ack = %q, want ack", ack)
	}
	if err := <-serverResult; err != nil {
		t.Fatalf("server exchange: %v", err)
	}
}

func TestDTLSSendReturnsPeerReceipt(t *testing.T) {
	serverIdentity, _, _, err := NewDTLSIdentity()
	if err != nil {
		t.Fatalf("create server identity: %v", err)
	}
	clientIdentity, _, _, err := NewDTLSIdentity()
	if err != nil {
		t.Fatalf("create client identity: %v", err)
	}
	listener, err := DTLSListen("127.0.0.1:0", serverIdentity)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	serverResult := make(chan error, 1)
	received := make(chan string, 1)
	go func() {
		serverResult <- DTLSServe(ctx, listener, func(payload []byte) ([]byte, error) {
			received <- string(payload)
			return []byte("accepted"), nil
		})
	}()

	receipt, err := DTLSSend(listener.Addr().String(), clientIdentity, serverIdentity.DER, serverIdentity.Fingerprint, []byte("signed-envelope"))
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if string(receipt) != "accepted" {
		t.Fatalf("receipt = %q, want accepted", receipt)
	}
	if payload := <-received; payload != "signed-envelope" {
		t.Fatalf("payload = %q, want signed-envelope", payload)
	}
	cancel()
	if closeErr := listener.Close(); closeErr != nil {
		t.Fatalf("close listener: %v", closeErr)
	}
	if err := <-serverResult; err != nil {
		t.Fatalf("serve: %v", err)
	}
}

func TestLoadOrCreateDTLSIdentityKeepsFingerprint(t *testing.T) {
	certificatePath := filepath.Join(t.TempDir(), "dtls", "certificate.pem")
	privateKeyPath := filepath.Join(filepath.Dir(certificatePath), "private-key.pem")
	first, err := LoadOrCreateDTLSIdentity(certificatePath, privateKeyPath)
	if err != nil {
		t.Fatalf("create identity: %v", err)
	}
	second, err := LoadOrCreateDTLSIdentity(certificatePath, privateKeyPath)
	if err != nil {
		t.Fatalf("reload identity: %v", err)
	}
	if first.Fingerprint != second.Fingerprint {
		t.Fatalf("fingerprint changed: %s != %s", first.Fingerprint, second.Fingerprint)
	}
}

func TestNewDTLSClientConfigRejectsMissingFingerprint(t *testing.T) {
	identity, _, _, err := NewDTLSIdentity()
	if err != nil {
		t.Fatalf("create identity: %v", err)
	}
	if _, err := NewDTLSClientConfig(identity, nil, ""); err == nil {
		t.Fatal("missing fingerprint was accepted")
	}
}
