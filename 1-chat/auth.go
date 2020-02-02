package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
		// provider := params["provider"]

		switch action {
		case "login":
			log.Println("TODO: handle login action")
		case "callback":
			log.Println("TODO: handle callback action")
		default:
			http.Error(w, fmt.Sprintf("Auth action %s not supported", action), http.StatusNotFound)
		}
	}
}
