package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/earnstein/GO/snippetsAPI/cmd/controllers"
)

var PORT = ":5000"

func main() {

	addr := flag.String("addr", PORT, "HTTP network port address")
	flag.Parse()

	
	server := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	server.Handle("/static/", http.StripPrefix("/static", fileServer))
	server.HandleFunc("/", controllers.HomeHandler)
	server.HandleFunc("/create", controllers.HandleSnippetCreate)
	server.HandleFunc("/snippet/view", controllers.HandleSnippetView)
	log.Printf("server is listening on port %s", *addr)
	err := http.ListenAndServe(*addr, server)
	log.Fatal(err)
}
