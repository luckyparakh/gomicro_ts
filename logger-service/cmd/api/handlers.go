package main

import (
	"errors"
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) writeLog(w http.ResponseWriter, r *http.Request) {
	var reqPayload JSONPayload
	err := app.readJson(w, r, &reqPayload)
	if err != nil {
		app.errorJson(w, errors.New("not able to read log payload"))
		return
	}
	event := data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}
	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}
	app.writeJson(w, http.StatusAccepted, resp)
}
