package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Response map[string]string

type Application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewApplicaton(infoLogger, errorLogger *log.Logger) *Application {
	return &Application{
		infoLog: infoLogger,
		errorLog: errorLogger,
	}
}
func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/pages/home.tmpl.html", "./ui/html/base.tmpl.html", "./ui/html/partials/nav.tmpl.html",
	}

	templates, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Printf("Error parsing file: %v",  err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err = templates.ExecuteTemplate(w, "base", nil); err != nil {
		app.errorLog.Printf("Error Executing templates: %v",  err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}


func(app *Application) HandleSnippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Allow", http.MethodPost)
		jsonData, _ := json.Marshal(Response{"error": "method not allowed"})
		http.Error(w, string(jsonData), http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("OK"))
}

func(app *Application) HandleSnippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "You sent an id: %d...", id)
}
