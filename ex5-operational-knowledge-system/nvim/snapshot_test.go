package nvim

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNeovimSnapshotCreatesDurableRevisionFromLiveDraft(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	type livePayload struct {
		UpdateBody bool   `json:"update_body"`
		Body       string `json:"body"`
	}
	type revisionPayload struct {
		Actor   string   `json:"actor"`
		Title   string   `json:"title"`
		Summary string   `json:"summary"`
		Body    string   `json:"body"`
		Tags    []string `json:"tags"`
	}

	var pushedBody string
	var liveHTTPBodyPushSeen bool
	var revisionSeen bool
	socketBodies := make(chan string, 4)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/items/ITEM-0001/live":
			if _, err := fmt.Fprint(response, `{
				"item_id":"ITEM-0001",
				"title":"Start line A",
				"status":"draft",
				"version":4,
				"current_revision":2,
				"body":"old line",
				"participants":[]
			}`); err != nil {
				t.Fatalf("write live get response: %v", err)
			}
		case request.Method == http.MethodPost && request.URL.Path == "/api/items/ITEM-0001/live":
			var payload livePayload
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode live payload: %v", err)
			}
			if payload.UpdateBody {
				liveHTTPBodyPushSeen = true
				pushedBody = payload.Body
			}
			if _, err := fmt.Fprintf(response, `{
				"item_id":"ITEM-0001",
				"title":"Start line A",
				"status":"draft",
				"version":5,
				"current_revision":2,
				"body":%q,
				"participants":[]
			}`, map[bool]string{true: payload.Body, false: "old line"}[payload.UpdateBody]); err != nil {
				t.Fatalf("write live post response: %v", err)
			}
		case request.Method == http.MethodGet && request.URL.Path == "/api/items/ITEM-0001/live/socket":
			socket, err := acceptTestWebSocket(response, request)
			if err != nil {
				t.Fatalf("accept test websocket: %v", err)
			}
			defer func() {
				_ = socket.Close()
			}()
			if err := socket.WriteJSON(`{"type":"live-state","state":{"item_id":"ITEM-0001","title":"Start line A","status":"draft","version":4,"current_revision":2,"body":"old line","participants":[]}}`); err != nil {
				t.Fatalf("write websocket state: %v", err)
			}
			for {
				payload, err := socket.ReadJSON()
				if err != nil {
					return
				}
				if strings.Contains(payload, `"update_body":true`) {
					select {
					case socketBodies <- payload:
					default:
					}
					return
				}
			}
		case request.Method == http.MethodGet && request.URL.Path == "/api/items/ITEM-0001":
			if _, err := fmt.Fprint(response, `{
				"id":"ITEM-0001",
				"kind":"procedure",
				"title":"Start line A",
				"summary":"Startup checklist",
				"status":"draft",
				"current_revision":2,
				"working_version":5,
				"tags":["startup","audit"],
				"responsibility_ids":[],
				"revisions":[],
				"approvals":[],
				"related_runs":[],
				"links":[]
			}`); err != nil {
				t.Fatalf("write item get response: %v", err)
			}
		case request.Method == http.MethodPost && request.URL.Path == "/api/items/ITEM-0001/revisions":
			var payload revisionPayload
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode revision payload: %v", err)
			}
			if payload.Actor != "Boss" {
				t.Fatalf("unexpected actor: %#v", payload)
			}
			if payload.Title != "Start line A" || payload.Summary != "Startup checklist" {
				t.Fatalf("unexpected title/summary: %#v", payload)
			}
			if payload.Body != "new line 1\nnew line 2" {
				t.Fatalf("unexpected revision body: %#v", payload)
			}
			if strings.Join(payload.Tags, ",") != "startup,audit" {
				t.Fatalf("unexpected tags: %#v", payload)
			}
			revisionSeen = true
			if _, err := fmt.Fprint(response, `{
				"id":"ITEM-0001",
				"kind":"procedure",
				"title":"Start line A",
				"summary":"Startup checklist",
				"status":"draft",
				"current_revision":3,
				"working_version":5,
				"tags":["startup","audit"],
				"responsibility_ids":[],
				"revisions":[
					{"number":1,"title":"Start line A","summary":"Initial","created_at":"2026-07-21T09:00:00Z"},
					{"number":2,"title":"Start line A","summary":"Updated","created_at":"2026-07-21T10:00:00Z"},
					{"number":3,"title":"Start line A","summary":"Startup checklist","created_at":"2026-07-22T01:15:00Z"}
				],
				"approvals":[],
				"related_runs":[],
				"links":[]
			}`); err != nil {
				t.Fatalf("write revision response: %v", err)
			}
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "snapshot.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
vim.env.OKS_DISPLAY_NAME = "Boss"
local oks = require("oks")
oks.setup()

