package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/earnstein/GO/greenlight/internal/data"
	"github.com/earnstein/GO/greenlight/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	if err := app.readJSONResponse(w, r, &reqBody); err != nil {
		app.badRequestErrorHandler(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   reqBody.Title,
		Year:    reqBody.Year,
		Runtime: reqBody.Runtime,
		Genres:  reqBody.Genres,
	}
	// Validation
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err := app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorHandler(w, r, err)
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	response := envelope{"status": "ok", "movie": reqBody}
	app.writeJSONResponse(w, http.StatusCreated, response, headers)
}

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundHandler(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundHandler(w, r)
		default:
			app.serverErrorHandler(w, r, err)
		}
		return
	}

	response := envelope{"movie": movie}
	err = app.writeJSONResponse(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorHandler(w, r, err)
	}

}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundHandler(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundHandler(w, r)
		default:
			app.serverErrorHandler(w, r, err)
		}
		return
	}

	var reqBody struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	if err := app.readJSONResponse(w, r, &reqBody); err != nil {
		app.badRequestErrorHandler(w, r, err)
		return
	}

	movie.Title = reqBody.Title
	movie.Year = reqBody.Year
	movie.Runtime = reqBody.Runtime
	movie.Genres = reqBody.Genres

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorHandler(w, r, err)
	}

	response := envelope{"movie": movie}
	err = app.writeJSONResponse(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorHandler(w, r, err)
	}

}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundHandler(w, r)
		return
	}

	if err = app.models.Movies.Delete(id); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundHandler(w, r)
		default:
			app.serverErrorHandler(w, r, err)
		}
	}

	response := envelope{"message": fmt.Sprintf("Movie with id %d is deleted successfully.", id)}
	err = app.writeJSONResponse(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorHandler(w, r, err)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "list movies")
}
