//go:build !windows

package ipc

import (
	"net"
	"os"
	"path/filepath"
)

func endpoint() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pass", "agent.sock")
}

func listen() (net.Listener, error) {
	path := endpoint()
	os.MkdirAll(filepath.Dir(path), 0700)
	os.Remove(path)
	return net.Listen("unix", path)
}

func dial() (net.Conn, error) {
	return net.Dial("unix", endpoint())
}
