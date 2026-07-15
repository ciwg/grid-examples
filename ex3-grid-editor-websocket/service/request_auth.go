package service

import (
	"net"
	"net/http"
)

func requestIsLoopback(request *http.Request) bool {
	if request == nil || request.RemoteAddr == "" {
		return true
	}
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		host = request.RemoteAddr
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return host == "localhost"
	}
	return ip.IsLoopback()
}
