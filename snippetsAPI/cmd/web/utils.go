package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) serverError(w http.ResponseWriter,  err error){
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}


func (app *Application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *Application) NotFound(w http.ResponseWriter){
	app.clientError(w, http.StatusNotFound)
}


