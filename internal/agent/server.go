package agent

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/subrat-dwi/passman-cli/internal/ipc"
)

func Run() {
	startServer(10 * time.Minute)
}

func startServer(ttl time.Duration) {
	// Check if agent is already running
	if conn, err := ipc.Dial(); err == nil {
		conn.Close()
		log.Println("Agent already running, exiting")
		return
	}

	// Start listening
	l, err := ipc.Listen()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("Agent listening on", ipc.Endpoint())

	// Initialize state
	state := NewState(ttl)
	go state.AutoLock()

	//Accept loop
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn, state)
	}
}

func handleConn(conn net.Conn, state *State) {
	defer conn.Close()

	var req Request
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		json.NewEncoder(conn).Encode(Response{
			OK:    false,
			Error: "invalid request",
		})
		return
	}

	log.Printf("Received command: %s\n", req.Cmd)
	resp := state.Handle(req)
	log.Printf("Response: OK=%v, Error=%s\n", resp.OK, resp.Error)
	json.NewEncoder(conn).Encode(resp)
}
