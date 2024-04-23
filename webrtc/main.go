package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	setUpAPI()
	fmt.Println("Listening on localhost port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setUpAPI() {
	manager := NewManager()
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.ServeWS)
}
