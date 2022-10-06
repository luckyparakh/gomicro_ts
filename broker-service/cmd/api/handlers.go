package main

import (
	"net/http"
)

func (app *Config) broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit Broker",
	}
	app.writeJson(w, http.StatusOK, payload)
}
