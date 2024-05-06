package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var PORT = "5000"

type Response map[string]string

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(Response{"message": "Hello World"})
	w.Write(jsonData)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Allow", http.MethodPost)
		jsonData, _ := json.Marshal(Response{"error": "method not allowed"})
		http.Error(w, string(jsonData), http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("OK"))
}



func main(){
	server := http.NewServeMux()
	server.HandleFunc("/", home)
	server.HandleFunc("/create", handleCreate)

	log.Printf("server is listening on port: %s", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), server)
	log.Fatal(err)
}