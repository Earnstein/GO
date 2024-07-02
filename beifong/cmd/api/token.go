package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/earnstein/GO/greenlight/internal/data"
	"github.com/earnstein/GO/greenlight/internal/validator"
)

func (app *application) signInHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSONResponse(w, r, &input)
	if err != nil {
		app.badRequestErrorHandler(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidateLoginPassword(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorHandler(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 4*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}

	err = app.writeJSONResponse(w, http.StatusCreated, envelope{"authentication": token}, nil)
	if err != nil {
		app.serverErrorHandler(w, r, err)
	}
}
