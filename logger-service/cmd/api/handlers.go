package main

import (
	"log"
	"net/http"
	"time"

	"github.com/puneet222/go-micro/logger-service/data"
)

type RequestPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) logEntry(w http.ResponseWriter, r *http.Request) {
	var rp RequestPayload
	err := app.readJSON(w, r, &rp)
	if err != nil {
		log.Println("unable to read json from request", err)
	}

	// log event
	event := data.LogEntry{
		Name:      rp.Name,
		Data:      rp.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
		Data:    event,
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
