package main

import "net/http"


func (app *Application) routes() *http.ServeMux {
	server := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	server.Handle("/static/", http.StripPrefix("/static", fileServer))

	server.HandleFunc("/", app.HomeHandler)
	server.HandleFunc("/create", app.HandleSnippetCreate)
	server.HandleFunc("/snippet/view", app.HandleSnippetView)
	
	return server
}