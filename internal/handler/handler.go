package handler

import (
	"encoding/json"
	"fmt"
	"gopher-chat/internal/server"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectToRoom establishes a WebSocket connection with the server for a specific
// room (roomName) and user (userName), and launches a background goroutine
// that continuously listens for incoming messages, parses them from JSON, and passes them
// to the onMessage callback. Returns an open connection for subsequent use.
func ConnectToRoom(roomName string, userName string, onMessage func(server.Message)) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/room/"+roomName+"?username="+userName, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading from websocket:", err)
				return
			}
			var msg server.Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				fmt.Println("Error unmarshalling from websocket:", err)
			}
			fmt.Println(msg.From, msg.Text)
			onMessage(msg)
		}
	}()

	return conn, nil
}

func SendMessage(conn *websocket.Conn, text string) error {
	return conn.WriteMessage(websocket.TextMessage, []byte(text))
}

// ConnectionPool stores already open WebSocket connections by room name
type ConnectionPool struct {
	mu    sync.Mutex
	conns map[string]*websocket.Conn
}

// NewConnectionPool creates a connection pool with an initialized map
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{mu: sync.Mutex{}, conns: make(map[string]*websocket.Conn)}
}

func (p *ConnectionPool) GetOrConnect(name string, userName string, onMessage func(message server.Message)) (*websocket.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn, ok := p.conns[name]
	if ok {
		return conn, nil
	}

	newConn, err := ConnectToRoom(name, userName, onMessage)
	if err != nil {
		return nil, err
	}
	p.conns[name] = newConn

	return newConn, nil
}
