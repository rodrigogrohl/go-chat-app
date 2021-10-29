package http_server

import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:       "auth",
		Value:      "",
		Path:       "/",
		MaxAge:     -1,
	})
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
