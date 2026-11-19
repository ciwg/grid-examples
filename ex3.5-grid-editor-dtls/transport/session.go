package transport

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const maxDTLSMessageBytes = 1 << 20

// DTLSMessageHandler consumes one complete direct-peer message and returns a
// receipt. Intent: Make direct document delivery observable and acknowledgeable
// without assigning durable message meaning to the DTLS transport. Source:
// DI-higom; DI-pirov.
type DTLSMessageHandler func([]byte) ([]byte, error)

// DTLSServe accepts direct peer sessions until the context is cancelled.
func DTLSServe(ctx context.Context, listener net.Listener, handler DTLSMessageHandler) error {
	for {
		connection, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return fmt.Errorf("accept dtls peer: %w", err)
			}
		}
		go handleDTLSConnection(connection, handler)
	}
}

// DTLSSend delivers one framed message to a known direct peer and returns the
// peer's receipt. Callers retain the message for bounded retry on failure.
func DTLSSend(endpoint string, identity DTLSIdentity, peerCertificateDER []byte, peerFingerprint string, payload []byte) ([]byte, error) {
	connection, err := DTLSDial(endpoint, identity, peerCertificateDER, peerFingerprint)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = connection.Close()
	}()
	if err := writeDTLSFrame(connection, payload); err != nil {
		return nil, fmt.Errorf("write direct peer message: %w", err)
	}
	receipt, err := readDTLSFrame(connection)
	if err != nil {
		return nil, fmt.Errorf("read direct peer receipt: %w", err)
	}
	return receipt, nil
}

func handleDTLSConnection(connection net.Conn, handler DTLSMessageHandler) {
	defer func() {
		_ = connection.Close()
	}()
	payload, err := readDTLSFrame(connection)
	if err != nil {
		return
	}
	receipt, err := handler(payload)
	if err != nil {
		receipt = []byte("error: " + err.Error())
	}
	_ = writeDTLSFrame(connection, receipt)
}

func writeDTLSFrame(writer io.Writer, payload []byte) error {
	if len(payload) == 0 || len(payload) > maxDTLSMessageBytes {
		return fmt.Errorf("dtls frame size %d is outside 1..%d", len(payload), maxDTLSMessageBytes)
	}
	var length [4]byte
	binary.BigEndian.PutUint32(length[:], uint32(len(payload)))
	if _, err := writer.Write(length[:]); err != nil {
		return err
	}
	_, err := writer.Write(payload)
	return err
}

func readDTLSFrame(reader io.Reader) ([]byte, error) {
	var length [4]byte
	if _, err := io.ReadFull(reader, length[:]); err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(length[:])
	if size == 0 || size > maxDTLSMessageBytes {
		return nil, fmt.Errorf("invalid dtls frame size %d", size)
	}
	payload := make([]byte, size)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
