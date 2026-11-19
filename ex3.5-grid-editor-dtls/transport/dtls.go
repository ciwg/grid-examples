// Package transport provides the direct DTLS carriage used between Grid Editor
// nodes. It deliberately owns only encrypted peer connectivity; signed grid
// envelopes retain their application-level identity and meaning.
package transport

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/pion/dtls/v3"
)

const dtlsServerName = "grid-editor-dtls"

const dtlsHandshakeTimeout = 5 * time.Second

// DTLSIdentity is the locally owned transport certificate and private key.
// Intent: A node presents a stable, fingerprintable transport identity to a
// directly configured peer instead of relying on a central relay. Source:
// DI-pirov.
type DTLSIdentity struct {
	Certificate tls.Certificate
	Fingerprint string
	DER         []byte
}

// NewDTLSIdentity creates a self-signed certificate for one local Grid Editor
// node. Callers persist the returned PEM material through their node storage.
func NewDTLSIdentity() (DTLSIdentity, []byte, []byte, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return DTLSIdentity{}, nil, nil, fmt.Errorf("generate dtls private key: %w", err)
	}
	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		return DTLSIdentity{}, nil, nil, fmt.Errorf("generate certificate serial: %w", err)
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: dtlsServerName},
		DNSNames:     []string{dtlsServerName},
		NotBefore:    time.Now().Add(-time.Minute),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return DTLSIdentity{}, nil, nil, fmt.Errorf("create dtls certificate: %w", err)
	}
	certificatePEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return DTLSIdentity{}, nil, nil, fmt.Errorf("marshal dtls private key: %w", err)
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER})
	certificate, err := tls.X509KeyPair(certificatePEM, privateKeyPEM)
	if err != nil {
		return DTLSIdentity{}, nil, nil, fmt.Errorf("load dtls certificate: %w", err)
	}
	return DTLSIdentity{
		Certificate: certificate,
		Fingerprint: dtlsFingerprint(der),
		DER:         der,
	}, certificatePEM, privateKeyPEM, nil
}

// LoadDTLSIdentity loads a previously persisted certificate and private key.
func LoadDTLSIdentity(certificatePEM []byte, privateKeyPEM []byte) (DTLSIdentity, error) {
	certificate, err := tls.X509KeyPair(certificatePEM, privateKeyPEM)
	if err != nil {
		return DTLSIdentity{}, fmt.Errorf("load dtls certificate: %w", err)
	}
	if len(certificate.Certificate) == 0 {
		return DTLSIdentity{}, fmt.Errorf("dtls certificate chain is empty")
	}
	return DTLSIdentity{
		Certificate: certificate,
		Fingerprint: dtlsFingerprint(certificate.Certificate[0]),
		DER:         certificate.Certificate[0],
	}, nil
}

// LoadOrCreateDTLSIdentity keeps a node's transport fingerprint stable across
// restarts so an out-of-band peer invitation remains valid.
func LoadOrCreateDTLSIdentity(certificatePath string, privateKeyPath string) (DTLSIdentity, error) {
	certificatePEM, certificateErr := os.ReadFile(certificatePath)
	privateKeyPEM, privateKeyErr := os.ReadFile(privateKeyPath)
	if certificateErr == nil && privateKeyErr == nil {
		identity, err := LoadDTLSIdentity(certificatePEM, privateKeyPEM)
		if err != nil {
			return DTLSIdentity{}, fmt.Errorf("load persisted dtls identity: %w", err)
		}
		return identity, nil
	}
	if certificateErr == nil || privateKeyErr == nil {
		return DTLSIdentity{}, fmt.Errorf("incomplete persisted dtls identity")
	}
	if !os.IsNotExist(certificateErr) || !os.IsNotExist(privateKeyErr) {
		return DTLSIdentity{}, fmt.Errorf("read persisted dtls identity: certificate=%v private-key=%v", certificateErr, privateKeyErr)
	}
	identity, certificatePEM, privateKeyPEM, err := NewDTLSIdentity()
	if err != nil {
		return DTLSIdentity{}, err
	}
	if err := os.MkdirAll(filepath.Dir(certificatePath), 0700); err != nil {
		return DTLSIdentity{}, fmt.Errorf("create dtls identity directory: %w", err)
	}
	if err := os.WriteFile(certificatePath, certificatePEM, 0600); err != nil {
		return DTLSIdentity{}, fmt.Errorf("write dtls certificate: %w", err)
	}
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return DTLSIdentity{}, fmt.Errorf("write dtls private key: %w", err)
	}
	return identity, nil
}

