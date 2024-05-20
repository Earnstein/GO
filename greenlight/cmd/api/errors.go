package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app *application) errorHandler(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	response := envelope{"error": message}

	err := app.writeJSONResponse(w, status, response, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorHandler(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorHandler(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorHandler(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	app.errorHandler(w, r, http.StatusBadRequest, err.Error())
}