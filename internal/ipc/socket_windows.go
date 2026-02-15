//go:build windows

package ipc

import (
	"net"

	"github.com/Microsoft/go-winio"
)

func endpoint() string {
	return `\\.\pipe\passman-agent`
}

func listen() (net.Listener, error) {
	return winio.ListenPipe(endpoint(), &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;WD)",
	})
}

func dial() (net.Conn, error) {
	return winio.DialPipe(endpoint(), nil)
}
