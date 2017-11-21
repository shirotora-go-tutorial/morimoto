package main

import (
	"net/http"
	"strings"
	"log"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"crypto/md5"
	"io"
	gomniauthcommn "github.com/stretchr/gomniauth/common"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommn.User
	uniqueID string
}

func (u chatUser) UniqueID() string{
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == ""{
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("failed to get auth provider", provider, "-", err)
		}

		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("failed to error while calling GetBeginAuthURL", provider, "-", err)
		}

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatal("failed to get auth provider", provider, "-", err)
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatal("failed to finishsed auth", provider, "-", err)
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatal("failed to get user", provider, "-", err)
		}

		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Name()))
		userID := fmt.Sprintf("%x", m.Sum(nil))

		authCookieValue := objx.New(map[string]interface{}{
			"userid" : userID,
			"name": user.Name(),
			"avatar_url": user.AvatarURL(),
			"email" : user.Email(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: authCookieValue,
			Path: "/",
		})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "action %s is not incompatible")
	}
}