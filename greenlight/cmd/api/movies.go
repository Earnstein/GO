package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/earnstein/GO/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	if err := app.readJSONResponse(w, r, &reqBody); err != nil {
		app.badRequestErrorHandler(w, r, err)
		return
	}
	data := envelope{"status": "ok", "movie": reqBody}
	app.writeJSONResponse(w, http.StatusCreated, data, nil)
}

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundHandler(w, r)
		return
	}

	Movie := data.Movie{
		ID:        id,
		Title:     "Casablanca",
		Year:      1942,
		Runtime:   102,
		CreatedAt: time.Now(),
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	data := envelope{"movie": Movie}
	err = app.writeJSONResponse(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorHandler(w, r, err)
	}

}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "list movies")
}
