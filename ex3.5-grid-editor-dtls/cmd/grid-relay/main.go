package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/computerscienceiscool/grid-examples/ex3.5-grid-editor-dtls/service"
)

type peerFlags []string

func (flags *peerFlags) String() string {
	return ""
}

func (flags *peerFlags) Set(value string) error {
	*flags = append(*flags, value)
	return nil
}

func main() {
	var (
		// Intent: Move the copied example off ex2's default loopback ports so the
		// two examples can run side by side on one machine. Source: DI-vatub
		listen = flag.String("listen", "127.0.0.1:7025", "listen address")
		// Intent: Keep each decentralized Ex3.5 node's persisted identity, DTLS
		// keys, and signed-object log separate from the WebSocket example. Source:
		// DI-pirov.
		dataRoot          = flag.String("data-root", ".grid-editor-dtls", "local runtime data root")
		remoteAccessToken = flag.String("remote-access-token", "", "optional bootstrap token for remote mutation sessions")
		dtlsListen        = flag.String("dtls-listen", ":7026", "direct DTLS/UDP listen address")
		dtlsAdvertise     = flag.String("dtls-advertise", "", "public direct DTLS/UDP endpoint to place in a peer invitation")
		browserTransport  = flag.String("browser-transport", "webrtc", "browser live transport: webrtc, websocket, or polling")
		peers             peerFlags
	)
	flag.Var(&peers, "dtls-peer", "base64url JSON direct-peer invitation (repeatable)")
	flag.Parse()

	app, err := service.NewApp(*dataRoot, service.AppOptions{RemoteAccessToken: *remoteAccessToken})
	if err != nil {
		log.Fatalf("new relay: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dtlsPeers := make([]service.DTLSPeerInvite, 0, len(peers))
	for _, encodedInvite := range peers {
		peerInvite, err := parseDTLSPeerInvite(encodedInvite)
		if err != nil {
			log.Fatalf("parse direct peer invitation: %v", err)
		}
		dtlsPeers = append(dtlsPeers, peerInvite)
	}
	node, err := app.StartDTLSPeerDelivery(ctx, service.DTLSPeerConfig{
		ListenEndpoint: *dtlsListen,
		Peers:          dtlsPeers,
	})
	if err != nil {
		log.Fatalf("start direct dtls peer delivery: %v", err)
	}

	server := service.NewServerWithBrowserTransport(app, *browserTransport)
	log.Printf("grid-relay listening on %s", *listen)
	log.Printf("relay author=%s", app.Meta().LocalID)
	log.Printf("dtls endpoint=%s", *dtlsListen)
	log.Printf("dtls fingerprint=%s", node.Identity.Fingerprint)
	if *dtlsAdvertise != "" {
		invitation := service.DTLSPeerInvite{
			Endpoint:               *dtlsAdvertise,
			CertificateBase64:      base64.StdEncoding.EncodeToString(node.Identity.DER),
			CertificateFingerprint: node.Identity.Fingerprint,
		}
		invitationBytes, err := json.Marshal(invitation)
		if err != nil {
			log.Fatalf("marshal direct peer invitation: %v", err)
		}
		log.Printf("dtls invitation=%s", base64.RawURLEncoding.EncodeToString(invitationBytes))
	} else {
		log.Printf("dtls invitation unavailable: set --dtls-advertise to a reachable host:port")
	}
	log.Printf("live-document pCID=%s", app.Meta().DocumentPCID)
	log.Printf("live-awareness pCID=%s", app.Meta().AwarenessPCID)
	log.Printf("remote mutation bootstrap enabled=%t", app.RemoteAccessEnabled())
	log.Fatal(http.ListenAndServe(*listen, server.Handler()))
}

func parseDTLSPeerInvite(encodedInvite string) (service.DTLSPeerInvite, error) {
	inviteBytes, err := base64.RawURLEncoding.DecodeString(encodedInvite)
	if err != nil {
		return service.DTLSPeerInvite{}, err
	}
	var peerInvite service.DTLSPeerInvite
	if err := json.Unmarshal(inviteBytes, &peerInvite); err != nil {
		return service.DTLSPeerInvite{}, err
	}
	if peerInvite.Endpoint == "" || peerInvite.CertificateBase64 == "" || peerInvite.CertificateFingerprint == "" {
		return service.DTLSPeerInvite{}, fmt.Errorf("invitation is incomplete")
	}
	return peerInvite, nil
}
