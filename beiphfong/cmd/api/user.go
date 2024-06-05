package main

import (
	"errors"
	"net/http"
	"time"

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

	_, err = app.models.Users.GetByEmail(newUser.Email)
	if nil == err {
		v.AddError("email", "a user with this email address already exists")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(newUser)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}

	// generate a token
	token, err := app.models.Tokens.New(newUser.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          newUser.ID,
		}
		err = app.mailer.Send(newUser.Email, "user_email.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
			return
		}
	})

	response := envelope{"user": newUser}
	err = app.writeJSONResponse(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}
}

func (app *application) activeUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSONResponse(w, r, &input)
	if err != nil {
		app.badRequestErrorHandler(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorHandler(w, r, err)
		}
		return
	}

	user.Activated = true
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorHandler(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}

	response := envelope{"message": "user successfully activated", "user": user}
	err = app.writeJSONResponse(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}

}
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.GetUsers()
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}

	err = app.writeJSONResponse(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}
}
