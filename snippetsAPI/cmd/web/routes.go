package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	sessionMiddleware := alice.New(app.sessionManager.LoadAndSave)

	// snippets routes
	router.Handler(http.MethodGet, "/", sessionMiddleware.ThenFunc(app.homeHandler))
	router.Handler(http.MethodPost, "/snippet/create", sessionMiddleware.ThenFunc(app.handleSnippetCreate))
	router.Handler(http.MethodGet, "/snippet/view/:id", sessionMiddleware.ThenFunc(app.handleSnippetView))
	router.Handler(http.MethodGet, "/snippet/latest", sessionMiddleware.ThenFunc(app.handleSnippetList))

	// user routes

	router.Handler(http.MethodPost, "/user/sign-up", sessionMiddleware.ThenFunc(app.handleUserSignup))
	router.Handler(http.MethodPost, "/user/sign-in", sessionMiddleware.ThenFunc(app.handleUserSignin))
	router.Handler(http.MethodGet, "/user/logout", sessionMiddleware.ThenFunc(app.handleUserLogout))

	middlewareChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return middlewareChain.Then(router)
}
