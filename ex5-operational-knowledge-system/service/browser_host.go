package service

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type BrowserHostEnvelope struct {
	RequestID  string                 `json:"request_id,omitempty"`
	SocketPath string                 `json:"socket_path,omitempty"`
	Request    LocalEmbodimentRequest `json:"request"`
}

type BrowserHostResponse struct {
	RequestID string                  `json:"request_id,omitempty"`
	Response  LocalEmbodimentResponse `json:"response,omitempty"`
	Error     string                  `json:"error,omitempty"`
}

type BrowserHost struct{}

func NewBrowserHost() *BrowserHost {
	return &BrowserHost{}
}

// Intent: Keep the Chrome/Chromium native host as a thin carriage bridge into
// the existing direct local runtime contract instead of creating a second
// browser-specific semantic layer. Source: DI-punek
func (host *BrowserHost) ServeSession(stdin io.Reader, stdout io.Writer) error {
	reader := newNativeMessageReader(stdin)
	writer := newNativeMessageWriter(stdout)
	firstMessage, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	var envelope BrowserHostEnvelope
	if err := json.Unmarshal(firstMessage, &envelope); err != nil {
		return err
	}
	if envelope.Request.Type == "live-open" {
		return host.serveLiveSession(reader, writer, envelope)
	}
	return host.serveRPC(writer, envelope)
}

func (host *BrowserHost) serveRPC(writer *nativeMessageWriter, envelope BrowserHostEnvelope) error {
	response, err := host.forwardOneShot(envelope)
	if err != nil {
		return writer.Write(BrowserHostResponse{RequestID: envelope.RequestID, Error: err.Error()})
	}
	return writer.Write(BrowserHostResponse{RequestID: envelope.RequestID, Response: response})
}

func (host *BrowserHost) serveLiveSession(reader *nativeMessageReader, writer *nativeMessageWriter, envelope BrowserHostEnvelope) error {
	conn, err := host.dialRuntimeSocket(envelope.SocketPath)
	if err != nil {
		return writer.Write(BrowserHostResponse{RequestID: envelope.RequestID, Error: err.Error()})
	}
	defer func() {
		_ = conn.Close()
	}()
	if err := json.NewEncoder(conn).Encode(envelope.Request); err != nil {
		return writer.Write(BrowserHostResponse{RequestID: envelope.RequestID, Error: err.Error()})
	}

	readDone := make(chan error, 1)
	go func() {
		readDone <- host.forwardSocketStream(conn, writer, envelope.RequestID)
	}()

	for {
		payload, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		var next BrowserHostEnvelope
		if err := json.Unmarshal(payload, &next); err != nil {
			return writer.Write(BrowserHostResponse{RequestID: envelope.RequestID, Error: err.Error()})
		}
		if next.Request.Type == "live-close" {
			return nil
		}
		if err := json.NewEncoder(conn).Encode(next.Request); err != nil {
			return writer.Write(BrowserHostResponse{RequestID: envelope.RequestID, Error: err.Error()})
		}
		select {
		case err := <-readDone:
			if err != nil && !errors.Is(err, io.EOF) {
				return err
			}
			return nil
		default:
		}
	}
}

func (host *BrowserHost) forwardSocketStream(conn net.Conn, writer *nativeMessageWriter, requestID string) error {
	decoder := json.NewDecoder(bufio.NewReader(conn))
	for {
		var response LocalEmbodimentResponse
		if err := decoder.Decode(&response); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			writeErr := writer.Write(BrowserHostResponse{RequestID: requestID, Error: err.Error()})
			if writeErr != nil {
				return writeErr
			}
			return err
		}
		if err := writer.Write(BrowserHostResponse{RequestID: requestID, Response: response}); err != nil {
			return err
		}
	}
}

func (host *BrowserHost) forwardOneShot(envelope BrowserHostEnvelope) (LocalEmbodimentResponse, error) {
	conn, err := host.dialRuntimeSocket(envelope.SocketPath)
	if err != nil {
		return LocalEmbodimentResponse{}, err
	}
	defer func() {
		_ = conn.Close()
	}()
	if err := json.NewEncoder(conn).Encode(envelope.Request); err != nil {
		return LocalEmbodimentResponse{}, err
	}
	var response LocalEmbodimentResponse
	if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&response); err != nil {
		return LocalEmbodimentResponse{}, err
	}
	return response, nil
}

func (host *BrowserHost) dialRuntimeSocket(socketPath string) (net.Conn, error) {
	if socketPath == "" {
		return nil, fmt.Errorf("browser host requires socket_path")
	}
	return net.DialTimeout("unix", socketPath, 2*time.Second)
}

type nativeMessageReader struct {
	reader io.Reader
}

func newNativeMessageReader(reader io.Reader) *nativeMessageReader {
	return &nativeMessageReader{reader: reader}
}

func (reader *nativeMessageReader) Read() ([]byte, error) {
	var size uint32
	if err := binary.Read(reader.reader, binary.LittleEndian, &size); err != nil {
		return nil, err
	}
	payload := make([]byte, size)
	if _, err := io.ReadFull(reader.reader, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

type nativeMessageWriter struct {
	writer io.Writer
	mu     sync.Mutex
}

func newNativeMessageWriter(writer io.Writer) *nativeMessageWriter {
	return &nativeMessageWriter{writer: writer}
}

func (writer *nativeMessageWriter) Write(message BrowserHostResponse) error {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	if err := binary.Write(writer.writer, binary.LittleEndian, uint32(len(payload))); err != nil {
		return err
	}
	_, err = writer.writer.Write(payload)
	return err
}