vim.cmd("OksOpen ITEM-0001")
vim.wait(2000, function()
  return oks.state.transport == "websocket"
end, 50)
if oks.state.transport ~= "websocket" then
  error("websocket transport did not connect")
end
vim.api.nvim_buf_set_lines(0, 0, -1, false, { "new line 1", "new line 2" })
vim.cmd("OksSnapshot")

if oks.state.current_revision ~= 3 then
  error("unexpected revision " .. tostring(oks.state.current_revision))
end
local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if body ~= "new line 1\nnew line 2" then
  error("unexpected live body " .. body)
end
vim.cmd("qa!")
`, server.URL)
	if err := os.WriteFile(script, []byte(scriptBody), 0o644); err != nil {
		t.Fatalf("write script: %v", err)
	}

	command := exec.Command(
		nvimPath,
		"--headless",
		"-u", "NONE",
		"-c", "set runtimepath+=.",
		"-l", script,
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("nvim snapshot regression: %v\n%s", err, string(output))
	}
	select {
	case payload := <-socketBodies:
		if !strings.Contains(payload, `"body":"new line 1\nnew line 2"`) && !strings.Contains(payload, "new line 1\\nnew line 2") {
			t.Fatalf("websocket live draft push did not carry updated body: %q", payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("expected websocket live draft push before revision snapshot")
	}
	if liveHTTPBodyPushSeen {
		t.Fatalf("unexpected HTTP live draft body push while websocket transport was available")
	}
	if pushedBody != "" {
		t.Fatalf("expected websocket live draft push instead of HTTP body push, got %q", pushedBody)
	}
	if !revisionSeen {
		t.Fatalf("revision POST was not observed")
	}
	if strings.Contains(string(output), "unexpected revision") || strings.Contains(string(output), "unexpected live body") {
		t.Fatalf("unexpected snapshot output: %s", string(output))
	}
}

type testWebSocketConn struct {
	raw   interface{ Close() error }
	read  func() (string, error)
	write func(string) error
}

func (socket *testWebSocketConn) Close() error {
	return socket.raw.Close()
}

func (socket *testWebSocketConn) ReadJSON() (string, error) {
	return socket.read()
}

func (socket *testWebSocketConn) WriteJSON(payload string) error {
	return socket.write(payload)
}

func acceptTestWebSocket(writer http.ResponseWriter, request *http.Request) (*testWebSocketConn, error) {
	hijacker, ok := writer.(http.Hijacker)
	if !ok {
		return nil, http.ErrNotSupported
	}
	conn, buffer, err := hijacker.Hijack()
	if err != nil {
		return nil, err
	}
	key := request.Header.Get("Sec-WebSocket-Key")
	sum := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	accept := base64.StdEncoding.EncodeToString(sum[:])
	response := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
	if _, err := buffer.WriteString(response); err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := buffer.Flush(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &testWebSocketConn{
		raw: conn,
		read: func() (string, error) {
			return readTestWebSocketFrame(conn)
		},
		write: func(payload string) error {
			return writeTestWebSocketFrame(conn, payload)
		},
	}, nil
}

func readTestWebSocketFrame(conn interface {
	Read([]byte) (int, error)
}) (string, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		return "", err
	}
	length := int(header[1] & 0x7f)
	if length == 126 {
		extended := make([]byte, 2)
		if _, err := io.ReadFull(conn, extended); err != nil {
			return "", err
		}
		length = int(binary.BigEndian.Uint16(extended))
	} else if length == 127 {
		extended := make([]byte, 8)
		if _, err := io.ReadFull(conn, extended); err != nil {
			return "", err
		}
		length = int(binary.BigEndian.Uint64(extended))
	}
	masked := header[1]&0x80 != 0
	maskKey := make([]byte, 4)
	if masked {
		if _, err := io.ReadFull(conn, maskKey); err != nil {
			return "", err
		}
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return "", err
	}
	if masked {
		for i := range payload {
			payload[i] ^= maskKey[i%4]
		}
	}
	return string(payload), nil
}

func writeTestWebSocketFrame(conn interface {
	Write([]byte) (int, error)
}, payload string) error {
	body := []byte(payload)
	frame := []byte{0x81}
	switch {
	case len(body) < 126:
		frame = append(frame, byte(len(body)))
	case len(body) <= 0xffff:
		frame = append(frame, 126, 0, 0)
		binary.BigEndian.PutUint16(frame[len(frame)-2:], uint16(len(body)))
	default:
		frame = append(frame, 127, 0, 0, 0, 0, 0, 0, 0, 0)
		binary.BigEndian.PutUint64(frame[len(frame)-8:], uint64(len(body)))
	}
	frame = append(frame, body...)
	_, err := conn.Write(frame)
	return err
}
