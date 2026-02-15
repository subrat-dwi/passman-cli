package ipc

import (
	"net"
)

func Listen() (net.Listener, error) { return listen() }

func Dial() (net.Conn, error) { return dial() }

func Endpoint() string { return endpoint() }
