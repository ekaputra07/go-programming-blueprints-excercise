package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stretchr/gomniauth"
)

func loginRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("auth")

		// not authenticated
		if err == http.ErrNoCookie {
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
			return

		case "callback":
			log.Println("TODO: handle callback action")
		default:
			http.Error(w, fmt.Sprintf("Auth action %s not supported", action), http.StatusNotFound)
		}
	}
}
