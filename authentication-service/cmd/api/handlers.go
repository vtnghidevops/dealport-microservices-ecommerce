package main

import (
	// "encoding/json"
	"fmt"
	"errors"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload) 
	if ( err != nil ) {
		app.errorJSON(w,err,http.StatusBadRequest)
		return
	}

	// validate the user against the db
	user, err := app.Models.User.GetByEmail(requestPayload.Email) 
	if ( err != nil ) {
		app.errorJSON(w,errors.New("invalid credentials"),http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if ( err != nil && !valid ){
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse {
		Error: false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data: user,
	}

	app.writeJson(w, http.StatusAccepted, payload)

}