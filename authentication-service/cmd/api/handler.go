package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	var reqPlayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJson(w, r, &reqPlayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	user, err := app.Models.User.GetByEmail(reqPlayload.Email)
	if err != nil {
		app.errorJson(w, errors.New("invalid Credentials"), http.StatusBadRequest)
		return
	}
	valid, err := user.PasswordMatches(reqPlayload.Password)
	if err != nil || !valid {
		app.errorJson(w, errors.New("invalid Credentials"), http.StatusBadRequest)
		return
	}
	palyload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	app.writeJson(w, http.StatusAccepted, palyload)
}
