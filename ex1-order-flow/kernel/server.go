package kernel

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ipfs/go-cid"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/token"
)

type Server struct {
	Address string

	mu          sync.RWMutex
	subscribers map[string]map[*clientConn]bool
}

type clientConn struct {
	role string
	conn net.Conn
	send sync.Mutex
}

func (server *Server) Run(ctx context.Context) (runErr error) {
	// Intent: The kernel only inspects registration and slot 0 routing state, then
	// forwards exact envelope bytes to every other subscriber for that pCID. Source: DI-sabol; DI-nonad
	listener, err := net.Listen("tcp", server.Address)
	if err != nil {
		return fmt.Errorf("listen kernel: %w", err)
	}
	defer func() {
		if closeErr := listener.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) && runErr == nil {
			runErr = closeErr
		}
	}()
	server.subscribers = map[string]map[*clientConn]bool{}
	go func() {
		<-ctx.Done()
		if closeErr := listener.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
			fmt.Fprintf(os.Stderr, "kernel listener close: %v\n", closeErr)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("accept kernel conn: %w", err)
		}
		go server.handleConn(ctx, conn)
	}
}

func (server *Server) handleConn(ctx context.Context, conn net.Conn) {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "kernel conn close: %v\n", closeErr)
		}
	}()
	registerCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	registerBytes, err := readFrame(registerCtx, conn)
	if err != nil {
		return
	}
	registerEnvelope, err := protocol.ParseEnvelope(registerBytes)
	if err != nil {
		return
	}
	if registerEnvelope.PCID.String() != protocol.KernelRegisterProfile.CID.String() {
		return
	}
	var register protocol.KernelRegisterMessage
	if err := protocol.Unmarshal(registerEnvelope.PayloadBytes, &register); err != nil {
		return
	}
	signable, err := registerEnvelope.SignableBytes()
	if err != nil {
		return
	}
	if err := token.VerifyProof(register.Role, signable, registerEnvelope.ProofBytes); err != nil {
		return
	}
	client := &clientConn{role: register.Role, conn: conn}
	server.addClient(client, register.ReceivePCIDs)
	defer server.removeClient(client)
	for {
		envelopeBytes, err := readFrame(ctx, conn)
		if err != nil {
			return
		}
		envelope, err := protocol.ParseEnvelope(envelopeBytes)
		if err != nil {
			return
		}
		server.broadcast(client, envelope.PCID, envelopeBytes)
	}
}

func (server *Server) addClient(client *clientConn, receivePCIDs [][]byte) {
	server.mu.Lock()
	defer server.mu.Unlock()
	for _, raw := range receivePCIDs {
		pcid, err := cid.Cast(raw)
		if err != nil {
			continue
		}
		key := pcid.String()
		if server.subscribers[key] == nil {
			server.subscribers[key] = map[*clientConn]bool{}
		}
		server.subscribers[key][client] = true
	}
}

func (server *Server) removeClient(client *clientConn) {
	server.mu.Lock()
	defer server.mu.Unlock()
	for key, subscribers := range server.subscribers {
		delete(subscribers, client)
		if len(subscribers) == 0 {
			delete(server.subscribers, key)
		}
	}
}

func (server *Server) broadcast(sender *clientConn, pcid cid.Cid, envelopeBytes []byte) {
	server.mu.RLock()
	targets := make([]*clientConn, 0, len(server.subscribers[pcid.String()]))
	for client := range server.subscribers[pcid.String()] {
		if client != sender {
			targets = append(targets, client)
		}
	}
	server.mu.RUnlock()
	for _, target := range targets {
		target.send.Lock()
		if err := writeFrame(context.Background(), target.conn, envelopeBytes); err != nil {
			fmt.Fprintf(os.Stderr, "kernel broadcast to %s failed: %v\n", target.role, err)
		}
		target.send.Unlock()
	}
}

func writeFrame(ctx context.Context, conn net.Conn, payload []byte) error {
	return agentWriteFrame(ctx, conn, payload)
}

func readFrame(ctx context.Context, conn net.Conn) ([]byte, error) {
	return agentReadFrame(ctx, conn)
}

func agentWriteFrame(ctx context.Context, conn net.Conn, payload []byte) error {
	return writeFrameImpl(ctx, conn, payload)
}

func agentReadFrame(ctx context.Context, conn net.Conn) ([]byte, error) {
	return readFrameImpl(ctx, conn)
}
