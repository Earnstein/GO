package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// NOT FOUND AND METHOD NOT ALLOWED ROUTES
	router.NotFound = http.HandlerFunc(app.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedHandler)

	// SYSTEM CHECK HANDLER
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// AUTHENTICATION ROUTES
	router.HandlerFunc(http.MethodPost, "/v1/user", app.createUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/token/authentication", app.signInHandler)

	// USER ROUTES
	router.HandlerFunc(http.MethodGet, "/v1/user/:id", app.getUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/user/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/user/activated", app.activeUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users", app.getAllUsers)

	// MOVIES ROUTES
	protectedMiddleWare := alice.New(app.requireActivateUser)
	router.Handler(http.MethodGet, "/v1/movies", protectedMiddleWare.ThenFunc(app.listMoviesHandler))
	router.Handler(http.MethodPost, "/v1/movie", protectedMiddleWare.ThenFunc(app.createMovieHandler))
	router.Handler(http.MethodGet, "/v1/movie/:id", protectedMiddleWare.ThenFunc(app.getMovieHandler))
	router.Handler(http.MethodPut, "/v1/movie/:id", protectedMiddleWare.ThenFunc(app.updateMovieHandler))
	router.Handler(http.MethodPatch, "/v1/movie/:id", protectedMiddleWare.ThenFunc(app.patchUpdateMovieHandler))
	router.Handler(http.MethodDelete, "/v1/movie/:id", protectedMiddleWare.ThenFunc(app.deleteMovieHandler))

	// GENERAL MIDDLEWARE
	middlewareChain := alice.New(app.recoverPanicMiddleware, app.rateLimitMiddleware, app.requestInfoMiddleware, app.authenticate)
	return middlewareChain.Then(router)
}
