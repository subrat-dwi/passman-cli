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
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return net.Listen("unix", path)
}

func dial() (net.Conn, error) {
	return net.Dial("unix", endpoint())
}
