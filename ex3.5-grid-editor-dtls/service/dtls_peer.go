package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex3.5-grid-editor-dtls/transport"
)

const dtlsPeerDeliveryInterval = 250 * time.Millisecond

// DTLSPeerInvite is the out-of-band material a node needs to establish a
// direct, fingerprint-pinned DTLS connection to another Grid Editor node.
type DTLSPeerInvite struct {
	Endpoint               string
	CertificateBase64      string
	CertificateFingerprint string
}

// DTLSPeerConfig configures one decentralized node's direct DTLS listener and
// known peers.
type DTLSPeerConfig struct {
	ListenEndpoint string
	Peers          []DTLSPeerInvite
}

// DTLSPeerNode describes the running node listener and its invitation identity.
type DTLSPeerNode struct {
	Identity transport.DTLSIdentity
	Endpoint string
}

// StartDTLSPeerDelivery starts direct peer ingress and durable outbound retry.
// Intent: Keep every node responsible for its own signed-envelope persistence
// and direct peer delivery rather than delegating live collaboration to a
// central relay. Source: DI-higom; DI-pirov.
func (app *App) StartDTLSPeerDelivery(ctx context.Context, peerConfig DTLSPeerConfig) (DTLSPeerNode, error) {
	if peerConfig.ListenEndpoint == "" {
		return DTLSPeerNode{}, fmt.Errorf("dtls listener endpoint is required")
	}
	identity, err := transport.LoadOrCreateDTLSIdentity(
		filepath.Join(app.dataRoot, "dtls", "certificate.pem"),
		filepath.Join(app.dataRoot, "dtls", "private-key.pem"),
	)
	if err != nil {
		return DTLSPeerNode{}, fmt.Errorf("load dtls node identity: %w", err)
	}
	listener, err := transport.DTLSListen(peerConfig.ListenEndpoint, identity)
	if err != nil {
		return DTLSPeerNode{}, err
	}
	go func() {
		<-ctx.Done()
		if closeErr := listener.Close(); closeErr != nil {
			log.Printf("grid dtls listener close: %v", closeErr)
		}
	}()
	go func() {
		serveErr := transport.DTLSServe(ctx, listener, func(payload []byte) ([]byte, error) {
			if err := app.IngestRawBase64(string(payload)); err != nil {
				return nil, err
			}
			return []byte("accepted"), nil
		})
		if serveErr != nil {
			log.Printf("grid dtls peer serve: %v", serveErr)
		}
	}()
	for _, peerInvite := range peerConfig.Peers {
		peerInvite := peerInvite
		go app.deliverDTLSPeer(ctx, identity, peerInvite)
	}
	return DTLSPeerNode{Identity: identity, Endpoint: listener.Addr().String()}, nil
}

func (app *App) deliverDTLSPeer(ctx context.Context, identity transport.DTLSIdentity, peerInvite DTLSPeerInvite) {
	peerCertificateDER, err := base64.StdEncoding.DecodeString(peerInvite.CertificateBase64)
	if err != nil {
		return
	}
	var deliveryOffset uint64
	deliveryTicker := time.NewTicker(dtlsPeerDeliveryInterval)
	defer deliveryTicker.Stop()
	for {
		messages, nextOffset := app.PeerMessagesSince(deliveryOffset, defaultFeedLimit)
		deliveredAll := true
		for _, message := range messages {
			receipt, sendErr := transport.DTLSSend(peerInvite.Endpoint, identity, peerCertificateDER, peerInvite.CertificateFingerprint, []byte(message))
			if sendErr != nil || string(receipt) != "accepted" {
				deliveredAll = false
				break
			}
		}
		if len(messages) > 0 && deliveredAll {
			deliveryOffset = nextOffset
		}
		select {
		case <-ctx.Done():
			return
		case <-deliveryTicker.C:
		}
	}
}
