package collector

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/artifact"
)

type Service struct {
	Address  string
	DataDir  string
	listener net.Listener
	mu       sync.Mutex
}

type DAGRecord struct {
	SourceRole string   `json:"source_role"`
	ExactCID   string   `json:"exact_cid"`
	PCID       string   `json:"pcid"`
	ParentCIDs []string `json:"parent_cids"`
}

func (service *Service) Run(ctx context.Context) (runErr error) {
	if err := os.MkdirAll(filepath.Join(service.DataDir, "message-cas"), 0o755); err != nil {
		return fmt.Errorf("mkdir collector data dir: %w", err)
	}
	listener, err := net.Listen("tcp", service.Address)
	if err != nil {
		return fmt.Errorf("listen collector: %w", err)
	}
	service.listener = listener
	defer func() {
		if closeErr := listener.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) && runErr == nil {
			runErr = closeErr
		}
	}()
	go func() {
		<-ctx.Done()
		if closeErr := listener.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
			fmt.Fprintf(os.Stderr, "collector listener close: %v\n", closeErr)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) || ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("accept collector conn: %w", err)
		}
		go service.handleConn(conn)
	}
}

func (service *Service) handleConn(conn net.Conn) {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "collector conn close: %v\n", closeErr)
		}
	}()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := append([]byte(nil), scanner.Bytes()...)
		var record artifact.CollectorArtifactRecord
		if err := json.Unmarshal(line, &record); err != nil {
			continue
		}
		if err := service.recordArtifact(record); err != nil {
			continue
		}
	}
}

func (service *Service) recordArtifact(record artifact.CollectorArtifactRecord) (err error) {
	envelopeBytes, err := base64.StdEncoding.DecodeString(record.EnvelopeBase64)
	if err != nil {
		return fmt.Errorf("decode artifact envelope: %w", err)
	}
	artifactPath := filepath.Join(service.DataDir, "message-cas", record.ExactCID+".cbor")
	if statErr := os.WriteFile(artifactPath, envelopeBytes, 0o644); statErr != nil {
		return fmt.Errorf("write collector artifact: %w", statErr)
	}
	dagRecord := DAGRecord{
		SourceRole: record.SourceRole,
		ExactCID:   record.ExactCID,
		PCID:       record.PCID,
		ParentCIDs: append([]string(nil), record.ParentCIDs...),
	}
	recordBytes, err := json.Marshal(dagRecord)
	if err != nil {
		return fmt.Errorf("marshal dag record: %w", err)
	}
	service.mu.Lock()
	defer service.mu.Unlock()
	file, err := os.OpenFile(filepath.Join(service.DataDir, "message-dag.jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open dag file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if _, err := file.Write(append(recordBytes, '\n')); err != nil {
		return fmt.Errorf("append dag record: %w", err)
	}
	return nil
}

type Client struct {
	address string
	conn    net.Conn
}

func Dial(ctx context.Context, address string) (*Client, error) {
	deadline := time.Now().Add(15 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
		if err == nil {
			return &Client{address: address, conn: conn}, nil
		}
		lastErr = err
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
	return nil, fmt.Errorf("dial collector %s: %w", address, lastErr)
}

func (client *Client) Close() error {
	if client == nil || client.conn == nil {
		return nil
	}
	return client.conn.Close()
}

func (client *Client) Send(record artifact.CollectorArtifactRecord) error {
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal collector record: %w", err)
	}
	if _, err := client.conn.Write(append(recordBytes, '\n')); err != nil {
		return fmt.Errorf("write collector record: %w", err)
	}
	return nil
}
