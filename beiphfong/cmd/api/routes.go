package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// NOT FOUND AND METHOD NOT ALLOWED ROUTES
	router.NotFound = http.HandlerFunc(app.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedHandler)

	// SYSTEM CHECK HANDLER
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// USER ROUTES
	router.HandlerFunc(http.MethodPost, "/v1/user", app.createUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/user/:id", app.getUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/user/:id", app.updateUserHandler)

	// MOVIES ROUTES
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movie", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getMovieHandler)
	router.HandlerFunc(http.MethodPut, "/v1/movie/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movie/:id", app.patchUpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movie/:id", app.deleteMovieHandler)

	return app.rateLimit(app.recoverPanic(router))
}