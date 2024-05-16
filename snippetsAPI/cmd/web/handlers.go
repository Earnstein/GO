package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/earnstein/GO/snippetsAPI/internal/models"
	"github.com/julienschmidt/httprouter"
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
		return
	}
}

func (app *Application) handleSnippetCreate(w http.ResponseWriter, r *http.Request) {
	var reqBody models.SnippetBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	_, err = app.snippets.Insert(reqBody.Title, reqBody.Content, reqBody.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	app.sessionManager.Put(r.Context(), "title", reqBody.Title)

	s, err := json.Marshal(reqBody)
	if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Fprint(w, string(s))
}

func (app *Application) handleSnippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	flash := app.sessionManager.Get(r.Context(), "flash")
	title := app.sessionManager.PopString(r.Context(), "title")
	fmt.Println(flash, title)
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	s, err := json.Marshal(snippet)
	if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Fprint(w, string(s))
}

func (app *Application) handleSnippetList(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.GetLatest()

	if err != nil {
		app.serverError(w, err)
		return
	}
	data, _ := json.Marshal(snippets)
	w.Write(data)

}

func (app *Application) handleUserSignup(w http.ResponseWriter, r *http.Request) {
	var reqBody models.UserRequestBody

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if err := app.userManager.Insert(reqBody.Name, reqBody.Email, reqBody.Password); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			errMsg := map[string]string{"error": "Email is already in use"}
			errBytes, err := json.Marshal(errMsg)
			if err != nil {
				app.serverError(w, err)
				return
			}

			fmt.Fprintf(w, "%+v", string(errBytes))

		}
		return
	}
	app.sessionManager.Put(r.Context(), "msg", "your signup was successful")
	fmt.Fprint(w, "success user created")
}

func (app *Application) handleUserSignin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "authenticate user")
}

func (app *Application) handleUserLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
