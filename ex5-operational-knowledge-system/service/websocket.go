package service

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

const websocketGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

type websocketConn struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer

	writeMu sync.Mutex
}

// Intent: Keep ex5 websocket collaboration transport on the existing local
// adapter by upgrading the current HTTP server in place instead of creating a
// separate live-collaboration daemon. Source: DI-noruv
func upgradeWebSocket(writer http.ResponseWriter, request *http.Request) (*websocketConn, error) {
	if !headerContainsToken(request.Header, "Connection", "upgrade") || !headerContainsToken(request.Header, "Upgrade", "websocket") {
		return nil, fmt.Errorf("websocket upgrade required")
	}
	if strings.TrimSpace(request.Header.Get("Sec-WebSocket-Version")) != "13" {
		return nil, fmt.Errorf("unsupported websocket version")
	}
	key := strings.TrimSpace(request.Header.Get("Sec-WebSocket-Key"))
	if key == "" {
		return nil, fmt.Errorf("missing websocket key")
	}
	hijacker, ok := writer.(http.Hijacker)
	if !ok {
		return nil, fmt.Errorf("websocket hijack unsupported")
	}
	conn, readWriter, err := hijacker.Hijack()
	if err != nil {
		return nil, fmt.Errorf("websocket hijack: %w", err)
	}
	accept := websocketAccept(key)
	if _, err := readWriter.WriteString(
		"HTTP/1.1 101 Switching Protocols\r\n" +
			"Upgrade: websocket\r\n" +
			"Connection: Upgrade\r\n" +
			"Sec-WebSocket-Accept: " + accept + "\r\n\r\n",
	); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("write websocket handshake: %w", err)
	}
	if err := readWriter.Flush(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("flush websocket handshake: %w", err)
	}
	return &websocketConn{
		conn:   conn,
		reader: readWriter.Reader,
		writer: readWriter.Writer,
	}, nil
}

func (socket *websocketConn) Close() error {
	return socket.conn.Close()
}

func (socket *websocketConn) ReadJSON(value any) error {
	for {
		payload, opcode, err := socket.readFrame()
		if err != nil {
			return err
		}
		switch opcode {
		case 0x1:
			if err := json.Unmarshal(payload, value); err != nil {
				return fmt.Errorf("decode websocket json: %w", err)
			}
			return nil
		case 0x8:
			return io.EOF
		case 0x9:
			if err := socket.writeFrame(0xA, payload, false); err != nil {
				return err
			}
		case 0xA:
			continue
		default:
			return fmt.Errorf("unsupported websocket opcode %d", opcode)
		}
	}
}

func (socket *websocketConn) WriteJSON(value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal websocket json: %w", err)
	}
	return socket.writeFrame(0x1, payload, false)
}

func (socket *websocketConn) readFrame() ([]byte, byte, error) {
	first, err := socket.reader.ReadByte()
	if err != nil {
		return nil, 0, err
	}
	second, err := socket.reader.ReadByte()
	if err != nil {
		return nil, 0, err
	}
	if first&0x80 == 0 {
		return nil, 0, fmt.Errorf("fragmented websocket frames are unsupported")
	}
	opcode := first & 0x0f
	masked := second&0x80 != 0
	lengthValue := uint64(second & 0x7f)
	switch lengthValue {
	case 126:
		var extended uint16
		if err := binary.Read(socket.reader, binary.BigEndian, &extended); err != nil {
			return nil, 0, err
		}
		lengthValue = uint64(extended)
	case 127:
		if err := binary.Read(socket.reader, binary.BigEndian, &lengthValue); err != nil {
			return nil, 0, err
		}
	}
	var maskKey [4]byte
	if masked {
		if _, err := io.ReadFull(socket.reader, maskKey[:]); err != nil {
			return nil, 0, err
		}
	}
	payload := make([]byte, lengthValue)
	if _, err := io.ReadFull(socket.reader, payload); err != nil {
		return nil, 0, err
	}
	if masked {
		for index := range payload {
			payload[index] ^= maskKey[index%len(maskKey)]
		}
	}
	return payload, opcode, nil
}

func (socket *websocketConn) writeFrame(opcode byte, payload []byte, masked bool) error {
	socket.writeMu.Lock()
	defer socket.writeMu.Unlock()

	if err := socket.writer.WriteByte(0x80 | (opcode & 0x0f)); err != nil {
		return err
	}
	maskBit := byte(0)
	if masked {
		maskBit = 0x80
	}
	switch {
	case len(payload) < 126:
		if err := socket.writer.WriteByte(maskBit | byte(len(payload))); err != nil {
			return err
		}
	case len(payload) <= 0xffff:
		if err := socket.writer.WriteByte(maskBit | 126); err != nil {
			return err
		}
		if err := binary.Write(socket.writer, binary.BigEndian, uint16(len(payload))); err != nil {
			return err
		}
	default:
		if err := socket.writer.WriteByte(maskBit | 127); err != nil {
			return err
		}
		if err := binary.Write(socket.writer, binary.BigEndian, uint64(len(payload))); err != nil {
			return err
		}
	}
	if masked {
		maskKey := [4]byte{0x13, 0x57, 0x9b, 0xdf}
		if _, err := socket.writer.Write(maskKey[:]); err != nil {
			return err
		}
		maskedPayload := append([]byte(nil), payload...)
		for index := range maskedPayload {
			maskedPayload[index] ^= maskKey[index%len(maskKey)]
		}
		if _, err := socket.writer.Write(maskedPayload); err != nil {
			return err
		}
	} else {
		if _, err := socket.writer.Write(payload); err != nil {
			return err
		}
	}
	return socket.writer.Flush()
}

func websocketAccept(key string) string {
	sum := sha1.Sum([]byte(key + websocketGUID))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func headerContainsToken(header http.Header, name string, want string) bool {
	for _, value := range header.Values(name) {
		for _, token := range strings.Split(value, ",") {
			if strings.EqualFold(strings.TrimSpace(token), want) {
				return true
			}
		}
	}
	return false
}
