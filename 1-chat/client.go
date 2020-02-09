package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return // this will break the loop and exit the function
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			return
		}
	}
}

func (c *client) start() {
	go c.write()
	c.read()
}

func newClient(socket *websocket.Conn, room *room, userData map[string]interface{}) *client {
	return &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     room,
		userData: userData,
	}
}
