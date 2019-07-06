package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return // this will break the loop and exit the function
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (c *client) start() {
	go c.write()
	c.read()
}

func newClient(socket *websocket.Conn, room *room) *client {
	return &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   room,
	}
}
