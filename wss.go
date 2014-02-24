package main

import (
	"github.com/gorilla/mux"

	"code.google.com/p/go.net/websocket"
)

var sessions map[string]students

type students map[*websocket.Conn]bool

func init() {
	sessions = make(map[string]students)
}

// Echo the data received on the WebSocket.
func studServer(ws *websocket.Conn) {
	for {
		id := mux.Vars(ws.Request())["id"]
		s, ok := sessions[id]

		if !ok {
			ws.Close()
			continue
		}

		if _, ok := s[ws]; !ok {
			s[ws] = true
		}

		var data map[string]interface{}
		err := websocket.JSON.Receive(ws, &data)
		if err != nil { // io.EOF = disconnect
			break
		}
		websocket.JSON.Send(ws, &data)
	}
}

func teachServer(ws *websocket.Conn) {
	for {
		id := mux.Vars(ws.Request())["id"]
		if _, ok := sessions[id]; !ok {
			//accepting == true
			sessions[id] = nil
		}

		var data map[string]interface{}
		err := websocket.JSON.Receive(ws, &data)
		if err != nil { // io.EOF = disconnect
			break
		}
	}
}
