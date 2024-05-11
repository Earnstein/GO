package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

type ErrorMessge map[string]string

func NewJsonResponse(app *Application, w http.ResponseWriter, statusCode int) []byte {
	newJsonmessage := ErrorMessge{"error": http.StatusText(statusCode)}
	errMsg, err := json.Marshal(newJsonmessage)
	if err != nil {
		app.serverError(w, err)
		return nil
	}
	return errMsg
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	errMsg := NewJsonResponse(app, w, http.StatusInternalServerError)
	http.Error(w, string(errMsg), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, statusCode int) {
	errMsg := NewJsonResponse(app, w, statusCode)
	http.Error(w, string(errMsg), statusCode)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
