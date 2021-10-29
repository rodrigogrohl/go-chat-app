package main

import (
	"github.com/rodrigogrohl/go-chat-app/configs"
	"github.com/rodrigogrohl/go-chat-app/internal/chat"
	"github.com/rodrigogrohl/go-chat-app/internal/primary/http_server"
	"github.com/rodrigogrohl/go-chat-app/internal/utils/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"log"
	"net/http"
	"os"
)

func main() {
	gomniauth.SetSecurityKey("20170403063245-04a5b1e1910d")
	gomniauth.WithProviders(
		google.New("555164498142-ecebtbfosg4gv33iiilr80mc93ujfhch.apps.googleusercontent.com",
			"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",
			"http://localhost:8080/auth/callback/google"),
		github.New("effa1d148cd04e15e553",
			"4777604c3f43a7533d38845e28d2fe8e0e573140",
			"http://localhost:8080/auth/callback/github"),
	)

	r := chat.NewRoom(chat.UseFileSystemAvatar)
	r.Tracer = trace.New(os.Stdout)

	// TODO: refactor
	http.Handle("/", http_server.MustAuth(http_server.New("chat.html")))
	http.HandleFunc("/auth/", http_server.LoginHandler)
	http.HandleFunc("/logout", http_server.Logout)
	http.Handle("/login", http_server.New("login.html"))
	http.Handle("/room", &http_server.RoomHandler{Room: r})
	http.Handle("/upload", http_server.New("upload.html"))
	http.HandleFunc("/uploader", http_server.UploaderHandler)
	http.Handle(configs.HttpAvatarGet, http.StripPrefix("/avatars/", http.FileServer(http.Dir("./web/avatars"))))
	go r.Run()

	// Start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listen and Serve: ", err)
	}
}
