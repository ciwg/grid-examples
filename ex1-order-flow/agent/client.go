package agent

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/artifact"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/collector"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
)

type Client struct {
	role      string
	conn      net.Conn
	store     *artifact.Store
	collector *collector.Client
}

func DialKernel(ctx context.Context, cfg Config) (*Client, error) {
	store, err := artifact.NewStore(cfg.Role, cfg.DataDir)
	if err != nil {
		return nil, err
	}
	collectorClient, err := collector.Dial(ctx, cfg.CollectorAddr)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(15 * time.Second)
	var conn net.Conn
	for time.Now().Before(deadline) {
		conn, err = (&net.Dialer{}).DialContext(ctx, "tcp", cfg.KernelAddr)
		if err == nil {
			return &Client{
				role:      cfg.Role,
				conn:      conn,
				store:     store,
				collector: collectorClient,
			}, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
	return nil, fmt.Errorf("dial kernel %s: %w", cfg.KernelAddr, err)
}

func (client *Client) Close() error {
	var err error
	if client.collector != nil {
		if closeErr := client.collector.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if client.conn != nil {
		if closeErr := client.conn.Close(); closeErr != nil {
			if err != nil {
				return fmt.Errorf("close collector: %v; close kernel conn: %w", err, closeErr)
			}
			return closeErr
		}
	}
	return err
}

func (client *Client) Register(ctx context.Context, envelope protocol.Envelope) error {
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return err
	}
	return writeFrame(ctx, client.conn, envelopeBytes)
}

func (client *Client) Send(ctx context.Context, envelope protocol.Envelope, parents []string) (string, error) {
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		return "", fmt.Errorf("marshal envelope: %w", err)
	}
	exactCID, err := client.store.SaveEnvelope("send", envelopeBytes, parents, envelope.PCID.String())
	if err != nil {
		return "", err
	}
	if err := client.collector.Send(artifact.CollectorArtifactRecord{
		SourceRole:     client.role,
		ExactCID:       exactCID,
		PCID:           envelope.PCID.String(),
		ParentCIDs:     append([]string(nil), parents...),
		EnvelopeBase64: artifact.EnvelopeBase64(envelopeBytes),
	}); err != nil {
		return "", err
	}
	if err := writeFrame(ctx, client.conn, envelopeBytes); err != nil {
		return "", err
	}
	return exactCID, nil
}

func (client *Client) Receive(ctx context.Context) (protocol.Envelope, string, error) {
	envelopeBytes, err := readFrame(ctx, client.conn)
	if err != nil {
		return protocol.Envelope{}, "", err
	}
	envelope, err := protocol.ParseEnvelope(envelopeBytes)
	if err != nil {
		return protocol.Envelope{}, "", err
	}
	exactCID, err := client.store.SaveEnvelope("recv", envelopeBytes, nil, envelope.PCID.String())
	if err != nil {
		return protocol.Envelope{}, "", err
	}
	return envelope, exactCID, nil
}

func writeFrame(ctx context.Context, conn net.Conn, payload []byte) error {
	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetWriteDeadline(deadline); err != nil {
			return err
		}
	} else {
		if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
			return err
		}
	}
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(payload)))
	if _, err := conn.Write(header); err != nil {
		return err
	}
	_, err := conn.Write(payload)
	return err
}

func readFrame(ctx context.Context, conn net.Conn) ([]byte, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetReadDeadline(deadline); err != nil {
			return nil, err
		}
	} else {
		if err := conn.SetReadDeadline(time.Time{}); err != nil {
			return nil, err
		}
	}
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(header)
	if length == 0 {
		return nil, fmt.Errorf("empty frame")
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
