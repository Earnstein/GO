package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	setUpAPI()
	fmt.Println("Listening on localhost port: 8080")
	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
}

func setUpAPI() {
	ctx := context.Background()
	manager := NewManager(ctx)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.ServeWS)
	http.HandleFunc("/login", manager.LoginHandler)
}
