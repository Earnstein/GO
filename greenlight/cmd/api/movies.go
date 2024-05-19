package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/earnstein/GO/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		http.NotFound(w, r)
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
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "list movies")
}
