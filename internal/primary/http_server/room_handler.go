package http_server

import (
	"github.com/gorilla/websocket"
	"github.com/rodrigogrohl/go-chat-app/internal/chat"
	"github.com/stretchr/objx"
	"log"
	"net/http"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrades = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

type RoomHandler struct {
	Room *chat.Room
}

func (r *RoomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrades.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie: ", err)
		return
	}

	client := &chat.Client{
		Socket:   socket,
		Send:     make(chan *chat.Message, messageBufferSize),
		Room:     r.Room,
		UserData: objx.MustFromBase64(authCookie.Value),
	}

	r.Room.Join <- client
	defer func() { r.Room.Leave <- client }()
	go client.Write()
	client.Read()
}
