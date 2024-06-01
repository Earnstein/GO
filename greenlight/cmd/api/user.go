package main

import (
	"errors"
	"net/http"

	"github.com/earnstein/GO/greenlight/internal/data"
	"github.com/earnstein/GO/greenlight/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSONResponse(w, r, &reqBody)
	if err != nil {
		app.badRequestErrorHandler(w, r, err)
		return
	}

	newUser := data.NewUser(reqBody.Username, reqBody.Email, false)
	err = newUser.Password.Set(reqBody.Password)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, newUser); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(newUser)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorHandler(w, r, err)
		}
		return
	}
	response := envelope{"user": newUser}
	err = app.writeJSONResponse(w, http.StatusCreated, response, nil)

	if err != nil {
		app.serverErrorHandler(w, r, err)
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

}
