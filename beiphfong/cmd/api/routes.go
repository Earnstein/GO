package main

import (
	"expvar"
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
	router.Handler(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", (app.listMoviesHandler)))
	router.Handler(http.MethodPost, "/v1/movie", app.requirePermission("movies:write", (app.createMovieHandler)))
	router.Handler(http.MethodGet, "/v1/movie/:id", app.requirePermission("movies:read", (app.getMovieHandler)))
	router.Handler(http.MethodPut, "/v1/movie/:id", app.requirePermission("movies:write", (app.updateMovieHandler)))
	router.Handler(http.MethodPatch, "/v1/movie/:id", app.requirePermission("movies:write", (app.patchUpdateMovieHandler)))
	router.Handler(http.MethodDelete, "/v1/movie/:id", app.requirePermission("movies:write", (app.deleteMovieHandler)))

	// METRICS ROUTES
	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	// GENERAL MIDDLEWARE
	middlewareChain := alice.New(app.recoverPanicMiddleware, app.enableCorsMiddleware, app.rateLimitMiddleware, app.requestInfoMiddleware, app.authenticate)
	return middlewareChain.Then(router)
}
