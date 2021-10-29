package http_server

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"io"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// Not authenticated
		w.Header().Set("Location", "login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// success - call next handler
	fmt.Println(cookie)
	h.next.ServeHTTP(w, r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// LoginHandler handles the third-party login process.
// format: /auth/{action}/{provider}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	platform := segs[3]
	switch action {
	case "login":
		log.Println("handling login with ", platform)
		provider, err := gomniauth.Provider(platform)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Error when trying to get provider %s: %s", provider, err),
				http.StatusBadRequest)
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Error when trying to get provider %s: %s", provider, err),
				http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(platform)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Error when trying to get provider %s: %s", provider, err),
				http.StatusBadRequest)
			return
		}
		credentials, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error trying to complete auth for %s: %s", platform, err),
				http.StatusInternalServerError)
			return
		}
		user, err := provider.GetUser(credentials)
		// From Google: id, email, verified_email, name, given_name, family_name, picture, locale
		// From Github: id, url, location, name, avatar_url, email, login,
		if err != nil {
			http.Error(w, fmt.Sprintf("Error trying to get User from %s: %s", platform, err),
				http.StatusInternalServerError)
			return
		}

		m := md5.New()
		_, _ = io.WriteString(m, strings.ToLower(user.Email()))
		userId := fmt.Sprintf("%x", m.Sum(nil))

		authCookieValue := objx.New(map[string]interface{}{
			"user_id": userId,
			"name": user.Name(),
			"email": user.Email(),
			"avatar_url": user.AvatarURL(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:       "auth",
			Value:      authCookieValue,
			Path:       "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
