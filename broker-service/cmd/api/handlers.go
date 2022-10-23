package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit Broker",
	}
	app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) handleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPlayload RequestPayload
	err := app.readJson(w, r, &reqPlayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	switch reqPlayload.Action {
	case "auth":
		app.authenticate(w, reqPlayload.Auth)
	case "log":
		app.log(w, reqPlayload.Log)
	default:
		app.errorJson(w, errors.New("unknown Action"))
	}
}
func (app *Config) log(w http.ResponseWriter, l LogPayload) {
	jsonData, err := json.Marshal(l)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	req, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling log service"))
		return
	}
	var jsonFromService jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJson(w, err)
		return
	}
	var response jsonResponse
	response.Error = jsonFromService.Error
	response.Message = jsonFromService.Message

	app.writeJson(w, http.StatusAccepted, response)
}
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.errorJson(w, err)
		return
	}
	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentails"))
		return
	} else if resp.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling auth service"))
		return
	}
	var jsonFromService jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJson(w, err)
		return
	}
	var response jsonResponse
	response.Error = false
	response.Message = "Authenticated"
	response.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, response)
}
