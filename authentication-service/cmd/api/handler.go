package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	var reqPlayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	fmt.Println("Auth service")
	err := app.readJson(w, r, &reqPlayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	fmt.Println("Reading payload")
	user, err := app.Models.User.GetByEmail(reqPlayload.Email)
	if err != nil {
		app.errorJson(w, errors.New("invalid Credentials"), http.StatusBadRequest)
		return
	}
	fmt.Println("Reading from DB")
	valid, err := user.PasswordMatches(reqPlayload.Password)
	if err != nil || !valid {
		app.errorJson(w, errors.New("invalid Credentials"), http.StatusBadRequest)
		return
	}
	fmt.Println("Matching pwd")
	err = app.log("auth", fmt.Sprintf("Logged in user: %s", user.Email))
	fmt.Println(err.Error())
	if err != nil {
		app.errorJson(w, err)
		return
	}
	fmt.Println("logging")
	palyload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	app.writeJson(w, http.StatusAccepted, palyload)
}

func (app *Config) log(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data
	jsonData, _ := json.Marshal(&entry)

	req, err := http.NewRequest("POST", "http://logger-service/log/", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	cl := http.Client{}
	req.Header.Set("Content-Type", "application/json")
	_, err = cl.Do(req)
	if err != nil {
		return err
	}
	return nil
}
