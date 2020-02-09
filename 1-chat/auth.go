package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

func loginRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")

		// not authenticated
		if err == http.ErrNoCookie || cookie.Value == "" {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		// some other error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		action := params["action"]
		provider := params["provider"]

		switch action {
		case "login":
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			loginURL, err := provider.GetBeginAuthURL(nil, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Add("Location", loginURL)
			w.WriteHeader(http.StatusTemporaryRedirect)

		case "callback":
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			user, err := provider.GetUser(creds)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			authCookieValue := objx.Map(map[string]interface{}{
				"name":       user.Name(),
				"avatar_url": user.AvatarURL(),
			}).MustBase64()
			http.SetCookie(w, &http.Cookie{
				Name:  "auth",
				Value: authCookieValue,
				Path:  "/",
			})
			w.Header().Add("Location", "/chat")
			w.WriteHeader(http.StatusTemporaryRedirect)

		default:
			http.Error(w, fmt.Sprintf("Auth action %s not supported", action), http.StatusNotFound)
		}
	}
}

func logout(redirectURI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "", // clear the value
			Path:   "/",
			MaxAge: -1, // so browser will delete this cookie automatically
		})
		w.Header().Set("Location", redirectURI)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
