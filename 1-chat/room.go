package main

import (
	"log"
	"net/http"

	"github.com/ekaputra07/goblueprints-excercise/tracer"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  tracer.Tracer
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining client
			r.clients[client] = true
			r.tracer.Trace("A client joined the room")
		case client := <-r.leave:
			// leaving client
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("A client left the room")
		case msg := <-r.forward:
			// message received, forward to all clients
			r.tracer.Trace("A message forwarded to all clients, msg: ", string(msg))
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("room.ServeHTTP error: ", err)
		return
	}
	client := newClient(socket, r)
	defer func() { r.leave <- client }()

	r.join <- client
	client.start()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  tracer.Off(),
	}
}