// NewDTLSServerConfig configures a node to authenticate itself to direct peers.
func NewDTLSServerConfig(identity DTLSIdentity) *dtls.Config {
	return &dtls.Config{
		Certificates:       []tls.Certificate{identity.Certificate},
		MTU:                1200,
		SupportedProtocols: []string{"grid-editor-dtls/1"},
	}
}

// NewDTLSClientConfig verifies a direct peer's self-signed certificate by the
// fingerprint carried in the out-of-band invitation.
func NewDTLSClientConfig(identity DTLSIdentity, peerCertificateDER []byte, peerFingerprint string) (*dtls.Config, error) {
	if peerFingerprint == "" {
		return nil, fmt.Errorf("peer dtls fingerprint is required")
	}
	peerCertificate, err := x509.ParseCertificate(peerCertificateDER)
	if err != nil {
		return nil, fmt.Errorf("parse invited peer certificate: %w", err)
	}
	peerRoots := x509.NewCertPool()
	peerRoots.AddCert(peerCertificate)
	return &dtls.Config{
		Certificates:       []tls.Certificate{identity.Certificate},
		ServerName:         dtlsServerName,
		MTU:                1200,
		RootCAs:            peerRoots,
		SupportedProtocols: []string{"grid-editor-dtls/1"},
		VerifyPeerCertificate: func(rawCertificates [][]byte, verifiedChains [][]*x509.Certificate) error {
			if len(rawCertificates) == 0 {
				return fmt.Errorf("peer did not present a dtls certificate")
			}
			if dtlsFingerprint(rawCertificates[0]) != peerFingerprint {
				return fmt.Errorf("peer dtls certificate fingerprint does not match invitation")
			}
			return nil
		},
	}, nil
}

// DTLSListen accepts direct peer sessions on a node-owned UDP endpoint.
func DTLSListen(endpoint string, identity DTLSIdentity) (net.Listener, error) {
	address, err := net.ResolveUDPAddr("udp", endpoint)
	if err != nil {
		return nil, fmt.Errorf("resolve dtls listener endpoint: %w", err)
	}
	listener, err := dtls.Listen("udp", address, NewDTLSServerConfig(identity))
	if err != nil {
		return nil, fmt.Errorf("listen for dtls peers: %w", err)
	}
	return listener, nil
}

// DTLSDial opens one direct, fingerprint-pinned DTLS session to a known peer.
func DTLSDial(endpoint string, identity DTLSIdentity, peerCertificateDER []byte, peerFingerprint string) (*dtls.Conn, error) {
	address, err := net.ResolveUDPAddr("udp", endpoint)
	if err != nil {
		return nil, fmt.Errorf("resolve dtls peer endpoint: %w", err)
	}
	config, err := NewDTLSClientConfig(identity, peerCertificateDER, peerFingerprint)
	if err != nil {
		return nil, err
	}
	packetConnection, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, fmt.Errorf("listen for dtls peer: %w", err)
	}
	connection, err := dtls.Client(packetConnection, address, config)
	if err != nil {
		if closeErr := packetConnection.Close(); closeErr != nil {
			return nil, fmt.Errorf("create dtls peer connection: %w (close packet connection: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("dial dtls peer: %w", err)
	}
	handshakeContext, cancel := context.WithTimeout(context.Background(), dtlsHandshakeTimeout)
	defer cancel()
	if err := connection.HandshakeContext(handshakeContext); err != nil {
		if closeErr := connection.Close(); closeErr != nil {
			return nil, fmt.Errorf("handshake dtls peer: %w (close connection: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("handshake dtls peer: %w", err)
	}
	return connection, nil
}

func dtlsFingerprint(certificateDER []byte) string {
	digest := sha256.Sum256(certificateDER)
	return hex.EncodeToString(digest[:])
}
