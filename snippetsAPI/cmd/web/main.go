package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/earnstein/GO/snippetsAPI/cmd/controllers"
)



func main() {
	// command flags
	addr := flag.String("addr", ":5000", "HTTP network port address")
	flag.Parse()

	// loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := controllers.NewApplicaton(infoLog, errorLog)

	// server
	server := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	server.Handle("/static/", http.StripPrefix("/static", fileServer))
	server.HandleFunc("/", app.HomeHandler)
	server.HandleFunc("/create", app.HandleSnippetCreate)
	server.HandleFunc("/snippet/view", app.HandleSnippetView)
	infoLog.Printf("server is listening on port %s", *addr)
	err := http.ListenAndServe(*addr, server)
	errorLog.Fatal(err)
}
