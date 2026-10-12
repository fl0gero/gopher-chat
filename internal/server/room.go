package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	clients map[*client]bool
	join    chan *client
	leave   chan *client
	forward chan Message
}

func newRoom() *room {
	return &room{
		clients: make(map[*client]bool),
		join:    make(chan *client),
		leave:   make(chan *client),
		forward: make(chan Message),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case message := <-r.forward:
			for client := range r.clients {
				select {
				case client.receive <- message:
				default:

				}
			}

		}

	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (reg *Registry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	name := req.PathValue("name")
	r := reg.getOrCreate(name)

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &client{
		name:    req.URL.Query().Get("username"),
		socket:  socket,
		receive: make(chan Message, messageBufferSize),
		room:    r,
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
