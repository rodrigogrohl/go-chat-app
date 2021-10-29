package chat

import (
	"github.com/rodrigogrohl/go-chat-app/internal/utils/trace"
)

type Room struct {
	// Forward is a channel that holds incoming messages that should be forwarded to the other Clients
	Forward chan *Message

	// Join is a channel for Clients wishing to Join the Room
	Join chan *Client

	// Leave is a channel for Clients wishing to Leave the Room
	Leave chan *Client

	// Clients holds all Clients in this Room
	Clients map[*Client]bool

	// Tracer will receive trace information of activity in the Room
	Tracer trace.Tracer

	// avatar is how avatar information will be obtained
	avatar Avatar
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
			r.Tracer.Trace("New client joined: ", &client)
		case client := <-r.Leave:
			delete(r.Clients, client)
			close(client.Send)
			r.Tracer.Trace("Client left: ", &client)
		case msg := <-r.Forward:
			r.Tracer.Trace("Message received: ", msg.Message)
			for client := range r.Clients {
				client.Send <- msg
				r.Tracer.Trace(" -- sent to client: ", client)
			}
		}
	}
}

func NewRoom(avatar Avatar) *Room {
	return &Room{
		Forward: make(chan *Message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
		Tracer:  trace.Off(),
		avatar:  avatar,
	}
}
