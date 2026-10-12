package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Message struct {
	From string
	Text string
}

type client struct {
	name string

	socket *websocket.Conn

	receive chan Message

	room *room
}

func (c *client) read() {

	defer c.socket.Close()
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		c.room.forward <- Message{
			From: c.name,
			Text: string(message),
		}
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.receive {
		data, err := json.Marshal(msg)
		if err != nil {
			return
		}
		err = c.socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
	}
}
