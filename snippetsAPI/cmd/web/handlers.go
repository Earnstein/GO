package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/earnstein/GO/snippetsAPI/internal/models"
)

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
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
	
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func(app *Application) handleSnippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	s , err := json.Marshal(snippet)
	 if err != nil {
		app.serverError(w, err)
	 }
	fmt.Fprint(w, string(s))
}




func(app *Application) handleLatestSnippet(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.GetLatest()

	if err != nil {
		app.serverError(w, err)
		return
	}
	data, _ := json.Marshal(snippets)

	// for _, snippet := range snippets {
	// 	data, err := json.Marshal(snippet)
	// 	if err != nil {
	// 		app.serverError(w, err)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(data)
	// }
	w.Write(data)

}