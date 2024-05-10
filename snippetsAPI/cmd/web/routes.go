package main

import "net/http"


func (app *Application) routes() *http.ServeMux {
	server := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	server.Handle("/static/", http.StripPrefix("/static", fileServer))

	server.HandleFunc("/", app.homeHandler)
	server.HandleFunc("/snippet/create", app.handleSnippetCreate)
	server.HandleFunc("/snippet/view", app.handleSnippetView)
	server.HandleFunc("/snippet/latest", app.handleSnippetList)
	
	return server
}