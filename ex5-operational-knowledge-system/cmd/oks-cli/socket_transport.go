package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/service"
)

var errLocalSocketUnavailable = errors.New("local unix socket unavailable")

// Intent: Prefer the direct local Unix-socket contract for the CLI while
// preserving HTTP fallback so shell workflows remain usable during mixed local
// deployments. Source: DI-favel
func (cli *CLI) roundTrip(method string, path string, contentType string, body []byte) (int, []byte, error) {
	if strings.TrimSpace(cli.SocketPath) != "" {
		status, payload, err := cli.localSocketRoundTrip(method, path, contentType, body)
		if err == nil {
			return status, payload, nil
		}
		if !errors.Is(err, errLocalSocketUnavailable) {
			return 0, nil, err
		}
	}
	return cli.httpRoundTrip(method, path, contentType, body)
}

func (cli *CLI) localSocketRoundTrip(method string, path string, contentType string, body []byte) (int, []byte, error) {
	if _, err := os.Stat(cli.SocketPath); err != nil {
		return 0, nil, errLocalSocketUnavailable
	}
	conn, err := net.DialTimeout("unix", cli.SocketPath, 2*time.Second)
	if err != nil {
		return 0, nil, errLocalSocketUnavailable
	}
	defer func() {
		_ = conn.Close()
	}()
	request := service.LocalEmbodimentRequest{
		Type:    "request",
		Method:  method,
		Path:    path,
		Headers: map[string]string{},
	}
	if contentType != "" {
		request.Headers["Content-Type"] = contentType
	}
	if len(body) > 0 {
		if strings.HasPrefix(contentType, "multipart/form-data") {
			request.BodyBase64 = base64.StdEncoding.EncodeToString(body)
		} else {
			request.Body = string(body)
		}
	}
	if err := json.NewEncoder(conn).Encode(request); err != nil {
		return 0, nil, err
	}
	var response service.LocalEmbodimentResponse
	if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&response); err != nil {
		return 0, nil, err
	}
	if response.Type == "error" {
		return 0, nil, errors.New(response.Message)
	}
	return response.Status, []byte(response.Body), nil
}

func (cli *CLI) httpRoundTrip(method string, path string, contentType string, body []byte) (int, []byte, error) {
	request, err := http.NewRequest(method, cli.ServerURL+path, strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, err
	}
	payload, err := readResponseBody(response)
	if err != nil {
		return 0, nil, err
	}
	return response.StatusCode, payload, nil
}
