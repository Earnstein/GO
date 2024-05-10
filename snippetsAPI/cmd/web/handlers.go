package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/earnstein/GO/snippetsAPI/internal/models"
)

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}

	files := []string{
		"./ui/html/pages/home.tmpl.html", "./ui/html/base.tmpl.html", "./ui/html/partials/nav.tmpl.html",
	}

	templates, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err = templates.ExecuteTemplate(w, "base", nil); err != nil {
		app.serverError(w, err)
	}
}


func(app *Application) handleSnippetCreate(w http.ResponseWriter, r *http.Request) {
	var reqBody models.SnippetBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	
	id , err := app.snippets.Insert(reqBody.Title, reqBody.Content, reqBody.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	jsonResponse := map[string]string{"message": "movie snippets created successfully", "id": fmt.Sprintf("%v", id), "title": reqBody.Title}
	msg, err := json.Marshal(jsonResponse)
	if err != nil {
		app.serverError(w, err)
		return
	}
	
	fmt.Fprintln(w, string(msg))
}

func(app *Application) handleSnippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	fmt.Fprintf(w, "You sent an id: %d...", id)
}
