package chat

import (
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	// Socket is the websocket for this Client
	Socket *websocket.Conn

	// Send is a channel on which messages are sent
	Send chan *Message

	// Room is the Room this Client is chatting it
	Room *Room

	// UserData holds information about the user
	UserData map[string] interface{}
}

func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		var msg *Message
		err := c.Socket.ReadJSON(&msg)
		//_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.UserData["name"].(string)
		msg.AvatarURL, _ = c.Room.avatar.GetAvatarURL(c)
		c.Room.Forward <- msg
	}
}

func (c *Client) Write(){
	defer c.Socket.Close()
	for msg := range c.Send {
		//err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}