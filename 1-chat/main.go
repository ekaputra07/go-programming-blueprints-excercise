package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ekaputra07/goblueprints-excercise/tracer"
	"github.com/gorilla/mux"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
)

const (
	// ServerPort is the port of our chat app will running
	ServerPort = 8080
)

func main() {
	// Setup gomniauth
	gomniauth.SetSecurityKey("b53533b9-c113-461b-8944-ac9be26d12c5")
	gomniauth.WithProviders(
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), "http://localhost:8080/auth/callback/github"),
	)

	r := newRoom()
	r.tracer = tracer.New(os.Stdout)

	router := mux.NewRouter()
	router.HandleFunc("/chat", loginRequired(render("chat.html")))
	router.HandleFunc("/login", render("login.html"))
	router.HandleFunc("/auth/{action}/{provider}", handleLogin())
	router.Handle("/room", r)

	http.Handle("/", router)

	go r.run()

	log.Println("Listening on port 8080...")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", ServerPort), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
