package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ekaputra07/goblueprints-excercise/tracer"
)

const (
	// ServerPort is the port of our chat app will running
	ServerPort = 8080
)

func main() {
	r := newRoom()
	r.tracer = tracer.New(os.Stdout)

	http.Handle("/", renderTemplate("index.html"))
	http.Handle("/room", r)

	go r.run()

	log.Println("Listening on port 8080...")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", ServerPort), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
