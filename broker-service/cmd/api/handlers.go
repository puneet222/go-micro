package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type requestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "OK",
	}

	out, _ := json.MarshalIndent(payload, "", " ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}

func (app *Config) handleSubmission(w http.ResponseWriter, r *http.Request) {
	var rp requestPayload
	err := app.readJSON(w, r, &rp)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
	}

	switch rp.Action {
	case "auth":
		app.authenticate(w, rp.Auth)
	default:
		app.errorJSON(w, fmt.Errorf("invalid action %s", rp.Action), http.StatusBadRequest)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json and send it to auth microservice
	jsonPayload, err := json.Marshal(a)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// call the microservice
	request, err := http.NewRequest(http.MethodPost, "http://authentication-service/authenticate", bytes.NewBuffer(jsonPayload))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("err calling auth service"))
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	log.Println("error while decoding response body", err)
	if err != nil {
		app.errorJSON(w, err, http.StatusUnauthorized)
	}

	log.Println("broker service", "json from service", jsonFromService)

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	app.writeJSON(w, http.StatusOK, jsonFromService)
}
