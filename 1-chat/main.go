package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ekaputra07/goblueprints-excercise/tracer"
	"github.com/gorilla/mux"
)

const (
	// ServerPort is the port of our chat app will running
	ServerPort = 8080
)

func main() {
	r := newRoom()
	r.tracer = tracer.New(os.Stdout)

	router := mux.NewRouter()
	router.HandleFunc("/chat", loginRequired(render("chat.html")))
	router.HandleFunc("/login", render("login.html"))
	router.HandleFunc("/auth/{action}/{provider}", handleLogin())
	router.HandleFunc("/room", r.start())
	http.Handle("/", router)

	go r.run()

	log.Println("Listening on port 8080...")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", ServerPort), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
